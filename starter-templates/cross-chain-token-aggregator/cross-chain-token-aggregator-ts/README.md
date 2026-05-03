# Cross-Chain Token Aggregator — CRE Starter Template (TypeScript)

Let users aggregate their tokens received on scattered chains to a single chain by utilizing the Transfer event of ERC-20 compatible token, powered by event-driven CRE workflows.

**⚠️ DISCLAIMER**

This template is an educational example to demonstrate how to interact with Chainlink systems, products, and services. It is provided **"AS IS"** and **"AS AVAILABLE"** without warranties of any kind, has **not** been audited, and may omit checks or error handling for clarity. **Do not use this code in production** without performing your own audits and applying best practices.

---

## Overview

This template demonstrates the **detect → notify → bridge** pattern using Chainlink CRE (Compute Runtime Environment). It is designed around a **user-centric aggregation flow**: a normal user can send ERC-20 tokens from any wallet to a single configured address, and the CRE workflow automatically detects the inbound transfer, notifies them via Telegram, and bridges the tokens to their wallet on a destination chain — without any manual intervention.

### Use Cases

- **Personal cross-chain aggregator**: User receives token scattered on different chain, which gets accumulated to a single chain.
- **Bridge abstraction layer**: Abstract away the choice of bridge (Across vs CCIP) behind a per-token config — users never interact with the bridge directly
- **Multi-source token consolidation**: Tokens arriving from multiple senders/contracts all funnel to one destination address on Unichain
- **On-chain event alerting + action**: Combine real-time Telegram notifications with automated on-chain bridging for any ERC-20 transfer

## Architecture

```
  User receives tokens to their address (any source)
                          │
                          │  Transfer(from, to=userAddress, amount)
                          v
┌─────────────────────────────────────────────────────────────────────┐
│                            CRE DON                                  │
│                                                                     │
│  ┌──────────────┐   ┌────────────────┐   ┌──────────────────────┐   │
│  │  LogTrigger  │──>│ Decode Transfer│──>│  Read ERC-20         │   │
│  │  (Transfer   │   │ Event (viem)   │   │  Decimals, Allwances │   │
│  │   to user)   │   └────────────────┘   └──────────┬───────────┘   │
│  └──────────────┘                                   │               │
│                                                     v               │
│                                        ┌────────────────────────┐   │
│                                        │  Telegram Notification │   │
│                                        │  "Token received from  │   │
│                                        │   0x... amount X"      │   │
│                                        └────────────┬───────────┘   │
│                                                     │               │
│                                        ┌────────────v───────────┐   │
│                                        │  BridgeFactory         │   │
│                                        │  decides token bridge  │   │
│                                        │  → "across" | "ccip"   │   │
│                                        └────────────┬───────────┘   │
│                                                     │               │
│              ┌──────────────────────────────────────┤               │
│              │                                      │               │
│   ┌──────────v──────────┐            ┌──────────────v─────────┐     │
│   │  Across Bridge      │            │  Chainlink CCIP Bridge │     │
│   │  GET /swap/approval │            │  Encode CCIP params    │     │
│   │  Encode + sign      │            │  Sign + writeReport    │     │
│   │  writeReport(0x01…) │            │  writeReport(0x02…)    │     │
│   └──────────┬──────────┘            └──────────────┬─────────┘     │
└──────────────┼──────────────────────────────────────┼───────────────┘
               │                                      │
               └─────────────────┬────────────────────┘
                                 │
                      ┌──────────v──────────┐
                      │  KeystoneForwarder  │
                      │  → Uniflow.onReport │
                      └──────────┬──────────┘
                                 │
                      ┌──────────v──────────┐
                      │  Uniflow.sol        │
                      │  0x01 → Across dep. │
                      │  0x02 → CCIP send   │
                      └──────────┬──────────┘
                                 │
                      ┌──────────v──────────┐
                      │  Tokens arrive in   │
                      │  user's wallet on   │
                      │  base chain         │
                      └─────────────────────┘
```

## Components

### CRE Workflow (`my-workflow/`)

The TypeScript workflow runs off-chain inside the CRE DON:

1. **LogTrigger** fires when any configured ERC-20 token emits a `Transfer(from, to, amount)` event where `to` matches the configured `targetUserAddress`, at `CONFIDENCE_LEVEL_FINALIZED`
2. **Decodes** the transfer event using `viem`'s `decodeEventLog`
3. **Reads** the token's `decimals()` on-chain via `evmClient.callContract`
4. **Sends a Telegram message** notifying the user of the inbound transfer (sender, amount, token contract)
5. **Looks up** the token in `config.tokenMap` to determine which bridge to use
6. **Routes to bridge**:
   - **Across**: Calls the Across REST API (`/api/swap/approval`) to get approval and deposit calldata, encodes bridge parameters, signs a CRE report prefixed `0x01`, and writes it on-chain
   - **Chainlink CCIP**: Encodes receiver, token, amount, and destination chain selector, signs a CRE report prefixed `0x02`, and writes it on-chain
