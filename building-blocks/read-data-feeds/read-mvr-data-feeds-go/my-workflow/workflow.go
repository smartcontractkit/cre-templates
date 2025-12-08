package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"my-project/contracts/evm/src/generated/bundle_aggregator_proxy"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

// Feed is the config for a single MVR data feed (BundleAggregatorProxy).
type Feed struct {
	Name    string `json:"name"`    // e.g. "S&P Global SSA EURC"
	Address string `json:"address"` // proxy address on the target chain
}

type Config struct {
	// 6-field cron; this fires on the 0th second every 10 minutes
	Schedule  string `json:"schedule"`  // e.g. "0 */10 * * * *"
	ChainName string `json:"chainName"` // e.g. "ethereum-mainnet-base-1"
	Feeds     []Feed `json:"feeds"`     // list of MVR feeds to read
}

// DecodedBundle represents the typed, decoded contents of latestBundle().
type DecodedBundle struct {
	// Raw integers as strings (so JSON consumers don't lose precision)
	LastModifiedDateTimeRaw string `json:"lastModifiedDateTimeRaw"`
	// RFC3339 representation if LastModifiedDateTime is seconds since epoch.
	LastModifiedDateTimeRFC3339 string `json:"lastModifiedDateTimeRfc3339"`

	SecurityID   string `json:"securityId"`
	SecurityName string `json:"securityName"`

	SSARaw     string `json:"ssaRaw"`
	SSAScaled  string `json:"ssaScaled"`
	SSADesc    string `json:"ssaDesc"`
	SSADecimal uint8  `json:"ssaDecimal"`
}

// BundleResult is what we output per feed.
type BundleResult struct {
	Name           string        `json:"name"`
	Address        string        `json:"address"`
	Bundle         DecodedBundle `json:"bundle"`
	BundleDecimals []uint8       `json:"bundleDecimals"`
}

func InitWorkflow(cfg *Config, logger *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[*Config], error) {
	return cre.Workflow[*Config]{
		cre.Handler(cron.Trigger(&cron.Config{Schedule: cfg.Schedule}), onTick),
	}, nil
}

