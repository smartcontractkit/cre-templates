export type Claim = {
  version: string
  claimId: string
  asset: string
  claimedAmount: string
  authorizationState: string
  verifiedAmount: string
  authorizedAmount: string
  settledAmount: string
  fundingSource: string | null
}

export function evaluateClaim(claim: Claim): 'AUTHORIZATION_REQUIRED' | 'REJECTED' {
  if (
    claim.authorizationState !== 'UNAUTHORIZED' ||
    claim.verifiedAmount !== '0' ||
    claim.authorizedAmount !== '0' ||
    claim.settledAmount !== '0' ||
    claim.fundingSource !== null
  ) return 'REJECTED'

  return 'AUTHORIZATION_REQUIRED'
}
