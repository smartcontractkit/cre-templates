# Confidential Random — CRE Starter Template (Go)

Generate a fresh random `uint32` on every cron tick inside an AWS Nitro enclave
using the TEE-local randomness API, `rt.Rand()`.

**⚠️ DISCLAIMER**

This template is an educational example to demonstrate how to interact with Chainlink systems, products, and services. It is provided **"AS IS"** and **"AS AVAILABLE"** without warranties of any kind, has **not** been audited, and may omit checks or error handling for clarity. **Do not use this code in production** without performing your own audits and applying best practices. Neither Chainlink Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for unintended outputs generated due to errors in code.

**⚠️ PRIVATE BETA**

[Confidential Workflows](https://docs.chain.link/cre/concepts/confidential-workflows) is in **private beta** and requires enrollment through your Chainlink account team — see [Requesting Confidential Workflows Access](https://docs.chain.link/cre/account/confidential-workflows-access). TEE randomness currently runs in **AWS Nitro in `us-west-2`** (staging).

---

## Overview

The workflow registers a cron trigger with `cre.HandlerInTee`, so the callback
executes inside an attested enclave. On each tick it calls `rt.Rand()` and
returns the generated `uint32`.

```
┌──────────────┐
│  CronTrigger │  fires on schedule (runs on the Workflow DON)
└──────┬───────┘
       │  DON hands the triggered request to an enclave
       v
╔══════════════════════════════════════════════════════════════╗
║  ENCLAVE (AWS Nitro, us-west-2)                              ║
║                                                              ║
║   rt.Rand()  ──>  TEE-local seed for this execution          ║
║   rand.Uint32()  ──>  fresh random value                     ║
╚══════════════════════════════════════════════════════════════╝
```

### Why `rt.Rand()` is recommended

- **Unique per execution** — the enclave derives a new seed for each run.
- **Non-repeating** — consecutive and concurrent executions produce different
  values; workflows running at the same time do not receive the same number.
- **Different from DON randomness** — `rt.Rand()` does not use the DON-seeded
  deterministic PRNG, so it avoids the repeatability constraints of DON
  randomness.

> TEE randomness is useful for raffles, sampling, matchmaking, and similar
> non-cryptographic draws. It is not a replacement for a verifiable random
> function (VRF) when you need an on-chain, publicly verifiable proof.

## Structure

```text
confidential-random-go/
├── .cre/template.yaml        # template metadata for the CLI
├── project.yaml              # RPC settings (not used by this workflow)
├── go.mod
├── README.md
└── my-workflow/
    ├── main.go               # wasip1 entrypoint
    ├── workflow.go           # config, InitWorkflow, onCronTrigger
    ├── config.staging.json
    ├── config.production.json
    └── workflow.yaml
```

## Getting Started

### Prerequisites

- Go 1.25.3 or later
- [CRE CLI](https://docs.chain.link/cre) installed
- Enrollment in the Confidential Workflows private beta, with access to the
  staging region (`us-west-2`)

### 1. Build and test

From `starter-templates/confidential-workflows/confidential-random/confidential-random-go`:

```bash
GOOS=wasip1 GOARCH=wasm go build ./my-workflow
go vet ./...
```

### 2. Simulate

```bash
cre workflow simulate my-workflow --target staging-settings --non-interactive --trigger-index 0
```

The simulator prints a notice that it is not a real TEE and resolves the TEE
requirement (`AWS Nitro in us-west-2`). Run the command several times: the
returned value should differ on every run. Running several simulations at the
same time should also produce distinct values.

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

The template returns `rand.Uint32()`. Limit the range in `onCronTrigger`:

```go
// Uniform in [0, max).
const max = 100
value := rand.Uint32() % max
```

`rand.Uint32() % N` introduces **modulo bias** when `N` does not divide
`2^32` evenly: the lowest `2^32 % N` values are slightly more likely. For
exact uniformity, rejection-sample:

```go
func uniformUint32(r *rand.Rand, max uint32) uint32 {
	limit := math.MaxUint32 - (math.MaxUint32 % max)
	for {
		if v := r.Uint32(); v < limit {
			return v % max
		}
	}
}
```

The returned `*rand.Rand` also supports `Uint64`, `Int63n`, `Intn`, `Float64`,
and the other standard `math/rand` methods. Use one draw per independent
outcome rather than deriving several outcomes from a single value.

## TEE requirements

The handler is constrained to a specific TEE and region:

```go
var teeRequirements = cre.OneOfTees{cre.Nitro{Regions: []cre.NitroRegion{cre.NitroUsWest2}}}
```

This is **staging-specific**. AWS Nitro in `us-west-2` is currently the only
registered TEE type and region. When additional TEEs or regions are registered,
update `teeRequirements` or relax it with `cre.AnyTee{}` /
`cre.AnyTeeInRegions{...}`. Each TEE binding has its own region enum, so using a
region that the TEE does not support is a compile-time error.

## SDK version

`go.mod` pins:

- `github.com/smartcontractkit/cre-sdk-go` to **`v1.22.0-beta.0`** — the first
  release whose `TeeRuntime.Rand()` uses TEE-local randomness. That tag is
  commit `c9308ff292c98d4067518ec9387d2ec89b446560`; earlier releases returned
  DON randomness from `TeeRuntime.Rand()`.
- `github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron` to the
  pseudo-version `v1.3.1-0.20260806122424-24d22c8dc5bd`, because the submodule
  has not had a tagged release since that commit.

If you are testing against a local SDK checkout, add a `go.work` file rather
than editing `go.mod`:

```bash
go work init .
go work edit \
  -replace github.com/smartcontractkit/cre-sdk-go=/path/to/cre-sdk-go \
  -replace github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron=/path/to/cre-sdk-go/capabilities/scheduler/cron
```

Delete `go.work` to return to the pinned versions. `go mod tidy` ignores the
workspace.

The TypeScript sibling in [`../confidential-random-ts`](../confidential-random-ts)
uses the same feature through `Math.random()`; inside a `handlerInTee` handler
the SDK points it at the TEE-local seed. See that variant's README for details.

## Further Reading

- [Confidential Workflows in CRE](https://docs.chain.link/cre/concepts/confidential-workflows)
- [Making a Workflow Confidential](https://docs.chain.link/cre/guides/workflow/using-confidential-workflows)
- [Confidential Workflows Client SDK Reference](https://docs.chain.link/cre/reference/sdk/confidential-workflows-client)

## License

MIT
