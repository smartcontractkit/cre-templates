import {
  cre,
  getNetwork,
  TxStatus,
  bytesToHex,
  type Runtime,
} from "@chainlink/cre-sdk"
import {
  type Address,
  encodeAbiParameters,
  parseAbiParameters,
  formatUnits,
} from "viem"
import { z } from "zod"
import { PredictionMarket, type DecodedLog, type DisputeRaisedDecoded } from "../contracts/evm/ts/generated/PredictionMarket"
import {
  fetchLatestReport,
  decodeFullReportV3,
  toPriceFeedDecimals,
  PRICE_FEED_DECIMALS,
} from "../data-streams/client"

// ─── Config Schema ──────────────────────────────────────────
export const configSchema = z.object({
  evms: z.array(
    z.object({
      chainSelectorName: z.string(),
      predictionMarketAddress: z.string(),
      gasLimit: z.string().optional(),
    })
  ),
  dataStreams: z.object({
    apiUrl: z.string(),
    feedId: z.string(),
    owner: z.string(),
  }),
})
type Config = z.infer<typeof configSchema>

const ACTION_RESOLVE_DISPUTE = 3

// ─── Helpers ────────────────────────────────────────────────

const safeJsonStringify = (obj: unknown): string =>
  JSON.stringify(obj, (_, v) => (typeof v === "bigint" ? v.toString() : v), 2)

// ─── Callback ───────────────────────────────────────────────
export const onDisputeRaised = (runtime: Runtime<Config>, log: DecodedLog<DisputeRaisedDecoded>): string => {
  const evmConfig = runtime.config.evms[0]
  const dsConfig = runtime.config.dataStreams

  // 1. Extract decoded event data
  const marketId = log.data.marketId
  const disputor = log.data.disputor
  const reason = log.data.reason

  runtime.log(`DisputeRaised: market ${marketId}, disputor ${disputor}, reason: "${reason}"`)

  // 2. Get network and create EVM client
  const network = getNetwork({
    chainFamily: "evm",
    chainSelectorName: evmConfig.chainSelectorName,
  })
  if (!network) throw new Error("Network not found")

  const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)

  const predictionMarket = new PredictionMarket(
    evmClient,
    evmConfig.predictionMarketAddress as Address,
  )

  // 3. Verify market is in Disputed status
  const status = predictionMarket.getMarketStatus(runtime, marketId)
  if (status !== 2) { // 2 = Disputed
    runtime.log(`Market ${marketId}: status is ${status}, not Disputed (2). Skipping.`)
    return `Skipped — market ${marketId} is not in Disputed status`
  }

  // 4. Fetch fresh BTC/USD price from Chainlink Data Streams
  const apiResponse = fetchLatestReport(runtime, dsConfig.apiUrl, dsConfig.feedId)
  const decoded = decodeFullReportV3(apiResponse.report.fullReport as `0x${string}`)

  // Convert from Data Streams 18 decimals to on-chain feed 8 decimals
  const latestAnswer = toPriceFeedDecimals(decoded.price)
  const priceScaled = formatUnits(latestAnswer, PRICE_FEED_DECIMALS)

  runtime.log(`Fresh BTC/USD price: $${priceScaled} (raw: ${latestAnswer}, source: Data Streams)`)

  // 5. Encode the RESOLVE_DISPUTE action
  const innerData = encodeAbiParameters(
    parseAbiParameters("uint256 marketId, int256 newPrice"),
    [marketId, latestAnswer],
  )

  const reportData = encodeAbiParameters(
    parseAbiParameters("uint8 action, bytes data"),
    [ACTION_RESOLVE_DISPUTE, innerData],
  )

  // 6. Write: Resolve dispute via signed report
  const resp = predictionMarket.writeReport(runtime, reportData, {
    gasLimit: evmConfig.gasLimit || "500000",
  })

  if (resp.txStatus !== TxStatus.SUCCESS) {
    throw new Error(`Dispute resolution TX failed: ${resp.errorMessage || resp.txStatus}`)
  }

  if (!resp.txHash) {
    runtime.log("Warning: transaction succeeded but no tx hash returned")
  }
  const txHash = bytesToHex(resp.txHash || new Uint8Array(32))
  runtime.log(`Dispute resolved for market ${marketId}! Fresh price: $${priceScaled}, TX: ${txHash}`)

  return safeJsonStringify({
    marketId: marketId.toString(),
    freshPrice: priceScaled,
    txHash,
  })
}

// ─── Workflow Init ──────────────────────────────────────────
export function initWorkflow(config: Config) {
  const evmConfig = config.evms[0]

  const network = getNetwork({
    chainFamily: "evm",
    chainSelectorName: evmConfig.chainSelectorName,
  })
  if (!network) throw new Error(`Network not found: ${evmConfig.chainSelectorName}`)

  const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
  const predictionMarket = new PredictionMarket(evmClient, evmConfig.predictionMarketAddress as Address)

  return [cre.handler(predictionMarket.logTriggerDisputeRaised(), onDisputeRaised)]
}
