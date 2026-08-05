import {
  CronCapability,
  HTTPClient,
  NITRO_REGIONS,
  Runner,
  handlerInTee,
  type TeeRuntime,
  type Workflow,
} from "@chainlink/cre-sdk";

type Holdings = Record<string, number>;
type Prices = Record<string, number>;

type Policy = {
  target_allocations: Record<string, number>;
  drift_threshold_pct: number;
  max_trade_usd: number;
  reserve_floor_usdc: number;
  max_slippage_bps: number;
  preferred_venues: string[];
  order_sequence_preference: "model-order" | "sells-first" | "buys-first";
};

type SecretsConfig = {
  exchange_api_key_id: string;
  openai_api_key_id: string;
  target_allocation_btc_pct_secret_id: string;
  target_allocation_eth_pct_secret_id: string;
  target_allocation_usdc_pct_secret_id: string;
  drift_threshold_pct_secret_id: string;
  max_trade_usd_secret_id: string;
  reserve_floor_usdc_secret_id: string;
  max_slippage_bps_secret_id: string;
  preferred_venues_secret_id: string;
  order_sequence_preference_secret_id: string;
};

export type Config = {
  schedule: string;
  mock_base_url: string;
  openai_url: string;
  openai_model: string;
  secrets_ids: SecretsConfig;
};

type PortfolioState = {
  holdings: Holdings;
  prices_usd: Prices;
  usdc_reserve: number;
  cash_balance_usd: number;
  volatility_index: number;
};

type Allocation = {
  symbol: string;
  usdValue: number;
  currentWeight: number;
  targetWeight: number;
  driftPct: number;
  deltaUsdToTarget: number;
};

type RebalanceTrade = {
  symbol: string;
  side: "buy" | "sell";
  notionalUsd: number;
  maxSlippageBps?: number;
  preferredVenues?: string[];
};

type LlmDecision = {
  shouldRebalance: boolean;
  reasoning: string;
  trades: RebalanceTrade[];
};

type ExecutableTrade = RebalanceTrade & {
  venue: string;
  chunkIndex: number;
  chunkCount: number;
};

type MarketSnapshot = {
  portfolio: PortfolioState;
  prices: Prices;
  volatilityIndex: number;
};

type PortfolioMetrics = {
  totalUsd: number;
  maxDriftPct: number;
  reserveBufferUsd: number;
  reserveCoverageRatio: number;
  stablecoinDepthUsd: number;
  averageVolatility: number;
};

type ExecutionPlan = {
  onchainTrades: ExecutableTrade[];
  offchainTrades: ExecutableTrade[];
};

const JSON_HEADERS = { "Content-Type": "application/json" };

const asObject = (value: unknown): Record<string, unknown> => {
  if (!value || typeof value !== "object" || Array.isArray(value)) {
    return {};
  }
  return value as Record<string, unknown>;
};

const asNumber = (value: unknown, fallback = 0): number => {
  if (typeof value === "number" && Number.isFinite(value)) {
    return value;
  }
  if (typeof value === "string" && value.trim() !== "") {
    const parsed = Number(value);
    if (Number.isFinite(parsed)) {
      return parsed;
    }
  }
  return fallback;
};

const asString = (value: unknown, fallback = ""): string => {
  return typeof value === "string" ? value : fallback;
};

const asStringArray = (value: unknown): string[] => {
  if (!Array.isArray(value)) {
    return [];
  }
  return value.filter((item): item is string => typeof item === "string" && item.trim() !== "");
};

const parseRequiredSecretNumber = (secretValue: string, secretId: string): number => {
  const parsed = asNumber(secretValue, Number.NaN);
  if (!Number.isFinite(parsed)) {
    throw new Error(`secret ${secretId} must be a finite number`);
  }
  return parsed;
};

const parseVenueListSecret = (secretValue: string): string[] => {
  return secretValue
    .split(",")
    .map((entry) => entry.trim())
    .filter((entry) => entry !== "");
};

const parseOrderSequencePreference = (
  value: string,
): Policy["order_sequence_preference"] => {
  if (value === "sells-first" || value === "buys-first") {
    return value;
  }
  return "model-order";
};

