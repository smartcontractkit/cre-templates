<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE\&message=Home\&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE\&message=Docs\&color=blue)](https://docs.chain.link/cre)

</div>

# Read & Write Solana - CRE Building Block (TypeScript)

A minimal example that, on a cron schedule, writes a `{key, value}` report to a
Solana program through the **keystone forwarder** — the CRE analogue of an EVM
chain write, but for Solana — and, on a second cron schedule, reads that same
account directly back off-chain to confirm the write landed.

---

**⚠️ DISCLAIMER**

This template is an educational example to demonstrate how to interact with
Chainlink systems, products, and services. It is provided **"AS IS"** and **"AS
AVAILABLE"** without warranties of any kind, has **not** been audited, and may
omit checks or error handling for clarity. **Do not use this code in production**
without performing your own audits and applying best practices. Neither Chainlink
Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for
unintended outputs generated due to errors in code.

---

## Table of Contents

- [What This Does](#what-this-does)
- [Architecture](#architecture)
- [Get Started](#get-started)
  - [1. Update .env](#1-update-env)
  - [2. Install dependencies](#2-install-dependencies)
  - [3. Run a simulation](#3-run-a-simulation)
  - [Deploying your own receiver program (optional)](#deploying-your-own-receiver-program-optional)
- [Security Considerations](#security-considerations)
- [License](#license)

---

## What This Does

The workflow registers two independent cron triggers (same schedule by default —
override the read one with `readSchedule`):

**Write cron** — every 5 minutes:

1. Builds a Borsh-encoded `KvEntry { key, value }` report.
2. Submits it as a CRE report through `SolanaClient.writeReport`.
3. The **keystone forwarder** (already deployed and operated by Chainlink — you
   never deploy or configure it yourself) verifies the report and CPIs into your
   `kv_store_receiver` program's `on_report` instruction.
4. `kv_store_receiver` stores the latest `{key, value}` in its `kv_store` account
   and emits a `KvUpdated` event.
5. The workflow logs the transaction signature and a Solana devnet block-explorer
   link so you can watch it land on-chain.

**Read cron** — same schedule by default:

1. Calls `SolanaClient.getAccountInfoWithOpts` directly — a plain, DON-consensus'd
   RPC read, no forwarder or report involved (the read-side counterpart to the
   write).
2. Decodes the returned account bytes (Anchor discriminator + Borsh `KvStore`
   struct) and logs the current `{key, value, updatedAt, updateCount}`.

## Architecture

```
CRE workflow (main.ts)
  │
  ├─ write cron ──► SolanaClient.writeReport ──► keystone forwarder (Chainlink-operated, devnet)
  │                                                       │  CPI: on_report(metadata, payload)
  │                                                       ▼
  │                                         your kv_store_receiver program
  │                                                       │
  │                                                       ▼
  │                                             kv_store account (on-chain)
  │                                                       ▲
  └─ read cron ──► SolanaClient.getAccountInfoWithOpts ────┘
```

The only piece **you** would deploy yourself is `kv_store_receiver` (see
[`contracts/solana`](./contracts/solana)) — but a `kv_store_receiver` is already
deployed and initialized on Solana devnet for this example, so `config.staging.json`
works as-is. The forwarder is shared, Chainlink-run infrastructure — the same one
every CRE Solana workflow in this environment CPIs through.

## Get Started

### 1. Update .env

Copy `.env.example` to `.env` and set a funded devnet keypair (base58-encoded
secret key) so `cre workflow simulate` can sign and submit the write transaction:

```bash
cp .env.example .env
```

```bash
CRE_SOLANA_PRIVATE_KEY=<your base58 devnet secret key>
```

This key only needs devnet SOL for transaction fees — it is not a DON key and has
no special authority; the forwarder is what authorizes the write.

### 2. Install dependencies

If **Bun** is not already installed, follow the instructions at:
[https://bun.com/docs/installation](https://bun.com/docs/installation)

From your project root:

```bash
bun install --cwd ./my-workflow
```

`my-workflow/config.staging.json` already points at the pre-deployed
`kv_store_receiver` for this example:

```json
{
  "schedule": "0 */5 * * * *",
  "solana": {
    "chainSelectorName": "solana-devnet",
    "receiverProgramId": "aPkfGwVg9Lj3yF2CSgxeRmMn7uU9tKWT6PSHTDf2ZX6",
    "forwarderState": "726bmjnxLeYkvSWUBRby5eUXds96xR58MFfcZoBzuuW3",
    "forwarderAuthority": "DBmhqGAE3eEjbVUUX45npzen8mEVY3HmQHqaUzn7PmYW",
    "receiverAccounts": [
      { "publicKey": "3pAmDKMDtsXtYtJJFHEqqc5UGGkfgFnn76crfNvCZNqX", "isWritable": true }
    ]
  }
}
```

* `schedule` is a 6-field cron expression — the default runs every 5 minutes.
* `forwarderState` is Chainlink's shared Solana-devnet keystone forwarder for the
  staging environment.
* `receiverProgramId`, `forwarderAuthority`, and `receiverAccounts[0].publicKey`
  identify the pre-deployed `kv_store_receiver` program and its `kv_store` account.

No changes are needed to run the golden path — see
[Deploying your own receiver program](#deploying-your-own-receiver-program-optional)
below if you'd rather write to a program you control.

### 3. Run a simulation

From your project root:

```bash
cre workflow simulate my-workflow --target staging-settings
```

`cre workflow simulate` prompts you to pick which trigger to run (or pass
`--trigger-index 0` for the write cron, `--trigger-index 1` for the read cron,
non-interactively). Write, then read, to see the round trip:

```
Workflow compiled
2026-08-19T09:24:27Z [SIMULATION] Simulator Initialized
2026-08-19T09:24:27Z [SIMULATION] Running trigger trigger=cron-trigger@1.0.0
2026-08-19T09:24:28Z [USER LOG] msg="Writing KvEntry {key: \"cre-example\", value: \"hello from CRE @ 2026-08-19T09:24:28.000Z\"} to kv_store_receiver..."
2026-08-19T09:24:30Z [USER LOG] msg="Wrote to Solana devnet, tx=5h1q... explorer=https://explorer.solana.com/tx/5h1q...?cluster=devnet"

Workflow Simulation Result:
 "{\"Key\":\"cre-example\",\"Value\":\"hello from CRE @ 2026-08-19T09:24:28.000Z\",\"TxSignature\":\"5h1q...\",\"ExplorerUrl\":\"https://explorer.solana.com/tx/5h1q...?cluster=devnet\"}"
```

```
2026-08-19T09:25:01Z [SIMULATION] Running trigger trigger=cron-trigger@1.0.0
2026-08-19T09:25:02Z [USER LOG] msg="Read kv_store: key=\"cre-example\" value=\"hello from CRE @ 2026-08-19T09:24:28.000Z\" updatedAt=1755595468 updateCount=1"

Workflow Simulation Result:
 "{\"Found\":true,\"Key\":\"cre-example\",\"Value\":\"hello from CRE @ 2026-08-19T09:24:28.000Z\",\"UpdatedAt\":\"1755595468\",\"UpdateCount\":\"1\"}"
```

Open the printed `explorer.solana.com` link (or the same signature on
[solscan.io](https://solscan.io), passing `?cluster=devnet`) to see the real
transaction and the `KvUpdated` event it emitted.

### Deploying your own receiver program (optional)

The pre-deployed program above is shared and fine for trying this out, but you may
want your own so you're not sharing state with anyone else running this template.
Follow **[contracts/solana/README.md](./contracts/solana/README.md)** — a
step-by-step guide (including Solana/Anchor install steps) for deploying
`kv_store_receiver` to devnet and initializing it against Chainlink's shared
keystone forwarder (`726bmjnxLeYkvSWUBRby5eUXds96xR58MFfcZoBzuuW3` — you do not
deploy the forwarder itself). It ends with four addresses; paste them into
`my-workflow/config.staging.json` in place of the pre-deployed ones.

If you're new to Solana, budget \~20–30 minutes for the toolchain install; the
actual deploy + initialize is a handful of commands.

## Security Considerations

**⚠️ Important Notes**

1. **Demo project** — not production-ready.
2. **Demo contract** — `kv_store_receiver` is intentionally minimal (no access
   control beyond the forwarder-authority check); do not use as-is in production.
3. **Devnet only** — the forwarder address wired up here is Chainlink's staging
   Solana-devnet deployment. Do not reuse it as a mainnet/production address.
4. **Secrets hygiene** — keep real private keys out of version control; use a
   secure secret manager for anything beyond local devnet experimentation.

## License

MIT — see the repository's [LICENSE](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE).
