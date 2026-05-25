// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./ReceiverTemplate.sol";

/**
 * @title AutomationReceiver
 * @notice Generic bridge to execute Automation-style upkeeps via CRE
 */
contract AutomationReceiver is ReceiverTemplate {
    event CallExecuted(address indexed target, bool success, bytes returnData);
    event CallFailed(address indexed target, bytes reason);

    error InvalidTargetAddress();
    error CallExecutionFailed(address target, bytes reason);

    constructor(address _forwarder) ReceiverTemplate(_forwarder) {}

    /**
     * @notice Decodes and executes the call on the target contract
     * @param report ABI-encoded (address target, bytes data)
     */
    function _processReport(bytes calldata report) internal override {
        (address target, bytes memory data) = abi.decode(report, (address, bytes));

        if (target == address(0)) {
            revert InvalidTargetAddress();
        }

        (bool success, bytes memory returnData) = target.call(data);

        if (success) {
            emit CallExecuted(target, success, returnData);
        } else {
            emit CallFailed(target, returnData);
            revert CallExecutionFailed(target, returnData);
        }
    }
}
