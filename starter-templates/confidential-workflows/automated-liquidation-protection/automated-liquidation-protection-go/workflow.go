package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

// ─── Config ─────────────────────────────────────────────────

// SecretsConfig maps each policy input to the secret ID that holds it. The
// policy itself is confidential: publishing thresholds on-chain or in config
// would tell an adversary exactly when this vault defends and how much it will
// spend doing it.
type SecretsConfig struct {
	ExchangeAPIKeyID                           string `json:"exchange_api_key_id"`
	OpenAIAPIKeyID                             string `json:"openai_api_key_id"`
	LiquidationWarningActionThresholdSecretID  string `json:"liquidation_warning_action_threshold_secret_id"`
	MinimumHealthFactorSecretID                string `json:"minimum_health_factor_secret_id"`
	TargetHealthFactorSecretID                 string `json:"target_health_factor_secret_id"`
	MaximumStablecoinReserveDeploymentSecretID string `json:"maximum_stablecoin_reserve_deployment_secret_id"`
	MinimumStablecoinReserveBalanceSecretID    string `json:"minimum_stablecoin_reserve_balance_secret_id"`
	MaximumCollateralAllocationSecretID        string `json:"maximum_collateral_allocation_secret_id"`
	MaximumPartialDebtRepaymentSecretID        string `json:"maximum_partial_debt_repayment_secret_id"`
	PreferredVenuesSecretID                    string `json:"preferred_venues_secret_id"`
}

type Config struct {
	Schedule    string `json:"schedule"`
	MockBaseURL string `json:"mock_base_url"`
	OpenAIURL   string `json:"openai_url"`
	OpenAIModel string `json:"openai_model"`
	// DefensiveActionSequencingPreference is plain config, not a secret: it only
	// reveals the ordering of defensive actions, never the thresholds or amounts
	// that decide when and how much the vault defends.
	DefensiveActionSequencingPreference string        `json:"defensive_action_sequencing_preference"`
	SecretsIDs                          SecretsConfig `json:"secrets_ids"`
}

// ─── Domain types ───────────────────────────────────────────

// Policy is assembled inside the enclave from Vault DON secrets. It never
// leaves the enclave — only the resulting actions do.
type Policy struct {
	MaxReserveDeploymentUSDC    float64  `json:"max_reserve_deployment_usdc"`
	MinReserveBalanceUSDC       float64  `json:"min_reserve_balance_usdc"`
	LiquidationWarningThreshold float64  `json:"liquidation_warning_action_threshold"`
	MinimumHealthFactor         float64  `json:"minimum_health_factor"`
	TargetHealthFactor          float64  `json:"target_health_factor"`
	MaxCollateralAllocationPct  float64  `json:"max_collateral_allocation_pct"`
	MaxPartialDebtRepaymentPct  float64  `json:"max_partial_debt_repayment_pct"`
	ExecutionSequencePreference string   `json:"execution_sequence_preference"`
	PreferredVenues             []string `json:"preferred_venues"`
}

type RiskState struct {
	CollateralAssetSymbol   string  `json:"collateral_asset_symbol"`
	BorrowedAssetSymbol     string  `json:"borrowed_asset_symbol"`
	CollateralAssetPriceUSD float64 `json:"collateral_asset_price_usd"`
	BorrowedAssetPriceUSD   float64 `json:"borrowed_asset_price_usd"`
	CollateralBalanceUSD    float64 `json:"collateral_balance_usd"`
	OutstandingDebtUSD      float64 `json:"outstanding_debt_usd"`
	CollateralHealthFactor  float64 `json:"collateral_health_factor"`
	LiquidationProximityPct float64 `json:"liquidation_proximity_pct"`
	LoanToValuePct          float64 `json:"loan_to_value_pct"`
	LiquidationThresholdPct float64 `json:"liquidation_threshold_pct"`
	VolatilityIndex         float64 `json:"volatility_index"`
	USDCReserve             float64 `json:"usdc_reserve"`
	CashBalanceUSD          float64 `json:"cash_balance_usd"`
}

