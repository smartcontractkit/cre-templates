// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { SportsMarketABI } from './SportsMarket'

export type SportsMarketMock = {
  getGame?: (gameId: bigint) => { gameId: bigint; description: string; outcome: number; settledAt: bigint }
  getForwarderAddress?: () => `0x${string}`
  getExpectedAuthor?: () => `0x${string}`
} & Pick<ContractMock<typeof SportsMarketABI>, 'writeReport'>

export function newSportsMarketMock(address: Address, evmMock: EvmMock): SportsMarketMock {
  return addContractMock(evmMock, { address, abi: SportsMarketABI }) as SportsMarketMock
}