export const calculateAllocations = (
  holdings: Holdings,
  prices: Prices,
  targets: Record<string, number>,
): { allocations: Allocation[]; totalUsd: number; maxDriftPct: number } => {
  const symbols = [...new Set([...Object.keys(targets), ...Object.keys(holdings)])];
  const usdBySymbol: Record<string, number> = {};

  for (const symbol of symbols) {
    const qty = asNumber(holdings[symbol], 0);
    const px = asNumber(prices[symbol], 0);
    usdBySymbol[symbol] = qty * px;
  }

  const totalUsd = Object.values(usdBySymbol).reduce((sum, value) => sum + value, 0);
  const safeTotal = totalUsd > 0 ? totalUsd : 1;

  const allocations = symbols.map((symbol) => {
    const targetWeight = asNumber(targets[symbol], 0);
    const usdValue = usdBySymbol[symbol] ?? 0;
    const currentWeight = usdValue / safeTotal;
    const driftPct = Math.abs(currentWeight - targetWeight) * 100;

    return {
      symbol,
      usdValue,
      currentWeight,
      targetWeight,
      driftPct,
      deltaUsdToTarget: targetWeight * totalUsd - usdValue,
    } satisfies Allocation;
  });

  const maxDriftPct = allocations.reduce((max, next) => Math.max(max, next.driftPct), 0);
  return { allocations, totalUsd, maxDriftPct };
};

const decodeBody = (raw: Uint8Array): string => new TextDecoder().decode(raw);

const parseJson = (body: string): Record<string, unknown> => {
  try {
    return asObject(JSON.parse(body));
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    throw new Error(`invalid json response: ${message}; body=${body}`);
  }
};

const parseHoldings = (payload: unknown): Holdings => {
  const obj = asObject(payload);
  const out: Holdings = {};
  for (const [symbol, value] of Object.entries(obj)) {
    out[symbol] = asNumber(value, 0);
  }
  return out;
};

const parsePrices = (payload: unknown): Prices => {
  const obj = asObject(payload);
  const out: Prices = {};
  for (const [symbol, value] of Object.entries(obj)) {
    out[symbol] = asNumber(value, 0);
  }
  return out;
};

const postJson = (
  runtime: TeeRuntime<Config>,
  client: HTTPClient,
  url: string,
  body: Record<string, unknown>,
  headers: Record<string, string>,
): Record<string, unknown> => {
  const bodyBytes = new TextEncoder().encode(JSON.stringify(body));
  const encodedBody = Buffer.from(bodyBytes).toString("base64");

  const response = client
    .sendRequest(runtime, {
      url,
      method: "POST",
      body: encodedBody,
      headers,
    })
    .result();

  const raw = decodeBody(response.body);
  if (response.statusCode >= 400) {
    throw new Error(`request failed status=${response.statusCode} body=${raw}`);
  }
  return parseJson(raw);
};

const getJson = (
  runtime: TeeRuntime<Config>,
  client: HTTPClient,
  url: string,
  headers: Record<string, string>,
): Record<string, unknown> => {
  const response = client
    .sendRequest(runtime, {
      url,
      method: "GET",
      headers,
    })
    .result();

  const raw = decodeBody(response.body);
  if (response.statusCode >= 400) {
    throw new Error(`request failed status=${response.statusCode} body=${raw}`);
  }
  return parseJson(raw);
};

const extractOpenAiText = (payload: Record<string, unknown>): string => {
  const direct = asString(payload.output_text);
  if (direct) {
    return direct;
  }

  const outputItems = Array.isArray(payload.output) ? payload.output : [];
  for (const item of outputItems) {
    const content = Array.isArray(asObject(item).content) ? (asObject(item).content as unknown[]) : [];
    for (const part of content) {
      const text = asString(asObject(part).text);
      if (text) {
        return text;
      }
    }
  }
  throw new Error("openai response did not contain output_text");
};

const parseLlmDecision = (rawText: string): LlmDecision => {
  const parsed = parseJson(rawText);
  const shouldRebalance = Boolean(parsed.shouldRebalance);
  const reasoning = asString(parsed.reasoning, "No reasoning provided by model");
  const tradesRaw = Array.isArray(parsed.trades) ? parsed.trades : [];

  const trades = tradesRaw
    .map((item) => {
      const row = asObject(item);
      const side = asString(row.side).toLowerCase() === "sell" ? "sell" : "buy";
      return {
        symbol: asString(row.symbol).toUpperCase(),
        side,
        notionalUsd: asNumber(row.notionalUsd, 0),
        maxSlippageBps: asNumber(row.maxSlippageBps, 0),
        preferredVenues: asStringArray(row.preferredVenues),
      } satisfies RebalanceTrade;
    })
    .filter((trade) => trade.symbol !== "" && trade.notionalUsd > 0);

  return {
    shouldRebalance,
    reasoning,
    trades,
  };
};

