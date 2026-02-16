import {
	bytesToHex,
	type CronPayload,
	cre,
	encodeCallMsg,
	EVMClient,
	getNetwork,
	hexToBase64,
	LATEST_BLOCK_NUMBER,
	Runner,
	type Runtime,
	TxStatus,
} from '@chainlink/cre-sdk'
import { type Address, decodeFunctionResult, encodeAbiParameters, encodeFunctionData, formatUnits, Hex, parseAbiParameters, zeroAddress } from 'viem'
import { z } from 'zod'
import { MockPool, ProtocolSmartWallet } from '../contracts/abi'

const configSchema = z.object({
	schedule: z.string(),
	minBPSDeltaForRebalance: z.number().min(0),
	evms: z.array(
		z.object({
			assetAddress: z.string(),
			poolAddress: z.string(),
			protocolSmartWalletAddress: z.string(),
			chainName: z.string(),
			gasLimit: z.string(),
		}),
	),
})

type Config = z.infer<typeof configSchema>
type EVMConfig = z.infer<typeof configSchema.shape.evms.element>

/** ===== Math Helpers ===== **/

// Converts APR (in RAY units) to floating APR
const aprInRAYToAPR = (apr: bigint): number => parseFloat(formatUnits(apr, 27))

// Converts APR (in RAY units) to APY (compounded)
const aprInRAYToAPY = (apr: bigint): number => Math.exp(parseFloat(formatUnits(apr, 27))) - 1

// Convert APR diff (RAY) to basis points (divide by 1e23)
const RAY_TO_BPS_DIVISOR = 100000000000000000000000n // 1e23
const aprDiffToBps = (diffRay: bigint): bigint => diffRay / RAY_TO_BPS_DIVISOR

/** ===== Domain Types ===== **/

type Pool = {
	chainName: string
	APR: bigint
	APY: number
	protocolSmartWalletAddress: string
	balance: bigint
}

/** ===== EVM/Contract Helpers ===== **/

const getEvmClientForChain = (evmCfg: EVMConfig) => {
	const network = getNetwork({
		chainFamily: 'evm',
		chainSelectorName: evmCfg.chainName,
		isTestnet: true,
	})
	if (!network) {
		throw new Error(`Network not found for chain selector name: ${evmCfg.chainName}`)
	}
	return new cre.capabilities.EVMClient(network.chainSelector.selector)
}

const readCurrentLiquidityRate = (
	runtime: Runtime<Config>,
	evmCfg: EVMConfig,
	evmClient: EVMClient,
): bigint => {
	const callData = encodeFunctionData({
		abi: MockPool,
		functionName: 'getReserveData',
		args: [evmCfg.assetAddress as Hex],
	})

	const callResult = evmClient
		.callContract(runtime, {
			call: encodeCallMsg({
				from: zeroAddress,
				to: evmCfg.poolAddress as Address,
				data: callData,
			}),
			blockNumber: LATEST_BLOCK_NUMBER, // warn: use finalized or safe in prod
		})
		.result()

	const reserveData = decodeFunctionResult({
		abi: MockPool,
		functionName: 'getReserveData',
		data: bytesToHex(callResult.data),
	})

	return reserveData.currentLiquidityRate as bigint
}

const readBalanceInPool = (
	runtime: Runtime<Config>,
	evmCfg: EVMConfig,
	evmClient: EVMClient,
): bigint => {
	const callData = encodeFunctionData({
		abi: MockPool,
		functionName: 'balanceOf',
		args: [evmCfg.protocolSmartWalletAddress as Hex, evmCfg.assetAddress as Hex],
	})

	const callResult = evmClient
		.callContract(runtime, {
			call: encodeCallMsg({
				from: zeroAddress,
				to: evmCfg.poolAddress as Address,
				data: callData,
			}),
			blockNumber: LATEST_BLOCK_NUMBER, // warn: use finalized or safe in prod
		})
		.result()

	const decoded = decodeFunctionResult({
		abi: MockPool,
		functionName: 'balanceOf',
		data: bytesToHex(callResult.data),
	})

	return decoded as bigint
}

const buildPoolForChain = (
	runtime: Runtime<Config>,
	evmCfg: EVMConfig,
): Pool => {
	runtime.log(
		`Reading APY for chain ${evmCfg.chainName} | pool ${evmCfg.poolAddress} | asset ${evmCfg.assetAddress}`,
	)

	const evmClient = getEvmClientForChain(evmCfg)

	const currentLiquidityRate = readCurrentLiquidityRate(runtime, evmCfg, evmClient)
	runtime.log(`Reserve data [${evmCfg.chainName}] currentLiquidityRate: ${currentLiquidityRate}`)

	const balanceInPool = readBalanceInPool(runtime, evmCfg, evmClient)
	runtime.log(`Balance in pool [${evmCfg.chainName}]: ${balanceInPool}`)

	const apy = aprInRAYToAPY(currentLiquidityRate)
	const apr = aprInRAYToAPR(currentLiquidityRate)
	runtime.log(
		`Supply yield [${evmCfg.chainName}] APY%: ${(apy * 100).toFixed(6)}, APR%: ${(apr * 100).toFixed(6)}`,
	)

	return {
		chainName: evmCfg.chainName,
		APR: currentLiquidityRate,
		APY: apy,
		protocolSmartWalletAddress: evmCfg.protocolSmartWalletAddress,
		balance: balanceInPool,
	}
}

