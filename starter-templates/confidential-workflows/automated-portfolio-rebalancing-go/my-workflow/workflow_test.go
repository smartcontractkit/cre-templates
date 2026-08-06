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
	mockBaseURL    = "http://127.0.0.1:8787/rebalancing"
	openAIURL      = "http://127.0.0.1:8787/rebalancing/v1/responses"
)

func policyFixture() Policy {
	return Policy{
		TargetAllocations:       map[string]float64{"BTC": 0.5, "ETH": 0.3, "USDC": 0.2},
		DriftThresholdPct:       5,
		MaxTradeUsd:             5000,
		ReserveFloorUsdc:        2000,
		MaxSlippageBps:          50,
		PreferredVenues:         []string{"binance", "onchain"},
		OrderSequencePreference: SequenceModelOrder,
	}
}

func testConfig() *Config {
	return &Config{
		Schedule:    "0 */5 * * * *",
		MockBaseURL: mockBaseURL,
		OpenAIURL:   openAIURL,
		OpenAIModel: "gpt-4.1-mini",
		SecretsIDs: SecretsConfig{
			ExchangeAPIKeyID:                "exchange_api_key",
			OpenAIAPIKeyID:                  "openai_api_key",
			TargetAllocationBtcPctSecretID:  "rebalancing_target_allocation_btc_pct",
			TargetAllocationEthPctSecretID:  "rebalancing_target_allocation_eth_pct",
			TargetAllocationUsdcPctSecretID: "rebalancing_target_allocation_usdc_pct",
			DriftThresholdPctSecretID:       "rebalancing_drift_threshold_pct",
			MaxTradeUsdSecretID:             "rebalancing_max_trade_usd",
			ReserveFloorUsdcSecretID:        "rebalancing_reserve_floor_usdc",
			MaxSlippageBpsSecretID:          "rebalancing_max_slippage_bps",
			PreferredVenuesSecretID:         "rebalancing_preferred_venues",
			OrderSequencePreferenceSecretID: "rebalancing_order_sequence_preference",
		},
	}
}

// testSecrets mirrors the values .env.example seeds, so the policy the enclave
// assembles matches policyFixture().
func testSecrets() testutils.Secrets {
	return testutils.Secrets{
		cre.DefaultSecretNamespace: {
			"exchange_api_key":                       exchangeAPIKey,
			"openai_api_key":                         openAIAPIKey,
			"rebalancing_target_allocation_btc_pct":  "0.5",
			"rebalancing_target_allocation_eth_pct":  "0.3",
			"rebalancing_target_allocation_usdc_pct": "0.2",
			"rebalancing_drift_threshold_pct":        "5",
			"rebalancing_max_trade_usd":              "5000",
			"rebalancing_reserve_floor_usdc":         "2000",
			"rebalancing_max_slippage_bps":           "50",
			"rebalancing_preferred_venues":           "binance,onchain",
			"rebalancing_order_sequence_preference":  "model-order",
		},
	}
}

type capturedCall struct {
	URL     string
	Method  string
	Headers map[string][]string
	Body    string
}

// driftedPortfolioJSON is far enough from the 50/30/20 target to clear the 5%
// drift threshold, so the workflow proceeds to execution. Values: BTC 20000,
// ETH 40000, USDC 20000 (total 80000) against targets 40000/24000/16000, giving
// a 20000 BTC buy and a 16000 ETH sell. usdc_reserve is set high enough that the
// buy leg does not trip the reserve floor.
func driftedPortfolioJSON() string {
	return `{"holdings":{"BTC":0.2,"ETH":20,"USDC":20000},` +
		`"prices_usd":{"BTC":100000,"ETH":2000,"USDC":1},` +
		`"usdc_reserve":60000,"cash_balance_usd":5000,"volatility_index":0.2}`
}

// balancedPortfolioJSON sits exactly on target, so drift is below threshold.
func balancedPortfolioJSON() string {
	return `{"holdings":{"BTC":0.5,"ETH":15,"USDC":20000},` +
		`"prices_usd":{"BTC":100000,"ETH":2000,"USDC":1},` +
		`"usdc_reserve":20000,"cash_balance_usd":5000,"volatility_index":0.2}`
}

