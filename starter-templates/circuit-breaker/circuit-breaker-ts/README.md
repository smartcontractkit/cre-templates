# Circuit Breaker — CRE Starter Template (TypeScript)

Monitor on-chain events for anomalies, automatically pause contracts when thresholds are breached.

**⚠️ DISCLAIMER**

This template is an educational example to demonstrate how to interact with Chainlink systems, products, and services. It is provided **"AS IS"** and **"AS AVAILABLE"** without warranties of any kind, has **not** been audited, and may omit checks or error handling for clarity. **Do not use this code in production** without performing your own audits and applying best practices. Neither Chainlink Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for unintended outputs generated due to errors in code.

---

## Overview

This template demonstrates the **dual-trigger circuit breaker** pattern using Chainlink CRE (Compute Runtime Environment). It combines a **LogTrigger** (react to price update events) with a **Cron trigger** (periodic health checks) to provide comprehensive protocol safety monitoring.

### Use Cases

- **DeFi protocol safety**: Price deviation > threshold -> pause deposits
- **Price deviation monitoring**: Oracle update deviates > 10% from last known -> pause trading
- **Bridge anomaly**: Unusual bridge volume -> pause bridge contract
- **Treasury drain**: Large unexpected outflows -> freeze treasury
- **NFT mint abuse**: Mint rate exceeds threshold -> pause minting

## Architecture

```
                    ┌──────────────────────┐
                    │  Protocol Contract   │
                    │  emits PriceUpdated  │
                    └──────────┬───────────┘
                               │ events
                               v
┌─────────────────────────────────────────────────────────────────┐
│                          CRE DON                                │
│                                                                 │
│  ┌──────────────┐   ┌──────────────┐   ┌──────────────────┐    │
│  │  LogTrigger  │-->│  Decode      │-->│  Anomaly Check   │    │
│  │  (events)    │   │  Event Data  │   │  deviation > X%? │    │
│  └──────────────┘   └──────────────┘   └────────┬─────────┘    │
│                                                  │              │
│  ┌──────────────┐                                │              │
│  │  Cron        │--> Health check ───────────────┤              │
│  │  (periodic)  │    (read state, verify ok)     │              │
│  └──────────────┘                                │              │
│                                        ┌─────────v──────────┐  │
│                                        │   ANOMALY?         │  │
│                                        │   YES -> Pause     │  │
│                                        │   NO  -> Log OK    │  │
│                                        └─────────┬──────────┘  │
│                                                  │              │
└──────────────────────────────────────────────────┼──────────────┘
                                                   │
                                        ┌──────────v──────────┐
                                        │  KeystoneForwarder  │
                                        │  -> _processReport  │
                                        │    (pause protocol) │
                                        └─────────────────────┘
```

## Components

### CRE Workflow (`my-workflow/`)

Two triggers work together:

**Trigger 0 — LogTrigger (event-driven):**
1. Fires when `ProtocolWithBreaker` emits a `PriceUpdated` event
2. Decodes the event (newPrice, oldPrice, timestamp) via generated typed bindings
3. Reads on-chain state (paused?, tripCount, threshold)
4. Calculates price deviation in basis points
5. If deviation > threshold: trips the circuit breaker (pauses protocol on-chain)

**Trigger 1 — Cron (periodic health check):**
1. Fires every 10 minutes (configurable)
2. Reads protocol state (paused, lastPrice, lastUpdate, tripCount)
3. Logs the current status for monitoring

### Consumer Contract (`contracts/evm/src/ProtocolWithBreaker.sol`)

A protocol contract with built-in circuit breaker:

- `updatePrice(uint256)` — simulates price updates, emits `PriceUpdated` events
- `paused()` — whether the protocol is currently paused
- `lastPrice()`, `lastPriceTimestamp()` — current price state
- `priceDeviationThresholdBps()` — deviation threshold in basis points (1000 = 10%)
- `tripCount()` — how many times the breaker has been tripped
- `_processReport(bytes)` — decodes action (PAUSE/UNPAUSE) and reason, updates state

## Getting Started

A demo `ProtocolWithBreaker` contract is pre-deployed on Sepolia with an anomalous price event already emitted — this template works out of the box.

### Prerequisites

- [Bun](https://bun.sh/) runtime installed
- [CRE CLI](https://docs.chain.link/cre) installed

### 1. Install Dependencies

```bash
cd my-workflow && bun install && cd ..
cd contracts && bun install && cd ..
```

### 2. Run Tests

```bash
cd my-workflow && bun test
```

### 3. Simulate

**Cron health check (trigger index 1):**
```bash
cre workflow simulate my-workflow --target staging-settings
```

**LogTrigger with pre-emitted anomalous price event (trigger index 0):**
```bash
cre workflow simulate my-workflow --target staging-settings --non-interactive \
  --trigger-index 0 \
  --evm-tx-hash 0x98dfb56db57cb0689e5e4092949d41431a283f52128066bc941b0c219d3f6203 \
  --evm-event-index 0
```

### 4. Emit Your Own Price Events

Emit a normal price update (within 10% threshold):
```bash
cast send 0xaCb13C9940cB61367b45eEd504E410D4B4d7A6e4 \
  "updatePrice(uint256)" 950000000000000000 \
  --private-key $CRE_ETH_PRIVATE_KEY \
  --rpc-url https://ethereum-sepolia-rpc.publicnode.com
```

Emit an anomalous price update (>10% deviation):
```bash
cast send 0xaCb13C9940cB61367b45eEd504E410D4B4d7A6e4 \
  "updatePrice(uint256)" 100000000000000000 \
  --private-key $CRE_ETH_PRIVATE_KEY \
  --rpc-url https://ethereum-sepolia-rpc.publicnode.com
```

Then simulate with the resulting tx hash:
```bash
cre workflow simulate my-workflow --target staging-settings --non-interactive \
  --trigger-index 0 \
  --evm-tx-hash <YOUR_TX_HASH> \
  --evm-event-index 0
```

### 5. Deploy Your Own Contract (Optional)

Deploy `contracts/evm/src/ProtocolWithBreaker.sol` using Foundry:

```bash
forge create src/ProtocolWithBreaker.sol:ProtocolWithBreaker --broadcast \
  --private-key <KEY> --rpc-url <RPC> \
  --constructor-args 0x15fc6ae953e024d975e77382eeec56a9101f9f88 1000000000000000000 1000
```

Constructor arguments:
- `forwarder`: CRE KeystoneForwarder on Sepolia (`0x15fc6ae953e024d975e77382eeec56a9101f9f88`)
- `_initialPrice`: Initial price in wei (e.g., `1000000000000000000` = 1e18)
- `_deviationThresholdBps`: Deviation threshold in basis points (`1000` = 10%)

## Customization

- **Change the threshold**: Deploy with a different `_deviationThresholdBps` value
- **Change the schedule**: Edit `schedule` in `config.staging.json`
- **Add webhook alerts**: Extend `tripCircuitBreaker()` with `cre.capabilities.HTTPClient` POST to Slack/PagerDuty
- **Monitor different events**: Change `logTriggerPriceUpdated()` to monitor other contract events
- **Add unpause logic**: Create a separate workflow or cron handler that calls ACTION_UNPAUSE when conditions normalize

## Security

- The contract is a **demo** — audit and customize before production use
- `ReceiverTemplate` validates that only CRE Forwarder can call `onReport()`
- The circuit breaker is **one-way by default** — only CRE can unpause via ACTION_UNPAUSE report
- Never commit `.env` files or secrets

## License

MIT