// Action type discriminators emitted by the model.
const (
	ActionAddCollateral                 = "add_collateral"
	ActionBridgeAndAddCollateral        = "bridge_and_add_collateral"
	ActionSwapReserveToCollateral       = "swap_reserve_to_collateral"
	ActionRepayWithReserves             = "repay_with_reserves"
	ActionSwapReserveToBorrowedAndRepay = "swap_reserve_to_borrowed_and_repay"
	ActionPartialDebtRepayment          = "partial_debt_repayment"
	ActionFullDebtRepayment             = "full_debt_repayment"
)

// Execution sequence preferences.
const (
	SequenceCollateralFirst = "collateral-first"
	SequenceDebtFirst       = "debt-first"
	SequenceBalanced        = "balanced"
)

var actionTypes = []string{
	ActionAddCollateral,
	ActionBridgeAndAddCollateral,
	ActionSwapReserveToCollateral,
	ActionRepayWithReserves,
	ActionSwapReserveToBorrowedAndRepay,
	ActionPartialDebtRepayment,
	ActionFullDebtRepayment,
}

func isActionType(value string) bool {
	for _, t := range actionTypes {
		if t == value {
			return true
		}
	}
	return false
}

// DefensiveAction is one step the model proposes. After enforcePolicy has
// clamped it, AmountUsd/RepayPct/Venue are always populated — the alias below
// marks that post-enforcement state at call sites.
type DefensiveAction struct {
	Type      string  `json:"type"`
	AmountUsd float64 `json:"amountUsd"`
	RepayPct  float64 `json:"repayPct"`
	FromChain string  `json:"fromChain"`
	ToChain   string  `json:"toChain"`
	Venue     string  `json:"venue"`
}

// ExecutableAction is a DefensiveAction that has passed policy enforcement.
type ExecutableAction = DefensiveAction

type LiquidationDecision struct {
	ShouldDefend bool
	Reasoning    string
	Actions      []DefensiveAction
}

// ─── Lenient JSON coercion ──────────────────────────────────
// The model's output is untrusted: it may emit numbers as strings, omit keys, or
// return the wrong shape entirely. These helpers coerce rather than fail so a
// sloppy-but-usable response still produces a decision, mirroring the
// TypeScript template's asNumber/asString/asObject behaviour.

func asObject(value any) map[string]any {
	obj, ok := value.(map[string]any)
	if !ok {
		return map[string]any{}
	}
	return obj
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

func asString(value any, fallback string) string {
	if s, ok := value.(string); ok {
		return s
	}
	return fallback
}

func parseJSONObject(body string) (map[string]any, error) {
	var parsed any
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		return nil, fmt.Errorf("invalid json response: %w; body=%s", err, body)
	}
	return asObject(parsed), nil
}

// ─── Policy inputs from secrets ─────────────────────────────

func parseRequiredSecretNumber(secretValue, secretID string) (float64, error) {
	parsed, err := strconv.ParseFloat(strings.TrimSpace(secretValue), 64)
	if err != nil || math.IsInf(parsed, 0) || math.IsNaN(parsed) {
		return 0, fmt.Errorf("secret %s must be a finite number", secretID)
	}
	return parsed, nil
}

// parseExecutionSequencePreference falls back to collateral-first for anything
// unrecognised, so a corrupted config value cannot silently invent a new ordering.
func parseExecutionSequencePreference(value string) string {
	switch strings.TrimSpace(value) {
	case SequenceDebtFirst:
		return SequenceDebtFirst
	case SequenceBalanced:
		return SequenceBalanced
	default:
		return SequenceCollateralFirst
	}
}

func parseVenueListSecret(secretValue string) []string {
	venues := []string{}
	for _, entry := range strings.Split(secretValue, ",") {
		if trimmed := strings.TrimSpace(entry); trimmed != "" {
			venues = append(venues, trimmed)
		}
	}
	return venues
}

