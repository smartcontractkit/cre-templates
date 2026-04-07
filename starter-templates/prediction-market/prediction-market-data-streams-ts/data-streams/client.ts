import { ConfidentialHTTPClient, ok, type Runtime } from "@chainlink/cre-sdk"
import { decodeAbiParameters } from "viem"

// ─── Types ──────────────────────────────────────────────────

/** Top-level response from Data Streams REST API */
export interface DataStreamsApiResponse {
  report: {
    feedID: string
    validFromTimestamp: number
    observationsTimestamp: number
    fullReport: string
  }
}

/** Decoded fields from a v3 (Crypto Advanced) report body */
export interface DecodedReportV3 {
  feedId: `0x${string}`
  validFromTimestamp: number
  observationsTimestamp: number
  nativeFee: bigint
  linkFee: bigint
  expiresAt: number
  price: bigint    // int192, 18 decimals for crypto feeds
  bid: bigint      // int192, 18 decimals
  ask: bigint      // int192, 18 decimals
}

// ─── Pure-JS SHA-256 Implementation ─────────────────────────
// Needed because CRE WASM runtime does not have Node.js crypto.

const K: number[] = [
  0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
  0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
  0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
  0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
  0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
  0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
  0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
  0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2,
]

function rotr(n: number, x: number): number { return (x >>> n) | (x << (32 - n)) }

function sha256(message: Uint8Array): Uint8Array {
  let h0 = 0x6a09e667, h1 = 0xbb67ae85, h2 = 0x3c6ef372, h3 = 0xa54ff53a
  let h4 = 0x510e527f, h5 = 0x9b05688c, h6 = 0x1f83d9ab, h7 = 0x5be0cd19

  const msgLen = message.length
  const bitLen = msgLen * 8

  // Pre-processing: pad to 512-bit blocks
  const padLen = ((msgLen + 8) >> 6 << 6) + 64
  const padded = new Uint8Array(padLen)
  padded.set(message)
  padded[msgLen] = 0x80

  // Append length as big-endian 64-bit
  const view = new DataView(padded.buffer)
  view.setUint32(padLen - 4, bitLen, false)

  const W = new Int32Array(64)

  for (let offset = 0; offset < padLen; offset += 64) {
    for (let i = 0; i < 16; i++) {
      W[i] = view.getInt32(offset + i * 4, false)
    }
    for (let i = 16; i < 64; i++) {
      const s0 = rotr(7, W[i - 15] >>> 0) ^ rotr(18, W[i - 15] >>> 0) ^ (W[i - 15] >>> 3)
      const s1 = rotr(17, W[i - 2] >>> 0) ^ rotr(19, W[i - 2] >>> 0) ^ (W[i - 2] >>> 10)
      W[i] = (W[i - 16] + s0 + W[i - 7] + s1) | 0
    }

    let a = h0, b = h1, c = h2, d = h3, e = h4, f = h5, g = h6, h = h7

    for (let i = 0; i < 64; i++) {
      const S1 = rotr(6, e >>> 0) ^ rotr(11, e >>> 0) ^ rotr(25, e >>> 0)
      const ch = (e & f) ^ (~e & g)
      const temp1 = (h + S1 + ch + K[i] + W[i]) | 0
      const S0 = rotr(2, a >>> 0) ^ rotr(13, a >>> 0) ^ rotr(22, a >>> 0)
      const maj = (a & b) ^ (a & c) ^ (b & c)
      const temp2 = (S0 + maj) | 0

      h = g; g = f; f = e; e = (d + temp1) | 0
      d = c; c = b; b = a; a = (temp1 + temp2) | 0
    }

    h0 = (h0 + a) | 0; h1 = (h1 + b) | 0; h2 = (h2 + c) | 0; h3 = (h3 + d) | 0
    h4 = (h4 + e) | 0; h5 = (h5 + f) | 0; h6 = (h6 + g) | 0; h7 = (h7 + h) | 0
  }

  const result = new Uint8Array(32)
  const rv = new DataView(result.buffer)
  rv.setUint32(0, h0, false); rv.setUint32(4, h1, false)
  rv.setUint32(8, h2, false); rv.setUint32(12, h3, false)
  rv.setUint32(16, h4, false); rv.setUint32(20, h5, false)
  rv.setUint32(24, h6, false); rv.setUint32(28, h7, false)
  return result
}

function hmacSha256(key: Uint8Array, message: Uint8Array): Uint8Array {
  const BLOCK_SIZE = 64

  // If key is longer than block size, hash it
  let keyBlock = key.length > BLOCK_SIZE ? sha256(key) : key

  // Pad key to block size
  const paddedKey = new Uint8Array(BLOCK_SIZE)
  paddedKey.set(keyBlock)

  const oKeyPad = new Uint8Array(BLOCK_SIZE)
  const iKeyPad = new Uint8Array(BLOCK_SIZE)
  for (let i = 0; i < BLOCK_SIZE; i++) {
    oKeyPad[i] = paddedKey[i] ^ 0x5c
    iKeyPad[i] = paddedKey[i] ^ 0x36
  }

  // Inner hash: SHA-256(iKeyPad || message)
  const inner = new Uint8Array(BLOCK_SIZE + message.length)
  inner.set(iKeyPad)
  inner.set(message, BLOCK_SIZE)
  const innerHash = sha256(inner)

  // Outer hash: SHA-256(oKeyPad || innerHash)
  const outer = new Uint8Array(BLOCK_SIZE + 32)
  outer.set(oKeyPad)
  outer.set(innerHash, BLOCK_SIZE)
  return sha256(outer)
}

