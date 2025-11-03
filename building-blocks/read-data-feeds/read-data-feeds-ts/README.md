<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE\&message=Home\&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE\&message=Docs\&color=blue)](https://docs.chain.link/cre)

</div>

# Read Data Feeds - CRE Building Blocks (TypeScript)

A minimal TypeScript example that, on a cron schedule (every 10 minutes), reads `decimals()` and `latestAnswer()` from **Chainlink Data Feeds** on **Arbitrum One (mainnet)** using the CRE chain reader. It logs the scaled values and returns a JSON array of results.

**Production contract (BTC/USD on Arbitrum One):**
`0x6ce185860a4963106506C203335A2910413708e9`
Source: [https://arbiscan.io/address/0x6ce185860a4963106506C203335A2910413708e9#code](https://arbiscan.io/address/0x6ce185860a4963106506C203335A2910413708e9#code)

---

## Quick start

### 1) Add the ABI (TypeScript)

Place your ABI under `contracts/abi` as a `.ts` module and export it as `as const`. Then optionally re-export it from `contracts/abi/index.ts` for clean imports.

```ts
// contracts/abi/PriceFeedAggregator.ts
import type { Abi } from 'viem';

export const PriceFeedAggregator = [
  // ... ABI array contents from the contract page ...
] as const;
```

```ts
// contracts/abi/index.ts
export * from './PriceFeedAggregator';
// add more as needed:
// export * from './IERC20';

```

> You can create additional ABI files the same way (e.g., `IERC20.ts`), exporting them as `as const`.

### 2) Configure RPC in `project.yaml`

Add an RPC for the chain you want to read from. For Arbitrum One mainnet:

```yaml
rpcs:
  - chain-name: ethereum-mainnet-arbitrum-1
    url: <YOUR_ARBITRUM_MAINNET_RPC_URL>
```

### 3) Configure the workflow

Create or update `config.json`:

```json
{
  "schedule": "0 */10 * * * *",
  "chainName": "ethereum-mainnet-arbitrum-1",
  "feeds": [
    {
      "name": "BTC/USD",
      "address": "0x6ce185860a4963106506C203335A2910413708e9"
    },
    {
      "name": "ETH/USD",
      "address": "0x639Fe6ab55C921f74e7fac1ee960C0B6293ba612"
    }
  ]
}
```

* `schedule` uses a 6-field cron expression — this runs every 10 minutes at second 0.
* `chainName` must match the RPC entry in `project.yaml`.
* `feeds` is a list of (name, address) pairs to read.

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

### 5) Install dependencies

From your project root:

```bash
bun install --cwd ./my-workflow
```

### 6) Run a local simulation

From your project root:

```bash
cre workflow simulate my-workflow
```

You should see output similar to:

```
Workflow compiled
2025-10-30T09:24:27Z [SIMULATION] Simulator Initialized

2025-10-30T09:24:27Z [SIMULATION] Running trigger trigger=cron-trigger@1.0.0
2025-10-30T09:24:28Z [USER LOG] msg="Data feed read" chain=ethereum-mainnet-arbitrum-1 feed=BTC/USD address=0x6ce185860a4963106506C203335A2910413708e9 decimals=8 latestAnswerRaw=10803231994131 latestAnswerScaled=108032.31994131
2025-10-30T09:24:29Z [USER LOG] msg="Data feed read" chain=ethereum-mainnet-arbitrum-1 feed=ETH/USD address=0x639Fe6ab55C921f74e7fac1ee960C0B6293ba612 decimals=8 latestAnswerRaw=378968000000 latestAnswerScaled=3789.68

Workflow Simulation Result:
 "[{\"name\":\"BTC/USD\",\"address\":\"0x6ce185860a4963106506C203335A2910413708e9\",\"decimals\":8,\"latestAnswerRaw\":\"10803231994131\",\"scaled\":\"108032.31994131\"},{\"name\":\"ETH/USD\",\"address\":\"0x639Fe6ab55C921f74e7fac1ee960C0B6293ba612\",\"decimals\":8,\"latestAnswerRaw\":\"378968000000\",\"scaled\":\"3789.68\"}]"
```