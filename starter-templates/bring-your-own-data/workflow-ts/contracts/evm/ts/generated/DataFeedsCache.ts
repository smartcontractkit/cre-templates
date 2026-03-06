// Code generated — DO NOT EDIT.
import {
  decodeEventLog,
  decodeFunctionResult,
  encodeEventTopics,
  encodeFunctionData,
  zeroAddress,
} from 'viem'
import type { Address, Hex } from 'viem'
import {
  bytesToHex,
  encodeCallMsg,
  EVMClient,
  hexToBase64,
  LAST_FINALIZED_BLOCK_NUMBER,
  prepareReportRequest,
  type EVMLog,
  type Runtime,
} from '@chainlink/cre-sdk'

export interface DecodedLog<T> extends Omit<EVMLog, 'data'> { data: T }





/**
 * Filter params for AnswerUpdated. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type AnswerUpdatedTopics = {
  current?: bigint
  roundId?: bigint
}

/**
 * Decoded AnswerUpdated event data.
 */
export type AnswerUpdatedDecoded = {
  current: bigint
  roundId: bigint
  updatedAt: bigint
}


/**
 * Filter params for BundleFeedConfigSet. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type BundleFeedConfigSetTopics = {
  dataId?: `0x${string}`
}

/**
 * Decoded BundleFeedConfigSet event data.
 */
export type BundleFeedConfigSetDecoded = {
  dataId: `0x${string}`
  decimals: readonly number[]
  description: string
  workflowMetadata: readonly { allowedSender: `0x${string}`; allowedWorkflowOwner: `0x${string}`; allowedWorkflowName: `0x${string}` }[]
}


/**
 * Filter params for BundleReportUpdated. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type BundleReportUpdatedTopics = {
  dataId?: `0x${string}`
  timestamp?: bigint
}

/**
 * Decoded BundleReportUpdated event data.
 */
export type BundleReportUpdatedDecoded = {
  dataId: `0x${string}`
  timestamp: bigint
  bundle: `0x${string}`
}


/**
 * Filter params for DecimalFeedConfigSet. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type DecimalFeedConfigSetTopics = {
  dataId?: `0x${string}`
}

/**
 * Decoded DecimalFeedConfigSet event data.
 */
export type DecimalFeedConfigSetDecoded = {
  dataId: `0x${string}`
  decimals: number
  description: string
  workflowMetadata: readonly { allowedSender: `0x${string}`; allowedWorkflowOwner: `0x${string}`; allowedWorkflowName: `0x${string}` }[]
}


/**
 * Filter params for DecimalReportUpdated. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type DecimalReportUpdatedTopics = {
  dataId?: `0x${string}`
  roundId?: bigint
  timestamp?: bigint
}

/**
 * Decoded DecimalReportUpdated event data.
 */
export type DecimalReportUpdatedDecoded = {
  dataId: `0x${string}`
  roundId: bigint
  timestamp: bigint
  answer: bigint
}


/**
 * Filter params for FeedAdminSet. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type FeedAdminSetTopics = {
  feedAdmin?: `0x${string}`
  isAdmin?: boolean
}

/**
 * Decoded FeedAdminSet event data.
 */
export type FeedAdminSetDecoded = {
  feedAdmin: `0x${string}`
  isAdmin: boolean
}


/**
 * Filter params for FeedConfigRemoved. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type FeedConfigRemovedTopics = {
  dataId?: `0x${string}`
}

/**
 * Decoded FeedConfigRemoved event data.
 */
export type FeedConfigRemovedDecoded = {
  dataId: `0x${string}`
}


/**
 * Filter params for InvalidUpdatePermission. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type InvalidUpdatePermissionTopics = {
  dataId?: `0x${string}`
}

/**
 * Decoded InvalidUpdatePermission event data.
 */
export type InvalidUpdatePermissionDecoded = {
  dataId: `0x${string}`
  sender: `0x${string}`
  workflowOwner: `0x${string}`
  workflowName: `0x${string}`
}


/**
 * Filter params for NewRound. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type NewRoundTopics = {
  roundId?: bigint
  startedBy?: `0x${string}`
}

/**
 * Decoded NewRound event data.
 */
export type NewRoundDecoded = {
  roundId: bigint
  startedBy: `0x${string}`
  startedAt: bigint
}


/**
 * Filter params for OwnershipTransferRequested. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type OwnershipTransferRequestedTopics = {
  from?: `0x${string}`
  to?: `0x${string}`
}

/**
 * Decoded OwnershipTransferRequested event data.
 */
export type OwnershipTransferRequestedDecoded = {
  from: `0x${string}`
  to: `0x${string}`
}


/**
 * Filter params for OwnershipTransferred. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type OwnershipTransferredTopics = {
  from?: `0x${string}`
  to?: `0x${string}`
}

/**
 * Decoded OwnershipTransferred event data.
 */
export type OwnershipTransferredDecoded = {
  from: `0x${string}`
  to: `0x${string}`
}


/**
 * Filter params for ProxyDataIdRemoved. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type ProxyDataIdRemovedTopics = {
  proxy?: `0x${string}`
  dataId?: `0x${string}`
}

/**
 * Decoded ProxyDataIdRemoved event data.
 */
export type ProxyDataIdRemovedDecoded = {
  proxy: `0x${string}`
  dataId: `0x${string}`
}


/**
 * Filter params for ProxyDataIdUpdated. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type ProxyDataIdUpdatedTopics = {
  proxy?: `0x${string}`
  dataId?: `0x${string}`
}

/**
 * Decoded ProxyDataIdUpdated event data.
 */
export type ProxyDataIdUpdatedDecoded = {
  proxy: `0x${string}`
  dataId: `0x${string}`
}


/**
 * Filter params for StaleBundleReport. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type StaleBundleReportTopics = {
  dataId?: `0x${string}`
}

/**
 * Decoded StaleBundleReport event data.
 */
export type StaleBundleReportDecoded = {
  dataId: `0x${string}`
  reportTimestamp: bigint
  latestTimestamp: bigint
}


