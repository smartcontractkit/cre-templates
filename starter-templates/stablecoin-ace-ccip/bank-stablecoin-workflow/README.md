# Phase 1: Basic Stablecoin Deployment

Deploy a stablecoin with mint and redeem functionality using Chainlink Runtime Environment (CRE).

**What you'll build:**
- ERC20 stablecoin contract
- CRE consumer for mint/redeem operations
- HTTP-triggered workflow for bank integration

**Time required:** ~20 minutes

---

## Prerequisites

### Required Tools
- Foundry (`forge --version`)
- CRE CLI (`cre --version`)
- Bun v1.2.21+ (`bun --version`)

### Required Funds
- Sepolia testnet ETH ([faucet](https://faucets.chain.link))

---

## Step 1.1: Initial Setup

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

## Step 1.2: Install Contract Dependencies

```bash
bun add @openzeppelin/contracts
```

## Step 1.3: Set Environment Variables

```bash
# Load all environment variables from .env (includes PRIVATE_KEY, RPCs, etc.)
source .env

# Export private key for Foundry commands
export PRIVATE_KEY=$CRE_ETH_PRIVATE_KEY

# Set project root for payload file paths
export CRE_PROJECT_ROOT=$(pwd)
```

**Note:** `.env` includes `SEPOLIA_RPC` and `FUJI_RPC`. If experiencing rate limits, update these in `.env` with your own Alchemy/Infura endpoints.

## Step 1.4: Deploy StablecoinERC20 on Sepolia

```bash
forge create contracts/StablecoinERC20.sol:StablecoinERC20 \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY --broadcast \
  --constructor-args "BankStablecoin" "creUSD"
```

Copy the `Deployed to:` address from output:

```bash
export STABLECOIN_SEPOLIA=<paste_deployed_address_here>
```

## Step 1.5: Deploy MintingConsumer on Sepolia

**Constructor arguments explained:**
- `$STABLECOIN_SEPOLIA` = Your deployed stablecoin address
- `0x0000000000000000000000000000000000000000` = Expected author (use address(0) for testing)
- `0x64756d6d790000000000` = Expected workflow name (we're using bytes10("dummy") for testing)

```bash
forge create contracts/MintingConsumer.sol:MintingConsumer \
  --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY --broadcast \
  --constructor-args $STABLECOIN_SEPOLIA 0x0000000000000000000000000000000000000000 0x64756d6d790000000000
```

Copy the `Deployed to:` address from output:

```bash
export MINTING_CONSUMER_SEPOLIA=<paste_deployed_address_here>
```

## Step 1.6: Grant Roles to MintingConsumer

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

## Step 1.7: Update Workflow Configuration

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

## Step 1.8: Install Workflow Dependencies

```bash
cd bank-stablecoin-workflow && bun install && cd ..
```

## Step 1.9: Test MINT

```bash
cre workflow simulate bank-stablecoin-workflow \
  --target local-simulation \
  --broadcast \
  --trigger-index 0 \
  --non-interactive \
  --http-payload @$CRE_PROJECT_ROOT/bank-stablecoin-workflow/http_trigger_payload.json
```

**Check the transaction hash in the output** to verify the mint was submitted on-chain.

## Step 1.10: Test REDEEM

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

## Phase 1 Complete

**Next:** [Phase 2: Add Cross-Chain CCIP Transfers](../ccip-transfer-workflow/README.md)

**See also:** [Complete Deployment Guide](../DEPLOYMENT_GUIDE.md) | [Main README](../README.md)

