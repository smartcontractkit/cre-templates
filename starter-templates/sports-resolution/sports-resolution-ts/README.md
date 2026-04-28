# Sports Resolution Template

A CRE starter template that resolves a sports game outcome on-chain using the full CRE workflow pattern: an **EVM Log Trigger** fires when a prediction market emits `SettlementRequested`, 2–3 sports APIs are queried via the **HTTP Capability** (each BFT-verified across DON nodes), workflow-level aggregation determines consensus, and a **cryptographically signed report** is written back on-chain.

Sports is the illustrative domain. The workflow pattern — *EVM log trigger → HTTP fetches → aggregation → signed onchain write* — is fully domain-agnostic and directly reusable for financial data, election outcomes, weather events, or any binary/categorical resolution use case.

> **Important:** This template uses user-supplied HTTP APIs for resolution. For resolution using Chainlink Data Feeds (on-chain price sources), see the [Prediction Market template](../prediction-market/).

**⚠️ DISCLAIMER**

This template is an educational example to demonstrate how to interact with Chainlink systems, products, and services. It is provided **"AS IS"** and **"AS AVAILABLE"** without warranties of any kind, has **not** been audited, and may omit checks or error handling for clarity. **Do not use this code in production** without performing your own audits and applying best practices. Neither Chainlink Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for unintended outputs generated due to errors in code.

---

## Architecture

```
Prediction market contract emits SettlementRequested(gameId)
                    │
                    ▼ EVM Log Trigger
        ┌───────────────────────┐
        │   CRE Workflow (WASM) │
        │                       │
        │  ┌─────────────────┐  │
        │  │ HTTP Capability │  │  ← BFT-verified per source:
        │  │                 │  │    each DON node fetches
        │  │ Source 1 fetch  │  │    independently; results
        │  │ Source 2 fetch  │  │    are consensus-verified
        │  │ Source 3 fetch  │  │    before callback sees them
        │  │ (optional)      │  │
        │  └────────┬────────┘  │
        │           │           │
        │  ┌────────▼────────┐  │
        │  │  Aggregation    │  │  ← TypeScript business logic
        │  │  majority /     │  │    in the callback function
        │  │  unanimous      │  │
        │  └────────┬────────┘  │
        │           │           │
        │  ┌────────▼────────┐  │
        │  │  Signed Report  │  │  ← Cryptographically signed
        │  │  EVM Write      │  │    by DON; contract verifies
        │  └─────────────────┘  │    Forwarder signature
        └───────────────────────┘
                    │
                    ▼
        SportsMarket.sol → GameSettled(gameId, outcome)
```

### CRE Concepts Demonstrated

| Concept | How It's Used |
|---|---|
| **EVM Log Trigger** | Fires on `SettlementRequested` — event-driven, no polling |
| **HTTP Capability** | Fetches from 2–3 sports APIs; each fetch BFT-verified across DON nodes |
| **BFT Consensus (built-in)** | Every `HTTPClient.sendRequest()` is verified across DON nodes automatically |
| **Workflow aggregation** | Custom majority/unanimous logic in TypeScript callback |
| **EVM Write + signed report** | Cryptographically verified result delivered on-chain |
| **Simulation** | Full workflow testable locally via `cre workflow simulate` |

---

## Project Structure

```
sports-resolution-ts/
├── .cre/template.yaml          # Template metadata
├── project.yaml                # CRE project settings (RPCs)
├── secrets.yaml                # Optional secret name declarations
├── README.md
├── contracts/                  # Smart contract + generated TypeScript bindings
│   ├── package.json
│   ├── abi/                    # Re-exported ABI for external use
│   └── evm/
│       ├── src/                # Solidity source + compiled ABI
│       │   ├── SportsMarket.sol
│       │   ├── ReceiverTemplate.sol
│       │   ├── IReceiver.sol
│       │   ├── IERC165.sol
│       │   └── abi/SportsMarket.abi
│       └── ts/generated/       # Generated TypeScript bindings
│           ├── SportsMarket.ts
│           └── SportsMarket_mock.ts
└── game-resolution/            # The CRE workflow
    ├── main.ts                 # Runner entrypoint
    ├── workflow.ts             # Core logic: trigger → fetch → aggregate → write
    ├── workflow.test.ts        # Unit and integration tests
    ├── workflow.yaml           # Staging/production workflow settings
    ├── config.staging.json
    ├── config.production.json
    ├── package.json
    └── tsconfig.json
```

---

## Quick Start

### Prerequisites

