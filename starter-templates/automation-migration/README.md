# CRE Automation Migration Template

This template provides a modernized, type-safe path for migrating existing **Chainlink Automation** upkeeps to **Chainlink Runtime Environment (CRE)**.

## Architecture

This migration uses a "Bridge" pattern:
1. **AutomationReceiver.sol**: A generic bridge contract deployed once on-chain. It receives CRE reports and forwards the execution to your existing legacy contracts.
2. **Type-Safe Workflows**: CRE workflows that poll your legacy `checkUpkeep` or `checkLog` functions and trigger the bridge when needed.

**Result: You can usually migrate to CRE without rewriting your original upkeep logic.** If your legacy contract checks `msg.sender`, uses an Automation Forwarder allowlist, or has role-based permissions, authorize the deployed `AutomationReceiver` before relying on the workflow.

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

### 2. Deploy the Bridge
Deploy `AutomationReceiver.sol` to your target chain. You only need to do this once to support multiple upkeeps.
- Pass the address of the **CRE Forwarder** for your DON to the constructor.
- For production, configure workflow identity checks on `ReceiverTemplate` (`setExpectedWorkflowId`, `setExpectedAuthor`, or both) or narrow this bridge to known target contracts and function selectors.

### 3. Configure the Workflow
Update `my-workflow/config.test.json`:
- `receiverAddress`: Your deployed `AutomationReceiver`.
- `targetAddress`: Your existing Automation contract.
- `migrationType`: `CRON`, `CUSTOM`, or `LOG`.
- `schedule`: Required for `CRON` and `CUSTOM`.
- `targetFunction` and `targetInputs`: Required for `CRON`.
- `logTriggerAddress`, `logTriggerEventSignature`, and optional `topic1`/`topic2`/`topic3`: Required for `LOG`.

Install workflow dependencies once:

```bash
cd my-workflow
bun install
cd ..
```

### 4. Run Simulation
```bash
# For Custom Logic
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
Legacy Automation contracts expect a specific `Log` struct in `checkLog`. This template includes a utility that automatically maps the CRE `EVMLog` to this legacy format, ensuring compatibility with your existing on-chain decoding logic.

### Generic Execution
The `AutomationReceiver` uses `target.call(data)` to execute actions. This means it can handle any function signature, not just `performUpkeep`. This is particularly useful for migrating time-based upkeeps that call custom functions like `performAction(uint256)`.

Because the bridge is generic, production deployments should pair it with explicit authorization: keep the CRE forwarder check enabled, configure expected workflow identity fields, and only grant the receiver permissions that the migrated upkeep needs.

---

## Resources
- [Chainlink Automation Documentation](https://docs.chain.link/automation)
- [CRE Getting Started Guide](https://docs.chain.link/cre)
- [Legacy Migration Toolkit](https://github.com/smartcontractkit/cla-cre-migration) (Alternative manual-encoding approach)
