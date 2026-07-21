// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { DemoAutomationContractABI } from './DemoAutomationContract'

export type DemoAutomationContractMock = {
  bytes32ToAddress?: (address: `0x${string}`) => `0x${string}`
  checkFlag?: () => boolean
  checkLog?: (log: { index: bigint; timestamp: bigint; txHash: `0x${string}`; blockNumber: bigint; blockHash: `0x${string}`; source: `0x${string}`; topics: readonly `0x${string}`[]; data: `0x${string}` }, arg1: `0x${string}`) => readonly [boolean, `0x${string}`]
  checkUpkeep?: (data: `0x${string}`) => readonly [boolean, `0x${string}`]
  counter?: () => bigint
  interval?: () => bigint
  lastTimeStamp?: () => bigint
  performFlag?: () => boolean
} & Pick<ContractMock<typeof DemoAutomationContractABI>, 'writeReport'>

export function newDemoAutomationContractMock(address: Address, evmMock: EvmMock): DemoAutomationContractMock {
  return addContractMock(evmMock, { address, abi: DemoAutomationContractABI }) as DemoAutomationContractMock
}

