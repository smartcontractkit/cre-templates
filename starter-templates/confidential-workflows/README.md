# Confidential Workflows

Starter templates for building confidential workflows using Chainlink CRE.

Each workflow is available in both TypeScript and Go. The two implementations are
behaviourally equivalent; pick whichever language you build in.

Both live in the same CRE project, as sibling workflow directories sharing one
`project.yaml`, `secrets.yaml` and `.env.example`:

```
ai-audit-firewall/            <- project root
  project.yaml  secrets.yaml  .env.example  contracts/
  ai-audit-firewall-ts/       <- TypeScript workflow
  ai-audit-firewall-go/       <- Go workflow (go.mod rooted here)
```

## Available Workflows

| Workflow | TypeScript | Go |
|----------|------------|-----|
| AI Audit Firewall | [ai-audit-firewall-ts](./ai-audit-firewall/ai-audit-firewall-ts) | [ai-audit-firewall-go](./ai-audit-firewall/ai-audit-firewall-go) |
| Automated Liquidation Protection | [automated-liquidation-protection-ts](./automated-liquidation-protection/automated-liquidation-protection-ts) | [automated-liquidation-protection-go](./automated-liquidation-protection/automated-liquidation-protection-go) |
| Automated Portfolio Rebalancing | [automated-portfolio-rebalancing-ts](./automated-portfolio-rebalancing/automated-portfolio-rebalancing-ts) | [automated-portfolio-rebalancing-go](./automated-portfolio-rebalancing/automated-portfolio-rebalancing-go) |

## Confidential execution

All six templates register their handler with the TEE variant of the handler API
(`handlerInTee` in TypeScript, `cre.HandlerInTee` in Go), so the handler receives
a TEE runtime rather than a regular one. Secrets are released by the Vault DON
directly into the attested enclave, and HTTP calls made from the handler keep
their URLs, headers and response bodies confidential from node operators.

The workflow binary itself is *not* confidential — it is provided to the enclave
by the Workflow DON. What the enclave protects is the data the logic computes
over. See each template's README for its specific confidentiality boundary.

## Note on the Go templates

The Go SDK versions in each `go.mod` are pseudo-versions pinned to an unreleased
commit, because Confidential Workflows is not yet in a tagged release. Repin them
once it ships. To build against a local SDK checkout instead, use a `go.work`
file (gitignored) rather than editing `go.mod` — see the "Testing against an
unreleased SDK" section in any Go template's README.
