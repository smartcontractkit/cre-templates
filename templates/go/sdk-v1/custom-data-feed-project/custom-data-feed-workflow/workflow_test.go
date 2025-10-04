package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	evmmock "github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/mock"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	httpmock "github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http/mock"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre/testutils"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"cre-v1-project/contracts/evm/src/generated/balance_reader"
	"cre-v1-project/contracts/evm/src/generated/ierc20"
	"cre-v1-project/contracts/evm/src/generated/message_emitter"
)

var anyExecutionTime = time.Unix(1752514917, 0)

func TestInitWorkflow(t *testing.T) {
	config := makeTestConfig(t)
	runtime := testutils.NewRuntime(t, map[string]string{})

	workflow, err := InitWorkflow(config, runtime.Logger(), nil)
	require.NoError(t, err)

	require.Len(t, workflow, 3) // cron, log, and http triggers
	require.Equal(t, cron.Trigger(&cron.Config{}).CapabilityID(), workflow[0].CapabilityID())
}

func TestOnCronTrigger(t *testing.T) {
	config := makeTestConfig(t)
	runtime := testutils.NewRuntime(t, map[string]string{
		"/SECRET_ADDRESS": "0x1234567890123456789012345678901234567890",
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
	assertLogContains(t, logs, `msg="Secret address balance"`)
}

func TestOnLogTrigger(t *testing.T) {
	config := makeTestConfig(t)
	runtime := testutils.NewRuntime(t, map[string]string{})

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

	// Create a mock log payload
	mockLog := &evm.Log{
		Topics: [][]byte{
			common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678901234").Bytes(), // event signature
			common.HexToHash("0x000000000000000000000000abcdefabcdefabcdefabcdefabcdefabcdefabcd").Bytes(), // emitter address (padded)
			common.HexToHash("0x5678567856785678567856785678567856785678567856785678567856785678").Bytes(), // additional topic
		},
		Data: []byte{},
	}

	result, err := onLogTrigger(config, runtime, mockLog)
	require.NoError(t, err)
	require.Equal(t, "Test message from contract", result)

	// Verify expected log messages
	logs := runtime.GetLogs()
	assertLogContains(t, logs, `msg="Message retrieved from the contract"`)
}

func TestOnHTTPTrigger(t *testing.T) {
	config := makeTestConfig(t)
	runtime := testutils.NewRuntime(t, map[string]string{
		"/SECRET_ADDRESS": "0x9876543210987654321098765432109876543210",
	})

	// Mock HTTP client for POR data
	httpMock, err := httpmock.NewClientCapability(t)
	require.NoError(t, err)
	httpMock.SendRequest = func(ctx context.Context, input *http.Request) (*http.Response, error) {
		porResponse := `{
			"accountName": "TrueUSD",
			"totalTrust": 2000000.0,
			"totalToken": 2000000.0,
			"ripcord": false,
			"updatedAt": "2023-01-01T00:00:00Z"
		}`
		return &http.Response{Body: []byte(porResponse)}, nil
	}

	// Mock EVM client (same pattern as cron test)
	chainSelector, err := config.EVMs[0].GetChainSelector()
	require.NoError(t, err)
	evmMock, err := evmmock.NewClientCapability(chainSelector, t)
	require.NoError(t, err)

	// Set up contract mocks (same pattern as cron test)
	evmCfg := config.EVMs[0]

	balanceReaderMock := balance_reader.NewBalanceReaderMock(
		common.HexToAddress(evmCfg.BalanceReaderAddress),
		evmMock,
	)
	balanceReaderMock.GetNativeBalances = func(input balance_reader.GetNativeBalancesInput) ([]*big.Int, error) {
		// Return mock balance for each address (same number as input addresses)
		balances := make([]*big.Int, len(input.Addresses))
		for i := range input.Addresses {
			balances[i] = big.NewInt(750000000000000000) // 0.75 ETH in wei (higher for HTTP test)
		}
		return balances, nil
	}

	ierc20Mock := ierc20.NewIERC20Mock(
		common.HexToAddress(evmCfg.TokenAddress),
		evmMock,
	)
	ierc20Mock.TotalSupply = func() (*big.Int, error) {
		return big.NewInt(2000000000000000000), nil // 2 tokens with 18 decimals (higher for HTTP test)
	}

	evmMock.WriteReport = func(ctx context.Context, input *evm.WriteReportRequest) (*evm.WriteReportReply, error) {
		return &evm.WriteReportReply{
			TxHash: common.HexToHash("0x9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba").Bytes(),
		}, nil
	}

	// Test with empty payload (should default to now)
	result, err := onHTTPTrigger(config, runtime, &http.Payload{Input: []byte{}})
	require.NoError(t, err)
	require.Equal(t, "2000000", result)

	// Test with JSON payload
	httpPayload := HTTPTriggerPayload{
		ExecutionTime: anyExecutionTime,
	}
	payloadBytes, err := json.Marshal(httpPayload)
	require.NoError(t, err)

	result2, err := onHTTPTrigger(config, runtime, &http.Payload{Input: payloadBytes})
	require.NoError(t, err)
	require.Equal(t, "2000000", result2)
}

//go:embed config.json
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
