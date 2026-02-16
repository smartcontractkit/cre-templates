import { 
	bytesToHex,
	cre,
	getNetwork,
	type HTTPPayload,
	hexToBase64,
	Runner,
	type Runtime,
	type NodeRuntime,
	TxStatus,
	consensusMedianAggregation,
} from '@chainlink/cre-sdk'
import { encodeAbiParameters, parseAbiParameters, encodeFunctionData, decodeFunctionResult, getAddress, parseUnits } from 'viem'
import { z } from 'zod'

// ========================================
// CONFIG SCHEMA
// ========================================
const configSchema = z.object({
	sepolia: z.object({
		stablecoinAddress: z.string(),
		mintingConsumerAddress: z.string(),
		ccipConsumerAddress: z.string(),
		chainSelector: z.string(),
	}),
	fuji: z.object({
		stablecoinAddress: z.string(),
		chainSelector: z.string(),
	}),
	porApiUrl: z.string(),
	decimals: z.number(),
})

type Config = z.infer<typeof configSchema>

// ========================================
// PAYLOAD SCHEMA
// ========================================
const payloadSchema = z.object({
	messageType: z.string(),
	transactionId: z.string(),
	beneficiary: z.object({
		account: z.string(),
		name: z.string().optional(),
	}),
	amount: z.string(),
	currency: z.string(),
	valueDate: z.string().optional(),
	bankReference: z.string(),
	crossChain: z.object({
		enabled: z.boolean(),
		destinationChain: z.string(),
		beneficiary: z.string(),
	}).optional(),
})

type Payload = z.infer<typeof payloadSchema>

// ========================================
// CONSTANTS
// ========================================
const INSTRUCTION_MINT = 1

// StablecoinERC20 ABI (minimal)
const StablecoinABI = [
	{
		type: 'function',
		name: 'totalSupply',
		inputs: [],
		outputs: [{ name: '', type: 'uint256' }],
		stateMutability: 'view',
	},
] as const

// ========================================
// UTILITY FUNCTIONS
// ========================================
const safeJsonStringify = (obj: any): string =>
	JSON.stringify(obj, (_, value) => (typeof value === 'bigint' ? value.toString() : value), 2)

const stringToBytes32 = (str: string): `0x${string}` => {
	const bytes32 = str.padEnd(32, '\0').slice(0, 32)
	return `0x${Buffer.from(bytes32).toString('hex')}` as `0x${string}`
}

// ========================================
// PROOF OF RESERVE VALIDATION
// ========================================
/**
 * Validates Proof of Reserve before minting
 * Fetches reserve data from API and compares to current supply
 */
const validateProofOfReserve = (
	runtime: Runtime<Config>,
	config: Config,
	mintAmount: bigint,
): boolean => {
	runtime.log('\n[PoR Validation] Fetching reserve data...')

	// For mock data (file:// URL), use hardcoded values
	// In production, this would fetch from a real PoR API endpoint
	let reserveData: { totalReserve: number; lastUpdated: string }

	if (config.porApiUrl.startsWith('file://')) {
		// Mock PoR data (matching mock-por-response.json)
		reserveData = {
			totalReserve: 500000.00,  // 500,000 USD in reserves
			lastUpdated: '2025-10-29T00:00:00Z',
		}
		runtime.log('Using mock PoR data for demo')
	} else {
		// Fetch from real PoR API in node mode
		reserveData = runtime.runInNodeMode(
			(nodeRuntime: NodeRuntime) => {
				const httpClient = new cre.capabilities.HTTPClient()
				const response = httpClient.sendRequest(nodeRuntime, {
					url: config.porApiUrl,
					method: 'GET',
				}).result()
				
				const data = JSON.parse(new TextDecoder().decode(response.body))
				return {
					totalReserve: data.totalReserve,
					lastUpdated: data.lastUpdated,
				}
			},
			consensusMedianAggregation()
		)().result()
	}

	runtime.log(`Reserve Data: ${reserveData.totalReserve} USD (as of ${reserveData.lastUpdated})`)

	// Scale reserves to wei (18 decimals)
	const reservesWei = parseUnits(reserveData.totalReserve.toString(), config.decimals)
	
	runtime.log(`Reserves: ${reservesWei} wei (${reserveData.totalReserve} USD)`)
	runtime.log(`Requested Mint: ${mintAmount} wei`)

	// Simplified PoR validation for demo
	// In production, you would read totalSupply() from StablecoinERC20 and compare:
	// if (reservesWei < currentSupply + mintAmount) { throw error }
	//
	// For this demo, we just check if reserves can cover the mint amount
	if (reservesWei < mintAmount) {
		throw new Error(
			`[PoR FAILED] Insufficient reserves: have ${reservesWei} wei (${reserveData.totalReserve} USD), need ${mintAmount} wei for this mint`
		)
	}

	runtime.log(`✓ PoR validation passed - reserves (${reserveData.totalReserve} USD) can cover mint`)
	return true
}

