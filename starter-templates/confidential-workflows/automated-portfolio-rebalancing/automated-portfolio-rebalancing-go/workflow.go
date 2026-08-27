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

// SecretsConfig maps each policy input to the secret ID that holds it. Target
// weights, drift tolerance and trade caps are all confidential: publishing them
// would let an adversary front-run every rebalance this portfolio performs.
type SecretsConfig struct {
	ExchangeAPIKeyID                string `json:"exchange_api_key_id"`
	OpenAIAPIKeyID                  string `json:"openai_api_key_id"`
	TargetAllocationBtcPctSecretID  string `json:"target_allocation_btc_pct_secret_id"`
	TargetAllocationEthPctSecretID  string `json:"target_allocation_eth_pct_secret_id"`
	TargetAllocationUsdcPctSecretID string `json:"target_allocation_usdc_pct_secret_id"`
	DriftThresholdPctSecretID       string `json:"drift_threshold_pct_secret_id"`
	MaxTradeUsdSecretID             string `json:"max_trade_usd_secret_id"`
	ReserveFloorUsdcSecretID        string `json:"reserve_floor_usdc_secret_id"`
	MaxSlippageBpsSecretID          string `json:"max_slippage_bps_secret_id"`
	PreferredVenuesSecretID         string `json:"preferred_venues_secret_id"`
	OrderSequencePreferenceSecretID string `json:"order_sequence_preference_secret_id"`
}

type Config struct {
	Schedule    string        `json:"schedule"`
	MockBaseURL string        `json:"mock_base_url"`
	OpenAIURL   string        `json:"openai_url"`
	OpenAIModel string        `json:"openai_model"`
	SecretsIDs  SecretsConfig `json:"secrets_ids"`
}

// ─── Domain types ───────────────────────────────────────────

type Holdings map[string]float64
type Prices map[string]float64

// Order sequence preferences.
const (
	SequenceModelOrder = "model-order"
	SequenceSellsFirst = "sells-first"
	SequenceBuysFirst  = "buys-first"
)

// Trade sides.
const (
	SideBuy  = "buy"
	SideSell = "sell"
)

// Policy is assembled inside the enclave from Vault DON secrets and never
// leaves it — only the resulting trades do.
type Policy struct {
	TargetAllocations       map[string]float64 `json:"target_allocations"`
	DriftThresholdPct       float64            `json:"drift_threshold_pct"`
	MaxTradeUsd             float64            `json:"max_trade_usd"`
	ReserveFloorUsdc        float64            `json:"reserve_floor_usdc"`
	MaxSlippageBps          float64            `json:"max_slippage_bps"`
	PreferredVenues         []string           `json:"preferred_venues"`
	OrderSequencePreference string             `json:"order_sequence_preference"`
}

type PortfolioState struct {
	Holdings        Holdings `json:"holdings"`
	PricesUsd       Prices   `json:"prices_usd"`
	UsdcReserve     float64  `json:"usdc_reserve"`
	CashBalanceUsd  float64  `json:"cash_balance_usd"`
	VolatilityIndex float64  `json:"volatility_index"`
}

type MarketSnapshot struct {
	Portfolio       PortfolioState `json:"portfolio"`
	Prices          Prices         `json:"prices"`
	VolatilityIndex float64        `json:"volatilityIndex"`
}

type Allocation struct {
	Symbol           string  `json:"symbol"`
	UsdValue         float64 `json:"usdValue"`
	CurrentWeight    float64 `json:"currentWeight"`
	TargetWeight     float64 `json:"targetWeight"`
	DriftPct         float64 `json:"driftPct"`
	DeltaUsdToTarget float64 `json:"deltaUsdToTarget"`
}

// jsNumber marshals like a JavaScript number: non-finite values become null
// rather than failing the encode. ReserveCoverageRatio is +Inf when the reserve
// floor is zero, and Go's encoding/json rejects Inf outright.
type jsNumber float64

func (n jsNumber) MarshalJSON() ([]byte, error) {
	f := float64(n)
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return []byte("null"), nil
	}
	return json.Marshal(f)
}

