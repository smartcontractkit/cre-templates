# Prediction Market Template

A CRE starter template implementing a full prediction market lifecycle with three workflows: **Market Creation**, **Market Resolution**, and **Dispute Management**. All three workflows share a single `PredictionMarket` smart contract and use **Chainlink Data Feeds** (BTC/USD) for deterministic on-chain resolution.

> **Important:** This template uses Chainlink Data Feeds for deterministic resolution. For AI-based resolution, see the [CRE Prediction Market Demo](https://github.com/smartcontractkit/cre-gcp-prediction-market-demo/tree/main/cre-workflow).

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  1. CREATION     │────>│  2. RESOLUTION   │────>│  3. DISPUTE      │
│                  │     │                  │     │                  │
│  Cron trigger    │     │  Cron trigger    │     │  LogTrigger      │
│  Creates new     │     │  Checks expired  │     │  (DisputeRaised) │
│  markets         │     │  markets, reads  │     │  Re-reads price  │
│                  │     │  BTC/USD feed,   │     │  feed, resolves  │
│                  │     │  resolves        │     │  dispute         │
└─────────────────┘     └─────────────────┘     └─────────────────┘
        │                        │                        │
        └────────────────────────┼────────────────────────┘
                                 │
                    ┌────────────▼────────────┐
                    │   PredictionMarket.sol   │
                    │   (shared contract)      │
                    │                          │
                    │  - createMarket()        │
                    │  - resolveMarket()       │
                    │  - resolveDispute()      │
                    │  - raiseDispute()        │
                    └────────────┬────────────┘
                                 │
                    ┌────────────▼────────────┐
                    │  Chainlink BTC/USD Feed  │
                    │  (on-chain price data)   │
                    └─────────────────────────┘
```

## Project Structure

```
prediction-market-ts/
├── .cre/template.yaml          # Template metadata
├── project.yaml                # CRE project settings (RPCs)
├── secrets.yaml                # Secrets configuration
├── .env                        # Private key for deployment
├── README.md
├── contracts/                  # Shared contracts and bindings
│   ├── package.json
│   ├── abi/                    # Human-readable ABI (viem parseAbi)
│   └── evm/
│       ├── src/                # Solidity source + ABI JSON
│       │   ├── PredictionMarket.sol
│       │   ├── ReceiverTemplate.sol
│       │   ├── IReceiver.sol
│       │   ├── IERC165.sol
│       │   └── abi/            # Compiled ABI JSON files
│       └── ts/generated/       # Generated TypeScript bindings
├── market-creation/            # Workflow 1: Create markets
│   ├── main.ts
│   ├── workflow.ts
│   ├── workflow.test.ts
│   ├── workflow.yaml
│   ├── config.staging.json
│   └── config.production.json
├── market-resolution/          # Workflow 2: Resolve expired markets
│   ├── main.ts
│   ├── workflow.ts
│   ├── workflow.test.ts
│   ├── workflow.yaml
│   ├── config.staging.json
│   └── config.production.json
└── market-dispute/             # Workflow 3: Handle disputes
    ├── main.ts
    ├── workflow.ts
    ├── workflow.test.ts
    ├── workflow.yaml
    ├── config.staging.json
    └── config.production.json
```

## Quick Start

### Prerequisites

- [Bun](https://bun.sh/) installed
- [CRE CLI](https://docs.chain.link/cre) installed

### 1. Install Dependencies

```bash
# Install contract dependencies
cd contracts && bun install && cd ..

# Install workflow dependencies
cd market-creation && bun install && cd ..
cd market-resolution && bun install && cd ..
cd market-dispute && bun install && cd ..
```

### 2. Run Tests

```bash
cd market-creation && bun test && cd ..
cd market-resolution && bun test && cd ..
cd market-dispute && bun test && cd ..
```

### 3. Simulate Workflows

```bash
# Simulate market creation (creates a new market)
cre workflow simulate market-creation --target staging-settings

# Simulate market resolution (checks and resolves expired markets)
cre workflow simulate market-resolution --target staging-settings

# Simulate market dispute (requires a real DisputeRaised tx hash)
# First call raiseDispute() on the contract, then:
cre workflow simulate market-dispute --target staging-settings \
  --non-interactive --trigger-index 0 \
  --evm-tx-hash 0xTX_THAT_EMITTED_DISPUTE_RAISED \
  --evm-event-index 0
```

### 4. Broadcast (write on-chain)

```bash
cre workflow simulate market-creation --target staging-settings --broadcast
```

## Workflow Details

### Workflow 1: Market Creation

**Trigger:** Cron (hourly by default)

Creates binary prediction markets on a schedule. On each cron tick:

1. Uses `runtime.now()` to get the current consensus timestamp
2. Computes `expirationTime = now + durationSeconds` (default 86400 = 24 hours)
3. Formats the question by replacing `{expirationDate}` in the config template (e.g., "Will BTC be above $100,000 by 2026-03-18?")
4. Scales the strike price from USD to Chainlink feed decimals using viem's `parseUnits` (e.g., $100,000 becomes `10000000000000` with 8 decimals)
5. Reads `nextMarketId` from the contract for logging
6. Writes on-chain: sends `ACTION_CREATE` report with `(question, strikePrice, expirationTime)` — the contract creates the market struct and emits `MarketCreated`

### Workflow 2: Market Resolution

**Trigger:** Cron (every 10 minutes by default)

Polls a list of market IDs from config and resolves any that are expired. On each cron tick:

1. Loops through `marketIdsToCheck` from config (default `[0, 1, 2]`)
2. For each market, calls `isResolvable()` on the contract — returns `true` only if status is `Open` AND `block.timestamp >= expirationTime`
3. If not resolvable (not expired yet or already resolved): logs and skips
4. If resolvable: reads `latestAnswer()` from the Chainlink BTC/USD Data Feed on Sepolia (live price, e.g., $74,406)
5. Writes on-chain: sends `ACTION_RESOLVE` report with `(marketId, price)` — the contract compares `price >= strikePrice`, sets outcome to `Yes` or `No`, and emits `MarketResolved`

### Workflow 3: Dispute Management

**Trigger:** LogTrigger on `DisputeRaised` event

Fires when anyone calls `raiseDispute(marketId, reason)` on the contract. The contract sets the market status to `Disputed` and emits the `DisputeRaised` event. The workflow then:

1. Decodes the event using generated typed bindings — receives `marketId`, `disputor` address, and `reason` string
2. Verifies the market is still in `Disputed` status (status == 2) via on-chain read
3. Re-reads `latestAnswer()` from the Chainlink BTC/USD feed (fresh price at dispute time)
4. Writes on-chain: sends `ACTION_RESOLVE_DISPUTE` report with `(marketId, newPrice)` — the contract re-evaluates the outcome against the strike price, sets an `overturned` flag if the outcome changed, and emits `DisputeResolved`

## Contract Deployment

A demo `PredictionMarket` contract is pre-deployed on Sepolia at `0xEb792aF46AB2c2f1389A774AB806423DB43aA425`.

To deploy your own instance:

### Option A: Deploy with Foundry (Forge)

```bash
forge create --broadcast \
  --private-key $PRIVATE_KEY \
  --rpc-url https://ethereum-sepolia-rpc.publicnode.com \
  src/PredictionMarket.sol:PredictionMarket \
  --constructor-args \
    0x15fc6ae953e024d975e77382eeec56a9101f9f88 \
    0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43 \
    86400
```

### Option B: Deploy with Hardhat Ignition

1. Set up a Hardhat project:
```bash
mkdir deploy-prediction-market && cd deploy-prediction-market
npm init -y
npm install --save-dev hardhat @nomicfoundation/hardhat-ignition @nomicfoundation/hardhat-ignition-ethers
npx hardhat init
```

2. Copy the Solidity files (`PredictionMarket.sol`, `ReceiverTemplate.sol`, `IReceiver.sol`, `IERC165.sol`) into `contracts/`.

3. Create an Ignition module at `ignition/modules/PredictionMarket.ts`:
```typescript
import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

const FORWARDER = "0x15fc6ae953e024d975e77382eeec56a9101f9f88";
const BTC_USD_FEED = "0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43";
const DISPUTE_WINDOW = 86400; // 24 hours

const PredictionMarketModule = buildModule("PredictionMarketModule", (m) => {
  const forwarder = m.getParameter("forwarder", FORWARDER);
  const priceFeed = m.getParameter("priceFeed", BTC_USD_FEED);
  const disputeWindow = m.getParameter("disputeWindow", DISPUTE_WINDOW);

  const predictionMarket = m.contract("PredictionMarket", [
    forwarder,
    priceFeed,
    disputeWindow,
  ]);

  return { predictionMarket };
});

export default PredictionMarketModule;
```

4. Configure Sepolia in `hardhat.config.ts`:
```typescript
import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-ignition-ethers";

const config: HardhatUserConfig = {
  solidity: "0.8.19",
  networks: {
    sepolia: {
      url: "https://ethereum-sepolia-rpc.publicnode.com",
      accounts: [process.env.PRIVATE_KEY!],
    },
  },
};

export default config;
```

5. Deploy:
```bash
npx hardhat ignition deploy ignition/modules/PredictionMarket.ts --network sepolia
```

### Post-Deployment Steps

1. Update all `config.staging.json` files with the new contract address
2. Create a test market: `cre workflow simulate market-creation --target staging-settings --broadcast`
3. Wait for Sepolia finalization (5-15 minutes) before running resolution simulations

### Constructor Arguments

| Argument | Value | Description |
|---|---|---|
| `forwarder` | `0x15fc6ae953e024d975e77382eeec56a9101f9f88` | Sepolia MockKeystoneForwarder |
| `priceFeed` | `0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43` | Sepolia BTC/USD Chainlink Data Feed |
| `disputeWindow` | `86400` | 24 hours in seconds |

## Alternative: Create Markets via HTTP Trigger

The `market-creation` workflow uses a Cron trigger to create markets on a schedule. In many production scenarios, you want markets created on-demand -- triggered by a frontend, a DAO governance action, or an external API call. You can achieve this by replacing the Cron trigger with an HTTP trigger.

```
┌──────────────────┐     HTTP POST (signed)
│  Frontend / API  │─────────────────────────────────┐
│  or DAO action   │     {                           │
│                  │       "question": "Will BTC...",│
│                  │       "strikePriceUsd": 100000, │
│                  │       "durationSeconds": 86400  │
│                  │     }                           │
└──────────────────┘                                 │
                                                     v
┌────────────────────────────────────────────────────────────┐
│                         CRE DON                            │
│  ┌────────────────┐   ┌─────────────┐   ┌──────────────┐  │
│  │  HTTP Trigger  │-->│ Parse JSON  │-->│ Write report │  │
│  │  (authorized   │   │ market      │   │ ACTION_CREATE│  │
│  │   keys)        │   │ params      │   │              │  │
│  └────────────────┘   └─────────────┘   └──────────────┘  │
└────────────────────────────────────────────────────────────┘
```

**Key changes from the Cron version:**
```typescript
// Instead of:
const cronTrigger = new cre.capabilities.CronCapability()
cre.handler(cronTrigger.trigger({ schedule: config.schedule }), onCronTrigger)

// You would use:
const httpTrigger = new cre.capabilities.HTTPCapability()
cre.handler(
  httpTrigger.trigger({ authorizedKeys: config.authorizedKeys }),
  onHttpTrigger,
)

// The callback receives market params from the HTTP payload:
const onHttpTrigger = (runtime: Runtime<Config>, payload: HTTPPayload): string => {
  const params = decodeJson(payload.input)
  // params.question, params.strikePriceUsd, params.durationSeconds
  // ... same create logic as the Cron version, but params come from the request
}
```

**When to use this instead of Cron:**
- Your frontend lets users propose new markets
- A DAO governance vote triggers market creation
- An external system (e.g., sports data provider) pushes new events

> **Note:** `authorizedKeys` must be empty for simulation but populated with real signing keys before deployment. See the [HTTP Trigger documentation](https://docs.chain.link/cre/guides/workflow/using-triggers/http-trigger) for details.

## Alternative: Resolve Markets via HTTP Trigger

The `market-resolution` workflow uses a Cron trigger to periodically check for expired markets. In some scenarios, you may want to trigger resolution immediately -- from a keeper bot, a frontend "Resolve" button, or a third-party settlement service.

```
┌──────────────────┐     HTTP POST (signed)
│  Keeper bot /    │─────────────────────────────────┐
│  Frontend /      │     {                           │
│  Settlement svc  │       "marketIds": [0, 1, 5]   │
│                  │     }                           │
└──────────────────┘                                 │
                                                     v
┌────────────────────────────────────────────────────────────┐
│                         CRE DON                            │
│  ┌────────────────┐   ┌──────────────┐   ┌─────────────┐  │
│  │  HTTP Trigger  │-->│ Read BTC/USD │-->│ Resolve     │  │
│  │  (authorized   │   │ from Chainlink│  │ each market │  │
│  │   keys)        │   │ Data Feed    │   │             │  │
│  └────────────────┘   └──────────────┘   └─────────────┘  │
└────────────────────────────────────────────────────────────┘
```

**Key changes from the Cron version:**
```typescript
// Instead of:
const cronTrigger = new cre.capabilities.CronCapability()
cre.handler(cronTrigger.trigger({ schedule: config.schedule }), onCronTrigger)

// You would use:
const httpTrigger = new cre.capabilities.HTTPCapability()
cre.handler(
  httpTrigger.trigger({ authorizedKeys: config.authorizedKeys }),
  onHttpTrigger,
)

// The callback receives market IDs from the HTTP payload:
const onHttpTrigger = (runtime: Runtime<Config>, payload: HTTPPayload): string => {
  const { marketIds } = decodeJson(payload.input)
  // ... same resolution logic, but marketIds come from the request instead of config
}
```

**When to use this instead of Cron:**
- A frontend "Resolve Now" button for individual markets
- An external keeper service that monitors expiration off-chain and triggers resolution precisely
- Batch settlement at a specific time (e.g., end of trading day)

> **Note:** The on-chain price read from the Chainlink Data Feed is identical in both the Cron and HTTP versions -- only the trigger mechanism changes. The resolution logic and DON consensus verification remain the same.

## Customization Guide

### Use a Different Asset (e.g., ETH/USD)

1. Deploy a new `PredictionMarket` contract with a different `priceFeed` address (e.g., Sepolia ETH/USD: `0x694AA1769357215DE4FAC081bf1f309aDC325306`)
2. Update `priceFeedAddress` in all `config.staging.json` files
3. Update the question template in `market-creation/config.staging.json`

### Change Market Duration

Update `durationSeconds` in `market-creation/config.staging.json`:
- 1 hour: `3600`
- 1 day: `86400`
- 1 week: `604800`

### Change Resolution Frequency

Update `schedule` in `market-resolution/config.staging.json`:
- Every minute: `"0 * * * * *"`
- Every 10 minutes: `"0 */10 * * * *"`
- Every hour: `"0 0 * * * *"`

### Change Dispute Window

Deploy a new contract with a different `disputeWindow` constructor argument.

## Security Notes

- The `PredictionMarket` contract inherits from `ReceiverTemplate`, which validates that only the trusted Chainlink Forwarder can call `onReport()`
- Resolution uses on-chain Chainlink Data Feeds only -- no off-chain API calls or AI models
- The dispute window ensures finality: after the deadline, no more disputes can be raised
- Each workflow uses `LAST_FINALIZED_BLOCK_NUMBER` for on-chain reads, ensuring data consistency

## What's Not Included

- **Betting/staking:** No token deposits or payouts. This is infrastructure, not a full dApp.
- **AI resolution:** Not included in this template. Chainlink Data Feeds provide deterministic, verifiable resolution. For AI-based resolution, see the [CRE Prediction Market Demo](https://github.com/smartcontractkit/cre-gcp-prediction-market-demo/tree/main/cre-workflow).
- **Frontend:** No UI. These are backend workflow templates.
- **Cross-chain:** Single chain only. CCIP integration would be a separate advanced template.
