# Keeper Bot — CRE Starter Template (TypeScript)

Cron-based smart contract maintenance — read state, check condition, execute if needed.

**⚠️ DISCLAIMER**

This template is an educational example to demonstrate how to interact with Chainlink systems, products, and services. It is provided **"AS IS"** and **"AS AVAILABLE"** without warranties of any kind, has **not** been audited, and may omit checks or error handling for clarity. **Do not use this code in production** without performing your own audits and applying best practices. Neither Chainlink Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for unintended outputs generated due to errors in code.

---

## Overview

This template demonstrates the universal **cron -> read -> check -> write** pattern using the Chainlink CRE (Compute Runtime Environment). A cron trigger fires on a configurable schedule, reads on-chain state to check if upkeep is needed, and conditionally writes back to the contract.

### Use Cases

- **Upkeep tasks**: Any periodic "check and execute" pattern
- **Rebasing tokens**: Call `rebase()` on elastic supply tokens on schedule
- **Oracle updates**: Push price updates at fixed intervals
- **Liquidation checks**: Monitor collateral ratios and trigger liquidations
- **Counter/heartbeat**: Simple liveness proof for monitoring

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      CRE DON                                │
│                                                             │
│  ┌──────────┐    ┌───────────────┐    ┌──────────────────┐  │
│  │   CRON   │───>│  Read State   │───>│  Check Condition │  │
│  │ Trigger  │    │  (EVMClient)  │    │  (if threshold)  │  │
│  │ (5 min)  │    │               │    │                  │  │
│  └──────────┘    └───────────────┘    └────────┬─────────┘  │
│                                                │            │
│                                     ┌──────────v──────────┐ │
│                                     │  Condition met?     │ │
│                                     │  YES -> Write report│ │
│                                     │  NO  -> Log & skip  │ │
│                                     └──────────┬──────────┘ │
│                                                │            │
└────────────────────────────────────────────────┼────────────┘
                                                 │
                                      ┌──────────v──────────┐
                                      │  KeystoneForwarder  │
                                      │  -> Consumer.onReport│
                                      └─────────────────────┘
```

## Components

### CRE Workflow (`my-workflow/`)

The TypeScript workflow runs off-chain inside CRE DON:

1. **Cron trigger** fires every 5 minutes (configurable)
2. **Reads** `needsUpkeep()` from the on-chain `KeeperConsumer` contract
3. **If upkeep is needed**: encodes the new counter value, generates a DON-signed report, and writes it on-chain via the KeystoneForwarder
4. **If not needed**: logs and skips

### Consumer Contract (`contracts/evm/src/KeeperConsumer.sol`)

A simple Solidity contract extending `ReceiverTemplate`:

- `needsUpkeep()` — view function that checks if `interval` seconds have elapsed
- `_processReport(bytes)` — called by CRE Forwarder via ReceiverTemplate, decodes and stores the new counter
- `counter`, `lastExecuted`, `interval` — public state variables

## Getting Started

A demo `KeeperConsumer` contract is pre-deployed on Sepolia — this template works out of the box.

### Prerequisites

- [Bun](https://bun.sh/) runtime installed
- [CRE CLI](https://docs.chain.link/cre) installed

### 1. Install Dependencies

```bash
cd my-workflow && bun install && cd ..
```

### 2. Run Tests

```bash
cd my-workflow && bun test
```

### 3. Simulate

```bash
# Dry run (no broadcast)
cre workflow simulate my-workflow --target staging-settings

# With actual on-chain transaction
cre workflow simulate my-workflow --target staging-settings --broadcast
```

### 4. Deploy Your Own Contract (Optional)

To use your own contract, deploy `contracts/evm/src/KeeperConsumer.sol` to Sepolia using [Remix](https://remix.ethereum.org/) or Foundry. Constructor arguments:

- `forwarder`: CRE KeystoneForwarder address on Sepolia (`0x15fc6ae953e024d975e77382eeec56a9101f9f88`)
- `_interval`: Seconds between allowed executions (e.g., `300` for 5 minutes)

Then update the `contractAddress` in `my-workflow/config.staging.json` with your deployed address.

## Customization

- **Change the schedule**: Edit `schedule` in `config.staging.json` (cron syntax)
- **Change the condition**: Replace `needsUpkeep()` with any view function that returns whether to execute
- **Change the action**: Modify the `_processReport()` logic in `KeeperConsumer.sol` to perform any state mutation

## Migration Guides

### Coming from Gelato?

This template replaces Gelato's `Web3Function.onRun()` with a CRE cron trigger. Gelato shut down Web3 Functions in March 2026 — this is the direct migration path.

| Gelato | CRE |
|--------|-----|
| `Web3Function.onRun()` | CRE cron callback (`onCronTrigger`) |
| `multiChainProvider` | `EVMClient` |
| `userArgs` | `config.json` |
| `Storage.get/set()` | On-chain state or Key-Value Store template |

### Coming from Chainlink Automation?

This template replaces the `checkUpkeep() + performUpkeep()` pattern with off-chain logic + `onReport()`. No more gas overhead for `checkUpkeep` — the condition check runs off-chain in CRE DON.

| Chainlink Automation | CRE |
|---------------------|-----|
| `checkUpkeep()` on-chain | `needsUpkeep()` read via `EVMClient` (off-chain) |
| `performUpkeep()` | `onReport()` via KeystoneForwarder |
| LINK funding | No per-call LINK costs |

## Security

- The `KeeperConsumer.sol` contract is a **demo** — audit and customize before production use
- The `ReceiverTemplate` base contract validates that only CRE Forwarder can call `onReport()`
- Never commit `.env` files or secrets

## License

MIT
