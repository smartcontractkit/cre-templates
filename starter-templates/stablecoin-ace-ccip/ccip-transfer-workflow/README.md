# Phase 2: Cross-Chain CCIP Integration

Add cross-chain token transfers using Chainlink CCIP.

**Builds on:** [Phase 1](../bank-stablecoin-workflow/README.md) - Complete Phase 1 first

**What you'll add:**
- CCIP infrastructure (TokenPools, routes)
- Cross-chain transfer consumers
- Bidirectional transfers (Sepolia ↔ Fuji)

**Time required:** ~30-60 minutes

**Important:** Phase 2 involves ~15+ blockchain transactions across 2 chains. With public RPCs, expect possible timeouts. Use your own RPC endpoint for faster/more reliable execution.

---

## Step 2.1: Install CCIP and ACE Dependencies

**Note:** Installing ACE now (even though we use it in Phase 3) because the ACE contracts need to compile alongside the rest of the codebase.

```bash
# Install CCIP contracts (required for TokenPools)
forge install smartcontractkit/ccip@ccip-develop

# Install ACE contracts (we'll use this in Phase 3)
forge install smartcontractkit/chainlink-ace

# Install ACE dependencies
cd lib/chainlink-ace && pnpm install && cd ../..
```

## Step 2.2: Deploy StablecoinERC20 on Fuji

**Note:** FUJI_RPC was already set in Phase 1

```bash
forge create contracts/StablecoinERC20.sol:StablecoinERC20 \
  --rpc-url $FUJI_RPC --private-key $PRIVATE_KEY --broadcast \
  --constructor-args "BankStablecoin" "creUSD"
```

```bash
export STABLECOIN_FUJI=<deployed_address>
```

## Step 2.3: Deploy MintingConsumer on Fuji

```bash
forge create contracts/MintingConsumer.sol:MintingConsumer \
  --rpc-url $FUJI_RPC --private-key $PRIVATE_KEY --broadcast \
  --constructor-args $STABLECOIN_FUJI 0x0000000000000000000000000000000000000000 0x64756d6d790000000000
```

Copy the `Deployed to:` address from output:

```bash
export MINTING_CONSUMER_FUJI=<paste_deployed_address_here>
```

## Step 2.4: Grant Roles on Fuji

```bash
cast send $STABLECOIN_FUJI \
  "grantMintRole(address)" \
  $MINTING_CONSUMER_FUJI \
  --rpc-url $FUJI_RPC \
  --private-key $PRIVATE_KEY

cast send $STABLECOIN_FUJI \
  "grantBurnRole(address)" \
  $MINTING_CONSUMER_FUJI \
  --rpc-url $FUJI_RPC \
  --private-key $PRIVATE_KEY
```

## Step 2.5: Deploy & Configure CCIP Infrastructure (All-in-One Script)

**This single script does everything:**
1. Deploys BurnMintTokenPool on Sepolia
2. Deploys BurnMintTokenPool on Fuji
3. Grants mint/burn roles on both chains
4. Registers pools with TokenAdminRegistry
5. Configures bidirectional routes (Sepolia ↔ Fuji)

**Run the all-in-one script:**
```bash
# Pass your deployed stablecoin addresses directly
STABLECOIN_SEPOLIA=<YOUR_SEPOLIA_STABLECOIN> \
STABLECOIN_FUJI=<YOUR_FUJI_STABLECOIN> \
ETHERSCAN_API_KEY=dummy \
forge script script/DeployCCIP.s.sol:DeployCCIP \
  --broadcast \
  --slow
```

