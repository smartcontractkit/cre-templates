import { describe, expect } from 'bun:test'
import { TxStatus } from '@chainlink/cre-sdk'
import { ConfidentialHttpMock, EvmMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import { type Address, encodeAbiParameters, parseAbiParameters } from 'viem'
import { newPredictionMarketMock } from '../contracts/evm/ts/generated/PredictionMarket_mock'
import { initWorkflow, onCronTrigger } from './workflow'

const CHAIN_SELECTOR = 16015286601757825753n // ethereum-testnet-sepolia
const PREDICTION_MARKET = '0xEb792aF46AB2c2f1389A774AB806423DB43aA425' as Address
const BTC_FEED_ID = '0x00027bbaff688c906a3e20a34fe951715d1018d262a5b66e38edd64e0fdd0d00'

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
    feedId: BTC_FEED_ID,
  },
  marketIdsToCheck: [0, 1, 2],
})

const makeSecrets = () => {
  const inner = new Map<string, string>()
  inner.set('DATA_STREAMS_API_KEY', 'test-api-key')
  inner.set('DATA_STREAMS_API_SECRET', 'test-api-secret')
  const outer = new Map<string, Map<string, string>>()
  outer.set('default', inner)
  return outer
}

// Build a valid v3 (Crypto Advanced) fullReport with the given price (18 decimals).
const encodeFullReport = (priceWith18Decimals: bigint): `0x${string}` => {
  const reportData = encodeAbiParameters(
    parseAbiParameters(
      'bytes32 feedId, uint32 validFromTimestamp, uint32 observationsTimestamp, uint192 nativeFee, uint192 linkFee, uint32 expiresAt, int192 price, int192 bid, int192 ask',
    ),
    [
      BTC_FEED_ID as `0x${string}`,
      1_700_000_000,
      1_700_000_000,
      0n,
      0n,
      1_800_000_000,
      priceWith18Decimals,
      priceWith18Decimals,
      priceWith18Decimals,
    ],
  )
  return encodeAbiParameters(
    parseAbiParameters(
      'bytes32[3] reportContext, bytes reportData, bytes32[] rawRs, bytes32[] rawSs, bytes32 rawVs',
    ),
    [
      [
        '0x0000000000000000000000000000000000000000000000000000000000000000',
        '0x0000000000000000000000000000000000000000000000000000000000000000',
        '0x0000000000000000000000000000000000000000000000000000000000000000',
      ],
      reportData,
      [],
      [],
      '0x0000000000000000000000000000000000000000000000000000000000000000',
    ],
  )
}

const mockDataStreamsResponse = (priceWith18Decimals: bigint) => {
  const fullReport = encodeFullReport(priceWith18Decimals)
  const body = JSON.stringify({
    report: {
      feedID: BTC_FEED_ID,
      validFromTimestamp: 1_700_000_000,
      observationsTimestamp: 1_700_000_000,
      fullReport,
    },
  })
  const httpMock = ConfidentialHttpMock.testInstance()
  httpMock.sendRequest = () => ({
    statusCode: 200,
    body: new TextEncoder().encode(body),
  })
}

describe('market-resolution', () => {
  test('resolves a market when resolvable', () => {
    const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
    const pmMock = newPredictionMarketMock(PREDICTION_MARKET, evmMock)

    pmMock.isResolvable = (marketId: bigint) => marketId === 0n
    pmMock.writeReport = () => ({
      txStatus: TxStatus.SUCCESS,
      txHash: new Uint8Array(32),
    })

    // BTC at $105,000 with 18 decimals
    mockDataStreamsResponse(105_000n * 10n ** 18n)

    const runtime = newTestRuntime(makeSecrets())
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
