# Trying out the Developer PoR example

This template provides an end-to-end Proof-of-Reserve (PoR) example (including precompiled smart contracts). It's designed to showcase key CRE capabilities and help you get started with local simulation quickly.

Follow the steps below to run the example:

## 1. Initialize CRE project

Start by initializing a new CRE project. This will scaffold the necessary project structure and a template workflow. Run cre init in the directory where you'd like your CRE project to live.

Example output:

```
Project name?: my_cre_project
✔ Custom data feed: Typescript updating on-chain data periodically using offchain API data
✔ Workflow name?: workflow01
```

## 2. Update .env file

You need to add a private key to the .env file. This is specifically required if you want to simulate chain writes. For that to work the key should be valid and funded.
If your workflow does not do any chain write then you can keep a dummy key as a private key. e.g.

```
CRE_ETH_PRIVATE_KEY=0000000000000000000000000000000000000000000000000000000000000001
```

## 3. Install dependencies
```
cd workflowName && bun install
```

## 4. Configure RPC endpoints

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

## 5. Set up workflow secrets

This template workflow demonstrates the use of secrets. To configure them, a separate file containing the secret names is required. For this template, a secrets.yaml file is already provided in the workflow root directory with all necessary entries.

To simulate how secrets are accessed, you need to define them in your environment. Below is an example based on the generated secrets.yaml.

```bash
export SECRET_ADDRESS_ALL="0xMySecretAddress"
```

Or you may append the secret value to your existing `.env` file and then it will be automatically loaded by the CLI:

```bash
SECRET_ADDRESS_ALL=0xMySecretAddress
```

Note in our example the secret address is used to read a balance from. You can use the address of our deployed contract 0x4700A50d858Cb281847ca4Ee0938F80DEfB3F1dd.

```bash
export SECRET_ADDRESS_ALL="0x4700A50d858Cb281847ca4Ee0938F80DEfB3F1dd"
```

## 6. Deploy contracts

Deploy the BalanceReader, MessageEmitter, ReserveManager and SimpleERC20 contracts. You can either do this on a local chain or on a testnet using tools like cast/foundry.

For a quick start, you can also use the pre-deployed contract addresses on Ethereum Sepolia—no action required on your part if you're just trying things out.

## 7. Configure workflow

Configure `config.json` for the workflow

- `schedule` should be set to `"*/30 * * * * *"` for every 30 seconds or any other cron expression you prefer
- `url` should be set to existing reserves HTTP endpoint API
- `tokenAddress` should be the SimpleERC20 contract address
- `porAddress` should be the ReserveManager contract address
- `proxyAddress` should be the UpdateReservesProxySimplified contract address
- `balanceReaderAddress` should be the BalanceReader contract address
- `messageEmitterAddress` should be the MessageEmitter contract address
- `chainSelectorName` should be human-readable chain name of selected chain (refer to https://github.com/smartcontractkit/chain-selectors/blob/main/selectors.yml)
- `gasLimit` should be the gas limit of chain write

The config is already populated with deployed contracts in template.

Note: Make sure your `workflow.yaml` file is pointing to the config.json, example:

```yaml
local-simulation:
  user-workflow:
    workflow-name: "workflow01"
  workflow-artifacts:
    workflow-path: "./main.ts"
    config-path: "./config.json"
    secrets-path: "../secrets.yaml"
```

## 8. Simulate the workflow

Run the command from <b>project root directory</b> and pass in the path to the workflow directory.

```bash
cre workflow simulate <path-to-workflow-directory>
```

For a workflow directory named `workflow01` the exact command would be:

```bash
cre workflow simulate ./workflow01
```

After this you will get a set of options similar to:

```
🚀 Workflow simulation ready. Please select a trigger:
1. cron-trigger@1.0.0 Trigger
2. evm:ChainSelector:16015286601757825753@1.0.0 LogTrigger
3. http-trigger@1.0.0-alpha Trigger

Enter your choice (1-3):
```

You can simulate each of the following triggers types as follows

### 8a. Simulating Cron Trigger Workflows

Select option 1, and the workflow should immediately execute.

### 8b. Simulating Log Trigger Workflows

Select option 2, and then two additional prompts will come up and you can pass in the example inputs:

Transaction Hash: 0x420721d7d00130a03c5b525b2dbfd42550906ddb3075e8377f9bb5d1a5992f8e
Log Event Index: 0

The output will look like:

```
🔗 EVM Trigger Configuration:
Please provide the transaction hash and event index for the EVM log event.
Enter transaction hash (0x...): 0x420721d7d00130a03c5b525b2dbfd42550906ddb3075e8377f9bb5d1a5992f8e
Enter event index (0-based): 0
Fetching transaction receipt for transaction 0x420721d7d00130a03c5b525b2dbfd42550906ddb3075e8377f9bb5d1a5992f8e...
Found log event at index 0: contract=0x1d598672486ecB50685Da5497390571Ac4E93FDc, topics=3
Created EVM trigger log for transaction 0x420721d7d00130a03c5b525b2dbfd42550906ddb3075e8377f9bb5d1a5992f8e, event 0
```

### 8c. Simulating HTTP Trigger Workflows

Select option 3, and then an additional prompt will come up where you can pass in:

File Path: ./http_trigger_payload.json

The output will look like:

```
🔍 HTTP Trigger Configuration:
Please provide JSON input for the HTTP trigger.
You can enter a file path or JSON directly

Enter your input: ./http_trigger_payload.json
Loaded JSON from file: ./http_trigger_payload.json
Created HTTP trigger payload with 1 fields
```
