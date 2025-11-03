// This workflow compares supply APYs across multiple EVM chains and
// rebalances supply to the chain with the highest APY via cross-chain
// transfers.
//
// It requires at least two EVM configurations to compare APYs.
//
// In the event of tie, the first chain with the highest APY is chosen.
package main

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"

	"cre-multi-chain-token-manager-go/contracts/evm/src/generated/mock_pool"
	"cre-multi-chain-token-manager-go/contracts/evm/src/generated/protocol_smart_wallet"
)

// EVMConfig holds per-chain configuration.
type EVMConfig struct {
	AssetAddress               string `json:"assetAddress"`
	PoolAddress                string `json:"poolAddress"`
	ProtocolSmartWalletAddress string `json:"protocolSmartWalletAddress"`
	ChainName                  string `json:"chainName"`
	GasLimit                   uint64 `json:"gasLimit"`
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
	Schedule                string      `json:"schedule"`
	MinBPSDeltaForRebalance int64       `json:"minBPSDeltaForRebalance"`
	EVMs                    []EVMConfig `json:"evms"`
}

func InitWorkflow(config *Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*Config], error) {
	if len(config.EVMs) < 2 {
		return nil, fmt.Errorf("at least two EVM configurations are required to compare supply APYs")
	}
	if config.MinBPSDeltaForRebalance < 0 {
		return nil, fmt.Errorf("minBPSDeltaForRebalance must be >= 0")
	}
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
	return doHighestSupplyAPY(config, runtime, outputs.ScheduledExecutionTime.AsTime())
}

