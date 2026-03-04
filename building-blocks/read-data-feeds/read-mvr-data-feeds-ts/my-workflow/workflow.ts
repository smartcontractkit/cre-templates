import {
	cre,
	getNetwork,
	type Runtime,
	type CronPayload,
} from '@chainlink/cre-sdk';
import { decodeAbiParameters, formatUnits, type Address } from 'viem';
import { z } from 'zod';
import { BundleAggregatorProxy } from '../contracts/evm/ts/generated/BundleAggregatorProxy';

// ---------- Config ----------

export const configSchema = z.object({
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
	const proxy = new BundleAggregatorProxy(evmClient, address as Address);
	const bundleDecimals = [...proxy.bundleDecimals(runtime)];
	const latestBundleBytes = proxy.latestBundle(runtime);

	// Decode the bundle bytes according to the known layout
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
	const ssaScaled = formatUnits(ssa, ssaDecimal);

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

export function onCron(runtime: Runtime<Config>, _payload: CronPayload): string {
	const evmClient = getEvmClient(runtime.config.chainName);

	const results: BundleResult[] = runtime.config.feeds.map((f) =>
		readMvrBundle(runtime, evmClient, f.name, f.address),
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
