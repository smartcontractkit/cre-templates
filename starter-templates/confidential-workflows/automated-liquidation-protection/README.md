# Automated Liquidation Protection CRE Project

This standalone CRE project implements a confidential liquidation-defense workflow for DeFi lending positions.

## Description

The workflow continuously evaluates borrower risk and takes defensive action before a position becomes unsafe. During periods of high volatility, it can increase collateral, reduce debt, or combine both strategies based on policy constraints. Sensitive operational data remains protected in confidential execution, including exchange credentials, model credentials, and user-defined risk thresholds.

## Target Customer

- Professional retail traders
- Developer shops
- Founders building trading products

## Structure

- `project.yaml`: project-level target settings
- `secrets.yaml`: secret ID mappings used by the workflow
- `mock-server.js`: local deterministic API server
- `automated-liquidation-protection/`: workflow implementation

## Private Inputs

The following inputs are treated as confidential:

- Exchange API credentials used to fetch account context such as stablecoin reserves and available cash balance.
- LLM reasoning API credentials and policy parameters used to govern defense behavior, including minimum and target health factors, reserve deployment caps, minimum reserve balance, and collateral allocation limits.
- Execution preferences that define how defense actions should be sequenced.

## Workflow Notes

1. Observe liquidation risk signals.
   The workflow tracks collateral and debt asset pricing, health factor, liquidation proximity, LTV, liquidation threshold, and market volatility.
2. Enforce user policy constraints.
   Confidential reasoning evaluates how much capital can be deployed, whether debt reduction should be prioritized, which reserve assets are eligible, and what execution sequence is preferred.
3. Select defense actions.
   The workflow builds a response plan that may include collateral-focused moves (deposit, bridge, swap-then-deposit) and debt-focused moves (repay, swap-then-repay, partial payoff, full payoff).
4. Execute the approved defense plan.

Note: Every reasoning stage can be implemented with deterministic rule-based logic instead of an LLM, if your deployment requires a fully rules-driven policy engine.

## Required Environment Variables

Copy `.env.example` to `.env` and provide values for:

- `CRE_ETH_PRIVATE_KEY` (optional for local simulate)
- `MOCK_PORT`
- `MOCK_EXCHANGE_API_KEY`
- `MOCK_OPENAI_API_KEY`
- `MOCK_LIQUIDATION_WARNING_ACTION_THRESHOLD`
- `MOCK_LIQUIDATION_MINIMUM_HEALTH_FACTOR`
- `MOCK_LIQUIDATION_TARGET_HEALTH_FACTOR`
- `MOCK_LIQUIDATION_MAX_STABLECOIN_RESERVE_DEPLOYMENT`
- `MOCK_LIQUIDATION_MIN_STABLECOIN_RESERVE_BALANCE`
- `MOCK_LIQUIDATION_MAX_COLLATERAL_ALLOCATION`
- `MOCK_LIQUIDATION_MAX_PARTIAL_DEBT_REPAYMENT`
- `MOCK_LIQUIDATION_DEFENSIVE_ACTION_SEQUENCING_PREFERENCE`
- `MOCK_LIQUIDATION_PREFERRED_VENUES`

The local mock server for this project only exposes routes under `/liquidation/*`.

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
cre workflow simulate ./automated-liquidation-protection --project-root ./ --target=staging-settings --env ./.env
```
