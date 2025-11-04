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

### 2. Install dependencies

If **Bun** is not already installed, follow the instructions at: [https://bun.com/docs/installation](https://bun.com/docs/installation)

From your project root, run:

```bash
bun install --cwd ./my-workflow
```

### 3. Update .env file

You need to add a private key to the .env file. This is specifically required if you want to simulate chain writes. For that to work the key should be valid and funded.
If your workflow does not do any chain write then you can just put any dummy key as a private key. e.g.
```
CRE_ETH_PRIVATE_KEY=0000000000000000000000000000000000000000000000000000000000000001
```

### 4. Configure RPC endpoints

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

### 5. [Optional] Deploy contracts

This step can be skipped if you are only going to test against local simulation.

Follow instructions in [../contracts/README.md](../contracts/README.md) to deploy the PoR or NAV contracts.

### 6. [Optional] Configure workflow

Only required if you would like to test different configurations or if you deployed
your own contracts in step 4.

Configure `config.json` for the workflow you would like to run. Refer to instructions for [PoR](./por/README.md) and [NAV](./nav/README.md).

### 7. Simulate the workflow

Run the command from the **project root** and pass the **path to the workflow directory**:

```bash
cre workflow simulate <path-to-workflow-directory>
```

For Proof-of-Reserve:

```bash
cre workflow simulate por
```

For Net-Asset-Value:

```bash
cre workflow simulate nav
```
