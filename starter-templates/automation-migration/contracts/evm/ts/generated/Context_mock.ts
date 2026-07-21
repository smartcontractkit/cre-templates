// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { ContextABI } from './Context'

export type ContextMock = {
} & Pick<ContractMock<typeof ContextABI>, 'writeReport'>

export function newContextMock(address: Address, evmMock: EvmMock): ContextMock {
  return addContractMock(evmMock, { address, abi: ContextABI }) as ContextMock
}

