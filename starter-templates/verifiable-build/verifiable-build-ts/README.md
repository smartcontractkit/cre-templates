# Verifiable Workflow Template (TypeScript)

A TypeScript workflow template with reproducible builds that work on both Windows and macOS. This template enables third-party verification of deployed workflows using the CRE CLI.

> **Note:** Go workflows are verifiable by default and do not require this template. This template is for TypeScript workflows only. For more details, see the [Verifying Workflows](https://docs.chain.link/cre/guides/operations/verifying-workflows) guide.

## Prerequisites

- [CRE CLI](https://github.com/smartcontractkit/cre-cli/releases) installed
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) running

## Project Structure

```
.
├── project.yaml
├── secrets.yaml
└── workflow/
    ├── Dockerfile
    ├── Makefile
    ├── bun.lock
    ├── config.production.json
    ├── config.staging.json
    ├── main.ts
    ├── main.test.ts
    ├── package.json
    ├── tsconfig.json
    └── workflow.yaml
```

## Getting Started

### Simulate the workflow

Run from the **project root directory**:

```bash
cre workflow simulate workflow --target=staging-settings
```

## Verifying a Workflow Build

Workflow verification lets anyone independently confirm that a deployed workflow matches its source code. The `cre workflow hash` command computes the workflow hash locally and compares it against the onchain workflow ID.

### Computing the workflow hash

From the **project root directory**, run:

```bash
cre workflow hash workflow --public_key <DEPLOYER_ADDRESS> --target production-settings
```

Replace `<DEPLOYER_ADDRESS>` with the deployer's public address (e.g. `0xb0f2D38245dD6d397ebBDB5A814b753D56c30715`).

Example output:

```
Compiling workflow...
✓ Workflow compiled
  Binary hash:   03c77e16354e5555f9a74e787f9a6aa0d939e9b8e4ddff06542b7867499c58ea
  Config hash:   3bdaebcc2f639d77cb248242c1d01c8651f540cdbf423d26fe3128516fd225b6
  Workflow hash: 001de36f9d689b57f2e4f1eaeda1db5e79f7991402e3611e13a5c930599c2297
```

The **Workflow hash** is the onchain workflow ID. If it matches the workflow ID observed onchain, the deployed workflow matches this source code.

### For verifiers (third-party auditors)

1. Install the [CRE CLI](https://github.com/smartcontractkit/cre-cli/releases). No login or deploy access is required.
2. Clone or unzip the shared workflow repository.
3. Run `cre workflow hash` as shown above, using the deployer's public address.
4. Compare the `Workflow hash` output with the onchain workflow ID. A match confirms the deployed workflow is built from this source.

### Generating the lockfile

If `bun.lock` is missing, generate it before building:

```bash
cd workflow
make lock
```

This runs the lockfile generation inside Docker to ensure consistency across platforms.

## How Reproducible Builds Work

The build process uses Docker to ensure identical output on any machine:

1. `make build` on the host starts a Docker build (`linux/amd64`)
2. Inside the container, `bun install --frozen-lockfile` installs exact dependencies from `bun.lock`
3. The workflow is compiled to a WASM binary (`workflow.wasm`)
4. The binary is exported back to the host

Because the build runs inside a pinned Docker image with a locked dependency tree, the same source always produces the same binary hash.

## Learn More

- [Verifying Workflows](https://docs.chain.link/cre/guides/operations/verifying-workflows) - Full verification guide
- [Deploying Workflows](https://docs.chain.link/cre/guides/operations/deploying-workflows)
- [Building Consumer Contracts](https://docs.chain.link/cre/guides/workflow/using-evm-client/onchain-write/building-consumer-contracts)
