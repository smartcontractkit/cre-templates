/**
 * Browser-side validation of a Chainlink CRE off-chain report.
 *
 * This mirrors the verification logic of the CRE TypeScript SDK
 * (cre-sdk/src/sdk/report.ts), but instead of reading the DON configuration
 * through the CRE runtime's EVM capability, it reads it directly on-chain via
 * a plain JSON-RPC `eth_call` — which is all a browser needs.
 *
 * Verification steps:
 *  1. Parse the 109-byte metadata header of the raw report.
 *  2. Compute the report hash: keccak256(keccak256(rawReport) ++ reportContext).
 *  3. Read the DON config (f + signer addresses) from the capabilities registry
 *     contract via `getDON(donId)` and `getNodesByP2PIds(bytes32[])`.
 *  4. Recover the signer of every attached signature and require at least f+1
 *     unique signers that belong to the DON.
 */
import {
	type Address,
	concatHex,
	createPublicClient,
	defineChain,
	getAddress,
	type Hex,
	hexToBytes,
	http,
	keccak256,
	recoverAddress,
	toHex,
} from 'viem'

// --- Capabilities registry contract (function selectors) ---------------------

// getDON(uint32)
const GET_DON_SELECTOR = new Uint8Array([0x23, 0x53, 0x74, 0x05])
// getNodesByP2PIds(bytes32[])
const GET_NODES_BY_P2P_IDS_SELECTOR = new Uint8Array([0x05, 0xa5, 0x19, 0x66])

/** Production environment: capabilities registry on Ethereum mainnet. */
export const PRODUCTION_ENVIRONMENT: Environment = {
	chainId: 1,
	registryAddress: '0x76c9cf548b4179F8901cda1f8623568b58215E62',
	rpcUrl: 'https://ethereum-rpc.publicnode.com',
}

export type Environment = {
	/** EIP-155 chain ID where the capabilities registry lives. */
	chainId: number
	/** Capabilities registry contract address. */
	registryAddress: string
	/** JSON-RPC endpoint used for the eth_call reads. */
	rpcUrl: string
}

export type DonInfo = {
	f: number
	/** checksumed signer address -> node operator id */
	signers: Map<Address, number>
}

export type ReportInput = {
	/** The full raw report (metadata header + body), hex encoded. */
	rawReport: Hex
	/** The report context attached to the report, hex encoded. */
	reportContext: Hex
	/** 65-byte ECDSA signatures (r,s,v) over the report hash, hex encoded. */
	signatures: Hex[]
}

// --- Errors (mirroring the SDK) ----------------------------------------------

export class NullReportError extends Error {
	constructor() {
		super('null report')
	}
}

export class RawReportTooShortError extends Error {
	constructor(need: number, got: number) {
		super(`raw report too short to contain metadata header: need ${need} bytes, got ${got}`)
	}
}

// --- Report metadata header ---------------------------------------------------

export const REPORT_METADATA_HEADER_LENGTH = 109

const REPORT_METADATA_OFFSETS = {
	version: 0,
	versionSize: 1,
	executionId: 1,
	executionIdSize: 32,
	timestamp: 33,
	timestampSize: 4,
	donId: 37,
	donIdSize: 4,
	donConfigVersion: 41,
	donConfigVersionSize: 4,
	workflowId: 45,
	workflowIdSize: 32,
	workflowName: 77,
	workflowNameSize: 10,
	workflowOwner: 87,
	workflowOwnerSize: 20,
	reportId: 107,
	reportIdSize: 2,
	bodyStart: 109,
} as const

export type ReportMetadataHeader = {
	version: number
	executionId: string
	timestamp: number
	donId: number
	donConfigVersion: number
	workflowId: string
	workflowName: string
	workflowOwner: string
	reportId: string
	body: Uint8Array
}

function encodeHexLower(bytes: Uint8Array): string {
	let out = ''
	for (let i = 0; i < bytes.length; i++) {
		out += bytes[i].toString(16).padStart(2, '0')
	}
	return out
}

