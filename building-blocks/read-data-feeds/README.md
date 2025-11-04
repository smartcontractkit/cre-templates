<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE&message=Home&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE&message=Docs&color=blue)](https://docs.chain.link/cre)

</div>

# Read Data Feeds - CRE Building Blocks

A minimal TypeScript example that, on a cron schedule (every 10 minutes), reads decimals() and latestAnswer() from Chainlink Data Feeds on Arbitrum One (mainnet) using the CRE chain reader. It logs the scaled values and returns a JSON array of results.

Production contract (BTC/USD on Arbitrum One): [0x6ce185860a4963106506C203335A2910413708e9](https://arbiscan.io/address/0x6ce185860a4963106506C203335A2910413708e9#code)
Production contract (ETH/USD on Arbitrum One): [0x639Fe6ab55C921f74e7fac1ee960C0B6293ba612](https://arbiscan.io/address/0x639Fe6ab55C921f74e7fac1ee960C0B6293ba612#code)

---

### Get Started

- To get started with Go, see the [Go README](https://github.com/smartcontractkit/cre-templates/blob/main/building-blocks/read-data-feeds/read-data-feeds-go/README.md).

- To get started with TypeScript, see the [TypeScript README](https://github.com/smartcontractkit/cre-templates/blob/main/building-blocks/read-data-feeds/read-data-feeds-ts/README.md).