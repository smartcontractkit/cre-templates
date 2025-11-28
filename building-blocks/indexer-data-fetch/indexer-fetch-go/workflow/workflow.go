package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

type Config struct {
	Schedule        string                 `json:"schedule"`
	GraphqlEndpoint string                 `json:"graphqlEndpoint"`
	Query           string                 `json:"query"`
	Variables       map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []interface{}   `json:"errors,omitempty"`
}

func InitWorkflow(config *Config, logger *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[*Config], error) {
	cronTriggerCfg := &cron.Config{
		Schedule: config.Schedule,
	}

	return cre.Workflow[*Config]{
		cre.Handler(
			cron.Trigger(cronTriggerCfg),
			onIndexerCronTrigger,
		),
	}, nil
}

func onIndexerCronTrigger(config *Config, runtime cre.Runtime, _ *cron.Payload) (string, error) {
	logger := runtime.Logger()
	timestamp := time.Now().UTC().Format(time.RFC3339)

	logger.Info("Cron triggered", "timestamp", timestamp)
	logger.Info("Querying The Graph indexer", "endpoint", config.GraphqlEndpoint)

	// Fetch data from The Graph using SendRequest pattern
	client := &http.Client{}
	logger.Info("setup client")
	result, err := http.SendRequest(config, runtime, client, fetchGraphData, cre.ConsensusIdenticalAggregation[string]()).Await()
	if err != nil {
		logger.Error("Failed to fetch indexer data", "err", err)
		return "", err
	}

	logger.Info("Indexer data fetched successfully", "timestamp", timestamp)

	// Format output
	output := map[string]interface{}{
		"timestamp": timestamp,
		"endpoint":  config.GraphqlEndpoint,
		"data":      json.RawMessage(result),
	}

	// Return a JSON string
	out, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func fetchGraphData(config *Config, logger *slog.Logger, sendRequester *http.SendRequester) (string, error) {

	// Prepare GraphQL request
	gqlRequest := GraphQLRequest{
		Query:     config.Query,
		Variables: config.Variables,
	}

	requestBody, err := json.Marshal(gqlRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make POST request
	httpResp, err := sendRequester.SendRequest(&http.Request{
		Method: "POST",
		Url:    config.GraphqlEndpoint,
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer bca58895bc60dcb319e3cbdfd989b964",
		},
		Body: requestBody,
	}).Await()

	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}

	// Parse response
	var gqlResponse GraphQLResponse
	if err := json.Unmarshal(httpResp.Body, &gqlResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for GraphQL errors
	if len(gqlResponse.Errors) > 0 {
		errJSON, _ := json.Marshal(gqlResponse.Errors)
		logger.Error("GraphQL errors", "errors", string(errJSON))
		return "", fmt.Errorf("GraphQL query failed: %s", string(errJSON))
	}

	if gqlResponse.Data == nil {
		return "", errors.New("no data returned from GraphQL query")
	}

	logger.Info("Successfully fetched data from indexer")

	return string(gqlResponse.Data), nil
}
