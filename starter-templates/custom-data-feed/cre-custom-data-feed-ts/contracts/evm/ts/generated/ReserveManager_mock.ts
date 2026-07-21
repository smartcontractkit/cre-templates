// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { ReserveManagerABI } from './ReserveManager'

export type ReserveManagerMock = {
  lastTotalMinted?: () => bigint
  lastTotalReserve?: () => bigint
} & Pick<ContractMock<typeof ReserveManagerABI>, 'writeReport'>

export function newReserveManagerMock(address: Address, evmMock: EvmMock): ReserveManagerMock {
  return addContractMock(evmMock, { address, abi: ReserveManagerABI }) as ReserveManagerMock
}

