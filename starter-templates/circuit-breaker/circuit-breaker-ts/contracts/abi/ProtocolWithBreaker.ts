import { parseAbi } from "viem"

export const ProtocolWithBreakerAbi = parseAbi([
  "event PriceUpdated(uint256 indexed newPrice, uint256 indexed oldPrice, uint256 timestamp)",
  "function paused() view returns (bool)",
  "function lastPrice() view returns (uint256)",
  "function lastPriceTimestamp() view returns (uint256)",
  "function priceDeviationThresholdBps() view returns (uint256)",
  "function tripCount() view returns (uint256)",
])
