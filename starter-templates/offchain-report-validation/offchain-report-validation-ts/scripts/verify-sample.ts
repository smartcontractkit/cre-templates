/**
 * CLI sanity check: verifies the bundled sample report against live on-chain
 * DON state using the exact same code the browser app uses.
 *
 *   npm run verify:sample
 */
import { PRODUCTION_ENVIRONMENT, verifyReport } from '../src/report'
import { sampleReport } from '../src/sample-report'

const result = await verifyReport(sampleReport, PRODUCTION_ENVIRONMENT)

console.log('report hash:       ', result.reportHash)
console.log('don id:            ', result.header.donId)
console.log('workflow name:     ', result.header.workflowName)
console.log('workflow owner:    ', result.header.workflowOwner)
console.log('body (utf8):       ', new TextDecoder().decode(result.header.body))
console.log('f:                 ', result.f, `(need f+1 = ${result.requiredSignatures} sigs)`)
console.log('don signers:       ', result.donSigners.length)
for (const sig of result.signatures) {
	console.log(`  sig ${sig.status.padEnd(15)} ${sig.recovered ?? ''} ${sig.signerId ?? ''}`)
}

if (!result.valid) {
	throw new Error(`Sample report failed validation: ${result.failures.join('; ')}`)
}
console.log('\nSample report is VALID')
