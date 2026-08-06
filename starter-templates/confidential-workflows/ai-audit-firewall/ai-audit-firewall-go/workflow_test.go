package main

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	evmmock "github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/mock"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	httpmock "github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http/mock"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/testutils"
)

const (
	scannerAPIKey      = "test-scanner-key"
	primaryLLMAPIKey   = "test-primary-key"
	secondaryLLMAPIKey = "test-secondary-key"
	mockBaseURL        = "http://127.0.0.1:8787/audit-firewall"
	scannerURL         = "http://127.0.0.1:8787/audit-firewall/scanner"
	primaryLLMURL      = "http://127.0.0.1:8787/audit-firewall/v1/analysis/primary"
	secondaryLLMURL    = "http://127.0.0.1:8787/audit-firewall/v1/analysis/secondary"
	consumerAddress    = "0x00000000000000000000000000000000000000aa"
	sepoliaChainName   = "ethereum-testnet-sepolia"
)

func testConfig() *Config {
	return &Config{
		Schedule:        "0 */5 * * * *",
		MockBaseURL:     mockBaseURL,
		ScannerURL:      scannerURL,
		PrimaryLLMURL:   primaryLLMURL,
		SecondaryLLMURL: secondaryLLMURL,
		SecretsIDs: SecretsConfig{
			ScannerAPIKeyID:      "scanner_api_key",
			PrimaryLLMAPIKeyID:   "primary_llm_api_key",
			SecondaryLLMAPIKeyID: "secondary_llm_api_key",
		},
	}
}

// configWithEVM opts into the onchain leg, which the base config leaves off.
func configWithEVM() *Config {
	config := testConfig()
	config.EVMs = []EvmWriteConfig{{
		ChainSelectorName: sepoliaChainName,
		ConsumerAddress:   consumerAddress,
		GasLimit:          "500000",
	}}
	return config
}

func testSecrets() testutils.Secrets {
	return testutils.Secrets{
		cre.DefaultSecretNamespace: {
			"scanner_api_key":       scannerAPIKey,
			"primary_llm_api_key":   primaryLLMAPIKey,
			"secondary_llm_api_key": secondaryLLMAPIKey,
		},
	}
}

func auditFixture(recommendation string, confidence float64, flags RiskFlags) LlmAuditResponse {
	return LlmAuditResponse{
		RiskFlags:      flags,
		Recommendation: recommendation,
		Confidence:     confidence,
		Reasoning:      "fixture reasoning",
	}
}

type capturedCall struct {
	URL     string
	Method  string
	Headers map[string][]string
	Body    string
}

// stubOptions controls the fixture responses the enclave HTTP stub returns.
type stubOptions struct {
	tokenVerified    bool
	protocolVerified bool
	scopes           []string
	credentialsValid bool
	primaryAudit     LlmAuditResponse
	secondaryAudit   LlmAuditResponse
}

func defaultStubOptions() stubOptions {
	clean := RiskFlags{}
	return stubOptions{
		tokenVerified:    true,
		protocolVerified: true,
		credentialsValid: true,
		scopes:           []string{"verification:read", "contracts:read"},
		primaryAudit:     auditFixture(RecommendationAllow, 0.95, clean),
		secondaryAudit:   auditFixture(RecommendationAllow, 0.92, clean),
	}
}

