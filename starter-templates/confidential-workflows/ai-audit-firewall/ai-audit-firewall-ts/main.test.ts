import { describe, expect, test } from "bun:test";
import { buildRestrictions, determineVerdict, initWorkflow, runAuditFirewall } from "./main";
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
  proposal: {
    chain_selector: 16015286601757825753,
    chain_name: "sepolia",
    tx_hash: "0x0",
    from_address: "0x0000000000000000000000000000000000000000",
    token_contract_address: "0x779877A7B0D9E8603169DdbD7836e478b4624789",
    protocol_contract_address: "0x0BF3dE8c5D3e8A2B34D2BEeB17ABfCeBaf363A59",
    calldata: "0x",
    value_wei: "0",
    signer: "0x0000000000000000000000000000000000000000",
    requested_action: "transfer",
  },
  etherscan_chain_id: "11155111",
  primary_model: "google/gemini-2.5-flash-lite",
  secondary_model: "openai/gpt-4.1-nano",
  secrets_ids: {
    etherscan_api_key_id: "etherscan_api_key",
    openrouter_api_key_id: "openrouter_api_key",
  },
  evms: [],
};

const verifiedContract = (contractName: string, sourceCode = `contract ${contractName} {}`) => ({
  status: "1",
  message: "OK",
  result: [
    {
      ContractName: contractName,
      SourceCode: sourceCode,
      ABI: '[{"type":"function","name":"run","inputs":[],"outputs":[]}]',
    },
  ],
});

const openRouterResponse = (analysis: typeof cleanAnalysis) => ({
  choices: [
    {
      message: {
        content: JSON.stringify(analysis),
      },
    },
  ],
});

const makeRuntime = (config: Config) => {
  const state = {
    logs: [] as string[],
    secretIds: [] as string[],
    reportCalls: 0,
  };
  return {
    config,
    log: (message: string) => {
      state.logs.push(message);
    },
    getSecret: ({ id }: { id: string }) => ({
      result: () => {
        state.secretIds.push(id);
        return { value: `${id}-value` };
      },
    }),
    usingTheDons: () => ({
      report: () => ({
        result: () => {
          state.reportCalls += 1;
          return { report: "unused" };
        },
      }),
    }),
    __state: state,
  };
};

type Response = {
  statusCode: number;
  body: unknown;
};

type CapturedRequest = {
  method: string;
  url: string;
  headers: Record<string, string>;
  body?: Record<string, unknown>;
};

