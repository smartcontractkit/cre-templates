<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE\&message=Home\&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE\&message=Docs\&color=blue)](https://docs.chain.link/cre)

</div>

## Quick start

This example shows how to use a **TypeScript** CRE workflow to read a **Multiple Variable Response (MVR)** data feed via the `BundleAggregatorProxy`, decode the bundle, and log/return the structured result.

The example uses the **S&P Global SSA EURC MVR feed on Base mainnet**.

* Production MVR feed overview: [https://data.chain.link/feeds/base/base/sp-global-ssa-eurc](https://data.chain.link/feeds/base/base/sp-global-ssa-eurc)
* Smart contract (BundleAggregatorProxy): [https://basescan.org/address/0xcF99622B5440a338f45daEE134d531A4BE64251F#code](https://basescan.org/address/0xcF99622B5440a338f45daEE134d531A4BE64251F#code)
* Solidity decoding guide (EVM & MVR): [https://docs.chain.link/data-feeds/mvr-feeds/guides/evm-solidity](https://docs.chain.link/data-feeds/mvr-feeds/guides/evm-solidity)

---

### 1) Add the ABI (TypeScript)

Place your `BundleAggregatorProxy` ABI under `contracts/abi` as a `.ts` module and export it as `as const`. Then optionally re-export it from `contracts/abi/index.ts` for clean imports.

```ts
// contracts/abi/BundleAggregatorProxy.ts
import type { Abi } from 'viem';

export const BundleAggregatorProxy = [
  // ... ABI array contents from the contract page ...
] as const;
```

```ts
// contracts/abi/index.ts
export * from './BundleAggregatorProxy';
// add more as needed:
// export * from './IERC20';
```

> You can create additional ABI files the same way (e.g., `IERC20.ts`), exporting them as `as const`.

---

### 2) Configure RPC in `project.yaml`

Add an RPC for the chain you want to read from. For the S&P Global SSA EURC MVR feed on **Base mainnet**:

```yaml
rpcs:
  - chain-name: ethereum-mainnet-base-1
    url: <YOUR_BASE_MAINNET_RPC_URL>
```

---

### 3) Configure the workflow

Create or update `config.json`:

```json
{
  "schedule": "0 */10 * * * *",
  "chainName": "ethereum-mainnet-base-1",
  "feeds": [
    {
      "name": "S&P Global SSA EURC",
      "address": "0xcF99622B5440a338f45daEE134d531A4BE64251F"
    }
  ]
}
```

* `schedule` uses a 6-field cron expression — this runs every 10 minutes at second 0.
* `chainName` must match the RPC entry in `project.yaml`.
* `feeds` is a list of MVR feeds, each pointing to a `BundleAggregatorProxy` contract.

The workflow decodes the bundle into:

* `lastModifiedDateTimeRaw` / `lastModifiedDateTimeRfc3339`
* `securityId`
* `securityName`
* `ssaRaw` / `ssaScaled`
* `ssaDesc`
* `ssaDecimal` (from `bundleDecimals[3]`)

---

### 4) Ensure `workflow.yaml` points to your config

```yaml
staging-settings:
  user-workflow:
    workflow-name: "my-workflow"
  workflow-artifacts:
    workflow-path: "."
    config-path: "./config.json"
    secrets-path: ""
```

---

### 5) Install dependencies

From your project root:

```bash
bun install --cwd ./my-workflow
```

(or use `npm` / `pnpm` if that’s how your project is set up.)

---

### 6) Run a local simulation

From your project root:

```bash
cre workflow simulate my-workflow
```

You should see output similar to:

```text
Workflow compiled
2025-12-08T10:49:19Z [SIMULATION] Simulator Initialized

2025-12-08T10:49:19Z [SIMULATION] Running trigger trigger=cron-trigger@1.0.0
2025-12-08T10:49:20Z [USER LOG] MVR bundle read | chain=ethereum-mainnet-base-1 feed="S&P Global SSA EURC" address=0xcF99622B5440a338f45daEE134d531A4BE64251F lastModifiedDateTimeRaw=1737072000 lastModifiedDateTimeRFC3339=2025-01-17T00:00:00.000Z securityId=SP-EURC securityName=EURC ssaRaw=2 ssaScaled=2 ssaDesc=Strong ssaDecimal=0

Workflow Simulation Result:
 "[\n  {\n    \"name\": \"S\u0026P Global SSA EURC\",\n    \"address\": \"0xcF99622B5440a338f45daEE134d531A4BE64251F\",\n    \"bundle\": {\n      \"lastModifiedDateTimeRaw\": \"1737072000\",\n      \"lastModifiedDateTimeRfc3339\": \"2025-01-17T00:00:00.000Z\",\n      \"securityId\": \"SP-EURC\",\n      \"securityName\": \"EURC\",\n      \"ssaRaw\": \"2\",\n      \"ssaScaled\": \"2\",\n      \"ssaDesc\": \"Strong\",\n      \"ssaDecimal\": 0\n    },\n    \"bundleDecimals\": [\n      0,\n      0,\n      0,\n      0,\n      0\n    ]\n  }\n]"
```

The JSON result is an array of objects like:

```json
{
  "name": "S&P Global SSA EURC",
  "address": "0xcF99622B5440a338f45daEE134d531A4BE64251F",
  "bundle": {
    "lastModifiedDateTimeRaw": "1737072000",
    "lastModifiedDateTimeRfc3339": "2025-01-17T00:00:00.000Z",
    "securityId": "SP-EURC",
    "securityName": "EURC",
    "ssaRaw": "2",
    "ssaScaled": "2",
    "ssaDesc": "Strong",
    "ssaDecimal": 0
  },
  "bundleDecimals": [0, 0, 0, 0, 0]
}
```

This makes it easy for downstream systems to consume the full MVR response from the Chainlink feed.
