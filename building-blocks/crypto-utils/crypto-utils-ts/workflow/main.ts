import { CronCapability, handler, Runner, type Runtime } from "@chainlink/cre-sdk";

// ============================================================================
// Noble Hashes - Cryptographic Hash Functions
// Alternative to: crypto.createHash()
// ============================================================================
import { sha256, sha384, sha512 } from "@noble/hashes/sha2";
import { sha3_256, sha3_512, keccak_256 } from "@noble/hashes/sha3";
import { blake2b } from "@noble/hashes/blake2b";
import { blake3 } from "@noble/hashes/blake3";
import { ripemd160 } from "@noble/hashes/ripemd160";

// ============================================================================
// Noble Hashes - HMAC
// Alternative to: crypto.createHmac()
// ============================================================================
import { hmac } from "@noble/hashes/hmac";

// ============================================================================
// Noble Hashes - Key Derivation Functions
// Alternative to: crypto.pbkdf2(), crypto.scrypt(), crypto.hkdf()
// ============================================================================
import { pbkdf2 } from "@noble/hashes/pbkdf2";
import { scrypt } from "@noble/hashes/scrypt";
import { hkdf } from "@noble/hashes/hkdf";

// ============================================================================
// Noble Hashes - Utilities
// Alternative to: Buffer.from(hex, 'hex'), buffer.toString('hex')
// NOTE: randomBytes() is NOT available in QuickJS (no crypto.getRandomValues)
// ============================================================================
import { bytesToHex, hexToBytes, utf8ToBytes } from "@noble/hashes/utils";

// ============================================================================
// Noble Curves - Elliptic Curve Cryptography
// Alternative to: crypto.generateKeyPair(), crypto.sign(), crypto.verify()
// ============================================================================
import { secp256k1 } from "@noble/curves/secp256k1";
import { ed25519 } from "@noble/curves/ed25519";
import { x25519 } from "@noble/curves/ed25519";

// ============================================================================
// Noble Ciphers - Symmetric Encryption
// Alternative to: crypto.createCipheriv(), crypto.createDecipheriv()
// ============================================================================
import { gcm } from "@noble/ciphers/aes";
import { chacha20poly1305 } from "@noble/ciphers/chacha";

type Config = {
  schedule: string;
};

// ============================================================================
// Random Bytes Generation for CRE
//
// Noble's randomBytes() uses crypto.getRandomValues (Web Crypto API) which is
// NOT available in QuickJS. However, CRE provides a Math.random() polyfill:
//
// - Implementation: ChaCha8Rng (cryptographic PRNG) in the Javy WASM plugin
// - Seeding: Host runtime provides seed via random_seed() import
// - Consensus-safe: Same seed produces same sequence across all DON nodes
// - Mode isolation: Different execution modes get independent RNG streams
//
// Below we provide randomBytes() using CRE's Math.random() polyfill.
// ============================================================================

/**
 * Generate random bytes using CRE's consensus-safe Math.random() polyfill.
 *
 * CRE replaces Math.random() with a deterministic ChaCha8Rng PRNG seeded by
 * the host runtime. This ensures all DON nodes generate the same sequence
 * of random values for consensus.
 */
function randomBytes(length: number): Uint8Array {
  const bytes = new Uint8Array(length);
  for (let i = 0; i < length; i++) {
    bytes[i] = Math.floor(Math.random() * 256);
  }
  return bytes;
}

/**
 * Derive deterministic bytes from a seed (for reproducible demo output).
 * Used in address derivation examples so output matches documentation.
 */
function deriveBytes(seed: string, length: number): Uint8Array {
  const hash = sha512(utf8ToBytes(seed));
  if (length <= 64) {
    return hash.slice(0, length);
  }
  const result = new Uint8Array(length);
  let offset = 0;
  let counter = 0;
  while (offset < length) {
    const block = sha512(utf8ToBytes(`${seed}-${counter}`));
    const toCopy = Math.min(64, length - offset);
    result.set(block.slice(0, toCopy), offset);
    offset += toCopy;
    counter++;
  }
  return result;
}

// ============================================================================
// Demo Functions for Each Crypto Category
// ============================================================================

