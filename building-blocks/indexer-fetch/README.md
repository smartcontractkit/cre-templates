# CRE Indexer Data Feed Workflows

Workflows for pulling data from The Graph indexer with scheduled cron triggers. These workflows demonstrate the **pull pattern** where the workflow initiates and fetches data on a schedule.

## Directory Structure

```
building-blocks/indexer-fetch/
├── README.md (this file)
├── indexer-fetch-go/     (Go-based workflow)
│   └── my-workflow/
│       ├── workflow.go
│       ├── main.go
│       ├── config.staging.json
│       ├── config.production.json
│       └── workflow.yaml
└── indexer-fetch-ts/     (TypeScript-based workflow)
    └── workflow/
        ├── main.ts
        ├── config.staging.json
        ├── config.production.json
        ├── package.json
        └── workflow.yaml
```

## Overview

These workflows demonstrate how to:
- Query The Graph indexer using GraphQL
- Use cron triggers to schedule periodic data fetching
- Process and return JSON-formatted indexer data
- Implement the same functionality in both Go and TypeScript

Both workflows query the Uniswap V4 subgraph on The Graph and fetch:
- Pool manager statistics (pool count, transaction count, total volume)
- ETH price data from bundles

## Workflows

### 1. indexer-fetch-go (Go Implementation)

**Language:** Go

**Features:**
- Uses `http.SendRequest` pattern from CRE Go SDK
- Implements `ConsensusIdenticalAggregation` for deterministic data
- Returns formatted JSON with timestamp and endpoint info

**Running the workflow:**
```bash
cd building-blocks/indexer-fetch/indexer-fetch-go
cre workflow simulate my-workflow --target staging-settings
```

### 2. indexer-fetch-ts (TypeScript Implementation)

**Language:** TypeScript

**Features:**
- Uses `runInNodeMode` pattern from CRE TypeScript SDK
- Implements custom first-result aggregation for deterministic data
- Returns formatted JSON with timestamp and endpoint info

**Running the workflow:**
```bash
cd building-blocks/indexer-fetch/indexer-fetch-ts
cre workflow simulate workflow --target staging-settings
```

## Configuration

Both workflows use the same configuration structure in their respective `config.staging.json` files:

```json
{
  "schedule": "0 * * * * *",
  "graphqlEndpoint": "https://gateway.thegraph.com/api/bca58895bc60dcb319e3cbdfd989b964/subgraphs/id/Gqm2b5J85n1bhCyDMpGbtbVn4935EvvdyHdHrx3dibyj",
  "query": "{ poolManagers(first: 5) { id poolCount txCount totalVolumeUSD } bundles(first: 5) { id ethPriceUSD } }",
  "variables": {}
}
```

### Configuration Options

- **schedule**: Cron expression in 6-field format (second minute hour day month weekday)
  - `"0 * * * * *"` - Every minute at second 0
  - `"*/30 * * * * *"` - Every 30 seconds
  - `"0 */5 * * * *"` - Every 5 minutes at second 0

- **graphqlEndpoint**: The Graph API endpoint URL
  - Gateway endpoint: `https://gateway.thegraph.com/api/{api-key}/subgraphs/id/{subgraph-id}`
  - Studio endpoint: `https://api.studio.thegraph.com/query/{id}/{name}/version/latest`

- **query**: GraphQL query string
  - Simple queries without variables work best
  - See The Graph documentation for query syntax

- **variables**: Object with variables for the GraphQL query (optional)

## Key Implementation Details

### Go Implementation

```go
// Uses HTTP SendRequest pattern with consensus aggregation
client := &http.Client{}
result, err := http.SendRequest(
    config,
    runtime,
    client,
    fetchGraphData,
    cre.ConsensusIdenticalAggregation[string](),
).Await()

// GraphQL request structure
gqlRequest := GraphQLRequest{
    Query:     config.Query,
    Variables: config.Variables,
}

// Makes POST request
httpResp, err := sendRequester.SendRequest(&http.Request{
    Method: "POST",
    Url:    config.GraphqlEndpoint,
    Headers: map[string]string{
        "Content-Type": "application/json",
    },
    Body: requestBody,
}).Await()
```

### TypeScript Implementation

