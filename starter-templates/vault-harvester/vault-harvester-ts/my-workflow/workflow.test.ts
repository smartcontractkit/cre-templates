import { describe, expect } from 'bun:test'
import { TxStatus } from '@chainlink/cre-sdk'
import { EvmMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import type { Address } from 'viem'
import { newVaultHarvesterMock } from '../contracts/evm/ts/generated/VaultHarvester_mock'
import { initWorkflow, onCronTrigger } from './workflow'

const CHAIN_SELECTOR = 16015286601757825753n // ethereum-testnet-sepolia
const CONTRACT_ADDRESS = '0x7109709ECfa91a80626fF3989D68f67F5b1DD12D' as Address

const makeConfig = () => ({
	schedule: '0 */5 * * * *',
	evms: [
		{
			chainSelectorName: 'ethereum-testnet-sepolia',
			contractAddress: CONTRACT_ADDRESS,
		},
	],
})

describe('onCronTrigger', () => {
	test('harvests when shouldHarvest returns true', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const vaultMock = newVaultHarvesterMock(CONTRACT_ADDRESS, evmMock)

		vaultMock.shouldHarvest = () => true
		vaultMock.pendingYield = () => 2000000000000000000n
		vaultMock.totalHarvested = () => 10000000000000000000n
		vaultMock.harvestCount = () => 5n
		vaultMock.lastHarvest = () => 1000n

		evmMock.writeReport = () => ({
			txStatus: TxStatus.SUCCESS,
			txHash: new Uint8Array(32),
		})

		const runtime = newTestRuntime()
		;(runtime as any).config = makeConfig()

		const result = onCronTrigger(runtime as any)
		expect(result).toContain('Harvested')
	})

	test('skips when shouldHarvest returns false', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const vaultMock = newVaultHarvesterMock(CONTRACT_ADDRESS, evmMock)

		vaultMock.shouldHarvest = () => false
		vaultMock.pendingYield = () => 500000000000000n
		vaultMock.totalHarvested = () => 10000000000000000000n
		vaultMock.harvestCount = () => 5n
		vaultMock.lastHarvest = () => 1000n

		const runtime = newTestRuntime()
		;(runtime as any).config = makeConfig()

		const result = onCronTrigger(runtime as any)
		expect(result).toContain('Skipped')
	})
})

describe('initWorkflow', () => {
	test('returns a handler subscribed to cron trigger', () => {
		const config = makeConfig()
		const handlers = initWorkflow(config)

		expect(handlers).toHaveLength(1)
		expect(handlers[0].fn).toBe(onCronTrigger)

		const cronTrigger = handlers[0].trigger as { config?: { schedule?: string } }
		expect(cronTrigger.config?.schedule).toBe(config.schedule)
	})
})
