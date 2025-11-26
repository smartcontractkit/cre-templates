import { cre, Runner, type NodeRuntime, type Runtime } from "@chainlink/cre-sdk"

type Config = {
  schedule: string
  graphqlEndpoint: string
  query: string
  variables?: Record<string, unknown>
}

type GraphQLRequest = {
  query: string
  variables?: Record<string, unknown>
}

type GraphQLResponse = {
  data?: unknown
  errors?: unknown[]
}

const initWorkflow = (config: Config) => {
  const cron = new cre.capabilities.CronCapability()

  return [cre.handler(cron.trigger({ schedule: config.schedule }), onIndexerCronTrigger)]
}

// fetchGraphData is the function passed to the runInNodeMode helper.
// It contains the logic for making the GraphQL request and parsing the response.
const fetchGraphData = (nodeRuntime: NodeRuntime<Config>): string => {
  const httpClient = new cre.capabilities.HTTPClient()

  // Prepare GraphQL request
  const gqlRequest: GraphQLRequest = {
    query: nodeRuntime.config.query,
    variables: nodeRuntime.config.variables,
  }

  const requestBody = JSON.stringify(gqlRequest)

  const req = {
    url: nodeRuntime.config.graphqlEndpoint,
    method: "POST" as const,
    headers: {
      "Content-Type": "application/json",
    },
    body: new TextEncoder().encode(requestBody),
  }

  // Send the request using the HTTP client
  // @ts-ignore
    const resp = httpClient.sendRequest(nodeRuntime, req).result()

  // Parse the GraphQL response
  const bodyText = new TextDecoder().decode(resp.body)
  const gqlResponse: GraphQLResponse = JSON.parse(bodyText)

  // Check for GraphQL errors
  if (gqlResponse.errors && gqlResponse.errors.length > 0) {
    nodeRuntime.log(`GraphQL errors: ${JSON.stringify(gqlResponse.errors)}`)
    throw new Error(`GraphQL query failed: ${JSON.stringify(gqlResponse.errors)}`)
  }

  if (!gqlResponse.data) {
    throw new Error("No data returned from GraphQL query")
  }

  nodeRuntime.log("Successfully fetched data from indexer")

  // Return the data as a JSON string
  return JSON.stringify(gqlResponse.data)
}

const onIndexerCronTrigger = (runtime: Runtime<Config>): string => {
  const timestamp = new Date().toISOString()

  runtime.log(`Cron triggered | timestamp=${timestamp}`)
  runtime.log(`Querying The Graph indexer | endpoint=${runtime.config.graphqlEndpoint}`)

  // Use runInNodeMode to execute the offchain fetch.
  // The Graph returns deterministic data across all nodes.
  // We define a simple aggregation that takes the first result since all nodes
  // should return identical data from The Graph.
  const firstResultAggregation = (results: string[]) => results[0]
  // @ts-ignore
    const result = runtime.runInNodeMode(fetchGraphData, firstResultAggregation)().result()

  runtime.log(`Indexer data fetched successfully | timestamp=${timestamp}`)

  // Format output
  const output = {
    timestamp,
    endpoint: runtime.config.graphqlEndpoint,
    data: JSON.parse(result),
  }

  // Return a formatted JSON string
  return JSON.stringify(output, null, 2)
}

export async function main() {
  const runner = await Runner.newRunner<Config>()
  await runner.run(initWorkflow)
}

main()
