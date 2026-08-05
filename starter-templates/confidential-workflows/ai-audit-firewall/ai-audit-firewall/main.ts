import {
  CronCapability,
  EVMClient,
  HTTPClient,
  NITRO_REGIONS,
  Runner,
  handlerInTee,
  TxStatus,
  bytesToHex,
  getNetwork,
  hexToBase64,
  type TeeRuntime,
  type Workflow,
} from "@chainlink/cre-sdk";
import { encodeAbiParameters, parseAbiParameters } from "viem";

type SecretsConfig = {
  scanner_api_key_id: string;
  primary_llm_api_key_id: string;
  secondary_llm_api_key_id: string;
};

type EvmWriteConfig = {
  chain_selector_name: string;
  consumer_address: string;
  gas_limit: string;
};

export type Config = {
  schedule: string;
  mock_base_url: string;
  scanner_url: string;
  primary_llm_url: string;
  secondary_llm_url: string;
  secrets_ids: SecretsConfig;
  evms?: EvmWriteConfig[];
};

type TransactionProposal = {
  chain_selector: number;
  chain_name: string;
  tx_hash: string;
  from_address: string;
  token_contract_address: string;
  protocol_contract_address: string;
  calldata: string;
  value_wei: string;
  signer: string;
  requested_action: string;
};

type ContractArtifact = {
  address: string;
  contract_name: string;
  verified: boolean;
  abi: string[];
  source_code: string;
  compiler_version: string;
  suspicious_notes: string[];
};

type ScannerCredentialValidation = {
  valid: boolean;
  provider: string;
  scopes: string[];
};

type RiskFlags = {
  obfuscatedTax: boolean;
  privilegeEscalation: boolean;
  externalCallRisk: boolean;
  logicBomb: boolean;
};

type LlmAuditResponse = {
  riskFlags: RiskFlags;
  recommendation: "allow" | "deny" | "review";
  confidence: number;
  reasoning: string;
};

type FirewallVerdict = "ALLOW" | "DENY" | "MANUAL_REVIEW";

type FinalAuditResult = {
  verdict: FirewallVerdict;
  reasoning: string;
  riskFlags: RiskFlags;
  proposal: TransactionProposal;
  tokenContract: ContractArtifact;
  protocolContract: ContractArtifact;
  analyses: {
    primary: LlmAuditResponse;
    secondary: LlmAuditResponse;
  };
  auditLogId: string;
  firewallActionId: string;
  onchainTxHash?: string;
};

const JSON_HEADERS = { "Content-Type": "application/json" };

const asObject = (value: unknown): Record<string, unknown> => {
  if (!value || typeof value !== "object" || Array.isArray(value)) {
    return {};
  }
  return value as Record<string, unknown>;
};

const asString = (value: unknown, fallback = ""): string => {
  return typeof value === "string" ? value : fallback;
};

const asNumber = (value: unknown, fallback = 0): number => {
  if (typeof value === "number" && Number.isFinite(value)) {
    return value;
  }
  if (typeof value === "string" && value.trim() !== "") {
    const parsed = Number(value);
    if (Number.isFinite(parsed)) {
      return parsed;
    }
  }
  return fallback;
};

const asBoolean = (value: unknown, fallback = false): boolean => {
  if (typeof value === "boolean") {
    return value;
  }
  if (typeof value === "string") {
    if (value.toLowerCase() === "true") {
      return true;
    }
    if (value.toLowerCase() === "false") {
      return false;
    }
  }
  return fallback;
};

const asStringArray = (value: unknown): string[] => {
  if (!Array.isArray(value)) {
    return [];
  }
  return value.filter((item): item is string => typeof item === "string" && item.trim() !== "");
};

const decodeBody = (raw: Uint8Array): string => new TextDecoder().decode(raw);

const parseJson = (body: string): Record<string, unknown> => {
  try {
    return asObject(JSON.parse(body));
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    throw new Error(`invalid json response: ${message}; body=${body}`);
  }
};

const getJson = (
  runtime: TeeRuntime<Config>,
  client: HTTPClient,
  url: string,
  headers: Record<string, string>,
): Record<string, unknown> => {
  const response = client
    .sendRequest(runtime, {
      url,
      method: "GET",
      headers,
    })
    .result();

  const raw = decodeBody(response.body);
  if (response.statusCode >= 400) {
    throw new Error(`request failed status=${response.statusCode} body=${raw}`);
  }
  return parseJson(raw);
};