**Expected output:**
```
============================================================
DEPLOYING CCIP INFRASTRUCTURE (SEPOLIA <> FUJI)
============================================================

PHASE 1: SEPOLIA SETUP
-----------------------------------------------------------
[ Sepolia ] Deploying BurnMintTokenPool...
[ Sepolia ] Pool deployed: 0x...
  Granting mint/burn roles to pool...
  Registering admin via RegistryModuleOwnerCustom...
  Accepting admin role...
  Setting pool in TokenAdminRegistry...
Sepolia TokenPool: 0x...

PHASE 2: FUJI SETUP
-----------------------------------------------------------
[ Fuji ] Deploying BurnMintTokenPool...
[ Fuji ] Pool deployed: 0x...
  (same registration steps)
Fuji TokenPool: 0x...

PHASE 3: CONFIGURE BIDIRECTIONAL ROUTES
-----------------------------------------------------------
[Route] Sepolia -> Fuji
  Route configured
[Route] Fuji -> Sepolia
  Route configured

============================================================
CCIP DEPLOYMENT COMPLETE
============================================================

Add to .env:
  POOL_SEPOLIA= 0x...
  POOL_FUJI= 0x...
```

Copy the pool addresses from output:
```bash
export POOL_SEPOLIA=<sepolia_pool_from_output>
export POOL_FUJI=<fuji_pool_from_output>
```

**Verify routes:**
```bash
# Check Sepolia pool knows about Fuji
cast call $POOL_SEPOLIA "isSupportedChain(uint64)(bool)" 14767482510784806043 --rpc-url $SEPOLIA_RPC
# Should return: true

# Check Fuji pool knows about Sepolia
cast call $POOL_FUJI "isSupportedChain(uint64)(bool)" 16015286601757825753 --rpc-url $FUJI_RPC
# Should return: true
```

## Step 2.6: Deploy CCIPTransferConsumer

**Set LINK token addresses:**
```bash
export LINK_SEPOLIA=0x779877A7B0D9E8603169DdbD7836e478b4624789
export LINK_FUJI=0x0b9d5D9136855f6FEc3c0993feE6E9CE8a297846
```

**Set CCIP Router addresses:**
```bash
export SEPOLIA_ROUTER=0x0BF3dE8c5D3e8A2B34D2BEeB17ABfCeBaf363A59
export FUJI_ROUTER=0xF694E193200268f9a4868e4Aa017A0118C9a8177
```

**Constructor arguments explained:**
- `$STABLECOIN` = Your stablecoin address
- `$ROUTER` = CCIP Router address
- `$LINK` = LINK token address (for fees)
- `0x0000...` = Expected author (use address(0) for testing)
- `0x6475...` = Expected workflow name (bytes10("dummy") for testing)

**Deploy on Sepolia:**
```bash
forge create contracts/CCIPTransferConsumer.sol:CCIPTransferConsumer \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY --broadcast \
  --constructor-args $STABLECOIN_SEPOLIA $SEPOLIA_ROUTER $LINK_SEPOLIA 0x0000000000000000000000000000000000000000 0x64756d6d790000000000
```

```bash
export CCIP_CONSUMER_SEPOLIA=<paste_deployed_address_here>
```

**Deploy on Fuji:**
```bash
forge create contracts/CCIPTransferConsumer.sol:CCIPTransferConsumer \
  --rpc-url $FUJI_RPC --private-key $PRIVATE_KEY --broadcast \
  --constructor-args $STABLECOIN_FUJI $FUJI_ROUTER $LINK_FUJI 0x0000000000000000000000000000000000000000 0x64756d6d790000000000
```

```bash
export CCIP_CONSUMER_FUJI=<paste_deployed_address_here>
```

## Step 2.7: Fund Consumers with LINK

**Sepolia:**
```bash
cast send $LINK_SEPOLIA \
  "transfer(address,uint256)" \
  $CCIP_CONSUMER_SEPOLIA \
  5000000000000000000 \
  --rpc-url $SEPOLIA_RPC \
  --private-key $PRIVATE_KEY
```

**Fuji:**
```bash
cast send $LINK_FUJI \
  "transfer(address,uint256)" \
  $CCIP_CONSUMER_FUJI \
  5000000000000000000 \
  --rpc-url $FUJI_RPC \
  --private-key $PRIVATE_KEY
```

## Step 2.8: Update CCIP Workflow Configuration