func parseRiskState(payload any) RiskState {
	row := asObject(payload)
	return RiskState{
		CollateralAssetSymbol:   asString(row["collateral_asset_symbol"], "ETH"),
		BorrowedAssetSymbol:     asString(row["borrowed_asset_symbol"], "USDC"),
		CollateralAssetPriceUSD: asNumber(row["collateral_asset_price_usd"], 0),
		BorrowedAssetPriceUSD:   asNumber(row["borrowed_asset_price_usd"], 1),
		CollateralBalanceUSD:    asNumber(row["collateral_balance_usd"], 0),
		OutstandingDebtUSD:      asNumber(row["outstanding_debt_usd"], 0),
		CollateralHealthFactor:  asNumber(row["collateral_health_factor"], 0),
		LiquidationProximityPct: asNumber(row["liquidation_proximity_pct"], 0),
		LoanToValuePct:          asNumber(row["loan_to_value_pct"], 0),
		LiquidationThresholdPct: asNumber(row["liquidation_threshold_pct"], 0),
		VolatilityIndex:         asNumber(row["volatility_index"], 0),
		USDCReserve:             asNumber(row["usdc_reserve"], 0),
		CashBalanceUSD:          asNumber(row["cash_balance_usd"], 0),
	}
}

// ─── Risk scoring ───────────────────────────────────────────

// ComputeRiskScore blends four independent pressures into a single score. It is
// deliberately pure and deterministic: the enclave result is attested and
// verified by DON consensus, so the same inputs must always yield the same score.
func ComputeRiskScore(risk RiskState, policy Policy) float64 {
	proximityRisk := math.Max(0, policy.LiquidationWarningThreshold-risk.LiquidationProximityPct) * 5
	ltvBufferRisk := math.Max(0, risk.LoanToValuePct-(risk.LiquidationThresholdPct-5)) * 2
	healthRisk := math.Max(0, policy.MinimumHealthFactor-risk.CollateralHealthFactor) * 100
	volatilityRisk := risk.VolatilityIndex * 25

	return proximityRisk + ltvBufferRisk + healthRisk + volatilityRisk
}

// ─── Confidential HTTP helpers ──────────────────────────────
// Both helpers go through SendRequestInTee, so URLs, headers (including the API
// keys) and response bodies stay inside the enclave.

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

func multiHeaders(headers map[string]string) map[string]*http.HeaderValues {
	out := make(map[string]*http.HeaderValues, len(headers))
	for key, value := range headers {
		out[key] = &http.HeaderValues{Values: []string{value}}
	}
	return out
}

func jsonHeaders() map[string]string {
	return map[string]string{"Content-Type": "application/json"}
}

// ─── Model plumbing ─────────────────────────────────────────

// extractOpenAiText handles both the flat output_text shape and the nested
// output[].content[].text shape the Responses API can return.
func extractOpenAiText(payload map[string]any) (string, error) {
	if direct := asString(payload["output_text"], ""); direct != "" {
		return direct, nil
	}

	outputItems, _ := payload["output"].([]any)
	for _, item := range outputItems {
		content, _ := asObject(item)["content"].([]any)
		for _, part := range content {
			if text := asString(asObject(part)["text"], ""); text != "" {
				return text, nil
			}
		}
	}
	return "", fmt.Errorf("openai response did not contain output_text")
}

