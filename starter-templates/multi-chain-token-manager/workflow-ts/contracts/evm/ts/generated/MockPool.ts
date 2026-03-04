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
 * Filter params for CurrentLiquidityRateUpdated. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type CurrentLiquidityRateUpdatedTopics = {
  asset?: `0x${string}`
}

/**
 * Decoded CurrentLiquidityRateUpdated event data.
 */
export type CurrentLiquidityRateUpdatedDecoded = {
  asset: `0x${string}`
  currentLiquidityRate: bigint
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
 * Filter params for Supply. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type SupplyTopics = {
  reserve?: `0x${string}`
  onBehalfOf?: `0x${string}`
  referralCode?: number
}

/**
 * Decoded Supply event data.
 */
export type SupplyDecoded = {
  reserve: `0x${string}`
  user: `0x${string}`
  onBehalfOf: `0x${string}`
  amount: bigint
  referralCode: number
}


/**
 * Filter params for Withdraw. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type WithdrawTopics = {
  reserve?: `0x${string}`
  user?: `0x${string}`
  to?: `0x${string}`
}

/**
 * Decoded Withdraw event data.
 */
export type WithdrawDecoded = {
  reserve: `0x${string}`
  user: `0x${string}`
  to: `0x${string}`
  amount: bigint
}


export const MockPoolABI = [{"type":"function","name":"acceptOwnership","inputs":[],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"balanceOf","inputs":[{"name":"","type":"address","internalType":"address"},{"name":"","type":"address","internalType":"address"}],"outputs":[{"name":"","type":"uint256","internalType":"uint256"}],"stateMutability":"view"},{"type":"function","name":"currentLiquidityRate","inputs":[{"name":"","type":"address","internalType":"address"}],"outputs":[{"name":"","type":"uint128","internalType":"uint128"}],"stateMutability":"view"},{"type":"function","name":"getReserveData","inputs":[{"name":"asset","type":"address","internalType":"address"}],"outputs":[{"name":"","type":"tuple","internalType":"structDataTypes.ReserveDataLegacy","components":[{"name":"configuration","type":"tuple","internalType":"structDataTypes.ReserveConfigurationMap","components":[{"name":"data","type":"uint256","internalType":"uint256"}]},{"name":"liquidityIndex","type":"uint128","internalType":"uint128"},{"name":"currentLiquidityRate","type":"uint128","internalType":"uint128"},{"name":"variableBorrowIndex","type":"uint128","internalType":"uint128"},{"name":"currentVariableBorrowRate","type":"uint128","internalType":"uint128"},{"name":"currentStableBorrowRate","type":"uint128","internalType":"uint128"},{"name":"lastUpdateTimestamp","type":"uint40","internalType":"uint40"},{"name":"id","type":"uint16","internalType":"uint16"},{"name":"aTokenAddress","type":"address","internalType":"address"},{"name":"stableDebtTokenAddress","type":"address","internalType":"address"},{"name":"variableDebtTokenAddress","type":"address","internalType":"address"},{"name":"interestRateStrategyAddress","type":"address","internalType":"address"},{"name":"accruedToTreasury","type":"uint128","internalType":"uint128"},{"name":"unbacked","type":"uint128","internalType":"uint128"},{"name":"isolationModeTotalDebt","type":"uint128","internalType":"uint128"}]}],"stateMutability":"view"},{"type":"function","name":"owner","inputs":[],"outputs":[{"name":"","type":"address","internalType":"address"}],"stateMutability":"view"},{"type":"function","name":"setCurrentLiquidityRate","inputs":[{"name":"asset","type":"address","internalType":"address"},{"name":"rate","type":"uint128","internalType":"uint128"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"supply","inputs":[{"name":"asset","type":"address","internalType":"address"},{"name":"amount","type":"uint256","internalType":"uint256"},{"name":"onBehalfOf","type":"address","internalType":"address"},{"name":"","type":"uint16","internalType":"uint16"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"transferOwnership","inputs":[{"name":"to","type":"address","internalType":"address"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"withdraw","inputs":[{"name":"asset","type":"address","internalType":"address"},{"name":"amount","type":"uint256","internalType":"uint256"},{"name":"to","type":"address","internalType":"address"}],"outputs":[{"name":"","type":"uint256","internalType":"uint256"}],"stateMutability":"nonpayable"},{"type":"event","name":"CurrentLiquidityRateUpdated","inputs":[{"name":"asset","type":"address","indexed":true,"internalType":"address"},{"name":"currentLiquidityRate","type":"uint256","indexed":false,"internalType":"uint256"}],"anonymous":false},{"type":"event","name":"OwnershipTransferRequested","inputs":[{"name":"from","type":"address","indexed":true,"internalType":"address"},{"name":"to","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},{"type":"event","name":"OwnershipTransferred","inputs":[{"name":"from","type":"address","indexed":true,"internalType":"address"},{"name":"to","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},{"type":"event","name":"Supply","inputs":[{"name":"reserve","type":"address","indexed":true,"internalType":"address"},{"name":"user","type":"address","indexed":false,"internalType":"address"},{"name":"onBehalfOf","type":"address","indexed":true,"internalType":"address"},{"name":"amount","type":"uint256","indexed":false,"internalType":"uint256"},{"name":"referralCode","type":"uint16","indexed":true,"internalType":"uint16"}],"anonymous":false},{"type":"event","name":"Withdraw","inputs":[{"name":"reserve","type":"address","indexed":true,"internalType":"address"},{"name":"user","type":"address","indexed":true,"internalType":"address"},{"name":"to","type":"address","indexed":true,"internalType":"address"},{"name":"amount","type":"uint256","indexed":false,"internalType":"uint256"}],"anonymous":false},{"type":"error","name":"AmountZero","inputs":[]},{"type":"error","name":"InsufficientBalance","inputs":[]},{"type":"error","name":"NoBalance","inputs":[]},{"type":"error","name":"SafeERC20FailedOperation","inputs":[{"name":"token","type":"address","internalType":"address"}]}] as const

export class MockPool {
  constructor(
    private readonly client: EVMClient,
    public readonly address: Address,
  ) {}

  balanceOf(
    runtime: Runtime<unknown>,
    arg0: `0x${string}`,
    arg1: `0x${string}`,
  ): bigint {
    const callData = encodeFunctionData({
      abi: MockPoolABI,
      functionName: 'balanceOf' as const,
      args: [arg0, arg1],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: MockPoolABI,
      functionName: 'balanceOf' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  currentLiquidityRate(
    runtime: Runtime<unknown>,
    arg0: `0x${string}`,
  ): bigint {
    const callData = encodeFunctionData({
      abi: MockPoolABI,
      functionName: 'currentLiquidityRate' as const,
      args: [arg0],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: MockPoolABI,
      functionName: 'currentLiquidityRate' as const,
      data: bytesToHex(result.data),
    }) as bigint
  }

  getReserveData(
    runtime: Runtime<unknown>,
    asset: `0x${string}`,
  ): { configuration: { data: bigint }; liquidityIndex: bigint; currentLiquidityRate: bigint; variableBorrowIndex: bigint; currentVariableBorrowRate: bigint; currentStableBorrowRate: bigint; lastUpdateTimestamp: number; id: number; aTokenAddress: `0x${string}`; stableDebtTokenAddress: `0x${string}`; variableDebtTokenAddress: `0x${string}`; interestRateStrategyAddress: `0x${string}`; accruedToTreasury: bigint; unbacked: bigint; isolationModeTotalDebt: bigint } {
    const callData = encodeFunctionData({
      abi: MockPoolABI,
      functionName: 'getReserveData' as const,
      args: [asset],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: MockPoolABI,
      functionName: 'getReserveData' as const,
      data: bytesToHex(result.data),
    }) as { configuration: { data: bigint }; liquidityIndex: bigint; currentLiquidityRate: bigint; variableBorrowIndex: bigint; currentVariableBorrowRate: bigint; currentStableBorrowRate: bigint; lastUpdateTimestamp: number; id: number; aTokenAddress: `0x${string}`; stableDebtTokenAddress: `0x${string}`; variableDebtTokenAddress: `0x${string}`; interestRateStrategyAddress: `0x${string}`; accruedToTreasury: bigint; unbacked: bigint; isolationModeTotalDebt: bigint }
  }

  owner(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: MockPoolABI,
      functionName: 'owner' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: MockPoolABI,
      functionName: 'owner' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  writeReportFromSetCurrentLiquidityRate(
    runtime: Runtime<unknown>,
    asset: `0x${string}`,
    rate: bigint,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: MockPoolABI,
      functionName: 'setCurrentLiquidityRate' as const,
      args: [asset, rate],
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

  writeReportFromSupply(
    runtime: Runtime<unknown>,
    asset: `0x${string}`,
    amount: bigint,
    onBehalfOf: `0x${string}`,
    arg3: number,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: MockPoolABI,
      functionName: 'supply' as const,
      args: [asset, amount, onBehalfOf, arg3],
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
      abi: MockPoolABI,
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

  writeReportFromWithdraw(
    runtime: Runtime<unknown>,
    asset: `0x${string}`,
    amount: bigint,
    to: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: MockPoolABI,
      functionName: 'withdraw' as const,
      args: [asset, amount, to],
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
   * Creates a log trigger for CurrentLiquidityRateUpdated events.
   * The returned trigger's adapt method decodes the raw log into CurrentLiquidityRateUpdatedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerCurrentLiquidityRateUpdated(
    filters?: CurrentLiquidityRateUpdatedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: MockPoolABI,
        eventName: 'CurrentLiquidityRateUpdated' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        asset: f.asset,
      }
      const encoded = encodeEventTopics({
        abi: MockPoolABI,
        eventName: 'CurrentLiquidityRateUpdated' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          asset: f.asset,
        }
        return encodeEventTopics({
          abi: MockPoolABI,
          eventName: 'CurrentLiquidityRateUpdated' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<CurrentLiquidityRateUpdatedDecoded> => contract.decodeCurrentLiquidityRateUpdated(rawOutput),
    }
  }

  /**
   * Decodes a log into CurrentLiquidityRateUpdated data, preserving all log metadata.
   */
  decodeCurrentLiquidityRateUpdated(log: EVMLog): DecodedLog<CurrentLiquidityRateUpdatedDecoded> {
    const decoded = decodeEventLog({
      abi: MockPoolABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as CurrentLiquidityRateUpdatedDecoded }
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
        abi: MockPoolABI,
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
        abi: MockPoolABI,
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
          abi: MockPoolABI,
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
      abi: MockPoolABI,
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
        abi: MockPoolABI,
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
        abi: MockPoolABI,
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
          abi: MockPoolABI,
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
      abi: MockPoolABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as OwnershipTransferredDecoded }
  }

  /**
   * Creates a log trigger for Supply events.
   * The returned trigger's adapt method decodes the raw log into SupplyDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerSupply(
    filters?: SupplyTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: MockPoolABI,
        eventName: 'Supply' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        reserve: f.reserve,
        onBehalfOf: f.onBehalfOf,
        referralCode: f.referralCode,
      }
      const encoded = encodeEventTopics({
        abi: MockPoolABI,
        eventName: 'Supply' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          reserve: f.reserve,
          onBehalfOf: f.onBehalfOf,
          referralCode: f.referralCode,
        }
        return encodeEventTopics({
          abi: MockPoolABI,
          eventName: 'Supply' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<SupplyDecoded> => contract.decodeSupply(rawOutput),
    }
  }

  /**
   * Decodes a log into Supply data, preserving all log metadata.
   */
  decodeSupply(log: EVMLog): DecodedLog<SupplyDecoded> {
    const decoded = decodeEventLog({
      abi: MockPoolABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as SupplyDecoded }
  }

  /**
   * Creates a log trigger for Withdraw events.
   * The returned trigger's adapt method decodes the raw log into WithdrawDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerWithdraw(
    filters?: WithdrawTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: MockPoolABI,
        eventName: 'Withdraw' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        reserve: f.reserve,
        user: f.user,
        to: f.to,
      }
      const encoded = encodeEventTopics({
        abi: MockPoolABI,
        eventName: 'Withdraw' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          reserve: f.reserve,
          user: f.user,
          to: f.to,
        }
        return encodeEventTopics({
          abi: MockPoolABI,
          eventName: 'Withdraw' as const,
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
      adapt: (rawOutput: EVMLog): DecodedLog<WithdrawDecoded> => contract.decodeWithdraw(rawOutput),
    }
  }

  /**
   * Decodes a log into Withdraw data, preserving all log metadata.
   */
  decodeWithdraw(log: EVMLog): DecodedLog<WithdrawDecoded> {
    const decoded = decodeEventLog({
      abi: MockPoolABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as WithdrawDecoded }
  }
}