const chooseVenue = (trade: RebalanceTrade, preferredVenues: string[]): string => {
  if (trade.preferredVenues && trade.preferredVenues.length > 0) {
    const overlap = trade.preferredVenues.find((venue) => preferredVenues.includes(venue));
    if (overlap) {
      return overlap;
    }
    return trade.preferredVenues[0];
  }
  if (preferredVenues.length > 0) {
    return preferredVenues[0];
  }
  return "default-venue";
};

const computePortfolioMetrics = (
  snapshot: MarketSnapshot,
  allocations: Allocation[],
  policy: Policy,
  totalUsd: number,
  maxDriftPct: number,
): PortfolioMetrics => {
  const usdcPrice = asNumber(snapshot.prices.USDC, 1);
  const usdcHoldingsUsd = asNumber(snapshot.portfolio.holdings.USDC, 0) * usdcPrice;
  const reserveBufferUsd = snapshot.portfolio.usdc_reserve - policy.reserve_floor_usdc;
  const reserveCoverageRatio =
    policy.reserve_floor_usdc > 0
      ? snapshot.portfolio.usdc_reserve / policy.reserve_floor_usdc
      : Number.POSITIVE_INFINITY;
  const averageVolatility = (snapshot.portfolio.volatility_index + snapshot.volatilityIndex) / 2;

  return {
    totalUsd,
    maxDriftPct,
    reserveBufferUsd,
    reserveCoverageRatio,
    stablecoinDepthUsd: usdcHoldingsUsd + snapshot.portfolio.cash_balance_usd,
    averageVolatility,
  };
};

const isTradableRiskAsset = (symbol: string): boolean => symbol !== "USDC";

export const buildTargetAlignedTrades = (
  allocations: Allocation[],
): RebalanceTrade[] => {
  return allocations
    .filter((allocation) => isTradableRiskAsset(allocation.symbol))
    .filter((allocation) => Math.abs(allocation.deltaUsdToTarget) >= 1)
    .map((allocation) => ({
      symbol: allocation.symbol,
      side: allocation.deltaUsdToTarget > 0 ? "buy" : "sell",
      notionalUsd: Math.abs(allocation.deltaUsdToTarget),
    } satisfies RebalanceTrade));
};

const indexTradesBySymbol = (trades: RebalanceTrade[]): Map<string, RebalanceTrade> => {
  return new Map(trades.map((trade) => [trade.symbol, trade]));
};

const reconcileDecisionWithTargets = (
  decision: LlmDecision,
  targetTrades: RebalanceTrade[],
  policy: Policy,
  metrics: PortfolioMetrics,
): LlmDecision => {
  if (metrics.maxDriftPct < policy.drift_threshold_pct || targetTrades.length === 0) {
    return {
      shouldRebalance: false,
      reasoning: `Drift ${metrics.maxDriftPct.toFixed(3)} is below threshold ${policy.drift_threshold_pct}`,
      trades: [],
    };
  }

  const llmTradesBySymbol = indexTradesBySymbol(decision.trades);
  const totalPlannedSellUsd = targetTrades.reduce(
    (sum, trade) => sum + (trade.side === "sell" ? trade.notionalUsd : 0),
    0,
  );
  let remainingBuyCapacityUsd = Math.max(metrics.reserveBufferUsd + totalPlannedSellUsd, 0);

  const reconciledTrades = targetTrades.flatMap((targetTrade) => {
    const llmTrade = llmTradesBySymbol.get(targetTrade.symbol);
    const preferredVenues =
      llmTrade?.preferredVenues && llmTrade.preferredVenues.length > 0
        ? llmTrade.preferredVenues
        : policy.preferred_venues;
    const slippage = Math.min(
      Math.max(asNumber(llmTrade?.maxSlippageBps, policy.max_slippage_bps), 0),
      policy.max_slippage_bps,
    );

    if (llmTrade && llmTrade.side !== targetTrade.side) {
      return [];
    }

    let cappedNotional = targetTrade.notionalUsd;
    if (targetTrade.side === "buy") {
      cappedNotional = Math.min(cappedNotional, remainingBuyCapacityUsd);
      remainingBuyCapacityUsd -= cappedNotional;
    }

    if (cappedNotional < 1) {
      return [];
    }

    return [
      {
        symbol: targetTrade.symbol,
        side: targetTrade.side,
        notionalUsd: cappedNotional,
        maxSlippageBps: slippage,
        preferredVenues,
      } satisfies RebalanceTrade,
    ];
  });

  return {
    shouldRebalance: reconciledTrades.length > 0,
    reasoning: decision.reasoning,
    trades: reconciledTrades,
  };
};

