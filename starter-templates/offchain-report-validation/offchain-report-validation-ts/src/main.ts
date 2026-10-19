import { isHex } from 'viem'
import {
	type Environment,
	PRODUCTION_ENVIRONMENT,
	type ReportInput,
	type VerificationResult,
	verifyReport,
} from './report'
import { sampleReport } from './sample-report'

const $ = <T extends HTMLElement>(id: string): T => {
	const el = document.getElementById(id)
	if (!el) throw new Error(`#${id} not found`)
	return el as T
}

const reportJsonEl = $<HTMLTextAreaElement>('report-json')
const rpcUrlEl = $<HTMLInputElement>('rpc-url')
const chainIdEl = $<HTMLInputElement>('chain-id')
const registryEl = $<HTMLInputElement>('registry-address')
const verifyBtn = $<HTMLButtonElement>('verify')
const resultEl = $<HTMLElement>('result')

rpcUrlEl.value = PRODUCTION_ENVIRONMENT.rpcUrl
chainIdEl.value = String(PRODUCTION_ENVIRONMENT.chainId)
registryEl.value = PRODUCTION_ENVIRONMENT.registryAddress

$('load-sample').addEventListener('click', () => {
	reportJsonEl.value = JSON.stringify(sampleReport, null, 2)
})

function parseReportInput(raw: string): ReportInput {
	let parsed: unknown
	try {
		parsed = JSON.parse(raw)
	} catch {
		throw new Error('Input is not valid JSON')
	}
	const obj = parsed as Record<string, unknown>
	if (!obj || typeof obj !== 'object') throw new Error('Expected a JSON object')
	if (typeof obj.rawReport !== 'string' || !isHex(obj.rawReport)) {
		throw new Error('"rawReport" must be a 0x-prefixed hex string')
	}
	if (typeof obj.reportContext !== 'string' || !isHex(obj.reportContext)) {
		throw new Error('"reportContext" must be a 0x-prefixed hex string')
	}
	if (
		!Array.isArray(obj.signatures) ||
		obj.signatures.length === 0 ||
		obj.signatures.some((s) => typeof s !== 'string' || !isHex(s))
	) {
		throw new Error('"signatures" must be a non-empty array of 0x-prefixed hex strings')
	}
	return {
		rawReport: obj.rawReport,
		reportContext: obj.reportContext,
		signatures: obj.signatures as `0x${string}`[],
	}
}

function readEnvironment(): Environment {
	const chainId = Number(chainIdEl.value)
	if (!Number.isInteger(chainId) || chainId <= 0) {
		throw new Error('Chain ID must be a positive integer')
	}
	return {
		chainId,
		registryAddress: registryEl.value.trim(),
		rpcUrl: rpcUrlEl.value.trim(),
	}
}

function escapeHtml(s: string): string {
	return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
}

function renderError(err: unknown): void {
	const msg =
		err instanceof AggregateError
			? err.errors.map((e) => (e instanceof Error ? e.message : String(e))).join('\n')
			: err instanceof Error
				? err.message
				: String(err)
	resultEl.innerHTML = `
		<h2>Verification failed</h2>
		<div class="error-box">${escapeHtml(msg)}</div>
	`
	resultEl.classList.remove('hidden')
}

function renderResult(r: VerificationResult): void {
	const sigRows = r.signatures
		.map(
			(s, i) => `
		<tr>
			<td>${i + 1}</td>
			<td class="mono">${s.recovered ? escapeHtml(s.recovered) : '—'}</td>
			<td>${s.signerId !== undefined ? s.signerId : '—'}</td>
			<td>${escapeHtml(s.status)}</td>
		</tr>`,
		)
		.join('')

	const failures =
		r.failures.length > 0
			? `<div class="error-box">${r.failures.map(escapeHtml).join('\n')}</div>`
			: ''

	const bodyText = (() => {
		try {
			return new TextDecoder('utf-8', { fatal: true }).decode(r.header.body)
		} catch {
			return null
		}
	})()

	resultEl.innerHTML = `
		<h2>Result <span class="badge ${r.valid ? 'ok' : 'fail'}">${r.valid ? 'VALID' : 'INVALID'}</span></h2>
		${failures}
		<dl class="kv">
			<dt>Report hash</dt><dd>${r.reportHash}</dd>
			<dt>DON ID</dt><dd>${r.header.donId}</dd>
			<dt>DON config version</dt><dd>${r.header.donConfigVersion}</dd>
			<dt>Fault tolerance f</dt><dd>${r.f} (requires f+1 = ${r.requiredSignatures} unique DON signatures, got ${r.acceptedSignatures})</dd>
			<dt>Version</dt><dd>${r.header.version}</dd>
			<dt>Execution ID</dt><dd>${r.header.executionId}</dd>
			<dt>Timestamp</dt><dd>${new Date(r.header.timestamp * 1000).toISOString()}</dd>
			<dt>Workflow ID</dt><dd>${r.header.workflowId}</dd>
			<dt>Workflow name</dt><dd>${escapeHtml(r.header.workflowName)}</dd>
			<dt>Workflow owner</dt><dd>0x${r.header.workflowOwner}</dd>
			<dt>Report ID</dt><dd>0x${r.header.reportId}</dd>
			<dt>Body</dt><dd>${bodyText !== null ? escapeHtml(bodyText) : `${r.header.body.length} bytes (binary)`}</dd>
		</dl>
		<h2>Signatures</h2>
		<table>
			<thead>
				<tr><th>#</th><th>Recovered signer</th><th>Node op. ID</th><th>Status</th></tr>
			</thead>
			<tbody>${sigRows}</tbody>
		</table>
		<p class="muted">DON signer set (${r.donSigners.length} addresses) read from registry
		<code>${escapeHtml(registryEl.value)}</code> on chain ${escapeHtml(chainIdEl.value)}.</p>
	`
	resultEl.classList.remove('hidden')
}

verifyBtn.addEventListener('click', async () => {
	verifyBtn.disabled = true
	verifyBtn.textContent = 'Verifying…'
	try {
		const input = parseReportInput(reportJsonEl.value)
		const env = readEnvironment()
		renderResult(await verifyReport(input, env))
	} catch (err) {
		renderError(err)
	} finally {
		verifyBtn.disabled = false
		verifyBtn.textContent = 'Verify report'
	}
})
