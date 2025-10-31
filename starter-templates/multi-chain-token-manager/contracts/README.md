<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE&message=Home&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE&message=Docs&color=blue)](https://docs.chain.link/cre)

</div>

# Multi-Chain Token Manager - Smart Contracts - CRE Template

Instructions on how to deploy your own contracts for use with the Multi-Chain Token Manager template.

## Prerequisites

Install [Foundry](https://getfoundry.sh/introduction/installation/).

## Deploy and configure contracts

### Fund Deployer Wallet

First your deployer wallet needs:
- Native gas tokens to deploy the contracts (e.g., SepoliaETH)
- LINK tokens that will be transferred to the ProtocolSmartWallet to pay CCIP fees for cross-chain token transfers
- CCIP BnM tokens (testnet only)
  - Another CCIP Cross-Chain Token (CCT) can be used. Simply update the BNM token addresses in the [deploy script](./scripts/multi-chain-token-manager/DeployTokenManagerContracts.s.sol) and [workflow config](../workflow/workflow/config.json) to the target CCT.

Faucets:
- Native gas tokens and LINK tokens https://faucets.chain.link
- CCIP BnM tokens https://docs.chain.link/ccip/test-tokens

### Deploy

Deploy multi-chain token manager contracts. Needs to be run once per chain you are targeting.

This script will:
- Deploy the ProtocolSmartWallet (PSW) and MockPool contracts
- Transfer 1 LINK to the PSW contract for CCIP fees
- Transfer 1 CCIP BnM to the PSW contract that is then deposited into to the MockPool contract

```
ENABLE_WORKFLOW_SIMULATION=true \
forge script ./scripts/multi-chain-token-manager/DeployTokenManagerContracts.s.sol \
--rpc-url "<RPC URL for target chain>" \
--private-key <EOA funded on target chain> \
--broadcast
```

## Configure

Configure ProtocolSmartWallet contracts for CCIP.

This step configures each ProtocolSmartWallet contract with the addresses of the ProtocolSmartWallet contracts deployed on each other chain, which allows it to verify that the CCIP sender is a ProtocolSmartWallet on another chain when it receives a cross-chain token transfer with data

First, update the `getProtocolSmartWallet` function in [ConfigureTokenManagerContracts.s.sol](./scripts/multi-chain-token-manager/ConfigureTokenManagerContracts.s.sol) to return the addresses of the contracts you deployed in the previous step.

Now configure your deployed ProtocolSmartWallet contracts. Needs to be run once per chain you deployed on in the previous step.

```
ENABLE_WORKFLOW_SIMULATION=true \
forge script ./scripts/multi-chain-token-manager/ConfigureTokenManagerContracts.s.sol \
--rpc-url "<RPC URL for target chain>" \
--private-key <EOA funded on target chain> \
--broadcast
```

## Live APY Rebalancing

Each pool is deployed with a preconfigured default APY, but you can update these
to test the workflow rebalancing tokens to the chain with the highest APY.

First, update the hardcoded APY in [SetPoolAPY.s.sol](./scripts/multi-chain-token-manager/SetPoolAPY.s.sol).

Finally, configure the APY on-chain.
```
ENABLE_WORKFLOW_SIMULATION=true \
forge script ./scripts/multi-chain-token-manager/SetPoolAPY.s.sol \
--rpc-url "<RPC URL for target chain>" \
--private-key <EOA funded on target chain> \
--broadcast
```