Edit `ccip-transfer-workflow/config.json`:
```json
{
  "chains": {
    "ethereum-testnet-sepolia": {
      "chainSelector": "16015286601757825753",
      "consumerAddress": "<your_CCIP_CONSUMER_SEPOLIA>",
      "stablecoinAddress": "<your_STABLECOIN_SEPOLIA>"
    },
    "avalanche-testnet-fuji": {
      "chainSelector": "14767482510784806043",
      "consumerAddress": "<your_CCIP_CONSUMER_FUJI>",
      "stablecoinAddress": "<your_STABLECOIN_FUJI>"
    }
  },
  "gasLimit": "500000"
}
```

## Step 2.9: Approve Consumer to Spend Tokens

**CRITICAL STEP - DO NOT SKIP!**

The CCIPTransferConsumer needs permission to spend your tokens before initiating cross-chain transfers. Without this approval, transfers will fail silently.

**On Sepolia:**
```bash
cast send $STABLECOIN_SEPOLIA \
  "approve(address,uint256)" \
  $CCIP_CONSUMER_SEPOLIA \
  1000000000000000000000 \
  --rpc-url $SEPOLIA_RPC \
  --private-key $PRIVATE_KEY
```

**Expected output:** Transaction hash with `status: 1 (success)`

**On Fuji:**
```bash
cast send $STABLECOIN_FUJI \
  "approve(address,uint256)" \
  $CCIP_CONSUMER_FUJI \
  1000000000000000000000 \
  --rpc-url $FUJI_RPC \
  --private-key $PRIVATE_KEY
```

**Why this matters:** The consumer calls `transferFrom(sender, consumer, amount)` which requires prior approval from the sender.

**Important:** The `sender.account` in your payload file must have:
- Sufficient token balance
- Approval granted to the consumer (Step 2.9)
- Match the private key you're using

In this demo, the deployer wallet (your private key) is both minting recipient and CCIP sender.

## Step 2.10: Test Cross-Chain Transfer

**Prerequisites Check:**
Before testing the transfer, ensure your sender account has tokens. If you used the beneficiary address in Phase 1, you can either:
- Option A: Use the beneficiary address (has 750 creUSD from Phase 1)
- Option B: Mint tokens to your deployer address (recommended for testing)

**If using deployer as sender (recommended):**
```bash
# Edit the deployer payload with your address
vim bank-stablecoin-workflow/http_trigger_payload_deployer.json
# Change "account": "YOUR_DEPLOYER_ADDRESS_HERE" to your deployer address
```

```bash
# Mint tokens to your deployer
cre workflow simulate bank-stablecoin-workflow \
  --target local-simulation \
  --broadcast \
  --trigger-index 0 \
  --non-interactive \
  --http-payload @$CRE_PROJECT_ROOT/bank-stablecoin-workflow/http_trigger_payload_deployer.json
```

**Install dependencies:**
```bash
cd ccip-transfer-workflow
bun install
cd ..
```

**Transfer Sepolia to Fuji:**
```bash
cre workflow simulate ccip-transfer-workflow \
  --target local-simulation \
  --broadcast \
  --trigger-index 0 \
  --non-interactive \
  --http-payload @$CRE_PROJECT_ROOT/ccip-transfer-workflow/http_trigger_payload.json
```

**Monitor transfer:**
- Copy messageId from logs
- Visit: `https://ccip.chain.link/msg/<messageId>`
- Wait 10-20 minutes for delivery

**Verify on Fuji:**
```bash
# After 10-20 minutes, check beneficiary's balance
# (Use the beneficiary address from your http_trigger_payload.json)
cast call $STABLECOIN_FUJI \
  "balanceOf(address)(uint256)" \
  <YOUR_BENEFICIARY_ADDRESS> \
  --rpc-url $FUJI_RPC
```

Should show 100 creUSD received on the beneficiary address.

---

## Phase 2 Complete

**Previous:** [Phase 1: Basic Stablecoin](../bank-stablecoin-workflow/README.md)

**Next:** [Phase 3: Add PoR + ACE](../bank-stablecoin-por-ace-ccip-workflow/README.md)

**See also:** [Complete Deployment Guide](../DEPLOYMENT_GUIDE.md) | [Main README](../README.md)

