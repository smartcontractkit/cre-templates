// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ReceiverTemplate} from "./ReceiverTemplate.sol";

/**
 * @title ReactorConsumer
 * @notice Receives CRE reports and takes action (e.g., flag an address).
 *         A monitored contract emits LargeTransfer events; CRE workflow
 *         checks off-chain data and writes back a flagging action.
 *
 *         Customize _processReport() with your own reaction logic.
 */
contract ReactorConsumer is ReceiverTemplate {
    mapping(address => bool) public flaggedAddresses;
    uint256 public totalFlags;

    event AddressFlagged(address indexed target, string reason, uint256 timestamp);
    event ActionTaken(uint8 indexed actionType, address indexed target, uint256 timestamp);

    constructor(address forwarder) ReceiverTemplate(forwarder) {}

    /// @notice Called by CRE Forwarder via ReceiverTemplate after metadata validation
    /// @param report ABI-encoded (uint8 actionType, address target, string reason)
    function _processReport(bytes calldata report) internal override {
        (uint8 actionType, address target, string memory reason) = abi.decode(
            report, (uint8, address, string)
        );

        if (actionType == 1) {
            // FLAG action
            flaggedAddresses[target] = true;
            totalFlags++;
            emit AddressFlagged(target, reason, block.timestamp);
        }

        emit ActionTaken(actionType, target, block.timestamp);
    }
}
