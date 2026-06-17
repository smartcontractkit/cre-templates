<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE\&message=Home\&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE\&message=Docs\&color=blue)](https://docs.chain.link/cre)

</div>

## Webhook Alerting (TypeScript)

Read a Chainlink Data Feed on-chain and trigger a **PagerDuty alert** via the ConfidentialHTTPClient.

This building block demonstrates **VaultDON secret injection** — the PagerDuty routing key is referenced as `{{.pagerdutyRoutingKey}}` in the JSON request body and resolved at runtime by the secure enclave (from env vars during simulation, from VaultDON in production).

> **Looking for Slack/Telegram notifications?** See the sibling [`webhook-notification`](../webhook-notification/) building block, which demonstrates enclave privacy for URL-embedded credentials.

### Capabilities used

- **EVM Client** — read on-chain price feed data
- **Confidential HTTP Client** — POST the alert with secret injection (`vaultDonSecrets`)
- **Cron Scheduler** — fire on a configurable schedule

## Quick start

### 1) Configure RPC in `project.yaml`

Add an RPC for the chain you want to read from. For Arbitrum One mainnet:

```yaml
rpcs:
  - chain-name: ethereum-mainnet-arbitrum-1
    url: <YOUR_ARBITRUM_MAINNET_RPC_URL>
```

### 2) Set your PagerDuty routing key

Add your routing key to `.env`:

```
PAGERDUTY_ROUTING_KEY=YOUR_PAGERDUTY_ROUTING_KEY_HERE
```

The `secrets.yaml` file maps the vault secret name `pagerdutyRoutingKey` to the `PAGERDUTY_ROUTING_KEY` env var. During simulation the CLI reads the value from `.env`; in production the enclave fetches it from VaultDON.

### 3) Configure the workflow

Update `my-workflow/config.production.json` with your settings:

```json
{
  "schedule": "0 */10 * * * *",
  "chainName": "ethereum-mainnet-arbitrum-1",
  "feed": {
    "name": "ETH/USD",
    "address": "0x639Fe6ab55C921f74e7fac1ee960C0B6293ba612"
  },
  "endpoint": "https://events.pagerduty.com/v2/enqueue",
  "severity": "critical",
  "source": "cre-workflow"
}
```

**Configuration fields:**

| Field | Description |
|---|---|
| `schedule` | 6-field cron expression (e.g., every 10 minutes at second 0) |
| `chainName` | Must match the RPC entry in `project.yaml` |
| `feed.name` | Human-readable feed name (e.g., `"ETH/USD"`) |
| `feed.address` | Chainlink Data Feed proxy address on the target chain |
| `endpoint` | PagerDuty Events API v2 URL (or `https://httpbin.org/post` for testing) |
| `severity` | PagerDuty severity: `"critical"`, `"error"`, `"warning"`, or `"info"` |
| `source` | Source identifier included in the alert payload |

### 4) Install dependencies

From your project root:

```bash
bun install --cwd ./my-workflow
```

### 5) Run a local simulation

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
[USER LOG] msg="Sending PagerDuty alert"
[USER LOG] msg="Alert response" statusCode=200

Workflow Simulation Result:
 "{\"feed\":\"ETH/USD\",\"address\":\"0x639Fe6ab55C921f74e7fac1ee960C0B6293ba612\",\"decimals\":8,...}"
```

### How secret injection works

1. `secrets.yaml` declares `pagerdutyRoutingKey` mapped to the `PAGERDUTY_ROUTING_KEY` env var.
2. The workflow builds a JSON body containing the literal string `{{.pagerdutyRoutingKey}}`.
3. The `ConfidentialHTTPClient.sendRequest()` call includes `vaultDonSecrets: [{ key: 'pagerdutyRoutingKey' }]`.
4. Before sending the request, the enclave resolves `{{.pagerdutyRoutingKey}}` in the body with the actual secret value.
5. The secret never appears in logs or leaves the enclave boundary.
