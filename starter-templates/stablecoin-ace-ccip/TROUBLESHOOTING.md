# Troubleshooting Guide

This guide covers common issues and debugging steps for the CRE Stablecoin Demo.

---

## Workflow Shows "Timeout" But Transactions Succeed

**Common with CRE workflows using public RPCs**

If CRE workflows show "Timeout waiting for execution" error:
- **Most transactions still succeed** on-chain before timeout
- Look for `[Step Name] Success! TX: 0x...` lines before timeout message
- Copy transaction hashes and verify on block explorer
- **Public RPC latency varies** - free RPCs can be slow during peak times

**What to do:**
1. Check which transactions succeeded (look for TX hashes in output)
2. Verify on Etherscan/Snowtrace that transactions confirmed
3. **For reliable execution:** Use your own RPC endpoint (Alchemy, Infura, etc.)

**Note:** Steps 2.7-2.8 now use Foundry scripts (no timeout issues!)

---

## Using Your Own RPC Endpoint (Recommended for Phase 2)

**To avoid timeout issues and speed up deployment:**

**Get free RPC from:**
- [Alchemy](https://www.alchemy.com/) - Free tier available
- [Infura](https://www.infura.io/) - Free tier available
- [QuickNode](https://www.quicknode.com/) - Free trial

**Update Step 1.2:**
```bash
export SEPOLIA_RPC=https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY
export FUJI_RPC=https://avalanche-fuji.infura.io/v3/YOUR_API_KEY
```

**Benefits:**
- ✅ Faster transaction confirmations
- ✅ More reliable (no timeout issues)
- ✅ Better for production use

---

## Deployment Fails
- Check RPC URL is accessible
- Verify wallet has testnet ETH
- Ensure OpenZeppelin installed: `bun add @openzeppelin/contracts`
- Ensure CCIP contracts installed: `forge install smartcontractkit/ccip@ccip-develop`

---

## "OnlyMinter" Error
```bash
# Check role granted
cast call $STABLECOIN "isMinter(address)(bool)" $CONSUMER --rpc-url $RPC
```
The command should return `true`.

---

## CCIP Setup Fails
- Verify all 4 addresses in payload are correct
- Ensure pool was deployed successfully
- Check wallet is token owner

---

## CCIP Transfer Not Arriving
- Wait full 10-20 minutes
- Check CCIP Explorer link
- Verify consumer has LINK balance
- Ensure sender approved consumer to spend tokens

---

## ACE Not Blocking Blacklisted Addresses

**Symptom:** Workflow succeeds even for blacklisted addresses

**🔍 Debugging Steps:**

### 1. Verify Function Selector

```bash
# The correct selector for onReport(bytes,bytes) is:
cast sig "onReport(bytes,bytes)"
# Should output: 0x805f2132

# Check what's configured:
export ON_REPORT_SELECTOR=$(cast sig "onReport(bytes,bytes)")
cast call $POLICY_ENGINE \
  "getExtractor(bytes4)(address)" \
  $ON_REPORT_SELECTOR \
  --rpc-url $SEPOLIA_RPC
# Should return your extractor address
```

**If extractor is NOT set or returns 0x00..., the selector is wrong!**

### 2. Verify Policies Are Attached

```bash
cast call $POLICY_ENGINE \
  "getPolicies(address,bytes4)(address[])" \
  $MINTING_CONSUMER_ACE \
  $ON_REPORT_SELECTOR \
  --rpc-url $SEPOLIA_RPC
# Should return: [0x971297f96F14884FD5D450B724bbe386ec37616F] (or your policy address)
```

**If this returns `[]`, policies are NOT attached! Re-run:**
```bash
ETHERSCAN_API_KEY=dummy forge script script/ConfigureACEWithConstants.s.sol:ConfigureACEWithConstants \
  --rpc-url $SEPOLIA_RPC \
  --broadcast \
  --legacy
```

### 3. Verify Parameter Encoding

ACE requires `keccak256("beneficiary")` NOT `bytes32("beneficiary")`!

```bash
# Correct: keccak256("beneficiary")
cast keccak "beneficiary"
# Output: 0xd49a1fdc7abf32173a5edf07d59d27b9172ad056cb5c008dc2bbefe702123a55

# ❌ WRONG: bytes32("beneficiary") - DO NOT USE THIS!
```

**Always use the configuration script** which correctly handles keccak256.

### 4. Verify Address Is Actually Blacklisted

```bash
cast call $BLACKLIST_POLICY \
  "isBlacklisted(address)(bool)" \
  <address_to_check> \
  --rpc-url $SEPOLIA_RPC
# Should return: true for blacklisted addresses
```

### 5. Understanding CRE Forwarder Behavior

**⚠️ CRITICAL INSIGHT:** CRE workflow will report `"reportDelivered": true` even when ACE blocks the transaction!

**Why?**
```
User → CRE Workflow → Forwarder Contract (TX succeeds ✅) 
                          ↓
                       Consumer.onReport() (can revert internally ❌)
```

The Forwarder transaction succeeds even if the consumer's internal logic reverts!

**What "reportDelivered" Means:**
- ✅ Report was successfully delivered to the consumer contract
- ❌ Does NOT mean the consumer's logic succeeded
- ⚠️  Consumer can internally revert (ACE blocking, insufficient balance, etc.)

**How to Verify ACE Actually Blocked:**

Check for Transfer events, NOT just TX status:
```bash
TX_HASH=<your_mint_transaction_hash>

# Get Transfer events
TRANSFER_SIG=$(cast sig-event "Transfer(address,address,uint256)")
cast logs --from-block <block> --to-block <block> \
  --address $STABLECOIN_SEPOLIA \
  $TRANSFER_SIG | grep $TX_HASH

# If NO output = ACE BLOCKED! ✅
# If Transfer event found = ACE did not block ❌
```

Or check balance:
```bash
# Before workflow
BAL_BEFORE=$(cast call $STABLECOIN_SEPOLIA "balanceOf(address)(uint256)" <beneficiary> --rpc-url $SEPOLIA_RPC)

# After workflow (with blacklisted address)
BAL_AFTER=$(cast call $STABLECOIN_SEPOLIA "balanceOf(address)(uint256)" <beneficiary> --rpc-url $SEPOLIA_RPC)

# If BAL_BEFORE == BAL_AFTER: ACE BLOCKED! ✅
```

### 6. Direct Test (Bypass CRE)

Test ACE directly to isolate the issue:
```bash
# Create test report
REPORT=$(cast abi-encode "f(uint8,address,uint256,bytes32)" 1 <blacklisted_addr> 50000000000000000000 0x5445535400000000000000000000000000000000000000000000000000000000)

# Call consumer directly
cast send $MINTING_CONSUMER_ACE \
  "onReport(bytes,bytes)" \
  0x \
  $REPORT \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY

# Expected if ACE works: 
# Error: execution reverted: PolicyRunRejected("address is blacklisted")
```

### 7. Verify ACE Blocking

```bash
# Check beneficiary balance (should be 0 if ACE blocked)
cast call $STABLECOIN_SEPOLIA \
  "balanceOf(address)(uint256)" \
  <blacklisted_address> \
  --rpc-url $SEPOLIA_RPC
# Should return: 0 (ACE blocked the mint)
```

---

## Need More Help?

- Check the [DEPLOYMENT_GUIDE.md](./DEPLOYMENT_GUIDE.md) for full deployment steps
- Review the [README.md](./README.md) for architecture overview
- Visit [Chainlink Documentation](https://docs.chain.link)

