import { describe, expect } from 'bun:test'
import { EvmMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import { decodeFunctionData } from 'viem'
import { StorageABI } from '../contracts/evm/ts/generated'
import { newStorageMock } from '../contracts/evm/ts/generated/Storage_mock'
import { type Config, initWorkflow, onCronTrigger } from './main'

const STORAGE_ADDRESS = '0x1234567890123456789012345678901234567890'
const SEPOLIA_CHAIN_SELECTOR = 16015286601757825753n // ethereum-testnet-sepolia

const testConfig: Config = {
	schedule: '*/30 * * * * *',
	evms: [
		{
			storageAddress: STORAGE_ADDRESS,
			chainSelectorName: 'ethereum-testnet-sepolia',
			gasLimit: '500000',
		},
	],
}

function extractReportJson(logs: string[]): Record<string, unknown> {
	const begin = logs.indexOf('REPORT_JSON_BEGIN')
	const end = logs.indexOf('REPORT_JSON_END')
	expect(begin).toBeGreaterThanOrEqual(0)
	expect(end).toBeGreaterThan(begin)
	return JSON.parse(logs.slice(begin + 1, end).join('\n'))
}

describe('onCronTrigger', () => {
	test('writes the value via the chain writer and logs copy-pasteable report JSON', () => {
		const evmMock = EvmMock.testInstance(SEPOLIA_CHAIN_SELECTOR)
		const mock = newStorageMock(STORAGE_ADDRESS, evmMock)

		let writtenValue: bigint | undefined
		mock.writeReport = (req) => {
			// The report payload is the encoded onReport(bytes, bytes) call;
			// decode it and unwrap the uint256 value the consumer would store.
			const rawReport = (req.report as { rawReport: Uint8Array }).rawReport
			// rawReport = 109-byte metadata header + onReport calldata
			const calldata = `0x${Array.from(rawReport.slice(109))
				.map((b) => b.toString(16).padStart(2, '0'))
				.join('')}` as `0x${string}`
			const decoded = decodeFunctionData({ abi: StorageABI, data: calldata })
			const innerReport = decoded.args[1] as `0x${string}`
			writtenValue = BigInt(innerReport)
			return { txStatus: 'TX_STATUS_SUCCESS', txHash: '', errorMessage: '' }
		}

		const runtime = newTestRuntime<Config>(null, {}, testConfig)

		onCronTrigger(runtime)

		// the consumer receives abi.encode(uint256) — the timestamp value
		expect(writtenValue).toBeDefined()
		expect(writtenValue as bigint).toBeGreaterThan(0n)

		const reportJson = extractReportJson(runtime.getLogs())
		expect(reportJson.rawReport).toMatch(/^0x[0-9a-fA-F]+$/)
		expect(reportJson.reportContext).toMatch(/^0x[0-9a-fA-F]*$/)
		expect(Array.isArray(reportJson.signatures)).toBe(true)
		for (const sig of reportJson.signatures as string[]) {
			expect(sig).toMatch(/^0x[0-9a-fA-F]+$/)
		}
	})

	test('throws when the onchain write fails', () => {
		const evmMock = EvmMock.testInstance(SEPOLIA_CHAIN_SELECTOR)
		const mock = newStorageMock(STORAGE_ADDRESS, evmMock)

		mock.writeReport = () => {
			return { txStatus: 'TX_STATUS_REVERTED', txHash: '', errorMessage: 'execution reverted' }
		}

		const runtime = newTestRuntime<Config>(null, {}, testConfig)

		expect(() => onCronTrigger(runtime)).toThrow(/Failed to write report/)
	})
})

describe('initWorkflow', () => {
	test('returns one handler with the configured cron schedule', () => {
		const handlers = initWorkflow(testConfig)

		expect(handlers).toBeArray()
		expect(handlers).toHaveLength(1)
		expect(handlers[0].trigger.config.schedule).toBe(testConfig.schedule)
	})
})
