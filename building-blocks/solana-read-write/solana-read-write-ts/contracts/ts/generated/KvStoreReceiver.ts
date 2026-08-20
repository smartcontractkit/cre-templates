// Code generated — DO NOT EDIT.
//
// Regenerate after (re)deploying `contracts/solana` with:
//   cre generate-bindings solana --language typescript
//
// This file mirrors that generator's output shape (see the official
// `solana-onchain-write` example in `@chainlink/cre-sdk-examples`) so it can be
// swapped for the real generated file once you've deployed your own program and
// have its IDL under `contracts/solana/target/idl/kv_store_receiver.json`.
import {
	bytesToHex,
	calculateAccountsHash,
	encodeForwarderReport,
	hexToBase64,
	prepareSolanaReportRequest,
	type Runtime,
	type SolanaAccountMeta,
	type SolanaClient,
	type SolanaComputeConfig,
	solanaAccountMetasToJson,
	solanaAddressToBytes,
} from '@chainlink/cre-sdk'
import {
	getStructCodec,
	addCodecSizePrefix,
	getUtf8Codec,
	getU32Codec,
	getI64Codec,
	getU64Codec,
	getBytesCodec,
	fixCodecSize,
} from '@solana/codecs'

// Pre-deployed on Solana devnet for this example (see contracts/solana/README.md
// to deploy your own instead — swap this constant, or pass your program id as the
// second constructor argument below).
export const KV_STORE_RECEIVER_PROGRAM_ID = 'aPkfGwVg9Lj3yF2CSgxeRmMn7uU9tKWT6PSHTDf2ZX6'

export const KV_STORE_RECEIVER_IDL = {
	address: KV_STORE_RECEIVER_PROGRAM_ID,
	metadata: {
		name: 'kv_store_receiver',
		version: '0.1.0',
		spec: '0.1.0',
		description: 'Minimal CRE keystone receiver that stores a key/value pair on Solana',
	},
	instructions: [
		{
			name: 'initialize',
			discriminator: [175, 175, 109, 31, 13, 152, 155, 237],
			accounts: [
				{ name: 'kv_store', writable: true, signer: true },
				{ name: 'signer', writable: true, signer: true },
				{ name: 'system_program', address: '11111111111111111111111111111111' },
			],
			args: [{ name: 'forwarder_program', type: 'pubkey' }],
		},
		{
			name: 'on_report',
			discriminator: [214, 173, 18, 221, 173, 148, 151, 208],
			accounts: [
				{ name: 'state' },
				{ name: 'forwarder_authority', signer: true },
				{ name: 'kv_store', writable: true },
			],
			args: [
				{ name: '_metadata', type: 'bytes' },
				{ name: 'report', type: 'bytes' },
			],
		},
	],
	accounts: [{ name: 'KvStore', discriminator: [202, 10, 189, 60, 188, 50, 44, 70] }],
	events: [{ name: 'KvUpdated', discriminator: [231, 113, 74, 73, 163, 93, 4, 161] }],
	errors: [
		{ code: 6000, name: 'InvalidForwarderProgram', msg: 'forwarder_program must be a non-default pubkey' },
		{
			code: 6001,
			name: 'ForwarderProgramNotConfigured',
			msg: 'Call initialize with the keystone forwarder program id before on_report',
		},
		{
			code: 6002,
			name: 'MismatchedForwarderProgram',
			msg: 'Forwarder state account owner does not match the forwarder_program stored at initialize',
		},
		{
			code: 6003,
			name: 'InvalidForwarderAuthority',
			msg: 'forwarder_authority is not the PDA for this state, receiver, and forwarder_program',
		},
		{ code: 6004, name: 'InvalidReportPayload', msg: 'Report payload must be Borsh KvEntry (key: string, value: string)' },
	],
	types: [
		{
			name: 'KvStore',
			type: {
				kind: 'struct',
				fields: [
					{ name: 'forwarder_program', type: 'pubkey' },
					{ name: 'key', type: 'string' },
					{ name: 'value', type: 'string' },
					{ name: 'updated_at', type: 'i64' },
					{ name: 'update_count', type: 'u64' },
				],
			},
		},
		{
			name: 'KvUpdated',
			type: {
				kind: 'struct',
				fields: [
					{ name: 'key', type: 'string' },
					{ name: 'value', type: 'string' },
					{ name: 'updated_at', type: 'i64' },
					{ name: 'update_count', type: 'u64' },
				],
			},
		},
	],
} as const

const DISCRIMINATOR_SIZE = 8

const expectDiscriminator = (label: string, expected: Uint8Array, data: Uint8Array): Uint8Array => {
	if (data.length < DISCRIMINATOR_SIZE) {
		throw new Error(`${label}: data too short for discriminator (${data.length} bytes)`)
	}
	for (let i = 0; i < DISCRIMINATOR_SIZE; i++) {
		if (data[i] !== expected[i]) {
			throw new Error(`${label}: discriminator mismatch`)
		}
	}
	return data.subarray(DISCRIMINATOR_SIZE)
}

/** Report payload the CRE workflow writes (Borsh: string key + string value). */
export type KvEntry = {
	key: string
	value: string
}

export const kvEntryCodec = getStructCodec([
	['key', addCodecSizePrefix(getUtf8Codec(), getU32Codec())],
	['value', addCodecSizePrefix(getUtf8Codec(), getU32Codec())],
])

export type KvStoreAccount = {
	forwarderProgram: Uint8Array
	key: string
	value: string
	updatedAt: bigint
	updateCount: bigint
}

