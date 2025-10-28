package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"cre-custom-data-feed-go/contracts/evm/src/generated/balance_reader"
	"cre-custom-data-feed-go/contracts/evm/src/generated/ierc20"
	"cre-custom-data-feed-go/contracts/evm/src/generated/message_emitter"
	"cre-custom-data-feed-go/contracts/evm/src/generated/reserve_manager"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	pbvalues "github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/bindings"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

// EVMConfig holds per-chain configuration.
type EVMConfig struct {
	TokenAddress          string `json:"tokenAddress"`
	ReserveManagerAddress string `json:"reserveManagerAddress"`
	BalanceReaderAddress  string `json:"balanceReaderAddress"`
	MessageEmitterAddress string `json:"messageEmitterAddress"`
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
	Schedule string      `json:"schedule"`
	URL      string      `json:"url"`
	EVMs     []EVMConfig `json:"evms"`
}

type HTTPTriggerPayload struct {
	ExecutionTime time.Time `json:"executionTime"`
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
			onPORCronTrigger,
		),
	}

	for _, evmCfg := range config.EVMs {
		msgEmitter, err := prepareMessageEmitter(logger, evmCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to prepare message emitter: %w", err)
		}
		chainSelector, err := evmCfg.GetChainSelector()
		if err != nil {
			return nil, fmt.Errorf("failed to get chain selector: %w", err)
		}
		trigger, err := msgEmitter.LogTriggerMessageEmittedLog(chainSelector, evm.ConfidenceLevel_CONFIDENCE_LEVEL_LATEST, []message_emitter.MessageEmittedTopics{})
		if err != nil {
			return nil, fmt.Errorf("failed to create message emitted trigger: %w", err)
		}
		workflow = append(workflow, cre.Handler(trigger, onLogTrigger))
	}

	return workflow, nil
}

func onPORCronTrigger(config *Config, runtime cre.Runtime, outputs *cron.Payload) (string, error) {
	return doPOR(config, runtime)
}

func onLogTrigger(config *Config, runtime cre.Runtime, payload *bindings.DecodedLog[message_emitter.MessageEmittedDecoded]) (string, error) {
	logger := runtime.Logger()

	// use the decoded event log to get the event message
	message := payload.Data.Message
	logger.Info("Message retrieved from the event log", "message", message)

	// the event message can also be retrieved from the contract itself
	// below is an example of how to read from the contract
	messageEmitter, err := prepareMessageEmitter(logger, config.EVMs[0])
	if err != nil {
		return "", fmt.Errorf("failed to prepare message emitter: %w", err)
	}

	// use the decoded event log to get the emitter address
	// the emitter address is not a dynamic type, so it can be decoded from log even though its indexed
	emitter := payload.Data.Emitter
	lastMessageInput := message_emitter.GetLastMessageInput{
		Emitter: common.Address(emitter),
	}

	blockNumber := pbvalues.ProtoToBigInt(payload.Log.BlockNumber)
	logger.Info("Block number of event log", "blockNumber", blockNumber)
	message, err = messageEmitter.GetLastMessage(runtime, lastMessageInput, blockNumber).Await()
	if err != nil {
		logger.Error("Could not read from contract", "contract_chain", config.EVMs[0].ChainName, "err", err.Error())
		return "", err
	}
	logger.Info("Message retrieved from the contract", "message", message)

	return message, nil
}

func doPOR(config *Config, runtime cre.Runtime) (string, error) {
	logger := runtime.Logger()
	// Fetch PoR
	logger.Info("fetching por", "url", config.URL, "evms", config.EVMs)
	client := &http.Client{}
	reserveInfo, err := http.SendRequest(config, runtime, client, fetchPOR, cre.ConsensusAggregationFromTags[*ReserveInfo]()).Await()
	if err != nil {
		logger.Error("error fetching por", "err", err)
		return "", err
	}

	logger.Info("ReserveInfo", "reserveInfo", reserveInfo)

	totalSupply, err := getTotalSupply(config, runtime)
	if err != nil {
		return "", err
	}

	logger.Info("TotalSupply", "totalSupply", totalSupply)
	totalReserveScaled := reserveInfo.TotalReserve.Mul(decimal.NewFromUint64(1e18)).BigInt()
	logger.Info("TotalReserveScaled", "totalReserveScaled", totalReserveScaled)

	nativeTokenBalance, err := fetchNativeTokenBalance(runtime, config.EVMs[0], config.EVMs[0].TokenAddress)
	if err != nil {
		return "", fmt.Errorf("failed to fetch native token balance: %w", err)
	}
	logger.Info("Native token balance", "token", config.EVMs[0].TokenAddress, "balance", nativeTokenBalance)

	// Update reserves
	if err := updateReserves(config, runtime, totalSupply, totalReserveScaled); err != nil {
		return "", fmt.Errorf("failed to update reserves: %w", err)
	}

	return reserveInfo.TotalReserve.String(), nil
}

func prepareMessageEmitter(logger *slog.Logger, evmCfg EVMConfig) (*message_emitter.MessageEmitter, error) {
	evmClient, err := evmCfg.NewEVMClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create EVM client for %s: %w", evmCfg.ChainName, err)
	}

	address := common.HexToAddress(evmCfg.MessageEmitterAddress)

	messageEmitter, err := message_emitter.NewMessageEmitter(evmClient, address, nil)
	if err != nil {
		logger.Error("failed to create message emitter", "address", evmCfg.MessageEmitterAddress, "err", err)
		return nil, fmt.Errorf("failed to create message emitter for address %s: %w", evmCfg.MessageEmitterAddress, err)
	}

	return messageEmitter, nil
}

