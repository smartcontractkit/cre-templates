# Verifiable Randomness — CRE Starter Template (TypeScript)

VRF-style request/fulfill randomness, powered by a Chainlink CRE workflow. An
on-chain event requests a random number; a CRE workflow generates it and writes
a signed result back on-chain.

**⚠️ DISCLAIMER**

This template is an educational example to demonstrate how to interact with Chainlink systems, products, and services. It is provided **"AS IS"** and **"AS AVAILABLE"** without warranties of any kind, has **not** been audited, and may omit checks or error handling for clarity. **Do not use this code in production** without performing your own audits and applying best practices. Neither Chainlink Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for unintended outputs generated due to errors in code.

---

## Overview

This template reproduces the **request → fulfill** shape of Chainlink VRF using CRE:

1. A caller invokes `requestRandomness()` on `RandomnessConsumer`, which emits `RandomnessRequested(requestId, requester)`.
2. A CRE workflow watches that event via a **LogTrigger**, generates a consensus-safe random word, and writes a signed report back to the consumer.
3. `ReceiverTemplate` verifies the report came from the CRE Forwarder, then `_processReport()` records the random word for that `requestId` and emits `RandomnessFulfilled`.

The random draw itself is one line (`Math.random()` inside the CRE runtime). The
value of this template is the surrounding request/fulfill scaffolding — the
LogTrigger, the signed write-back, request-ID tracking, and the safety patterns
below — which is what you actually port when moving off VRF.

### How the randomness works

Inside the CRE WASM runtime, `Math.random()` is overridden with a DON-seeded
generator: every node in the DON derives the same per-execution seed and produces
the same value, so the network reaches consensus on a single result. No single
node operator can choose or bias the number, and the result is delivered on-chain
only as a DON-signed report accepted by the Forwarder.

> **Not cryptographically secure.** `Math.random()` in the runtime is a seeded
> pseudo-random generator for consensus, not a CSPRNG. It suits loot boxes,
> raffles, sampling, and matchmaking — not key generation or signature nonces.
> See [Using Randomness in Workflows](https://docs.chain.link/cre/guides/workflow/using-randomness).

## Architecture

```
┌────────────────────┐  requestRandomness()
│  Caller / dApp      │───────────────────────────┐
└────────────────────┘                            │
                                                   v
┌────────────────────┐  emits RandomnessRequested(requestId, requester)
│ RandomnessConsumer  │──────────────────────────────────┐
└────────────────────┘                                   │
                                                          v
┌───────────────────────────────────────────────────────────────┐
│                            CRE DON                              │
│  ┌────────────┐   ┌──────────────┐   ┌──────────────────────┐  │
│  │ LogTrigger │──>│ Check pending │──>│ Math.random()        │  │
│  │ (event)    │   │ (idempotency) │   │ (DON-seeded, agreed) │  │
│  └────────────┘   └──────────────┘   └──────────┬───────────┘  │
│                                                  │              │
│                                      ┌───────────v───────────┐  │
│                                      │ encode (id, word) →   │  │
│                                      │ signed report         │  │
│                                      └───────────┬───────────┘  │
└──────────────────────────────────────────────────┼─────────────┘
                                                    v
                                      ┌─────────────────────────┐
                                      │ KeystoneForwarder       │
                                      │ -> Consumer._process    │
                                      │    (records randomWord) │
                                      └─────────────────────────┘
```

## Components

### CRE Workflow (`my-workflow/`)

1. **LogTrigger** fires when `RandomnessConsumer` emits `RandomnessRequested`.
2. **Decodes** the event via generated typed bindings (`requestId`, `requester`).
3. **Reads** `pendingRequests(requestId)` from a finalized block to skip unknown or already-fulfilled requests.
4. **Generates** a consensus-safe `randomWord` with `Math.random()`.
5. **Writes** a signed report `(requestId, randomWord)` back to the consumer.

### Contracts

**RandomnessConsumer** (`contracts/evm/src/RandomnessConsumer.sol`):
- `requestRandomness()` — mints a `requestId`, marks it pending, emits `RandomnessRequested`.
- `pendingRequests(uint256)` / `randomWords(uint256)` — request state and fulfilled values.
- `_processReport(bytes)` — decodes `(requestId, randomWord)`, marks fulfilled first, stores the word.

## Getting Started

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

### 3. Deploy the Consumer

Deploy `RandomnessConsumer` with your chain's KeystoneForwarder address as the
constructor argument (see the [Forwarder Directory](https://docs.chain.link/cre/guides/workflow/using-evm-client/forwarder-directory-ts)):

```bash
forge create src/RandomnessConsumer.sol:RandomnessConsumer --broadcast \
  --private-key <KEY> --rpc-url <RPC> \
  --constructor-args <KEYSTONE_FORWARDER_ADDRESS>
```

Set `consumerAddress` in `my-workflow/config.staging.json` to the deployed address.

### 4. Emit a Request and Simulate

LogTrigger simulation requires the transaction hash that emitted the event:

```bash
cast send <CONSUMER_ADDRESS> "requestRandomness()" \
  --private-key $CRE_ETH_PRIVATE_KEY \
  --rpc-url https://ethereum-sepolia-rpc.publicnode.com

# Copy the transactionHash from the output, then:
cre workflow simulate my-workflow --target staging-settings --non-interactive \
  --trigger-index 0 \
  --evm-tx-hash <YOUR_TX_HASH> \
  --evm-event-index 0
```

## Security patterns (carry over from VRF)

These protect any outcome that depends on a value known only after the request:

- **Commit before reveal** — lock in user choices before the draw; reveal only after fulfillment.
- **Separate request from consumption** — never request and consume randomness in the same transaction.
- **Track request IDs** — the consumer marks each `requestId` fulfilled **before** any downstream logic, closing re-entrancy and reroll/replay windows.
- **Trigger choice matters** — this template uses a LogTrigger, whose seed derives from on-chain event data no requester controls. Avoid HTTP triggers with caller-chosen request IDs for high-value draws, or add commit-reveal.
- **One value per independent outcome** — request a separate word per independent outcome rather than deriving several from one.

## Customization

- **Reveal logic**: edit `_processReport()` in `RandomnessConsumer.sol`.
- **Random range / shape**: edit the `Math.random()` mapping in `workflow.ts`.
- **Add off-chain context**: call `cre.capabilities.HTTPClient` before generating, e.g. per-user odds or a pity system.

## Security

- The contract is a **demo** — audit and customize before production use.
- `ReceiverTemplate` validates that only the CRE Forwarder can call `onReport()`.
- Never commit `.env` files or secrets.

## License

MIT
