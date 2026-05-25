// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { AutomationReceiverABI } from './AutomationReceiver'

export type AutomationReceiverMock = {
  getExpectedAuthor?: () => `0x${string}`
  getExpectedWorkflowId?: () => `0x${string}`
  getExpectedWorkflowName?: () => `0x${string}`
  getForwarderAddress?: () => `0x${string}`
  owner?: () => `0x${string}`
  supportsInterface?: (interfaceId: `0x${string}`) => boolean
} & Pick<ContractMock<typeof AutomationReceiverABI>, 'writeReport'>

export function newAutomationReceiverMock(address: Address, evmMock: EvmMock): AutomationReceiverMock {
  return addContractMock(evmMock, { address, abi: AutomationReceiverABI }) as AutomationReceiverMock
}

