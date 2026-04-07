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
 * Filter params for DisputeRaised. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type DisputeRaisedTopics = {
  marketId?: bigint
  disputor?: `0x${string}`
}

/**
 * Decoded DisputeRaised event data.
 */
export type DisputeRaisedDecoded = {
  marketId: bigint
  disputor: `0x${string}`
  reason: string
}


/**
 * Filter params for DisputeResolved. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type DisputeResolvedTopics = {
  marketId?: bigint
}

/**
 * Decoded DisputeResolved event data.
 */
export type DisputeResolvedDecoded = {
  marketId: bigint
  outcome: number
  newPrice: bigint
  overturned: boolean
}


/**
 * Filter params for ExpectedAuthorUpdated. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type ExpectedAuthorUpdatedTopics = {
  previousAuthor?: `0x${string}`
  newAuthor?: `0x${string}`
}

/**
 * Decoded ExpectedAuthorUpdated event data.
 */
export type ExpectedAuthorUpdatedDecoded = {
  previousAuthor: `0x${string}`
  newAuthor: `0x${string}`
}


/**
 * Filter params for ExpectedWorkflowIdUpdated. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type ExpectedWorkflowIdUpdatedTopics = {
  previousId?: `0x${string}`
  newId?: `0x${string}`
}

/**
 * Decoded ExpectedWorkflowIdUpdated event data.
 */
export type ExpectedWorkflowIdUpdatedDecoded = {
  previousId: `0x${string}`
  newId: `0x${string}`
}


/**
 * Filter params for ExpectedWorkflowNameUpdated. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type ExpectedWorkflowNameUpdatedTopics = {
  previousName?: `0x${string}`
  newName?: `0x${string}`
}

/**
 * Decoded ExpectedWorkflowNameUpdated event data.
 */
export type ExpectedWorkflowNameUpdatedDecoded = {
  previousName: `0x${string}`
  newName: `0x${string}`
}


/**
 * Filter params for ForwarderAddressUpdated. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type ForwarderAddressUpdatedTopics = {
  previousForwarder?: `0x${string}`
  newForwarder?: `0x${string}`
}

/**
 * Decoded ForwarderAddressUpdated event data.
 */
export type ForwarderAddressUpdatedDecoded = {
  previousForwarder: `0x${string}`
  newForwarder: `0x${string}`
}


/**
 * Filter params for MarketCreated. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type MarketCreatedTopics = {
  marketId?: bigint
}

/**
 * Decoded MarketCreated event data.
 */
export type MarketCreatedDecoded = {
  marketId: bigint
  question: string
  strikePrice: bigint
  expirationTime: bigint
  disputeDeadline: bigint
}


/**
 * Filter params for MarketResolved. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type MarketResolvedTopics = {
  marketId?: bigint
}

/**
 * Decoded MarketResolved event data.
 */
export type MarketResolvedDecoded = {
  marketId: bigint
  outcome: number
  resolutionPrice: bigint
  resolvedAt: bigint
}


/**
 * Filter params for OwnershipTransferred. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type OwnershipTransferredTopics = {
  previousOwner?: `0x${string}`
  newOwner?: `0x${string}`
}

/**
 * Decoded OwnershipTransferred event data.
 */
export type OwnershipTransferredDecoded = {
  previousOwner: `0x${string}`
  newOwner: `0x${string}`
}


/**
 * Filter params for SecurityWarning. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type SecurityWarningTopics = {
}

/**
 * Decoded SecurityWarning event data.
 */
export type SecurityWarningDecoded = {
  message: string
}


