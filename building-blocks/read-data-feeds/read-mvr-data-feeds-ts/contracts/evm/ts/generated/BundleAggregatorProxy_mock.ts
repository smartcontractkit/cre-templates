// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { BundleAggregatorProxyABI } from './BundleAggregatorProxy'

export type BundleAggregatorProxyMock = {
  aggregator?: () => `0x${string}`
  bundleDecimals?: () => readonly number[]
  description?: () => string
  latestBundle?: () => `0x${string}`
  latestBundleTimestamp?: () => bigint
  owner?: () => `0x${string}`
  proposedAggregator?: () => `0x${string}`
  typeAndVersion?: () => string
  version?: () => bigint
} & Pick<ContractMock<typeof BundleAggregatorProxyABI>, 'writeReport'>

export function newBundleAggregatorProxyMock(address: Address, evmMock: EvmMock): BundleAggregatorProxyMock {
  return addContractMock(evmMock, { address, abi: BundleAggregatorProxyABI }) as BundleAggregatorProxyMock
}

