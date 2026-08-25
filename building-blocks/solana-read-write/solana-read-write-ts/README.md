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

**A note on `cre workflow simulate` and forwarders:** `cre workflow simulate`
(with or without `--broadcast`) always constructs Solana writes through CRE's own
**simulator mock forwarder** — a real, stable devnet program Chainlink runs
specifically for this — never through the real DON forwarder your workflow would
use once deployed. Because of that, this template's `kv_store_receiver` program
has **two independent `kv_store` accounts**, each `initialize`d against a
different forwarder:

| Use case | forwarderState | kv_store account |
|---|---|---|
| `cre workflow simulate` (this README's golden path, incl. `--broadcast`) | simulator mock forwarder | used by `config.staging.json` |
| Real `cre workflow deploy` to a live DON | Chainlink's real staging devnet forwarder | see [Deploying your own receiver program](#deploying-your-own-receiver-program-optional) |

`config.staging.json` ships pointed at the mock-forwarder account so the golden
path below produces a real, finalized devnet transaction out of the box.

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
3. A **keystone forwarder** (already deployed and operated by Chainlink — you
   never deploy or configure it yourself) verifies the report and CPIs into your
   `kv_store_receiver` program's `on_report` instruction.
4. `kv_store_receiver` stores the latest `{key, value}` in its `kv_store` account
   and emits a `KvUpdated` event.
5. The workflow logs the transaction signature and a Solana devnet block-explorer
   link so you can watch it land on-chain.

With `--broadcast`, this is a genuine on-chain write — see the forwarder note
above for which `kv_store` account is used and why.

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
deployed and initialized on Solana devnet for this example (twice — once per
forwarder, see the table above), so `config.staging.json` works as-is. Forwarders
are shared, Chainlink-run infrastructure — you never deploy one yourself.

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
`kv_store_receiver`, initialized against CRE's simulator mock forwarder so
`cre workflow simulate` works out of the box:

```json
{
  "schedule": "0 */5 * * * *",
  "solana": {
    "chainSelectorName": "solana-devnet",
    "receiverProgramId": "aPkfGwVg9Lj3yF2CSgxeRmMn7uU9tKWT6PSHTDf2ZX6",
    "forwarderState": "5Tipz3yhTBdVsDbaBxZkrp7Gjf3brGq5SKkxReefPMP7",
    "forwarderProgramId": "7kuEAA3mSC1Tz8gQjnvH7bKFda9xSPRRin9SZbH49cNK",
    "receiverAccounts": [
      { "publicKey": "7d1kczQb6B2fMcnyCECyLitjvEHt6jLpEc8swPvscQSr", "isWritable": true }
    ]
  }
}
```

* `schedule` is a 6-field cron expression — the default runs every 5 minutes.
* `forwarderState`/`forwarderProgramId` here are **CRE's simulator mock
  forwarder** — the one `cre workflow simulate` actually uses, not the real DON
  forwarder (see the forwarder note above). `forwarderAuthority` is not part of
  the config — the workflow derives it at runtime as the PDA of
  `["forwarder", forwarderState, receiverProgramId]` under `forwarderProgramId`
  (see [Solana Forwarder Directory](https://docs.chain.link/cre/guides/workflow/using-solana-client/solana-forwarder-directory-ts)).
* `receiverProgramId` and `receiverAccounts[0].publicKey` identify the
  pre-deployed `kv_store_receiver` program and the `kv_store` account tied to
  that mock forwarder.

No changes are needed to run the golden path — see
[Deploying your own receiver program](#deploying-your-own-receiver-program-optional)
below if you'd rather write to a program you control, or to point at the real
forwarder for an actual `cre workflow deploy`.

#### Mock forwarder vs. production forwarder

This template's `kv_store_receiver` is `initialize`d twice on devnet — once
against each forwarder — so both configs below are ready to use as-is; you
only need to switch `config.staging.json` between them.

**Mock forwarder** — what `config.staging.json` ships with. Use this whenever
you run `cre workflow simulate` (with or without `--broadcast`), since the CLI
always builds the transaction through this forwarder regardless of what's
configured — see the "note on `cre workflow simulate` and forwarders" near the
top of this README. This is the config to develop and test against locally.

```json
{
  "schedule": "0 */5 * * * *",
  "solana": {
    "chainSelectorName": "solana-devnet",
    "receiverProgramId": "aPkfGwVg9Lj3yF2CSgxeRmMn7uU9tKWT6PSHTDf2ZX6",
    "forwarderState": "5Tipz3yhTBdVsDbaBxZkrp7Gjf3brGq5SKkxReefPMP7",
    "forwarderProgramId": "7kuEAA3mSC1Tz8gQjnvH7bKFda9xSPRRin9SZbH49cNK",
    "receiverAccounts": [
      { "publicKey": "7d1kczQb6B2fMcnyCECyLitjvEHt6jLpEc8swPvscQSr", "isWritable": true }
    ]
  }
}
```

**Production forwarder** — Chainlink's real staging Solana-devnet keystone
forwarder, the one a live DON actually signs through (see the [Solana Forwarder
Directory](https://docs.chain.link/cre/guides/workflow/using-solana-client/solana-forwarder-directory-ts)
for the full, tenant-scoped list). `cre workflow simulate` cannot write
successfully against this config (it'll fail with `MismatchedForwarderProgram`
— expected, not a bug, see above); this config is what you deploy with
`cre workflow deploy` to a real DON.

```json
{
  "schedule": "0 */5 * * * *",
  "solana": {
    "chainSelectorName": "solana-devnet",
    "receiverProgramId": "aPkfGwVg9Lj3yF2CSgxeRmMn7uU9tKWT6PSHTDf2ZX6",
    "forwarderState": "8QoomCQyPSkJ8WopJbX9B4HyvrFzziwvJdU8hZE6DCr9",
    "forwarderProgramId": "CXsKEJcs25TQEYU2e5jZ8QTPE3ffMLZhH6BWHrdcCCB5",
    "receiverAccounts": [
      { "publicKey": "<kvStoreAccount initialized against this forwarder>", "isWritable": true }
    ]
  }
}
```

The `receiverAccounts[0]` value above must come from re-running
[`contracts/solana/scripts/initialize.ts`](./contracts/solana/README.md#6-initialize-your-receiver)
against `forwarderProgramId` `CXsKEJcs25TQEYU2e5jZ8QTPE3ffMLZhH6BWHrdcCCB5` — the
account previously shipped here (`3pAmDKMDtsXtYtJJFHEqqc5UGGkfgFnn76crfNvCZNqX`)
was initialized against a stale, incorrect forwarder program id and will fail
`on_report` with `MismatchedForwarderProgram` against the real Keystone
forwarder.

| | Mock forwarder (default) | Production forwarder |
|---|---|---|
| `cre workflow simulate` write (`--broadcast`) | ✅ succeeds, real finalized tx | ❌ `MismatchedForwarderProgram` (expected) |
| `cre workflow simulate` read | ✅ works | ✅ works |
| Real `cre workflow deploy` to a live DON | n/a — mock forwarder is simulate-only | ✅ this is what a live DON signs through |

Both `kv_store` accounts belong to the same `receiverProgramId`, so once you've
initialized one against each forwarder (see
[Deploying your own receiver program](#deploying-your-own-receiver-program-optional))
you can freely switch between the two JSON blocks above without redeploying
anything.

### 3. Run a simulation

From your project root:

```bash
cre workflow simulate my-workflow --target staging-settings
```

`cre workflow simulate` prompts you to pick which trigger to run (or pass
`--trigger-index 0` for the write cron, `--trigger-index 1` for the read cron,
non-interactively). `--broadcast` additionally requires `CRE_ETH_PRIVATE_KEY` in
`.env` — the CLI's settings loader checks for it regardless of chain family; any
well-formed placeholder works since this workflow never touches EVM (already in
`.env.example`).

**Write cron with `--broadcast`** produces a real, finalized devnet transaction:

```bash
cre workflow simulate my-workflow --target staging-settings --trigger-index 0 --broadcast
```

```
2026-08-19T15:28:47Z [USER LOG] Writing KvEntry {key: "cre-example", value: "hello from CRE @ 2026-08-19T19:28:47.922Z"} to kv_store_receiver...
2026-08-19T15:28:48Z [USER LOG] Wrote to Solana devnet, tx=4JLvLioF1T23BUW5RJSSwBHVtDywYErcd9JvMzuKSPExzXYSiMBHxzMQxzwM3KBm3TMMqPC9cKLeWHVYVe6CwDC explorer=https://explorer.solana.com/tx/4JLvLioF1T23BUW5RJSSwBHVtDywYErcd9JvMzuKSPExzXYSiMBHxzMQxzwM3KBm3TMMqPC9cKLeWHVYVe6CwDC?cluster=devnet

Workflow Simulation Result:
 "{\"Key\":\"cre-example\",\"Value\":\"hello from CRE @ 2026-08-19T19:28:47.922Z\",\"TxSignature\":\"4JLvLioF1T23BUW5RJSSwBHVtDywYErcd9JvMzuKSPExzXYSiMBHxzMQxzwM3KBm3TMMqPC9cKLeWHVYVe6CwDC\",\"ExplorerUrl\":\"https://explorer.solana.com/tx/4JLvLioF1T23BUW5RJSSwBHVtDywYErcd9JvMzuKSPExzXYSiMBHxzMQxzwM3KBm3TMMqPC9cKLeWHVYVe6CwDC?cluster=devnet\"}"
```

Open the printed `explorer.solana.com` link to see the real, finalized
transaction and its `KvUpdated` event. Without `--broadcast`, the same command
dry-runs the transaction (no fees spent, nothing lands on-chain) — useful for
checking your report/account setup before spending devnet SOL.

If you point `config.staging.json` at the **real** DON forwarder instead (see
[Deploying your own receiver program](#deploying-your-own-receiver-program-optional)),
this same command instead fails with `MismatchedForwarderProgram` — expected,
since `cre workflow simulate` always builds the transaction through the mock
forwarder, never a real one. That failure happens *after* correct Borsh
encoding, account ordering, PDA derivation, and CPI dispatch, so it still
confirms everything is wired correctly up to real DON verification; actually
landing a write through the real forwarder requires `cre workflow deploy` to a
live DON (see [CRE docs](https://docs.chain.link/cre)).

**Read cron** (`--trigger-index 1`) is a real, consensus'd RPC call against the
account above:

```
2026-08-19T15:29:28Z [USER LOG] kv_store account 7d1kczQb6B2fMcnyCECyLitjvEHt6jLpEc8swPvscQSr exists but its data was not returned by this read

Workflow Simulation Result:
 "{\"Exists\":true,\"Found\":false,\"Key\":\"\",\"Value\":\"\",\"UpdatedAt\":\"\",\"UpdateCount\":\"\"}"
```

`Exists: true` confirms the account is live on devnet (via `lamports > 0`). As of
the current CRE engine build, `SolanaClient.getAccountInfoWithOpts` reliably
returns `lamports`/`owner`/`space` but doesn't yet populate the account's byte
contents (`data`) inside `cre workflow simulate` — so `Found`/`Key`/`Value` stay
empty even though the on-chain account genuinely has the `{key, value}` the
write cron just stored. `readKvStore` decodes `data` whenever the capability
does return it, so this becomes a full `{key, value, updatedAt, updateCount}`
readout with no code changes once that's populated upstream. You can
independently verify the account's real bytes any time with a plain RPC call:

```bash
curl -s https://api.devnet.solana.com -X POST -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"getAccountInfo","params":["7d1kczQb6B2fMcnyCECyLitjvEHt6jLpEc8swPvscQSr",{"encoding":"base64"}]}'
```

### Deploying your own receiver program (optional)

The pre-deployed program above is shared and fine for trying this out, but you may
want your own so you're not sharing state with anyone else running this template.
Follow **[contracts/solana/README.md](./contracts/solana/README.md)** — a
step-by-step guide (including Solana/Anchor install steps) for deploying
`kv_store_receiver` to devnet and initializing it against Chainlink's shared
keystone forwarder (`forwarderProgramId` `CXsKEJcs25TQEYU2e5jZ8QTPE3ffMLZhH6BWHrdcCCB5`,
`forwarderState` `8QoomCQyPSkJ8WopJbX9B4HyvrFzziwvJdU8hZE6DCr9` on Solana devnet
— you do not deploy the forwarder itself; see the [Solana Forwarder
Directory](https://docs.chain.link/cre/guides/workflow/using-solana-client/solana-forwarder-directory-ts)
for other networks). It ends with three addresses; paste them into
`my-workflow/config.staging.json` in place of the pre-deployed ones.

If you're new to Solana, budget \~20–30 minutes for the toolchain install; the
actual deploy + initialize is a handful of commands.

Note that once your `config.staging.json` points at this real forwarder,
`cre workflow simulate --broadcast` will fail the write with
`MismatchedForwarderProgram` (see the note above) — that's expected. To
actually see `--broadcast` succeed end-to-end while developing, keep a second
`kv_store` account initialized against the mock forwarder (as described above)
and switch back to the real forwarder's account only once you're ready to
`cre workflow deploy`.

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
