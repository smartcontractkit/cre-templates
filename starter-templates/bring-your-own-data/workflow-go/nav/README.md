<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE&message=Home&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE&message=Docs&color=blue)](https://docs.chain.link/cre)

</div>

# Bring Your Own Data - Net Asset Value (NAV) CRE Template

This template provides an end-to-end Net-Asset-Value (NAV) example.

You can either run against the predeployed contracts with local simulation
or deploy your own contracts to run the example end-to-end.

## Trying it out

See instructions in [../README.md](../README.md).

## [Optional] Targeting your own contracts

1. Deploy your DataFeedsCache contracts following the instructions in [../../contracts/README.md](../../contracts/README.md)
2. Update the target chains/contracts in [./config.json](./config.json)

## [Optional] Configure workflow

Configure [config.json](./config.json) for the workflow
- `schedule` should be set to `"0 */5 * * * *"` for every 5 minutes or any other cron expression you prefer
- `url` API url that returns NAV data; if updated, then you need to modify the workflow code in [./workflow.go](./workflow.go) to handle your API schema
- `dataIdHex` Data feed ID; any 32 hex characters
- `dataFeedsCacheAddress` should be the DataFeedsCache contract address; either the default or the one deployed from [../../contracts/](../../contracts/)
- `chainName` should be name of selected chain (refer to https://github.com/smartcontractkit/chain-selectors/blob/main/selectors.yml)
- `gasLimit` should be the gas limit of chain write
