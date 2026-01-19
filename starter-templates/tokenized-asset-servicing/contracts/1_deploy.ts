import * as fs from 'fs'
import { join } from 'path'
import {Abi} from 'viem'
import { createSepoliaClients } from './viemUtils'

const { walletClient, publicClient } = createSepoliaClients()

// ================================================
// get contract's abi and bytecode
// abi and bin have already been generated to dir output
// run command to generate abi and bin if you need regenerate
// solc --abi --bin --optimize \
//     --base-path . \
//     --include-path node_modules/ \
//     --overwrite \
//     -o output/ \
//     TokenizedAssetPlatform.sol
// ================================================
const abiPath = join(process.cwd(), 'output', 'TokenizedAssetPlatform.abi')
if (!fs.existsSync(abiPath)) {
  throw new Error(`cannot find ABI file at ${abiPath}`)
}

let abi: Abi
try {
    const abiRaw = fs.readFileSync(abiPath, 'utf8').trim()
    abi = JSON.parse(abiRaw) as Abi
} catch (e) {
  throw new Error(`ABI file format error, not a valid JSON: ${abiPath}\n${e}`)
}

const bytecodePath = join(process.cwd(), 'output', `TokenizedAssetPlatform.bin`)
if (!fs.existsSync(bytecodePath)) {
  throw new Error(`cannot find bytecode file at ${bytecodePath}`)
}

let bytecodeRaw = fs.readFileSync(bytecodePath, 'utf8').trim()
if (!bytecodeRaw.startsWith('0x')) {
  bytecodeRaw = '0x' + bytecodeRaw
}

const bytecode = bytecodeRaw as `0x${string}`


async function main() {
  // find forwarder address for different chains at doc: 
  // https://docs.chain.link/cre/guides/workflow/using-evm-client/supported-networks-ts#simulation-testnets
  // Eth sepolia is used in this script
  const forwarderAddr = "0x15fC6ae953E024d975e77382eEeC56A9101f9F88";

  try {
    const hash = await walletClient.deployContract({
      abi,
      bytecode,
      args: [forwarderAddr],
    })

    console.log('deploy tx Hash:', hash)

    const receipt = await publicClient.waitForTransactionReceipt({ hash })
    
    console.log('Contract TokenizedAssetPlatform deployed successfully...')
    console.log('Contact deployed at:', receipt.contractAddress)
    console.log('Block number:', receipt.blockNumber)
  } catch (error) {
    console.error('Failed to deploy:', error)
    process.exit(1)
  }
}

main()