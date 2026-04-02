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
	ProtocolWithBreaker,
	type DecodedLog,
	type PriceUpdatedDecoded,
} from '../contracts/evm/ts/generated/ProtocolWithBreaker'

// ─── Config Schema ──────────────────────────────────────────
export const configSchema = z.object({
	schedule: z.string(),
	chainSelectorName: z.string(),
	protocolContractAddress: z.string(),
})
type Config = z.infer<typeof configSchema>

const ACTION_PAUSE = 1

// ─── Shared: Trip the circuit breaker ───────────────────────
function tripCircuitBreaker(
	runtime: Runtime<Config>,
	protocol: ProtocolWithBreaker,
	reason: string,
): string {
	runtime.log(`ANOMALY DETECTED: ${reason}`)

	const reportData = encodeAbiParameters(
		parseAbiParameters('uint8 action, string reason'),
		[ACTION_PAUSE, reason],
	)

	const writeResult = protocol.writeReport(runtime, reportData)

	if (writeResult.txStatus !== TxStatus.SUCCESS) {
		throw new Error(`Pause TX failed: ${writeResult.errorMessage || writeResult.txStatus}`)
	}

	const txHash = bytesToHex(writeResult.txHash || new Uint8Array(32))
	runtime.log(`Circuit breaker TRIPPED! TX: ${txHash}`)

	return `TRIPPED: ${reason} — tx: ${txHash}`
}

// ─── Helper: Create protocol instance ───────────────────────
function getProtocol(config: Config): ProtocolWithBreaker {
	const network = getNetwork({
		chainFamily: 'evm',
		chainSelectorName: config.chainSelectorName,
		isTestnet: true,
	})
	if (!network) throw new Error(`Network not found: ${config.chainSelectorName}`)

	const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
	return new ProtocolWithBreaker(evmClient, config.protocolContractAddress as Address)
}

// ─── Trigger 1: LogTrigger (Event-Driven Price Monitoring) ──
export const onPriceUpdated = (
	runtime: Runtime<Config>,
	payload: DecodedLog<PriceUpdatedDecoded>,
): string => {
	const protocol = getProtocol(runtime.config)

	// 1. Read decoded event data
	const { newPrice, oldPrice, timestamp } = payload.data
	runtime.log(`PriceUpdated: ${oldPrice} -> ${newPrice} at ${timestamp}`)

	// 2. Read on-chain state
	const isPaused = protocol.paused(runtime)
	const tripCount = protocol.tripCount(runtime)

	runtime.log(`Protocol state: paused=${isPaused}, tripCount=${tripCount}`)

	if (isPaused) {
		runtime.log('Protocol already paused. No action needed.')
		return 'Already paused'
	}

	// 3. Calculate price deviation in basis points
	if (oldPrice === 0n) {
		runtime.log('Old price is zero, skipping deviation check')
		return 'Skipped — zero price'
	}

	const diff = newPrice > oldPrice ? newPrice - oldPrice : oldPrice - newPrice
	const deviationBps = (diff * 10000n) / oldPrice
	const thresholdBps = protocol.priceDeviationThresholdBps(runtime)

	runtime.log(`Price deviation: ${deviationBps} bps (threshold: ${thresholdBps} bps)`)

	// 4. Check if deviation exceeds threshold
	if (deviationBps <= thresholdBps) {
		runtime.log('Price within normal range. No action needed.')
		return `Normal — deviation ${deviationBps} bps`
	}

	// 5. ANOMALY — trip the circuit breaker
	return tripCircuitBreaker(
		runtime,
		protocol,
		`Price deviation ${deviationBps}bps exceeds ${thresholdBps}bps threshold (${oldPrice} -> ${newPrice})`,
	)
}

// ─── Trigger 2: Cron (Periodic Health Check) ────────────────
export const onHealthCheck = (runtime: Runtime<Config>): string => {
	const protocol = getProtocol(runtime.config)

	// Read on-chain state
	const isPaused = protocol.paused(runtime)
	const lastPrice = protocol.lastPrice(runtime)
	const lastPriceTimestamp = protocol.lastPriceTimestamp(runtime)
	const tripCount = protocol.tripCount(runtime)

	runtime.log(`Health check — paused=${isPaused}, lastPrice=${lastPrice}, lastUpdate=${lastPriceTimestamp}, tripCount=${tripCount}`)

	if (isPaused) {
		runtime.log('Protocol is paused. Monitoring only.')
		return 'Paused — monitoring'
	}

	runtime.log('Health check passed. Protocol operating normally.')
	return 'Healthy'
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
	const protocol = new ProtocolWithBreaker(evmClient, config.protocolContractAddress as Address)
	const cronTrigger = new cre.capabilities.CronCapability()

	return [
		// Trigger 0: React to PriceUpdated events
		cre.handler(
			protocol.logTriggerPriceUpdated(),
			onPriceUpdated,
		),
		// Trigger 1: Periodic health check
		cre.handler(
			cronTrigger.trigger({ schedule: config.schedule }),
			onHealthCheck,
		),
	]
}
