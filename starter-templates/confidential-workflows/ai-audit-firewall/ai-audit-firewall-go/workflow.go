package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

// ─── Config ─────────────────────────────────────────────────

type SecretsConfig struct {
	ScannerAPIKeyID      string `json:"scanner_api_key_id"`
	PrimaryLLMAPIKeyID   string `json:"primary_llm_api_key_id"`
	SecondaryLLMAPIKeyID string `json:"secondary_llm_api_key_id"`
}

type EvmWriteConfig struct {
	ChainSelectorName string `json:"chain_selector_name"`
	ConsumerAddress   string `json:"consumer_address"`
	GasLimit          string `json:"gas_limit"`
}

type Config struct {
	Schedule        string           `json:"schedule"`
	MockBaseURL     string           `json:"mock_base_url"`
	ScannerURL      string           `json:"scanner_url"`
	PrimaryLLMURL   string           `json:"primary_llm_url"`
	SecondaryLLMURL string           `json:"secondary_llm_url"`
	SecretsIDs      SecretsConfig    `json:"secrets_ids"`
	EVMs            []EvmWriteConfig `json:"evms,omitempty"`
}

// ─── Domain types ───────────────────────────────────────────

type TransactionProposal struct {
	ChainSelector           float64 `json:"chain_selector"`
	ChainName               string  `json:"chain_name"`
	TxHash                  string  `json:"tx_hash"`
	FromAddress             string  `json:"from_address"`
	TokenContractAddress    string  `json:"token_contract_address"`
	ProtocolContractAddress string  `json:"protocol_contract_address"`
	Calldata                string  `json:"calldata"`
	ValueWei                string  `json:"value_wei"`
	Signer                  string  `json:"signer"`
	RequestedAction         string  `json:"requested_action"`
}

type ContractArtifact struct {
	Address         string   `json:"address"`
	ContractName    string   `json:"contract_name"`
	Verified        bool     `json:"verified"`
	ABI             []string `json:"abi"`
	SourceCode      string   `json:"source_code"`
	CompilerVersion string   `json:"compiler_version"`
	SuspiciousNotes []string `json:"suspicious_notes"`
}

type ScannerCredentialValidation struct {
	Valid    bool     `json:"valid"`
	Provider string   `json:"provider"`
	Scopes   []string `json:"scopes"`
}

type RiskFlags struct {
	ObfuscatedTax       bool `json:"obfuscatedTax"`
	PrivilegeEscalation bool `json:"privilegeEscalation"`
	ExternalCallRisk    bool `json:"externalCallRisk"`
	LogicBomb           bool `json:"logicBomb"`
}

// Model recommendations.
const (
	RecommendationAllow  = "allow"
	RecommendationDeny   = "deny"
	RecommendationReview = "review"
)

type LlmAuditResponse struct {
	RiskFlags      RiskFlags `json:"riskFlags"`
	Recommendation string    `json:"recommendation"`
	Confidence     float64   `json:"confidence"`
	Reasoning      string    `json:"reasoning"`
}

// Firewall verdicts.
const (
	VerdictAllow        = "ALLOW"
	VerdictDeny         = "DENY"
	VerdictManualReview = "MANUAL_REVIEW"
)

// minConfidence is the bar both models must clear independently; below it the
// transaction is escalated to a human rather than auto-allowed.
const minConfidence = 0.7

type AuditAnalyses struct {
	Primary   LlmAuditResponse `json:"primary"`
	Secondary LlmAuditResponse `json:"secondary"`
}

type FinalAuditResult struct {
	Verdict          string              `json:"verdict"`
	Reasoning        string              `json:"reasoning"`
	RiskFlags        RiskFlags           `json:"riskFlags"`
	Proposal         TransactionProposal `json:"proposal"`
	TokenContract    ContractArtifact    `json:"tokenContract"`
	ProtocolContract ContractArtifact    `json:"protocolContract"`
	Analyses         AuditAnalyses       `json:"analyses"`
	AuditLogID       string              `json:"auditLogId"`
	FirewallActionID string              `json:"firewallActionId"`
	OnchainTxHash    string              `json:"onchainTxHash,omitempty"`
}

// ─── Lenient JSON coercion ──────────────────────────────────
// Both the scanner and the models are untrusted sources: they may emit numbers
// as strings, booleans as strings, omit keys, or return the wrong shape. These
// helpers coerce rather than fail.

