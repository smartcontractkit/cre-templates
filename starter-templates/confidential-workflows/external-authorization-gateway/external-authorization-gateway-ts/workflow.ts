import {
  cre,
  hexToBase64,
  ok,
  text,
  type TeeRuntime,
} from '@chainlink/cre-sdk'
import { encodeAbiParameters, parseAbiParameters } from 'viem'
import { z } from 'zod'
import { evaluateClaim, type Claim } from './src/policy'

export const configSchema = z.object({
  schedule: z.string(),
  url: z.string().url(),
  secretId: z.string().min(1),
})
type Config = z.infer<typeof configSchema>

const GatewayClaimSchema = z.object({
  ok: z.literal(true),
  record: z.object({
    claim: z.object({
      version: z.string(),
      claimId: z.string(),
      asset: z.string(),
      claimedAmount: z.string(),
      authorizationState: z.string(),
      verifiedAmount: z.string(),
      authorizedAmount: z.string(),
      settledAmount: z.string(),
      fundingSource: z.string().nullable(),
    }),
    status: z.string(),
  }),
})

export const onCronTrigger = (runtime: TeeRuntime<Config>): string => {
  const token = runtime.getSecret({ id: runtime.config.secretId }).result().value

  const response = new cre.capabilities.HTTPClient()
    .sendRequest(runtime, {
      url: runtime.config.url,
      method: 'GET',
      multiHeaders: {
        Authorization: { values: [`Bearer ${token}`] },
        Accept: { values: ['application/json'] },
      },
    })
    .result()

  if (!ok(response)) {
    throw new Error(`gateway request failed with status ${response.statusCode}`)
  }

  const parsed = GatewayClaimSchema.parse(JSON.parse(text(response)))
  const decision = evaluateClaim(parsed.record.claim as Claim)

  // Cross only the non-sensitive decision back to the Workflow DON.
  const donRuntime = runtime.usingTheDons()
  const encodedPayload = encodeAbiParameters(
    parseAbiParameters('string claimId, string decision'),
    [parsed.record.claim.claimId, decision],
  )

  donRuntime
    .report({
      encodedPayload: hexToBase64(encodedPayload),
      encoderName: 'evm',
      signingAlgo: 'ecdsa',
      hashingAlgo: 'keccak256',
    })
    .result()

  // This template intentionally performs no EVM write and no settlement.
  return JSON.stringify({
    claimId: parsed.record.claim.claimId,
    decision,
  })
}

export function initWorkflow(config: Config) {
  const cronTrigger = new cre.capabilities.CronCapability()

  return [
    cre.handlerInTee(
      cronTrigger.trigger({ schedule: config.schedule }),
      onCronTrigger,
      [{ tee: 'nitro', regions: ['us-west-2'] }],
    ),
  ]
}
