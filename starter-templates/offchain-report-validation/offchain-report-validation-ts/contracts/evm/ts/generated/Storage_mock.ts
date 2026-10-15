// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { StorageABI } from './Storage'

export type StorageMock = {
  getExpectedAuthor?: () => `0x${string}`
  getExpectedWorkflowId?: () => `0x${string}`
  getExpectedWorkflowName?: () => `0x${string}`
  getForwarderAddress?: () => `0x${string}`
  owner?: () => `0x${string}`
  supportsInterface?: (interfaceId: `0x${string}`) => boolean
  value?: () => bigint
} & Pick<ContractMock<typeof StorageABI>, 'writeReport'>

export function newStorageMock(address: Address, evmMock: EvmMock): StorageMock {
  return addContractMock(evmMock, { address, abi: StorageABI }) as StorageMock
}

