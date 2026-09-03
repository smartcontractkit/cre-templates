package httpapi

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/example/cre-confidential-authorization-template/gateway/internal/authorization"
	"github.com/example/cre-confidential-authorization-template/gateway/internal/model"
	"github.com/example/cre-confidential-authorization-template/gateway/internal/store"
)

const testToken = "test-bearer-token"

func httpClaim() model.ClaimEnvelope {
	return model.ClaimEnvelope{Version: "x402.example.claim.v1", ClaimID: "example-1", Asset: "TEST_TOKEN", ClaimedAmount: "1000", CanonicalClaimSHA256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", EVMSettlementAddress: "0x1111111111111111111111111111111111111111", AuthorizationState: "UNAUTHORIZED", VerifiedAmount: "0", AuthorizedAmount: "0", SettledAmount: "0"}
}
func readyStore(t *testing.T) *store.Store {
	t.Helper()
	s := store.New()
	if err := s.AddClaim(httpClaim()); err != nil {
		t.Fatal(err)
	}
	if err := s.SetClaimStatus("example-1", "AUTHORIZATION_REQUIRED", "required"); err != nil {
		t.Fatal(err)
	}
	return s
}
func signedHTTPArtifact(t *testing.T) (authorization.Artifact, string) {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	signer := crypto.PubkeyToAddress(key.PublicKey).Hex()
	now := time.Now().Unix()
	a := authorization.Artifact{Version: authorization.VersionV1, AuthorizationID: "auth-http-1", ClaimID: "example-1", CanonicalClaimSHA256: httpClaim().CanonicalClaimSHA256, Asset: "TEST_TOKEN", Amount: "1000", Recipient: httpClaim().EVMSettlementAddress, FundingSource: "example-counterparty-source", Scope: authorization.ScopeSettlementOnly, Signer: signer, SigningScheme: authorization.SigningSchemeEIP191SHA256BindingV2, IssuedAt: now - 5, ExpiresAt: now + 300}
	b, _ := hex.DecodeString(authorization.BindingDigest(a))
	sig, err := crypto.Sign(accounts.TextHash(b), key)
	if err != nil {
		t.Fatal(err)
	}
	a.Signature = "0x" + hex.EncodeToString(sig)
	return a, signer
}
func postVerify(t *testing.T, api *API, body any) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/authorization/verify", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+testToken)
	rr := httptest.NewRecorder()
	api.Handler().ServeHTTP(rr, req)
	return rr
}
func TestAuthorizationVerifyFailsClosedWithoutTrustedSigner(t *testing.T) {
	rr := postVerify(t, New(store.New(), testToken), map[string]any{})
	if rr.Code != 503 {
		t.Fatalf("expected 503 got %d", rr.Code)
	}
}
func TestAuthorizationVerifyValidSignatureIsReadOnly(t *testing.T) {
	s := readyStore(t)
	a, signer := signedHTTPArtifact(t)
	before, _ := s.GetClaim("example-1")
	rr := postVerify(t, NewWithTrustedSigner(s, testToken, signer), map[string]any{"claimId": "example-1", "artifact": a})
	if rr.Code != 200 {
		t.Fatalf("expected 200 got %d body=%s", rr.Code, rr.Body.String())
	}
	after, _ := s.GetClaim("example-1")
	if !reflect.DeepEqual(before, after) {
		t.Fatal("claim mutated")
	}
	if after.Claim.AuthorizationState != "UNAUTHORIZED" || after.Claim.AuthorizedAmount != "0" || after.Claim.SettledAmount != "0" {
		t.Fatal("financial state mutated")
	}
}