func fetchNativeTokenBalance(runtime cre.Runtime, evmCfg EVMConfig, tokenHolderAddress string) (*big.Int, error) {
	logger := runtime.Logger()
	evmClient, err := evmCfg.NewEVMClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create EVM client for %s: %w", evmCfg.ChainName, err)
	}

	balanceReaderAddress := common.HexToAddress(evmCfg.BalanceReaderAddress)
	balanceReader, err := balance_reader.NewBalanceReader(evmClient, balanceReaderAddress, nil)
	if err != nil {
		logger.Error("failed to create balance reader", "address", evmCfg.BalanceReaderAddress, "err", err)
		return nil, fmt.Errorf("failed to create balance reader for address %s: %w", evmCfg.BalanceReaderAddress, err)
	}
	tokenAddress, err := hexToBytes(tokenHolderAddress)
	if err != nil {
		logger.Error("failed to decode token address", "address", tokenHolderAddress, "err", err)
		return nil, fmt.Errorf("failed to decode token address %s: %w", tokenHolderAddress, err)
	}

	logger.Info("Getting native balances", "address", evmCfg.BalanceReaderAddress, "tokenAddress", tokenHolderAddress)
	balances, err := balanceReader.GetNativeBalances(runtime, balance_reader.GetNativeBalancesInput{
		Addresses: []common.Address{common.Address(tokenAddress)},
	}, big.NewInt(rpc.FinalizedBlockNumber.Int64())).Await()

	if err != nil {
		logger.Error("Could not read from contract", "contract_chain", evmCfg.ChainName, "err", err.Error())
		return nil, err
	}

	if len(balances) < 1 {
		logger.Error("No balances returned from contract", "contract_chain", evmCfg.ChainName)
		return nil, fmt.Errorf("no balances returned from contract for chain %s", evmCfg.ChainName)
	}

	return balances[0], nil
}

func getTotalSupply(config *Config, runtime cre.Runtime) (*big.Int, error) {
	evms := config.EVMs
	logger := runtime.Logger()
	// Fetch supply from all EVMs in parallel
	supplyPromises := make([]cre.Promise[*big.Int], len(evms))
	for i, evmCfg := range evms {
		evmClient, err := evmCfg.NewEVMClient()
		if err != nil {
			logger.Error("failed to create EVM client", "chainName", evmCfg.ChainName, "err", err)
			return nil, fmt.Errorf("failed to create EVM client for %s: %w", evmCfg.ChainName, err)
		}

		address := common.HexToAddress(evmCfg.TokenAddress)
		token, err := ierc20.NewIERC20(evmClient, address, nil)
		if err != nil {
			logger.Error("failed to create token", "address", evmCfg.TokenAddress, "err", err)
			return nil, fmt.Errorf("failed to create token for address %s: %w", evmCfg.TokenAddress, err)
		}
		evmTotalSupplyPromise := token.TotalSupply(runtime, big.NewInt(rpc.FinalizedBlockNumber.Int64()))
		supplyPromises[i] = evmTotalSupplyPromise
	}

	// We can add cre.AwaitAll that takes []cre.Promise[T] and returns ([]T, error)
	totalSupply := big.NewInt(0)
	for i, promise := range supplyPromises {
		supply, err := promise.Await()
		if err != nil {
			chainName := evms[i].ChainName
			logger.Error("Could not read from contract", "contract_chain", chainName, "err", err.Error())
			return nil, err
		}

		totalSupply = totalSupply.Add(totalSupply, supply)
	}

	return totalSupply, nil
}

func updateReserves(config *Config, runtime cre.Runtime, totalSupply *big.Int, totalReserveScaled *big.Int) error {
	evmCfg := config.EVMs[0]
	logger := runtime.Logger()
	logger.Info("Updating reserves", "totalSupply", totalSupply, "totalReserveScaled", totalReserveScaled)

	evmClient, err := evmCfg.NewEVMClient()
	if err != nil {
		return fmt.Errorf("failed to create EVM client for %s: %w", evmCfg.ChainName, err)
	}

	reserveManager, err := reserve_manager.NewReserveManager(evmClient, common.HexToAddress(evmCfg.ReserveManagerAddress), nil)
	if err != nil {
		return fmt.Errorf("failed to create reserve manager: %w", err)
	}

	logger.Info("Writing report", "totalSupply", totalSupply, "totalReserveScaled", totalReserveScaled)
	resp, err := reserveManager.WriteReportFromUpdateReserves(runtime, reserve_manager.UpdateReserves{
		TotalMinted:  totalSupply,
		TotalReserve: totalReserveScaled,
	}, nil).Await()

	if err != nil {
		logger.Error("WriteReport await failed", "error", err, "errorType", fmt.Sprintf("%T", err))
		return fmt.Errorf("failed to write report: %w", err)
	}
	logger.Info("Write report succeeded", "response", resp)
	logger.Info("Write report transaction succeeded at", "txHash", common.BytesToHash(resp.TxHash).Hex())
	return nil
}

func fetchPOR(config *Config, logger *slog.Logger, sendRequester *http.SendRequester) (*ReserveInfo, error) {
	httpActionOut, err := sendRequester.SendRequest(&http.Request{
		Method: "GET",
		Url:    config.URL,
	}).Await()
	if err != nil {
		return nil, err
	}

	porResp := &PORResponse{}
	if err = json.Unmarshal(httpActionOut.Body, porResp); err != nil {
		return nil, err
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

func hexToBytes(hexStr string) ([]byte, error) {
	if len(hexStr) < 2 || hexStr[:2] != "0x" {
		return nil, fmt.Errorf("invalid hex string: %s", hexStr)
	}
	return hex.DecodeString(hexStr[2:])
}
