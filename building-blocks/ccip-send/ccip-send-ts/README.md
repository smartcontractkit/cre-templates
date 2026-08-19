<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE\&message=Home\&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE\&message=Docs\&color=blue)](https://docs.chain.link/cre)

</div>

# Send a CCIP Message (TypeScript)

A minimal example: on a cron schedule, this workflow sends a CCIP message from Ethereum Sepolia to Arbitrum Sepolia
using [`cre-ccip-sdk`](https://github.com/smartcontractkit/cre-libraries-ts/tree/main/packages/cre-ccip-sdk) — no
tokens, just a text payload — to demonstrate the minimum wiring needed to use the library from a real workflow.

For token transfers, pull-based funding, and the full API, see the
[`cre-ccip-sdk` README](https://github.com/smartcontractkit/cre-libraries-ts/tree/main/packages/cre-ccip-sdk).

## How it works

`cre-ccip-sdk` needs a `CCIPSenderReceiver` contract deployed onchain — CRE workflows submit writes as signed
reports to a receiver contract you control, not as arbitrary calls, so this contract is what actually calls the
CCIP Router's `ccipSend()` on your workflow's behalf. See
["Why this needs a contract, not just a package"](https://github.com/smartcontractkit/cre-libraries-ts/tree/main/packages/cre-ccip-sdk#why-this-needs-a-contract-not-just-a-package)
in the library's README for the full explanation.

## Quick start

> **This template already points at a deployed, funded demo contract** (`0x4A2E52B12B4d0b35Bf12AF535f753DE26C2a10e8`
> on Ethereum Sepolia, configured with the `MockKeystoneForwarder` for local simulation). You can skip to
> [step 4](#4-install-dependencies) to try a real send immediately. Steps 1 and 3 below show how to point the
> workflow at your own contract instead — do this once the demo contract's funding runs low (each send costs
> ~0.0003 Sepolia ETH in CCIP fees), or before adapting this template beyond a first test.

### 1) (Optional) Deploy and fund your own `CCIPSenderReceiver`

The contract source ships inside the `cre-ccip-sdk` npm package (`contracts/CCIPSenderReceiver.sol`). Copy it,
along with `contracts/vendor/`, into a Foundry or Hardhat project, and deploy it to Ethereum Sepolia:

```solidity
new CCIPSenderReceiver(
  0xF8344CFd5c43616a4366C34E3EEE75af79a74482, // Sepolia KeystoneForwarder (production)
  0x0BF3dE8c5D3e8A2B34D2BEeB17ABfCeBaf363A59  // Sepolia CCIP Router
)
```

> **Local simulation uses a different forwarder than production.** `cre workflow simulate --broadcast` (step 5)
> submits the write through a `MockKeystoneForwarder`, not the production `KeystoneForwarder` — deploying against
> the production forwarder address above will make every local-simulation write revert with no useful error, since
> the local run's signatures aren't registered with the production forwarder's config. For a contract you'll only
> use to test locally, deploy with Ethereum Sepolia's mock forwarder instead:
> `0x15fC6ae953E024d975e77382eEeC56A9101f9F88`. Redeploy against the production address (or call
> `setForwarderAddress()`) before using this contract with a real, deployed-and-registered workflow. See the
> [forwarder directory](https://docs.chain.link/cre/guides/workflow/using-evm-client/forwarder-directory-ts) for
> other networks' mock/production addresses.

This example pays the CCIP fee in native currency (Sepolia ETH), so send some to the deployed contract address —
it has a `receive()` function. Since this example doesn't transfer tokens, that's the only funding step needed.
`checkReadiness()` reads balances at the last **finalized** block, which lags the chain head by several minutes on
Sepolia — if you just funded the contract, wait a few minutes before simulating or you may see a funding-check
failure even though the funds have technically arrived.

Full deploy/fund/lock-down steps (including pull-based funding and token transfers) are in the
[`cre-ccip-sdk` README](https://github.com/smartcontractkit/cre-libraries-ts/tree/main/packages/cre-ccip-sdk#setup).

### 2) Configure RPC in `project.yaml`

Already set to a public Sepolia RPC — replace it with your own if you hit rate limits:

```yaml
rpcs:
  - chain-name: ethereum-testnet-sepolia
    url: <YOUR_SEPOLIA_RPC_URL>
```

### 3) (Optional) Point the workflow at your own contract

`my-workflow/config.staging.json` already has a working configuration:

```json
{
  "schedule": "0 */10 * * * *",
  "chainName": "ethereum-testnet-sepolia",
  "routerAddress": "0x0BF3dE8c5D3e8A2B34D2BEeB17ABfCeBaf363A59",
  "receiverContractAddress": "0x4A2E52B12B4d0b35Bf12AF535f753DE26C2a10e8",
  "destinationChainSelector": "3478487238524512106",
  "destinationReceiver": "0x000000000000000000000000000000000000dEaD",
  "message": "Hello from CRE!"
}
```

- `receiverContractAddress`: the shared demo `CCIPSenderReceiver` from step 1. Once you deploy your own (or if
  this one runs out of funding), replace it with your deployed address.
- `destinationChainSelector`: `3478487238524512106` is Arbitrum Sepolia. See the
  [CCIP Directory](https://docs.chain.link/ccip/directory) for other destinations.
- `destinationReceiver`: any address on the destination chain — this example sends a plain text payload, not a
  message meant to be decoded by a specific contract, so any address works for a first test. The placeholder
  (`0x000000000000000000000000000000000000dEaD`, a well-known Ethereum burn address) is fine to leave as-is.

### 4) Install dependencies

From the project root:

```bash
bun install --cwd ./my-workflow
```

> `cre-ccip-sdk` isn't published to npm yet, so `my-workflow/package.json` points `cre-ccip-sdk` at a local
> `file:` path into a `cre-libraries-ts` checkout (`file:../../../../../cre-libraries-ts/packages/cre-ccip-sdk`).
> Adjust that path if your `cre-libraries-ts` checkout lives somewhere else, and switch it to a real version
> range (e.g. `^0.1.0`) once the package is published.

### 5) Run a local simulation

From the project root. `--broadcast` is required — without it, the workflow still computes `getFee()` and
`checkReadiness()`, but there's no real forwarder to deliver the write to, so `send()` fails on the write step:

```bash
cre workflow simulate my-workflow --target staging-settings --broadcast
```

This has been run end-to-end against a real deployed `CCIPSenderReceiver` on Sepolia — you should see output like:

```
Workflow compiled
2026-08-18T15:32:46Z [SIMULATION] Simulator Initialized

2026-08-18T15:32:46Z [SIMULATION] Running trigger trigger=cron-trigger@1.0.0
2026-08-18T15:40:24Z [USER LOG] CCIP message sent: txHash=0x293c80883f4501cfb5b3824c1529f3676d4fc6f625dbce2025ec82cb77bae8a8

Workflow Simulation Result:
 "sent:0x293c80883f4501cfb5b3824c1529f3676d4fc6f625dbce2025ec82cb77bae8a8"
```

If the receiver contract isn't funded yet (or the funding transaction hasn't finalized — see the note in step 1),
you'll see a `checkReadiness()` failure instead — that's the library catching the problem before submitting
anything onchain:

```
2026-08-18T15:10:00Z [USER LOG] CCIP send failed: Error: CCIP send aborted, request is not funded: native balance of receiver 0x4A2E52B12B4d0b35Bf12AF535f753DE26C2a10e8 is 0, need 280445790582210 for the CCIP fee

Workflow Simulation Result:
 "failed:Error: CCIP send aborted, request is not funded: ..."
```
