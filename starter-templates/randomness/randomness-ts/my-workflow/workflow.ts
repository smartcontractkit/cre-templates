import {
	bytesToHex,
	cre,
	getNetwork,
	TxStatus,
	type Runtime,
} from '@chainlink/cre-sdk'
import {
	type Address,
	encodeAbiParameters,
	parseAbiParameters,
} from 'viem'
import { z } from 'zod'
import {
	RandomnessConsumer,
	type DecodedLog,
	type RandomnessRequestedDecoded,
} from '../contracts/evm/ts/generated/RandomnessConsumer'

// ─── Config Schema ──────────────────────────────────────────
export const configSchema = z.object({
	chainSelectorName: z.string(),
	consumerAddress: z.string(),
})
type Config = z.infer<typeof configSchema>

// ─── Log Trigger Callback ───────────────────────────────────
// Fires on each RandomnessRequested event. Generates a consensus-safe random
// word and writes it back to the consumer as a signed report. This is the CRE
// equivalent of VRF's requestRandomWords -> fulfillRandomWords round trip.
export const onRandomnessRequested = (
	runtime: Runtime<Config>,
	payload: DecodedLog<RandomnessRequestedDecoded>,
): string => {
	const config = runtime.config

	// 1. Read decoded event data (typed via generated bindings)
	const { requestId, requester } = payload.data
	runtime.log(`RandomnessRequested: requestId=${requestId}, requester=${requester}`)

	const network = getNetwork({
		chainFamily: 'evm',
		chainSelectorName: config.chainSelectorName,
		isTestnet: true,
	})
	if (!network) throw new Error(`Network not found: ${config.chainSelectorName}`)

	const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
	const consumer = new RandomnessConsumer(evmClient, config.consumerAddress as Address)

	// 2. Idempotency: skip requests that are unknown or already fulfilled.
	//    (Reading finalized state prevents fulfilling the same requestId twice.)
	const isPending = consumer.pendingRequests(runtime, requestId)
	if (!isPending) {
		runtime.log(`Request ${requestId} is not pending (unknown or already fulfilled). Skipping.`)
		return `Skipped — request ${requestId} not pending`
	}

	// 3. Generate a consensus-safe random word. Inside the CRE runtime,
	//    Math.random() is DON-seeded, so every node produces the same value.
	//    (Not cryptographically secure — see README.)
	const randomWord = BigInt(Math.floor(Math.random() * Number.MAX_SAFE_INTEGER))
	runtime.log(`Generated randomWord=${randomWord} for request ${requestId}`)

	// 4. Encode the fulfillment payload exactly as RandomnessConsumer._processReport decodes it.
	const reportData = encodeAbiParameters(
		parseAbiParameters('uint256 requestId, uint256 randomWord'),
		[requestId, randomWord],
	)

	// 5. Write the signed report back to the consumer via the Keystone Forwarder.
	const writeResult = consumer.writeReport(runtime, reportData)
	if (writeResult.txStatus !== TxStatus.SUCCESS) {
		throw new Error(`Fulfillment TX failed: ${writeResult.errorMessage || writeResult.txStatus}`)
	}

	const txHash = bytesToHex(writeResult.txHash || new Uint8Array(32))
	runtime.log(`Fulfilled request ${requestId} on-chain! TX: ${txHash}`)

	return `Fulfilled ${requestId} — tx: ${txHash}`
}

// ─── Workflow Init ──────────────────────────────────────────
export function initWorkflow(config: Config) {
	const network = getNetwork({
		chainFamily: 'evm',
		chainSelectorName: config.chainSelectorName,
		isTestnet: true,
	})
	if (!network) throw new Error(`Network not found: ${config.chainSelectorName}`)

	const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
	const consumer = new RandomnessConsumer(evmClient, config.consumerAddress as Address)

	return [
		cre.handler(
			consumer.logTriggerRandomnessRequested(),
			onRandomnessRequested,
		),
	]
}