// stubEnclaveHTTP routes every endpoint the workflow calls through a single
// in-enclave stub and records each request.
func stubEnclaveHTTP(t *testing.T, opts stubOptions) *[]capturedCall {
	t.Helper()

	capability, err := httpmock.NewClientCapability(t)
	require.NoError(t, err)

	contractJSON := func(address string, verified bool) []byte {
		payload, err := json.Marshal(map[string]any{
			"address":          address,
			"contract_name":    "FixtureContract",
			"verified":         verified,
			"abi":              []string{"function transfer(address,uint256)"},
			"source_code":      "contract Fixture {}",
			"compiler_version": "solc-0.8.24",
			"suspicious_notes": []string{},
		})
		require.NoError(t, err)
		return payload
	}

	auditJSON := func(audit LlmAuditResponse) []byte {
		inner, err := json.Marshal(audit)
		require.NoError(t, err)
		payload, err := json.Marshal(map[string]any{"output_text": string(inner)})
		require.NoError(t, err)
		return payload
	}

	calls := &[]capturedCall{}
	capability.SendRequest = func(_ context.Context, input *http.Request) (*http.Response, error) {
		headers := map[string][]string{}
		for key, value := range input.MultiHeaders {
			headers[key] = value.Values
		}
		*calls = append(*calls, capturedCall{
			URL:     input.Url,
			Method:  input.Method,
			Headers: headers,
			Body:    string(input.Body),
		})

		switch {
		case strings.HasSuffix(input.Url, "/transaction-proposal"):
			return &http.Response{StatusCode: 200, Body: []byte(
				`{"chain_selector":16015286601757825753,"chain_name":"ethereum-testnet-sepolia",` +
					`"tx_hash":"0xabc","from_address":"0x1111111111111111111111111111111111111111",` +
					`"token_contract_address":"0x2222222222222222222222222222222222222222",` +
					`"protocol_contract_address":"0x3333333333333333333333333333333333333333",` +
					`"calldata":"0xdeadbeef","value_wei":"0","signer":"0x1111111111111111111111111111111111111111",` +
					`"requested_action":"review"}`)}, nil

		case strings.HasSuffix(input.Url, "/credentials/verify"):
			payload, err := json.Marshal(map[string]any{
				"valid":    opts.credentialsValid,
				"provider": "fixture-scanner",
				"scopes":   opts.scopes,
			})
			require.NoError(t, err)
			return &http.Response{StatusCode: 200, Body: payload}, nil

		case strings.Contains(input.Url, "/contracts/0x2222"):
			return &http.Response{StatusCode: 200, Body: contractJSON(
				"0x2222222222222222222222222222222222222222", opts.tokenVerified)}, nil

		case strings.Contains(input.Url, "/contracts/0x3333"):
			return &http.Response{StatusCode: 200, Body: contractJSON(
				"0x3333333333333333333333333333333333333333", opts.protocolVerified)}, nil

		case strings.HasSuffix(input.Url, "/analysis/primary"):
			return &http.Response{StatusCode: 200, Body: auditJSON(opts.primaryAudit)}, nil

		case strings.HasSuffix(input.Url, "/analysis/secondary"):
			return &http.Response{StatusCode: 200, Body: auditJSON(opts.secondaryAudit)}, nil

		case strings.HasSuffix(input.Url, "/audit-log"):
			return &http.Response{StatusCode: 200, Body: []byte(`{"audit_log_id":"log-1"}`)}, nil

		case strings.HasSuffix(input.Url, "/firewall-action"):
			return &http.Response{StatusCode: 200, Body: []byte(`{"firewall_action_id":"act-1"}`)}, nil
		}
		return nil, assert.AnError
	}
	return calls
}

// stubEVMWrite registers the EVM capability and records the write request.
func stubEVMWrite(t *testing.T, status evm.TxStatus) **evm.WriteReportRequest {
	t.Helper()

	chainSelector, err := evm.ChainSelectorFromName(sepoliaChainName)
	require.NoError(t, err)

	capability, err := evmmock.NewClientCapability(chainSelector, t)
	require.NoError(t, err)

	captured := new(*evm.WriteReportRequest)
	capability.WriteReport = func(_ context.Context, input *evm.WriteReportRequest) (*evm.WriteReportReply, error) {
		*captured = input
		txHash := make([]byte, 32)
		txHash[31] = 0x7b
		return &evm.WriteReportReply{TxStatus: status, TxHash: txHash}, nil
	}
	return captured
}

func runFirewall(t *testing.T, config *Config, runtime cre.TeeRuntime) (FinalAuditResult, error) {
	t.Helper()

	raw, err := RunAuditFirewall(config, runtime, &http.Client{})
	if err != nil {
		return FinalAuditResult{}, err
	}

	var result FinalAuditResult
	require.NoError(t, json.Unmarshal([]byte(raw), &result))
	return result, nil
}

