package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	pb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/bindings"
	evmmock "github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/mock"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	httpmock "github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http/mock"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre/testutils"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"indexer-workflow-ts/contracts/evm/src/generated/balance_reader"
	"indexer-workflow-ts/contracts/evm/src/generated/ierc20"
	"indexer-workflow-ts/contracts/evm/src/generated/message_emitter"
)

var anyExecutionTime = time.Unix(1752514917, 0)

func TestInitWorkflow(t *testing.T) {
	config := makeTestConfig(t)
	runtime := testutils.NewRuntime(t, testutils.Secrets{})

	workflow, err := InitWorkflow(config, runtime.Logger(), nil)
	require.NoError(t, err)

	require.Len(t, workflow, 2) // cron, log triggers
	require.Equal(t, cron.Trigger(&cron.Config{}).CapabilityID(), workflow[0].CapabilityID())
}

func TestOnCronTrigger(t *testing.T) {
	config := makeTestConfig(t)
	runtime := testutils.NewRuntime(t, testutils.Secrets{
		"": {},
	})

	// Mock HTTP client for POR data
	httpMock, err := httpmock.NewClientCapability(t)
	require.NoError(t, err)
	httpMock.SendRequest = func(ctx context.Context, input *http.Request) (*http.Response, error) {
		// Return mock POR response
		porResponse := `{
			"accountName": "TrueUSD",
			"totalTrust": 1000000.0,
			"totalToken": 1000000.0,
			"ripcord": false,
			"updatedAt": "2023-01-01T00:00:00Z"
		}`
		return &http.Response{Body: []byte(porResponse)}, nil
	}

	// Mock EVM client
	chainSelector, err := config.EVMs[0].GetChainSelector()
	require.NoError(t, err)
	evmMock, err := evmmock.NewClientCapability(chainSelector, t)
	require.NoError(t, err)

	// Set up contract mocks using generated mock contracts
	evmCfg := config.EVMs[0]

	// Mock BalanceReader for fetchNativeTokenBalance
	balanceReaderMock := balance_reader.NewBalanceReaderMock(
		common.HexToAddress(evmCfg.BalanceReaderAddress),
		evmMock,
	)
	balanceReaderMock.GetNativeBalances = func(input balance_reader.GetNativeBalancesInput) ([]*big.Int, error) {
		// Return mock balance for each address (same number as input addresses)
		balances := make([]*big.Int, len(input.Addresses))
		for i := range input.Addresses {
			balances[i] = big.NewInt(500000000000000000) // 0.5 ETH in wei
		}
		return balances, nil
	}

	// Mock IERC20 for getTotalSupply
	ierc20Mock := ierc20.NewIERC20Mock(
		common.HexToAddress(evmCfg.TokenAddress),
		evmMock,
	)
	ierc20Mock.TotalSupply = func() (*big.Int, error) {
		return big.NewInt(1000000000000000000), nil // 1 token with 18 decimals
	}

	// Note: ReserveManager WriteReportFromUpdateReserves is not a read method,
	// so it's handled by the EVM mock transaction system directly
	evmMock.WriteReport = func(ctx context.Context, input *evm.WriteReportRequest) (*evm.WriteReportReply, error) {
		return &evm.WriteReportReply{
			TxHash: common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef").Bytes(),
		}, nil
	}

	result, err := onPORCronTrigger(config, runtime, &cron.Payload{
		ScheduledExecutionTime: timestamppb.New(anyExecutionTime),
	})

	require.NoError(t, err)
	require.NotNil(t, result)

	// Check that the result contains the expected reserve value
	require.Equal(t, "1000000", result) // Should match the totalToken from mock response

	// Verify expected log messages
	logs := runtime.GetLogs()
	assertLogContains(t, logs, `msg="fetching por"`)
	assertLogContains(t, logs, `msg=ReserveInfo`)
	assertLogContains(t, logs, `msg=TotalSupply`)
	assertLogContains(t, logs, `msg=TotalReserveScaled`)
	assertLogContains(t, logs, `msg="Native token balance"`)
}

func TestOnLogTrigger(t *testing.T) {
	config := makeTestConfig(t)
	runtime := testutils.NewRuntime(t, testutils.Secrets{})

	// Mock EVM client
	chainSelector, err := config.EVMs[0].GetChainSelector()
	require.NoError(t, err)
	evmMock, err := evmmock.NewClientCapability(chainSelector, t)
	require.NoError(t, err)

	// Mock MessageEmitter for log trigger
	evmCfg := config.EVMs[0]
	messageEmitterMock := message_emitter.NewMessageEmitterMock(
		common.HexToAddress(evmCfg.MessageEmitterAddress),
		evmMock,
	)
	messageEmitterMock.GetLastMessage = func(input message_emitter.GetLastMessageInput) (string, error) {
		return "Test message from contract", nil
	}

	msgEmitterAbi, err := message_emitter.MessageEmitterMetaData.GetAbi()
	require.NoError(t, err)
	eventData, err := abi.Arguments{msgEmitterAbi.Events["MessageEmitted"].Inputs[2]}.Pack("Test message from contract")
	require.NoError(t, err, "Encoding event data should not return an error")
	// Create a mock log payload
	mockLog := &evm.Log{
		Topics: [][]byte{
			common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678901234").Bytes(), // event signature
			common.HexToHash("0x000000000000000000000000abcdefabcdefabcdefabcdefabcdefabcdefabcd").Bytes(), // emitter address (padded)
			common.HexToHash("0x000000000000000000000000000000000000000000000000000000006716eb80").Bytes(), // additional topic
		},
		Data:        eventData, // this is not used by the test as we pass in mockLogDecoded, but encoding here for consistency
		BlockNumber: pb.NewBigIntFromInt(big.NewInt(100)),
	}

	mockLogDecoded := &bindings.DecodedLog[message_emitter.MessageEmittedDecoded]{
		Log: mockLog,
		Data: message_emitter.MessageEmittedDecoded{
			Emitter:   common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
			Message:   "Test message from contract",
			Timestamp: big.NewInt(100),
		},
	}

	result, err := onLogTrigger(config, runtime, mockLogDecoded)
	require.NoError(t, err)
	require.Equal(t, "Test message from contract", result)

	// Verify expected log messages
	logs := runtime.GetLogs()
	assertLogContains(t, logs, `msg="Message retrieved from the contract"`)
	assertLogContains(t, logs, `blockNumber=100`)
}

//go:embed config/config.production.json
var configJson []byte

func makeTestConfig(t *testing.T) *Config {
	config := &Config{}
	require.NoError(t, json.Unmarshal(configJson, config))
	return config
}

func assertLogContains(t *testing.T, logs [][]byte, substr string) {
	for _, line := range logs {
		if strings.Contains(string(line), substr) {
			return
		}
	}
	t.Fatalf("Expected logs to contain substring %q, but it was not found in logs:\n%s",
		substr, strings.Join(func() []string {
			var logStrings []string
			for _, log := range logs {
				logStrings = append(logStrings, string(log))
			}
			return logStrings
		}(), "\n"))
}
