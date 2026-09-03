package authorization

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/example/cre-confidential-authorization-template/gateway/internal/model"
)

func testClaim() model.ClaimEnvelope {
	return model.ClaimEnvelope{
		Version: "x402.example.claim.v1", ClaimID: "example-1", Asset: "TEST_TOKEN", ClaimedAmount: "1000",
		CanonicalClaimSHA256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		EVMSettlementAddress: "0x1111111111111111111111111111111111111111", AuthorizationState: "UNAUTHORIZED",
		VerifiedAmount: "0", AuthorizedAmount: "0", SettledAmount: "0",
	}
}

func signedArtifact(t *testing.T) (Artifact, string) {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	signer := crypto.PubkeyToAddress(key.PublicKey).Hex()
	now := time.Now().Unix()
	a := Artifact{Version: VersionV1, AuthorizationID: "auth-test-1", ClaimID: "example-1", CanonicalClaimSHA256: testClaim().CanonicalClaimSHA256,
		Asset: "TEST_TOKEN", Amount: "1000", Recipient: testClaim().EVMSettlementAddress, FundingSource: "example-counterparty-source",
		Scope: ScopeSettlementOnly, Signer: signer, SigningScheme: SigningSchemeEIP191SHA256BindingV2, IssuedAt: now - 5, ExpiresAt: now + 300}
	binding, _ := hex.DecodeString(BindingDigest(a))
	sig, err := crypto.Sign(accounts.TextHash(binding), key)
	if err != nil {
		t.Fatal(err)
	}
	a.Signature = "0x" + hex.EncodeToString(sig)
	return a, signer
}

func TestEIP191VerifierAcceptsMatchingSigner(t *testing.T) {
	a, signer := signedArtifact(t)
	if err := VerifyClaim(a, testClaim(), time.Now().Unix(), signer, EIP191Verifier{ExpectedSigner: signer}); err != nil {
		t.Fatal(err)
	}
}

func TestEIP191VerifierRejectsTamperedExpiry(t *testing.T) {
	a, signer := signedArtifact(t)
	a.ExpiresAt += 3600
	err := EIP191Verifier{ExpectedSigner: signer}.Verify(a)
	if err == nil {
		t.Fatal("expected tampered expiry to invalidate signature")
	}
}

func TestVerifyClaimDoesNotMutateClaim(t *testing.T) {
	a, signer := signedArtifact(t)
	c := testClaim()
	if err := VerifyClaim(a, c, time.Now().Unix(), signer, EIP191Verifier{ExpectedSigner: signer}); err != nil {
		t.Fatal(err)
	}
	if c.AuthorizationState != "UNAUTHORIZED" || c.VerifiedAmount != "0" || c.AuthorizedAmount != "0" || c.SettledAmount != "0" || c.FundingSource != nil {
		t.Fatal("claim mutated")
	}
}
