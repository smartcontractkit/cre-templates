#!/bin/bash

set -e -u -x

cp ../contracts/lib/chainlink-evm/contracts/src/v0.8/data-feeds/DataFeedsCache.sol ./contracts/evm/src/
cat ../contracts/foundry-artifacts/DataFeedsCache.sol/DataFeedsCache.json | jq .abi > ./contracts/evm/src/abi/DataFeedsCache.abi

cre generate-bindings evm

