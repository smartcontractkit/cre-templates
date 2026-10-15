// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ReceiverTemplate} from "./ReceiverTemplate.sol";

/**
 * @title Storage
 * @notice Minimal consumer that stores a uint256 delivered by a CRE workflow.
 *
 *         The workflow ABI-encodes a single uint256 as the report body; the
 *         KeystoneForwarder delivers the DON-signed report via onReport, and
 *         ReceiverTemplate validates the forwarder (and optionally the expected
 *         author / workflow name / workflow ID) before _processReport runs.
 */
contract Storage is ReceiverTemplate {
    uint256 public value;

    event ValueUpdated(uint256 value);

    constructor(address forwarder) ReceiverTemplate(forwarder) {}

    /**
     * @notice Processes the CRE-signed report.
     * @dev    report = abi.encode(uint256 value)
     */
    function _processReport(bytes calldata report) internal override {
        uint256 newValue = abi.decode(report, (uint256));
        value = newValue;
        emit ValueUpdated(newValue);
    }
}
