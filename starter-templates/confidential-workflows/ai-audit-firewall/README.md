# AI Audit Firewall CRE Project

This standalone CRE project implements a confidential pre-execution security firewall for smart contract interactions.

## Description

The workflow screens proposed transactions before they are allowed to proceed. It fetches and validates contract intelligence, runs confidential reasoning to classify risk, and then enforces a firewall decision path. Scanner and model credentials remain protected inside confidential execution throughout the process.

## Target Customer

- Professional retail traders
- Developer shops
- Founders building trading products

## Structure

- `project.yaml`: project-level target settings
- `secrets.yaml`: secret ID mappings used by the workflow
- `mock-server.js`: local deterministic API server
- `ai-audit-firewall/`: workflow implementation

## Private Inputs

The following inputs are handled as confidential:

- Chain scanner API credentials used for contract metadata retrieval and verification checks.
- LLM reasoning API credentials used for independent audit analysis.

## Workflow Notes

1. Monitor and ingest the proposed interaction.
	The workflow receives candidate transaction context, including token and protocol contract addresses.
2. Fetch and validate contract data confidentially.
	It retrieves source and ABI artifacts through the scanner and verifies scanner credential permissions before trusting fetched data.
3. Run smart contract audit analysis.
	The workflow submits context to multiple reasoning models and classifies behavior into structured risk signals:
	- `obfuscatedTax`
	- `privilegeEscalation`
	- `externalCallRisk`
	- `logicBomb`
4. Enforce firewall action and record outcomes.
	Based on aggregate risk, the workflow allows execution, blocks malicious interactions, or routes the attempt for manual review while preserving audit and action logs.

Note: Any reasoning stage can be replaced with deterministic rule-based logic if a purely policy-engine implementation is preferred.

## Required Environment Variables

Copy `.env.example` to `.env` and provide values for:

- `CRE_ETH_PRIVATE_KEY` (optional for local simulate, required for real chain writes)
- `MOCK_PORT`
- `MOCK_SCANNER_API_KEY`
- `MOCK_PRIMARY_LLM_API_KEY`
- `MOCK_SECONDARY_LLM_API_KEY`

The local mock server for this project only exposes routes under `/audit-firewall/*`.

## Quick Start

1. Install dependencies

```bash
bun install
```

2. Create environment file

```bash
cp .env.example .env
```

3. Start mock server

```bash
bun run mock:server
```

4. In another terminal, run checks

```bash
bun run typecheck
bun run test
```

5. Simulate workflow

```bash
cre workflow simulate ./ai-audit-firewall --target=staging-settings
```
