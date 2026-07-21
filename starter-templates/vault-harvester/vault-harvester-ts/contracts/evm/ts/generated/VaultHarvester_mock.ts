// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { VaultHarvesterABI } from './VaultHarvester'

export type VaultHarvesterMock = {
  getExpectedAuthor?: () => `0x${string}`
  getExpectedWorkflowId?: () => `0x${string}`
  getExpectedWorkflowName?: () => `0x${string}`
  getForwarderAddress?: () => `0x${string}`
  harvestCount?: () => bigint
  harvestInterval?: () => bigint
  lastHarvest?: () => bigint
  minYieldThreshold?: () => bigint
  owner?: () => `0x${string}`
  pendingYield?: () => bigint
  shouldHarvest?: () => boolean
  supportsInterface?: (interfaceId: `0x${string}`) => boolean
  totalHarvested?: () => bigint
} & Pick<ContractMock<typeof VaultHarvesterABI>, 'writeReport'>

export function newVaultHarvesterMock(address: Address, evmMock: EvmMock): VaultHarvesterMock {
  return addContractMock(evmMock, { address, abi: VaultHarvesterABI }) as VaultHarvesterMock
}

