import {
	consensusIdenticalAggregation,
	cre,
	type HTTPSendRequester,
	json,
	Runner,
	type Runtime,
} from '@chainlink/cre-sdk'

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

// fetchGraphData is the function passed to the HTTP capability's sendRequest helper.
// It contains the logic for making the GraphQL request and parsing the response.
const fetchGraphData = (sendRequester: HTTPSendRequester, config: Config): string => {
	// Prepare GraphQL request
	const gqlRequest: GraphQLRequest = {
		query: config.query,
		variables: config.variables,
	}

	const req = {
		url: config.graphqlEndpoint,
		method: 'POST' as const,
		headers: {
			'Content-Type': 'application/json',
		},
		body: Buffer.from(JSON.stringify(gqlRequest)).toString('base64'),
	}

	// Send the request using the HTTP client
	const resp = sendRequester.sendRequest(req).result()

	// Parse the GraphQL response
	const gqlResponse = json(resp) as GraphQLResponse

	// Check for GraphQL errors
	if (gqlResponse.errors && gqlResponse.errors.length > 0) {
		throw new Error(`GraphQL query failed: ${JSON.stringify(gqlResponse.errors)}`)
	}

	if (!gqlResponse.data) {
		throw new Error('No data returned from GraphQL query')
	}

	// Return the data as a JSON string
	return JSON.stringify(gqlResponse.data)
}

const onIndexerCronTrigger = (runtime: Runtime<Config>): string => {
	const timestamp = new Date().toISOString()

	runtime.log(`Cron triggered | timestamp=${timestamp}`)
	runtime.log(`Querying The Graph indexer | endpoint=${runtime.config.graphqlEndpoint}`)

	const httpClient = new cre.capabilities.HTTPClient()

	// Use sendRequest sugar to execute the offchain fetch.
	// The Graph returns deterministic data across all nodes.
	// We use identical aggregation since all nodes should return identical data from The Graph.
	const result = httpClient
		.sendRequest(
			runtime,
			fetchGraphData,
			consensusIdenticalAggregation<string>(),
		)(runtime.config)
		.result()

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