type PortfolioMetrics struct {
	TotalUsd             float64  `json:"totalUsd"`
	MaxDriftPct          float64  `json:"maxDriftPct"`
	ReserveBufferUsd     float64  `json:"reserveBufferUsd"`
	ReserveCoverageRatio jsNumber `json:"reserveCoverageRatio"`
	StablecoinDepthUsd   float64  `json:"stablecoinDepthUsd"`
	AverageVolatility    float64  `json:"averageVolatility"`
}

type RebalanceTrade struct {
	Symbol          string   `json:"symbol"`
	Side            string   `json:"side"`
	NotionalUsd     float64  `json:"notionalUsd"`
	MaxSlippageBps  float64  `json:"maxSlippageBps"`
	PreferredVenues []string `json:"preferredVenues,omitempty"`
}

// ExecutableTrade is a RebalanceTrade that has passed policy enforcement and
// been split into at most max_trade_usd sized chunks.
type ExecutableTrade struct {
	RebalanceTrade
	Venue      string `json:"venue"`
	ChunkIndex int    `json:"chunkIndex"`
	ChunkCount int    `json:"chunkCount"`
}

type LlmDecision struct {
	ShouldRebalance bool
	Reasoning       string
	Trades          []RebalanceTrade
}

type ExecutionPlan struct {
	OnchainTrades  []ExecutableTrade
	OffchainTrades []ExecutableTrade
}

// ─── Lenient JSON coercion ──────────────────────────────────
// The model's output is untrusted: it may emit numbers as strings, omit keys, or
// return the wrong shape entirely. These helpers coerce rather than fail so a
// sloppy-but-usable response still produces a decision.

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

// ─── Policy inputs from secrets ─────────────────────────────

