<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE&message=Home&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE&message=Docs&color=blue)](https://docs.chain.link/cre)

</div>

# CRE Templates

A curated set of ready-to-run examples for the **Chainlink Runtime Environment (CRE)**:

- **Building Blocks** – tiny, focused workflows that teach one concept at a time (on-chain reads, off-chain calls, secrets, scheduling, etc.).
- **Starter Templates** – opinionated, end-to-end workflows that combine multiple capabilities and look closer to real-world use cases.

Use these as references or starting points to compose your own production workflows.

---

## Table of Contents

- [Repository Structure](#repository-structure)
- [When to Use Which](#when-to-use-which)
- [License](#license)

---

## Repository Structure

### Building Blocks
Small, focused examples. Each directory includes its own README.

- **`building-blocks/kv-store`** – Read/modify/write a value in **AWS S3** using SigV4-signed HTTP requests, CRE secrets, and a **consensus read → single write** flow.
- **`building-blocks/read-data-feeds`** – Read `decimals()` and `latestAnswer()` from **Chainlink Data Feeds** on a schedule; includes ABI/bindings and RPC config examples.

### Starter Templates
More complex, end-to-end workflows. Each directory includes its own README (some marked **WIP**).

- **`starter-template/custom-data-feed`** – Fetch off-chain data (HTTP) and **push updates on-chain**; shows cron scheduling, secrets, bindings, and chain writes.
- **`starter-template/bring-your-own-data`** – NAV (Net Asset Value) & PoR (Proof of Reserve) templates for publishing your own data on-chain.
- **`starter-template/multi-chain-token-manager`** – Orchestrate token operations and state across **multiple chains**.

---

## When to Use Which

* **Building Blocks**
  Use these when you want to learn a **single concept quickly** (e.g., secrets + HTTP signing, reading a data feed, cron triggers). Great for copy/paste into your project.

* **Starter Templates**
  Choose these when you want a **runnable reference architecture** that strings together multiple steps (off-chain + on-chain), includes contracts/bindings, and mirrors real workflows.

---

## License

MIT — see the repository’s [LICENSE](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE).