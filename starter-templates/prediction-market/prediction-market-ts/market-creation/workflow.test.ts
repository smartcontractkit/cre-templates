import { describe, expect } from 'bun:test'
import { cre, getNetwork, TxStatus } from '@chainlink/cre-sdk'
import { EvmMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import type { Address } from 'viem'
import { PredictionMarket } from '../contracts/evm/ts/generated/PredictionMarket'
import { newPredictionMarketMock } from '../contracts/evm/ts/generated/PredictionMarket_mock'
import { initWorkflow, onCronTrigger } from './workflow'

const CHAIN_SELECTOR = 16015286601757825753n // ethereum-testnet-sepolia
const PREDICTION_MARKET = '0xEb792aF46AB2c2f1389A774AB806423DB43aA425' as Address
const PRICE_FEED = '0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43' as Address

const makeConfig = () => ({
  schedule: '0 0 * * * *',
  evms: [
    {
      chainSelectorName: 'ethereum-testnet-sepolia',
      predictionMarketAddress: PREDICTION_MARKET,
      priceFeedAddress: PRICE_FEED,
      gasLimit: '500000',
    },
  ],
  marketDefaults: {
    question: 'Will BTC be above $100,000 by {expirationDate}?',
    strikePriceUsd: 100000,
    durationSeconds: 86400,
  },
})

describe('market-creation', () => {
  test('creates a market successfully', () => {
    const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
    const pmMock = newPredictionMarketMock(PREDICTION_MARKET, evmMock)
    pmMock.getNextMarketId = () => 0n
    pmMock.writeReport = () => ({
      txStatus: TxStatus.SUCCESS,
      txHash: new Uint8Array(32),
    })

    const runtime = newTestRuntime()
    ;(runtime as any).config = makeConfig()

    const result = onCronTrigger(runtime as any, { scheduledExecutionTime: new Date().toISOString() } as any)
    const parsed = JSON.parse(result)

    expect(parsed.marketId).toBe('0')
    expect(parsed.question).toContain('Will BTC be above $100,000')
    expect(parsed.txHash).toBeDefined()
  })

  test('reads next market ID from contract via mock', () => {
    const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
    const pmMock = newPredictionMarketMock(PREDICTION_MARKET, evmMock)
    pmMock.getNextMarketId = () => 5n

    const runtime = newTestRuntime()
    const network = getNetwork({
      chainFamily: 'evm',
      chainSelectorName: 'ethereum-testnet-sepolia',
    })
    expect(network).toBeDefined()
    if (!network) return

    const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
    const pm = new PredictionMarket(evmClient, PREDICTION_MARKET)
    const nextId = pm.getNextMarketId(runtime)
    expect(nextId).toBe(5n)
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
