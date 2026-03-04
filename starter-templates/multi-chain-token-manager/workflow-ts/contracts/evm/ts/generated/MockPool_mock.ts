// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { MockPoolABI } from './MockPool'

export type MockPoolMock = {
  balanceOf?: (arg0: `0x${string}`, arg1: `0x${string}`) => bigint
  currentLiquidityRate?: (arg0: `0x${string}`) => bigint
  getReserveData?: (asset: `0x${string}`) => { configuration: { data: bigint }; liquidityIndex: bigint; currentLiquidityRate: bigint; variableBorrowIndex: bigint; currentVariableBorrowRate: bigint; currentStableBorrowRate: bigint; lastUpdateTimestamp: number; id: number; aTokenAddress: `0x${string}`; stableDebtTokenAddress: `0x${string}`; variableDebtTokenAddress: `0x${string}`; interestRateStrategyAddress: `0x${string}`; accruedToTreasury: bigint; unbacked: bigint; isolationModeTotalDebt: bigint }
  owner?: () => `0x${string}`
} & Pick<ContractMock<typeof MockPoolABI>, 'writeReport'>

export function newMockPoolMock(address: Address, evmMock: EvmMock): MockPoolMock {
  return addContractMock(evmMock, { address, abi: MockPoolABI }) as MockPoolMock
}

