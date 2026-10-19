import {
	bytesToBase64,
	bytesToHex,
	consensusIdenticalAggregation,
	CronCapability,
	EVMClient,
	getNetwork,
	handler,
	HTTPClient,
	type HTTPSendRequester,
	prepareReportRequest,
	type Report,
	Runner,
	type Runtime,
	TxStatus,
} from '@chainlink/cre-sdk'
import { type Address, encodeAbiParameters, encodeFunctionData, type Hex, parseAbiParameters } from 'viem'
import { Storage, StorageABI } from '../contracts/evm/ts/generated'

export type Config = {
	schedule: string
	backendUrl: string
	evms: Array<{
		storageAddress: string
		chainSelectorName: string
		gasLimit: string
	}>
}

/**
 * Runs in node mode: POSTs the signed report to the backend as JSON using the
 * HTTP capability's sendReport helper. The sendReport callback receives the
 * ReportResponse (rawReport, reportContext, sigs) — the supported way to read a
 * report's fields for off-chain transport.
 *
 * Returns a JSON string with the push outcome and the exact payload sent, so the
 * DON-mode handler can log them.
 */
const pushReport = (
	sendRequester: HTTPSendRequester,
	report: Report,
	backendUrl: string,
): string => {
	// Built inside the sendReport callback — it runs before the request is sent.
	let payload = ''
	let statusCode = -1
	try {
		const resp = sendRequester
			.sendReport(report, (reportResponse) => {
				payload = JSON.stringify({
					rawReport: bytesToHex(reportResponse.rawReport),
					reportContext: bytesToHex(reportResponse.reportContext),
					signatures: reportResponse.sigs.map((s) => bytesToHex(s.signature)),
				})
				return {
					url: backendUrl,
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: bytesToBase64(new TextEncoder().encode(payload)),
				}
			})
			.result()
		statusCode = resp.statusCode
	} catch {
		// Endpoint unreachable — the payload is still returned for logging.
	}
	return JSON.stringify({ statusCode, payload })
}

export const onCronTrigger = (runtime: Runtime<Config>): string => {
	const evmConfig = runtime.config.evms[0]
	const network = getNetwork({
		chainFamily: 'evm',
		chainSelectorName: evmConfig.chainSelectorName,
		isTestnet: true,
	})
	if (!network) {
		throw new Error(`Network not found for chain selector name: ${evmConfig.chainSelectorName}`)
	}

	const evmClient = new EVMClient(network.chainSelector.selector)
	const storage = new Storage(evmClient, evmConfig.storageAddress as Address)

	// 1. The value to write onchain: the current execution timestamp.
	const value = BigInt(Math.floor(runtime.now().getTime() / 1000))

	// 2. Encode the write as a call to the consumer's onReport(bytes, bytes),
	//    exactly like the generated binding's writeReportFromOnReport does.
	//    The consumer decodes `report` as abi.encode(uint256 value).
	const reportBody = encodeAbiParameters(parseAbiParameters('uint256 value'), [value])
	const callData = encodeFunctionData({
		abi: StorageABI,
		functionName: 'onReport',
		args: ['0x' as Hex, reportBody],
	})

	// 3. Every chain write through the chain writer is delivered as a DON-signed
	//    report — the report comes with the write by default. We generate it
	//    explicitly here (this is precisely what the generated binding's
	//    writeReportFromOnReport does internally) so that we keep a reference to
	//    the signed report and can export it for offchain verification.
	const signedReport = runtime.report(prepareReportRequest(callData)).result()

	// 4. Submit the write onchain via the chain writer.
	const resp = evmClient
		.writeReport(runtime, {
			receiver: storage.address,
			report: signedReport,
			gasConfig: { gasLimit: evmConfig.gasLimit },
		})
		.result()

	if (resp.txStatus !== TxStatus.SUCCESS) {
		throw new Error(`Failed to write report: ${resp.errorMessage || resp.txStatus}`)
	}
	runtime.log(`Value ${value} written onchain`)

	// 5. Publish the signed report to your offchain backend over HTTP (via the
	//    sendReport helper), and log it field by field (one line per field, one
	//    signature per line, so nothing is truncated in the logs). The lines
	//    between REPORT_JSON_BEGIN / REPORT_JSON_END reassemble into the exact
	//    JSON shape the offchain validator app expects.
	const exportResult = new HTTPClient()
		.sendRequest(
			runtime,
			pushReport,
			consensusIdenticalAggregation(),
		)(signedReport, runtime.config.backendUrl)
		.result()

	const { statusCode, payload } = JSON.parse(exportResult) as {
		statusCode: number
		payload: string
	}
	if (statusCode >= 200 && statusCode < 300) {
		runtime.log(`Report pushed to ${runtime.config.backendUrl} (HTTP ${statusCode})`)
	} else {
		runtime.log(
			`Could not push report to ${runtime.config.backendUrl} (status ${statusCode}) — the report is still logged below`,
		)
	}

	const reportPayload = JSON.parse(payload) as {
		rawReport: string
		reportContext: string
		signatures: string[]
	}
	runtime.log('REPORT_JSON_BEGIN')
	runtime.log(`{"rawReport":"${reportPayload.rawReport}",`)
	runtime.log(`"reportContext":"${reportPayload.reportContext}",`)
	runtime.log('"signatures":[')
	reportPayload.signatures.forEach((sig, i) => {
		runtime.log(`"${sig}"${i < reportPayload.signatures.length - 1 ? ',' : ''}`)
	})
	runtime.log(']}')
	runtime.log('REPORT_JSON_END')

	// Return values are only visible in simulation, not once the workflow is
	// deployed — everything observable must go through runtime.log.
	const result = bytesToHex(resp.txHash || new Uint8Array(32))
	runtime.log(`Workflow result (tx hash): ${result}`)
	return result
}

export const initWorkflow = (config: Config) => {
	const cron = new CronCapability()

	return [handler(cron.trigger({ schedule: config.schedule }), onCronTrigger)]
}

export async function main() {
	const runner = await Runner.newRunner<Config>()
	await runner.run(initWorkflow)
}
