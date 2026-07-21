// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { DataFeedsCacheABI } from './DataFeedsCache'

export type DataFeedsCacheMock = {
  bundleDecimals?: () => readonly number[]
  checkFeedPermission?: (dataId: `0x${string}`, workflowMetadata: { allowedSender: `0x${string}`; allowedWorkflowOwner: `0x${string}`; allowedWorkflowName: `0x${string}` }) => boolean
  decimals?: () => number
  description?: () => string
  getAnswer?: (roundId: bigint) => bigint
  getBundleDecimals?: (dataId: `0x${string}`) => readonly number[]
  getDataIdForProxy?: (proxy: `0x${string}`) => `0x${string}`
  getDecimals?: (dataId: `0x${string}`) => number
  getDescription?: (dataId: `0x${string}`) => string
  getFeedMetadata?: (dataId: `0x${string}`, startIndex: bigint, maxCount: bigint) => readonly { allowedSender: `0x${string}`; allowedWorkflowOwner: `0x${string}`; allowedWorkflowName: `0x${string}` }[]
  getLatestAnswer?: (dataId: `0x${string}`) => bigint
  getLatestBundle?: (dataId: `0x${string}`) => `0x${string}`
  getLatestBundleTimestamp?: (dataId: `0x${string}`) => bigint
  getLatestRoundData?: (dataId: `0x${string}`) => readonly [bigint, bigint, bigint, bigint, bigint]
  getLatestTimestamp?: (dataId: `0x${string}`) => bigint
  getRoundData?: (roundId: bigint) => readonly [bigint, bigint, bigint, bigint, bigint]
  getTimestamp?: (roundId: bigint) => bigint
  isFeedAdmin?: (feedAdmin: `0x${string}`) => boolean
  latestAnswer?: () => bigint
  latestBundle?: () => `0x${string}`
  latestBundleTimestamp?: () => bigint
  latestRound?: () => bigint
  latestRoundData?: () => readonly [bigint, bigint, bigint, bigint, bigint]
  latestTimestamp?: () => bigint
  owner?: () => `0x${string}`
  supportsInterface?: (interfaceId: `0x${string}`) => boolean
  typeAndVersion?: () => string
  version?: () => bigint
} & Pick<ContractMock<typeof DataFeedsCacheABI>, 'writeReport'>

export function newDataFeedsCacheMock(address: Address, evmMock: EvmMock): DataFeedsCacheMock {
  return addContractMock(evmMock, { address, abi: DataFeedsCacheABI }) as DataFeedsCacheMock
}

