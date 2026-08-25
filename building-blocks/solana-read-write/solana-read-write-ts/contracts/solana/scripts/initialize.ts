// Deploy helper — NOT part of the CRE workflow. Run this once after `anchor deploy`
// to create the `kv_store` account and record the keystone forwarder program id.
// It also derives `forwarderAuthority`, the PDA the forwarder signs with when it
// CPIs into this program's `on_report` instruction — printed only so you can sanity
// check it; the workflow (main.ts) derives it itself at runtime and it is not part
// of config.staging.json.
//
// Usage (from contracts/solana):
//   bun run scripts/initialize.ts \
//     --forwarder-program <FORWARDER_PROGRAM_ID> \
//     --forwarder-state <FORWARDER_STATE_PUBKEY> \
//     --keypair ~/.config/solana/id.json \
//     --url https://api.devnet.solana.com
//
// Prints receiverProgramId / forwarderState / forwarderProgramId / kvStoreAccount /
// a devnet explorer link — copy these into my-workflow/config.staging.json.
import { AnchorProvider, Program, Wallet, type Idl } from '@coral-xyz/anchor'
import { Connection, Keypair, PublicKey } from '@solana/web3.js'
import { readFileSync } from 'node:fs'

function arg(name: string, fallback?: string): string {
	const flag = `--${name}`
	const idx = process.argv.indexOf(flag)
	if (idx !== -1 && process.argv[idx + 1]) return process.argv[idx + 1]
	if (fallback !== undefined) return fallback
	throw new Error(`missing required flag ${flag}`)
}

function loadKeypair(path: string): Keypair {
	const expanded = path.replace(/^~/, process.env.HOME ?? '~')
	const secret = JSON.parse(readFileSync(expanded, 'utf8'))
	return Keypair.fromSecretKey(Uint8Array.from(secret))
}

async function main() {
	const forwarderProgram = new PublicKey(arg('forwarder-program'))
	const forwarderState = new PublicKey(arg('forwarder-state'))
	const rpcUrl = arg('url', 'https://api.devnet.solana.com')
	const keypairPath = arg('keypair', '~/.config/solana/id.json')

	const idl = JSON.parse(
		readFileSync(new URL('../target/idl/kv_store_receiver.json', import.meta.url), 'utf8'),
	) as Idl & { address: string }
	const receiverProgramId = new PublicKey(idl.address)

	const payer = loadKeypair(keypairPath)
	const connection = new Connection(rpcUrl, 'confirmed')
	const provider = new AnchorProvider(connection, new Wallet(payer), {
		commitment: 'confirmed',
	})
	const program = new Program(idl, provider)

	const kvStore = Keypair.generate()

	const [forwarderAuthority] = PublicKey.findProgramAddressSync(
		[Buffer.from('forwarder'), forwarderState.toBuffer(), receiverProgramId.toBuffer()],
		forwarderProgram,
	)

	console.log('Sending initialize()...')
	const sig = await program.methods
		.initialize(forwarderProgram)
		.accounts({
			kvStore: kvStore.publicKey,
			signer: payer.publicKey,
			systemProgram: new PublicKey('11111111111111111111111111111111'),
		})
		.signers([kvStore])
		.rpc()

	console.log('\nInitialized kv_store_receiver')
	console.log('--------------------------------')
	console.log(`receiverProgramId  : ${receiverProgramId.toBase58()}`)
	console.log(`forwarderState     (input, echo) : ${forwarderState.toBase58()}`)
	console.log(`forwarderProgramId (input, echo) : ${forwarderProgram.toBase58()}`)
	console.log(`forwarderAuthority (derived, sanity-check only, not part of config) : ${forwarderAuthority.toBase58()}`)
	console.log(`kvStoreAccount     (receiverAccounts[0]) : ${kvStore.publicKey.toBase58()}`)
	console.log(`\ninitialize() tx: https://explorer.solana.com/tx/${sig}?cluster=devnet`)
	console.log(
		'\nCopy receiverProgramId, forwarderState, forwarderProgramId, and kvStoreAccount into my-workflow/config.staging.json.',
	)
}

main().catch((err) => {
	console.error(err)
	process.exit(1)
})
