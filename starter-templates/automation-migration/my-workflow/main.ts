import {
  cre,
  EVMClient,
  getNetwork,
  hexToBase64,
  bytesToHex,
  bytesToBigint,
  TxStatus,
  type Runtime,
  Runner,
  type EVMLog,
} from "@chainlink/cre-sdk";
import {
  encodeFunctionData,
  encodeAbiParameters,
  parseAbiItem,
  type Hex,
  type Address,
  type AbiFunction,
} from "viem";
import { z } from "zod";
import { IAutomationCompatibleABI } from "../contracts/evm/ts/generated/IAutomationCompatible";
import { AutomationReceiver } from "../contracts/evm/ts/generated/AutomationReceiver";
import { IAutomationCompatible } from "../contracts/evm/ts/generated/IAutomationCompatible";

// ─── Config Schema ──────────────────────────────────────────
export const configSchema = z.object({
  chainSelectorName: z.string(),
  receiverAddress: z.string(),
  targetAddress: z.string(),
  migrationType: z.enum(["CRON", "CUSTOM", "LOG"]),
  // Shared config
  checkData: z.string().optional().default("0x"),
  writeGasLimit: z.string().optional().default("500000"),
  
  // Cron/Custom specific
  schedule: z.string().optional(),
  
  // Cron specific (Time-based migration)
  targetFunction: z.string().optional(),
  targetInputs: z.array(z.unknown()).optional(),
  
  // Log specific
  logTriggerAddress: z.string().optional(),
  logTriggerEventSignature: z.string().optional(),
  topic1: z.string().optional(),
  topic2: z.string().optional(),
  topic3: z.string().optional(),
});
type Config = z.infer<typeof configSchema>;

// ─── Utility: Map CRE Log to Automation Log Struct ──────────
function mapLogToAutomation(log: EVMLog): any {
  return {
    index: BigInt(log.index || 0),
    timestamp: 0n, 
    txHash: (log.txHash ? bytesToHex(log.txHash) : "0x" + "0".repeat(64)) as Hex,
    blockNumber: log.blockNumber ? bytesToBigint(log.blockNumber.absVal) : 0n,
    blockHash: (log.blockHash ? bytesToHex(log.blockHash) : "0x" + "0".repeat(64)) as Hex,
    source: (log.address ? bytesToHex(log.address) : "0x" + "0".repeat(40)) as Address,
    topics: (log.topics || []).map(t => bytesToHex(t) as Hex),
    data: log.data ? bytesToHex(log.data) as Hex : "0x" as Hex,
  };
}

// ─── Handler ────────────────────────────────────────────────
const onTrigger = (runtime: Runtime<Config>, triggerLog?: EVMLog): string => {
  const config = runtime.config;
  runtime.log(`=== Migration Workflow Started: ${config.migrationType} ===`);

  const network = getNetwork({ chainFamily: "evm", chainSelectorName: config.chainSelectorName });
  if (!network) throw new Error(`Network not found: ${config.chainSelectorName}`);

  const evmClient = new EVMClient(network.chainSelector.selector);
  const target = new IAutomationCompatible(evmClient, config.targetAddress as Address);
  const receiver = new AutomationReceiver(evmClient, config.receiverAddress as Address);

  let performData: Hex = "0x";
  let upkeepNeeded = false;
  let finalCallData: Hex = "0x";

  if (config.migrationType === "CRON") {
    // Time-based migration: Always true (it triggered)
    upkeepNeeded = true;
    
    // Encode the target function call
    if (!config.targetFunction) throw new Error("targetFunction is required for CRON migration");
    const abiItem = parseAbiItem(`function ${config.targetFunction}`) as AbiFunction;
    finalCallData = encodeFunctionData({
      abi: [abiItem],
      functionName: abiItem.name,
      args: (config.targetInputs || []) as readonly unknown[],
    });
    
  } else if (config.migrationType === "CUSTOM") {
    runtime.log("Checking custom logic upkeep...");
    const [needed, data] = target.checkUpkeep(runtime, config.checkData as Hex);
    upkeepNeeded = needed;
    performData = data;
    
    finalCallData = encodeFunctionData({
      abi: IAutomationCompatibleABI,
      functionName: 'performUpkeep',
      args: [performData],
    });
    
  } else if (config.migrationType === "LOG") {
    if (!triggerLog) throw new Error("Log data missing for LOG migration");
    runtime.log("Checking log trigger upkeep...");
    const [needed, data] = target.checkLog(runtime, mapLogToAutomation(triggerLog), config.checkData as Hex);
    upkeepNeeded = needed;
    performData = data;
    
    finalCallData = encodeFunctionData({
      abi: IAutomationCompatibleABI,
      functionName: 'performUpkeep',
      args: [performData],
    });
  }

  if (!upkeepNeeded) {
    runtime.log("Upkeep not needed. Skipping.");
    return "No upkeep needed";
  }

  runtime.log(`Upkeep needed! Target: ${config.targetAddress}`);

  // Encode the payload for the generic AutomationReceiver
  const reportPayload = encodeAbiParameters(
    [{ name: "target", type: "address" }, { name: "data", type: "bytes" }],
    [config.targetAddress as Address, finalCallData]
  );

  // Write the report to the Receiver
  runtime.log("Writing report to AutomationReceiver...");
  const writeResult = receiver.writeReport(runtime, reportPayload, {
    gasLimit: config.writeGasLimit
  });

  if (writeResult.txStatus !== TxStatus.SUCCESS) {
    throw new Error(`Transaction failed: ${writeResult.errorMessage || writeResult.txStatus}`);
  }

  const txHash = bytesToHex(writeResult.txHash || new Uint8Array(32));
  runtime.log(`Success! TX: ${txHash}`);
  return txHash;
};

// ─── Workflow Init ──────────────────────────────────────────
export function initWorkflow(config: Config) {
  const network = getNetwork({ chainFamily: "evm", chainSelectorName: config.chainSelectorName });
  if (!network) throw new Error(`Network not found: ${config.chainSelectorName}`);
  const evmClient = new EVMClient(network.chainSelector.selector);

  if (config.migrationType === "LOG") {
    // Build topics filter if needed (Topic 0 is always the event hash)
    // For simplicity in this template, we assume Topic 0 is provided via event signature
    return [
      cre.handler(
        evmClient.logTrigger({
          addresses: [hexToBase64(config.logTriggerAddress as Address)],
          // In a real scenario, you'd build the topics array here
          // This template uses the basic trigger and lets the user refine in code
        }),
        onTrigger,
      ),
    ];
  } else {
    // CRON or CUSTOM use the scheduler
    if (!config.schedule) throw new Error("schedule is required for CRON/CUSTOM migration");
    return [
      cre.handler(
        new cre.capabilities.CronCapability().trigger({ schedule: config.schedule }),
        onTrigger,
      ),
    ];
  }
}

/**
 * Entry point
 */
export async function main() {
  const runner = await Runner.newRunner<Config>({ configSchema });
  await runner.run(initWorkflow);
}

main();
