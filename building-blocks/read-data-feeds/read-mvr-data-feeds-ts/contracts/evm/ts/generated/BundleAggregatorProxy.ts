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
 * Filter params for AggregatorConfirmed. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type AggregatorConfirmedTopics = {
  previous?: `0x${string}`
  latest?: `0x${string}`
}

/**
 * Decoded AggregatorConfirmed event data.
 */
export type AggregatorConfirmedDecoded = {
  previous: `0x${string}`
  latest: `0x${string}`
}


/**
 * Filter params for AggregatorProposed. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type AggregatorProposedTopics = {
  current?: `0x${string}`
  proposed?: `0x${string}`
}

/**
 * Decoded AggregatorProposed event data.
 */
export type AggregatorProposedDecoded = {
  current: `0x${string}`
  proposed: `0x${string}`
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


export const BundleAggregatorProxyABI = [{"inputs":[{"internalType":"address","name":"aggregatorAddress","type":"address"},{"internalType":"address","name":"owner","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[{"internalType":"address","name":"aggregator","type":"address"}],"name":"AggregatorNotProposed","type":"error"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previous","type":"address"},{"indexed":true,"internalType":"address","name":"latest","type":"address"}],"name":"AggregatorConfirmed","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"current","type":"address"},{"indexed":true,"internalType":"address","name":"proposed","type":"address"}],"name":"AggregatorProposed","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"}],"name":"OwnershipTransferRequested","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"inputs":[],"name":"acceptOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"aggregator","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"bundleDecimals","outputs":[{"internalType":"uint8[]","name":"decimals","type":"uint8[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"aggregatorAddress","type":"address"}],"name":"confirmAggregator","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"description","outputs":[{"internalType":"string","name":"aggregatorDescription","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"latestBundle","outputs":[{"internalType":"bytes","name":"bundle","type":"bytes"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"latestBundleTimestamp","outputs":[{"internalType":"uint256","name":"timestamp","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"aggregatorAddress","type":"address"}],"name":"proposeAggregator","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"proposedAggregator","outputs":[{"internalType":"address","name":"proposedAggregatorAddress","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"typeAndVersion","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"uint256","name":"aggregatorVersion","type":"uint256"}],"stateMutability":"view","type":"function"}] as const

export class BundleAggregatorProxy {
  constructor(
    private readonly client: EVMClient,
    public readonly address: Address,
  ) {}

  aggregator(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: BundleAggregatorProxyABI,
      functionName: 'aggregator' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: BundleAggregatorProxyABI,
      functionName: 'aggregator' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  bundleDecimals(
    runtime: Runtime<unknown>,
  ): readonly number[] {
    const callData = encodeFunctionData({
      abi: BundleAggregatorProxyABI,
      functionName: 'bundleDecimals' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: BundleAggregatorProxyABI,
      functionName: 'bundleDecimals' as const,
      data: bytesToHex(result.data),
    }) as readonly number[]
  }

  description(
    runtime: Runtime<unknown>,
  ): string {
    const callData = encodeFunctionData({
      abi: BundleAggregatorProxyABI,
      functionName: 'description' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: BundleAggregatorProxyABI,
      functionName: 'description' as const,
      data: bytesToHex(result.data),
    }) as string
  }

  latestBundle(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: BundleAggregatorProxyABI,
      functionName: 'latestBundle' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: BundleAggregatorProxyABI,
      functionName: 'latestBundle' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  latestBundleTimestamp(
    runtime: Runtime<unknown>,
  ): bigint {
    const callData = encodeFunctionData({
      abi: BundleAggregatorProxyABI,
      functionName: 'latestBundleTimestamp' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: BundleAggregatorProxyABI,
      functionName: 'latestBundleTimestamp' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  owner(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: BundleAggregatorProxyABI,
      functionName: 'owner' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: BundleAggregatorProxyABI,
      functionName: 'owner' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  proposedAggregator(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: BundleAggregatorProxyABI,
      functionName: 'proposedAggregator' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: BundleAggregatorProxyABI,
      functionName: 'proposedAggregator' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  typeAndVersion(
    runtime: Runtime<unknown>,
  ): string {
    const callData = encodeFunctionData({
      abi: BundleAggregatorProxyABI,
      functionName: 'typeAndVersion' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: BundleAggregatorProxyABI,
      functionName: 'typeAndVersion' as const,
      data: bytesToHex(result.data),
    }) as string
  }

  version(
    runtime: Runtime<unknown>,
  ): bigint {
    const callData = encodeFunctionData({
      abi: BundleAggregatorProxyABI,
      functionName: 'version' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: BundleAggregatorProxyABI,
      functionName: 'version' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  writeReportFromConfirmAggregator(
    runtime: Runtime<unknown>,
    aggregatorAddress: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: BundleAggregatorProxyABI,
      functionName: 'confirmAggregator' as const,
      args: [aggregatorAddress],
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

  writeReportFromProposeAggregator(
    runtime: Runtime<unknown>,
    aggregatorAddress: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: BundleAggregatorProxyABI,
      functionName: 'proposeAggregator' as const,
      args: [aggregatorAddress],
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
      abi: BundleAggregatorProxyABI,
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
   * Creates a log trigger for AggregatorConfirmed events.
   * The returned trigger's adapt method decodes the raw log into AggregatorConfirmedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerAggregatorConfirmed(
    filters?: AggregatorConfirmedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: BundleAggregatorProxyABI,
        eventName: 'AggregatorConfirmed' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        previous: f.previous,
        latest: f.latest,
      }
      const encoded = encodeEventTopics({
        abi: BundleAggregatorProxyABI,
        eventName: 'AggregatorConfirmed' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          previous: f.previous,
          latest: f.latest,
        }
        return encodeEventTopics({
          abi: BundleAggregatorProxyABI,
          eventName: 'AggregatorConfirmed' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<AggregatorConfirmedDecoded> => contract.decodeAggregatorConfirmed(rawOutput),
    }
  }

  /**
   * Decodes a log into AggregatorConfirmed data, preserving all log metadata.
   */
  decodeAggregatorConfirmed(log: EVMLog): DecodedLog<AggregatorConfirmedDecoded> {
    const decoded = decodeEventLog({
      abi: BundleAggregatorProxyABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as AggregatorConfirmedDecoded }
  }

  /**
   * Creates a log trigger for AggregatorProposed events.
   * The returned trigger's adapt method decodes the raw log into AggregatorProposedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerAggregatorProposed(
    filters?: AggregatorProposedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: BundleAggregatorProxyABI,
        eventName: 'AggregatorProposed' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        current: f.current,
        proposed: f.proposed,
      }
      const encoded = encodeEventTopics({
        abi: BundleAggregatorProxyABI,
        eventName: 'AggregatorProposed' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          current: f.current,
          proposed: f.proposed,
        }
        return encodeEventTopics({
          abi: BundleAggregatorProxyABI,
          eventName: 'AggregatorProposed' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<AggregatorProposedDecoded> => contract.decodeAggregatorProposed(rawOutput),
    }
  }

  /**
   * Decodes a log into AggregatorProposed data, preserving all log metadata.
   */
  decodeAggregatorProposed(log: EVMLog): DecodedLog<AggregatorProposedDecoded> {
    const decoded = decodeEventLog({
      abi: BundleAggregatorProxyABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as AggregatorProposedDecoded }
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
        abi: BundleAggregatorProxyABI,
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
        abi: BundleAggregatorProxyABI,
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
          abi: BundleAggregatorProxyABI,
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
      abi: BundleAggregatorProxyABI,
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
        abi: BundleAggregatorProxyABI,
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
        abi: BundleAggregatorProxyABI,
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
          abi: BundleAggregatorProxyABI,
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
      abi: BundleAggregatorProxyABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as OwnershipTransferredDecoded }
  }
}

