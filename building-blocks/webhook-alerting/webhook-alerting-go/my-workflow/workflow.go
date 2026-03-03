package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"

	"my-project/contracts/evm/src/generated/price_feed_aggregator"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/confidentialhttp"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

type Feed struct {
	Name    string `json:"name"`    // e.g. "ETH/USD"
	Address string `json:"address"` // proxy address on the target chain
}

type Config struct {
	Schedule  string `json:"schedule"`  // 6-field cron; e.g. "0 */10 * * * *"
	ChainName string `json:"chainName"` // e.g. "ethereum-mainnet-arbitrum-1"
	Feed      Feed   `json:"feed"`      // single feed to monitor
	Endpoint  string `json:"endpoint"`  // PagerDuty Events API v2 endpoint
	Severity  string `json:"severity"`  // "critical", "error", "warning", or "info"
	Source    string `json:"source"`    // source identifier for the alert
}

type AlertResult struct {
	Feed            string `json:"feed"`
	Address         string `json:"address"`
	Decimals        uint8  `json:"decimals"`
	LatestAnswerRaw string `json:"latestAnswerRaw"`
	Scaled          string `json:"scaled"`
	FormattedPrice  string `json:"formattedPrice"`
	AlertEndpoint   string `json:"alertEndpoint"`
	AlertStatusCode uint32 `json:"alertStatusCode"`
}

func InitWorkflow(cfg *Config, logger *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[*Config], error) {
	return cre.Workflow[*Config]{
		cre.Handler(cron.Trigger(&cron.Config{Schedule: cfg.Schedule}), onTick),
	}, nil
}

func onTick(cfg *Config, runtime cre.Runtime, _ *cron.Payload) (string, error) {
	lg := runtime.Logger()

	// 1. Read the data feed
	selector, err := evm.ChainSelectorFromName(cfg.ChainName)
	if err != nil {
		return "", fmt.Errorf("chain selector: %w", err)
	}
	client := &evm.Client{ChainSelector: selector}

	addr := common.HexToAddress(cfg.Feed.Address)
	feed, err := price_feed_aggregator.NewPriceFeedAggregator(client, addr, nil)
	if err != nil {
		return "", fmt.Errorf("binding failed: %w", err)
	}

	decimals, err := feed.Decimals(runtime, nil).Await()
	if err != nil {
		return "", fmt.Errorf("decimals() failed: %w", err)
	}

	raw, err := feed.LatestAnswer(runtime, nil).Await()
	if err != nil {
		return "", fmt.Errorf("latestAnswer() failed: %w", err)
	}

	scale := decimal.New(1, int32(decimals))
	scaled := decimal.NewFromBigInt(raw, 0).Div(scale)

	lg.Info("Data feed read",
		"chain", cfg.ChainName,
		"feed", cfg.Feed.Name,
		"address", cfg.Feed.Address,
		"decimals", decimals,
		"latestAnswerRaw", raw.String(),
		"latestAnswerScaled", scaled.String(),
	)

	// 2. Format price and build PagerDuty alert body
	formattedPrice := formatPrice(raw, int32(decimals))
	alertBody, err := buildPagerDutyBody(cfg, formattedPrice)
	if err != nil {
		return "", fmt.Errorf("build alert body: %w", err)
	}

	lg.Info("Sending PagerDuty alert",
		"endpoint", cfg.Endpoint,
	)

	// 3. Send alert via ConfidentialHTTPClient with secret injection
	//    The {{.pagerdutyRoutingKey}} template in the body is resolved by the
	//    enclave from VaultDON secrets (or from env vars during simulation).
	confHttpClient := &confidentialhttp.Client{}
	resp, err := confHttpClient.SendRequest(runtime, &confidentialhttp.ConfidentialHTTPRequest{
		VaultDonSecrets: []*confidentialhttp.SecretIdentifier{
			{Key: "pagerdutyRoutingKey"},
		},
		Request: &confidentialhttp.HTTPRequest{
			Url:    cfg.Endpoint,
			Method: "POST",
			Body:   &confidentialhttp.HTTPRequest_BodyString{BodyString: alertBody},
			MultiHeaders: map[string]*confidentialhttp.HeaderValues{
				"Content-Type": {Values: []string{"application/json"}},
			},
		},
	}).Await()
	if err != nil {
		return "", fmt.Errorf("alert send failed: %w", err)
	}

	lg.Info("Alert response", "statusCode", resp.StatusCode)

	// 4. Return summary
	result := AlertResult{
		Feed:            cfg.Feed.Name,
		Address:         cfg.Feed.Address,
		Decimals:        decimals,
		LatestAnswerRaw: raw.String(),
		Scaled:          scaled.String(),
		FormattedPrice:  formattedPrice,
		AlertEndpoint:   cfg.Endpoint,
		AlertStatusCode: resp.StatusCode,
	}

	out, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func formatPrice(raw *big.Int, decimals int32) string {
	scale := decimal.New(1, decimals)
	price := decimal.NewFromBigInt(raw, 0).Div(scale)
	return price.StringFixed(2)
}

func buildPagerDutyBody(cfg *Config, formattedPrice string) (string, error) {
	payload := map[string]interface{}{
		"routing_key":  "{{.pagerdutyRoutingKey}}",
		"event_action": "trigger",
		"payload": map[string]string{
			"summary":  fmt.Sprintf("%s price: $%s on %s", cfg.Feed.Name, formattedPrice, cfg.ChainName),
			"severity": cfg.Severity,
			"source":   cfg.Source,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
