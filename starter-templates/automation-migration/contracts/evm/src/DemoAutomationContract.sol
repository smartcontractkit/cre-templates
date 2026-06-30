// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title CLACompatible (UpkeepCounterFire)
 * @notice DISCLAIMER: A test contract compatible with Chainlink Automation
 *         that demonstrates both time and custom logic (checkUpkeep) and log-based
 *         (checkLog) trigger patterns for integration with CRE.
 *         This contract is an educational example to demonstrate how to interact
 *         with Chainlink systems, products, and services. It is provided “AS IS” 
 *         and “AS AVAILABLE” without warranties of any kind, 
 *         has not been audited, and may omit checks or error handling for clarity.
 *         Do not use this code in production without performing your own audits 
 *         and applying best practices. Neither Chainlink Labs, the Chainlink Foundation,
 *         nor Chainlink node operators are responsible for unintended outputs generated
 *         due to errors in code.
---
 * @dev Deploy this contract to test your CRE migration workflows.
 *      - For cron-trigger: Uses checkUpkeep with interval-based logic
 *      - For custom-logic: Uses checkUpkeep/performUpkeep pattern
 *      - For log-trigger: Uses checkLog/performUpkeep pattern
 */
contract DemoAutomationContract {
    bool public checkFlag = false;
    bool public performFlag = true;
    uint256 public immutable interval;
    uint256 public lastTimeStamp;
    uint256 public counter = 0;

    bool comingFromCheck = false;

    event CounterIncreased(address indexed sender, uint256 newCounter);
    event CounterReset(address indexed sender, uint256 newCounter);
    event UpkeepPerformed(
        address indexed sender,
        uint256 timestamp,
        uint256 lastTimeStamp,
        bytes performData
    );

    struct Log {
        uint256 index; // Index of the log in the block
        uint256 timestamp; // Timestamp of the block containing the log
        bytes32 txHash; // Hash of the transaction containing the log
        uint256 blockNumber; // Number of the block containing the log
        bytes32 blockHash; // Hash of the block containing the log
        address source; // Address of the contract that emitted the log
        bytes32[] topics; // Indexed topics of the log
        bytes data; // Data of the log
    }

    /**
     * @notice Constructs the CLACompatible contract
     * @param updateInterval The minimum interval (in seconds) between upkeeps
     */
    constructor(uint256 updateInterval) {
        interval = updateInterval;
        lastTimeStamp = block.timestamp;
    }

    /**
     * @notice Sets the check flag for testing
     * @param _checkFlag The new check flag value
     */
    function setCheckFlag(bool _checkFlag) external {
        checkFlag = _checkFlag;
    }

    /**
     * @notice Sets the perform flag for testing
     * @param _performFlag The new perform flag value
     */
    function setPerformFlag(bool _performFlag) external {
        performFlag = _performFlag;
    }

    /**
     * @notice Manually increases the counter (useful for testing log triggers)
     */
    function increaseCounter() external {
        counter = counter + 1;
        emit CounterIncreased(msg.sender, counter);
    }

    /**
     * @notice Resets the counter to zero
     */
    function resetCounter() external {
        counter = 0;
        emit CounterReset(msg.sender, counter);
    }

    /**
     * @notice Checks if upkeep is needed based on time interval
     * @dev Called by CRE custom-logic workflow to determine if performUpkeep should be called
     * @param data Optional check data passed from the workflow
     * @return upkeepNeeded True if the interval has passed since last upkeep
     * @return performData The data to pass to performUpkeep
     */
    function checkUpkeep(
        bytes calldata data
    ) external  view returns (bool, bytes memory) {
        bool upkeepNeeded = (block.timestamp - lastTimeStamp) > interval;
        return (upkeepNeeded, data);
    }

    /**
     * @notice Checks if upkeep is needed based on a log event
     * @dev Called by CRE log-trigger workflow when a matching log is detected
     * @param log The log data from the triggering event
     * @return upkeepNeeded Always returns true (log was emitted, so action needed)
     * @return performData Encoded address extracted from log topics
     */
    function checkLog(
        Log calldata log,
        bytes memory
    ) external view returns (bool upkeepNeeded, bytes memory performData) {
        upkeepNeeded = true;
        address logSender = bytes32ToAddress(log.topics[1]);
        performData = abi.encode(logSender);
        return (upkeepNeeded, performData);
    }

    /**
     * @notice Performs the upkeep action
     * @dev Called by the Receiver contract after CRE workflow sends a report
     * @param performData Data from checkUpkeep or checkLog
     */
    function performUpkeep(bytes calldata performData) external {
        if (!performFlag) {
            require(comingFromCheck == true, "only succeed in simulation");
        }
        uint256 previousTimestamp = lastTimeStamp;
        lastTimeStamp = block.timestamp;
        emit UpkeepPerformed(
            msg.sender,
            lastTimeStamp,
            previousTimestamp,
            performData
        );
    }

    /**
     * @notice Converts a bytes32 value to an address
     * @param _address The bytes32 value to convert
     * @return The address extracted from the bytes32
     */
    function bytes32ToAddress(bytes32 _address) public pure returns (address) {
        return address(uint160(uint256(_address)));
    }
}