func asObject(value any) map[string]any {
	obj, ok := value.(map[string]any)
	if !ok {
		return map[string]any{}
	}
	return obj
}

func asString(value any, fallback string) string {
	if s, ok := value.(string); ok {
		return s
	}
	return fallback
}

func asNumber(value any, fallback float64) float64 {
	switch v := value.(type) {
	case float64:
		if !math.IsInf(v, 0) && !math.IsNaN(v) {
			return v
		}
	case string:
		if strings.TrimSpace(v) != "" {
			if parsed, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil &&
				!math.IsInf(parsed, 0) && !math.IsNaN(parsed) {
				return parsed
			}
		}
	}
	return fallback
}

func asBoolean(value any, fallback bool) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		switch strings.ToLower(v) {
		case "true":
			return true
		case "false":
			return false
		}
	}
	return fallback
}

func asStringSlice(value any) []string {
	items, ok := value.([]any)
	if !ok {
		return []string{}
	}
	out := []string{}
	for _, item := range items {
		if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
			out = append(out, s)
		}
	}
	return out
}

func parseJSONObject(body string) (map[string]any, error) {
	var parsed any
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		return nil, fmt.Errorf("invalid json response: %w; body=%s", err, body)
	}
	return asObject(parsed), nil
}

const zeroAddress = "0x0000000000000000000000000000000000000000"

func parseTransactionProposal(payload any) TransactionProposal {
	row := asObject(payload)
	return TransactionProposal{
		ChainSelector:           asNumber(row["chain_selector"], 0),
		ChainName:               asString(row["chain_name"], "unknown-chain"),
		TxHash:                  asString(row["tx_hash"], "0x0"),
		FromAddress:             asString(row["from_address"], zeroAddress),
		TokenContractAddress:    asString(row["token_contract_address"], zeroAddress),
		ProtocolContractAddress: asString(row["protocol_contract_address"], zeroAddress),
		Calldata:                asString(row["calldata"], "0x"),
		ValueWei:                asString(row["value_wei"], "0"),
		Signer:                  asString(row["signer"], zeroAddress),
		RequestedAction:         asString(row["requested_action"], RecommendationReview),
	}
}

func parseContractArtifact(payload any) ContractArtifact {
	row := asObject(payload)
	return ContractArtifact{
		Address:         asString(row["address"], zeroAddress),
		ContractName:    asString(row["contract_name"], "UnknownContract"),
		Verified:        asBoolean(row["verified"], false),
		ABI:             asStringSlice(row["abi"]),
		SourceCode:      asString(row["source_code"], ""),
		CompilerVersion: asString(row["compiler_version"], "solc-unknown"),
		SuspiciousNotes: asStringSlice(row["suspicious_notes"]),
	}
}

func parseScannerCredentialValidation(payload any) ScannerCredentialValidation {
	row := asObject(payload)
	return ScannerCredentialValidation{
		Valid:    asBoolean(row["valid"], false),
		Provider: asString(row["provider"], "unknown-scanner"),
		Scopes:   asStringSlice(row["scopes"]),
	}
}

func parseRiskFlags(value any) RiskFlags {
	row := asObject(value)
	return RiskFlags{
		ObfuscatedTax:       asBoolean(row["obfuscatedTax"], false),
		PrivilegeEscalation: asBoolean(row["privilegeEscalation"], false),
		ExternalCallRisk:    asBoolean(row["externalCallRisk"], false),
		LogicBomb:           asBoolean(row["logicBomb"], false),
	}
}

// parseAuditResponse normalises anything that is not an explicit allow/deny to
// "review" — an unparseable recommendation must never become an allow.
func parseAuditResponse(payload any) LlmAuditResponse {
	row := asObject(payload)

	recommendation := asString(row["recommendation"], RecommendationReview)
	if recommendation != RecommendationAllow && recommendation != RecommendationDeny {
		recommendation = RecommendationReview
	}

	return LlmAuditResponse{
		RiskFlags:      parseRiskFlags(row["riskFlags"]),
		Recommendation: recommendation,
		Confidence:     asNumber(row["confidence"], 0),
		Reasoning:      asString(row["reasoning"], "No reasoning provided by model"),
	}
}

// ─── Verdict logic ──────────────────────────────────────────

