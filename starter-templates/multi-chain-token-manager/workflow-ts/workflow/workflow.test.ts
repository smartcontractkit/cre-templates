import { describe, expect } from 'bun:test'
import { cre, getNetwork } from '@chainlink/cre-sdk'
import { EvmMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import { getAddress, type Address } from 'viem'
import { MockPool } from '../contracts/evm/ts/generated/MockPool'
import { newMockPoolMock } from '../contracts/evm/ts/generated/MockPool_mock'
import { initWorkflow, onCronTrigger } from './workflow'

const CHAIN_SELECTOR = 16015286601757825753n // ethereum-testnet-sepolia
const POOL_ADDRESS = getAddress('0x6ae43d3271ff6888e7fc43fd7321a503ff738951')
const USDC_ASSET = getAddress('0x1f9840a85d5af5bf1d1762f925bdaddc4201f984')
const PROTOCOL_WALLET = getAddress('0x7109709ecfa91a80626ff3989d68f67f5b1dd12d')
const ATOKEN_ADDRESS = getAddress('0xbcf7c21f0b2f4114b5c764b064b97df773c93b44')
const DEBT_TOKEN_ADDRESS = getAddress('0xd5c3e3b477a458bfd721ca957d6978f4e71b9f23')
const VARIABLE_DEBT_ADDRESS = getAddress('0x3e0437898a5667a4769b1ca5a34aab1ae7e81377')
const INTEREST_RATE_STRATEGY = getAddress('0xa9f3c3cae095527061e6d270dbe163693e6fda9d')

describe('onCronTrigger', () => {
	test('reads lending pool reserve data and current liquidity rate', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const mockPoolMock = newMockPoolMock(POOL_ADDRESS, evmMock)
		mockPoolMock.getReserveData = () => ({
			configuration: { data: 0n },
			liquidityIndex: 0n,
			currentLiquidityRate: 1n * 10n ** 27n, // 1 RAY
			variableBorrowIndex: 0n,
			currentVariableBorrowRate: 0n,
			currentStableBorrowRate: 0n,
			lastUpdateTimestamp: 0,
			id: 0,
			aTokenAddress: ATOKEN_ADDRESS,
			stableDebtTokenAddress: DEBT_TOKEN_ADDRESS,
			variableDebtTokenAddress: VARIABLE_DEBT_ADDRESS,
			interestRateStrategyAddress: INTEREST_RATE_STRATEGY,
			accruedToTreasury: 0n,
			unbacked: 0n,
			isolationModeTotalDebt: 0n,
		})

		const runtime = newTestRuntime()
		const network = getNetwork({
			chainFamily: 'evm',
			chainSelectorName: 'ethereum-testnet-sepolia',
			isTestnet: true,
		})
		expect(network).toBeDefined()
		if (!network) return

		const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
		const pool = new MockPool(evmClient, POOL_ADDRESS)
		const reserveData = pool.getReserveData(runtime, USDC_ASSET)
		expect(reserveData.currentLiquidityRate).toBe(10n ** 27n)
	})

	test('reads token balance in lending pool for protocol wallet', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const mockPoolMock = newMockPoolMock(POOL_ADDRESS, evmMock)
		mockPoolMock.balanceOf = () => 5000000000000000000n

		const runtime = newTestRuntime()
		const network = getNetwork({
			chainFamily: 'evm',
			chainSelectorName: 'ethereum-testnet-sepolia',
			isTestnet: true,
		})
		expect(network).toBeDefined()
		if (!network) return

		const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
		const pool = new MockPool(evmClient, POOL_ADDRESS)
		const balance = pool.balanceOf(runtime, PROTOCOL_WALLET, USDC_ASSET)
		expect(balance).toBe(5000000000000000000n)
	})

	test('throws when scheduledExecutionTime is missing', () => {
		const runtime = newTestRuntime()
		expect(() => onCronTrigger(runtime as any, {} as any)).toThrow(
			'Scheduled execution time is required',
		)
	})

	test('throws when fewer than 2 EVM chains are configured', () => {
		const runtime = newTestRuntime()
		;(runtime as any).config = {
			schedule: '0 */10 * * * *',
			minBPSDeltaForRebalance: 50,
			evms: [
				{
					assetAddress: USDC_ASSET,
					poolAddress: POOL_ADDRESS,
					protocolSmartWalletAddress: PROTOCOL_WALLET,
					chainName: 'ethereum-testnet-sepolia',
					gasLimit: '500000',
				},
			],
		}
		expect(() =>
			onCronTrigger(runtime as any, { scheduledExecutionTime: Date.now() } as any),
		).toThrow('At least two EVM configurations are required')
	})
})

describe('initWorkflow', () => {
	test('subscribes onCronTrigger to the configured cron schedule', () => {
		const config = {
			schedule: '0 */10 * * * *',
			minBPSDeltaForRebalance: 50,
			evms: [
				{
					assetAddress: USDC_ASSET,
					poolAddress: POOL_ADDRESS,
					protocolSmartWalletAddress: PROTOCOL_WALLET,
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