func parseRequiredSecretNumber(secretValue, secretID string) (float64, error) {
	parsed, err := strconv.ParseFloat(strings.TrimSpace(secretValue), 64)
	if err != nil || math.IsInf(parsed, 0) || math.IsNaN(parsed) {
		return 0, fmt.Errorf("secret %s must be a finite number", secretID)
	}
	return parsed, nil
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

// parseOrderSequencePreference falls back to model-order for anything
// unrecognised, so a corrupted secret cannot invent a new ordering.
func parseOrderSequencePreference(value string) string {
	switch strings.TrimSpace(value) {
	case SequenceSellsFirst:
		return SequenceSellsFirst
	case SequenceBuysFirst:
		return SequenceBuysFirst
	default:
		return SequenceModelOrder
	}
}

func parseNumberMap(payload any) map[string]float64 {
	out := map[string]float64{}
	for symbol, value := range asObject(payload) {
		out[symbol] = asNumber(value, 0)
	}
	return out
}

// ─── Allocation maths ───────────────────────────────────────

// CalculateAllocations values every position, derives its drift from target and
// the USD delta needed to close that drift.
//
// Symbols are sorted so the returned slice is deterministic. This matters beyond
// tidiness: allocation order flows through to trade order, and EnforcePolicyAndChunk
// walks trades in order while tracking the projected reserve — so a different
// order can produce a different reserve-floor outcome. Go map iteration order is
// randomised, and the enclave result is attested and verified by DON consensus,
// so an unordered walk here would be non-deterministic across nodes.
func CalculateAllocations(
	holdings Holdings,
	prices Prices,
	targets map[string]float64,
) (allocations []Allocation, totalUsd float64, maxDriftPct float64) {
	seen := map[string]bool{}
	symbols := []string{}
	for symbol := range targets {
		if !seen[symbol] {
			seen[symbol] = true
			symbols = append(symbols, symbol)
		}
	}
	for symbol := range holdings {
		if !seen[symbol] {
			seen[symbol] = true
			symbols = append(symbols, symbol)
		}
	}
	sort.Strings(symbols)

	usdBySymbol := map[string]float64{}
	for _, symbol := range symbols {
		usdBySymbol[symbol] = holdings[symbol] * prices[symbol]
		totalUsd += usdBySymbol[symbol]
	}

	safeTotal := totalUsd
	if safeTotal <= 0 {
		safeTotal = 1
	}

	allocations = make([]Allocation, 0, len(symbols))
	for _, symbol := range symbols {
		targetWeight := targets[symbol]
		usdValue := usdBySymbol[symbol]
		currentWeight := usdValue / safeTotal
		driftPct := math.Abs(currentWeight-targetWeight) * 100

		allocations = append(allocations, Allocation{
			Symbol:           symbol,
			UsdValue:         usdValue,
			CurrentWeight:    currentWeight,
			TargetWeight:     targetWeight,
			DriftPct:         driftPct,
			DeltaUsdToTarget: targetWeight*totalUsd - usdValue,
		})

		maxDriftPct = math.Max(maxDriftPct, driftPct)
	}

	return allocations, totalUsd, maxDriftPct
}

func computePortfolioMetrics(
	snapshot MarketSnapshot,
	policy Policy,
	totalUsd float64,
	maxDriftPct float64,
) PortfolioMetrics {
	usdcPrice := 1.0
	if price, ok := snapshot.Prices["USDC"]; ok {
		usdcPrice = price
	}
	usdcHoldingsUsd := snapshot.Portfolio.Holdings["USDC"] * usdcPrice

	reserveCoverageRatio := math.Inf(1)
	if policy.ReserveFloorUsdc > 0 {
		reserveCoverageRatio = snapshot.Portfolio.UsdcReserve / policy.ReserveFloorUsdc
	}

	return PortfolioMetrics{
		TotalUsd:             totalUsd,
		MaxDriftPct:          maxDriftPct,
		ReserveBufferUsd:     snapshot.Portfolio.UsdcReserve - policy.ReserveFloorUsdc,
		ReserveCoverageRatio: jsNumber(reserveCoverageRatio),
		StablecoinDepthUsd:   usdcHoldingsUsd + snapshot.Portfolio.CashBalanceUsd,
		AverageVolatility:    (snapshot.Portfolio.VolatilityIndex + snapshot.VolatilityIndex) / 2,
	}
}

func isTradableRiskAsset(symbol string) bool { return symbol != "USDC" }

// BuildTargetAlignedTrades derives the trades that would close the drift, from
// the allocations alone. These are the ground truth the model's suggestions get
// reconciled against — the model cannot invent a symbol or flip a side.
func BuildTargetAlignedTrades(allocations []Allocation) []RebalanceTrade {
	trades := []RebalanceTrade{}
	for _, allocation := range allocations {
		if !isTradableRiskAsset(allocation.Symbol) {
			continue
		}
		if math.Abs(allocation.DeltaUsdToTarget) < 1 {
			continue
		}

		side := SideSell
		if allocation.DeltaUsdToTarget > 0 {
			side = SideBuy
		}

		trades = append(trades, RebalanceTrade{
			Symbol:      allocation.Symbol,
			Side:        side,
			NotionalUsd: math.Abs(allocation.DeltaUsdToTarget),
		})
	}
	return trades
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

func parseLlmDecision(rawText string) (LlmDecision, error) {
	parsed, err := parseJSONObject(rawText)
	if err != nil {
		return LlmDecision{}, err
	}

	shouldRebalance, _ := parsed["shouldRebalance"].(bool)
	tradesRaw, _ := parsed["trades"].([]any)

	trades := []RebalanceTrade{}
	for _, item := range tradesRaw {
		row := asObject(item)

		side := SideBuy
		if strings.ToLower(asString(row["side"], "")) == SideSell {
			side = SideSell
		}

		trade := RebalanceTrade{
			Symbol:          strings.ToUpper(asString(row["symbol"], "")),
			Side:            side,
			NotionalUsd:     asNumber(row["notionalUsd"], 0),
			MaxSlippageBps:  asNumber(row["maxSlippageBps"], 0),
			PreferredVenues: asStringSlice(row["preferredVenues"]),
		}

		if trade.Symbol == "" || trade.NotionalUsd <= 0 {
			continue
		}
		trades = append(trades, trade)
	}

	return LlmDecision{
		ShouldRebalance: shouldRebalance,
		Reasoning:       asString(parsed["reasoning"], "No reasoning provided by model"),
		Trades:          trades,
	}, nil
}

func createOpenAiPrompt(
	allocations []Allocation,
	metrics PortfolioMetrics,
	policy Policy,
	snapshot MarketSnapshot,
) (string, error) {
	prompt := struct {
		Objective           string           `json:"objective"`
		Policy              Policy           `json:"policy"`
		Metrics             PortfolioMetrics `json:"metrics"`
		Allocations         []Allocation     `json:"allocations"`
		TargetAlignedTrades []RebalanceTrade `json:"targetAlignedTrades"`
		Snapshot            MarketSnapshot   `json:"snapshot"`
	}{
		Objective: "Return strict JSON with shouldRebalance:boolean, reasoning:string, " +
			"trades:[{symbol, side, notionalUsd, maxSlippageBps, preferredVenues}] to restore " +
			"portfolio target weights while respecting policy. Use only the provided " +
			"targetAlignedTrades symbols and sides, and keep notionalUsd within those bounds.",
		Policy:              policy,
		Metrics:             metrics,
		Allocations:         allocations,
		TargetAlignedTrades: BuildTargetAlignedTrades(allocations),
		Snapshot:            snapshot,
	}

	encoded, err := json.MarshalIndent(prompt, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to encode prompt: %w", err)
	}
	return string(encoded), nil
}

// ─── Reconciliation and policy enforcement ──────────────────

func chooseVenue(trade RebalanceTrade, preferredVenues []string) string {
	if len(trade.PreferredVenues) > 0 {
		for _, venue := range trade.PreferredVenues {
			for _, allowed := range preferredVenues {
				if venue == allowed {
					return venue
				}
			}
		}
		return trade.PreferredVenues[0]
	}
	if len(preferredVenues) > 0 {
		return preferredVenues[0]
	}
	return "default-venue"
}

// ReconcileDecisionWithTargets is the guardrail between the model and the
// executor. The notionals come from the deterministic target-aligned trades, not
// from the model — the model only gets to influence slippage and venue, and to
// veto a trade by disagreeing on its side. That way a hallucinated size cannot
// move more capital than closing the drift requires.
func ReconcileDecisionWithTargets(
	decision LlmDecision,
	targetTrades []RebalanceTrade,
	policy Policy,
	metrics PortfolioMetrics,
) LlmDecision {
	if metrics.MaxDriftPct < policy.DriftThresholdPct || len(targetTrades) == 0 {
		return LlmDecision{
			ShouldRebalance: false,
			Reasoning: fmt.Sprintf("Drift %.3f is below threshold %v",
				metrics.MaxDriftPct, policy.DriftThresholdPct),
			Trades: []RebalanceTrade{},
		}
	}

	llmTradesBySymbol := map[string]RebalanceTrade{}
	for _, trade := range decision.Trades {
		llmTradesBySymbol[trade.Symbol] = trade
	}

	totalPlannedSellUsd := 0.0
	for _, trade := range targetTrades {
		if trade.Side == SideSell {
			totalPlannedSellUsd += trade.NotionalUsd
		}
	}
	remainingBuyCapacityUsd := math.Max(metrics.ReserveBufferUsd+totalPlannedSellUsd, 0)

	reconciled := []RebalanceTrade{}
	for _, targetTrade := range targetTrades {
		llmTrade, hasLlmTrade := llmTradesBySymbol[targetTrade.Symbol]

		// The model disagreeing on direction vetoes the trade entirely.
		if hasLlmTrade && llmTrade.Side != targetTrade.Side {
			continue
		}

		preferredVenues := policy.PreferredVenues
		if hasLlmTrade && len(llmTrade.PreferredVenues) > 0 {
			preferredVenues = llmTrade.PreferredVenues
		}

		slippage := policy.MaxSlippageBps
		if hasLlmTrade {
			slippage = llmTrade.MaxSlippageBps
		}
		slippage = math.Min(math.Max(slippage, 0), policy.MaxSlippageBps)

		cappedNotional := targetTrade.NotionalUsd
		if targetTrade.Side == SideBuy {
			cappedNotional = math.Min(cappedNotional, remainingBuyCapacityUsd)
			remainingBuyCapacityUsd -= cappedNotional
		}

		if cappedNotional < 1 {
			continue
		}

		reconciled = append(reconciled, RebalanceTrade{
			Symbol:          targetTrade.Symbol,
			Side:            targetTrade.Side,
			NotionalUsd:     cappedNotional,
			MaxSlippageBps:  slippage,
			PreferredVenues: preferredVenues,
		})
	}

	return LlmDecision{
		ShouldRebalance: len(reconciled) > 0,
		Reasoning:       decision.Reasoning,
		Trades:          reconciled,
	}
}

func classifyExecutionVenue(venue string) string {
	normalized := strings.ToLower(strings.TrimSpace(venue))
	if normalized == "onchain" ||
		strings.Contains(normalized, "uniswap") ||
		strings.Contains(normalized, "1inch") ||
		strings.Contains(normalized, "dex") {
		return "onchain"
	}
	return "offchain"
}

// BuildExecutionPlan splits chunks by settlement domain so the executor can
// route on-chain swaps and venue orders down different paths.
func BuildExecutionPlan(trades []ExecutableTrade) ExecutionPlan {
	plan := ExecutionPlan{
		OnchainTrades:  []ExecutableTrade{},
		OffchainTrades: []ExecutableTrade{},
	}
	for _, trade := range trades {
		if classifyExecutionVenue(trade.Venue) == "onchain" {
			plan.OnchainTrades = append(plan.OnchainTrades, trade)
		} else {
			plan.OffchainTrades = append(plan.OffchainTrades, trade)
		}
	}
	return plan
}

// EnforcePolicyAndChunk applies the hard caps and splits each trade into chunks
// no larger than max_trade_usd. It walks trades in order tracking the projected
// reserve, so a sequence that would dip below the floor aborts the whole run
// rather than executing partially.
func EnforcePolicyAndChunk(
	decision LlmDecision,
	policy Policy,
	usdcReserve float64,
) ([]ExecutableTrade, error) {
	if !decision.ShouldRebalance || len(decision.Trades) == 0 {
		return []ExecutableTrade{}, nil
	}

	orderedTrades := make([]RebalanceTrade, len(decision.Trades))
	copy(orderedTrades, decision.Trades)

	switch policy.OrderSequencePreference {
	case SequenceSellsFirst:
		sort.SliceStable(orderedTrades, func(i, j int) bool {
			return orderedTrades[i].Side == SideSell && orderedTrades[j].Side != SideSell
		})
	case SequenceBuysFirst:
		sort.SliceStable(orderedTrades, func(i, j int) bool {
			return orderedTrades[i].Side == SideBuy && orderedTrades[j].Side != SideBuy
		})
	}

	projectedReserve := usdcReserve
	chunks := []ExecutableTrade{}

	for _, trade := range orderedTrades {
		slippage := trade.MaxSlippageBps
		if slippage > policy.MaxSlippageBps {
			return nil, fmt.Errorf("trade %s slippage %v exceeds max %v",
				trade.Symbol, slippage, policy.MaxSlippageBps)
		}

		if trade.Side == SideBuy {
			projectedReserve -= trade.NotionalUsd
		} else {
			projectedReserve += trade.NotionalUsd
		}

		if projectedReserve < policy.ReserveFloorUsdc {
			return nil, fmt.Errorf("trade %s breaches reserve floor: projected %.2f < floor %v",
				trade.Symbol, projectedReserve, policy.ReserveFloorUsdc)
		}

		chunkCount := int(math.Ceil(trade.NotionalUsd / policy.MaxTradeUsd))
		remainder := trade.NotionalUsd
		for idx := 0; idx < chunkCount; idx++ {
			chunkSize := math.Min(remainder, policy.MaxTradeUsd)
			remainder -= chunkSize

			chunked := trade
			chunked.NotionalUsd = chunkSize
			chunked.MaxSlippageBps = slippage

			chunks = append(chunks, ExecutableTrade{
				RebalanceTrade: chunked,
				Venue:          chooseVenue(trade, policy.PreferredVenues),
				ChunkIndex:     idx + 1,
				ChunkCount:     chunkCount,
			})
		}
	}

	return chunks, nil
}

// ─── Enclave data collection ────────────────────────────────

func collectMarketSnapshot(
	runtime cre.TeeRuntime,
	client *http.Client,
	baseURL string,
	exchangeAPIKey string,
) (MarketSnapshot, error) {
	authHeaders := jsonHeaders()
	authHeaders["x-exchange-api-key"] = exchangeAPIKey

	portfolioResponse, err := getJSON(runtime, client, baseURL+"/portfolio-state", authHeaders)
	if err != nil {
		return MarketSnapshot{}, err
	}
	pricesResponse, err := getJSON(runtime, client, baseURL+"/prices", authHeaders)
	if err != nil {
		return MarketSnapshot{}, err
	}
	volatilityResponse, err := getJSON(runtime, client, baseURL+"/volatility", authHeaders)
	if err != nil {
		return MarketSnapshot{}, err
	}

	return MarketSnapshot{
		Portfolio: PortfolioState{
			Holdings:        parseNumberMap(portfolioResponse["holdings"]),
			PricesUsd:       parseNumberMap(portfolioResponse["prices_usd"]),
			UsdcReserve:     asNumber(portfolioResponse["usdc_reserve"], 0),
			CashBalanceUsd:  asNumber(portfolioResponse["cash_balance_usd"], 0),
			VolatilityIndex: asNumber(portfolioResponse["volatility_index"], 0),
		},
		Prices:          parseNumberMap(pricesResponse["prices"]),
		VolatilityIndex: asNumber(volatilityResponse["volatility_index"], 0),
	}, nil
}

// buildPolicy assembles the policy from the batch-fetched Vault DON secrets.
// The secrets are released only into an attested enclave and decrypted at the
// moment GetSecrets runs, so node operators never see the target weights or
// trade caps.
func buildPolicy(secretValues map[string]string, ids SecretsConfig) (Policy, error) {
	numeric := func(secretID string) (float64, error) {
		return parseRequiredSecretNumber(secretValues[secretID], secretID)
	}

	targetBtc, err := numeric(ids.TargetAllocationBtcPctSecretID)
	if err != nil {
		return Policy{}, err
	}
	targetEth, err := numeric(ids.TargetAllocationEthPctSecretID)
	if err != nil {
		return Policy{}, err
	}
	targetUsdc, err := numeric(ids.TargetAllocationUsdcPctSecretID)
	if err != nil {
		return Policy{}, err
	}
	driftThreshold, err := numeric(ids.DriftThresholdPctSecretID)
	if err != nil {
		return Policy{}, err
	}
	maxTradeUsd, err := numeric(ids.MaxTradeUsdSecretID)
	if err != nil {
		return Policy{}, err
	}
	reserveFloor, err := numeric(ids.ReserveFloorUsdcSecretID)
	if err != nil {
		return Policy{}, err
	}
	maxSlippage, err := numeric(ids.MaxSlippageBpsSecretID)
	if err != nil {
		return Policy{}, err
	}

	return Policy{
		TargetAllocations: map[string]float64{
			"BTC":  targetBtc,
			"ETH":  targetEth,
			"USDC": targetUsdc,
		},
		DriftThresholdPct:       driftThreshold,
		MaxTradeUsd:             maxTradeUsd,
		ReserveFloorUsdc:        reserveFloor,
		MaxSlippageBps:          maxSlippage,
		PreferredVenues:         parseVenueListSecret(secretValues[ids.PreferredVenuesSecretID]),
		OrderSequencePreference: parseOrderSequencePreference(secretValues[ids.OrderSequencePreferenceSecretID]),
	}, nil
}

// ─── TEE Cron Callback ──────────────────────────────────────
// Receives a cre.TeeRuntime, not a cre.Runtime: the secrets, the portfolio
// snapshot, the model prompt and the execution request all stay in the enclave.
func onCronTrigger(config *Config, runtime cre.TeeRuntime, _ *cron.Payload) (string, error) {
	ids := config.SecretsIDs

	// Batch fetch every secret in a single call; if any secret fails, the error
	// reports every failed secret at once. Results are keyed by ID afterwards so
	// no value depends on slice position.
	secrets, err := runtime.GetSecrets([]*cre.SecretRequest{
		{Id: ids.ExchangeAPIKeyID},
		{Id: ids.OpenAIAPIKeyID},
		{Id: ids.TargetAllocationBtcPctSecretID},
		{Id: ids.TargetAllocationEthPctSecretID},
		{Id: ids.TargetAllocationUsdcPctSecretID},
		{Id: ids.DriftThresholdPctSecretID},
		{Id: ids.MaxTradeUsdSecretID},
		{Id: ids.ReserveFloorUsdcSecretID},
		{Id: ids.MaxSlippageBpsSecretID},
		{Id: ids.PreferredVenuesSecretID},
		{Id: ids.OrderSequencePreferenceSecretID},
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

	policy, err := buildPolicy(secretValues, ids)
	if err != nil {
		return "", err
	}

	// ⚠️ Logs are for simulation only and MUST be removed before deploying to
	// production — anything logged inside the enclave weakens the
	// confidentiality guarantee. This records only that the fetch succeeded.
	runtime.Logger().Info("rebalance-getsecrets-ok")

	client := &http.Client{}

	snapshot, err := collectMarketSnapshot(runtime, client, config.MockBaseURL, exchangeAPIKey)
	if err != nil {
		return "", err
	}

	allocations, totalUsd, maxDriftPct := CalculateAllocations(
		snapshot.Portfolio.Holdings,
		snapshot.Prices,
		policy.TargetAllocations,
	)

	metrics := computePortfolioMetrics(snapshot, policy, totalUsd, maxDriftPct)

	prompt, err := createOpenAiPrompt(allocations, metrics, policy, snapshot)
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
				"content": "You are a portfolio rebalancing policy engine. Emit strict JSON only " +
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

	targetTrades := BuildTargetAlignedTrades(allocations)
	reconciledDecision := ReconcileDecisionWithTargets(decision, targetTrades, policy, metrics)

	chunks, err := EnforcePolicyAndChunk(reconciledDecision, policy, snapshot.Portfolio.UsdcReserve)
	if err != nil {
		return "", err
	}

	if maxDriftPct < policy.DriftThresholdPct || len(chunks) == 0 {
		runtime.Logger().Info("rebalance-skip",
			"drift", fmt.Sprintf("%.3f", maxDriftPct),
			"threshold", policy.DriftThresholdPct,
			"reason", reconciledDecision.Reasoning,
		)
		return "NOOP", nil
	}

	executionPlan := BuildExecutionPlan(chunks)

	executionHeaders := jsonHeaders()
	executionHeaders["x-exchange-api-key"] = exchangeAPIKey

	executionResponse, err := postJSON(runtime, client, config.MockBaseURL+"/execute-rebalance",
		map[string]any{
			"totalPortfolioUsd":  totalUsd,
			"maxDriftPct":        maxDriftPct,
			"reserveBufferUsd":   metrics.ReserveBufferUsd,
			"stablecoinDepthUsd": metrics.StablecoinDepthUsd,
			"volatilityIndex":    metrics.AverageVolatility,
			"reasoning":          reconciledDecision.Reasoning,
			"trades":             chunks,
			"onchain_trades":     executionPlan.OnchainTrades,
			"offchain_trades":    executionPlan.OffchainTrades,
		}, executionHeaders)
	if err != nil {
		return "", err
	}

	executionID := asString(executionResponse["execution_id"], "unknown")

	runtime.Logger().Info("rebalance-executed",
		"trade_count", len(chunks),
		"execution_id", executionID,
	)

	result, err := json.Marshal(map[string]any{
		"status":             "EXECUTED",
		"tradeCount":         len(chunks),
		"onchainTradeCount":  len(executionPlan.OnchainTrades),
		"offchainTradeCount": len(executionPlan.OffchainTrades),
		"maxDriftPct":        maxDriftPct,
		"executionId":        executionID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to encode result: %w", err)
	}
	return string(result), nil
}

// ─── Workflow Init ──────────────────────────────────────────

func InitWorkflow(config *Config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[*Config], error) {
	if config.Schedule == "" || config.MockBaseURL == "" || config.OpenAIURL == "" || config.OpenAIModel == "" {
		return nil, fmt.Errorf("config requires schedule, mock_base_url, openai_url, and openai_model")
	}

	ids := config.SecretsIDs
	if ids.ExchangeAPIKeyID == "" ||
		ids.OpenAIAPIKeyID == "" ||
		ids.TargetAllocationBtcPctSecretID == "" ||
		ids.TargetAllocationEthPctSecretID == "" ||
		ids.TargetAllocationUsdcPctSecretID == "" ||
		ids.DriftThresholdPctSecretID == "" ||
		ids.MaxTradeUsdSecretID == "" ||
		ids.ReserveFloorUsdcSecretID == "" ||
		ids.MaxSlippageBpsSecretID == "" ||
		ids.PreferredVenuesSecretID == "" ||
		ids.OrderSequencePreferenceSecretID == "" {
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
			cre.OneOfTees{cre.Nitro{Regions: []cre.NitroRegion{cre.NitroUsWest2}}},
		),
	}, nil
}
