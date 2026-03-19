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
import { KeeperConsumer } from '../contracts/evm/ts/generated/KeeperConsumer'

// ─── Config Schema ──────────────────────────────────────────
export const configSchema = z.object({
	schedule: z.string(),
	evms: z.array(
		z.object({
			chainSelectorName: z.string(),
			contractAddress: z.string(),
		}),
	),
})
type Config = z.infer<typeof configSchema>

// ─── Callback ───────────────────────────────────────────────
export const onCronTrigger = (runtime: Runtime<Config>): string => {
	const evmConfig = runtime.config.evms[0]

	// 1. Get network and create EVM client
	const network = getNetwork({
		chainFamily: 'evm',
		chainSelectorName: evmConfig.chainSelectorName,
		isTestnet: true,
	})
	if (!network) throw new Error(`Network not found: ${evmConfig.chainSelectorName}`)

	const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
	const keeper = new KeeperConsumer(evmClient, evmConfig.contractAddress as Address)

	// 2. Read on-chain state
	const upkeepNeeded = keeper.needsUpkeep(runtime)
	const currentCounter = keeper.counter(runtime)
	const lastExecuted = keeper.lastExecuted(runtime)
	const interval = keeper.interval(runtime)

	runtime.log(`On-chain state: counter=${currentCounter}, lastExecuted=${lastExecuted}, interval=${interval}s`)
	runtime.log(`Upkeep needed: ${upkeepNeeded}`)

	// 3. Conditional: Only write if upkeep is needed
	if (!upkeepNeeded) {
		runtime.log('No upkeep needed. Skipping execution.')
		return 'Skipped — no upkeep needed'
	}

	// 4. Write: Send execution signal via signed report
	const reportData = encodeAbiParameters(parseAbiParameters('bool shouldExecute'), [true])

	const writeResult = keeper.writeReport(runtime, reportData)

	if (writeResult.txStatus !== TxStatus.SUCCESS) {
		throw new Error(`Keeper TX failed: ${writeResult.errorMessage || writeResult.txStatus}`)
	}

	const txHash = bytesToHex(writeResult.txHash || new Uint8Array(32))
	runtime.log(`Keeper executed! Counter will increment to ${currentCounter + 1n}. TX: ${txHash}`)

	return `Executed — tx: ${txHash}`
}

// ─── Workflow Init ──────────────────────────────────────────
export function initWorkflow(config: Config) {
	const cronTrigger = new cre.capabilities.CronCapability()

	return [
		cre.handler(
			cronTrigger.trigger({ schedule: config.schedule }),
			onCronTrigger,
		),
	]
}
