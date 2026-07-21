import { describe, expect } from 'bun:test'
import { cre, getNetwork } from '@chainlink/cre-sdk'
import { EvmMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import type { Address } from 'viem'
import { PriceFeedAggregator } from '../contracts/evm/ts/generated/PriceFeedAggregator'
import { newPriceFeedAggregatorMock } from '../contracts/evm/ts/generated/PriceFeedAggregator_mock'
import { initWorkflow, onCron } from './workflow'

const CHAIN_SELECTOR = 4949039107694359620n // ethereum-mainnet-arbitrum-1
const ETH_USD_FEED = '0x639Fe6ab55C921f74e7fac1ee960C0B6293ba612' as Address

describe('onCron', () => {
	test('reads price feed decimals and latest answer for each configured feed', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const mock = newPriceFeedAggregatorMock(ETH_USD_FEED, evmMock)
		mock.decimals = () => 8
		mock.latestAnswer = () => 500000000000n // 5000 * 1e8

		const runtime = newTestRuntime()
		const network = getNetwork({
			chainFamily: 'evm',
			chainSelectorName: 'ethereum-mainnet-arbitrum-1',
			isTestnet: false,
		})
		expect(network).toBeDefined()
		if (!network) return

		const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
		const aggregator = new PriceFeedAggregator(evmClient, ETH_USD_FEED)
		expect(aggregator.decimals(runtime)).toBe(8)
		expect(aggregator.latestAnswer(runtime)).toBe(500000000000n)
	})
})

describe('initWorkflow', () => {
	test('subscribes onCron to the configured cron schedule', () => {
		const config = {
			schedule: '0 */10 * * * *',
			chainName: 'ethereum-mainnet-arbitrum-1',
			feeds: [{ name: 'ETH/USD', address: ETH_USD_FEED }],
		}
		const handlers = initWorkflow(config)

		expect(handlers).toHaveLength(1)
		expect(handlers[0].fn).toBe(onCron)
		const cronTrigger = handlers[0].trigger as { config?: { schedule?: string } }
		expect(cronTrigger.config?.schedule).toBe(config.schedule)
	})
})