// parseLlmDecision silently drops actions with an unrecognised type rather than
// failing the run — an unknown action is one the executor could not carry out
// anyway, and the remaining valid actions are still worth executing.
func parseLlmDecision(rawText string) (LiquidationDecision, error) {
	parsed, err := parseJSONObject(rawText)
	if err != nil {
		return LiquidationDecision{}, err
	}

	shouldDefend, _ := parsed["shouldDefend"].(bool)
	actionsRaw, _ := parsed["actions"].([]any)

	actions := []DefensiveAction{}
	for _, item := range actionsRaw {
		row := asObject(item)
		actionType := asString(row["type"], "")
		if !isActionType(actionType) {
			continue
		}

		actions = append(actions, DefensiveAction{
			Type:      actionType,
			AmountUsd: asNumber(row["amountUsd"], 0),
			RepayPct:  asNumber(row["repayPct"], 0),
			FromChain: asString(row["fromChain"], ""),
			ToChain:   asString(row["toChain"], ""),
			Venue:     asString(row["venue"], ""),
		})
	}

	return LiquidationDecision{
		ShouldDefend: shouldDefend,
		Reasoning:    asString(parsed["reasoning"], "No reasoning provided by model"),
		Actions:      actions,
	}, nil
}

func createOpenAiPrompt(risk RiskState, riskScore float64, policy Policy) (string, error) {
	prompt := struct {
		Objective string    `json:"objective"`
		Policy    Policy    `json:"policy"`
		Risk      RiskState `json:"risk"`
		RiskScore float64   `json:"riskScore"`
	}{
		Objective: "Return strict JSON with shouldDefend:boolean, reasoning:string, " +
			"actions:[{type, amountUsd, repayPct, fromChain, toChain, venue}] for liquidation " +
			"defense via collateral and debt-repayment actions while respecting policy.",
		Policy:    policy,
		Risk:      risk,
		RiskScore: riskScore,
	}

	encoded, err := json.MarshalIndent(prompt, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to encode prompt: %w", err)
	}
	return string(encoded), nil
}

// ─── Policy enforcement ─────────────────────────────────────

func chooseVenue(action DefensiveAction, preferredVenues []string) string {
	if action.Venue != "" {
		return action.Venue
	}
	if len(preferredVenues) > 0 {
		return preferredVenues[0]
	}
	return "default-venue"
}

// sequencePriorities ranks action types per preference. Ordering matters: adding
// collateral before repaying debt has a different risk profile than the reverse.
var sequencePriorities = map[string]map[string]int{
	SequenceCollateralFirst: {
		ActionAddCollateral:                 1,
		ActionBridgeAndAddCollateral:        2,
		ActionSwapReserveToCollateral:       3,
		ActionRepayWithReserves:             4,
		ActionSwapReserveToBorrowedAndRepay: 5,
		ActionPartialDebtRepayment:          6,
		ActionFullDebtRepayment:             7,
	},
	SequenceDebtFirst: {
		ActionRepayWithReserves:             1,
		ActionSwapReserveToBorrowedAndRepay: 2,
		ActionPartialDebtRepayment:          3,
		ActionFullDebtRepayment:             4,
		ActionAddCollateral:                 5,
		ActionBridgeAndAddCollateral:        6,
		ActionSwapReserveToCollateral:       7,
	},
	SequenceBalanced: {
		ActionAddCollateral:                 4,
		ActionBridgeAndAddCollateral:        2,
		ActionSwapReserveToCollateral:       6,
		ActionRepayWithReserves:             3,
		ActionSwapReserveToBorrowedAndRepay: 5,
		ActionPartialDebtRepayment:          1,
		ActionFullDebtRepayment:             7,
	},
}

func orderActions(actions []ExecutableAction, preference string) []ExecutableAction {
	priorities, ok := sequencePriorities[preference]
	if !ok {
		priorities = sequencePriorities[SequenceCollateralFirst]
	}

	ordered := make([]ExecutableAction, len(actions))
	copy(ordered, actions)
	sort.SliceStable(ordered, func(i, j int) bool {
		return priorities[ordered[i].Type] < priorities[ordered[j].Type]
	})
	return ordered
}

