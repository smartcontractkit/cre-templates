import { parseAbi } from "viem"

export const PredictionMarketAbi = parseAbi([
  "function getMarket(uint256 marketId) view returns ((uint256 marketId, string question, uint256 strikePrice, uint256 expirationTime, uint256 disputeDeadline, uint8 status, uint8 outcome, int256 resolutionPrice, uint256 resolvedAt))",
  "function getMarketStatus(uint256 marketId) view returns (uint8)",
  "function getNextMarketId() view returns (uint256)",
  "function isExpired(uint256 marketId) view returns (bool)",
  "function isResolvable(uint256 marketId) view returns (bool)",
  "function priceFeed() view returns (address)",
  "function disputeWindow() view returns (uint256)",
  "event MarketCreated(uint256 indexed marketId, string question, uint256 strikePrice, uint256 expirationTime, uint256 disputeDeadline)",
  "event MarketResolved(uint256 indexed marketId, uint8 outcome, int256 resolutionPrice, uint256 resolvedAt)",
  "event DisputeRaised(uint256 indexed marketId, address indexed disputor, string reason)",
  "event DisputeResolved(uint256 indexed marketId, uint8 outcome, int256 newPrice, bool overturned)",
])
