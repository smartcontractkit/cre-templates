# Trying out the Developer PoR example

This template provides an end-to-end Proof-of-Reserve (PoR) example (including precompiled smart contracts). It's designed to showcase key CRE capabilities and help you get started with local simulation quickly.

Follow the steps below to run the example:

## 1. Initialize CRE project

Start by initializing a new CRE project. This will scaffold the necessary project structure and a template workflow. Run cre init in the directory where you'd like your CRE project to live. Note that workflow names must be exactly 10 characters long (we will relax this requirement in the future).

Example output:
```
Project name?: my_cre_project
✔ Development PoR Example to understand capabilities and simulate workflows
✔ Workflow name?: workflow01
```

## 2. Update .env file

You need to add a private key to the .env file. This is specifically required if you want to simulate chain writes. For that to work the key should be valid and funded.
If your workflow does not do any chain write then you can just put any dummy key as a private key. e.g.
```
CRE_ETH_PRIVATE_KEY=0000000000000000000000000000000000000000000000000000000000000001
```

## 3. Configure RPC endpoints

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

## 4. Deploy contracts

Deploy the BalanceReader, MessageEmitter, ReserveManager and SimpleERC20 contracts. You can either do this on a local chain or on a testnet using tools like cast/foundry.

For a quick start, you can also use the pre-deployed contract addresses on Ethereum Sepolia—no action required on your part if you're just trying things out.

For completeness, the Solidity source code for these contracts is located under projectRoot/contracts/evm/src.
- chain: `ethereum-testnet-sepolia`
- ReserveManager contract address: `0x073671aE6EAa2468c203fDE3a79dEe0836adF032`
- SimpleERC20 contract address: `0x4700A50d858Cb281847ca4Ee0938F80DEfB3F1dd`
- BalanceReader contract address: `0x4b0739c94C1389B55481cb7506c62430cA7211Cf`
- MessageEmitter contract address: `0x1d598672486ecB50685Da5497390571Ac4E93FDc`

## 5. [Optional] Generate contract bindings

To enable seamless interaction between the workflow and the contracts, Go bindings need to be generated from the contract ABIs. These ABIs are located in projectRoot/contracts/src/abi. Use the cre generate-bindings command to generate the bindings.

Note: Bindings for the template is pre-generated, so you can skip this step if there is no abi/contract changes. This command must be run from the <b>project root directory</b> where project.yaml is located. The CLI looks for a contracts folder and a go.mod file in this directory.

```bash
# Navigate to your project root (where project.yaml is located)
# Generate bindings for all contracts
cre generate-bindings evm

# The bindings will be generated in contracts/evm/src/generated/
# Each contract gets its own package subdirectory:
# - contracts/evm/src/generated/ierc20/IERC20.go
# - contracts/evm/src/generated/reserve_manager/ReserveManager.go
# - contracts/evm/src/generated/balance_reader/BalanceReader.go
# - etc.
```

This will create Go binding files for all the contracts (ReserveManager, SimpleERC20, BalanceReader, MessageEmitter, etc.) that can be imported and used in your workflow.

## 6. Configure workflow

Configure `config.json` for the workflow
- `schedule` should be set to `"0 */1 * * * *"` for every 1 minute(s) or any other cron expression you prefer, note [CRON service quotas](https://docs.chain.link/cre/service-quotas)
- `url` should be set to existing reserves HTTP endpoint API
- `tokenAddress` should be the SimpleERC20 contract address
- `reserveManagerAddress` should be the ReserveManager contract address
- `balanceReaderAddress` should be the BalanceReader contract address
- `messageEmitterAddress` should be the MessageEmitter contract address
- `chainName` should be name of selected chain (refer to https://github.com/smartcontractkit/chain-selectors/blob/main/selectors.yml)
- `gasLimit` should be the gas limit of chain write

The config is already populated with deployed contracts in template.

Note: Make sure your `workflow.yaml` file is pointing to the config.json, example:

```yaml
staging-settings:
  user-workflow:
    workflow-name: "workflow01"
  workflow-artifacts:
    workflow-path: "."
    config-path: "./config.json"
    secrets-path: ""
```


## 7. Simulate the workflow

> **Note:** Run `go mod tidy` to update dependencies after generating bindings.
```bash
go mod tidy

cre workflow simulate <path-to-workflow>
```

After this you will get a set of options similar to:

```
🚀 Workflow simulation ready. Please select a trigger:
1. cron-trigger@1.0.0 Trigger
2. evm:ChainSelector:16015286601757825753@1.0.0 LogTrigger

Enter your choice (1-2):
```

You can simulate each of the following triggers types as follows

### 7a. Simulating Cron Trigger Workflows

Select option 1, and the workflow should immediately execute.

### 7b. Simulating Log Trigger Workflows

Select option 2, and then two additional prompts will come up and you can pass in the example inputs:

Transaction Hash: 0x9394cc015736e536da215c31e4f59486a8d85f4cfc3641e309bf00c34b2bf410
Log Event Index: 0

The output will look like:
```
🔗 EVM Trigger Configuration:
Please provide the transaction hash and event index for the EVM log event.
Enter transaction hash (0x...): 0x9394cc015736e536da215c31e4f59486a8d85f4cfc3641e309bf00c34b2bf410
Enter event index (0-based): 0
Fetching transaction receipt for transaction 0x9394cc015736e536da215c31e4f59486a8d85f4cfc3641e309bf00c34b2bf410...
Found log event at index 0: contract=0x1d598672486ecB50685Da5497390571Ac4E93FDc, topics=3
Created EVM trigger log for transaction 0x9394cc015736e536da215c31e4f59486a8d85f4cfc3641e309bf00c34b2bf410, event 0
```
