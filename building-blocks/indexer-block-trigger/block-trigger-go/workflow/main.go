//go:build wasip1

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

// Workflow configuration loaded from the config.json file
type Config struct {
	WatchedAddresses []string `json:"watchedAddresses"`
}

// AlchemyWebhookPayload represents the webhook data from Alchemy
type AlchemyWebhookPayload struct {
	WebhookID string `json:"webhookId"`
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	Type      string `json:"type"`
	Event     struct {
		Data struct {
			Block struct {
				Hash      string `json:"hash"`
				Number    int64  `json:"number"`
				Timestamp int64  `json:"timestamp"`
				Logs      []struct {
					Data    string   `json:"data"`
					Topics  []string `json:"topics"`
					Index   int      `json:"index"`
					Account struct {
						Address string `json:"address"`
					} `json:"account"`
					Transaction struct {
						Hash                 string  `json:"hash"`
						Nonce                int     `json:"nonce"`
						Index                int     `json:"index"`
						From                 Address `json:"from"`
						To                   *Address `json:"to"`
						Value                string  `json:"value"`
						GasPrice             string  `json:"gasPrice"`
						MaxFeePerGas         *string `json:"maxFeePerGas"`
						MaxPriorityFeePerGas *string `json:"maxPriorityFeePerGas"`
						Gas                  int     `json:"gas"`
						Status               int     `json:"status"`
						GasUsed              int     `json:"gasUsed"`
						CumulativeGasUsed    int     `json:"cumulativeGasUsed"`
						EffectiveGasPrice    string  `json:"effectiveGasPrice"`
						CreatedContract      *Address `json:"createdContract"`
					} `json:"transaction"`
				} `json:"logs"`
			} `json:"block"`
		} `json:"data"`
	} `json:"event"`
}

type Address struct {
	Address string `json:"address"`
}

// Transaction represents a processed transaction
type Transaction struct {
	Hash              string  `json:"hash"`
	Nonce             int     `json:"nonce"`
	Index             int     `json:"index"`
	From              string  `json:"from"`
	To                *string `json:"to"`
	Value             string  `json:"value"`
	GasPrice          string  `json:"gasPrice"`
	Gas               int     `json:"gas"`
	Status            int     `json:"status"`
	GasUsed           int     `json:"gasUsed"`
	BlockNumber       int64   `json:"blockNumber"`
	BlockHash         string  `json:"blockHash"`
	Timestamp         int64   `json:"timestamp"`
}

// ExecutionResult represents the workflow output
type ExecutionResult struct {
	BlockNumber         int64         `json:"blockNumber"`
	BlockHash           string        `json:"blockHash"`
	Timestamp           int64         `json:"timestamp"`
	TotalLogs           int           `json:"totalLogs"`
	UniqueTransactions  int           `json:"uniqueTransactions"`
	MatchedTransactions int           `json:"matchedTransactions"`
	Transactions        []Transaction `json:"transactions"`
}

// TransactionStore is a simple in-memory database mock
type TransactionStore struct {
	watchedAddresses map[string]bool
	transactions     map[string]Transaction
}

func createTransactionStore(addresses []string) *TransactionStore {
	store := &TransactionStore{
		watchedAddresses: make(map[string]bool),
		transactions:     make(map[string]Transaction),
	}

	for _, addr := range addresses {
		store.watchedAddresses[strings.ToLower(addr)] = true
	}

	return store
}

func (s *TransactionStore) isWatchedAddress(address *string) bool {
	if address == nil {
		return false
	}
	return s.watchedAddresses[strings.ToLower(*address)]
}

func (s *TransactionStore) saveTransaction(tx Transaction) {
	s.transactions[tx.Hash] = tx
}

// Workflow implementation with HTTP trigger
func InitWorkflow(config *Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*Config], error) {
	httpTrigger := http.Trigger(&http.Config{})

	return cre.Workflow[*Config]{
		cre.Handler(httpTrigger, onHttpTrigger),
	}, nil
}

func onHttpTrigger(config *Config, runtime cre.Runtime, payload *http.Payload) (*ExecutionResult, error) {
	logger := runtime.Logger()

	// Parse the webhook payload
	var blockData AlchemyWebhookPayload
	if err := json.Unmarshal(payload.Input, &blockData); err != nil {
		logger.Error("Failed to parse webhook payload", "error", err)
		return nil, fmt.Errorf("invalid webhook payload: %w", err)
	}

	block := blockData.Event.Data.Block
	logger.Info("Processing block", "blockNumber", block.Number, "hash", block.Hash)

	// Initialize store with watched addresses from config
	store := createTransactionStore(config.WatchedAddresses)

	// Extract unique transactions from logs
	processedHashes := make(map[string]bool)
	var matchedTransactions []Transaction

	for _, log := range block.Logs {
		tx := log.Transaction

		// Skip if we've already processed this transaction
		if processedHashes[tx.Hash] {
			continue
		}
		processedHashes[tx.Hash] = true

		// Check if the 'to' address matches any watched addresses
		var toAddress *string
		if tx.To != nil {
			toAddress = &tx.To.Address
		}

		if store.isWatchedAddress(toAddress) {
			logger.Info("Match found!", "txHash", tx.Hash, "toAddress", *toAddress)

			// Create transaction record
			transaction := Transaction{
				Hash:        tx.Hash,
				Nonce:       tx.Nonce,
				Index:       tx.Index,
				From:        tx.From.Address,
				To:          toAddress,
				Value:       tx.Value,
				GasPrice:    tx.GasPrice,
				Gas:         tx.Gas,
				Status:      tx.Status,
				GasUsed:     tx.GasUsed,
				BlockNumber: block.Number,
				BlockHash:   block.Hash,
				Timestamp:   block.Timestamp,
			}

			// Save to store
			store.saveTransaction(transaction)
			matchedTransactions = append(matchedTransactions, transaction)
		}
	}

	logger.Info("Block processing complete",
		"block", block.Number,
		"totalLogs", len(block.Logs),
		"uniqueTransactions", len(processedHashes),
		"matchedTransactions", len(matchedTransactions))

	// Return summary
	return &ExecutionResult{
		BlockNumber:         block.Number,
		BlockHash:           block.Hash,
		Timestamp:           block.Timestamp,
		TotalLogs:           len(block.Logs),
		UniqueTransactions:  len(processedHashes),
		MatchedTransactions: len(matchedTransactions),
		Transactions:        matchedTransactions,
	}, nil
}

func main() {
	wasm.NewRunner(cre.ParseJSON[Config]).Run(InitWorkflow)
}