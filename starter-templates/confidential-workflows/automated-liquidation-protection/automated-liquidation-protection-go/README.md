# Automated Liquidation Protection CRE Project (Go)

This standalone CRE project implements a confidential liquidation-defense workflow for DeFi lending positions. It is the Go port of `automated-liquidation-protection/` (TypeScript) and behaves identically.

## Description

The workflow continuously evaluates borrower risk and takes defensive action before a position becomes unsafe. During periods of high volatility, it can increase collateral, reduce debt, or combine both strategies based on policy constraints. Sensitive operational data remains protected in confidential execution, including exchange credentials, model credentials, and user-defined risk thresholds.

## Target Customer

- Professional retail traders
- Developer shops
- Founders building trading products

## Structure

- `go.mod`: module definition (rooted in this workflow directory), with the SDK pinned to an unreleased commit (see below)
- `../project.yaml`: project-level target settings (shared with the TypeScript template)
- `../secrets.yaml`: secret ID mappings (shared with the TypeScript template)
- `main.go`: WASM entry point (`//go:build wasip1`)
- `workflow.go`: workflow implementation
- `workflow_test.go`: unit tests, including in-enclave HTTP stubs
- `mock-server.js`: local deterministic API server

## Private Inputs

The following inputs are treated as confidential:

- Exchange API credentials used to fetch account context such as stablecoin reserves and available cash balance.
- LLM reasoning API credentials and policy parameters used to govern defense behavior, including minimum and target health factors, reserve deployment caps, minimum reserve balance, and collateral allocation limits.
- Execution preferences that define how defense actions should be sequenced.

## Workflow Notes

1. Observe liquidation risk signals.
   The workflow tracks collateral and debt asset pricing, health factor, liquidation proximity, LTV, liquidation threshold, and market volatility.
2. Enforce user policy constraints.
   Confidential reasoning evaluates how much capital can be deployed, whether debt reduction should be prioritized, which reserve assets are eligible, and what execution sequence is preferred.
3. Select defense actions.
   The workflow builds a response plan that may include collateral-focused moves (deposit, bridge, swap-then-deposit) and debt-focused moves (repay, swap-then-repay, partial payoff, full payoff).
4. Execute the approved defense plan.

Note: Every reasoning stage can be implemented with deterministic rule-based logic instead of an LLM, if your deployment requires a fully rules-driven policy engine.

## Confidentiality boundary

The handler is registered with `cre.HandlerInTee`, so it receives a `cre.TeeRuntime` rather than a `cre.Runtime`. Everything in the handler runs inside the enclave:

- Policy secrets are released by the Vault DON directly into the attested enclave and decrypted at the moment `GetSecret()` runs.
- HTTP calls go through `client.SendRequestInTee(runtime, ...)`, keeping URLs, headers (including API keys) and response bodies confidential from node operators.

What is *not* confidential: the workflow binary itself, including this logic, is provided to the enclave by the Workflow DON and is therefore revealed. What the enclave protects is the *data* the logic computes over. Crossing back with `runtime.UsingTheDons()` leaves the enclave, so anything passed to a capability call there is no longer confidential.

Logging inside the enclave weakens the guarantee. This template logs only non-sensitive markers (`liquidation-getsecret-ok`, action counts) and `workflow_test.go` asserts that no secret reaches the logs — but you should remove the log lines entirely before deploying to production.

## Policy enforcement

The model only ever *proposes* actions. `EnforcePolicy` is the guardrail between the model and the executor: every amount is clamped to the policy caps, and the running reserve is checked against the floor, so a hallucinated or adversarial suggestion cannot drain the vault. A floor breach aborts the whole run rather than executing a partial sequence. Unrecognised action types are dropped.

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

- `CRE_ETH_PRIVATE_KEY` (optional for local simulate)
- `MOCK_PORT`
- `MOCK_EXCHANGE_API_KEY`
- `MOCK_OPENAI_API_KEY`
- `MOCK_LIQUIDATION_WARNING_ACTION_THRESHOLD`
- `MOCK_LIQUIDATION_MINIMUM_HEALTH_FACTOR`
- `MOCK_LIQUIDATION_TARGET_HEALTH_FACTOR`
- `MOCK_LIQUIDATION_MAX_STABLECOIN_RESERVE_DEPLOYMENT`
- `MOCK_LIQUIDATION_MIN_STABLECOIN_RESERVE_BALANCE`
- `MOCK_LIQUIDATION_MAX_COLLATERAL_ALLOCATION`
- `MOCK_LIQUIDATION_MAX_PARTIAL_DEBT_REPAYMENT`
- `MOCK_LIQUIDATION_DEFENSIVE_ACTION_SEQUENCING_PREFERENCE`
- `MOCK_LIQUIDATION_PREFERRED_VENUES`

The local mock server for this project only exposes routes under `/liquidation/*`.

## Configuration

`config.staging.json`:

| Field | Description |
|-------|-------------|
| `schedule` | Cron expression (6 fields, seconds first) |
| `mock_base_url` | Base URL for the risk-state and execute-defense endpoints |
| `openai_url` | Reasoning endpoint called from inside the enclave |
| `openai_model` | Model name passed to the reasoning endpoint |
| `secrets_ids` | Maps each policy input to a secret ID in `secrets.yaml` |

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
cd .. && cre workflow simulate ./automated-liquidation-protection-go --target=staging-settings
```

## Testing against an unreleased SDK

`go.mod` pins the SDK to a commit, which is enough to build and simulate. If you are iterating on the Go SDK itself and want the template to compile against your **local working tree**, add a `go.work` file (already gitignored) rather than editing `go.mod`:

```bash
go work init .
go work edit \
  -replace github.com/smartcontractkit/cre-sdk-go=/path/to/cre-sdk-go \
  -replace github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http=/path/to/cre-sdk-go/capabilities/networking/http \
  -replace github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron=/path/to/cre-sdk-go/capabilities/scheduler/cron
```

`go build` and `go test` honour the workspace; **`go mod tidy` deliberately ignores it**, so run tidy only after you have real versions in `go.mod`. Delete `go.work` to go back to the pinned versions.
