import { parseAbi } from "viem"

export const PriceFeedAggregatorAbi = parseAbi([
  "function decimals() view returns (uint8)",
  "function latestAnswer() view returns (int256)",
  "function latestRoundData() view returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)",
])
