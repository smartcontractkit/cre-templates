package model

type ClaimEnvelope struct {
	Version              string  `json:"version"`
	ClaimID              string  `json:"claimId"`
	Asset                string  `json:"asset"`
	ClaimedAmount        string  `json:"claimedAmount"`
	CanonicalClaimSHA256 string  `json:"canonicalClaimSha256"`
	EVMSettlementAddress string  `json:"evmSettlementAddress"`
	AuthorizationState   string  `json:"authorizationState"`
	VerifiedAmount       string  `json:"verifiedAmount"`
	AuthorizedAmount     string  `json:"authorizedAmount"`
	SettledAmount        string  `json:"settledAmount"`
	FundingSource        *string `json:"fundingSource"`
}

type ClaimRecord struct {
	Claim      ClaimEnvelope `json:"claim"`
	Status     string        `json:"status"`
	Reason     string        `json:"reason,omitempty"`
	ReceivedAt int64         `json:"receivedAt"`
	UpdatedAt  int64         `json:"updatedAt,omitempty"`
}
