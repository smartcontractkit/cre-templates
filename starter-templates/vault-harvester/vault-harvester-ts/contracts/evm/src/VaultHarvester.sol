// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ReceiverTemplate} from "./ReceiverTemplate.sol";

/**
 * @title VaultHarvester
 * @notice A DeFi vault that accumulates yield and can be harvested via CRE.
 *         CRE workflow calls `shouldHarvest()` to check if enough time has passed
 *         and yield exceeds a minimum threshold, then sends a report to trigger harvest.
 *         The contract handles its own state updates.
 *
 *         Customize shouldHarvest() and _processReport() for your vault strategy.
 */
contract VaultHarvester is ReceiverTemplate {
    uint256 public totalHarvested;
    uint256 public lastHarvest;
    uint256 public harvestCount;
    uint256 public harvestInterval;    // min seconds between harvests
    uint256 public minYieldThreshold;  // min yield to justify gas cost
    uint256 public pendingYield;       // simulated pending yield

    event Harvested(uint256 indexed amount, uint256 indexed harvestCount, uint256 timestamp);
    event YieldAccrued(uint256 amount, uint256 newPending);

    constructor(
        address forwarder,
        uint256 _harvestInterval,
        uint256 _minYieldThreshold
    ) ReceiverTemplate(forwarder) {
        harvestInterval = _harvestInterval;
        minYieldThreshold = _minYieldThreshold;
        lastHarvest = block.timestamp;
    }

    /// @notice Called by CRE Forwarder via ReceiverTemplate after metadata validation
    /// @param report ABI-encoded (bool shouldExecute)
    function _processReport(bytes calldata report) internal override {
        (bool shouldExecute) = abi.decode(report, (bool));
        require(shouldExecute, "Execution not requested");

        uint256 harvested = pendingYield;
        totalHarvested += harvested;
        pendingYield = 0;
        lastHarvest = block.timestamp;
        harvestCount++;

        emit Harvested(harvested, harvestCount, block.timestamp);
    }

    /// @notice View function for CRE workflow to check if harvest is profitable
    function shouldHarvest() external view returns (bool harvestNeeded) {
        uint256 timeSinceLastHarvest = block.timestamp - lastHarvest;
        harvestNeeded = (timeSinceLastHarvest >= harvestInterval) && (pendingYield >= minYieldThreshold);
    }

    /// @notice Simulate yield accrual (in production, yield comes from the underlying strategy)
    function accrueYield(uint256 amount) external {
        pendingYield += amount;
        emit YieldAccrued(amount, pendingYield);
    }
}
