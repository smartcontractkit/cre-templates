# offchain-report-write

CRE TypeScript workflow that writes a simple `uint256` (the execution timestamp) to a
`Storage` consumer contract through the chain writer, and logs the DON-signed report as
copy-pasteable JSON for offchain verification.

See [../README.md](../README.md) for the full guide.

## Quick reference

```bash
bun install                 # from this directory (also run bun install in ../contracts)

# from the project root (the template root, where project.yaml lives):
cre workflow simulate offchain-report-write --target staging-settings
```

Config (`config.staging.json` / `config.production.json`):

| Key | Description |
|-----|-------------|
| `schedule` | Cron schedule (6-field, with seconds) |
| `backendUrl` | HTTP endpoint the signed report is POSTed to (via `sendReport`) |
| `evms[0].storageAddress` | Deployed `Storage` contract address |
| `evms[0].chainSelectorName` | CRE chain selector name (default `ethereum-testnet-sepolia`) |
| `evms[0].gasLimit` | Gas limit for the write transaction |
