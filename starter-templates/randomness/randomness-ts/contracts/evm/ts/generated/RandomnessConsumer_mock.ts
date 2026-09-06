// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { RandomnessConsumerABI } from './RandomnessConsumer'

export type RandomnessConsumerMock = {
  getExpectedAuthor?: () => `0x${string}`
  getExpectedWorkflowId?: () => `0x${string}`
  getExpectedWorkflowName?: () => `0x${string}`
  getForwarderAddress?: () => `0x${string}`
  nextRequestId?: () => bigint
  owner?: () => `0x${string}`
  pendingRequests?: (arg0: bigint) => boolean
  randomWords?: (arg0: bigint) => bigint
  supportsInterface?: (interfaceId: `0x${string}`) => boolean
} & Pick<ContractMock<typeof RandomnessConsumerABI>, 'writeReport'>

export function newRandomnessConsumerMock(address: Address, evmMock: EvmMock): RandomnessConsumerMock {
  return addContractMock(evmMock, { address, abi: RandomnessConsumerABI }) as RandomnessConsumerMock
}

