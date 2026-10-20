package main

import (
	"fmt"
	"log/slog"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

// Config mirrors config.staging.json / config.production.json.
type Config struct {
	Schedule string `json:"schedule"`
}

// Run in an AWS Nitro enclave in us-west-2 (staging confidential-workflows
// region). Update this when other TEE types or regions become available.
var teeRequirements = cre.OneOfTees{cre.Nitro{Regions: []cre.NitroRegion{cre.NitroUsWest2}}}

// InitWorkflow registers the cron handler as a TEE handler.
func InitWorkflow(cfg *Config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[*Config], error) {
	if cfg.Schedule == "" {
		return cre.Workflow[*Config]{}, fmt.Errorf("config requires schedule")
	}
	return cre.Workflow[*Config]{
		cre.HandlerInTee(
			cron.Trigger(&cron.Config{Schedule: cfg.Schedule}),
			onCronTrigger,
			teeRequirements,
		),
	}, nil
}

// onCronTrigger runs inside the enclave and generates a fresh random uint32.
func onCronTrigger(cfg *Config, rt cre.TeeRuntime, _ *cron.Payload) (string, error) {
	// Recommended: TEE-local randomness. Each execution gets its own seed, so
	// values do not repeat across runs or across concurrent executions.
	rand, err := rt.Rand()
	if err != nil {
		return "", fmt.Errorf("failed to get random value: %w", err)
	}
	value := rand.Uint32()
	rt.Logger().Info("random-number-generated", "random_value", value)
	return fmt.Sprintf("%d", value), nil
}