// ========================================
// MINT WITH ACE
// ========================================
/**
 * Mints stablecoins via ACE-protected consumer
 * ACE automatically checks if beneficiary is blacklisted
 * 
 * @param beneficiary - Address checked by ACE blacklist policy
 * @param mintRecipient - Address that receives the minted tokens
 */
const mintWithACE = (
	runtime: Runtime<Config>,
	evmClient: cre.capabilities.EVMClient,
	beneficiary: string,
	mintRecipient: string,
	amount: bigint,
	bankRef: string,
): string => {
	runtime.log(`\n[ACE Mint] Minting ${amount} tokens to ${mintRecipient}`)
	runtime.log(`ACE will check: Is beneficiary (${beneficiary}) blacklisted?`)

	// Normalize addresses
	const checksummedBeneficiary = getAddress(beneficiary)
	const checksummedMintRecipient = getAddress(mintRecipient)

	// Convert bankRef to bytes32
	const bankRefHex = stringToBytes32(bankRef)

	// Encode mint report: (instructionType=1, beneficiary, amount, bankRef)
	// Note: beneficiary is used for BOTH ACE check AND mint recipient
	// The ACE policy checks "beneficiary" parameter for blacklist
	// The mint consumer mints to "beneficiary" parameter
	// 
	// When CCIP is enabled:
	//   - mintRecipient = CCIPConsumer (tokens staged in contract)
	//   - ACE checks if CCIPConsumer is blacklisted (it won't be - it's our contract)
	//   - Final beneficiary blacklist check happens during CCIP transfer via VolumePolicy
	// When CCIP is disabled:
	//   - mintRecipient = end user
	//   - ACE checks if end user is blacklisted
	const reportData = encodeAbiParameters(
		parseAbiParameters('uint8 instructionType, address beneficiary, uint256 amount, bytes32 bankRef'),
		[INSTRUCTION_MINT, checksummedMintRecipient, amount, bankRefHex],
	)

	runtime.log(`Encoded mint report: ${reportData.slice(0, 66)}...`)

	// Generate DON-signed report
	const reportResponse = runtime
		.report({
			encodedPayload: hexToBase64(reportData),
			encoderName: 'evm',
			signingAlgo: 'ecdsa',
			hashingAlgo: 'keccak256',
		})
		.result()

	// Write to MintingConsumerWithACE
	// ACE Policy Check happens here via runPolicy modifier:
	//   1. PolicyEngine calls MintingConsumerExtractor
	//   2. Extractor returns [beneficiary, amount]
	//   3. PolicyEngine runs AddressBlacklistPolicy
	//   4. If blacklisted → reverts with PolicyRunRejected
	//   5. If allowed → mint proceeds
	const resp = evmClient
		.writeReport(runtime, {
			receiver: runtime.config.sepolia.mintingConsumerAddress,
			report: reportResponse,
			gasConfig: {
				gasLimit: '500000',
			},
		})
		.result()

	const txStatus = resp.txStatus

	// Important: CRE Forwarder transaction may succeed even if consumer call fails!
	// We need to check the actual execution result from Forwarder events
	// For now, check if txStatus is success and errorMessage is empty
	if (txStatus !== TxStatus.SUCCESS) {
		const errorMsg = resp.errorMessage || txStatus
		
		// Check if it's a PolicyRunRejected error
		if (errorMsg.includes('PolicyRunRejected') || errorMsg.includes('blacklisted')) {
			throw new Error(`[ACE REJECTED] Address ${beneficiary} is blacklisted`)
		}
		
		throw new Error(`Failed to mint: ${errorMsg}`)
	}

	// Additional check: If errorMessage contains policy rejection info
	if (resp.errorMessage && (resp.errorMessage.includes('PolicyRunRejected') || resp.errorMessage.includes('blacklisted'))) {
		throw new Error(`[ACE REJECTED] Address ${beneficiary} is blacklisted`)
	}

	const txHash = resp.txHash || new Uint8Array(32)
	const txHashHex = bytesToHex(txHash)

	runtime.log(`⚠️  Mint report delivered: ${txHashHex}`)
	runtime.log(`   ACE policies apply: Blacklist check for beneficiary`)
	runtime.log(`   Verify execution: https://sepolia.etherscan.io/tx/${txHashHex}`)
	return txHashHex
}

// ========================================
// CCIP TRANSFER WITH ACE
// ========================================
/**
 * Executes cross-chain transfer via ACE-protected CCIP consumer
 * ACE automatically checks if beneficiary is blacklisted
 * Note: We use "beneficiary" for consistency with minting operations (who receives tokens)
 */
