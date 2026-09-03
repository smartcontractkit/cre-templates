import {
  CronCapability,
  HTTPClient,
  NITRO_REGIONS,
  Runner,
  handlerInTee,
  type TeeRuntime,
  type Workflow,
} from "@chainlink/cre-sdk";

type Policy = {
  max_reserve_deployment_usdc: number;
  min_reserve_balance_usdc: number;
  liquidation_warning_action_threshold: number;
  minimum_health_factor: number;
  target_health_factor: number;
  max_collateral_allocation_pct: number;
  max_partial_debt_repayment_pct: number;
  execution_sequence_preference: "collateral-first" | "debt-first" | "balanced";
  preferred_venues: string[];
};

type SecretsConfig = {
  exchange_api_key_id: string;
  openai_api_key_id: string;
  liquidation_warning_action_threshold_secret_id: string;
  minimum_health_factor_secret_id: string;
  target_health_factor_secret_id: string;
  maximum_stablecoin_reserve_deployment_secret_id: string;
  minimum_stablecoin_reserve_balance_secret_id: string;
  maximum_collateral_allocation_secret_id: string;
  maximum_partial_debt_repayment_secret_id: string;
  preferred_venues_secret_id: string;
};

export type Config = {
  schedule: string;
  mock_base_url: string;
  openai_url: string;
  openai_model: string;
  defensive_action_sequencing_preference: string;
  secrets_ids: SecretsConfig;
};

type RiskState = {
  collateral_asset_symbol: string;
  borrowed_asset_symbol: string;
  collateral_asset_price_usd: number;
  borrowed_asset_price_usd: number;
  collateral_balance_usd: number;
  outstanding_debt_usd: number;
  collateral_health_factor: number;
  liquidation_proximity_pct: number;
  loan_to_value_pct: number;
  liquidation_threshold_pct: number;
  volatility_index: number;
  usdc_reserve: number;
  cash_balance_usd: number;
};

type DefensiveAction = {
  type:
    | "add_collateral"
    | "bridge_and_add_collateral"
    | "swap_reserve_to_collateral"
    | "repay_with_reserves"
    | "swap_reserve_to_borrowed_and_repay"
    | "partial_debt_repayment"
    | "full_debt_repayment";
  amountUsd?: number;
  repayPct?: number;
  fromChain?: string;
  toChain?: string;
  venue?: string;
};

type LiquidationDecision = {
  shouldDefend: boolean;
  reasoning: string;
  actions: DefensiveAction[];
};

type ExecutableAction = DefensiveAction & {
  amountUsd: number;
  repayPct: number;
  venue: string;
};

const ACTION_TYPES: DefensiveAction["type"][] = [
  "add_collateral",
  "bridge_and_add_collateral",
  "swap_reserve_to_collateral",
  "repay_with_reserves",
  "swap_reserve_to_borrowed_and_repay",
  "partial_debt_repayment",
  "full_debt_repayment",
];

