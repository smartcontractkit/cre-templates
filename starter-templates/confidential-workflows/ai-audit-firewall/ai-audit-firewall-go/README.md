# AI Audit Firewall CRE Project (Go)

This standalone CRE project implements a confidential pre-execution security firewall for smart contract interactions. It is the Go port of `ai-audit-firewall/` (TypeScript) and behaves identically.

## Description

The workflow screens proposed transactions before they are allowed to proceed. It fetches and validates contract intelligence, runs confidential reasoning to classify risk, and then enforces a firewall decision path. Scanner and model credentials remain protected inside confidential execution throughout the process.

## Target Customer

- Professional retail traders
- Developer shops
- Founders building trading products

## Structure

- `go.mod`: module definition (rooted in this workflow directory), with the SDK pinned to an unreleased commit (see below)
- `../project.yaml`: project-level target settings (shared with the TypeScript template)
- `../secrets.yaml`: secret ID mappings (shared with the TypeScript template)
- `../contracts/`: `AuditFirewallConsumer.sol` and `ReceiverTemplate.sol` (shared)
- `main.go`: WASM entry point (`//go:build wasip1`)
- `workflow.go`: workflow implementation
- `workflow_test.go`: unit tests, including in-enclave HTTP and EVM stubs
- `mock-server.js`: local deterministic API server

## Private Inputs

The following inputs are handled as confidential:

- Chain scanner API credentials used for contract metadata retrieval and verification checks.
- LLM reasoning API credentials used for independent audit analysis.

## Workflow Notes

1. Monitor and ingest the proposed interaction.
   The workflow receives candidate transaction context, including token and protocol contract addresses.
2. Fetch and validate contract data confidentially.
   It retrieves source and ABI artifacts through the scanner and verifies scanner credential permissions before trusting fetched data.
3. Run smart contract audit analysis.
   The workflow submits context to two reasoning models and classifies behavior into structured risk signals:
   - `obfuscatedTax`
   - `privilegeEscalation`
   - `externalCallRisk`
   - `logicBomb`
4. Enforce the firewall decision, log the audit, and optionally deliver the verdict on-chain.

## Confidentiality boundary

The handler is registered with `cre.HandlerInTee`, so it receives a `cre.TeeRuntime` rather than a `cre.Runtime`. Everything in the handler runs inside the enclave:

- Scanner and model credentials are released by the Vault DON directly into the attested enclave and decrypted at the moment `GetSecret()` runs.
- HTTP calls go through `client.SendRequestInTee(runtime, ...)`, keeping URLs, headers (including API keys), contract source and model reasoning confidential from node operators.

The one place that leaves the enclave is `writeVerdictOnChain`, which calls `runtime.UsingTheDons()` to sign and deliver the report. Only three values cross that boundary — the verdict code, the risk-flag bitmask and the chain selector. The contract source, model reasoning and scanner credentials never do.

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
| `maxTotalCalls` | 10 | Overall ceiling across all capabilities |
| `http-actions@1.0.0-alpha` / `SendRequest` | 8 | Steady state is 7 calls (proposal, credential check, 2 contract fetches, 2 model calls, audit log, firewall action); 8 leaves one call of headroom |
| `consensus@1.0.0-alpha` / `Report` | 1 | One signed report per run |
| `evm:ChainSelector:<selector>@1.0.0` / `WriteReport` | 1 | One on-chain write per run, scoped to the configured chain |
| `maxSecrets` | 3 | Exact-match only, in the `main` namespace |

The EVM restriction is added only when an EVM target is configured, keeping the on-chain leg opt-in. An unresolvable chain name is skipped here rather than treated as fatal — the write path reports that error at execution time.

Because the set is closed, **adding a capability call to the workflow means raising the matching limit here**, or the run will be cut off by its own restrictions. `TestBuildRestrictions_HTTPBudgetCoversASuccessfulRun` guards against that by asserting a full successful run stays within the declared HTTP budget.

## Default config

`ParseConfig` (wired into `wasm.NewRunner` in `main.go` in place of `cre.ParseJSON`) falls back to `DefaultConfig` when the config payload is empty or whitespace. The DON hands the runner an empty payload when it invokes the pre-hook before a real config is attached; without the fallback that would fail to unmarshal and take the pre-hook down with it.

