<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE&message=Home&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE&message=Docs&color=blue)](https://docs.chain.link/cre)

</div>

# Compression Utils - CRE Building Block

This building block demonstrates how to perform compression and decompression operations in CRE workflows using **fflate** as a pure JavaScript alternative to the Node.js `zlib` module.

## The Problem

The CRE TypeScript SDK runs on **QuickJS**, a lightweight JavaScript engine that does not support Node.js native modules. This means the standard `zlib` module from Node.js is **not available** in CRE workflows. Additionally, CRE's HTTP capability has a **100KB response body limit**, making compression essential for working with large datasets.

## The Solution

The [fflate library](https://github.com/101arrowz/fflate) provides pure JavaScript implementations of compression algorithms:

- **Pure JavaScript** - No native dependencies, works in any JS environment including QuickJS
- **Standards-compliant** - Output is interoperable with native gzip, zlib, and deflate tools
- **Small** - ~8KB minified, tree-shakeable
- **Sync-capable** - Provides synchronous APIs (`gzipSync`, `deflateSync`, etc.) that work reliably in QuickJS
- **Practical** - Enable workflows to handle large datasets despite the 100KB HTTP limit

---

## What's Covered

This template demonstrates:

| Category | fflate API | Use Cases |
|----------|------------|-----------|
| **Gzip Compression** | `gzipSync` / `gunzipSync` | Compress/decompress with header (standard `.gz` format) |
| **Raw Deflate** | `deflateSync` / `inflateSync` | Lightweight compression without headers |
| **Zlib Format** | `zlibSync` / `unzlibSync` | zlib wrapper compression (common in APIs) |
| **Auto-detect Decompression** | `decompressSync` | Handle gzip, zlib, or deflate automatically |
| **Compression Levels** | `gzipSync(data, { level })` | Tune speed vs. compression (0-9) |
| **ZIP Archives** | `zipSync` / `unzipSync` | Bundle multiple files into a single archive |
| **String Utilities** | `strToU8` / `strFromU8` | Convert between strings and byte arrays |

## Get Started

- **TypeScript**: See the [TypeScript README](./compression-utils-ts/README.md) for detailed setup, usage examples, and a complete Node.js zlib-to-fflate mapping table.

## Quick Example

```typescript
// Instead of Node.js zlib:
// const compressed = require('zlib').gzipSync(data);

// Use fflate:
import { gzipSync, gunzipSync, strToU8, strFromU8 } from "fflate";

const data = strToU8(JSON.stringify(largeObject));
const compressed = gzipSync(data);
const restored = JSON.parse(strFromU8(gunzipSync(compressed)));
```

## Reference Documentation

- [fflate GitHub](https://github.com/101arrowz/fflate)
- [fflate npm](https://www.npmjs.com/package/fflate)
- [CRE Documentation](https://docs.chain.link/cre)
- [QuickJS Engine](https://bellard.org/quickjs/)