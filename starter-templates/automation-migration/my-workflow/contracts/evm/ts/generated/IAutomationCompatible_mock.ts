// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { IAutomationCompatibleABI } from './IAutomationCompatible'

export type IAutomationCompatibleMock = {
  checkLog?: (log: { index: bigint; timestamp: bigint; txHash: `0x${string}`; blockNumber: bigint; blockHash: `0x${string}`; source: `0x${string}`; topics: readonly `0x${string}`[]; data: `0x${string}` }, checkData: `0x${string}`) => readonly [boolean, `0x${string}`]
  checkUpkeep?: (checkData: `0x${string}`) => readonly [boolean, `0x${string}`]
} & Pick<ContractMock<typeof IAutomationCompatibleABI>, 'writeReport'>

export function newIAutomationCompatibleMock(address: Address, evmMock: EvmMock): IAutomationCompatibleMock {
  return addContractMock(evmMock, { address, abi: IAutomationCompatibleABI }) as IAutomationCompatibleMock
}