function toHex(bytes: Uint8Array): string {
  let hex = ""
  for (let i = 0; i < bytes.length; i++) {
    hex += bytes[i].toString(16).padStart(2, "0")
  }
  return hex
}

function textToBytes(text: string): Uint8Array {
  const bytes = new Uint8Array(text.length)
  for (let i = 0; i < text.length; i++) {
    bytes[i] = text.charCodeAt(i)
  }
  return bytes
}

// ─── HMAC Authentication ────────────────────────────────────

function generateHMAC(
  method: string,
  path: string,
  body: string,
  apiKey: string,
  apiSecret: string,
  timestamp: number,
): string {
  const bodyHash = toHex(sha256(textToBytes(body || "")))
  const stringToSign = `${method} ${path} ${bodyHash} ${apiKey} ${timestamp}`
  return toHex(hmacSha256(textToBytes(apiSecret), textToBytes(stringToSign)))
}

// ─── Report Decoding ────────────────────────────────────────

/**
 * Decode a v3 (Crypto Advanced) fullReport blob into its constituent fields.
 *
 * The fullReport is ABI-encoded as:
 *   (bytes32[3] reportContext, bytes reportData, bytes32[] rawRs, bytes32[] rawSs, bytes32 rawVs)
 *
 * The reportData inside is ABI-encoded as:
 *   (bytes32 feedId, uint32 validFromTimestamp, uint32 observationsTimestamp,
 *    uint192 nativeFee, uint192 linkFee, uint32 expiresAt,
 *    int192 price, int192 bid, int192 ask)
 */
export function decodeFullReportV3(fullReportHex: `0x${string}`): DecodedReportV3 {
  // Step 1: Decode the outer wrapper to extract reportData
  const [, reportData] = decodeAbiParameters(
    [
      { type: "bytes32[3]", name: "reportContext" },
      { type: "bytes", name: "reportData" },
      { type: "bytes32[]", name: "rawRs" },
      { type: "bytes32[]", name: "rawSs" },
      { type: "bytes32", name: "rawVs" },
    ],
    fullReportHex,
  )

  // Step 2: Decode the report body (v3 Crypto Advanced schema)
  const [feedId, validFromTimestamp, observationsTimestamp, nativeFee, linkFee, expiresAt, price, bid, ask] =
    decodeAbiParameters(
      [
        { type: "bytes32", name: "feedId" },
        { type: "uint32", name: "validFromTimestamp" },
        { type: "uint32", name: "observationsTimestamp" },
        { type: "uint192", name: "nativeFee" },
        { type: "uint192", name: "linkFee" },
        { type: "uint32", name: "expiresAt" },
        { type: "int192", name: "price" },
        { type: "int192", name: "bid" },
        { type: "int192", name: "ask" },
      ],
      reportData as `0x${string}`,
    )

  return {
    feedId: feedId as `0x${string}`,
    validFromTimestamp,
    observationsTimestamp,
    nativeFee,
    linkFee,
    expiresAt,
    price,
    bid,
    ask,
  }
}

// ─── Fetch Latest Report ────────────────────────────────────

/**
 * Fetch the latest Data Streams report for a given feed ID.
 *
 * Uses the CRE ConfidentialHTTPClient which executes in a secure enclave
 * and supports the Authorization header. The HMAC signature and timestamp
 * are pre-computed and passed as templatePublicValues so all nodes produce
 * the same request. The API key is injected via vaultDonSecrets.
 */
export function fetchLatestReport(
  runtime: Runtime<any>,
  apiUrl: string,
  feedId: string,
): DataStreamsApiResponse {
  const method = "GET"
  const path = `/api/v1/reports/latest?feedID=${feedId}`
  const timestamp = Math.floor(runtime.now().getTime())

  // Read secrets to compute HMAC (needed for the signature only)
  const apiKey = runtime.getSecret({ id: "DATA_STREAMS_API_KEY" }).result().value
  const apiSecret = runtime.getSecret({ id: "DATA_STREAMS_API_SECRET" }).result().value
  const signature = generateHMAC(method, path, "", apiKey, apiSecret, timestamp)

  const confHTTPClient = new ConfidentialHTTPClient()
  const response = confHTTPClient.sendRequest(runtime, {
    request: {
      url: `${apiUrl}${path}`,
      method,
      multiHeaders: {
        "Authorization": { values: ["{{.DATA_STREAMS_API_KEY}}"] },
        "X-Authorization-Timestamp": { values: [String(timestamp)] },
        "X-Authorization-Signature-SHA256": { values: [signature] },
      },
    },
    vaultDonSecrets: [{ key: "DATA_STREAMS_API_KEY" }],
  }).result()

  if (!ok(response)) {
    throw new Error(`Data Streams API error (${response.statusCode}): ${Buffer.from(response.body).toString("utf-8")}`)
  }

  const body = Buffer.from(response.body).toString("utf-8")
  return JSON.parse(body) as DataStreamsApiResponse
}

// ─── Price Extraction Helper ────────────────────────────────

/**
 * Data Streams crypto feeds use 18 decimal places.
 * On-chain Chainlink Data Feeds use 8 decimal places for USD pairs.
 * This helper converts from 18 to 8 decimals to match the PredictionMarket
 * contract's strike price format.
 */
export const DATA_STREAMS_DECIMALS = 18
export const PRICE_FEED_DECIMALS = 8
const SCALE_FACTOR = BigInt(10 ** (DATA_STREAMS_DECIMALS - PRICE_FEED_DECIMALS))

export function toPriceFeedDecimals(priceWith18Decimals: bigint): bigint {
  return priceWith18Decimals / SCALE_FACTOR
}
