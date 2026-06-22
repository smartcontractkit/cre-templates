// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { PredictionMarketABI } from './PredictionMarket'

export type PredictionMarketMock = {
  disputeWindow?: () => bigint
  getExpectedAuthor?: () => `0x${string}`
  getExpectedWorkflowId?: () => `0x${string}`
  getExpectedWorkflowName?: () => `0x${string}`
  getForwarderAddress?: () => `0x${string}`
  getMarket?: (marketId: bigint) => { marketId: bigint; question: string; strikePrice: bigint; expirationTime: bigint; disputeDeadline: bigint; status: number; outcome: number; resolutionPrice: bigint; resolvedAt: bigint }
  getMarketStatus?: (marketId: bigint) => number
  getNextMarketId?: () => bigint
  isExpired?: (marketId: bigint) => boolean
  isResolvable?: (marketId: bigint) => boolean
  markets?: (arg0: bigint) => readonly [bigint, string, bigint, bigint, bigint, number, number, bigint, bigint]
  nextMarketId?: () => bigint
  owner?: () => `0x${string}`
  priceFeed?: () => `0x${string}`
  supportsInterface?: (interfaceId: `0x${string}`) => boolean
} & Pick<ContractMock<typeof PredictionMarketABI>, 'writeReport'>

export function newPredictionMarketMock(address: Address, evmMock: EvmMock): PredictionMarketMock {
  return addContractMock(evmMock, { address, abi: PredictionMarketABI }) as PredictionMarketMock
}

