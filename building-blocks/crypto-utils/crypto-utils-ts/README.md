# Crypto Utils - CRE Building Block (TypeScript)

**⚠️ DISCLAIMER**

This tutorial represents an educational example to use a Chainlink system, product, or service and is provided to demonstrate how to interact with Chainlink's systems, products, and services to integrate them into your own. This template is provided "AS IS" and "AS AVAILABLE" without warranties of any kind, it has not been audited, and it may be missing key checks or error handling to make the usage of the system, product or service more clear. Do not use the code in this example in a production environment without completing your own audits and application of best practices. Neither Chainlink Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for unintended outputs that are generated due to errors in code.

---

This building block demonstrates how to use cryptographic operations in CRE TypeScript workflows. Since the CRE TypeScript SDK runs on QuickJS (a lightweight JavaScript engine), the standard Node.js `crypto` module is not available. This template shows how to use the **Noble crypto libraries** as a drop-in alternative.

## Why Noble Libraries?

The [Noble crypto libraries](https://paulmillr.com/noble/) are:
- **Pure JavaScript** - No native dependencies, works in any JS environment
- **Audited** - Security-audited implementations
- **Minimal** - Small bundle size, tree-shakeable
- **Standards-compliant** - Implements widely-used cryptographic standards

## Features Demonstrated

This workflow demonstrates the following cryptographic operations:

| Category | Operations |
|----------|------------|
| **Hashing** | SHA-256, SHA-384, SHA-512, SHA3-256, SHA3-512, Keccak-256, BLAKE2b, BLAKE3, RIPEMD-160 |
| **HMAC** | HMAC-SHA256, HMAC-SHA512, HMAC-SHA3-256 |
| **Key Derivation** | PBKDF2, Scrypt, HKDF |
| **Symmetric Encryption** | AES-256-GCM, ChaCha20-Poly1305 |
| **Digital Signatures** | ECDSA (secp256k1), Ed25519 |
| **Key Exchange** | ECDH (secp256k1), X25519 |
| **Utilities** | Hex encoding/decoding, UTF-8 encoding |
| **Random Bytes** | Using CRE's ChaCha8Rng-based `Math.random()` polyfill |
| **Real-world Examples** | Bitcoin address derivation, Ethereum address derivation |

---

## Random Bytes in CRE

**Noble's `randomBytes()` won't work** because it requires `crypto.getRandomValues` (Web Crypto API) which QuickJS doesn't support.

**However, CRE provides a `Math.random()` polyfill** that IS consensus-safe:

- **Implementation**: ChaCha8Rng (cryptographic PRNG) in the Javy WASM plugin
- **Seeding**: Host runtime provides seed via `random_seed()` import
- **Consensus-safe**: Same seed produces same sequence across all DON nodes
- **Mode isolation**: Different execution modes get independent RNG streams

### Using Random Bytes in CRE

Create your own `randomBytes()` function using CRE's `Math.random()`:

```typescript
function randomBytes(length: number): Uint8Array {
  const bytes = new Uint8Array(length);
  for (let i = 0; i < length; i++) {
    bytes[i] = Math.floor(Math.random() * 256);
  }
  return bytes;
}

// Use it for keys and nonces
const aesKey = randomBytes(32);    // 256-bit key
const nonce = randomBytes(12);      // 96-bit nonce
const privateKey = randomBytes(32); // EC private key
```

---

## Node.js crypto to Noble Library Mapping

Use this table to migrate from Node.js `crypto` module to Noble libraries:

| Node.js crypto | Noble Library | Import | Code Reference |
|----------------|---------------|--------|----------------|
| `crypto.createHash('sha256')` | `@noble/hashes/sha2` | `import { sha256 } from "@noble/hashes/sha2"` | [main.ts:89](./workflow/main.ts#L89) |
| `crypto.createHash('sha512')` | `@noble/hashes/sha2` | `import { sha512 } from "@noble/hashes/sha2"` | [main.ts:95](./workflow/main.ts#L95) |
| `crypto.createHash('sha3-256')` | `@noble/hashes/sha3` | `import { sha3_256 } from "@noble/hashes/sha3"` | [main.ts:99](./workflow/main.ts#L99) |
| `crypto.createHash('ripemd160')` | `@noble/hashes/ripemd160` | `import { ripemd160 } from "@noble/hashes/ripemd160"` | [main.ts:117](./workflow/main.ts#L117) |
| `crypto.createHmac('sha256', key)` | `@noble/hashes/hmac` | `import { hmac } from "@noble/hashes/hmac"` | [main.ts:135](./workflow/main.ts#L135) |
| `crypto.pbkdf2()` | `@noble/hashes/pbkdf2` | `import { pbkdf2 } from "@noble/hashes/pbkdf2"` | [main.ts:155](./workflow/main.ts#L155) |
| `crypto.scrypt()` | `@noble/hashes/scrypt` | `import { scrypt } from "@noble/hashes/scrypt"` | [main.ts:160](./workflow/main.ts#L160) |
| `crypto.hkdf()` | `@noble/hashes/hkdf` | `import { hkdf } from "@noble/hashes/hkdf"` | [main.ts:167](./workflow/main.ts#L167) |
| `crypto.createCipheriv('aes-256-gcm')` | `@noble/ciphers/aes` | `import { gcm } from "@noble/ciphers/aes"` | [main.ts:180](./workflow/main.ts#L180) |
| `crypto.createCipheriv('chacha20-poly1305')` | `@noble/ciphers/chacha` | `import { chacha20poly1305 } from "@noble/ciphers/chacha"` | [main.ts:192](./workflow/main.ts#L192) |
| `crypto.sign()` with ECDSA | `@noble/curves/secp256k1` | `secp256k1.sign(msgHash, privateKey)` | [main.ts:227](./workflow/main.ts#L227) |
| `crypto.verify()` with ECDSA | `@noble/curves/secp256k1` | `secp256k1.verify(signature, msgHash, publicKey)` | [main.ts:238](./workflow/main.ts#L238) |
| `crypto.sign()` with Ed25519 | `@noble/curves/ed25519` | `ed25519.sign(message, privateKey)` | [main.ts:261](./workflow/main.ts#L261) |
| `crypto.diffieHellman()` | `@noble/curves` | `secp256k1.getSharedSecret()` or `x25519.getSharedSecret()` | [main.ts:282](./workflow/main.ts#L282) |
| `crypto.randomBytes()` | Custom using `Math.random()` | See [Random Bytes in CRE](#random-bytes-in-cre) | [main.ts:68](./workflow/main.ts#L68) |
| `Buffer.from(hex, 'hex')` | `@noble/hashes/utils` | `import { hexToBytes } from "@noble/hashes/utils"` | [main.ts:320](./workflow/main.ts#L320) |
| `buffer.toString('hex')` | `@noble/hashes/utils` | `import { bytesToHex } from "@noble/hashes/utils"` | [main.ts:321](./workflow/main.ts#L321) |

---

## Dependencies

The following Noble packages are used:

```json
{
  "@noble/hashes": "^1.7.1",
  "@noble/curves": "^1.8.1",
  "@noble/ciphers": "^1.2.1"
}
```

---

## Setup and Prerequisites

1. **Install CRE CLI**
   ```bash
   # See https://docs.chain.link/cre for installation instructions
   ```

2. **Login to CRE**
   ```bash
   cre login
   ```

3. **Install Bun** (if not already installed)
   ```bash
   # See https://bun.sh/docs/installation
   ```

4. **Install dependencies**
   ```bash
   cd building-blocks/crypto-utils/crypto-utils-ts/workflow
   bun install
   ```

---

## Running the Workflow

### Simulate the workflow

From the project root directory (`crypto-utils-ts`):

```bash
cre workflow simulate workflow
```

---

## Example Output

When the workflow runs, it logs the output of each cryptographic operation:

```
========================================
NOBLE CRYPTO LIBRARIES DEMO
Alternative to Node.js crypto module
Compatible with QuickJS / CRE Workflows
========================================

=== HASHING FUNCTIONS ===
SHA-256: 828ce562b0a732143cddc2895bdb9b4ea98293805b24491853c2970186a36f88
SHA-384: 0a5d403a2d3ca36332a533ea6dfb85fd5facaa3120f8a41cb63cf8ea4da2f3fa...
SHA-512: 81d3babc839e8894eb89e7579050d14ee36e6a8599696e551a89e576ab31e2e5...
SHA3-256: d2a7297f24e4da526fb97568b69f8c1bbe53574d1383d2533700e16c551dce45
Keccak-256 (Ethereum): ec4fa95ae783668d4468af54f9ec9630729ac6a9b4550acb26278b6507f7cf51
BLAKE2b-256: 8a3030f77cae72dc512c41ab1645a179f6e4d31533ffe8a417c37c1b004eaddb
BLAKE3: 3cc4857ee80f9820a775feed897ab77d6bbc6202c44a885f1281c616df6b67c2
RIPEMD-160: 60a86f1115366d6f78a212b5239e8c8f135d6bed

=== HMAC (Hash-based Message Authentication Code) ===
HMAC-SHA256: e3c9d42ece6cbc90bef8b2ef280bd073505941800b250a1d1121125d608311a9
HMAC-SHA512: 6cf6da999feb5cfdcfd5517b7e2ad5bc13caeb492ba7613a8fbaf0744947d731...

=== KEY DERIVATION FUNCTIONS ===
PBKDF2-SHA256 (1k iterations): d484e54ed63ba6b9ad7f1750e566081efae6a9ef93986a03982956da1570b98e
Scrypt (N=1024, r=8, p=1): b7ba286855e8a2f502b47aed41702123b490e4688c5056e90802b2ab7988dba4
HKDF-SHA256: 21b286fe4324bdbe19319aa7ac3d894757a1c66b3752c4a294b5e91165260ac4

=== SYMMETRIC ENCRYPTION (AEAD) ===
NOTE: Using deterministic keys for demo. In production, use secure random keys.
AES-256-GCM Ciphertext: 237133d47a1324a7f4e7144e56e8a985f8a8d0dfa42403d06d9383e9b9b61353...
AES-256-GCM Decrypted: Sensitive data to encrypt
ChaCha20-Poly1305 Ciphertext: 0e95af78139c7050d5daa94f2d252ab7aed590e5f0174de932c36e3a6e19cf20...
ChaCha20-Poly1305 Decrypted: Sensitive data to encrypt

=== ECDSA with secp256k1 (Bitcoin/Ethereum) ===
Private Key: 67211ff271555e843683256b426a4453e4104e09dff3064d59a9486044987a89
Public Key (compressed): 03544cfb0431489a1d4a551b1cd9eca0445b68a15d9a1a02c36aceec17850a4559
Signature Valid: true
Recovered Public Key: 03544cfb0431489a1d4a551b1cd9eca0445b68a15d9a1a02c36aceec17850a4559

=== Ed25519 Signatures (Modern, Fast) ===
Private Key: c6dc9f6ae2eb26ba9c9ea78407fb16d36cab14af6b22e94f08ddc706cef782cb
Public Key: 4895153ca98d10fda4e516ea74575e33a0ba54634cea73c8c1c2b39f6b124521
Signature Valid: true

=== ECDH Key Exchange ===
Alice's Shared Secret: 02529bb8e1a3106f7c7097211aee7598d37f1f5735218e8fff926810c984f55c10
Bob's Shared Secret: 02529bb8e1a3106f7c7097211aee7598d37f1f5735218e8fff926810c984f55c10
Secrets Match: true

=== UTILITY FUNCTIONS ===
NOTE: randomBytes() not available in QuickJS (no crypto.getRandomValues)
Original Hex: 48656c6c6f2c20576f726c6421
Decoded Text: Hello, World!

=== BITCOIN ADDRESS DERIVATION EXAMPLE ===
Hash160 (RIPEMD160(SHA256(pubkey))): 01ec9bb0efbb084ff780954eed1ca941d9567b84

=== ETHEREUM ADDRESS DERIVATION EXAMPLE ===
Ethereum Address: 0xa181b2fcc624d496c7e250c96449f55e716efe12

========================================
DEMO COMPLETE
========================================
```

---

## Use Cases

### 1. Message Signing and Verification
Sign messages or transactions for blockchain interactions:
```typescript
import { secp256k1 } from "@noble/curves/secp256k1";
import { sha256 } from "@noble/hashes/sha2";
import { utf8ToBytes } from "@noble/hashes/utils";

const messageHash = sha256(utf8ToBytes("Transaction data"));
const signature = secp256k1.sign(messageHash, privateKey);
const isValid = secp256k1.verify(signature, messageHash, publicKey);
```

### 2. Password Hashing
Securely hash passwords for storage:
```typescript
import { scrypt } from "@noble/hashes/scrypt";

// Use higher N (2^14 or more) in production for better security
const hashedPassword = scrypt(password, salt, { N: 2 ** 14, r: 8, p: 1, dkLen: 32 });
```

### 3. Data Encryption
Encrypt sensitive data before storage or transmission:
```typescript
import { gcm } from "@noble/ciphers/aes";

// Key and nonce should come from secure sources (see Random Bytes Limitation)
const cipher = gcm(key, nonce);
const encrypted = cipher.encrypt(plaintext);
const decrypted = cipher.decrypt(encrypted);
```

### 4. Ethereum/Bitcoin Address Generation
Derive addresses from private keys:
```typescript
import { secp256k1 } from "@noble/curves/secp256k1";
import { keccak_256 } from "@noble/hashes/sha3";
import { bytesToHex } from "@noble/hashes/utils";

const publicKey = secp256k1.getPublicKey(privateKey, false).slice(1);
const address = "0x" + bytesToHex(keccak_256(publicKey).slice(-20));
```

### 5. Secure Key Exchange (ECDH)
Establish shared secrets between parties:
```typescript
import { x25519 } from "@noble/curves/ed25519";

const sharedSecret = x25519.getSharedSecret(myPrivateKey, theirPublicKey);
// Use HKDF to derive encryption keys from the shared secret
```

---

## Reference Documentation

- [CRE Documentation](https://docs.chain.link/cre)
- [Noble Hashes](https://github.com/paulmillr/noble-hashes)
- [Noble Curves](https://github.com/paulmillr/noble-curves)
- [Noble Ciphers](https://github.com/paulmillr/noble-ciphers)
- [Noble Overview](https://paulmillr.com/noble/)

---

## License

MIT - see the repository's [LICENSE](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE).
