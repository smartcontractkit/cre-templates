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
	MonitoredToken,
	type DecodedLog,
	type LargeTransferDecoded,
} from '../contracts/evm/ts/generated/MonitoredToken'
import { ReactorConsumer } from '../contracts/evm/ts/generated/ReactorConsumer'

// ─── Config Schema ──────────────────────────────────────────
export const configSchema = z.object({
	chainSelectorName: z.string(),
	monitoredContractAddress: z.string(),
	reactorContractAddress: z.string(),
})
type Config = z.infer<typeof configSchema>

const ACTION_FLAG = 1

// ─── Log Trigger Callback ───────────────────────────────────
export const onLargeTransfer = (
	runtime: Runtime<Config>,
	payload: DecodedLog<LargeTransferDecoded>,
): string => {
	const config = runtime.config

	// 1. Read decoded event data (typed via generated bindings)
	const { from, to, amount, timestamp } = payload.data

	runtime.log(`LargeTransfer detected: ${from} -> ${to}, amount: ${amount}, timestamp: ${timestamp}`)

	// 2. Read on-chain state from reactor contract
	const network = getNetwork({
		chainFamily: 'evm',
		chainSelectorName: config.chainSelectorName,
		isTestnet: true,
	})
	if (!network) throw new Error(`Network not found: ${config.chainSelectorName}`)

	const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
	const reactor = new ReactorConsumer(evmClient, config.reactorContractAddress as Address)

	const alreadyFlagged = reactor.flaggedAddresses(runtime, from)
	const totalFlags = reactor.totalFlags(runtime)

	runtime.log(`Reactor state: totalFlags=${totalFlags}, sender already flagged=${alreadyFlagged}`)

	// 3. Decision: skip if already flagged
	if (alreadyFlagged) {
		runtime.log('Sender already flagged. No action needed.')
		return 'Skipped — already flagged'
	}

	// 4. Write: Flag the sender on-chain
	const reportData = encodeAbiParameters(
		parseAbiParameters('uint8 actionType, address target, string reason'),
		[ACTION_FLAG, from, `Large transfer of ${amount} detected`],
	)

	const writeResult = reactor.writeReport(runtime, reportData)

	if (writeResult.txStatus !== TxStatus.SUCCESS) {
		throw new Error(`Flag TX failed: ${writeResult.errorMessage || writeResult.txStatus}`)
	}

	const txHash = bytesToHex(writeResult.txHash || new Uint8Array(32))
	runtime.log(`Address flagged on-chain! TX: ${txHash}`)

	return `Flagged ${from} — tx: ${txHash}`
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
	const monitoredToken = new MonitoredToken(evmClient, config.monitoredContractAddress as Address)

	return [
		cre.handler(
			monitoredToken.logTriggerLargeTransfer(),
			onLargeTransfer,
		),
	]
}