function readUint32BE(raw: Uint8Array, offset: number): number {
	return new DataView(raw.buffer, raw.byteOffset + offset, 4).getUint32(0, false)
}

export function parseReportMetadataHeader(raw: Uint8Array | null | undefined): ReportMetadataHeader {
	if (raw === undefined || raw === null) {
		throw new NullReportError()
	}
	if (raw.length < REPORT_METADATA_HEADER_LENGTH) {
		throw new RawReportTooShortError(REPORT_METADATA_HEADER_LENGTH, raw.length)
	}

	const o = REPORT_METADATA_OFFSETS
	const workflowNameBytes = raw.subarray(o.workflowName, o.workflowName + o.workflowNameSize)

	return {
		version: raw[o.version],
		executionId: encodeHexLower(raw.subarray(o.executionId, o.executionId + o.executionIdSize)),
		timestamp: readUint32BE(raw, o.timestamp),
		donId: readUint32BE(raw, o.donId),
		donConfigVersion: readUint32BE(raw, o.donConfigVersion),
		workflowId: encodeHexLower(raw.subarray(o.workflowId, o.workflowId + o.workflowIdSize)),
		workflowName: new TextDecoder('utf-8', { fatal: false }).decode(workflowNameBytes),
		workflowOwner: encodeHexLower(
			raw.subarray(o.workflowOwner, o.workflowOwner + o.workflowOwnerSize),
		),
		reportId: encodeHexLower(raw.subarray(o.reportId, o.reportId + o.reportIdSize)),
		body: raw.subarray(REPORT_METADATA_OFFSETS.bodyStart),
	}
}

// --- Report hash ----------------------------------------------------------------

export function computeReportHash(rawReport: Uint8Array, reportContext: Uint8Array): Hex {
	const innerHash = keccak256(toHex(rawReport))
	return keccak256(concatHex([innerHash, toHex(reportContext)]))
}

// --- ABI helpers (identical byte layout to what the SDK produces/expects) --------

function padUint256(v: bigint | number): Uint8Array {
	const n = typeof v === 'bigint' ? v : BigInt(v)
	const b = new Uint8Array(32)
	const view = new DataView(b.buffer)
	view.setBigUint64(24, n, false)
	return b
}

function bytesToBigIntBE(word: Uint8Array): bigint {
	let x = 0n
	for (let i = 0; i < word.length; i++) {
		x = (x << 8n) | BigInt(word[i])
	}
	return x
}

function readUint256AsInt(word: Uint8Array): number {
	const b = bytesToBigIntBE(word)
	if (b > BigInt(Number.MAX_SAFE_INTEGER)) {
		throw new Error('ABI uint256 value too large for Number')
	}
	return Number(b)
}

function concatBytes(parts: readonly Uint8Array[]): Uint8Array {
	const total = parts.reduce((n, p) => n + p.length, 0)
	const r = new Uint8Array(total)
	let o = 0
	for (const p of parts) {
		r.set(p, o)
		o += p.length
	}
	return r
}

function encodeGetDONCalldata(donID: number): Hex {
	const padded = new Uint8Array(32)
	const view = new DataView(padded.buffer)
	view.setUint32(28, donID >>> 0, false)
	return toHex(concatBytes([GET_DON_SELECTOR, padded]))
}

function encodeGetNodesByP2PIdsCalldata(p2pIds: Uint8Array[]): Hex {
	const chunks: Uint8Array[] = [
		GET_NODES_BY_P2P_IDS_SELECTOR,
		padUint256(32),
		padUint256(p2pIds.length),
	]
	for (const id of p2pIds) {
		if (id.length !== 32) {
			throw new Error('p2p id must be 32 bytes')
		}
		chunks.push(new Uint8Array(id))
	}
	return toHex(concatBytes(chunks))
}

// --- On-chain DON info via RPC ----------------------------------------------------