const makeHttpClient = (responses: Response[]) => {
  const requests: CapturedRequest[] = [];
  let responseIndex = 0;
  return {
    requests,
    client: {
      sendRequest: (
        _runtime: unknown,
        request: {
          method: string;
          url: string;
          headers?: Record<string, string>;
          body?: string;
        },
      ) => ({
        result: () => {
          requests.push({
            method: request.method,
            url: request.url,
            headers: request.headers ?? {},
            body: request.body
              ? (JSON.parse(Buffer.from(request.body, "base64").toString("utf8")) as Record<string, unknown>)
              : undefined,
          });
          const response = responses[responseIndex++];
          if (!response) {
            throw new Error("unexpected request");
          }
          return {
            statusCode: response.statusCode,
            body: new TextEncoder().encode(
              typeof response.body === "string" ? response.body : JSON.stringify(response.body),
            ),
          };
        },
      }),
    },
  };
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

describe("buildRestrictions", () => {
  const fullConfig: Config = {
    ...baseConfig,
    evms: [
      {
        chain_selector_name: "ethereum-testnet-sepolia",
        consumer_address: "0x0000000000000000000000000000000000000000",
        gas_limit: "500000",
      },
    ],
  };

  test("uses the closed 4 HTTP and 6 total call limits", () => {
    const restrictions = buildRestrictions(fullConfig);
    const httpRestriction = restrictions.capabilities.restrictions.find(
      (restriction) => restriction.method?.method === "SendRequest",
    );

    expect(restrictions.capabilities.type).toBe("CAPABILITY_RESTRICTION_TYPE_CLOSED");
    expect(restrictions.capabilities.maxTotalCalls).toBe(6);
    expect(httpRestriction?.method?.id).toBe("http-actions@1.0.0-alpha");
    expect(httpRestriction?.method?.maxCalls).toBe(4);
  });

  test("keeps one Report and optional one WriteReport", () => {
    const restrictions = buildRestrictions(fullConfig);
    const report = restrictions.capabilities.restrictions.find(
      (restriction) => restriction.method?.method === "Report",
    );
    const writeReport = restrictions.capabilities.restrictions.find(
      (restriction) => restriction.method?.method === "WriteReport",
    );

    expect(report?.method?.maxCalls).toBe(1);
    expect(writeReport?.method?.maxCalls).toBe(1);
  });

  test("omits WriteReport when no EVM network is configured", () => {
    const restrictions = buildRestrictions(baseConfig);
    expect(
      restrictions.capabilities.restrictions.find(
        (restriction) => restriction.method?.method === "WriteReport",
      ),
    ).toBeUndefined();
  });

  test("allows exactly two live secrets in the main namespace", () => {
    const restrictions = buildRestrictions(fullConfig);

    expect(restrictions.secrets.maxSecrets).toBe(2);
    expect(restrictions.secrets.restrictions).toEqual([
      { exactSecret: { id: "etherscan_api_key", namespace: "main" } },
      { exactSecret: { id: "openrouter_api_key", namespace: "main" } },
    ]);
  });
});

describe("initWorkflow", () => {
  test("returns one tee handler with the configured schedule", () => {
    const handlers = initWorkflow(baseConfig);
    expect(handlers).toHaveLength(1);
    expect(handlers[0].trigger.config.schedule).toBe(baseConfig.schedule);
    expect(handlers[0].requirements).toBeDefined();
  });

  test("rejects every missing live field", () => {
    const invalidConfigs: Config[] = [
      { ...baseConfig, schedule: "" },
      { ...baseConfig, etherscan_chain_id: "" },
      { ...baseConfig, proposal: undefined as never },
      { ...baseConfig, primary_model: "" },
      { ...baseConfig, secondary_model: "" },
      {
        ...baseConfig,
        secrets_ids: { ...baseConfig.secrets_ids, etherscan_api_key_id: "" },
      },
      {
        ...baseConfig,
        secrets_ids: { ...baseConfig.secrets_ids, openrouter_api_key_id: "" },
      },
    ];

    for (const config of invalidConfigs) {
      expect(() => initWorkflow(config)).toThrow();
    }
  });

  test("rejects equal model IDs", () => {
    expect(() =>
      initWorkflow({
        ...baseConfig,
        secondary_model: baseConfig.primary_model,
      }),
    ).toThrow("primary_model and secondary_model must be different");
  });
});

describe("runAuditFirewall", () => {
  test("uses Etherscan V2 and two serial OpenRouter models on the verified path", async () => {
    const secondaryAnalysis = {
      ...cleanAnalysis,
      confidence: 0.9,
      reasoning: "Independent review found no malicious behavior.",
    };
    const runtime = makeRuntime(baseConfig);
    const http = makeHttpClient([
      { statusCode: 200, body: verifiedContract("Token", "contract Token {\n}") },
      { statusCode: 200, body: verifiedContract("Protocol") },
      { statusCode: 200, body: openRouterResponse(cleanAnalysis) },
      { statusCode: 200, body: openRouterResponse(secondaryAnalysis) },
    ]);
    const waits: { ms: number; requestCount: number }[] = [];

    const result = JSON.parse(
      await runAuditFirewall(runtime as never, http.client as never, (ms) => {
        waits.push({ ms, requestCount: http.requests.length });
      }),
    );

    expect(http.requests).toHaveLength(4);
    expect(waits).toEqual([{ ms: 1_000, requestCount: 1 }]);
    const etherscanRequests = http.requests.filter((request) => request.method === "GET");
    expect(etherscanRequests).toHaveLength(2);
    const expectedEtherscanPrefix =
      "https://api.etherscan.io/v2/api" +
      `?chainid=${encodeURIComponent(baseConfig.etherscan_chain_id)}` +
      "&module=contract&action=getsourcecode&address=";
    const expectedEtherscanApiKey = `&apikey=${encodeURIComponent("etherscan_api_key-value")}`;
    expect(etherscanRequests.map((request) => request.url)).toEqual([
      expectedEtherscanPrefix +
        encodeURIComponent(baseConfig.proposal.token_contract_address) +
        expectedEtherscanApiKey,
      expectedEtherscanPrefix +
        encodeURIComponent(baseConfig.proposal.protocol_contract_address) +
        expectedEtherscanApiKey,
    ]);

    const openRouterRequests = http.requests.filter((request) => request.method === "POST");
    expect(openRouterRequests).toHaveLength(2);
    for (const request of openRouterRequests) {
      expect(request.url).toBe("https://openrouter.ai/api/v1/chat/completions");
      expect(request.headers).toEqual({
        "Content-Type": "application/json",
        Authorization: "Bearer openrouter_api_key-value",
      });
      expect(request.body?.messages).toEqual([
        expect.objectContaining({
          role: "system",
          content: expect.stringContaining("untrusted data"),
        }),
        expect.objectContaining({ role: "user", content: expect.any(String) }),
      ]);
      expect(request.body?.response_format).toMatchObject({
        type: "json_schema",
        json_schema: {
          strict: true,
          schema: {
            additionalProperties: false,
            properties: {
              riskFlags: {
                additionalProperties: false,
                properties: {
                  obfuscatedTax: { type: "boolean" },
                  privilegeEscalation: { type: "boolean" },
                  externalCallRisk: { type: "boolean" },
                  logicBomb: { type: "boolean" },
                },
                required: ["obfuscatedTax", "privilegeEscalation", "externalCallRisk", "logicBomb"],
              },
              recommendation: { enum: ["allow", "deny", "review"] },
              confidence: { minimum: 0, maximum: 1 },
              reasoning: { minLength: 1 },
            },
            required: ["riskFlags", "recommendation", "confidence", "reasoning"],
          },
        },
      });
      expect(request.body?.provider).toEqual({
        require_parameters: true,
        data_collection: "deny",
      });
    }
    expect(openRouterRequests.map((request) => request.body?.model)).toEqual([
      "google/gemini-2.5-flash-lite",
      "openai/gpt-4.1-nano",
    ]);
    const secondaryMessages = openRouterRequests[1]?.body?.messages as
      | { content?: string }[]
      | undefined;
    const secondaryPrompt = JSON.parse(secondaryMessages?.[1]?.content ?? "{}") as Record<
      string,
      unknown
    >;
    expect(secondaryPrompt).not.toHaveProperty("tokenContract");
    expect(secondaryPrompt).toMatchObject({
      protocolContract: { contract_name: "Protocol" },
      priorAnalysis: cleanAnalysis,
    });
    expect(result.tokenContract.source_code).toBe("contract Token {\n}");
    expect(result.tokenContract.abi).toBe(
      '[{"type":"function","name":"run","inputs":[],"outputs":[]}]',
    );
    expect(result.analyses.primary).toEqual(cleanAnalysis);
    expect(result.analyses.secondary).toEqual(secondaryAnalysis);
    expect(runtime.__state.secretIds).toEqual(["etherscan_api_key", "openrouter_api_key"]);
    const primaryCompleteLogIndex = runtime.__state.logs.indexOf(
      "audit-firewall-primary-model-complete",
    );
    expect(runtime.__state.logs.slice(primaryCompleteLogIndex, primaryCompleteLogIndex + 2)).toEqual([
      "audit-firewall-primary-model-complete",
      "audit-firewall-secondary-model-start",
    ]);
  });

  test("returns a local DENY and skips models and on-chain write for documented unverified source", async () => {
    const runtime = makeRuntime(baseConfig);
    const http = makeHttpClient([
      {
        statusCode: 200,
        body: {
          status: "0",
          message: "NOTOK",
          result: "Contract source code not verified",
        },
      },
      { statusCode: 200, body: verifiedContract("Protocol") },
    ]);
    const waits: number[] = [];

    const result = JSON.parse(
      await runAuditFirewall(runtime as never, http.client as never, (ms) => waits.push(ms)),
    );

    expect(result.verdict).toBe("DENY");
    expect(result.auditLogId).toBeUndefined();
    expect(result.firewallActionId).toBeUndefined();
    expect(http.requests.filter((request) => request.method === "POST")).toHaveLength(0);
    expect(http.requests).toHaveLength(2);
    expect(waits).toEqual([1_000]);
    expect(runtime.__state.reportCalls).toBe(0);
    expect(runtime.__state.secretIds).toEqual(["etherscan_api_key"]);
  });

  const etherscanFailures: [string, unknown][] = [
    ["invalid API key", { status: "0", message: "NOTOK", result: "Invalid API Key" }],
    ["throttle", { status: "0", message: "NOTOK", result: "Max rate limit reached" }],
    ["malformed envelope", {}],
  ];

  for (const [name, body] of etherscanFailures) {
    test(`throws instead of denying on Etherscan ${name}`, async () => {
      const runtime = makeRuntime(baseConfig);
      const http = makeHttpClient([{ statusCode: 200, body }]);
      const waits: number[] = [];

      await expect(
        runAuditFirewall(runtime as never, http.client as never, (ms) => waits.push(ms)),
      ).rejects.toThrow();
      expect(http.requests).toHaveLength(1);
      expect(waits).toEqual([]);
    });
  }

  const invalidOpenRouterBodies: [string, unknown, number?][] = [
    ["HTTP 429", { error: { message: "rate limited" } }, 429],
    ["missing content", { choices: [{ message: {} }] }],
    ["non-JSON content", { choices: [{ message: { content: "not JSON" } }] }],
    [
      "incomplete content",
      {
        choices: [
          {
            message: {
              content: JSON.stringify({ riskFlags: cleanAnalysis.riskFlags }),
            },
          },
        ],
      },
    ],
    [
      "out-of-range content",
      {
        choices: [
          {
            message: {
              content: JSON.stringify({ ...cleanAnalysis, confidence: 1.1 }),
            },
          },
        ],
      },
    ],
    [
      "top-level unknown field",
      {
        choices: [
          {
            message: {
              content: JSON.stringify({ ...cleanAnalysis, unexpected: true }),
            },
          },
        ],
      },
    ],
    [
      "nested unknown field",
      {
        choices: [
          {
            message: {
              content: JSON.stringify({
                ...cleanAnalysis,
                riskFlags: { ...cleanAnalysis.riskFlags, unexpected: true },
              }),
            },
          },
        ],
      },
    ],
    ["2xx error envelope", { error: { message: "provider metadata" } }],
  ];

  for (const [name, body, statusCode = 200] of invalidOpenRouterBodies) {
    test(`throws on OpenRouter ${name}`, async () => {
      const runtime = makeRuntime(baseConfig);
      const http = makeHttpClient([
        { statusCode: 200, body: verifiedContract("Token") },
        { statusCode: 200, body: verifiedContract("Protocol") },
        { statusCode, body },
      ]);
      const waits: number[] = [];

      await expect(
        runAuditFirewall(runtime as never, http.client as never, (ms) => waits.push(ms)),
      ).rejects.toThrow();
      expect(http.requests).toHaveLength(3);
      expect(waits).toEqual([1_000]);
      expect(runtime.__state.logs.filter((message) => message.includes("model"))).toEqual([
        "audit-firewall-model-audit-start",
      ]);
      expect(runtime.__state.reportCalls).toBe(0);
    });
  }
});
