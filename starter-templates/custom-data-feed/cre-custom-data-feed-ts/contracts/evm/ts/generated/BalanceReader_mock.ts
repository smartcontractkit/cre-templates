// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { BalanceReaderABI } from './BalanceReader'

export type BalanceReaderMock = {
  getNativeBalances?: (addresses: readonly `0x${string}`[]) => readonly bigint[]
  typeAndVersion?: () => string
} & Pick<ContractMock<typeof BalanceReaderABI>, 'writeReport'>

export function newBalanceReaderMock(address: Address, evmMock: EvmMock): BalanceReaderMock {
  return addContractMock(evmMock, { address, abi: BalanceReaderABI }) as BalanceReaderMock
}

