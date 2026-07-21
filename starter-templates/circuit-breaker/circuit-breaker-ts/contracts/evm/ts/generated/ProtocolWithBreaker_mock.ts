// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { ProtocolWithBreakerABI } from './ProtocolWithBreaker'

export type ProtocolWithBreakerMock = {
  getExpectedAuthor?: () => `0x${string}`
  getExpectedWorkflowId?: () => `0x${string}`
  getExpectedWorkflowName?: () => `0x${string}`
  getForwarderAddress?: () => `0x${string}`
  lastPrice?: () => bigint
  lastPriceTimestamp?: () => bigint
  owner?: () => `0x${string}`
  paused?: () => boolean
  priceDeviationThresholdBps?: () => bigint
  supportsInterface?: (interfaceId: `0x${string}`) => boolean
  tripCount?: () => bigint
} & Pick<ContractMock<typeof ProtocolWithBreakerABI>, 'writeReport'>

export function newProtocolWithBreakerMock(address: Address, evmMock: EvmMock): ProtocolWithBreakerMock {
  return addContractMock(evmMock, { address, abi: ProtocolWithBreakerABI }) as ProtocolWithBreakerMock
}

