<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE\&message=Home\&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE\&message=Docs\&color=blue)](https://docs.chain.link/cre)

</div>

# CRE Starter Templates

End-to-end **CRE (Chainlink Runtime Environment)** templates that demonstrate richer workflows than the simple building blocks. Each template combines multiple capabilities (on-chain + off-chain), production-like configuration, and a runnable simulation flow.

---

## Table of Contents

- [What are Starter Templates?](#what-are-starter-templates)
- [Available Templates](#available-templates)
- [License](#license)

---

## What are Starter Templates?

Starter templates help you bootstrap real-world patterns:
- Multiple capabilities in one workflow (scheduler, HTTP, EVM calls/writes, secrets).
- Clear configuration via `config.json`.
- Optional precompiled smart contracts and generated bindings.
- Ready to **simulate locally** with `cre workflow simulate`.

They are more comprehensive than **building-blocks**, and can be adapted into your own projects.

---

## Available Templates

1. **Custom Data Feed** — [`./custom-data-feed`](./custom-data-feed)  
   Periodically fetch off-chain data (HTTP) and **push updates on-chain**, demonstrating scheduling, secrets, bindings, and chain writes.

2. **Bring Your Own Data (BYOD)** — NAV & PoR — [`./bring-your-own-data`](./bring-your-own-data)  
   End-to-end examples for publishing **Net Asset Value** and **Proof of Reserve** data on-chain using CRE workflows and demo contracts.

3. **Multi-Chain Token Manager** — [`./multi-chain-token-manager`](./multi-chain-token-manager)  
   Orchestrate token operations and state across multiple chains, showing RPC configuration, bindings, and cross-chain patterns.

> Each subdirectory includes its own README with template-specific steps and example logs.

## License

MIT — see the repository’s [LICENSE](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE).