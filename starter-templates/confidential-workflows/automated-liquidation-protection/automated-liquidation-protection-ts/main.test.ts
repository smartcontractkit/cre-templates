import { describe, expect, test } from "bun:test";
import { computeRiskScore, enforcePolicy, initWorkflow } from "./main";
import type { Config } from "./main";

const policyFixture = {
  max_reserve_deployment_usdc: 5000,
  min_reserve_balance_usdc: 2000,
  liquidation_warning_action_threshold: 18,
  minimum_health_factor: 1.25,
  target_health_factor: 1.5,
  max_collateral_allocation_pct: 80,
  max_partial_debt_repayment_pct: 25,
  execution_sequence_preference: "collateral-first" as const,
  preferred_venues: ["binance", "onchain"],
};

const riskFixture = {
  collateral_asset_symbol: "ETH",
  borrowed_asset_symbol: "USDC",
  collateral_asset_price_usd: 2500,
  borrowed_asset_price_usd: 1,
  collateral_balance_usd: 45000,
  outstanding_debt_usd: 32000,
  collateral_health_factor: 1.12,
  liquidation_proximity_pct: 11,
  loan_to_value_pct: 71,
  liquidation_threshold_pct: 78,
  volatility_index: 0.35,
  usdc_reserve: 10000,
  cash_balance_usd: 12000,
};

const baseConfig: Config = {
  schedule: "0 */5 * * * *",
  mock_base_url: "http://127.0.0.1:8787/liquidation",
  openai_url: "http://127.0.0.1:8787/liquidation/v1/responses",
  openai_model: "gpt-4.1-mini",
  defensive_action_sequencing_preference: "collateral-first",
  secrets_ids: {
    exchange_api_key_id: "exchange_api_key",
    openai_api_key_id: "openai_api_key",
    liquidation_warning_action_threshold_secret_id: "liquidation_liquidation_warning_action_threshold",
    minimum_health_factor_secret_id: "liquidation_minimum_health_factor",
    target_health_factor_secret_id: "liquidation_target_health_factor",
    maximum_stablecoin_reserve_deployment_secret_id: "liquidation_maximum_stablecoin_reserve_deployment",
    minimum_stablecoin_reserve_balance_secret_id: "liquidation_minimum_stablecoin_reserve_balance",
    maximum_collateral_allocation_secret_id: "liquidation_maximum_collateral_allocation",
    maximum_partial_debt_repayment_secret_id: "liquidation_maximum_partial_debt_repayment",
    preferred_venues_secret_id: "liquidation_preferred_venues",
  },
};

describe("computeRiskScore", () => {
  test("increases score when liquidation and debt risk rise", () => {
    const score = computeRiskScore(riskFixture, policyFixture);
    expect(score).toBeGreaterThan(0);

    const safer = computeRiskScore(
      {
        ...riskFixture,
        liquidation_proximity_pct: 30,
        loan_to_value_pct: 55,
        collateral_health_factor: 1.5,
        volatility_index: 0.05,
      },
      policyFixture,
    );

    expect(score).toBeGreaterThan(safer);
  });
});

describe("enforcePolicy", () => {
  test("caps reserve deployment and honors sequence preference", () => {
    const actions = enforcePolicy(
      {
        shouldDefend: true,
        reasoning: "liquidation risk elevated",
        actions: [
          {
            type: "partial_debt_repayment",
            repayPct: 70,
          },
          {
            type: "add_collateral",
            amountUsd: 7000,
          },
        ],
      },
      policyFixture,
      riskFixture,
    );

    expect(actions).toHaveLength(2);
    expect(actions[0].type).toBe("add_collateral");
    expect(actions[0].amountUsd).toBe(5000);
    expect(actions[1].type).toBe("partial_debt_repayment");
    expect(actions[1].repayPct).toBe(25);
    expect(actions[1].amountUsd).toBe(8000);
    expect(actions[0].venue).toBe("binance");
  });

  test("rejects reserve floor breaches", () => {
    expect(() =>
      enforcePolicy(
        {
          shouldDefend: true,
          reasoning: "critical",
          actions: [{ type: "add_collateral", amountUsd: 9500 }],
        },
        policyFixture,
        {
          ...riskFixture,
          usdc_reserve: 6000,
        },
      ),
    ).toThrow("breaches reserve floor");
  });
});

describe("initWorkflow", () => {
  test("returns one tee handler with configured schedule", () => {
    const handlers = initWorkflow(baseConfig);
    expect(handlers).toHaveLength(1);
    expect(handlers[0].trigger.config.schedule).toBe(baseConfig.schedule);
    expect(handlers[0].requirements).toBeDefined();
  });

  test("throws when core fields are missing", () => {
    expect(() =>
      initWorkflow({
        ...baseConfig,
        openai_url: "",
      }),
    ).toThrow(
      "config requires schedule, mock_base_url, openai_url, openai_model, and defensive_action_sequencing_preference",
    );
  });
});
