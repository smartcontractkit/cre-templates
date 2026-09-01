# Automated Portfolio Rebalancing CRE Project

This standalone CRE project implements a confidential portfolio rebalancing workflow for crypto allocations.

## Description

The workflow continuously tracks allocation drift and triggers rebalancing when policy thresholds are exceeded. It is designed to restore user-defined target weights while protecting sensitive operational inputs, including exchange credentials, model credentials, policy thresholds, and execution preferences inside confidential execution.

## Target Customer

- Professional retail traders
- Developer shops
- Founders building trading products

## Structure

- `project.yaml`: project-level target settings
- `secrets.yaml`: secret ID mappings used by the workflow
- `mock-server.js`: local deterministic API server
- `automated-portfolio-rebalancing/`: workflow implementation

## Private Inputs

The following inputs are handled as confidential:

- Exchange API credentials used to read holdings, reserve data, and execution context, including stablecoin reserve depth and cash balance.
- LLM reasoning API credentials.
- Portfolio policy settings such as target allocation mix, minimum drift threshold, maximum trade size per execution, and required stablecoin reserve floor.
- Execution preferences such as venue priority, slippage limits, and trade chunking/ordering behavior.

## Workflow Notes

1. Monitor portfolio state.
	The workflow gathers market prices, current asset weights, drift from target allocations, reserve health, and volatility signals.
2. Enforce user-defined portfolio constraints.
	Confidential reasoning validates weight constraints, drift triggers, rebalance sizing limits, reserve protection requirements, and slippage controls.
3. Build a rebalance action plan.
	The plan can include buying underweight assets, selling overweight assets, enforcing reserve floors, capping per-trade notionals, and optimizing execution through chunking and smart venue routing.
4. Execute rebalance actions across venues.
	Depending on route selection, the workflow can execute both on-chain operations (such as swaps) and off-chain operations (such as centralized exchange API trades).

Note: Reasoning stages can be implemented with deterministic rule-based logic instead of an LLM when a fully rules-driven execution model is preferred.

## Required Environment Variables

Copy `.env.example` to `.env` and provide values for:

- `CRE_ETH_PRIVATE_KEY` (optional for local simulate)
- `MOCK_PORT`
- `MOCK_EXCHANGE_API_KEY`
- `MOCK_OPENAI_API_KEY`
- `MOCK_REBALANCING_TARGET_ALLOCATION_BTC_PCT`
- `MOCK_REBALANCING_TARGET_ALLOCATION_ETH_PCT`
- `MOCK_REBALANCING_TARGET_ALLOCATION_USDC_PCT`
- `MOCK_REBALANCING_DRIFT_THRESHOLD_PCT`
- `MOCK_REBALANCING_MAX_TRADE_USD`
- `MOCK_REBALANCING_RESERVE_FLOOR_USDC`
- `MOCK_REBALANCING_MAX_SLIPPAGE_BPS`
- `MOCK_REBALANCING_PREFERRED_VENUES`

The order sequence preference is not a secret; it is set via the `order_sequence_preference` field in each workflow's config file.

The local mock server for this project only exposes routes under `/rebalancing/*`.

## Quick Start

1. Install dependencies

```bash
bun install
```

2. Create environment file

```bash
cp .env.example .env
```

3. Start mock server

```bash
bun run mock:server
```

4. In another terminal, run checks

```bash
bun run typecheck
bun run test
```

5. Simulate workflow

```bash
cre workflow simulate ./automated-portfolio-rebalancing --target=staging-settings
```
