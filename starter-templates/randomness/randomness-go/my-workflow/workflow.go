// Verifiable Randomness — CRE workflow (Go).
//
// VRF-style request/fulfill randomness. An on-chain event requests a random
// number; this workflow watches that event via a LogTrigger, generates a
// consensus-safe random word with runtime.Rand(), and writes a signed report
// back to the consumer. This is the CRE equivalent of VRF's
// requestRandomWords -> fulfillRandomWords round trip.
package main

import (
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/bindings"
	"github.com/smartcontractkit/cre-sdk-go/cre"

	"randomness-go/contracts/evm/src/generated/randomness_consumer"
)

// Config mirrors config.staging.json / config.production.json.
type Config struct {
	ChainSelectorName string `json:"chainSelectorName"`
	ConsumerAddress   string `json:"consumerAddress"`
	GasLimit          uint64 `json:"gasLimit"`
}

// newConsumer builds an EVM client + typed contract binding from config.
func newConsumer(config *Config) (*randomness_consumer.RandomnessConsumer, uint64, error) {
	chainSelector, err := evm.ChainSelectorFromName(config.ChainSelectorName)
	if err != nil {
		return nil, 0, fmt.Errorf("unknown chain %q: %w", config.ChainSelectorName, err)
	}
	evmClient := &evm.Client{ChainSelector: chainSelector}
	consumer, err := randomness_consumer.NewRandomnessConsumer(
		evmClient,
		common.HexToAddress(config.ConsumerAddress),
		nil,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create RandomnessConsumer binding: %w", err)
	}
	return consumer, chainSelector, nil
}

// InitWorkflow subscribes to RandomnessRequested via a LogTrigger.
func InitWorkflow(config *Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*Config], error) {
	consumer, chainSelector, err := newConsumer(config)
	if err != nil {
		return nil, err
	}

	// Read event data from a finalized block so the seed derives from state no
	// requester can reorg away. Passing nil filters matches every request.
	trigger, err := consumer.LogTriggerRandomnessRequestedLog(
		chainSelector,
		evm.ConfidenceLevel_CONFIDENCE_LEVEL_FINALIZED,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build RandomnessRequested log trigger: %w", err)
	}

	return cre.Workflow[*Config]{
		cre.Handler(trigger, onRandomnessRequested),
	}, nil
}

// onRandomnessRequested fires on each RandomnessRequested event. It generates a
// consensus-safe random word and writes it back to the consumer as a signed report.
func onRandomnessRequested(
	config *Config,
	runtime cre.Runtime,
	payload *bindings.DecodedLog[randomness_consumer.RandomnessRequestedDecoded],
) (string, error) {
	logger := runtime.Logger()

	// 1. Read decoded event data (typed via generated bindings).
	requestId := payload.Data.RequestId
	requester := payload.Data.Requester
	logger.Info("RandomnessRequested", "requestId", requestId, "requester", requester.Hex())

	consumer, _, err := newConsumer(config)
	if err != nil {
		return "", err
	}

	// 2. Idempotency: skip requests that are unknown or already fulfilled.
	//    Reading finalized state prevents fulfilling the same requestId twice.
	isPending, err := consumer.PendingRequests(
		runtime,
		randomness_consumer.PendingRequestsInput{Arg0: requestId},
		nil, // nil block number => finalized
	).Await()
	if err != nil {
		return "", fmt.Errorf("failed to read pendingRequests: %w", err)
	}
	if !isPending {
		logger.Info("request not pending (unknown or already fulfilled), skipping", "requestId", requestId)
		return fmt.Sprintf("Skipped — request %s not pending", requestId), nil
	}

	// 3. Generate a consensus-safe random word. runtime.Rand() is the CRE-specified
	//    randomness source: in DON mode every node derives the same seed, so the
	//    network agrees on a single value. Not a CSPRNG — see README.
	randSource, err := runtime.Rand()
	if err != nil {
		return "", fmt.Errorf("failed to get random source: %w", err)
	}
	// Draw across the full uint256 range [0, 2^256).
	maxWord := new(big.Int).Lsh(big.NewInt(1), 256)
	randomWord := new(big.Int).Rand(randSource, maxWord)
	logger.Info("generated randomWord", "requestId", requestId, "randomWord", randomWord)

	// 4. Encode the fulfillment payload exactly as RandomnessConsumer._processReport decodes it.
	encoded, err := encodeRandomnessReport(requestId, randomWord)
	if err != nil {
		return "", fmt.Errorf("failed to encode report payload: %w", err)
	}

	// 5. Generate a DON-signed report, then write it back via the Keystone Forwarder.
	report, err := runtime.GenerateReport(&cre.ReportRequest{
		EncodedPayload: encoded,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}).Await()
	if err != nil {
		return "", fmt.Errorf("failed to generate report: %w", err)
	}

	resp, err := consumer.WriteReport(runtime, report, &evm.GasConfig{GasLimit: config.GasLimit}).Await()
	if err != nil {
		return "", fmt.Errorf("failed to write report: %w", err)
	}
	if resp.TxStatus != evm.TxStatus_TX_STATUS_SUCCESS {
		errorMsg := "unknown error"
		if resp.ErrorMessage != nil {
			errorMsg = *resp.ErrorMessage
		}
		return "", fmt.Errorf("fulfillment tx failed with status %v: %s", resp.TxStatus, errorMsg)
	}
	if resp.ReceiverContractExecutionStatus != nil &&
		*resp.ReceiverContractExecutionStatus != evm.ReceiverContractExecutionStatus_RECEIVER_CONTRACT_EXECUTION_STATUS_SUCCESS {
		return "", fmt.Errorf("receiver contract execution failed for request %s", requestId)
	}

	txHash := common.BytesToHash(resp.TxHash).Hex()
	logger.Info("fulfilled request on-chain", "requestId", requestId, "txHash", txHash)
	return fmt.Sprintf("Fulfilled %s — tx: %s", requestId, txHash), nil
}

// encodeRandomnessReport ABI-encodes (uint256 requestId, uint256 randomWord),
// matching abi.decode(report, (uint256, uint256)) in _processReport.
func encodeRandomnessReport(requestId, randomWord *big.Int) ([]byte, error) {
	uint256Ty, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, err
	}
	args := abi.Arguments{{Type: uint256Ty}, {Type: uint256Ty}}
	return args.Pack(requestId, randomWord)
}
