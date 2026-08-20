# AI Audit Firewall CRE Project (Go)

This standalone CRE project implements a confidential pre-execution security firewall for smart contract interactions. It is the Go port of the TypeScript workflow and uses the same live Etherscan and OpenRouter flow.

## Description

The workflow screens a configured proposed transaction before it is allowed to proceed. It fetches verified contract source from Etherscan, runs confidential OpenRouter reasoning to classify risk, and then enforces a firewall decision path. Both API credentials remain protected inside confidential execution.

## Target Customer

- Professional retail traders
- Developer shops
- Founders building trading products

## Structure

- `go.mod`: module definition (rooted in this workflow directory), pinned to the `cre-sdk-go` `v1.18.0` release (see below)
- `../project.yaml`: project-level target settings (shared with the TypeScript template)
- `../secrets.yaml`: secret ID mappings (shared with the TypeScript template)
- `../contracts/`: `AuditFirewallConsumer.sol` and `ReceiverTemplate.sol` (shared)
- `main.go`: WASM entry point (`//go:build wasip1`)
- `workflow.go`: workflow implementation
- `workflow_test.go`: unit tests, including in-enclave HTTP and EVM stubs

## Private Inputs

The following inputs are handled as confidential:

- The Etherscan API credential used to retrieve verified contract source.
- The OpenRouter API credential used for both audit models.

## Workflow Notes

1. Read the proposed interaction from workflow config.
   Proposal config is visible to the DON; do not put secrets in it.
2. Fetch and validate contract data confidentially.
   The workflow makes two calls to `https://api.etherscan.io/v2/api`, one for each configured contract, and denies the proposal if either contract lacks verified source.
   The workflow waits one second between these requests to stay within Etherscan's free-tier rate limit, using CRE DON time rather than a Go local-clock sleep.
3. Run smart contract audit analysis.
   The workflow makes two calls to `https://openrouter.ai/api/v1/chat/completions`. The primary model receives the proposal and verified token contract artifact. The secondary model receives the proposal, verified protocol contract artifact, and primary analysis. The models classify:
   During simulation, the non-sensitive `audit-firewall-primary-model-complete` and `audit-firewall-secondary-model-start` markers identify whether the primary or secondary model call failed.
   - `obfuscatedTax`
   - `privilegeEscalation`
   - `externalCallRisk`
   - `logicBomb`
4. Apply the firewall decision and optionally deliver the verdict on-chain.

Model output is only an advisory signal used by the workflow's deterministic decision logic. OpenRouter provider data collection is denied.

## Confidentiality boundary

The handler is registered with `cre.HandlerInTeeWithPreHook`, so it receives a `cre.TeeRuntime` rather than a `cre.Runtime`. Everything in the handler runs inside the enclave:

- Etherscan and OpenRouter credentials are released by the Vault DON directly into the attested enclave and decrypted when `GetSecret()` runs.
- HTTP calls go through `client.SendRequestInTee(runtime, ...)`, keeping headers, verified contract source and model reasoning confidential from node operators.

The one place that leaves the enclave is `writeVerdictOnChain`, which calls `runtime.UsingTheDons()` to sign and deliver the report. Only three values cross that boundary — the verdict code, the risk-flag bitmask and the chain selector. The contract source, model reasoning and API credentials never do.

What is *not* confidential: the workflow binary itself, including this logic, is provided to the enclave by the Workflow DON and is therefore revealed. What the enclave protects is the *data* the logic computes over.

Logging inside the enclave weakens the guarantee. This template logs only non-sensitive markers and `workflow_test.go` asserts that no credential or contract source reaches the logs — but you should remove the log lines entirely before deploying to production.

## Decision logic

`DetermineVerdict` fails safe in most respects — any raised risk flag denies outright, and low confidence, a `review` recommendation, or disagreement between the models all escalate to `MANUAL_REVIEW` rather than allowing.

> **⚠️ Inherited quirk.** The function branches on whether the models *agree* and how confident they are, but never on *what* they agreed. Two models that confidently and unanimously recommend `deny`, without raising any risk flag, therefore fall through to `ALLOW`. `TestDetermineVerdict_UnanimousDenyWithoutFlagsFallsThroughToAllow` pins this behaviour so the Go port stays faithful to the TypeScript original. If you want to change it, fix both templates by adding an explicit deny check before the agreement check.

An unverified contract short-circuits to `DENY` without consulting a model at all — with no source there is nothing meaningful to audit.

## Restrictions (pre-hook)

The handler is registered with `cre.HandlerInTeeWithPreHook`. The fourth argument is a pre-hook that runs **in the DON, before the enclave executes**, and returns the restrictions for that execution. `BuildRestrictions` declares them:

| Restriction | Value | Why |
|-------------|-------|-----|
| Capability set type | `CLOSED` | Any capability call not listed below is rejected outright |
| `maxTotalCalls` | 6 | Four HTTP calls, one report and one optional on-chain write |
| `http-actions@1.0.0-alpha` / `SendRequest` | 4 | Two Etherscan source fetches followed by two OpenRouter model calls |
| `consensus@1.0.0-alpha` / `Report` | 1 | One signed report per run |
| `evm:ChainSelector:<selector>@1.0.0` / `WriteReport` | 1 | One on-chain write per run, scoped to the configured chain |
| `maxSecrets` | 2 | Exact matches for `etherscan_api_key` and `openrouter_api_key` in the `main` namespace |