func doHighestSupplyAPY(config *Config, runtime cre.Runtime, runTime time.Time) (string, error) {
	logger := runtime.Logger()

	if len(config.EVMs) < 2 {
		return "", fmt.Errorf("at least two EVM configurations are required to compare supply APYs")
	}

	// First determine the chain with the highest supply APY.

	logger.Info("Reading supply APYs...")

	type pool struct {
		ChainName                  string
		APR                        *big.Int
		APY                        float64
		ProtocolSmartWalletAddress common.Address
		Balance                    *big.Int
	}

	chainNameToPool := make(map[string]pool)
	maxAPYPool := pool{}
	blockNumber := big.NewInt(rpc.LatestBlockNumber.Int64()) // warn: use finalized or safe in prod

	for _, evmCfg := range config.EVMs {

		logger.Info("Reading APY", "chainName", evmCfg.ChainName, "poolAddress", evmCfg.PoolAddress, "assetAddress", common.HexToAddress(evmCfg.AssetAddress).Hex())

		evmClient, err := evmCfg.NewEVMClient()
		if err != nil {
			return "", fmt.Errorf("failed to create EVM client for %s: %w", evmCfg.ChainName, err)
		}

		mockPool, err := mock_pool.NewMockPool(
			evmClient,
			common.HexToAddress(evmCfg.PoolAddress),
			nil)
		if err != nil {
			return "", fmt.Errorf("failed to create mock pool for %s: %w", evmCfg.ChainName, err)
		}

		reserveData, err := mockPool.GetReserveData(
			runtime,
			mock_pool.GetReserveDataInput{
				Asset: common.HexToAddress(evmCfg.AssetAddress),
			},
			blockNumber,
		).Await()
		if err != nil {
			return "", fmt.Errorf("failed to get reserve data: %w", err)
		}

		logger.Info("Reserve data", "chainName", evmCfg.ChainName, "currentLiquidityRate", reserveData.CurrentLiquidityRate)

		balanceInPool, err := mockPool.BalanceOf(
			runtime,
			mock_pool.BalanceOfInput{
				Arg0: common.HexToAddress(evmCfg.ProtocolSmartWalletAddress),
				Arg1: common.HexToAddress(evmCfg.AssetAddress),
			},
			blockNumber,
		).Await()
		if err != nil {
			return "", fmt.Errorf("failed to get supply APY: %w", err)
		}

		logger.Info("Balance in pool", "chainName", evmCfg.ChainName, "balanceInPool", balanceInPool)

		currentLiquidity := reserveData.CurrentLiquidityRate
		supplyAPY := aprInRAYToAPY(currentLiquidity)
		logger.Info(
			"Supply yield",
			"chainName",
			evmCfg.ChainName,
			"apy_percent", supplyAPY*100,
			"apr_percent", aprInRAYToAPR(currentLiquidity)*100)

		curPool := pool{
			ChainName:                  evmCfg.ChainName,
			APR:                        currentLiquidity,
			APY:                        supplyAPY,
			ProtocolSmartWalletAddress: common.HexToAddress(evmCfg.ProtocolSmartWalletAddress),
			Balance:                    balanceInPool,
		}
		chainNameToPool[evmCfg.ChainName] = curPool

		if maxAPYPool.APR == nil || curPool.APR.Cmp(maxAPYPool.APR) > 0 {
			// New best APY found.
			maxAPYPool = curPool
		} else if maxAPYPool.APR != nil && curPool.APR.Cmp(maxAPYPool.APR) == 0 {
			// Tie: keep the first one found.
			logger.Info("Found tie in APY, keeping existing best chain", "existing_chain", maxAPYPool.ChainName, "new_chain", curPool.ChainName)
		}
	}

	if maxAPYPool.APR == nil || maxAPYPool.APR.Cmp(big.NewInt(0)) != 1 {
		return "", fmt.Errorf("best APY unset or <= 0")
	}

	logger.Info(
		"Found best APY",
		"best_apy", maxAPYPool.APY,
		"best_apy_chain_name", maxAPYPool.ChainName,
	)

	// Second, withdraw on each chain with suboptimal APY and deposit to the
	// best APY chain via CCIP.

	bestChainSelector, err := evm.ChainSelectorFromName(maxAPYPool.ChainName)
	if err != nil {
		return "", fmt.Errorf("failed to get chain selector for best APY chain %s: %w", maxAPYPool.ChainName, err)
	}

	newMaxAPYPoolBalance := new(big.Int).Set(maxAPYPool.Balance)

	for _, evmCfg := range config.EVMs {
		if evmCfg.ChainName == maxAPYPool.ChainName {
			// Skip the winning chain.
			continue
		}

		logger.Info(
			"Rebalancing assets to best APY chain",
			"from_chain", evmCfg.ChainName,
			"to_chain", maxAPYPool.ChainName,
			"to_chain_selector", bestChainSelector,
			"to_protocol_smart_wallet", maxAPYPool.ProtocolSmartWalletAddress.Hex())

		// First, check the balance to see if there is anything to rebalance.

		evmClient, err := evmCfg.NewEVMClient()
		if err != nil {
			return "", fmt.Errorf("failed to create EVM client for %s: %w", evmCfg.ChainName, err)
		}

		mockPool, err := mock_pool.NewMockPool(
			evmClient,
			common.HexToAddress(evmCfg.PoolAddress),
			nil)
		if err != nil {
			return "", fmt.Errorf("failed to create mock pool for %s: %w", evmCfg.ChainName, err)
		}

		balanceInPool, err := mockPool.BalanceOf(
			runtime,
			mock_pool.BalanceOfInput{
				Arg0: common.HexToAddress(evmCfg.ProtocolSmartWalletAddress),
				Arg1: common.HexToAddress(evmCfg.AssetAddress),
			},
			blockNumber,
		).Await()
		if err != nil {
			return "", fmt.Errorf("failed to get supply APY: %w", err)
		}

		if balanceInPool.Cmp(big.NewInt(0)) == 0 {
			logger.Info("No balance to rebalance from this chain, skipping...", "chain", evmCfg.ChainName)
			continue
		}

		curPool, ok := chainNameToPool[evmCfg.ChainName]
		if !ok {
			return "", fmt.Errorf("pool info not found for chain %s", evmCfg.ChainName)
		}

		// Check if the APY difference is above the minimum delta to rebalance
		// threshold.

		diff := new(big.Int).Sub(maxAPYPool.APR, curPool.APR)
		// APR/(1e27/1e4)=APR/1e23 to convert from ray to basis points.
		decimals, ok := new(big.Int).SetString("100000000000000000000000", 10)
		if !ok {
			return "", errors.New("failed to parse decimals")
		}
		// Convert from ray to basis points; e.g., 5% = 500 bps.
		diffBps := new(big.Int).Div(diff, decimals)

		if diffBps.Cmp(big.NewInt(config.MinBPSDeltaForRebalance)) < 0 {
			logger.Info(
				"APY difference below minimum delta for rebalance, skipping...",
				"from_chain", evmCfg.ChainName,
				"to_chain", maxAPYPool.ChainName,
				"diff_bps", diffBps,
				"min_bps_delta_for_rebalance", config.MinBPSDeltaForRebalance,
			)
			continue
		}

		// Withdraw from the suboptimal APY chain and deposit on the the best
		// APY chain via CCIP.

		logger.Info(
			"Rebalancing supply",
			"from_chain", evmCfg.ChainName,
			"balanceInPool", balanceInPool,
			"to_chain", maxAPYPool.ChainName,
			"diff_bps", diffBps,
			"min_bps_delta_for_rebalance", config.MinBPSDeltaForRebalance,
		)

		protocolSmartWallet, err := protocol_smart_wallet.NewProtocolSmartWallet(
			evmClient,
			common.HexToAddress(evmCfg.ProtocolSmartWalletAddress),
			nil,
		)
		if err != nil {
			return "", fmt.Errorf("failed to create protocol smart wallet for %s: %w", evmCfg.ChainName, err)
		}

		resp, err := protocolSmartWallet.WriteReportFromRebalanceParams(
			runtime,
			protocol_smart_wallet.RebalanceParams{
				Asset:                          common.HexToAddress(evmCfg.AssetAddress),
				Amount:                         balanceInPool,
				DestinationChainSelector:       bestChainSelector,
				DestinationProtocolSmartWallet: maxAPYPool.ProtocolSmartWalletAddress,
			},
			&evm.GasConfig{GasLimit: evmCfg.GasLimit},
		).Await()
		if err != nil {
			return "", fmt.Errorf("failed to submit report: %w", err)
		}

		// Check transaction status
		if resp.TxStatus != evm.TxStatus_TX_STATUS_SUCCESS {
			errorMsg := "unknown error"
			if resp.ErrorMessage != nil {
				errorMsg = *resp.ErrorMessage
			}
			return "", fmt.Errorf("transaction failed with status %v: %s", resp.TxStatus, errorMsg)
		} else if resp.ReceiverContractExecutionStatus != nil &&
			*resp.ReceiverContractExecutionStatus != evm.ReceiverContractExecutionStatus_RECEIVER_CONTRACT_EXECUTION_STATUS_SUCCESS {
			return "", fmt.Errorf("failed to execute receiver contract")
		}

		logger.Info("Write report transaction succeeded at", "txHash", common.BytesToHash(resp.TxHash).Hex(), "chainName", evmCfg.ChainName)

		newMaxAPYPoolBalance.Add(newMaxAPYPoolBalance, balanceInPool)
	}

	logger.Info(
		"Rebalancing completed",
		"old_balance", maxAPYPool.Balance,
		"new_balance", newMaxAPYPoolBalance,
		"chain_name", maxAPYPool.ChainName,
		"amount_rebalanced", new(big.Int).Sub(newMaxAPYPoolBalance, maxAPYPool.Balance),
	)

	return "", nil
}

func aprInRAYToAPR(apr *big.Int) float64 {
	aprFloat, _ := new(big.Float).Quo(new(big.Float).SetInt(apr), big.NewFloat(1e27)).Float64()
	return aprFloat
}

func aprInRAYToAPY(apr *big.Int) float64 {
	// APY = e^(APR) - 1
	exp, _ := new(big.Float).Quo(new(big.Float).SetInt(apr), big.NewFloat(1e27)).Float64()
	apy := math.Pow(math.E, exp) - 1
	return apy
}
