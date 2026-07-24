# Verifiable Randomness — CRE Starter Template (Go)

VRF-style request/fulfill randomness, powered by a Chainlink CRE workflow. An
on-chain event requests a random number; a CRE workflow generates it with
`runtime.Rand()` and writes a signed result back on-chain.

**⚠️ DISCLAIMER**

This template is an educational example to demonstrate how to interact with Chainlink systems, products, and services. It is provided **"AS IS"** and **"AS AVAILABLE"** without warranties of any kind, has **not** been audited, and may omit checks or error handling for clarity. **Do not use this code in production** without performing your own audits and applying best practices. Neither Chainlink Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for unintended outputs generated due to errors in code.

---

## Overview

This template reproduces the **request → fulfill** shape of Chainlink VRF using CRE:

1. A caller invokes `requestRandomness()` on `RandomnessConsumer`, which emits `RandomnessRequested(requestId, requester)`.
2. A CRE workflow watches that event via a **LogTrigger**, generates a consensus-safe random word, and writes a signed report back to the consumer.
3. `ReceiverTemplate` verifies the report came from the CRE Forwarder, then `_processReport()` records the random word for that `requestId` and emits `RandomnessFulfilled`.

The random draw itself is a few lines (`runtime.Rand()` inside the CRE runtime).
The value of this template is the surrounding request/fulfill scaffolding — the
LogTrigger, the signed write-back, request-ID tracking, and the safety patterns
below — which is what you actually port when moving off VRF.

### How the randomness works

`runtime.Rand()` is the **CRE-specified** randomness source. It returns a
`*math/rand.Rand` seeded by the runtime: in DON mode every node derives the same
per-execution seed, so the network reaches consensus on a single result. No
single node operator can choose or bias the number, and the result is delivered
on-chain only as a DON-signed report accepted by the Forwarder.

> **Consensus-safe, not a VRF-grade CSPRNG.** `runtime.Rand()` is a
> consensus-deterministic PRNG. It suits loot boxes, raffles, sampling, and
> matchmaking — not key generation or signature nonces. Never use Go's global
> `math/rand` or `crypto/rand` in a workflow; they break consensus.
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
│  │ LogTrigger │──>│ Check pending │──>│ runtime.Rand()       │  │
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
4. **Generates** a consensus-safe `randomWord` with `runtime.Rand()` across the full `uint256` range.
5. **Writes** a signed report `(requestId, randomWord)` back to the consumer.

### Contracts

**RandomnessConsumer** (`contracts/evm/src/RandomnessConsumer.sol`):
- `requestRandomness()` — mints a `requestId`, marks it pending, emits `RandomnessRequested`.
- `pendingRequests(uint256)` / `randomWords(uint256)` — request state and fulfilled values.
- `_processReport(bytes)` — decodes `(requestId, randomWord)`, marks fulfilled first, stores the word.

## Getting Started

### Prerequisites

- [Go](https://go.dev/) (see `go.mod` for the required version)
- [CRE CLI](https://docs.chain.link/cre) installed

### 1. Deploy the Consumer

Deploy `RandomnessConsumer` with your chain's KeystoneForwarder address as the
constructor argument (see the [Forwarder Directory](https://docs.chain.link/cre/guides/workflow/using-evm-client/forwarder-directory)):

```bash
forge create src/RandomnessConsumer.sol:RandomnessConsumer --broadcast \
  --private-key <KEY> --rpc-url <RPC> \
  --constructor-args <KEYSTONE_FORWARDER_ADDRESS>
```

Set `consumerAddress` in `my-workflow/config.staging.json` to the deployed address.

### 2. Emit a Request and Simulate

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
- **Random range / shape**: edit the `runtime.Rand()` mapping in `workflow.go` (e.g. `randSource.Int63n(n)` for a bounded draw, or `big.Int.Rand`).
- **Add off-chain context**: call the HTTP capability before generating, e.g. per-user odds or a pity system.

## Security

- The contract is a **demo** — audit and customize before production use.
- `ReceiverTemplate` validates that only the CRE Forwarder can call `onReport()`.
- Never commit `.env` files or secrets.

## License

MIT
