import { parseAbi } from "viem"

export const RandomnessConsumerAbi = parseAbi([
  "function nextRequestId() view returns (uint256)",
  "function pendingRequests(uint256) view returns (bool)",
  "function randomWords(uint256) view returns (uint256)",
  "function requestRandomness() returns (uint256)",
])
