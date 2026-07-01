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
    // Sum of non-target gas costs inside _processReport (EIP-2929 cold-slot / cold-address):
    //   SLOAD s_consumerGasLimit (cold slot)          2,100
    //   Pre-call ops (GAS, ADD, LT, JUMPI, stack)        50
    //   CALL opcode dispatch to target (cold addr)    2,600
    //   Post-call (success flag, JUMPI, LOG2 min)     1,200
    //   Misc stack / memory                              50
    //                                        Total:   6,000
    uint256 private constant GAS_OVERHEAD = 6_000;

    /// @notice Closed-by-default allowlist of callable (target, selector) pairs.
    mapping(address target => mapping(bytes4 selector => bool allowed)) private s_callAllowed;

    /// @notice Per-(target, selector) minimum gas the consumer needs to execute. 0 = no limit.
    mapping(address target => mapping(bytes4 selector => uint256 gasLimit)) private s_consumerGasLimit;

    /// @notice Emitted when a target call succeeds.
    event CallExecuted(address indexed target, bytes4 indexed selector, bytes returnData);
    /// @notice Emitted when an allowed target call reverts. The report is still consumed.
    event CallFailed(address indexed target, bytes4 indexed selector, bytes reason);
    /// @notice Emitted when the owner updates the outbound allowlist.
    event CallAllowedSet(address indexed target, bytes4 indexed selector, bool allowed);
    /// @notice Emitted when the owner updates the consumer gas limit for a (target, selector) pair.
    event ConsumerGasLimitSet(address indexed target, bytes4 indexed selector, uint256 previousLimit, uint256 newLimit);

    /// @notice Thrown when the decoded target is the zero address.
    error InvalidTargetAddress();
    /// @notice Thrown when the target address has no deployed code (EOA, mistyped address, or never-deployed contract).
    error TargetHasNoCode(address target);
    /// @notice Thrown when the report carries fewer than 4 bytes of calldata (no selector).
    error MissingSelector();
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
    /// @dev When non-zero, `_processReport` will revert with `InsufficientGas` before calling the
    ///      target if available gas is below `gasLimit + GAS_OVERHEAD`. This causes the CRE
    ///      Forwarder to record the transmission as failed (retryable) rather than permanently
    ///      consuming the report. Set this to the `performGasLimit` tuned in Automation for the
    ///      specific function being migrated. Each (target, selector) pair has its own limit.
    ///      Zero (the default) disables the guard for that pair and preserves fire-and-forget.
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

    /// @notice Decodes and executes the call on the target contract.
    /// @param report ABI-encoded (address target, bytes data), where `data` is a full
    ///        function call (4-byte selector followed by its arguments).
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
    ///      Authorization failures (zero target, missing selector, not-allowlisted) revert
    ///      loudly — they indicate misconfiguration or a malformed report. Execution failures
    ///      (an allowed call that reverts) are swallowed: `CallFailed` is emitted and the
    ///      report is consumed, matching Chainlink Automation's fire-and-forget semantics
    ///      where the next trigger re-evaluates eligibility.
    ///      Gas-guard: when `s_consumerGasLimit[target][selector]` is non-zero, the function
    ///      reverts with `InsufficientGas` before the target call if available gas is below
    ///      `gasLimit + GAS_OVERHEAD`. This ensures a low-gas delivery is recorded as failed
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

        (address target, bytes memory data) = abi.decode(report, (address, bytes));

        if (target == address(0)) {
            revert InvalidTargetAddress();
        }
        if (data.length < 4) {
            revert MissingSelector();
        }

        // Read the leading 4-byte selector from the in-memory `data`. Safe: length >= 4 is
        // checked above, and a 32-byte mload assigned to bytes4 keeps the high (first) 4 bytes.
        bytes4 selector;
        assembly {
            selector := mload(add(data, 0x20))
        }
        if (!s_callAllowed[target][selector]) {
            revert CallNotAllowed(target, selector);
        }

        uint256 consumerGasLimit = s_consumerGasLimit[target][selector];
        bool success;
        bytes memory returnData;
        if (consumerGasLimit > 0) {
            uint256 required = consumerGasLimit + GAS_OVERHEAD;
            if (gasleft() < required) {
                revert InsufficientGas(gasleft(), required);
            }
            (success, returnData) = target.call{gas: consumerGasLimit}(data);
        } else {
            (success, returnData) = target.call(data);
        }

        if (success) {
            emit CallExecuted(target, selector, returnData);
        } else {
            emit CallFailed(target, selector, returnData);
        }
    }
}
