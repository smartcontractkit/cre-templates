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
	Schedule         string `json:"schedule"`         // 6-field cron; e.g. "0 */10 * * * *"
	ChainName        string `json:"chainName"`        // e.g. "ethereum-mainnet-arbitrum-1"
	Feed             Feed   `json:"feed"`             // single feed to monitor
	WebhookURL       string `json:"webhookUrl"`       // Slack or Telegram webhook URL
	NotificationType string `json:"notificationType"` // "slack" or "telegram"
	TelegramChatID   string `json:"telegramChatId"`   // Telegram chat ID (only for telegram)
}

type NotificationResult struct {
	Feed              string `json:"feed"`
	Address           string `json:"address"`
	Decimals          uint8  `json:"decimals"`
	LatestAnswerRaw   string `json:"latestAnswerRaw"`
	Scaled            string `json:"scaled"`
	FormattedPrice    string `json:"formattedPrice"`
	NotificationType  string `json:"notificationType"`
	WebhookStatusCode uint32 `json:"webhookStatusCode"`
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

	// 2. Format price and build webhook body
	formattedPrice := formatPrice(raw, int32(decimals))
	webhookBody, err := buildWebhookBody(cfg, formattedPrice)
	if err != nil {
		return "", fmt.Errorf("build webhook body: %w", err)
	}

	lg.Info("Sending notification",
		"type", cfg.NotificationType,
		"url", cfg.WebhookURL,
	)

	// 3. Send webhook via ConfidentialHTTPClient
	confHttpClient := &confidentialhttp.Client{}
	resp, err := confHttpClient.SendRequest(runtime, &confidentialhttp.ConfidentialHTTPRequest{
		Request: &confidentialhttp.HTTPRequest{
			Url:    cfg.WebhookURL,
			Method: "POST",
			Body:   &confidentialhttp.HTTPRequest_BodyString{BodyString: webhookBody},
			MultiHeaders: map[string]*confidentialhttp.HeaderValues{
				"Content-Type": {Values: []string{"application/json"}},
			},
		},
		VaultDonSecrets: []*confidentialhttp.SecretIdentifier{},
	}).Await()
	if err != nil {
		return "", fmt.Errorf("webhook send failed: %w", err)
	}

	lg.Info("Webhook response", "statusCode", resp.StatusCode)

	// 4. Return summary
	result := NotificationResult{
		Feed:              cfg.Feed.Name,
		Address:           cfg.Feed.Address,
		Decimals:          decimals,
		LatestAnswerRaw:   raw.String(),
		Scaled:            scaled.String(),
		FormattedPrice:    formattedPrice,
		NotificationType:  cfg.NotificationType,
		WebhookStatusCode: resp.StatusCode,
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
	// Format with 2 decimal places
	return price.StringFixed(2)
}

func buildWebhookBody(cfg *Config, formattedPrice string) (string, error) {
	var payload interface{}

	if cfg.NotificationType == "telegram" {
		payload = map[string]string{
			"chat_id":    cfg.TelegramChatID,
			"text":       fmt.Sprintf("*%s*: $%s", cfg.Feed.Name, formattedPrice),
			"parse_mode": "Markdown",
		}
	} else {
		// Default: Slack
		payload = map[string]string{
			"text": fmt.Sprintf(":chart_with_upwards_trend: *%s*: $%s", cfg.Feed.Name, formattedPrice),
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
