package authorization

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	"github.com/example/cre-confidential-authorization-template/gateway/internal/model"
)

var (
	ErrTrustedSignerMissing   = errors.New("trusted counterparty signer not configured")
	ErrVersion                = errors.New("invalid authorization version")
	ErrClaimBinding           = errors.New("authorization claim binding mismatch")
	ErrAssetBinding           = errors.New("authorization asset mismatch")
	ErrAmountBinding          = errors.New("authorization amount mismatch")
	ErrRecipientBinding       = errors.New("authorization recipient mismatch")
	ErrFundingSource          = errors.New("authorization funding source missing")
	ErrScope                  = errors.New("authorization scope invalid")
	ErrSigner                 = errors.New("authorization signer is not trusted")
	ErrSignatureMissing       = errors.New("authorization signature missing")
	ErrSigningSchemeMissing   = errors.New("authorization signing scheme missing")
	ErrTimeWindow             = errors.New("authorization time window invalid")
	ErrClaimAlreadyAuthorized = errors.New("claim authorization state is not pristine")
)

func PreValidate(a Artifact, c model.ClaimEnvelope, now int64, trustedSigner string) error {
	if trustedSigner == "" {
		return ErrTrustedSignerMissing
	}
	if a.Version != VersionV1 {
		return ErrVersion
	}
	if c.AuthorizationState != "UNAUTHORIZED" || c.VerifiedAmount != "0" || c.AuthorizedAmount != "0" || c.SettledAmount != "0" || c.FundingSource != nil {
		return ErrClaimAlreadyAuthorized
	}
	if a.AuthorizationID == "" || a.ClaimID != c.ClaimID || a.CanonicalClaimSHA256 != c.CanonicalClaimSHA256 {
		return ErrClaimBinding
	}
	if a.Asset != c.Asset {
		return ErrAssetBinding
	}
	if a.Amount != c.ClaimedAmount {
		return ErrAmountBinding
	}
	if !strings.EqualFold(a.Recipient, c.EVMSettlementAddress) {
		return ErrRecipientBinding
	}
	if a.FundingSource == "" {
		return ErrFundingSource
	}
	if a.Scope != ScopeSettlementOnly {
		return ErrScope
	}
	if !strings.EqualFold(a.Signer, trustedSigner) {
		return ErrSigner
	}
	if a.SigningScheme == "" {
		return ErrSigningSchemeMissing
	}
	if a.Signature == "" {
		return ErrSignatureMissing
	}
	if a.IssuedAt <= 0 || a.ExpiresAt <= a.IssuedAt || now < a.IssuedAt || now > a.ExpiresAt {
		return ErrTimeWindow
	}
	return nil
}

func BindingDigest(a Artifact) string {
	payload := strings.Join([]string{
		"x402.authorization.binding.v2", a.Version, a.AuthorizationID, a.ClaimID, a.CanonicalClaimSHA256,
		a.Asset, a.Amount, strings.ToLower(a.Recipient), a.FundingSource, a.Scope,
		strings.ToLower(a.Signer), a.SigningScheme, strconv.FormatInt(a.IssuedAt, 10), strconv.FormatInt(a.ExpiresAt, 10),
	}, "\x1f")
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}
