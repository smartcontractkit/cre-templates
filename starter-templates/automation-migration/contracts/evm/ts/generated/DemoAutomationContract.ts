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
 * Filter params for CounterIncreased. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type CounterIncreasedTopics = {
  sender?: `0x${string}`
}

/**
 * Decoded CounterIncreased event data.
 */
export type CounterIncreasedDecoded = {
  sender: `0x${string}`
  newCounter: bigint
}


/**
 * Filter params for CounterReset. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type CounterResetTopics = {
  sender?: `0x${string}`
}

/**
 * Decoded CounterReset event data.
 */
export type CounterResetDecoded = {
  sender: `0x${string}`
  newCounter: bigint
}


/**
 * Filter params for UpkeepPerformed. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type UpkeepPerformedTopics = {
  sender?: `0x${string}`
}

/**
 * Decoded UpkeepPerformed event data.
 */
export type UpkeepPerformedDecoded = {
  sender: `0x${string}`
  timestamp: bigint
  lastTimeStamp: bigint
  performData: `0x${string}`
}


export const DemoAutomationContractABI = [{"inputs":[{"internalType":"uint256","name":"updateInterval","type":"uint256"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"uint256","name":"newCounter","type":"uint256"}],"name":"CounterIncreased","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"uint256","name":"newCounter","type":"uint256"}],"name":"CounterReset","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"uint256","name":"timestamp","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"lastTimeStamp","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"performData","type":"bytes"}],"name":"UpkeepPerformed","type":"event"},{"inputs":[{"internalType":"bytes32","name":"_address","type":"bytes32"}],"name":"bytes32ToAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"pure","type":"function"},{"inputs":[],"name":"checkFlag","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"components":[{"internalType":"uint256","name":"index","type":"uint256"},{"internalType":"uint256","name":"timestamp","type":"uint256"},{"internalType":"bytes32","name":"txHash","type":"bytes32"},{"internalType":"uint256","name":"blockNumber","type":"uint256"},{"internalType":"bytes32","name":"blockHash","type":"bytes32"},{"internalType":"address","name":"source","type":"address"},{"internalType":"bytes32[]","name":"topics","type":"bytes32[]"},{"internalType":"bytes","name":"data","type":"bytes"}],"internalType":"structDemoAutomationContract.Log","name":"log","type":"tuple"},{"internalType":"bytes","name":"","type":"bytes"}],"name":"checkLog","outputs":[{"internalType":"bool","name":"upkeepNeeded","type":"bool"},{"internalType":"bytes","name":"performData","type":"bytes"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes","name":"data","type":"bytes"}],"name":"checkUpkeep","outputs":[{"internalType":"bool","name":"","type":"bool"},{"internalType":"bytes","name":"","type":"bytes"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"counter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"increaseCounter","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"interval","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"lastTimeStamp","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"performFlag","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes","name":"performData","type":"bytes"}],"name":"performUpkeep","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"resetCounter","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bool","name":"_checkFlag","type":"bool"}],"name":"setCheckFlag","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bool","name":"_performFlag","type":"bool"}],"name":"setPerformFlag","outputs":[],"stateMutability":"nonpayable","type":"function"}] as const

export class DemoAutomationContract {
  constructor(
    private readonly client: EVMClient,
    public readonly address: Address,
  ) {}