const postJson = (
  runtime: TeeRuntime<Config>,
  client: HTTPClient,
  url: string,
  body: Record<string, unknown>,
  headers: Record<string, string>,
): Record<string, unknown> => {
  const bodyBytes = new TextEncoder().encode(JSON.stringify(body));
  const encodedBody = Buffer.from(bodyBytes).toString("base64");

  const response = client
    .sendRequest(runtime, {
      url,
      method: "POST",
      body: encodedBody,
      headers,
    })
    .result();

  const raw = decodeBody(response.body);
  if (response.statusCode >= 400) {
    throw new Error(`request failed status=${response.statusCode} body=${raw}`);
  }
  return parseJson(raw);
};

const parseTransactionProposal = (payload: unknown): TransactionProposal => {
  const row = asObject(payload);
  return {
    chain_selector: asNumber(row.chain_selector),
    chain_name: asString(row.chain_name, "unknown-chain"),
    tx_hash: asString(row.tx_hash, "0x0"),
    from_address: asString(row.from_address, "0x0000000000000000000000000000000000000000"),
    token_contract_address: asString(row.token_contract_address, "0x0000000000000000000000000000000000000000"),
    protocol_contract_address: asString(row.protocol_contract_address, "0x0000000000000000000000000000000000000000"),
    calldata: asString(row.calldata, "0x"),
    value_wei: asString(row.value_wei, "0"),
    signer: asString(row.signer, "0x0000000000000000000000000000000000000000"),
    requested_action: asString(row.requested_action, "review"),
  };
};

const parseContractArtifact = (payload: unknown): ContractArtifact => {
  const row = asObject(payload);
  return {
    address: asString(row.address, "0x0000000000000000000000000000000000000000"),
    contract_name: asString(row.contract_name, "UnknownContract"),
    verified: asBoolean(row.verified),
    abi: asStringArray(row.abi),
    source_code: asString(row.source_code, ""),
    compiler_version: asString(row.compiler_version, "solc-unknown"),
    suspicious_notes: asStringArray(row.suspicious_notes),
  };
};

const parseScannerCredentialValidation = (payload: unknown): ScannerCredentialValidation => {
  const row = asObject(payload);
  return {
    valid: asBoolean(row.valid),
    provider: asString(row.provider, "unknown-scanner"),
    scopes: asStringArray(row.scopes),
  };
};

const parseRiskFlags = (value: unknown): RiskFlags => {
  const row = asObject(value);
  return {
    obfuscatedTax: asBoolean(row.obfuscatedTax),
    privilegeEscalation: asBoolean(row.privilegeEscalation),
    externalCallRisk: asBoolean(row.externalCallRisk),
    logicBomb: asBoolean(row.logicBomb),
  };
};

const parseAuditResponse = (payload: unknown): LlmAuditResponse => {
  const row = asObject(payload);
  const recommendation = asString(row.recommendation, "review");
  const normalizedRecommendation: LlmAuditResponse["recommendation"] =
    recommendation === "allow" || recommendation === "deny" ? recommendation : "review";

  return {
    riskFlags: parseRiskFlags(row.riskFlags),
    recommendation: normalizedRecommendation,
    confidence: asNumber(row.confidence, 0),
    reasoning: asString(row.reasoning, "No reasoning provided by model"),
  };
};

const mergeFlags = (primary: RiskFlags, secondary: RiskFlags): RiskFlags => ({
  obfuscatedTax: primary.obfuscatedTax || secondary.obfuscatedTax,
  privilegeEscalation: primary.privilegeEscalation || secondary.privilegeEscalation,
  externalCallRisk: primary.externalCallRisk || secondary.externalCallRisk,
  logicBomb: primary.logicBomb || secondary.logicBomb,
});

const hasMaliciousRisk = (flags: RiskFlags): boolean => {
  return flags.obfuscatedTax || flags.privilegeEscalation || flags.externalCallRisk || flags.logicBomb;
};

const verdictToCode = (verdict: FirewallVerdict): bigint => {
  if (verdict === "ALLOW") {
    return 1n;
  }
  if (verdict === "DENY") {
    return 2n;
  }
  return 3n;
};