// MergeFlags unions the two models' findings: either model raising a flag is
// enough to raise it overall.
func MergeFlags(primary, secondary RiskFlags) RiskFlags {
	return RiskFlags{
		ObfuscatedTax:       primary.ObfuscatedTax || secondary.ObfuscatedTax,
		PrivilegeEscalation: primary.PrivilegeEscalation || secondary.PrivilegeEscalation,
		ExternalCallRisk:    primary.ExternalCallRisk || secondary.ExternalCallRisk,
		LogicBomb:           primary.LogicBomb || secondary.LogicBomb,
	}
}

func hasMaliciousRisk(flags RiskFlags) bool {
	return flags.ObfuscatedTax || flags.PrivilegeEscalation || flags.ExternalCallRisk || flags.LogicBomb
}

// DetermineVerdict fails safe: ALLOW requires both models to agree, both to be
// confident, and neither to have raised a risk flag. Anything short of unanimous
// confident agreement escalates rather than allowing.
func DetermineVerdict(primary, secondary LlmAuditResponse) string {
	if hasMaliciousRisk(MergeFlags(primary.RiskFlags, secondary.RiskFlags)) {
		return VerdictDeny
	}

	reviewRequested := primary.Recommendation == RecommendationReview ||
		secondary.Recommendation == RecommendationReview
	lowConfidence := primary.Confidence < minConfidence || secondary.Confidence < minConfidence

	if reviewRequested || lowConfidence || primary.Recommendation != secondary.Recommendation {
		return VerdictManualReview
	}

	return VerdictAllow
}

func verdictToCode(verdict string) uint8 {
	switch verdict {
	case VerdictAllow:
		return 1
	case VerdictDeny:
		return 2
	default:
		return 3
	}
}

// RiskFlagsToMask packs the four flags into a bitmask so the consumer contract
// can inspect individual findings without decoding a struct.
func RiskFlagsToMask(flags RiskFlags) uint8 {
	var mask uint8
	if flags.ObfuscatedTax {
		mask |= 1 << 0
	}
	if flags.PrivilegeEscalation {
		mask |= 1 << 1
	}
	if flags.ExternalCallRisk {
		mask |= 1 << 2
	}
	if flags.LogicBomb {
		mask |= 1 << 3
	}
	return mask
}

// EncodeVerdictReport ABI-encodes (uint8 verdictCode, uint8 riskMask,
// uint64 chainSelector) so a consumer contract can decode it with
// abi.decode(report, (uint8, uint8, uint64)).
//
// Note only the verdict and the flag mask cross to the DON — never the contract
// source, the model reasoning or the scanner credentials.
func EncodeVerdictReport(result FinalAuditResult, chainSelector uint64) ([]byte, error) {
	uint8Type, err := abi.NewType("uint8", "", nil)
	if err != nil {
		return nil, err
	}
	uint64Type, err := abi.NewType("uint64", "", nil)
	if err != nil {
		return nil, err
	}

	args := abi.Arguments{
		{Name: "verdictCode", Type: uint8Type},
		{Name: "riskMask", Type: uint8Type},
		{Name: "chainSelector", Type: uint64Type},
	}
	return args.Pack(verdictToCode(result.Verdict), RiskFlagsToMask(result.RiskFlags), chainSelector)
}

// ─── Confidential HTTP helpers ──────────────────────────────
// Both helpers go through SendRequestInTee, so URLs, headers (including the
// scanner and model API keys) and response bodies stay inside the enclave.

func getJSON(
	runtime cre.TeeRuntime,
	client *http.Client,
	url string,
	headers map[string]string,
) (map[string]any, error) {
	response, err := client.SendRequestInTee(runtime, &http.Request{
		Url:          url,
		Method:       "GET",
		MultiHeaders: multiHeaders(headers),
	}).Await()
	if err != nil {
		return nil, fmt.Errorf("confidential request failed: %w", err)
	}

	raw := string(response.Body)
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed status=%d body=%s", response.StatusCode, raw)
	}
	return parseJSONObject(raw)
}

func postJSON(
	runtime cre.TeeRuntime,
	client *http.Client,
	url string,
	body any,
	headers map[string]string,
) (map[string]any, error) {
	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %w", err)
	}

	response, err := client.SendRequestInTee(runtime, &http.Request{
		Url:          url,
		Method:       "POST",
		Body:         encoded,
		MultiHeaders: multiHeaders(headers),
	}).Await()
	if err != nil {
		return nil, fmt.Errorf("confidential request failed: %w", err)
	}

	raw := string(response.Body)
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed status=%d body=%s", response.StatusCode, raw)
	}
	return parseJSONObject(raw)
}

func multiHeaders(headers map[string]string) map[string]*http.HeaderValues {
	out := make(map[string]*http.HeaderValues, len(headers))
	for key, value := range headers {
		out[key] = &http.HeaderValues{Values: []string{value}}
	}
	return out
}

func scannerHeaders(scannerAPIKey string) map[string]string {
	return map[string]string{
		"Content-Type":      "application/json",
		"x-scanner-api-key": scannerAPIKey,
	}
}

// ─── Prompts ────────────────────────────────────────────────

var auditChecks = []string{"obfuscatedTax", "privilegeEscalation", "externalCallRisk", "logicBomb"}

func buildPrimaryPrompt(proposal TransactionProposal, tokenContract ContractArtifact) (string, error) {
	prompt := struct {
		Objective   string              `json:"objective"`
		Transaction TransactionProposal `json:"transaction"`
		Contract    ContractArtifact    `json:"contract"`
		Checks      []string            `json:"checks"`
	}{
		Objective:   "Analyze the token contract and proposed transaction for malicious token behavior.",
		Transaction: proposal,
		Contract:    tokenContract,
		Checks:      auditChecks,
	}

	encoded, err := json.MarshalIndent(prompt, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to encode primary prompt: %w", err)
	}
	return string(encoded), nil
}

// buildSecondaryPrompt feeds the primary analysis in as prior context, so the
// second model reviews the protocol contract knowing what the first found.
func buildSecondaryPrompt(
	proposal TransactionProposal,
	tokenContract ContractArtifact,
	protocolContract ContractArtifact,
	primaryAnalysis LlmAuditResponse,
) (string, error) {
	prompt := struct {
		Objective        string              `json:"objective"`
		Transaction      TransactionProposal `json:"transaction"`
		TokenContract    ContractArtifact    `json:"tokenContract"`
		ProtocolContract ContractArtifact    `json:"protocolContract"`
		PriorAnalysis    LlmAuditResponse    `json:"priorAnalysis"`
		Checks           []string            `json:"checks"`
	}{
		Objective:        "Analyze the protocol contract using the token analysis as prior context.",
		Transaction:      proposal,
		TokenContract:    tokenContract,
		ProtocolContract: protocolContract,
		PriorAnalysis:    primaryAnalysis,
		Checks:           auditChecks,
	}

	encoded, err := json.MarshalIndent(prompt, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to encode secondary prompt: %w", err)
	}
	return string(encoded), nil
}

// ─── Enclave data collection ────────────────────────────────

func collectTransactionProposal(
	runtime cre.TeeRuntime,
	client *http.Client,
	baseURL string,
	scannerAPIKey string,
) (TransactionProposal, error) {
	response, err := getJSON(runtime, client, baseURL+"/transaction-proposal", scannerHeaders(scannerAPIKey))
	if err != nil {
		return TransactionProposal{}, err
	}
	return parseTransactionProposal(response), nil
}

func collectContractArtifact(
	runtime cre.TeeRuntime,
	client *http.Client,
	scannerURL string,
	scannerAPIKey string,
	address string,
) (ContractArtifact, error) {
	response, err := getJSON(runtime, client, scannerURL+"/contracts/"+address, scannerHeaders(scannerAPIKey))
	if err != nil {
		return ContractArtifact{}, err
	}
	return parseContractArtifact(response), nil
}

// validateScannerCredentials refuses to proceed on a credential that lacks the
// scopes this workflow depends on: a silently degraded scanner response would
// otherwise look like a clean audit.
func validateScannerCredentials(
	runtime cre.TeeRuntime,
	client *http.Client,
	scannerURL string,
	scannerAPIKey string,
) (ScannerCredentialValidation, error) {
	response, err := getJSON(runtime, client, scannerURL+"/credentials/verify", scannerHeaders(scannerAPIKey))
	if err != nil {
		return ScannerCredentialValidation{}, err
	}

	validation := parseScannerCredentialValidation(response)

	hasVerificationScope := false
	hasContractScope := false
	for _, scope := range validation.Scopes {
		switch scope {
		case "verification:read":
			hasVerificationScope = true
		case "contracts:read":
			hasContractScope = true
		}
	}

	if !validation.Valid || !hasVerificationScope || !hasContractScope {
		return ScannerCredentialValidation{}, fmt.Errorf(
			"scanner credentials failed validation provider=%s valid=%t scopes=%s",
			validation.Provider, validation.Valid, strings.Join(validation.Scopes, ","),
		)
	}

	return validation, nil
}

func requestAuditModel(
	runtime cre.TeeRuntime,
	client *http.Client,
	url string,
	apiKey string,
	model string,
	prompt string,
) (LlmAuditResponse, error) {
	response, err := postJSON(runtime, client, url, map[string]any{
		"model": model,
		"input": []map[string]string{
			{"role": "system", "content": "You are an AI smart contract audit engine. Emit strict JSON only."},
			{"role": "user", "content": prompt},
		},
	}, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + apiKey,
	})
	if err != nil {
		return LlmAuditResponse{}, err
	}

	outputText := asString(response["output_text"], "")
	if outputText == "" {
		return LlmAuditResponse{}, fmt.Errorf("LLM response missing output_text")
	}

	parsed, err := parseJSONObject(outputText)
	if err != nil {
		return LlmAuditResponse{}, err
	}
	return parseAuditResponse(parsed), nil
}

func executeAuditLog(
	runtime cre.TeeRuntime,
	client *http.Client,
	baseURL string,
	scannerAPIKey string,
	payload map[string]any,
) (map[string]any, error) {
	return postJSON(runtime, client, baseURL+"/audit-log", payload, scannerHeaders(scannerAPIKey))
}

func executeFirewallAction(
	runtime cre.TeeRuntime,
	client *http.Client,
	baseURL string,
	scannerAPIKey string,
	payload map[string]any,
) (map[string]any, error) {
	return postJSON(runtime, client, baseURL+"/firewall-action", payload, scannerHeaders(scannerAPIKey))
}

// ─── Onchain delivery ───────────────────────────────────────

// writeVerdictOnChain crosses back to the DON to sign and deliver the verdict.
// This is the one place that leaves the enclave, so only the verdict code, the
// risk mask and the chain selector are passed over — never the secrets, contract
// source or model reasoning.
//
// Returns an empty hash when no EVM target is configured, which makes the
// onchain leg opt-in.
func writeVerdictOnChain(
	runtime cre.TeeRuntime,
	config *Config,
	result FinalAuditResult,
) (string, error) {
	if len(config.EVMs) == 0 {
		return "", nil
	}

	evmConfig := config.EVMs[0]
	if evmConfig.ChainSelectorName == "" || evmConfig.ConsumerAddress == "" || evmConfig.GasLimit == "" {
		return "", nil
	}

	chainSelector, err := evm.ChainSelectorFromName(evmConfig.ChainSelectorName)
	if err != nil {
		return "", fmt.Errorf("network not found: %s: %w", evmConfig.ChainSelectorName, err)
	}

	gasLimit, err := strconv.ParseUint(evmConfig.GasLimit, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid gas_limit %q: %w", evmConfig.GasLimit, err)
	}

	encodedPayload, err := EncodeVerdictReport(result, chainSelector)
	if err != nil {
		return "", fmt.Errorf("failed to encode report payload: %w", err)
	}

	// UsingTheDons() returns a regular cre.Runtime. Anything passed into a
	// capability call on it executes on Workflow DON nodes and is NO LONGER
	// confidential.
	donRuntime := runtime.UsingTheDons()

	report, err := donRuntime.GenerateReport(&cre.ReportRequest{
		EncodedPayload: encodedPayload,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}).Await()
	if err != nil {
		return "", fmt.Errorf("failed to generate report on the DON: %w", err)
	}

	evmClient := &evm.Client{ChainSelector: chainSelector}
	writeResult, err := evmClient.WriteReport(donRuntime, &evm.WriteCreReportRequest{
		Receiver:  common.HexToAddress(evmConfig.ConsumerAddress).Bytes(),
		Report:    report,
		GasConfig: &evm.GasConfig{GasLimit: gasLimit},
	}).Await()
	if err != nil {
		return "", fmt.Errorf("onchain write failed: %w", err)
	}

	if writeResult.TxStatus != evm.TxStatus_TX_STATUS_SUCCESS {
		return "", fmt.Errorf("onchain write failed with status %v", writeResult.TxStatus)
	}

	txHash := writeResult.TxHash
	if txHash == nil {
		txHash = make([]byte, 32)
	}
	return hexutil.Encode(txHash), nil
}

// ─── Audit pipeline ─────────────────────────────────────────

// RunAuditFirewall screens one proposed transaction end to end, entirely inside
// the enclave until the optional onchain write at the end.
func RunAuditFirewall(config *Config, runtime cre.TeeRuntime, client *http.Client) (string, error) {
	ids := config.SecretsIDs

	// Batch fetch returns the secrets in the same order as the requests; if any
	// secret fails, the error reports every failed secret at once.
	secrets, err := runtime.GetSecrets([]*cre.SecretRequest{
		{Id: ids.ScannerAPIKeyID},
		{Id: ids.PrimaryLLMAPIKeyID},
		{Id: ids.SecondaryLLMAPIKeyID},
	}).Await()
	if err != nil {
		return "", fmt.Errorf("failed to fetch secrets inside the enclave: %w", err)
	}
	scannerAPIKey := secrets[0].Value
	primaryLLMAPIKey := secrets[1].Value
	secondaryLLMAPIKey := secrets[2].Value

	// ⚠️ Logs are for simulation only and MUST be removed before deploying to
	// production — anything logged inside the enclave weakens the
	// confidentiality guarantee. These record only non-sensitive markers.
	runtime.Logger().Info("audit-firewall-getsecrets-ok")

	proposal, err := collectTransactionProposal(runtime, client, config.MockBaseURL, scannerAPIKey)
	if err != nil {
		return "", err
	}

	validation, err := validateScannerCredentials(runtime, client, config.ScannerURL, scannerAPIKey)
	if err != nil {
		return "", err
	}
	runtime.Logger().Info("audit-firewall-scanner-credentials-ok",
		"provider", validation.Provider,
		"scopes", strings.Join(validation.Scopes, ","),
	)

	tokenContract, err := collectContractArtifact(
		runtime, client, config.ScannerURL, scannerAPIKey, proposal.TokenContractAddress)
	if err != nil {
		return "", err
	}

	protocolContract, err := collectContractArtifact(
		runtime, client, config.ScannerURL, scannerAPIKey, proposal.ProtocolContractAddress)
	if err != nil {
		return "", err
	}

	// An unverified contract short-circuits to DENY: without source there is
	// nothing meaningful to audit, so it fails closed rather than asking a model
	// to guess.
	if !tokenContract.Verified || !protocolContract.Verified {
		const reason = "One or more contracts are not verified by the scanner."

		auditLog, err := executeAuditLog(runtime, client, config.MockBaseURL, scannerAPIKey, map[string]any{
			"proposal":         proposal,
			"tokenContract":    tokenContract,
			"protocolContract": protocolContract,
			"verdict":          VerdictDeny,
			"reason":           reason,
			"mode":             "verification-failure",
		})
		if err != nil {
			return "", err
		}
		auditLogID := asString(auditLog["audit_log_id"], "unknown")

		firewallAction, err := executeFirewallAction(runtime, client, config.MockBaseURL, scannerAPIKey,
			map[string]any{
				"verdict":    VerdictDeny,
				"reason":     "Verification failure",
				"proposal":   proposal,
				"auditLogId": auditLogID,
			})
		if err != nil {
			return "", err
		}

		encoded, err := json.Marshal(map[string]any{
			"verdict":          VerdictDeny,
			"reason":           reason,
			"auditLogId":       auditLogID,
			"firewallActionId": asString(firewallAction["firewall_action_id"], "unknown"),
		})
		if err != nil {
			return "", fmt.Errorf("failed to encode result: %w", err)
		}
		return string(encoded), nil
	}

	primaryPrompt, err := buildPrimaryPrompt(proposal, tokenContract)
	if err != nil {
		return "", err
	}
	primaryAnalysis, err := requestAuditModel(
		runtime, client, config.PrimaryLLMURL, primaryLLMAPIKey, "audit-primary", primaryPrompt)
	if err != nil {
		return "", err
	}

	secondaryPrompt, err := buildSecondaryPrompt(proposal, tokenContract, protocolContract, primaryAnalysis)
	if err != nil {
		return "", err
	}
	secondaryAnalysis, err := requestAuditModel(
		runtime, client, config.SecondaryLLMURL, secondaryLLMAPIKey, "audit-secondary", secondaryPrompt)
	if err != nil {
		return "", err
	}

	verdict := DetermineVerdict(primaryAnalysis, secondaryAnalysis)
	riskFlags := MergeFlags(primaryAnalysis.RiskFlags, secondaryAnalysis.RiskFlags)
	reasoning := strings.Join([]string{primaryAnalysis.Reasoning, secondaryAnalysis.Reasoning}, " | ")
	analyses := AuditAnalyses{Primary: primaryAnalysis, Secondary: secondaryAnalysis}

	auditLog, err := executeAuditLog(runtime, client, config.MockBaseURL, scannerAPIKey, map[string]any{
		"proposal":         proposal,
		"tokenContract":    tokenContract,
		"protocolContract": protocolContract,
		"verdict":          verdict,
		"reasoning":        reasoning,
		"riskFlags":        riskFlags,
		"analyses":         analyses,
	})
	if err != nil {
		return "", err
	}
	auditLogID := asString(auditLog["audit_log_id"], "unknown")

	firewallAction, err := executeFirewallAction(runtime, client, config.MockBaseURL, scannerAPIKey,
		map[string]any{
			"verdict":    verdict,
			"reasoning":  reasoning,
			"proposal":   proposal,
			"riskFlags":  riskFlags,
			"auditLogId": auditLogID,
		})
	if err != nil {
		return "", err
	}

	result := FinalAuditResult{
		Verdict:          verdict,
		Reasoning:        reasoning,
		RiskFlags:        riskFlags,
		Proposal:         proposal,
		TokenContract:    tokenContract,
		ProtocolContract: protocolContract,
		Analyses:         analyses,
		AuditLogID:       auditLogID,
		FirewallActionID: asString(firewallAction["firewall_action_id"], "unknown"),
	}

	runtime.Logger().Info("audit-firewall-onchain-report-start")

	onchainTxHash, err := writeVerdictOnChain(runtime, config, result)
	if err != nil {
		return "", err
	}
	if onchainTxHash != "" {
		result.OnchainTxHash = onchainTxHash
		runtime.Logger().Info("audit-firewall-onchain", "tx_hash", onchainTxHash)
	}

	runtime.Logger().Info("audit-firewall-complete",
		"verdict", verdict,
		"audit_log_id", result.AuditLogID,
	)

	encoded, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to encode result: %w", err)
	}
	return string(encoded), nil
}

