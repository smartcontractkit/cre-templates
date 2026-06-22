// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ReceiverTemplate} from "./ReceiverTemplate.sol";

/**
 * @title PredictionMarket
 * @notice A simple binary prediction market resolved by Chainlink Data Feeds via CRE.
 *         Markets ask: "Will BTC be above $X by timestamp Y?"
 *         Resolution reads the on-chain Chainlink BTC/USD price feed — no AI, no off-chain APIs.
 *
 *         Three CRE workflows interact with this contract:
 *         1. Creation — opens new markets
 *         2. Resolution — settles expired markets using price feed data
 *         3. Dispute — re-checks resolution if disputed within the dispute window
 */
contract PredictionMarket is ReceiverTemplate {

    // ─── Types ──────────────────────────────────────
    enum MarketStatus { Open, Resolved, Disputed, DisputeResolved }
    enum Outcome { Unresolved, Yes, No }

    struct Market {
        uint256 marketId;
        string question;           // e.g. "Will BTC be above $100,000 by 2026-04-01?"
        uint256 strikePrice;       // price threshold in feed decimals (e.g. 100000 * 1e8)
        uint256 expirationTime;    // unix timestamp when market can be resolved
        uint256 disputeDeadline;   // unix timestamp after which disputes are no longer accepted
        MarketStatus status;
        Outcome outcome;
        int256 resolutionPrice;    // the price used to resolve (from Chainlink feed)
        uint256 resolvedAt;
    }

    // ─── State ──────────────────────────────────────
    uint256 public nextMarketId;
    mapping(uint256 => Market) public markets;
    address public priceFeed;       // Chainlink BTC/USD aggregator proxy address
    uint256 public disputeWindow;   // seconds after resolution during which disputes are accepted

    // ─── Events ─────────────────────────────────────
    event MarketCreated(
        uint256 indexed marketId,
        string question,
        uint256 strikePrice,
        uint256 expirationTime,
        uint256 disputeDeadline
    );
    event MarketResolved(
        uint256 indexed marketId,
        Outcome outcome,
        int256 resolutionPrice,
        uint256 resolvedAt
    );
    event DisputeRaised(
        uint256 indexed marketId,
        address indexed disputor,
        string reason
    );
    event DisputeResolved(
        uint256 indexed marketId,
        Outcome outcome,
        int256 newPrice,
        bool overturned
    );

    // ─── Action Types (decoded from CRE report) ────
    uint8 constant ACTION_CREATE = 1;
    uint8 constant ACTION_RESOLVE = 2;
    uint8 constant ACTION_RESOLVE_DISPUTE = 3;

    constructor(
        address forwarder,
        address _priceFeed,
        uint256 _disputeWindow
    ) ReceiverTemplate(forwarder) {
        priceFeed = _priceFeed;
        disputeWindow = _disputeWindow;
    }

    // ─── CRE Entry Point ────────────────────────────

    function _processReport(bytes calldata report) internal override {
        (uint8 action, bytes memory data) = abi.decode(report, (uint8, bytes));

        if (action == ACTION_CREATE) {
            _createMarket(data);
        } else if (action == ACTION_RESOLVE) {
            _resolveMarket(data);
        } else if (action == ACTION_RESOLVE_DISPUTE) {
            _resolveDispute(data);
        } else {
            revert("Unknown action");
        }
    }

    // ─── Action: Create Market ──────────────────────

    function _createMarket(bytes memory data) internal {
        (
            string memory question,
            uint256 strikePrice,
            uint256 expirationTime
        ) = abi.decode(data, (string, uint256, uint256));

        require(expirationTime > block.timestamp, "Expiration must be in the future");

        uint256 marketId = nextMarketId++;
        uint256 disputeDeadline = expirationTime + disputeWindow;

        markets[marketId] = Market({
            marketId: marketId,
            question: question,
            strikePrice: strikePrice,
            expirationTime: expirationTime,
            disputeDeadline: disputeDeadline,
            status: MarketStatus.Open,
            outcome: Outcome.Unresolved,
            resolutionPrice: 0,
            resolvedAt: 0
        });

        emit MarketCreated(marketId, question, strikePrice, expirationTime, disputeDeadline);
    }

    // ─── Action: Resolve Market ─────────────────────

    function _resolveMarket(bytes memory data) internal {
        (uint256 marketId, int256 price) = abi.decode(data, (uint256, int256));

        Market storage m = markets[marketId];
        require(m.status == MarketStatus.Open, "Market not open");
        require(block.timestamp >= m.expirationTime, "Market not expired");

        m.outcome = price >= int256(m.strikePrice) ? Outcome.Yes : Outcome.No;
        m.resolutionPrice = price;
        m.resolvedAt = block.timestamp;
        m.status = MarketStatus.Resolved;

        emit MarketResolved(marketId, m.outcome, price, block.timestamp);
    }

    // ─── Action: Resolve Dispute ────────────────────

    function _resolveDispute(bytes memory data) internal {
        (uint256 marketId, int256 newPrice) = abi.decode(data, (uint256, int256));

        Market storage m = markets[marketId];
        require(m.status == MarketStatus.Disputed, "Market not disputed");
        require(block.timestamp <= m.disputeDeadline, "Dispute window closed");

        Outcome newOutcome = newPrice >= int256(m.strikePrice) ? Outcome.Yes : Outcome.No;
        bool overturned = newOutcome != m.outcome;

        m.outcome = newOutcome;
        m.resolutionPrice = newPrice;
        m.resolvedAt = block.timestamp;
        m.status = MarketStatus.DisputeResolved;

        emit DisputeResolved(marketId, newOutcome, newPrice, overturned);
    }

    // ─── Public: Raise Dispute (called by users, not CRE) ──

    function raiseDispute(uint256 marketId, string calldata reason) external {
        Market storage m = markets[marketId];
        require(m.status == MarketStatus.Resolved, "Market not resolved");
        require(block.timestamp <= m.disputeDeadline, "Dispute window closed");

        m.status = MarketStatus.Disputed;
        emit DisputeRaised(marketId, msg.sender, reason);
    }

    // ─── View Functions (read by CRE workflows) ─────

    function getMarket(uint256 marketId) external view returns (Market memory) {
        return markets[marketId];
    }

    function getMarketStatus(uint256 marketId) external view returns (MarketStatus) {
        return markets[marketId].status;
    }

    function getNextMarketId() external view returns (uint256) {
        return nextMarketId;
    }

    function isExpired(uint256 marketId) external view returns (bool) {
        return block.timestamp >= markets[marketId].expirationTime;
    }

    function isResolvable(uint256 marketId) external view returns (bool) {
        Market storage m = markets[marketId];
        return m.status == MarketStatus.Open && block.timestamp >= m.expirationTime;
    }
}
