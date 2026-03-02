import {
	bytesToHex,
	ConfidentialHTTPClient,
	cre,
	encodeCallMsg,
	getNetwork,
	LAST_FINALIZED_BLOCK_NUMBER,
	Runner,
	type Runtime,
	type CronPayload,
} from '@chainlink/cre-sdk';
import { encodeFunctionData, decodeFunctionResult, formatUnits, type Address, zeroAddress } from 'viem';
import { z } from 'zod';
import { PriceFeedAggregator } from '../contracts/abi';

// ---------- Config ----------

const configSchema = z.object({
	// 6-field cron; e.g. "0 */10 * * * *" (every 10 minutes)
	schedule: z.string(),
	// e.g. "ethereum-mainnet-arbitrum-1"
	chainName: z.string(),
	// single data feed to monitor
	feed: z.object({
		name: z.string(),    // "ETH/USD"
		address: z.string(), // proxy address
	}),
	// webhook endpoint URL
	webhookUrl: z.string(),
	// "slack" or "telegram"
	notificationType: z.enum(['slack', 'telegram']),
	// Telegram chat ID (only used when notificationType is "telegram")
	telegramChatId: z.string().optional().default(''),
});

type Config = z.infer<typeof configSchema>;

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

const safeJsonStringify = (obj: unknown) =>
	JSON.stringify(obj, (_, v) => (typeof v === 'bigint' ? v.toString() : v), 2);

function formatPrice(raw: bigint, decimals: number): string {
	const scaled = formatUnits(raw, decimals);
	const num = parseFloat(scaled);
	return num.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

// ---------- Read Feed ----------

function readFeed(
	runtime: Runtime<Config>,
	evmClient: InstanceType<typeof cre.capabilities.EVMClient>,
	name: string,
	address: string,
): { decimals: number; latestAnswer: bigint; scaled: string } {
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

	const scaled = formatUnits(latestAnswer, decimals);

	runtime.log(
		`Price feed read | chain=${runtime.config.chainName} feed="${name}" address=${address} decimals=${decimals} latestAnswerRaw=${latestAnswer.toString()} latestAnswerScaled=${scaled}`,
	);

	return { decimals, latestAnswer, scaled };
}

// ---------- Webhook ----------

function buildWebhookBody(config: Config, feedName: string, formattedPrice: string): string {
	if (config.notificationType === 'telegram') {
		return JSON.stringify({
			chat_id: config.telegramChatId,
			text: `*${feedName}*: $${formattedPrice}`,
			parse_mode: 'Markdown',
		});
	}

	// Default: Slack
	return JSON.stringify({
		text: `:chart_with_upwards_trend: *${feedName}*: $${formattedPrice}`,
	});
}

// ---------- Handler ----------

function onCron(runtime: Runtime<Config>, _payload: CronPayload): string {
	const { feed, chainName } = runtime.config;
	const evmClient = getEvmClient(chainName);

	// 1. Read the data feed
	const result = readFeed(runtime, evmClient, feed.name, feed.address);
	const formattedPrice = formatPrice(result.latestAnswer, result.decimals);

	runtime.log(`Formatted price | feed="${feed.name}" price=$${formattedPrice}`);

	// 2. Build webhook payload
	const webhookBody = buildWebhookBody(runtime.config, feed.name, formattedPrice);

	runtime.log(
		`Sending ${runtime.config.notificationType} notification | url=${runtime.config.webhookUrl}`,
	);

	// 3. Send webhook via ConfidentialHTTPClient
	const confHttpClient = new ConfidentialHTTPClient();
	const response = confHttpClient
		.sendRequest(runtime, {
			request: {
				url: runtime.config.webhookUrl,
				method: 'POST',
				bodyString: webhookBody,
				multiHeaders: {
					'Content-Type': { values: ['application/json'] },
				},
			},
			vaultDonSecrets: [],
		})
		.result();

	runtime.log(`Webhook response | statusCode=${response.statusCode}`);

	// 4. Return summary
	return safeJsonStringify({
		feed: feed.name,
		address: feed.address,
		decimals: result.decimals,
		latestAnswerRaw: result.latestAnswer,
		scaled: result.scaled,
		formattedPrice,
		notificationType: runtime.config.notificationType,
		webhookStatusCode: response.statusCode,
	});
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
