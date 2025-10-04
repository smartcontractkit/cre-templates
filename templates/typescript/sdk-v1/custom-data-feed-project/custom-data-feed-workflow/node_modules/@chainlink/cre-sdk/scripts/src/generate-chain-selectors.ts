#!/usr/bin/env bun

import { mkdirSync, readFileSync, rmSync, writeFileSync } from 'node:fs'
import { join } from 'node:path'
import { parse } from 'yaml'
import { z } from 'zod'

// Zod schemas for validation
const SelectorValueSchema = z.union([z.bigint(), z.string(), z.number()])

const EvmSelectorSchema = z.object({
	selector: SelectorValueSchema,
	name: z.string().optional(),
})

const NonEvmSelectorSchema = z.object({
	selector: SelectorValueSchema,
	name: z.string().optional(),
})

const EvmSelectorsSchema = z.object({
	selectors: z.record(z.string(), EvmSelectorSchema),
})

const NonEvmSelectorsSchema = z.object({
	selectors: z.record(z.string(), NonEvmSelectorSchema),
})

type ChainFamily = 'evm' | 'solana' | 'aptos' | 'sui' | 'ton' | 'tron'
type NetworkType = 'mainnet' | 'testnet'

interface NetworkInfo {
	chainId: string
	chainSelector: {
		name: string
		selector: bigint
	}
	chainFamily: ChainFamily
	networkType: NetworkType
}

interface ChainSelectorConfig {
	family: ChainFamily
	filename: string
	schema: typeof EvmSelectorsSchema | typeof NonEvmSelectorsSchema
}

const CHAIN_CONFIGS: ChainSelectorConfig[] = [
	{
		family: 'evm',
		filename: 'selectors.yml',
		schema: EvmSelectorsSchema,
	},
	{
		family: 'solana',
		filename: 'selectors_solana.yml',
		schema: NonEvmSelectorsSchema,
	},
	{
		family: 'aptos',
		filename: 'selectors_aptos.yml',
		schema: NonEvmSelectorsSchema,
	},
	{
		family: 'sui',
		filename: 'selectors_sui.yml',
		schema: NonEvmSelectorsSchema,
	},
	{
		family: 'ton',
		filename: 'selectors_ton.yml',
		schema: NonEvmSelectorsSchema,
	},
	{
		family: 'tron',
		filename: 'selectors_tron.yml',
		schema: NonEvmSelectorsSchema,
	},
]

const readYamlFile = (filename: string): string => {
	// Look for chain-selectors in workspace root node_modules (turborepo setup)
	const workspaceRoot = join(process.cwd(), '..', '..')
	const chainSelectorsPath = join(workspaceRoot, 'node_modules', 'chain-selectors', filename)
	try {
		return readFileSync(chainSelectorsPath, 'utf-8')
	} catch (error) {
		throw new Error(`Failed to read ${filename}: ${error}`)
	}
}

const parseChainSelectors = (): NetworkInfo[] => {
	const allNetworks: NetworkInfo[] = []

	for (const config of CHAIN_CONFIGS) {
		try {
			console.log(`📂 Reading ${config.family} selectors from ${config.filename}...`)
			const yamlContent = readYamlFile(config.filename)
			const parsed = parse(yamlContent, { intAsBigInt: true })
			const validated = config.schema.parse(parsed)

			for (const [chainId, selectorData] of Object.entries(validated.selectors)) {
				const typedSelectorData = selectorData as {
					selector: bigint | string | number
					name?: string
				}
				// Skip entries without names (they might be test or incomplete entries)
				if (!typedSelectorData.name) {
					console.log(`⚠️  Skipping ${config.family} chain ${chainId} - no name provided`)
					continue
				}

				const selectorBigInt: bigint =
					typeof typedSelectorData.selector === 'bigint'
						? typedSelectorData.selector
						: BigInt(typedSelectorData.selector)

				allNetworks.push({
					chainId,
					chainSelector: {
						name: typedSelectorData.name,
						selector: selectorBigInt,
					},
					chainFamily: config.family,
					networkType: detectNetworkType(typedSelectorData.name),
				})
			}

			console.log(`✅ Parsed ${Object.keys(validated.selectors).length} ${config.family} networks`)
		} catch (error) {
			console.error(`❌ Failed to process ${config.family}:`, error)
		}
	}

	return allNetworks
}

const detectNetworkType = (networkName: string): NetworkType => {
	const name = networkName.toLowerCase()

	// Check for testnet indicators
	const testnetIndicators = [
		'testnet',
		'test',
		'devnet',
		'dev',
		'localnet',
		'local',
		'goerli',
		'sepolia',
		'holesky',
		'mumbai',
		'amoy',
		'fuji',
		'chapel',
		'shasta',
		'nile',
		'kairos',
		'shibuya',
		'pangoro',
		'cardona',
		'puppynet',
		'alfajores',
		'moonbase',
		'tatara',
	]

	// Check if any testnet indicator is present
	for (const indicator of testnetIndicators) {
		if (name.includes(indicator)) {
			return 'testnet'
		}
	}

	return 'mainnet'
}

