# Off-Chain Report Validation — CRE Starter Template (TypeScript + Web App)

Write a value on-chain through the CRE chain writer, export the DON-signed report, and verify it entirely client-side in a browser app.

**⚠️ DISCLAIMER**

This template is an educational example to demonstrate how to interact with Chainlink systems, products, and services. It is provided **"AS IS"** and **"AS AVAILABLE"** without warranties of any kind, has **not** been audited, and may omit checks or error handling for clarity. **Do not use this code in production** without performing your own audits and applying best practices. Neither Chainlink Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for unintended outputs generated due to errors in code.

---

## Overview

Whenever a CRE workflow writes on-chain, the DON produces a **report**: the encoded
execution result plus metadata, signed by the DON nodes. The `KeystoneForwarder`
verifies those signatures before the report ever reaches your consumer contract.

That same report can also be published off-chain — for example to your own backend —
so that your UI can display it to users as cryptographic proof. This template shows
**both ends** of that flow:

1. **`offchain-report-write/`** — a workflow that writes a `uint256` to a `Storage`
   consumer contract via the chain writer (`writeReport`) and emits the signed report
   in its logs, ready to be copy-pasted or POSTed to your backend.
2. **The web app at this template's root** — a minimal browser app (Vite + TypeScript +
   [viem](https://viem.sh)) that validates such a report with **no backend of your own**:
   it reads the DON configuration directly from the capabilities registry contract over
   a public JSON-RPC endpoint and checks that `f + 1` unique DON nodes signed the report
   hash.

It is the only template that ships a small web app — the app exists because off-chain
report validation is, by definition, something your users do outside the DON.

## Architecture

```
                          CRE DON
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│  ┌──────────────┐   encode value    ┌────────────────────┐  │
│  │ Cron trigger │ ────────────────> │ runtime.report()   │  │
│  └──────────────┘                   │ (DON-signed report)│  │
│                                     └─────────┬──────────┘  │
│                                               │             │
│                                     ┌─────────v──────────┐  │
│                                     │ writeReport        │  │
│                                     │ (chain writer)     │  │
│                                     └─────────┬──────────┘  │
│                                               │             │
│                            runtime.log / HTTP POST (TODO)   │
└───────────────────────────────┬───────────────┬─────────────┘
                                │               │
                ┌───────────────v───┐   ┌───────v──────────────┐
                │ KeystoneForwarder │   │ your backend / logs  │
                │ -> Storage        │   │ (report JSON)        │
                │   ._processReport │   └───────┬──────────────┘
                └───────────────────┘           │ paste / fetch
                                                ▼
                                   ┌────────────────────────────┐
                                   │  Browser validator app     │
                                   │  1. parse report header    │
                                   │  2. keccak report hash     │
                                   │  3. read DON config from   │
                                   │     capabilities registry  │
                                   │     via eth_call (RPC)     │
                                   │  4. recover f+1 signers    │
                                   └────────────────────────────┘
```

## Components

### CRE Workflow (`offchain-report-write/`)

Cron-triggered workflow:

1. ABI-encodes the current timestamp as `uint256`, wrapped in the consumer's
   `onReport(bytes, bytes)` calldata.
2. Generates the DON-signed report with `runtime.report(...)`. Every chain write
   through the chain writer is delivered as a signed report — the report comes with
   the write by default; the workflow generates it explicitly (exactly what the
   generated binding's `writeReportFromOnReport` does internally) so it can keep the
   report object and export it.
3. Submits the write on-chain via `evmClient.writeReport(...)` — the
   `KeystoneForwarder` delivers it to the consumer contract.
4. Logs the report field by field between `REPORT_JSON_BEGIN` / `REPORT_JSON_END`
   markers (one value per line, so nothing is truncated in the logs). The lines
   reassemble into:

```json
{
	"rawReport": "0x...",
	"reportContext": "0x...",
	"signatures": ["0x...", "..."]
}
```

A `TODO`-style comment in `main.ts` marks the hand-off point where you would instead
push the report to your offchain backend with the `HTTPClient` (intentionally not
implemented in this template).

### Consumer Contract (`contracts/evm/src/Storage.sol`)

Minimal `ReceiverTemplate`-based consumer that stores the reported `uint256`.
TypeScript bindings (`contracts/evm/ts/generated/`) are generated from
`src/abi/Storage.abi` with `cre generate-bindings evm` and used by both the workflow
and its unit tests (`newStorageMock`).

### Validator Web App (template root)

- `src/report.ts` — the validator. It mirrors the verification logic of the CRE
  TypeScript SDK, with one swap: the SDK reads the DON configuration through the CRE
  runtime's EVM capability, while the browser app reads it directly on-chain with
  `eth_call` via any JSON-RPC endpoint.
- `src/sample-report.ts` — a real report fixture (produced by the production DON) so
  the app verifies something out of the box.
- `scripts/verify-sample.ts` — CLI sanity check of the same code path.

Validation steps performed in the browser:

1. Parse the 109-byte metadata header of `rawReport` (version, execution ID,
   timestamp, DON ID, DON config version, workflow ID/name/owner, report ID).
2. Compute the report hash: `keccak256(keccak256(rawReport) ++ reportContext)`.
3. Read the DON config from the capabilities registry (production:
   `0x76c9cf548b4179F8901cda1f8623568b58215E62` on Ethereum mainnet):
   `getDON(donId)` returns the fault tolerance `f` and node p2p IDs,
   `getNodesByP2PIds(...)` returns each node's signer address.
4. Recover the signer of every attached ECDSA signature and require at least `f + 1`
   unique signers that are members of the DON.

The RPC URL, chain ID, and registry address are editable in the UI.

## Getting Started

### Prerequisites

- [Bun](https://bun.sh/) (workflow + bindings) and Node.js (web app)
- [CRE CLI](https://docs.chain.link/cre) installed and authenticated

### 1. Install dependencies

```bash
cd contracts && bun install && cd ..
cd offchain-report-write && bun install && cd ..
npm install
```

### 2. Run the validator app

```bash
npm run dev
```

Open the printed URL, click **Load sample report**, then **Verify report**. The bundled
sample is a real report signed by the production DON and validates against live
mainnet registry state. `npm run verify:sample` runs the same check headlessly.

### 3. Simulate the workflow

From this directory (the CRE project root):

```bash
cre workflow simulate offchain-report-write --target staging-settings
```

You will see the `REPORT_JSON_BEGIN ... REPORT_JSON_END` block in the `[USER LOG]`
output. To exercise a real testnet write from simulation, deploy `Storage.sol` with
Sepolia's mock forwarder (`0x15fc6ae953e024d975e77382eeec56a9101f9f88`), set
`evms[0].storageAddress` in `config.staging.json`, add `CRE_ETH_PRIVATE_KEY` to
`.env` (see `.env.example`), and run with `--broadcast`.

> **Important:** simulation signs reports with a local mock DON, so simulated reports
> will **not** validate in the app (you will see every signer flagged
> `unknown-signer`). Reports from a deployed workflow are signed by the real DON and
> validate. Deploy with `cre workflow deploy` — see the operations docs linked below.

### 4. Workflow unit tests

```bash
cd offchain-report-write && bun test
```

## Caveats

> **⚠️ Signer rotation:** validation runs against the **current** on-chain DON
> configuration. DON membership can change over time (key rotation, security
> response, node changes). A report produced **before** such a change will no
> longer validate — its signers will show as `unknown-signer`. To verify a
> historical report you need the DON signer set **as of the time the report was
> generated**: archive the DON config (or a registry state reference, e.g. block
> number) alongside each report, or read the registry at the report's block.
- The default RPC endpoint is a public, rate-limited one — use your own RPC for
  anything beyond a demo.
- The workflow checks `txStatus` of the forwarder transaction; for stricter delivery
  guarantees also inspect `receiverContractExecutionStatus` on the write reply.

## Customization

- **Push reports to your backend**: replace the log lines in `main.ts` with an
  `HTTPClient` POST of the report JSON (the comment marks the spot), and point the
  app's fetch logic at your endpoint instead of pasting.
- **Change the payload**: edit the `encodeAbiParameters` call in `main.ts` and the
  `abi.decode` in `Storage.sol` together, then regenerate the ABI + bindings.
- **Verify on another chain/environment**: override RPC URL, chain ID, and registry
  address in the app UI.

## References

- [Writing data on-chain (chain writer / `writeReport`)](https://docs.chain.link/cre/guides/workflow/using-evm-client/onchain-write/writing-data-onchain)
- [Building consumer contracts (`IReceiver` / `ReceiverTemplate`)](https://docs.chain.link/cre/guides/workflow/using-evm-client/onchain-write/building-consumer-contracts)
- [Forwarder directory (per-network forwarder addresses)](https://docs.chain.link/cre/guides/workflow/using-evm-client/forwarder-directory-ts)
- [Project configuration (TypeScript)](https://docs.chain.link/cre/reference/project-configuration-ts)
- [CLI reference (`simulate`, `deploy`, ...)](https://docs.chain.link/cre/reference/cli/workflow)

## License

MIT
