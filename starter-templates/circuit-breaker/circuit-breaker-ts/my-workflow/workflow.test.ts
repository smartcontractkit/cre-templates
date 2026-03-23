import { describe, expect } from 'bun:test'
import { TxStatus } from '@chainlink/cre-sdk'
import { EvmMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import type { Address } from 'viem'
import { newProtocolWithBreakerMock } from '../contracts/evm/ts/generated/ProtocolWithBreaker_mock'
import { initWorkflow, onPriceUpdated, onHealthCheck } from './workflow'

const CHAIN_SELECTOR = 16015286601757825753n // ethereum-testnet-sepolia
const CONTRACT_ADDRESS = '0xaCb13C9940cB61367b45eEd504E410D4B4d7A6e4' as Address

const makeConfig = () => ({
	schedule: '0 */10 * * * *',
	chainSelectorName: 'ethereum-testnet-sepolia',
	protocolContractAddress: CONTRACT_ADDRESS,
})

const makePricePayload = (newPrice: bigint, oldPrice: bigint) => ({
	topics: [],
	data: {
		newPrice,
		oldPrice,
		timestamp: 1773700000n,
	},
})

describe('onPriceUpdated', () => {
	test('trips circuit breaker when deviation exceeds threshold', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const protocolMock = newProtocolWithBreakerMock(CONTRACT_ADDRESS, evmMock)

		protocolMock.paused = () => false
		protocolMock.tripCount = () => 0n
		protocolMock.priceDeviationThresholdBps = () => 1000n // 10%

		evmMock.writeReport = () => ({
			txStatus: TxStatus.SUCCESS,
			txHash: new Uint8Array(32),
		})

		const runtime = newTestRuntime()
		;(runtime as any).config = makeConfig()

		// 50% drop: 1e18 -> 0.5e18
		const result = onPriceUpdated(
			runtime as any,
			makePricePayload(500000000000000000n, 1000000000000000000n) as any,
		)
		expect(result).toContain('TRIPPED')
	})

	test('does nothing when deviation is within threshold', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const protocolMock = newProtocolWithBreakerMock(CONTRACT_ADDRESS, evmMock)

		protocolMock.paused = () => false
		protocolMock.tripCount = () => 0n
		protocolMock.priceDeviationThresholdBps = () => 1000n // 10%

		const runtime = newTestRuntime()
		;(runtime as any).config = makeConfig()

		// 5% drop: 1e18 -> 0.95e18
		const result = onPriceUpdated(
			runtime as any,
			makePricePayload(950000000000000000n, 1000000000000000000n) as any,
		)
		expect(result).toContain('Normal')
	})

	test('skips when protocol already paused', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const protocolMock = newProtocolWithBreakerMock(CONTRACT_ADDRESS, evmMock)

		protocolMock.paused = () => true
		protocolMock.tripCount = () => 1n

		const runtime = newTestRuntime()
		;(runtime as any).config = makeConfig()

		const result = onPriceUpdated(
			runtime as any,
			makePricePayload(500000000000000000n, 1000000000000000000n) as any,
		)
		expect(result).toContain('Already paused')
	})
})

describe('onHealthCheck', () => {
	test('reports healthy when protocol is operating normally', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const protocolMock = newProtocolWithBreakerMock(CONTRACT_ADDRESS, evmMock)

		protocolMock.paused = () => false
		protocolMock.lastPrice = () => 1000000000000000000n
		protocolMock.lastPriceTimestamp = () => 1773700000n
		protocolMock.tripCount = () => 0n

		const runtime = newTestRuntime()
		;(runtime as any).config = makeConfig()

		const result = onHealthCheck(runtime as any)
		expect(result).toContain('Healthy')
	})

	test('reports monitoring when protocol is paused', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const protocolMock = newProtocolWithBreakerMock(CONTRACT_ADDRESS, evmMock)

		protocolMock.paused = () => true
		protocolMock.lastPrice = () => 500000000000000000n
		protocolMock.lastPriceTimestamp = () => 1773700000n
		protocolMock.tripCount = () => 1n

		const runtime = newTestRuntime()
		;(runtime as any).config = makeConfig()

		const result = onHealthCheck(runtime as any)
		expect(result).toContain('Paused')
	})
})

describe('initWorkflow', () => {
	test('returns two handlers: log trigger + cron', () => {
		const config = makeConfig()
		const handlers = initWorkflow(config)

		expect(handlers).toHaveLength(2)
		expect(handlers[0].fn).toBe(onPriceUpdated)
		expect(handlers[1].fn).toBe(onHealthCheck)

		// Log trigger has adapt function
		const logTrigger = handlers[0].trigger as { adapt: Function; configAsAny: Function }
		expect(typeof logTrigger.adapt).toBe('function')

		// Cron trigger has schedule config
		const cronTrigger = handlers[1].trigger as { config?: { schedule?: string } }
		expect(cronTrigger.config?.schedule).toBe(config.schedule)
	})
})
