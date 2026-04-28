import {
  bytesToHex,
  ConsensusAggregationByFields,
  cre,
  getNetwork,
  median,
  TxStatus,
  type HTTPSendRequester,
  type Runtime,
} from '@chainlink/cre-sdk'
import { type Address, encodeAbiParameters, parseAbiParameters } from 'viem'
import { z } from 'zod'
import {
  SportsMarket,
  type DecodedLog,
  type SettlementRequestedDecoded,
} from '../contracts/evm/ts/generated/SportsMarket'

// ─── Config Schema ──────────────────────────────────────────
export const configSchema = z.object({
  evms: z.array(
    z.object({
      chainSelectorName: z.string(),
      sportsMarketAddress: z.string(),
      gasLimit: z.string().optional(),
    })
  ).min(1),
  // Base URLs for each sports data source. The default ESPN example uses
  // a public endpoint; adapt fetchGameResult() when using API-key providers.
  dataSourceUrls: z.array(z.string()).min(2).max(3),
  // "majority"  — strictly more than half of sources must agree (2-of-2 or 2-of-3)
  // "unanimous" — all sources must agree; any disagreement returns null
  aggregationMode: z.enum(['majority', 'unanimous']).default('majority'),
})
export type Config = z.infer<typeof configSchema>

// Outcome values match SportsMarket.sol Outcome enum
const OUTCOME_HOME_WINS = 1
const OUTCOME_AWAY_WINS = 2
const OUTCOME_DRAW = 3

// ─── Types ───────────────────────────────────────────────────

interface GameResult {
  homeScore: number
  awayScore: number
}

interface FetchConfig {
  url: string
  gameId: string
}

// ─── Helpers ─────────────────────────────────────────────────

const safeJsonStringify = (obj: unknown): string =>
  JSON.stringify(obj, (_, v) => (typeof v === 'bigint' ? v.toString() : v), 2)

export const computeOutcome = (result: GameResult): number => {
  if (result.homeScore > result.awayScore) return OUTCOME_HOME_WINS
  if (result.awayScore > result.homeScore) return OUTCOME_AWAY_WINS
  return OUTCOME_DRAW
}

export const applyAggregation = (
  outcomes: number[],
  mode: 'majority' | 'unanimous',
): number | null => {
  if (outcomes.length === 0) return null

  if (mode === 'unanimous') {
    return outcomes.every((o) => o === outcomes[0]) ? outcomes[0] : null
  }

  // majority: strictly more than half must agree
  const threshold = Math.floor(outcomes.length / 2) + 1
  const counts = new Map<number, number>()
  for (const o of outcomes) counts.set(o, (counts.get(o) ?? 0) + 1)
  for (const [outcome, count] of counts) {
    if (count >= threshold) return outcome
  }
  return null
}

// ─── HTTP Fetch ───────────────────────────────────────────────

/**
 * Fetches a game result from the ESPN per-game scoreboard endpoint.
 *
 * config.url is the base scoreboard URL, e.g.:
 *   https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard
 *
 * The game ID is appended as a path segment: {url}/{gameId}
 * This returns ~20KB of per-game data — well within CRE's 100KB HTTP limit.
 * (The date-based scoreboard returns all games for a day and can exceed 100KB.)
 *
 * Response path: competitions[0].competitors[].{homeAway, score}
 * Scores are returned as strings — parsed to integers before returning.
 *
 * This function runs inside the DON's BFT consensus layer: each node calls
 * the same URL, and ConsensusAggregationByFields produces a verified result
 * before the value is returned to the callback.
 */
export const fetchGameResult = (
  sendRequester: HTTPSendRequester,
  config: FetchConfig,
): GameResult => {
  const response = sendRequester
    .sendRequest({
      url: `${config.url}/${config.gameId}`,
      method: 'GET',
    })
    .result()

  if (response.statusCode !== 200) {
    throw new Error(`ESPN returned HTTP ${response.statusCode}`)
  }

  const body = JSON.parse(Buffer.from(response.body).toString('utf-8'))

  const competitors: Array<{ homeAway: string; score: string }> =
    body.competitions?.[0]?.competitors

  if (!Array.isArray(competitors) || competitors.length < 2) {
    throw new Error(`Unexpected competitor structure for game ${config.gameId}`)
  }

  const home = competitors.find((c) => c.homeAway === 'home')
  const away = competitors.find((c) => c.homeAway === 'away')

  if (!home || !away) {
    throw new Error(`Could not identify home/away for game ${config.gameId}`)
  }

  return {
    homeScore: parseInt(home.score, 10),
    awayScore: parseInt(away.score, 10),
  }
}

