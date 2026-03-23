import { parseAbi } from "viem"

export const VaultHarvesterAbi = parseAbi([
  "function shouldHarvest() view returns (bool harvestNeeded)",
  "function totalHarvested() view returns (uint256)",
  "function lastHarvest() view returns (uint256)",
  "function harvestCount() view returns (uint256)",
  "function pendingYield() view returns (uint256)",
  "function harvestInterval() view returns (uint256)",
  "function minYieldThreshold() view returns (uint256)",
])
