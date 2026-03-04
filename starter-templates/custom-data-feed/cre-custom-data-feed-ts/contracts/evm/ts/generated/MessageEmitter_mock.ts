// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { MessageEmitterABI } from './MessageEmitter'

export type MessageEmitterMock = {
  getLastMessage?: (emitter: `0x${string}`) => string
  getMessage?: (emitter: `0x${string}`, timestamp: bigint) => string
  typeAndVersion?: () => string
} & Pick<ContractMock<typeof MessageEmitterABI>, 'writeReport'>

export function newMessageEmitterMock(address: Address, evmMock: EvmMock): MessageEmitterMock {
  return addContractMock(evmMock, { address, abi: MessageEmitterABI }) as MessageEmitterMock
}