// ─── Callback ────────────────────────────────────────────────

export const onSettlementRequested = (
  runtime: Runtime<Config>,
  log: DecodedLog<SettlementRequestedDecoded>,
): string => {
  const evmConfig = runtime.config.evms[0]
  const { gameId, description } = log.data

  runtime.log(`SettlementRequested: gameId=${gameId}, description="${description}"`)

  // 1. Get network and create EVM client
  const network = getNetwork({
    chainFamily: 'evm',
    chainSelectorName: evmConfig.chainSelectorName,
  })
  if (!network) throw new Error(`Network not found: ${evmConfig.chainSelectorName}`)

  const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
  const sportsMarket = new SportsMarket(evmClient, evmConfig.sportsMarketAddress as Address)

  const game = sportsMarket.getGame(runtime, gameId)
  if (game.outcome !== 0 || game.settledAt !== 0n) {
    runtime.log(`Game ${gameId} is already settled; skipping resolution`)
    return safeJsonStringify({
      gameId: gameId.toString(),
      status: 'already_settled',
      outcome: game.outcome,
      settledAt: game.settledAt.toString(),
    })
  }

  // TODO: support per-source game IDs for providers that use different ID schemes.
  // When dataSourceUrls point to different providers (e.g. ESPN + Sportradar), each
  // may use a different game ID for the same match. A structured description field
  // (e.g. JSON encoding an array of IDs by source index) could carry this mapping
  // without requiring a config redeploy.

  // 2. Fetch from each data source (BFT-verified per source via DON consensus)
  const httpClient = new cre.capabilities.HTTPClient()
  const outcomes: number[] = []

  for (const url of runtime.config.dataSourceUrls) {
    runtime.log(`Fetching from: ${url}`)

    const result = httpClient
      .sendRequest(
        runtime,
        fetchGameResult,
        ConsensusAggregationByFields<GameResult>({
          homeScore: median,
          awayScore: median,
        }),
      )({ url, gameId: gameId.toString() })
      .result()

    const outcome = computeOutcome(result)
    runtime.log(
      `Source ${url}: home=${result.homeScore}, away=${result.awayScore}, outcome=${outcome}`,
    )
    outcomes.push(outcome)
  }

  // 3. Apply workflow-level aggregation across sources
  const consensusOutcome = applyAggregation(outcomes, runtime.config.aggregationMode)

  if (consensusOutcome === null) {
    runtime.log(
      `No consensus (${runtime.config.aggregationMode}) across sources [${outcomes.join(', ')}]. Returning without write.`,
    )
    return safeJsonStringify({
      gameId: gameId.toString(),
      status: 'no_consensus',
      outcomes,
      aggregationMode: runtime.config.aggregationMode,
    })
  }

  runtime.log(`Consensus outcome for game ${gameId}: ${consensusOutcome}`)

  // 4. Encode and submit signed report — contract verifies Forwarder signature before accepting
  const reportData = encodeAbiParameters(
    parseAbiParameters('uint256 gameId, uint8 outcome'),
    [gameId, consensusOutcome],
  )

  const resp = sportsMarket.writeReport(runtime, reportData, {
    gasLimit: evmConfig.gasLimit ?? '500000',
  })

  if (resp.txStatus !== TxStatus.SUCCESS) {
    throw new Error(`Settlement TX failed: ${resp.errorMessage ?? resp.txStatus}`)
  }

  const txHash = bytesToHex(resp.txHash ?? new Uint8Array(32))
  runtime.log(`Game ${gameId} settled! Outcome: ${consensusOutcome}, TX: ${txHash}`)

  return safeJsonStringify({
    gameId: gameId.toString(),
    outcome: consensusOutcome,
    sources: outcomes,
    txHash,
  })
}

// ─── Workflow Init ────────────────────────────────────────────

export function initWorkflow(config: Config) {
  const evmConfig = config.evms[0]

  const network = getNetwork({
    chainFamily: 'evm',
    chainSelectorName: evmConfig.chainSelectorName,
  })
  if (!network) throw new Error(`Network not found: ${evmConfig.chainSelectorName}`)

  const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
  const sportsMarket = new SportsMarket(evmClient, evmConfig.sportsMarketAddress as Address)

  return [cre.handler(sportsMarket.logTriggerSettlementRequested(), onSettlementRequested)]
}
