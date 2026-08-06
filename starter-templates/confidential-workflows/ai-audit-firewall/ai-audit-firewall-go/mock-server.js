import express from "express";

const app = express();
app.use(express.json());

const PORT = Number(process.env.MOCK_PORT ?? 8787);
const SCANNER_KEY = process.env.MOCK_SCANNER_API_KEY ?? "mock-scanner-key";
const PRIMARY_LLM_KEY = process.env.MOCK_PRIMARY_LLM_API_KEY ?? "mock-primary-llm-key";
const SECONDARY_LLM_KEY = process.env.MOCK_SECONDARY_LLM_API_KEY ?? "mock-secondary-llm-key";

const auditFirewallTransaction = {
  chain_selector: 16015286601757825753,
  chain_name: "ethereum-testnet-sepolia",
  tx_hash: "0x9b4b3f6fd2d9a1e7d5f0c1c3d2ef0af2b9f0ab7f7a9d4e1ef7a8d2e3c4b5a6f7",
  from_address: "0x1111111111111111111111111111111111111111",
  token_contract_address: "0x2222222222222222222222222222222222222222",
  protocol_contract_address: "0x3333333333333333333333333333333333333333",
  calldata: "0xa9059cbb000000000000000000000000444444444444444444444444444444444444444400000000000000000000000000000000000000000000000000000000000003e8",
  value_wei: "0",
  signer: "0x5555555555555555555555555555555555555555",
  requested_action: "transfer",
};

const auditFirewallContracts = {
  "0x2222222222222222222222222222222222222222": {
    address: "0x2222222222222222222222222222222222222222",
    contract_name: "MockERC20Token",
    verified: true,
    abi: [
      "function transfer(address recipient, uint256 amount) external returns (bool)",
      "function balanceOf(address account) external view returns (uint256)",
    ],
    source_code: "contract MockERC20Token { function transfer(address,uint256) external returns (bool) { return true; } }",
    compiler_version: "solc-0.8.26",
    suspicious_notes: [],
  },
  "0x3333333333333333333333333333333333333333": {
    address: "0x3333333333333333333333333333333333333333",
    contract_name: "MockProtocolRouter",
    verified: true,
    abi: [
      "function execute(address token, address recipient, uint256 amount) external",
      "function owner() external view returns (address)",
    ],
    source_code: "contract MockProtocolRouter { function execute(address,address,uint256) external {} }",
    compiler_version: "solc-0.8.26",
    suspicious_notes: ["performs external calls to token and router targets"],
  },
};

const requireHeader = (req, res, expected, key) => {
  const got = req.header(key);
  if (got !== expected) {
    res.status(401).json({ error: `invalid ${key}` });
    return false;
  }
  return true;
};

app.get("/audit-firewall/transaction-proposal", (req, res) => {
  if (!requireHeader(req, res, SCANNER_KEY, "x-scanner-api-key")) {
    return;
  }

  res.json(auditFirewallTransaction);
});

app.get("/audit-firewall/scanner/contracts/:address", (req, res) => {
  if (!requireHeader(req, res, SCANNER_KEY, "x-scanner-api-key")) {
    return;
  }

  const artifact = auditFirewallContracts[req.params.address];
  if (!artifact) {
    res.status(404).json({ error: "contract not found" });
    return;
  }

  res.json(artifact);
});

app.get("/audit-firewall/scanner/credentials/verify", (req, res) => {
  if (!requireHeader(req, res, SCANNER_KEY, "x-scanner-api-key")) {
    return;
  }

  res.json({
    valid: true,
    provider: "mock-scanner",
    scopes: ["contracts:read", "verification:read"],
  });
});

const parseAuditPrompt = (body) => {
  const input = Array.isArray(body?.input) ? body.input : [];
  const joined = input
    .map((item) => (typeof item?.content === "string" ? item.content : ""))
    .join("\n");
  try {
    return JSON.parse(joined);
  } catch {
    return {};
  }
};

const auditDecision = (prompt) => {
  const serialized = JSON.stringify(prompt);
  const malicious = serialized.includes("logicBomb") || serialized.includes("privilegeEscalation");
  const suspicious = serialized.includes("performs external calls") || serialized.includes("delegatecall");

  return {
    riskFlags: {
      obfuscatedTax: serialized.includes("tax"),
      privilegeEscalation: malicious,
      externalCallRisk: suspicious,
      logicBomb: serialized.includes("timelock") || serialized.includes("upgrade"),
    },
    recommendation: malicious ? "deny" : suspicious ? "review" : "allow",
    confidence: malicious ? 0.98 : suspicious ? 0.74 : 0.93,
    reasoning: malicious
      ? "The audit detected privilege escalation or logic-bomb behavior in the contract path."
      : suspicious
        ? "The contract performs external calls and should be manually reviewed before approval."
        : "The contract and transaction look consistent with the declared behavior.",
  };
};

app.post("/audit-firewall/v1/analysis/primary", (req, res) => {
  const auth = req.header("authorization") ?? "";
  if (auth !== `Bearer ${PRIMARY_LLM_KEY}`) {
    res.status(401).json({ error: "invalid authorization" });
    return;
  }

  const prompt = parseAuditPrompt(req.body);
  const decision = auditDecision(prompt);

  res.json({
    id: "audit_primary_1",
    output_text: JSON.stringify(decision),
  });
});

app.post("/audit-firewall/v1/analysis/secondary", (req, res) => {
  const auth = req.header("authorization") ?? "";
  if (auth !== `Bearer ${SECONDARY_LLM_KEY}`) {
    res.status(401).json({ error: "invalid authorization" });
    return;
  }

  const prompt = parseAuditPrompt(req.body);
  const decision = auditDecision(prompt);

  res.json({
    id: "audit_secondary_1",
    output_text: JSON.stringify(decision),
  });
});

app.post("/audit-firewall/audit-log", (req, res) => {
  if (!requireHeader(req, res, SCANNER_KEY, "x-scanner-api-key")) {
    return;
  }

  res.json({
    status: "ok",
    audit_log_id: `audit_${Date.now()}`,
  });
});

app.post("/audit-firewall/firewall-action", (req, res) => {
  if (!requireHeader(req, res, SCANNER_KEY, "x-scanner-api-key")) {
    return;
  }

  res.json({
    status: "ok",
    firewall_action_id: `firewall_${Date.now()}`,
  });
});

app.listen(PORT, () => {
  console.log(`Mock server running at http://127.0.0.1:${PORT}`);
});
