import { describe, expect } from 'bun:test'
import { cre, getNetwork, TxStatus } from '@chainlink/cre-sdk'
import { EvmMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import type { Address } from 'viem'
import { PredictionMarket } from '../contracts/evm/ts/generated/PredictionMarket'
import { newPredictionMarketMock } from '../contracts/evm/ts/generated/PredictionMarket_mock'
import { newPriceFeedAggregatorMock } from '../contracts/evm/ts/generated/PriceFeedAggregator_mock'
import { initWorkflow, onDisputeRaised } from './workflow'

const CHAIN_SELECTOR = 16015286601757825753n // ethereum-testnet-sepolia
const PREDICTION_MARKET = '0xEb792aF46AB2c2f1389A774AB806423DB43aA425' as Address
const PRICE_FEED = '0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43' as Address
const DISPUTOR = '0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984' as Address

const makeConfig = () => ({
  evms: [
    {
      chainSelectorName: 'ethereum-testnet-sepolia',
      predictionMarketAddress: PREDICTION_MARKET,
      priceFeedAddress: PRICE_FEED,
      gasLimit: '500000',
    },
  ],
})

describe('market-dispute', () => {
  test('resolves dispute when market is in Disputed status', () => {
    const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
    const pmMock = newPredictionMarketMock(PREDICTION_MARKET, evmMock)
    const pfMock = newPriceFeedAggregatorMock(PRICE_FEED, evmMock)

    pmMock.getMarketStatus = () => 2 // Disputed
    pfMock.latestAnswer = () => 10500000000000n // $105,000
    pfMock.decimals = () => 8
    pmMock.writeReport = () => ({
      txStatus: TxStatus.SUCCESS,
      txHash: new Uint8Array(32),
    })

    const runtime = newTestRuntime()
    ;(runtime as any).config = makeConfig()

    const payload = {
      topics: [],
      data: {
        marketId: 0n,
        disputor: DISPUTOR as `0x${string}`,
        reason: 'Price was stale',
      },
    }

    const result = onDisputeRaised(runtime as any, payload as any)
    const parsed = JSON.parse(result)

    expect(parsed.marketId).toBe('0')
    expect(parsed.freshPrice).toBeDefined()
    expect(parsed.txHash).toBeDefined()
  })

  test('skips when market is not in Disputed status', () => {
    const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
    const pmMock = newPredictionMarketMock(PREDICTION_MARKET, evmMock)

    pmMock.getMarketStatus = () => 1 // Resolved, not Disputed

    const runtime = newTestRuntime()
    ;(runtime as any).config = makeConfig()

    const payload = {
      topics: [],
      data: {
        marketId: 0n,
        disputor: DISPUTOR as `0x${string}`,
        reason: 'Price was stale',
      },
    }

    const result = onDisputeRaised(runtime as any, payload as any)
    expect(result).toContain('Skipped')
    expect(result).toContain('not in Disputed status')
  })
})

describe('initWorkflow', () => {
  test('returns a single log trigger handler', () => {
    const config = makeConfig()
    const handlers = initWorkflow(config)

    expect(handlers).toHaveLength(1)
    expect(handlers[0].fn).toBe(onDisputeRaised)
    const logTrigger = handlers[0].trigger as {
      adapt: (raw: any) => any
      configAsAny: () => any
    }
    expect(typeof logTrigger.adapt).toBe('function')
    expect(typeof logTrigger.configAsAny).toBe('function')
  })
})