The EVM restriction is added only when an EVM target is configured, keeping the on-chain leg opt-in. An unresolvable chain name is skipped here rather than treated as fatal — the write path reports that error at execution time.

Because the set is closed, **adding a capability call to the workflow means raising the matching limit here**, or the run will be cut off by its own restrictions.

## Default config

`ParseConfig` (wired into `wasm.NewRunner` in `main.go` in place of `cre.ParseJSON`) falls back to `DefaultConfig` when the config payload is empty or whitespace. The DON hands the runner an empty payload when it invokes the pre-hook before a real config is attached; without the fallback that would fail to unmarshal and take the pre-hook down with it.

The defaults use Etherscan chain ID `11155111`, low-cost paid primary model `google/gemini-2.5-flash-lite`, low-cost paid secondary model `openai/gpt-4.1-nano`, the two live secret IDs, and an empty `evms` array. Direct tests with the workflow's real prompts returned valid structured output from both models within CRE's 10-second standard HTTP limit. Provider pricing, availability, and latency can change. The service endpoints are fixed in the workflow rather than configurable.

## Onchain delivery

The `evms` array defaults to empty, so the workflow returns the verdict without writing on-chain. To opt in, add a complete EVM target. The report is ABI-encoded as `(uint8 verdictCode, uint8 riskMask, uint64 chainSelector)`, decodable with `abi.decode(report, (uint8, uint8, uint64))`.

Verdict codes: `ALLOW` = 1, `DENY` = 2, `MANUAL_REVIEW` = 3. Risk mask bits: `obfuscatedTax` = 1, `privilegeEscalation` = 2, `externalCallRisk` = 4, `logicBomb` = 8.

To use on-chain delivery, deploy `../contracts/AuditFirewallConsumer.sol`, then add its address, the chain selector name and a gas limit under `evms[0]` in `config.staging.json`.

## TEE constraints

The third argument to `cre.HandlerInTeeWithPreHook` declares which enclaves the handler accepts:

```go
cre.AnyTee{}                                               // any registered TEE, any region
cre.AnyTeeInRegions{Regions: []cre.Region{cre.AwsUsWest2}} // any TEE, restricted to a region
cre.OneOfTees{cre.Nitro{Regions: []cre.NitroRegion{cre.NitroUsWest2}}} // specific TEEs and regions
```

Each TEE binding owns its own region enum (`Nitro` ↔ `NitroRegion`), so passing a region that the TEE does not support is a **compile-time** error. AWS Nitro in `us-west-2` is currently the only registered TEE type and region.

## Required Environment Variables

Copy `../.env.example` to `../.env` and provide the two API secrets mapped by `../secrets.yaml`:

- `ETHERSCAN_API_KEY` → `etherscan_api_key`
- `OPENROUTER_API_KEY` → `openrouter_api_key`

`CRE_ETH_PRIVATE_KEY` is optional for simulation and required only when you opt into an on-chain write.

## Configuration

`config.staging.json`:

| Field | Description |
|-------|-------------|
| `schedule` | Cron expression (6 fields, seconds first) |
| `proposal` | Transaction and contract context to audit; visible to the DON |
| `etherscan_chain_id` | Etherscan V2 chain ID used for both source requests |
| `primary_model` | First OpenRouter audit model; defaults to `google/gemini-2.5-flash-lite` |
| `secondary_model` | Second OpenRouter audit model; defaults to `openai/gpt-4.1-nano` |
| `secrets_ids` | The `etherscan_api_key` and `openrouter_api_key` secret IDs |
| `evms` | Optional on-chain delivery targets; empty by default |

The Etherscan and OpenRouter endpoints are fixed at `https://api.etherscan.io/v2/api` and `https://openrouter.ai/api/v1/chat/completions`.

## Quick Start

1. Create the shared environment file:

```bash
cp ../.env.example ../.env
```

2. Add `ETHERSCAN_API_KEY` and `OPENROUTER_API_KEY` to `../.env`.

3. Run checks:

```bash
go vet ./...
go test ./...
```

4. Simulate the workflow:

```bash
cd .. && cre workflow simulate ./ai-audit-firewall-go --target=staging-settings
```

## Testing against a local SDK checkout

`go.mod` pins `cre-sdk-go` to the `v1.18.0` tag; the `capabilities/networking/http`, `capabilities/scheduler/cron`, and `capabilities/blockchain/evm` submodules are pinned by pseudo-version to that same commit, since they haven't had a tagged release since. If you are iterating on the Go SDK itself and want the template to compile against your **local working tree**, add a `go.work` file (already gitignored) rather than editing `go.mod`:

```bash
go work init .
go work edit \
  -replace github.com/smartcontractkit/cre-sdk-go=/path/to/cre-sdk-go \
  -replace github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http=/path/to/cre-sdk-go/capabilities/networking/http \
  -replace github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron=/path/to/cre-sdk-go/capabilities/scheduler/cron \
  -replace github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm=/path/to/cre-sdk-go/capabilities/blockchain/evm
```

`go build` and `go test` honour the workspace; **`go mod tidy` deliberately ignores it**, so run tidy only after you have real versions in `go.mod`. Delete `go.work` to go back to the pinned versions.