// EnforcePolicy is the guardrail between the model and the executor. The model
// only ever proposes; every amount is clamped to the policy caps and the running
// reserve is checked against the floor, so a hallucinated or adversarial
// suggestion cannot drain the vault. A breach aborts the whole run rather than
// executing a partial sequence.
func EnforcePolicy(decision LiquidationDecision, policy Policy, risk RiskState) ([]ExecutableAction, error) {
	if !decision.ShouldDefend || len(decision.Actions) == 0 {
		return []ExecutableAction{}, nil
	}

	projectedReserve := risk.USDCReserve
	executable := []ExecutableAction{}

	for _, action := range decision.Actions {
		amount := math.Max(0, action.AmountUsd)
		cappedRepayPct := math.Min(math.Max(0, action.RepayPct), policy.MaxPartialDebtRepaymentPct)

		switch action.Type {
		case ActionAddCollateral, ActionBridgeAndAddCollateral, ActionSwapReserveToCollateral,
			ActionRepayWithReserves, ActionSwapReserveToBorrowedAndRepay:
			capped := math.Min(amount, policy.MaxReserveDeploymentUSDC)
			projectedReserve -= capped

			if projectedReserve < policy.MinReserveBalanceUSDC {
				return nil, fmt.Errorf(
					"action %s breaches reserve floor: projected %.2f < floor %v",
					action.Type, projectedReserve, policy.MinReserveBalanceUSDC,
				)
			}

			settled := action
			settled.AmountUsd = capped
			settled.RepayPct = 0
			settled.Venue = chooseVenue(action, policy.PreferredVenues)
			executable = append(executable, settled)

		case ActionPartialDebtRepayment:
			boundedRepayPct := math.Min(cappedRepayPct, 100)

			settled := action
			settled.AmountUsd = risk.OutstandingDebtUSD * boundedRepayPct / 100
			settled.RepayPct = boundedRepayPct
			settled.Venue = chooseVenue(action, policy.PreferredVenues)
			executable = append(executable, settled)

		case ActionFullDebtRepayment:
			amountUsd := risk.OutstandingDebtUSD
			projectedReserve -= amountUsd

			if projectedReserve < policy.MinReserveBalanceUSDC {
				return nil, fmt.Errorf(
					"action %s breaches reserve floor: projected %.2f < floor %v",
					action.Type, projectedReserve, policy.MinReserveBalanceUSDC,
				)
			}

			settled := action
			settled.AmountUsd = amountUsd
			settled.RepayPct = 100
			settled.Venue = chooseVenue(action, policy.PreferredVenues)
			executable = append(executable, settled)

		default:
			settled := action
			settled.AmountUsd = amount
			settled.RepayPct = 0
			settled.Venue = chooseVenue(action, policy.PreferredVenues)
			executable = append(executable, settled)
		}
	}

	return orderActions(executable, policy.ExecutionSequencePreference), nil
}

// ─── Enclave data collection ────────────────────────────────

func collectRiskSnapshot(
	runtime cre.TeeRuntime,
	client *http.Client,
	baseURL string,
	exchangeAPIKey string,
) (RiskState, error) {
	authHeaders := jsonHeaders()
	authHeaders["x-exchange-api-key"] = exchangeAPIKey

	riskResponse, err := getJSON(runtime, client, baseURL+"/risk-state", authHeaders)
	if err != nil {
		return RiskState{}, err
	}
	return parseRiskState(riskResponse), nil
}