const classifyExecutionVenue = (venue: string): "onchain" | "offchain" => {
  const normalized = venue.trim().toLowerCase();
  if (
    normalized === "onchain" ||
    normalized.includes("uniswap") ||
    normalized.includes("1inch") ||
    normalized.includes("dex")
  ) {
    return "onchain";
  }
  return "offchain";
};

export const buildExecutionPlan = (trades: ExecutableTrade[]): ExecutionPlan => {
  const onchainTrades: ExecutableTrade[] = [];
  const offchainTrades: ExecutableTrade[] = [];

  for (const trade of trades) {
    if (classifyExecutionVenue(trade.venue) === "onchain") {
      onchainTrades.push(trade);
    } else {
      offchainTrades.push(trade);
    }
  }

  return { onchainTrades, offchainTrades };
};

export const enforcePolicyAndChunk = (
  decision: LlmDecision,
  policy: Policy,
  usdcReserve: number,
): ExecutableTrade[] => {
  if (!decision.shouldRebalance || decision.trades.length === 0) {
    return [];
  }

  const maxTrade = policy.max_trade_usd;
  const maxSlippage = policy.max_slippage_bps;
  let projectedReserve = usdcReserve;
  const chunks: ExecutableTrade[] = [];

  const orderedTrades = [...decision.trades];
  if (policy.order_sequence_preference === "sells-first") {
    orderedTrades.sort((a, b) => (a.side === b.side ? 0 : a.side === "sell" ? -1 : 1));
  } else if (policy.order_sequence_preference === "buys-first") {
    orderedTrades.sort((a, b) => (a.side === b.side ? 0 : a.side === "buy" ? -1 : 1));
  }

  for (const trade of orderedTrades) {
    const slippage = trade.maxSlippageBps ?? maxSlippage;
    if (slippage > maxSlippage) {
      throw new Error(`trade ${trade.symbol} slippage ${slippage} exceeds max ${maxSlippage}`);
    }

    if (trade.side === "buy") {
      projectedReserve -= trade.notionalUsd;
    } else {
      projectedReserve += trade.notionalUsd;
    }

    if (projectedReserve < policy.reserve_floor_usdc) {
      throw new Error(
        `trade ${trade.symbol} breaches reserve floor: projected ${projectedReserve.toFixed(2)} < floor ${policy.reserve_floor_usdc}`,
      );
    }

    const chunkCount = Math.ceil(trade.notionalUsd / maxTrade);
    let remainder = trade.notionalUsd;
    for (let idx = 0; idx < chunkCount; idx++) {
      const chunkSize = Math.min(remainder, maxTrade);
      remainder -= chunkSize;
      chunks.push({
        ...trade,
        notionalUsd: chunkSize,
        maxSlippageBps: slippage,
        venue: chooseVenue(trade, policy.preferred_venues),
        chunkIndex: idx + 1,
        chunkCount,
      });
    }
  }

  return chunks;
};

const collectMarketSnapshot = (
  runtime: TeeRuntime<Config>,
  client: HTTPClient,
  baseUrl: string,
  exchangeApiKey: string,
): MarketSnapshot => {
  const authHeaders = {
    ...JSON_HEADERS,
    "x-exchange-api-key": exchangeApiKey,
  };

  const portfolioResponse = getJson(runtime, client, `${baseUrl}/portfolio-state`, authHeaders);
  const pricesResponse = getJson(runtime, client, `${baseUrl}/prices`, authHeaders);
  const volatilityResponse = getJson(runtime, client, `${baseUrl}/volatility`, authHeaders);

  const portfolioObj = asObject(portfolioResponse);
  const portfolio: PortfolioState = {
    holdings: parseHoldings(portfolioObj.holdings),
    prices_usd: parsePrices(portfolioObj.prices_usd),
    usdc_reserve: asNumber(portfolioObj.usdc_reserve),
    cash_balance_usd: asNumber(portfolioObj.cash_balance_usd),
    volatility_index: asNumber(portfolioObj.volatility_index),
  };

  return {
    portfolio,
    prices: parsePrices(pricesResponse.prices),
    volatilityIndex: asNumber(volatilityResponse.volatility_index),
  };
};

