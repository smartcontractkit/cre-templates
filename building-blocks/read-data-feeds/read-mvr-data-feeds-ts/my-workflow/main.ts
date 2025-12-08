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
import {
	encodeFunctionData,
	decodeFunctionResult,
	decodeAbiParameters,
	type Address,
	zeroAddress,
} from 'viem';
import { z } from 'zod';
import { BundleAggregatorProxy } from '../contracts/abi';

// ---------- Config ----------

const configSchema = z.object({
	// e.g. "0 */10 * * * *" (every 10 minutes, at second 0)
	schedule: z.string(),
	// e.g. "ethereum-mainnet-base-1"
	chainName: z.string(),
	// list of MVR feeds (BundleAggregatorProxy contracts)
	feeds: z.array(
		z.object({
			name: z.string(),    // "S&P Global SSA EURC"
			address: z.string(), // proxy address
		}),
	),
});

type Config = z.infer<typeof configSchema>;

type DecodedBundle = {
	lastModifiedDateTimeRaw: string;
	lastModifiedDateTimeRfc3339: string;
	securityId: string;
	securityName: string;
	ssaRaw: string;
	ssaScaled: string;
	ssaDesc: string;
	ssaDecimal: number;
};

type BundleResult = {
	name: string;
	address: string;
	bundle: DecodedBundle;
	bundleDecimals: number[];
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

// ---------- MVR Reader ----------
//
// Bundle layout (known for this MVR feed):
//   0: LastModifiedDateTime : uint256
//   1: SecurityID           : string
//   2: SecurityName         : string
//   3: SSA                  : uint256
//   4: SSADesc              : string
//
// bundleDecimals is uint8[], where index 3 corresponds to SSA.
//

function readMvrBundle(
	runtime: Runtime<Config>,
	evmClient: InstanceType<typeof cre.capabilities.EVMClient>,
	name: string,
	address: string,
): BundleResult {
	// 1) bundleDecimals()
	const decCallData = encodeFunctionData({
		abi: BundleAggregatorProxy,
		functionName: 'bundleDecimals',
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

	const decodedDecimals = decodeFunctionResult({
		abi: BundleAggregatorProxy,
		functionName: 'bundleDecimals',
		data: bytesToHex(decResp.data),
	}) as readonly (bigint | number)[];

	const bundleDecimals = decodedDecimals.map((d) => Number(d));

	// 2) latestBundle()
	const bundleCallData = encodeFunctionData({
		abi: BundleAggregatorProxy,
		functionName: 'latestBundle',
	});

	const bundleResp = evmClient
		.callContract(runtime, {
			call: encodeCallMsg({
				from: zeroAddress,
				to: address as Address,
				data: bundleCallData,
			}),
			blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
		})
		.result();

	const latestBundleBytes = decodeFunctionResult({
		abi: BundleAggregatorProxy,
		functionName: 'latestBundle',
		data: bytesToHex(bundleResp.data),
	}) as `0x${string}`;

	// 3) Decode the bundle bytes according to the known layout
	const [lastModified, securityId, securityName, ssa, ssaDesc] = decodeAbiParameters(
		[
			{ type: 'uint256', name: 'lastModifiedDateTime' },
			{ type: 'string', name: 'securityId' },
			{ type: 'string', name: 'securityName' },
			{ type: 'uint256', name: 'ssa' },
			{ type: 'string', name: 'ssaDesc' },
		],
		latestBundleBytes,
	) as [bigint, string, string, bigint, string];

	const lastModifiedRaw = lastModified.toString();

	let lastModifiedIso = '';
	const secNumber = Number(lastModified);
	if (Number.isSafeInteger(secNumber)) {
		lastModifiedIso = new Date(secNumber * 1000).toISOString();
	}

	const ssaRaw = ssa.toString();
	const ssaDecimal = bundleDecimals[3] ?? 0;
	const ssaScaled = formatScaled(ssa, ssaDecimal);

	runtime.log(
		`MVR bundle read | ` +
		`chain=${runtime.config.chainName} ` +
		`feed="${name}" ` +
		`address=${address} ` +
		`lastModifiedDateTimeRaw=${lastModifiedRaw} ` +
		`lastModifiedDateTimeRFC3339=${lastModifiedIso} ` +
		`securityId=${securityId} ` +
		`securityName=${securityName} ` +
		`ssaRaw=${ssaRaw} ` +
		`ssaScaled=${ssaScaled} ` +
		`ssaDesc=${ssaDesc} ` +
		`ssaDecimal=${ssaDecimal}`,
	);

	const bundle: DecodedBundle = {
		lastModifiedDateTimeRaw: lastModifiedRaw,
		lastModifiedDateTimeRfc3339: lastModifiedIso,
		securityId,
		securityName,
		ssaRaw,
		ssaScaled,
		ssaDesc,
		ssaDecimal,
	};

	return {
		name,
		address,
		bundle,
		bundleDecimals,
	};
}

// ---------- Handlers ----------

function onCron(runtime: Runtime<Config>, _payload: CronPayload): string {
	const evmClient = getEvmClient(runtime.config.chainName);

	const results: BundleResult[] = runtime.config.feeds.map((f) =>
		readMvrBundle(runtime, evmClient, f.name, f.address),
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
