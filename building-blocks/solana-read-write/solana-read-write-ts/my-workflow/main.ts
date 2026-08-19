import {
	CronCapability,
	getNetwork,
	handler,
	Runner,
	type Runtime,
	type SolanaAccountMeta,
	SolanaClient,
	SolanaTxStatus,
	solanaAccountMeta,
} from '@chainlink/cre-sdk'
import { address } from '@solana/addresses'
import { getBase58Decoder } from '@solana/codecs'
import { z } from 'zod'
import { KvStoreReceiver, type KvEntry } from '../contracts/ts/generated'

const BASE58_DECODER = getBase58Decoder()

// Validates base58-encoded Solana addresses at config-parse time.
const base58Address = z.string().refine(
	(value) => {
		try {
			address(value)
			return true
		} catch {
			return false
		}
	},
	{ message: 'Invalid base58-encoded Solana address' },
)

const configSchema = z.object({
	// e.g. "0 */5 * * * *" (every 5 minutes, at second 0)
	schedule: z.string(),
	// Cron for the read-back check; defaults to `schedule` if omitted.
	readSchedule: z.string().optional(),
	solana: z.object({
		// e.g. "solana-devnet"
		chainSelectorName: z.string(),
		// Your kv_store_receiver program id — see contracts/solana/README.md
		receiverProgramId: base58Address,
		// Chainlink's shared keystone forwarder for this environment (do not deploy your own)
		forwarderState: base58Address,
		// PDA of ["forwarder", forwarderState, receiverProgramId] under the forwarder program
		forwarderAuthority: base58Address,
		// Accounts kv_store_receiver's on_report instruction needs, in order
		receiverAccounts: z.array(
			z.object({
				publicKey: base58Address,
				isWritable: z.boolean().optional(),
			}),
		),
	}),
})

type Config = z.infer<typeof configSchema>
type ConfiguredAccount = Config['solana']['receiverAccounts'][number]

const onCronTrigger = (runtime: Runtime<Config>) => {
	const solanaConfig = runtime.config.solana

	const network = getNetwork({
		chainFamily: 'solana',
		chainSelectorName: solanaConfig.chainSelectorName,
		isTestnet: true,
	})
	if (!network) {
		throw new Error(`Network not found for chain selector name: ${solanaConfig.chainSelectorName}`)
	}

	const kvStoreReceiver = new KvStoreReceiver(
		new SolanaClient(network.chainSelector.selector),
		solanaConfig.receiverProgramId,
	)

	// keystone-forwarder account layout: forwarder state first, the forwarder
	// authority PDA second, then the accounts kv_store_receiver's on_report
	// instruction needs (here: just the kv_store account). Order matters — the
	// full list is hashed into the report and verified on-chain.
	const accounts: SolanaAccountMeta[] = [
		solanaAccountMeta(solanaConfig.forwarderState, true),
		solanaAccountMeta(solanaConfig.forwarderAuthority),
		...solanaConfig.receiverAccounts.map((account: ConfiguredAccount) =>
			solanaAccountMeta(account.publicKey, account.isWritable ?? false),
		),
	]

	const entry: KvEntry = {
		key: 'cre-example',
		value: `hello from CRE @ ${runtime.now().toISOString()}`,
	}

	runtime.log(`Writing KvEntry {key: "${entry.key}", value: "${entry.value}"} to kv_store_receiver...`)

	const resp = kvStoreReceiver.writeReportFromKvEntry(runtime, entry, accounts)

	if (resp.txStatus !== SolanaTxStatus.SUCCESS) {
		throw new Error(`Failed to write report: ${resp.errorMessage || resp.txStatus}`)
	}

	const txSignature = resp.txSignature ? BASE58_DECODER.decode(resp.txSignature) : ''

	runtime.log(
		`Wrote to Solana devnet, tx=${txSignature} ` +
			`explorer=https://explorer.solana.com/tx/${txSignature}?cluster=devnet`,
	)

	return {
		Key: entry.key,
		Value: entry.value,
		TxSignature: txSignature,
		ExplorerUrl: `https://explorer.solana.com/tx/${txSignature}?cluster=devnet`,
	}
}

const onReadCronTrigger = (runtime: Runtime<Config>) => {
	const solanaConfig = runtime.config.solana

	const network = getNetwork({
		chainFamily: 'solana',
		chainSelectorName: solanaConfig.chainSelectorName,
		isTestnet: true,
	})
	if (!network) {
		throw new Error(`Network not found for chain selector name: ${solanaConfig.chainSelectorName}`)
	}

	const kvStoreReceiver = new KvStoreReceiver(
		new SolanaClient(network.chainSelector.selector),
		solanaConfig.receiverProgramId,
	)

	// The kv_store account is always receiverAccounts[0] for this program — the
	// same account the write cron passes to on_report.
	const kvStoreAccount = solanaConfig.receiverAccounts[0]?.publicKey
	if (!kvStoreAccount) {
		throw new Error('solana.receiverAccounts[0] (the kv_store account) is required to read')
	}

	const account = kvStoreReceiver.readKvStore(runtime, kvStoreAccount)

	if (!account) {
		runtime.log(`kv_store account ${kvStoreAccount} not found yet (has on_report ever run?)`)
		return { Found: false, Key: '', Value: '', UpdatedAt: '', UpdateCount: '' }
	}

	runtime.log(
		`Read kv_store: key="${account.key}" value="${account.value}" ` +
			`updatedAt=${account.updatedAt} updateCount=${account.updateCount}`,
	)

	return {
		Found: true,
		Key: account.key,
		Value: account.value,
		UpdatedAt: account.updatedAt.toString(),
		UpdateCount: account.updateCount.toString(),
	}
}

const initWorkflow = (config: Config) => {
	const cron = new CronCapability()
	return [
		handler(cron.trigger({ schedule: config.schedule }), onCronTrigger),
		handler(
			cron.trigger({ schedule: config.readSchedule ?? config.schedule }),
			onReadCronTrigger,
		),
	]
}

export async function main() {
	const runner = await Runner.newRunner<Config>({ configSchema })
	await runner.run(initWorkflow)
}

main()
