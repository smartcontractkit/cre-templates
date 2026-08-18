import { cre, getNetwork, type CronPayload, type Runtime } from '@chainlink/cre-sdk';
import { CcipClient } from 'cre-ccip-sdk';
import { toHex } from 'viem';
import { z } from 'zod';

// ---------- Config ----------

export const configSchema = z.object({
	// e.g. "0 */10 * * * *" (every 10 minutes, at second 0)
	schedule: z.string(),
	// Source chain name, must match the RPC entry in project.yaml. e.g. "ethereum-testnet-sepolia"
	chainName: z.string(),
	// CCIP Router address on the source chain. See https://docs.chain.link/ccip/directory
	routerAddress: z.string(),
	// Your deployed CCIPSenderReceiver address on the source chain. See the README.
	receiverContractAddress: z.string(),
	// CCIP chain selector of the destination chain, as a string (JSON has no bigint).
	// e.g. "3478487238524512106" for Arbitrum Sepolia.
	destinationChainSelector: z.string(),
	// Recipient address on the destination chain.
	destinationReceiver: z.string(),
	// Plain text payload to send. This example ABI-encodes it as bytes; your destination
	// receiver contract must know how to decode whatever shape you send.
	message: z.string(),
});

type Config = z.infer<typeof configSchema>;

// ---------- Handler ----------

export function onCronTrigger(runtime: Runtime<Config>, _payload: CronPayload): string {
	const config = runtime.config;
	const network = getNetwork({ chainFamily: 'evm', chainSelectorName: config.chainName, isTestnet: true });
	if (!network) throw new Error(`Network not found for chain name: ${config.chainName}`);

	const ccip = new CcipClient({
		chainSelector: network.chainSelector.selector,
		routerAddress: config.routerAddress as `0x${string}`,
		receiverContractAddress: config.receiverContractAddress as `0x${string}`,
	});

	const request = {
		destinationChainSelector: BigInt(config.destinationChainSelector),
		receiver: config.destinationReceiver as `0x${string}`,
		data: toHex(config.message),
	} as const;

	// checkReadiness() runs automatically inside send() and throws if the receiver contract
	// (or a configured tokenSource/feeTokenSource) isn't funded — but funding can still change
	// between that check and report delivery, so send() can also throw or the eventual onchain
	// transaction can still revert. Always handle both cases here, not just the happy path.
	try {
		const result = ccip.send(runtime, request);
		runtime.log(`CCIP message sent: txHash=${result.txHash ?? 'unknown'}`);
		return `sent:${result.txHash ?? 'unknown'}`;
	} catch (error) {
		runtime.log(`CCIP send failed: ${error}`);
		return `failed:${String(error)}`;
	}
}

// ---------- Init ----------

export function initWorkflow(config: Config) {
	const cron = new cre.capabilities.CronCapability();
	return [cre.handler(cron.trigger({ schedule: config.schedule }), onCronTrigger)];
}
