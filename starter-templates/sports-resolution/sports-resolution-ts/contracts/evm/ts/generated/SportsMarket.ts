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
 * Filter params for SettlementRequested. Only indexed fields can be used for filtering.
 */
export type SettlementRequestedTopics = {
  gameId?: bigint
}

/**
 * Decoded SettlementRequested event data.
 */
export type SettlementRequestedDecoded = {
  gameId: bigint
  description: string
}

/**
 * Filter params for GameSettled. Only indexed fields can be used for filtering.
 */
export type GameSettledTopics = {
  gameId?: bigint
}

/**
 * Decoded GameSettled event data.
 */
export type GameSettledDecoded = {
  gameId: bigint
  outcome: number
  settledAt: bigint
}

export const SportsMarketABI = [
  {"type":"constructor","inputs":[{"name":"forwarder","type":"address","internalType":"address"}],"stateMutability":"nonpayable"},
  {"type":"function","name":"requestSettlement","inputs":[{"name":"gameId","type":"uint256","internalType":"uint256"},{"name":"description","type":"string","internalType":"string"}],"outputs":[],"stateMutability":"nonpayable"},
  {"type":"function","name":"getGame","inputs":[{"name":"gameId","type":"uint256","internalType":"uint256"}],"outputs":[{"name":"","type":"tuple","internalType":"structSportsMarket.Game","components":[{"name":"gameId","type":"uint256","internalType":"uint256"},{"name":"description","type":"string","internalType":"string"},{"name":"outcome","type":"uint8","internalType":"enumSportsMarket.Outcome"},{"name":"settledAt","type":"uint256","internalType":"uint256"}]}],"stateMutability":"view"},
  {"type":"function","name":"games","inputs":[{"name":"","type":"uint256","internalType":"uint256"}],"outputs":[{"name":"gameId","type":"uint256","internalType":"uint256"},{"name":"description","type":"string","internalType":"string"},{"name":"outcome","type":"uint8","internalType":"enumSportsMarket.Outcome"},{"name":"settledAt","type":"uint256","internalType":"uint256"}],"stateMutability":"view"},
  {"type":"function","name":"onReport","inputs":[{"name":"metadata","type":"bytes","internalType":"bytes"},{"name":"report","type":"bytes","internalType":"bytes"}],"outputs":[],"stateMutability":"nonpayable"},
  {"type":"function","name":"getForwarderAddress","inputs":[],"outputs":[{"name":"","type":"address","internalType":"address"}],"stateMutability":"view"},
  {"type":"function","name":"getExpectedAuthor","inputs":[],"outputs":[{"name":"","type":"address","internalType":"address"}],"stateMutability":"view"},
  {"type":"function","name":"getExpectedWorkflowId","inputs":[],"outputs":[{"name":"","type":"bytes32","internalType":"bytes32"}],"stateMutability":"view"},
  {"type":"function","name":"getExpectedWorkflowName","inputs":[],"outputs":[{"name":"","type":"bytes10","internalType":"bytes10"}],"stateMutability":"view"},
  {"type":"function","name":"setForwarderAddress","inputs":[{"name":"_forwarder","type":"address","internalType":"address"}],"outputs":[],"stateMutability":"nonpayable"},
  {"type":"function","name":"setExpectedAuthor","inputs":[{"name":"_author","type":"address","internalType":"address"}],"outputs":[],"stateMutability":"nonpayable"},
  {"type":"function","name":"setExpectedWorkflowName","inputs":[{"name":"_name","type":"string","internalType":"string"}],"outputs":[],"stateMutability":"nonpayable"},
  {"type":"function","name":"setExpectedWorkflowId","inputs":[{"name":"_id","type":"bytes32","internalType":"bytes32"}],"outputs":[],"stateMutability":"nonpayable"},
  {"type":"function","name":"supportsInterface","inputs":[{"name":"interfaceId","type":"bytes4","internalType":"bytes4"}],"outputs":[{"name":"","type":"bool","internalType":"bool"}],"stateMutability":"view"},
  {"type":"function","name":"owner","inputs":[],"outputs":[{"name":"","type":"address","internalType":"address"}],"stateMutability":"view"},
  {"type":"function","name":"transferOwnership","inputs":[{"name":"newOwner","type":"address","internalType":"address"}],"outputs":[],"stateMutability":"nonpayable"},
  {"type":"function","name":"renounceOwnership","inputs":[],"outputs":[],"stateMutability":"nonpayable"},
  {"type":"event","name":"SettlementRequested","inputs":[{"name":"gameId","type":"uint256","indexed":true,"internalType":"uint256"},{"name":"description","type":"string","indexed":false,"internalType":"string"}],"anonymous":false},
  {"type":"event","name":"GameSettled","inputs":[{"name":"gameId","type":"uint256","indexed":true,"internalType":"uint256"},{"name":"outcome","type":"uint8","indexed":false,"internalType":"enumSportsMarket.Outcome"},{"name":"settledAt","type":"uint256","indexed":false,"internalType":"uint256"}],"anonymous":false},
  {"type":"event","name":"ForwarderAddressUpdated","inputs":[{"name":"previousForwarder","type":"address","indexed":true,"internalType":"address"},{"name":"newForwarder","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},
  {"type":"event","name":"ExpectedAuthorUpdated","inputs":[{"name":"previousAuthor","type":"address","indexed":true,"internalType":"address"},{"name":"newAuthor","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},
  {"type":"event","name":"ExpectedWorkflowNameUpdated","inputs":[{"name":"previousName","type":"bytes10","indexed":true,"internalType":"bytes10"},{"name":"newName","type":"bytes10","indexed":true,"internalType":"bytes10"}],"anonymous":false},
  {"type":"event","name":"ExpectedWorkflowIdUpdated","inputs":[{"name":"previousId","type":"bytes32","indexed":true,"internalType":"bytes32"},{"name":"newId","type":"bytes32","indexed":true,"internalType":"bytes32"}],"anonymous":false},
  {"type":"event","name":"SecurityWarning","inputs":[{"name":"message","type":"string","indexed":false,"internalType":"string"}],"anonymous":false},
  {"type":"event","name":"OwnershipTransferred","inputs":[{"name":"previousOwner","type":"address","indexed":true,"internalType":"address"},{"name":"newOwner","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},
  {"type":"error","name":"AlreadySettled","inputs":[{"name":"gameId","type":"uint256","internalType":"uint256"}]},
  {"type":"error","name":"InvalidOutcome","inputs":[{"name":"outcome","type":"uint8","internalType":"uint8"}]},
  {"type":"error","name":"InvalidForwarderAddress","inputs":[]},
  {"type":"error","name":"InvalidSender","inputs":[{"name":"sender","type":"address","internalType":"address"},{"name":"expected","type":"address","internalType":"address"}]},
  {"type":"error","name":"InvalidAuthor","inputs":[{"name":"received","type":"address","internalType":"address"},{"name":"expected","type":"address","internalType":"address"}]},
  {"type":"error","name":"InvalidWorkflowName","inputs":[{"name":"received","type":"bytes10","internalType":"bytes10"},{"name":"expected","type":"bytes10","internalType":"bytes10"}]},
  {"type":"error","name":"InvalidWorkflowId","inputs":[{"name":"received","type":"bytes32","internalType":"bytes32"},{"name":"expected","type":"bytes32","internalType":"bytes32"}]},
  {"type":"error","name":"WorkflowNameRequiresAuthorValidation","inputs":[]},
  {"type":"error","name":"OwnableInvalidOwner","inputs":[{"name":"owner","type":"address","internalType":"address"}]},
  {"type":"error","name":"OwnableUnauthorizedAccount","inputs":[{"name":"account","type":"address","internalType":"address"}]},
] as const

export class SportsMarket {
  constructor(
    private readonly client: EVMClient,
    public readonly address: Address,
  ) {}

  getGame(
    runtime: Runtime<unknown>,
    gameId: bigint,
  ): { gameId: bigint; description: string; outcome: number; settledAt: bigint } {
    const callData = encodeFunctionData({
      abi: SportsMarketABI,
      functionName: 'getGame' as const,
      args: [gameId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: SportsMarketABI,
      functionName: 'getGame' as const,
      data: bytesToHex(result.data),
    }) as { gameId: bigint; description: string; outcome: number; settledAt: bigint }
  }

  getForwarderAddress(runtime: Runtime<unknown>): `0x${string}` {
    const callData = encodeFunctionData({ abi: SportsMarketABI, functionName: 'getForwarderAddress' as const })
    const result = this.client.callContract(runtime, { call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }), blockNumber: LAST_FINALIZED_BLOCK_NUMBER }).result()
    return decodeFunctionResult({ abi: SportsMarketABI, functionName: 'getForwarderAddress' as const, data: bytesToHex(result.data) }) as `0x${string}`
  }

  getExpectedAuthor(runtime: Runtime<unknown>): `0x${string}` {
    const callData = encodeFunctionData({ abi: SportsMarketABI, functionName: 'getExpectedAuthor' as const })
    const result = this.client.callContract(runtime, { call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }), blockNumber: LAST_FINALIZED_BLOCK_NUMBER }).result()
    return decodeFunctionResult({ abi: SportsMarketABI, functionName: 'getExpectedAuthor' as const, data: bytesToHex(result.data) }) as `0x${string}`
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
   * Creates a log trigger for SettlementRequested events.
   * The returned trigger's adapt method decodes the raw log into SettlementRequestedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerSettlementRequested(
    filters?: SettlementRequestedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: SportsMarketABI,
        eventName: 'SettlementRequested' as const,
      })
      topics = encoded.map((t) => ({ values: encodeTopicValue(t) }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const encoded = encodeEventTopics({
        abi: SportsMarketABI,
        eventName: 'SettlementRequested' as const,
        args: { gameId: f.gameId },
      })
      topics = encoded.map((t) => ({ values: encodeTopicValue(t) }))
    } else {
      const allEncoded = filters.map((f) =>
        encodeEventTopics({
          abi: SportsMarketABI,
          eventName: 'SettlementRequested' as const,
          args: { gameId: f.gameId },
        })
      )
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.flatMap((row) => encodeTopicValue(row[i])))],
      }))
    }

    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
      confidence: 'CONFIDENCE_LEVEL_LATEST',
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<SettlementRequestedDecoded> =>
        contract.decodeSettlementRequested(rawOutput),
    }
  }

  /**
   * Decodes a raw log into SettlementRequested data, preserving all log metadata.
   */
  decodeSettlementRequested(log: EVMLog): DecodedLog<SettlementRequestedDecoded> {
    const decoded = decodeEventLog({
      abi: SportsMarketABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as [Hex, ...Hex[]],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as SettlementRequestedDecoded }
  }

  /**
   * Creates a log trigger for GameSettled events.
   */
  logTriggerGameSettled(filters?: GameSettledTopics[]) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({ abi: SportsMarketABI, eventName: 'GameSettled' as const })
      topics = encoded.map((t) => ({ values: encodeTopicValue(t) }))
    } else if (filters.length === 1) {
      const encoded = encodeEventTopics({ abi: SportsMarketABI, eventName: 'GameSettled' as const, args: { gameId: filters[0].gameId } })
      topics = encoded.map((t) => ({ values: encodeTopicValue(t) }))
    } else {
      const allEncoded = filters.map((f) => encodeEventTopics({ abi: SportsMarketABI, eventName: 'GameSettled' as const, args: { gameId: f.gameId } }))
      topics = allEncoded[0].map((_, i) => ({ values: [...new Set(allEncoded.flatMap((row) => encodeTopicValue(row[i])))] }))
    }
    const baseTrigger = this.client.logTrigger({ addresses: [hexToBase64(this.address)], topics })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<GameSettledDecoded> => contract.decodeGameSettled(rawOutput),
    }
  }

  decodeGameSettled(log: EVMLog): DecodedLog<GameSettledDecoded> {
    const decoded = decodeEventLog({ abi: SportsMarketABI, data: bytesToHex(log.data), topics: log.topics.map((t) => bytesToHex(t)) as [Hex, ...Hex[]] })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as GameSettledDecoded }
  }
}
