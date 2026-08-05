import { describe, expect } from 'bun:test'
import type { TeeRuntime } from '@chainlink/cre-sdk'
import { test } from '@chainlink/cre-sdk/test'
import { initWorkflow, onCronTrigger } from './workflow'

const API_TOKEN = 'test-token'

const makeConfig = () => ({
	schedule: '0 */1 * * * *',
	url: 'https://postman-echo.com/headers',
	secretId: 'API_TOKEN',
	scoreThreshold: 500,
})

// The public test surface does not yet ship a TEE runtime factory
// (`newTestRuntime` returns a DON `Runtime`), so we stand up the small slice of
// `TeeRuntime` the handler actually uses: config, getSecret, callCapability
// (which HTTPClient.sendRequest goes through), log, and usingTheDons.
type FakeTeeRuntimeOptions = {
	statusCode?: number
	body?: string
}

const makeFakeTeeRuntime = ({ statusCode = 200, body = 'hello' }: FakeTeeRuntimeOptions = {}) => {
	const capturedHeaders: string[] = []
	const reports: unknown[] = []
	const logs: string[] = []

	const runtime = {
		config: makeConfig(),
		getSecret: (request: { id?: string }) => ({
			result: () => ({ id: request.id, value: API_TOKEN }),
		}),
		callCapability: ({ payload }: { payload: { multiHeaders?: Record<string, unknown> } }) => {
			const auth = payload.multiHeaders?.Authorization as { values?: string[] } | undefined
			capturedHeaders.push(...(auth?.values ?? []))
			return {
				result: () => ({
					statusCode,
					body: new TextEncoder().encode(body),
				}),
			}
		},
		log: (message: string) => logs.push(message),
		usingTheDons: () => ({
			report: (input: unknown) => {
				reports.push(input)
				return { result: () => ({}) }
			},
		}),
	}

	return { runtime: runtime as unknown as TeeRuntime<ReturnType<typeof makeConfig>>, capturedHeaders, reports, logs }
}

describe('onCronTrigger', () => {
	test('injects the enclave-fetched secret into the outbound request', () => {
		const { runtime, capturedHeaders } = makeFakeTeeRuntime()

		onCronTrigger(runtime)

		expect(capturedHeaders).toEqual([`Bearer ${API_TOKEN}`])
	})

	test('confirms the secret reached the API when the response echoes it back', () => {
		const { runtime } = makeFakeTeeRuntime({ body: `{"authorization":"Bearer ${API_TOKEN}"}` })

		expect(onCronTrigger(runtime)).toContain('secret reached API: true')
	})

	test('reports the secret did not reach the API when it is absent', () => {
		const { runtime } = makeFakeTeeRuntime({ body: '{"authorization":"Bearer other"}' })

		expect(onCronTrigger(runtime)).toContain('secret reached API: false')
	})

	test('crosses back to the DON to generate a report', () => {
		const { runtime, reports } = makeFakeTeeRuntime()

		onCronTrigger(runtime)

		expect(reports).toHaveLength(1)
		expect(reports[0]).toMatchObject({
			encoderName: 'evm',
			signingAlgo: 'ecdsa',
			hashingAlgo: 'keccak256',
		})
	})

	test('APPROVEs when the confidential score clears the threshold', () => {
		// 'zzzzzzzz' sums to 976, above the 500 threshold.
		const { runtime } = makeFakeTeeRuntime({ body: 'zzzzzzzz' })

		expect(onCronTrigger(runtime)).toContain('APPROVE')
	})

	test('REJECTs when the confidential score is below the threshold', () => {
		// 'a' sums to 97, below the 500 threshold.
		const { runtime } = makeFakeTeeRuntime({ body: 'a' })

		expect(onCronTrigger(runtime)).toContain('REJECT')
	})

	test('throws on a non-2xx response and never reaches the DON', () => {
		const { runtime, reports } = makeFakeTeeRuntime({ statusCode: 401 })

		expect(() => onCronTrigger(runtime)).toThrow('status: 401')
		expect(reports).toHaveLength(0)
	})

	test('does not log the secret or the raw response body', () => {
		const { runtime, logs } = makeFakeTeeRuntime({ body: 'sensitive-response' })

		onCronTrigger(runtime)

		for (const line of logs) {
			expect(line).not.toContain(API_TOKEN)
			expect(line).not.toContain('sensitive-response')
		}
	})
})

describe('initWorkflow', () => {
	test('registers the cron handler with a Nitro TEE constraint', () => {
		const handlers = initWorkflow(makeConfig())

		expect(handlers).toHaveLength(1)
		expect(handlers[0].fn).toBe(onCronTrigger)

		// handlerInTee attaches TEE requirements; cre.handler does not.
		expect(handlers[0].requirements).toBeDefined()
	})
})
