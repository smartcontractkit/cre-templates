<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE&message=Home&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE&message=Docs&color=blue)](https://docs.chain.link/cre)

</div>

## Trying out the Multi-Chain Token Manager Template

This template provides end-to-end example of how to maximize supply APY by rebalancing tokens cross-chain with CRE via CCIP. It serves as a starting point for writing your own multi-chain token manager with the Chainlink Runtime Environment (CRE).

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

Follow instructions in [../contracts/README.md](../contracts/README.md) to deploy your own versions of the contracts.

### 6. [Optional] Configure workflow

Only required if you would like to test different configurations or if you deployed
your own contracts in step 4.

Configure [config.json](./workflow/config.json) for the workflow
- `schedule` should be set to `"0 */5 * * * *"` for every 5 minutes or any other cron expression you prefer
- `minBPSDeltaForRebalance` minimum basis points difference in APR for tokens to be rebalanced cross-chain
- `assetAddress` CCIP CCT address; currently set to CCIP BnM token
- `poolAddress` should be the MockPool contract address
- `protocolSmartWalletAddress` should be the ProtocolSmartWallet contract address
- `chainName` should be name of selected chain (refer to https://github.com/smartcontractkit/chain-selectors/blob/main/selectors.yml)
- `gasLimit` should be the gas limit of chain write

### 7. Simulate the workflow

Run the command from the **project root** and pass the **path to the workflow directory**:

```bash
cre workflow simulate workflow
```