7. Returns `{ success, txHash, error }` — bridging is complete once the on-chain transaction confirms

### Smart Contracts (`contracts/`)

**`Uniflow.sol`** — The on-chain receiver and bridge dispatcher:

- Extends `ReceiverTemplate` — validates the CRE Forwarder and optional workflow identity before processing any report
- `setupToken(token, tokenConfig)` — owner configures which tokens are supported with a minimum bridging amount and receiver address
- `allowlistDestinationChainForCCIP(selector, enable)` — enable/disable destination chains for CCIP bridging
- `_processReport(bytes)` — dispatches on the first-byte opcode: `0x01` → Across, `0x02` → CCIP
- `_performAcrossBridgeOp(bytes)` — pulls tokens from owner via `safeTransferFrom`, approves the Across deposit contract, and calls `depositContract.call(depositData)` to initiate the bridge
- `_performChainlinkCCIPBridgeOp(bytes)` — pulls tokens from owner, pays CCIP fees in LINK, and calls `ccipSend` on the CCIP Router to send tokens to the destination chain

**`ReceiverTemplate.sol`** — Abstract base with layered security:

- Validates that `msg.sender == s_forwarderAddress` (the Chainlink KeystoneForwarder)
- Optionally validates workflow ID, workflow owner address, and workflow name
- `setForwarderAddress`, `setExpectedAuthor`, `setExpectedWorkflowName`, `setExpectedWorkflowId` — all owner-configurable post-deployment

### Bridge Integrations (`my-workflow/bridge/`)

New bridges can be added by implementing the `IBridge` interface and registering a key in `BridgeFactory`.

| Bridge | Config Key | Mechanism |
|--------|-----------|-----------|
| Across Protocol | `"across"` | GET `/api/swap/approval` → encode approval + deposit calldata → report prefix `0x01` |
| Chainlink CCIP | `"chainlink_ccip"` | Encode `(receiver, token, amount, chainSelector)` → report prefix `0x02` |

## Getting Started

### Prerequisites

