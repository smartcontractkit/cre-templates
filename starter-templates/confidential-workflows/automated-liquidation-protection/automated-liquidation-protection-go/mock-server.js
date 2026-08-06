import express from "express";

const app = express();
app.use(express.json());

const PORT = Number(process.env.MOCK_PORT ?? 8787);
const EXCHANGE_KEY = process.env.MOCK_EXCHANGE_API_KEY ?? "mock-exchange-key";
const OPENAI_KEY = process.env.MOCK_OPENAI_API_KEY ?? "mock-openai-key";

const liquidationRiskState = {
  collateral_asset_symbol: "ETH",
  borrowed_asset_symbol: "USDC",
  collateral_asset_price_usd: 2510,
  borrowed_asset_price_usd: 1,
  collateral_balance_usd: 45000,
  outstanding_debt_usd: 32000,
  collateral_health_factor: 1.14,
  liquidation_proximity_pct: 12,
  loan_to_value_pct: 71,
  liquidation_threshold_pct: 78,
  volatility_index: 0.37,
  usdc_reserve: 10000,
  cash_balance_usd: 12500,
};

const requireHeader = (req, res, expected, key) => {
  const got = req.header(key);
  if (got !== expected) {
    res.status(401).json({ error: `invalid ${key}` });
    return false;
  }
  return true;
};

app.get("/liquidation/risk-state", (req, res) => {
  if (!requireHeader(req, res, EXCHANGE_KEY, "x-exchange-api-key")) {
    return;
  }
  res.json(liquidationRiskState);
});

app.post("/liquidation/v1/responses", (req, res) => {
  const auth = req.header("authorization") ?? "";
  if (auth !== `Bearer ${OPENAI_KEY}`) {
    res.status(401).json({ error: "invalid authorization" });
    return;
  }

  const input = Array.isArray(req.body?.input) ? req.body.input : [];
  const inputText = input
    .map((item) => (typeof item?.content === "string" ? item.content : ""))
    .join("\n");

  const shouldDefend = !inputText.includes('"liquidation_proximity_pct": 30');
  const decision = {
    shouldDefend,
    reasoning: shouldDefend
      ? "Liquidation proximity is elevated; add collateral and reduce debt exposure within policy limits."
      : "Position is sufficiently healthy; no defensive action required.",
    actions: shouldDefend
      ? [
          {
            type: "add_collateral",
            amountUsd: 3200,
            venue: "onchain",
          },
          {
            type: "partial_debt_repayment",
            repayPct: 18,
            venue: "binance",
          },
        ]
      : [],
  };

  res.json({
    id: "resp_mock_liquidation_1",
    output_text: JSON.stringify(decision),
  });
});

app.post("/liquidation/execute-defense", (req, res) => {
  if (!requireHeader(req, res, EXCHANGE_KEY, "x-exchange-api-key")) {
    return;
  }

  const actions = Array.isArray(req.body?.actions) ? req.body.actions : [];

  res.json({
    status: "ok",
    execution_id: `defense_${Date.now()}`,
    accepted_actions: actions.length,
  });
});

app.listen(PORT, () => {
  console.log(`Mock server running at http://127.0.0.1:${PORT}`);
});
