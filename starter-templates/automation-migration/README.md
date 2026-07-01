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

The `my-workflow/contracts/evm` directory ships a ready-to-use [Foundry](https://book.getfoundry.sh/) project: a `foundry.toml` and a vendored copy of the only OpenZeppelin files used (`Ownable`, `Context`). No `forge install` or local node is required — just build and deploy:

```bash
cd my-workflow/contracts/evm
forge build
forge test          # 24 tests covering the receiver + permission template

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
- `target`: The address of your existing upkeep contract. Must be a deployed contract — `setCallAllowed` reverts with `TargetHasNoCode` if the address has no code (EOA, mistyped address, or never-deployed contract).
- `selector`: The 4-byte function selector (computed from the function signature via `cast sig`).
- `allowed`: Set to `true` to allow, `false` to revoke.

The selector must match the function the workflow encodes (`performUpkeep` for `CUSTOM`/`LOG`, or your `targetFunction` for `CRON`). A mismatch makes `onReport` revert with `CallNotAllowed`.

#### 3b. Configure the Consumer Gas Limit (Recommended)

Set the minimum gas the receiver must have available before forwarding the call to your upkeep consumer contract. When configured, `_processReport` reverts with `InsufficientGas` — causing the CRE Forwarder to record the delivery as **failed and retryable** — if the incoming gas is below `consumerGasLimit + 6,000` (the on-chain overhead).

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
The `AutomationReceiver` executes `target.call(data)`, so it can drive any function signature — `performUpkeep(bytes)` for custom-logic/log upkeeps, or a custom function like `performAction(uint256)` for time-based ones. Every call is gated by the closed-by-default `(target, selector)` allowlist (Step 3).

The receiver distinguishes three failure modes:
- **Authorization failure** — a zero target, a target with no deployed code, calldata shorter than a 4-byte selector, or a `(target, selector)` that is not allowlisted — **reverts** (`InvalidTargetAddress` / `TargetHasNoCode` / `MissingSelector` / `CallNotAllowed`). These indicate misconfiguration or a malformed report and must surface loudly.
- **Gas guard failure** — when `setConsumerGasLimit` has been configured for the specific `(target, selector)` pair and the incoming gas is below `consumerGasLimit + 6,000`, `_processReport` **reverts** with `InsufficientGas(available, required)`. The forwarder records the transmission as failed and it can be retried with higher gas. This closes a griefing attack where a report is delivered with just enough gas to pass the forwarder's minimum check but not enough for `performUpkeep` to execute, which would otherwise permanently consume the transmission ID. The 6,000-gas overhead covers the EIP-2929 cold storage read, call-opcode dispatch, post-call event emission, and bookkeeping. Each `(target, selector)` pair has its own independent limit; pairs with no configured limit retain fire-and-forget semantics.
- **Execution failure** — an allowed call that itself reverts — does **not** revert `onReport`. The receiver emits `CallFailed(target, selector, reason)` and the report is consumed. This mirrors Chainlink Automation's fire-and-forget behavior, where a failed `performUpkeep` simply ends that round and the next trigger re-evaluates eligibility.

### Gas Limit
`writeGasLimit` (default `"500000"`) caps the on-chain execution. Carry over the `performGasLimit` you tuned for your existing upkeep rather than relying on the default.

---

## Compatibility

Migrating **without redeploying** your upkeep contract requires that contract to let you re-point who is authorized to call it:
- Contracts that expose a setter for the Automation Forwarder / caller (the recommended Automation pattern) can simply point it at the deployed `AutomationReceiver`.
- Contracts that **hardcode** the Automation registry/forwarder with no setter, or that have an immutable role for it, cannot be migrated in place and must be redeployed.

After migration, the `msg.sender` your target sees is the `AutomationReceiver` address (not the Automation Forwarder), so authorize that address in whatever permission check your contract uses.

## Out of Scope

- **Off-chain `offchainConfig` / gas-price-threshold controls.** In CRE these are expressed in the workflow, not on the registry.

> **Bindings note:** if you change `AutomationReceiver.sol`, regenerate the TypeScript bindings in `my-workflow/contracts/evm/ts/generated/` and the committed ABI in `my-workflow/contracts/evm/src/abi/`. The workflow itself only uses the generic `writeReport` entrypoint, so it is unaffected by the receiver's other ABI changes.

---

## Resources
- [Chainlink Automation Documentation](https://docs.chain.link/automation)
- [CRE Getting Started Guide](https://docs.chain.link/cre)
- [Legacy Migration Toolkit](https://github.com/smartcontractkit/cla-cre-migration) (Alternative manual-encoding approach)
