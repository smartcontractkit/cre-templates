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
 * Filter params for RequestReserveUpdate. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type RequestReserveUpdateTopics = {
}

/**
 * Decoded RequestReserveUpdate event data.
 */
export type RequestReserveUpdateDecoded = {
  requestId: bigint
}


export const ReserveManagerABI = [{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"requestId","type":"uint256"}],"name":"RequestReserveUpdate","type":"event"},{"inputs":[],"name":"lastTotalMinted","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"lastTotalReserve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"components":[{"internalType":"uint256","name":"totalMinted","type":"uint256"},{"internalType":"uint256","name":"totalReserve","type":"uint256"}],"internalType":"structUpdateReserves","name":"updateReserves","type":"tuple"}],"name":"updateReserves","outputs":[],"stateMutability":"nonpayable","type":"function"}] as const

export class ReserveManager {
  constructor(
    private readonly client: EVMClient,
    public readonly address: Address,
  ) {}

  lastTotalMinted(
    runtime: Runtime<unknown>,
  ): bigint {
    const callData = encodeFunctionData({
      abi: ReserveManagerABI,
      functionName: 'lastTotalMinted' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: ReserveManagerABI,
      functionName: 'lastTotalMinted' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  lastTotalReserve(
    runtime: Runtime<unknown>,
  ): bigint {
    const callData = encodeFunctionData({
      abi: ReserveManagerABI,
      functionName: 'lastTotalReserve' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: ReserveManagerABI,
      functionName: 'lastTotalReserve' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  writeReportFromUpdateReserves(
    runtime: Runtime<unknown>,
    updateReserves: { totalMinted: bigint; totalReserve: bigint },
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: ReserveManagerABI,
      functionName: 'updateReserves' as const,
      args: [updateReserves],
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
   * Creates a log trigger for RequestReserveUpdate events.
   * The returned trigger's adapt method decodes the raw log into RequestReserveUpdateDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerRequestReserveUpdate(
    filters?: RequestReserveUpdateTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: ReserveManagerABI,
        eventName: 'RequestReserveUpdate' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
      }
      const encoded = encodeEventTopics({
        abi: ReserveManagerABI,
        eventName: 'RequestReserveUpdate' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
        }
        return encodeEventTopics({
          abi: ReserveManagerABI,
          eventName: 'RequestReserveUpdate' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<RequestReserveUpdateDecoded> => contract.decodeRequestReserveUpdate(rawOutput),
    }
  }

  /**
   * Decodes a log into RequestReserveUpdate data, preserving all log metadata.
   */
  decodeRequestReserveUpdate(log: EVMLog): DecodedLog<RequestReserveUpdateDecoded> {
    const decoded = decodeEventLog({
      abi: ReserveManagerABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as RequestReserveUpdateDecoded }
  }
}

