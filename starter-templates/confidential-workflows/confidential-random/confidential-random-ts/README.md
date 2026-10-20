# Confidential Random — CRE Starter Template (TypeScript)

Generate a fresh random `uint32` on every cron tick inside an AWS Nitro enclave
using the confidential workflow's TEE-local `Math.random()` seed.

**⚠️ DISCLAIMER**

This template is an educational example to demonstrate how to interact with Chainlink systems, products, and services. It is provided **"AS IS"** and **"AS AVAILABLE"** without warranties of any kind, has **not** been audited, and may omit checks or error handling for clarity. **Do not use this code in production** without performing your own audits and applying best practices. Neither Chainlink Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for unintended outputs generated due to errors in code.

**⚠️ PRIVATE BETA**

[Confidential Workflows](https://docs.chain.link/cre/concepts/confidential-workflows) is in **private beta** and requires enrollment through your Chainlink account team — see [Requesting Confidential Workflows Access](https://docs.chain.link/cre/account/confidential-workflows-access). TEE randomness currently runs in **AWS Nitro in `us-west-2`** (staging).

---

## Overview

The workflow registers a cron trigger with `handlerInTee`, so the callback
executes inside an attested enclave. On each tick it generates a random
`uint32` from `Math.random()` and returns it as a string.

```
┌──────────────┐
│  CronTrigger │  fires on schedule (runs on the Workflow DON)
└──────┬───────┘
       │  DON hands the triggered request to an enclave
       v
╔══════════════════════════════════════════════════════════════╗
║  ENCLAVE (AWS Nitro, us-west-2)                              ║
║                                                              ║
║   Math.random()  ──>  TEE-local seed for this execution      ║
║   >>> 0          ──>  fresh random uint32                    ║
╚══════════════════════════════════════════════════════════════╝
```

### Why `Math.random()` is recommended here

The TypeScript SDK does not expose a dedicated `runtime.rand()` method. Instead,
when a handler is registered with `handlerInTee`, the SDK switches the host's
random mode to the TEE (`Mode.NODE`) before invoking the callback. `Math.random()`
inside that callback therefore uses the **TEE-local seed** rather than the
deterministic DON seed.

- **Unique per execution** — the enclave derives a new seed for each run.
- **Non-repeating** — consecutive and concurrent executions produce different
  values; workflows running at the same time do not receive the same number.
- **Different from DON randomness** — the TEE seed is not the DON-seeded
  deterministic PRNG used by regular (non-TEE) handlers.

> In Go, the equivalent is `rt.Rand()` on a `cre.TeeRuntime`. See the
> [`../confidential-random-go`](../confidential-random-go) sibling.
>
> TEE randomness is useful for raffles, sampling, matchmaking, and similar
> non-cryptographic draws. It is not a replacement for a verifiable random
> function (VRF) when you need an on-chain, publicly verifiable proof.

## Structure

```text
confidential-random-ts/
├── .cre/template.yaml        # template metadata for the CLI
├── project.yaml              # RPC settings (not used by this workflow)
├── README.md
└── my-workflow/
    ├── main.ts               # config, initWorkflow, onCronTrigger
    ├── config.staging.json
    ├── config.production.json
    ├── package.json
    ├── tsconfig.json
    └── workflow.yaml
```

## Getting Started

### Prerequisites

- [Bun](https://bun.sh/) 1.2.21 or later (the SDK requires it for compilation)
- Node.js 24 or later for the CRE CLI tooling
- [CRE CLI](https://docs.chain.link/cre) installed
- Enrollment in the Confidential Workflows private beta, with access to the
  staging region (`us-west-2`)

### 1. Install and typecheck

From `starter-templates/confidential-workflows/confidential-random/confidential-random-ts/my-workflow`:

```bash
bun install
bun run typecheck
```

### 2. Simulate

From `starter-templates/confidential-workflows/confidential-random/confidential-random-ts`:

```bash
cre workflow simulate my-workflow --target staging-settings --non-interactive --trigger-index 0
```

The simulator prints a notice that it is not a real TEE and resolves the TEE
requirement (`AWS Nitro in us-west-2`). Run the command several times: the
returned value should differ on every run. Concurrent executions should also
produce distinct values.

### 3. Deploy and run

Deploy the workflow to the Workflow Registry:

```bash
cre workflow deploy my-workflow --target staging-settings
```

Activate it (if your deployment did not already do so):

```bash
cre workflow activate my-workflow --target staging-settings
```

Once active, the cron trigger schedules executions automatically. Check
deployment health and recent executions with:

```bash
cre workflow get my-workflow --target staging-settings
```

## Configuration

`my-workflow/config.staging.json`:

```json
{ "schedule": "*/30 * * * * *" }
```

| Field | Description |
|-------|-------------|
| `schedule` | Cron expression with six fields (seconds first). The default fires every 30 seconds. |

## Customizing the random value

The template returns a full-range `uint32`:

```ts
const value = Math.floor(Math.random() * 0x1_0000_0000) >>> 0;
```

Limit the range in `onCronTrigger`:

```ts
// Uniform in [0, max).
const max = 100;
const value = Math.floor(Math.random() * max);
```

The modulo variant `Math.floor(Math.random() * 0x1_0000_0000) % N` introduces
**modulo bias** when `N` does not divide `2^32` evenly: the lowest `2^32 % N`
values are slightly more likely. For exact uniformity, rejection-sample:

```ts
const uniformUint32 = (max: number): number => {
  const limit = 0x1_0000_0000 - (0x1_0000_0000 % max);
  for (;;) {
    const value = Math.floor(Math.random() * 0x1_0000_0000);
    if (value < limit) {
      return value % max;
    }
  }
};
```

Draw one value per independent outcome rather than deriving several outcomes
from a single value.

## TEE requirements

The handler is constrained to a specific TEE and region:

```ts
const teeRequirements: TeeConstraint = [
  { tee: "nitro", regions: [...NITRO_REGIONS] },
];
```

This is **staging-specific**. AWS Nitro in `us-west-2` is currently the only
registered TEE type and region. When additional TEEs or regions are registered,
update `teeRequirements` — each binding is validated at runtime, so an
unsupported TEE/region pair fails before the handler runs.

## SDK version

`package.json` pins `@chainlink/cre-sdk` to **`1.24.0-alpha.1`**, the first
release that switches `Math.random()` to the TEE seed for `handlerInTee`
handlers (the `latest` tag, `1.23.0`, still uses the DON seed). It also includes
the current `TeeConstraint` types (`NITRO_REGIONS`, `AnyTeeConstraint`, etc.).

Repin to the `latest` tag once the TEE random-mode behavior ships outside the
alpha channel.

## Further Reading

- [Confidential Workflows in CRE](https://docs.chain.link/cre/concepts/confidential-workflows)
- [Making a Workflow Confidential](https://docs.chain.link/cre/guides/workflow/using-confidential-workflows)
- [Confidential Workflows Client SDK Reference](https://docs.chain.link/cre/reference/sdk/confidential-workflows-client)

## License

MIT
