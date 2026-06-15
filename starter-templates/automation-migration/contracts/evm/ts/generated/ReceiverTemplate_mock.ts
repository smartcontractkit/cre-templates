// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { ReceiverTemplateABI } from './ReceiverTemplate'

export type ReceiverTemplateMock = {
  getExpectedAuthor?: () => `0x${string}`
  getExpectedWorkflowId?: () => `0x${string}`
  getExpectedWorkflowName?: () => `0x${string}`
  getForwarderAddress?: () => `0x${string}`
  owner?: () => `0x${string}`
  supportsInterface?: (interfaceId: `0x${string}`) => boolean
} & Pick<ContractMock<typeof ReceiverTemplateABI>, 'writeReport'>

export function newReceiverTemplateMock(address: Address, evmMock: EvmMock): ReceiverTemplateMock {
  return addContractMock(evmMock, { address, abi: ReceiverTemplateABI }) as ReceiverTemplateMock
}

