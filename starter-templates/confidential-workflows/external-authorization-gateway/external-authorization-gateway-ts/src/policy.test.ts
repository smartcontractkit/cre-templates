import { expect, test } from 'bun:test'
import { evaluateClaim } from './policy'

test('pristine claim requires independent authorization', () => {
  expect(evaluateClaim({
    version: 'x402.example.claim.v1', claimId: 'example-1', asset: 'TEST_TOKEN', claimedAmount: '1000',
    authorizationState: 'UNAUTHORIZED', verifiedAmount: '0', authorizedAmount: '0', settledAmount: '0', fundingSource: null,
  })).toBe('AUTHORIZATION_REQUIRED')
})

test('claim cannot self-authorize', () => {
  expect(evaluateClaim({
    version: 'x402.example.claim.v1', claimId: 'example-1', asset: 'TEST_TOKEN', claimedAmount: '1000',
    authorizationState: 'AUTHORIZED', verifiedAmount: '1000', authorizedAmount: '1000', settledAmount: '0', fundingSource: 'self',
  })).toBe('REJECTED')
})
