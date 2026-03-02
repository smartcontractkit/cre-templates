<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE\&message=Home\&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE\&message=Docs\&color=blue)](https://docs.chain.link/cre)

</div>

## Webhook Notification (TypeScript)

Read a Chainlink Data Feed on-chain and send a price notification to **Slack** or **Telegram** via webhook.

This building block combines two CRE capabilities:
- **EVM Client** to read on-chain data (price feed)
- **HTTP Client** to POST a notification to an external webhook

## Quick start

### 1) Configure RPC in `project.yaml`

Add an RPC for the chain you want to read from. For Arbitrum One mainnet:

```yaml
rpcs:
  - chain-name: ethereum-mainnet-arbitrum-1
    url: <YOUR_ARBITRUM_MAINNET_RPC_URL>
```

### 2) Configure the workflow

Update `my-workflow/config.production.json` with your webhook settings:

```json
{
  "schedule": "0 */10 * * * *",
  "chainName": "ethereum-mainnet-arbitrum-1",
  "feed": {
    "name": "ETH/USD",
    "address": "0x639Fe6ab55C921f74e7fac1ee960C0B6293ba612"
  },
  "webhookUrl": "https://hooks.slack.com/services/YOUR/WEBHOOK/URL",
  "notificationType": "slack",
  "telegramChatId": ""
}
```

**Configuration fields:**

| Field | Description |
|---|---|
| `schedule` | 6-field cron expression (e.g., every 10 minutes at second 0) |
| `chainName` | Must match the RPC entry in `project.yaml` |
| `feed.name` | Human-readable feed name (e.g., `"ETH/USD"`) |
| `feed.address` | Chainlink Data Feed proxy address on the target chain |
| `webhookUrl` | Slack incoming webhook URL or Telegram bot API URL |
| `notificationType` | `"slack"` or `"telegram"` |
| `telegramChatId` | Required when `notificationType` is `"telegram"` |

**Telegram example:**

```json
{
  "webhookUrl": "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/sendMessage",
  "notificationType": "telegram",
  "telegramChatId": "123456789"
}
```

### 3) Install dependencies

From your project root:

```bash
bun install --cwd ./my-workflow
```

### 4) Run a local simulation

The staging config uses `https://httpbin.org/post` as a test echo endpoint:

```bash
cre workflow simulate my-workflow
```

You should see output similar to:

```
Workflow compiled
[SIMULATION] Simulator Initialized

[SIMULATION] Running trigger trigger=cron-trigger@1.0.0
[USER LOG] msg="Price feed read" chain=ethereum-mainnet-arbitrum-1 feed="ETH/USD" address=0x639Fe6ab55C921f74e7fac1ee960C0B6293ba612 decimals=8 latestAnswerRaw=378968000000 latestAnswerScaled=3789.68
[USER LOG] msg="Formatted price" feed="ETH/USD" price=$3,789.68
[USER LOG] msg="Sending slack notification" url=https://httpbin.org/post
[USER LOG] msg="Webhook response" statusCode=200

Workflow Simulation Result:
 "{\"feed\":\"ETH/USD\",\"address\":\"0x639Fe6ab55C921f74e7fac1ee960C0B6293ba612\",\"decimals\":8,...}"
```
