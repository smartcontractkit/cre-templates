<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE&message=Home&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE&message=Docs&color=blue)](https://docs.chain.link/cre)

</div>

# Deploying your own Bring Your Own Data contracts

Instructions on how to deploy your own contracts for use with the Bring Your Own Data templates.

## Prerequisites

Install [Foundry](https://getfoundry.sh/introduction/installation/).

### Configure RPC URLs

Populate the RPC URLs for sepolia and base-sepolia in [foundry.toml](./foundry.toml).

## Proof-of-Reserve (PoR)

```
ENABLE_WORKFLOW_SIMULATION=true \
forge script ./scripts/por/DeployPoRContracts.s.sol \
--rpc-url "<RPC URL for target chain>" \
--private-key <EOA funded on target chain> \
--broadcast
```

Update the PoR workflow config with the newly deployed contract addresses as per [../workflow/por/README.md](../workflow/por/README.md).

## Net Asset Value (NAV)

```
ENABLE_WORKFLOW_SIMULATION=true \
forge script ./scripts/por/DeployPoRContracts.s.sol \
--private-key <EOA funded on target chain> \
--broadcast
```

Update the NAV workflow config with the newly deployed contract addresses as per [../workflow/por/README.md](../workflow/nav/README.md).

## Verifying Contracts

Add `--verify` and populate the Etherscan API keys in [foundry.toml](./foundry.toml) to automatically verify the contracts on Etherscan.