const sanitizeFilename = (name: string): string => {
	// Replace invalid characters with dots and ensure it's a valid filename
	return name.replace(/[^a-zA-Z0-9.-]/g, '.').replace(/\.+/g, '.')
}

const generateNetworkFiles = (networks: NetworkInfo[]): void => {
	const baseDir = 'src/generated/chain-selectors'

	// Clean up existing directory
	try {
		rmSync(baseDir, { recursive: true, force: true })
	} catch {
		// Directory might not exist, that's ok
	}

	// Create base directory
	mkdirSync(baseDir, { recursive: true })

	// Group networks by network type and then by family
	const networksByTypeAndFamily = networks.reduce(
		(acc, network) => {
			if (!acc[network.networkType]) {
				acc[network.networkType] = {} as Record<ChainFamily, NetworkInfo[]>
			}
			if (!acc[network.networkType][network.chainFamily]) {
				acc[network.networkType][network.chainFamily] = []
			}
			acc[network.networkType][network.chainFamily].push(network)
			return acc
		},
		{} as Record<NetworkType, Record<ChainFamily, NetworkInfo[]>>,
	)

	// Create directories and files for each network type and family
	for (const [networkType, familiesInType] of Object.entries(networksByTypeAndFamily)) {
		for (const [family, familyNetworks] of Object.entries(familiesInType)) {
			const familyDir = join(baseDir, networkType, family)
			mkdirSync(familyDir, { recursive: true })

			for (const network of familyNetworks) {
				const filename = sanitizeFilename(network.chainSelector.name)
				const filepath = join(familyDir, `${filename}.ts`)

				const fileContent = `// This file is auto-generated. Do not edit manually.
// Generated from: https://github.com/smartcontractkit/chain-selectors

import type { NetworkInfo } from "@cre/sdk/utils/chain-selectors/types";

const network: NetworkInfo = {
	chainId: "${network.chainId}",
	chainSelector: {
		name: "${network.chainSelector.name}",
		selector: ${network.chainSelector.selector}n,
	},
	chainFamily: "${network.chainFamily}",
	networkType: "${network.networkType}",
} as const;

export default network;
`

				writeFileSync(filepath, fileContent)
			}

			console.log(`📁 Created ${familyNetworks.length} ${networkType}/${family} network files`)
		}
	}
}

