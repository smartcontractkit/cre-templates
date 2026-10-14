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
	transaction: {
		from: '0x1111111111111111111111111111111111111111',
		to: '0x2222222222222222222222222222222222222222',
		data: '0x',
	},
})

describe('configSchema', () => {
	test('applies tenderly defaults and validates the transaction', () => {
		expect(config.tenderly.explainReverts).toBe(true)
		expect(config.tenderly.deleteEnvironment).toBe(true)
		expect(config.transaction.value).toBe('0x0')
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
				transaction: {
					from: '0x1111111111111111111111111111111111111111',
					to: '0x2222222222222222222222222222222222222222',
					data: '0x',
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
