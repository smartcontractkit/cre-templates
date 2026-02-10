import { 
  cre, 
  Runner, 
  type Runtime,
  getNetwork,
  type HTTPPayload,
  HTTPSendRequester,
  ok,
  consensusIdenticalAggregation,
  EVMLog,
  hexToBase64
 } from "@chainlink/cre-sdk";
import { bytesToHex, decodeEventLog, encodeAbiParameters, parseAbi, parseAbiParameters } from "viem";
import { z } from 'zod';

const configSchema = z.object({
  	schedule: z.string(),
    url: z.string(),
	  evms: z.array(
		z.object({
			assetAddress: z.string(),
			chainSelectorName: z.string(),
			gasLimit: z.string(),
		}),
	),
})

type Config = z.infer<typeof configSchema>

type PostResponse = {
  statusCode: number
}

// Multiple events CRE can process are: AssetRegistered, AssetVerified, TokensMinted and TokensRedeemed
// Use string for bigint because bigint cannot be serialized
type AssetRegisterParams = {
  action: "AssetRegistered"
  assetId: string
  issuer: string
  initialSupply: string
  assetName: string
}

type AssetVerifiedParams = {
  action: "AssetVerified"
  assetId: string
  isValid: boolean
}

type TokensMintedParams = {
  action: "TokensMinted"
  assetId: string
  amount: string
}

type TokensRedeemedParams = {
  action: "TokensRedeemed"
  assetId: string
  amount: string
}

type AssetParams = AssetRegisterParams | AssetVerifiedParams | TokensMintedParams | TokensRedeemedParams;


// function to be used to send multiple request to AWS lambda function
const postData = (sendRequester: HTTPSendRequester, config: Config, assetParams: AssetParams): PostResponse => {
  // 1. Prepare the payload for POST request
  let dataToSend: any = { ...assetParams };

  // 2. Serialize the data to JSON and encode as bytes
  const bodyBytes = new TextEncoder().encode(JSON.stringify(dataToSend))

  // 3. Convert to base64 for the request
  const body = Buffer.from(bodyBytes).toString("base64")

  // 4. Construct the POST request with cacheSettings
  const req = {
    url: config.url,
    method: "POST" as const,
    body,
    headers: {
      "Content-Type": "application/json",
    },
    cacheSettings: {
      readFromCache: true, // Enable reading from cache
      maxAgeMs: 60000, // Accept cached responses up to 60 seconds old
    },
  }

  // 5. Send the request and wait for the response
  const resp = sendRequester.sendRequest(req).result()
  if (!ok(resp)) {
    throw new Error(`HTTP request failed with status: ${resp.statusCode}`)
  }
  return { statusCode: resp.statusCode }
}


const eventAbi = parseAbi([
  "event AssetRegistered(uint256 indexed assetId, address indexed issuer, uint256 indexed initialSupply, string name, string symbol, string assetType)",
  "event AssetVerified(uint256 indexed assetId, bool indexed isValid, string verificationDetails)",
  "event TokensMinted(uint256 indexed assetId, uint256 indexed amount, address indexed to, string reason)",
  "event TokensRedeemed(uint256 indexed assetId, uint256 indexed amount, address indexed account, string settlementDetails)",
])

const onLogTrigger = (runtime: Runtime<Config>, log: EVMLog): string => {
  
  const topics = log.topics.map((topic) => bytesToHex(topic)) as [`0x${string}`, ...`0x${string}`[]]
  const data = bytesToHex(log.data)

  // Decode the event
  const decodedLog = decodeEventLog({
    abi: eventAbi,
    data,
    topics
  })

  runtime.log(`Event name: ${decodedLog.eventName}`)
  let assetParams: AssetParams
  
  const httpClient = new cre.capabilities.HTTPClient()

  // extract info from event and prepare the parameters for function postData
  switch(decodedLog.eventName) {
    case "AssetRegistered":
      const { assetId: assetIdReg, issuer, initialSupply, name } = decodedLog.args
      assetParams = {
        action: "AssetRegistered",
        assetId: assetIdReg.toString(),
        issuer,
        initialSupply: initialSupply.toString(),
        assetName: name,
      };    
      runtime.log(`Event AssetRegistered detected: assetId ${assetIdReg} | issuer ${issuer} initialSupply ${initialSupply} | name ${name}`)
      break;
    case "AssetVerified":
      const { assetId: assetIdVer, isValid } = decodedLog.args
      assetParams = {
        action: "AssetVerified",
        assetId: assetIdVer.toString(),
        isValid
      }
      runtime.log(`Event AssetVerified detected: assetId ${assetIdVer} | isValid ${isValid}`)
      break
    case "TokensMinted":
      const { assetId: assetIdMint, amount: amountMint } = decodedLog.args
      assetParams = {
        action: "TokensMinted",
        assetId: assetIdMint.toString(),
        amount: amountMint.toString()
      }
      runtime.log(`Event TokensMinted detected: assetId ${assetIdMint} | amount ${amountMint}`)
      break
    case "TokensRedeemed":
      const {assetId: assetIdRedeem, amount: amountRedeem } = decodedLog.args
      assetParams = {
        action: "TokensRedeemed",
        assetId: assetIdRedeem.toString(),
        amount: amountRedeem.toString()
      }
      runtime.log(`Event TokensRedeemed detected: assetId ${assetIdRedeem} | amount ${amountRedeem}`)
      break
    default:
      return "No key event detected"
  }

  const result = httpClient.sendRequest(
    runtime,
    postData,
    consensusIdenticalAggregation<PostResponse>()
  )(
    runtime.config,
    assetParams
  ).result()

  runtime.log(`Successfully sent data to url. Status ${result.statusCode}`)
  return "Success"
};

