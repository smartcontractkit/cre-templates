import { parseAbi } from "viem"

export const ReactorConsumerAbi = parseAbi([
  "function flaggedAddresses(address) view returns (bool)",
  "function totalFlags() view returns (uint256)",
])
