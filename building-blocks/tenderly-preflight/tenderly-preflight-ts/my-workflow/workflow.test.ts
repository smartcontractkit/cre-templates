import { describe, expect } from 'bun:test'
import { test } from '@chainlink/cre-sdk/test'
import { configSchema, initWorkflow, onCron } from './workflow'

const config = configSchema.parse({
	schedule: '0 */5 * * * *',
	tenderly: {
		accountSlug: 'my-account',
		projectSlug: 'my-project',
		accessKeySecretId: 'tenderlyaccesskey',
		fork: {
			networkId: '11155111', // ethereum-testnet-sepolia
			at: 'latest',
		},
	},
})

describe('configSchema', () => {
	test('applies tenderly defaults', () => {
		expect(config.tenderly.explainReverts).toBe(true)
		expect(config.tenderly.deleteEnvironment).toBe(true)
	})

	test('rejects a misspelled tenderly key', () => {
		expect(() =>
			configSchema.parse({
				schedule: '0 */5 * * * *',
				tenderly: {
					accountSlug: 'my-account',
					projectSlug: 'my-project',
					accessKeySecretId: 'tenderlyaccesskey',
					deleteEnviroment: false, // typo: strict objects reject unknown keys
					fork: { networkId: '11155111', at: 'latest' },
				},
			}),
		).toThrow()
	})
})

describe('initWorkflow', () => {
	test('subscribes onCron to the configured cron schedule', () => {
		const handlers = initWorkflow(config)

		expect(handlers).toHaveLength(1)
		expect(handlers[0].fn).toBe(onCron)
		const cronTrigger = handlers[0].trigger as { config?: { schedule?: string } }
		expect(cronTrigger.config?.schedule).toBe(config.schedule)
	})
})
