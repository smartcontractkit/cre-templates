import type { ReportInput } from './report'

/**
 * A real CRE report (from the CRE TypeScript SDK test fixtures), produced by
 * the production DON (donId 1) whose config lives in the capabilities registry
 * on Ethereum mainnet. The report body is "Hello, World!".
 *
 * This verifies against live on-chain DON state, so if the DON membership ever
 * rotates enough the sample may stop validating — the UI will show exactly why.
 */
export const sampleReport: ReportInput = {
	rawReport:
		'0x010d46e9e05b28a322cd96c94f1f60c236a1b52bd446d98396f1a84afae388836a69c96faa000000010000000100a3bb3b17a1053fe409147c2b170af4866c182070905e5a0f314be04c8296a336613263393832646464b0f2d38245dd6d397ebbdb5a814b753d56c30715000148656c6c6f2c20576f726c6421',
	reportContext:
		'0x000e8ce31db48e5e44619d24d9dadfc5f22a34db8205b2b25cd831eab02244c500000000000000000000000000000000000000000000000000000001612bcb000000000000000000000000000000000000000000000000000000000000000000',
	signatures: [
		'0x452d9e1fef501946b5e0a52e32dd4ec48a796c8c30a4f62332006a1fb5a7d4b578b6ad0b7a15f3c4ef46bffb891884b50ac23f5fcfa5cf75e7a7329d6fed7a6400',
		'0x3b48aa3db3d6398f1d69f649e266d4388d3d176b76ff58de72ae4ca30e27559b5789a491b8f94578f793bc5781ad5c50dbf48bb79e61b929f625162a4c79172300',
		'0xa0fa2428e44e7ffa9452433b3275bd79819db0620b7d2090fea21065a1dfed1a7c5135540c6e44c20be4d36d18c1859253eca53476f45ea3f8a54ddc97d4b42c01',
		'0xf868744afd82e657110ec1e214fe91e6abc10ea84edcf69b5076ee5bf9126e056ba50eb078d1b594140ddaac11aa0b52fe9ca4f629b334b243e38a345e0d519501',
	],
}
