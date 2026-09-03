package authorization

type Artifact struct {
	Version              string `json:"version"`
	AuthorizationID      string `json:"authorizationId"`
	ClaimID              string `json:"claimId"`
	CanonicalClaimSHA256 string `json:"canonicalClaimSha256"`
	Asset                string `json:"asset"`
	Amount               string `json:"amount"`
	Recipient            string `json:"recipient"`
	FundingSource        string `json:"fundingSource"`
	Scope                string `json:"scope"`
	Signer               string `json:"signer"`
	SigningScheme        string `json:"signingScheme"`
	Signature            string `json:"signature"`
	IssuedAt             int64  `json:"issuedAt"`
	ExpiresAt            int64  `json:"expiresAt"`
}

const (
	VersionV1                          = "x402.example.authorization.v1"
	ScopeSettlementOnly                = "SETTLEMENT_ONLY"
	SigningSchemeEIP191SHA256BindingV2 = "EIP191_SHA256_BINDING_V2"
)
