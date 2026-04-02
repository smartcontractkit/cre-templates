// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { ReactorConsumerABI } from './ReactorConsumer'

export type ReactorConsumerMock = {
  flaggedAddresses?: (arg0: `0x${string}`) => boolean
  getExpectedAuthor?: () => `0x${string}`
  getExpectedWorkflowId?: () => `0x${string}`
  getExpectedWorkflowName?: () => `0x${string}`
  getForwarderAddress?: () => `0x${string}`
  owner?: () => `0x${string}`
  supportsInterface?: (interfaceId: `0x${string}`) => boolean
  totalFlags?: () => bigint
} & Pick<ContractMock<typeof ReactorConsumerABI>, 'writeReport'>

export function newReactorConsumerMock(address: Address, evmMock: EvmMock): ReactorConsumerMock {
  return addContractMock(evmMock, { address, abi: ReactorConsumerABI }) as ReactorConsumerMock
}

