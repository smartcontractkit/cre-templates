package main

import (
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

// ─── Config ─────────────────────────────────────────────────
type Config struct {
	Schedule       string `json:"schedule"`
	URL            string `json:"url"`
	SecretID       string `json:"secretId"`
	ScoreThreshold uint64 `json:"scoreThreshold"`
}

// ─── Logic to be executed over confidential data ────────────
// Some logic needs to be computed over sensitive data while preserving the
// confidentiality of that data from node operators: risk thresholds, API
// credentials, centralised exchange stablecoin reserves for reasoning, identity
// details. Leaking this data could have adverse effects, including enabling
// front-running attacks, exposing sensitive financial information, and
// compromising individual privacy.
//
// Note what is and is not confidential here: a confidential workflow, despite
// running inside the enclave, is part of the binary the Workflow DON provides to
// the enclave — so the binary, including this logic, is revealed. What the
// enclave keeps confidential is the data this logic computes over: Vault DON
// secrets, the request and response payloads of HTTP calls made from the
// enclave, and other intermediate values.
//
// Keep it deterministic for a given input — the enclave result is attested and
// verified by DON consensus before the workflow completes.
func scoreResponse(body string) uint64 {
	var score uint64
	for _, b := range []byte(body) {
		score = (score + uint64(b)) % 1000
	}
	return score
}

// ─── TEE Cron Callback ──────────────────────────────────────
// Receives a cre.TeeRuntime, not a cre.Runtime. Everything here runs inside the
// enclave until we explicitly cross back with UsingTheDons().
func onCronTrigger(config *Config, runtime cre.TeeRuntime, _ *cron.Payload) (string, error) {
	// ── Step 2: Fetch a secret inside the enclave ──
	// The Vault DON releases this secret only into an attested enclave, and it is
	// decrypted at the moment GetSecret() runs. There is nothing to declare
	// upfront (unlike Confidential HTTP's VaultDonSecrets).
	secret, err := runtime.GetSecret(&cre.SecretRequest{Id: config.SecretID}).Await()
	if err != nil {
		return "", fmt.Errorf("failed to fetch secret %q inside the enclave: %w", config.SecretID, err)
	}
	apiToken := secret.Value

	// ── Step 3: Make a capability call from inside the enclave ──
	// SendRequestInTee takes a cre.TeeRuntime, so the request executes from inside
	// the enclave, keeping the request and response payloads confidential from node
	// operators. The Workflow DON offers consensus verification of enclave
	// attestations, proving the integrity of the logic executed within the enclave.
	//
	// Note: do NOT reach for the confidentialhttp client here — its SendRequest
	// only accepts a cre.Runtime and it is not meant to be called from a TEE
	// handler.
	client := &http.Client{}
	response, err := client.SendRequestInTee(runtime, &http.Request{
		Url:    config.URL,
		Method: "GET",
		MultiHeaders: map[string]*http.HeaderValues{
			"Authorization": {Values: []string{"Bearer " + apiToken}},
		},
	}).Await()
	if err != nil {
		return "", fmt.Errorf("confidential request failed: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return "", fmt.Errorf("confidential request failed with status: %d", response.StatusCode)
	}

	body := strings.TrimSpace(string(response.Body))

	// The default endpoint echoes the request headers back, so we can confirm the
	// secret really was injected inside the enclave — as a boolean, never by
	// logging the token itself. Drop this once URL points at a real API.
	secretReachedAPI := strings.Contains(body, apiToken)

	// Decision logic executed over the confidential response payload.
	score := scoreResponse(body)
	verdict := "REJECT"
	if score >= config.ScoreThreshold {
		verdict = "APPROVE"
	}

	// ⚠️ Logs should be used for simulations only, and MUST be removed before
	// deploying to production to preserve the confidentiality offered by enclaves.
	// Avoid logging inside the enclave in general — sensitive or not.
	runtime.Logger().Info("Enclave computation complete", "verdict", verdict)

	// ── Step 4: Cross back to the DON for anything that needs consensus ──
	// UsingTheDons() returns a regular cre.Runtime. Anything passed into a
	// capability call on it executes on Workflow DON nodes and is NO LONGER
	// confidential — so we cross over the verdict and score only, never the
	// secret or the raw response body.
	donRuntime := runtime.UsingTheDons()

	encodedPayload, err := encodeVerdict(verdict, score)
	if err != nil {
		return "", fmt.Errorf("failed to encode report payload: %w", err)
	}

	if _, err := donRuntime.GenerateReport(&cre.ReportRequest{
		EncodedPayload: encodedPayload,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}).Await(); err != nil {
		return "", fmt.Errorf("failed to generate report on the DON: %w", err)
	}

	// The signed report is now a normal CRE report. To deliver it on-chain, pass
	// it to evmClient.WriteReport(donRuntime, ...) — see the Custom Data Feed
	// template for the full write path.
	return fmt.Sprintf("%s (score: %d, secret reached API: %t)", verdict, score, secretReachedAPI), nil
}

// encodeVerdict ABI-encodes (string verdict, uint256 score) so a consumer
// contract can decode it with abi.decode(report, (string, uint256)).
func encodeVerdict(verdict string, score uint64) ([]byte, error) {
	stringType, err := abi.NewType("string", "", nil)
	if err != nil {
		return nil, err
	}
	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, err
	}
	args := abi.Arguments{{Type: stringType}, {Type: uint256Type}}
	return args.Pack(verdict, new(big.Int).SetUint64(score))
}

// ─── Workflow Init ──────────────────────────────────────────
func InitWorkflow(config *Config, _ *slog.Logger, _ cre.SecretsProvider) (cre.Workflow[*Config], error) {
	return cre.Workflow[*Config]{
		// ── Step 1: Register a TEE handler ──
		// cre.HandlerInTee instead of cre.Handler. The third argument is a
		// cre.TeeConstraint describing which enclaves this handler will accept.
		//
		// Alternatives:
		//   cre.AnyTee{}                                              — any registered TEE, any region
		//   cre.AnyTeeInRegions{Regions: []cre.Region{cre.AwsUsWest2}} — any TEE, restricted to a region
		//
		// Each TEE binding owns its own region enum (Nitro <-> NitroRegion), so a
		// region for the wrong TEE is a compile-time error. AWS Nitro in us-west-2
		// is currently the only registered TEE type and region.
		cre.HandlerInTee(
			cron.Trigger(&cron.Config{Schedule: config.Schedule}),
			onCronTrigger,
			cre.OneOfTees{cre.Nitro{Regions: []cre.NitroRegion{cre.NitroUsWest2}}},
		),
	}, nil
}
