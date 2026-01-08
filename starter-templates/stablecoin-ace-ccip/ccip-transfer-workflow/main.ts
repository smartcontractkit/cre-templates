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
import { encodeAbiParameters, parseAbiParameters, getAddress } from 'viem'
import { z } from 'zod'

const configSchema = z.object({
	chains: z.record(
		z.string(),
		z.object({
			chainSelector: z.string(),
			consumerAddress: z.string(),
			stablecoinAddress: z.string(),
		})
	),
	gasLimit: z.string(),
})

type Config = z.infer<typeof configSchema>

// SWIFT-style cross-chain transfer payload schema
const transferPayloadSchema = z.object({
	messageType: z.string(),
	bankReference: z.string(),
	valueDate: z.string(),
	amount: z.string(),
	currency: z.string(),
	sourceChain: z.string(),
	destinationChain: z.string(),
	sender: z.object({
		account: z.string(),
		name: z.string(),
	}),
	beneficiary: z.object({
		account: z.string(),
		name: z.string(),
	}),
	instructionCode: z.string(),
})

type TransferPayload = z.infer<typeof transferPayloadSchema>

// Utility function to safely stringify objects with bigints
const safeJsonStringify = (obj: any): string =>
	JSON.stringify(obj, (_, value) => (typeof value === 'bigint' ? value.toString() : value), 2)

const submitCCIPTransfer = (
	runtime: Runtime<Config>,
	evmClient: cre.capabilities.EVMClient,
	consumerAddress: string,
	sender: string,
	recipient: string,
	amount: bigint,
	destinationChainSelector: bigint,
	bankRef: string,
): string => {
	runtime.log(`Submitting CCIP transfer - Sender: ${sender}, Recipient: ${recipient}, Amount: ${amount.toString()}, Destination: ${destinationChainSelector.toString()}, BankRef: ${bankRef}`)

	// 1. Normalize addresses with checksum
	const checksummedSender = getAddress(sender)
	const checksummedRecipient = getAddress(recipient)

	// 2. Convert bankRef string to bytes32
	const bankRefBytes32 = bankRef.padEnd(32, '\0').slice(0, 32)
	const bankRefHex = `0x${Buffer.from(bankRefBytes32).toString('hex')}` as `0x${string}`

	// 3. Encode CCIP transfer instruction with sender
	const reportData = encodeAbiParameters(
		parseAbiParameters('uint64 destinationChainSelector, address sender, address recipient, uint256 amount, bytes32 bankRef'),
		[destinationChainSelector, checksummedSender, checksummedRecipient, amount, bankRefHex],
	)

	runtime.log(`Encoded transfer data: ${reportData}`)

	// 4. Generate DON-signed report
	const reportResponse = runtime
		.report({
			encodedPayload: hexToBase64(reportData),
			encoderName: 'evm',
			signingAlgo: 'ecdsa',
			hashingAlgo: 'keccak256',
		})
		.result()

	// 5. Submit report to CCIPTransferConsumer
	const resp = evmClient
		.writeReport(runtime, {
			receiver: consumerAddress,
			report: reportResponse,
			gasConfig: {
				gasLimit: runtime.config.gasLimit,
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

const processCCIPTransfer = (runtime: Runtime<Config>, evmClient: cre.capabilities.EVMClient, transferData: TransferPayload): string => {
	runtime.log(`Processing CCIP transfer from bank: ${transferData.bankReference}`)

	// Look up source and destination chain configs
	const sourceChainConfig = runtime.config.chains[transferData.sourceChain]
	const destChainConfig = runtime.config.chains[transferData.destinationChain]

	if (!sourceChainConfig) {
		throw new Error(`Source chain not found in config: ${transferData.sourceChain}`)
	}
	if (!destChainConfig) {
		throw new Error(`Destination chain not found in config: ${transferData.destinationChain}`)
	}

	// Convert amount from string to scaled bigint
	const amountScaled = BigInt(parseFloat(transferData.amount) * 1e18)
	runtime.log(`Amount scaled: ${amountScaled.toString()}`)
	runtime.log(`Source: ${transferData.sourceChain}`)
	runtime.log(`Destination: ${transferData.destinationChain}`)
	runtime.log(`Recipient: ${transferData.beneficiary.account}`)

	// Get destination chain selector
	const destChainSelector = BigInt(destChainConfig.chainSelector)

	submitCCIPTransfer(
		runtime,
		evmClient,
		sourceChainConfig.consumerAddress,
		transferData.sender.account,
		transferData.beneficiary.account,
		amountScaled,
		destChainSelector,
		transferData.bankReference,
	)

	return `CCIP transfer instruction processed: ${transferData.amount} ${transferData.currency} from ${transferData.sourceChain} to ${transferData.destinationChain}`
}

const onHTTPTrigger = (runtime: Runtime<Config>, payload: HTTPPayload): string => {
	runtime.log('Raw HTTP trigger received for CCIP transfer')

	// Require payload
	if (!payload.input || payload.input.length === 0) {
		throw new Error('HTTP trigger payload is required')
	}

	// Log the raw JSON
	runtime.log(`Payload bytes: ${payload.input.toString()}`)

	try {
		// Parse cross-chain transfer payload
		const payloadJson = JSON.parse(payload.input.toString())
		const transferData = transferPayloadSchema.parse(payloadJson)
		
		runtime.log(`Parsed CCIP transfer payload: ${safeJsonStringify(transferData)}`)
		runtime.log(`${transferData.instructionCode} instruction - From: ${transferData.sourceChain}, To: ${transferData.destinationChain}, Amount: ${transferData.amount} ${transferData.currency}, BankRef: ${transferData.bankReference}`)
		
		// Get source chain network and create evmClient dynamically
		const sourceNetwork = getNetwork({
			chainFamily: 'evm',
			chainSelectorName: transferData.sourceChain,
			isTestnet: true,
		})

		if (!sourceNetwork) {
			throw new Error(`Network not found for source chain: ${transferData.sourceChain}`)
		}

		const evmClient = new cre.capabilities.EVMClient(sourceNetwork.chainSelector.selector)
		
		return processCCIPTransfer(runtime, evmClient, transferData)
	} catch (error) {
		runtime.log(`Failed to parse HTTP trigger payload: ${error}`)
		throw new Error('Failed to parse HTTP trigger payload')
	}
}

const initWorkflow = (config: Config) => {
	const httpTrigger = new cre.capabilities.HTTPCapability()

	return [
		cre.handler(httpTrigger.trigger({}), onHTTPTrigger),
	]
}

export async function main() {
	const runner = await Runner.newRunner<Config>({
		configSchema,
	})
	await runner.run(initWorkflow)
}

main()
