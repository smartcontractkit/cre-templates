import {
	bytesToHex,
	cre,
	encodeCallMsg,
	getNetwork,
	LAST_FINALIZED_BLOCK_NUMBER,
	Runner,
	type Runtime,
	type CronPayload,
} from '@chainlink/cre-sdk';
import { encodeFunctionData, decodeFunctionResult, type Address, zeroAddress } from 'viem';
import { z } from 'zod';
import { PriceFeedAggregator } from '../contracts/abi';

// ---------- Config ----------

const configSchema = z.object({
	// e.g. "0 */10 * * * *" (every 10 minutes, at second 0)
	schedule: z.string(),
	// e.g. "ethereum-mainnet-arbitrum-1"
	chainName: z.string(),
	// list of feeds (BTC/USD, ETH/USD, ...)
	feeds: z.array(
		z.object({
			name: z.string(),    // "BTC/USD"
			address: z.string(), // proxy address
		}),
	),
});

type Config = z.infer<typeof configSchema>;

type PriceResult = {
	name: string;
	address: string;
	decimals: number;
	latestAnswerRaw: string;
	scaled: string;
};

// ---------- Helpers ----------

function getEvmClient(chainName: string) {
	const net = getNetwork({
		chainFamily: 'evm',
		chainSelectorName: chainName,
		isTestnet: false,
	});
	if (!net) throw new Error(`Network not found for chain name: ${chainName}`);
	return new cre.capabilities.EVMClient(net.chainSelector.selector);
}

function formatScaled(raw: bigint, decimals: number): string {
	if (decimals === 0) return raw.toString();
	const s = raw.toString();
	if (s.length <= decimals) {
		return `0.${s.padStart(decimals, '0')}`;
	}
	const i = s.length - decimals;
	return `${s.slice(0, i)}.${s.slice(i)}`;
}

// Safely stringify BigInt
const safeJsonStringify = (obj: unknown) =>
	JSON.stringify(obj, (_, v) => (typeof v === 'bigint' ? v.toString() : v), 2);

// ---------- Reader ----------

function readFeed(
	runtime: Runtime<Config>,
	evmClient: InstanceType<typeof cre.capabilities.EVMClient>,
	name: string,
	address: string,
): PriceResult {
	// decimals()
	const decCallData = encodeFunctionData({
		abi: PriceFeedAggregator,
		functionName: 'decimals',
	});

	const decResp = evmClient
		.callContract(runtime, {
			call: encodeCallMsg({
				from: zeroAddress,
				to: address as Address,
				data: decCallData,
			}),
			blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
		})
		.result();

	const decimals = decodeFunctionResult({
		abi: PriceFeedAggregator,
		functionName: 'decimals',
		data: bytesToHex(decResp.data),
	}) as number;

	// latestAnswer()
	const ansCallData = encodeFunctionData({
		abi: PriceFeedAggregator,
		functionName: 'latestAnswer',
	});

	const ansResp = evmClient
		.callContract(runtime, {
			call: encodeCallMsg({
				from: zeroAddress,
				to: address as Address,
				data: ansCallData,
			}),
			blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
		})
		.result();

	const latestAnswer = decodeFunctionResult({
		abi: PriceFeedAggregator,
		functionName: 'latestAnswer',
		data: bytesToHex(ansResp.data),
	}) as bigint;

	const scaled = formatScaled(latestAnswer, decimals);

	runtime.log(
		`Price feed read | chain=${runtime.config.chainName} feed="${name}" address=${address} decimals=${decimals} latestAnswerRaw=${latestAnswer.toString()} latestAnswerScaled=${scaled}`,
	);

	return {
		name,
		address,
		decimals,
		latestAnswerRaw: latestAnswer.toString(),
		scaled,
	};
}

// ---------- Handlers ----------

function onCron(runtime: Runtime<Config>, _payload: CronPayload): string {
	const evmClient = getEvmClient(runtime.config.chainName);

	const results: PriceResult[] = runtime.config.feeds.map((f) =>
		readFeed(runtime, evmClient, f.name, f.address),
	);

	// Return JSON
	return safeJsonStringify(results);
}

// ---------- Init ----------

function initWorkflow(config: Config) {
	const cron = new cre.capabilities.CronCapability();
	return [
		cre.handler(
			cron.trigger({ schedule: config.schedule }),
			onCron,
		),
	];
}

export async function main() {
	const runner = await Runner.newRunner<Config>({ configSchema });
	await runner.run(initWorkflow);
}

main();
