package httpapi

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/example/cre-confidential-authorization-template/gateway/internal/model"
	"github.com/example/cre-confidential-authorization-template/gateway/internal/store"
)

type API struct {
	store         *store.Store
	token         string
	trustedSigner string
}

func New(s *store.Store, token string) *API { return NewWithTrustedSigner(s, token, "") }
func NewWithTrustedSigner(s *store.Store, token, signer string) *API {
	return &API{store: s, token: token, trustedSigner: signer}
}

func (a *API) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", a.health)
	mux.HandleFunc("/v1/claims", a.auth(a.claims))
	mux.HandleFunc("/v1/claims/", a.auth(a.claimByID))
	mux.HandleFunc("/v1/authorization/verify", a.auth(a.authorizationVerify))
	return mux
}

func (a *API) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if len(got) != len(a.token) || subtle.ConstantTimeCompare([]byte(got), []byte(a.token)) != 1 {
			writeJSON(w, 401, map[string]any{"ok": false})
			return
		}
		next(w, r)
	}
}
func (a *API) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"ok": true, "service": "cre-confidential-authorization-template", "time": time.Now().Unix()})
}
func (a *API) claims(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	var c model.ClaimEnvelope
	dec := json.NewDecoder(io.LimitReader(r.Body, 64<<10))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&c); err != nil {
		writeError(w, 400, err)
		return
	}
	if err := validateClaim(c); err != nil {
		writeError(w, 422, err)
		return
	}
	if err := a.store.AddClaim(c); err != nil {
		writeError(w, 409, err)
		return
	}
	writeJSON(w, 201, map[string]any{"ok": true, "claim": c})
}
func validateClaim(c model.ClaimEnvelope) error {
	if c.Version != "x402.example.claim.v1" || c.ClaimID == "" || c.Asset == "" || c.ClaimedAmount == "" || len(c.CanonicalClaimSHA256) != 64 {
		return errors.New("invalid claim")
	}
	if c.AuthorizationState != "UNAUTHORIZED" || c.VerifiedAmount != "0" || c.AuthorizedAmount != "0" || c.SettledAmount != "0" || c.FundingSource != nil {
		return errors.New("claim must enter pristine and unauthorized")
	}
	return nil
}
func (a *API) claimByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/claims/")
	if strings.HasSuffix(path, "/status") {
		id := strings.TrimSuffix(path, "/status")
		a.claimStatus(w, r, id)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		return
	}
	rec, ok := a.store.GetClaim(path)
	if !ok {
		w.WriteHeader(404)
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "record": rec})
}
func (a *API) claimStatus(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	var req struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, 16<<10)).Decode(&req); err != nil {
		writeError(w, 400, err)
		return
	}
	if err := a.store.SetClaimStatus(id, req.Status, req.Reason); err != nil {
		writeError(w, 422, err)
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true})
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]any{"ok": false, "error": err.Error()})
}