- [Bun](https://bun.sh/) runtime installed
- [CRE CLI](https://docs.chain.link/cre) installed
- [Foundry](https://getfoundry.sh/) installed (for contract deployment)
- A Telegram bot token and chat ID (for notifications - optional)

### 1. Install Dependencies

```bash
cd my-workflow && bun install && cd ..
```

### 2. Set Up Secrets

> **How to get Telegram credentials**: Message `@BotFather` on Telegram, create a bot with `/newbot`, and copy the token. To find your chat ID, message the bot once and then call `https://api.telegram.org/bot<TOKEN>/getUpdates` — your chat ID is in the `message.chat.id` field.

The `secrets.yaml` at the project root maps secret names to the workflow's environment variable names:

```yaml
secretsNames:
  TELEGRAM_BOT_ACCESS_TOKEN:
    - TELEGRAM_BOT_ACCESS_TOKEN_VAR
  TELEGRAM_CHAT_ID:
    - TELEGRAM_CHAT_ID_VAR
```

### 3. Review and Customize Config

`my-workflow/config.staging.json` controls which user address, tokens, and destination chain the workflow uses:

```json
{
  "networks": {
    "eth": {
      "chainId": "11155111",
      "creNetworkConfig": [{ "chainFamily": "evm", "chainSelectorName": "ethereum-testnet-sepolia", "isTestnet": true }],
      "configContract": "0x4d236Ec82Af6c73835F29BBdc4b34574a0E0FdaE",
      "targetUserAddress": "0xF1c8170181364DeD1C56c4361DED2eB47f2eef1b",
      "tokenArr": ["0x1c7d4b196cb0c7b01d743fbc6116a902379c7238"],
      "tokenMap": {
        "0x1c7d4b196cb0c7b01d743fbc6116a902379c7238": {
          "unichainToken": "0x036CbD53842c5426634e7929541eC2318f3dCF7e",
          "bridge": "chainlink_ccip"
        }
      }
    }
  },
  "unichain": {
    "chainId": "84532",
    "unichainDestinationAddress": "0xF1c8170181364DeD1C56c4361DED2eB47f2eef1b",
    "chainlinkCCIPSelector": "10344971235874465080"
  },
  "telegramMessageApi": "https://api.telegram.org/bot{{TELEGRAM_BOT_ACCESS_TOKEN}}/sendMessage?chat_id={{TELEGRAM_CHAT_ID}}&text={{MESSAGE}}",
  "acrossApiUrl": "https://testnet.across.to"
}
```

Key fields to update for your own setup:

| Field | Description |
|-------|-------------|
| `targetUserAddress` | The user's Ethereum sepolia address whose inbound transfers are monitored |
| `tokenArr` | ERC-20 token contract addresses to watch on Ethereum |
| `tokenMap[token].bridge` | `"chainlink_ccip"` or `"across"` per token |
| `tokenMap[token].unichainToken` | Corresponding token address on the destination chain |
| `unichain.unichainDestinationAddress` | The user's wallet address on base sepolia chain (where tokens gets routed) |
| `configContract` | Address of your deployed `Uniflow.sol` |

### 4. Simulate

LogTrigger simulation requires a real transaction hash that emitted a `Transfer` to your `targetUserAddress`. Send a small ERC-20 transfer on Sepolia to that address, then simulate:

```bash
cre workflow simulate my-workflow --broadcast
```
Enter the hash of the transaction containing the Transfer event, and the workflow will do its work.

### 5. Deploy Contracts

Deploy `Uniflow.sol` using Foundry:

```bash
forge script script/DeployUniflow.s.sol --rpc-url <rpc-url> --account <account> --broadcast
```
Now on the deployed contract:
1. Set up the token config.
2. Allowlist the destination chain selector for chainlink ccip bridging.
3. Set token approval for the desired token to the Uniflow contract, so that it can spend on your behalf.
4. Send some link to the contract to facilitate ccip fees.

> **Sepolia contract addresses** (verify on [Chainlink docs](https://docs.chain.link/ccip/directory/testnet)):
> - CRE KeystoneForwarder: `0x15fC6ae953E024d975e77382eEeC56A9101f9F88`
> - CCIP Router: `0x0BF3dE8c5D3e8A2B34D2BEeB17ABfCeBaf363A59`
> - LINK Token: `0x779877A7B0D9E8603169DdbD7836e478b4624789`

**1. Register a supported token:**
```bash
cast send <UNIFLOW_ADDRESS> \
  "setupToken(address,(address,uint256))" \
  <TOKEN_ADDRESS> "(<RECEIVER_ON_UNICHAIN>,1000000)" \
  --private-key <KEY> --rpc-url https://ethereum-sepolia-rpc.publicnode.com
```

**2. Allowlist the Unichain destination chain for CCIP:**
```bash
cast send <UNIFLOW_ADDRESS> \
  "allowlistDestanationChainForCCIP(uint64,bool)" \
  10344971235874465080 true \
  --private-key <KEY> --rpc-url https://ethereum-sepolia-rpc.publicnode.com
```

**3. Approve `Uniflow.sol` to pull the user's tokens** (required before any bridge operation executes):
```bash
cast send <TOKEN_ADDRESS> \
  "approve(address,uint256)" <UNIFLOW_ADDRESS> 115792089237316195423570985008687907853269984665640564039457584007913129639935 \
  --private-key <KEY> --rpc-url https://ethereum-sepolia-rpc.publicnode.com
```

**4. Fund the contract with LINK to cover CCIP fees:**
```bash
cast send <LINK_TOKEN> \
  "transfer(address,uint256)" <UNIFLOW_ADDRESS> 5000000000000000000 \
  --private-key <KEY> --rpc-url https://ethereum-sepolia-rpc.publicnode.com
```

**5. Update `configContract`** in your config JSON to the deployed `Uniflow.sol` address.

## Customization

- **Add a new bridge**: Implement `IBridge` in `my-workflow/bridge/integrations/` and register the key in `BridgeFactory.getBridge()`
- **Monitor tokens on multiple source chains**: Add more entries to the `networks` map in config and create additional `evmClient` + `logTrigger` instances in `main.ts`
- **Change the destination chain**: Update `unichain.chainId`, `chainlinkCCIPSelector`, and each `tokenMap[token].unichainToken` to your target chain's values
- **Per-token bridge selection**: Already supported — set `"bridge": "across"` or `"bridge": "chainlink_ccip"` individually per token in `tokenMap`
- **Richer notifications**: Extend `telegramMessageService.ts` to include bridge type, estimated arrival time, or destination tx hash in the Telegram message
- **Use Confidential HTTP**: Replace `new cre.capabilities.HTTPClient()` with the Confidential HTTP capability in `across.ts` to keep bridge quote data private from node operators

## Security

- The contracts are **demos** — audit and customize before production use
- `ReceiverTemplate` ensures only the CRE Forwarder can call `onReport()`. After deployment, call `setExpectedAuthor()` and `setExpectedWorkflowId()` to lock reports to your specific workflow and prevent spoofed reports
- `Uniflow.sol` pulls tokens from the `owner` via `safeTransferFrom` — the owner must `approve` the contract for at least the bridging amount before any transfer event is processed
- The CCIP bridge explicitly checks the destination chain selector against an allowlist before sending — no unexpected chains can be targeted
- Never commit `secrets.yaml` with real values or any `.env` files to version control
