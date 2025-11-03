//go:build wasip1

package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	protos "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	crehttp "github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
	"google.golang.org/protobuf/types/known/durationpb"
)

// Config struct defines the parameters that can be passed to the workflow.
type Config struct {
	Schedule  string `json:"schedule"`
	AWSRegion string `json:"aws_region"`
	S3Bucket  string `json:"s3_bucket"`
	S3Key     string `json:"s3_key"`
}

// The result of our workflow
type MyResult struct {
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}

// emptyPayloadHash is the SHA256 hash of an empty string.
// AWS SigV4 requires this for GET requests that do not have a body.
// See: https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-header-based-auth.html
const emptyPayloadHash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

// ---------------------------
// S3 Client
// ---------------------------

type S3Client struct {
	cfg         *Config
	nodeRuntime cre.NodeRuntime
	creds       aws.Credentials
	signer      *v4.Signer
	httpClient  *crehttp.Client
	ctx         context.Context
}

func NewS3Client(cfg *Config, nodeRuntime cre.NodeRuntime, creds aws.Credentials, signer *v4.Signer) *S3Client {
	return &S3Client{
		cfg:         cfg,
		nodeRuntime: nodeRuntime,
		creds:       creds,
		signer:      signer,
		httpClient:  &crehttp.Client{},
		ctx:         context.Background(),
	}
}

func (c *S3Client) ReadValue() (string, error) {
	host, path := c.endpoint()
	fullUrl := "https://" + host + path
	method := "GET"
	consensusTime := c.nodeRuntime.Now()

	req, err := http.NewRequest(method, fullUrl, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("x-amz-content-sha256", emptyPayloadHash)

	if err := c.signer.SignHTTP(c.ctx, c.creds, req, emptyPayloadHash, "s3", c.cfg.AWSRegion, consensusTime); err != nil {
		return "", fmt.Errorf("failed to sign request with AWS SDK: %w", err)
	}

	creReq := c.convertReq(req, nil)

	resp, err := c.httpClient.SendRequest(c.nodeRuntime, creReq).Await()
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 404 {
		return "", fmt.Errorf("404 not found")
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("s3 GET failed: %d %s", resp.StatusCode, string(resp.Body))
	}

	return string(resp.Body), nil
}

func (c *S3Client) WriteValue(value string) error {
	host, path := c.endpoint()
	fullUrl := "https://" + host + path
	method := "PUT"
	body := []byte(value)
	consensusTime := c.nodeRuntime.Now()

	req, err := http.NewRequest(method, fullUrl, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(body)))

	payloadHash := sha256Hex(body)
	req.Header.Set("x-amz-content-sha256", payloadHash)

	if err := c.signer.SignHTTP(c.ctx, c.creds, req, payloadHash, "s3", c.cfg.AWSRegion, consensusTime); err != nil {
		return fmt.Errorf("failed to sign request with AWS SDK: %w", err)
	}

	creReq := c.convertReq(req, body)
	creReq.CacheSettings = &crehttp.CacheSettings{
		Store:  true,
		MaxAge: durationpb.New(60 * time.Second),
	}

	resp, err := c.httpClient.SendRequest(c.nodeRuntime, creReq).Await()
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("s3 PUT failed: %d %s", resp.StatusCode, string(resp.Body))
	}
	return nil
}

func (c *S3Client) convertReq(req *http.Request, body []byte) *crehttp.Request {
	creReq := &crehttp.Request{
		Url:     req.URL.String(),
		Method:  req.Method,
		Headers: make(map[string]string),
		Body:    body,
	}

	for k, v := range req.Header {
		if len(v) > 0 {
			creReq.Headers[k] = v[0]
		}
	}
	return creReq
}

func (c *S3Client) endpoint() (string, string) {
	u := url.URL{Path: c.cfg.S3Key}
	escapedPath := u.EscapedPath()
	if !strings.HasPrefix(escapedPath, "/") {
		escapedPath = "/" + escapedPath
	}

	host := fmt.Sprintf("%s.s3.%s.amazonaws.com", c.cfg.S3Bucket, c.cfg.AWSRegion)

	if strings.Contains(c.cfg.S3Bucket, ".") {
		host = fmt.Sprintf("s3.%s.amazonaws.com", c.cfg.AWSRegion)
		escapedPath = "/" + c.cfg.S3Bucket + escapedPath
	}

	return host, escapedPath
}

