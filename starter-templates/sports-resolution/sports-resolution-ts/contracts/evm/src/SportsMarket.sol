// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ReceiverTemplate} from "./ReceiverTemplate.sol";

/**
 * @title SportsMarket
 * @notice Prediction market that resolves game outcomes via CRE.
 *
 *         Flow:
 *           1. Any caller invokes requestSettlement(gameId, description).
 *           2. The contract emits SettlementRequested, which triggers the CRE workflow.
 *           3. The workflow fetches from 2–3 sports APIs, aggregates results, and
 *              calls back via a cryptographically signed report.
 *           4. _processReport decodes the report and records the final outcome.
 *
 *         The pattern is domain-agnostic — swap sports APIs for financial data,
 *         election results, or weather endpoints without changing the contract.
 */
contract SportsMarket is ReceiverTemplate {

    /// @notice Outcome enum — values must match workflow constants
    enum Outcome { Unresolved, HomeWins, AwayWins, Draw }

    struct Game {
        uint256 gameId;
        string  description;
        Outcome outcome;
        uint256 settledAt;
    }

    mapping(uint256 => Game) public games;

    event SettlementRequested(uint256 indexed gameId, string description);
    event GameSettled(uint256 indexed gameId, Outcome outcome, uint256 settledAt);

    error AlreadySettled(uint256 gameId);
    error InvalidOutcome(uint8 outcome);

    constructor(address forwarder) ReceiverTemplate(forwarder) {}

    /**
     * @notice Request CRE to resolve a game's outcome.
     * @param gameId      External identifier matching your sports data provider (e.g. ESPN game ID).
     * @param description Human-readable label, e.g. "Lakers vs Celtics 2026-04-20".
     */
    function requestSettlement(uint256 gameId, string calldata description) external {
        if (games[gameId].settledAt != 0) revert AlreadySettled(gameId);
        games[gameId].gameId = gameId;
        games[gameId].description = description;
        emit SettlementRequested(gameId, description);
    }

    /**
     * @notice Receives the CRE-signed resolution report.
     * @dev    report = abi.encode(uint256 gameId, uint8 outcome)
     *         outcome: 1 = HomeWins, 2 = AwayWins, 3 = Draw
     *         The ReceiverTemplate base contract validates the Forwarder signature
     *         before this function is reached.
     */
    function _processReport(bytes calldata report) internal override {
        (uint256 gameId, uint8 outcome) = abi.decode(report, (uint256, uint8));
        if (games[gameId].settledAt != 0) revert AlreadySettled(gameId);
        if (outcome < 1 || outcome > 3) revert InvalidOutcome(outcome);

        games[gameId].outcome = Outcome(outcome);
        games[gameId].settledAt = block.timestamp;

        emit GameSettled(gameId, Outcome(outcome), block.timestamp);
    }

    /// @notice Read game state — outcome 0 means unresolved
    function getGame(uint256 gameId) external view returns (Game memory) {
        return games[gameId];
    }
}
