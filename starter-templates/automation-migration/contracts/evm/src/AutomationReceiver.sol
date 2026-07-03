// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import "./ReceiverTemplate.sol";

/**
 * @title AutomationReceiver
 * @notice Generic bridge that executes Automation-style upkeeps delivered by a CRE workflow.
 *
 * @dev Two independent authorization layers protect this contract:
 *
 *      1. INBOUND — answers "who may deliver a report?":
 *         a) {ReceiverTemplate} enforces the CRE Forwarder address check and optional
 *            workflowId / workflowName / workflowOwner identity checks.
 *         b) {_processReport} additionally requires that the forwarder is non-zero (closing
 *            the gap left by `ReceiverTemplate.setForwarderAddress` which permits address(0))
 *            and that at least one complete workflow identity option is configured before any
 *            report is accepted: either (1) workflowId is set, or (2) both workflowOwner and
 *            workflowName are set. Neither option alone is sufficient for option 2.
 *
 *      2. OUTBOUND (this contract) — answers "what may a report make this contract do?":
 *         a closed-by-default allowlist of (target, function-selector) pairs. The inbound checks
 *         do NOT constrain the `(target, data)` a report carries, so without this layer any
 *         authorized report could call any contract/function this receiver is trusted by. The
 *         owner must explicitly allow each (target, selector) before it can be executed.
 *
 *      Migration rule of thumb: inbound authorizes the workflow; outbound authorizes the action.
 */
