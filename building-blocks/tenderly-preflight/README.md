<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE&message=Home&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE&message=Docs&color=blue)](https://docs.chain.link/cre)

</div>

# Tenderly Preflight - CRE Building Block

This building block demonstrates how to simulate an EVM transaction on a [Tenderly Virtual Environment](https://docs.tenderly.co/virtual-environments) from inside a CRE workflow using the [`@tenderly/cre-preflight`](https://www.npmjs.com/package/@tenderly/cre-preflight) library, and gate an on-chain write on the **consensus-verified verdict**.

## The Problem

A workflow that writes on-chain commits real gas to a transaction that might revert. Simulating the transaction off-chain first is not enough on its own: each node would need to run its own simulation, and the results would have to agree before the DON can act on them.

## The Solution

`@tenderly/cre-preflight` forks the chain on a Tenderly Virtual Environment **per node, at one agreed block**, mines the transaction on each fork, and runs the per-node receipts through CRE consensus. The workflow gets back a single verdict — `success`, `reverted`, `rejected`, `misconfigured`, `unavailable`, `oversized`, or `indeterminate` — and can skip the on-chain write when the transaction would fail.

- **Consensus-verified** — the verdict is aggregated across the DON, not one node's opinion
- **Nothing reaches the real chain** — the transaction is mined on a throwaway fork, deleted afterwards by default
- **Credentials stay in the Vault** — the Tenderly access key is a CRE secret name, never a config value

## Get Started

- **TypeScript**: See the [TypeScript README](./tenderly-preflight-ts/README.md) for setup, configuration, and expected simulation output.

## Quick Example

```typescript
import { TenderlyPreflight } from '@tenderly/cre-preflight';

const tenderly = new TenderlyPreflight(runtime.config.tenderly);
const verdict = tenderly.sendTransaction(runtime, { from, to, data });

if (verdict.outcome !== 'success') {
  runtime.log(`skipping write: ${verdict.outcome} ${verdict.reason}`);
  return 'skipped';
}
// ... safe to write on-chain
```

## Reference Documentation

- [@tenderly/cre-preflight on npm](https://www.npmjs.com/package/@tenderly/cre-preflight)
- [tenderly-cre-preflight GitHub](https://github.com/Tenderly/tenderly-cre-preflight)
- [Tenderly Virtual Environments](https://docs.tenderly.co/virtual-environments)
- [CRE Documentation](https://docs.chain.link/cre)