// ---------------------------
// Workflow Definition
// ---------------------------

func InitWorkflow(config *Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*Config], error) {
	return cre.Workflow[*Config]{
		cre.Handler(
			cron.Trigger(&cron.Config{Schedule: config.Schedule}),
			onCronTrigger,
		),
	}, nil
}

// fetchAWSCredentials uses the runtime to get secrets and returns them in the AWS SDK's struct.
func fetchAWSCredentials(runtime cre.Runtime) (aws.Credentials, error) {
	accessKeyPromise := runtime.GetSecret(&protos.SecretRequest{Id: "AWS_ACCESS_KEY_ID"})
	accessKey, err := accessKeyPromise.Await()
	if err != nil {
		return aws.Credentials{}, fmt.Errorf("failed to get AWS_ACCESS_KEY_ID: %w", err)
	}

	secretKeyPromise := runtime.GetSecret(&protos.SecretRequest{Id: "AWS_SECRET_ACCESS_KEY"})
	secretKey, err := secretKeyPromise.Await()
	if err != nil {
		return aws.Credentials{}, fmt.Errorf("failed to get AWS_SECRET_ACCESS_KEY: %w", err)
	}

	return aws.Credentials{
		AccessKeyID:     accessKey.Value,
		SecretAccessKey: secretKey.Value,
	}, nil
}

func onCronTrigger(config *Config, runtime cre.Runtime, trigger *cron.Payload) (*MyResult, error) {
	logger := runtime.Logger()
	logger.Info("Cron trigger fired. Fetching AWS credentials...")

	creds, err := fetchAWSCredentials(runtime)
	if err != nil {
		logger.Error("Workflow failed: could not fetch AWS credentials", "err", err)
		return nil, err
	}
	logger.Info("AWS credentials fetched. Running S3 Read/Increment/Write.")

	signer := v4.NewSigner()

	resultPromise := cre.RunInNodeMode(config, runtime,
		func(config *Config, nodeRuntime cre.NodeRuntime) (*MyResult, error) {
			return s3ReadIncrementWrite(config, nodeRuntime, creds, signer)
		},
		cre.ConsensusIdenticalAggregation[*MyResult](),
	)

	result, err := resultPromise.Await()
	if err != nil {
		logger.Error("Workflow failed", "err", err)
		return nil, err
	}

	logger.Info("Workflow finished successfully.", "old", result.OldValue, "new", result.NewValue)
	return result, nil
}

func s3ReadIncrementWrite(config *Config, nodeRuntime cre.NodeRuntime, creds aws.Credentials, signer *v4.Signer) (*MyResult, error) {
	logger := nodeRuntime.Logger()

	s3Client := NewS3Client(config, nodeRuntime, creds, signer)

	oldVal, err := s3Client.ReadValue()
	if err != nil {
		if isNotFound(err) {
			logger.Info("S3 object not found, starting count at 0.")
			oldVal = "0"
		} else {
			return nil, fmt.Errorf("s3ReadValue failed: %w", err)
		}
	}
	oldVal = strings.TrimSpace(oldVal)

	// Business logic - increment value
	newVal := incStringInt(oldVal)
	logger.Info("Value incremented", "old", oldVal, "new", newVal)

	if err := s3Client.WriteValue(newVal); err != nil {
		return nil, fmt.Errorf("s3WriteValue failed: %w", err)
	}

	return &MyResult{OldValue: oldVal, NewValue: newVal}, nil
}

func isNotFound(err error) bool {
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "404") || strings.Contains(errStr, "not found")
}

func incStringInt(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "1"
	}
	var n int64
	_, err := fmt.Sscan(s, &n)
	if err != nil {
		return "1"
	}
	n++
	return fmt.Sprintf("%d", n)
}

func sha256Hex(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

func main() {
	wasm.NewRunner(cre.ParseJSON[Config]).Run(InitWorkflow)
}