const riskFlagsToMask = (flags: RiskFlags): bigint => {
  let mask = 0n;
  if (flags.obfuscatedTax) mask |= 1n << 0n;
  if (flags.privilegeEscalation) mask |= 1n << 1n;
  if (flags.externalCallRisk) mask |= 1n << 2n;
  if (flags.logicBomb) mask |= 1n << 3n;
  return mask;
};

const encodeVerdictReport = (result: FinalAuditResult, chainSelector: bigint): string => {
  return encodeAbiParameters(
    parseAbiParameters("uint8 verdictCode, uint8 riskMask, uint64 chainSelector"),
    [Number(verdictToCode(result.verdict)), Number(riskFlagsToMask(result.riskFlags)), chainSelector],
  );
};

const writeVerdictOnChain = async (
  runtime: TeeRuntime<Config>,
  result: FinalAuditResult,
): Promise<string | undefined> => {

  const evmConfig = runtime.config.evms?.[0];

  if (!evmConfig) {
    return undefined;
  }

  const { chain_selector_name, consumer_address, gas_limit } = evmConfig;
  
  if (!chain_selector_name || !consumer_address || !gas_limit) {
    return undefined;
  }

  const network = getNetwork({
    chainFamily: "evm",
    chainSelectorName: chain_selector_name,
  });

  if (!network) {
    throw new Error(`Network not found: ${chain_selector_name}`);
  }

  const evmClient = new EVMClient(network.chainSelector.selector);
  const reportPayload = encodeVerdictReport(result, BigInt(network.chainSelector.selector));
  const donRuntime = runtime.usingTheDons();

  const reportResponse = donRuntime
    .report({
      encodedPayload: hexToBase64(reportPayload),
      encoderName: "evm",
      signingAlgo: "ecdsa",
      hashingAlgo: "keccak256",
    })
    .result();

  const writeResult = evmClient
    .writeReport(donRuntime, {
      receiver: consumer_address,
      report: reportResponse,
      gasConfig: {
        gasLimit: gas_limit,
      },
    })
    .result();

  if (writeResult.txStatus !== TxStatus.SUCCESS) {
    throw new Error(`onchain write failed with status ${writeResult.txStatus}`);
  }

  return bytesToHex(writeResult.txHash || new Uint8Array(32));
};

export const determineVerdict = (primary: LlmAuditResponse, secondary: LlmAuditResponse): FirewallVerdict => {
  const combinedFlags = mergeFlags(primary.riskFlags, secondary.riskFlags);
  if (hasMaliciousRisk(combinedFlags)) {
    return "DENY";
  }

  const reviewRequested = primary.recommendation === "review" || secondary.recommendation === "review";
  const lowConfidence = primary.confidence < 0.7 || secondary.confidence < 0.7;
  if (reviewRequested || lowConfidence || primary.recommendation !== secondary.recommendation) {
    return "MANUAL_REVIEW";
  }

  return "ALLOW";
};

const buildPrimaryPrompt = (proposal: TransactionProposal, tokenContract: ContractArtifact): string => {
  return JSON.stringify(
    {
      objective: "Analyze the token contract and proposed transaction for malicious token behavior.",
      transaction: proposal,
      contract: tokenContract,
      checks: ["obfuscatedTax", "privilegeEscalation", "externalCallRisk", "logicBomb"],
    },
    null,
    2,
  );
};

const buildSecondaryPrompt = (
  proposal: TransactionProposal,
  tokenContract: ContractArtifact,
  protocolContract: ContractArtifact,
  primaryAnalysis: LlmAuditResponse,
): string => {
  return JSON.stringify(
    {
      objective: "Analyze the protocol contract using the token analysis as prior context.",
      transaction: proposal,
      tokenContract,
      protocolContract,
      priorAnalysis: primaryAnalysis,
      checks: ["obfuscatedTax", "privilegeEscalation", "externalCallRisk", "logicBomb"],
    },
    null,
    2,
  );
};

const collectTransactionProposal = (
  runtime: TeeRuntime<Config>,
  client: HTTPClient,
  baseUrl: string,
  scannerApiKey: string,
): TransactionProposal => {
  const response = getJson(runtime, client, `${baseUrl}/transaction-proposal`, {
    ...JSON_HEADERS,
    "x-scanner-api-key": scannerApiKey,
  });
  return parseTransactionProposal(response);
};