/**
 * Filter params for StaleDecimalReport. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type StaleDecimalReportTopics = {
  dataId?: `0x${string}`
}

/**
 * Decoded StaleDecimalReport event data.
 */
export type StaleDecimalReportDecoded = {
  dataId: `0x${string}`
  reportTimestamp: bigint
  latestTimestamp: bigint
}


/**
 * Filter params for TokenRecovered. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type TokenRecoveredTopics = {
  token?: `0x${string}`
  to?: `0x${string}`
}

/**
 * Decoded TokenRecovered event data.
 */
export type TokenRecoveredDecoded = {
  token: `0x${string}`
  to: `0x${string}`
  amount: bigint
}


export const DataFeedsCacheABI = [{"type":"function","name":"acceptOwnership","inputs":[],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"bundleDecimals","inputs":[],"outputs":[{"name":"bundleFeedDecimals","type":"uint8[]","internalType":"uint8[]"}],"stateMutability":"view"},{"type":"function","name":"checkFeedPermission","inputs":[{"name":"dataId","type":"bytes16","internalType":"bytes16"},{"name":"workflowMetadata","type":"tuple","internalType":"structDataFeedsCache.WorkflowMetadata","components":[{"name":"allowedSender","type":"address","internalType":"address"},{"name":"allowedWorkflowOwner","type":"address","internalType":"address"},{"name":"allowedWorkflowName","type":"bytes10","internalType":"bytes10"}]}],"outputs":[{"name":"hasPermission","type":"bool","internalType":"bool"}],"stateMutability":"view"},{"type":"function","name":"decimals","inputs":[],"outputs":[{"name":"feedDecimals","type":"uint8","internalType":"uint8"}],"stateMutability":"view"},{"type":"function","name":"description","inputs":[],"outputs":[{"name":"feedDescription","type":"string","internalType":"string"}],"stateMutability":"view"},{"type":"function","name":"getAnswer","inputs":[{"name":"roundId","type":"uint256","internalType":"uint256"}],"outputs":[{"name":"answer","type":"int256","internalType":"int256"}],"stateMutability":"view"},{"type":"function","name":"getBundleDecimals","inputs":[{"name":"dataId","type":"bytes16","internalType":"bytes16"}],"outputs":[{"name":"bundleFeedDecimals","type":"uint8[]","internalType":"uint8[]"}],"stateMutability":"view"},{"type":"function","name":"getDataIdForProxy","inputs":[{"name":"proxy","type":"address","internalType":"address"}],"outputs":[{"name":"dataId","type":"bytes16","internalType":"bytes16"}],"stateMutability":"view"},{"type":"function","name":"getDecimals","inputs":[{"name":"dataId","type":"bytes16","internalType":"bytes16"}],"outputs":[{"name":"feedDecimals","type":"uint8","internalType":"uint8"}],"stateMutability":"pure"},{"type":"function","name":"getDescription","inputs":[{"name":"dataId","type":"bytes16","internalType":"bytes16"}],"outputs":[{"name":"feedDescription","type":"string","internalType":"string"}],"stateMutability":"view"},{"type":"function","name":"getFeedMetadata","inputs":[{"name":"dataId","type":"bytes16","internalType":"bytes16"},{"name":"startIndex","type":"uint256","internalType":"uint256"},{"name":"maxCount","type":"uint256","internalType":"uint256"}],"outputs":[{"name":"workflowMetadata","type":"tuple[]","internalType":"structDataFeedsCache.WorkflowMetadata[]","components":[{"name":"allowedSender","type":"address","internalType":"address"},{"name":"allowedWorkflowOwner","type":"address","internalType":"address"},{"name":"allowedWorkflowName","type":"bytes10","internalType":"bytes10"}]}],"stateMutability":"view"},{"type":"function","name":"getLatestAnswer","inputs":[{"name":"dataId","type":"bytes16","internalType":"bytes16"}],"outputs":[{"name":"answer","type":"int256","internalType":"int256"}],"stateMutability":"view"},{"type":"function","name":"getLatestBundle","inputs":[{"name":"dataId","type":"bytes16","internalType":"bytes16"}],"outputs":[{"name":"bundle","type":"bytes","internalType":"bytes"}],"stateMutability":"view"},{"type":"function","name":"getLatestBundleTimestamp","inputs":[{"name":"dataId","type":"bytes16","internalType":"bytes16"}],"outputs":[{"name":"timestamp","type":"uint256","internalType":"uint256"}],"stateMutability":"view"},{"type":"function","name":"getLatestRoundData","inputs":[{"name":"dataId","type":"bytes16","internalType":"bytes16"}],"outputs":[{"name":"id","type":"uint80","internalType":"uint80"},{"name":"answer","type":"int256","internalType":"int256"},{"name":"startedAt","type":"uint256","internalType":"uint256"},{"name":"updatedAt","type":"uint256","internalType":"uint256"},{"name":"answeredInRound","type":"uint80","internalType":"uint80"}],"stateMutability":"view"},{"type":"function","name":"getLatestTimestamp","inputs":[{"name":"dataId","type":"bytes16","internalType":"bytes16"}],"outputs":[{"name":"timestamp","type":"uint256","internalType":"uint256"}],"stateMutability":"view"},{"type":"function","name":"getRoundData","inputs":[{"name":"roundId","type":"uint80","internalType":"uint80"}],"outputs":[{"name":"id","type":"uint80","internalType":"uint80"},{"name":"answer","type":"int256","internalType":"int256"},{"name":"startedAt","type":"uint256","internalType":"uint256"},{"name":"updatedAt","type":"uint256","internalType":"uint256"},{"name":"answeredInRound","type":"uint80","internalType":"uint80"}],"stateMutability":"view"},{"type":"function","name":"getTimestamp","inputs":[{"name":"roundId","type":"uint256","internalType":"uint256"}],"outputs":[{"name":"timestamp","type":"uint256","internalType":"uint256"}],"stateMutability":"view"},{"type":"function","name":"isFeedAdmin","inputs":[{"name":"feedAdmin","type":"address","internalType":"address"}],"outputs":[{"name":"","type":"bool","internalType":"bool"}],"stateMutability":"view"},{"type":"function","name":"latestAnswer","inputs":[],"outputs":[{"name":"answer","type":"int256","internalType":"int256"}],"stateMutability":"view"},{"type":"function","name":"latestBundle","inputs":[],"outputs":[{"name":"bundle","type":"bytes","internalType":"bytes"}],"stateMutability":"view"},{"type":"function","name":"latestBundleTimestamp","inputs":[],"outputs":[{"name":"timestamp","type":"uint256","internalType":"uint256"}],"stateMutability":"view"},{"type":"function","name":"latestRound","inputs":[],"outputs":[{"name":"round","type":"uint256","internalType":"uint256"}],"stateMutability":"view"},{"type":"function","name":"latestRoundData","inputs":[],"outputs":[{"name":"id","type":"uint80","internalType":"uint80"},{"name":"answer","type":"int256","internalType":"int256"},{"name":"startedAt","type":"uint256","internalType":"uint256"},{"name":"updatedAt","type":"uint256","internalType":"uint256"},{"name":"answeredInRound","type":"uint80","internalType":"uint80"}],"stateMutability":"view"},{"type":"function","name":"latestTimestamp","inputs":[],"outputs":[{"name":"timestamp","type":"uint256","internalType":"uint256"}],"stateMutability":"view"},{"type":"function","name":"onReport","inputs":[{"name":"metadata","type":"bytes","internalType":"bytes"},{"name":"report","type":"bytes","internalType":"bytes"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"owner","inputs":[],"outputs":[{"name":"","type":"address","internalType":"address"}],"stateMutability":"view"},{"type":"function","name":"recoverTokens","inputs":[{"name":"token","type":"address","internalType":"contractIERC20"},{"name":"to","type":"address","internalType":"address"},{"name":"amount","type":"uint256","internalType":"uint256"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"removeDataIdMappingsForProxies","inputs":[{"name":"proxies","type":"address[]","internalType":"address[]"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"removeFeedConfigs","inputs":[{"name":"dataIds","type":"bytes16[]","internalType":"bytes16[]"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"setBundleFeedConfigs","inputs":[{"name":"dataIds","type":"bytes16[]","internalType":"bytes16[]"},{"name":"descriptions","type":"string[]","internalType":"string[]"},{"name":"decimalsMatrix","type":"uint8[][]","internalType":"uint8[][]"},{"name":"workflowMetadata","type":"tuple[]","internalType":"structDataFeedsCache.WorkflowMetadata[]","components":[{"name":"allowedSender","type":"address","internalType":"address"},{"name":"allowedWorkflowOwner","type":"address","internalType":"address"},{"name":"allowedWorkflowName","type":"bytes10","internalType":"bytes10"}]}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"setDecimalFeedConfigs","inputs":[{"name":"dataIds","type":"bytes16[]","internalType":"bytes16[]"},{"name":"descriptions","type":"string[]","internalType":"string[]"},{"name":"workflowMetadata","type":"tuple[]","internalType":"structDataFeedsCache.WorkflowMetadata[]","components":[{"name":"allowedSender","type":"address","internalType":"address"},{"name":"allowedWorkflowOwner","type":"address","internalType":"address"},{"name":"allowedWorkflowName","type":"bytes10","internalType":"bytes10"}]}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"setFeedAdmin","inputs":[{"name":"feedAdmin","type":"address","internalType":"address"},{"name":"isAdmin","type":"bool","internalType":"bool"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"supportsInterface","inputs":[{"name":"interfaceId","type":"bytes4","internalType":"bytes4"}],"outputs":[{"name":"","type":"bool","internalType":"bool"}],"stateMutability":"pure"},{"type":"function","name":"transferOwnership","inputs":[{"name":"to","type":"address","internalType":"address"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"typeAndVersion","inputs":[],"outputs":[{"name":"","type":"string","internalType":"string"}],"stateMutability":"view"},{"type":"function","name":"updateDataIdMappingsForProxies","inputs":[{"name":"proxies","type":"address[]","internalType":"address[]"},{"name":"dataIds","type":"bytes16[]","internalType":"bytes16[]"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"version","inputs":[],"outputs":[{"name":"","type":"uint256","internalType":"uint256"}],"stateMutability":"view"},{"type":"event","name":"AnswerUpdated","inputs":[{"name":"current","type":"int256","indexed":true,"internalType":"int256"},{"name":"roundId","type":"uint256","indexed":true,"internalType":"uint256"},{"name":"updatedAt","type":"uint256","indexed":false,"internalType":"uint256"}],"anonymous":false},{"type":"event","name":"BundleFeedConfigSet","inputs":[{"name":"dataId","type":"bytes16","indexed":true,"internalType":"bytes16"},{"name":"decimals","type":"uint8[]","indexed":false,"internalType":"uint8[]"},{"name":"description","type":"string","indexed":false,"internalType":"string"},{"name":"workflowMetadata","type":"tuple[]","indexed":false,"internalType":"structDataFeedsCache.WorkflowMetadata[]","components":[{"name":"allowedSender","type":"address","internalType":"address"},{"name":"allowedWorkflowOwner","type":"address","internalType":"address"},{"name":"allowedWorkflowName","type":"bytes10","internalType":"bytes10"}]}],"anonymous":false},{"type":"event","name":"BundleReportUpdated","inputs":[{"name":"dataId","type":"bytes16","indexed":true,"internalType":"bytes16"},{"name":"timestamp","type":"uint256","indexed":true,"internalType":"uint256"},{"name":"bundle","type":"bytes","indexed":false,"internalType":"bytes"}],"anonymous":false},{"type":"event","name":"DecimalFeedConfigSet","inputs":[{"name":"dataId","type":"bytes16","indexed":true,"internalType":"bytes16"},{"name":"decimals","type":"uint8","indexed":false,"internalType":"uint8"},{"name":"description","type":"string","indexed":false,"internalType":"string"},{"name":"workflowMetadata","type":"tuple[]","indexed":false,"internalType":"structDataFeedsCache.WorkflowMetadata[]","components":[{"name":"allowedSender","type":"address","internalType":"address"},{"name":"allowedWorkflowOwner","type":"address","internalType":"address"},{"name":"allowedWorkflowName","type":"bytes10","internalType":"bytes10"}]}],"anonymous":false},{"type":"event","name":"DecimalReportUpdated","inputs":[{"name":"dataId","type":"bytes16","indexed":true,"internalType":"bytes16"},{"name":"roundId","type":"uint256","indexed":true,"internalType":"uint256"},{"name":"timestamp","type":"uint256","indexed":true,"internalType":"uint256"},{"name":"answer","type":"uint224","indexed":false,"internalType":"uint224"}],"anonymous":false},{"type":"event","name":"FeedAdminSet","inputs":[{"name":"feedAdmin","type":"address","indexed":true,"internalType":"address"},{"name":"isAdmin","type":"bool","indexed":true,"internalType":"bool"}],"anonymous":false},{"type":"event","name":"FeedConfigRemoved","inputs":[{"name":"dataId","type":"bytes16","indexed":true,"internalType":"bytes16"}],"anonymous":false},{"type":"event","name":"InvalidUpdatePermission","inputs":[{"name":"dataId","type":"bytes16","indexed":true,"internalType":"bytes16"},{"name":"sender","type":"address","indexed":false,"internalType":"address"},{"name":"workflowOwner","type":"address","indexed":false,"internalType":"address"},{"name":"workflowName","type":"bytes10","indexed":false,"internalType":"bytes10"}],"anonymous":false},{"type":"event","name":"NewRound","inputs":[{"name":"roundId","type":"uint256","indexed":true,"internalType":"uint256"},{"name":"startedBy","type":"address","indexed":true,"internalType":"address"},{"name":"startedAt","type":"uint256","indexed":false,"internalType":"uint256"}],"anonymous":false},{"type":"event","name":"OwnershipTransferRequested","inputs":[{"name":"from","type":"address","indexed":true,"internalType":"address"},{"name":"to","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},{"type":"event","name":"OwnershipTransferred","inputs":[{"name":"from","type":"address","indexed":true,"internalType":"address"},{"name":"to","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},{"type":"event","name":"ProxyDataIdRemoved","inputs":[{"name":"proxy","type":"address","indexed":true,"internalType":"address"},{"name":"dataId","type":"bytes16","indexed":true,"internalType":"bytes16"}],"anonymous":false},{"type":"event","name":"ProxyDataIdUpdated","inputs":[{"name":"proxy","type":"address","indexed":true,"internalType":"address"},{"name":"dataId","type":"bytes16","indexed":true,"internalType":"bytes16"}],"anonymous":false},{"type":"event","name":"StaleBundleReport","inputs":[{"name":"dataId","type":"bytes16","indexed":true,"internalType":"bytes16"},{"name":"reportTimestamp","type":"uint256","indexed":false,"internalType":"uint256"},{"name":"latestTimestamp","type":"uint256","indexed":false,"internalType":"uint256"}],"anonymous":false},{"type":"event","name":"StaleDecimalReport","inputs":[{"name":"dataId","type":"bytes16","indexed":true,"internalType":"bytes16"},{"name":"reportTimestamp","type":"uint256","indexed":false,"internalType":"uint256"},{"name":"latestTimestamp","type":"uint256","indexed":false,"internalType":"uint256"}],"anonymous":false},{"type":"event","name":"TokenRecovered","inputs":[{"name":"token","type":"address","indexed":true,"internalType":"address"},{"name":"to","type":"address","indexed":true,"internalType":"address"},{"name":"amount","type":"uint256","indexed":false,"internalType":"uint256"}],"anonymous":false},{"type":"error","name":"ArrayLengthMismatch","inputs":[]},{"type":"error","name":"EmptyConfig","inputs":[]},{"type":"error","name":"ErrorSendingNative","inputs":[{"name":"to","type":"address","internalType":"address"},{"name":"amount","type":"uint256","internalType":"uint256"},{"name":"data","type":"bytes","internalType":"bytes"}]},{"type":"error","name":"FeedNotConfigured","inputs":[{"name":"dataId","type":"bytes16","internalType":"bytes16"}]},{"type":"error","name":"InsufficientBalance","inputs":[{"name":"balance","type":"uint256","internalType":"uint256"},{"name":"requiredBalance","type":"uint256","internalType":"uint256"}]},{"type":"error","name":"InvalidAddress","inputs":[{"name":"addr","type":"address","internalType":"address"}]},{"type":"error","name":"InvalidDataId","inputs":[]},{"type":"error","name":"InvalidWorkflowName","inputs":[{"name":"workflowName","type":"bytes10","internalType":"bytes10"}]},{"type":"error","name":"NoMappingForSender","inputs":[{"name":"proxy","type":"address","internalType":"address"}]},{"type":"error","name":"SafeERC20FailedOperation","inputs":[{"name":"token","type":"address","internalType":"address"}]},{"type":"error","name":"UnauthorizedCaller","inputs":[{"name":"caller","type":"address","internalType":"address"}]}] as const

export class DataFeedsCache {
  constructor(
    private readonly client: EVMClient,
    public readonly address: Address,
  ) {}

  bundleDecimals(
    runtime: Runtime<unknown>,
  ): readonly number[] {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'bundleDecimals' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'bundleDecimals' as const,
      data: bytesToHex(result.data),
    }) as readonly number[]
  }

  checkFeedPermission(
    runtime: Runtime<unknown>,
    dataId: `0x${string}`,
    workflowMetadata: { allowedSender: `0x${string}`; allowedWorkflowOwner: `0x${string}`; allowedWorkflowName: `0x${string}` },
  ): boolean {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'checkFeedPermission' as const,
      args: [dataId, workflowMetadata],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'checkFeedPermission' as const,
      data: bytesToHex(result.data),
    }) as boolean
  }

  decimals(
    runtime: Runtime<unknown>,
  ): number {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'decimals' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'decimals' as const,
      data: bytesToHex(result.data),
    }) as number
  }

  description(
    runtime: Runtime<unknown>,
  ): string {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'description' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'description' as const,
      data: bytesToHex(result.data),
    }) as string
  }

  getAnswer(
    runtime: Runtime<unknown>,
    roundId: bigint,
  ): bigint {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'getAnswer' as const,
      args: [roundId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'getAnswer' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  getBundleDecimals(
    runtime: Runtime<unknown>,
    dataId: `0x${string}`,
  ): readonly number[] {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'getBundleDecimals' as const,
      args: [dataId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'getBundleDecimals' as const,
      data: bytesToHex(result.data),
    }) as readonly number[]
  }

  getDataIdForProxy(
    runtime: Runtime<unknown>,
    proxy: `0x${string}`,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'getDataIdForProxy' as const,
      args: [proxy],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'getDataIdForProxy' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  getDecimals(
    runtime: Runtime<unknown>,
    dataId: `0x${string}`,
  ): number {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'getDecimals' as const,
      args: [dataId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'getDecimals' as const,
      data: bytesToHex(result.data),
    }) as number
  }

  getDescription(
    runtime: Runtime<unknown>,
    dataId: `0x${string}`,
  ): string {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'getDescription' as const,
      args: [dataId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'getDescription' as const,
      data: bytesToHex(result.data),
    }) as string
  }

  getFeedMetadata(
    runtime: Runtime<unknown>,
    dataId: `0x${string}`,
    startIndex: bigint,
    maxCount: bigint,
  ): readonly { allowedSender: `0x${string}`; allowedWorkflowOwner: `0x${string}`; allowedWorkflowName: `0x${string}` }[] {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'getFeedMetadata' as const,
      args: [dataId, startIndex, maxCount],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'getFeedMetadata' as const,
      data: bytesToHex(result.data),
    }) as readonly { allowedSender: `0x${string}`; allowedWorkflowOwner: `0x${string}`; allowedWorkflowName: `0x${string}` }[]
  }

  getLatestAnswer(
    runtime: Runtime<unknown>,
    dataId: `0x${string}`,
  ): bigint {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'getLatestAnswer' as const,
      args: [dataId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'getLatestAnswer' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  getLatestBundle(
    runtime: Runtime<unknown>,
    dataId: `0x${string}`,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'getLatestBundle' as const,
      args: [dataId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'getLatestBundle' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  getLatestBundleTimestamp(
    runtime: Runtime<unknown>,
    dataId: `0x${string}`,
  ): bigint {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'getLatestBundleTimestamp' as const,
      args: [dataId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'getLatestBundleTimestamp' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  getLatestRoundData(
    runtime: Runtime<unknown>,
    dataId: `0x${string}`,
  ): readonly [bigint, bigint, bigint, bigint, bigint] {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'getLatestRoundData' as const,
      args: [dataId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'getLatestRoundData' as const,
      data: bytesToHex(result.data),
    }) as readonly [bigint, bigint, bigint, bigint, bigint]
  }

  getLatestTimestamp(
    runtime: Runtime<unknown>,
    dataId: `0x${string}`,
  ): bigint {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'getLatestTimestamp' as const,
      args: [dataId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'getLatestTimestamp' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  getRoundData(
    runtime: Runtime<unknown>,
    roundId: bigint,
  ): readonly [bigint, bigint, bigint, bigint, bigint] {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'getRoundData' as const,
      args: [roundId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'getRoundData' as const,
      data: bytesToHex(result.data),
    }) as readonly [bigint, bigint, bigint, bigint, bigint]
  }

  getTimestamp(
    runtime: Runtime<unknown>,
    roundId: bigint,
  ): bigint {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'getTimestamp' as const,
      args: [roundId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'getTimestamp' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  isFeedAdmin(
    runtime: Runtime<unknown>,
    feedAdmin: `0x${string}`,
  ): boolean {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'isFeedAdmin' as const,
      args: [feedAdmin],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'isFeedAdmin' as const,
      data: bytesToHex(result.data),
    }) as boolean
  }

  latestAnswer(
    runtime: Runtime<unknown>,
  ): bigint {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'latestAnswer' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'latestAnswer' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  latestBundle(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'latestBundle' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'latestBundle' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  latestBundleTimestamp(
    runtime: Runtime<unknown>,
  ): bigint {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'latestBundleTimestamp' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'latestBundleTimestamp' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  latestRound(
    runtime: Runtime<unknown>,
  ): bigint {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'latestRound' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'latestRound' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  latestRoundData(
    runtime: Runtime<unknown>,
  ): readonly [bigint, bigint, bigint, bigint, bigint] {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'latestRoundData' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'latestRoundData' as const,
      data: bytesToHex(result.data),
    }) as readonly [bigint, bigint, bigint, bigint, bigint]
  }

  latestTimestamp(
    runtime: Runtime<unknown>,
  ): bigint {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'latestTimestamp' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'latestTimestamp' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  owner(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'owner' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'owner' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  supportsInterface(
    runtime: Runtime<unknown>,
    interfaceId: `0x${string}`,
  ): boolean {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'supportsInterface' as const,
      args: [interfaceId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'supportsInterface' as const,
      data: bytesToHex(result.data),
    }) as boolean
  }

  typeAndVersion(
    runtime: Runtime<unknown>,
  ): string {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'typeAndVersion' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'typeAndVersion' as const,
      data: bytesToHex(result.data),
    }) as string
  }

  version(
    runtime: Runtime<unknown>,
  ): bigint {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'version' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DataFeedsCacheABI,
      functionName: 'version' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  writeReportFromOnReport(
    runtime: Runtime<unknown>,
    metadata: `0x${string}`,
    report: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'onReport' as const,
      args: [metadata, report],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromRecoverTokens(
    runtime: Runtime<unknown>,
    token: `0x${string}`,
    to: `0x${string}`,
    amount: bigint,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'recoverTokens' as const,
      args: [token, to, amount],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromRemoveDataIdMappingsForProxies(
    runtime: Runtime<unknown>,
    proxies: readonly `0x${string}`[],
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'removeDataIdMappingsForProxies' as const,
      args: [proxies],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromRemoveFeedConfigs(
    runtime: Runtime<unknown>,
    dataIds: readonly `0x${string}`[],
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'removeFeedConfigs' as const,
      args: [dataIds],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromSetBundleFeedConfigs(
    runtime: Runtime<unknown>,
    dataIds: readonly `0x${string}`[],
    descriptions: readonly string[],
    decimalsMatrix: readonly readonly number[][],
    workflowMetadata: readonly { allowedSender: `0x${string}`; allowedWorkflowOwner: `0x${string}`; allowedWorkflowName: `0x${string}` }[],
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'setBundleFeedConfigs' as const,
      args: [dataIds, descriptions, decimalsMatrix, workflowMetadata],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromSetDecimalFeedConfigs(
    runtime: Runtime<unknown>,
    dataIds: readonly `0x${string}`[],
    descriptions: readonly string[],
    workflowMetadata: readonly { allowedSender: `0x${string}`; allowedWorkflowOwner: `0x${string}`; allowedWorkflowName: `0x${string}` }[],
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'setDecimalFeedConfigs' as const,
      args: [dataIds, descriptions, workflowMetadata],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromSetFeedAdmin(
    runtime: Runtime<unknown>,
    feedAdmin: `0x${string}`,
    isAdmin: boolean,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'setFeedAdmin' as const,
      args: [feedAdmin, isAdmin],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromTransferOwnership(
    runtime: Runtime<unknown>,
    to: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'transferOwnership' as const,
      args: [to],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromUpdateDataIdMappingsForProxies(
    runtime: Runtime<unknown>,
    proxies: readonly `0x${string}`[],
    dataIds: readonly `0x${string}`[],
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: DataFeedsCacheABI,
      functionName: 'updateDataIdMappingsForProxies' as const,
      args: [proxies, dataIds],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReport(
    runtime: Runtime<unknown>,
    callData: Hex,
    gasConfig?: { gasLimit?: string },
  ) {
    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  /**
   * Creates a log trigger for AnswerUpdated events.
   * The returned trigger's adapt method decodes the raw log into AnswerUpdatedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerAnswerUpdated(
    filters?: AnswerUpdatedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'AnswerUpdated' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        current: f.current,
        roundId: f.roundId,
      }
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'AnswerUpdated' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          current: f.current,
          roundId: f.roundId,
        }
        return encodeEventTopics({
          abi: DataFeedsCacheABI,
          eventName: 'AnswerUpdated' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<AnswerUpdatedDecoded> => contract.decodeAnswerUpdated(rawOutput),
    }
  }

  /**
   * Decodes a log into AnswerUpdated data, preserving all log metadata.
   */
  decodeAnswerUpdated(log: EVMLog): DecodedLog<AnswerUpdatedDecoded> {
    const decoded = decodeEventLog({
      abi: DataFeedsCacheABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as AnswerUpdatedDecoded }
  }

  /**
   * Creates a log trigger for BundleFeedConfigSet events.
   * The returned trigger's adapt method decodes the raw log into BundleFeedConfigSetDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerBundleFeedConfigSet(
    filters?: BundleFeedConfigSetTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'BundleFeedConfigSet' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        dataId: f.dataId,
      }
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'BundleFeedConfigSet' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          dataId: f.dataId,
        }
        return encodeEventTopics({
          abi: DataFeedsCacheABI,
          eventName: 'BundleFeedConfigSet' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<BundleFeedConfigSetDecoded> => contract.decodeBundleFeedConfigSet(rawOutput),
    }
  }

  /**
   * Decodes a log into BundleFeedConfigSet data, preserving all log metadata.
   */
  decodeBundleFeedConfigSet(log: EVMLog): DecodedLog<BundleFeedConfigSetDecoded> {
    const decoded = decodeEventLog({
      abi: DataFeedsCacheABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as BundleFeedConfigSetDecoded }
  }

  /**
   * Creates a log trigger for BundleReportUpdated events.
   * The returned trigger's adapt method decodes the raw log into BundleReportUpdatedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerBundleReportUpdated(
    filters?: BundleReportUpdatedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'BundleReportUpdated' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        dataId: f.dataId,
        timestamp: f.timestamp,
      }
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'BundleReportUpdated' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          dataId: f.dataId,
          timestamp: f.timestamp,
        }
        return encodeEventTopics({
          abi: DataFeedsCacheABI,
          eventName: 'BundleReportUpdated' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<BundleReportUpdatedDecoded> => contract.decodeBundleReportUpdated(rawOutput),
    }
  }

  /**
   * Decodes a log into BundleReportUpdated data, preserving all log metadata.
   */
  decodeBundleReportUpdated(log: EVMLog): DecodedLog<BundleReportUpdatedDecoded> {
    const decoded = decodeEventLog({
      abi: DataFeedsCacheABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as BundleReportUpdatedDecoded }
  }

  /**
   * Creates a log trigger for DecimalFeedConfigSet events.
   * The returned trigger's adapt method decodes the raw log into DecimalFeedConfigSetDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerDecimalFeedConfigSet(
    filters?: DecimalFeedConfigSetTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'DecimalFeedConfigSet' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        dataId: f.dataId,
      }
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'DecimalFeedConfigSet' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          dataId: f.dataId,
        }
        return encodeEventTopics({
          abi: DataFeedsCacheABI,
          eventName: 'DecimalFeedConfigSet' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<DecimalFeedConfigSetDecoded> => contract.decodeDecimalFeedConfigSet(rawOutput),
    }
  }

  /**
   * Decodes a log into DecimalFeedConfigSet data, preserving all log metadata.
   */
  decodeDecimalFeedConfigSet(log: EVMLog): DecodedLog<DecimalFeedConfigSetDecoded> {
    const decoded = decodeEventLog({
      abi: DataFeedsCacheABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as DecimalFeedConfigSetDecoded }
  }

  /**
   * Creates a log trigger for DecimalReportUpdated events.
   * The returned trigger's adapt method decodes the raw log into DecimalReportUpdatedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerDecimalReportUpdated(
    filters?: DecimalReportUpdatedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'DecimalReportUpdated' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        dataId: f.dataId,
        roundId: f.roundId,
        timestamp: f.timestamp,
      }
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'DecimalReportUpdated' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          dataId: f.dataId,
          roundId: f.roundId,
          timestamp: f.timestamp,
        }
        return encodeEventTopics({
          abi: DataFeedsCacheABI,
          eventName: 'DecimalReportUpdated' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<DecimalReportUpdatedDecoded> => contract.decodeDecimalReportUpdated(rawOutput),
    }
  }

  /**
   * Decodes a log into DecimalReportUpdated data, preserving all log metadata.
   */
  decodeDecimalReportUpdated(log: EVMLog): DecodedLog<DecimalReportUpdatedDecoded> {
    const decoded = decodeEventLog({
      abi: DataFeedsCacheABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as DecimalReportUpdatedDecoded }
  }

  /**
   * Creates a log trigger for FeedAdminSet events.
   * The returned trigger's adapt method decodes the raw log into FeedAdminSetDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerFeedAdminSet(
    filters?: FeedAdminSetTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'FeedAdminSet' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        feedAdmin: f.feedAdmin,
        isAdmin: f.isAdmin,
      }
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'FeedAdminSet' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          feedAdmin: f.feedAdmin,
          isAdmin: f.isAdmin,
        }
        return encodeEventTopics({
          abi: DataFeedsCacheABI,
          eventName: 'FeedAdminSet' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<FeedAdminSetDecoded> => contract.decodeFeedAdminSet(rawOutput),
    }
  }

  /**
   * Decodes a log into FeedAdminSet data, preserving all log metadata.
   */
  decodeFeedAdminSet(log: EVMLog): DecodedLog<FeedAdminSetDecoded> {
    const decoded = decodeEventLog({
      abi: DataFeedsCacheABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as FeedAdminSetDecoded }
  }

  /**
   * Creates a log trigger for FeedConfigRemoved events.
   * The returned trigger's adapt method decodes the raw log into FeedConfigRemovedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerFeedConfigRemoved(
    filters?: FeedConfigRemovedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'FeedConfigRemoved' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        dataId: f.dataId,
      }
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'FeedConfigRemoved' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          dataId: f.dataId,
        }
        return encodeEventTopics({
          abi: DataFeedsCacheABI,
          eventName: 'FeedConfigRemoved' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<FeedConfigRemovedDecoded> => contract.decodeFeedConfigRemoved(rawOutput),
    }
  }

  /**
   * Decodes a log into FeedConfigRemoved data, preserving all log metadata.
   */
  decodeFeedConfigRemoved(log: EVMLog): DecodedLog<FeedConfigRemovedDecoded> {
    const decoded = decodeEventLog({
      abi: DataFeedsCacheABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as FeedConfigRemovedDecoded }
  }

  /**
   * Creates a log trigger for InvalidUpdatePermission events.
   * The returned trigger's adapt method decodes the raw log into InvalidUpdatePermissionDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerInvalidUpdatePermission(
    filters?: InvalidUpdatePermissionTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'InvalidUpdatePermission' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        dataId: f.dataId,
      }
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'InvalidUpdatePermission' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          dataId: f.dataId,
        }
        return encodeEventTopics({
          abi: DataFeedsCacheABI,
          eventName: 'InvalidUpdatePermission' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<InvalidUpdatePermissionDecoded> => contract.decodeInvalidUpdatePermission(rawOutput),
    }
  }

  /**
   * Decodes a log into InvalidUpdatePermission data, preserving all log metadata.
   */
  decodeInvalidUpdatePermission(log: EVMLog): DecodedLog<InvalidUpdatePermissionDecoded> {
    const decoded = decodeEventLog({
      abi: DataFeedsCacheABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as InvalidUpdatePermissionDecoded }
  }

  /**
   * Creates a log trigger for NewRound events.
   * The returned trigger's adapt method decodes the raw log into NewRoundDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerNewRound(
    filters?: NewRoundTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'NewRound' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        roundId: f.roundId,
        startedBy: f.startedBy,
      }
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'NewRound' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          roundId: f.roundId,
          startedBy: f.startedBy,
        }
        return encodeEventTopics({
          abi: DataFeedsCacheABI,
          eventName: 'NewRound' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<NewRoundDecoded> => contract.decodeNewRound(rawOutput),
    }
  }

  /**
   * Decodes a log into NewRound data, preserving all log metadata.
   */
  decodeNewRound(log: EVMLog): DecodedLog<NewRoundDecoded> {
    const decoded = decodeEventLog({
      abi: DataFeedsCacheABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as NewRoundDecoded }
  }

  /**
   * Creates a log trigger for OwnershipTransferRequested events.
   * The returned trigger's adapt method decodes the raw log into OwnershipTransferRequestedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerOwnershipTransferRequested(
    filters?: OwnershipTransferRequestedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'OwnershipTransferRequested' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        from: f.from,
        to: f.to,
      }
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'OwnershipTransferRequested' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          from: f.from,
          to: f.to,
        }
        return encodeEventTopics({
          abi: DataFeedsCacheABI,
          eventName: 'OwnershipTransferRequested' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<OwnershipTransferRequestedDecoded> => contract.decodeOwnershipTransferRequested(rawOutput),
    }
  }

  /**
   * Decodes a log into OwnershipTransferRequested data, preserving all log metadata.
   */
  decodeOwnershipTransferRequested(log: EVMLog): DecodedLog<OwnershipTransferRequestedDecoded> {
    const decoded = decodeEventLog({
      abi: DataFeedsCacheABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as OwnershipTransferRequestedDecoded }
  }

  /**
   * Creates a log trigger for OwnershipTransferred events.
   * The returned trigger's adapt method decodes the raw log into OwnershipTransferredDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerOwnershipTransferred(
    filters?: OwnershipTransferredTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'OwnershipTransferred' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        from: f.from,
        to: f.to,
      }
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'OwnershipTransferred' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          from: f.from,
          to: f.to,
        }
        return encodeEventTopics({
          abi: DataFeedsCacheABI,
          eventName: 'OwnershipTransferred' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<OwnershipTransferredDecoded> => contract.decodeOwnershipTransferred(rawOutput),
    }
  }

  /**
   * Decodes a log into OwnershipTransferred data, preserving all log metadata.
   */
  decodeOwnershipTransferred(log: EVMLog): DecodedLog<OwnershipTransferredDecoded> {
    const decoded = decodeEventLog({
      abi: DataFeedsCacheABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as OwnershipTransferredDecoded }
  }

  /**
   * Creates a log trigger for ProxyDataIdRemoved events.
   * The returned trigger's adapt method decodes the raw log into ProxyDataIdRemovedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerProxyDataIdRemoved(
    filters?: ProxyDataIdRemovedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'ProxyDataIdRemoved' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        proxy: f.proxy,
        dataId: f.dataId,
      }
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'ProxyDataIdRemoved' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          proxy: f.proxy,
          dataId: f.dataId,
        }
        return encodeEventTopics({
          abi: DataFeedsCacheABI,
          eventName: 'ProxyDataIdRemoved' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<ProxyDataIdRemovedDecoded> => contract.decodeProxyDataIdRemoved(rawOutput),
    }
  }

  /**
   * Decodes a log into ProxyDataIdRemoved data, preserving all log metadata.
   */
  decodeProxyDataIdRemoved(log: EVMLog): DecodedLog<ProxyDataIdRemovedDecoded> {
    const decoded = decodeEventLog({
      abi: DataFeedsCacheABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as ProxyDataIdRemovedDecoded }
  }

  /**
   * Creates a log trigger for ProxyDataIdUpdated events.
   * The returned trigger's adapt method decodes the raw log into ProxyDataIdUpdatedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerProxyDataIdUpdated(
    filters?: ProxyDataIdUpdatedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'ProxyDataIdUpdated' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        proxy: f.proxy,
        dataId: f.dataId,
      }
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'ProxyDataIdUpdated' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          proxy: f.proxy,
          dataId: f.dataId,
        }
        return encodeEventTopics({
          abi: DataFeedsCacheABI,
          eventName: 'ProxyDataIdUpdated' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<ProxyDataIdUpdatedDecoded> => contract.decodeProxyDataIdUpdated(rawOutput),
    }
  }

  /**
   * Decodes a log into ProxyDataIdUpdated data, preserving all log metadata.
   */
  decodeProxyDataIdUpdated(log: EVMLog): DecodedLog<ProxyDataIdUpdatedDecoded> {
    const decoded = decodeEventLog({
      abi: DataFeedsCacheABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as ProxyDataIdUpdatedDecoded }
  }

  /**
   * Creates a log trigger for StaleBundleReport events.
   * The returned trigger's adapt method decodes the raw log into StaleBundleReportDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerStaleBundleReport(
    filters?: StaleBundleReportTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'StaleBundleReport' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        dataId: f.dataId,
      }
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'StaleBundleReport' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          dataId: f.dataId,
        }
        return encodeEventTopics({
          abi: DataFeedsCacheABI,
          eventName: 'StaleBundleReport' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<StaleBundleReportDecoded> => contract.decodeStaleBundleReport(rawOutput),
    }
  }

  /**
   * Decodes a log into StaleBundleReport data, preserving all log metadata.
   */
  decodeStaleBundleReport(log: EVMLog): DecodedLog<StaleBundleReportDecoded> {
    const decoded = decodeEventLog({
      abi: DataFeedsCacheABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as StaleBundleReportDecoded }
  }

  /**
   * Creates a log trigger for StaleDecimalReport events.
   * The returned trigger's adapt method decodes the raw log into StaleDecimalReportDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerStaleDecimalReport(
    filters?: StaleDecimalReportTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'StaleDecimalReport' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        dataId: f.dataId,
      }
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'StaleDecimalReport' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          dataId: f.dataId,
        }
        return encodeEventTopics({
          abi: DataFeedsCacheABI,
          eventName: 'StaleDecimalReport' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<StaleDecimalReportDecoded> => contract.decodeStaleDecimalReport(rawOutput),
    }
  }

  /**
   * Decodes a log into StaleDecimalReport data, preserving all log metadata.
   */
  decodeStaleDecimalReport(log: EVMLog): DecodedLog<StaleDecimalReportDecoded> {
    const decoded = decodeEventLog({
      abi: DataFeedsCacheABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as StaleDecimalReportDecoded }
  }

  /**
   * Creates a log trigger for TokenRecovered events.
   * The returned trigger's adapt method decodes the raw log into TokenRecoveredDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerTokenRecovered(
    filters?: TokenRecoveredTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'TokenRecovered' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        token: f.token,
        to: f.to,
      }
      const encoded = encodeEventTopics({
        abi: DataFeedsCacheABI,
        eventName: 'TokenRecovered' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          token: f.token,
          to: f.to,
        }
        return encodeEventTopics({
          abi: DataFeedsCacheABI,
          eventName: 'TokenRecovered' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<TokenRecoveredDecoded> => contract.decodeTokenRecovered(rawOutput),
    }
  }

  /**
   * Decodes a log into TokenRecovered data, preserving all log metadata.
   */
  decodeTokenRecovered(log: EVMLog): DecodedLog<TokenRecoveredDecoded> {
    const decoded = decodeEventLog({
      abi: DataFeedsCacheABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as TokenRecoveredDecoded }
  }
}