func stubEnclaveHTTP(t *testing.T, portfolioJSON string, llmText string) *[]capturedCall {
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
		case strings.HasSuffix(input.Url, "/portfolio-state"):
			return &http.Response{StatusCode: 200, Body: []byte(portfolioJSON)}, nil
		case strings.HasSuffix(input.Url, "/prices"):
			return &http.Response{
				StatusCode: 200,
				Body:       []byte(`{"prices":{"BTC":100000,"ETH":2000,"USDC":1}}`),
			}, nil
		case strings.HasSuffix(input.Url, "/volatility"):
			return &http.Response{StatusCode: 200, Body: []byte(`{"volatility_index":0.3}`)}, nil
		case strings.HasSuffix(input.Url, "/v1/responses"):
			payload, err := json.Marshal(map[string]any{"output_text": llmText})
			require.NoError(t, err)
			return &http.Response{StatusCode: 200, Body: payload}, nil
		case strings.HasSuffix(input.Url, "/execute-rebalance"):
			return &http.Response{StatusCode: 200, Body: []byte(`{"execution_id":"exec-789"}`)}, nil
		}
		return nil, assert.AnError
	}
	return calls
}

// rebalanceDecisionJSON agrees with the target-aligned sides so nothing is vetoed.
func rebalanceDecisionJSON() string {
	return `{"shouldRebalance":true,"reasoning":"drift above threshold","trades":[` +
		`{"symbol":"BTC","side":"buy","notionalUsd":9999999,"maxSlippageBps":40},` +
		`{"symbol":"ETH","side":"sell","notionalUsd":1,"maxSlippageBps":40}]}`
}

// ─── CalculateAllocations ───────────────────────────────────

func TestCalculateAllocations_ComputesDriftAndTotal(t *testing.T) {
	allocations, totalUsd, maxDriftPct := CalculateAllocations(
		Holdings{"BTC": 1, "ETH": 10, "USDC": 1000},
		Prices{"BTC": 50000, "ETH": 3000, "USDC": 1},
		map[string]float64{"BTC": 0.5, "ETH": 0.3, "USDC": 0.2},
	)

	assert.Equal(t, 81000.0, totalUsd)
	require.Len(t, allocations, 3)
	assert.Greater(t, maxDriftPct, 0.0)

	var usdc *Allocation
	for i := range allocations {
		if allocations[i].Symbol == "USDC" {
			usdc = &allocations[i]
		}
	}
	require.NotNil(t, usdc)
	assert.InDelta(t, 15200.0, usdc.DeltaUsdToTarget, 1e-9)
}

// Allocation order feeds trade order, which EnforcePolicyAndChunk walks while
// tracking the projected reserve — so it must not depend on Go's randomised map
// iteration order.
func TestCalculateAllocations_IsDeterministicallyOrdered(t *testing.T) {
	holdings := Holdings{"BTC": 1, "ETH": 10, "USDC": 1000, "SOL": 5, "AVAX": 3}
	prices := Prices{"BTC": 50000, "ETH": 3000, "USDC": 1, "SOL": 200, "AVAX": 40}
	targets := map[string]float64{"BTC": 0.5, "ETH": 0.3, "USDC": 0.2}

	first, _, _ := CalculateAllocations(holdings, prices, targets)
	for i := 0; i < 20; i++ {
		next, _, _ := CalculateAllocations(holdings, prices, targets)
		assert.Equal(t, first, next, "allocation order must be stable across runs")
	}

	symbols := make([]string, len(first))
	for i, allocation := range first {
		symbols[i] = allocation.Symbol
	}
	assert.Equal(t, []string{"AVAX", "BTC", "ETH", "SOL", "USDC"}, symbols)
}

func TestCalculateAllocations_HandlesZeroTotalValue(t *testing.T) {
	allocations, totalUsd, _ := CalculateAllocations(
		Holdings{"BTC": 0},
		Prices{"BTC": 0},
		map[string]float64{"BTC": 1},
	)

	assert.Equal(t, 0.0, totalUsd)
	require.Len(t, allocations, 1)
	assert.Equal(t, 0.0, allocations[0].CurrentWeight, "must not divide by zero")
}