function demoHashing(runtime: Runtime<Config>): void {
  runtime.log("=== HASHING FUNCTIONS ===");

  const message = utf8ToBytes("Hello, CRE Workflow!");

  // SHA-2 Family (most common)
  const sha256Hash = sha256(message);
  runtime.log(`SHA-256: ${bytesToHex(sha256Hash)}`);

  const sha384Hash = sha384(message);
  runtime.log(`SHA-384: ${bytesToHex(sha384Hash)}`);

  const sha512Hash = sha512(message);
  runtime.log(`SHA-512: ${bytesToHex(sha512Hash)}`);

  // SHA-3 Family
  const sha3_256Hash = sha3_256(message);
  runtime.log(`SHA3-256: ${bytesToHex(sha3_256Hash)}`);

  const sha3_512Hash = sha3_512(message);
  runtime.log(`SHA3-512: ${bytesToHex(sha3_512Hash)}`);

  // Keccak-256 (used in Ethereum for addresses and tx hashes)
  const keccak256Hash = keccak_256(message);
  runtime.log(`Keccak-256 (Ethereum): ${bytesToHex(keccak256Hash)}`);

  // BLAKE2b (fast, secure, used in many blockchains)
  const blake2bHash = blake2b(message, { dkLen: 32 });
  runtime.log(`BLAKE2b-256: ${bytesToHex(blake2bHash)}`);

  // BLAKE3 (fastest secure hash)
  const blake3Hash = blake3(message);
  runtime.log(`BLAKE3: ${bytesToHex(blake3Hash)}`);

  // RIPEMD-160 (used in Bitcoin addresses)
  const ripemd160Hash = ripemd160(message);
  runtime.log(`RIPEMD-160: ${bytesToHex(ripemd160Hash)}`);

  // Incremental hashing (for large data)
  const incrementalHash = sha256.create()
    .update(utf8ToBytes("Hello, "))
    .update(utf8ToBytes("CRE Workflow!"))
    .digest();
  runtime.log(`SHA-256 (incremental): ${bytesToHex(incrementalHash)}`);
}

function demoHmac(runtime: Runtime<Config>): void {
  runtime.log("=== HMAC (Hash-based Message Authentication Code) ===");

  const key = utf8ToBytes("my-secret-key-for-hmac");
  const message = utf8ToBytes("Message to authenticate");

  // HMAC-SHA256 (most common)
  const hmacSha256 = hmac(sha256, key, message);
  runtime.log(`HMAC-SHA256: ${bytesToHex(hmacSha256)}`);

  // HMAC-SHA512
  const hmacSha512 = hmac(sha512, key, message);
  runtime.log(`HMAC-SHA512: ${bytesToHex(hmacSha512)}`);

  // HMAC-SHA3-256
  const hmacSha3 = hmac(sha3_256, key, message);
  runtime.log(`HMAC-SHA3-256: ${bytesToHex(hmacSha3)}`);
}

function demoKeyDerivation(runtime: Runtime<Config>): void {
  runtime.log("=== KEY DERIVATION FUNCTIONS ===");

  const password = utf8ToBytes("user-password");
  const salt = utf8ToBytes("random-salt-value");

  // PBKDF2 - Password-Based Key Derivation Function 2
  // Good for: Password hashing, deriving encryption keys from passwords
  // NOTE: Using lower iterations (1000) for demo speed. Use 100000+ in production.
  const pbkdf2Key = pbkdf2(sha256, password, salt, { c: 1000, dkLen: 32 });
  runtime.log(`PBKDF2-SHA256 (1k iterations): ${bytesToHex(pbkdf2Key)}`);

  // Scrypt - Memory-hard password hashing
  // Good for: Password storage (resistant to GPU/ASIC attacks)
  // NOTE: Using lower N (2^10) for demo speed. Use 2^14 or higher in production.
  const scryptKey = scrypt(password, salt, { N: 2 ** 10, r: 8, p: 1, dkLen: 32 });
  runtime.log(`Scrypt (N=1024, r=8, p=1): ${bytesToHex(scryptKey)}`);

  // HKDF - HMAC-based Key Derivation Function
  // Good for: Deriving multiple keys from a single secret
  const inputKeyMaterial = utf8ToBytes("shared-secret-from-ecdh");
  const info = utf8ToBytes("encryption-key");
  const hkdfKey = hkdf(sha256, inputKeyMaterial, salt, info, 32);
  runtime.log(`HKDF-SHA256: ${bytesToHex(hkdfKey)}`);
}