contract AutomationReceiver is ReceiverTemplate {
    // Non-target gas costs reserved by the guard (gasleft() < required).
    // The SLOAD is the last cost incurred BEFORE the check; everything else is AFTER.
    //
    // [Pre-check — already spent at the guard point:]
    //   SLOAD s_consumerGasLimit (cold mapping)         2,100
    //
    // [Post-check — must complete with remaining gas:]
    //   Pre-call ops (GAS, ADD, LT, JUMPI, stack)          50
    //   CALL opcode dispatch to target (cold addr)      2,600
    //     (EIP-2929 replaces the pre-Berlin 700 base; 2,600 is the full cold-access cost)
    //   Post-call (success flag, JUMPI, LOG3 min)       2,200
    //     ├─ LOG3 base + 3 topics (375 + 3×375)         1,500
    //     ├─ LOG3 data  (64 B ABI-encoded empty bytes)    512
    //     └─ misc (returnData mem, JUMPI, stack)          188
    //   Misc stack / memory                                50
    //                                          Total:   7,000
    //
    // EIP-150 (63/64 rule): a CALL can forward at most 63/64 of available gas.
    // A fixed GAS_OVERHEAD buffer is only sufficient when consumerGasLimit ≤ 63 × GAS_OVERHEAD
    // (~441,000). Above that threshold the 63/64 cap would deliver less than consumerGasLimit
    // to the target. _processReport therefore adds consumerGasLimit / 63 to `required`
    // dynamically, ensuring the available gas at the CALL satisfies:
    //   63/64 × available  ≥  consumerGasLimit
    //
    // Pre-guard overhead (excluded from GAS_OVERHEAD, paid before the check point):
    // Five cold SLOADs in ReceiverTemplate.onReport (s_forwarderAddress, s_expectedWorkflowId ×2,
    // s_expectedAuthor, s_expectedWorkflowName: 5 × 2,100 = 10,500 gas), two STATICCALL frames to
    // this in _processReport (~1,000 gas), abi.decode of the report (~300 gas), and the cold SLOAD
    // of s_callAllowed (~2,100 gas) add up to ~14,000 gas on top of consumerGasLimit + GAS_OVERHEAD.
    // When the block-number monotonicity check is enabled for a (target, selector) pair, two extra
    // cold SLOADs (s_blockNumberCheckEnabled, s_lastReportBlock) plus one SSTORE add a further
    // ~10,000 gas before the guard. Callers must budget for this in writeGasLimit (see README Gas
    // Limit section).
    uint256 private constant GAS_OVERHEAD = 7_000;

    /// @notice Closed-by-default allowlist of callable (target, selector) pairs.
    mapping(address target => mapping(bytes4 selector => bool allowed)) private s_callAllowed;

    /// @notice Per-(target, selector) minimum gas the consumer needs to execute. 0 = no limit.
    mapping(address target => mapping(bytes4 selector => uint256 gasLimit)) private s_consumerGasLimit;

    /// @notice Opt-in per-(target, selector) block-number monotonicity switch. Closed by default.
    mapping(address target => mapping(bytes4 selector => bool enabled)) private s_blockNumberCheckEnabled;

    /// @notice Highest report block number accepted so far for a (target, selector) pair. Doubles as
    ///         the configurable initial floor set by {setBlockNumberCheck}.
    mapping(address target => mapping(bytes4 selector => uint256 blockNumber)) private s_lastReportBlock;

    /// @notice Emitted when a target call succeeds.
    event CallExecuted(address indexed target, bytes4 indexed selector, bytes returnData);
    /// @notice Emitted when an allowed target call reverts. The report is still consumed.
    event CallFailed(address indexed target, bytes4 indexed selector, bytes reason);
    /// @notice Emitted when the owner updates the outbound allowlist.
    event CallAllowedSet(address indexed target, bytes4 indexed selector, bool allowed);
    /// @notice Emitted when the owner updates the consumer gas limit for a (target, selector) pair.
    event ConsumerGasLimitSet(address indexed target, bytes4 indexed selector, uint256 previousLimit, uint256 newLimit);
    /// @notice Emitted when the owner enables/disables the block-number monotonicity check for a pair.
    ///         When enabled, `initialBlockNumber` is the floor the next report must meet or exceed.
    event BlockNumberCheckSet(address indexed target, bytes4 indexed selector, bool enabled, uint256 initialBlockNumber);
    /// @notice Emitted when a report is skipped because its block number is older than the last
    ///         accepted one for the pair. The report is consumed (not reverted), mirroring
    ///         Chainlink Automation's fire-and-forget semantics.
    event StaleReportSkipped(address indexed target, bytes4 indexed selector, uint256 reportBlockNumber, uint256 lastReportBlock);

    /// @notice Thrown when the decoded target is the zero address.
    error InvalidTargetAddress();
    /// @notice Thrown when the target address has no deployed code (EOA, mistyped address, or never-deployed contract).
    error TargetHasNoCode(address target);
    /// @notice Thrown when (target, selector) is not on the outbound allowlist.
    error CallNotAllowed(address target, bytes4 selector);
    /// @notice Thrown when there is not enough gas to safely forward consumerGasLimit to the target.
    ///         Causes the forwarder to record the transmission as failed so it can be retried.
    error InsufficientGas(uint256 available, uint256 required);
    /// @notice Thrown when onReport is called without a complete workflow identity configuration.
    ///         The receiver requires exactly one of the two valid options to be satisfied:
    ///         (1) workflowId is set, or (2) both workflowOwner and workflowName are set.
    ///         Without at least one complete option the receiver cannot be bound to a specific
    ///         workflow and would accept reports from any DON-signed payload.
    error WorkflowIdentityNotConfigured();

    constructor(address _forwarder) ReceiverTemplate(_forwarder) {}

    /// @notice Allow or disallow the receiver to call `selector` on `target`.
    /// @dev Closed by default. Register every (target, selector) the migrated upkeep needs,
    ///      e.g. `performUpkeep(bytes)` for custom-logic/log upkeeps, or your specific
    ///      time-based function. Owner-only.
    ///      Validates that `target` has deployed code at the time of registration; passing an EOA,
    ///      a mistyped address, or a never-deployed address reverts with `TargetHasNoCode`.
    /// @param target The contract the receiver is permitted to call.
    /// @param selector The 4-byte function selector permitted on `target`.
    /// @param allowed True to permit, false to revoke.
    function setCallAllowed(address target, bytes4 selector, bool allowed) external onlyOwner {
        if (target == address(0)) {
            revert InvalidTargetAddress();
        }
        if (target.code.length == 0) {
            revert TargetHasNoCode(target);
        }
        s_callAllowed[target][selector] = allowed;
        emit CallAllowedSet(target, selector, allowed);
    }

    /// @notice Returns whether the receiver may call `selector` on `target`.
    function isCallAllowed(address target, bytes4 selector) external view returns (bool) {
        return s_callAllowed[target][selector];
    }

    /// @notice Set the minimum gas required to execute `selector` on `target`.
    /// @dev When non-zero, `_processReport` will revert with `InsufficientGas` before calling
    ///      the target if available gas is below `gasLimit + gasLimit / 63 + GAS_OVERHEAD`.
    ///      The `gasLimit / 63` term compensates for the EIP-150 (63/64) rule: a CALL can
    ///      forward at most 63/64 of the available gas, so without this buffer a high gas limit
    ///      would cause the target to receive less than `gasLimit`. This causes the CRE
    ///      Forwarder to record the transmission as failed (retryable) rather than permanently
    ///      consuming the report. Set this to the `performGasLimit` tuned in Automation for the
    ///      specific function being migrated. Each (target, selector) pair has its own limit.
    ///      Zero (the default) disables the guard for that pair and preserves fire-and-forget.
    ///      Note: the on-chain formula only covers costs from the guard check onward. The
    ///      workflow's writeGasLimit must also budget for ~14,000 gas of pre-guard overhead
    ///      (five cold SLOADs in ReceiverTemplate, two STATICCALL frames, abi.decode, and
    ///      the cold s_callAllowed SLOAD) on top of this limit.
    /// @param target  The contract the limit applies to. Must not be the zero address.
    /// @param selector The 4-byte function selector the limit applies to.
    /// @param gasLimit Minimum gas required by the consumer. 0 = no guard.
    function setConsumerGasLimit(address target, bytes4 selector, uint256 gasLimit) external onlyOwner {
        if (target == address(0)) revert InvalidTargetAddress();
        uint256 previous = s_consumerGasLimit[target][selector];
        s_consumerGasLimit[target][selector] = gasLimit;
        emit ConsumerGasLimitSet(target, selector, previous, gasLimit);
    }

    /// @notice Returns the configured consumer gas limit for a (target, selector) pair (0 = no guard).
    function getConsumerGasLimit(address target, bytes4 selector) external view returns (uint256) {
        return s_consumerGasLimit[target][selector];
    }

    /// @notice Enable or disable the block-number monotonicity (staleness) check for a (target, selector).
    /// @dev Closed by default. When enabled, `_processReport` rejects any report whose encoded block
    ///      number is strictly less than the last accepted one for the pair (CLA-parity: a report at
    ///      the same block is still accepted). Stale reports are consumed, not reverted, so the
    ///      forwarder does not pointlessly retry an always-stale delivery.
    ///
    ///      Under Chainlink Automation the registry rejected any report whose trigger block preceded
    ///      the last performed block. The CRE delivery path does not preserve this ordering, so this
    ///      opt-in check restores it for conditional-trigger / log upkeeps that relied on it. The
    ///      workflow always encodes a block number in the report payload, so this can be turned on
    ///      later without any workflow change.
    ///
    ///      When enabling, the initial floor is configurable:
    ///        - `initialBlockNumber == 0`: snapshot the current `block.number` as the floor. Note this
    ///          does NOT reject reports minted between deployment/config and the moment this check is
    ///          enabled if they carry a block number above the snapshot; prefer an explicit floor when
    ///          that matters.
    ///        - `initialBlockNumber != 0`: use the provided value as the floor.
    ///      The very first accepted report must have a block number >= this floor.
    ///      Each (target, selector) pair is configured independently. Owner-only.
    /// @param target The contract the check applies to. Must not be the zero address.
    /// @param selector The 4-byte function selector the check applies to.
    /// @param enabled True to enable the check, false to disable it (and clear the stored floor).
    /// @param initialBlockNumber Floor for the first report when enabling; 0 snapshots `block.number`.
    function setBlockNumberCheck(address target, bytes4 selector, bool enabled, uint256 initialBlockNumber)
        external
        onlyOwner
    {
        if (target == address(0)) revert InvalidTargetAddress();
        s_blockNumberCheckEnabled[target][selector] = enabled;
        uint256 floor = enabled ? (initialBlockNumber == 0 ? block.number : initialBlockNumber) : 0;
        s_lastReportBlock[target][selector] = floor;
        emit BlockNumberCheckSet(target, selector, enabled, floor);
    }

    /// @notice Returns the block-number monotonicity configuration for a (target, selector) pair.
    /// @return enabled Whether the check is active.
    /// @return lastReportBlock The highest report block number accepted so far (or the floor set when
    ///         the check was enabled, if no report has been accepted since).
    function getBlockNumberCheck(address target, bytes4 selector)
        external
        view
        returns (bool enabled, uint256 lastReportBlock)
    {
        return (s_blockNumberCheckEnabled[target][selector], s_lastReportBlock[target][selector]);
    }

    /// @notice Decodes and executes the call on the target contract.
    /// @param report ABI-encoded (address target, bytes4 selector, uint256 blockNumber, bytes data),
    ///        where `selector` is the 4-byte function selector, `blockNumber` is the block the report
    ///        was produced for (used only by the opt-in monotonicity check), and `data` is the
    ///        ABI-encoded function arguments (without the selector). The call executed on `target` is
    ///        `bytes.concat(selector, data)`.
    /// @dev Two pre-conditions are enforced before any decoding:
    ///      1. The forwarder address must not be zero. `ReceiverTemplate.setForwarderAddress`
    ///         does not block address(0), so this guard closes that gap: if the owner ever
    ///         sets the forwarder to zero (disabling the caller check in onReport), every
    ///         subsequent report delivery is rejected here instead.
    ///      2. A complete workflow identity option must be configured. Two options are accepted:
    ///         (a) workflowId is set — binds the receiver to one specific workflow instance; or
    ///         (b) both workflowOwner and workflowName are set — binds the receiver to a named
    ///         workflow from a specific owner. Either piece of option (b) alone is insufficient:
    ///         owner alone allows any workflow from that owner; name alone is globally ambiguous.
    ///         Requiring a complete option closes the cross-receiver replay vector from audit M-02.
    ///      Authorization failures (zero target, not-allowlisted) revert loudly — they indicate
    ///      misconfiguration or a malformed report. Execution failures (an allowed call that
    ///      reverts) are swallowed: `CallFailed` is emitted and the report is consumed, matching
    ///      Chainlink Automation's fire-and-forget semantics where the next trigger re-evaluates
    ///      eligibility.
    ///      Staleness: when the block-number monotonicity check is enabled for the pair (see
    ///      {setBlockNumberCheck}), a report whose block number is older than the last accepted one
    ///      is skipped — `StaleReportSkipped` is emitted and the report is consumed without executing
    ///      the target call or reverting.
    ///      Gas-guard: when `s_consumerGasLimit[target][selector]` is non-zero, the function
    ///      reverts with `InsufficientGas` before the target call if available gas is below
    ///      `gasLimit + gasLimit / 63 + GAS_OVERHEAD`. The `gasLimit / 63` term accounts for
    ///      the EIP-150 (63/64) rule: a CALL forwards at most 63/64 of available gas, so
    ///      without this buffer a high gas limit (above ~441,000) would cause the target to
    ///      receive less than configured. This ensures a low-gas delivery is recorded as failed
    ///      by the forwarder and can be retried, preventing griefing attacks. Each
    ///      (target, selector) pair has its own configurable limit.
    function _processReport(bytes calldata report) internal override {
        if (this.getForwarderAddress() == address(0)) {
            revert InvalidForwarderAddress();
        }
        if (this.getExpectedWorkflowId() == bytes32(0) &&
            (this.getExpectedAuthor() == address(0) || this.getExpectedWorkflowName() == bytes10(0))) {
            revert WorkflowIdentityNotConfigured();
        }

        (address target, bytes4 selector, uint256 reportBlockNumber, bytes memory data) =
            abi.decode(report, (address, bytes4, uint256, bytes));

        if (target == address(0)) {
            revert InvalidTargetAddress();
        }
        if (!s_callAllowed[target][selector]) {
            revert CallNotAllowed(target, selector);
        }

        // Opt-in staleness protection: reject reports older than the last accepted one for this pair.
        // CLA-parity comparison (>=): a report at the same block is still accepted. The stored block
        // is updated before the gas guard, but an InsufficientGas revert rolls it back so a retry at
        // the same block stays valid.
        if (s_blockNumberCheckEnabled[target][selector]) {
            uint256 lastReportBlock = s_lastReportBlock[target][selector];
            if (reportBlockNumber < lastReportBlock) {
                emit StaleReportSkipped(target, selector, reportBlockNumber, lastReportBlock);
                return;
            }
            s_lastReportBlock[target][selector] = reportBlockNumber;
        }

        bytes memory callData = bytes.concat(selector, data);

        uint256 consumerGasLimit = s_consumerGasLimit[target][selector];
        bool success;
        bytes memory returnData;
        if (consumerGasLimit > 0) {
            // consumerGasLimit / 63 compensates for EIP-150: a CALL forwards at most
            // 63/64 of available gas. Without this term, limits above ~441,000
            // (63 × GAS_OVERHEAD) would cause the target to receive less than requested.
            uint256 required = consumerGasLimit + consumerGasLimit / 63 + GAS_OVERHEAD;
            if (gasleft() < required) {
                revert InsufficientGas(gasleft(), required);
            }
            (success, returnData) = target.call{gas: consumerGasLimit}(callData);
        } else {
            (success, returnData) = target.call(callData);
        }

        if (success) {
            emit CallExecuted(target, selector, returnData);
        } else {
            emit CallFailed(target, selector, returnData);
        }
    }
}
