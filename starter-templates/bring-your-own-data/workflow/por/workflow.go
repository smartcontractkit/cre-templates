package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

// EVMConfig holds per-chain configuration.
type EVMConfig struct {
	DataFeedsCacheAddress string `json:"dataFeedsCacheAddress"`
	ChainName             string `json:"chainName"`
	GasLimit              uint64 `json:"gasLimit"`
}

func (e *EVMConfig) GetChainSelector() (uint64, error) {
	return evm.ChainSelectorFromName(e.ChainName)
}

func (e *EVMConfig) NewEVMClient() (*evm.Client, error) {
	chainSelector, err := e.GetChainSelector()
	if err != nil {
		return nil, err
	}
	return &evm.Client{
		ChainSelector: chainSelector,
	}, nil
}

type Config struct {
	Schedule  string      `json:"schedule"`
	URL       string      `json:"url"`
	DataIdHex string      `json:"dataIdHex"`
	EVMs      []EVMConfig `json:"evms"`
}

type ReserveInfo struct {
	LastUpdated  time.Time       `consensus_aggregation:"median" json:"lastUpdated"`
	TotalReserve decimal.Decimal `consensus_aggregation:"median" json:"totalReserve"`
}

type PORResponse struct {
	AccountName string    `json:"accountName"`
	TotalTrust  float64   `json:"totalTrust"`
	TotalToken  float64   `json:"totalToken"`
	Ripcord     bool      `json:"ripcord"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func InitWorkflow(config *Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*Config], error) {
	cronTriggerCfg := &cron.Config{
		Schedule: config.Schedule,
	}
	workflow := cre.Workflow[*Config]{
		cre.Handler(
			cron.Trigger(cronTriggerCfg),
			onCronTrigger,
		),
	}
	return workflow, nil
}

func onCronTrigger(config *Config, runtime cre.Runtime, outputs *cron.Payload) (string, error) {
	logger := runtime.Logger()
	s, err := doPOR(config, runtime, outputs.ScheduledExecutionTime.AsTime())
	if err != nil {
		logger.Error("doPOR failed", "err", err)
		return "", err
	}
	return s, nil
}

func doPOR(config *Config, runtime cre.Runtime, runTime time.Time) (string, error) {

	logger := runtime.Logger()
	logger.Info("fetching por", "url", config.URL, "evms", config.EVMs, "time", runTime)

	reserveInfo, err := http.SendRequest(
		config,
		runtime,
		&http.Client{},
		fetchPOR,
		cre.ConsensusAggregationFromTags[*ReserveInfo]()).Await()
	if err != nil {
		return "", fmt.Errorf("error fetching por: %w", err)
	}

	logger.Info("ReserveInfo", "reserveInfo", reserveInfo)
	totalReserveScaled := reserveInfo.TotalReserve.Mul(decimal.NewFromUint64(1e18)).BigInt()
	logger.Info("TotalReserveScaled", "totalReserveScaled", totalReserveScaled)

	// Update reserves on each target chain
	for _, evmCfg := range config.EVMs {
		logger.Info("updating reserves", "totalReserveScaled", totalReserveScaled, "chain", evmCfg.ChainName)
		if err := updateReserves(evmCfg, runtime, config.DataIdHex, totalReserveScaled, reserveInfo.LastUpdated); err != nil {
			return "", fmt.Errorf("failed to update reserves on chain \"%s\": %w", evmCfg.ChainName, err)
		}
	}

	return reserveInfo.TotalReserve.String(), nil
}

func updateReserves(evmCfg EVMConfig, runtime cre.Runtime, dataIdHex string, totalReserveScaled *big.Int, lastUpdated time.Time) error {

	logger := runtime.Logger()

	evmClient, err := evmCfg.NewEVMClient()
	if err != nil {
		return fmt.Errorf("failed to create EVM client for %s: %w", evmCfg.ChainName, err)
	}

	dataId, err := hexToBytes32(dataIdHex)
	if err != nil {
		return fmt.Errorf("failed to convert data ID %s: %w", dataIdHex, err)
	}

	bundle, err := encodeBundleStruct(Bundle{
		TotalReserve: totalReserveScaled,
		LastUpdated:  uint32(lastUpdated.UTC().Unix()),
	})
	if err != nil {
		return fmt.Errorf("failed to encode Bundle struct: %w", err)
	}

	encodedStruct, err := encodeReceivedBundledReports(
		[]receivedBundledReport{
			{
				DataId:    dataId,
				Timestamp: uint32(runtime.Now().UTC().Unix()),
				Bundle:    bundle,
			},
		})
	if err != nil {
		return fmt.Errorf("failed to encode ReceivedBundledReport: %w", err)
	}

	logger.Info("Generating report...")
	reportPromise := runtime.GenerateReport(&cre.ReportRequest{
		EncodedPayload: encodedStruct,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	})
	report, err := reportPromise.Await()
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}
	logger.Info("Report generated successfully")

	logger.Info("Writing report", "totalReserveScaled", totalReserveScaled)
	writePromise := evmClient.WriteReport(runtime, &evm.WriteCreReportRequest{
		Receiver:  common.HexToAddress(evmCfg.DataFeedsCacheAddress).Bytes(),
		Report:    report,
		GasConfig: &evm.GasConfig{GasLimit: evmCfg.GasLimit},
	})
	resp, err := writePromise.Await()
	if err != nil {
		return fmt.Errorf("failed to submit report: %w", err)
	}
	if resp.TxStatus != evm.TxStatus_TX_STATUS_SUCCESS {
		errorMsg := "unknown error"
		if resp.ErrorMessage != nil {
			errorMsg = *resp.ErrorMessage
		}
		return fmt.Errorf("transaction failed with status %v: %s", resp.TxStatus, errorMsg)
	} else if resp.ReceiverContractExecutionStatus != nil &&
		*resp.ReceiverContractExecutionStatus != evm.ReceiverContractExecutionStatus_RECEIVER_CONTRACT_EXECUTION_STATUS_SUCCESS {
		return fmt.Errorf("failed to execute receiver contract")
	}
	logger.Info("Write report transaction succeeded at", "txHash", common.BytesToHash(resp.TxHash).Hex(), "chainName", evmCfg.ChainName, "response", resp)

	return nil
}

func fetchPOR(config *Config, logger *slog.Logger, sendRequester *http.SendRequester) (*ReserveInfo, error) {

	httpActionOut, err := sendRequester.SendRequest(&http.Request{
		Method: "GET",
		Url:    config.URL,
	}).Await()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch POR: %w", err)
	}

	porResp := &PORResponse{}
	if err = json.Unmarshal(httpActionOut.Body, porResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal POR response: %w", err)
	}

	if porResp.Ripcord {
		return nil, errors.New("ripcord is true")
	}

	res := &ReserveInfo{
		LastUpdated:  porResp.UpdatedAt.UTC(),
		TotalReserve: decimal.NewFromFloat(porResp.TotalToken),
	}

	return res, nil
}

type Bundle struct {
	TotalReserve *big.Int
	LastUpdated  uint32
}

func encodeBundleStruct(in Bundle) ([]byte, error) {
	tupleType, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "lastUpdated", Type: "uint32"},
			{Name: "totalReserve", Type: "uint256"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tuple type for ReserveInfo: %w", err)
	}
	args := abi.Arguments{
		{Name: "reserveInfo", Type: tupleType},
	}

	return args.Pack(in)
}

// hexToBytes32 decodes hexString as a hexadecimal to a [32]byte array. If
// hexString represents fewer than 32 bytes, the resulting array is
// right-padded with zeros.
func hexToBytes32(hexString string) ([32]byte, error) {
	b, err := hex.DecodeString(hexString)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to decode hex string: %w", err)
	}
	if len(b) > 32 {
		return [32]byte{}, fmt.Errorf("hex string of length %d decodes to %d bytes, which exceeds 32 bytes", len(hexString), len(b))
	} else if len(b) < 32 {
		nb := [32]byte{}
		copy(nb[:], b[:])
		return nb, nil
	}
	return [32]byte(b), nil
}