// ─── TEE Cron Callback ──────────────────────────────────────
// Receives a cre.TeeRuntime, not a cre.Runtime: the scanner credentials, the
// contract source and the model reasoning all stay inside the enclave.
func onCronTrigger(config *Config, runtime cre.TeeRuntime, _ *cron.Payload) (string, error) {
	return RunAuditFirewall(config, runtime, &http.Client{})
}

// ─── Default config ─────────────────────────────────────────

// DefaultConfig backs ParseConfig when the runner is handed an empty config
// payload. The pre-hook runs in the DON to compute restrictions and only needs
// the *shape* of the config — the secret IDs and the EVM target — so the URLs
// here are deliberately obvious placeholders rather than real endpoints.
var DefaultConfig = Config{
	Schedule:        "0 */5 * * * *",
	MockBaseURL:     "prehook-default",
	ScannerURL:      "prehook-default",
	PrimaryLLMURL:   "prehook-default",
	SecondaryLLMURL: "prehook-default",
	SecretsIDs: SecretsConfig{
		ScannerAPIKeyID:      "scanner_api_key",
		PrimaryLLMAPIKeyID:   "primary_llm_api_key",
		SecondaryLLMAPIKeyID: "secondary_llm_api_key",
	},
	EVMs: []EvmWriteConfig{{
		ChainSelectorName: "ethereum-testnet-sepolia",
		ConsumerAddress:   "0x0000000000000000000000000000000000000000",
		GasLimit:          "500000",
	}},
}

// ParseConfig decodes the workflow config, falling back to DefaultConfig when the
// payload is empty or whitespace. Without the fallback an empty payload would
// fail to unmarshal and take the pre-hook down with it.
func ParseConfig(configBytes []byte) (*Config, error) {
	if len(bytes.TrimSpace(configBytes)) == 0 {
		defaults := DefaultConfig
		return &defaults, nil
	}
	return cre.ParseJSON[Config](configBytes)
}

// ─── Restrictions ───────────────────────────────────────────

// consensusCapabilityID is spelled out here because the Go SDK only exposes a
// ConsensusRestrictor under internal_testing/, which is not a path a template
// should depend on.
const consensusCapabilityID = "consensus@1.0.0-alpha"

// Call budget for one execution. The set is CLOSED: any capability call the
// workflow makes that is not listed below is rejected outright, which caps the
// blast radius if the enclave logic is ever coerced into making extra calls.
const (
	// Steady state is 7 in-enclave HTTP calls: proposal, credential check, two
	// contract fetches, two model calls, audit log and firewall action. 8 leaves
	// exactly one call of headroom.
	maxHTTPSendRequestCalls = 8
	maxReportCalls          = 1
	maxWriteReportCalls     = 1
	maxTotalCapabilityCalls = 10
	maxSecretsFetched       = 3
)