func TestBuildTargetAlignedTrades_DerivesBuysAndSellsFromDrift(t *testing.T) {
	allocations, _, _ := CalculateAllocations(
		Holdings{"BTC": 0.2, "ETH": 20, "USDC": 20000},
		Prices{"BTC": 100000, "ETH": 2000, "USDC": 1},
		map[string]float64{"BTC": 0.4, "ETH": 0.2, "USDC": 0.4},
	)

	trades := BuildTargetAlignedTrades(allocations)
	require.Len(t, trades, 2)

	bySymbol := map[string]RebalanceTrade{}
	for _, trade := range trades {
		bySymbol[trade.Symbol] = trade
	}
	assert.Equal(t, SideBuy, bySymbol["BTC"].Side)
	assert.Equal(t, SideSell, bySymbol["ETH"].Side)
}

func TestBuildTargetAlignedTrades_SkipsUSDCAndDustDeltas(t *testing.T) {
	trades := BuildTargetAlignedTrades([]Allocation{
		{Symbol: "USDC", DeltaUsdToTarget: 10000},
		{Symbol: "BTC", DeltaUsdToTarget: 0.5},
		{Symbol: "ETH", DeltaUsdToTarget: 250},
	})

	require.Len(t, trades, 1)
	assert.Equal(t, "ETH", trades[0].Symbol)
}

// ─── EnforcePolicyAndChunk ──────────────────────────────────

func TestEnforcePolicyAndChunk_ChunksByMaxSizeAndAssignsVenue(t *testing.T) {
	chunks, err := EnforcePolicyAndChunk(LlmDecision{
		ShouldRebalance: true,
		Reasoning:       "drift above threshold",
		Trades: []RebalanceTrade{
			{Symbol: "BTC", Side: SideBuy, NotionalUsd: 12000, MaxSlippageBps: 40},
		},
	}, policyFixture(), 20000)
	require.NoError(t, err)

	require.Len(t, chunks, 3)
	assert.Equal(t, 5000.0, chunks[0].NotionalUsd)
	assert.Equal(t, 5000.0, chunks[1].NotionalUsd)
	assert.Equal(t, 2000.0, chunks[2].NotionalUsd)
	assert.Equal(t, "binance", chunks[0].Venue)
	assert.Equal(t, 3, chunks[0].ChunkCount)
	assert.Equal(t, 1, chunks[0].ChunkIndex)
}

func TestEnforcePolicyAndChunk_RejectsReserveFloorBreach(t *testing.T) {
	_, err := EnforcePolicyAndChunk(LlmDecision{
		ShouldRebalance: true,
		Reasoning:       "buy signal",
		Trades:          []RebalanceTrade{{Symbol: "ETH", Side: SideBuy, NotionalUsd: 4000}},
	}, policyFixture(), 5000)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "breaches reserve floor")
}

func TestEnforcePolicyAndChunk_RejectsExcessiveSlippage(t *testing.T) {
	_, err := EnforcePolicyAndChunk(LlmDecision{
		ShouldRebalance: true,
		Trades:          []RebalanceTrade{{Symbol: "ETH", Side: SideSell, NotionalUsd: 100, MaxSlippageBps: 500}},
	}, policyFixture(), 20000)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds max")
}

func TestEnforcePolicyAndChunk_SupportsSellsFirstSequencing(t *testing.T) {
	policy := policyFixture()
	policy.OrderSequencePreference = SequenceSellsFirst

	chunks, err := EnforcePolicyAndChunk(LlmDecision{
		ShouldRebalance: true,
		Reasoning:       "drift above threshold",
		Trades: []RebalanceTrade{
			{Symbol: "ETH", Side: SideBuy, NotionalUsd: 1000},
			{Symbol: "BTC", Side: SideSell, NotionalUsd: 1000},
		},
	}, policy, 10000)
	require.NoError(t, err)

	require.Len(t, chunks, 2)
	assert.Equal(t, SideSell, chunks[0].Side)
}

func TestEnforcePolicyAndChunk_SupportsBuysFirstSequencing(t *testing.T) {
	policy := policyFixture()
	policy.OrderSequencePreference = SequenceBuysFirst

	chunks, err := EnforcePolicyAndChunk(LlmDecision{
		ShouldRebalance: true,
		Trades: []RebalanceTrade{
			{Symbol: "BTC", Side: SideSell, NotionalUsd: 1000},
			{Symbol: "ETH", Side: SideBuy, NotionalUsd: 1000},
		},
	}, policy, 10000)
	require.NoError(t, err)

	require.Len(t, chunks, 2)
	assert.Equal(t, SideBuy, chunks[0].Side)
}

