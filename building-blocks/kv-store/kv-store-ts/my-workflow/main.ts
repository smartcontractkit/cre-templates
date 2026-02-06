import { sha256 } from '@noble/hashes/sha2'
import { hmac } from '@noble/hashes/hmac'
import { bytesToHex as nobleToHex } from '@noble/hashes/utils'
import {
	consensusIdenticalAggregation,
	consensusMedianAggregation,
	type CronPayload,
	cre,
	type HTTPSendRequester,
	Runner,
	type Runtime,
	text,
} from '@chainlink/cre-sdk'

// ---------------------------
// Types
// ---------------------------

type Config = {
	schedule: string
	aws_region: string
	s3_bucket: string
	s3_key: string
}

type AWSCredentials = {
	accessKeyId: string
	secretAccessKey: string
}

// ---------------------------
// AWS SigV4 Signing (pure JS)
// ---------------------------

const EMPTY_PAYLOAD_HASH = 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855'

function sha256Hex(data: string): string {
	return nobleToHex(sha256(new TextEncoder().encode(data)))
}

function sha256HexBytes(data: Uint8Array): string {
	return nobleToHex(sha256(data))
}

function hmacSha256(key: Uint8Array, data: string): Uint8Array {
	return hmac(sha256, key, new TextEncoder().encode(data))
}

function deriveSigningKey(
	secretKey: string,
	dateStamp: string,
	region: string,
	service: string,
): Uint8Array {
	const kDate = hmacSha256(new TextEncoder().encode('AWS4' + secretKey), dateStamp)
	const kRegion = hmacSha256(kDate, region)
	const kService = hmacSha256(kRegion, service)
	return hmacSha256(kService, 'aws4_request')
}

function formatTimestamp(date: Date): { amzDate: string; dateStamp: string } {
	const iso = date.toISOString().replace(/[-:]/g, '').replace(/\.\d{3}/, '')
	return {
		amzDate: iso, // e.g. "20250101T120000Z"
		dateStamp: iso.slice(0, 8), // e.g. "20250101"
	}
}

function signRequest(params: {
	method: string
	host: string
	path: string
	queryString: string
	headers: Record<string, string>
	payloadHash: string
	credentials: AWSCredentials
	region: string
	timestamp: Date
}): Record<string, string> {
	const { method, host, path, queryString, payloadHash, credentials, region, timestamp } = params
	const { amzDate, dateStamp } = formatTimestamp(timestamp)
	const service = 's3'
	const scope = `${dateStamp}/${region}/${service}/aws4_request`

	// Headers to sign (must be sorted)
	const signedHeaderEntries: [string, string][] = [
		['host', host],
		['x-amz-content-sha256', payloadHash],
		['x-amz-date', amzDate],
	]

	// Add any extra headers from params (e.g. content-type, content-length)
	for (const [k, v] of Object.entries(params.headers)) {
		const lower = k.toLowerCase()
		if (lower !== 'host' && lower !== 'x-amz-content-sha256' && lower !== 'x-amz-date') {
			signedHeaderEntries.push([lower, v])
		}
	}

	signedHeaderEntries.sort((a, b) => a[0].localeCompare(b[0]))

	const canonicalHeaders = signedHeaderEntries.map(([k, v]) => `${k}:${v.trim()}\n`).join('')
	const signedHeaders = signedHeaderEntries.map(([k]) => k).join(';')

	// Canonical request
	const canonicalRequest = [
		method,
		path,
		queryString,
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	].join('\n')

	// String to sign
	const stringToSign = ['AWS4-HMAC-SHA256', amzDate, scope, sha256Hex(canonicalRequest)].join(
		'\n',
	)

	// Signature
	const signingKey = deriveSigningKey(credentials.secretAccessKey, dateStamp, region, service)
	const signature = nobleToHex(hmacSha256(signingKey, stringToSign))

	const authorization = `AWS4-HMAC-SHA256 Credential=${credentials.accessKeyId}/${scope}, SignedHeaders=${signedHeaders}, Signature=${signature}`

	// Return all headers needed for the request
	const resultHeaders: Record<string, string> = {
		...params.headers,
		host: host,
		'x-amz-content-sha256': payloadHash,
		'x-amz-date': amzDate,
		Authorization: authorization,
	}

	return resultHeaders
}

// ---------------------------
// S3 Helpers
// ---------------------------