- [Bun](https://bun.sh/) installed
- [Foundry](https://getfoundry.sh/) installed (`forge`, `cast`)
- [CRE CLI](https://docs.chain.link/cre) installed and authenticated (`cre login`)
- A deployed `SportsMarket` contract (see [Contract Deployment](#contract-deployment))

### Environment Variables

| Variable | Required | Description |
|---|---|---|
| `PRIVATE_KEY` | Yes | EOA private key for deploying the contract and calling `requestSettlement` |
| `CRE_ETH_PRIVATE_KEY` | Yes | Private key used by the CRE CLI to sign broadcast transactions (`--broadcast`) |
| `ETHERSCAN_API_KEY` | Optional | Etherscan API key for contract verification |

Set them in your shell before running any commands:

```bash
export PRIVATE_KEY=0x...
export CRE_ETH_PRIVATE_KEY=0x...
export ETHERSCAN_API_KEY=...  # optional
```

### 1. Install Dependencies

```bash
# Contract bindings (needed for test module resolution)
cd contracts && bun install && cd ..

# Workflow
cd game-resolution && bun install && cd ..
```

### 2. Configure

Edit `game-resolution/config.staging.json`:

```json
{
  "evms": [{
    "chainSelectorName": "ethereum-testnet-sepolia",
    "sportsMarketAddress": "0xYOUR_DEPLOYED_CONTRACT",
    "gasLimit": "500000"
  }],
  "dataSourceUrls": [
    "https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard",
    "https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard"
  ],
  "aggregationMode": "majority"
}
```

### 3. Run Tests

```bash
cd game-resolution && bun test && cd ..
```

### 4. Simulate

First, call `requestSettlement(gameId, description)` on your deployed contract to emit a `SettlementRequested` event. Then simulate the workflow against that transaction:

```bash
cre workflow simulate game-resolution --target staging-settings \
  --non-interactive --trigger-index 0 \
  --evm-tx-hash 0xTX_THAT_EMITTED_SETTLEMENT_REQUESTED \
  --evm-event-index 0
```

### 5. Broadcast (write on-chain)

```bash
cre workflow simulate game-resolution --target staging-settings \
  --non-interactive --trigger-index 0 \
  --evm-tx-hash 0xTX_THAT_EMITTED_SETTLEMENT_REQUESTED \
  --evm-event-index 0 \
  --broadcast
```

---

## Workflow Details

**Trigger:** EVM Log Trigger on `SettlementRequested(gameId)`

The workflow fires once per event and executes the following steps:

1. **Decode event** — extract `gameId` and `description` from the typed log payload
2. **Guard read** — call `sportsMarket.getGame(gameId)` on-chain; skip if already settled (`outcome != 0`)
3. **Fetch per source** — for each URL in `dataSourceUrls`, call `HTTPClient.sendRequest()`:
   - The DON's nodes each fetch the URL independently
   - `ConsensusAggregationByFields` with `median` produces a BFT-verified `{ homeScore, awayScore }` result
4. **Compute outcomes** — apply `computeOutcome()` to each verified result: `homeScore > awayScore → 1`, `awayScore > homeScore → 2`, equal → `3`
5. **Aggregate** — apply `applyAggregation()` across source outcomes:
   - `majority`: strictly more than half must agree (2-of-2 or 2-of-3)
   - `unanimous`: all sources must agree
   - Returns `null` (no consensus) if threshold not met — no on-chain write
6. **Write report** — encode `(uint256 gameId, uint8 outcome)` and call `sportsMarket.writeReport()`, which submits a signed report via the CRE Forwarder

### Outcome Values

| Value | Meaning |
|---|---|
| `1` | HomeWins |
| `2` | AwayWins |
| `3` | Draw |
| `null` | No consensus — no on-chain write; retry or escalate |

---

## API Integration

The workflow appends `/{gameId}` to each `dataSourceUrls` entry to form the full request URL. The default implementation uses ESPN's per-game scoreboard endpoint (already configured in `config.staging.json`):

```
GET https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard/{gameId}
```

This returns ~20KB per request — well within CRE's 100KB HTTP response limit. Adapt `fetchGameResult()` in `workflow.ts` for other providers:

### SportsDataIO

URL pattern: `https://api.sportsdata.io/v3/nba/scores/json/BoxScore/{gameId}?key=YOUR_KEY`

```typescript
const body = JSON.parse(...)
return {
  homeScore: body.Game.HomeTeamScore,
  awayScore: body.Game.AwayTeamScore,
}
```

### Sportradar

URL pattern: `https://api.sportradar.com/nba/production/v8/en/games/{gameId}/summary.json?api_key=YOUR_KEY`

```typescript
const body = JSON.parse(...)
return {
  homeScore: body.game.home.points,
  awayScore: body.game.away.points,
}
```

### ESPN (public API — pre-configured)

URL pattern: `https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard/{gameId}`

```typescript
const body = JSON.parse(...)
const competitors = body.competitions[0].competitors
const home = competitors.find((c: any) => c.homeAway === 'home')
const away = competitors.find((c: any) => c.homeAway === 'away')
return {
  homeScore: parseInt(home.score, 10),
  awayScore: parseInt(away.score, 10),
}
```

### Using API Keys

The default ESPN example uses a public endpoint, so `secrets.yaml` and `workflow.yaml` do not declare secrets by default.

If you adapt this template to a private sports data provider, do not hardcode API keys in committed config files or URLs. Add the secret names to `secrets.yaml`, set `secrets-path` in `game-resolution/workflow.yaml`, and extend `fetchGameResult()` to inject the key according to your provider's authentication model, such as a request header or signed URL.

---

## Contract Deployment

The `SportsMarket` contract accepts CRE-signed reports via the `ReceiverTemplate` base class, which validates the Chainlink Forwarder signature before processing.

### Option A: Deploy with Foundry (Forge)

```bash
cd contracts/evm

forge create --broadcast \
  --private-key $PRIVATE_KEY \
  --rpc-url https://ethereum-sepolia-rpc.publicnode.com \
  src/SportsMarket.sol:SportsMarket \
  --constructor-args 0xf8344cfd5c43616a4366c34e3eee75af79a74482
```

### Option B: Deploy with Hardhat Ignition

1. Copy `SportsMarket.sol`, `ReceiverTemplate.sol`, `IReceiver.sol`, `IERC165.sol` into `contracts/`.

2. Create `ignition/modules/SportsMarket.ts`:
```typescript
import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

const FORWARDER = "0xf8344cfd5c43616a4366c34e3eee75af79a74482"; // Sepolia CRE Forwarder

const SportsMarketModule = buildModule("SportsMarketModule", (m) => {
  const forwarder = m.getParameter("forwarder", FORWARDER);
  const sportsMarket = m.contract("SportsMarket", [forwarder]);
  return { sportsMarket };
});

export default SportsMarketModule;
```

3. Deploy:
```bash
npx hardhat ignition deploy ignition/modules/SportsMarket.ts --network sepolia
```

### Post-Deployment

1. Update `sportsMarketAddress` in `config.staging.json`
2. Emit a `SettlementRequested` event by calling `requestSettlement(gameId, description)` on the contract
3. Simulate the workflow against that transaction hash

### Constructor Arguments

| Argument | Value | Description |
|---|---|---|
| `forwarder` | `0xf8344cfd5c43616a4366c34e3eee75af79a74482` | Sepolia CRE Forwarder |

---

## Aggregation Modes

| Mode | Logic | Use When |
|---|---|---|
| `majority` | >50% of sources agree — 2-of-2 or 2-of-3 | Standard resolution; tolerates one outlier source |
| `unanimous` | 100% of sources must agree | High-stakes markets; any disagreement triggers retry |

When no consensus is reached, the workflow logs the source breakdown and returns without writing on-chain. Implement retry logic in your prediction market contract (e.g., re-emit `SettlementRequested` after a delay) or set up a monitoring alert.

---

## Customization Guide

### Add a Third Data Source

Add a third URL to `dataSourceUrls` in config. Majority logic automatically upgrades to 2-of-3.

### Change to a Different Sport or Domain

1. Update `dataSourceUrls` to point to your data provider endpoints
2. Adapt `fetchGameResult()` to parse the provider's response format
3. No contract or trigger changes needed

### Non-Binary Outcomes (e.g., multi-team tournaments)

Extend the `Outcome` enum in `SportsMarket.sol` and update `computeOutcome()` in the workflow.

### Different Chain

Update `chainSelectorName` in config and the RPC in `project.yaml`. The contract needs to be deployed on the target chain with the correct Forwarder address for that network.

---

## Security Notes

- `SportsMarket` inherits `ReceiverTemplate`, which enforces that only the trusted Chainlink Forwarder can call `onReport()`. No third party can inject a fake settlement.
- Each `HTTPClient.sendRequest()` call is BFT-verified across DON nodes — a single compromised node cannot influence the result.
- The workflow-level aggregation (`applyAggregation`) provides additional protection against a single API source returning incorrect data.
- Each workflow execution uses `LAST_FINALIZED_BLOCK_NUMBER` for on-chain reads to ensure data consistency.
- `AlreadySettled` guards in both `requestSettlement()` and `_processReport()` prevent double-settlement even if the workflow fires twice.

## What's Not Included

- **Betting/staking:** No token deposits or payouts — this is resolution infrastructure, not a full dApp.
- **Retry logic:** When no consensus is reached, the workflow returns without writing. Retry is left to the market contract's design.
- **Frontend:** No UI. These are backend workflow templates.
- **Cross-chain:** Single chain only. See the [Multi-Chain Token Manager](../multi-chain-token-manager/) template for CCIP patterns.
