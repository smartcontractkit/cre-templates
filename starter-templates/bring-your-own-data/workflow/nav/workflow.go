package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/durationpb"

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

type NAVInfo struct {
	NAV decimal.Decimal `consensus_aggregation:"median" json:"nav"`
}

type NAVResponse struct {
	AggregatedCollateral float64 `json:"_aggregatedCollateral"`
	TotalOwedM           float64 `json:"_totalOwedM"`
	TotalCollateral      float64 `json:"totalCollateral"`
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
	s, err := doNAV(config, runtime, outputs.ScheduledExecutionTime.AsTime())
	if err != nil {
		logger.Error("doNAV failed", "err", err)
		return "", err
	}
	return s, nil
}

func doNAV(config *Config, runtime cre.Runtime, runTime time.Time) (string, error) {

	logger := runtime.Logger()
	logger.Info("fetching nav", "url", config.URL, "evms", config.EVMs, "time", runTime)

	navResp, err := http.SendRequest(
		config,
		runtime,
		&http.Client{},
		fetchNAV,
		cre.ConsensusAggregationFromTags[*NAVInfo]()).Await()
	if err != nil {
		return "", fmt.Errorf("error fetching nav: %w", err)
	}

	logger.Info("NAVResponse", "navResponse", navResp)
	// TODO: do we need to scale here? Need to scale for at least 6 decimal
	// places.
	navScaled := navResp.NAV.Mul(decimal.NewFromUint64(1e18)).BigInt()
	logger.Info("NAVScaled", "navScaled", navScaled)

	// Update NAV on each target chain
	for _, evmCfg := range config.EVMs {
		logger.Info("updating nav", "navScaled", navScaled, "chain", evmCfg.ChainName)
		if err := updateNAV(evmCfg, runtime, config.DataIdHex, navScaled); err != nil {
			return "", fmt.Errorf("failed to update nav on chain \"%s\": %w", evmCfg.ChainName, err)
		}
	}

	return navResp.NAV.String(), nil
}

func updateNAV(evmCfg EVMConfig, runtime cre.Runtime, dataIdHex string, navScaled *big.Int) error {

	logger := runtime.Logger()

	evmClient, err := evmCfg.NewEVMClient()
	if err != nil {
		return fmt.Errorf("failed to create EVM client for %s: %w", evmCfg.ChainName, err)
	}

	dataId, err := hexToBytes32(dataIdHex)
	if err != nil {
		return fmt.Errorf("failed to convert data ID %s: %w", dataIdHex, err)
	}

	encodedStruct, err := encodeReceivedDecimalReport(
		[]receivedDecimalReport{
			{
				DataId:    dataId,
				Timestamp: uint32(runtime.Now().UTC().Unix()),
				Answer:    navScaled,
			},
		})
	if err != nil {
		return fmt.Errorf("failed to encode ReceivedDecimalReport: %w", err)
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

	logger.Info("Writing report", "navScaled", navScaled)
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
	logger.Info("Write report succeeded", "response", resp, "txHash", common.BytesToHash(resp.TxHash).Hex())

	return nil
}

func fetchNAV(config *Config, logger *slog.Logger, sendRequester *http.SendRequester) (*NAVInfo, error) {

	type body struct {
		Method string `json:"method"`
	}

	reqBody, err := json.Marshal(&body{
		Method: "navDetails",
	})
	httpActionOut, err := sendRequester.SendRequest(&http.Request{
		Method:        "POST",
		Url:           config.URL,
		Headers:       map[string]string{},
		Body:          reqBody,
		Timeout:       &durationpb.Duration{},
		CacheSettings: &http.CacheSettings{},
	}).Await()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch NAV: %w", err)
	}

	navResp := &NAVResponse{}
	if err = json.Unmarshal(httpActionOut.Body, navResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal NAV response: %w", err)
	}

	logger.Info("NAV Response", "navResp", navResp)

	nav := decimal.NewFromFloat(navResp.AggregatedCollateral).Sub(decimal.NewFromFloat(navResp.TotalOwedM))

	logger.Info("Computed NAV", "nav", nav)

	return &NAVInfo{
		NAV: nav,
	}, nil
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
