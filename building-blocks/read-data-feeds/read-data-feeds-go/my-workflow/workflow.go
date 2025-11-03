package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"my-project/contracts/evm/src/generated/price_feed_aggregator"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/cre"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
)

type Feed struct {
	Name    string `json:"name"`    // e.g. "BTC/USD", "ETH/USD"
	Address string `json:"address"` // proxy address on the target chain
}

type Config struct {
	// 6-field cron; this fires on the 0th second every 10 minutes
	Schedule  string `json:"schedule"`  // e.g. "0 */10 * * * *"
	ChainName string `json:"chainName"` // e.g. "ethereum-mainnet-arbitrum-1"
	Feeds     []Feed `json:"feeds"`     // list of feeds to read
}

type PriceResult struct {
	Name            string `json:"name"`
	Address         string `json:"address"`
	Decimals        uint8  `json:"decimals"`
	LatestAnswerRaw string `json:"latestAnswerRaw"`
	Scaled          string `json:"scaled"`
}

func InitWorkflow(cfg *Config, logger *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[*Config], error) {
	return cre.Workflow[*Config]{
		cre.Handler(cron.Trigger(&cron.Config{Schedule: cfg.Schedule}), onTick),
	}, nil
}

func onTick(cfg *Config, runtime cre.Runtime, _ *cron.Payload) (string, error) {
	lg := runtime.Logger()

	selector, err := evm.ChainSelectorFromName(cfg.ChainName)
	if err != nil {
		return "", fmt.Errorf("chain selector: %w", err)
	}
	client := &evm.Client{ChainSelector: selector}

	results := make([]PriceResult, 0, len(cfg.Feeds))

	for _, f := range cfg.Feeds {
		addr := common.HexToAddress(f.Address)
		feed, err := price_feed_aggregator.NewPriceFeedAggregator(client, addr, nil)
		if err != nil {
			lg.Error("binding failed", "feed", f.Name, "address", f.Address, "err", err)
			return "", err
		}

		// Passing nil blockNumber uses the finalized block
		decimals, err := feed.Decimals(runtime, nil).Await()
		if err != nil {
			lg.Error("decimals() failed", "feed", f.Name, "address", f.Address, "err", err)
			return "", err
		}

		raw, err := feed.LatestAnswer(runtime, nil).Await()
		if err != nil {
			lg.Error("latestAnswer() failed", "feed", f.Name, "address", f.Address, "err", err)
			return "", err
		}

		scale := decimal.New(1, int32(decimals)) // 10^decimals
		scaled := decimal.NewFromBigInt(raw, 0).Div(scale)

		lg.Info("Data feed read",
			"chain", cfg.ChainName,
			"feed", f.Name,
			"address", f.Address,
			"decimals", decimals,
			"latestAnswerRaw", raw.String(),
			"latestAnswerScaled", scaled.String(),
		)

		results = append(results, PriceResult{
			Name:            f.Name,
			Address:         f.Address,
			Decimals:        decimals,
			LatestAnswerRaw: raw.String(),
			Scaled:          scaled.String(),
		})
	}

	// Return a JSON string so the simulation result is easy to consume.
	out, err := json.Marshal(results)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
