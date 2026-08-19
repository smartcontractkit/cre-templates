# Deploying `kv_store_receiver` (new to Solana? start here)

This is the **only** thing you need to deploy yourself. The **keystone forwarder** —
the program that verifies the CRE DON's signed report and CPIs into your receiver —
is already deployed and operated by Chainlink. You never deploy or configure it.

By the end of this guide you'll have:

- Your own `kv_store_receiver` program live on Solana **devnet**.
- A `kv_store` account it can write into.
- The four addresses (`receiverProgramId`, `forwarderState`, `forwarderAuthority`,
  `kvStoreAccount`) to paste into [`../../my-workflow/config.staging.json`](../../my-workflow/config.staging.json).

Do this in a fresh terminal — no CRE-specific tooling required for this part, it's
pure Solana/Anchor.

## 1. Install the toolchain

You need **Rust**, the **Solana CLI**, and **Anchor**.

```bash
# Rust (skip if you already have it)
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Solana CLI (Agave)
sh -c "$(curl -sSfL https://release.anza.xyz/stable/install)"
# Add the printed path to your shell profile, e.g.:
#   export PATH="$HOME/.local/share/solana/install/active_release/bin:$PATH"

# Anchor version manager, then Anchor 0.31.0
cargo install --git https://github.com/coral-xyz/anchor avm --locked --force
avm install 0.31.0
avm use 0.31.0
```

Verify:

```bash
solana --version   # solana-cli 2.x
anchor --version   # anchor-cli 0.31.0
```

## 2. Create a devnet wallet and fund it

```bash
solana-keygen new --outfile ~/.config/solana/id.json
solana config set --url https://api.devnet.solana.com --keypair ~/.config/solana/id.json
solana airdrop 2
solana balance   # should show ~2 SOL
```

If the public faucet is rate-limited, use https://faucet.solana.com with your
`solana address` output.

## 3. Build the program

From this directory (`contracts/solana`):

```bash
anchor build
```

This compiles `programs/kv_store_receiver` and generates:

- `target/deploy/kv_store_receiver.so` — the program binary.
- `target/deploy/kv_store_receiver-keypair.json` — your program's on-chain address (keypair).
- `target/idl/kv_store_receiver.json` — the IDL, used by the deploy script and by
  `cre generate-bindings` if you regenerate the TypeScript client.

## 4. Align the program id and rebuild

The program currently declares a placeholder id. Sync it to the keypair Anchor just
generated for you, then rebuild so the binary embeds the real id:

```bash
anchor keys sync
anchor build
```

## 5. Deploy to devnet

```bash
anchor deploy --provider.cluster devnet
```

Note the program id it prints (also visible via `solana address -k target/deploy/kv_store_receiver-keypair.json`)
— that's your `receiverProgramId`.

## 6. Initialize your receiver

Install the small helper script's dependencies once, then run it:

```bash
bun install
bun run scripts/initialize.ts \
  --forwarder-program Ey9KNM4ny1S2Sx2g9NiUJUsTy7e6ta1bcP28cT1xmKj5 \
  --forwarder-state 726bmjnxLeYkvSWUBRby5eUXds96xR58MFfcZoBzuuW3 \
  --keypair ~/.config/solana/id.json \
  --url https://api.devnet.solana.com
```

`--forwarder-program` and `--forwarder-state` above are Chainlink's shared
Solana-devnet keystone forwarder for the **staging** CRE environment — the same
values already filled into [`my-workflow/config.staging.json`](../../my-workflow/config.staging.json).
Leave them as-is unless Chainlink gives you different addresses for your environment.

The script creates your `kv_store` account, calls `initialize(forwarder_program)`,
derives `forwarderAuthority` (the PDA the forwarder signs with), and prints all four
values plus a devnet explorer link for the `initialize` transaction. Example output:

```
Initialized kv_store_receiver
--------------------------------
receiverProgramId : <your program id>
forwarderState     (input, echo) : 726bmjnxLeYkvSWUBRby5eUXds96xR58MFfcZoBzuuW3
forwarderAuthority (derived)      : <derived PDA>
kvStoreAccount     (receiverAccounts[0]) : <new account>

initialize() tx: https://explorer.solana.com/tx/<signature>?cluster=devnet
```

## 7. Wire it into the workflow config

Copy the four values into [`my-workflow/config.staging.json`](../../my-workflow/config.staging.json):

```json
{
  "solana": {
    "chainSelectorName": "solana-devnet",
    "receiverProgramId": "<receiverProgramId>",
    "forwarderState": "726bmjnxLeYkvSWUBRby5eUXds96xR58MFfcZoBzuuW3",
    "forwarderAuthority": "<forwarderAuthority>",
    "receiverAccounts": [
      { "publicKey": "<kvStoreAccount>", "isWritable": true }
    ]
  }
}
```

You're done — see the top-level [README.md](../../README.md) to run the workflow.

## Troubleshooting

- **`Error: Insufficient funds`** on deploy or initialize — airdrop more devnet SOL
  (`solana airdrop 2`), or wait a minute and retry (devnet faucet rate limits).
- **`DeclaredProgramIdMismatch`** — you built before running `anchor keys sync`, or
  edited `declare_id!` by hand. Run `anchor keys sync && anchor build` again, then
  redeploy.
- **`Program is not deployed` right after `anchor deploy`** — devnet RPC lag; wait
  10–20s and rerun `scripts/initialize.ts`.
- Want to redeploy from scratch (new program id)? Delete
  `target/deploy/kv_store_receiver-keypair.json`, then repeat from step 3.
