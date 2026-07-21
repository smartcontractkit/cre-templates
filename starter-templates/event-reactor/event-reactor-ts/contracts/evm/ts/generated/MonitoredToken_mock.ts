// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { MonitoredTokenABI } from './MonitoredToken'

export type MonitoredTokenMock = {
  largeTransferThreshold?: () => bigint
} & Pick<ContractMock<typeof MonitoredTokenABI>, 'writeReport'>

export function newMonitoredTokenMock(address: Address, evmMock: EvmMock): MonitoredTokenMock {
  return addContractMock(evmMock, { address, abi: MonitoredTokenABI }) as MonitoredTokenMock
}