export const PredictionMarketABI = [{"type":"constructor","inputs":[{"name":"forwarder","type":"address","internalType":"address"},{"name":"_priceFeed","type":"address","internalType":"address"},{"name":"_disputeWindow","type":"uint256","internalType":"uint256"}],"stateMutability":"nonpayable"},{"type":"function","name":"disputeWindow","inputs":[],"outputs":[{"name":"","type":"uint256","internalType":"uint256"}],"stateMutability":"view"},{"type":"function","name":"getExpectedAuthor","inputs":[],"outputs":[{"name":"","type":"address","internalType":"address"}],"stateMutability":"view"},{"type":"function","name":"getExpectedWorkflowId","inputs":[],"outputs":[{"name":"","type":"bytes32","internalType":"bytes32"}],"stateMutability":"view"},{"type":"function","name":"getExpectedWorkflowName","inputs":[],"outputs":[{"name":"","type":"bytes10","internalType":"bytes10"}],"stateMutability":"view"},{"type":"function","name":"getForwarderAddress","inputs":[],"outputs":[{"name":"","type":"address","internalType":"address"}],"stateMutability":"view"},{"type":"function","name":"getMarket","inputs":[{"name":"marketId","type":"uint256","internalType":"uint256"}],"outputs":[{"name":"","type":"tuple","internalType":"structPredictionMarket.Market","components":[{"name":"marketId","type":"uint256","internalType":"uint256"},{"name":"question","type":"string","internalType":"string"},{"name":"strikePrice","type":"uint256","internalType":"uint256"},{"name":"expirationTime","type":"uint256","internalType":"uint256"},{"name":"disputeDeadline","type":"uint256","internalType":"uint256"},{"name":"status","type":"uint8","internalType":"enumPredictionMarket.MarketStatus"},{"name":"outcome","type":"uint8","internalType":"enumPredictionMarket.Outcome"},{"name":"resolutionPrice","type":"int256","internalType":"int256"},{"name":"resolvedAt","type":"uint256","internalType":"uint256"}]}],"stateMutability":"view"},{"type":"function","name":"getMarketStatus","inputs":[{"name":"marketId","type":"uint256","internalType":"uint256"}],"outputs":[{"name":"","type":"uint8","internalType":"enumPredictionMarket.MarketStatus"}],"stateMutability":"view"},{"type":"function","name":"getNextMarketId","inputs":[],"outputs":[{"name":"","type":"uint256","internalType":"uint256"}],"stateMutability":"view"},{"type":"function","name":"isExpired","inputs":[{"name":"marketId","type":"uint256","internalType":"uint256"}],"outputs":[{"name":"","type":"bool","internalType":"bool"}],"stateMutability":"view"},{"type":"function","name":"isResolvable","inputs":[{"name":"marketId","type":"uint256","internalType":"uint256"}],"outputs":[{"name":"","type":"bool","internalType":"bool"}],"stateMutability":"view"},{"type":"function","name":"markets","inputs":[{"name":"","type":"uint256","internalType":"uint256"}],"outputs":[{"name":"marketId","type":"uint256","internalType":"uint256"},{"name":"question","type":"string","internalType":"string"},{"name":"strikePrice","type":"uint256","internalType":"uint256"},{"name":"expirationTime","type":"uint256","internalType":"uint256"},{"name":"disputeDeadline","type":"uint256","internalType":"uint256"},{"name":"status","type":"uint8","internalType":"enumPredictionMarket.MarketStatus"},{"name":"outcome","type":"uint8","internalType":"enumPredictionMarket.Outcome"},{"name":"resolutionPrice","type":"int256","internalType":"int256"},{"name":"resolvedAt","type":"uint256","internalType":"uint256"}],"stateMutability":"view"},{"type":"function","name":"nextMarketId","inputs":[],"outputs":[{"name":"","type":"uint256","internalType":"uint256"}],"stateMutability":"view"},{"type":"function","name":"onReport","inputs":[{"name":"metadata","type":"bytes","internalType":"bytes"},{"name":"report","type":"bytes","internalType":"bytes"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"owner","inputs":[],"outputs":[{"name":"","type":"address","internalType":"address"}],"stateMutability":"view"},{"type":"function","name":"priceFeed","inputs":[],"outputs":[{"name":"","type":"address","internalType":"address"}],"stateMutability":"view"},{"type":"function","name":"raiseDispute","inputs":[{"name":"marketId","type":"uint256","internalType":"uint256"},{"name":"reason","type":"string","internalType":"string"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"renounceOwnership","inputs":[],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"setExpectedAuthor","inputs":[{"name":"_author","type":"address","internalType":"address"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"setExpectedWorkflowId","inputs":[{"name":"_id","type":"bytes32","internalType":"bytes32"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"setExpectedWorkflowName","inputs":[{"name":"_name","type":"string","internalType":"string"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"setForwarderAddress","inputs":[{"name":"_forwarder","type":"address","internalType":"address"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"supportsInterface","inputs":[{"name":"interfaceId","type":"bytes4","internalType":"bytes4"}],"outputs":[{"name":"","type":"bool","internalType":"bool"}],"stateMutability":"view"},{"type":"function","name":"transferOwnership","inputs":[{"name":"newOwner","type":"address","internalType":"address"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"event","name":"DisputeRaised","inputs":[{"name":"marketId","type":"uint256","indexed":true,"internalType":"uint256"},{"name":"disputor","type":"address","indexed":true,"internalType":"address"},{"name":"reason","type":"string","indexed":false,"internalType":"string"}],"anonymous":false},{"type":"event","name":"DisputeResolved","inputs":[{"name":"marketId","type":"uint256","indexed":true,"internalType":"uint256"},{"name":"outcome","type":"uint8","indexed":false,"internalType":"enumPredictionMarket.Outcome"},{"name":"newPrice","type":"int256","indexed":false,"internalType":"int256"},{"name":"overturned","type":"bool","indexed":false,"internalType":"bool"}],"anonymous":false},{"type":"event","name":"ExpectedAuthorUpdated","inputs":[{"name":"previousAuthor","type":"address","indexed":true,"internalType":"address"},{"name":"newAuthor","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},{"type":"event","name":"ExpectedWorkflowIdUpdated","inputs":[{"name":"previousId","type":"bytes32","indexed":true,"internalType":"bytes32"},{"name":"newId","type":"bytes32","indexed":true,"internalType":"bytes32"}],"anonymous":false},{"type":"event","name":"ExpectedWorkflowNameUpdated","inputs":[{"name":"previousName","type":"bytes10","indexed":true,"internalType":"bytes10"},{"name":"newName","type":"bytes10","indexed":true,"internalType":"bytes10"}],"anonymous":false},{"type":"event","name":"ForwarderAddressUpdated","inputs":[{"name":"previousForwarder","type":"address","indexed":true,"internalType":"address"},{"name":"newForwarder","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},{"type":"event","name":"MarketCreated","inputs":[{"name":"marketId","type":"uint256","indexed":true,"internalType":"uint256"},{"name":"question","type":"string","indexed":false,"internalType":"string"},{"name":"strikePrice","type":"uint256","indexed":false,"internalType":"uint256"},{"name":"expirationTime","type":"uint256","indexed":false,"internalType":"uint256"},{"name":"disputeDeadline","type":"uint256","indexed":false,"internalType":"uint256"}],"anonymous":false},{"type":"event","name":"MarketResolved","inputs":[{"name":"marketId","type":"uint256","indexed":true,"internalType":"uint256"},{"name":"outcome","type":"uint8","indexed":false,"internalType":"enumPredictionMarket.Outcome"},{"name":"resolutionPrice","type":"int256","indexed":false,"internalType":"int256"},{"name":"resolvedAt","type":"uint256","indexed":false,"internalType":"uint256"}],"anonymous":false},{"type":"event","name":"OwnershipTransferred","inputs":[{"name":"previousOwner","type":"address","indexed":true,"internalType":"address"},{"name":"newOwner","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},{"type":"event","name":"SecurityWarning","inputs":[{"name":"message","type":"string","indexed":false,"internalType":"string"}],"anonymous":false},{"type":"error","name":"InvalidAuthor","inputs":[{"name":"received","type":"address","internalType":"address"},{"name":"expected","type":"address","internalType":"address"}]},{"type":"error","name":"InvalidForwarderAddress","inputs":[]},{"type":"error","name":"InvalidSender","inputs":[{"name":"sender","type":"address","internalType":"address"},{"name":"expected","type":"address","internalType":"address"}]},{"type":"error","name":"InvalidWorkflowId","inputs":[{"name":"received","type":"bytes32","internalType":"bytes32"},{"name":"expected","type":"bytes32","internalType":"bytes32"}]},{"type":"error","name":"InvalidWorkflowName","inputs":[{"name":"received","type":"bytes10","internalType":"bytes10"},{"name":"expected","type":"bytes10","internalType":"bytes10"}]},{"type":"error","name":"OwnableInvalidOwner","inputs":[{"name":"owner","type":"address","internalType":"address"}]},{"type":"error","name":"OwnableUnauthorizedAccount","inputs":[{"name":"account","type":"address","internalType":"address"}]},{"type":"error","name":"WorkflowNameRequiresAuthorValidation","inputs":[]}] as const

export class PredictionMarket {
  constructor(
    private readonly client: EVMClient,
    public readonly address: Address,
  ) {}

  disputeWindow(
    runtime: Runtime<unknown>,
  ): bigint {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'disputeWindow' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: PredictionMarketABI,
      functionName: 'disputeWindow' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  getExpectedAuthor(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'getExpectedAuthor' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: PredictionMarketABI,
      functionName: 'getExpectedAuthor' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  getExpectedWorkflowId(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'getExpectedWorkflowId' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: PredictionMarketABI,
      functionName: 'getExpectedWorkflowId' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  getExpectedWorkflowName(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'getExpectedWorkflowName' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: PredictionMarketABI,
      functionName: 'getExpectedWorkflowName' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  getForwarderAddress(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'getForwarderAddress' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: PredictionMarketABI,
      functionName: 'getForwarderAddress' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  getMarket(
    runtime: Runtime<unknown>,
    marketId: bigint,
  ): { marketId: bigint; question: string; strikePrice: bigint; expirationTime: bigint; disputeDeadline: bigint; status: number; outcome: number; resolutionPrice: bigint; resolvedAt: bigint } {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'getMarket' as const,
      args: [marketId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: PredictionMarketABI,
      functionName: 'getMarket' as const,
      data: bytesToHex(result.data),
    }) as { marketId: bigint; question: string; strikePrice: bigint; expirationTime: bigint; disputeDeadline: bigint; status: number; outcome: number; resolutionPrice: bigint; resolvedAt: bigint }
  }

  getMarketStatus(
    runtime: Runtime<unknown>,
    marketId: bigint,
  ): number {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'getMarketStatus' as const,
      args: [marketId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: PredictionMarketABI,
      functionName: 'getMarketStatus' as const,
      data: bytesToHex(result.data),
    }) as number
  }

  getNextMarketId(
    runtime: Runtime<unknown>,
  ): bigint {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'getNextMarketId' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: PredictionMarketABI,
      functionName: 'getNextMarketId' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  isExpired(
    runtime: Runtime<unknown>,
    marketId: bigint,
  ): boolean {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'isExpired' as const,
      args: [marketId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: PredictionMarketABI,
      functionName: 'isExpired' as const,
      data: bytesToHex(result.data),
    }) as boolean
  }

  isResolvable(
    runtime: Runtime<unknown>,
    marketId: bigint,
  ): boolean {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'isResolvable' as const,
      args: [marketId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: PredictionMarketABI,
      functionName: 'isResolvable' as const,
      data: bytesToHex(result.data),
    }) as boolean
  }

  markets(
    runtime: Runtime<unknown>,
    arg0: bigint,
  ): readonly [bigint, string, bigint, bigint, bigint, number, number, bigint, bigint] {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'markets' as const,
      args: [arg0],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: PredictionMarketABI,
      functionName: 'markets' as const,
      data: bytesToHex(result.data),
    }) as readonly [bigint, string, bigint, bigint, bigint, number, number, bigint, bigint]
  }

  nextMarketId(
    runtime: Runtime<unknown>,
  ): bigint {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'nextMarketId' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: PredictionMarketABI,
      functionName: 'nextMarketId' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  owner(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'owner' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: PredictionMarketABI,
      functionName: 'owner' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  priceFeed(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'priceFeed' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: PredictionMarketABI,
      functionName: 'priceFeed' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  supportsInterface(
    runtime: Runtime<unknown>,
    interfaceId: `0x${string}`,
  ): boolean {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
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
      abi: PredictionMarketABI,
      functionName: 'supportsInterface' as const,
      data: bytesToHex(result.data),
    }) as boolean
  }

  writeReportFromOnReport(
    runtime: Runtime<unknown>,
    metadata: `0x${string}`,
    report: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
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

  writeReportFromRaiseDispute(
    runtime: Runtime<unknown>,
    marketId: bigint,
    reason: string,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'raiseDispute' as const,
      args: [marketId, reason],
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

  writeReportFromSetExpectedAuthor(
    runtime: Runtime<unknown>,
    author: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'setExpectedAuthor' as const,
      args: [author],
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

  writeReportFromSetExpectedWorkflowId(
    runtime: Runtime<unknown>,
    id: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'setExpectedWorkflowId' as const,
      args: [id],
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

  writeReportFromSetExpectedWorkflowName(
    runtime: Runtime<unknown>,
    name: string,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'setExpectedWorkflowName' as const,
      args: [name],
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

  writeReportFromSetForwarderAddress(
    runtime: Runtime<unknown>,
    forwarder: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'setForwarderAddress' as const,
      args: [forwarder],
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
    newOwner: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: PredictionMarketABI,
      functionName: 'transferOwnership' as const,
      args: [newOwner],
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
   * Creates a log trigger for DisputeRaised events.
   * The returned trigger's adapt method decodes the raw log into DisputeRaisedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerDisputeRaised(
    filters?: DisputeRaisedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'DisputeRaised' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        marketId: f.marketId,
        disputor: f.disputor,
      }
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'DisputeRaised' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          marketId: f.marketId,
          disputor: f.disputor,
        }
        return encodeEventTopics({
          abi: PredictionMarketABI,
          eventName: 'DisputeRaised' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<DisputeRaisedDecoded> => contract.decodeDisputeRaised(rawOutput),
    }
  }

  /**
   * Decodes a log into DisputeRaised data, preserving all log metadata.
   */
  decodeDisputeRaised(log: EVMLog): DecodedLog<DisputeRaisedDecoded> {
    const decoded = decodeEventLog({
      abi: PredictionMarketABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as DisputeRaisedDecoded }
  }

  /**
   * Creates a log trigger for DisputeResolved events.
   * The returned trigger's adapt method decodes the raw log into DisputeResolvedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerDisputeResolved(
    filters?: DisputeResolvedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'DisputeResolved' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        marketId: f.marketId,
      }
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'DisputeResolved' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          marketId: f.marketId,
        }
        return encodeEventTopics({
          abi: PredictionMarketABI,
          eventName: 'DisputeResolved' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<DisputeResolvedDecoded> => contract.decodeDisputeResolved(rawOutput),
    }
  }

  /**
   * Decodes a log into DisputeResolved data, preserving all log metadata.
   */
  decodeDisputeResolved(log: EVMLog): DecodedLog<DisputeResolvedDecoded> {
    const decoded = decodeEventLog({
      abi: PredictionMarketABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as DisputeResolvedDecoded }
  }

  /**
   * Creates a log trigger for ExpectedAuthorUpdated events.
   * The returned trigger's adapt method decodes the raw log into ExpectedAuthorUpdatedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerExpectedAuthorUpdated(
    filters?: ExpectedAuthorUpdatedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'ExpectedAuthorUpdated' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        previousAuthor: f.previousAuthor,
        newAuthor: f.newAuthor,
      }
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'ExpectedAuthorUpdated' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          previousAuthor: f.previousAuthor,
          newAuthor: f.newAuthor,
        }
        return encodeEventTopics({
          abi: PredictionMarketABI,
          eventName: 'ExpectedAuthorUpdated' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<ExpectedAuthorUpdatedDecoded> => contract.decodeExpectedAuthorUpdated(rawOutput),
    }
  }

  /**
   * Decodes a log into ExpectedAuthorUpdated data, preserving all log metadata.
   */
  decodeExpectedAuthorUpdated(log: EVMLog): DecodedLog<ExpectedAuthorUpdatedDecoded> {
    const decoded = decodeEventLog({
      abi: PredictionMarketABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as ExpectedAuthorUpdatedDecoded }
  }

  /**
   * Creates a log trigger for ExpectedWorkflowIdUpdated events.
   * The returned trigger's adapt method decodes the raw log into ExpectedWorkflowIdUpdatedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerExpectedWorkflowIdUpdated(
    filters?: ExpectedWorkflowIdUpdatedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'ExpectedWorkflowIdUpdated' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        previousId: f.previousId,
        newId: f.newId,
      }
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'ExpectedWorkflowIdUpdated' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          previousId: f.previousId,
          newId: f.newId,
        }
        return encodeEventTopics({
          abi: PredictionMarketABI,
          eventName: 'ExpectedWorkflowIdUpdated' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<ExpectedWorkflowIdUpdatedDecoded> => contract.decodeExpectedWorkflowIdUpdated(rawOutput),
    }
  }

  /**
   * Decodes a log into ExpectedWorkflowIdUpdated data, preserving all log metadata.
   */
  decodeExpectedWorkflowIdUpdated(log: EVMLog): DecodedLog<ExpectedWorkflowIdUpdatedDecoded> {
    const decoded = decodeEventLog({
      abi: PredictionMarketABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as ExpectedWorkflowIdUpdatedDecoded }
  }

  /**
   * Creates a log trigger for ExpectedWorkflowNameUpdated events.
   * The returned trigger's adapt method decodes the raw log into ExpectedWorkflowNameUpdatedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerExpectedWorkflowNameUpdated(
    filters?: ExpectedWorkflowNameUpdatedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'ExpectedWorkflowNameUpdated' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        previousName: f.previousName,
        newName: f.newName,
      }
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'ExpectedWorkflowNameUpdated' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          previousName: f.previousName,
          newName: f.newName,
        }
        return encodeEventTopics({
          abi: PredictionMarketABI,
          eventName: 'ExpectedWorkflowNameUpdated' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<ExpectedWorkflowNameUpdatedDecoded> => contract.decodeExpectedWorkflowNameUpdated(rawOutput),
    }
  }

  /**
   * Decodes a log into ExpectedWorkflowNameUpdated data, preserving all log metadata.
   */
  decodeExpectedWorkflowNameUpdated(log: EVMLog): DecodedLog<ExpectedWorkflowNameUpdatedDecoded> {
    const decoded = decodeEventLog({
      abi: PredictionMarketABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as ExpectedWorkflowNameUpdatedDecoded }
  }

  /**
   * Creates a log trigger for ForwarderAddressUpdated events.
   * The returned trigger's adapt method decodes the raw log into ForwarderAddressUpdatedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerForwarderAddressUpdated(
    filters?: ForwarderAddressUpdatedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'ForwarderAddressUpdated' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        previousForwarder: f.previousForwarder,
        newForwarder: f.newForwarder,
      }
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'ForwarderAddressUpdated' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          previousForwarder: f.previousForwarder,
          newForwarder: f.newForwarder,
        }
        return encodeEventTopics({
          abi: PredictionMarketABI,
          eventName: 'ForwarderAddressUpdated' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<ForwarderAddressUpdatedDecoded> => contract.decodeForwarderAddressUpdated(rawOutput),
    }
  }

  /**
   * Decodes a log into ForwarderAddressUpdated data, preserving all log metadata.
   */
  decodeForwarderAddressUpdated(log: EVMLog): DecodedLog<ForwarderAddressUpdatedDecoded> {
    const decoded = decodeEventLog({
      abi: PredictionMarketABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as ForwarderAddressUpdatedDecoded }
  }

  /**
   * Creates a log trigger for MarketCreated events.
   * The returned trigger's adapt method decodes the raw log into MarketCreatedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerMarketCreated(
    filters?: MarketCreatedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'MarketCreated' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        marketId: f.marketId,
      }
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'MarketCreated' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          marketId: f.marketId,
        }
        return encodeEventTopics({
          abi: PredictionMarketABI,
          eventName: 'MarketCreated' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<MarketCreatedDecoded> => contract.decodeMarketCreated(rawOutput),
    }
  }

  /**
   * Decodes a log into MarketCreated data, preserving all log metadata.
   */
  decodeMarketCreated(log: EVMLog): DecodedLog<MarketCreatedDecoded> {
    const decoded = decodeEventLog({
      abi: PredictionMarketABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as MarketCreatedDecoded }
  }

  /**
   * Creates a log trigger for MarketResolved events.
   * The returned trigger's adapt method decodes the raw log into MarketResolvedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerMarketResolved(
    filters?: MarketResolvedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'MarketResolved' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        marketId: f.marketId,
      }
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'MarketResolved' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          marketId: f.marketId,
        }
        return encodeEventTopics({
          abi: PredictionMarketABI,
          eventName: 'MarketResolved' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<MarketResolvedDecoded> => contract.decodeMarketResolved(rawOutput),
    }
  }

  /**
   * Decodes a log into MarketResolved data, preserving all log metadata.
   */
  decodeMarketResolved(log: EVMLog): DecodedLog<MarketResolvedDecoded> {
    const decoded = decodeEventLog({
      abi: PredictionMarketABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as MarketResolvedDecoded }
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
        abi: PredictionMarketABI,
        eventName: 'OwnershipTransferred' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        previousOwner: f.previousOwner,
        newOwner: f.newOwner,
      }
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'OwnershipTransferred' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          previousOwner: f.previousOwner,
          newOwner: f.newOwner,
        }
        return encodeEventTopics({
          abi: PredictionMarketABI,
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
      abi: PredictionMarketABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as OwnershipTransferredDecoded }
  }

  /**
   * Creates a log trigger for SecurityWarning events.
   * The returned trigger's adapt method decodes the raw log into SecurityWarningDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerSecurityWarning(
    filters?: SecurityWarningTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'SecurityWarning' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
      }
      const encoded = encodeEventTopics({
        abi: PredictionMarketABI,
        eventName: 'SecurityWarning' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
        }
        return encodeEventTopics({
          abi: PredictionMarketABI,
          eventName: 'SecurityWarning' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<SecurityWarningDecoded> => contract.decodeSecurityWarning(rawOutput),
    }
  }

  /**
   * Decodes a log into SecurityWarning data, preserving all log metadata.
   */
  decodeSecurityWarning(log: EVMLog): DecodedLog<SecurityWarningDecoded> {
    const decoded = decodeEventLog({
      abi: PredictionMarketABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as SecurityWarningDecoded }
  }
}

