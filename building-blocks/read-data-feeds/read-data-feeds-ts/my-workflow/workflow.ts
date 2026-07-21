import {
	cre,
	getNetwork,
	type Runtime,
	type CronPayload,
} from '@chainlink/cre-sdk';
import { formatUnits, type Address } from 'viem';
import { z } from 'zod';
import { PriceFeedAggregator } from '../contracts/evm/ts/generated/PriceFeedAggregator';

// ---------- Config ----------

export const configSchema = z.object({
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
	const aggregator = new PriceFeedAggregator(evmClient, address as Address);
	const decimals = aggregator.decimals(runtime);
	const latestAnswer = aggregator.latestAnswer(runtime);
	const scaled = formatUnits(latestAnswer, decimals);

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

export function onCron(runtime: Runtime<Config>, _payload: CronPayload): string {
	const evmClient = getEvmClient(runtime.config.chainName);

	const results: PriceResult[] = runtime.config.feeds.map((f) =>
		readFeed(runtime, evmClient, f.name, f.address),
	);

	return safeJsonStringify(results);
}

// ---------- Init ----------

export function initWorkflow(config: Config) {
	const cron = new cre.capabilities.CronCapability();
	return [
		cre.handler(
			cron.trigger({ schedule: config.schedule }),
			onCron,
		),
	];
}
