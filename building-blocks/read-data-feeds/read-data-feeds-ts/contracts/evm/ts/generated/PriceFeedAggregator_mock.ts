// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { PriceFeedAggregatorABI } from './PriceFeedAggregator'

export type PriceFeedAggregatorMock = {
  accessController?: () => `0x${string}`
  aggregator?: () => `0x${string}`
  decimals?: () => number
  description?: () => string
  getAnswer?: (roundId: bigint) => bigint
  getRoundData?: (roundId: bigint) => readonly [bigint, bigint, bigint, bigint, bigint]
  getTimestamp?: (roundId: bigint) => bigint
  latestAnswer?: () => bigint
  latestRound?: () => bigint
  latestRoundData?: () => readonly [bigint, bigint, bigint, bigint, bigint]
  latestTimestamp?: () => bigint
  owner?: () => `0x${string}`
  phaseAggregators?: (arg0: number) => `0x${string}`
  phaseId?: () => number
  proposedAggregator?: () => `0x${string}`
  proposedGetRoundData?: (roundId: bigint) => readonly [bigint, bigint, bigint, bigint, bigint]
  proposedLatestRoundData?: () => readonly [bigint, bigint, bigint, bigint, bigint]
  version?: () => bigint
} & Pick<ContractMock<typeof PriceFeedAggregatorABI>, 'writeReport'>

export function newPriceFeedAggregatorMock(address: Address, evmMock: EvmMock): PriceFeedAggregatorMock {
  return addContractMock(evmMock, { address, abi: PriceFeedAggregatorABI }) as PriceFeedAggregatorMock
}

