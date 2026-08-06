import {
	cre,
	hexToBase64,
	ok,
	text,
	type TeeRuntime,
} from '@chainlink/cre-sdk'
import { encodeAbiParameters, parseAbiParameters } from 'viem'
import { z } from 'zod'

// ─── Config Schema ──────────────────────────────────────────
export const configSchema = z.object({
	schedule: z.string(),
	url: z.string(),
	secretId: z.string(),
	scoreThreshold: z.number(),
})
type Config = z.infer<typeof configSchema>

// ─── Logic to be executed over confidential data ────────────
// Some logic needs to be computed over sensitive data while preserving the
// confidentiality of that data from node operators: risk thresholds, API
// credentials, centralised exchange stablecoin reserves for reasoning, identity
// details. Leaking this data could have adverse effects, including enabling
// front-running attacks, exposing sensitive financial information, and
// compromising individual privacy.
//
// Note what is and is not confidential here: a confidential workflow, despite
// running inside the enclave, is part of the binary the Workflow DON provides to
// the enclave — so the binary, including this logic, is revealed. What the
// enclave keeps confidential is the data this logic computes over: Vault DON
// secrets, the request and response payloads of HTTP calls made from the
// enclave, and other intermediate values.
//
// Keep it deterministic for a given input — the enclave result is attested and
// verified by DON consensus before the workflow completes.
const scoreResponse = (body: string): number => {
	let score = 0
	for (let i = 0; i < body.length; i++) {
		score = (score + body.charCodeAt(i)) % 1000
	}
	return score
}

// ─── TEE Cron Callback ──────────────────────────────────────
// Receives a `TeeRuntime`, not a `Runtime`. Everything here runs inside the
// enclave until we explicitly cross back with `usingTheDons()`.
export const onCronTrigger = (runtime: TeeRuntime<Config>): string => {
	const config = runtime.config

	// ── Step 2: Fetch a secret inside the enclave ──
	// The Vault DON releases this secret only into an attested enclave, and it is
	// decrypted at the moment `getSecret()` runs. There is nothing to declare
	// upfront (unlike Confidential HTTP's `vaultDonSecrets`).
	const apiToken = runtime.getSecret({ id: config.secretId }).result().value

	// ── Step 3: Make a capability call from inside the enclave ──
	// `HTTPClient.sendRequest()` has a `TeeRuntime` overload, so passing the TEE
	// runtime executes the request from inside the enclave, keeping the request
	// and response payloads confidential from node operators. The Workflow DON
	// offers consensus verification of enclave attestations, proving the integrity
	// of the logic executed within the enclave.
	//
	// Note: do NOT reach for `ConfidentialHTTPClient` here — it has no
	// `TeeRuntime` overload and is not meant to be called from a TEE handler.
	const response = new cre.capabilities.HTTPClient()
		.sendRequest(runtime, {
			url: config.url,
			method: 'GET',
			multiHeaders: {
				Authorization: { values: [`Bearer ${apiToken}`] },
			},
		})
		.result()

	if (!ok(response)) {
		throw new Error(`Confidential request failed with status: ${response.statusCode}`)
	}

	const body = text(response)

	// The default endpoint echoes the request headers back, so we can confirm the
	// secret really was injected inside the enclave — as a boolean, never by
	// logging the token itself. Drop this once `url` points at a real API.
	const secretReachedApi = body.includes(apiToken)

	// Decision logic executed over the confidential response payload.
	const score = scoreResponse(body)
	const verdict = score >= config.scoreThreshold ? 'APPROVE' : 'REJECT'

	// ⚠️ Logs should be used for simulations only, and MUST be removed before
	// deploying to production to preserve the confidentiality offered by enclaves.
	// Avoid logging inside the enclave in general — sensitive or not.
	runtime.log(`Enclave computation complete. verdict=${verdict}`)

	// ── Step 4: Cross back to the DON for anything that needs consensus ──
	// `usingTheDons()` returns a regular `Runtime`. Anything passed into a
	// capability call on it executes on Workflow DON nodes and is NO LONGER
	// confidential — so we cross over the verdict and score only, never the
	// secret or the raw response body.
	const donRuntime = runtime.usingTheDons()

	const encodedPayload = encodeAbiParameters(
		parseAbiParameters('string verdict, uint256 score'),
		[verdict, BigInt(score)],
	)

	donRuntime
		.report({
			encodedPayload: hexToBase64(encodedPayload),
			encoderName: 'evm',
			signingAlgo: 'ecdsa',
			hashingAlgo: 'keccak256',
		})
		.result()

	// The signed report is now a normal CRE report. To deliver it on-chain, pass
	// it to `evmClient.writeReport(donRuntime, report)` — see the Keeper Bot or
	// Event Reactor templates for the full write path.
	return `${verdict} (score: ${score}, secret reached API: ${secretReachedApi})`
}

// ─── Workflow Init ──────────────────────────────────────────
export function initWorkflow(config: Config) {
	const cronTrigger = new cre.capabilities.CronCapability()

	return [
		// ── Step 1: Register a TEE handler ──
		// `cre.handlerInTee` instead of `cre.handler`. The third argument is a
		// `TeeConstraint` describing which enclaves this handler will accept.
		//
		// Alternatives:
		//   {}                        — any registered TEE, any region
		//   { regions: ['us-west-2'] } — any TEE, restricted to a region
		//
		// AWS Nitro in us-west-2 is currently the only registered TEE type and
		// region; check your SDK version if you expect otherwise.
		cre.handlerInTee(cronTrigger.trigger({ schedule: config.schedule }), onCronTrigger, [
			{ tee: 'nitro', regions: ['us-west-2'] },
		]),
	]
}