const NODE_TUPLE_HEAD_SIZE = 288

async function ethCall(env: Environment, callData: Hex): Promise<Uint8Array> {
	const chain = defineChain({
		id: env.chainId,
		name: `chain-${env.chainId}`,
		nativeCurrency: { name: 'Ether', symbol: 'ETH', decimals: 18 },
		rpcUrls: { default: { http: [env.rpcUrl] } },
	})
	const client = createPublicClient({ chain, transport: http(env.rpcUrl) })
	const registryAddress = getAddress(env.registryAddress)
	const result = await client.call({ to: registryAddress, data: callData })
	if (!result.data || result.data === '0x') {
		throw new Error(
			`empty response from ${env.registryAddress} on chain ${env.chainId} (wrong registry/RPC?)`,
		)
	}
	return hexToBytes(result.data)
}

/**
 * Fetches the DON configuration (fault tolerance f and the set of signer
 * addresses) from the capabilities registry contract over JSON-RPC.
 */
export async function fetchDonInfo(env: Environment, donID: number): Promise<DonInfo> {
	const getDONABI = await ethCall(env, encodeGetDONCalldata(donID))

	if (getDONABI.length < 224) {
		throw new Error(`getDON ABI response too short: ${getDONABI.length} bytes`)
	}

	const f = readUint256AsInt(getDONABI.subarray(96, 128))

	const tupleStart = 32
	const nodeP2PIdsPtr = readUint256AsInt(getDONABI.subarray(192, 224))
	const nodeCountOff = tupleStart + nodeP2PIdsPtr
	if (nodeCountOff + 32 > getDONABI.length) {
		throw new Error('getDON ABI: nodeP2PIds pointer out of range')
	}
	const nodeCount = readUint256AsInt(getDONABI.subarray(nodeCountOff, nodeCountOff + 32))
	if (nodeCountOff + 32 + nodeCount * 32 > getDONABI.length) {
		throw new Error('getDON ABI: nodeP2PIds data out of range')
	}

	const nodeP2PIds: Uint8Array[] = []
	for (let i = 0; i < nodeCount; i++) {
		const start = nodeCountOff + 32 + i * 32
		nodeP2PIds.push(getDONABI.subarray(start, start + 32))
	}

	if (nodeCount === 0) {
		return { f, signers: new Map() }
	}

	const getNodesABI = await ethCall(env, encodeGetNodesByP2PIdsCalldata(nodeP2PIds))

	if (getNodesABI.length < 64) {
		throw new Error(`getNodesByP2PIds ABI response too short: ${getNodesABI.length} bytes`)
	}

	const outerPtr = readUint256AsInt(getNodesABI.subarray(0, 32))
	if (outerPtr + 32 > getNodesABI.length) {
		throw new Error('getNodesByP2PIds ABI: outer pointer out of range')
	}
	const returnedCount = readUint256AsInt(getNodesABI.subarray(outerPtr, outerPtr + 32))

	const signers = new Map<Address, number>()
	for (let i = 0; i < returnedCount; i++) {
		const elemPtrOff = outerPtr + 32 + i * 32
		if (elemPtrOff + 32 > getNodesABI.length) {
			break
		}
		const elemPtr = readUint256AsInt(getNodesABI.subarray(elemPtrOff, elemPtrOff + 32))
		const tupleBase = outerPtr + 32 + elemPtr
		if (tupleBase + NODE_TUPLE_HEAD_SIZE > getNodesABI.length) {
			break
		}
		const nodeOperatorId = Number(
			bytesToBigIntBE(getNodesABI.subarray(tupleBase, tupleBase + 32)) & 0xffffffffn,
		)
		const signerSlot = tupleBase + 3 * 32
		const addrBytes = getNodesABI.subarray(signerSlot, signerSlot + 20)
		const addr = getAddress(toHex(addrBytes))
		signers.set(addr, nodeOperatorId)
	}

	return { f, signers }
}

// --- Signature verification -------------------------------------------------------

