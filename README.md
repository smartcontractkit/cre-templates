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
- [E2E Tests](#e2e-tests)
- [Contributing](#contributing)
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

## E2E Tests

The `e2e/` directory contains end-to-end tests that scaffold every template using the real `cre init` CLI and validate the resulting project structure and build.

### Prerequisites

- [`cre` CLI](https://docs.chain.link/cre) installed (or a pre-built binary)
- Go toolchain
- `bun` or `npm` (for TypeScript templates)

### Running the tests

```bash
cd e2e

# Run all template tests:
go test -v -count=1 -timeout 300s

# Run a single template:
go test -run "TestAllTemplates/kv-store-go" -v -count=1
```

### Environment variables

| Variable | Description | Default |
|---|---|---|
| `CRE_CLI_PATH` | Path to the `cre` binary | `cre` (looked up on `$PATH`) |
| `CRE_TEMPLATE_REPO` | Local repo path used for template discovery | Parent of `e2e/` directory |
| `CRE_TEMPLATE_REPO_REF` | GitHub `owner/repo@ref` to add as a template source before running | Default remote repos |

### Testing a feature branch

By default the CLI fetches templates from the published remote repo. To test templates from an unpublished branch, push your branch and run:

```bash
CRE_TEMPLATE_REPO_REF="smartcontractkit/cre-templates@my-branch" go test -v -count=1 -timeout 300s
```

### What the tests verify

For each template discovered in the repo:

1. **Scaffold** — runs `cre init` with the template ID in an isolated temp directory
2. **Structure** — asserts `project.yaml`, workflow directories, and entrypoint files (`main.go`/`main.ts`) exist
3. **Build (Go)** — runs `go mod tidy` and `go build` (with WASM cross-compilation fallback for `wasip1`-tagged code)
4. **Build (TypeScript)** — finds all `package.json` files and runs `bun install` (or `npm install`)

---

## Contributing

We welcome contributions of new CRE workflow templates! If you have a useful workflow that others could benefit from, here's how to contribute:

1. **Fork the repository** and create a new branch for your template.
2. **Add your workflow template** in the appropriate directory:
   - `building-blocks/` – for small, focused examples that teach a single concept
   - `starter-template/` – for more complete, end-to-end workflows
3. **Include a README** in your template directory explaining:
   - What the workflow does
   - Prerequisites and setup instructions
   - How to run and test it
4. **Submit a Pull Request** with a clear description of your template and its use case.

Our team will review your PR and provide feedback. Once approved, we'll merge it into the main branch.

---

## License

MIT — see the repository's [LICENSE](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE).