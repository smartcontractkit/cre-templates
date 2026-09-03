package store

import (
	"path/filepath"
	"testing"

	"github.com/example/cre-confidential-authorization-template/gateway/internal/model"
)

func sampleClaim() model.ClaimEnvelope {
	return model.ClaimEnvelope{
		Version: "x402.example.claim.v1", ClaimID: "example-1", Asset: "TEST_TOKEN", ClaimedAmount: "1000",
		CanonicalClaimSHA256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		EVMSettlementAddress: "0x1111111111111111111111111111111111111111", AuthorizationState: "UNAUTHORIZED",
		VerifiedAmount: "0", AuthorizedAmount: "0", SettledAmount: "0",
	}
}

func TestPersistentClaimSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s, err := NewPersistent(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.AddClaim(sampleClaim()); err != nil {
		t.Fatal(err)
	}
	if err := s.SetClaimStatus("example-1", "AUTHORIZATION_REQUIRED", "COUNTERPARTY_AUTHORIZATION_ABSENT"); err != nil {
		t.Fatal(err)
	}
	reopened, err := NewPersistent(path)
	if err != nil {
		t.Fatal(err)
	}
	r, ok := reopened.GetClaim("example-1")
	if !ok {
		t.Fatal("missing claim")
	}
	if r.Status != "AUTHORIZATION_REQUIRED" || r.Claim.AuthorizationState != "UNAUTHORIZED" {
		t.Fatalf("unexpected record: %+v", r)
	}
}

func TestTerminalStatusCannotBecomeAuthorized(t *testing.T) {
	s := New()
	_ = s.AddClaim(sampleClaim())
	_ = s.SetClaimStatus("example-1", "AUTHORIZATION_REQUIRED", "required")
	if err := s.SetClaimStatus("example-1", "AUTHORIZED", "bad"); err == nil {
		t.Fatal("expected rejection")
	}
}