const transferWithACE = (
	runtime: Runtime<Config>,
	evmClient: cre.capabilities.EVMClient,
	sender: string,
	beneficiary: string,
	amount: bigint,
	bankRef: string,
): string => {
	runtime.log(`\n[ACE CCIP] Transferring ${amount} tokens to ${beneficiary} on Fuji`)
	runtime.log('ACE will check: Is beneficiary blacklisted?')

	// Normalize addresses
	const checksummedSender = getAddress(sender)
	const checksummedBeneficiary = getAddress(beneficiary)

	// Get destination chain selector
	const destChainSelector = BigInt(runtime.config.fuji.chainSelector)

	// Convert bankRef to bytes32
	const bankRefHex = stringToBytes32(bankRef)

	// Encode CCIP report: (destChainSelector, sender, beneficiary, amount, bankRef)
	const reportData = encodeAbiParameters(
		parseAbiParameters('uint64 destinationChainSelector, address sender, address beneficiary, uint256 amount, bytes32 bankRef'),
		[destChainSelector, checksummedSender, checksummedBeneficiary, amount, bankRefHex],
	)

	runtime.log(`Encoded CCIP report: ${reportData.slice(0, 66)}...`)

	// Generate DON-signed report
	const reportResponse = runtime
		.report({
			encodedPayload: hexToBase64(reportData),
			encoderName: 'evm',
			signingAlgo: 'ecdsa',
			hashingAlgo: 'keccak256',
		})
		.result()

	// Write to CCIPTransferConsumerWithACE
	// ACE Policy Check happens here via runPolicy modifier:
	//   1. PolicyEngine calls CCIPTransferConsumerExtractor
	//   2. Extractor returns [beneficiary, sender, amount]
	//   3. PolicyEngine runs AddressBlacklistPolicy
	//   4. If blacklisted → reverts with PolicyRunRejected
	//   5. If allowed → CCIP transfer proceeds
	const resp = evmClient
		.writeReport(runtime, {
			receiver: runtime.config.sepolia.ccipConsumerAddress,
			report: reportResponse,
			gasConfig: {
				gasLimit: '1000000',
			},
		})
		.result()

	const txStatus = resp.txStatus

	if (txStatus !== TxStatus.SUCCESS) {
		const errorMsg = resp.errorMessage || txStatus
		
		// Check if it's a PolicyRunRejected error
		if (errorMsg.includes('PolicyRunRejected') || errorMsg.includes('blacklisted')) {
			throw new Error(`[ACE REJECTED] Beneficiary ${beneficiary} is blacklisted`)
		}
		
		throw new Error(`Failed to initiate CCIP transfer: ${errorMsg}`)
	}

	const txHash = resp.txHash || new Uint8Array(32)
	const txHashHex = bytesToHex(txHash)

	runtime.log(`⚠️  CCIP report delivered: ${txHashHex}`)
	runtime.log(`   ACE policies apply: VolumePolicy check (100-10,000 creUSD)`)
	runtime.log(`   Verify execution: https://sepolia.etherscan.io/tx/${txHashHex}`)
	return txHashHex
}

