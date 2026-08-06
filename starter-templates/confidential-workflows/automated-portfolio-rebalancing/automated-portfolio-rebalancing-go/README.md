# Automated Portfolio Rebalancing CRE Project (Go)

This standalone CRE project implements a confidential portfolio rebalancing workflow for crypto allocations. It is the Go port of `automated-portfolio-rebalancing/` (TypeScript) and behaves identically.

## Description

The workflow continuously tracks allocation drift and triggers rebalancing when policy thresholds are exceeded. It is designed to restore user-defined target weights while protecting sensitive operational inputs, including exchange credentials, model credentials, policy thresholds, and execution preferences inside confidential execution.

## Target Customer

- Professional retail traders
- Developer shops
- Founders building trading products

## Structure

- `go.mod`: module definition (rooted in this workflow directory), pinned to the `cre-sdk-go` `v1.18.0` release (see below)
- `../project.yaml`: project-level target settings (shared with the TypeScript template)
- `../secrets.yaml`: secret ID mappings (shared with the TypeScript template)
- `main.go`: WASM entry point (`//go:build wasip1`)
- `workflow.go`: workflow implementation
- `workflow_test.go`: unit tests, including in-enclave HTTP stubs
- `mock-server.js`: local deterministic API server

## Private Inputs

The following inputs are handled as confidential:

- Exchange API credentials used to read holdings, reserve data, and execution context, including stablecoin reserve depth and cash balance.
- LLM reasoning API credentials.
- Portfolio policy settings such as target allocation mix, minimum drift threshold, maximum trade size per execution, and required stablecoin reserve floor.
- Execution preferences such as venue priority, slippage limits, and trade chunking/ordering behavior.

## Workflow Notes

1. Monitor portfolio state.
   The workflow gathers market prices, current asset weights, drift from target allocations, reserve health, and volatility signals.
2. Enforce user-defined portfolio constraints.
   Confidential reasoning validates weight constraints, drift triggers, rebalance sizing limits, reserve protection requirements, and slippage controls.
3. Build a rebalance action plan.
   The plan can include buying underweight assets, selling overweight assets, enforcing reserve floors, capping per-trade notionals, and optimizing execution through chunking and smart venue routing.
4. Execute rebalance actions across venues.
   Depending on route selection, the workflow can execute both on-chain operations (such as swaps) and off-chain operations (such as centralized exchange API trades).

Note: Reasoning stages can be implemented with deterministic rule-based logic instead of an LLM when a fully rules-driven execution model is preferred.

## Confidentiality boundary

The handler is registered with `cre.HandlerInTee`, so it receives a `cre.TeeRuntime` rather than a `cre.Runtime`. Everything in the handler runs inside the enclave:

- Policy secrets are released by the Vault DON directly into the attested enclave and decrypted at the moment `GetSecret()` runs.
- HTTP calls go through `client.SendRequestInTee(runtime, ...)`, keeping URLs, headers (including API keys) and response bodies confidential from node operators.

What is *not* confidential: the workflow binary itself, including this logic, is provided to the enclave by the Workflow DON and is therefore revealed. What the enclave protects is the *data* the logic computes over. Crossing back with `runtime.UsingTheDons()` leaves the enclave, so anything passed to a capability call there is no longer confidential.

Logging inside the enclave weakens the guarantee. This template logs only non-sensitive markers (`rebalance-getsecret-ok`, trade counts) and `workflow_test.go` asserts that no secret reaches the logs — but you should remove the log lines entirely before deploying to production.

## How the model is constrained

The model does **not** decide trade sizes. `BuildTargetAlignedTrades` derives the required buys and sells deterministically from the drift, and `ReconcileDecisionWithTargets` uses *those* notionals — the model may only:

- influence slippage (clamped to `max_slippage_bps`) and venue selection, and
- veto a trade by disagreeing on its side.

Note the model's `shouldRebalance` flag is deliberately not a veto; the target-aligned set drives execution. `EnforcePolicyAndChunk` then applies the hard caps, splits each trade into `max_trade_usd` chunks, and walks them in order tracking the projected reserve — a sequence that would dip below `reserve_floor_usdc` aborts the whole run rather than executing partially.

### Determinism

`CalculateAllocations` sorts symbols before walking them. This matters beyond tidiness: allocation order flows through to trade order, and `EnforcePolicyAndChunk` checks the running reserve as it goes, so a different order can produce a different reserve-floor outcome. Go map iteration order is randomised and the enclave result is attested and verified by DON consensus, so an unordered walk would be non-deterministic across nodes.

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
- `MOCK_REBALANCING_TARGET_ALLOCATION_BTC_PCT`
- `MOCK_REBALANCING_TARGET_ALLOCATION_ETH_PCT`
- `MOCK_REBALANCING_TARGET_ALLOCATION_USDC_PCT`
- `MOCK_REBALANCING_DRIFT_THRESHOLD_PCT`
- `MOCK_REBALANCING_MAX_TRADE_USD`
- `MOCK_REBALANCING_RESERVE_FLOOR_USDC`
- `MOCK_REBALANCING_MAX_SLIPPAGE_BPS`
- `MOCK_REBALANCING_PREFERRED_VENUES`
- `MOCK_REBALANCING_ORDER_SEQUENCE_PREFERENCE`

The local mock server for this project only exposes routes under `/rebalancing/*`.

## Configuration

`config.staging.json`:

| Field | Description |
|-------|-------------|
| `schedule` | Cron expression (6 fields, seconds first) |
| `mock_base_url` | Base URL for the portfolio, prices, volatility and execution endpoints |
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
cd .. && cre workflow simulate ./automated-portfolio-rebalancing-go --target=staging-settings
```

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
