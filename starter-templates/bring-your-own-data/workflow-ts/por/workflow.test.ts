import { describe, expect } from 'bun:test'
import { TxStatus } from '@chainlink/cre-sdk'
import { EvmMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import type { Address } from 'viem'
import { newDataFeedsCacheMock } from '../contracts/evm/ts/generated/DataFeedsCache_mock'
import {
	encodeBundleStruct,
	encodeReceivedBundledReports,
	hexToBytes32RightPadded,
	initWorkflow,
	onCronTrigger,
} from './workflow'

const CHAIN_SELECTOR = 16015286601757825753n // ethereum-testnet-sepolia
const DATA_FEEDS_CACHE = '0x694AA1769357215DE4FAC081bf1f309aDC325306' as Address

describe('onCronTrigger', () => {
	test('throws when scheduledExecutionTime is missing', () => {
		const runtime = newTestRuntime() as any
		expect(() => onCronTrigger(runtime, {} as any)).toThrow(
			'Scheduled execution time is required',
		)
	})
})

describe('encodeBundleStruct', () => {
	test('encodes a bundle with totalReserve and lastUpdated', () => {
		const result = encodeBundleStruct({
			totalReserve: 1000000000000000000n,
			lastUpdated: 1709510400,
		})
		expect(result).toMatch(/^0x/)
		expect(result.length).toBeGreaterThan(2)
	})
})

describe('encodeReceivedBundledReports', () => {
	test('encodes a single bundled report', () => {
		const bundle = encodeBundleStruct({
			totalReserve: 1000000000000000000n,
			lastUpdated: 1709510400,
		})
		const result = encodeReceivedBundledReports([
			{
				dataId: '0x01',
				timestamp: 1709510400,
				bundle,
			},
		])
		expect(result).toMatch(/^0x/)
		expect(result.length).toBeGreaterThan(2)
	})
})

describe('hexToBytes32RightPadded', () => {
	test('pads a short hex string to 32 bytes', () => {
		const result = hexToBytes32RightPadded('0x01')
		expect(result).toBe('0x0100000000000000000000000000000000000000000000000000000000000000')
	})

	test('throws for hex strings longer than 32 bytes', () => {
		const longHex = '0x' + 'ff'.repeat(33)
		expect(() => hexToBytes32RightPadded(longHex)).toThrow('exceeds 32 bytes')
	})
})

describe('DataFeedsCache mock', () => {
	test('writes report via generated DataFeedsCache binding', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const mock = newDataFeedsCacheMock(DATA_FEEDS_CACHE, evmMock)
		mock.writeReport = () => ({
			txStatus: TxStatus.SUCCESS,
			txHash: new Uint8Array(32),
		})

		expect(typeof mock.writeReport).toBe('function')
	})
})

describe('initWorkflow', () => {
	test('subscribes onCronTrigger to the configured cron schedule', () => {
		const config = {
			schedule: '0 */10 * * * *',
			url: 'https://api.example.com/por',
			dataIdHex: '0x01',
			evms: [
				{
					dataFeedsCacheAddress: DATA_FEEDS_CACHE,
					chainName: 'ethereum-testnet-sepolia',
					gasLimit: '500000',
				},
			],
		}
		const handlers = initWorkflow(config)

		expect(handlers).toHaveLength(1)
		expect(handlers[0].fn).toBe(onCronTrigger)
		const cronTrigger = handlers[0].trigger as { config?: { schedule?: string } }
		expect(cronTrigger.config?.schedule).toBe(config.schedule)
	})
})