function demoSymmetricEncryption(runtime: Runtime<Config>): void {
  runtime.log("=== SYMMETRIC ENCRYPTION (AEAD) ===");
  runtime.log("Using randomBytes() with CRE's ChaCha8Rng-based Math.random() polyfill.");

  const plaintext = utf8ToBytes("Sensitive data to encrypt");

  // AES-256-GCM - Industry standard authenticated encryption
  const aesKey = randomBytes(32); // 256-bit key
  const aesNonce = randomBytes(12); // 96-bit nonce for GCM

  const aesGcm = gcm(aesKey, aesNonce);
  const aesCiphertext = aesGcm.encrypt(plaintext);
  const aesDecrypted = aesGcm.decrypt(aesCiphertext);

  runtime.log(`AES-256-GCM Key: ${bytesToHex(aesKey)}`);
  runtime.log(`AES-256-GCM Nonce: ${bytesToHex(aesNonce)}`);
  runtime.log(`AES-256-GCM Ciphertext: ${bytesToHex(aesCiphertext)}`);
  runtime.log(`AES-256-GCM Decrypted: ${new TextDecoder().decode(aesDecrypted)}`);

  // ChaCha20-Poly1305 - Modern alternative to AES-GCM
  // Faster in software, no timing attacks, used by TLS 1.3
  const chachaKey = randomBytes(32); // 256-bit key
  const chachaNonce = randomBytes(12); // 96-bit nonce

  const chacha = chacha20poly1305(chachaKey, chachaNonce);
  const chachaCiphertext = chacha.encrypt(plaintext);
  const chachaDecrypted = chacha.decrypt(chachaCiphertext);

  runtime.log(`ChaCha20-Poly1305 Key: ${bytesToHex(chachaKey)}`);
  runtime.log(`ChaCha20-Poly1305 Nonce: ${bytesToHex(chachaNonce)}`);
  runtime.log(`ChaCha20-Poly1305 Ciphertext: ${bytesToHex(chachaCiphertext)}`);
  runtime.log(`ChaCha20-Poly1305 Decrypted: ${new TextDecoder().decode(chachaDecrypted)}`);
}

function demoEcdsaSecp256k1(runtime: Runtime<Config>): void {
  runtime.log("=== ECDSA with secp256k1 (Bitcoin/Ethereum) ===");
  runtime.log("Using randomBytes() with CRE's ChaCha8Rng-based Math.random() polyfill.");

  // Generate random private key using CRE's consensus-safe PRNG
  const privateKey = randomBytes(32);
  runtime.log(`Private Key: ${bytesToHex(privateKey)}`);

  // Derive public key (compressed format - 33 bytes)
  const publicKey = secp256k1.getPublicKey(privateKey);
  runtime.log(`Public Key (compressed): ${bytesToHex(publicKey)}`);

  // Derive public key (uncompressed format - 65 bytes)
  const publicKeyUncompressed = secp256k1.getPublicKey(privateKey, false);
  runtime.log(`Public Key (uncompressed): ${bytesToHex(publicKeyUncompressed)}`);

  // Sign a message (message should be a 32-byte hash)
  const messageHash = sha256(utf8ToBytes("Transaction data to sign"));
  const signature = secp256k1.sign(messageHash, privateKey);

  runtime.log(`Signature (r): ${signature.r.toString(16)}`);
  runtime.log(`Signature (s): ${signature.s.toString(16)}`);
  runtime.log(`Signature (recovery): ${signature.recovery}`);
  runtime.log(`Signature (compact): ${bytesToHex(signature.toCompactRawBytes())}`);
  runtime.log(`Signature (DER): ${bytesToHex(signature.toDERRawBytes())}`);

  // Verify the signature
  const isValid = secp256k1.verify(signature, messageHash, publicKey);
  runtime.log(`Signature Valid: ${isValid}`);

  // Recover public key from signature (useful for Ethereum)
  const recoveredPublicKey = signature.recoverPublicKey(messageHash);
  runtime.log(`Recovered Public Key: ${bytesToHex(recoveredPublicKey.toRawBytes())}`);
}

