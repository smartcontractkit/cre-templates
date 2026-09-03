package authorization

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ErrMalformedSignature      = errors.New("malformed EIP-191 signature")
	ErrMalformedExpectedSigner = errors.New("malformed expected signer")
	ErrRecoveredSignerMismatch = errors.New("recovered signer mismatch")
)

type EIP191Verifier struct{ ExpectedSigner string }

func (v EIP191Verifier) Verify(a Artifact) error {
	if a.SigningScheme != SigningSchemeEIP191SHA256BindingV2 {
		return ErrUnsupportedSigningScheme
	}
	if !common.IsHexAddress(v.ExpectedSigner) || !common.IsHexAddress(a.Signer) {
		return ErrMalformedExpectedSigner
	}
	if !strings.EqualFold(v.ExpectedSigner, a.Signer) {
		return ErrRecoveredSignerMismatch
	}
	binding, err := hex.DecodeString(BindingDigest(a))
	if err != nil || len(binding) != 32 {
		return fmt.Errorf("decode binding digest: %w", ErrMalformedSignature)
	}
	sigText := strings.TrimPrefix(a.Signature, "0x")
	sig, err := hex.DecodeString(sigText)
	if err != nil || len(sig) != 65 {
		return ErrMalformedSignature
	}
	if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27
	}
	if sig[64] > 1 {
		return ErrMalformedSignature
	}
	pub, err := crypto.SigToPub(accounts.TextHash(binding), sig)
	if err != nil {
		return ErrMalformedSignature
	}
	recovered := crypto.PubkeyToAddress(*pub)
	if !strings.EqualFold(recovered.Hex(), v.ExpectedSigner) {
		return ErrRecoveredSignerMismatch
	}
	return nil
}
