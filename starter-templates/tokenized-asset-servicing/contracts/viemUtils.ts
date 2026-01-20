import { createPublicClient, createWalletClient, http } from 'viem'
import { privateKeyToAccount } from 'viem/accounts'
import { sepolia } from 'viem/chains'
import * as dotenv from 'dotenv'
import * as yaml from 'js-yaml'
import * as fs from 'fs'
import { join } from 'path'

export function loadTokenizedAssetPlatformAddr() {
  const parentDir = join(process.cwd(), '..')
  const configPath = join(parentDir, 'asset-log-trigger-workflow/config.json')

  if (!fs.existsSync(configPath)) {
    throw new Error(`Cannot find config file at path ${configPath}`)
  }

  const content = fs.readFileSync(configPath, 'utf8')
  let config

  try {
    config = JSON.parse(content)
  } catch (err) {
    throw new Error(`failed to parse config.json ${(err as Error).message}`)
  }

  const assetAddress = config?.evms?.[0]?.assetAddress

  if (!assetAddress) {
    throw new Error(
      'evms[0].assetAddress does not exist in config.json\n' +
      'please check the assetAddress exist evm in config.json'
    )
  }

  if (!assetAddress.startsWith('0x') || assetAddress.length !== 42) {
    throw new Error(`invalid assetAddress: ${assetAddress}`)
  }

  return assetAddress as `0x${string}`
}

/**
 * Load private key from ../.env
 * Expects CRE_ETH_PRIVATE_KEY in .env file
 */
export function loadPrivateKeyAndAccount() {
  const parentDir = join(process.cwd(), '..')
  const envPath = join(parentDir, '.env')

  if (!fs.existsSync(envPath)) {
    throw new Error(`Cannot find .env file at: ${envPath}`)
  }

  dotenv.config({ path: envPath })

  const privateKey = process.env.CRE_ETH_PRIVATE_KEY

  if (!privateKey) {
    throw new Error('CRE_ETH_PRIVATE_KEY not found in .env file')
  }

  const account = privateKeyToAccount(
    privateKey.startsWith('0x')
      ? (privateKey as `0x${string}`)
      : (`0x${privateKey}` as `0x${string}`)
  )

  return { account, privateKey }
}

/**
 * Load Sepolia RPC URL from ../project.yaml
 * Looks for local-simulation.rpcs where chain-name === 'ethereum-testnet-sepolia'
 */
export function loadSepoliaRpcUrl() {
  const parentDir = join(process.cwd(), '..')
  const yamlPath = join(parentDir, 'project.yaml')

  if (!fs.existsSync(yamlPath)) {
    throw new Error(`Cannot find project.yaml at: ${yamlPath}`)
  }

  const yamlContent = fs.readFileSync(yamlPath, 'utf8')
  const config = yaml.load(yamlContent) as any

  const rpcs = config?.['local-simulation']?.rpcs || []

  let rpcUrl: string | undefined

  for (const rpc of rpcs) {
    if (rpc['chain-name'] === 'ethereum-testnet-sepolia') {
      rpcUrl = rpc.url
      break
    }
  }

  if (!rpcUrl || typeof rpcUrl !== 'string') {
    throw new Error(
      'No valid Sepolia RPC URL found in project.yaml\n' +
      'Please check local-simulation.rpcs for chain-name: ethereum-testnet-sepolia'
    )
  }

  return rpcUrl
}

/**
 * Create viem clients for Sepolia network
 * Returns both walletClient and publicClient
 */
export function createSepoliaClients() {
  const { account } = loadPrivateKeyAndAccount()
  const rpcUrl = loadSepoliaRpcUrl()

  const walletClient = createWalletClient({
    account,
    chain: sepolia,
    transport: http(rpcUrl),
  })

  const publicClient = createPublicClient({
    chain: sepolia,
    transport: http(rpcUrl),
  })

  return {
    walletClient,
    publicClient,
    account,
    rpcUrl,
  }
}