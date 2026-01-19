import { parseAbi } from 'viem'
import { createSepoliaClients, loadTokenizedAssetPlatformAddr } from './viemUtils'

const { walletClient, publicClient, account } = createSepoliaClients()


// ABI - Only including the needed function (you can also load full ABI from file)
const abi = parseAbi([
  'function redeem(uint256 assetId, uint256 amount, string settlementDetail) public',
])

// ================================================
// Function call parameters
// ================================================
const redeemParams = {
  assetId: 1n,
  amount: 10n,
  settlementDetail: "",
} as const

async function main() {
    const CONTRACT_ADDRESS = loadTokenizedAssetPlatformAddr()
    if (!CONTRACT_ADDRESS || 
        CONTRACT_ADDRESS.trim() === '' || !CONTRACT_ADDRESS.startsWith('0x') || 
        CONTRACT_ADDRESS.length !== 42
    ) {
        console.error('\nERROR: Tokenized Asset Platform address is not set or invalid!')
        console.error('Please update the assetAddress config.json under workflow directory.')
        console.error('Script will exit now.\n')
        process.exit(1)
    }

    console.log('Preparing to call verifyAsset function...')
    console.log('Contract address:', CONTRACT_ADDRESS)
    console.log('Caller:', account.address)
    console.log('Parameters:', redeemParams)

  try {
    // 1. Send transaction
    const hash = await walletClient.writeContract({
      address: CONTRACT_ADDRESS,
      abi,
      functionName: 'redeem',
      args: [
        redeemParams.assetId,
        redeemParams.amount,
        redeemParams.settlementDetail,
      ],
    })

    console.log('Redeem Transaction sent! Hash:', hash)

    // 2. Wait for confirmation
    const receipt = await publicClient.waitForTransactionReceipt({
      hash,
      confirmations: 1,
    })

    console.log('Transaction successful!')
    console.log('Block number:', receipt.blockNumber.toString())
    console.log('Gas used:', receipt.gasUsed.toString())
  } catch (error: any) {
    console.error('Call failed:', error)
    
    if (error.shortMessage) {
      console.error('Error message:', error.shortMessage)
    }
    if (error.cause) {
      console.error('Detailed reason:', error.cause.message)
    }
  }
}

main().catch(console.error)