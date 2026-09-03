package authorization

import "errors"

var (
	ErrVerifierMissing          = errors.New("authorization verifier is required")
	ErrUnsupportedSigningScheme = errors.New("unsupported authorization signing scheme")
)

type Verifier interface {
	Verify(Artifact) error
}