function demoEd25519(runtime: Runtime<Config>): void {
  runtime.log("=== Ed25519 Signatures (Modern, Fast) ===");
  runtime.log("Using randomBytes() with CRE's ChaCha8Rng-based Math.random() polyfill.");

  // Generate random private key using CRE's consensus-safe PRNG
  const privateKey = randomBytes(32);
  runtime.log(`Private Key: ${bytesToHex(privateKey)}`);

  // Derive public key
  const publicKey = ed25519.getPublicKey(privateKey);
  runtime.log(`Public Key: ${bytesToHex(publicKey)}`);

  // Sign a message (Ed25519 hashes internally, no pre-hashing needed)
  const message = utf8ToBytes("Message to sign with Ed25519");
  const signature = ed25519.sign(message, privateKey);
  runtime.log(`Signature: ${bytesToHex(signature)}`);

  // Verify the signature
  const isValid = ed25519.verify(signature, message, publicKey);
  runtime.log(`Signature Valid: ${isValid}`);
}

function demoEcdh(runtime: Runtime<Config>): void {
  runtime.log("=== ECDH Key Exchange ===");
  runtime.log("Using randomBytes() with CRE's ChaCha8Rng-based Math.random() polyfill.");

  // ECDH with secp256k1 (Bitcoin/Ethereum compatible)
  runtime.log("--- secp256k1 ECDH ---");
  const alicePrivKey = randomBytes(32);
  const alicePubKey = secp256k1.getPublicKey(alicePrivKey);

  const bobPrivKey = randomBytes(32);
  const bobPubKey = secp256k1.getPublicKey(bobPrivKey);

  runtime.log(`Alice Public Key: ${bytesToHex(alicePubKey)}`);
  runtime.log(`Bob Public Key: ${bytesToHex(bobPubKey)}`);

  // Both parties derive the same shared secret
  const aliceSharedSecret = secp256k1.getSharedSecret(alicePrivKey, bobPubKey);
  const bobSharedSecret = secp256k1.getSharedSecret(bobPrivKey, alicePubKey);

  runtime.log(`Alice's Shared Secret: ${bytesToHex(aliceSharedSecret)}`);
  runtime.log(`Bob's Shared Secret: ${bytesToHex(bobSharedSecret)}`);
  runtime.log(`Secrets Match: ${bytesToHex(aliceSharedSecret) === bytesToHex(bobSharedSecret)}`);

  // X25519 ECDH (Modern, recommended for new applications)
  runtime.log("--- X25519 ECDH (Curve25519) ---");
  const aliceX25519Priv = randomBytes(32);
  const aliceX25519Pub = x25519.getPublicKey(aliceX25519Priv);

  const bobX25519Priv = randomBytes(32);
  const bobX25519Pub = x25519.getPublicKey(bobX25519Priv);

  runtime.log(`Alice X25519 Public Key: ${bytesToHex(aliceX25519Pub)}`);
  runtime.log(`Bob X25519 Public Key: ${bytesToHex(bobX25519Pub)}`);

  const aliceX25519Shared = x25519.getSharedSecret(aliceX25519Priv, bobX25519Pub);
  const bobX25519Shared = x25519.getSharedSecret(bobX25519Priv, aliceX25519Pub);

  runtime.log(`Alice's X25519 Shared Secret: ${bytesToHex(aliceX25519Shared)}`);
  runtime.log(`Bob's X25519 Shared Secret: ${bytesToHex(bobX25519Shared)}`);
  runtime.log(`X25519 Secrets Match: ${bytesToHex(aliceX25519Shared) === bytesToHex(bobX25519Shared)}`);
}

function demoUtilities(runtime: Runtime<Config>): void {
  runtime.log("=== UTILITY FUNCTIONS ===");

  // Random bytes using CRE's Math.random() polyfill (ChaCha8Rng PRNG)
  // Note: Noble's randomBytes() won't work, but our custom one does!
  runtime.log("Random bytes using CRE's ChaCha8Rng-based Math.random():");
  const randomData = randomBytes(16);
  runtime.log(`Random Bytes (16): ${bytesToHex(randomData)}`);

  // Hex encoding/decoding
  const originalHex = "48656c6c6f2c20576f726c6421"; // "Hello, World!" in hex
  const decodedBytes = hexToBytes(originalHex);
  const reEncodedHex = bytesToHex(decodedBytes);
  runtime.log(`Original Hex: ${originalHex}`);
  runtime.log(`Decoded Text: ${new TextDecoder().decode(decodedBytes)}`);
  runtime.log(`Re-encoded Hex: ${reEncodedHex}`);

  // UTF-8 encoding
  const text = "Hello, World!";
  const utf8Bytes = utf8ToBytes(text);
  runtime.log(`UTF-8 Text: ${text}`);
  runtime.log(`UTF-8 Bytes (hex): ${bytesToHex(utf8Bytes)}`);
  runtime.log(`UTF-8 Byte Length: ${utf8Bytes.length}`);
}

