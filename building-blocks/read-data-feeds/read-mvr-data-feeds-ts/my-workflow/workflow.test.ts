import { describe, expect } from 'bun:test'
import { cre, getNetwork } from '@chainlink/cre-sdk'
import { EvmMock, newTestRuntime, test } from '@chainlink/cre-sdk/test'
import { encodeAbiParameters, type Address } from 'viem'
import { BundleAggregatorProxy } from '../contracts/evm/ts/generated/BundleAggregatorProxy'
import { newBundleAggregatorProxyMock } from '../contracts/evm/ts/generated/BundleAggregatorProxy_mock'
import { initWorkflow, onCron } from './workflow'

const CHAIN_SELECTOR = 15971525489660198786n // ethereum-mainnet-base-1
const BUNDLE_PROXY = '0x7109709ECfa91a80626fF3989D68f67F5b1DD12D' as Address

const MOCK_BUNDLE = encodeAbiParameters(
	[
		{ type: 'uint256', name: 'lastModifiedDateTime' },
		{ type: 'string', name: 'securityId' },
		{ type: 'string', name: 'securityName' },
		{ type: 'uint256', name: 'ssa' },
		{ type: 'string', name: 'ssaDesc' },
	],
	[1709510400n, 'US9128284T69', 'US Treasury 5Y', 4875000n, 'Investment Grade'],
)

describe('onCron', () => {
	test('reads on-chain MVR bundle decimals and decodes latest bundle data', async () => {
		const evmMock = EvmMock.testInstance(CHAIN_SELECTOR)
		const mock = newBundleAggregatorProxyMock(BUNDLE_PROXY, evmMock)
		mock.bundleDecimals = () => [0, 0, 0, 6]
		mock.latestBundle = () => MOCK_BUNDLE

		const runtime = newTestRuntime()
		const network = getNetwork({
			chainFamily: 'evm',
			chainSelectorName: 'ethereum-mainnet-base-1',
			isTestnet: false,
		})
		expect(network).toBeDefined()
		if (!network) return

		const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
		const proxy = new BundleAggregatorProxy(evmClient, BUNDLE_PROXY)
		expect(proxy.bundleDecimals(runtime)).toEqual([0, 0, 0, 6])
		expect(proxy.latestBundle(runtime)).toBe(MOCK_BUNDLE)
	})
})

describe('initWorkflow', () => {
	test('subscribes onCron to the configured cron schedule', () => {
		const config = {
			schedule: '0 */10 * * * *',
			chainName: 'ethereum-mainnet-base-1',
			feeds: [{ name: 'S&P Global SSA EURC', address: BUNDLE_PROXY }],
		}
		const handlers = initWorkflow(config)

		expect(handlers).toHaveLength(1)
		expect(handlers[0].fn).toBe(onCron)
		const cronTrigger = handlers[0].trigger as { config?: { schedule?: string } }
		expect(cronTrigger.config?.schedule).toBe(config.schedule)
	})
})
