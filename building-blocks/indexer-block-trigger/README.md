# CRE Indexer Block Trigger Workflows

Workflows for processing new blocks and transactions using block-triggered webhooks (from Alchemy Notify) and matching 
against watched addresses. These workflows demonstrate the **block trigger pattern** where the workflow reacts to 
incoming block data and extracts relevant transactions.

## Directory Structure

```
building-blocks/indexer-block-trigger/
├── block-trigger-go/     (Go-based workflow)
│   └── workflow/
│       ├── main.go
│       ├── config.staging.json
│       ├── config.production.json
│       ├── workflow.yaml
│       └── README.md
├── block-trigger-ts/     (TypeScript-based workflow)
    └── workflow/
        ├── main.ts
        ├── package.json
        ├── config.staging.json (optional)
        ├── workflow.yaml
        └── README.md

```

## Overview

These workflows demonstrate how to:
- React to block events via HTTP webhook triggers 
(We use Alchemy Notify for this workflow)
- Match transactions to a list of watched addresses
- Process and return JSON-formatted block and transaction data
- Implement the same logic in both Go and TypeScript

Both workflows process incoming block data and extract:
- Block number, hash, timestamp
- All transactions in the block
- Transactions where the `to` address matches a watched address

## Workflows

### 1. block-trigger-go (Go Implementation)

**Language:** Go

**Features:**
- Uses `http.Trigger` from CRE Go SDK
- Matches transactions to watched addresses from config
- Returns formatted JSON summary of block and matched transactions

**Running the workflow:**
```bash
cd building-blocks/indexer-block-trigger/block-trigger-go
cre workflow simulate workflow --non-interactive --trigger-index 0 --http-payload test-block.json --target staging-settings
```

### 2. block-trigger-ts (TypeScript Implementation)

**Language:** TypeScript

**Features:**
- Uses HTTP trigger from CRE TypeScript SDK
- Matches transactions to watched addresses from config
- Returns formatted JSON summary of block and matched transactions

**Running the workflow:**
```bash
cd building-blocks/indexer-block-trigger/block-trigger-ts/workflow
bun install
cre workflow simulate workflow --non-interactive --trigger-index 0 --http-payload test-block.json --target staging-settings
```

## Setup and Testing

### Prerequisites

**For Go workflow:**
1. Install CRE CLI
2. Login: `cre login`
3. Install Go

**For TypeScript workflow:**
1. Install CRE CLI
2. Login: `cre login`
3. Install Bun (or Node.js)
4. Run `bun install` in the workflow directory

### Running the Workflows

**Go Workflow:**
```bash
cd building-blocks/indexer-block-trigger/block-trigger-go
cre workflow simulate workflow --non-interactive --trigger-index 0 --http-payload test-block.json --target staging-settings
```

**TypeScript Workflow:**
```bash
cd building-blocks/indexer-block-trigger/block-trigger-ts
cre workflow simulate workflow --non-interactive --trigger-index 0 --http-payload test-block.json --target staging-settings
```

### Example Output

Both workflows return JSON output like:

```json
{
  "blockNumber": 12345678,
  "blockHash": "0xabc...",
  "timestamp": 1700000000,
  "totalLogs": 42,
  "uniqueTransactions": 10,
  "matchedTransactions": 2,
  "transactions": [
    {
      "hash": "0xdef...",
      "from": "0x...",
      "to": "0x73b668d8374ddb42c9e2f46fd5b754ac215495bc",
      "value": "1000000000000000000",
      ...
    }
  ]
}
```

## Example Use Cases

### 1. Monitoring High-Value Addresses
Track transactions to specific addresses in real time:
```json
{
  "watchedAddresses": ["0x...", "0x..."]
}
```

### 2. Contract Interaction Tracking
Detect when contracts of interest receive transactions:
```json
{
  "watchedAddresses": ["0xContract1", "0xContract2"]
}
```

### 3. Block-Level Analytics
Summarize block activity and matched transactions for analytics dashboards.

## Reference Documentation

- [CRE Documentation](https://docs.chain.link/cre)
- [Alchemy Webhooks](https://www.alchemy.com/docs/reference/custom-webhook)