function demoBitcoinAddressDerivation(runtime: Runtime<Config>): void {
  runtime.log("=== BITCOIN ADDRESS DERIVATION EXAMPLE ===");

  // Use a deterministic private key for demo
  const privateKey = deriveBytes("demo-bitcoin-private-key", 32);
  runtime.log(`Private Key (WIF would be derived from this): ${bytesToHex(privateKey)}`);

  // Get compressed public key
  const publicKey = secp256k1.getPublicKey(privateKey, true);
  runtime.log(`Compressed Public Key: ${bytesToHex(publicKey)}`);

  // Bitcoin P2PKH address derivation: RIPEMD160(SHA256(publicKey))
  const sha256Hash = sha256(publicKey);
  const hash160 = ripemd160(sha256Hash);
  runtime.log(`Hash160 (RIPEMD160(SHA256(pubkey))): ${bytesToHex(hash160)}`);

  // Note: To get actual Bitcoin address, you'd need Base58Check encoding
  // which requires additional libraries (like @scure/base)
}

function demoEthereumAddressDerivation(runtime: Runtime<Config>): void {
  runtime.log("=== ETHEREUM ADDRESS DERIVATION EXAMPLE ===");

  // Use a deterministic private key for demo
  const privateKey = deriveBytes("demo-ethereum-private-key", 32);
  runtime.log(`Private Key: ${bytesToHex(privateKey)}`);

  // Get uncompressed public key (remove 0x04 prefix)
  const publicKeyUncompressed = secp256k1.getPublicKey(privateKey, false);
  const publicKeyWithoutPrefix = publicKeyUncompressed.slice(1); // Remove 0x04 prefix
  runtime.log(`Public Key (without prefix): ${bytesToHex(publicKeyWithoutPrefix)}`);

  // Ethereum address: last 20 bytes of Keccak256(publicKey)
  const keccakHash = keccak_256(publicKeyWithoutPrefix);
  const ethereumAddress = keccakHash.slice(-20); // Last 20 bytes
  runtime.log(`Ethereum Address: 0x${bytesToHex(ethereumAddress)}`);
}

// ============================================================================
// Main Workflow Handler
// ============================================================================

const onCronTrigger = (runtime: Runtime<Config>): string => {
  runtime.log("========================================");
  runtime.log("NOBLE CRYPTO LIBRARIES DEMO");
  runtime.log("Alternative to Node.js crypto module");
  runtime.log("Compatible with QuickJS / CRE Workflows");
  runtime.log("========================================");
  runtime.log("");

  demoHashing(runtime);
  runtime.log("");

  demoHmac(runtime);
  runtime.log("");

  demoKeyDerivation(runtime);
  runtime.log("");

  demoSymmetricEncryption(runtime);
  runtime.log("");

  demoEcdsaSecp256k1(runtime);
  runtime.log("");

  demoEd25519(runtime);
  runtime.log("");

  demoEcdh(runtime);
  runtime.log("");

  demoUtilities(runtime);
  runtime.log("");

  demoBitcoinAddressDerivation(runtime);
  runtime.log("");

  demoEthereumAddressDerivation(runtime);
  runtime.log("");

  runtime.log("========================================");
  runtime.log("DEMO COMPLETE");
  runtime.log("========================================");

  return "Crypto demo completed successfully";
};

const initWorkflow = (config: Config) => {
  const cron = new CronCapability();

  return [
    handler(
      cron.trigger(
        { schedule: config.schedule }
      ),
      onCronTrigger
    ),
  ];
};

export async function main() {
  const runner = await Runner.newRunner<Config>();
  await runner.run(initWorkflow);
}
