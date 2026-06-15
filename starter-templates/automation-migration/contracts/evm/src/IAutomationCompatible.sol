// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title IAutomationCompatible
 * @notice Interface for Chainlink Automation compatible contracts
 */
interface IAutomationCompatible {
    struct Log {
        uint256 index;
        uint256 timestamp;
        bytes32 txHash;
        uint256 blockNumber;
        bytes32 blockHash;
        address source;
        bytes32[] topics;
        bytes data;
    }

    function checkUpkeep(bytes calldata checkData) external view returns (bool upkeepNeeded, bytes memory performData);

    function checkLog(Log calldata log, bytes calldata checkData) external view returns (bool upkeepNeeded, bytes memory performData);

    function performUpkeep(bytes calldata performData) external;
}