```typescript
// Uses runInNodeMode pattern with custom aggregation
const firstResultAggregation = (results: string[]) => results[0]
const result = runtime.runInNodeMode(
    fetchGraphData,
    firstResultAggregation
)().result()

// GraphQL request structure
const gqlRequest: GraphQLRequest = {
  query: nodeRuntime.config.query,
  variables: nodeRuntime.config.variables,
}

// Makes POST request
const resp = httpClient.sendRequest(nodeRuntime, {
  url: nodeRuntime.config.graphqlEndpoint,
  method: "POST" as const,
  headers: {
    "Content-Type": "application/json",
  },
  body: new TextEncoder().encode(requestBody),
}).result()
```

## Setup and Testing

### Prerequisites

**For Go workflow:**
1. Install CRE CLI
2. Login: `cre login`
3. Go 1.23+ installed

**For TypeScript workflow:**
1. Install CRE CLI
2. Login: `cre login`
3. Bun installed (or Node.js)
4. Run `bun install` in the workflow directory

### Running the Workflows

**Go Workflow:**
```bash
cd building-blocks/indexer-fetch/indexer-fetch-go
cre workflow simulate my-workflow --target staging-settings
```

**TypeScript Workflow:**
```bash
cd building-blocks/indexer-fetch/indexer-fetch-ts
cre workflow simulate workflow --target staging-settings
```

### Expected Output

Both workflows return JSON output like:

```json
{
  "timestamp": "2025-11-18T18:43:08.452Z",
  "endpoint": "https://gateway.thegraph.com/api/.../subgraphs/id/...",
  "data": {
    "bundles": [
      {
        "ethPriceUSD": "3157.000458184067393927942592490315",
        "id": "1"
      }
    ],
    "poolManagers": [
      {
        "id": "0x498581ff718922c3f8e6a244956af099b2652b2b",
        "poolCount": "5123368",
        "totalVolumeUSD": "5611562100.854190095192400782985064",
        "txCount": "480580367"
      }
    ]
  }
}
```

## Status

### ✅ Both Workflows Working

- ✅ Workflow structure and configuration
- ✅ Compilation and WASM generation
- ✅ Cron trigger setup
- ✅ GraphQL request formatting
- ✅ HTTP requests to The Graph
- ✅ Error handling
- ✅ Successful simulation and execution

## Example Use Cases

### 1. Monitoring Uniswap V4 Pools
Query pool statistics every minute:
```json
{
  "schedule": "0 * * * * *",
  "graphqlEndpoint": "https://gateway.thegraph.com/api/{key}/subgraphs/id/{id}",
  "query": "{ poolManagers(first: 5) { id poolCount totalVolumeUSD } }",
  "variables": {}
}
```

### 2. Tracking Token Prices
Monitor token prices every 30 seconds:
```json
{
  "schedule": "*/30 * * * * *",
  "graphqlEndpoint": "https://gateway.thegraph.com/api/{key}/subgraphs/id/{id}",
  "query": "{ tokens(first: 10, orderBy: volumeUSD, orderDirection: desc) { id symbol volumeUSD } }",
  "variables": {}
}
```

### 3. DeFi Protocol Metrics
Check protocol statistics every 5 minutes:
```json
{
  "schedule": "0 */5 * * * *",
  "graphqlEndpoint": "https://gateway.thegraph.com/api/{key}/subgraphs/id/{id}",
  "query": "{ protocols(first: 1) { totalValueLockedUSD totalVolumeUSD txCount } }",
  "variables": {}
}
```

## Comparison: Go vs TypeScript

| Feature | Go | TypeScript |
|---------|-----|------------|
| **HTTP Pattern** | `http.SendRequest` | `runInNodeMode` with `HTTPClient` |
| **Consensus** | `ConsensusIdenticalAggregation[string]()` | Custom `firstResultAggregation` |
| **Type Safety** | Compile-time with structs | Compile-time with TypeScript types |
| **Error Handling** | Go error returns | JavaScript try/catch |
| **Entry Point** | `main.go` | `main.ts` |
| **Workflow Directory** | `my-workflow/` | `workflow/` |

## Reference Documentation

- [CRE Documentation](https://docs.chain.link/cre)
- [The Graph Documentation](https://thegraph.com/docs/)
- [Cron Expression Reference](https://en.wikipedia.org/wiki/Cron)
- [CRE TypeScript SDK](https://www.npmjs.com/package/@chainlink/cre-sdk)

## Related Patterns

This is a **pull pattern** workflow where the workflow initiates data fetching on a schedule. For the complementary **push pattern** (event-driven workflows triggered by indexer events), see the `indexer-events` building block.

## Learn More

For other workflow examples:
- See `building-blocks/read-data-feeds` for reading on-chain data feeds
- See `building-blocks/kv-store` for key-value storage patterns
- Check CRE CLI help: `cre workflow simulate --help`
