import { describe, expect } from 'bun:test'
import { EvmMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import type { Address } from 'viem'
import { newSportsMarketMock } from '../contracts/evm/ts/generated/SportsMarket_mock'
import { onSettlementRequested } from './workflow'

const CHAIN_SELECTOR = 16015286601757825753n // ethereum-testnet-sepolia
const SPORTS_MARKET = '0x000000000000000000000000000000000000dead' as Address

const makeRuntime = () => {
  const runtime = newTestRuntime()
  ;(runtime as any).config = {
    evms: [
      {
        chainSelectorName: 'ethereum-testnet-sepolia',
        sportsMarketAddress: SPORTS_MARKET,
        gasLimit: '500000',
      },
    ],
    dataSourceUrls: ['https://api.example1.com/v1', 'https://api.example2.com/v1'],
    aggregationMode: 'majority',
  }
  return runtime
}

describe('onSettlementRequested', () => {
  test('skips API resolution when the game is already settled on-chain', () => {
    const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
    const sportsMarketMock = newSportsMarketMock(SPORTS_MARKET, evmMock)
    sportsMarketMock.getGame = (gameId: bigint) => ({
      gameId,
      description: 'Pacers at Thunder',
      outcome: 1,
      settledAt: 1710000000n,
    })

    const result = JSON.parse(
      onSettlementRequested(makeRuntime() as any, {
        data: { gameId: 401766123n, description: 'Pacers at Thunder' },
      } as any),
    )

    expect(result).toEqual({
      gameId: '401766123',
      status: 'already_settled',
      outcome: 1,
      settledAt: '1710000000',
    })
  })
})