func TestEnforcePolicyAndChunk_ReturnsNothingWhenModelDeclines(t *testing.T) {
	chunks, err := EnforcePolicyAndChunk(LlmDecision{
		ShouldRebalance: false,
		Trades:          []RebalanceTrade{{Symbol: "BTC", Side: SideBuy, NotionalUsd: 100}},
	}, policyFixture(), 20000)

	require.NoError(t, err)
	assert.Empty(t, chunks)
}

// ─── Reconciliation ─────────────────────────────────────────

// The notional must come from the deterministic target-aligned trade, never from
// the model — otherwise a hallucinated size could move far more capital than
// closing the drift requires.
func TestReconcileDecisionWithTargets_IgnoresModelNotional(t *testing.T) {
	targetTrades := []RebalanceTrade{{Symbol: "ETH", Side: SideSell, NotionalUsd: 3000}}
	decision := LlmDecision{
		ShouldRebalance: true,
		Trades:          []RebalanceTrade{{Symbol: "ETH", Side: SideSell, NotionalUsd: 999999}},
	}
	metrics := PortfolioMetrics{MaxDriftPct: 10, ReserveBufferUsd: 5000}

	reconciled := ReconcileDecisionWithTargets(decision, targetTrades, policyFixture(), metrics)

	require.Len(t, reconciled.Trades, 1)
	assert.Equal(t, 3000.0, reconciled.Trades[0].NotionalUsd)
}

func TestReconcileDecisionWithTargets_VetoesOnSideDisagreement(t *testing.T) {
	targetTrades := []RebalanceTrade{{Symbol: "ETH", Side: SideSell, NotionalUsd: 3000}}
	decision := LlmDecision{
		ShouldRebalance: true,
		Trades:          []RebalanceTrade{{Symbol: "ETH", Side: SideBuy, NotionalUsd: 3000}},
	}
	metrics := PortfolioMetrics{MaxDriftPct: 10, ReserveBufferUsd: 5000}

	reconciled := ReconcileDecisionWithTargets(decision, targetTrades, policyFixture(), metrics)

	assert.False(t, reconciled.ShouldRebalance)
	assert.Empty(t, reconciled.Trades)
}

func TestReconcileDecisionWithTargets_ClampsSlippageToPolicy(t *testing.T) {
	targetTrades := []RebalanceTrade{{Symbol: "ETH", Side: SideSell, NotionalUsd: 3000}}
	decision := LlmDecision{
		ShouldRebalance: true,
		Trades: []RebalanceTrade{
			{Symbol: "ETH", Side: SideSell, NotionalUsd: 3000, MaxSlippageBps: 9999},
		},
	}
	metrics := PortfolioMetrics{MaxDriftPct: 10, ReserveBufferUsd: 5000}

	reconciled := ReconcileDecisionWithTargets(decision, targetTrades, policyFixture(), metrics)

	require.Len(t, reconciled.Trades, 1)
	assert.Equal(t, 50.0, reconciled.Trades[0].MaxSlippageBps)
}

func TestReconcileDecisionWithTargets_StandsDownBelowDriftThreshold(t *testing.T) {
	targetTrades := []RebalanceTrade{{Symbol: "ETH", Side: SideSell, NotionalUsd: 3000}}
	metrics := PortfolioMetrics{MaxDriftPct: 1, ReserveBufferUsd: 5000}

	reconciled := ReconcileDecisionWithTargets(LlmDecision{ShouldRebalance: true}, targetTrades, policyFixture(), metrics)

	assert.False(t, reconciled.ShouldRebalance)
	assert.Contains(t, reconciled.Reasoning, "below threshold")
}

// Buys are limited by the reserve buffer plus whatever the planned sells free up.
func TestReconcileDecisionWithTargets_CapsBuysByAvailableCapacity(t *testing.T) {
	targetTrades := []RebalanceTrade{{Symbol: "BTC", Side: SideBuy, NotionalUsd: 10000}}
	metrics := PortfolioMetrics{MaxDriftPct: 10, ReserveBufferUsd: 2500}

	reconciled := ReconcileDecisionWithTargets(LlmDecision{ShouldRebalance: true}, targetTrades, policyFixture(), metrics)

	require.Len(t, reconciled.Trades, 1)
	assert.Equal(t, 2500.0, reconciled.Trades[0].NotionalUsd)
}

