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
import { address, getAddressEncoder, getProgramDerivedAddress } from '@solana/addresses'
import { getBase58Decoder } from '@solana/codecs'
import { z } from 'zod'
import { KvStoreReceiver, type KvEntry } from '../contracts/ts/generated'

const BASE58_DECODER = getBase58Decoder()
const ADDRESS_ENCODER = getAddressEncoder()

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
		// The keystone forwarder program id that owns forwarderState — see
		// https://docs.chain.link/cre/guides/workflow/using-solana-client/solana-forwarder-directory-ts
		forwarderProgramId: base58Address,
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

// PDA the forwarder program signs with when it CPIs into on_report — must match
// kv_store_receiver's own derivation (see contracts/solana/programs/kv_store_receiver/src/lib.rs).
const deriveForwarderAuthority = async (solanaConfig: Config['solana']) => {
	const [authority] = await getProgramDerivedAddress({
		programAddress: address(solanaConfig.forwarderProgramId),
		seeds: [
			'forwarder',
			ADDRESS_ENCODER.encode(address(solanaConfig.forwarderState)),
			ADDRESS_ENCODER.encode(address(solanaConfig.receiverProgramId)),
		],
	})
	return authority
}

const logSolanaConfig = (runtime: Runtime<Config>, solanaConfig: Config['solana']) => {
	runtime.log(
		`Solana config: chainSelectorName=${solanaConfig.chainSelectorName} ` +
			`receiverProgramId=${solanaConfig.receiverProgramId} ` +
			`forwarderState=${solanaConfig.forwarderState} ` +
			`forwarderProgramId=${solanaConfig.forwarderProgramId}`,
	)
}

const onCronTrigger = async (runtime: Runtime<Config>) => {
	const solanaConfig = runtime.config.solana
	logSolanaConfig(runtime, solanaConfig)

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

	const forwarderAuthority = await deriveForwarderAuthority(solanaConfig)

	// keystone-forwarder account layout: forwarder state first, the forwarder
	// authority PDA second, then the accounts kv_store_receiver's on_report
	// instruction needs (here: just the kv_store account). Order matters — the
	// full list is hashed into the report and verified on-chain.
	const accounts: SolanaAccountMeta[] = [
		solanaAccountMeta(solanaConfig.forwarderState, true),
		solanaAccountMeta(forwarderAuthority),
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
	logSolanaConfig(runtime, solanaConfig)

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

	const { exists, account } = kvStoreReceiver.readKvStore(runtime, kvStoreAccount)

	if (!exists) {
		runtime.log(`kv_store account ${kvStoreAccount} not found yet (has on_report ever run?)`)
		return { Exists: false, Found: false, Key: '', Value: '', UpdatedAt: '', UpdateCount: '' }
	}

	if (!account) {
		// The account exists on-chain (confirmed via lamports), but this CRE
		// engine build didn't return its byte contents this time — see
		// KvStoreReceiver.readKvStore's doc comment.
		runtime.log(`kv_store account ${kvStoreAccount} exists but its data was not returned by this read`)
		return { Exists: true, Found: false, Key: '', Value: '', UpdatedAt: '', UpdateCount: '' }
	}

	runtime.log(
		`Read kv_store: key="${account.key}" value="${account.value}" ` +
			`updatedAt=${account.updatedAt} updateCount=${account.updateCount}`,
	)

	return {
		Exists: true,
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
