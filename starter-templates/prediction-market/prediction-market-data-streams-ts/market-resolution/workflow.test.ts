import { describe, expect } from 'bun:test'
import { cre, getNetwork, TxStatus } from '@chainlink/cre-sdk'
import { EvmMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import type { Address } from 'viem'
import { newPredictionMarketMock } from '../contracts/evm/ts/generated/PredictionMarket_mock'
import { initWorkflow, onCronTrigger } from './workflow'

const CHAIN_SELECTOR = 16015286601757825753n // ethereum-testnet-sepolia
const PREDICTION_MARKET = '0xEb792aF46AB2c2f1389A774AB806423DB43aA425' as Address

const makeConfig = () => ({
  schedule: '0 */10 * * * *',
  evms: [
    {
      chainSelectorName: 'ethereum-testnet-sepolia',
      predictionMarketAddress: PREDICTION_MARKET,
      gasLimit: '500000',
    },
  ],
  dataStreams: {
    apiUrl: 'https://api.testnet-dataengine.chain.link',
    feedId: '0x00027bbaff688c906a3e20a34fe951715d1018d262a5b66e38edd64e0fdd0d00',
  },
  marketIdsToCheck: [0, 1, 2],
})

describe('market-resolution', () => {
  test('resolves a market when resolvable', () => {
    const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
    const pmMock = newPredictionMarketMock(PREDICTION_MARKET, evmMock)

    pmMock.isResolvable = (marketId: bigint) => marketId === 0n
    pmMock.writeReport = () => ({
      txStatus: TxStatus.SUCCESS,
      txHash: new Uint8Array(32),
    })

    const runtime = newTestRuntime()
    ;(runtime as any).config = makeConfig()

    // Mock runtime.getSecret to return test credentials
    ;(runtime as any).getSecret = (name: string) => {
      if (name === 'DATA_STREAMS_API_KEY') return 'test-api-key'
      if (name === 'DATA_STREAMS_API_SECRET') return 'test-api-secret'
      return ''
    }

    // Note: In a full test, the HTTP client mock would intercept the Data Streams
    // API call and return a mock fullReport. For unit tests, the HTTP capability
    // needs to be mocked at the CRE SDK level.
    // This test verifies the workflow structure and contract interactions.

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
