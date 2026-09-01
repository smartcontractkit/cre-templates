package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	httpmock "github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http/mock"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/testutils"
)

const (
	exchangeAPIKey = "test-exchange-key"
	openAIAPIKey   = "test-openai-key"
	mockBaseURL    = "http://127.0.0.1:8787/liquidation"
	openAIURL      = "http://127.0.0.1:8787/liquidation/v1/responses"
)

func policyFixture() Policy {
	return Policy{
		MaxReserveDeploymentUSDC:    5000,
		MinReserveBalanceUSDC:       2000,
		LiquidationWarningThreshold: 18,
		MinimumHealthFactor:         1.25,
		TargetHealthFactor:          1.5,
		MaxCollateralAllocationPct:  80,
		MaxPartialDebtRepaymentPct:  25,
		ExecutionSequencePreference: SequenceCollateralFirst,
		PreferredVenues:             []string{"binance", "onchain"},
	}
}

func riskFixture() RiskState {
	return RiskState{
		CollateralAssetSymbol:   "ETH",
		BorrowedAssetSymbol:     "USDC",
		CollateralAssetPriceUSD: 2500,
		BorrowedAssetPriceUSD:   1,
		CollateralBalanceUSD:    45000,
		OutstandingDebtUSD:      32000,
		CollateralHealthFactor:  1.12,
		LiquidationProximityPct: 11,
		LoanToValuePct:          71,
		LiquidationThresholdPct: 78,
		VolatilityIndex:         0.35,
		USDCReserve:             10000,
		CashBalanceUSD:          12000,
	}
}

func testConfig() *Config {
	return &Config{
		Schedule:                            "0 */5 * * * *",
		MockBaseURL:                         mockBaseURL,
		OpenAIURL:                           openAIURL,
		OpenAIModel:                         "gpt-4.1-mini",
		DefensiveActionSequencingPreference: "collateral-first",
		SecretsIDs: SecretsConfig{
			ExchangeAPIKeyID: "exchange_api_key",
			OpenAIAPIKeyID:   "openai_api_key",
			LiquidationWarningActionThresholdSecretID:  "liquidation_liquidation_warning_action_threshold",
			MinimumHealthFactorSecretID:                "liquidation_minimum_health_factor",
			TargetHealthFactorSecretID:                 "liquidation_target_health_factor",
			MaximumStablecoinReserveDeploymentSecretID: "liquidation_maximum_stablecoin_reserve_deployment",
			MinimumStablecoinReserveBalanceSecretID:    "liquidation_minimum_stablecoin_reserve_balance",
			MaximumCollateralAllocationSecretID:        "liquidation_maximum_collateral_allocation",
			MaximumPartialDebtRepaymentSecretID:        "liquidation_maximum_partial_debt_repayment",
			PreferredVenuesSecretID:                    "liquidation_preferred_venues",
		},
	}
}

// testSecrets mirrors the values .env.example seeds, so the policy the enclave
// assembles here matches policyFixture().
func testSecrets() testutils.Secrets {
	return testutils.Secrets{
		cre.DefaultSecretNamespace: {
			"exchange_api_key": exchangeAPIKey,
			"openai_api_key":   openAIAPIKey,
			"liquidation_liquidation_warning_action_threshold":  "18",
			"liquidation_minimum_health_factor":                 "1.25",
			"liquidation_target_health_factor":                  "1.5",
			"liquidation_maximum_stablecoin_reserve_deployment": "5000",
			"liquidation_minimum_stablecoin_reserve_balance":    "2000",
			"liquidation_maximum_collateral_allocation":         "80",
			"liquidation_maximum_partial_debt_repayment":        "25",
			"liquidation_preferred_venues":                      "binance,onchain",
		},
	}
}

// capturedCall records one in-enclave HTTP call so tests can assert on what
// crossed the enclave boundary.
type capturedCall struct {
	URL     string
	Method  string
	Headers map[string][]string
	Body    string
}