function s3Endpoint(config: Config): { host: string; path: string } {
	let escapedPath = encodeURI(config.s3_key)
	if (!escapedPath.startsWith('/')) {
		escapedPath = '/' + escapedPath
	}

	// Path-style for dotted buckets (SSL cert constraints), virtual-hosted otherwise
	if (config.s3_bucket.includes('.')) {
		return {
			host: `s3.${config.aws_region}.amazonaws.com`,
			path: `/${config.s3_bucket}${escapedPath}`,
		}
	}

	return {
		host: `${config.s3_bucket}.s3.${config.aws_region}.amazonaws.com`,
		path: escapedPath,
	}
}

const s3Read = (
	sendRequester: HTTPSendRequester,
	config: Config,
	creds: AWSCredentials,
): number => {
	const { host, path } = s3Endpoint(config)
	const fullURL = `https://${host}${path}`
	const timestamp = new Date()

	const headers = signRequest({
		method: 'GET',
		host,
		path,
		queryString: '',
		headers: {},
		payloadHash: EMPTY_PAYLOAD_HASH,
		credentials: creds,
		region: config.aws_region,
		timestamp,
	})

	const resp = sendRequester
		.sendRequest({
			url: fullURL,
			method: 'GET',
			headers,
		})
		.result()

	if (resp.statusCode === 404) {
		return 0
	}

	if (resp.statusCode >= 300) {
		throw new Error(`S3 GET failed: ${resp.statusCode} ${text(resp)}`)
	}

	const bodyStr = text(resp).trim()
	const value = Number(bodyStr)
	return Number.isNaN(value) ? 0 : value
}

const s3Write = (
	sendRequester: HTTPSendRequester,
	config: Config,
	creds: AWSCredentials,
	newValue: string,
): string => {
	const { host, path } = s3Endpoint(config)
	const fullURL = `https://${host}${path}`
	const timestamp = new Date()
	const bodyBytes = new TextEncoder().encode(newValue)
	const payloadHash = sha256HexBytes(bodyBytes)

	const headers = signRequest({
		method: 'PUT',
		host,
		path,
		queryString: '',
		headers: {
			'content-type': 'text/plain',
			'content-length': String(bodyBytes.length),
		},
		payloadHash,
		credentials: creds,
		region: config.aws_region,
		timestamp,
	})

	const resp = sendRequester
		.sendRequest({
			url: fullURL,
			method: 'PUT',
			headers,
			body: Buffer.from(newValue).toString('base64'),
			cacheSettings: { store: true, maxAge: '60s' },
		})
		.result()

	if (resp.statusCode >= 300) {
		throw new Error(`S3 PUT failed: ${resp.statusCode} ${text(resp)}`)
	}

	return newValue
}

// ---------------------------
// Workflow
// ---------------------------

const onCronTrigger = (runtime: Runtime<Config>, _payload: CronPayload): string => {
	runtime.log('Cron trigger fired. Fetching AWS credentials...')

	const accessKey = runtime.getSecret({ id: 'AWS_ACCESS_KEY_ID' }).result()
	const secretKey = runtime.getSecret({ id: 'AWS_SECRET_ACCESS_KEY' }).result()
	const creds: AWSCredentials = {
		accessKeyId: accessKey.value,
		secretAccessKey: secretKey.value,
	}

	runtime.log('AWS credentials fetched. Performing consensus read, then write.')

	const httpClient = new cre.capabilities.HTTPClient()

	// Phase 1: Read — median consensus on the numeric value
	const oldValue = httpClient
		.sendRequest(runtime, s3Read, consensusMedianAggregation<number>())
		(runtime.config, creds)
		.result()

	const newValue = (oldValue + 1).toString()
	runtime.log(`Consensus old value computed. Incrementing. old=${oldValue} new=${newValue}`)

	// Phase 2: Write — identical consensus (cache ensures same response across nodes)
	httpClient
		.sendRequest(runtime, s3Write, consensusIdenticalAggregation<string>())
		(runtime.config, creds, newValue)
		.result()

	runtime.log(`Workflow finished successfully. old=${oldValue} new=${newValue}`)

	return JSON.stringify({ oldValue, newValue })
}

const initWorkflow = (config: Config) => {
	const cron = new cre.capabilities.CronCapability()
	return [cre.handler(cron.trigger({ schedule: config.schedule }), onCronTrigger)]
}

export async function main() {
	const runner = await Runner.newRunner<Config>()
	await runner.run(initWorkflow)
}

main()
