import { parseAbi } from "viem"

export const KeeperConsumerAbi = parseAbi([
  "function needsUpkeep() view returns (bool upkeepNeeded)",
  "function counter() view returns (uint256)",
  "function lastExecuted() view returns (uint256)",
  "function interval() view returns (uint256)",
])
