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
 *            and that both workflowId and workflowOwner are configured before any report is
 *            accepted. Either field alone does not fully bind the receiver to a single workflow.
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
    /// @notice Closed-by-default allowlist of callable (target, selector) pairs.
    mapping(address target => mapping(bytes4 selector => bool allowed)) private s_callAllowed;

    /// @notice Emitted when a target call succeeds.
    event CallExecuted(address indexed target, bytes4 indexed selector, bytes returnData);
    /// @notice Emitted when an allowed target call reverts. The report is still consumed.
    event CallFailed(address indexed target, bytes4 indexed selector, bytes reason);
    /// @notice Emitted when the owner updates the outbound allowlist.
    event CallAllowedSet(address indexed target, bytes4 indexed selector, bool allowed);

    /// @notice Thrown when the decoded target is the zero address.
    error InvalidTargetAddress();
    /// @notice Thrown when the target address has no deployed code (EOA, mistyped address, or never-deployed contract).
    error TargetHasNoCode(address target);
    /// @notice Thrown when the report carries fewer than 4 bytes of calldata (no selector).
    error MissingSelector();
    /// @notice Thrown when (target, selector) is not on the outbound allowlist.
    error CallNotAllowed(address target, bytes4 selector);
    /// @notice Thrown when onReport is called but neither workflowId nor workflowOwner has been
    ///         configured. At least one must be set so the receiver is bound to a specific
    ///         workflow and cannot be triggered by an arbitrary DON-signed report.
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

    /// @notice Decodes and executes the call on the target contract.
    /// @param report ABI-encoded (address target, bytes data), where `data` is a full
    ///        function call (4-byte selector followed by its arguments).
    /// @dev Two pre-conditions are enforced before any decoding:
    ///      1. The forwarder address must not be zero. `ReceiverTemplate.setForwarderAddress`
    ///         does not block address(0), so this guard closes that gap: if the owner ever
    ///         sets the forwarder to zero (disabling the caller check in onReport), every
    ///         subsequent report delivery is rejected here instead.
    ///      2. Both workflow identity fields (workflowId and workflowOwner) must be
    ///         configured. Either field alone leaves an attack surface: workflowId alone does
    ///         not bind the report to a specific owner, and workflowOwner alone does not bind
    ///         it to a specific workflow instance. Requiring both closes the cross-receiver
    ///         replay vector described in audit finding M-02.
    ///      Authorization failures (zero target, missing selector, not-allowlisted) revert
    ///      loudly — they indicate misconfiguration or a malformed report. Execution failures
    ///      (an allowed call that reverts) are swallowed: `CallFailed` is emitted and the
    ///      report is consumed, matching Chainlink Automation's fire-and-forget semantics
    ///      where the next trigger re-evaluates eligibility.
    function _processReport(bytes calldata report) internal override {
        if (this.getForwarderAddress() == address(0)) {
            revert InvalidForwarderAddress();
        }
        if (this.getExpectedWorkflowId() == bytes32(0) || this.getExpectedAuthor() == address(0)) {
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

        (bool success, bytes memory returnData) = target.call(data);

        if (success) {
            emit CallExecuted(target, selector, returnData);
        } else {
            emit CallFailed(target, selector, returnData);
        }
    }
}
