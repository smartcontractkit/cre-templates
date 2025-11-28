# CRE Indexer Block Trigger Workflow (Go)

This workflow processes new blocks and transactions using block-triggered webhooks (e.g., Alchemy Notify) and matches against watched addresses. It demonstrates the **block trigger pattern** in Go.

## Features
- Uses `http.Trigger` from CRE Go SDK
- Matches transactions to watched addresses from config
- Returns formatted JSON summary of block and matched transactions

## Setup and Prerequisites
1. Install CRE CLI
2. Login: `cre login`
3. Install Go

## Running the Workflow
```bash
cd building-blocks/indexer-block-trigger/block-trigger-go
cre workflow simulate workflow --non-interactive --trigger-index 0 --http-payload test-block.json --target staging-settings
```

## Example Output
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
      "value": "1000000000000000000"
    }
  ]
}
```

## Example Use Cases
- Monitoring high-value addresses
- Contract interaction tracking
- Block-level analytics

## Reference Documentation
- [CRE Documentation](https://docs.chain.link/cre)
- [Alchemy Webhooks](https://www.alchemy.com/docs/reference/custom-webhook)
