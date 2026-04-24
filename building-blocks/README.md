<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE&message=Home&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE&message=Docs&color=blue)](https://docs.chain.link/cre)

</div>

# CRE Templates – Building Blocks

Small, focused examples showing how to use the **Chainlink Runtime Environment (CRE)** to interact with **on-chain** and/or **off-chain** components. Each block is a self-contained workflow with a short README and runnable config so you can learn by doing.

---

## Table of Contents

- [What are Building Blocks?](#what-are-building-blocks)
- [Available Blocks](#available-blocks)
- [When to Use Which Block](#when-to-use-which-block)
- [License](#license)

---

## What are Building Blocks?

Each block demonstrates a single CRE capability or pattern (e.g., reading on-chain data, writing to an off-chain service, using secrets, scheduling with cron). Use them as learning references or as starting points to compose your own workflows.

**Common traits:**
- Minimal code focused on one task
- Clear configuration via `config.json`
- Runnable locally with `cre workflow simulate`
- Documented prerequisites and expected output

---

## Available Blocks

### 1) **Key-Value Store (AWS S3)**
Path: [`./kv-store`](./kv-store)

- Reads a value from an **AWS S3 object**, increments it, and writes it back.
- Uses CRE secrets (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`) and **SigV4** signing.
- Performs a **consensus read** of the current value across nodes, then writes the agreed next value.

👉 See the block’s README for setup, config, and sample logs.

---

### 2) **Read Data Feeds**
Path: [`./read-data-feeds`](./read-data-feeds)

- On a cron schedule, reads `decimals()` and `latestAnswer()` from **Chainlink Data Feeds** (example targets BTC/USD and ETH/USD on Arbitrum One).
- Shows how to add an ABI, generate Go bindings, configure chain RPC, and log scaled answers.

👉 See the block’s README for setup, config, and sample logs.

---

### 3) **Compression Utils**
Path: [`./compression-utils`](./compression-utils)

- Demonstrates how to perform compression and decompression operations in CRE workflows using **fflate** as a pure JavaScript alternative to the Node.js `zlib` module.
- Covers gzip, deflate, zlib formats, compression levels, ZIP archives, and string utilities.
- Shows how to fetch a large JSON payload via CRE HTTP and compress it to stay within the 100KB response body limit.

👉 See the block’s README for setup, config, and sample logs.

---

### 4) **XML Utils**
Path: [`./xml-utils`](./xml-utils)

- Demonstrates how to use XML in CRE workflows using **fast-xml-parser** as a pure JavaScript alternative to the browser's `DOMParser` and the Node.js `xml2js` module.
- Covers parsing, validating, and building XML.
- Shows how to fetch a XML payload via CRE HTTP, validate it, parse it with attributes and namespaces, and rebuild a filtered XML document.

👉 See the block’s README for setup, config, and sample logs.

---

## When to Use Which Block

* **kv-store**: You want to see an **off-chain write** pattern (AWS S3), secrets usage, SigV4 signing, and a **consensus read → single write** flow.
* **read-data-feeds**: You want to **read on-chain data** via contract calls, manage ABIs/bindings, and configure **RPC** access.
* **compression-utils**: You want to handle **data compression** in CRE Typescript workflows, especially when working with large datasets and the CRE HTTP response limit.
* **xml-utils**: You want to work with **XML data** in CRE workflows, including parsing, validating, and building XML documents.
---

## License

MIT — see the repository’s [LICENSE](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE).