func onTick(cfg *Config, runtime cre.Runtime, _ *cron.Payload) (string, error) {
	lg := runtime.Logger()

	selector, err := evm.ChainSelectorFromName(cfg.ChainName)
	lg.Info("chain selector obtained", "chainName", cfg.ChainName)
	lg.Info("selector details", "selector", selector)
	if err != nil {
		return "", fmt.Errorf("chain selector: %w", err)
	}
	client := &evm.Client{ChainSelector: selector}

	results := make([]BundleResult, 0, len(cfg.Feeds))

	for _, f := range cfg.Feeds {
		addr := common.HexToAddress(f.Address)
		proxy, err := bundle_aggregator_proxy.NewBundleAggregatorProxy(client, addr, nil)
		if err != nil {
			lg.Error("binding failed", "feed", f.Name, "address", f.Address, "err", err)
			return "", err
		}

		// Get decimals for each field in the bundle.
		decimals, err := proxy.BundleDecimals(runtime, nil).Await()
		if err != nil {
			lg.Error("bundleDecimals() failed", "feed", f.Name, "address", f.Address, "err", err)
			return "", err
		}

		// Get the raw bundle bytes.
		rawBundle, err := proxy.LatestBundle(runtime, nil).Await()
		if err != nil {
			lg.Error("latestBundle() failed", "feed", f.Name, "address", f.Address, "err", err)
			return "", err
		}

		bundle, err := decodeBundleMVR(rawBundle, decimals)
		if err != nil {
			lg.Error("decode bundle failed", "feed", f.Name, "address", f.Address, "err", err)
			return "", err
		}

		lg.Info("MVR bundle read",
			"chain", cfg.ChainName,
			"feed", f.Name,
			"address", f.Address,
			"lastModifiedDateTimeRaw", bundle.LastModifiedDateTimeRaw,
			"lastModifiedDateTimeRFC3339", bundle.LastModifiedDateTimeRFC3339,
			"securityId", bundle.SecurityID,
			"securityName", bundle.SecurityName,
			"ssaRaw", bundle.SSARaw,
			"ssaScaled", bundle.SSAScaled,
			"ssaDesc", bundle.SSADesc,
			"ssaDecimal", bundle.SSADecimal,
		)

		results = append(results, BundleResult{
			Name:           f.Name,
			Address:        f.Address,
			Bundle:         *bundle,
			BundleDecimals: decimals,
		})
	}

	out, err := json.Marshal(results)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// decodeBundleMVR decodes the bytes returned by latestBundle() into the known
// 5-field structure:
//
//	0: LastModifiedDateTime : uint256
//	1: SecurityID           : string
//	2: SecurityName         : string
//	3: SSA                  : uint256
//	4: SSADesc              : string
//
// The decimals array is used to scale SSA; it is assumed that
// decimals[3] corresponds to SSA.
func decodeBundleMVR(data []byte, decimals []uint8) (*DecodedBundle, error) {
	const wordSize = 32
	if len(data) < 5*wordSize {
		return nil, fmt.Errorf("bundle too short: got %d bytes", len(data))
	}

	getWord := func(i int) ([]byte, error) {
		start := i * wordSize
		end := start + wordSize
		if end > len(data) {
			return nil, fmt.Errorf("out of range word index %d", i)
		}
		return data[start:end], nil
	}

	// Parse head
	w0, err := getWord(0)
	if err != nil {
		return nil, err
	}
	w1, err := getWord(1)
	if err != nil {
		return nil, err
	}
	w2, err := getWord(2)
	if err != nil {
		return nil, err
	}
	w3, err := getWord(3)
	if err != nil {
		return nil, err
	}
	w4, err := getWord(4)
	if err != nil {
		return nil, err
	}

	lastModified := new(big.Int).SetBytes(w0)
	offsetID := new(big.Int).SetBytes(w1).Uint64()
	offsetName := new(big.Int).SetBytes(w2).Uint64()
	ssa := new(big.Int).SetBytes(w3)
	offsetDesc := new(big.Int).SetBytes(w4).Uint64()

	// Parse strings at the given offsets.
	securityID, err := parseABIString(data, int(offsetID))
	if err != nil {
		return nil, fmt.Errorf("parse securityId: %w", err)
	}
	securityName, err := parseABIString(data, int(offsetName))
	if err != nil {
		return nil, fmt.Errorf("parse securityName: %w", err)
	}
	ssaDesc, err := parseABIString(data, int(offsetDesc))
	if err != nil {
		return nil, fmt.Errorf("parse ssaDesc: %w", err)
	}

	// Scale SSA according to decimals[3] (if present).
	var ssaScaled decimal.Decimal
	var ssaDecimals uint8
	if len(decimals) > 3 {
		ssaDecimals = decimals[3]
		scale := decimal.New(1, int32(ssaDecimals)) // 10^decimals
		ssaScaled = decimal.NewFromBigInt(ssa, 0).Div(scale)
	} else {
		ssaScaled = decimal.NewFromBigInt(ssa, 0)
	}

	// Interpret LastModifiedDateTime as seconds since Unix epoch (common pattern).
	var lastModifiedRFC3339 string
	if lastModified.IsInt64() {
		t := time.Unix(lastModified.Int64(), 0).UTC()
		lastModifiedRFC3339 = t.Format(time.RFC3339)
	}

	return &DecodedBundle{
		LastModifiedDateTimeRaw:     lastModified.String(),
		LastModifiedDateTimeRFC3339: lastModifiedRFC3339,
		SecurityID:                  securityID,
		SecurityName:                securityName,
		SSARaw:                      ssa.String(),
		SSAScaled:                   ssaScaled.String(),
		SSADesc:                     ssaDesc,
		SSADecimal:                  ssaDecimals,
	}, nil
}

// parseABIString reads a Solidity ABI-encoded string at the given byte offset.
// At offset:
//
//	[0..31]   length (uint256)
//	[32..]    data (length bytes) padded to 32-byte boundary
func parseABIString(data []byte, offset int) (string, error) {
	const wordSize = 32
	if offset < 0 || offset+wordSize > len(data) {
		return "", fmt.Errorf("invalid offset %d", offset)
	}

	lengthBytes := data[offset : offset+wordSize]
	length := new(big.Int).SetBytes(lengthBytes).Uint64()

	dataStart := offset + wordSize
	dataEnd := dataStart + int(length)
	if dataEnd > len(data) {
		return "", fmt.Errorf("string length out of range: len=%d offset=%d", length, offset)
	}

	return string(data[dataStart:dataEnd]), nil
}