export type SignatureCheckStatus =
	| 'ok'
	| 'invalid-length'
	| 'unrecoverable'
	| 'duplicate'
	| 'unknown-signer'

export type SignatureCheck = {
	signature: Hex
	status: SignatureCheckStatus
	/** Recovered signer address (checksummed) when recovery succeeded. */
	recovered?: Address
	/** Node operator id within the DON, when the signer is a DON member. */
	signerId?: number
}

export type VerificationResult = {
	valid: boolean
	header: ReportMetadataHeader
	reportHash: Hex
	/** Required number of unique valid DON signatures: f + 1. */
	requiredSignatures: number
	/** Fault tolerance f read from the registry. */
	f: number
	/** Number of accepted signatures. */
	acceptedSignatures: number
	/** Per-signature outcome, in input order. */
	signatures: SignatureCheck[]
	/** All signer addresses of the DON as read from the registry. */
	donSigners: Address[]
	/** Reasons signatures were skipped, when the f+1 threshold was not met. */
	failures: string[]
}

function addressKeyNo0x(addr: Address): string {
	return addr.slice(2).toLowerCase()
}

/**
 * Verifies a CRE off-chain report end-to-end.
 *
 * Throws when the report is malformed or the chain read fails. Returns a
 * VerificationResult with `valid = false` (plus per-signature detail) when
 * the f+1 signature threshold is not met.
 */
export async function verifyReport(
	input: ReportInput,
	env: Environment = PRODUCTION_ENVIRONMENT,
): Promise<VerificationResult> {
	const rawReport = hexToBytes(input.rawReport)
	const reportContext = hexToBytes(input.reportContext)
	const header = parseReportMetadataHeader(rawReport)
	const reportHash = computeReportHash(rawReport, reportContext)

	const info = await fetchDonInfo(env, header.donId)
	const signers = new Map<string, number>()
	for (const [addr, id] of info.signers) {
		signers.set(addr.slice(2).toLowerCase(), id)
	}

	const required = info.f + 1
	const seen = new Set<string>()
	const checks: SignatureCheck[] = []
	let accepted = 0
	const failures: string[] = []

	if (input.signatures.length < required) {
		failures.push(
			`wrong signature count: got ${input.signatures.length}, need at least f+1 = ${required}`,
		)
	}

	for (const sig of input.signatures) {
		if (accepted === required) {
			break
		}
		const sigBytes = new Uint8Array(hexToBytes(sig))

		if (sigBytes.length !== 65) {
			failures.push('failed to parse signature: not 65 bytes')
			checks.push({ signature: sig, status: 'invalid-length' })
			continue
		}

		const normalized = new Uint8Array(sigBytes)
		if (normalized[64] === 27 || normalized[64] === 28) {
			normalized[64] -= 27
		}

		let recovered: Address
		try {
			recovered = await recoverAddress({
				hash: reportHash,
				signature: toHex(normalized),
			})
		} catch {
			failures.push('failed to recover signer from signature')
			checks.push({ signature: sig, status: 'unrecoverable' })
			continue
		}

		const key = addressKeyNo0x(recovered)
		if (seen.has(key)) {
			failures.push(`duplicate signer ${recovered}`)
			checks.push({ signature: sig, status: 'duplicate', recovered })
			continue
		}
		seen.add(key)

		const nodeOperatorId = signers.get(key)
		if (nodeOperatorId === undefined) {
			failures.push(`unknown signer ${recovered} (not in DON)`)
			checks.push({ signature: sig, status: 'unknown-signer', recovered })
			continue
		}

		checks.push({ signature: sig, status: 'ok', recovered, signerId: nodeOperatorId })
		accepted++
	}

	return {
		valid: accepted >= required,
		header,
		reportHash,
		requiredSignatures: required,
		f: info.f,
		acceptedSignatures: accepted,
		signatures: checks,
		donSigners: [...info.signers.keys()],
		failures,
	}
}
