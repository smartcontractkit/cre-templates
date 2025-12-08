<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE\&message=Home\&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE\&message=Docs\&color=blue)](https://docs.chain.link/cre)

</div>

# Read Data Feeds - CRE Building Blocks

A set of minimal examples that, on a cron schedule (every 10 minutes), read from Chainlink Data Feeds using the CRE chain reader.

There are **two types of examples**:

1. **Regular (single-value) price feeds**

    * Reads `decimals()` and `latestAnswer()` from Chainlink Data Feeds (e.g. BTC/USD, ETH/USD), logs the scaled values, and returns a JSON array of results.
    * Example network: **Arbitrum One (mainnet)**.

   Production contracts:

    * BTC/USD on Arbitrum One: [0x6ce185860a4963106506C203335A2910413708e9](https://arbiscan.io/address/0x6ce185860a4963106506C203335A2910413708e9#code)
    * ETH/USD on Arbitrum One: [0x639Fe6ab55C921f74e7fac1ee960C0B6293ba612](https://arbiscan.io/address/0x639Fe6ab55C921f74e7fac1ee960C0B6293ba612#code)

2. **MVR (Multiple Variable Response) bundle feeds**

    * Reads `latestBundle()` and `bundleDecimals()` from a Chainlink **BundleAggregatorProxy**, decodes the returned bytes into multiple typed fields, and returns a structured JSON object (including scaled numeric values).
    * Example network: **Base mainnet**.
    * Example MVR feed: **S&P Global SSA EURC**.

   References:

    * Production MVR feed overview (S&P Global SSA EURC on Base):
      [https://data.chain.link/feeds/base/base/sp-global-ssa-eurc](https://data.chain.link/feeds/base/base/sp-global-ssa-eurc)
    * BundleAggregatorProxy contract on Base:
      [https://basescan.org/address/0xcF99622B5440a338f45daEE134d531A4BE64251F#code](https://basescan.org/address/0xcF99622B5440a338f45daEE134d531A4BE64251F#code)
    * MVR decoding guide (Solidity / EVM):
      [https://docs.chain.link/data-feeds/mvr-feeds/guides/evm-solidity](https://docs.chain.link/data-feeds/mvr-feeds/guides/evm-solidity)

---

## Get Started

Language-specific examples (both regular and MVR) live under the `building-blocks/read-data-feeds` directory.

* To get started with **Go (regular data feeds)**, see the [Go README](https://github.com/smartcontractkit/cre-templates/blob/main/building-blocks/read-data-feeds/read-data-feeds-go/README.md).

* To get started with **TypeScript (regular data feeds)**, see the [TypeScript README](https://github.com/smartcontractkit/cre-templates/blob/main/building-blocks/read-data-feeds/read-data-feeds-ts/README.md).

* To get started with **Go MVR bundle feeds**, use the Go example that reads `latestBundle()` / `bundleDecimals()` from the S&P Global SSA EURC feed on Base (see the corresponding Go README in the MVR example directory).

* To get started with **TypeScript MVR bundle feeds**, use the TypeScript example that reads and decodes the same MVR feed on Base (see the corresponding TypeScript README in the MVR example directory).
