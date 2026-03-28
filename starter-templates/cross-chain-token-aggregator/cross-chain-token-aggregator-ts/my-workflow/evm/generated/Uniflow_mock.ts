// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { UniflowABI } from './Uniflow'

export type UniflowMock = {
  getExpectedAuthor?: () => `0x${string}`
  getExpectedWorkflowId?: () => `0x${string}`
  getExpectedWorkflowName?: () => `0x${string}`
  getForwarderAddress?: () => `0x${string}`
  owner?: () => `0x${string}`
  sAllowlistedChains?: (arg0: bigint) => boolean
  sCcipRouter?: () => `0x${string}`
  sLinkToken?: () => `0x${string}`
  sTokenConfig?: (token: `0x${string}`) => readonly [`0x${string}`, bigint]
  supportsInterface?: (interfaceId: `0x${string}`) => boolean
} & Pick<ContractMock<typeof UniflowABI>, 'writeReport'>

export function newUniflowMock(address: Address, evmMock: EvmMock): UniflowMock {
  return addContractMock(evmMock, { address, abi: UniflowABI }) as UniflowMock
}