// stubEnclaveHTTP routes the three endpoints the workflow calls through a single
// in-enclave stub, and records every request.
func stubEnclaveHTTP(t *testing.T, riskBody string, llmText string) *[]capturedCall {
	t.Helper()

	capability, err := httpmock.NewClientCapability(t)
	require.NoError(t, err)

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
		case strings.HasSuffix(input.Url, "/risk-state"):
			return &http.Response{StatusCode: 200, Body: []byte(riskBody)}, nil
		case strings.HasSuffix(input.Url, "/v1/responses"):
			payload, err := json.Marshal(map[string]any{"output_text": llmText})
			require.NoError(t, err)
			return &http.Response{StatusCode: 200, Body: payload}, nil
		case strings.HasSuffix(input.Url, "/execute-defense"):
			return &http.Response{
				StatusCode: 200,
				Body:       []byte(`{"execution_id":"exec-123"}`),
			}, nil
		}
		return nil, assert.AnError
	}
	return calls
}

func riskStateJSON(proximityPct float64) string {
	risk := riskFixture()
	risk.LiquidationProximityPct = proximityPct
	encoded, err := json.Marshal(risk)
	if err != nil {
		panic(err)
	}
	return string(encoded)
}

func defendDecisionJSON() string {
	return `{"shouldDefend":true,"reasoning":"liquidation risk elevated","actions":[` +
		`{"type":"partial_debt_repayment","repayPct":70},` +
		`{"type":"add_collateral","amountUsd":7000}]}`
}

// ─── ComputeRiskScore ───────────────────────────────────────

func TestComputeRiskScore_RisesWithLiquidationAndDebtRisk(t *testing.T) {
	score := ComputeRiskScore(riskFixture(), policyFixture())
	assert.Greater(t, score, 0.0)

	safer := riskFixture()
	safer.LiquidationProximityPct = 30
	safer.LoanToValuePct = 55
	safer.CollateralHealthFactor = 1.5
	safer.VolatilityIndex = 0.05

	assert.Greater(t, score, ComputeRiskScore(safer, policyFixture()))
}

func TestComputeRiskScore_IsDeterministic(t *testing.T) {
	assert.Equal(t,
		ComputeRiskScore(riskFixture(), policyFixture()),
		ComputeRiskScore(riskFixture(), policyFixture()),
	)
}

// ─── EnforcePolicy ──────────────────────────────────────────

func TestEnforcePolicy_CapsReserveDeploymentAndHonoursSequence(t *testing.T) {
	decision, err := parseLlmDecision(defendDecisionJSON())
	require.NoError(t, err)

	actions, err := EnforcePolicy(decision, policyFixture(), riskFixture())
	require.NoError(t, err)

	require.Len(t, actions, 2)
	// collateral-first reorders add_collateral ahead of the repayment.
	assert.Equal(t, ActionAddCollateral, actions[0].Type)
	assert.Equal(t, 5000.0, actions[0].AmountUsd, "7000 must clamp to the 5000 reserve cap")
	assert.Equal(t, "binance", actions[0].Venue)
	assert.Equal(t, ActionPartialDebtRepayment, actions[1].Type)
	assert.Equal(t, 25.0, actions[1].RepayPct, "70%% must clamp to the 25%% policy cap")
	assert.Equal(t, 8000.0, actions[1].AmountUsd)
}

func TestEnforcePolicy_RejectsReserveFloorBreach(t *testing.T) {
	risk := riskFixture()
	risk.USDCReserve = 6000

	decision := LiquidationDecision{
		ShouldDefend: true,
		Reasoning:    "critical",
		Actions:      []DefensiveAction{{Type: ActionAddCollateral, AmountUsd: 9500}},
	}

	_, err := EnforcePolicy(decision, policyFixture(), risk)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "breaches reserve floor")
}

func TestEnforcePolicy_ReturnsNothingWhenModelDeclinesToDefend(t *testing.T) {
	decision := LiquidationDecision{
		ShouldDefend: false,
		Actions:      []DefensiveAction{{Type: ActionAddCollateral, AmountUsd: 100}},
	}

	actions, err := EnforcePolicy(decision, policyFixture(), riskFixture())
	require.NoError(t, err)
	assert.Empty(t, actions)
}

