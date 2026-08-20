# AI Audit Firewall CRE Project

This CRE project screens proposed smart contract interactions before execution. It fetches verified contract artifacts, asks two independent models for structured risk analysis, and returns `ALLOW`, `DENY`, or `MANUAL_REVIEW`. TypeScript and Go implementations share the same configuration and secret mappings.

## Workflow

Each successful audit makes four fixed, in-enclave HTTP calls:

1. Fetch the token contract source and ABI from `https://api.etherscan.io/v2/api`.
2. Fetch the protocol contract source and ABI from the same Etherscan endpoint.
3. Send the proposal and verified token contract artifact to the configured OpenRouter primary model at `https://openrouter.ai/api/v1/chat/completions`.
4. Send the proposal, verified protocol contract artifact, and primary analysis to the configured OpenRouter secondary model at the same endpoint.

The workflow waits one second between the two Etherscan source requests to stay within Etherscan's free-tier rate limit.

During simulation, the non-sensitive `audit-firewall-primary-model-complete` and `audit-firewall-secondary-model-start` markers identify whether the primary or secondary model call failed.

The models check for `obfuscatedTax`, `privilegeEscalation`, `externalCallRisk`, and `logicBomb`. An unverified contract is denied without a model call. Model output is an advisory security signal, not authorization to execute a transaction.

## Configuration

Edit `config.staging.json` or `config.production.json` in the implementation you plan to run:

- `proposal` contains the candidate interaction, including both contract addresses.
- `proposal.from_address` is the account that would send the proposed transaction.
- `proposal.signer` is the account or system expected to authorize it. This can differ from `from_address` for multisigs, delegated wallets, or relayers.
- `from_address` and `signer` are input context sent to the models and included in the returned JSON. The workflow does not verify them, check a signature, or send the proposed transaction. The DON signs any verdict report.
- `etherscan_chain_id` selects the Etherscan chain.
- `primary_model` and `secondary_model` must be two different OpenRouter model IDs.
- `secrets_ids` must remain `etherscan_api_key` and `openrouter_api_key` to match `secrets.yaml`.
- `evms` is empty by default, so the workflow returns the verdict as JSON without requiring a deployed receiver. To write verdicts on-chain, deploy `AuditFirewallConsumer`, then configure `evms[0]` with its `consumer_address`, the matching `chain_selector_name`, and a `gas_limit`.

The low-cost paid defaults are primary model `google/gemini-2.5-flash-lite` and secondary model `openai/gpt-4.1-nano`. Direct tests with the workflow's real prompts returned valid structured output from both models within CRE's 10-second standard HTTP limit. Provider pricing, availability, and latency can change.

## Secrets

Copy the environment template and set both values:

```bash
cp .env.example .env
```

```dotenv
ETHERSCAN_API_KEY=your-etherscan-key
OPENROUTER_API_KEY=your-openrouter-key
```

Do not commit `.env`. For deployment, provide these same two values through the CRE secret workflow referenced by `secrets.yaml`.

## Simulate or Deploy

Run commands from this directory. Choose one implementation:

```bash
# TypeScript
cre workflow simulate ./ai-audit-firewall-ts --target=staging-settings

# Go
cre workflow simulate ./ai-audit-firewall-go --target=staging-settings
```

After reviewing the staging result and production config, deploy the chosen implementation:

```bash
# TypeScript
cre workflow deploy ./ai-audit-firewall-ts --target=production-settings

# Go
cre workflow deploy ./ai-audit-firewall-go --target=production-settings
```

No local service is required.

## Trust Boundaries

- Etherscan and OpenRouter credentials are fetched inside confidential execution and are not included in workflow output.
- Verified contract source and model reasoning remain protected from DON node operators, but Etherscan receives the configured chain ID and contract addresses.
- OpenRouter receives the proposal and verified token contract artifact for the primary call. For the secondary call, it receives the proposal, verified protocol contract artifact, and primary analysis. Its provider policy sets data collection to `deny`.
- The configured `proposal` is visible to the DON. Do not place secrets in it.
- With `evms` configured, only the encoded verdict, risk flags, and chain selector are written on-chain; credentials, contract source, and model reasoning are not.
