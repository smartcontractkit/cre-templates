# cre/ — Report Producer Workflow

CRE project containing the `offchain-report-write` workflow and the `Storage` consumer
contract (+ generated bindings). It writes a `uint256` timestamp on-chain via the chain
writer and logs the DON-signed report between `REPORT_JSON_BEGIN` / `REPORT_JSON_END`
markers, in the JSON format the validator app accepts.

**See [../README.md](../README.md) for the full template guide.**

## Quick reference

```bash
# install (from this template's root)
cd cre/contracts && bun install && cd ../..
cd cre/offchain-report-write && bun install && cd ../..

# simulate (from this directory, the CRE project root)
cre workflow simulate offchain-report-write --target staging-settings

# unit tests
cd offchain-report-write && bun test
```

## Regenerating bindings

`contracts/evm/src/abi/Storage.abi` is the source the bindings are generated from. If
you change `Storage.sol`, regenerate the ABI with your toolchain of choice (e.g.
`forge inspect Storage abi --json`), overwrite `src/abi/Storage.abi`, then re-run from
this directory:

```bash
cre generate-bindings evm --language typescript --non-interactive
```
