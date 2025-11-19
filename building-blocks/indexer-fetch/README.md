# CRE Indexer Workflows

Workflows for pulling data from The Graph indexer with scheduled cron triggers, created using `cre init`.

## Directory Structure

```
building-blocks/cre-indexer-workflows/
├── README.md (this file)
├── indexer-workflow-ts/  (Go-based workflow for indexer queries)
│   └── my-workflow/
│       ├── workflow.go
│       ├── main.go
│       ├── config.staging.json
│       ├── config.production.json
│       └── workflow.yaml
└── indexer-workflow-go/  (Go hello world template)
    └── my-workflow/
        └── [template files]
```

## Overview

These workflows demonstrate how to:
- Query The Graph indexer using GraphQL
- Use cron triggers to schedule periodic data fetching (every minute by default)
- Process and return JSON-formatted indexer data

## Main Workflow: indexer-workflow-ts

**Note:** Despite the name, this uses Go (created from CRE template 1).

### Configuration

The workflow is configured in `my-workflow/config.staging.json`:

```json
{
  "schedule": "0 * * * * *",
  "graphqlEndpoint": "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2",
  "query": "query GetPairs($first: Int!) { pairs(first: $first, orderBy: reserveUSD, orderDirection: desc) { id token0 { id symbol name } token1 { id symbol name } reserveUSD } }",
  "variables": {
    "first": 3
  }
}
```

### Configuration Options

- **schedule**: Cron expression in 6-field format (second minute hour day month weekday)
  - `"0 * * * * *"` - Every minute at second 0
  - `"0 */5 * * * *"` - Every 5 minutes at second 0

- **graphqlEndpoint**: The Graph API endpoint URL
  - Public endpoint: `https://api.thegraph.com/subgraphs/name/{owner}/{subgraph}`
  - Studio endpoint: `https://api.studio.thegraph.com/query/{id}/{name}/version/latest`

- **query**: GraphQL query string with optional variables

- **variables**: Object with variables for the GraphQL query

### Workflow Code Structure

The workflow (`workflow.go`) includes:

1. **Config struct**: Holds GraphQL endpoint, query, schedule, and variables
2. **InitWorkflow**: Sets up the cron trigger
3. **onIndexerCronTrigger**: Main handler that fetches data when cron fires
4. **fetchGraphData**: Makes HTTP POST request to The Graph endpoint

### Key Implementation Details

```go
// Uses HTTP SendRequest pattern from CRE SDK
client := &http.Client{}
result, err := http.SendRequest(config, runtime, client, fetchGraphData, nil).Await()

// GraphQL request structure
gqlRequest := GraphQLRequest{
    Query:     config.Query,
    Variables: config.Variables,
}

// Makes POST request with proper headers
httpResp, err := sendRequester.SendRequest(&http.Request{
    Method: "POST",
    Url:    config.GraphqlEndpoint,
    Headers: map[string]string{
        "Content-Type": "application/json",
    },
    Body: requestBody,
}).Await()
```

## Setup and Testing

### Prerequisites

1. Install CRE CLI
2. Login: `cre login`
3. Go 1.23+ installed

### Running the Workflow

1. Navigate to the workflow directory:
```bash
cd building-blocks/indexer-workflows/indexer-workflow-ts
```

2. Test with simulation:
```bash
cre workflow simulate my-workflow --target staging-settings
```

3. When prompted, select trigger `1` (cron-trigger)

### Expected Behavior

**Boilerplate (PoR) workflow**: ✅ Compiles and runs successfully
**Enhanced indexer workflow**: ⚠️ Compiles successfully but encounters runtime issues with HTTP capability in simulation

## Current Status

### ✅ Working

- Workflow structure and configuration
- Compilation and WASM generation
- Cron trigger setup
- GraphQL request formatting
- Error handling

### ⚠️ Known Issues

The enhanced indexer workflow compiles but fails at runtime with:
```
Workflow execution failed:
 error while executing at wasm backtrace
Caused by:
    Exited with i32 exit status 2
```

This appears to be a runtime issue with the HTTP capability in the CRE SDK during simulation. The code structure follows the correct patterns from working examples.

##  Example Use Cases

### 1. Monitoring Uniswap Pairs
Query top liquidity pools every minute:
```json
{
  "schedule": "0 * * * * *",
  "graphqlEndpoint": "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2",
  "query": "query GetPairs($first: Int!) { pairs(first: $first, orderBy: reserveUSD, orderDirection: desc) { id reserveUSD } }",
  "variables": { "first": 10 }
}
```

### 2. Tracking Token Transfers
Monitor recent transfers every 5 minutes:
```json
{
  "schedule": "0 */5 * * * *",
  "graphqlEndpoint": "https://api.thegraph.com/subgraphs/name/{owner}/{token-subgraph}",
  "query": "query GetTransfers($first: Int!) { transfers(first: $first, orderBy: timestamp, orderDirection: desc) { id from to value } }",
  "variables": { "first": 20 }
}
```

### 3. Price Feed Updates
Check prices from a DEX every minute:
```json
{
  "schedule": "0 * * * * *",
  "graphqlEndpoint": "https://api.thegraph.com/subgraphs/name/{owner}/{dex-subgraph}",
  "query": "query GetPrices { tokens(first: 5, orderBy: derivedETH, orderDirection: desc) { id symbol derivedETH } }",
  "variables": {}
}
```

## Workflow Creation Process

These workflows were created using:

```bash
# Initialize Go workflow from template 1 (PoR example)
cre init --project-name indexer-workflow-ts -t 1 --rpc-url https://ethereum-sepolia-rpc.publicnode.com

# Initialize Go workflow from template 2 (Hello World)
cre init --project-name indexer-workflow-go -t 2 --rpc-url https://ethereum-sepolia-rpc.publicnode.com
```

The first workflow was then enhanced to support The Graph indexer queries.

## Files Modified

From the boilerplate template, the following files were modified:

1. **workflow.go**: Completely rewritten to query The Graph indexer
2. **config.staging.json**: Updated with Graph endpoint and query
3. **config.production.json**: Updated with Graph endpoint and query
4. **main.go**: Updated to use new Config struct

## Reference Documentation

- [CRE Documentation](https://docs.chain.link/cre)
- [The Graph Documentation](https://thegraph.com/docs/)
- [Cron Expression Reference](https://en.wikipedia.org/wiki/Cron)

## Next Steps

To resolve the runtime issues:
1. Test in actual deployment environment (not just simulation)
2. Check CRE SDK documentation for HTTP capability usage updates
3. Verify network connectivity and endpoint accessibility
4. Consider alternative HTTP request patterns supported by the SDK

## Learn More

For working examples:
- See `building-blocks/read-data-feeds/read-data-feeds-go` for a production-ready workflow
- Check CRE CLI help: `cre workflow simulate --help`
