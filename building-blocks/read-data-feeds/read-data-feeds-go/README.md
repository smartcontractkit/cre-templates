<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE\&message=Home\&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE\&message=Docs\&color=blue)](https://docs.chain.link/cre)

</div>

# Read Data Feeds - CRE Building Blocks

A minimal example that, on a cron schedule (every 10 minutes), reads `decimals()` and `latestAnswer()` from **Chainlink BTC/USD** Data Feed and **Chainlink ETH/USD** Data Feed on **Arbitrum One (mainnet)** and logs the scaled value.

**Production contract (Arbitrum One):**
`0x6ce185860a4963106506C203335A2910413708e9`
Source: [https://arbiscan.io/address/0x6ce185860a4963106506C203335A2910413708e9#code](https://arbiscan.io/address/0x6ce185860a4963106506C203335A2910413708e9#code)

---

## To read from another feed / chain

### 1) Add the ABI

Copy the feed’s ABI into your repo at:

```
contract/abi/PriceFeedAggregator.abi
```

### 2) Generate bindings

From your **project root** (where `project.yaml` lives):

```bash
cre generate-bindings evm
```

This creates Go bindings under something like:

```
contracts/evm/src/generated/price_feed_aggregator/...
```

> After generation, if your module picked up new deps, run:

```bash
go mod tidy
```

### 3) Configure RPC in `project.yaml` for your selected chain

The example data feed uses Arbitrum mainnet. Use the chain name shown below.

```yaml
rpcs:
  - chain-name: ethereum-mainnet-arbitrum-1
    url: <YOUR_ARBITRUM_MAINNET_RPC_URL>
```

### 4) Configure the workflow

Create/update `config.json` for the simple reader workflow:

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

* `schedule` uses a 6-field cron expression: run on the 0th second every 10 minutes.
* `chainName` must match your `project.yaml` RPC entry.
* `feeds` is an array of price pair and associated address.

### 5) Run a local simulation

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
