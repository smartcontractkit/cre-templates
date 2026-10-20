# Confidential Random

Generate a fresh random `uint32` on every cron tick inside an AWS Nitro enclave.
This starter template uses the **Confidential Workflows** TEE handler and the
TEE-local randomness seed, so values differ across runs and across concurrent
executions.

Both language variants are behaviourally equivalent; pick the one you build in:

| Variant | Path | Randomness API |
|---------|------|----------------|
| Go | [`confidential-random-go`](./confidential-random-go) | `rt.Rand()` on a `cre.TeeRuntime` |
| TypeScript | [`confidential-random-ts`](./confidential-random-ts) | `Math.random()` inside `handlerInTee` |

**⚠️ PRIVATE BETA**

[Confidential Workflows](https://docs.chain.link/cre/concepts/confidential-workflows) is in **private beta** and requires enrollment through your Chainlink account team — see [Requesting Confidential Workflows Access](https://docs.chain.link/cre/account/confidential-workflows-access). TEE randomness currently runs in **AWS Nitro in `us-west-2`** (staging).

## What it does

Each variant registers a cron handler with the TEE variant of the handler API
(`cre.HandlerInTee` in Go, `handlerInTee` in TypeScript). On every tick the
handler draws a random `uint32` inside the enclave and returns it as a string.

- **Go:** `rand, err := rt.Rand()` then `rand.Uint32()`.
- **TypeScript:** `Math.floor(Math.random() * 0x1_0000_0000) >>> 0`.

## Why TEE randomness

- **Unique per execution** — the enclave derives a new seed for each run.
- **Non-repeating** — consecutive and concurrent executions get different
  values; workflows running at the same time do not receive the same number.
- **Different from DON randomness** — it does not use the deterministic
  DON-seeded PRNG, so values are not constrained by the repeatability
  requirements of DON consensus.

TEE randomness suits raffles, sampling, matchmaking, and similar
non-cryptographic draws. It is not a replacement for a verifiable random
function (VRF) when you need an on-chain, publicly verifiable proof.

## Quick start

Go:

```bash
cd confidential-random-go
GOOS=wasip1 GOARCH=wasm go build ./my-workflow
cre workflow simulate my-workflow --target staging-settings --non-interactive --trigger-index 0
```

TypeScript:

```bash
cd confidential-random-ts/my-workflow
bun install
cd ..
cre workflow simulate my-workflow --target staging-settings --non-interactive --trigger-index 0
```

Deploy and activate each variant with:

```bash
cre workflow deploy my-workflow --target staging-settings
cre workflow activate my-workflow --target staging-settings
```

See each variant's README for configuration, customization, SDK pinning, and
the staging-specific TEE requirements.

## License

MIT