// buildPolicy assembles the policy from the batch-fetched Vault DON secrets.
// The secrets are released only into an attested enclave and decrypted at the
// moment GetSecrets runs, so node operators never see the thresholds.
func buildPolicy(secretValues map[string]string, ids SecretsConfig, sequencingPreference string) (Policy, error) {
	numeric := func(secretID string) (float64, error) {
		return parseRequiredSecretNumber(secretValues[secretID], secretID)
	}

	warningThreshold, err := numeric(ids.LiquidationWarningActionThresholdSecretID)
	if err != nil {
		return Policy{}, err
	}
	minimumHealthFactor, err := numeric(ids.MinimumHealthFactorSecretID)
	if err != nil {
		return Policy{}, err
	}
	targetHealthFactor, err := numeric(ids.TargetHealthFactorSecretID)
	if err != nil {
		return Policy{}, err
	}
	maxReserveDeployment, err := numeric(ids.MaximumStablecoinReserveDeploymentSecretID)
	if err != nil {
		return Policy{}, err
	}
	minReserveBalance, err := numeric(ids.MinimumStablecoinReserveBalanceSecretID)
	if err != nil {
		return Policy{}, err
	}
	maxCollateralAllocation, err := numeric(ids.MaximumCollateralAllocationSecretID)
	if err != nil {
		return Policy{}, err
	}
	maxPartialDebtRepayment, err := numeric(ids.MaximumPartialDebtRepaymentSecretID)
	if err != nil {
		return Policy{}, err
	}

	return Policy{
		LiquidationWarningThreshold: warningThreshold,
		MinimumHealthFactor:         minimumHealthFactor,
		TargetHealthFactor:          targetHealthFactor,
		MaxReserveDeploymentUSDC:    maxReserveDeployment,
		MinReserveBalanceUSDC:       minReserveBalance,
		MaxCollateralAllocationPct:  maxCollateralAllocation,
		MaxPartialDebtRepaymentPct:  maxPartialDebtRepayment,
		ExecutionSequencePreference: parseExecutionSequencePreference(sequencingPreference),
		PreferredVenues:             parseVenueListSecret(secretValues[ids.PreferredVenuesSecretID]),
	}, nil
}

