// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ReceiverTemplate} from "./ReceiverTemplate.sol";

/**
 * @title KeeperConsumer
 * @notice A generic keeper contract demonstrating CRE cron -> read -> check -> write pattern.
 *         CRE workflow calls `needsUpkeep()` to check if enough time has passed, and if so,
 *         sends a report that triggers execution. The contract handles its own state updates.
 *
 *         This is the simplest possible write-path template. Replace needsUpkeep() and
 *         _processReport() with your own protocol logic.
 */
contract KeeperConsumer is ReceiverTemplate {
    uint256 public counter;
    uint256 public lastExecuted;
    uint256 public interval; // seconds between allowed executions

    event KeeperExecuted(uint256 indexed counter, uint256 timestamp);

    constructor(address forwarder, uint256 _interval) ReceiverTemplate(forwarder) {
        interval = _interval;
        lastExecuted = block.timestamp;
    }

    /// @notice Called by CRE Forwarder via ReceiverTemplate after metadata validation
    /// @param report ABI-encoded (bool shouldExecute)
    function _processReport(bytes calldata report) internal override {
        (bool shouldExecute) = abi.decode(report, (bool));
        require(shouldExecute, "Execution not requested");
        counter += 1;
        lastExecuted = block.timestamp;
        emit KeeperExecuted(counter, block.timestamp);
    }

    /// @notice View function for CRE workflow to check if upkeep is needed
    function needsUpkeep() external view returns (bool upkeepNeeded) {
        upkeepNeeded = (block.timestamp - lastExecuted) >= interval;
    }
}