// ─── DetermineVerdict ───────────────────────────────────────

func TestDetermineVerdict_AllowsOnConfidentUnanimousAgreement(t *testing.T) {
	verdict := DetermineVerdict(
		auditFixture(RecommendationAllow, 0.9, RiskFlags{}),
		auditFixture(RecommendationAllow, 0.8, RiskFlags{}),
	)
	assert.Equal(t, VerdictAllow, verdict)
}

// Any raised flag denies outright, regardless of how confident the models were
// or whether they recommended allow.
func TestDetermineVerdict_DeniesWhenAnyRiskFlagIsRaised(t *testing.T) {
	for name, flags := range map[string]RiskFlags{
		"obfuscatedTax":       {ObfuscatedTax: true},
		"privilegeEscalation": {PrivilegeEscalation: true},
		"externalCallRisk":    {ExternalCallRisk: true},
		"logicBomb":           {LogicBomb: true},
	} {
		t.Run(name, func(t *testing.T) {
			verdict := DetermineVerdict(
				auditFixture(RecommendationAllow, 1.0, flags),
				auditFixture(RecommendationAllow, 1.0, RiskFlags{}),
			)
			assert.Equal(t, VerdictDeny, verdict)
		})
	}
}

func TestDetermineVerdict_EscalatesOnDisagreement(t *testing.T) {
	verdict := DetermineVerdict(
		auditFixture(RecommendationAllow, 0.9, RiskFlags{}),
		auditFixture(RecommendationDeny, 0.9, RiskFlags{}),
	)
	assert.Equal(t, VerdictManualReview, verdict)
}

func TestDetermineVerdict_EscalatesOnLowConfidence(t *testing.T) {
	verdict := DetermineVerdict(
		auditFixture(RecommendationAllow, 0.69, RiskFlags{}),
		auditFixture(RecommendationAllow, 0.99, RiskFlags{}),
	)
	assert.Equal(t, VerdictManualReview, verdict)
}

func TestDetermineVerdict_EscalatesWhenEitherModelRequestsReview(t *testing.T) {
	verdict := DetermineVerdict(
		auditFixture(RecommendationReview, 0.99, RiskFlags{}),
		auditFixture(RecommendationAllow, 0.99, RiskFlags{}),
	)
	assert.Equal(t, VerdictManualReview, verdict)
}

// ⚠️ Documents a quirk inherited from the TypeScript template: DetermineVerdict
// only branches on *agreement* and confidence, never on what was agreed. Two
// models that confidently and unanimously recommend DENY, having raised no risk
// flag, therefore fall through to ALLOW.
//
// This test pins the current behaviour so the port stays faithful. If you decide
// this is wrong, the fix belongs in both templates: add an explicit
// `if primary.Recommendation == RecommendationDeny { return VerdictDeny }` before
// the agreement check.
func TestDetermineVerdict_UnanimousDenyWithoutFlagsFallsThroughToAllow(t *testing.T) {
	verdict := DetermineVerdict(
		auditFixture(RecommendationDeny, 0.99, RiskFlags{}),
		auditFixture(RecommendationDeny, 0.99, RiskFlags{}),
	)
	assert.Equal(t, VerdictAllow, verdict,
		"agreement is checked but the agreed-upon recommendation is not")
}

func TestMergeFlags_UnionsFindings(t *testing.T) {
	merged := MergeFlags(
		RiskFlags{ObfuscatedTax: true, ExternalCallRisk: false},
		RiskFlags{ExternalCallRisk: true},
	)
	assert.Equal(t, RiskFlags{ObfuscatedTax: true, ExternalCallRisk: true}, merged)
}

// ─── Report encoding ────────────────────────────────────────

