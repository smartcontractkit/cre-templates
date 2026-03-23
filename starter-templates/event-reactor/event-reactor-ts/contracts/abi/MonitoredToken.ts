import { parseAbi } from "viem"

export const MonitoredTokenAbi = parseAbi([
  "event LargeTransfer(address indexed from, address indexed to, uint256 amount, uint256 timestamp)",
])
