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
  keccak256,
  padHex,
  toBytes,
  type Hex,
  type Address,
  type AbiFunction,
} from "viem";
import { z } from "zod";
import { IAutomationCompatibleABI } from "../contracts/evm/ts/generated/IAutomationCompatible";
import { AutomationReceiver } from "../contracts/evm/ts/generated/AutomationReceiver";
import { IAutomationCompatible } from "../contracts/evm/ts/generated/IAutomationCompatible";

// ─── Config Schema ──────────────────────────────────────────
const configSchema = z.object({
  chainSelectorName: z.string(),
  receiverAddress: z.string(),
  targetAddress: z.string(),
  migrationType: z.enum(["CRON", "CUSTOM", "LOG"]),
  // Shared config
  checkData: z.string().optional(),
  writeGasLimit: z.string().optional(),
  
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

const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000";

type AutomationLog = {
  index: bigint;
  timestamp: bigint;
  txHash: Hex;
  blockNumber: bigint;
  blockHash: Hex;
  source: Address;
  topics: readonly Hex[];
  data: Hex;
};

function assertConfiguredAddress(value: string | undefined, fieldName: string): asserts value is Address {
  if (!value || value.toLowerCase() === ZERO_ADDRESS) {
    throw new Error(`${fieldName} must be configured with a deployed contract address`);
  }
}

// ─── Utility: Map CRE Log to Automation Log Struct ──────────
function mapLogToAutomation(log: EVMLog): AutomationLog {
  return {
    index: BigInt(log.index || 0),
    // CRE's EVMLog does not carry the block timestamp; legacy checkLog logic that
    // depends on log.timestamp must source it elsewhere. See README "Log Mapping".
    timestamp: 0n,
    txHash: (log.txHash ? bytesToHex(log.txHash) : "0x" + "0".repeat(64)) as Hex,
    blockNumber: log.blockNumber ? bytesToBigint(log.blockNumber.absVal) : 0n,
    blockHash: (log.blockHash ? bytesToHex(log.blockHash) : "0x" + "0".repeat(64)) as Hex,
    source: (log.address ? bytesToHex(log.address) : "0x" + "0".repeat(40)) as Address,
    topics: (log.topics || []).map(t => bytesToHex(t) as Hex),
    data: log.data ? bytesToHex(log.data) as Hex : "0x" as Hex,
  };
}

function encodeTopicFilter(value?: string): string[] {
  if (!value) return [];
  return [hexToBase64(padHex(value as Hex, { size: 32 }))];
}

function buildTopicsFilter(config: Config): Array<{ values: string[] }> {
  if (!config.logTriggerEventSignature) {
    throw new Error("logTriggerEventSignature is required for LOG migration");
  }

  const topics: Array<{ values: string[] }> = [
    { values: [hexToBase64(keccak256(toBytes(config.logTriggerEventSignature)))] },
  ];

  if (config.topic1 || config.topic2 || config.topic3) {
    topics.push({ values: encodeTopicFilter(config.topic1) });
  }

  if (config.topic2 || config.topic3) {
    topics.push({ values: encodeTopicFilter(config.topic2) });
  }

  if (config.topic3) {
    topics.push({ values: encodeTopicFilter(config.topic3) });
  }

  return topics;
}

// ─── Handler ────────────────────────────────────────────────
const runMigration = (runtime: Runtime<Config>, triggerLog?: EVMLog): string => {
  const config = runtime.config;
  const checkData = (config.checkData ?? "0x") as Hex;
  const writeGasLimit = config.writeGasLimit ?? "500000";
  assertConfiguredAddress(config.receiverAddress, "receiverAddress");
  assertConfiguredAddress(config.targetAddress, "targetAddress");
  runtime.log(`=== Migration Workflow Started: ${config.migrationType} ===`);

  const network = getNetwork({ chainFamily: "evm", chainSelectorName: config.chainSelectorName });
  if (!network) throw new Error(`Network not found: ${config.chainSelectorName}`);

  const evmClient = new EVMClient(network.chainSelector.selector);
  const target = new IAutomationCompatible(evmClient, config.targetAddress);
  const receiver = new AutomationReceiver(evmClient, config.receiverAddress);

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
    const [needed, data] = target.checkUpkeep(runtime, checkData);
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
    const [needed, data] = target.checkLog(runtime, mapLogToAutomation(triggerLog), checkData);
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
    gasLimit: writeGasLimit
  });

  if (writeResult.txStatus !== TxStatus.SUCCESS) {
    throw new Error(`Transaction failed: ${writeResult.errorMessage || writeResult.txStatus}`);
  }

  const txHash = bytesToHex(writeResult.txHash || new Uint8Array(32));
  runtime.log(`Success! TX: ${txHash}`);
  return txHash;
};

const onCronTrigger = (runtime: Runtime<Config>): string => {
  return runMigration(runtime);
};

const onLogTrigger = (runtime: Runtime<Config>, triggerLog: EVMLog): string => {
  return runMigration(runtime, triggerLog);
};

// ─── Workflow Init ──────────────────────────────────────────
function initWorkflow(config: Config) {
  const network = getNetwork({ chainFamily: "evm", chainSelectorName: config.chainSelectorName });
  if (!network) throw new Error(`Network not found: ${config.chainSelectorName}`);
  const evmClient = new EVMClient(network.chainSelector.selector);

  if (config.migrationType === "LOG") {
    assertConfiguredAddress(config.logTriggerAddress, "logTriggerAddress");

    return [
      cre.handler(
        evmClient.logTrigger({
          addresses: [hexToBase64(config.logTriggerAddress)],
          topics: buildTopicsFilter(config),
          confidence: "CONFIDENCE_LEVEL_FINALIZED",
        }),
        onLogTrigger,
      ),
    ];
  } else {
    // CRON or CUSTOM use the scheduler
    if (!config.schedule) throw new Error("schedule is required for CRON/CUSTOM migration");
    return [
      cre.handler(
        new cre.capabilities.CronCapability().trigger({ schedule: config.schedule }),
        onCronTrigger,
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