const generateAllNetworksFile = (networks: NetworkInfo[]): void => {
	// Group networks by type and family for exports
	const mainnetNetworks = networks.filter((n) => n.networkType === 'mainnet')
	const testnetNetworks = networks.filter((n) => n.networkType === 'testnet')

	const mainnetByFamily = mainnetNetworks.reduce(
		(acc, network) => {
			if (!acc[network.chainFamily]) {
				acc[network.chainFamily] = []
			}
			acc[network.chainFamily].push(network)
			return acc
		},
		{} as Record<ChainFamily, NetworkInfo[]>,
	)

	const testnetByFamily = testnetNetworks.reduce(
		(acc, network) => {
			if (!acc[network.chainFamily]) {
				acc[network.chainFamily] = []
			}
			acc[network.chainFamily].push(network)
			return acc
		},
		{} as Record<ChainFamily, NetworkInfo[]>,
	)

	// Generate import statements and variable names
	const generateImportName = (network: NetworkInfo): string => {
		const sanitized = sanitizeFilename(network.chainSelector.name)
			.replace(/\./g, '_')
			.replace(/-/g, '_')
		return `${network.networkType}_${network.chainFamily}_${sanitized}`
	}

	const allImports = networks.map((network) => ({
		importName: generateImportName(network),
		importPath: `./chain-selectors/${network.networkType}/${
			network.chainFamily
		}/${sanitizeFilename(network.chainSelector.name)}`,
		network,
	}))

	const content = `// This file is auto-generated. Do not edit manually.
// Generated from: https://github.com/smartcontractkit/chain-selectors

import type { NetworkInfo, ChainFamily } from "@cre/sdk/utils/chain-selectors/types";

// Import all individual network files
${allImports
	.map(({ importName, importPath }) => `import ${importName} from "${importPath}";`)
	.join('\n')}

export const allNetworks: NetworkInfo[] = [
${allImports.map(({ importName }) => `	${importName},`).join('\n')}
] as const;

export const mainnet = {
${Object.entries(mainnetByFamily)
	.map(([family]) => {
		const familyImports = allImports.filter(
			({ network }) => network.chainFamily === family && network.networkType === 'mainnet',
		)
		return `	${family}: [
${familyImports.map(({ importName }) => `		${importName},`).join('\n')}
	] as const,`
	})
	.join('\n')}
} as const;

export const testnet = {
${Object.entries(testnetByFamily)
	.map(([family]) => {
		const familyImports = allImports.filter(
			({ network }) => network.chainFamily === family && network.networkType === 'testnet',
		)
		return `	${family}: [
${familyImports.map(({ importName }) => `		${importName},`).join('\n')}
	] as const,`
	})
	.join('\n')}
} as const;

// Optimized Maps for fast lookups by chain selector (bigint keys)
export const mainnetBySelector = new Map<bigint, NetworkInfo>([
${allImports
	.filter(({ network }) => network.networkType === 'mainnet')
	.map(({ importName, network }) => `	[${network.chainSelector.selector}n, ${importName}],`)
	.join('\n')}
]);

export const testnetBySelector = new Map<bigint, NetworkInfo>([
${allImports
	.filter(({ network }) => network.networkType === 'testnet')
	.map(({ importName, network }) => `	[${network.chainSelector.selector}n, ${importName}],`)
	.join('\n')}
]);

// Optimized Maps for fast lookups by chain selector name
export const mainnetByName = new Map<string, NetworkInfo>([
${allImports
	.filter(({ network }) => network.networkType === 'mainnet')
	.map(({ importName, network }) => `	["${network.chainSelector.name}", ${importName}],`)
	.join('\n')}
]);

export const testnetByName = new Map<string, NetworkInfo>([
${allImports
	.filter(({ network }) => network.networkType === 'testnet')
	.map(({ importName, network }) => `	["${network.chainSelector.name}", ${importName}],`)
	.join('\n')}
]);

// Maps by family and network type for chain selector lookups
export const mainnetBySelectorByFamily = {
${Object.entries(mainnetByFamily)
	.map(([family]) => {
		const familyImports = allImports.filter(
			({ network }) => network.chainFamily === family && network.networkType === 'mainnet',
		)
		return `	${family}: new Map<bigint, NetworkInfo>([
${familyImports
	.map(({ importName, network }) => `		[${network.chainSelector.selector}n, ${importName}],`)
	.join('\n')}
	]),`
	})
	.join('\n')}
} as const;

export const testnetBySelectorByFamily = {
${Object.entries(testnetByFamily)
	.map(([family]) => {
		const familyImports = allImports.filter(
			({ network }) => network.chainFamily === family && network.networkType === 'testnet',
		)
		return `	${family}: new Map<bigint, NetworkInfo>([
${familyImports
	.map(({ importName, network }) => `		[${network.chainSelector.selector}n, ${importName}],`)
	.join('\n')}
	]),`
	})
	.join('\n')}
} as const;

// Maps by family and network type for name lookups
export const mainnetByNameByFamily = {
${Object.entries(mainnetByFamily)
	.map(([family]) => {
		const familyImports = allImports.filter(
			({ network }) => network.chainFamily === family && network.networkType === 'mainnet',
		)
		return `	${family}: new Map<string, NetworkInfo>([
${familyImports
	.map(({ importName, network }) => `		["${network.chainSelector.name}", ${importName}],`)
	.join('\n')}
	]),`
	})
	.join('\n')}
} as const;

export const testnetByNameByFamily = {
${Object.entries(testnetByFamily)
	.map(([family]) => {
		const familyImports = allImports.filter(
			({ network }) => network.chainFamily === family && network.networkType === 'testnet',
		)
		return `	${family}: new Map<string, NetworkInfo>([
${familyImports
	.map(({ importName, network }) => `		["${network.chainSelector.name}", ${importName}],`)
	.join('\n')}
	]),`
	})
	.join('\n')}
} as const;
`

	writeFileSync('src/generated/networks.ts', content)
	console.log('📄 Created networks array file with mainnet/testnet grouping')
}

export const main = () => {
	try {
		console.log('🚀 Starting chain selectors generation...')

		const networks = parseChainSelectors()
		console.log(`📊 Total networks processed: ${networks.length}`)

		// Generate individual network files
		generateNetworkFiles(networks)

		// Generate main networks array
		generateAllNetworksFile(networks)

		// Generate summary stats
		const stats = networks.reduce(
			(acc, network) => {
				acc[network.chainFamily] = (acc[network.chainFamily] || 0) + 1
				return acc
			},
			{} as Record<ChainFamily, number>,
		)

		console.log('📈 Summary by family:')
		Object.entries(stats).forEach(([family, count]) => {
			console.log(`  ${family}: ${count} networks`)
		})
	} catch (error) {
		console.error('💥 Generation failed:', error)
		process.exit(1)
	}
}
