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

Deploy MockPool and ProtocolSmartWallet contracts. Needs to be run once per chain you are targeting.
```
ENABLE_WORKFLOW_SIMULATION=true \
forge script ./scripts/multi-chain-token-manager/DeployTokenManagerContracts.s.sol \
--rpc-url "<RPC URL for target chain>" \
--private-key <EOA funded on target chain> \
--broadcast
```

Configure ProtocolSmartWallet contracts for CCIP.

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
to test the workflow rebalancing to the chain with the highest APY.

First, update the hardcoded APY in [SetPoolAPY.s.sol](./scripts/multi-chain-token-manager/SetPoolAPY.s.sol).

Finally, configure the APY on-chain.
```
ENABLE_WORKFLOW_SIMULATION=true \
forge script ./scripts/multi-chain-token-manager/SetPoolAPY.s.sol \
--rpc-url "<RPC URL for target chain>" \
--private-key <EOA funded on target chain> \
--broadcast
```
