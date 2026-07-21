// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { OwnableABI } from './Ownable'

export type OwnableMock = {
  owner?: () => `0x${string}`
} & Pick<ContractMock<typeof OwnableABI>, 'writeReport'>

export function newOwnableMock(address: Address, evmMock: EvmMock): OwnableMock {
  return addContractMock(evmMock, { address, abi: OwnableABI }) as OwnableMock
}

