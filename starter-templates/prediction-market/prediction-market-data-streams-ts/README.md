# Prediction Market Template (Data Streams)

A CRE starter template implementing a full prediction market lifecycle with three workflows: **Market Creation**, **Market Resolution**, and **Dispute Management**. All three workflows share a single `PredictionMarket` smart contract. Price resolution uses **Chainlink Data Streams** — an off-chain, low-latency price feed accessed via the Data Streams REST API.

> **Data Feeds vs Data Streams:** This template fetches prices from the [Data Streams REST API](https://docs.chain.link/data-streams) using HMAC-authenticated HTTP requests. For the on-chain Data Feeds variant, see the sibling `prediction-market-ts` template.

**⚠️ DISCLAIMER**

This template is an educational example to demonstrate how to interact with Chainlink systems, products, and services. It is provided **"AS IS"** and **"AS AVAILABLE"** without warranties of any kind, has **not** been audited, and may omit checks or error handling for clarity. **Do not use this code in production** without performing your own audits and applying best practices. Neither Chainlink Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for unintended outputs generated due to errors in code.

---

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  1. CREATION     │────>│  2. RESOLUTION   │────>│  3. DISPUTE      │
│                  │     │                  │     │                  │
│  Cron trigger    │     │  Cron trigger    │     │  LogTrigger      │
│  Creates new     │     │  Checks expired  │     │  (DisputeRaised) │
│  markets         │     │  markets, calls  │     │  Calls Data      │
│                  │     │  Data Streams    │     │  Streams API,    │
│                  │     │  API, resolves   │     │  resolves dispute│
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
                    │  Chainlink Data Streams  │
                    │  REST API (off-chain)    │
                    │  HMAC-SHA256 auth        │
                    └─────────────────────────┘
```

## Project Structure

```
prediction-market-data-streams-ts/
├── .cre/template.yaml          # Template metadata
├── project.yaml                # CRE project settings (RPCs)
├── secrets.yaml                # Secrets for Data Streams API credentials
├── .env                        # Private key + API credentials (local)
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
├── data-streams/               # Shared Data Streams API client
│   └── client.ts               # HMAC auth, report decoding, price fetch
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

---

## Quick Start

### Prerequisites

- [Bun](https://bun.sh/) (v1.2.21+)
- [CRE CLI](https://docs.chain.link/cre/getting-started)
- Chainlink Data Streams API credentials ([request access](https://chain.link/data-streams))

### 1. Set Up API Credentials

Edit `.env` with your Data Streams credentials:

```bash
CRE_ETH_PRIVATE_KEY=<your-private-key>
DATA_STREAMS_API_KEY=<your-api-key>
DATA_STREAMS_API_SECRET=<your-api-secret>
```

### 2. Install Dependencies

```bash
cd market-creation && bun install && cd ..
cd market-resolution && bun install && cd ..
cd market-dispute && bun install && cd ..
```

### 3. Run Tests

```bash
cd market-creation && bun test && cd ..
cd market-resolution && bun test && cd ..
cd market-dispute && bun test && cd ..
```

### 4. Simulate

```bash
# Create a market
cre workflow simulate market-creation --target staging-settings

# Resolve expired markets (requires Data Streams API credentials)
cre workflow simulate market-resolution --target staging-settings

# Handle a dispute (triggered by on-chain DisputeRaised event)
cre workflow simulate market-dispute --target staging-settings
```

### 5. Broadcast (Sepolia Testnet)

```bash
cre workflow broadcast market-creation --target staging-settings
cre workflow broadcast market-resolution --target staging-settings
cre workflow broadcast market-dispute --target staging-settings
```

---

## Workflows

### 1. Market Creation (Cron Trigger)

Creates a new binary prediction market every hour. The market asks: "Will BTC be above $100,000 by {date}?"

**Config** (`market-creation/config.staging.json`):
- `schedule`: Cron expression (default: hourly)
- `evms[].predictionMarketAddress`: Deployed PredictionMarket contract
- `marketDefaults.question`: Market question template; `{expirationDate}` is replaced at creation
- `marketDefaults.strikePriceUsd`: Price threshold in USD
- `marketDefaults.durationSeconds`: Market duration (default: 24h)

**How it works:**
1. Reads `getNextMarketId()` from the contract
2. Computes strike price in 8 decimal places
3. Encodes a `CREATE` action and writes it on-chain via CRE report

### 2. Market Resolution (Cron Trigger + Data Streams)

Checks markets every 10 minutes. For each resolvable market, fetches the latest BTC/USD price from Chainlink Data Streams and resolves the market.

**Config** (`market-resolution/config.staging.json`):
- `schedule`: Cron expression (default: every 10 min)
- `dataStreams.apiUrl`: Data Streams REST API base URL
- `dataStreams.feedId`: The Data Streams feed ID (BTC/USD)
- `marketIdsToCheck`: Array of market IDs to monitor

**How it works:**
1. Calls `isResolvable(marketId)` on the contract
2. Fetches the latest report from the Data Streams REST API with HMAC authentication
3. Decodes the v3 (Crypto Advanced) report to extract the benchmark `price`
4. Converts from 18 decimal places (Data Streams) to 8 decimal places (on-chain)
5. Encodes a `RESOLVE` action and writes it on-chain via CRE report

### 3. Dispute Management (Log Trigger + Data Streams)

Listens for `DisputeRaised` events. When triggered, fetches a fresh price from Data Streams and re-resolves the market.

**Config** (`market-dispute/config.staging.json`):
- `dataStreams.apiUrl`: Data Streams REST API base URL
- `dataStreams.feedId`: The Data Streams feed ID (BTC/USD)

**How it works:**
1. Decodes `DisputeRaised` event (marketId, disputor, reason)
2. Verifies the market is in `Disputed` status
3. Fetches a fresh price from Data Streams
4. Encodes a `RESOLVE_DISPUTE` action and writes it on-chain

---

## Data Streams Integration

### Authentication

The Data Streams REST API requires HMAC-SHA256 authentication. Three headers are sent with every request:

| Header | Description |
|--------|-------------|
| `Authorization` | Your API key (UUID) |
| `X-Authorization-Timestamp` | Current timestamp in milliseconds |
| `X-Authorization-Signature-SHA256` | HMAC-SHA256 signature of the request |

The HMAC signature is computed over: `METHOD PATH BODY_HASH API_KEY TIMESTAMP`

See [`data-streams/client.ts`](./data-streams/client.ts) for the full implementation.

### Report Schema (v3 — Crypto Advanced)

Data Streams returns a `fullReport` blob that is ABI-decoded in two steps:

1. **Outer wrapper**: `(bytes32[3] reportContext, bytes reportData, bytes32[] rawRs, bytes32[] rawSs, bytes32 rawVs)`
2. **Report body (v3)**: `(bytes32 feedId, uint32 validFromTimestamp, uint32 observationsTimestamp, uint192 nativeFee, uint192 linkFee, uint32 expiresAt, int192 price, int192 bid, int192 ask)`

The `price` field uses **18 decimal places** for crypto feeds. This template converts to 8 decimals to match the on-chain PredictionMarket contract's strike price format.

### Feed IDs

| Feed | Testnet Feed ID | Network |
|------|----------------|---------|
| BTC/USD | `0x00027bbaff688c906a3e20a34fe951715d1018d262a5b66e38edd64e0fdd0d00` | Sepolia |

Find more feed IDs at [docs.chain.link/data-streams](https://docs.chain.link/data-streams).

### API Endpoints

| Environment | Base URL |
|------------|----------|
| Testnet | `https://api.testnet-dataengine.chain.link` |
| Mainnet | `https://api.dataengine.chain.link` |

---

## Smart Contract

The template uses the same `PredictionMarket.sol` contract as the Data Feeds variant. The contract is agnostic to the price source — it receives prices via CRE reports regardless of whether the workflow reads them from on-chain Data Feeds or the off-chain Data Streams API.

### Pre-deployed Contract (Sepolia)

| Component | Address |
|-----------|---------|
| PredictionMarket | `0xEb792aF46AB2c2f1389A774AB806423DB43aA425` |
| MockKeystoneForwarder | `0x15fc6ae953e024d975e77382eeec56a9101f9f88` |

### Constructor Arguments

| Argument | Value (Sepolia) | Description |
|----------|----------------|-------------|
| `forwarder` | `0x15fc6ae953e024d975e77382eeec56a9101f9f88` | MockKeystoneForwarder |
| `_priceFeed` | `0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43` | BTC/USD feed (stored for reference) |
| `_disputeWindow` | `86400` | 24-hour dispute window |

---

## Differences from Data Feeds Template

| Aspect | Data Feeds (`prediction-market-ts`) | Data Streams (this template) |
|--------|-------------------------------------|------------------------------|
| Price source | On-chain `PriceFeedAggregator` contract | Off-chain Data Streams REST API |
| Price read | `priceFeed.latestAnswer(runtime)` | HTTP GET + HMAC auth + ABI decode |
| Decimals | 8 (native) | 18 → converted to 8 |
| Latency | Block-finality dependent | Sub-second updates |
| Auth | None (public on-chain read) | HMAC-SHA256 with API key/secret |
| Capabilities | `chain-read`, `chain-write`, `cron`, `log-trigger` | + `http` |
| Secrets | None | `DATA_STREAMS_API_KEY`, `DATA_STREAMS_API_SECRET` |

---

## Customization

### Different asset
Change `dataStreams.feedId` in the config files and update the market question and strike price accordingly.

### Different duration
Adjust `marketDefaults.durationSeconds` in market-creation config.

### Different resolution frequency
Modify the `schedule` cron expression in market-resolution config.

### Different dispute window
Redeploy the contract with a different `_disputeWindow` constructor argument.

---

## Security Notes

- **Forwarder Validation**: The contract validates that only the trusted Chainlink Forwarder can call `onReport()`
- **HMAC Authentication**: Data Streams API access is secured with HMAC-SHA256 signatures
- **Secrets Management**: API credentials are stored as CRE secrets, never hardcoded in workflow code
- **Consensus**: CRE nodes reach consensus on the HTTP response before writing on-chain
- **Dispute Window**: Prevents disputes after `expirationTime + disputeWindow`

## Known Limitations

- **Spot-price resolution**: Markets resolve using the latest price at resolution time (not historical high/low)
- **Single feed**: All markets use the same BTC/USD feed
- **No betting**: This contract does not implement token deposits or payouts
- **API access required**: Data Streams requires approved API credentials

---

## Alternative Patterns

- **On-chain Data Feeds**: Use the sibling `prediction-market-ts` template for simpler zero-setup price resolution
- **HTTP trigger**: Replace the Cron trigger with an HTTP trigger for on-demand market creation/resolution
- **WebSocket**: For real-time price monitoring, Data Streams also offers a WebSocket API
