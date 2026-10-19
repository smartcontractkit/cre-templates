import { describe, expect } from 'bun:test'
import { EvmMock, HttpActionsMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import { decodeFunctionData } from 'viem'
import { StorageABI } from '../contracts/evm/ts/generated'
import { newStorageMock } from '../contracts/evm/ts/generated/Storage_mock'
import { type Config, initWorkflow, onCronTrigger } from './main'

const STORAGE_ADDRESS = '0x1234567890123456789012345678901234567890'
const SEPOLIA_CHAIN_SELECTOR = 16015286601757825753n // ethereum-testnet-sepolia
const BACKEND_URL = 'https://backend.example.com/reports'

const testConfig: Config = {
	schedule: '*/30 * * * * *',
	backendUrl: BACKEND_URL,
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

function setupMocks(pushStatus = 200) {
	const evmMock = EvmMock.testInstance(SEPOLIA_CHAIN_SELECTOR)
	const storageMock = newStorageMock(STORAGE_ADDRESS, evmMock)
	storageMock.writeReport = () => {
		return { txStatus: 'TX_STATUS_SUCCESS', txHash: '', errorMessage: '' }
	}

	const httpMock = HttpActionsMock.testInstance()
	let pushedBody: Record<string, unknown> | undefined
	let pushedUrl: string | undefined
	httpMock.sendRequest = (req) => {
		pushedUrl = req.url
		pushedBody = JSON.parse(new TextDecoder().decode(req.body))
		return { statusCode: pushStatus }
	}

	return { storageMock, getPushed: () => ({ url: pushedUrl, body: pushedBody }) }
}

describe('onCronTrigger', () => {
	test('writes onchain, pushes the report via sendReport, and logs it copy-pasteable', () => {
		const { getPushed } = setupMocks(200)
		const runtime = newTestRuntime<Config>(null, {}, testConfig)

		onCronTrigger(runtime)

		// the report was POSTed to the backend with the validator app's JSON shape
		const pushed = getPushed()
		expect(pushed.url).toBe(BACKEND_URL)
		expect(pushed.body).toBeDefined()
		expect(pushed.body?.rawReport).toMatch(/^0x[0-9a-fA-F]+$/)
		expect(pushed.body?.reportContext).toMatch(/^0x[0-9a-fA-F]*$/)
		expect(Array.isArray(pushed.body?.signatures)).toBe(true)

		expect(runtime.getLogs().some((l) => l.includes(`HTTP 200`))).toBe(true)

		// the same payload is logged between the markers
		const logged = extractReportJson(runtime.getLogs())
		expect(logged).toEqual(pushed.body as Record<string, unknown>)
	})

	test('decodes to the uint256 value the consumer would store', () => {
		const { storageMock } = setupMocks(200)

		let writtenValue: bigint | undefined
		storageMock.writeReport = (req) => {
			// The report payload is the encoded onReport(bytes, bytes) call;
			// rawReport = 109-byte metadata header + onReport calldata.
			const rawReport = (req.report as { rawReport: Uint8Array }).rawReport
			const calldata = `0x${Array.from(rawReport.slice(109))
				.map((b) => b.toString(16).padStart(2, '0'))
				.join('')}` as `0x${string}`
			const decoded = decodeFunctionData({ abi: StorageABI, data: calldata })
			writtenValue = BigInt(decoded.args[1] as `0x${string}`)
			return { txStatus: 'TX_STATUS_SUCCESS', txHash: '', errorMessage: '' }
		}

		const runtime = newTestRuntime<Config>(null, {}, testConfig)
		onCronTrigger(runtime)

		expect(writtenValue).toBeDefined()
		expect(writtenValue as bigint).toBeGreaterThan(0n)
	})

	test('still logs the report when the backend push fails', () => {
		setupMocks(500)
		const runtime = newTestRuntime<Config>(null, {}, testConfig)

		onCronTrigger(runtime)

		expect(runtime.getLogs().some((l) => l.includes('status 500'))).toBe(true)
		const logged = extractReportJson(runtime.getLogs())
		expect(logged.rawReport).toMatch(/^0x[0-9a-fA-F]+$/)
	})

	test('throws when the onchain write fails', () => {
		const evmMock = EvmMock.testInstance(SEPOLIA_CHAIN_SELECTOR)
		const storageMock = newStorageMock(STORAGE_ADDRESS, evmMock)
		storageMock.writeReport = () => {
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
