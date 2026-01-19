// scripts/getAssetUid.ts
import { parseAbi } from 'viem'
import { createSepoliaClients, loadTokenizedAssetPlatformAddr } from './viemUtils'

// ABI - Include the getter for the public mapping 'assets'
// Solidity automatically generates a getter for public mappings
const abi = parseAbi([
  'function assets(uint256) view returns (string name, string symbol, string assetType, address issuer, uint256 totalSupply, bool active, string uid)',
])

async function main() {
  // Get viem clients and account (only publicClient is needed for reading data)
  const { publicClient, account } = createSepoliaClients()

  // Use the contract address defined in utils
  const CONTRACT_ADDRESS = loadTokenizedAssetPlatformAddr()

  // Safety check for contract address
  if (
    !CONTRACT_ADDRESS ||
    CONTRACT_ADDRESS.trim() === '' ||
    !CONTRACT_ADDRESS.startsWith('0x') ||
    CONTRACT_ADDRESS.length !== 42
  ) {
    console.error('\nERROR: Tokenized Asset Platform address is not set or invalid!')
    console.error('Please update the assetAddress config.json under workflow directory.')
    console.error('Script will exit now.\n')
    process.exit(1)
  }

  // ================================================
  // The assetId you want to query (change this value as needed)
  // ================================================
  const assetId = 1n

  console.log('Preparing to read Asset data from contract...')
  console.log('Contract address      :', CONTRACT_ADDRESS)
  console.log('Query assetId          :', assetId.toString())
  console.log('Query account          :', account.address)

  try {
    // Call the public getter function assets(assetId)
    const assetData = await publicClient.readContract({
      address: CONTRACT_ADDRESS,
      abi,
      functionName: 'assets',
      args: [assetId],
    })

    // assetData is returned as a tuple in the order of the struct fields
    const [name, symbol, assetType, issuer, totalSupply, active, uid] = assetData

    console.log('\nRead successful! Asset information:')
    console.log('----------------------------------------')
    console.log('Asset ID            :', assetId.toString())
    console.log('Name                :', name)
    console.log('Symbol              :', symbol)
    console.log('Asset Type          :', assetType)
    console.log('Issuer              :', issuer)
    console.log('Total Supply        :', totalSupply.toString())
    console.log('Active (valid)      :', active)
    console.log('UID (assetUid)      :', uid || '(empty)')
    console.log('----------------------------------------')

    // If you only need the UID:
    console.log('\nExtracted UID:', uid)
  } catch (error: any) {
    console.error('\nRead failed:', error)

    if (error.shortMessage) {
      console.error('Error summary:', error.shortMessage)
    }

    // Common error hints
    if (error.message?.includes('execution reverted')) {
      console.error('Possible reason: This assetId does not exist, or the contract is not properly deployed')
    }
  }
}

main().catch((err) => {
  console.error('Script execution failed:', err)
  process.exit(1)
})