  bytes32ToAddress(
    runtime: Runtime<unknown>,
    address: `0x${string}`,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: DemoAutomationContractABI,
      functionName: 'bytes32ToAddress' as const,
      args: [address],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DemoAutomationContractABI,
      functionName: 'bytes32ToAddress' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  checkFlag(
    runtime: Runtime<unknown>,
  ): boolean {
    const callData = encodeFunctionData({
      abi: DemoAutomationContractABI,
      functionName: 'checkFlag' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DemoAutomationContractABI,
      functionName: 'checkFlag' as const,
      data: bytesToHex(result.data),
    }) as boolean
  }

  checkLog(
    runtime: Runtime<unknown>,
    log: { index: bigint; timestamp: bigint; txHash: `0x${string}`; blockNumber: bigint; blockHash: `0x${string}`; source: `0x${string}`; topics: readonly `0x${string}`[]; data: `0x${string}` },
    arg1: `0x${string}`,
  ): readonly [boolean, `0x${string}`] {
    const callData = encodeFunctionData({
      abi: DemoAutomationContractABI,
      functionName: 'checkLog' as const,
      args: [log, arg1],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DemoAutomationContractABI,
      functionName: 'checkLog' as const,
      data: bytesToHex(result.data),
    }) as readonly [boolean, `0x${string}`]
  }

  checkUpkeep(
    runtime: Runtime<unknown>,
    data: `0x${string}`,
  ): readonly [boolean, `0x${string}`] {
    const callData = encodeFunctionData({
      abi: DemoAutomationContractABI,
      functionName: 'checkUpkeep' as const,
      args: [data],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DemoAutomationContractABI,
      functionName: 'checkUpkeep' as const,
      data: bytesToHex(result.data),
    }) as readonly [boolean, `0x${string}`]
  }

  counter(
    runtime: Runtime<unknown>,
  ): bigint {
    const callData = encodeFunctionData({
      abi: DemoAutomationContractABI,
      functionName: 'counter' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DemoAutomationContractABI,
      functionName: 'counter' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  interval(
    runtime: Runtime<unknown>,
  ): bigint {
    const callData = encodeFunctionData({
      abi: DemoAutomationContractABI,
      functionName: 'interval' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DemoAutomationContractABI,
      functionName: 'interval' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  lastTimeStamp(
    runtime: Runtime<unknown>,
  ): bigint {
    const callData = encodeFunctionData({
      abi: DemoAutomationContractABI,
      functionName: 'lastTimeStamp' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DemoAutomationContractABI,
      functionName: 'lastTimeStamp' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  performFlag(
    runtime: Runtime<unknown>,
  ): boolean {
    const callData = encodeFunctionData({
      abi: DemoAutomationContractABI,
      functionName: 'performFlag' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: DemoAutomationContractABI,
      functionName: 'performFlag' as const,
      data: bytesToHex(result.data),
    }) as boolean
  }

  writeReportFromPerformUpkeep(
    runtime: Runtime<unknown>,
    performData: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: DemoAutomationContractABI,
      functionName: 'performUpkeep' as const,
      args: [performData],
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

  writeReportFromSetCheckFlag(
    runtime: Runtime<unknown>,
    checkFlag: boolean,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: DemoAutomationContractABI,
      functionName: 'setCheckFlag' as const,
      args: [checkFlag],
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

  writeReportFromSetPerformFlag(
    runtime: Runtime<unknown>,
    performFlag: boolean,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: DemoAutomationContractABI,
      functionName: 'setPerformFlag' as const,
      args: [performFlag],
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
   * Creates a log trigger for CounterIncreased events.
   * The returned trigger's adapt method decodes the raw log into CounterIncreasedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerCounterIncreased(
    filters?: CounterIncreasedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DemoAutomationContractABI,
        eventName: 'CounterIncreased' as const,
      })
      topics = encoded.map((t) => ({ values: encodeTopicValue(t) }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        sender: f.sender,
      }
      const encoded = encodeEventTopics({
        abi: DemoAutomationContractABI,
        eventName: 'CounterIncreased' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: encodeTopicValue(t) }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          sender: f.sender,
        }
        return encodeEventTopics({
          abi: DemoAutomationContractABI,
          eventName: 'CounterIncreased' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<CounterIncreasedDecoded> => contract.decodeCounterIncreased(rawOutput),
    }
  }

  /**
   * Decodes a log into CounterIncreased data, preserving all log metadata.
   */
  decodeCounterIncreased(log: EVMLog): DecodedLog<CounterIncreasedDecoded> {
    const decoded = decodeEventLog({
      abi: DemoAutomationContractABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as [Hex, ...Hex[]],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as CounterIncreasedDecoded }
  }

  /**
   * Creates a log trigger for CounterReset events.
   * The returned trigger's adapt method decodes the raw log into CounterResetDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerCounterReset(
    filters?: CounterResetTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DemoAutomationContractABI,
        eventName: 'CounterReset' as const,
      })
      topics = encoded.map((t) => ({ values: encodeTopicValue(t) }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        sender: f.sender,
      }
      const encoded = encodeEventTopics({
        abi: DemoAutomationContractABI,
        eventName: 'CounterReset' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: encodeTopicValue(t) }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          sender: f.sender,
        }
        return encodeEventTopics({
          abi: DemoAutomationContractABI,
          eventName: 'CounterReset' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<CounterResetDecoded> => contract.decodeCounterReset(rawOutput),
    }
  }

  /**
   * Decodes a log into CounterReset data, preserving all log metadata.
   */
  decodeCounterReset(log: EVMLog): DecodedLog<CounterResetDecoded> {
    const decoded = decodeEventLog({
      abi: DemoAutomationContractABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as [Hex, ...Hex[]],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as CounterResetDecoded }
  }

  /**
   * Creates a log trigger for UpkeepPerformed events.
   * The returned trigger's adapt method decodes the raw log into UpkeepPerformedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerUpkeepPerformed(
    filters?: UpkeepPerformedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: DemoAutomationContractABI,
        eventName: 'UpkeepPerformed' as const,
      })
      topics = encoded.map((t) => ({ values: encodeTopicValue(t) }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        sender: f.sender,
      }
      const encoded = encodeEventTopics({
        abi: DemoAutomationContractABI,
        eventName: 'UpkeepPerformed' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: encodeTopicValue(t) }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          sender: f.sender,
        }
        return encodeEventTopics({
          abi: DemoAutomationContractABI,
          eventName: 'UpkeepPerformed' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<UpkeepPerformedDecoded> => contract.decodeUpkeepPerformed(rawOutput),
    }
  }

  /**
   * Decodes a log into UpkeepPerformed data, preserving all log metadata.
   */
  decodeUpkeepPerformed(log: EVMLog): DecodedLog<UpkeepPerformedDecoded> {
    const decoded = decodeEventLog({
      abi: DemoAutomationContractABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as [Hex, ...Hex[]],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as UpkeepPerformedDecoded }
  }
}

