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
 * Filter params for MessageEmitted. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type MessageEmittedTopics = {
  emitter?: `0x${string}`
  timestamp?: bigint
}

/**
 * Decoded MessageEmitted event data.
 */
export type MessageEmittedDecoded = {
  emitter: `0x${string}`
  timestamp: bigint
  message: string
}


export const MessageEmitterABI = [{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"emitter","type":"address"},{"indexed":true,"internalType":"uint256","name":"timestamp","type":"uint256"},{"indexed":false,"internalType":"string","name":"message","type":"string"}],"name":"MessageEmitted","type":"event"},{"inputs":[{"internalType":"string","name":"message","type":"string"}],"name":"emitMessage","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"emitter","type":"address"}],"name":"getLastMessage","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"emitter","type":"address"},{"internalType":"uint256","name":"timestamp","type":"uint256"}],"name":"getMessage","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"typeAndVersion","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}] as const

export class MessageEmitter {
  constructor(
    private readonly client: EVMClient,
    public readonly address: Address,
  ) {}

  getLastMessage(
    runtime: Runtime<unknown>,
    emitter: `0x${string}`,
  ): string {
    const callData = encodeFunctionData({
      abi: MessageEmitterABI,
      functionName: 'getLastMessage' as const,
      args: [emitter],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: MessageEmitterABI,
      functionName: 'getLastMessage' as const,
      data: bytesToHex(result.data),
    }) as string
  }

  getMessage(
    runtime: Runtime<unknown>,
    emitter: `0x${string}`,
    timestamp: bigint,
  ): string {
    const callData = encodeFunctionData({
      abi: MessageEmitterABI,
      functionName: 'getMessage' as const,
      args: [emitter, timestamp],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: MessageEmitterABI,
      functionName: 'getMessage' as const,
      data: bytesToHex(result.data),
    }) as string
  }

  typeAndVersion(
    runtime: Runtime<unknown>,
  ): string {
    const callData = encodeFunctionData({
      abi: MessageEmitterABI,
      functionName: 'typeAndVersion' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: MessageEmitterABI,
      functionName: 'typeAndVersion' as const,
      data: bytesToHex(result.data),
    }) as string
  }

  writeReportFromEmitMessage(
    runtime: Runtime<unknown>,
    message: string,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: MessageEmitterABI,
      functionName: 'emitMessage' as const,
      args: [message],
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
   * Creates a log trigger for MessageEmitted events.
   * The returned trigger's adapt method decodes the raw log into MessageEmittedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerMessageEmitted(
    filters?: MessageEmittedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: MessageEmitterABI,
        eventName: 'MessageEmitted' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        emitter: f.emitter,
        timestamp: f.timestamp,
      }
      const encoded = encodeEventTopics({
        abi: MessageEmitterABI,
        eventName: 'MessageEmitted' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          emitter: f.emitter,
          timestamp: f.timestamp,
        }
        return encodeEventTopics({
          abi: MessageEmitterABI,
          eventName: 'MessageEmitted' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<MessageEmittedDecoded> => contract.decodeMessageEmitted(rawOutput),
    }
  }

  /**
   * Decodes a log into MessageEmitted data, preserving all log metadata.
   */
  decodeMessageEmitted(log: EVMLog): DecodedLog<MessageEmittedDecoded> {
    const decoded = decodeEventLog({
      abi: MessageEmitterABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as MessageEmittedDecoded }
  }
}