func TestEnforcePolicy_DebtFirstReordersActions(t *testing.T) {
	policy := policyFixture()
	policy.ExecutionSequencePreference = SequenceDebtFirst

	decision, err := parseLlmDecision(defendDecisionJSON())
	require.NoError(t, err)

	actions, err := EnforcePolicy(decision, policy, riskFixture())
	require.NoError(t, err)

	require.Len(t, actions, 2)
	assert.Equal(t, ActionPartialDebtRepayment, actions[0].Type)
	assert.Equal(t, ActionAddCollateral, actions[1].Type)
}

// ─── Untrusted model output ─────────────────────────────────

func TestParseLlmDecision_DropsUnknownActionTypes(t *testing.T) {
	decision, err := parseLlmDecision(
		`{"shouldDefend":true,"actions":[{"type":"drain_vault","amountUsd":1},` +
			`{"type":"add_collateral","amountUsd":2}]}`,
	)
	require.NoError(t, err)

	require.Len(t, decision.Actions, 1)
	assert.Equal(t, ActionAddCollateral, decision.Actions[0].Type)
}

func TestParseLlmDecision_CoercesStringNumbers(t *testing.T) {
	decision, err := parseLlmDecision(
		`{"shouldDefend":true,"actions":[{"type":"add_collateral","amountUsd":"1500"}]}`,
	)
	require.NoError(t, err)

	require.Len(t, decision.Actions, 1)
	assert.Equal(t, 1500.0, decision.Actions[0].AmountUsd)
}

func TestParseLlmDecision_DefaultsReasoning(t *testing.T) {
	decision, err := parseLlmDecision(`{"shouldDefend":false}`)
	require.NoError(t, err)
	assert.Equal(t, "No reasoning provided by model", decision.Reasoning)
}

func TestParseLlmDecision_ErrorsOnInvalidJSON(t *testing.T) {
	_, err := parseLlmDecision("not json")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid json response")
}

func TestExtractOpenAiText_ReadsNestedOutputShape(t *testing.T) {
	payload := map[string]any{
		"output": []any{
			map[string]any{"content": []any{map[string]any{"text": "nested"}}},
		},
	}

	text, err := extractOpenAiText(payload)
	require.NoError(t, err)
	assert.Equal(t, "nested", text)
}

func TestExtractOpenAiText_ErrorsWhenAbsent(t *testing.T) {
	_, err := extractOpenAiText(map[string]any{})
	require.Error(t, err)
}

func TestParseExecutionSequencePreference_FallsBackToCollateralFirst(t *testing.T) {
	assert.Equal(t, SequenceDebtFirst, parseExecutionSequencePreference("debt-first"))
	assert.Equal(t, SequenceBalanced, parseExecutionSequencePreference("balanced"))
	assert.Equal(t, SequenceCollateralFirst, parseExecutionSequencePreference("nonsense"))
}

func TestParseRequiredSecretNumber_RejectsNonNumeric(t *testing.T) {
	_, err := parseRequiredSecretNumber("abc", "some_secret")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be a finite number")
}

// ─── End-to-end through the enclave runtime ─────────────────

func TestOnCronTrigger_DefendsWhenActionsSurvivePolicy(t *testing.T) {
	stubEnclaveHTTP(t, riskStateJSON(11), defendDecisionJSON())
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	result, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal([]byte(result), &decoded))
	assert.Equal(t, "DEFENDED", decoded["status"])
	assert.Equal(t, 2.0, decoded["actionCount"])
	assert.Equal(t, "exec-123", decoded["executionId"])
}

// Proximity above the threshold means the position is still comfortably far
// from liquidation, so the workflow must stand down.
func TestOnCronTrigger_ReturnsSafeWhenProximityAboveThreshold(t *testing.T) {
	stubEnclaveHTTP(t, riskStateJSON(40), defendDecisionJSON())
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	result, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)
	assert.Equal(t, "SAFE", result)
}