// ─── Execution plan ─────────────────────────────────────────

func TestBuildExecutionPlan_SplitsByVenue(t *testing.T) {
	plan := BuildExecutionPlan([]ExecutableTrade{
		{
			RebalanceTrade: RebalanceTrade{Symbol: "BTC", Side: SideSell, NotionalUsd: 1000, MaxSlippageBps: 20},
			Venue:          "binance", ChunkIndex: 1, ChunkCount: 1,
		},
		{
			RebalanceTrade: RebalanceTrade{Symbol: "ETH", Side: SideBuy, NotionalUsd: 500, MaxSlippageBps: 20},
			Venue:          "onchain", ChunkIndex: 1, ChunkCount: 1,
		},
	})

	require.Len(t, plan.OnchainTrades, 1)
	require.Len(t, plan.OffchainTrades, 1)
	assert.Equal(t, "onchain", plan.OnchainTrades[0].Venue)
	assert.Equal(t, "binance", plan.OffchainTrades[0].Venue)
}

func TestClassifyExecutionVenue_RecognisesDexNames(t *testing.T) {
	assert.Equal(t, "onchain", classifyExecutionVenue("Uniswap V3"))
	assert.Equal(t, "onchain", classifyExecutionVenue("1inch"))
	assert.Equal(t, "onchain", classifyExecutionVenue("some-dex"))
	assert.Equal(t, "offchain", classifyExecutionVenue("coinbase"))
}

// ─── Untrusted model output ─────────────────────────────────

func TestParseLlmDecision_DropsTradesWithoutSymbolOrNotional(t *testing.T) {
	decision, err := parseLlmDecision(
		`{"shouldRebalance":true,"trades":[{"symbol":"","notionalUsd":100},` +
			`{"symbol":"BTC","notionalUsd":0},{"symbol":"eth","side":"SELL","notionalUsd":50}]}`,
	)
	require.NoError(t, err)

	require.Len(t, decision.Trades, 1)
	assert.Equal(t, "ETH", decision.Trades[0].Symbol, "symbols are upper-cased")
	assert.Equal(t, SideSell, decision.Trades[0].Side, "sides are lower-cased")
}

func TestParseLlmDecision_DefaultsSideToBuy(t *testing.T) {
	decision, err := parseLlmDecision(`{"shouldRebalance":true,"trades":[{"symbol":"BTC","notionalUsd":10}]}`)
	require.NoError(t, err)

	require.Len(t, decision.Trades, 1)
	assert.Equal(t, SideBuy, decision.Trades[0].Side)
}

func TestParseLlmDecision_ErrorsOnInvalidJSON(t *testing.T) {
	_, err := parseLlmDecision("not json")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid json response")
}

func TestParseOrderSequencePreference_FallsBackToModelOrder(t *testing.T) {
	assert.Equal(t, SequenceSellsFirst, parseOrderSequencePreference("sells-first"))
	assert.Equal(t, SequenceBuysFirst, parseOrderSequencePreference("buys-first"))
	assert.Equal(t, SequenceModelOrder, parseOrderSequencePreference("nonsense"))
}

// An infinite coverage ratio (zero reserve floor) must still encode, since it
// goes into the model prompt. Go's encoding/json rejects Inf by default.
func TestPortfolioMetrics_EncodesInfiniteCoverageRatioAsNull(t *testing.T) {
	snapshot := MarketSnapshot{
		Portfolio: PortfolioState{Holdings: Holdings{}, UsdcReserve: 100},
		Prices:    Prices{},
	}
	policy := policyFixture()
	policy.ReserveFloorUsdc = 0

	metrics := computePortfolioMetrics(snapshot, policy, 1000, 5)

	encoded, err := json.Marshal(metrics)
	require.NoError(t, err)
	assert.Contains(t, string(encoded), `"reserveCoverageRatio":null`)
}

// ─── End-to-end through the enclave runtime ─────────────────

func TestOnCronTrigger_ExecutesWhenDriftExceedsThreshold(t *testing.T) {
	stubEnclaveHTTP(t, driftedPortfolioJSON(), rebalanceDecisionJSON())
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	result, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal([]byte(result), &decoded))
	assert.Equal(t, "EXECUTED", decoded["status"])
	assert.Equal(t, "exec-789", decoded["executionId"])
	assert.Greater(t, decoded["tradeCount"], 0.0)
}

