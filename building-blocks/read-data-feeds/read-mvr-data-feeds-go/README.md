<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE\&message=Home\&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE\&message=Docs\&color=blue)](https://docs.chain.link/cre)

</div>

## Quick start

This example shows how to use CRE to read a **Multiple Variable Response (MVR)** data feed via the `BundleAggregatorProxy`, decode the bundle, and log/return the structured result.

The example uses the **S&P Global SSA EURC MVR feed on Base mainnet**.

* Production MVR feed overview: [https://data.chain.link/feeds/base/base/sp-global-ssa-eurc](https://data.chain.link/feeds/base/base/sp-global-ssa-eurc)
* Smart contract (BundleAggregatorProxy): [https://basescan.org/address/0xcF99622B5440a338f45daEE134d531A4BE64251F#code](https://basescan.org/address/0xcF99622B5440a338f45daEE134d531A4BE64251F#code)
* Solidity decoding guide (EVM & MVR): [https://docs.chain.link/data-feeds/mvr-feeds/guides/evm-solidity](https://docs.chain.link/data-feeds/mvr-feeds/guides/evm-solidity)

---

### 1) Add the ABI

Copy the MVR feed’s **BundleAggregatorProxy** ABI into your repo, for example:

```text
contracts/abi/BundleAggregatorProxy.abi
```

### 2) Generate bindings

From your **project root** (where `project.yaml` lives):

```bash
cre generate-bindings evm
```

This creates Go bindings under something like:

```text
contracts/evm/src/generated/bundle_aggregator_proxy/...
```

After generation, if your module picked up new deps, run:

```bash
go mod tidy
```

### 3) Configure RPC in `project.yaml` for Base mainnet

The example MVR feed is on **Base mainnet**. Use the chain name:

```yaml
rpcs:
  - chain-name: ethereum-mainnet-base-1
    url: <YOUR_BASE_MAINNET_RPC_URL>
```

### 4) Configure the workflow

Create/update `config.json` for the **MVR bundle reader** workflow:

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

* `schedule` uses a 6-field cron expression: run on the 0th second every 10 minutes.
* `chainName` must match your `project.yaml` RPC entry.
* `feeds` is an array of MVR feeds (each is a `BundleAggregatorProxy` address on the target chain).

### 5) Run a local simulation

From your project root:

```bash
cre workflow simulate my-workflow
```

You should see output similar to:

```text
Workflow compiled
2025-12-08T10:32:09Z [SIMULATION] Simulator Initialized

2025-12-08T10:32:09Z [SIMULATION] Running trigger trigger=cron-trigger@1.0.0
2025-12-08T10:32:09Z [USER LOG] msg="chain selector obtained" chainName=ethereum-mainnet-base-1
2025-12-08T10:32:09Z [USER LOG] msg="selector details" selector=15971525489660198786
2025-12-08T10:32:11Z [USER LOG] msg="MVR bundle read" chain=ethereum-mainnet-base-1 feed="S&P Global SSA EURC" address=0xcF99622B5440a338f45daEE134d531A4BE64251F lastModifiedDateTimeRaw=1737072000 lastModifiedDateTimeRFC3339=2025-01-17T00:00:00Z securityId=SP-EURC securityName=EURC ssaRaw=2 ssaScaled=2 ssaDesc=Strong ssaDecimal=0

Workflow Simulation Result:
 "[{\"name\":\"S\\u0026P Global SSA EURC\",\"address\":\"0xcF99622B5440a338f45daEE134d531A4BE64251F\",\"bundle\":{\"lastModifiedDateTimeRaw\":\"1737072000\",\"lastModifiedDateTimeRfc3339\":\"2025-01-17T00:00:00Z\",\"securityId\":\"SP-EURC\",\"securityName\":\"EURC\",\"ssaRaw\":\"2\",\"ssaScaled\":\"2\",\"ssaDesc\":\"Strong\",\"ssaDecimal\":0},\"bundleDecimals\":\"AAAAAAA=\"}]"
```

The workflow:

* Calls `latestBundle()` and `bundleDecimals()` on the `BundleAggregatorProxy`.

* Decodes the bundle into:

  ```text
  LastModifiedDateTime : uint256
  SecurityID           : string
  SecurityName         : string
  SSA                  : uint256
  SSADesc              : string
  ```

* Scales `SSA` using the corresponding `bundleDecimals` entry.

* Logs the decoded values and returns them as JSON for easy downstream consumption.