const collectContractArtifact = (
  runtime: TeeRuntime<Config>,
  client: HTTPClient,
  scannerUrl: string,
  scannerApiKey: string,
  address: string,
): ContractArtifact => {
  const response = getJson(runtime, client, `${scannerUrl}/contracts/${address}`, {
    ...JSON_HEADERS,
    "x-scanner-api-key": scannerApiKey,
  });
  return parseContractArtifact(response);
};

const validateScannerCredentials = (
  runtime: TeeRuntime<Config>,
  client: HTTPClient,
  scannerUrl: string,
  scannerApiKey: string,
): ScannerCredentialValidation => {
  const response = getJson(runtime, client, `${scannerUrl}/credentials/verify`, {
    ...JSON_HEADERS,
    "x-scanner-api-key": scannerApiKey,
  });

  const validation = parseScannerCredentialValidation(response);
  const hasVerificationScope = validation.scopes.includes("verification:read");
  const hasContractScope = validation.scopes.includes("contracts:read");

  if (!validation.valid || !hasVerificationScope || !hasContractScope) {
    throw new Error(
      `scanner credentials failed validation provider=${validation.provider} valid=${validation.valid} scopes=${validation.scopes.join(",")}`,
    );
  }

  return validation;
};

const requestAuditModel = (
  runtime: TeeRuntime<Config>,
  client: HTTPClient,
  url: string,
  apiKey: string,
  model: string,
  prompt: string,
): LlmAuditResponse => {
  const response = postJson(
    runtime,
    client,
    url,
    {
      model,
      input: [
        {
          role: "system",
          content: "You are an AI smart contract audit engine. Emit strict JSON only.",
        },
        {
          role: "user",
          content: prompt,
        },
      ],
    },
    {
      ...JSON_HEADERS,
      Authorization: `Bearer ${apiKey}`,
    },
  );

  const outputText = asString(response.output_text);
  if (!outputText) {
    throw new Error("LLM response missing output_text");
  }
  return parseAuditResponse(parseJson(outputText));
};

const executeAuditLog = (
  runtime: TeeRuntime<Config>,
  client: HTTPClient,
  baseUrl: string,
  scannerApiKey: string,
  payload: Record<string, unknown>,
): Record<string, unknown> => {
  return postJson(runtime, client, `${baseUrl}/audit-log`, payload, {
    ...JSON_HEADERS,
    "x-scanner-api-key": scannerApiKey,
  });
};

const executeFirewallAction = (
  runtime: TeeRuntime<Config>,
  client: HTTPClient,
  baseUrl: string,
  scannerApiKey: string,
  payload: Record<string, unknown>,
): Record<string, unknown> => {
  return postJson(runtime, client, `${baseUrl}/firewall-action`, payload, {
    ...JSON_HEADERS,
    "x-scanner-api-key": scannerApiKey,
  });
};

