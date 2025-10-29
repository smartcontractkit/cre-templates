package main

import (
	"fmt"
	"log/slog"

	"github.com/shopspring/decimal"

	"github.com/ethereum/go-ethereum/common"

	"my-project/contracts/evm/src/generated/btcusd_price_feed"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

type Config struct {
	// 6-field cron; this fires on the 0th second every 10 minutes
	Schedule          string `json:"schedule"`          // e.g. "0 */10 * * * *"
	ChainName         string `json:"chainName"`         // e.g. "ethereum-mainnet-arbitrum-1"
	AggregatorAddress string `json:"aggregatorAddress"` // e.g. "0x6ce185860a4963106506C203335A2910413708e9"
}

func InitWorkflow(cfg *Config, logger *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[*Config], error) {
	return cre.Workflow[*Config]{
		cre.Handler(cron.Trigger(&cron.Config{Schedule: cfg.Schedule}), onTick),
	}, nil
}

func onTick(cfg *Config, runtime cre.Runtime, _ *cron.Payload) (string, error) {
	lg := runtime.Logger()

	// Build EVM client from chain name
	selector, err := evm.ChainSelectorFromName(cfg.ChainName)
	if err != nil {
		return "", fmt.Errorf("chain selector: %w", err)
	}
	client := &evm.Client{ChainSelector: selector}

	// Bind the feed proxy.
	addr := common.HexToAddress(cfg.AggregatorAddress)
	feed, err := btcusd_price_feed.NewBTCUSDPriceFeed(client, addr, nil)
	if err != nil {
		return "", fmt.Errorf("binding: %w", err)
	}

	// Reads default to the finalized block when you pass nil
	decimals, err := feed.Decimals(runtime, nil).Await()
	if err != nil {
		lg.Error("decimals() failed", "err", err)
		return "", err
	}

	raw, err := feed.LatestAnswer(runtime, nil).Await()
	if err != nil {
		lg.Error("latestAnswer() failed", "err", err)
		return "", err
	}

	// Scale latestAnswer by 10^decimals.
	scale := decimal.New(1, int32(decimals)) // = 10^decimals
	scaled := decimal.NewFromBigInt(raw, 0).Div(scale)

	lg.Info("BTC/USD feed read",
		"chain", cfg.ChainName,
		"address", cfg.AggregatorAddress,
		"decimals", decimals,
		"latestAnswerRaw", raw.String(),
		"latestAnswerScaled", scaled.String(),
	)

	// Return the human-readable number
	return scaled.String(), nil
}
