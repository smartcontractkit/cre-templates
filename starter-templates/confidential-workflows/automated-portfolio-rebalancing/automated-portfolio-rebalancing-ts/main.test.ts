import { describe, expect, test } from "bun:test";
import {
  buildExecutionPlan,
  buildTargetAlignedTrades,
  calculateAllocations,
  enforcePolicyAndChunk,
  initWorkflow,
} from "./main";
import type { Config } from "./main";

const policyFixture = {
  target_allocations: {
    BTC: 0.5,
    ETH: 0.3,
    USDC: 0.2,
  },
  drift_threshold_pct: 5,
  max_trade_usd: 5000,
  reserve_floor_usdc: 2000,
  max_slippage_bps: 50,
  preferred_venues: ["binance", "onchain"],
  order_sequence_preference: "model-order" as const,
};

const baseConfig: Config = {
  schedule: "0 */5 * * * *",
  mock_base_url: "http://127.0.0.1:8787",
  openai_url: "http://127.0.0.1:8787/rebalancing/v1/responses",
  openai_model: "gpt-4.1-mini",
  order_sequence_preference: "model-order",
  secrets_ids: {
    exchange_api_key_id: "exchange_api_key",
    openai_api_key_id: "openai_api_key",
    target_allocation_btc_pct_secret_id: "rebalancing_target_allocation_btc_pct",
    target_allocation_eth_pct_secret_id: "rebalancing_target_allocation_eth_pct",
    target_allocation_usdc_pct_secret_id: "rebalancing_target_allocation_usdc_pct",
    drift_threshold_pct_secret_id: "rebalancing_drift_threshold_pct",
    max_trade_usd_secret_id: "rebalancing_max_trade_usd",
    reserve_floor_usdc_secret_id: "rebalancing_reserve_floor_usdc",
    max_slippage_bps_secret_id: "rebalancing_max_slippage_bps",
    preferred_venues_secret_id: "rebalancing_preferred_venues",
  },
};

describe("calculateAllocations", () => {
  test("computes drift and total value", () => {
    const result = calculateAllocations(
      { BTC: 1, ETH: 10, USDC: 1000 },
      { BTC: 50000, ETH: 3000, USDC: 1 },
      { BTC: 0.5, ETH: 0.3, USDC: 0.2 },
    );

    expect(result.totalUsd).toBe(81000);
    expect(result.allocations).toHaveLength(3);
    expect(result.maxDriftPct).toBeGreaterThan(0);
    expect(result.allocations.find((row) => row.symbol === "USDC")?.deltaUsdToTarget).toBe(15200);
  });

  test("derives target-aligned buy and sell trades from drift", () => {
    const result = calculateAllocations(
      { BTC: 0.2, ETH: 20, USDC: 20000 },
      { BTC: 100000, ETH: 2000, USDC: 1 },
      { BTC: 0.4, ETH: 0.2, USDC: 0.4 },
    );

    const trades = buildTargetAlignedTrades(result.allocations);

    expect(trades).toHaveLength(2);
    expect(trades.find((trade) => trade.symbol === "BTC")?.side).toBe("buy");
    expect(trades.find((trade) => trade.symbol === "ETH")?.side).toBe("sell");
  });
});

describe("enforcePolicyAndChunk", () => {
  test("chunks trades by max size and assigns preferred venue", () => {
    const chunks = enforcePolicyAndChunk(
      {
        shouldRebalance: true,
        reasoning: "drift above threshold",
        trades: [
          {
            symbol: "BTC",
            side: "buy",
            notionalUsd: 12000,
            maxSlippageBps: 40,
          },
        ],
      },
      policyFixture,
      20000,
    );

    expect(chunks).toHaveLength(3);
    expect(chunks[0].notionalUsd).toBe(5000);
    expect(chunks[2].notionalUsd).toBe(2000);
    expect(chunks[0].venue).toBe("binance");
  });

  test("rejects reserve floor breaches", () => {
    expect(() =>
      enforcePolicyAndChunk(
        {
          shouldRebalance: true,
          reasoning: "buy signal",
          trades: [{ symbol: "ETH", side: "buy", notionalUsd: 4000 }],
        },
        policyFixture,
        5000,
      ),
    ).toThrow("breaches reserve floor");
  });

  test("supports sell-first sequencing preference", () => {
    const chunks = enforcePolicyAndChunk(
      {
        shouldRebalance: true,
        reasoning: "drift above threshold",
        trades: [
          {
            symbol: "ETH",
            side: "buy",
            notionalUsd: 1000,
          },
          {
            symbol: "BTC",
            side: "sell",
            notionalUsd: 1000,
          },
        ],
      },
      {
        ...policyFixture,
        order_sequence_preference: "sells-first",
      },
      10000,
    );

    expect(chunks).toHaveLength(2);
    expect(chunks[0].side).toBe("sell");
  });
});

describe("buildExecutionPlan", () => {
  test("splits onchain and offchain trades by venue", () => {
    const plan = buildExecutionPlan([
      {
        symbol: "BTC",
        side: "sell",
        notionalUsd: 1000,
        maxSlippageBps: 20,
        venue: "binance",
        chunkIndex: 1,
        chunkCount: 1,
      },
      {
        symbol: "ETH",
        side: "buy",
        notionalUsd: 500,
        maxSlippageBps: 20,
        venue: "onchain",
        chunkIndex: 1,
        chunkCount: 1,
      },
    ]);

    expect(plan.offchainTrades).toHaveLength(1);
    expect(plan.onchainTrades).toHaveLength(1);
    expect(plan.onchainTrades[0].venue).toBe("onchain");
    expect(plan.offchainTrades[0].venue).toBe("binance");
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
      "config requires schedule, mock_base_url, openai_url, openai_model, and order_sequence_preference",
    );
  });
});
