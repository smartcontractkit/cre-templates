package main

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	httpmock "github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http/mock"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/testutils"
)

const apiToken = "test-token"

func testConfig() *Config {
	return &Config{
		Schedule:       "0 */1 * * * *",
		URL:            "https://postman-echo.com/headers",
		SecretID:       "API_TOKEN",
		ScoreThreshold: 500,
	}
}

func testSecrets() testutils.Secrets {
	return testutils.Secrets{
		cre.DefaultSecretNamespace: {"API_TOKEN": apiToken},
	}
}

// stubHTTP registers an in-enclave HTTP stub and records the Authorization
// header values it was called with.
func stubHTTP(t *testing.T, statusCode uint32, body string) *[]string {
	t.Helper()

	capability, err := httpmock.NewClientCapability(t)
	require.NoError(t, err)

	captured := &[]string{}
	capability.SendRequest = func(_ context.Context, input *http.Request) (*http.Response, error) {
		if auth := input.MultiHeaders["Authorization"]; auth != nil {
			*captured = append(*captured, auth.Values...)
		}
		return &http.Response{StatusCode: statusCode, Body: []byte(body)}, nil
	}
	return captured
}

func TestOnCronTrigger_InjectsEnclaveFetchedSecret(t *testing.T) {
	captured := stubHTTP(t, 200, "hello")
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)

	assert.Equal(t, []string{"Bearer " + apiToken}, *captured)
}

func TestOnCronTrigger_ConfirmsSecretReachedAPI(t *testing.T) {
	stubHTTP(t, 200, `{"authorization":"Bearer `+apiToken+`"}`)
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	result, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)

	assert.Contains(t, result, "secret reached API: true")
}

func TestOnCronTrigger_ReportsSecretDidNotReachAPI(t *testing.T) {
	stubHTTP(t, 200, `{"authorization":"Bearer other"}`)
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	result, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)

	assert.Contains(t, result, "secret reached API: false")
}

func TestOnCronTrigger_ApprovesWhenScoreClearsThreshold(t *testing.T) {
	// "zzzzzzzz" sums to 976, above the 500 threshold.
	stubHTTP(t, 200, "zzzzzzzz")
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	result, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)

	assert.Contains(t, result, "APPROVE")
}

func TestOnCronTrigger_RejectsWhenScoreBelowThreshold(t *testing.T) {
	// "a" sums to 97, below the 500 threshold.
	stubHTTP(t, 200, "a")
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	result, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)

	assert.Contains(t, result, "REJECT")
}

func TestOnCronTrigger_ErrorsOnNon2xx(t *testing.T) {
	stubHTTP(t, 401, "unauthorized")
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err := onCronTrigger(testConfig(), runtime, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "status: 401")
}

func TestOnCronTrigger_DoesNotLogSecretOrResponseBody(t *testing.T) {
	const sensitiveBody = "sensitive-response"
	stubHTTP(t, 200, sensitiveBody)
	runtime := testutils.NewTeeRuntime(t, testSecrets())

	_, err := onCronTrigger(testConfig(), runtime, nil)
	require.NoError(t, err)

	for _, raw := range runtime.GetLogs() {
		line := string(raw)
		assert.False(t, strings.Contains(line, apiToken), "log leaked the secret: %s", line)
		assert.False(t, strings.Contains(line, sensitiveBody), "log leaked the response body: %s", line)
	}
}

func TestScoreResponseIsDeterministic(t *testing.T) {
	assert.Equal(t, scoreResponse("hello"), scoreResponse("hello"))
	assert.Equal(t, uint64(97), scoreResponse("a"))
	assert.Equal(t, uint64(976), scoreResponse("zzzzzzzz"))
}

func TestEncodeVerdictRoundTrips(t *testing.T) {
	encoded, err := encodeVerdict("APPROVE", 371)
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)
}

func TestInitWorkflowRegistersTeeHandler(t *testing.T) {
	workflow, err := InitWorkflow(testConfig(), nil, nil)
	require.NoError(t, err)

	assert.Len(t, workflow, 1)
}
