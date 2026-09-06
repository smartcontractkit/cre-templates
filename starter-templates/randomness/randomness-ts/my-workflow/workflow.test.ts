import { describe, expect } from 'bun:test'
import { TxStatus } from '@chainlink/cre-sdk'
import { EvmMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import type { Address } from 'viem'
import { newRandomnessConsumerMock } from '../contracts/evm/ts/generated/RandomnessConsumer_mock'
import { initWorkflow, onRandomnessRequested } from './workflow'

const CHAIN_SELECTOR = 16015286601757825753n // ethereum-testnet-sepolia
const CONSUMER_ADDRESS = '0x899526269d89284aD454ADb3C3BC1277c0B5E09a' as Address
const REQUESTER = '0x3d89B571fe7D400f42Cb75ecffB76eea1eB96DB5' as Address

const makeConfig = () => ({
	chainSelectorName: 'ethereum-testnet-sepolia',
	consumerAddress: CONSUMER_ADDRESS,
})

const makeRequestPayload = (requestId: bigint) => ({
	topics: [],
	data: {
		requestId,
		requester: REQUESTER as `0x${string}`,
	},
})

describe('onRandomnessRequested', () => {
	test('fulfills a pending request', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const consumerMock = newRandomnessConsumerMock(CONSUMER_ADDRESS, evmMock)

		consumerMock.pendingRequests = () => true

		evmMock.writeReport = () => ({
			txStatus: TxStatus.SUCCESS,
			txHash: new Uint8Array(32),
		})

		const runtime = newTestRuntime()
		;(runtime as any).config = makeConfig()

		const result = onRandomnessRequested(runtime as any, makeRequestPayload(1n) as any)
		expect(result).toContain('Fulfilled')
	})

	test('skips a request that is not pending (already fulfilled or unknown)', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const consumerMock = newRandomnessConsumerMock(CONSUMER_ADDRESS, evmMock)

		consumerMock.pendingRequests = () => false

		const runtime = newTestRuntime()
		;(runtime as any).config = makeConfig()

		const result = onRandomnessRequested(runtime as any, makeRequestPayload(2n) as any)
		expect(result).toContain('Skipped')
	})
})

describe('initWorkflow', () => {
	test('returns a handler subscribed to RandomnessRequested log trigger', () => {
		const config = makeConfig()
		const handlers = initWorkflow(config)

		expect(handlers).toHaveLength(1)
		expect(handlers[0].fn).toBe(onRandomnessRequested)

		const trigger = handlers[0].trigger as { adapt: Function; configAsAny: Function }
		expect(typeof trigger.adapt).toBe('function')
		expect(typeof trigger.configAsAny).toBe('function')
	})
})