const createOpenAiPrompt = (
  allocations: Allocation[],
  metrics: PortfolioMetrics,
  policy: Policy,
  snapshot: MarketSnapshot,
): string => {
  const targetTrades = buildTargetAlignedTrades(allocations);

  return JSON.stringify(
    {
      objective:
        "Return strict JSON with shouldRebalance:boolean, reasoning:string, trades:[{symbol, side, notionalUsd, maxSlippageBps, preferredVenues}] to restore portfolio target weights while respecting policy. Use only the provided targetAlignedTrades symbols and sides, and keep notionalUsd within those bounds.",
      policy,
      metrics,
      allocations,
      targetAlignedTrades: targetTrades,
      snapshot,
    },
    null,
    2,
  );
};

export const onCronTrigger = async (runtime: TeeRuntime<Config>): Promise<string> => {
  const { mock_base_url, openai_url, openai_model, secrets_ids } = runtime.config;

  const exchangeApiKey = runtime.getSecret({ id: secrets_ids.exchange_api_key_id }).result().value;
  const openAiApiKey = runtime.getSecret({ id: secrets_ids.openai_api_key_id }).result().value;
  const targetAllocationBtcPct = parseRequiredSecretNumber(
    runtime.getSecret({ id: secrets_ids.target_allocation_btc_pct_secret_id }).result().value,
    secrets_ids.target_allocation_btc_pct_secret_id,
  );
  const targetAllocationEthPct = parseRequiredSecretNumber(
    runtime.getSecret({ id: secrets_ids.target_allocation_eth_pct_secret_id }).result().value,
    secrets_ids.target_allocation_eth_pct_secret_id,
  );
  const targetAllocationUsdcPct = parseRequiredSecretNumber(
    runtime.getSecret({ id: secrets_ids.target_allocation_usdc_pct_secret_id }).result().value,
    secrets_ids.target_allocation_usdc_pct_secret_id,
  );
  const driftThresholdPct = parseRequiredSecretNumber(
    runtime.getSecret({ id: secrets_ids.drift_threshold_pct_secret_id }).result().value,
    secrets_ids.drift_threshold_pct_secret_id,
  );
  const maxTradeUsd = parseRequiredSecretNumber(
    runtime.getSecret({ id: secrets_ids.max_trade_usd_secret_id }).result().value,
    secrets_ids.max_trade_usd_secret_id,
  );
  const reserveFloorUsdc = parseRequiredSecretNumber(
    runtime.getSecret({ id: secrets_ids.reserve_floor_usdc_secret_id }).result().value,
    secrets_ids.reserve_floor_usdc_secret_id,
  );
  const maxSlippageBps = parseRequiredSecretNumber(
    runtime.getSecret({ id: secrets_ids.max_slippage_bps_secret_id }).result().value,
    secrets_ids.max_slippage_bps_secret_id,
  );
  const preferredVenues = parseVenueListSecret(
    runtime.getSecret({ id: secrets_ids.preferred_venues_secret_id }).result().value,
  );
  const orderSequencePreference = parseOrderSequencePreference(
    runtime.getSecret({ id: secrets_ids.order_sequence_preference_secret_id }).result().value,
  );
  runtime.log("rebalance-getsecret-ok");

  const client = new HTTPClient();
  const snapshot = collectMarketSnapshot(runtime, client, mock_base_url, exchangeApiKey);
  const policy: Policy = {
    target_allocations: {
      BTC: targetAllocationBtcPct,
      ETH: targetAllocationEthPct,
      USDC: targetAllocationUsdcPct,
    },
    drift_threshold_pct: driftThresholdPct,
    max_trade_usd: maxTradeUsd,
    reserve_floor_usdc: reserveFloorUsdc,
    max_slippage_bps: maxSlippageBps,
    preferred_venues: preferredVenues,
    order_sequence_preference: orderSequencePreference,
  };

  const { allocations, maxDriftPct, totalUsd } = calculateAllocations(
    snapshot.portfolio.holdings,
    snapshot.prices,
    policy.target_allocations,
  );

  const metrics = computePortfolioMetrics(snapshot, allocations, policy, totalUsd, maxDriftPct);

  const prompt = createOpenAiPrompt(allocations, metrics, policy, snapshot);

  const llmResponse = postJson(
    runtime,
    client,
    openai_url,
    {
      model: openai_model,
      input: [
        {
          role: "system",
          content:
            "You are a portfolio rebalancing policy engine. Emit strict JSON only with key names exactly as requested.",
        },
        {
          role: "user",
          content: prompt,
        },
      ],
    },
    {
      ...JSON_HEADERS,
      Authorization: `Bearer ${openAiApiKey}`,
    },
  );

  const decision = parseLlmDecision(extractOpenAiText(llmResponse));
  const targetTrades = buildTargetAlignedTrades(allocations);
  const reconciledDecision = reconcileDecisionWithTargets(decision, targetTrades, policy, metrics);
  const chunks = enforcePolicyAndChunk(
    reconciledDecision,
    policy,
    snapshot.portfolio.usdc_reserve,
  );

  if (maxDriftPct < policy.drift_threshold_pct || chunks.length === 0) {
    runtime.log(
      `rebalance-skip drift=${maxDriftPct.toFixed(3)} threshold=${policy.drift_threshold_pct} reason=${reconciledDecision.reasoning}`,
    );
    return "NOOP";
  }

  const executionPlan = buildExecutionPlan(chunks);

  const executionResponse = postJson(
    runtime,
    client,
    `${mock_base_url}/execute-rebalance`,
    {
      totalPortfolioUsd: totalUsd,
      maxDriftPct,
      reserveBufferUsd: metrics.reserveBufferUsd,
      stablecoinDepthUsd: metrics.stablecoinDepthUsd,
      volatilityIndex: metrics.averageVolatility,
      reasoning: reconciledDecision.reasoning,
      trades: chunks,
      onchain_trades: executionPlan.onchainTrades,
      offchain_trades: executionPlan.offchainTrades,
    },
    {
      ...JSON_HEADERS,
      "x-exchange-api-key": exchangeApiKey,
    },
  );

  runtime.log(
    `rebalance-executed trade_count=${chunks.length} execution_id=${asString(executionResponse.execution_id, "unknown")}`,
  );
  return JSON.stringify({
    status: "EXECUTED",
    tradeCount: chunks.length,
    onchainTradeCount: executionPlan.onchainTrades.length,
    offchainTradeCount: executionPlan.offchainTrades.length,
    maxDriftPct,
    executionId: asString(executionResponse.execution_id, "unknown"),
  });
};

