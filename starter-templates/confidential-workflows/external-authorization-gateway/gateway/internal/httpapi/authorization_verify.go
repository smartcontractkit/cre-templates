package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/example/cre-confidential-authorization-template/gateway/internal/authorization"
)

type authorizationVerifyRequest struct {
	ClaimID  string                 `json:"claimId"`
	Artifact authorization.Artifact `json:"artifact"`
}

func (a *API) authorizationVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	if a.trustedSigner == "" {
		writeError(w, 503, authorization.ErrTrustedSignerMissing)
		return
	}
	dec := json.NewDecoder(io.LimitReader(r.Body, 128<<10))
	dec.DisallowUnknownFields()
	var req authorizationVerifyRequest
	if err := dec.Decode(&req); err != nil {
		writeError(w, 400, err)
		return
	}
	if req.ClaimID == "" || req.Artifact.ClaimID != req.ClaimID {
		writeError(w, 422, errors.New("claimId mismatch"))
		return
	}
	rec, ok := a.store.GetClaim(req.ClaimID)
	if !ok {
		w.WriteHeader(404)
		return
	}
	if rec.Status != "AUTHORIZATION_REQUIRED" {
		writeError(w, 409, errors.New("claim is not awaiting authorization"))
		return
	}
	v := authorization.EIP191Verifier{ExpectedSigner: a.trustedSigner}
	if err := authorization.VerifyClaim(req.Artifact, rec.Claim, time.Now().Unix(), a.trustedSigner, v); err != nil {
		writeError(w, 422, err)
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "claimId": req.ClaimID, "cryptographicallyValid": true, "claimMutated": false, "authorizationApplied": false, "settlementCreated": false, "signingScheme": req.Artifact.SigningScheme})
}
