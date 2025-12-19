import { createPublicClient, http, parseAbi } from 'viem';
import { sepolia } from 'viem/chains';
import config from './config.json' with {type: 'json'}

async function main() {

  const RPC_URL = 'https://por.bcy-p.metalhosts.com/cre-alpha/MvqtrdftrbxcP3ZgGBJb3bK5/ethereum/sepolia'
  const contractAddress = config.evms[0].assetAddress;
  const assetId = 1n
  
  const publicClient = createPublicClient({
    chain: sepolia, 
    transport: http(RPC_URL), 
  });

  const abi = parseAbi([
    'function uid(uint256 assetId) public view returns (string memory)'
  ]);

  try {
    const storedData = await publicClient.readContract({
      address: contractAddress,
      abi,
      functionName: 'uid',
      args: [assetId],
    });

    console.log('value of the assetId 1 uid:', storedData.toString());
  } catch (error) {
    console.error('failed to read the value:', error);
  }
}

main().catch(console.error);