func TestOnCronTrigger_ReturnsNoopWhenPortfolioOnTarget(t *testing.T) {
	stubEnclaveHTTP(t, balancedPortfolioJSON(), rebalanceDecisionJSON())
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	result, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)
	assert.Equal(t, "NOOP", result)
}

// The model's shouldRebalance flag is deliberately NOT a veto: the trades are
// derived from the deterministic target-aligned set, so an empty/declining
// response still rebalances. The model's only veto is disagreeing on a side.
func TestOnCronTrigger_ModelDeclineFlagDoesNotOverrideTargets(t *testing.T) {
	stubEnclaveHTTP(t, driftedPortfolioJSON(), `{"shouldRebalance":false,"trades":[]}`)
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	result, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)
	assert.Contains(t, result, "EXECUTED")
}

// Disagreeing on every side vetoes every trade, which is what actually produces
// a NOOP from a non-empty target set.
func TestOnCronTrigger_ReturnsNoopWhenModelDisagreesOnEverySide(t *testing.T) {
	stubEnclaveHTTP(t, driftedPortfolioJSON(),
		`{"shouldRebalance":true,"trades":[`+
			`{"symbol":"BTC","side":"sell","notionalUsd":10},`+
			`{"symbol":"ETH","side":"buy","notionalUsd":10}]}`)
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	result, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)
	assert.Equal(t, "NOOP", result)
}

func TestOnCronTrigger_InjectsEnclaveFetchedSecretsIntoRequests(t *testing.T) {
	calls := stubEnclaveHTTP(t, driftedPortfolioJSON(), rebalanceDecisionJSON())
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)

	// portfolio-state, prices, volatility, openai, execute-rebalance
	require.Len(t, *calls, 5)

	for _, call := range (*calls)[:3] {
		assert.Equal(t, "GET", call.Method)
		assert.Equal(t, []string{exchangeAPIKey}, call.Headers["x-exchange-api-key"])
	}
	assert.Equal(t, []string{"Bearer " + openAIAPIKey}, (*calls)[3].Headers["Authorization"])
	assert.Equal(t, []string{exchangeAPIKey}, (*calls)[4].Headers["x-exchange-api-key"])
}

func TestOnCronTrigger_SendsConfidentialPolicyToModel(t *testing.T) {
	calls := stubEnclaveHTTP(t, driftedPortfolioJSON(), rebalanceDecisionJSON())
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)

	llmBody := (*calls)[3].Body
	assert.Contains(t, llmBody, "target_allocations")
	assert.Contains(t, llmBody, "targetAlignedTrades")
}

func TestOnCronTrigger_ErrorsOnNon2xx(t *testing.T) {
	capability, err := httpmock.NewClientCapability(t)
	require.NoError(t, err)
	capability.SendRequest = func(_ context.Context, _ *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 503, Body: []byte("unavailable")}, nil
	}

	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err = onCronTrigger(testConfig(), runtime, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "status=503")
}

func TestOnCronTrigger_ErrorsWhenPolicySecretIsNotNumeric(t *testing.T) {
	stubEnclaveHTTP(t, driftedPortfolioJSON(), rebalanceDecisionJSON())

	secrets := testSecrets()
	secrets[cre.DefaultSecretNamespace]["rebalancing_max_trade_usd"] = "lots"
	runtime := testutils.NewTeeRuntime(t, secrets)

	_, err := onCronTrigger(testConfig(), runtime, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be a finite number")
}

// Logging inside an enclave weakens its confidentiality guarantee. The workflow
// logs only non-sensitive markers, never secrets.
func TestOnCronTrigger_DoesNotLogSecrets(t *testing.T) {
	stubEnclaveHTTP(t, driftedPortfolioJSON(), rebalanceDecisionJSON())
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
	assert.Contains(t, err.Error(), "config requires schedule, mock_base_url, openai_url, and openai_model")
}

func TestInitWorkflow_ErrorsWhenSecretIDsMissing(t *testing.T) {
	config := testConfig()
	config.SecretsIDs.OrderSequencePreferenceSecretID = ""

	_, err := InitWorkflow(config, nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "config requires secrets_ids fields")
}