// Field order matches the on-chain `KvStore` struct: forwarder_program comes
// first (fixed 32 bytes), then the Borsh-encoded key/value/updated_at/update_count.
export const kvStoreAccountCodec = getStructCodec([
	['forwarderProgram', fixCodecSize(getBytesCodec(), 32)],
	['key', addCodecSizePrefix(getUtf8Codec(), getU32Codec())],
	['value', addCodecSizePrefix(getUtf8Codec(), getU32Codec())],
	['updatedAt', getI64Codec()],
	['updateCount', getU64Codec()],
])

export const ACCOUNT_KV_STORE_DISCRIMINATOR = new Uint8Array([202, 10, 189, 60, 188, 50, 44, 70])

/**
 * Decodes raw KvStore account data (with its 8-byte discriminator) into KvStoreAccount.
 */
export const decodeKvStoreAccount = (data: Uint8Array): KvStoreAccount =>
	kvStoreAccountCodec.decode(
		expectDiscriminator('account KvStore', ACCOUNT_KV_STORE_DISCRIMINATOR, data),
	) as KvStoreAccount

export type KvUpdated = {
	key: string
	value: string
	updatedAt: bigint
	updateCount: bigint
}

export const kvUpdatedCodec = getStructCodec([
	['key', addCodecSizePrefix(getUtf8Codec(), getU32Codec())],
	['value', addCodecSizePrefix(getUtf8Codec(), getU32Codec())],
	['updatedAt', getI64Codec()],
	['updateCount', getU64Codec()],
])

export const EVENT_KV_UPDATED_DISCRIMINATOR = new Uint8Array([231, 113, 74, 73, 163, 93, 4, 161])

/**
 * Decodes raw KvUpdated event data (with its 8-byte discriminator) into KvUpdated.
 */
export const decodeKvUpdatedEvent = (data: Uint8Array): KvUpdated =>
	kvUpdatedCodec.decode(expectDiscriminator('event KvUpdated', EVENT_KV_UPDATED_DISCRIMINATOR, data)) as KvUpdated

// Solana compute units for the write tx. `computeConfig` is required by the
// capability (an omitted/undefined value is rejected), so this is applied
// whenever the caller doesn't pass their own.
const DEFAULT_COMPUTE_CONFIG: SolanaComputeConfig = { computeLimit: 290_000 }

export class KvStoreReceiver {
	readonly programId: Uint8Array

	// The program ID defaults to the placeholder above — pass the real deployed
	// address (see contracts/solana/README.md) as the second constructor argument.
	constructor(
		private readonly client: SolanaClient,
		programId: string | Uint8Array = KV_STORE_RECEIVER_PROGRAM_ID,
	) {
		this.programId = typeof programId === 'string' ? solanaAddressToBytes(programId) : programId
	}

	/**
	 * Publishes a Borsh-encoded `KvEntry` report through the CRE signer to this
	 * program's `on_report` entrypoint via the keystone forwarder.
	 *
	 * remainingAccounts must follow the keystone-forwarder account layout:
	 *   - Index 0: forwarderState — the forwarder program's state account.
	 *   - Index 1: forwarderAuthority — PDA derived from seeds
	 *     ["forwarder", forwarderState, receiverProgram] under the forwarder program ID.
	 *   - Index 2: the `kv_store` account (writable).
	 *
	 * The full account list is hashed (via calculateAccountsHash) into the report.
	 * The on-chain forwarder strips indices 0 and 1 before CPI-ing into the
	 * receiver, so they must be present and correctly ordered.
	 */
	writeReportFromKvEntry(
		runtime: Runtime<unknown>,
		input: KvEntry,
		remainingAccounts: SolanaAccountMeta[],
		computeConfig?: SolanaComputeConfig,
	) {
		const payload = new Uint8Array(kvEntryCodec.encode(input))

		const report = runtime
			.report(
				prepareSolanaReportRequest(
					encodeForwarderReport({
						accountHash: calculateAccountsHash(remainingAccounts),
						payload,
					}),
				),
			)
			.result()

		return this.client
			.writeReport(runtime, {
				remainingAccounts: solanaAccountMetasToJson(remainingAccounts),
				receiver: bytesToHex(this.programId),
				computeConfig: computeConfig ?? DEFAULT_COMPUTE_CONFIG,
				report,
			})
			.result()
	}

	/**
	 * Reads the `kv_store` account directly from the chain (no forwarder, no
	 * report — this is a plain consensus'd RPC read, the counterpart to
	 * `writeReportFromKvEntry`).
	 *
	 * `exists` reflects `lamports > 0` (account has been created by `initialize`).
	 * `account` is the decoded key/value — populated only when the capability's
	 * reply actually includes account bytes. As of the current CRE engine build,
	 * `GetAccountInfoWithOpts` reliably returns `lamports`/`owner`/`space`, but
	 * `data` is not always populated yet; `exists` lets callers confirm the
	 * account is live even if byte-level decoding isn't available yet.
	 */
	readKvStore(
		runtime: Runtime<unknown>,
		kvStoreAccount: string | Uint8Array,
	): { exists: boolean; account: KvStoreAccount | null } {
		const accountBytes =
			typeof kvStoreAccount === 'string' ? solanaAddressToBytes(kvStoreAccount) : kvStoreAccount

		const reply = this.client
			.getAccountInfoWithOpts(runtime, {
				account: hexToBase64(bytesToHex(accountBytes)),
				opts: {
					encoding: 'ENCODING_TYPE_BASE64',
					commitment: 'COMMITMENT_TYPE_CONFIRMED',
				},
			})
			.result()

		if (!reply.value || reply.value.lamports <= 0n) {
			return { exists: false, account: null }
		}

		const account =
			reply.value.data?.body.case === 'raw' ? decodeKvStoreAccount(reply.value.data.body.value) : null

		return { exists: true, account }
	}
}
