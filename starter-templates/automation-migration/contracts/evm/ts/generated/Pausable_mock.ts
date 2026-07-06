// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { PausableABI } from './Pausable'

export type PausableMock = {
  paused?: () => boolean
} & Pick<ContractMock<typeof PausableABI>, 'writeReport'>

export function newPausableMock(address: Address, evmMock: EvmMock): PausableMock {
  return addContractMock(evmMock, { address, abi: PausableABI }) as PausableMock
}

