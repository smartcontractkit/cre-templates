# Vault Harvester — CRE Starter Template (TypeScript)

Automated DeFi vault harvesting — check if yield is profitable, then harvest and compound.

## Overview

This template demonstrates the **cron -> read -> check profitability -> harvest** pattern using Chainlink CRE (Compute Runtime Environment). Same cron -> read -> check -> write pattern as the Keeper Bot, but with a DeFi-specific vault contract.

### Use Cases

- **Vault harvesting** (primary): Check if `timeSinceLastHarvest > threshold` AND `yield > minThreshold`, then harvest
- **Yield compounding**: Harvest + swap + redeposit accumulated rewards
- **Protocol revenue collection**: Claim and redistribute protocol fees on schedule
- **Target protocols**: Beefy, Yearn, Alchemix, and any vault with a `harvest()` function

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      CRE DON                                │
│                                                             │
│  ┌──────────┐    ┌───────────────┐    ┌──────────────────┐  │
│  │   CRON   │───>│  Read Vault   │───>│  Profitable to   │  │
│  │ Trigger  │    │  State        │    │  harvest?         │  │
│  │ (5 min)  │    │  (EVMClient)  │    │  (time + yield)  │  │
│  └──────────┘    └───────────────┘    └────────┬─────────┘  │
│                                                │            │
│                                     ┌──────────v──────────┐ │
│                                     │  Profitable?        │ │
│                                     │  YES -> Write report│ │
│                                     │  NO  -> Log & skip  │ │
│                                     └──────────┬──────────┘ │
│                                                │            │
└────────────────────────────────────────────────┼────────────┘
                                                 │
                                      ┌──────────v──────────┐
                                      │  KeystoneForwarder  │
                                      │  -> Vault._processReport│
                                      │    (harvest + compound) │
                                      └─────────────────────┘
```

## Components

### CRE Workflow (`my-workflow/`)

The TypeScript workflow runs off-chain inside CRE DON:

1. **Cron trigger** fires every 5 minutes (configurable)
2. **Reads** `shouldHarvest()`, `pendingYield()`, `totalHarvested()`, `harvestCount()` from the on-chain `VaultHarvester` contract
3. **If profitable**: sends a signed report to trigger harvest on-chain
4. **If not profitable**: logs and skips

### Consumer Contract (`contracts/evm/src/VaultHarvester.sol`)

A DeFi vault contract extending `ReceiverTemplate`:

- `shouldHarvest()` — view function that checks if `harvestInterval` has elapsed AND `pendingYield >= minYieldThreshold`
- `_processReport(bytes)` — called by CRE Forwarder, harvests pending yield, resets state, emits event
- `accrueYield(uint256)` — simulates yield accrual (in production, yield comes from the underlying strategy)
- `pendingYield`, `totalHarvested`, `harvestCount`, `lastHarvest` — public state variables

## Getting Started

A demo `VaultHarvester` contract is pre-deployed on Sepolia with simulated yield — this template works out of the box.

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

```bash
# Dry run (no broadcast)
cre workflow simulate my-workflow --target staging-settings

# With actual on-chain transaction
cre workflow simulate my-workflow --target staging-settings --broadcast
```

### 4. Deploy Your Own Contract (Optional)

To use your own contract, deploy `contracts/evm/src/VaultHarvester.sol` to Sepolia using [Remix](https://remix.ethereum.org/) or Foundry. Constructor arguments:

- `forwarder`: CRE KeystoneForwarder address on Sepolia (`0x15fc6ae953e024d975e77382eeec56a9101f9f88`)
- `_harvestInterval`: Seconds between allowed harvests (e.g., `300` for 5 minutes)
- `_minYieldThreshold`: Minimum yield in wei to justify harvest (e.g., `1000000000000000000` for 1 token)

Then update the `contractAddress` in `my-workflow/config.staging.json` with your deployed address.

To simulate yield accrual on your contract:
```bash
cast send <YOUR_CONTRACT> "accrueYield(uint256)" 2000000000000000000 \
  --private-key <YOUR_KEY> --rpc-url https://ethereum-sepolia-rpc.publicnode.com
```

## Customization

- **Change the schedule**: Edit `schedule` in `config.staging.json` (cron syntax)
- **Change the profitability check**: Modify `shouldHarvest()` in `VaultHarvester.sol`
- **Change the harvest logic**: Modify `_processReport()` to implement your vault's harvest + compound strategy

## Migration Guides

### Coming from Gelato?

This template replaces Gelato's vault harvesting Web3 Functions with CRE cron triggers. Gelato shut down Web3 Functions in March 2026.

| Gelato | CRE |
|--------|-----|
| `Web3Function.onRun()` | CRE cron callback (`onCronTrigger`) |
| `multiChainProvider` | `EVMClient` |
| `userArgs` | `config.json` |

### Coming from Chainlink Automation?

This replaces the Beefy/Yearn `checkUpkeep + performUpkeep` pattern. Off-chain profitability check means no wasted gas.

| Chainlink Automation | CRE |
|---------------------|-----|
| `checkUpkeep()` on-chain | `shouldHarvest()` read via `EVMClient` (off-chain) |
| `performUpkeep()` | `_processReport()` via KeystoneForwarder |
| LINK funding | No per-call LINK costs |

## Security

- The `VaultHarvester.sol` contract is a **demo** — audit and customize before production use
- The `ReceiverTemplate` base contract validates that only CRE Forwarder can call `onReport()`
- Never commit `.env` files or secrets

## License

MIT