const onHTTPTrigger = (runtime: Runtime<Config>, payload: HTTPPayload): string => {
	runtime.log('Raw HTTP trigger received')

	// Expect a HTTP request with metdata info, and it cannot be empty
	if (!payload.input || payload.input.length === 0) {
		runtime.log('HTTP trigger payload is empty')
    throw new Error("Json payload is empty")
	}

	// Log the raw JSON for debugging (human-readable).
	runtime.log(`Payload bytes payloadBytes ${payload.input.toString()}`)

	try {
    // fetch the assetId and newUri from http request
		runtime.log(`Parsed HTTP trigger received payload ${payload.input.toString()}`)
    const responseText = Buffer.from(payload.input).toString('utf-8')
    const {assetId, uid} = JSON.parse(responseText)

    runtime.log(`Asset ID is ${assetId}`)
    runtime.log(`Asset UID is ${uid}`)

    if(!assetId || !uid) {
      throw new Error("Failed to extract assetId or newUri from Http request, please check the json payload file")
    }

    // init an evmClient to send transaction 
    const evmConfig = runtime.config.evms[0]
    runtime.log(`Updating metadata for Asset State contract, address is: ${evmConfig.assetAddress}`)

    const network = getNetwork({
      chainFamily: "evm",
      chainSelectorName: evmConfig.chainSelectorName,
      isTestnet: true
    })

    if(!network) {
      throw new Error("Failed to get network config")
    }
    const evmClient = new cre.capabilities.EVMClient(network?.chainSelector.selector)

    const reportData = encodeAbiParameters(
    parseAbiParameters("uint256 assetId, string memory newUri"),
    [BigInt(assetId), uid as string]
  )

    // generate signed report
    const reportResponse = runtime.report({
      encodedPayload: hexToBase64(reportData),
			encoderName: 'evm',
      signingAlgo: 'ecdsa',
      hashingAlgo: 'keccak256',
    })
    .result()

    // submit report to the tokenized asset platform contract 
    const writeReportResult = evmClient
		.writeReport(runtime, {
			receiver: evmConfig.assetAddress,
			report: reportResponse,
			gasConfig: {
				gasLimit: evmConfig.gasLimit,
			},
		})
		.result()

    const txHash = bytesToHex(writeReportResult.txHash || new Uint8Array(32))

    runtime.log(`write report transaction succeeded: ${txHash}`)

    return txHash

	} catch (error) {
		runtime.log('Failed to parse HTTP trigger payload')
		throw new Error('Failed to parse HTTP trigger payload')
	}
}

const initWorkflow = (config: Config) => {
  const network = getNetwork({
		chainFamily: 'evm',
		chainSelectorName: config.evms[0].chainSelectorName,
		isTestnet: true,
	})

  if (!network) {
		throw new Error(
			`Network not found for chain selector name: ${config.evms[0].chainSelectorName}`,
		)
	}

  const evmClient = new cre.capabilities.EVMClient(network.chainSelector.selector)
  const httpTrigger = new cre.capabilities.HTTPCapability()

  return [
    cre.handler(
      evmClient.logTrigger({
        addresses: [hexToBase64(config.evms[0].assetAddress)],
      }),
      onLogTrigger,
    ),
    cre.handler(httpTrigger.trigger({}), onHTTPTrigger),
  ]
};

export async function main() {
  const runner = await Runner.newRunner<Config>();
  await runner.run(initWorkflow);
}

main();