func TestRiskFlagsToMask_PacksEachFlagIntoItsOwnBit(t *testing.T) {
	assert.Equal(t, uint8(0), RiskFlagsToMask(RiskFlags{}))
	assert.Equal(t, uint8(1), RiskFlagsToMask(RiskFlags{ObfuscatedTax: true}))
	assert.Equal(t, uint8(2), RiskFlagsToMask(RiskFlags{PrivilegeEscalation: true}))
	assert.Equal(t, uint8(4), RiskFlagsToMask(RiskFlags{ExternalCallRisk: true}))
	assert.Equal(t, uint8(8), RiskFlagsToMask(RiskFlags{LogicBomb: true}))
	assert.Equal(t, uint8(15), RiskFlagsToMask(RiskFlags{
		ObfuscatedTax: true, PrivilegeEscalation: true, ExternalCallRisk: true, LogicBomb: true,
	}))
}

func TestEncodeVerdictReport_ProducesThreeAbiWords(t *testing.T) {
	encoded, err := EncodeVerdictReport(FinalAuditResult{
		Verdict:   VerdictDeny,
		RiskFlags: RiskFlags{LogicBomb: true},
	}, 16015286601757825753)
	require.NoError(t, err)

	// uint8, uint8, uint64 each occupy a full 32-byte ABI word.
	require.Len(t, encoded, 96)
	assert.Equal(t, byte(2), encoded[31], "verdictCode DENY == 2")
	assert.Equal(t, byte(8), encoded[63], "riskMask logicBomb == 8")
}

func TestEncodeVerdictReport_MapsEachVerdictToItsCode(t *testing.T) {
	for verdict, code := range map[string]byte{
		VerdictAllow:        1,
		VerdictDeny:         2,
		VerdictManualReview: 3,
	} {
		encoded, err := EncodeVerdictReport(FinalAuditResult{Verdict: verdict}, 1)
		require.NoError(t, err)
		assert.Equal(t, code, encoded[31], "verdict %s", verdict)
	}
}

// ─── Untrusted model output ─────────────────────────────────

// An unparseable or unexpected recommendation must never become an allow.
func TestParseAuditResponse_NormalisesUnknownRecommendationToReview(t *testing.T) {
	parsed := parseAuditResponse(map[string]any{"recommendation": "definitely-fine"})
	assert.Equal(t, RecommendationReview, parsed.Recommendation)
}

func TestParseAuditResponse_CoercesStringBooleansAndNumbers(t *testing.T) {
	parsed := parseAuditResponse(map[string]any{
		"recommendation": RecommendationAllow,
		"confidence":     "0.85",
		"riskFlags":      map[string]any{"logicBomb": "true"},
	})

	assert.Equal(t, 0.85, parsed.Confidence)
	assert.True(t, parsed.RiskFlags.LogicBomb)
}

func TestParseAuditResponse_DefaultsMissingFields(t *testing.T) {
	parsed := parseAuditResponse(map[string]any{})

	assert.Equal(t, RecommendationReview, parsed.Recommendation)
	assert.Equal(t, 0.0, parsed.Confidence)
	assert.Equal(t, "No reasoning provided by model", parsed.Reasoning)
	assert.Equal(t, RiskFlags{}, parsed.RiskFlags)
}

// ─── End-to-end through the enclave runtime ─────────────────

func TestRunAuditFirewall_AllowsCleanVerifiedTransaction(t *testing.T) {
	stubEnclaveHTTP(t, defaultStubOptions())
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	result, err := runFirewall(t, testConfig(), runtime)
	require.NoError(t, err)

	assert.Equal(t, VerdictAllow, result.Verdict)
	assert.Equal(t, "log-1", result.AuditLogID)
	assert.Equal(t, "act-1", result.FirewallActionID)
	assert.Empty(t, result.OnchainTxHash, "no EVM configured means no onchain write")
}

func TestRunAuditFirewall_DeniesWhenModelsRaiseFlags(t *testing.T) {
	opts := defaultStubOptions()
	opts.secondaryAudit = auditFixture(RecommendationAllow, 0.99, RiskFlags{LogicBomb: true})
	stubEnclaveHTTP(t, opts)
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	result, err := runFirewall(t, testConfig(), runtime)
	require.NoError(t, err)

	assert.Equal(t, VerdictDeny, result.Verdict)
	assert.True(t, result.RiskFlags.LogicBomb)
}

