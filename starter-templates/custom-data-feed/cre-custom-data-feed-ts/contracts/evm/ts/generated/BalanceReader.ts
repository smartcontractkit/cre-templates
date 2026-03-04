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





export const BalanceReaderABI = [{"inputs":[{"internalType":"address[]","name":"addresses","type":"address[]"}],"name":"getNativeBalances","outputs":[{"internalType":"uint256[]","name":"","type":"uint256[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"typeAndVersion","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}] as const

export class BalanceReader {
  constructor(
    private readonly client: EVMClient,
    public readonly address: Address,
  ) {}

  getNativeBalances(
    runtime: Runtime<unknown>,
    addresses: readonly `0x${string}`[],
  ): readonly bigint[] {
    const callData = encodeFunctionData({
      abi: BalanceReaderABI,
      functionName: 'getNativeBalances' as const,
      args: [addresses],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: BalanceReaderABI,
      functionName: 'getNativeBalances' as const,
      data: bytesToHex(result.data),
    }) as readonly bigint[]
  }

  typeAndVersion(
    runtime: Runtime<unknown>,
  ): string {
    const callData = encodeFunctionData({
      abi: BalanceReaderABI,
      functionName: 'typeAndVersion' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: BalanceReaderABI,
      functionName: 'typeAndVersion' as const,
      data: bytesToHex(result.data),
    }) as string
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

