# Phase 3: Production-Ready with PoR + ACE

Add Proof of Reserve validation and Automated Compliance Engine to your stablecoin.

**Builds on:** [Phase 1](../bank-stablecoin-workflow/README.md) and [Phase 2](../ccip-transfer-workflow/README.md) - Complete both phases first

**What you'll add:**
- Proof of Reserve (PoR) validation
- Automated Compliance Engine (ACE)
- Address blacklist policy
- Volume limit policy for CCIP

**Time required:** ~45 minutes

---

## Overview

Phase 3 combines multiple Chainlink services:
- **Proof of Reserve (PoR)** - Validates sufficient off-chain reserves before minting
- **Automated Compliance Engine (ACE)** - Enforces address blacklist and volume limit policies
- **CCIP** - Cross-chain transfers with compliance checks (optional)

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

**Note:** ACE dependencies were already installed in Phase 2, Step 2.1

**Deploy PolicyEngine + BlacklistPolicy:**
```bash
ETHERSCAN_API_KEY=dummy forge script script/DeployACESystem.s.sol:DeployACESystem \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY --broadcast
```

Save the deployed addresses:
```bash
export POLICY_ENGINE=<PolicyEngine_proxy_address>
export BLACKLIST_POLICY=<BlacklistPolicy_proxy_address>
```

**Deploy UnifiedExtractor:**
```bash
ETHERSCAN_API_KEY=dummy forge script script/DeployUnifiedExtractor.s.sol:DeployUnifiedExtractor \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY --broadcast
```

Save the deployed address:
```bash
export UNIFIED_EXTRACTOR=<UnifiedExtractor_address>
```

**Deploy VolumePolicy:**
```bash
forge script script/DeployVolumePolicy.s.sol:DeployVolumePolicy \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY --broadcast
```

Save the deployed address:
```bash
export VOLUME_POLICY=<VolumePolicy_proxy_address>
```

**What we deployed:**
- `PolicyEngine` - Orchestrates all policy checks
- `BlacklistPolicy` - Blocks blacklisted addresses
- `UnifiedExtractor` - Parses both mint and CCIP reports
- `VolumePolicy` - Enforces amount limits (100-10,000 creUSD)

### Step 3.5: Deploy ACE Consumers

**Note:** Consumers are initialized with the PolicyEngine from Step 3.4. We won't attach policies until Step 3.9 - this allows us to test PoR first.

**Note:** `SEPOLIA_ROUTER` and `LINK_SEPOLIA` were already set in Phase 2.

**Deploy consumers:**
```bash
# Pass PolicyEngine address directly to ensure consumers use the correct one
POLICY_ENGINE=<YOUR_POLICY_ENGINE_FROM_STEP_3.4> \
ETHERSCAN_API_KEY=dummy \
forge script script/DeployACEConsumers.s.sol:DeployACEConsumers \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY --broadcast \
  --sig "run(bool)" true
```

Save the deployed addresses:
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

**Note:** Use the actual addresses from Step 3.4 for `policyEngineAddress` and `blacklistPolicyAddress`.

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

**CRITICAL:** ACE uses `keccak256("parameterName")` convention, NOT `bytes32("parameterName")`!

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

## Phase 3 Complete

**Previous:** [Phase 2: Cross-Chain CCIP](../ccip-transfer-workflow/README.md) | [Phase 1: Basic Stablecoin](../bank-stablecoin-workflow/README.md)

**See also:** [Complete Deployment Guide](../DEPLOYMENT_GUIDE.md) | [Main README](../README.md) | [Troubleshooting](../TROUBLESHOOTING.md)

