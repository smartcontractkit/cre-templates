import {
	bytesToHex,
	ConsensusAggregationByFields,
	type CronPayload,
	cre,
	getNetwork,
	type HTTPSendRequester,
	hexToBase64,
	median,
	Runner,
	type Runtime,
	TxStatus,
} from '@chainlink/cre-sdk'
import { encodeAbiParameters, Hex, parseUnits } from 'viem'
import { z } from 'zod'

const configSchema = z.object({
	schedule: z.string(),
	url: z.string(),
	dataIdHex: z.string(),
	evms: z.array(
		z.object({
			dataFeedsCacheAddress: z.string(),
			chainName: z.string(),
			gasLimit: z.string(),
		}),
	),
})

type Config = z.infer<typeof configSchema>
type EVMConfig = z.infer<typeof configSchema.shape.evms.element>

interface NAVResponse {
	_aggregatedCollateral: number
	_totalOwedM: number
	totalCollateral: number
}

interface NAVInfo {
	nav: number
}

// Utility function to safely stringify objects with bigints
const safeJsonStringify = (obj: any): string =>
	JSON.stringify(obj, (_, value) => (typeof value === 'bigint' ? value.toString() : value), 2)

const fetchNAVInfo = (sendRequester: HTTPSendRequester, config: Config): NAVInfo => {

	const body = JSON.stringify({
		method: 'navDetails',
	})

	const response = sendRequester.sendRequest(
		{
			url: config.url,
			method: 'POST',
			body: Buffer.from(body).toString("base64"),
		}
	).result()

	if (response.statusCode !== 200) {
		throw new Error(`HTTP request failed with status: ${response.statusCode}`)
	}

	const responseText = Buffer.from(response.body).toString('utf-8')
	const navResp: NAVResponse = JSON.parse(responseText)

	const nav = navResp._aggregatedCollateral - navResp._totalOwedM

	return {
		nav: nav,
	}
}

export type ReceivedDecimaldReport = {
  dataId: string
  timestamp: number
  answer: bigint
}

const encodeReceivedDecimaldReports = (
  reports: ReceivedDecimaldReport[],
): Hex => {
  return encodeAbiParameters(
    [
      {
        name: 'reports',
        type: 'tuple[]',
        components: [
          { name: 'dataId', type: 'bytes32' },
          { name: 'timestamp', type: 'uint32' },
          { name: 'answer', type: 'uint224' },
        ],
      },
    ],
    [
      reports.map((r) => ({
        dataId: hexToBytes32RightPadded(r.dataId),
        timestamp: r.timestamp,
        answer: r.answer,
      })),
    ],
  )
}

const hexToBytes32RightPadded = (input: string): Hex => {
  let hex = input.toLowerCase()
  if (hex.startsWith('0x')) hex = hex.slice(2)
  if (hex.length % 2 !== 0) hex = '0' + hex // normalize to full bytes

  const byteLen = hex.length / 2
  if (byteLen > 32) {
    throw new Error(
      `hex string decodes to ${byteLen} bytes, which exceeds 32 bytes`,
    )
  }

  // right-pad to 32 bytes
  const padded = hex.padEnd(64, '0')
  return ('0x' + padded) as Hex
}

const updateNAV = (
	evmConfig: EVMConfig,
	runtime: Runtime<Config>,
	navScaled: bigint,
): string => {

	const network = getNetwork({
		chainFamily: 'evm',
		chainSelectorName: evmConfig.chainName,
		isTestnet: true,
	})

	if (!network) {
		throw new Error(`Network not found for chain selector name: ${evmConfig.chainName}`)
	}

	const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)

	runtime.log(
		`Updating NAV navScaled ${navScaled.toString()}`,
	)

	const reportData = encodeReceivedDecimaldReports([
		{
			dataId: runtime.config.dataIdHex,
			timestamp: Math.floor(Date.now() / 1000),
			answer: navScaled,
		},
	])

	// Step 1: Generate report using consensus capability
	const reportResponse = runtime
		.report({
			encodedPayload: hexToBase64(reportData),
			encoderName: 'evm',
			signingAlgo: 'ecdsa',
			hashingAlgo: 'keccak256',
		})
		.result()

	const resp = evmClient
		.writeReport(runtime, {
			receiver: evmConfig.dataFeedsCacheAddress,
			report: reportResponse,
			gasConfig: {
				gasLimit: evmConfig.gasLimit,
			},
		})
		.result()

	const txStatus = resp.txStatus

	if (txStatus !== TxStatus.SUCCESS) {
		throw new Error(`Failed to write report: ${resp.errorMessage || txStatus}`)
	}

	const txHash = resp.txHash || new Uint8Array(32)

	runtime.log(`Write report transaction succeeded at txHash: ${bytesToHex(txHash)}`)

	return txHash.toString()
}

const doNAV = (runtime: Runtime<Config>): string => {
	runtime.log(`fetching nav url ${runtime.config.url}`)

	const httpCapability = new cre.capabilities.HTTPClient()
	const navInfo = httpCapability
		.sendRequest(
			runtime,
			fetchNAVInfo,
			ConsensusAggregationByFields<NAVInfo>({
				nav: median,
			}),
		)(runtime.config)
		.result()

	runtime.log(`NAVInfo ${safeJsonStringify(navInfo)}`)

	const navScaled = parseUnits(navInfo.nav.toString(), 18)
	runtime.log(`NavScaled ${navScaled.toString()}`)

	for (const evmConfig of runtime.config.evms) {
		runtime.log(`Updating NAV on chain ${evmConfig.chainName}`)

		updateNAV(evmConfig, runtime, navScaled)
	}

	return navInfo.nav.toString()
}

const onCronTrigger = (runtime: Runtime<Config>, payload: CronPayload): string => {
	if (!payload.scheduledExecutionTime) {
		throw new Error('Scheduled execution time is required')
	}

	runtime.log('Running CronTrigger')

	return doNAV(runtime)
}

const initWorkflow = (config: Config) => {
	const cronTrigger = new cre.capabilities.CronCapability()
	return [
		cre.handler(
			cronTrigger.trigger({
				schedule: config.schedule,
			}),
			onCronTrigger,
		),
	]
}

export async function main() {
	const runner = await Runner.newRunner<Config>({ configSchema });
	await runner.run(initWorkflow);
}

main()
