// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { KeeperConsumerABI } from './KeeperConsumer'

export type KeeperConsumerMock = {
  counter?: () => bigint
  getExpectedAuthor?: () => `0x${string}`
  getExpectedWorkflowId?: () => `0x${string}`
  getExpectedWorkflowName?: () => `0x${string}`
  getForwarderAddress?: () => `0x${string}`
  interval?: () => bigint
  lastExecuted?: () => bigint
  needsUpkeep?: () => boolean
  owner?: () => `0x${string}`
  supportsInterface?: (interfaceId: `0x${string}`) => boolean
} & Pick<ContractMock<typeof KeeperConsumerABI>, 'writeReport'>

export function newKeeperConsumerMock(address: Address, evmMock: EvmMock): KeeperConsumerMock {
  return addContractMock(evmMock, { address, abi: KeeperConsumerABI }) as KeeperConsumerMock
}

