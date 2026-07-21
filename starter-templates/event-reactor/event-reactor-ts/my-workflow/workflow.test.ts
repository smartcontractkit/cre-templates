import { describe, expect } from 'bun:test'
import { TxStatus } from '@chainlink/cre-sdk'
import { EvmMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import type { Address } from 'viem'
import { newReactorConsumerMock } from '../contracts/evm/ts/generated/ReactorConsumer_mock'
import { initWorkflow, onLargeTransfer } from './workflow'

const CHAIN_SELECTOR = 16015286601757825753n // ethereum-testnet-sepolia
const MONITORED_ADDRESS = '0x805A04e5C8b2dcb2B26C9e9C9aa12ce34374A35b' as Address
const REACTOR_ADDRESS = '0x899526269d89284aD454ADb3C3BC1277c0B5E09a' as Address
const SENDER_ADDRESS = '0x3d89B571fe7D400f42Cb75ecffB76eea1eB96DB5' as Address

const makeConfig = () => ({
	chainSelectorName: 'ethereum-testnet-sepolia',
	monitoredContractAddress: MONITORED_ADDRESS,
	reactorContractAddress: REACTOR_ADDRESS,
})

const makeLargeTransferPayload = () => ({
	topics: [],
	data: {
		from: SENDER_ADDRESS as `0x${string}`,
		to: '0x000000000000000000000000000000000000dEaD' as `0x${string}`,
		amount: 2000000000000000000000n,
		timestamp: 1773692000n,
	},
})

describe('onLargeTransfer', () => {
	test('flags sender when not already flagged', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const reactorMock = newReactorConsumerMock(REACTOR_ADDRESS, evmMock)

		reactorMock.flaggedAddresses = () => false
		reactorMock.totalFlags = () => 0n

		evmMock.writeReport = () => ({
			txStatus: TxStatus.SUCCESS,
			txHash: new Uint8Array(32),
		})

		const runtime = newTestRuntime()
		;(runtime as any).config = makeConfig()

		const result = onLargeTransfer(runtime as any, makeLargeTransferPayload() as any)
		expect(result).toContain('Flagged')
	})

	test('skips when sender already flagged', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const reactorMock = newReactorConsumerMock(REACTOR_ADDRESS, evmMock)

		reactorMock.flaggedAddresses = () => true
		reactorMock.totalFlags = () => 3n

		const runtime = newTestRuntime()
		;(runtime as any).config = makeConfig()

		const result = onLargeTransfer(runtime as any, makeLargeTransferPayload() as any)
		expect(result).toContain('Skipped')
	})
})

describe('initWorkflow', () => {
	test('returns a handler subscribed to LargeTransfer log trigger', () => {
		const config = makeConfig()
		const handlers = initWorkflow(config)

		expect(handlers).toHaveLength(1)
		expect(handlers[0].fn).toBe(onLargeTransfer)

		const trigger = handlers[0].trigger as { adapt: Function; configAsAny: Function }
		expect(typeof trigger.adapt).toBe('function')
		expect(typeof trigger.configAsAny).toBe('function')
	})
})
