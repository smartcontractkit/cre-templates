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

const encodeTopicValue = (t: Hex | Hex[] | null): string[] => {
  if (t == null) return []
  if (Array.isArray(t)) return t.map(hexToBase64)
  return [hexToBase64(t)]
}





/**
 * Filter params for LargeTransfer. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type LargeTransferTopics = {
  from?: `0x${string}`
  to?: `0x${string}`
}

/**
 * Decoded LargeTransfer event data.
 */
export type LargeTransferDecoded = {
  from: `0x${string}`
  to: `0x${string}`
  amount: bigint
  timestamp: bigint
}


export const MonitoredTokenABI = [{"type":"function","name":"largeTransferThreshold","inputs":[],"outputs":[{"name":"","type":"uint256","internalType":"uint256"}],"stateMutability":"view"},{"type":"function","name":"simulateLargeTransfer","inputs":[{"name":"to","type":"address","internalType":"address"},{"name":"amount","type":"uint256","internalType":"uint256"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"event","name":"LargeTransfer","inputs":[{"name":"from","type":"address","indexed":true,"internalType":"address"},{"name":"to","type":"address","indexed":true,"internalType":"address"},{"name":"amount","type":"uint256","indexed":false,"internalType":"uint256"},{"name":"timestamp","type":"uint256","indexed":false,"internalType":"uint256"}],"anonymous":false}] as const

export class MonitoredToken {
  constructor(
    private readonly client: EVMClient,
    public readonly address: Address,
  ) {}

  largeTransferThreshold(
    runtime: Runtime<unknown>,
  ): bigint {
    const callData = encodeFunctionData({
      abi: MonitoredTokenABI,
      functionName: 'largeTransferThreshold' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: MonitoredTokenABI,
      functionName: 'largeTransferThreshold' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  writeReportFromSimulateLargeTransfer(
    runtime: Runtime<unknown>,
    to: `0x${string}`,
    amount: bigint,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: MonitoredTokenABI,
      functionName: 'simulateLargeTransfer' as const,
      args: [to, amount],
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
   * Creates a log trigger for LargeTransfer events.
   * The returned trigger's adapt method decodes the raw log into LargeTransferDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerLargeTransfer(
    filters?: LargeTransferTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: MonitoredTokenABI,
        eventName: 'LargeTransfer' as const,
      })
      topics = encoded.map((t) => ({ values: encodeTopicValue(t) }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        from: f.from,
        to: f.to,
      }
      const encoded = encodeEventTopics({
        abi: MonitoredTokenABI,
        eventName: 'LargeTransfer' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: encodeTopicValue(t) }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          from: f.from,
          to: f.to,
        }
        return encodeEventTopics({
          abi: MonitoredTokenABI,
          eventName: 'LargeTransfer' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.flatMap((row) => encodeTopicValue(row[i])))],
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
      adapt: (rawOutput: EVMLog): DecodedLog<LargeTransferDecoded> => contract.decodeLargeTransfer(rawOutput),
    }
  }

  /**
   * Decodes a log into LargeTransfer data, preserving all log metadata.
   */
  decodeLargeTransfer(log: EVMLog): DecodedLog<LargeTransferDecoded> {
    const decoded = decodeEventLog({
      abi: MonitoredTokenABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as [Hex, ...Hex[]],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as LargeTransferDecoded }
  }
}

