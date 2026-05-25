# CRE Automation Migration Template

This template provides a modernized, type-safe path for migrating existing **Chainlink Automation** upkeeps to **Chainlink Runtime Environment (CRE)**.

## Architecture

This migration uses a "Bridge" pattern:
1. **AutomationReceiver.sol**: A generic bridge contract deployed once on-chain. It receives CRE reports and forwards the execution to your existing legacy contracts.
2. **Type-Safe Workflows**: CRE workflows that poll your legacy `checkUpkeep` or `checkLog` functions and trigger the bridge when needed.

**Result: You migrate to CRE without changing a single line of your original Solidity code.**

---

## Migration Path Mapping

| Your Automation Setup | CRE Equivalent | Migration Type |
|---|---|---|
| **Time-based upkeep** | Cron Trigger | `CRON` |
| **Custom logic upkeep** | Cron + EVM Read | `CUSTOM` |
| **Log trigger upkeep** | EVM Log Trigger | `LOG` |

---

## Getting Started

### 1. Deploy the Bridge
Deploy `AutomationReceiver.sol` to your target chain. You only need to do this once to support multiple upkeeps.
- Pass the address of the **CRE Forwarder** for your DON to the constructor.

### 2. Configure the Workflow
Update `my-workflow/config.test.json`:
- `receiverAddress`: Your deployed `AutomationReceiver`.
- `targetAddress`: Your existing Automation contract.
- `migrationType`: `CRON`, `CUSTOM`, or `LOG`.

### 3. Run Simulation
```bash
# For Custom Logic
cre workflow simulate my-workflow --target=test-settings

# For Log Trigger (requires a transaction hash containing the event)
EVM_TX_HASH=0x... cre workflow simulate my-workflow --target=test-settings
```

---

## Technical Details

### Log Mapping
Legacy Automation contracts expect a specific `Log` struct in `checkLog`. This template includes a utility that automatically maps the CRE `EVMLog` to this legacy format, ensuring compatibility with your existing on-chain decoding logic.

### Generic Execution
The `AutomationReceiver` uses `target.call(data)` to execute actions. This means it can handle ANY function signature, not just `performUpkeep`. This is particularly useful for migrating time-based upkeeps that call custom functions like `performAction(uint256)`.

---

## Resources
- [Chainlink Automation Documentation](https://docs.chain.link/automation)
- [CRE Getting Started Guide](https://docs.chain.link/cre)
- [Legacy Migration Toolkit](https://github.com/smartcontractkit/cla-cre-migration) (Alternative manual-encoding approach)
