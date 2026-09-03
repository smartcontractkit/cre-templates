# CRE Confidential Authorization Template

A sanitized starter template for a fail-closed confidential authorization flow.

This template demonstrates:

- authenticated claim ingress through a local Go gateway;
- persistent claim state;
- an explicit `AUTHORIZATION_REQUIRED` boundary;
- deterministic authorization binding with `IssuedAt` and `ExpiresAt` included in the signed digest;
- EIP-191 / secp256k1 signer recovery;
- a read-only `/v1/authorization/verify` endpoint;
- no claim self-authorization and no settlement side effects;
- a Chainlink CRE TypeScript workflow skeleton pinned to the same dependency versions as the current official confidential TypeScript starter.

## Safety boundary

A cryptographically valid authorization is evidence only in this template. Verification never mutates the claim and never creates a transfer, settlement, mint, or token movement.

The gateway fails closed when `X402_AUTH_TRUSTED_SIGNER` is unset.

## Local verification

Prerequisites:

- Go 1.23+
- Bun (optional for CRE workflow typecheck/tests)
- Chainlink CRE CLI (optional for `cre workflow simulate`)

Run:

```bash
./verify-local.sh
```

The script isolates Go caches under `/tmp`, runs Go formatting checks, unit tests, race tests, and builds the gateway. If Bun is installed, it also installs and typechecks/tests the CRE workflow package.

## Run the gateway

```bash
export X402_GATEWAY_TOKEN='local-development-token'
export X402_STORE_PATH="$PWD/.runtime/gateway-state.json"
# Leave X402_AUTH_TRUSTED_SIGNER unset to exercise fail-closed mode.
cd gateway

go run ./cmd/gateway
```

Health:

```bash
curl -s http://127.0.0.1:8402/healthz | jq .
```

## Example claim

The included claim is intentionally synthetic:

- asset: `TEST_TOKEN`
- amount: `1000`
- recipient: `0x1111111111111111111111111111111111111111`

It does not identify or assert a real counterparty, obligation, payment, or settlement.

## CRE workflow

The `external-authorization-gateway-ts/` directory follows the current confidential TypeScript starter shape with `main.ts`, `workflow.ts`, `workflow.yaml`, staging/production config, `package.json`, and `tsconfig.json`.

The workflow is deliberately limited to reading a claim from the gateway and reporting `AUTHORIZATION_REQUIRED`. It contains no EVM write capability.
