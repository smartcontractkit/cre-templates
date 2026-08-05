import express from "express";

const app = express();
app.use(express.json());

const PORT = Number(process.env.MOCK_PORT ?? 8787);
const EXCHANGE_KEY = process.env.MOCK_EXCHANGE_API_KEY ?? "mock-exchange-key";
const OPENAI_KEY = process.env.MOCK_OPENAI_API_KEY ?? "mock-openai-key";

const portfolioState = {
  holdings: {
    BTC: 0.55,
    ETH: 11,
    USDC: 12000,
  },
  prices_usd: {
    BTC: 102000,
    ETH: 3200,
    USDC: 1,
  },
  usdc_reserve: 12000,
  cash_balance_usd: 13500,
  volatility_index: 0.31,
};

const prices = {
  prices: {
    BTC: 102000,
    ETH: 3200,
    USDC: 1,
  },
};

const volatility = {
  volatility_index: 0.31,
};

const requireHeader = (req, res, expected, key) => {
  const got = req.header(key);
  if (got !== expected) {
    res.status(401).json({ error: `invalid ${key}` });
    return false;
  }
  return true;
};

app.get("/rebalancing/portfolio-state", (req, res) => {
  if (!requireHeader(req, res, EXCHANGE_KEY, "x-exchange-api-key")) {
    return;
  }
  res.json(portfolioState);
});

app.get("/rebalancing/prices", (req, res) => {
  if (!requireHeader(req, res, EXCHANGE_KEY, "x-exchange-api-key")) {
    return;
  }
  res.json(prices);
});

app.get("/rebalancing/volatility", (req, res) => {
  if (!requireHeader(req, res, EXCHANGE_KEY, "x-exchange-api-key")) {
    return;
  }
  res.json(volatility);
});

app.post("/rebalancing/v1/responses", (req, res) => {
  const auth = req.header("authorization") ?? "";
  if (auth !== `Bearer ${OPENAI_KEY}`) {
    res.status(401).json({ error: "invalid authorization" });
    return;
  }

  const input = Array.isArray(req.body?.input) ? req.body.input : [];
  const inputText = input
    .map((item) => (typeof item?.content === "string" ? item.content : ""))
    .join("\n");

  const shouldRebalance = !inputText.includes('"maxDriftPct": 0');
  const decision = {
    shouldRebalance,
    reasoning: shouldRebalance
      ? "Detected drift above threshold and selected constrained trade plan."
      : "No drift requiring action.",
    trades: shouldRebalance
      ? [
          {
            symbol: "BTC",
            side: "sell",
            notionalUsd: 7000,
            maxSlippageBps: 30,
            preferredVenues: ["binance", "onchain"],
          },
          {
            symbol: "ETH",
            side: "buy",
            notionalUsd: 7000,
            maxSlippageBps: 30,
            preferredVenues: ["coinbase", "binance"],
          },
        ]
      : [],
  };

  res.json({
    id: "resp_mock_1",
    output_text: JSON.stringify(decision),
  });
});

app.post("/rebalancing/execute-rebalance", (req, res) => {
  if (!requireHeader(req, res, EXCHANGE_KEY, "x-exchange-api-key")) {
    return;
  }

  const trades = Array.isArray(req.body?.trades) ? req.body.trades : [];
  const onchainTrades = Array.isArray(req.body?.onchain_trades) ? req.body.onchain_trades : [];
  const offchainTrades = Array.isArray(req.body?.offchain_trades) ? req.body.offchain_trades : [];

  res.json({
    status: "ok",
    execution_id: `exec_${Date.now()}`,
    accepted_trades: trades.length,
    accepted_onchain_trades: onchainTrades.length,
    accepted_offchain_trades: offchainTrades.length,
  });
});

app.listen(PORT, () => {
  console.log(`Mock server running at http://127.0.0.1:${PORT}`);
});
