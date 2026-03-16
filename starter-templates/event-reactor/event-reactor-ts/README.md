# Event Reactor — CRE Starter Template (TypeScript)

Listen for on-chain events, read contract state, and respond on-chain.

## Overview

This template demonstrates the **event -> read -> decide -> write** pattern using Chainlink CRE (Compute Runtime Environment). Unlike the cron-based Keeper Bot and Vault Harvester templates, this one uses a **LogTrigger** — it reacts to on-chain events in real-time.

### Use Cases

- **Large transfer monitoring**: Detect large transfers, check if sender is already flagged, flag on-chain
- **Compliance checks**: Token transfer event -> check compliance rules -> approve or reject
- **Governance reactor**: Proposal created on-chain -> read vote data -> submit result
- **Price deviation response**: Price update event -> check against threshold -> flag deviation
- **Agent response**: Agent emits request event -> CRE reads state -> writes response on-chain

## Architecture

```
┌──────────────┐
│ Monitored    │  emits LargeTransfer(from, to, amount, timestamp)
│ Contract     │──────────────────────────────────┐
└──────────────┘                                   │
                                                   v
┌─────────────────────────────────────────────────────────────┐
│                        CRE DON                              │
│                                                             │
│  ┌──────────────┐   ┌──────────────┐   ┌────────────────┐  │
│  │  LogTrigger  │-->│  Decode Log  │-->│  Check State   │  │
│  │  (event)     │   │  (typed via  │   │  (already      │  │
│  │              │   │   bindings)  │   │   flagged?)    │  │
│  └──────────────┘   └──────────────┘   └───────┬────────┘  │
│                                                 │           │
│                                      ┌──────────v────────┐ │
│                                      │ Decision Logic    │ │
│                                      │ not flagged ->    │ │
│                                      │   flag on-chain   │ │
│                                      │ already flagged ->│ │
│                                      │   skip            │ │
│                                      └──────────┬────────┘ │
└─────────────────────────────────────────────────┼───────────┘
                                                  │
                                       ┌──────────v──────────┐
                                       │ KeystoneForwarder   │
                                       │ -> Reactor._process │
                                       └─────────────────────┘
```

## Components

### CRE Workflow (`my-workflow/`)

The TypeScript workflow runs off-chain inside CRE DON:

1. **LogTrigger** fires when `MonitoredToken` emits a `LargeTransfer` event
2. **Decodes** the event using generated typed bindings (from, to, amount, timestamp)
3. **Reads** `ReactorConsumer` state to check if sender is already flagged
4. **If not flagged**: sends a signed report to flag the address on-chain
5. **If already flagged**: logs and skips

### Contracts

**MonitoredToken** (`contracts/evm/src/MonitoredToken.sol`):
- Emits `LargeTransfer(from, to, amount, timestamp)` events
- `simulateLargeTransfer()` — test function to emit events

**ReactorConsumer** (`contracts/evm/src/ReactorConsumer.sol`):
- `flaggedAddresses(address)` — check if an address is flagged
- `totalFlags()` — total number of flagged addresses
- `_processReport(bytes)` — decodes action type, target, reason and flags the address

## Getting Started

Demo contracts are pre-deployed on Sepolia with a test event already emitted — this template works out of the box.

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

LogTrigger simulation requires a real transaction hash that emitted the event:

```bash
# Use the pre-emitted event
cre workflow simulate my-workflow --target staging-settings \
  --trigger-index 0 \
  --evm-tx-hash 0xe5ac97df3f93c5fdc89690a71d68a8e5ebd35ec8ba9c6cb5d1069f5818222553 \
  --evm-event-index 0
```

### 4. Emit Your Own Event

To generate a fresh `LargeTransfer` event on the pre-deployed demo contract:

```bash
cast send 0x805A04e5C8b2dcb2B26C9e9C9aa12ce34374A35b \
  "simulateLargeTransfer(address,uint256)" \
  0x000000000000000000000000000000000000dEaD 2000000000000000000000 \
  --private-key $CRE_ETH_PRIVATE_KEY \
  --rpc-url https://ethereum-sepolia-rpc.publicnode.com
```

Copy the `transactionHash` from the output and simulate with it:

```bash
cre workflow simulate my-workflow --target staging-settings \
  --trigger-index 0 \
  --evm-tx-hash <YOUR_TX_HASH> \
  --evm-event-index 0
```

### 5. Deploy Your Own Contracts (Optional)

Deploy both contracts to Sepolia using Foundry:

**MonitoredToken** (no constructor args):
```bash
forge create src/MonitoredToken.sol:MonitoredToken --broadcast --private-key <KEY> --rpc-url <RPC>
```

**ReactorConsumer** (forwarder address):
```bash
forge create src/ReactorConsumer.sol:ReactorConsumer --broadcast --private-key <KEY> --rpc-url <RPC> \
  --constructor-args 0x15fc6ae953e024d975e77382eeec56a9101f9f88
```

Then emit an event to get a tx hash for simulation:
```bash
cast send <MONITORED_TOKEN> "simulateLargeTransfer(address,uint256)" 0x000000000000000000000000000000000000dEaD 2000000000000000000000 \
  --private-key <KEY> --rpc-url <RPC>
```

## Customization

- **Change the monitored event**: Modify `MonitoredToken.sol` or point to any contract that emits events
- **Change the reaction logic**: Modify `onLargeTransfer` in `workflow.ts` and `_processReport` in `ReactorConsumer.sol`
- **Add off-chain API calls**: Use `cre.capabilities.HTTPClient` to fetch external data before deciding

## Security

- The contracts are **demos** — audit and customize before production use
- `ReceiverTemplate` validates that only CRE Forwarder can call `onReport()`
- Never commit `.env` files or secrets

## License

MIT
