package authorization

import "github.com/example/cre-confidential-authorization-template/gateway/internal/model"

func VerifyClaim(a Artifact, c model.ClaimEnvelope, now int64, trustedSigner string, verifier Verifier) error {
	if err := PreValidate(a, c, now, trustedSigner); err != nil {
		return err
	}
	if verifier == nil {
		return ErrVerifierMissing
	}
	return verifier.Verify(a)
}