export const runAuditFirewall = async (
  runtime: TeeRuntime<Config>,
  client = new HTTPClient(),
): Promise<string> => {
  const { mock_base_url, scanner_url, primary_llm_url, secondary_llm_url, secrets_ids } = runtime.config;

  const scannerApiKey = runtime.getSecret({ id: secrets_ids.scanner_api_key_id }).result().value;
  const primaryLlmApiKey = runtime.getSecret({ id: secrets_ids.primary_llm_api_key_id }).result().value;
  const secondaryLlmApiKey = runtime.getSecret({ id: secrets_ids.secondary_llm_api_key_id }).result().value;

  runtime.log("audit-firewall-getsecret-ok");

  const proposal = collectTransactionProposal(runtime, client, mock_base_url, scannerApiKey);
  const scannerCredentialValidation = validateScannerCredentials(runtime, client, scanner_url, scannerApiKey);
  runtime.log(
    `audit-firewall-scanner-credentials-ok provider=${scannerCredentialValidation.provider} scopes=${scannerCredentialValidation.scopes.join(",")}`,
  );

  const tokenContract = collectContractArtifact(
    runtime,
    client,
    scanner_url,
    scannerApiKey,
    proposal.token_contract_address,
  );
  const protocolContract = collectContractArtifact(
    runtime,
    client,
    scanner_url,
    scannerApiKey,
    proposal.protocol_contract_address,
  );

  if (!tokenContract.verified || !protocolContract.verified) {
    const auditLog = executeAuditLog(runtime, client, mock_base_url, scannerApiKey, {
      proposal,
      tokenContract,
      protocolContract,
      verdict: "DENY",
      reason: "One or more contracts are not verified by the scanner.",
      mode: "verification-failure",
    });

    const firewallAction = executeFirewallAction(runtime, client, mock_base_url, scannerApiKey, {
      verdict: "DENY",
      reason: "Verification failure",
      proposal,
      auditLogId: asString(auditLog.audit_log_id, "unknown"),
    });

    return JSON.stringify({
      verdict: "DENY",
      reason: "One or more contracts are not verified by the scanner.",
      auditLogId: asString(auditLog.audit_log_id, "unknown"),
      firewallActionId: asString(firewallAction.firewall_action_id, "unknown"),
    });
  }

  const primaryAnalysis = requestAuditModel(
    runtime,
    client,
    primary_llm_url,
    primaryLlmApiKey,
    "audit-primary",
    buildPrimaryPrompt(proposal, tokenContract),
  );

  const secondaryAnalysis = requestAuditModel(
    runtime,
    client,
    secondary_llm_url,
    secondaryLlmApiKey,
    "audit-secondary",
    buildSecondaryPrompt(proposal, tokenContract, protocolContract, primaryAnalysis),
  );

  const verdict = determineVerdict(primaryAnalysis, secondaryAnalysis);
  const riskFlags = mergeFlags(primaryAnalysis.riskFlags, secondaryAnalysis.riskFlags);
  const reasoning = [primaryAnalysis.reasoning, secondaryAnalysis.reasoning].join(" | ");

  const auditLog = executeAuditLog(runtime, client, mock_base_url, scannerApiKey, {
    proposal,
    tokenContract,
    protocolContract,
    verdict,
    reasoning,
    riskFlags,
    analyses: {
      primary: primaryAnalysis,
      secondary: secondaryAnalysis,
    },
  });

  const firewallAction = executeFirewallAction(runtime, client, mock_base_url, scannerApiKey, {
    verdict,
    reasoning,
    proposal,
    riskFlags,
    auditLogId: asString(auditLog.audit_log_id, "unknown"),
  });

  const result: FinalAuditResult = {
    verdict,
    reasoning,
    riskFlags,
    proposal,
    tokenContract,
    protocolContract,
    analyses: {
      primary: primaryAnalysis,
      secondary: secondaryAnalysis,
    },
    auditLogId: asString(auditLog.audit_log_id, "unknown"),
    firewallActionId: asString(firewallAction.firewall_action_id, "unknown"),
  };

  runtime.log("audit-firewall-onchain-report-start");
  const onchainTxHash = await writeVerdictOnChain(runtime, result);

  if (onchainTxHash) {
    result.onchainTxHash = onchainTxHash;
    runtime.log(`audit-firewall-onchain tx_hash=${onchainTxHash}`);
  }

  runtime.log(`audit-firewall-complete verdict=${verdict} audit_log_id=${result.auditLogId}`);
  return JSON.stringify(result);
};

export const onCronTrigger = async (runtime: TeeRuntime<Config>): Promise<string> => {
  return runAuditFirewall(runtime);
};

export const initWorkflow = (config: Config): Workflow<Config> => {
  if (
    !config.schedule ||
    !config.mock_base_url ||
    !config.scanner_url ||
    !config.primary_llm_url ||
    !config.secondary_llm_url
  ) {
    throw new Error(
      "config requires schedule, mock_base_url, scanner_url, primary_llm_url, and secondary_llm_url",
    );
  }

  if (
    !config.secrets_ids?.scanner_api_key_id ||
    !config.secrets_ids?.primary_llm_api_key_id ||
    !config.secrets_ids?.secondary_llm_api_key_id
  ) {
    throw new Error("config requires secrets_ids fields");
  }

  const cron = new CronCapability();

  return [
    handlerInTee(
      cron.trigger({ schedule: config.schedule }),
      onCronTrigger,
      [{ tee: "nitro", regions: [NITRO_REGIONS[0]] }],
    ),
  ];
};

export async function main() {
  const runner = await Runner.newRunner<Config>();
  await runner.run(initWorkflow);
}
