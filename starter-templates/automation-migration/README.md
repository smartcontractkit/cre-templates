# CRE Automation Migration Template

This template provides a modernized, type-safe path for migrating existing **Chainlink Automation** upkeeps to **Chainlink Runtime Environment (CRE)**.

## Architecture

This migration uses a "Bridge" pattern:
1. **AutomationReceiver.sol**: A bridge contract deployed once on-chain. It receives CRE reports and forwards the execution to your existing legacy contracts. It enforces two independent authorization layers (see [Security Model](#security-model)).
2. **Workflows**: CRE workflows that poll your legacy `checkUpkeep` or `checkLog` functions and trigger the bridge when needed. The `CUSTOM`/`LOG` paths use generated, type-safe contract bindings; the `CRON` path ABI-encodes a configured function signature at runtime.

**Result: You can usually migrate to CRE without rewriting your original upkeep logic** — provided your legacy contract lets you re-point who is allowed to call it (see [Compatibility](#compatibility)). If it checks `msg.sender`, uses an Automation Forwarder allowlist, or has role-based permissions, authorize the deployed `AutomationReceiver` there before relying on the workflow.

### Security Model
The receiver authorizes reports on two separate layers, and **both** apply:
- **Inbound — who may deliver a report**: the CRE Forwarder address is set at construction and validated by `ReceiverTemplate`. `AutomationReceiver._processReport` adds two additional hard guards: (1) it rejects any delivery if the forwarder was ever set to `address(0)` post-deployment (closing the gap in `ReceiverTemplate.setForwarderAddress`), and (2) it requires at least one complete workflow identity option to be configured — an unconfigured receiver rejects all reports with `WorkflowIdentityNotConfigured`. Two options are accepted: **(a)** workflowId is set (binds the receiver to one specific workflow), or **(b)** both workflowOwner and workflowName are set (binds to a named workflow from a specific owner). Either piece of option (b) alone is insufficient.
- **Outbound — what a report may make the receiver do** (`AutomationReceiver`): a **closed-by-default allowlist** of `(target, function-selector)` pairs. Inbound checks only prove a report came from your workflow; they do **not** constrain the `(target, data)` it carries. Until you allowlist a pair with `setCallAllowed`, the receiver will reject it.

On top of these layers, the owner can trigger a global **emergency pause** (OpenZeppelin `Pausable`) via `pause(bool retryable)`. While paused, all deliveries are rejected; the `retryable` flag chosen at pause time decides whether they stay retryable (revert) or are dropped (consumed). See [Emergency Pause](#3d-emergency-pause-optional).

---

## Migration Path Mapping

| Your Automation Setup | CRE Equivalent | Migration Type |
|---|---|---|
| **Time-based upkeep** | Cron Trigger | `CRON` |
| **Custom logic upkeep** | Cron + EVM Read | `CUSTOM` |
| **Log trigger upkeep** | EVM Log Trigger | `LOG` |

---

## Getting Started

### 1. Initialize the Template
After this branch is merged into the template source, initialize it with:

```bash
cre init --template=automation-migration-ts --project-name my-automation-migration --workflow-name my-workflow
```

While working directly from this branch, copy or open `starter-templates/automation-migration`.

### 2. Build and Deploy the Bridge
Deploy `AutomationReceiver.sol` to your target chain. You only need to do this once to support multiple upkeeps.

The `contracts/evm` directory ships a ready-to-use [Foundry](https://book.getfoundry.sh/) project: a `foundry.toml` and a vendored copy of the only OpenZeppelin files used (`Ownable`, `Pausable`, `Context`). No `forge install` or local node is required — just build and deploy:

```bash
cd contracts/evm
forge build
forge test          # 51 tests covering the receiver + permission template

# Deploy, passing the CRE Forwarder address for your DON as the constructor arg
forge create src/AutomationReceiver.sol:AutomationReceiver \
  --rpc-url "$RPC_URL" \
  --private-key "$PRIVATE_KEY" \
  --constructor-args "$FORWARDER_ADDRESS"
```

- `FORWARDER_ADDRESS` is the CRE Forwarder for your target network — it is the **only** address allowed to call `onReport`. Look it up in the [CRE documentation](https://docs.chain.link/cre) for your network; do not guess it. A wrong value means the DON's reports are rejected. The constructor blocks `address(0)`, and even if the owner ever sets the forwarder to `address(0)` post-deployment (via `setForwarderAddress`), `AutomationReceiver._processReport` will reject all subsequent deliveries with `InvalidForwarderAddress`.

### 3. Configure and Authorize the Receiver

The receiver rejects every outbound call until you allowlist it. You must set up both the call allowlist and optional workflow identity checks.

#### 3a. Allow Upkeep Calls

For each migrated upkeep, allowlist the exact `(target, selector)` the workflow will invoke:

```bash
# Custom-logic / log-trigger upkeeps call performUpkeep(bytes)
cast send "$RECEIVER_ADDRESS" \
  "setCallAllowed(address,bytes4,bool)" \
  "$TARGET_ADDRESS" "$(cast sig 'performUpkeep(bytes)')" true \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"

# Time-based upkeeps call your specific function, e.g. performAction(uint256)
cast send "$RECEIVER_ADDRESS" \
  "setCallAllowed(address,bytes4,bool)" \
  "$TARGET_ADDRESS" "$(cast sig 'performAction(uint256)')" true \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"
```

**Parameters:**
- `target`: The address of your existing upkeep contract. When allowing (`allowed = true`) it must be a deployed contract — `setCallAllowed` reverts with `TargetHasNoCode` if the address has no code (EOA, mistyped address, or never-deployed contract). The code check is skipped when revoking (`allowed = false`) so an entry can always be removed, even if the target has since self-destructed.
- `selector`: The 4-byte function selector (computed from the function signature via `cast sig`).
- `allowed`: Set to `true` to allow, `false` to revoke.

The selector must match the function the workflow encodes (`performUpkeep` for `CUSTOM`/`LOG`, or your `targetFunction` for `CRON`). A mismatch makes `onReport` revert with `CallNotAllowed`.

#### 3b. Configure the Consumer Gas Limit (Recommended)

Set the minimum gas the receiver must have available before forwarding the call to your upkeep consumer contract. When configured, `_processReport` reverts with `InsufficientGas` — causing the CRE Forwarder to record the delivery as **failed and retryable** — if the incoming gas is below `consumerGasLimit + consumerGasLimit / 63 + 7,000` (the on-chain overhead).

The limit is configured **per `(target, selector)` pair** — each allowlisted call gets its own gas guard. Set it to the estimated gas limit for that specific function. Zero (the default for every pair) disables the guard and preserves fire-and-forget semantics.

Note that this helps to mirror Chainlink Automation's fire-and-forget behavior, where a failed `performUpkeep` simply ends that round and the next trigger re-evaluates eligibility.

```bash
# Custom-logic / log-trigger upkeeps: performUpkeep(bytes)
cast send "$RECEIVER_ADDRESS" \
  "setConsumerGasLimit(address,bytes4,uint256)" \
  "$TARGET_ADDRESS" "$(cast sig 'performUpkeep(bytes)')" "$PERFORM_GAS_LIMIT" \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"

# Time-based upkeeps: your specific function, e.g. performAction(uint256)
cast send "$RECEIVER_ADDRESS" \
  "setConsumerGasLimit(address,bytes4,uint256)" \
  "$TARGET_ADDRESS" "$(cast sig 'performAction(uint256)')" "$PERFORM_GAS_LIMIT" \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"
```

**Parameters:**
- `target`: The address of the upkeep contract the limit applies to.
- `selector`: The 4-byte function selector the limit applies to (same value used in `setCallAllowed`).
- `gasLimit`: The minimum gas required by that specific call. Set to `0` to disable the guard for that pair.

#### 3c. Set Workflow Identity Checks (Optional but Recommended for Production)

**At least one complete identity option must be configured before the receiver will accept any report.** Two options are supported:

- **Option A — workflowId**: set the workflow ID; owner and name are not required.
- **Option B — workflowOwner + workflowName**: set both the owner address and the workflow name together; either piece alone is insufficient, and workflowId is not required.

Without a complete option, `_processReport` reverts with `WorkflowIdentityNotConfigured` on every delivery attempt.

**Option A — identify by workflow ID:**

```bash
cast send "$RECEIVER_ADDRESS" \
  "setExpectedWorkflowId(bytes32)" \
  "$WORKFLOW_ID" \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"
```

**Option B — identify by owner and name (both required):**

```bash
# Step 1: set the workflow owner
cast send "$RECEIVER_ADDRESS" \
  "setExpectedAuthor(address)" \
  "$WORKFLOW_OWNER_ADDRESS" \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"

# Step 2: set the workflow name (requires author to also be set)
cast send "$RECEIVER_ADDRESS" \
  "setExpectedWorkflowName(string)" \
  "$WORKFLOW_NAME" \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"
```

You may also combine both options for the strongest guarantee (e.g. set all three fields).

**Parameters:**
- `setExpectedWorkflowId(_id)`: The exact workflow ID (available from `cre workflow info`). Satisfies option A on its own.
- `setExpectedAuthor(_author)`: The account that deploys or owns the workflow. Required for option B; does not satisfy the guard alone.
- `setExpectedWorkflowName(_name)`: The exact workflow name. Required for option B; must always be paired with `setExpectedAuthor` (workflow names are unique per owner, not globally).

#### 3e. Emergency Pause (Optional but Recommended for Production)

The receiver includes an owner-only global emergency stop built on OpenZeppelin's [`Pausable`](https://docs.openzeppelin.com/contracts/5.x/api/utils#Pausable). This is the CRE equivalent of Chainlink Automation's registry-wide pause: it halts report processing across the whole receiver (there is no per-upkeep `pauseUpkeep` equivalent — pausing stops every allowlisted `(target, selector)` at once).

`pause(bool retryable)` takes a flag that lets **you decide, per pause, what happens to reports delivered while paused**:

- **`pause(true)` — retryable:** `_processReport` reverts with `EnforcedPause`. The revert propagates through `onReport`, so the CRE Forwarder records the transmission as **failed and retryable** rather than consuming it. Reports that were already signed while paused (or that arrive during the pause) are **not lost** and resume delivery once you `unpause()`.
- **`pause(false)` — drop:** `_processReport` emits `ReportSkippedWhilePaused` and returns, **consuming (dropping)** the report. Dropped reports are **not** redelivered after `unpause()`. Use this when you don't want a backlog of retried reports to flush the moment you unpause.

```bash
# Halt processing, keeping reports retryable (they resume after unpause)
cast send "$RECEIVER_ADDRESS" "pause(bool)" true \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"

# Halt processing, dropping reports delivered while paused
cast send "$RECEIVER_ADDRESS" "pause(bool)" false \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"

# Resume report processing
cast send "$RECEIVER_ADDRESS" "unpause()" \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"

# Read the current pause state and the selected mode
cast call "$RECEIVER_ADDRESS" "paused()(bool)" --rpc-url "$RPC_URL"
cast call "$RECEIVER_ADDRESS" "retryableWhilePaused()(bool)" --rpc-url "$RPC_URL"
```

> **Why on-chain pause?** Pausing or deleting a workflow stops the DON from producing *new* reports, but reports that were already signed remain permissionlessly deliverable afterward. The on-chain `pause()` closes that gap by rejecting those in-flight reports at the receiver. Use `pause(true)` to keep them retryable until you unpause, or `pause(false)` to drop them.

#### 3e. Configure Block Number Monotonicity / Staleness Protection (Optional)

Under Chainlink Automation the registry rejected any report whose trigger block preceded the last performed block, which prevented an older report from executing after a newer one. The CRE delivery path does **not** preserve this ordering: the Forwarder deduplicates only exact retransmissions of the same report and enforces no ordering across separate executions.

If your upkeep relied on that implicit ordering (typical for conditional-trigger and log-trigger upkeeps), enable the opt-in monotonicity check. The workflow **always** stamps a block number into the report payload (the triggering log's block for `LOG`, the last finalized block for `CRON`/`CUSTOM`), so you can turn this on later without any workflow change. When enabled for a `(target, selector)` pair, `_processReport` rejects any report whose block number is **older** than the last accepted one (CLA-parity: a report at the same block is still accepted). A stale report is **consumed and skipped** — the receiver emits `StaleReportSkipped` and returns without executing the call or reverting, so the forwarder does not pointlessly retry an always-stale delivery.

Like `setConsumerGasLimit`, this is configured **per `(target, selector)` pair** and is **off by default**.

```bash
# Enable the check and snapshot the current block as the floor (initialBlockNumber = 0)
cast send "$RECEIVER_ADDRESS" \
  "setBlockNumberCheck(address,bytes4,bool,uint256)" \
  "$TARGET_ADDRESS" "$(cast sig 'performUpkeep(bytes)')" true 0 \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"

# Or enable with an explicit floor (recommended): the first accepted report must be >= this block
cast send "$RECEIVER_ADDRESS" \
  "setBlockNumberCheck(address,bytes4,bool,uint256)" \
  "$TARGET_ADDRESS" "$(cast sig 'performUpkeep(bytes)')" true "$INITIAL_BLOCK" \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"

# Disable the check (also clears the stored floor)
cast send "$RECEIVER_ADDRESS" \
  "setBlockNumberCheck(address,bytes4,bool,uint256)" \
  "$TARGET_ADDRESS" "$(cast sig 'performUpkeep(bytes)')" false 0 \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"
```

**Parameters:**
- `target` / `selector`: The pair the check applies to (same values used in `setCallAllowed`).
- `enabled`: `true` to enforce monotonicity, `false` to disable and clear the stored floor.
- `initialBlockNumber`: The floor the first report must meet or exceed. `0` snapshots the current `block.number`; any non-zero value is used as-is.

> **Security note:** prefer an explicit `initialBlockNumber` over the `0` snapshot. Snapshotting only records the block at the moment you enable the check, so a report generated between deployment/configuration and that moment could still carry a higher block number and would not be rejected. Because the floor is per `(target, selector)`, set an appropriate value for each pair. Use `getBlockNumberCheck(target, selector)` to read the current `(enabled, lastReportBlock)` state.

### 4. Configure the Workflow
Update `my-workflow/config.test.json`:
- `receiverAddress`: Your deployed `AutomationReceiver`.
- `targetAddress`: Your existing Automation contract.
- `migrationType`: `CRON`, `CUSTOM`, or `LOG`.
- `schedule`: Required for `CRON` and `CUSTOM`.
- `targetFunction` and `targetInputs`: Required for `CRON`. These are ABI-encoded at runtime from the config string (e.g. `"performAction(uint256)"`), so unlike the `CUSTOM`/`LOG` paths they are not statically type-checked — a malformed signature or mismatched inputs fails when the workflow runs.
- `logTriggerAddress`, `logTriggerEventSignature`, and optional `topic1`/`topic2`/`topic3`: Required for `LOG`.

Install workflow dependencies once:

```bash
cd my-workflow
bun install
cd ..
```

### 5. Run Simulation
> The workflow throws on startup if `receiverAddress` or `targetAddress` is still the zero address, so configure Step 4 before simulating.

```bash
# For CRON (time-based) or CUSTOM (custom logic) — both run on the cron scheduler
cre workflow simulate my-workflow --target=test-settings

# For Log Trigger (requires a transaction hash containing the event)
cre workflow simulate my-workflow \
  --target=test-settings \
  --non-interactive \
  --trigger-index=0 \
  --evm-tx-hash=0x... \
  --evm-event-index=0
```

---

## Technical Details

### Log Mapping
Legacy Automation contracts expect a specific `Log` struct in `checkLog`. This template includes a utility (`mapLogToAutomation`) that maps the CRE `EVMLog` into that struct so your existing on-chain decoding logic keeps working.

> **Limitation:** the CRE log does not carry the block timestamp, so `log.timestamp` is always mapped to `0`. If your `checkLog` relies on `log.timestamp`, read it from another source (e.g. an EVM read of the block) instead.

### Read Finality
The generated `checkUpkeep` / `checkLog` bindings read at the **last finalized block** so every DON node observes the same state and reaches consensus deterministically. This differs from Automation, which simulates against the chain head. On Ethereum mainnet finality lags the head by ~13 minutes; on most L2s it is seconds. Account for this latency when migrating time-sensitive interval checks.

### Execution
The report payload is ABI-encoded as `(address target, uint256 blockNumber, bytes data)`, where `data` is a full function call (4-byte selector followed by its arguments). The `AutomationReceiver` executes `target.call(data)`, so it can drive any function signature — `performUpkeep(bytes)` for custom-logic/log upkeeps, or a custom function like `performAction(uint256)` for time-based ones. The `blockNumber` field feeds the optional monotonicity check (Step 3d) and is otherwise ignored. Every call is gated by the closed-by-default `(target, selector)` allowlist (Step 3).

The receiver distinguishes four outcomes:
- **Authorization failure** — a zero target, a target with no deployed code, calldata shorter than a 4-byte selector, or a `(target, selector)` that is not allowlisted — **reverts** (`InvalidTargetAddress` / `TargetHasNoCode` / `MissingSelector` / `CallNotAllowed`). These indicate misconfiguration or a malformed report and must surface loudly.
- **Stale report** — when the block-number monotonicity check is enabled for the pair (Step 3d) and the report's block number is older than the last accepted one — does **not** revert. The receiver emits `StaleReportSkipped(target, selector, reportBlockNumber, lastReportBlock)` and consumes the report without executing the call. Retrying would not help (the report is permanently stale), so it is skipped rather than marked retryable.
- **Gas guard failure** — when `setConsumerGasLimit` has been configured for the specific `(target, selector)` pair and the incoming gas is below `consumerGasLimit + consumerGasLimit / 63 + 7,000`, `_processReport` **reverts** with `InsufficientGas(available, required)`. The forwarder records the transmission as failed and it can be retried with higher gas. This closes a griefing attack where a report is delivered with just enough gas to pass the forwarder's minimum check but not enough for `performUpkeep` to execute, which would otherwise permanently consume the transmission ID. The 7,000-gas constant covers the EIP-2929 cold storage read, call-opcode dispatch, post-call LOG3 event emission (3 topics), and bookkeeping. The `consumerGasLimit / 63` term compensates for the EIP-150 (63/64) rule: a `CALL` can forward at most 63/64 of available gas, so without this buffer a high gas limit (above ~441,000) would cause the target to receive less than configured. Each `(target, selector)` pair has its own independent limit; pairs with no configured limit retain fire-and-forget semantics.
- **Paused** — when the owner has called `pause(bool retryable)`, `_processReport` short-circuits before any decoding or execution. With `pause(true)` it **reverts** with `EnforcedPause` (the forwarder records the transmission as failed and retryable, so reports resume after `unpause()`); with `pause(false)` it emits `ReportSkippedWhilePaused` and **consumes** the report (dropped, not redelivered). See [Emergency Pause](#3d-emergency-pause-optional).
- **Execution failure** — an allowed call that itself reverts — does **not** revert `onReport`. The receiver emits `CallFailed(target, selector, reason)` and the report is consumed. This mirrors Chainlink Automation's fire-and-forget behavior, where a failed `performUpkeep` simply ends that round and the next trigger re-evaluates eligibility.

### Gas Limit
`writeGasLimit` (default `"500000"`) caps the on-chain execution gas forwarded to `onReport`. The on-chain guard formula is `consumerGasLimit + consumerGasLimit/63 + 7,000`, but an additional ~11,000 gas is consumed before that check point by pre-guard operations — four cold `SLOAD`s in `ReceiverTemplate.onReport` (forwarder address, workflow ID, owner, name: 4 × 2,100 = 8,400 gas), four warm re-reads of those same slots in `_processReport` (the identity guard reads the inherited permission fields directly rather than via self-`STATICCALL`s: ~400 gas), `abi.decode` of the report payload (~300 gas), and the cold `SLOAD` of `s_callAllowed` (~2,100 gas). When the block-number monotonicity check is enabled for a pair (Step 3d), budget a further ~10,000 gas for its two cold `SLOAD`s and the `SSTORE` that advances the stored block. Set `writeGasLimit` to at least your `performGasLimit` plus ~20,000 (or ~30,000 with the monotonicity check enabled) to ensure deliveries comfortably clear the guard. If deliveries still fail with `InsufficientGas`, increase `writeGasLimit` in 10,000-gas increments until they pass.

---

## Compatibility

Migrating **without redeploying** your upkeep contract requires that contract to let you re-point who is authorized to call it:
- Contracts that expose a setter for the Automation Forwarder / caller (the recommended Automation pattern) can simply point it at the deployed `AutomationReceiver`.
- Contracts that **hardcode** the Automation registry/forwarder with no setter, or that have an immutable role for it, cannot be migrated in place and must be redeployed.

After migration, the `msg.sender` your target sees is the `AutomationReceiver` address (not the Automation Forwarder), so authorize that address in whatever permission check your contract uses.

## Out of Scope

- **Off-chain `offchainConfig` / gas-price-threshold controls.** In CRE these are expressed in the workflow, not on the registry.

> **Bindings note:** if you change `AutomationReceiver.sol`, regenerate the TypeScript bindings in `contracts/evm/ts/generated/` and the committed ABI in `contracts/evm/src/abi/`. The workflow itself only uses the generic `writeReport` entrypoint, so it is unaffected by the receiver's other ABI changes.

---

## Resources
- [Chainlink Automation Documentation](https://docs.chain.link/automation)
- [CRE Getting Started Guide](https://docs.chain.link/cre)
- [Legacy Migration Toolkit](https://github.com/smartcontractkit/cla-cre-migration) (Alternative manual-encoding approach)