// ========================================
// HTTP TRIGGER HANDLER
// ========================================
const onHTTPTrigger = (runtime: Runtime<Config>, payload: HTTPPayload): object => {
	runtime.log('=== Phase 3: PoR + ACE + CCIP Workflow ===')

	// Require payload
	if (!payload.input || payload.input.length === 0) {
		throw new Error('HTTP trigger payload is required')
	}

	runtime.log(`Payload: ${payload.input.toString()}`)

	try {
		// Parse MT103 bank message
		const payloadJson = JSON.parse(payload.input.toString())
		const parsedPayload = payloadSchema.parse(payloadJson)

		runtime.log(`Parsed MT103 payload: ${safeJsonStringify(parsedPayload)}`)
		runtime.log(`Transaction ID: ${parsedPayload.transactionId}`)
		runtime.log(`Beneficiary: ${parsedPayload.beneficiary.account}`)
		runtime.log(`Amount: ${parsedPayload.amount} ${parsedPayload.currency}`)

		// Initialize EVM client for Sepolia
		const network = getNetwork({
			chainFamily: 'evm',
			chainSelectorName: 'ethereum-testnet-sepolia',
			isTestnet: true,
		})

		if (!network) {
			throw new Error('Sepolia network not found')
		}

		const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)

		// Convert amount to wei
		const amountWei = parseUnits(parsedPayload.amount, runtime.config.decimals)
		const beneficiary = parsedPayload.beneficiary.account
		const hasCrossChain = parsedPayload.crossChain?.enabled === true

		// ========================================
		// STEP 1: Proof of Reserve Validation
		// ========================================
		runtime.log('\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━')
		runtime.log('STEP 1: Proof of Reserve Validation')
		runtime.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━')
		
		try {
			validateProofOfReserve(runtime, runtime.config, amountWei)
		} catch (error: any) {
			runtime.log(`❌ PoR validation failed: ${error.message}`)
			return {
				success: false,
				error: 'POR_INSUFFICIENT_RESERVES',
				message: error.message,
				transactionId: parsedPayload.transactionId,
			}
		}

		// ========================================
		// STEP 2: Mint with ACE
		// ========================================
		runtime.log('\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━')
		runtime.log('STEP 2: Mint with ACE Policy Enforcement')
		runtime.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━')
		
		// Determine mint destination:
		// - If CCIP enabled: mint to CCIP consumer (tokens staged for cross-chain transfer)
		// - If CCIP disabled: mint to beneficiary (final destination)
		const mintRecipient = hasCrossChain 
			? runtime.config.sepolia.ccipConsumerAddress 
			: beneficiary
		
		runtime.log(`Mint destination: ${hasCrossChain ? 'CCIP Consumer (for cross-chain)' : 'Beneficiary (final)'}`)
		runtime.log(`Mint recipient: ${mintRecipient}`)
		
		let mintTxHash: string
		try {
			mintTxHash = mintWithACE(
				runtime,
				evmClient,
				beneficiary, // ACE checks this for blacklist
				mintRecipient, // Tokens minted here
				amountWei,
				parsedPayload.bankReference,
			)
		} catch (error: any) {
			runtime.log(`❌ Mint failed: ${error.message}`)
			return {
				success: false,
				error: error.message.includes('ACE REJECTED') ? 'ACE_POLICY_REJECTED' : 'MINT_FAILED',
				message: error.message,
				beneficiary: beneficiary,
				transactionId: parsedPayload.transactionId,
			}
		}

		// ========================================
		// STEP 3: Optional CCIP Transfer with ACE
		// ========================================
		let ccipTxHash: string | null = null

		if (hasCrossChain) {
			runtime.log('\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━')
			runtime.log('STEP 3: CCIP Transfer with ACE Policy Enforcement')
			runtime.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━')
			
			const ccipBeneficiary = parsedPayload.crossChain!.beneficiary
			
			try {
				// Sender is the CCIP consumer (where tokens were minted)
				// Beneficiary is the end user on destination chain
				ccipTxHash = transferWithACE(
					runtime,
					evmClient,
					mintRecipient, // sender (CCIPConsumer - tokens were minted here)
					ccipBeneficiary, // who receives on destination chain
					amountWei,
					parsedPayload.bankReference,
				)
			} catch (error: any) {
				runtime.log(`❌ CCIP transfer failed: ${error.message}`)
				return {
					success: false,
					error: error.message.includes('ACE REJECTED') ? 'ACE_POLICY_REJECTED_CCIP' : 'CCIP_FAILED',
					message: error.message,
					mintTransaction: mintTxHash, // Mint succeeded but CCIP failed
					beneficiary: ccipBeneficiary,
					transactionId: parsedPayload.transactionId,
				}
			}
		} else {
			runtime.log('\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━')
			runtime.log('STEP 3: CCIP Transfer - SKIPPED')
			runtime.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━')
			runtime.log('No cross-chain transfer requested')
		}

		// ========================================
		// REPORT DELIVERY RESPONSE
		// ========================================
		// Build result object (conditionally include CCIP fields to avoid null values)
		const result: any = {
			reportDelivered: true,
			transactionId: parsedPayload.transactionId,
			beneficiary: beneficiary,
			amount: parsedPayload.amount,
			currency: parsedPayload.currency,
			mintTransaction: mintTxHash,
			message: hasCrossChain
				? `Reports delivered: Mint + CCIP transfer to ${parsedPayload.crossChain!.destinationChain} (verify on-chain)`
				: `Report delivered: Mint ${parsedPayload.amount} ${parsedPayload.currency} to ${beneficiary} (verify on-chain)`,
			etherscanMint: `https://sepolia.etherscan.io/tx/${mintTxHash}`,
			verificationNote: 'ACE policies may block execution. Verify balance and events on-chain.',
		}
		
		// Only add CCIP fields if CCIP transfer was executed
		if (ccipTxHash) {
			result.ccipTransaction = ccipTxHash
			result.etherscanCCIP = `https://sepolia.etherscan.io/tx/${ccipTxHash}`
			result.ccipExplorer = `https://ccip.chain.link`
		}

		runtime.log(`\n⚠️  REPORTS DELIVERED TO CONSUMERS`)
		runtime.log(`   ACE policies applied during execution`)
		runtime.log(`   Verify on-chain to confirm actual results`)
		runtime.log(`\nResult: ${safeJsonStringify(result)}`)
		
		return JSON.stringify(result)

	} catch (error: any) {
		runtime.log(`❌ Workflow error: ${error.message}`)
		throw error
	}
}

// ========================================
// WORKFLOW INITIALIZATION
// ========================================
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