// An unverified contract has no source to audit, so the pipeline fails closed
// without consulting a model at all.
func TestRunAuditFirewall_ShortCircuitsToDenyOnUnverifiedContract(t *testing.T) {
	opts := defaultStubOptions()
	opts.protocolVerified = false
	calls := stubEnclaveHTTP(t, opts)
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	raw, err := RunAuditFirewall(testConfig(), runtime, &http.Client{})
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal([]byte(raw), &decoded))
	assert.Equal(t, VerdictDeny, decoded["verdict"])
	assert.Contains(t, decoded["reason"], "not verified")

	for _, call := range *calls {
		assert.NotContains(t, call.URL, "/analysis/", "must not consult a model on unverified source")
	}
}

func TestRunAuditFirewall_RejectsScannerCredentialsMissingScopes(t *testing.T) {
	opts := defaultStubOptions()
	opts.scopes = []string{"contracts:read"}
	stubEnclaveHTTP(t, opts)
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err := RunAuditFirewall(testConfig(), runtime, &http.Client{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "scanner credentials failed validation")
}

func TestRunAuditFirewall_RejectsInvalidScannerCredentials(t *testing.T) {
	opts := defaultStubOptions()
	opts.credentialsValid = false
	stubEnclaveHTTP(t, opts)
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err := RunAuditFirewall(testConfig(), runtime, &http.Client{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "scanner credentials failed validation")
}

func TestRunAuditFirewall_InjectsEnclaveFetchedSecretsIntoRequests(t *testing.T) {
	calls := stubEnclaveHTTP(t, defaultStubOptions())
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err := RunAuditFirewall(testConfig(), runtime, &http.Client{})
	require.NoError(t, err)

	sawScannerKey := false
	sawPrimaryKey := false
	sawSecondaryKey := false
	for _, call := range *calls {
		if values, ok := call.Headers["x-scanner-api-key"]; ok {
			assert.Equal(t, []string{scannerAPIKey}, values)
			sawScannerKey = true
		}
		for _, auth := range call.Headers["Authorization"] {
			if auth == "Bearer "+primaryLLMAPIKey {
				sawPrimaryKey = true
			}
			if auth == "Bearer "+secondaryLLMAPIKey {
				sawSecondaryKey = true
			}
		}
	}

	assert.True(t, sawScannerKey, "scanner key must be injected inside the enclave")
	assert.True(t, sawPrimaryKey, "primary model key must be injected inside the enclave")
	assert.True(t, sawSecondaryKey, "secondary model key must be injected inside the enclave")
}

// The second model must receive the first model's findings as prior context.
func TestRunAuditFirewall_ChainsPrimaryAnalysisIntoSecondaryPrompt(t *testing.T) {
	calls := stubEnclaveHTTP(t, defaultStubOptions())
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err := RunAuditFirewall(testConfig(), runtime, &http.Client{})
	require.NoError(t, err)

	var secondaryBody string
	for _, call := range *calls {
		if strings.HasSuffix(call.URL, "/analysis/secondary") {
			secondaryBody = call.Body
		}
	}

	require.NotEmpty(t, secondaryBody)
	assert.Contains(t, secondaryBody, "priorAnalysis")
	assert.Contains(t, secondaryBody, "protocolContract")
}

func TestRunAuditFirewall_ErrorsOnNon2xx(t *testing.T) {
	capability, err := httpmock.NewClientCapability(t)
	require.NoError(t, err)
	capability.SendRequest = func(_ context.Context, _ *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 500, Body: []byte("boom")}, nil
	}

	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err = RunAuditFirewall(testConfig(), runtime, &http.Client{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "status=500")
}

// Logging inside an enclave weakens its confidentiality guarantee. The workflow
// logs only non-sensitive markers, never credentials or contract source.
func TestRunAuditFirewall_DoesNotLogSecretsOrContractSource(t *testing.T) {
	stubEnclaveHTTP(t, defaultStubOptions())
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err := RunAuditFirewall(testConfig(), runtime, &http.Client{})
	require.NoError(t, err)

	for _, raw := range runtime.GetLogs() {
		line := string(raw)
		assert.NotContains(t, line, scannerAPIKey, "log leaked the scanner key")
		assert.NotContains(t, line, primaryLLMAPIKey, "log leaked the primary model key")
		assert.NotContains(t, line, secondaryLLMAPIKey, "log leaked the secondary model key")
		assert.NotContains(t, line, "contract Fixture", "log leaked contract source")
	}
}

// ─── Onchain delivery ───────────────────────────────────────

func TestRunAuditFirewall_WritesVerdictOnChainWhenConfigured(t *testing.T) {
	stubEnclaveHTTP(t, defaultStubOptions())
	captured := stubEVMWrite(t, evm.TxStatus_TX_STATUS_SUCCESS)
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	result, err := runFirewall(t, configWithEVM(), runtime)
	require.NoError(t, err)

	assert.NotEmpty(t, result.OnchainTxHash)
	assert.True(t, strings.HasPrefix(result.OnchainTxHash, "0x"))

	require.NotNil(t, *captured)
	assert.Equal(t, uint64(500000), (*captured).GasConfig.GetGasLimit())
	assert.Len(t, (*captured).Receiver, 20, "receiver is a 20-byte address")
}

func TestRunAuditFirewall_FailsWhenOnchainWriteIsNotSuccessful(t *testing.T) {
	stubEnclaveHTTP(t, defaultStubOptions())
	stubEVMWrite(t, evm.TxStatus_TX_STATUS_REVERTED)
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err := RunAuditFirewall(configWithEVM(), runtime, &http.Client{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "onchain write failed with status")
}

func TestWriteVerdictOnChain_SkipsWhenEvmConfigIncomplete(t *testing.T) {
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	config := configWithEVM()
	config.EVMs[0].ConsumerAddress = ""

	txHash, err := writeVerdictOnChain(runtime, config, FinalAuditResult{Verdict: VerdictAllow})
	require.NoError(t, err)
	assert.Empty(t, txHash)
}

func TestWriteVerdictOnChain_ErrorsOnUnknownChain(t *testing.T) {
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	config := configWithEVM()
	config.EVMs[0].ChainSelectorName = "not-a-real-chain"

	_, err := writeVerdictOnChain(runtime, config, FinalAuditResult{Verdict: VerdictAllow})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "network not found")
}

func TestWriteVerdictOnChain_ErrorsOnUnparseableGasLimit(t *testing.T) {
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	config := configWithEVM()
	config.EVMs[0].GasLimit = "plenty"

	_, err := writeVerdictOnChain(runtime, config, FinalAuditResult{Verdict: VerdictAllow})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid gas_limit")
}

// ─── Restrictions (pre-hook) ────────────────────────────────

// findMethodRestriction locates a restriction by method name, or nil.
func findMethodRestriction(restrictions *sdkpb.Restrictions, method string) *sdkpb.MethodRestriction {
	for _, restriction := range restrictions.GetCapabilities().GetRestrictions() {
		if m := restriction.GetMethod(); m != nil && m.GetMethod() == method {
			return m
		}
	}
	return nil
}

func findMethodRestrictionByID(restrictions *sdkpb.Restrictions, id string) *sdkpb.MethodRestriction {
	for _, restriction := range restrictions.GetCapabilities().GetRestrictions() {
		if m := restriction.GetMethod(); m != nil && m.GetId() == id {
			return m
		}
	}
	return nil
}

func TestBuildRestrictions_UsesClosedSetWithExactCallBudget(t *testing.T) {
	restrictions, err := BuildRestrictions(configWithEVM())
	require.NoError(t, err)

	assert.Equal(t,
		sdkpb.CapabilityRestrictionType_CAPABILITY_RESTRICTION_TYPE_CLOSED,
		restrictions.GetCapabilities().GetType())
	assert.Equal(t, uint32(10), restrictions.GetCapabilities().GetMaxTotalCalls())
}

func TestBuildRestrictions_LimitsHTTPSendRequestTo8Calls(t *testing.T) {
	restrictions, err := BuildRestrictions(configWithEVM())
	require.NoError(t, err)

	httpRestriction := findMethodRestriction(restrictions, "SendRequest")
	require.NotNil(t, httpRestriction)
	assert.Equal(t, "http-actions@1.0.0-alpha", httpRestriction.GetId())
	assert.Equal(t, uint32(8), httpRestriction.GetMaxCalls())
}

func TestBuildRestrictions_LimitsConsensusReportTo1Call(t *testing.T) {
	restrictions, err := BuildRestrictions(configWithEVM())
	require.NoError(t, err)

	consensusRestriction := findMethodRestrictionByID(restrictions, consensusCapabilityID)
	require.NotNil(t, consensusRestriction)
	assert.Equal(t, "Report", consensusRestriction.GetMethod())
	assert.Equal(t, uint32(1), consensusRestriction.GetMaxCalls())
}

func TestBuildRestrictions_LimitsEvmWriteReportTo1Call(t *testing.T) {
	restrictions, err := BuildRestrictions(configWithEVM())
	require.NoError(t, err)

	evmRestriction := findMethodRestriction(restrictions, "WriteReport")
	require.NotNil(t, evmRestriction)
	assert.Equal(t, uint32(1), evmRestriction.GetMaxCalls())

	chainSelector, err := evm.ChainSelectorFromName(sepoliaChainName)
	require.NoError(t, err)
	assert.Contains(t, evmRestriction.GetId(), strconv.FormatUint(chainSelector, 10),
		"the restriction must be scoped to the configured chain")
}

// The onchain leg is opt-in, so an unconfigured EVM target must not be granted
// WriteReport at all.
func TestBuildRestrictions_OmitsEvmRestrictionWhenNoEvmsConfigured(t *testing.T) {
	restrictions, err := BuildRestrictions(testConfig())
	require.NoError(t, err)

	assert.Nil(t, findMethodRestriction(restrictions, "WriteReport"))
}

func TestBuildRestrictions_SkipsEvmRestrictionOnUnknownChain(t *testing.T) {
	config := configWithEVM()
	config.EVMs[0].ChainSelectorName = "not-a-real-chain"

	restrictions, err := BuildRestrictions(config)
	require.NoError(t, err)

	assert.Nil(t, findMethodRestriction(restrictions, "WriteReport"),
		"an unresolvable chain is skipped here; the write path reports it at execution time")
}

func TestBuildRestrictions_AllowsExactlyThreeExactMatchSecrets(t *testing.T) {
	restrictions, err := BuildRestrictions(configWithEVM())
	require.NoError(t, err)

	assert.Equal(t, uint32(3), restrictions.GetSecrets().GetMaxSecrets())
	require.Len(t, restrictions.GetSecrets().GetRestrictions(), 3)

	ids := []string{}
	for _, restriction := range restrictions.GetSecrets().GetRestrictions() {
		exact := restriction.GetExactSecret()
		require.NotNil(t, exact, "restrictions must be exact-match, not prefixed")
		ids = append(ids, exact.GetId())
	}

	assert.Contains(t, ids, "scanner_api_key")
	assert.Contains(t, ids, "primary_llm_api_key")
	assert.Contains(t, ids, "secondary_llm_api_key")
}

func TestBuildRestrictions_UsesMainNamespaceForAllSecrets(t *testing.T) {
	restrictions, err := BuildRestrictions(configWithEVM())
	require.NoError(t, err)

	for _, restriction := range restrictions.GetSecrets().GetRestrictions() {
		assert.Equal(t, "main", restriction.GetExactSecret().GetNamespace())
	}
}

// The HTTP budget must cover the steady-state call count, otherwise a clean run
// would be cut off part-way through by its own restrictions.
func TestBuildRestrictions_HTTPBudgetCoversASuccessfulRun(t *testing.T) {
	calls := stubEnclaveHTTP(t, defaultStubOptions())
	stubEVMWrite(t, evm.TxStatus_TX_STATUS_SUCCESS)
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err := RunAuditFirewall(configWithEVM(), runtime, &http.Client{})
	require.NoError(t, err)

	restrictions, err := BuildRestrictions(configWithEVM())
	require.NoError(t, err)

	budget := findMethodRestriction(restrictions, "SendRequest").GetMaxCalls()
	assert.LessOrEqual(t, uint32(len(*calls)), budget,
		"a successful run made %d HTTP calls but the budget is %d", len(*calls), budget)
}

// ─── Config parsing ─────────────────────────────────────────

// An empty payload is what the DON hands the runner when it invokes the pre-hook
// before a real config is attached, so it must not fail to unmarshal.
func TestParseConfig_FallsBackToDefaultsOnEmptyPayload(t *testing.T) {
	for name, payload := range map[string][]byte{
		"nil":        nil,
		"empty":      {},
		"whitespace": []byte("  \n\t "),
	} {
		t.Run(name, func(t *testing.T) {
			config, err := ParseConfig(payload)
			require.NoError(t, err)
			assert.Equal(t, DefaultConfig.Schedule, config.Schedule)
			assert.Equal(t, "scanner_api_key", config.SecretsIDs.ScannerAPIKeyID)
		})
	}
}

// The defaults must be complete enough to build restrictions from, since that is
// the whole reason they exist.
func TestParseConfig_DefaultsAreSufficientForThePreHook(t *testing.T) {
	config, err := ParseConfig(nil)
	require.NoError(t, err)

	require.NoError(t, func() error { _, err := InitWorkflow(config, nil, nil); return err }())

	restrictions, err := BuildRestrictions(config)
	require.NoError(t, err)
	assert.Len(t, restrictions.GetSecrets().GetRestrictions(), 3)
	assert.NotNil(t, findMethodRestriction(restrictions, "WriteReport"))
}

func TestParseConfig_ParsesRealPayload(t *testing.T) {
	encoded, err := json.Marshal(testConfig())
	require.NoError(t, err)

	config, err := ParseConfig(encoded)
	require.NoError(t, err)
	assert.Equal(t, mockBaseURL, config.MockBaseURL)
	assert.Empty(t, config.EVMs, "must not inherit the default EVM target")
}

func TestParseConfig_ErrorsOnMalformedPayload(t *testing.T) {
	_, err := ParseConfig([]byte("{not json"))
	require.Error(t, err)
}

// Mutating a parsed default must not corrupt DefaultConfig for later calls.
func TestParseConfig_ReturnsIndependentCopyOfDefaults(t *testing.T) {
	first, err := ParseConfig(nil)
	require.NoError(t, err)
	first.Schedule = "mutated"

	second, err := ParseConfig(nil)
	require.NoError(t, err)
	assert.Equal(t, "0 */5 * * * *", second.Schedule)
}

// ─── InitWorkflow ───────────────────────────────────────────

func TestInitWorkflow_RegistersOneTeeHandler(t *testing.T) {
	workflow, err := InitWorkflow(testConfig(), nil, nil)
	require.NoError(t, err)
	assert.Len(t, workflow, 1)
}

func TestInitWorkflow_ErrorsWhenCoreFieldsMissing(t *testing.T) {
	config := testConfig()
	config.SecondaryLLMURL = ""

	_, err := InitWorkflow(config, nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(),
		"config requires schedule, mock_base_url, scanner_url, primary_llm_url, and secondary_llm_url")
}

func TestInitWorkflow_ErrorsWhenSecretIDsMissing(t *testing.T) {
	config := testConfig()
	config.SecretsIDs.SecondaryLLMAPIKeyID = ""

	_, err := InitWorkflow(config, nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "config requires secrets_ids fields")
}
