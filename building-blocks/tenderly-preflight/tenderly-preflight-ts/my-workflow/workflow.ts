import {
	cre,
	type CronPayload,
	type Runtime,
} from '@chainlink/cre-sdk';
import {
	TenderlyPreflight,
	tenderlyConfigSchema,
	transactionSchema,
} from '@tenderly/cre-preflight';
import { z } from 'zod';

// ---------- Config ----------

export const configSchema = z.object({
	// e.g. "0 */5 * * * *" (every 5 minutes, at second 0)
	schedule: z.string(),
	// Tenderly Virtual Environment settings. accessKeySecretId is the *name* of
	// the CRE secret holding your Tenderly access key (see secrets.yaml), never
	// the key itself. See https://github.com/Tenderly/tenderly-cre-preflight#configuration
	tenderly: tenderlyConfigSchema,
	// The transaction to preflight. In a real workflow you would build this from
	// whatever the execution just worked out (a fetched price, a report, an
	// on-chain decision) and pass it as an argument instead. It lives in config
	// here only so the template has something to simulate out of the box.
	transaction: transactionSchema,
});

type Config = z.infer<typeof configSchema>;

// ---------- Handler ----------

export function onCron(runtime: Runtime<Config>, _payload: CronPayload): string {
	const tenderly = new TenderlyPreflight(runtime.config.tenderly);

	// Every node forks the chain at the same agreed block, mines the transaction
	// on its own Virtual Environment, and the per-node verdicts go through
	// consensus. Nothing reaches the real chain.
	const verdict = tenderly.sendTransaction(runtime, runtime.config.transaction);

	switch (verdict.outcome) {
		case 'success':
			runtime.log(`preflight passed | gasUsed=${verdict.gasUsed}`);
			// This is where a real workflow would commit to the on-chain write,
			// e.g. runtime.report(...) followed by evmClient.writeReport(...).
			return 'preflight passed';

		case 'reverted':
		case 'rejected':
			// A real answer about the transaction: sending it on-chain would fail.
			// Skip the write and say why.
			runtime.log(`skipping write, would fail | outcome=${verdict.outcome} reason=${verdict.reason}`);
			return `skipped: ${verdict.outcome}`;

		case 'misconfigured':
			// The request can never succeed as written: the config, the Vault
			// secret, or the funding passed with the call. Every run fails until
			// someone fixes it — this is the outcome worth alerting on.
			runtime.log(`tenderly rejected the request: ${verdict.reason}`);
			return 'skipped: misconfigured';

		case 'unavailable':
		case 'oversized':
		case 'indeterminate':
			// No verdict can be trusted. Fail closed rather than writing blind.
			runtime.log(`skipping write, no verdict: ${verdict.outcome}`);
			return `skipped: ${verdict.outcome}`;
	}
}

// ---------- Init ----------

export function initWorkflow(config: Config) {
	const cron = new cre.capabilities.CronCapability();
	return [
		cre.handler(
			cron.trigger({ schedule: config.schedule }),
			onCron,
		),
	];
}
