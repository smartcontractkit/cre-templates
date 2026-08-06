# Hello Confidential Workflows — CRE Starter Template (Go)

Quickstart confidential workflow. Run a handler's callback inside a secure enclave: fetch a secret from the Vault DON, call an API from inside the enclave, execute decision logic over the confidential data such as Vault DON secrets or HTTP response payloads, then cross back to the Workflow DON for consensus and DON capability calls. 

**⚠️ DISCLAIMER**

This template is an educational example to demonstrate how to interact with Chainlink systems, products, and services. It is provided **"AS IS"** and **"AS AVAILABLE"** without warranties of any kind, has **not** been audited, and may omit checks or error handling for clarity. **Do not use this code in production** without performing your own audits and applying best practices. Neither Chainlink Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for unintended outputs generated due to errors in code.

**⚠️ PRIVATE BETA**

[Confidential Workflows](https://docs.chain.link/cre/concepts/confidential-workflows) is in **private beta** and requires enrollment through your Chainlink account team — see [Requesting Confidential Workflows Access](https://docs.chain.link/cre/account/confidential-workflows-access).

---

## Overview

By default, a CRE workflow's callback runs on Workflow DON nodes, where node operators can in principle inspect the data it is computing over. That's fine for most workflows — but some logic needs to be computed over sensitive data while preserving the confidentiality of that data: risk thresholds, API credentials, centralised exchange stablecoin reserves for reasoning, identity details. Leaking this data could have adverse effects, including enabling front-running attacks, exposing sensitive financial information, and compromising individual privacy.

A **Confidential Workflow** moves that computation into a hardware-isolated [enclave](https://docs.chain.link/cre/key-terms#enclave). This template is the minimal end-to-end shape of one, in four steps:

| Step | What it demonstrates | API |
|------|----------------------|-----|
| 1 | Registers a TEE handler to execute in the enclave | `cre.HandlerInTee(trigger, callback, tees)` |
| 2 | Securely fetches a secret inside the enclave | `runtime.GetSecret(&cre.SecretRequest{Id: ...})` |
| 3 | Executes a capability call from within the enclave | `client.SendRequestInTee(runtime, req)` |
| 4 | Returns to the DON for any operations requiring decentralized consensus | `runtime.UsingTheDons()` |

> This is the Go port of [`hello-confidential-workflows-ts`](../hello-confidential-workflows-ts). Both produce the same result for the same input.

### Use Cases

- **Automated liquidation protection**: Automatically protect DeFi lending positions by continuously monitoring liquidation risk and executing collateral management, debt repayment while preserving the confidentiality of centralized exchange as well as LLM API keys, proprietary risk management thresholds, and execution preferences.
- **Automated portfolio rebalancing**: Automatically rebalance crypto portfolios by continuously monitoring allocation drift and executing portfolio adjustments when predefined thresholds are exceeded, while preserving the confidentiality of exchange API keys, LLM reasoning, portfolio allocation thresholds, and execution preferences.
- **AI smart contract audit firewall**: Automatically analyze and screen smart contract interactions before execution to detect and block malicious transactions, while preserving the confidentiality of chain scanner and LLM reasoning API credentials.

## Architecture

```
┌──────────────┐
│  CronTrigger │  fires on schedule (runs on the Workflow DON)
└──────┬───────┘
       │  DON hands the triggered request to an enclave
       v
╔══════════════════════════════════════════════════════════════╗
║  ENCLAVE (TEE)                                               ║
║  Data below is kept confidential from node operators.        ║
║  The binary — including this logic — is NOT confidential.    ║
║                                                              ║
║   Step 2: runtime.GetSecret(&cre.SecretRequest{...})         ║
║             ▲                                                ║
║             └── released by Vault DON, decrypted in-enclave  ║
║                                                              ║
║   Step 3: client.SendRequestInTee(runtime, &http.Request{})  ║
║             Authorization: Bearer <secret>                   ║
║             ▲ request + response payloads stay confidential  ║
║                                                              ║
║   Logic over confidential data:                              ║
║             scoreResponse(body) vs. ScoreThreshold           ║
║             -> verdict = APPROVE | REJECT                    ║
╚═══════════════════════════╤══════════════════════════════════╝
                            │  Step 4: runtime.UsingTheDons()
                            │  ONLY the verdict + score cross out
                            v
┌──────────────────────────────────────────────────────────────┐
│  WORKFLOW DON — donRuntime.GenerateReport(...)                │
│  Consensus verifies the enclave attestations, proving the     │
│  integrity of the logic executed in the enclave, then signs   │
└──────────────────────────────────────────────────────────────┘
```

## What the workflow does

`my-workflow/workflow.go`:

1. **Registers the cron handler with `cre.HandlerInTee`**, constrained to `cre.OneOfTees{cre.Nitro{Regions: []cre.NitroRegion{cre.NitroUsWest2}}}`
2. **Fetches `API_TOKEN`** with `runtime.GetSecret()` — the Vault DON releases it only into an attested enclave, and it's decrypted at the moment the call runs
3. **Calls the configured URL** with `client.SendRequestInTee(runtime, ...)`, so the request executes from inside the enclave with the secret in the `Authorization` header
4. **Scores the response** against `ScoreThreshold` — decision logic executed over confidential data. The data it reads (the secret and the response payload) stays confidential from node operators; the logic itself is part of the binary and is not
5. **Crosses back with `UsingTheDons()`** and generates a signed report containing only the verdict and score — never the secret or the raw response

The default endpoint is `https://postman-echo.com/headers`, which echoes the request headers back — no signup or real API key needed. The workflow uses that to confirm the secret really was injected inside the enclave, reporting it as the boolean `secret reached API: true` rather than by logging the token. Note that it never logs the response body either; the confidentiality boundary is the reason, and it's worth keeping that habit even in simulation.

## Getting Started

### Prerequisites

- Go 1.25.3 or later
- [CRE CLI](https://docs.chain.link/cre) installed
- Enrollment in the Confidential Workflows private beta (required to **deploy**; see the note above)

### 1. Configure Secrets

```bash
cp .env.example .env
```

Then set `SECRET_API_TOKEN` in `.env`. `secrets.yaml` maps the workflow-facing secret ID `API_TOKEN` to that environment variable:

```yaml
secretsNames:
    API_TOKEN:
        - SECRET_API_TOKEN
```

With the default echo endpoint any non-empty value works.

### 2. Run Tests

```bash
go test ./my-workflow/
```

### 3. Simulate

```bash
cre workflow simulate my-workflow --target staging-settings --non-interactive --trigger-index 0
```

Expected output:

```
2026-01-01T00:00:00Z [SIMULATION] Running trigger trigger=cron-trigger@1.0.0
╭────────────────────────────────────────────────────────────────────────────────────────────────────╮
│ Trigger requested TEE Execution your trigger will run in one of the following Tees:                │
│     - AWS Nitro in us-west-2                                                                       │
│ The simulator is not a real TEE, and is meant to debug.                                            │
│ Do not use it for sensitive information.                                                           │
│ During real execution, user logs for this trigger will not be visible, and will not leave the TEE. │
│ They are presented in the simulator for debugging only.                                            │
╰────────────────────────────────────────────────────────────────────────────────────────────────────╯

2026-01-01T00:00:00Z [USER LOG] msg="Enclave computation complete" verdict=REJECT

✓ Workflow Simulation Result:
"REJECT (score: 371, secret reached API: true)"
```

Three things to notice:

- The simulator confirms the TEE constraint it resolved (`AWS Nitro in us-west-2`) and warns that **it is not a real enclave** — logs are shown for debugging only. In real execution those logs never leave the TEE.
- `secret reached API: true` means the Vault DON secret was fetched inside the enclave and arrived in the outbound request's `Authorization` header.
- The verdict flips between `APPROVE` and `REJECT` from run to run. That's expected: the score is derived from the live response body, and the echo endpoint includes a per-request trace ID. Lower `scoreThreshold` to see `APPROVE` consistently.

## Testing against a local SDK checkout

`go.mod` pins `cre-sdk-go` to the `v1.18.0` tag; the `capabilities/networking/http` and `capabilities/scheduler/cron` submodules are pinned by pseudo-version to that same commit, since they haven't had a tagged release since. If you are iterating on the Go SDK itself and want the template to compile against your **local working tree**, add a `go.work` file (already gitignored) rather than editing `go.mod`:

```bash
go work init .
go work edit \
  -replace github.com/smartcontractkit/cre-sdk-go=/path/to/cre-sdk-go \
  -replace github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http=/path/to/cre-sdk-go/capabilities/networking/http \
  -replace github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron=/path/to/cre-sdk-go/capabilities/scheduler/cron
```

`go build` and `go test` honour the workspace; **`go mod tidy` deliberately ignores it**, so run tidy only after you have real versions in `go.mod`. Delete `go.work` to go back to the pinned versions.

## Configuration

`my-workflow/config.staging.json`:

| Field | Description |
|-------|-------------|
| `schedule` | Cron expression (6 fields, seconds first) |
| `url` | Endpoint called from inside the enclave |
| `secretId` | Secret ID fetched with `GetSecret()`; must match `secrets.yaml` |
| `scoreThreshold` | Threshold the in-enclave scoring compares against |

## TEE constraints

The third argument to `cre.HandlerInTee` declares which enclaves the handler accepts:

```go
cre.AnyTee{}                                               // any registered TEE, any region
cre.AnyTeeInRegions{Regions: []cre.Region{cre.AwsUsWest2}} // any TEE, restricted to a region
cre.OneOfTees{cre.Nitro{Regions: []cre.NitroRegion{cre.NitroUsWest2}}} // specific TEEs and regions
```

Each TEE binding owns its own region enum (`Nitro` ↔ `NitroRegion`), so passing a region that the TEE does not support is a **compile-time** error. AWS Nitro in `us-west-2` is currently the only registered TEE type and region.

The constraint can also come from config: `cre.Tees` is a JSON wrapper whose `UnmarshalJSON` picks the matching variant and fails closed on unknown shapes.

## Confidentiality boundary

Understanding what is and isn't protected matters more here than in a regular workflow.

| Protected by default | **Not** automatically protected |
|----------------------|--------------------------------|
| Secrets the Vault DON releases into the enclave | Triggers, chain reads, and chain writes — these always run on Workflow DON nodes |
| Request and response payloads of HTTP calls made from the enclave | Your workflow's **source code and deployed binary, including the logic executed in the enclave** |
| Sensitive inputs and intermediate values you don't share outside the enclave | Capability calls not routed through the enclave |
| Enclave execution memory, while your computation runs | Reports, calldata, and any output you deliver outside the enclave |

Consequences worth internalizing:

- **The logic is not confidential — the data is.** A confidential workflow, despite running inside the enclave, is part of the binary the Workflow DON provides to the enclave, so that binary including the enclave logic is revealed. What running that logic in the enclave currently provides is confidentiality of the *data* it computes over: Vault DON secrets such as API keys, the request and response payloads of HTTP calls made from the enclave, and other intermediate values. Confidential logic is on the future roadmap, not part of the current beta.
- **`UsingTheDons()` is a one-way door.** Anything you pass into a capability call on that runtime executes on Workflow DON nodes like any non-confidential call. Cross over only the data that does not need to stay confidential.
- **Logs are for simulation only.** Every `runtime.Logger()` call inside the enclave MUST be removed before deploying to production to preserve the confidentiality offered by enclaves — and logging should be avoided for sensitive and non-sensitive values alike.
- **Keep enclave logic deterministic.** The Workflow DON verifies enclave attestations and reaches consensus before the workflow completes successfully.
- **Multiple confidential workflows may execute within the same enclave.** Workflows are isolated from one another by the wasmtime. Dedicated per workflow enclave isolation is planned as a future enhancement.

## Customization

- **Put your real logic in the enclave**: replace `scoreResponse` in `workflow.go` with the decision logic you need to run over confidential data. Remember that the logic itself is revealed as part of the binary — what the enclave preserves is the confidentiality of the secrets and payloads it reads
- **Deliver the report on-chain**: pass the report from Step 4 to `evmClient.WriteReport(donRuntime, ...)` — the RPCs in `project.yaml` are already set up for Sepolia. See the [Custom Data Feed](../../custom-data-feed) template for the full write path
- **Change the trigger**: `cre.HandlerInTee` accepts any CRE trigger, same as `cre.Handler` — swap cron for a log trigger to react to on-chain events confidentially
- **Fetch several secrets at once**: `runtime.GetSecrets([]*cre.SecretRequest{...})` batches them in a single call
- **Add a pre-hook**: `cre.HandlerInTeeWithPreHook` runs a function on the DON before execution to produce per-execution restrictions

## Which secrets belong in an enclave?

Not every secret needs enclave-level protection.

**Higher-value — consider enclave execution:** wallet and CA private keys; exchange, custody, payment-processor, banking, or LLM-provider credentials; OAuth client secrets, JWT signing keys, KMS keys; payment data, health data, other PII.

**Lower-value — regular DON execution is usually fine:** API keys for publicly available data (weather, explorers, public price feeds, public RPCs); public wallet addresses.

The common thread: a secret belongs in the enclave if disclosure would expose more than the workflow needs.

## Security

- Never commit `.env` files or secrets — `.gitignore` covers `*.env`
- Remove every `runtime.Logger()` call inside the TEE handler before deploying to production
- Audit what crosses `UsingTheDons()`; that data is no longer confidential
- Do not treat the enclave logic as secret — the binary that contains it is provided to the enclave by the Workflow DON and is revealed

## Further Reading

- [Confidential Workflows in CRE](https://docs.chain.link/cre/concepts/confidential-workflows) — concepts and use cases
- [Making a Workflow Confidential](https://docs.chain.link/cre/guides/workflow/using-confidential-workflows) — step-by-step guide
- [Confidential Workflows Client SDK Reference](https://docs.chain.link/cre/reference/sdk/confidential-workflows-client) — full API
- [Confidential HTTP](https://docs.chain.link/cre/capabilities/confidential-http) — for a single outbound request, without a full confidential handler
- [confidential-compute-examples](https://github.com/smartcontractkit/confidential-compute-examples) — production-shaped reference workflows

## License

MIT
