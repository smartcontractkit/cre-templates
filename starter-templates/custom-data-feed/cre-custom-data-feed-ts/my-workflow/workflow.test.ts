import { describe, expect } from 'bun:test'
import { cre, getNetwork, TxStatus } from '@chainlink/cre-sdk'
import { EvmMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import type { Address } from 'viem'
import { BalanceReader } from '../contracts/evm/ts/generated/BalanceReader'
import { newBalanceReaderMock } from '../contracts/evm/ts/generated/BalanceReader_mock'
import { MessageEmitter } from '../contracts/evm/ts/generated/MessageEmitter'
import { newMessageEmitterMock } from '../contracts/evm/ts/generated/MessageEmitter_mock'
import { ReserveManager } from '../contracts/evm/ts/generated/ReserveManager'
import { newReserveManagerMock } from '../contracts/evm/ts/generated/ReserveManager_mock'
import { initWorkflow, onCronTrigger, onLogTrigger } from './workflow'

const CHAIN_SELECTOR = 16015286601757825753n // ethereum-testnet-sepolia
const BALANCE_READER = '0x7109709ECfa91a80626fF3989D68f67F5b1DD12D' as Address
const MESSAGE_EMITTER = '0x694AA1769357215DE4FAC081bf1f309aDC325306' as Address
const RESERVE_MANAGER = '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48' as Address
const TOKEN_HOLDER = '0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984' as Address

describe('onCronTrigger', () => {
	test('reads native token balances on-chain via BalanceReader', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const balanceReaderMock = newBalanceReaderMock(BALANCE_READER, evmMock)
		balanceReaderMock.getNativeBalances = () => [1000000000000000000n]

		const runtime = newTestRuntime()
		const network = getNetwork({
			chainFamily: 'evm',
			chainSelectorName: 'ethereum-testnet-sepolia',
			isTestnet: true,
		})
		expect(network).toBeDefined()
		if (!network) return

		const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
		const reader = new BalanceReader(evmClient, BALANCE_READER)
		const balances = reader.getNativeBalances(runtime, [TOKEN_HOLDER])
		expect(balances).toEqual([1000000000000000000n])
	})

	test('writes updated reserve data on-chain via ReserveManager', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const reserveManagerMock = newReserveManagerMock(RESERVE_MANAGER, evmMock)
		reserveManagerMock.writeReport = () => ({
			txStatus: TxStatus.SUCCESS,
			txHash: new Uint8Array(32),
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
		const manager = new ReserveManager(evmClient, RESERVE_MANAGER)
		const resp = manager.writeReportFromUpdateReserves(
			runtime,
			{ totalMinted: 1000n, totalReserve: 2000n },
			{ gasLimit: '500000' },
		)
		expect(resp.txStatus).toBe(TxStatus.SUCCESS)
	})

	test('throws when scheduledExecutionTime is missing', () => {
		const runtime = newTestRuntime()
		expect(() => onCronTrigger(runtime as any, {} as any)).toThrow(
			'Scheduled execution time is required',
		)
	})
})

describe('onLogTrigger', () => {
	test('returns decoded message from adapted trigger payload', () => {
		const runtime = newTestRuntime() as any

		const payload = {
			topics: [],
			data: {
				emitter: TOKEN_HOLDER as `0x${string}`,
				timestamp: 1709510400n,
				message: 'hello-from-adapted-trigger',
			},
		}

		const result = onLogTrigger(runtime, payload as any)
		expect(result).toBe('hello-from-adapted-trigger')
	})

	test('reads last emitted message from MessageEmitter contract via mock', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const msgEmitterMock = newMessageEmitterMock(MESSAGE_EMITTER, evmMock)
		msgEmitterMock.getLastMessage = () => 'hello-from-mock'

		const runtime = newTestRuntime()
		const network = getNetwork({
			chainFamily: 'evm',
			chainSelectorName: 'ethereum-testnet-sepolia',
			isTestnet: true,
		})
		expect(network).toBeDefined()
		if (!network) return

		const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
		const emitter = new MessageEmitter(evmClient, MESSAGE_EMITTER)
		const msg = emitter.getLastMessage(runtime, TOKEN_HOLDER)
		expect(msg).toBe('hello-from-mock')
	})
})

describe('initWorkflow', () => {
	test('subscribes onCronTrigger to cron schedule and onLogTrigger to log filter', () => {
		const config = {
			schedule: '0 */10 * * * *',
			url: 'https://api.example.com/por',
			evms: [
				{
					tokenAddress: BALANCE_READER,
					porAddress: RESERVE_MANAGER,
					proxyAddress: RESERVE_MANAGER,
					balanceReaderAddress: BALANCE_READER,
					messageEmitterAddress: MESSAGE_EMITTER,
					chainSelectorName: 'ethereum-testnet-sepolia',
					gasLimit: '500000',
				},
			],
		}
		const handlers = initWorkflow(config)

		expect(handlers).toHaveLength(2)

		expect(handlers[0].fn).toBe(onCronTrigger)
		const cronTrigger = handlers[0].trigger as { config?: { schedule?: string } }
		expect(cronTrigger.config?.schedule).toBe(config.schedule)

		expect(handlers[1].fn).toBe(onLogTrigger)
		const logTrigger = handlers[1].trigger as {
			adapt: (raw: any) => any
			configAsAny: () => any
		}
		expect(typeof logTrigger.adapt).toBe('function')
		expect(typeof logTrigger.configAsAny).toBe('function')
	})
})
