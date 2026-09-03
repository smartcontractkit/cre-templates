<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE&message=Home&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE&message=Docs&color=blue)](https://docs.chain.link/cre)

</div>

# Crypto Utils - CRE Building Block

This building block demonstrates how to perform cryptographic operations in CRE workflows using the **Noble crypto libraries** as an alternative to the Node.js `crypto` module.

## The Problem

The CRE TypeScript SDK runs on **QuickJS**, a lightweight JavaScript engine that does not support Node.js native modules. This means the standard `crypto` module from Node.js is **not available** in CRE workflows.

## The Solution

The [Noble crypto libraries](https://paulmillr.com/noble/) provide pure JavaScript implementations of common cryptographic algorithms. They are:

- **Pure JavaScript** - No native dependencies, works in QuickJS
- **Audited** - Security-audited implementations
- **Minimal** - Small bundle size (~29KB gzipped for curves)
- **Standards-compliant** - Implements widely-used cryptographic standards

## Random Bytes

Noble's `randomBytes()` won't work (needs Web Crypto API), but **CRE provides a `Math.random()` polyfill** using ChaCha8Rng that IS consensus-safe. You can build your own `randomBytes()`:

```typescript
function randomBytes(length: number): Uint8Array {
  const bytes = new Uint8Array(length);
  for (let i = 0; i < length; i++) {
    bytes[i] = Math.floor(Math.random() * 256);
  }
  return bytes;
}
```

See the [TypeScript README](./crypto-utils-ts/README.md#random-bytes-in-cre) for details.

---

## What's Covered

This template demonstrates:

| Category | Noble Package | Operations |
|----------|---------------|------------|
| **Hashing** | `@noble/hashes` | SHA-256, SHA-512, SHA3, Keccak-256, BLAKE2b, BLAKE3, RIPEMD-160 |
| **HMAC** | `@noble/hashes` | HMAC with any hash function |
| **Key Derivation** | `@noble/hashes` | PBKDF2, Scrypt, HKDF |
| **Symmetric Encryption** | `@noble/ciphers` | AES-GCM, ChaCha20-Poly1305 |
| **Digital Signatures** | `@noble/curves` | ECDSA (secp256k1), Ed25519 |
| **Key Exchange** | `@noble/curves` | ECDH, X25519 |
| **Utilities** | `@noble/hashes` | Hex encoding, UTF-8 encoding |
| **Random Bytes** | `Math.random()` polyfill | Consensus-safe PRNG (ChaCha8Rng) |

## Get Started

- **TypeScript**: See the [TypeScript README](./crypto-utils-ts/README.md) for detailed setup, usage examples, and a complete Node.js-to-Noble mapping table.

## Quick Example

```typescript
// Instead of Node.js crypto:
// const hash = crypto.createHash('sha256').update('hello').digest('hex');

// Use Noble:
import { sha256 } from "@noble/hashes/sha2";
import { bytesToHex, utf8ToBytes } from "@noble/hashes/utils";

const hash = bytesToHex(sha256(utf8ToBytes("hello")));
```

## Reference Documentation

- [Noble Crypto Libraries](https://paulmillr.com/noble/)
- [Noble Hashes GitHub](https://github.com/paulmillr/noble-hashes)
- [Noble Curves GitHub](https://github.com/paulmillr/noble-curves)
- [Noble Ciphers GitHub](https://github.com/paulmillr/noble-ciphers)
- [CRE Documentation](https://docs.chain.link/cre)
