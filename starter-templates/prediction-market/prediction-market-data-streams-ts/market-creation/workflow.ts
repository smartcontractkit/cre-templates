import {
  cre,
  getNetwork,
  TxStatus,
  bytesToHex,
  type Runtime,
  type CronPayload,
} from "@chainlink/cre-sdk"
import { type Address, encodeAbiParameters, parseAbiParameters, parseUnits } from "viem"
import { z } from "zod"
import { PredictionMarket } from "../contracts/evm/ts/generated/PredictionMarket"

// ─── Config Schema ──────────────────────────────────────────
export const configSchema = z.object({
  schedule: z.string(),
  evms: z.array(
    z.object({
      chainSelectorName: z.string(),
      predictionMarketAddress: z.string(),
      gasLimit: z.string().optional(),
    })
  ),
  marketDefaults: z.object({
    question: z.string(),
    strikePriceUsd: z.number(),
    durationSeconds: z.number(),
  }),
})
type Config = z.infer<typeof configSchema>

const ACTION_CREATE = 1
const PRICE_FEED_DECIMALS = 8 as const // Strike prices use 8 decimals (matching Chainlink USD feeds)

// ─── Helpers ────────────────────────────────────────────────

const safeJsonStringify = (obj: unknown): string =>
  JSON.stringify(obj, (_, v) => (typeof v === "bigint" ? v.toString() : v), 2)

// ─── Callback ───────────────────────────────────────────────
export const onCronTrigger = (runtime: Runtime<Config>, _payload: CronPayload): string => {
  const evmConfig = runtime.config.evms[0]
  const defaults = runtime.config.marketDefaults

  // 1. Get network and create EVM client
  const network = getNetwork({
    chainFamily: "evm",
    chainSelectorName: evmConfig.chainSelectorName,
  })
  if (!network) throw new Error(`Network not found: ${evmConfig.chainSelectorName}`)

  const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
  const predictionMarket = new PredictionMarket(
    evmClient,
    evmConfig.predictionMarketAddress as Address,
  )

  // 2. Read current state — check next market ID (to log it)
  const nextMarketId = predictionMarket.getNextMarketId(runtime)
  runtime.log(`Next market ID will be: ${nextMarketId}`)

  // 3. Build market parameters
  const now = runtime.now()
  const expirationTime = BigInt(Math.floor(now.getTime() / 1000)) + BigInt(defaults.durationSeconds)
  const strikePriceScaled = parseUnits(defaults.strikePriceUsd.toString(), PRICE_FEED_DECIMALS)

  // Format question with expiration date
  const expirationDate = new Date(Number(expirationTime) * 1000).toISOString().split("T")[0]
  const question = defaults.question.replace("{expirationDate}", expirationDate)

  runtime.log(`Creating market: "${question}" | Strike: $${defaults.strikePriceUsd} | Expires: ${expirationDate}`)

  // 4. Encode the CREATE action
  const innerData = encodeAbiParameters(
    parseAbiParameters("string question, uint256 strikePrice, uint256 expirationTime"),
    [question, strikePriceScaled, expirationTime],
  )

  const reportData = encodeAbiParameters(
    parseAbiParameters("uint8 action, bytes data"),
    [ACTION_CREATE, innerData],
  )

  // 5. Write: Create market via signed report
  const resp = predictionMarket.writeReport(runtime, reportData, {
    gasLimit: evmConfig.gasLimit || "500000",
  })

  if (resp.txStatus !== TxStatus.SUCCESS) {
    throw new Error(`Create market TX failed: ${resp.errorMessage || resp.txStatus}`)
  }

  if (!resp.txHash) {
    runtime.log("Warning: transaction succeeded but no tx hash returned")
  }
  const txHash = bytesToHex(resp.txHash || new Uint8Array(32))
  runtime.log(`Market created! ID: ${nextMarketId}, TX: ${txHash}`)

  return safeJsonStringify({
    marketId: nextMarketId.toString(),
    question,
    strikePrice: strikePriceScaled.toString(),
    expirationTime: expirationTime.toString(),
    txHash,
  })
}

// ─── Workflow Init ──────────────────────────────────────────
export function initWorkflow(config: Config) {
  const cronTrigger = new cre.capabilities.CronCapability()
  return [
    cre.handler(
      cronTrigger.trigger({ schedule: config.schedule }),
      onCronTrigger,
    ),
  ]
}
