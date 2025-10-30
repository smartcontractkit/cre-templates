<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE&message=Home&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE&message=Docs&color=blue)](https://docs.chain.link/cre)

</div>

## Trying out the Bring Your Own Data Templates

These templates provide end-to-end bring your own data examples. They serve as a
starting point for bringing your own data on-chain with the Chainlink Runtime Environment (CRE).

Follow the steps below to run the examples:

### 1. CRE CLI

Install the [CRE CLI](https://docs.chain.link/cre).

### 2. Update .env file

You need to add a private key to the .env file. This is specifically required if you want to simulate chain writes. For that to work the key should be valid and funded.
If your workflow does not do any chain write then you can just put any dummy key as a private key. e.g.
```
CRE_ETH_PRIVATE_KEY=0000000000000000000000000000000000000000000000000000000000000001
```

### 3. Configure RPC endpoints

For local simulation to interact with a chain, you must specify RPC endpoints for the chains you interact with in the `project.yaml` file. This is required for submitting transactions and reading blockchain state.

Note: The following 7 chains are supported in local simulation (both testnet and mainnet variants):
- Ethereum (`ethereum-testnet-sepolia`, `ethereum-mainnet`)
- Base (`ethereum-testnet-sepolia-base-1`, `ethereum-mainnet-base-1`)
- Avalanche (`avalanche-testnet-fuji`, `avalanche-mainnet`)
- Polygon (`polygon-testnet-amoy`, `polygon-mainnet`)
- BNB Chain (`binance-smart-chain-testnet`, `binance-smart-chain-mainnet`)
- Arbitrum (`ethereum-testnet-sepolia-arbitrum-1`, `ethereum-mainnet-arbitrum-1`)
- Optimism (`ethereum-testnet-sepolia-optimism-1`, `ethereum-mainnet-optimism-1`)

Add your preferred RPCs under the `rpcs` section. For chain names, refer to https://github.com/smartcontractkit/chain-selectors/blob/main/selectors.yml

```yaml
rpcs:
  - chain-name: ethereum-testnet-sepolia
    url: <Your RPC endpoint to ETH Sepolia>
```
Ensure the provided URLs point to valid RPC endpoints for the specified chains. You may use public RPC providers or set up your own node.

## 4. [Optional] Deploy contracts

This step can be skipped if you are only going to test against local simulation.

Follow instructions in [../contracts/README.md](../contracts/README.md) to deploy the PoR or NAV contracts.

# 5. [Optional] generate contract bindings

To enable seamless interaction between the workflow and the contracts, Go bindings need to be generated from the contract ABIs. These ABIs are located in projectRoot/contracts/src/abi. Use the cre generate-bindings command to generate the bindings.

Note: Bindings for the template is pre-generated, so you can skip this step if there is no abi/contract changes. This command must be run from the <b>project root directory</b> where project.yaml is located. The CLI looks for a contracts folder and a go.mod file in this directory.

```bash
# Navigate to your project root (where project.yaml is located)
# Generate bindings for all contracts
cre generate-bindings evm

# The bindings will be generated in contracts/evm/src/generated/
```

This will create Go binding files for all the contracts (DataFeedsCache, etc.) that can be imported and used in your workflow.

## 6. [Optional] Configure workflow

Only required if you would like to test different configurations or if you deployed
your own contracts in step 4.

Configure `config.json` for the workflow you would like to run.

[./por/README.md](./por/README.md)

[./nav/README.md](./nav/README.md)

## 7. Simulate the workflow

> **Note:** Run `go mod tidy` to update dependencies after generating bindings.
```bash
go mod tidy

cre workflow simulate <nav|por>
```
