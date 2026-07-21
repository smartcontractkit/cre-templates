// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ReceiverTemplate} from "./ReceiverTemplate.sol";

/**
 * @title ProtocolWithBreaker
 * @notice A protocol contract with circuit breaker functionality.
 *         CRE monitors PriceUpdated events for anomalies and can pause the protocol
 *         via signed reports when deviation thresholds are breached.
 *
 *         In production, replace updatePrice() with your actual oracle/price feed,
 *         and add protocol-specific pause logic in _processReport().
 */
contract ProtocolWithBreaker is ReceiverTemplate {
    bool public paused;
    uint256 public lastPrice;
    uint256 public lastPriceTimestamp;
    uint256 public priceDeviationThresholdBps; // basis points (e.g., 1000 = 10%)
    uint256 public tripCount;

    event PriceUpdated(
        uint256 indexed newPrice,
        uint256 indexed oldPrice,
        uint256 timestamp
    );
    event CircuitBreakerTripped(string reason, uint256 indexed tripCount, uint256 timestamp);
    event CircuitBreakerReset(uint256 timestamp);

    uint8 constant ACTION_PAUSE = 1;
    uint8 constant ACTION_UNPAUSE = 2;

    constructor(
        address forwarder,
        uint256 _initialPrice,
        uint256 _deviationThresholdBps
    ) ReceiverTemplate(forwarder) {
        lastPrice = _initialPrice;
        lastPriceTimestamp = block.timestamp;
        priceDeviationThresholdBps = _deviationThresholdBps;
    }

    /// @notice Called by CRE Forwarder via ReceiverTemplate after metadata validation
    /// @param report ABI-encoded (uint8 action, string reason)
    function _processReport(bytes calldata report) internal override {
        (uint8 action, string memory reason) = abi.decode(report, (uint8, string));

        if (action == ACTION_PAUSE) {
            paused = true;
            tripCount++;
            emit CircuitBreakerTripped(reason, tripCount, block.timestamp);
        } else if (action == ACTION_UNPAUSE) {
            paused = false;
            emit CircuitBreakerReset(block.timestamp);
        }
    }

    /// @notice Simulates a price update (in production, called by oracle or price feed)
    function updatePrice(uint256 newPrice) external {
        uint256 oldPrice = lastPrice;
        lastPrice = newPrice;
        lastPriceTimestamp = block.timestamp;
        emit PriceUpdated(newPrice, oldPrice, block.timestamp);
    }
}
