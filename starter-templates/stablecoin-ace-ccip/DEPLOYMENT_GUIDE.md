# Demo Complete Walkthrough

This guide is for a developer to deploy your own contracts and build the app from bottom-up, If you want to try the demo quickly, you can use existing deployment and refer to: [README.md](./README.md)

---

## Prerequisites

### Required Tools
- Foundry (`forge --version`)
- CRE CLI (`cre --version`)
- Bun v1.2.21+ (`bun --version`)

### Required Funds
- Sepolia testnet ETH ([faucet](https://faucets.chain.link))
- Fuji testnet AVAX ([faucet](https://faucets.chain.link))
- LINK on both chains (from faucet)

---

## Phase 1: Single-Chain Bank Integration

### Step 1.1: Initial Setup

Clone the repository:
```bash
git clone https://github.com/smartcontractkit/cre-demo-dapps
cd cre-demo-dapps
git checkout bank-stablecoin
```

Copy configuration templates:
```bash
cp .env.example .env
cp secrets.yaml.example secrets.yaml
```

Edit .env with your private key:
```bash
vim .env  # or nano .env
```

### Step 1.2: Install Contract Dependencies

```bash
bun add @openzeppelin/contracts
```

### Step 1.3: Set Environment Variables

```bash
# Load all environment variables from .env (includes PRIVATE_KEY, RPCs, etc.)
source .env

# Export private key for Foundry commands
export PRIVATE_KEY=$CRE_ETH_PRIVATE_KEY

# Set project root for payload file paths
export CRE_PROJECT_ROOT=$(pwd)
```

**Note:** `.env` includes `SEPOLIA_RPC` and `FUJI_RPC`. If experiencing rate limits, update these in `.env` with your own Alchemy/Infura endpoints.

### Step 1.4: Deploy StablecoinERC20 on Sepolia

```bash
forge create contracts/StablecoinERC20.sol:StablecoinERC20 \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY --broadcast \
  --constructor-args "BankStablecoin" "creUSD"
```

**💡 Copy the `Deployed to:` address from output**

```bash
export STABLECOIN_SEPOLIA=<paste_deployed_address_here>
```

### Step 1.5: Deploy MintingConsumer on Sepolia

**💡 Constructor arguments explained:**
- `$STABLECOIN_SEPOLIA` = Your deployed stablecoin address
- `0x0000000000000000000000000000000000000000` = Expected author (use address(0) for testing)
- `0x64756d6d790000000000` = Expected workflow name (we're using bytes10("dummy") for testing)

```bash
forge create contracts/MintingConsumer.sol:MintingConsumer \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY --broadcast \
  --constructor-args $STABLECOIN_SEPOLIA 0x0000000000000000000000000000000000000000 0x64756d6d790000000000
```

**💡 Copy the `Deployed to:` address from output**

```bash
export MINTING_CONSUMER_SEPOLIA=<paste_deployed_address_here>
```

### Step 1.6: Grant Roles to MintingConsumer

```bash
# Grant minter role
cast send $STABLECOIN_SEPOLIA \
  "grantMintRole(address)" \
  $MINTING_CONSUMER_SEPOLIA \
  --rpc-url $SEPOLIA_RPC \
  --private-key $PRIVATE_KEY

# Grant burner role
cast send $STABLECOIN_SEPOLIA \
  "grantBurnRole(address)" \
  $MINTING_CONSUMER_SEPOLIA \
  --rpc-url $SEPOLIA_RPC \
  --private-key $PRIVATE_KEY
```

### Step 1.7: Update Workflow Configuration

Edit `bank-stablecoin-workflow/config.json`:
```json
{
  "evms": [
    {
      "stablecoinAddress": "<your_STABLECOIN_SEPOLIA>",
      "consumerAddress": "<your_MINTING_CONSUMER_SEPOLIA>",
      "chainSelectorName": "ethereum-testnet-sepolia",
      "gasLimit": "500000"
    }
  ]
}
```

### Step 1.8: Install Workflow Dependencies

```bash
cd bank-stablecoin-workflow && bun install && cd ..
```

### Step 1.9: Test MINT

```bash
cre workflow simulate bank-stablecoin-workflow \
  --target local-simulation \
  --broadcast \
  --trigger-index 0 \
  --non-interactive \
  --http-payload @$CRE_PROJECT_ROOT/bank-stablecoin-workflow/http_trigger_payload.json
```

**Check the transaction hash in the output** to verify the mint was submitted on-chain.

### Step 1.10: Test REDEEM

```bash
cre workflow simulate bank-stablecoin-workflow \
  --target local-simulation \
  --broadcast \
  --trigger-index 0 \
  --non-interactive \
  --http-payload @$CRE_PROJECT_ROOT/bank-stablecoin-workflow/http_trigger_payload_redeem.json
```

**Expected:** 250 creUSD burned (check transaction hash in output).


---

## Phase 2: Cross-Chain CCIP Integration

**⚠️ Important:** Phase 2 involves ~15+ blockchain transactions across 2 chains. With public RPCs, expect:
- **Total time:** 30-60 minutes
- **Possible timeouts:** Steps 2.7-2.8 may show timeout errors (transactions often still succeed)
- **Recommendation:** Use your own RPC endpoint for faster/more reliable execution

### Step 2.1: Install CCIP and ACE Dependencies

**💡 Note:** Installing ACE now (even though we use it in Phase 3) because the ACE contracts need to compile alongside the rest of the codebase.

```bash
# Install CCIP contracts (required for TokenPools)
forge install smartcontractkit/ccip@ccip-develop

# Install ACE contracts (we'll use this in Phase 3)
forge install smartcontractkit/chainlink-ace

# Install ACE dependencies
cd lib/chainlink-ace && pnpm install && cd ../..
```

### Step 2.2: Deploy StablecoinERC20 on Fuji

**💡 Note:** FUJI_RPC was already set in Step 1.1

```bash
forge create contracts/StablecoinERC20.sol:StablecoinERC20 \
  --rpc-url $FUJI_RPC --private-key $PRIVATE_KEY --broadcast \
  --constructor-args "BankStablecoin" "creUSD"
```

```bash
export STABLECOIN_FUJI=<deployed_address>
```

### Step 2.3: Deploy MintingConsumer on Fuji

```bash
forge create contracts/MintingConsumer.sol:MintingConsumer \
  --rpc-url $FUJI_RPC --private-key $PRIVATE_KEY --broadcast \
  --constructor-args $STABLECOIN_FUJI 0x0000000000000000000000000000000000000000 0x64756d6d790000000000
```

**💡 Copy the `Deployed to:` address from output**

```bash
export MINTING_CONSUMER_FUJI=<paste_deployed_address_here>
```

### Step 2.4: Grant Roles on Fuji

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

### Step 2.5: Deploy & Configure CCIP Infrastructure (All-in-One Script)

**💡 This single script does everything:**
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

**💡 Expected output:**
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

**💡 Copy the pool addresses from output:**
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

### Step 2.6: Deploy CCIPTransferConsumer

**Set LINK token addresses:**
```bash
export LINK_SEPOLIA=0x779877A7B0D9E8603169DdbD7836e478b4624789
export LINK_FUJI=0x0b9d5D9136855f6FEc3c0993feE6E9CE8a297846
```

**💡 Constructor arguments explained:**
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

### Step 2.7: Fund Consumers with LINK

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

### Step 2.8: Update CCIP Workflow Configuration

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

### Step 2.9: Approve Consumer to Spend Tokens

**⚠️ CRITICAL STEP - DO NOT SKIP!**

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

**💡 Why this matters:** The consumer calls `transferFrom(sender, consumer, amount)` which requires prior approval from the sender.

**💡 Important:** The `sender.account` in your payload file must have:
- Sufficient token balance
- Approval granted to the consumer (Step 2.12)
- Match the private key you're using

In this demo, the deployer wallet (your private key) is both minting recipient and CCIP sender.

### Step 2.10: Test Cross-Chain Transfer

**💡 Prerequisites Check:**
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

**Transfer Sepolia → Fuji:**
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

## Phase 3: Advanced Multi-Service Integration (PoR + ACE + CCIP)

**🎯 Overview:** Phase 3 demonstrates an advanced stablecoin demo combining:
- **Proof of Reserve (PoR)** - Validates sufficient off-chain reserves before minting
- **Automated Compliance Engine (ACE)** - Enforces address blacklist and volume limit policies
- **CCIP** - Cross-chain transfers with compliance checks (optional)

**Prerequisites:** Phase 1 & 2 must be completed successfully.

---

## Part 1: Proof of Reserve (PoR) Testing

### Step 3.1: Install Workflow Dependencies

```bash
cd bank-stablecoin-por-ace-ccip-workflow && bun install && cd ..
```

### Step 3.2: What is Proof of Reserve (PoR)?

**Proof of Reserve (PoR)** validates that sufficient off-chain reserves exist before minting new stablecoins. This prevents over-minting and ensures 1:1 backing.

**How it works in this demo:**
- CRE workflow fetches reserve data via HTTP capability
- Compares requested mint amount against available reserves
- If insufficient reserves → workflow fails before any on-chain transaction
- Demo uses a mock JSON file (`mock-por-response.json`) - real implementations would use a live API

### Step 3.3: What is ACE?

**Automated Compliance Engine (ACE)** enforces on-chain compliance policies before transactions execute. It acts as a policy layer that can block unauthorized operations.

**How it works:**
- `PolicyEngine` orchestrates policy checks
- `Extractors` parse transaction data to extract parameters (e.g., beneficiary address, amount)
- `Policies` evaluate parameters (e.g., is address blacklisted? is amount within limits?)
- If any policy rejects → transaction reverts before execution

**Policies in this demo:**
- `AddressBlacklistPolicy` - Blocks mints/transfers to blacklisted addresses
- `VolumePolicy` - Enforces min/max transfer amounts for CCIP (100-10,000 creUSD)

### Step 3.4: Deploy ACE Infrastructure

**💡 Note:** ACE dependencies were already installed in Phase 2, Step 2.1

**Deploy PolicyEngine + BlacklistPolicy:**
```bash
ETHERSCAN_API_KEY=dummy forge script script/DeployACESystem.s.sol:DeployACESystem \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY --broadcast
```

**💡 Save the deployed addresses:**
```bash
export POLICY_ENGINE=<PolicyEngine_proxy_address>
export BLACKLIST_POLICY=<BlacklistPolicy_proxy_address>
```

**Deploy UnifiedExtractor:**
```bash
ETHERSCAN_API_KEY=dummy forge script script/DeployUnifiedExtractor.s.sol:DeployUnifiedExtractor \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY --broadcast
```

**💡 Save the deployed address:**
```bash
export UNIFIED_EXTRACTOR=<UnifiedExtractor_address>
```

**Deploy VolumePolicy:**
```bash
forge script script/DeployVolumePolicy.s.sol:DeployVolumePolicy \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY --broadcast
```

**💡 Save the deployed address:**
```bash
export VOLUME_POLICY=<VolumePolicy_proxy_address>
```

**💡 What we deployed:**
- `PolicyEngine` - Orchestrates all policy checks
- `BlacklistPolicy` - Blocks blacklisted addresses
- `UnifiedExtractor` - Parses both mint and CCIP reports
- `VolumePolicy` - Enforces amount limits (100-10,000 creUSD)

### Step 3.5: Deploy ACE Consumers

**💡 Note:** Consumers are initialized with the PolicyEngine from Step 3.4. We won't attach policies until Step 3.9 - this allows us to test PoR first.

**Set required environment variables:**
```bash
export SEPOLIA_ROUTER=0x0BF3dE8c5D3e8A2B34D2BEeB17ABfCeBaf363A59
export LINK_SEPOLIA=0x779877A7B0D9E8603169DdbD7836e478b4624789
```

**Deploy consumers:**
```bash
# Pass PolicyEngine address directly to ensure consumers use the correct one
POLICY_ENGINE=<YOUR_POLICY_ENGINE_FROM_STEP_3.4> \
ETHERSCAN_API_KEY=dummy \
forge script script/DeployACEConsumers.s.sol:DeployACEConsumers \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY --broadcast \
  --sig "run(bool)" true
```

**💡 Save the deployed addresses:**
```bash
export MINTING_CONSUMER_ACE=<MintingConsumerWithACE_proxy_address>
export CCIP_CONSUMER_ACE=<CCIPTransferConsumerWithACE_proxy_address>
```

**Example output:**
```
MintingConsumerWithACE: 0x24c0f5C1A286Fbd27A730303a1a845b4cf85F0Cc
CCIPTransferConsumerWithACE: 0xFa031de805af3a9A72D37f57a01634ADF4a61cD5
```

### Step 3.6: Configure Workflow for PoR Testing

Edit `bank-stablecoin-por-ace-ccip-workflow/config.json`:

```json
{
  "sepolia": {
    "stablecoinAddress": "<your_STABLECOIN_SEPOLIA>",
    "mintingConsumerAddress": "<your_MINTING_CONSUMER_ACE>",
    "ccipConsumerAddress": "<your_CCIP_CONSUMER_ACE>",
    "policyEngineAddress": "<your_POLICY_ENGINE>",
    "blacklistPolicyAddress": "<your_BLACKLIST_POLICY>",
    "chainSelector": "16015286601757825753"
  },
  "fuji": {
    "stablecoinAddress": "<your_STABLECOIN_FUJI>",
    "chainSelector": "14767482510784806043"
  },
  "porApiUrl": "file://./mock-por-response.json",
  "decimals": 18
}
```

**💡 Note:** Use the actual addresses from Step 3.4 for `policyEngineAddress` and `blacklistPolicyAddress`.

### Step 3.7: Test PoR - Success Case

**Test with sufficient reserves:**

```bash
export CRE_PROJECT_ROOT=$(pwd)

cre workflow simulate bank-stablecoin-por-ace-ccip-workflow \
  --target local-simulation \
  --broadcast \
  --trigger-index 0 \
  --non-interactive \
  --http-payload @$CRE_PROJECT_ROOT/bank-stablecoin-por-ace-ccip-workflow/http_trigger_payload.json
```

**Expected output:**
```
[PoR Validation] Fetching reserve data...
Using mock PoR data for demo
Reserve Data: 500000 USD
✓ PoR validation passed - reserves (500000 USD) can cover mint
✓ Mint transaction submitted: 0x...
```

### Step 3.8: Test PoR - Failure Case (Overminting)

**Edit payload to exceed reserves:**

```bash
vim bank-stablecoin-por-ace-ccip-workflow/http_trigger_payload.json
# Change: "amount": "1000"
# To:     "amount": "600000"
```

**Run workflow (will fail before sending TX):**

```bash
cre workflow simulate bank-stablecoin-por-ace-ccip-workflow \
  --target local-simulation \
  --trigger-index 0 \
  --non-interactive \
  --http-payload @$CRE_PROJECT_ROOT/bank-stablecoin-por-ace-ccip-workflow/http_trigger_payload.json
```

**CRE should return with execution error** - workflow fails with `POR_INSUFFICIENT_RESERVES`. No transaction sent.

**Restore payload to default:**
```bash
vim bank-stablecoin-por-ace-ccip-workflow/http_trigger_payload.json
# Change back: "amount": "1000"
```

---

## Part 2: ACE Configuration & Testing

### Step 3.9: Configure ACE (Attach Policies)

**⚠️ CRITICAL:** ACE uses `keccak256("parameterName")` convention, NOT `bytes32("parameterName")`!

**Run configuration script:**
```bash
# Pass all ACE addresses directly to ensure correct configuration
POLICY_ENGINE=<YOUR_POLICY_ENGINE> \
BLACKLIST_POLICY=<YOUR_BLACKLIST_POLICY> \
UNIFIED_EXTRACTOR=<YOUR_UNIFIED_EXTRACTOR> \
VOLUME_POLICY=<YOUR_VOLUME_POLICY> \
MINTING_CONSUMER_ACE=<YOUR_MINTING_CONSUMER> \
CCIP_CONSUMER_ACE=<YOUR_CCIP_CONSUMER> \
ETHERSCAN_API_KEY=dummy \
forge script script/ConfigureACEWithConstants.s.sol:ConfigureACEWithConstants \
  --rpc-url $SEPOLIA_RPC \
  --broadcast \
  --legacy
```

The script will:
1. Verify `UnifiedExtractor` is attached to `onReport` selector
2. Attach `BlacklistPolicy` to MintingConsumer (checks `beneficiary` parameter)
3. Attach `VolumePolicy` to CCIPConsumer (checks `amount` parameter)

**Fund CCIP Consumer with LINK:**
```bash
cast send $LINK_SEPOLIA \
  "transfer(address,uint256)" \
  $CCIP_CONSUMER_ACE \
  5000000000000000000 \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY
```

### Step 3.10: Test ACE - Mint Blacklist Blocking

**Blacklist a test address:**
```bash
cast send $BLACKLIST_POLICY \
  "addToBlacklist(address)" \
  0x15ee0018efda0f7689abd704b20c8c1c9b8e52b7 \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY

# Verify blacklisted
cast call $BLACKLIST_POLICY \
  "isBlacklisted(address)(bool)" \
  0x15ee0018efda0f7689abd704b20c8c1c9b8e52b7 \
  --rpc-url $SEPOLIA_RPC
# Should return: true
```

**Check balance before:**
```bash
cast call $STABLECOIN_SEPOLIA \
  "balanceOf(address)(uint256)" \
  0x15ee0018efda0f7689abd704b20c8c1c9b8e52b7 \
  --rpc-url $SEPOLIA_RPC
# Should return: 0
```

**Edit payload to use blacklisted address:**
```bash
vim bank-stablecoin-por-ace-ccip-workflow/http_trigger_payload.json
# Change: "account": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"
# To: "account": "0x15ee0018efda0f7689abd704b20c8c1c9b8e52b7"
```

**Run workflow (ACE should block):**
```bash
cre workflow simulate bank-stablecoin-por-ace-ccip-workflow \
  --target local-simulation \
  --broadcast \
  --trigger-index 0 \
  --non-interactive \
  --http-payload @$CRE_PROJECT_ROOT/bank-stablecoin-por-ace-ccip-workflow/http_trigger_payload.json
```

**Verify ACE blocked:**
```bash
# Check balance (should still be 0)
cast call $STABLECOIN_SEPOLIA \
  "balanceOf(address)(uint256)" \
  0x15ee0018efda0f7689abd704b20c8c1c9b8e52b7 \
  --rpc-url $SEPOLIA_RPC
# Should return: 0 (ACE blocked the mint)
```

**Restore payload:**
```bash
vim bank-stablecoin-por-ace-ccip-workflow/http_trigger_payload.json
# Change back to: "account": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"
```

### Step 3.11: Test ACE - CCIP Volume Limits

**Test Scenario 1: Amount Within Range (500 creUSD) - Should ALLOW**

Run test:
```bash
cre workflow simulate bank-stablecoin-por-ace-ccip-workflow \
  --target local-simulation \
  --broadcast \
  --trigger-index 0 \
  --non-interactive \
  --http-payload @$CRE_PROJECT_ROOT/bank-stablecoin-por-ace-ccip-workflow/test_volume_within_range.json
```

Verify CCIP transfer succeeded:
```bash
# Check for burn event (Transfer to zero address)
cast receipt <mintTransaction_hash> --rpc-url $SEPOLIA_RPC --json | \
jq --arg addr "$STABLECOIN_SEPOLIA" \
   '.logs[] | select(
       (.address | ascii_downcase == ($addr | ascii_downcase)) and
       .topics[0] == "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" and
       .topics[2] == "0x0000000000000000000000000000000000000000000000000000000000000000"
   )'
# Non-empty output → Tokens burned for CCIP ✅
```

**Check CCIP Explorer:**
```bash
# Copy ccipTransaction from workflow output and paste into CCIP Explorer:
# https://ccip.chain.link/msg/<ccipTransaction_hash>
```

**Test Scenario 2: Amount Below Minimum (50 creUSD) - Should BLOCK**

Run test:
```bash
cre workflow simulate bank-stablecoin-por-ace-ccip-workflow \
  --target local-simulation \
  --broadcast \
  --trigger-index 0 \
  --non-interactive \
  --http-payload @$CRE_PROJECT_ROOT/bank-stablecoin-por-ace-ccip-workflow/test_volume_below_min.json
```

Verify blocking:
```bash
# Check for burn event (Transfer to zero address)
cast receipt <mintTransaction_hash> --rpc-url $SEPOLIA_RPC --json | \
jq --arg addr "$STABLECOIN_SEPOLIA" \
   '.logs[] | select(
       (.address | ascii_downcase == ($addr | ascii_downcase)) and
       .topics[0] == "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" and
       .topics[2] == "0x0000000000000000000000000000000000000000000000000000000000000000"
   )'
# Empty output → VolumePolicy blocked ✅
# (No CCIP Explorer check - transfer never happened)
```

**Test Scenario 3: Amount Above Maximum (15,000 creUSD) - Should BLOCK**

Run test:
```bash
cre workflow simulate bank-stablecoin-por-ace-ccip-workflow \
  --target local-simulation \
  --broadcast \
  --trigger-index 0 \
  --non-interactive \
  --http-payload @$CRE_PROJECT_ROOT/bank-stablecoin-por-ace-ccip-workflow/test_volume_above_max.json
```

Verify blocking:
```bash
# Check for burn event (Transfer to zero address)
cast receipt <mintTransaction_hash> --rpc-url $SEPOLIA_RPC --json | \
jq --arg addr "$STABLECOIN_SEPOLIA" \
   '.logs[] | select(
       (.address | ascii_downcase == ($addr | ascii_downcase)) and
       .topics[0] == "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" and
       .topics[2] == "0x0000000000000000000000000000000000000000000000000000000000000000"
   )'
# Empty output → VolumePolicy blocked ✅
# (No CCIP Explorer check - transfer never happened)
```

---

## Reference: Managing ACE Policies

After deployment, you can manage ACE policies using these commands:

**Add address to blacklist:**
```bash
cast send $BLACKLIST_POLICY \
  "addToBlacklist(address)" \
  <address_to_block> \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY
```

**Remove address from blacklist:**
```bash
cast send $BLACKLIST_POLICY \
  "removeFromBlacklist(address)" \
  <address_to_unblock> \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY
```

**Check if address is blacklisted:**
```bash
cast call $BLACKLIST_POLICY \
  "isBlacklisted(address)(bool)" \
  <address> \
  --rpc-url $SEPOLIA_RPC
```

**Transfer policy ownership:**
```bash
# PolicyEngine, BlacklistPolicy, VolumePolicy, and Consumers are all Ownable
cast send $BLACKLIST_POLICY \
  "transferOwnership(address)" \
  <new_owner> \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY
```

---

## Troubleshooting

For common issues and debugging steps, see [`TROUBLESHOOTING.md`](./TROUBLESHOOTING.md).

