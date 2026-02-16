import { 
	bytesToHex,
	cre,
	getNetwork,
	type HTTPPayload,
	hexToBase64,
	Runner,
	type Runtime,
	TxStatus,
} from '@chainlink/cre-sdk'
import { encodeAbiParameters, getAddress, parseAbiParameters, parseUnits } from 'viem'
import { z } from 'zod'

const configSchema = z.object({
	evms: z.array(
		z.object({
			stablecoinAddress: z.string(),
			consumerAddress: z.string(),
			chainSelectorName: z.string(),
			gasLimit: z.string(),
		}),
	),
})

type Config = z.infer<typeof configSchema>

// SWIFT MT103-style payload schema
const swiftPayloadSchema = z.object({
	messageType: z.string(),
	bankReference: z.string(),
	valueDate: z.string(),
	amount: z.string(),
	currency: z.string(),
	beneficiary: z.object({
		account: z.string(),
		name: z.string(),
	}),
	instructionCode: z.string(),
})

type SWIFTPayload = z.infer<typeof swiftPayloadSchema>

// Instruction type constants (match MintingConsumer.sol)
const INSTRUCTION_MINT = 1
const INSTRUCTION_REDEEM = 2

// Utility function to safely stringify objects with bigints
const safeJsonStringify = (obj: any): string =>
	JSON.stringify(obj, (_, value) => (typeof value === 'bigint' ? value.toString() : value), 2)

const submitBankInstruction = (
	runtime: Runtime<Config>,
	evmClient: cre.capabilities.EVMClient,
	instructionType: number,
	account: string,
	amount: bigint,
	bankRef: string,
): string => {
	const evmConfig = runtime.config.evms[0]
	const instructionName = instructionType === INSTRUCTION_MINT ? 'MINT' : 'REDEEM'

	runtime.log(`Submitting ${instructionName} instruction - Account: ${account}, Amount: ${amount.toString()}, BankRef: ${bankRef}`)

	// 1. Normalize address with checksum
	const checksummedAccount = getAddress(account)

	// 2. Convert bankRef string to bytes32 (pad/truncate to 32 bytes)
	const bankRefBytes32 = bankRef.padEnd(32, '\0').slice(0, 32)
	const bankRefHex = `0x${Buffer.from(bankRefBytes32).toString('hex')}` as `0x${string}`

	// 3. Encode instruction data with type identifier
	const reportData = encodeAbiParameters(
		parseAbiParameters('uint8 instructionType, address account, uint256 amount, bytes32 bankRef'),
		[instructionType, checksummedAccount, amount, bankRefHex],
	)

	runtime.log(`Encoded report data: ${reportData}`)

	// 4. Generate DON-signed report using consensus capability
	const reportResponse = runtime
		.report({
			encodedPayload: hexToBase64(reportData),
			encoderName: 'evm',
			signingAlgo: 'ecdsa',
			hashingAlgo: 'keccak256',
		})
		.result()

	// 5. Submit report to consumer contract via Forwarder
	const resp = evmClient
		.writeReport(runtime, {
			receiver: evmConfig.consumerAddress,
			report: reportResponse,
			gasConfig: {
				gasLimit: evmConfig.gasLimit,
			},
		})
		.result()

	const txStatus = resp.txStatus

	if (txStatus !== TxStatus.SUCCESS) {
		throw new Error(`Failed to write report: ${resp.errorMessage || txStatus}`)
	}

	const txHash = resp.txHash || new Uint8Array(32)
	const txHashHex = bytesToHex(txHash)

	runtime.log(`⚠️  Report delivered to consumer at txHash: ${txHashHex}`)
	runtime.log(`   Verify execution: https://sepolia.etherscan.io/tx/${txHashHex}`)

	return txHashHex
}

const processBankInstruction = (runtime: Runtime<Config>, evmClient: cre.capabilities.EVMClient, swiftData: SWIFTPayload): string => {
	runtime.log(`Processing ${swiftData.instructionCode} instruction from bank: ${swiftData.bankReference}`)

	// Convert amount from string to scaled bigint (e.g., "1000.00" -> 1000000000000000000000)
	const amountScaled = parseUnits(swiftData.amount, 18)
	runtime.log(`Amount scaled: ${amountScaled.toString()}`)
	runtime.log(`Account: ${swiftData.beneficiary.account}`)

	// Determine instruction type
	let instructionType: number
	if (swiftData.instructionCode === 'MINT') {
		instructionType = INSTRUCTION_MINT
	} else if (swiftData.instructionCode === 'REDEEM') {
		instructionType = INSTRUCTION_REDEEM
	} else {
		throw new Error(`Invalid instruction code: ${swiftData.instructionCode}. Must be MINT or REDEEM`)
	}

	submitBankInstruction(
		runtime,
		evmClient,
		instructionType,
		swiftData.beneficiary.account,
		amountScaled,
		swiftData.bankReference,
	)

	return `${swiftData.instructionCode} instruction processed: ${swiftData.amount} ${swiftData.currency}`
}

const onHTTPTrigger = (runtime: Runtime<Config>, evmClient: cre.capabilities.EVMClient, payload: HTTPPayload): string => {
	runtime.log('Raw HTTP trigger received')

	// Require payload
	if (!payload.input || payload.input.length === 0) {
		throw new Error('HTTP trigger payload is required')
	}

	// Log the raw JSON for debugging
	runtime.log(`Payload bytes: ${payload.input.toString()}`)

	try {
		// Parse SWIFT MT103 payload
		const payloadJson = JSON.parse(payload.input.toString())
		const swiftData = swiftPayloadSchema.parse(payloadJson)
		
		runtime.log(`Parsed SWIFT payload: ${safeJsonStringify(swiftData)}`)
		runtime.log(`${swiftData.instructionCode} instruction - Account: ${swiftData.beneficiary.account}, Amount: ${swiftData.amount} ${swiftData.currency}, BankRef: ${swiftData.bankReference}`)
		
		return processBankInstruction(runtime, evmClient, swiftData)
	} catch (error) {
		runtime.log(`Failed to parse HTTP trigger payload: ${error}`)
		throw new Error('Failed to parse HTTP trigger payload')
	}
}

const initWorkflow = (config: Config) => {
	const httpTrigger = new cre.capabilities.HTTPCapability()
	
	// Initialize EVM client for the configured chain
	const network = getNetwork({
		chainFamily: 'evm',
		chainSelectorName: config.evms[0].chainSelectorName,
		isTestnet: true,
	})

	if (!network) {
		throw new Error(
			`Network not found for chain selector name: ${config.evms[0].chainSelectorName}`,
		)
	}

	const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)

	return [
		cre.handler(httpTrigger.trigger({}), (runtime, payload) => 
			onHTTPTrigger(runtime, evmClient, payload)
		),
	]
}

export async function main() {
	const runner = await Runner.newRunner<Config>({
		configSchema,
	})
	await runner.run(initWorkflow)
}

main()
