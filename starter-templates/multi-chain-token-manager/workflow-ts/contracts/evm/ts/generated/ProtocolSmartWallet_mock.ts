// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { ProtocolSmartWalletABI } from './ProtocolSmartWallet'

export type ProtocolSmartWalletMock = {
  allowedCcipSenders?: (chainSelector: bigint) => `0x${string}`
  allowedKeystoneForwarders?: (keystoneForwarder: `0x${string}`) => boolean
  allowedWorkflowOwners?: (workflowOwner: `0x${string}`) => boolean
  getPoolAddress?: () => `0x${string}`
  getRouter?: () => `0x${string}`
  owner?: () => `0x${string}`
  pool?: () => `0x${string}`
  supportsInterface?: (interfaceId: `0x${string}`) => boolean
} & Pick<ContractMock<typeof ProtocolSmartWalletABI>, 'writeReport'>

export function newProtocolSmartWalletMock(address: Address, evmMock: EvmMock): ProtocolSmartWalletMock {
  return addContractMock(evmMock, { address, abi: ProtocolSmartWalletABI }) as ProtocolSmartWalletMock
}