// ─── TEE Cron Callback ──────────────────────────────────────
// Receives a cre.TeeRuntime, not a cre.Runtime: the secrets, the risk snapshot,
// the model prompt and the execution request all stay inside the enclave.
func onCronTrigger(config *Config, runtime cre.TeeRuntime, _ *cron.Payload) (string, error) {
	ids := config.SecretsIDs

	// Batch fetch every secret in a single call; if any secret fails, the error
	// reports every failed secret at once. Results are keyed by ID afterwards so
	// no value depends on slice position.
	secrets, err := runtime.GetSecrets([]*cre.SecretRequest{
		{Id: ids.ExchangeAPIKeyID},
		{Id: ids.OpenAIAPIKeyID},
		{Id: ids.LiquidationWarningActionThresholdSecretID},
		{Id: ids.MinimumHealthFactorSecretID},
		{Id: ids.TargetHealthFactorSecretID},
		{Id: ids.MaximumStablecoinReserveDeploymentSecretID},
		{Id: ids.MinimumStablecoinReserveBalanceSecretID},
		{Id: ids.MaximumCollateralAllocationSecretID},
		{Id: ids.MaximumPartialDebtRepaymentSecretID},
		{Id: ids.PreferredVenuesSecretID},
	}).Await()
	if err != nil {
		return "", fmt.Errorf("failed to fetch secrets inside the enclave: %w", err)
	}

	secretValues := make(map[string]string, len(secrets))
	for _, secret := range secrets {
		secretValues[secret.Id] = secret.Value
	}

	exchangeAPIKey := secretValues[ids.ExchangeAPIKeyID]
	openAIAPIKey := secretValues[ids.OpenAIAPIKeyID]

	policy, err := buildPolicy(secretValues, ids, config.DefensiveActionSequencingPreference)
	if err != nil {
		return "", err
	}

	// ⚠️ Logs are for simulation only and MUST be removed before deploying to
	// production — anything logged inside the enclave weakens the
	// confidentiality guarantee. Note this line records only that the fetch
	// succeeded, never a secret value.
	runtime.Logger().Info("liquidation-getsecrets-ok")

	client := &http.Client{}

	risk, err := collectRiskSnapshot(runtime, client, config.MockBaseURL, exchangeAPIKey)
	if err != nil {
		return "", err
	}

	riskScore := ComputeRiskScore(risk, policy)

	prompt, err := createOpenAiPrompt(risk, riskScore, policy)
	if err != nil {
		return "", err
	}

	openAIHeaders := jsonHeaders()
	openAIHeaders["Authorization"] = "Bearer " + openAIAPIKey

	llmResponse, err := postJSON(runtime, client, config.OpenAIURL, map[string]any{
		"model": config.OpenAIModel,
		"input": []map[string]string{
			{
				"role": "system",
				"content": "You are a liquidation-defense policy engine. Emit strict JSON only " +
					"with key names exactly as requested.",
			},
			{"role": "user", "content": prompt},
		},
	}, openAIHeaders)
	if err != nil {
		return "", err
	}

	rawText, err := extractOpenAiText(llmResponse)
	if err != nil {
		return "", err
	}

	decision, err := parseLlmDecision(rawText)
	if err != nil {
		return "", err
	}

	actions, err := EnforcePolicy(decision, policy, risk)
	if err != nil {
		return "", err
	}

	if risk.LiquidationProximityPct > policy.LiquidationWarningThreshold || len(actions) == 0 {
		runtime.Logger().Info("liquidation-no-action",
			"proximity", fmt.Sprintf("%.3f", risk.LiquidationProximityPct),
			"threshold", policy.LiquidationWarningThreshold,
			"reason", decision.Reasoning,
		)
		return "SAFE", nil
	}

	defenseHeaders := jsonHeaders()
	defenseHeaders["x-exchange-api-key"] = exchangeAPIKey

	defenseResponse, err := postJSON(runtime, client, config.MockBaseURL+"/execute-defense", map[string]any{
		"riskScore":               riskScore,
		"liquidationProximityPct": risk.LiquidationProximityPct,
		"reasoning":               decision.Reasoning,
		"actions":                 actions,
	}, defenseHeaders)
	if err != nil {
		return "", err
	}

	executionID := asString(defenseResponse["execution_id"], "unknown")

	runtime.Logger().Info("liquidation-defense-executed",
		"action_count", len(actions),
		"execution_id", executionID,
	)

	result, err := json.Marshal(map[string]any{
		"status":      "DEFENDED",
		"actionCount": len(actions),
		"riskScore":   riskScore,
		"executionId": executionID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to encode result: %w", err)
	}
	return string(result), nil
}

// ─── Workflow Init ──────────────────────────────────────────

func InitWorkflow(config *Config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[*Config], error) {
	if config.Schedule == "" || config.MockBaseURL == "" || config.OpenAIURL == "" ||
		config.OpenAIModel == "" || config.DefensiveActionSequencingPreference == "" {
		return nil, fmt.Errorf("config requires schedule, mock_base_url, openai_url, openai_model, and defensive_action_sequencing_preference")
	}

	ids := config.SecretsIDs
	if ids.ExchangeAPIKeyID == "" ||
		ids.OpenAIAPIKeyID == "" ||
		ids.LiquidationWarningActionThresholdSecretID == "" ||
		ids.MinimumHealthFactorSecretID == "" ||
		ids.TargetHealthFactorSecretID == "" ||
		ids.MaximumStablecoinReserveDeploymentSecretID == "" ||
		ids.MinimumStablecoinReserveBalanceSecretID == "" ||
		ids.MaximumCollateralAllocationSecretID == "" ||
		ids.MaximumPartialDebtRepaymentSecretID == "" ||
		ids.PreferredVenuesSecretID == "" {
		return nil, fmt.Errorf("config requires secrets_ids fields")
	}

	return cre.Workflow[*Config]{
		// cre.HandlerInTee instead of cre.Handler. The third argument declares
		// which enclaves this handler accepts. Alternatives:
		//   cre.AnyTee{}                                               — any registered TEE, any region
		//   cre.AnyTeeInRegions{Regions: []cre.Region{cre.AwsUsWest2}} — any TEE, one region
		cre.HandlerInTee(
			cron.Trigger(&cron.Config{Schedule: config.Schedule}),
			onCronTrigger,
			cre.AnyTee{},
			//cre.OneOfTees{cre.Nitro{Regions: []cre.NitroRegion{cre.NitroUsWest2}}},
		),
	}, nil
}
