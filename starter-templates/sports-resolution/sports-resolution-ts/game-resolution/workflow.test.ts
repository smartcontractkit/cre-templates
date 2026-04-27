import { describe, expect } from 'bun:test'
import { newTestRuntime, test } from '@chainlink/cre-sdk/test'
import {
  applyAggregation,
  computeOutcome,
  initWorkflow,
  onSettlementRequested,
} from './workflow'

const SPORTS_MARKET = '0x000000000000000000000000000000000000dead'

const makeConfig = () => ({
  evms: [
    {
      chainSelectorName: 'ethereum-testnet-sepolia',
      sportsMarketAddress: SPORTS_MARKET,
      gasLimit: '500000',
    },
  ],
  dataSourceUrls: ['https://api.example1.com/v1', 'https://api.example2.com/v1'],
  aggregationMode: 'majority' as const,
})

// ─── Unit: computeOutcome ────────────────────────────────────

describe('computeOutcome', () => {
  test('home wins when homeScore > awayScore', () => {
    expect(computeOutcome({ homeScore: 3, awayScore: 1 })).toBe(1)
  })

  test('away wins when awayScore > homeScore', () => {
    expect(computeOutcome({ homeScore: 0, awayScore: 2 })).toBe(2)
  })

  test('draw when scores are equal', () => {
    expect(computeOutcome({ homeScore: 1, awayScore: 1 })).toBe(3)
  })
})

// ─── Unit: applyAggregation ──────────────────────────────────

describe('applyAggregation — majority', () => {
  test('2-of-2 agree: returns consensus', () => {
    expect(applyAggregation([1, 1], 'majority')).toBe(1)
  })

  test('2 sources disagree: returns null', () => {
    expect(applyAggregation([1, 2], 'majority')).toBeNull()
  })

  test('2-of-3 agree: returns consensus', () => {
    expect(applyAggregation([1, 1, 2], 'majority')).toBe(1)
  })

  test('3-way split (1 each): returns null', () => {
    expect(applyAggregation([1, 2, 3], 'majority')).toBeNull()
  })

  test('3-of-3 agree: returns consensus', () => {
    expect(applyAggregation([2, 2, 2], 'majority')).toBe(2)
  })
})

describe('applyAggregation — unanimous', () => {
  test('all agree: returns consensus', () => {
    expect(applyAggregation([1, 1, 1], 'unanimous')).toBe(1)
  })

  test('any disagreement: returns null', () => {
    expect(applyAggregation([1, 1, 2], 'unanimous')).toBeNull()
  })

  test('2-of-2 agree: returns consensus', () => {
    expect(applyAggregation([2, 2], 'unanimous')).toBe(2)
  })
})

// ─── Integration: initWorkflow ───────────────────────────────

describe('initWorkflow', () => {
  test('returns exactly one log trigger handler bound to onSettlementRequested', () => {
    const handlers = initWorkflow(makeConfig())

    expect(handlers).toHaveLength(1)
    expect(handlers[0].fn).toBe(onSettlementRequested)

    const trigger = handlers[0].trigger as {
      adapt: (raw: any) => any
      configAsAny: () => any
    }
    expect(typeof trigger.adapt).toBe('function')
    expect(typeof trigger.configAsAny).toBe('function')
  })
})
