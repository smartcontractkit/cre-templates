import { describe, expect } from 'bun:test'
import { cre, getNetwork, TxStatus } from '@chainlink/cre-sdk'
import { EvmMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import type { Address } from 'viem'
import { newPredictionMarketMock } from '../contracts/evm/ts/generated/PredictionMarket_mock'
import { newPriceFeedAggregatorMock } from '../contracts/evm/ts/generated/PriceFeedAggregator_mock'
import { initWorkflow, onCronTrigger } from './workflow'

const CHAIN_SELECTOR = 16015286601757825753n // ethereum-testnet-sepolia
const PREDICTION_MARKET = '0xEb792aF46AB2c2f1389A774AB806423DB43aA425' as Address
const PRICE_FEED = '0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43' as Address

const makeConfig = () => ({
  schedule: '0 */10 * * * *',
  evms: [
    {
      chainSelectorName: 'ethereum-testnet-sepolia',
      predictionMarketAddress: PREDICTION_MARKET,
      priceFeedAddress: PRICE_FEED,
      gasLimit: '500000',
    },
  ],
  marketIdsToCheck: [0, 1, 2],
})

describe('market-resolution', () => {
  test('resolves a market when resolvable', () => {
    const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
    const pmMock = newPredictionMarketMock(PREDICTION_MARKET, evmMock)
    const pfMock = newPriceFeedAggregatorMock(PRICE_FEED, evmMock)

    pmMock.isResolvable = (marketId: bigint) => marketId === 0n
    pfMock.latestAnswer = () => 10500000000000n // $105,000
    pfMock.decimals = () => 8
    pmMock.writeReport = () => ({
      txStatus: TxStatus.SUCCESS,
      txHash: new Uint8Array(32),
    })

    const runtime = newTestRuntime()
    ;(runtime as any).config = makeConfig()

    const result = onCronTrigger(runtime as any, { scheduledExecutionTime: new Date().toISOString() } as any)
    expect(result).toBe('Resolved markets: 0')
  })

  test('skips markets that are not resolvable', () => {
    const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
    const pmMock = newPredictionMarketMock(PREDICTION_MARKET, evmMock)

    pmMock.isResolvable = () => false

    const runtime = newTestRuntime()
    ;(runtime as any).config = makeConfig()

    const result = onCronTrigger(runtime as any, { scheduledExecutionTime: new Date().toISOString() } as any)
    expect(result).toBe('No markets were resolvable')
  })
})

describe('initWorkflow', () => {
  test('returns a single cron handler', () => {
    const config = makeConfig()
    const handlers = initWorkflow(config)

    expect(handlers).toHaveLength(1)
    expect(handlers[0].fn).toBe(onCronTrigger)
    const cronTrigger = handlers[0].trigger as { config?: { schedule?: string } }
    expect(cronTrigger.config?.schedule).toBe(config.schedule)
  })
})
