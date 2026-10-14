<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE\&message=Home\&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE\&message=Docs\&color=blue)](https://docs.chain.link/cre)

</div>

# Tenderly Preflight (TypeScript)

Simulate an EVM transaction on a [Tenderly Virtual Environment](https://docs.tenderly.co/virtual-environments) from inside a CRE workflow, and get a **consensus-verified verdict** back — before your workflow commits to an on-chain write.

Each DON node forks the chain at the same agreed block, mines the transaction on its own fork, and reads the receipt. Nothing reaches the real chain. The workflow then switches on the verdict:

| `verdict.outcome` | Meaning | This template does |
|---|---|---|
| `success` | The transaction would succeed | Logs the gas used (a real workflow would write on-chain here) |
| `reverted` / `rejected` | It would fail on-chain | Skips the write, logs the reason |
| `misconfigured` | Config/secret/funding can never work | Skips; alert on this one — every run fails until fixed |
| `unavailable` / `oversized` / `indeterminate` | No verdict can be trusted | Skips — fail closed rather than writing blind |

Powered by [`@tenderly/cre-preflight`](https://www.npmjs.com/package/@tenderly/cre-preflight). See its README for the full configuration surface, funding/state overrides, and cost model.

---

## Prerequisites

- [CRE CLI](https://docs.chain.link/cre/getting-started) installed and logged in
- **Bun 1.4 or newer** — older versions silently produce a WASM binary that traps at handler registration (`wasm trap: unreachable`). Check with `bun --version`.
- A [Tenderly](https://tenderly.co) account with a project, and an **access key scoped to that one project**

## Quick start

### 1) Install dependencies

From this directory:

```bash
bun install --cwd ./my-workflow
```

### 2) Add your Tenderly access key

The key is a CRE **secret**, never config or code. Put it in `.env` at the project root:

```bash
TENDERLY_ACCESS_KEY=<your-tenderly-access-key>
```

`secrets.yaml` maps it to the secret name the workflow asks for:

```yaml
secretsNames:
  tenderlyaccesskey:
    - TENDERLY_ACCESS_KEY
```

### 3) Configure the workflow

Edit `my-workflow/config.staging.json`:

```json
{
  "schedule": "0 */5 * * * *",
  "tenderly": {
    "accountSlug": "my-account",
    "projectSlug": "my-project",
    "accessKeySecretId": "tenderlyaccesskey",
    "explainReverts": true,
    "fork": {
      "networkId": "11155111",
      "at": "latest"
    }
  },
  "transaction": {
    "from": "0x1111111111111111111111111111111111111111",
    "to": "0x2222222222222222222222222222222222222222",
    "data": "0x",
    "value": "0x0",
    "gas": "100000"
  }
}
```

* `accountSlug` / `projectSlug` — **replace with yours** (visible in your Tenderly dashboard URL).
* `fork.networkId` — chain id as a decimal string (`11155111` is Sepolia). Check [Tenderly's supported networks](https://docs.tenderly.co/supported-networks-and-languages); an unsupported one fails with `misconfigured`.
* `fork.at` — `"latest"`, `"finalized"`, or a pinned block number. A tag is resolved **once per execution** via a consensus-verified chain read (this is what the RPC in `project.yaml` is for), so every node forks at the same height.
* `transaction` — a placeholder so the template runs out of the box. In a real workflow you build this per execution and pass it as an argument to `sendTransaction`. **The shipped placeholder will revert or be rejected** — the `from` address has no funds — which still demonstrates the verdict flow.

### 4) Run a local simulation

From this directory:

```bash
cre workflow simulate my-workflow
```

You should see output similar to:

```
Workflow compiled
2025-10-30T09:24:27Z [SIMULATION] Simulator Initialized

2025-10-30T09:24:27Z [SIMULATION] Running trigger trigger=cron-trigger@1.0.0
2025-10-30T09:24:29Z [USER LOG] msg="tenderly: environment <id> tx 0x..."
2025-10-30T09:24:29Z [USER LOG] msg="skipping write, would fail | outcome=rejected reason=insufficient funds for gas ..."

Workflow Simulation Result:
 "skipped: rejected"
```

Point `transaction` at a call that succeeds (and give `from` a balance via the `fund` option) and you will see `preflight passed | gasUsed=...` instead.

### 5) Run the tests

```bash
bun test --cwd ./my-workflow
```

---

## Using this in your own workflow

The transaction is an argument, not config, because it usually depends on what the execution just worked out:

```ts
import { TenderlyPreflight } from '@tenderly/cre-preflight';

const tenderly = new TenderlyPreflight(runtime.config.tenderly);
const verdict = tenderly.sendTransaction(runtime, { from, to, data });

if (verdict.outcome !== 'success') {
  runtime.log(`skipping write: ${verdict.outcome} ${verdict.reason}`);
  return 'skipped';
}
// ... safe to write on-chain: runtime.report(...) + evmClient.writeReport(...)
```

Useful per-call options:

```ts
tenderly.sendTransaction(runtime, tx, {
  // fund accounts before the send (native funding is free; each ERC20 entry
  // costs one HTTP action)
  fund: [{ addresses: [sender], balance: 10n ** 18n }],
  // or fork from a block your workflow already knows (e.g. from an EVM log
  // trigger) instead of resolving `fork.at`
  forkBlockNumber: blockNumber,
});
```

**Note on funding:** the sender's real balance is part of what a preflight answers. Topping the sender up turns "would this succeed if I sent it now" into an easier question — fund deliberately.

## Gotchas

- **`misconfigured` is not a retryable error.** It means the config, the Vault secret, or the options are wrong; fix them rather than retry.
- **Every node builds its own environment.** `deleteEnvironment: true` (the default) cleans them up; turn it off only to inspect a fork in the Tenderly dashboard, then turn it back off.
- **Budget HTTP actions.** Default config costs 5 HTTP actions per node per execution (`tenderly.httpActionCost` reports it). Resolving a `fork.at` tag additionally costs one EVM chain read.
- **Time-dependent contracts can come back `indeterminate`.** Nodes execute at slightly different `block.timestamp` values; if gas or a revert reason embeds the timestamp, set `explainReverts: false` and decide on `outcome` alone.

## Reference

- [@tenderly/cre-preflight on npm](https://www.npmjs.com/package/@tenderly/cre-preflight) / [GitHub](https://github.com/Tenderly/tenderly-cre-preflight)
- [Tenderly Virtual Environments](https://docs.tenderly.co/virtual-environments)
- [CRE Documentation](https://docs.chain.link/cre)