const isActionType = (value: string): value is DefensiveAction["type"] => {
  return ACTION_TYPES.includes(value as DefensiveAction["type"]);
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

export const computeRiskScore = (risk: RiskState, policy: Policy): number => {
  const proximityRisk = Math.max(0, policy.liquidation_warning_action_threshold - risk.liquidation_proximity_pct) * 5;
  const ltvBufferRisk = Math.max(0, risk.loan_to_value_pct - (risk.liquidation_threshold_pct - 5)) * 2;
  const healthRisk = Math.max(0, policy.minimum_health_factor - risk.collateral_health_factor) * 100;
  const volatilityRisk = risk.volatility_index * 25;

  return proximityRisk + ltvBufferRisk + healthRisk + volatilityRisk;
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

const parseExecutionSequencePreference = (
  value: unknown,
  fallback: Policy["execution_sequence_preference"] = "collateral-first",
): Policy["execution_sequence_preference"] => {
  const sequence = asString(value, fallback);
  if (sequence === "debt-first" || sequence === "balanced") {
    return sequence;
  }
  return "collateral-first";
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

const parseRiskState = (payload: unknown): RiskState => {
  const row = asObject(payload);
  return {
    collateral_asset_symbol: asString(row.collateral_asset_symbol, "ETH"),
    borrowed_asset_symbol: asString(row.borrowed_asset_symbol, "USDC"),
    collateral_asset_price_usd: asNumber(row.collateral_asset_price_usd),
    borrowed_asset_price_usd: asNumber(row.borrowed_asset_price_usd, 1),
    collateral_balance_usd: asNumber(row.collateral_balance_usd),
    outstanding_debt_usd: asNumber(row.outstanding_debt_usd),
    collateral_health_factor: asNumber(row.collateral_health_factor),
    liquidation_proximity_pct: asNumber(row.liquidation_proximity_pct),
    loan_to_value_pct: asNumber(row.loan_to_value_pct),
    liquidation_threshold_pct: asNumber(row.liquidation_threshold_pct),
    volatility_index: asNumber(row.volatility_index),
    usdc_reserve: asNumber(row.usdc_reserve),
    cash_balance_usd: asNumber(row.cash_balance_usd),
  };
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

const parseLlmDecision = (rawText: string): LiquidationDecision => {
  const parsed = parseJson(rawText);
  const shouldDefend = Boolean(parsed.shouldDefend);
  const reasoning = asString(parsed.reasoning, "No reasoning provided by model");
  const actionsRaw = Array.isArray(parsed.actions) ? parsed.actions : [];

  const actions: DefensiveAction[] = [];
  for (const item of actionsRaw) {
    const row = asObject(item);
    const type = asString(row.type);
    if (!isActionType(type)) {
      continue;
    }

    actions.push({
      type,
      amountUsd: asNumber(row.amountUsd, 0),
      repayPct: asNumber(row.repayPct, 0),
      fromChain: asString(row.fromChain),
      toChain: asString(row.toChain),
      venue: asString(row.venue),
    });
  }

  return {
    shouldDefend,
    reasoning,
    actions,
  };
};

const chooseVenue = (action: DefensiveAction, preferredVenues: string[]): string => {
  if (action.venue) {
    return action.venue;
  }
  if (preferredVenues.length > 0) {
    return preferredVenues[0];
  }
  return "default-venue";
};

const orderActions = (actions: ExecutableAction[], preference: Policy["execution_sequence_preference"]): ExecutableAction[] => {
  const priorities: Record<Policy["execution_sequence_preference"], Record<DefensiveAction["type"], number>> = {
    "collateral-first": {
      add_collateral: 1,
      bridge_and_add_collateral: 2,
      swap_reserve_to_collateral: 3,
      repay_with_reserves: 4,
      swap_reserve_to_borrowed_and_repay: 5,
      partial_debt_repayment: 6,
      full_debt_repayment: 7,
    },
    "debt-first": {
      repay_with_reserves: 1,
      swap_reserve_to_borrowed_and_repay: 2,
      partial_debt_repayment: 3,
      full_debt_repayment: 4,
      add_collateral: 5,
      bridge_and_add_collateral: 6,
      swap_reserve_to_collateral: 7,
    },
    balanced: {
      add_collateral: 4,
      bridge_and_add_collateral: 2,
      swap_reserve_to_collateral: 6,
      repay_with_reserves: 3,
      swap_reserve_to_borrowed_and_repay: 5,
      partial_debt_repayment: 1,
      full_debt_repayment: 7,
    },
  };

  return [...actions].sort((a, b) => priorities[preference][a.type] - priorities[preference][b.type]);
};

export const enforcePolicy = (
  decision: LiquidationDecision,
  policy: Policy,
  risk: RiskState,
): ExecutableAction[] => {
  if (!decision.shouldDefend || decision.actions.length === 0) {
    return [];
  }

  let projectedReserve = risk.usdc_reserve;
  const executable: ExecutableAction[] = [];

  for (const action of decision.actions) {
    const amount = Math.max(0, action.amountUsd ?? 0);
    const repayPctRaw = Math.max(0, action.repayPct ?? 0);
    const cappedRepayPct = Math.min(repayPctRaw, policy.max_partial_debt_repayment_pct);

    if (
      action.type === "add_collateral" ||
      action.type === "bridge_and_add_collateral" ||
      action.type === "swap_reserve_to_collateral"
    ) {
      const capped = Math.min(amount, policy.max_reserve_deployment_usdc);
      projectedReserve -= capped;

      if (projectedReserve < policy.min_reserve_balance_usdc) {
        throw new Error(
          `action ${action.type} breaches reserve floor: projected ${projectedReserve.toFixed(2)} < floor ${policy.min_reserve_balance_usdc}`,
        );
      }

      executable.push({
        ...action,
        amountUsd: capped,
        repayPct: 0,
        venue: chooseVenue(action, policy.preferred_venues),
      });
      continue;
    }

    if (
      action.type === "repay_with_reserves" ||
      action.type === "swap_reserve_to_borrowed_and_repay"
    ) {
      const capped = Math.min(amount, policy.max_reserve_deployment_usdc);
      projectedReserve -= capped;

      if (projectedReserve < policy.min_reserve_balance_usdc) {
        throw new Error(
          `action ${action.type} breaches reserve floor: projected ${projectedReserve.toFixed(2)} < floor ${policy.min_reserve_balance_usdc}`,
        );
      }

      executable.push({
        ...action,
        amountUsd: capped,
        repayPct: 0,
        venue: chooseVenue(action, policy.preferred_venues),
      });
      continue;
    }

    if (action.type === "partial_debt_repayment") {
      const boundedRepayPct = Math.min(cappedRepayPct, 100);
      executable.push({
        ...action,
        amountUsd: (risk.outstanding_debt_usd * boundedRepayPct) / 100,
        repayPct: boundedRepayPct,
        venue: chooseVenue(action, policy.preferred_venues),
      });
      continue;
    }

    if (action.type === "full_debt_repayment") {
      const amountUsd = risk.outstanding_debt_usd;
      projectedReserve -= amountUsd;

      if (projectedReserve < policy.min_reserve_balance_usdc) {
        throw new Error(
          `action ${action.type} breaches reserve floor: projected ${projectedReserve.toFixed(2)} < floor ${policy.min_reserve_balance_usdc}`,
        );
      }

      executable.push({
        ...action,
        amountUsd,
        repayPct: 100,
        venue: chooseVenue(action, policy.preferred_venues),
      });
      continue;
    }

    executable.push({
      ...action,
      amountUsd: amount,
      repayPct: 0,
      venue: chooseVenue(action, policy.preferred_venues),
    });
  }

  return orderActions(executable, policy.execution_sequence_preference);
};

const collectRiskSnapshot = (
  runtime: TeeRuntime<Config>,
  client: HTTPClient,
  baseUrl: string,
  exchangeApiKey: string,
): RiskState => {
  const authHeaders = {
    ...JSON_HEADERS,
    "x-exchange-api-key": exchangeApiKey,
  };

  const riskResponse = getJson(runtime, client, `${baseUrl}/risk-state`, authHeaders);
  return parseRiskState(riskResponse);
};

const createOpenAiPrompt = (
  risk: RiskState,
  riskScore: number,
  policy: Policy,
): string => {
  return JSON.stringify(
    {
      objective:
        "Return strict JSON with shouldDefend:boolean, reasoning:string, actions:[{type, amountUsd, repayPct, fromChain, toChain, venue}] for liquidation defense via collateral and debt-repayment actions while respecting policy.",
      policy,
      risk,
      riskScore,
    },
    null,
    2,
  );
};

export const onCronTrigger = async (runtime: TeeRuntime<Config>): Promise<string> => {
  const { mock_base_url, openai_url, openai_model, secrets_ids } = runtime.config;

  const secrets = runtime
    .getSecrets([
      { id: secrets_ids.exchange_api_key_id },
      { id: secrets_ids.openai_api_key_id },
      { id: secrets_ids.liquidation_warning_action_threshold_secret_id },
      { id: secrets_ids.minimum_health_factor_secret_id },
      { id: secrets_ids.target_health_factor_secret_id },
      { id: secrets_ids.maximum_stablecoin_reserve_deployment_secret_id },
      { id: secrets_ids.minimum_stablecoin_reserve_balance_secret_id },
      { id: secrets_ids.maximum_collateral_allocation_secret_id },
      { id: secrets_ids.maximum_partial_debt_repayment_secret_id },
      { id: secrets_ids.preferred_venues_secret_id },
    ])
    .result();

  const exchangeApiKey = secrets[secrets_ids.exchange_api_key_id].value;
  const openAiApiKey = secrets[secrets_ids.openai_api_key_id].value;
  const liquidationWarningActionThreshold = parseRequiredSecretNumber(
    secrets[secrets_ids.liquidation_warning_action_threshold_secret_id].value,
    secrets_ids.liquidation_warning_action_threshold_secret_id,
  );
  const minimumHealthFactor = parseRequiredSecretNumber(
    secrets[secrets_ids.minimum_health_factor_secret_id].value,
    secrets_ids.minimum_health_factor_secret_id,
  );
  const targetHealthFactor = parseRequiredSecretNumber(
    secrets[secrets_ids.target_health_factor_secret_id].value,
    secrets_ids.target_health_factor_secret_id,
  );
  const maxStablecoinReserveDeployment = parseRequiredSecretNumber(
    secrets[secrets_ids.maximum_stablecoin_reserve_deployment_secret_id].value,
    secrets_ids.maximum_stablecoin_reserve_deployment_secret_id,
  );
  const minStablecoinReserveBalance = parseRequiredSecretNumber(
    secrets[secrets_ids.minimum_stablecoin_reserve_balance_secret_id].value,
    secrets_ids.minimum_stablecoin_reserve_balance_secret_id,
  );
  const maxCollateralAllocation = parseRequiredSecretNumber(
    secrets[secrets_ids.maximum_collateral_allocation_secret_id].value,
    secrets_ids.maximum_collateral_allocation_secret_id,
  );
  const maxPartialDebtRepayment = parseRequiredSecretNumber(
    secrets[secrets_ids.maximum_partial_debt_repayment_secret_id].value,
    secrets_ids.maximum_partial_debt_repayment_secret_id,
  );
  const defensiveSequencePreference = parseExecutionSequencePreference(
    runtime.config.defensive_action_sequencing_preference,
  );
  const preferredVenues = parseVenueListSecret(
    secrets[secrets_ids.preferred_venues_secret_id].value,
  );
  runtime.log("liquidation-getsecrets-ok");

  const client = new HTTPClient();
  const risk = collectRiskSnapshot(runtime, client, mock_base_url, exchangeApiKey);

  const policy: Policy = {
    liquidation_warning_action_threshold: liquidationWarningActionThreshold,
    minimum_health_factor: minimumHealthFactor,
    target_health_factor: targetHealthFactor,
    max_reserve_deployment_usdc: maxStablecoinReserveDeployment,
    min_reserve_balance_usdc: minStablecoinReserveBalance,
    max_collateral_allocation_pct: maxCollateralAllocation,
    max_partial_debt_repayment_pct: maxPartialDebtRepayment,
    execution_sequence_preference: defensiveSequencePreference,
    preferred_venues: preferredVenues,
  };

  const riskScore = computeRiskScore(risk, policy);

  const prompt = createOpenAiPrompt(risk, riskScore, policy);

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
            "You are a liquidation-defense policy engine. Emit strict JSON only with key names exactly as requested.",
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
  const actions = enforcePolicy(decision, policy, risk);

  if (risk.liquidation_proximity_pct > policy.liquidation_warning_action_threshold || actions.length === 0) {
    runtime.log(
      `liquidation-no-action proximity=${risk.liquidation_proximity_pct.toFixed(3)} threshold=${policy.liquidation_warning_action_threshold} reason=${decision.reasoning}`,
    );
    return "SAFE";
  }

  const defenseResponse = postJson(
    runtime,
    client,
    `${mock_base_url}/execute-defense`,
    {
      riskScore,
      liquidationProximityPct: risk.liquidation_proximity_pct,
      reasoning: decision.reasoning,
      actions,
    },
    {
      ...JSON_HEADERS,
      "x-exchange-api-key": exchangeApiKey,
    },
  );

  runtime.log(
    `liquidation-defense-executed action_count=${actions.length} execution_id=${asString(defenseResponse.execution_id, "unknown")}`,
  );

  return JSON.stringify({
    status: "DEFENDED",
    actionCount: actions.length,
    riskScore,
    executionId: asString(defenseResponse.execution_id, "unknown"),
  });
};

export const initWorkflow = (config: Config): Workflow<Config> => {
  if (
    !config.schedule ||
    !config.mock_base_url ||
    !config.openai_url ||
    !config.openai_model ||
    !config.defensive_action_sequencing_preference
  ) {
    throw new Error(
      "config requires schedule, mock_base_url, openai_url, openai_model, and defensive_action_sequencing_preference",
    );
  }

  if (
    !config.secrets_ids?.exchange_api_key_id ||
    !config.secrets_ids?.openai_api_key_id ||
    !config.secrets_ids?.liquidation_warning_action_threshold_secret_id ||
    !config.secrets_ids?.minimum_health_factor_secret_id ||
    !config.secrets_ids?.target_health_factor_secret_id ||
    !config.secrets_ids?.maximum_stablecoin_reserve_deployment_secret_id ||
    !config.secrets_ids?.minimum_stablecoin_reserve_balance_secret_id ||
    !config.secrets_ids?.maximum_collateral_allocation_secret_id ||
    !config.secrets_ids?.maximum_partial_debt_repayment_secret_id ||
    !config.secrets_ids?.preferred_venues_secret_id
  ) {
    throw new Error("config requires secrets_ids fields");
  }

  const cron = new CronCapability();

  return [
    handlerInTee(
      cron.trigger({ schedule: config.schedule }),
      onCronTrigger,
      {},
      //[{ tee: "nitro", regions: [NITRO_REGIONS[0]] }],
    ),
  ];
};

export async function main() {
  const runner = await Runner.newRunner<Config>();
  await runner.run(initWorkflow);
}