func exactSecretRestriction(id string) *sdkpb.SecretRestriction {
	return &sdkpb.SecretRestriction{
		Restriction: &sdkpb.SecretRestriction_ExactSecret{
			ExactSecret: &cre.Secret{Id: id, Namespace: cre.DefaultSecretNamespace},
		},
	}
}

// BuildRestrictions is the pre-hook body. It runs in the DON before the enclave
// executes and declares exactly which capabilities and secrets this run may
// touch, so the workflow cannot reach anything it was not provisioned for.
func BuildRestrictions(config *Config) (*sdkpb.Restrictions, error) {
	httpRestrictor := &http.ClientRestrictor{}

	capabilityRestrictions := []*sdkpb.CapabilityRestriction{
		httpRestrictor.LimitSendRequest(maxHTTPSendRequestCalls),
		{
			Restriction: &sdkpb.CapabilityRestriction_Method{
				Method: &sdkpb.MethodRestriction{
					Id:       consensusCapabilityID,
					Method:   "Report",
					MaxCalls: maxReportCalls,
				},
			},
		},
	}

	// The onchain leg is opt-in, so only grant WriteReport when an EVM target is
	// configured. An unresolvable chain name is skipped rather than fatal: the
	// write path itself reports that error at execution time.
	if len(config.EVMs) > 0 && config.EVMs[0].ChainSelectorName != "" {
		if chainSelector, err := evm.ChainSelectorFromName(config.EVMs[0].ChainSelectorName); err == nil {
			evmRestrictor := &evm.ClientRestrictor{ChainSelector: chainSelector}
			capabilityRestrictions = append(
				capabilityRestrictions,
				evmRestrictor.LimitWriteReport(maxWriteReportCalls),
			)
		}
	}

	ids := config.SecretsIDs

	return &sdkpb.Restrictions{
		Capabilities: &sdkpb.CapabilityRestrictions{
			Type:          sdkpb.CapabilityRestrictionType_CAPABILITY_RESTRICTION_TYPE_CLOSED,
			MaxTotalCalls: maxTotalCapabilityCalls,
			Restrictions:  capabilityRestrictions,
		},
		Secrets: &sdkpb.SecretsRestritions{
			MaxSecrets: maxSecretsFetched,
			Restrictions: []*sdkpb.SecretRestriction{
				exactSecretRestriction(ids.ScannerAPIKeyID),
				exactSecretRestriction(ids.PrimaryLLMAPIKeyID),
				exactSecretRestriction(ids.SecondaryLLMAPIKeyID),
			},
		},
	}, nil
}

// ─── Workflow Init ──────────────────────────────────────────

func InitWorkflow(config *Config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[*Config], error) {
	if config.Schedule == "" ||
		config.MockBaseURL == "" ||
		config.ScannerURL == "" ||
		config.PrimaryLLMURL == "" ||
		config.SecondaryLLMURL == "" {
		return nil, fmt.Errorf(
			"config requires schedule, mock_base_url, scanner_url, primary_llm_url, and secondary_llm_url")
	}

	ids := config.SecretsIDs
	if ids.ScannerAPIKeyID == "" || ids.PrimaryLLMAPIKeyID == "" || ids.SecondaryLLMAPIKeyID == "" {
		return nil, fmt.Errorf("config requires secrets_ids fields")
	}

	return cre.Workflow[*Config]{
		// cre.HandlerInTee instead of cre.Handler. The third argument declares
		// which enclaves this handler accepts. Alternatives:
		//   cre.AnyTee{}                                               — any registered TEE, any region
		//   cre.AnyTeeInRegions{Regions: []cre.Region{cre.AwsUsWest2}} — any TEE, multiple regions
		// The fourth argument is the pre-hook: it runs in the DON before the
		// enclave executes and returns the restrictions for that execution.
		cre.HandlerInTeeWithPreHook(
			cron.Trigger(&cron.Config{Schedule: config.Schedule}),
			onCronTrigger,
			cre.AnyTee{},
			func(cfg *Config, _ *cron.Payload) (*sdkpb.Restrictions, error) {
				return BuildRestrictions(cfg)
			},
		),
	}, nil
}
