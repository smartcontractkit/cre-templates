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
import { VaultHarvester } from '../contracts/evm/ts/generated/VaultHarvester'

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
	const vault = new VaultHarvester(evmClient, evmConfig.contractAddress as Address)

	// 2. Read on-chain state
	const harvestNeeded = vault.shouldHarvest(runtime)
	const pendingYield = vault.pendingYield(runtime)
	const totalHarvested = vault.totalHarvested(runtime)
	const harvestCount = vault.harvestCount(runtime)
	const lastHarvest = vault.lastHarvest(runtime)

	runtime.log(`On-chain state: pendingYield=${pendingYield}, totalHarvested=${totalHarvested}, harvestCount=${harvestCount}, lastHarvest=${lastHarvest}`)
	runtime.log(`Harvest needed: ${harvestNeeded}`)

	// 3. Conditional: Only harvest if profitable
	if (!harvestNeeded) {
		runtime.log('Harvest not profitable yet. Skipping.')
		return 'Skipped — not profitable'
	}

	// 4. Write: Trigger harvest via signed report
	const reportData = encodeAbiParameters(parseAbiParameters('bool shouldExecute'), [true])

	const writeResult = vault.writeReport(runtime, reportData)

	if (writeResult.txStatus !== TxStatus.SUCCESS) {
		throw new Error(`Harvest TX failed: ${writeResult.errorMessage || writeResult.txStatus}`)
	}

	const txHash = bytesToHex(writeResult.txHash || new Uint8Array(32))
	runtime.log(`Vault harvested! Yield: ${pendingYield}, TX: ${txHash}`)

	return `Harvested — tx: ${txHash}`
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
