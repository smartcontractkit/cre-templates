import { describe, expect, test } from "bun:test";
import { determineVerdict, initWorkflow, runAuditFirewall } from "./main";
import type { Config } from "./main";

const cleanAnalysis = {
  riskFlags: {
    obfuscatedTax: false,
    privilegeEscalation: false,
    externalCallRisk: false,
    logicBomb: false,
  },
  recommendation: "allow" as const,
  confidence: 0.95,
  reasoning: "No malicious patterns detected.",
};

const baseConfig: Config = {
  schedule: "0 */5 * * * *",
  mock_base_url: "http://127.0.0.1:8787/audit-firewall",
  scanner_url: "http://127.0.0.1:8787/audit-firewall/scanner",
  primary_llm_url: "http://127.0.0.1:8787/audit-firewall/v1/analysis/primary",
  secondary_llm_url: "http://127.0.0.1:8787/audit-firewall/v1/analysis/secondary",
  secrets_ids: {
    scanner_api_key_id: "scanner_api_key",
    primary_llm_api_key_id: "primary_llm_api_key",
    secondary_llm_api_key_id: "secondary_llm_api_key",
  },
};

describe("determineVerdict", () => {
  test("denies when either model reports malicious risk", () => {
    const verdict = determineVerdict(
      {
        ...cleanAnalysis,
        riskFlags: {
          ...cleanAnalysis.riskFlags,
          logicBomb: true,
        },
      },
      cleanAnalysis,
    );

    expect(verdict).toBe("DENY");
  });

  test("routes uncertain outcomes to manual review", () => {
    const verdict = determineVerdict(
      cleanAnalysis,
      {
        ...cleanAnalysis,
        recommendation: "review",
      },
    );

    expect(verdict).toBe("MANUAL_REVIEW");
  });

  test("allows clean, high-confidence agreements", () => {
    expect(determineVerdict(cleanAnalysis, cleanAnalysis)).toBe("ALLOW");
  });
});

describe("initWorkflow", () => {
  test("returns one tee handler with configured schedule", () => {
    const handlers = initWorkflow(baseConfig);
    expect(handlers).toHaveLength(1);
    expect(handlers[0].trigger.config.schedule).toBe(baseConfig.schedule);
    expect(handlers[0].requirements).toBeDefined();
  });

  test("throws when core fields are missing", () => {
    expect(() =>
      initWorkflow({
        ...baseConfig,
        primary_llm_url: "",
      }),
    ).toThrow("config requires schedule, mock_base_url, scanner_url, primary_llm_url, and secondary_llm_url");
  });
});

describe("runAuditFirewall", () => {
  const makeRuntime = (config: Config) => {
    const logs: string[] = [];
    return {
      config,
      log: (message: string) => {
        logs.push(message);
      },
      getSecret: ({ id }: { id: string }) => ({
        result: () => ({ value: `${id}-value` }),
      }),
      usingTheDons: () => ({
        report: () => ({
          result: () => ({ report: "unused" }),
        }),
      }),
      __logs: logs,
    };
  };

  const makeHttpClient = (handlers: Record<string, { statusCode: number; body: Record<string, unknown> }>) => ({
    sendRequest: (_runtime: unknown, request: { url: string; method: string }) => ({
      result: () => {
        const response = handlers[`${request.method} ${request.url}`];
        if (!response) {
          throw new Error(`unexpected request ${request.method} ${request.url}`);
        }

        return {
          statusCode: response.statusCode,
          body: new TextEncoder().encode(JSON.stringify(response.body)),
        };
      },
    }),
  });

  test("fails fast when scanner credentials do not include verification access", async () => {
    const runtime = makeRuntime(baseConfig);
    const httpClient = makeHttpClient({
      "GET http://127.0.0.1:8787/audit-firewall/transaction-proposal": {
        statusCode: 200,
        body: {
          chain_selector: 1,
          chain_name: "sepolia",
          tx_hash: "0x1",
          from_address: "0x1",
          token_contract_address: "0x2",
          protocol_contract_address: "0x3",
          calldata: "0x",
          value_wei: "0",
          signer: "0x4",
          requested_action: "transfer",
        },
      },
      "GET http://127.0.0.1:8787/audit-firewall/scanner/credentials/verify": {
        statusCode: 200,
        body: {
          valid: true,
          provider: "mock-scanner",
          scopes: ["contracts:read"],
        },
      },
    });

    await expect(runAuditFirewall(runtime as never, httpClient as never)).rejects.toThrow(
      "scanner credentials failed validation",
    );
  });

  test("denies immediately when the scanner marks a contract as unverified", async () => {
    const runtime = makeRuntime({
      ...baseConfig,
      evms: [],
    });
    const httpClient = makeHttpClient({
      "GET http://127.0.0.1:8787/audit-firewall/transaction-proposal": {
        statusCode: 200,
        body: {
          chain_selector: 1,
          chain_name: "sepolia",
          tx_hash: "0x1",
          from_address: "0x1",
          token_contract_address: "0x2",
          protocol_contract_address: "0x3",
          calldata: "0x",
          value_wei: "0",
          signer: "0x4",
          requested_action: "transfer",
        },
      },
      "GET http://127.0.0.1:8787/audit-firewall/scanner/credentials/verify": {
        statusCode: 200,
        body: {
          valid: true,
          provider: "mock-scanner",
          scopes: ["contracts:read", "verification:read"],
        },
      },
      "GET http://127.0.0.1:8787/audit-firewall/scanner/contracts/0x2": {
        statusCode: 200,
        body: {
          address: "0x2",
          contract_name: "Token",
          verified: false,
          abi: ["function transfer(address,uint256) external returns (bool)"],
          source_code: "contract Token {}",
          compiler_version: "solc-0.8.26",
          suspicious_notes: [],
        },
      },
      "GET http://127.0.0.1:8787/audit-firewall/scanner/contracts/0x3": {
        statusCode: 200,
        body: {
          address: "0x3",
          contract_name: "Protocol",
          verified: true,
          abi: ["function execute() external"],
          source_code: "contract Protocol {}",
          compiler_version: "solc-0.8.26",
          suspicious_notes: [],
        },
      },
      "POST http://127.0.0.1:8787/audit-firewall/audit-log": {
        statusCode: 200,
        body: {
          audit_log_id: "audit_1",
        },
      },
      "POST http://127.0.0.1:8787/audit-firewall/firewall-action": {
        statusCode: 200,
        body: {
          firewall_action_id: "firewall_1",
        },
      },
    });

    const result = JSON.parse(await runAuditFirewall(runtime as never, httpClient as never));

    expect(result.verdict).toBe("DENY");
    expect(result.auditLogId).toBe("audit_1");
    expect(result.firewallActionId).toBe("firewall_1");
  });
});
