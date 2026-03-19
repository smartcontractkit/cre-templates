import {
  cre,
  getNetwork,
  TxStatus,
  bytesToHex,
  type Runtime,
  type CronPayload,
} from "@chainlink/cre-sdk"
import { type Address, encodeAbiParameters, parseAbiParameters, formatUnits } from "viem"
import { z } from "zod"
import { PredictionMarket } from "../contracts/evm/ts/generated/PredictionMarket"
import { PriceFeedAggregator } from "../contracts/evm/ts/generated/PriceFeedAggregator"

// ─── Config Schema ──────────────────────────────────────────
export const configSchema = z.object({
  schedule: z.string(),
  evms: z.array(
    z.object({
      chainSelectorName: z.string(),
      predictionMarketAddress: z.string(),
      priceFeedAddress: z.string(),
      gasLimit: z.string().optional(),
    })
  ),
  marketIdsToCheck: z.array(z.number()),
})
type Config = z.infer<typeof configSchema>

const ACTION_RESOLVE = 2

// ─── Helpers ────────────────────────────────────────────────

const safeJsonStringify = (obj: unknown): string =>
  JSON.stringify(obj, (_, v) => (typeof v === "bigint" ? v.toString() : v), 2)

// ─── Callback ───────────────────────────────────────────────
export const onCronTrigger = (runtime: Runtime<Config>, _payload: CronPayload): string => {
  const evmConfig = runtime.config.evms[0]

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
  const priceFeed = new PriceFeedAggregator(
    evmClient,
    evmConfig.priceFeedAddress as Address,
  )

  const resolved: string[] = []

  // 2. Check each market ID for resolvability
  for (const marketId of runtime.config.marketIdsToCheck) {
    const isResolvable = predictionMarket.isResolvable(runtime, BigInt(marketId))

    if (!isResolvable) {
      runtime.log(`Market ${marketId}: not resolvable (either not expired or already resolved)`)
      continue
    }

    runtime.log(`Market ${marketId}: resolvable — reading price feed...`)

    // 3. Read BTC/USD price from Chainlink Data Feed on-chain
    const latestAnswer = priceFeed.latestAnswer(runtime)
    const decimals = priceFeed.decimals(runtime)
    const priceScaled = formatUnits(latestAnswer, decimals)

    runtime.log(`BTC/USD price: $${priceScaled} (raw: ${latestAnswer}, decimals: ${decimals})`)

    // 4. Encode the RESOLVE action with the price
    const innerData = encodeAbiParameters(
      parseAbiParameters("uint256 marketId, int256 price"),
      [BigInt(marketId), latestAnswer],
    )

    const reportData = encodeAbiParameters(
      parseAbiParameters("uint8 action, bytes data"),
      [ACTION_RESOLVE, innerData],
    )

    // 5. Write: Resolve market via signed report
    const resp = predictionMarket.writeReport(runtime, reportData, {
      gasLimit: evmConfig.gasLimit || "500000",
    })

    if (resp.txStatus !== TxStatus.SUCCESS) {
      runtime.log(`Market ${marketId}: resolve TX failed — ${resp.errorMessage || resp.txStatus}`)
      continue
    }

    if (!resp.txHash) {
      runtime.log(`Market ${marketId}: warning — transaction succeeded but no tx hash returned`)
    }
    const txHash = bytesToHex(resp.txHash || new Uint8Array(32))
    runtime.log(`Market ${marketId}: resolved! Price: $${priceScaled}, TX: ${txHash}`)

    resolved.push(`${marketId}`)
  }

  if (resolved.length === 0) {
    return "No markets were resolvable"
  }

  return `Resolved markets: ${resolved.join(", ")}`
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