const findBestPool = (pools: Pool[], log: (m: string) => void): Pool => {
	let best: Pool | null = null
	for (const p of pools) {
		if (!best || p.APR > best.APR) {
			best = p
		} else if (p.APR === best.APR) {
			log(`Found tie in APY between ${best.chainName} and ${p.chainName}, keeping existing best.`)
		}
	}
	if (!best || best.APR <= 0n) throw new Error('Best APY unset or <= 0')
	return best
}

const getChainSelectorFor = (chainName: string): bigint => {
	const network = getNetwork({
		chainFamily: 'evm',
		chainSelectorName: chainName,
		isTestnet: true,
	})
	if (!network) {
		throw new Error(`Could not find network for chain ${chainName}`)
	}
	return network.chainSelector.selector
}

const shouldRebalance = (
	maxAPR: bigint,
	curAPR: bigint,
	minBpsDelta: number,
): { ok: boolean; diffBps: bigint } => {
	const diff = maxAPR - curAPR
	const diffBps = aprDiffToBps(diff)
	return { ok: diffBps >= BigInt(minBpsDelta), diffBps }
}

const performRebalance = (
	runtime: Runtime<Config>,
	evmCfg: EVMConfig,
	amount: bigint,
	bestChainSelector: bigint,
	bestProtocolSmartWallet: string,
): void => {
	const evmClient = getEvmClientForChain(evmCfg)

	const reportData = encodeAbiParameters(
		parseAbiParameters("address asset, uint256 amount, uint64 destinationChainSelector, address destinationProtocolSmartWallet"),
		[evmCfg.assetAddress as Hex, amount, bestChainSelector, bestProtocolSmartWallet as Hex]
	)

	const reportResponse = runtime
		.report({
			encodedPayload: hexToBase64(reportData),
			encoderName: 'evm',
			signingAlgo: 'ecdsa',
			hashingAlgo: 'keccak256',
		})
		.result()

	const resp = evmClient
		.writeReport(runtime, {
			receiver: evmCfg.protocolSmartWalletAddress,
			report: reportResponse,
			gasConfig: {
				gasLimit: evmCfg.gasLimit,
			},
		})
		.result()

	if (resp.txStatus !== TxStatus.SUCCESS) {
		throw new Error(`Failed to write report: ${resp.errorMessage || resp.txStatus}`)
	}

	const txHash = resp.txHash || new Uint8Array(32)
	runtime.log(
		`Write report transaction succeeded on ${evmCfg.chainName} txHash: ${bytesToHex(txHash)}`,
	)
}

/** ===== Orchestration ===== **/

const doHighestSupplyAPY = (runtime: Runtime<Config>): string => {
	const config = runtime.config

	if (config.evms.length < 2) {
		throw new Error('At least two EVM configurations are required to compare supply APYs')
	}

	runtime.log('Reading supply APYs...')

	// 1) Build pool snapshots for all chains
	const pools: Pool[] = config.evms.map((e) => buildPoolForChain(runtime, e))

	// 2) Find best pool by APR (ties keep first)
	const bestPool = findBestPool(pools, runtime.log)
	runtime.log(`Found best APY: ${bestPool.APY} on chain ${bestPool.chainName}`)

	const bestChainSelector = getChainSelectorFor(bestPool.chainName)

	// Track new balance and amount rebalanced after rebalances
	let newBestBalance = bestPool.balance
	let totalRebalancedAmount = 0n

	// 3) Rebalance from suboptimal pools
	for (const evmCfg of config.evms) {
		if (evmCfg.chainName === bestPool.chainName) continue

		runtime.log(
			`Rebalancing from ${evmCfg.chainName} -> ${bestPool.chainName} (selector ${bestChainSelector})`,
		)

		const evmClient = getEvmClientForChain(evmCfg)
		const balance = readBalanceInPool(runtime, evmCfg, evmClient)

		if (balance === 0n) {
			runtime.log(`No balance to rebalance on chain ${evmCfg.chainName}, skipping.`)
			continue
		}

		const curPool = pools.find((p) => p.chainName === evmCfg.chainName)
		if (!curPool) throw new Error(`Pool info not found for chain ${evmCfg.chainName}`)

		const { ok, diffBps } = shouldRebalance(
			bestPool.APR,
			curPool.APR,
			config.minBPSDeltaForRebalance,
		)

		if (!ok) {
			runtime.log(
				`APY diff below threshold: diff=${diffBps}, min=${config.minBPSDeltaForRebalance} → skipping.`,
			)
			continue
		}

		runtime.log(
			`Rebalancing supply from ${evmCfg.chainName} to ${bestPool.chainName} | balance=${balance}`,
		)

		performRebalance(
			runtime,
			evmCfg,
			balance,
			bestChainSelector,
			bestPool.protocolSmartWalletAddress,
		)

		newBestBalance += balance
		totalRebalancedAmount += balance
	}

	runtime.log(
		`Rebalancing complete | Old balance: ${bestPool.balance} | New balance: ${newBestBalance} | Amount rebalanced: ${totalRebalancedAmount} | Chain: ${bestPool.chainName}`,
	)

	return ''
}

/** ===== Workflow ===== **/

const onCronTrigger = (runtime: Runtime<Config>, payload: CronPayload): string => {
	if (!payload.scheduledExecutionTime) {
		throw new Error('Scheduled execution time is required')
	}
	runtime.log('Running CronTrigger for supply APY rebalance')
	return doHighestSupplyAPY(runtime)
}

const initWorkflow = (config: Config) => {
	const cron = new cre.capabilities.CronCapability()
	return [
		cre.handler(
			cron.trigger({
				schedule: config.schedule,
			}),
			onCronTrigger,
		),
	]
}

export async function main() {
	const runner = await Runner.newRunner<Config>({ configSchema })
	await runner.run(initWorkflow)
}

main()