export const initWorkflow = (config: Config): Workflow<Config> => {
  if (!config.schedule || !config.mock_base_url || !config.openai_url || !config.openai_model) {
    throw new Error(
      "config requires schedule, mock_base_url, openai_url, and openai_model",
    );
  }

  if (
    !config.secrets_ids?.exchange_api_key_id ||
    !config.secrets_ids?.openai_api_key_id ||
    !config.secrets_ids?.target_allocation_btc_pct_secret_id ||
    !config.secrets_ids?.target_allocation_eth_pct_secret_id ||
    !config.secrets_ids?.target_allocation_usdc_pct_secret_id ||
    !config.secrets_ids?.drift_threshold_pct_secret_id ||
    !config.secrets_ids?.max_trade_usd_secret_id ||
    !config.secrets_ids?.reserve_floor_usdc_secret_id ||
    !config.secrets_ids?.max_slippage_bps_secret_id ||
    !config.secrets_ids?.preferred_venues_secret_id ||
    !config.secrets_ids?.order_sequence_preference_secret_id
  ) {
    throw new Error("config requires secrets_ids fields");
  }

  const cron = new CronCapability();

  return [
    handlerInTee(
      cron.trigger({ schedule: config.schedule }),
      onCronTrigger,
      [{ tee: "nitro", regions: [NITRO_REGIONS[0]] }],
    ),
  ];
};

export async function main() {
  const runner = await Runner.newRunner<Config>();
  await runner.run(initWorkflow);
}
