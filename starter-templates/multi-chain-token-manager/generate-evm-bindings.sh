#!/bin/bash

set -e -u -x

cp ../contracts/src/multi-chain-token-manager/MockPool.sol ./contracts/evm/src/
cat ../contracts/foundry-artifacts/MockPool.sol/MockPool.json | jq .abi > ./contracts/evm/src/abi/MockPool.abi
cp ../contracts/src/multi-chain-token-manager/ProtocolSmartWallet.sol ./contracts/evm/src/
cat ../contracts/foundry-artifacts/ProtocolSmartWallet.sol/ProtocolSmartWallet.json | jq .abi > ./contracts/evm/src/abi/ProtocolSmartWallet.abi

cre generate-bindings evm
