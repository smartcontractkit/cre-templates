<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE&message=Home&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE&message=Docs&color=blue)](https://docs.chain.link/cre)

</div>

# Custom Data Feed – CRE Template

A template for bringing your own **custom off-chain data feed** on-chain with the **Chainlink Runtime Environment (CRE)**.

---

**⚠️ DISCLAIMER**

This template is an educational example to demonstrate how to interact with Chainlink systems, products, and services. It is provided **“AS IS”** and **“AS AVAILABLE”** without warranties of any kind, has **not** been audited, and may omit checks or error handling for clarity. **Do not use this code in production** without performing your own audits and applying best practices. Neither Chainlink Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for unintended outputs generated due to errors in code.

---

## Table of Contents

* [What This Template Does](#what-this-template-does)
* [Getting Started](#getting-started)
  * [1. Update .env file](#1-update-env-file)
  * [2. Configure RPC endpoints](#2-configure-rpc-endpoints)
  * [3. Deploy contracts](#3-deploy-contracts)
  * [4. Generate contract bindings](#4-generate-contract-bindings)
  * [5. Configure workflow](#5-configure-workflow)
  * [6. Simulate the workflow](#6-simulate-the-workflow)
* [Security Considerations](#security-considerations)
* [License](#license)

---

## What This Template Does

This template provides an end-to-end starting point for bringing your own **custom data feed** on-chain with the **Chainlink Runtime Environment (CRE)**. It showcases local simulation and the core CRE workflow patterns.

**Components:**
- **Contracts** (Solidity) under `projectRoot/contracts/evm/src`  
  Example demo contracts used by the workflow:
  - `ReserveManager`
  - `SimpleERC20`
  - `BalanceReader`
  - `MessageEmitter`
- **CRE Workflow** (Go) that fetches your off-chain data and optionally performs chain writes based on configurable triggers (cron or EVM log).

**Key Technologies**
- **CRE (Chainlink Runtime Environment)** – orchestrates workflows with DON consensus.

---

## Getting Started

### 1. Update .env file

Add a private key to `.env`. This is required if your workflow performs **chain writes** (must be valid and funded).
If your workflow does **not** write on-chain, you can use a dummy key:

```bash
CRE_ETH_PRIVATE_KEY=0000000000000000000000000000000000000000000000000000000000000001
```

### 2. Configure RPC endpoints

Specify RPC endpoints for chains you interact with in `project.yaml`. This is needed to submit transactions and read state.

Supported for local simulation (testnet/mainnet variants):

* Ethereum (`ethereum-testnet-sepolia`, `ethereum-mainnet`)
* Base (`ethereum-testnet-sepolia-base-1`, `ethereum-mainnet-base-1`)
* Avalanche (`avalanche-testnet-fuji`, `avalanche-mainnet`)
* Polygon (`polygon-testnet-amoy`, `polygon-mainnet`)
* BNB Chain (`binance-smart-chain-testnet`, `binance-smart-chain-mainnet`)
* Arbitrum (`ethereum-testnet-sepolia-arbitrum-1`, `ethereum-mainnet-arbitrum-1`)
* Optimism (`ethereum-testnet-sepolia-optimism-1`, `ethereum-mainnet-optimism-1`)

Add your RPC under `rpcs`:

```yaml
rpcs:
  - chain-name: ethereum-testnet-sepolia
    url: <Your RPC endpoint to ETH Sepolia>
```

> For chain names, see the selectors list: [https://github.com/smartcontractkit/chain-selectors/blob/main/selectors.yml](https://github.com/smartcontractkit/chain-selectors/blob/main/selectors.yml)

### 3. Deploy contracts

Deploy the demo contracts: **BalanceReader**, **MessageEmitter**, **ReserveManager**, **SimpleERC20** (locally or to a testnet using tools like Foundry).

For quick trials, you can use the pre-deployed addresses on **Ethereum Sepolia**:

* chain: `ethereum-testnet-sepolia`
* `ReserveManager`: `0x073671aE6EAa2468c203fDE3a79dEe0836adF032`
* `SimpleERC20`: `0x4700A50d858Cb281847ca4Ee0938F80DEfB3F1dd`
* `BalanceReader`: `0x4b0739c94C1389B55481cb7506c62430cA7211Cf`
* `MessageEmitter`: `0x1d598672486ecB50685Da5497390571Ac4E93FDc`

Source code lives under `projectRoot/contracts/evm/src`.

### 4. Generate contract bindings

Generate Go bindings from ABIs (located in `projectRoot/contracts/src/abi`).
Run this **from your project root** (where `project.yaml` is):

```bash
# Generate bindings for all contracts
cre generate-bindings evm
```

Bindings will appear under `contracts/evm/src/generated/`, e.g.:

* `contracts/evm/src/generated/ierc20/IERC20.go`
* `contracts/evm/src/generated/reserve_manager/ReserveManager.go`
* `contracts/evm/src/generated/balance_reader/BalanceReader.go`
* etc.

(Template bindings are pre-generated; re-run only if ABIs/contracts change.)

### 5. Configure workflow

Edit `config.json`:

* `schedule`: e.g. `"*/3 * * * * *"` (every 3 seconds) or your preferred cron
* `url`: your off-chain data endpoint (custom data feed)
* `tokenAddress`: `SimpleERC20` address
* `reserveManagerAddress`: `ReserveManager` address
* `balanceReaderAddress`: `BalanceReader` address
* `messageEmitterAddress`: `MessageEmitter` address
* `chainName`: selected chain name (see selectors YAML linked above)
* `gasLimit`: gas limit used for chain writes

Ensure `workflow.yaml` points to the config:

```yaml
staging-settings:
  user-workflow:
    workflow-name: "my-workflow"
  workflow-artifacts:
    workflow-path: "."
    config-path: "./config.json"
    secrets-path: ""
```

### 6. Simulate the workflow

From your project root:

```bash
cre workflow simulate my-workflow
```

You’ll see trigger options, e.g.:

```
🚀 Workflow simulation ready. Please select a trigger:
1. cron-trigger@1.0.0 Trigger
2. evm:ChainSelector:16015286601757825753@1.0.0 LogTrigger
```

* **Cron Trigger**: choose `1` → the workflow executes immediately on the schedule.
* **Log Trigger**: choose `2` → provide the example inputs:

```
Transaction Hash: 0x420721d7d00130a03c5b525b2dbfd42550906ddb3075e8377f9bb5d1a5992f8e
Log Event Index: 0
```

Example output:

```
🔗 EVM Trigger Configuration:
Please provide the transaction hash and event index for the EVM log event.
Enter transaction hash (0x...): 0x420721...
Enter event index (0-based): 0
Fetching transaction receipt for transaction 0x420721...
Found log event at index 0: contract=0x1d5986..., topics=3
Created EVM trigger log for transaction 0x420721..., event 0
```

---

## Security Considerations

**⚠️ Important Notes**

1. **Demo project** – Not production-ready.
2. **Demo contracts** – Not audited; do not use as-is in production.
3. **Use your own RPCs** – For stability and performance, prefer private RPCs for deployment and chain writes.
4. **Secrets hygiene** – Keep real secrets out of version control; use secure secret managers for `.env` values.

---

## License

This project is licensed under the **MIT License** – see the [LICENSE](./LICENSE) file for details.