func TestOnCronTrigger_ReturnsSafeWhenModelDeclines(t *testing.T) {
	stubEnclaveHTTP(t, riskStateJSON(11), `{"shouldDefend":false,"actions":[]}`)
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	result, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)
	assert.Equal(t, "SAFE", result)
}

func TestOnCronTrigger_InjectsEnclaveFetchedSecretsIntoRequests(t *testing.T) {
	calls := stubEnclaveHTTP(t, riskStateJSON(11), defendDecisionJSON())
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)

	require.Len(t, *calls, 3)

	risk, llm, defense := (*calls)[0], (*calls)[1], (*calls)[2]

	assert.Equal(t, "GET", risk.Method)
	assert.Equal(t, []string{exchangeAPIKey}, risk.Headers["x-exchange-api-key"])

	assert.Equal(t, "POST", llm.Method)
	assert.Equal(t, []string{"Bearer " + openAIAPIKey}, llm.Headers["Authorization"])

	assert.Equal(t, "POST", defense.Method)
	assert.Equal(t, []string{exchangeAPIKey}, defense.Headers["x-exchange-api-key"])
}

// The confidential policy must reach the model, otherwise the enclave is doing
// no useful work over the secret data.
func TestOnCronTrigger_SendsConfidentialPolicyToModel(t *testing.T) {
	calls := stubEnclaveHTTP(t, riskStateJSON(11), defendDecisionJSON())
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)

	llmBody := (*calls)[1].Body
	assert.Contains(t, llmBody, "max_reserve_deployment_usdc")
	assert.Contains(t, llmBody, "riskScore")
}

func TestOnCronTrigger_ErrorsOnNon2xxRiskFetch(t *testing.T) {
	capability, err := httpmock.NewClientCapability(t)
	require.NoError(t, err)
	capability.SendRequest = func(_ context.Context, _ *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 401, Body: []byte("unauthorized")}, nil
	}

	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err = onCronTrigger(testConfig(), runtime, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "status=401")
}

func TestOnCronTrigger_ErrorsWhenPolicySecretIsNotNumeric(t *testing.T) {
	stubEnclaveHTTP(t, riskStateJSON(11), defendDecisionJSON())

	secrets := testSecrets()
	secrets[cre.DefaultSecretNamespace]["liquidation_minimum_health_factor"] = "not-a-number"
	runtime := testutils.NewTeeRuntime(t, secrets)

	_, err := onCronTrigger(testConfig(), runtime, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be a finite number")
}

// Logging inside an enclave weakens its confidentiality guarantee. The workflow
// logs only non-sensitive markers, never secrets or confidential payloads.
func TestOnCronTrigger_DoesNotLogSecretsOrConfidentialPayloads(t *testing.T) {
	stubEnclaveHTTP(t, riskStateJSON(11), defendDecisionJSON())
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)

	for _, raw := range runtime.GetLogs() {
		line := string(raw)
		assert.NotContains(t, line, exchangeAPIKey, "log leaked the exchange API key")
		assert.NotContains(t, line, openAIAPIKey, "log leaked the OpenAI API key")
	}
}

// ─── InitWorkflow ───────────────────────────────────────────

func TestInitWorkflow_RegistersOneTeeHandler(t *testing.T) {
	workflow, err := InitWorkflow(testConfig(), nil, nil)
	require.NoError(t, err)
	assert.Len(t, workflow, 1)
}

func TestInitWorkflow_ErrorsWhenCoreFieldsMissing(t *testing.T) {
	config := testConfig()
	config.OpenAIURL = ""

	_, err := InitWorkflow(config, nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "config requires schedule, mock_base_url, openai_url, openai_model, and defensive_action_sequencing_preference")
}

func TestInitWorkflow_ErrorsWhenSecretIDsMissing(t *testing.T) {
	config := testConfig()
	config.SecretsIDs.PreferredVenuesSecretID = ""

	_, err := InitWorkflow(config, nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "config requires secrets_ids fields")
}
