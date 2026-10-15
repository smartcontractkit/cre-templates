import {
	bytesToHex,
	CronCapability,
	EVMClient,
	getNetwork,
	handler,
	prepareReportRequest,
	Runner,
	type Runtime,
	TxStatus,
} from '@chainlink/cre-sdk'
import { type Address, encodeAbiParameters, encodeFunctionData, type Hex, parseAbiParameters } from 'viem'
import { Storage, StorageABI } from '../contracts/evm/ts/generated'

export type Config = {
	schedule: string
	evms: Array<{
		storageAddress: string
		chainSelectorName: string
		gasLimit: string
	}>
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

	// 5. Emit the signed report field by field (one log line each, signatures one
	//    per line) so long values don't get truncated in the logs. The lines
	//    between REPORT_JSON_BEGIN / REPORT_JSON_END reassemble into the exact
	//    JSON shape the offchain validator UI expects.
	const rr = signedReport.x_generatedCodeOnly_unwrap()
	const reportPayload = {
		rawReport: bytesToHex(rr.rawReport),
		reportContext: bytesToHex(rr.reportContext),
		signatures: rr.sigs.map((s) => bytesToHex(s.signature)),
	}

	// Here you would send the report to your offchain backend via HTTP
	// (e.g. HTTPClient POST <your-backend>/reports with JSON body
	// JSON.stringify(reportPayload)) so clients can verify it — intentionally
	// left as log lines in this example.
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