The default's URLs are deliberately the placeholder string `prehook-default` — the pre-hook only needs the config's *shape* (the secret IDs and the EVM target) to build the restriction set, never real endpoints.

## Onchain delivery

The `evms` array in config is optional; with no entry (or an incomplete one) the workflow skips the onchain leg and returns the verdict only. The report is ABI-encoded as `(uint8 verdictCode, uint8 riskMask, uint64 chainSelector)`, decodable with `abi.decode(report, (uint8, uint8, uint64))`.

Verdict codes: `ALLOW` = 1, `DENY` = 2, `MANUAL_REVIEW` = 3. Risk mask bits: `obfuscatedTax` = 1, `privilegeEscalation` = 2, `externalCallRisk` = 4, `logicBomb` = 8.

To use it, deploy `../contracts/AuditFirewallConsumer.sol` and set the deployed address in `config.staging.json` under `evms[0].consumer_address`. The shipped config uses the zero address as a placeholder — the write is still attempted against it, so set a real address before relying on this leg.

## TEE constraints

The third argument to `cre.HandlerInTee` declares which enclaves the handler accepts:

```go
cre.AnyTee{}                                               // any registered TEE, any region
cre.AnyTeeInRegions{Regions: []cre.Region{cre.AwsUsWest2}} // any TEE, restricted to a region
cre.OneOfTees{cre.Nitro{Regions: []cre.NitroRegion{cre.NitroUsWest2}}} // specific TEEs and regions
```

Each TEE binding owns its own region enum (`Nitro` ↔ `NitroRegion`), so passing a region that the TEE does not support is a **compile-time** error. AWS Nitro in `us-west-2` is currently the only registered TEE type and region.

## Required Environment Variables

Copy `.env.example` to `.env` and provide values for:

- `CRE_ETH_PRIVATE_KEY` (required for the onchain write; optional otherwise)
- `MOCK_PORT`
- `MOCK_SCANNER_API_KEY`
- `MOCK_PRIMARY_LLM_API_KEY`
- `MOCK_SECONDARY_LLM_API_KEY`

The local mock server for this project only exposes routes under `/audit-firewall/*`.

## Configuration

`config.staging.json`:

| Field | Description |
|-------|-------------|
| `schedule` | Cron expression (6 fields, seconds first) |
| `mock_base_url` | Base URL for the proposal, audit-log and firewall-action endpoints |
| `scanner_url` | Scanner base URL for credential and contract lookups |
| `primary_llm_url` | First audit model endpoint, called from inside the enclave |
| `secondary_llm_url` | Second audit model endpoint, given the first model's findings |
| `secrets_ids` | Maps each credential to a secret ID in `secrets.yaml` |
| `evms` | Optional onchain delivery target (chain, consumer address, gas limit) |

## Quick Start

1. Create environment file (at the shared project root)

```bash
cp ../.env.example ../.env
```

2. Start the mock server (requires Node or Bun)

```bash
bun mock-server.js
```

3. In another terminal, run checks

```bash
go vet ./...
go test ./...
```

4. Simulate workflow

```bash
cd .. && cre workflow simulate ./ai-audit-firewall-go --target=staging-settings
```

## Testing against an unreleased SDK

`go.mod` pins the SDK to a commit, which is enough to build and simulate. If you are iterating on the Go SDK itself and want the template to compile against your **local working tree**, add a `go.work` file (already gitignored) rather than editing `go.mod`:

```bash
go work init .
go work edit \
  -replace github.com/smartcontractkit/cre-sdk-go=/path/to/cre-sdk-go \
  -replace github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http=/path/to/cre-sdk-go/capabilities/networking/http \
  -replace github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron=/path/to/cre-sdk-go/capabilities/scheduler/cron \
  -replace github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm=/path/to/cre-sdk-go/capabilities/blockchain/evm
```

`go build` and `go test` honour the workspace; **`go mod tidy` deliberately ignores it**, so run tidy only after you have real versions in `go.mod`. Delete `go.work` to go back to the pinned versions.
