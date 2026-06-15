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





export const IAutomationCompatibleABI = [{"inputs":[{"components":[{"internalType":"uint256","name":"index","type":"uint256"},{"internalType":"uint256","name":"timestamp","type":"uint256"},{"internalType":"bytes32","name":"txHash","type":"bytes32"},{"internalType":"uint256","name":"blockNumber","type":"uint256"},{"internalType":"bytes32","name":"blockHash","type":"bytes32"},{"internalType":"address","name":"source","type":"address"},{"internalType":"bytes32[]","name":"topics","type":"bytes32[]"},{"internalType":"bytes","name":"data","type":"bytes"}],"internalType":"structIAutomationCompatible.Log","name":"log","type":"tuple"},{"internalType":"bytes","name":"checkData","type":"bytes"}],"name":"checkLog","outputs":[{"internalType":"bool","name":"upkeepNeeded","type":"bool"},{"internalType":"bytes","name":"performData","type":"bytes"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes","name":"checkData","type":"bytes"}],"name":"checkUpkeep","outputs":[{"internalType":"bool","name":"upkeepNeeded","type":"bool"},{"internalType":"bytes","name":"performData","type":"bytes"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes","name":"performData","type":"bytes"}],"name":"performUpkeep","outputs":[],"stateMutability":"nonpayable","type":"function"}] as const

export class IAutomationCompatible {
  constructor(
    private readonly client: EVMClient,
    public readonly address: Address,
  ) {}

  checkLog(
    runtime: Runtime<unknown>,
    log: { index: bigint; timestamp: bigint; txHash: `0x${string}`; blockNumber: bigint; blockHash: `0x${string}`; source: `0x${string}`; topics: readonly `0x${string}`[]; data: `0x${string}` },
    checkData: `0x${string}`,
  ): readonly [boolean, `0x${string}`] {
    const callData = encodeFunctionData({
      abi: IAutomationCompatibleABI,
      functionName: 'checkLog' as const,
      args: [log, checkData],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: IAutomationCompatibleABI,
      functionName: 'checkLog' as const,
      data: bytesToHex(result.data),
    }) as readonly [boolean, `0x${string}`]
  }

  checkUpkeep(
    runtime: Runtime<unknown>,
    checkData: `0x${string}`,
  ): readonly [boolean, `0x${string}`] {
    const callData = encodeFunctionData({
      abi: IAutomationCompatibleABI,
      functionName: 'checkUpkeep' as const,
      args: [checkData],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: IAutomationCompatibleABI,
      functionName: 'checkUpkeep' as const,
      data: bytesToHex(result.data),
    }) as readonly [boolean, `0x${string}`]
  }

  writeReportFromPerformUpkeep(
    runtime: Runtime<unknown>,
    performData: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: IAutomationCompatibleABI,
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
}

