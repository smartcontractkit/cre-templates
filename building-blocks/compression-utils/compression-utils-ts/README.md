# Compression Utils - CRE Building Block (TypeScript)

**⚠️ DISCLAIMER**

This tutorial represents an educational example to use a Chainlink system, product, or service and is provided to demonstrate how to interact with Chainlink's systems, products, and services to integrate them into your own. This template is provided "AS IS" and "AS AVAILABLE" without warranties of any kind, it has not been audited, and it may be missing key checks or error handling to make the usage of the system, product or service more clear. Do not use the code in this example in a production environment without completing your own audits and application of best practices. Neither Chainlink Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for unintended outputs that are generated due to errors in code.

---

This building block demonstrates how to use compression and decompression in CRE TypeScript workflows. Since the CRE TypeScript SDK runs on QuickJS (a lightweight JavaScript engine), the standard Node.js `zlib` module is not available. This template shows how to use **[fflate](https://github.com/101arrowz/fflate)** as a pure JavaScript drop-in alternative.

## Why fflate?

[fflate](https://github.com/101arrowz/fflate) is:
- **Pure JavaScript** - No native dependencies, works in any JS environment including QuickJS
- **Standards-compliant** - Output is interoperable with native gzip, zlib, and deflate tools
- **Small** - ~8KB minified, tree-shakeable
- **Sync-capable** - Provides synchronous APIs (`gzipSync`, `deflateSync`, etc.) that work reliably in QuickJS

## Features Demonstrated

This workflow fetches a large JSON payload (~65KB, 215 records) via HTTP, compresses it to pass consensus, decompresses it, and uses it as the input for each compression demo:

| Category | fflate API | Node.js Equivalent |
|----------|------------|--------------------|
| **Gzip** | `gzipSync` / `gunzipSync` | `zlib.gzipSync` / `zlib.gunzipSync` |
| **Raw Deflate** | `deflateSync` / `inflateSync` | `zlib.deflateRawSync` / `zlib.inflateRawSync` |
| **Zlib** | `zlibSync` / `unzlibSync` | `zlib.deflateSync` / `zlib.inflateSync` |
| **Auto-detect** | `decompressSync` | — |
| **Compression levels** | `gzipSync(data, { level })` | `zlib.gzipSync(data, { level })` |
| **ZIP archives** | `zipSync` / `unzipSync` | — (no single-call equivalent) |
| **String utilities** | `strToU8` / `strFromU8` | `Buffer.from()` / `TextEncoder` / `TextDecoder` |

---

## HTTP Fetch in CRE

The workflow fetches its payload using the CRE HTTP capability (`HTTPClient`), which routes requests through the DON's off-chain network layer. The standard `fetch` / `node:https` APIs are not available in QuickJS.

> **Note:** The CRE HTTP capability has a response body limit of **100KB** but the Consensus capability has a smaller limit **25KB**. This workflow fetches `/comments?_limit=215` (~65KB) and compresses it (~24KB) to stay within those limits.

```typescript
import { HTTPClient, consensusIdenticalAggregation, type HTTPSendRequester, text } from "@chainlink/cre-sdk";

const fetchComments = (sendRequester: HTTPSendRequester): string => {
  const resp = sendRequester
    .sendRequest({ url: "https://jsonplaceholder.typicode.com/comments?_limit=215", method: "GET" })
    .result();
  return text(resp);
};

// In the handler:
const httpClient = new HTTPClient();
const rawJson = httpClient
  .sendRequest(runtime, fetchComments, consensusIdenticalAggregation<string>())()
  .result();
```

`consensusIdenticalAggregation` ensures all DON nodes agree on an identical response before the workflow proceeds.

---

## Node.js zlib to fflate Mapping

| Node.js zlib | fflate | Import | Code Reference |
|---|---|---|---|
| `zlib.gzipSync(data)` | `gzipSync(data)` | `import { gzipSync } from "fflate"` | [main.ts:62](./workflow/main.ts#L62) |
| `zlib.gunzipSync(data)` | `gunzipSync(data)` | `import { gunzipSync } from "fflate"` | [main.ts:66](./workflow/main.ts#L66) |
| `zlib.deflateRawSync(data)` | `deflateSync(data)` | `import { deflateSync } from "fflate"` | [main.ts:75](./workflow/main.ts#L75) |
| `zlib.inflateRawSync(data)` | `inflateSync(data)` | `import { inflateSync } from "fflate"` | [main.ts:79](./workflow/main.ts#L79) |
| `zlib.deflateSync(data)` | `zlibSync(data)` | `import { zlibSync } from "fflate"` | [main.ts:88](./workflow/main.ts#L88) |
| `zlib.inflateSync(data)` | `unzlibSync(data)` | `import { unzlibSync } from "fflate"` | [main.ts:92](./workflow/main.ts#L92) |
| *(auto-detect format)* | `decompressSync(data)` | `import { decompressSync } from "fflate"` | [main.ts:108](./workflow/main.ts#L108) |
| `zlib.gzipSync(data, { level })` | `gzipSync(data, { level })` | `import { gzipSync } from "fflate"` | [main.ts:125](./workflow/main.ts#L125) |
| *(no equivalent)* | `zipSync(files)` | `import { zipSync } from "fflate"` | [main.ts:147](./workflow/main.ts#L147) |
| *(no equivalent)* | `unzipSync(data)` | `import { unzipSync } from "fflate"` | [main.ts:150](./workflow/main.ts#L150) |
| `Buffer.from(str)` / `TextEncoder` | `strToU8(str)` | `import { strToU8 } from "fflate"` | [main.ts:162](./workflow/main.ts#L162) |
| `buf.toString()` / `TextDecoder` | `strFromU8(bytes)` | `import { strFromU8 } from "fflate"` | [main.ts:165](./workflow/main.ts#L165) |

---

## Dependencies

```json
{
  "fflate": "^0.8.2"
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
   cd building-blocks/compression-utils/compression-utils-ts/workflow
   bun install
   ```

---

## Running the Workflow

### Simulate the workflow

From the project root directory (`compression-utils-ts`):

```bash
cre workflow simulate workflow
```

---

## Example Output

```
========================================
FFLATE COMPRESSION LIBRARIES DEMO
Alternative to Node.js zlib module
Compatible with QuickJS / CRE Workflows
========================================

Fetching payload from JSONPlaceholder /comments (compressed over consensus)...
Consensus payload: 23.58 KB (base64-encoded gzip)
Decompressed: 65.04 KB — 215 records

=== GZIP / GUNZIP ===
Alternative to: zlib.gzipSync() / zlib.gunzipSync()
Original:   65.04 KB
Compressed: 17.71 KB (72.8% smaller)
Verified:   215 records restored via gunzipSync

=== DEFLATE / INFLATE (raw) ===
Alternative to: zlib.deflateRawSync() / zlib.inflateRawSync()
Original:   65.04 KB
Compressed: 17.69 KB (72.8% smaller)
Verified:   215 records restored via inflateSync

=== ZLIB / UNZLIB ===
Alternative to: zlib.deflateSync() / zlib.inflateSync()
Original:   65.04 KB
Compressed: 17.70 KB (72.8% smaller)
Verified:   215 records restored via unzlibSync

=== AUTO-DETECT DECOMPRESSION ===
decompressSync() detects gzip / zlib / deflate automatically.
From gzip   : 215 records (17.71 KB)
From zlib   : 215 records (17.70 KB)
From deflate: 215 records (17.69 KB)

=== COMPRESSION LEVELS ===
Level 0 = store only  |  Level 4 |  Level 9 = max
Level 0 (store): 65.07 KB (-0.0% smaller)
Level 4: 17.99 KB (72.3% smaller)
Level 9 (max): 17.68 KB (72.8% smaller)

=== ZIP ARCHIVES ===
Multi-file archiving — no Node.js zlib equivalent.
Archive:   19.09 KB (4 files, 215 total records)
Extracted: batch-1.json, batch-2.json, batch-3.json, manifest.json
Verified:  215 records across 3 batches

=== STRING UTILITIES ===
strToU8 / strFromU8 — alternative to Buffer.from() / TextEncoder / TextDecoder
strToU8(rawJson):  65.04 KB Uint8Array
strFromU8(bytes):  215 records decoded
strToU8("Hello, CRE Workflow!"): [72, 101, 108, 108, 111, 44, 32, 67, 82, 69, 32, 87, 111, 114, 107, 102, 108, 111, 119, 33]
strFromU8(...):    "Hello, CRE Workflow!"

========================================
DEMO COMPLETE
========================================
```

---

## Use Cases

### 1. Compress large API responses before on-chain storage
Reduce gas costs by compressing data before writing it on-chain or to decentralized storage:
```typescript
import { gzipSync, gunzipSync, strToU8, strFromU8 } from "fflate";

const payload = JSON.stringify(largeDataObject);
const compressed = gzipSync(strToU8(payload));
// Store `compressed` — decompress on read
const restored = JSON.parse(strFromU8(gunzipSync(compressed)));
```

### 2. Decompress API responses to get around 25KB limit
The CRE HTTP capability has a 100KB response body limit but the Consensus capability has a 25KB limit. Compress large payloads before passing it to consensus and decompress in the workflow:
```typescript
import { gunzipSync, strFromU8 } from "fflate";

// Server sends data (~65KB)
// which compresses to ~24KB consensus payload
const compressedData = sendRequester.sendRequest({ 
  url: "https://api.example.com/data?format=gzip" 
}).result();

const raw = text(resp);
const compressed = gzipSync(strToU8(raw), { level: 9 });
const compressedBase64 = fromByteArray(compressed);

const compressedBytes = toByteArray(compressedBase64);
const fullData = JSON.parse(strFromU8(gunzipSync(compressedBytes)));
// Now you can work with the full dataset despite the 25KB limit
```

### 3. Auto-detect and decompress unknown formats
Handle compressed payloads from external sources without knowing the format upfront:
```typescript
import { decompressSync, strFromU8 } from "fflate";

// Works regardless of whether the source used gzip, zlib, or raw deflate
const decompressed = strFromU8(decompressSync(incomingBytes));
```

### 4. Tune compression for speed vs. size
Use lower levels for real-time workflows, higher levels when minimizing output size matters:
```typescript
import { gzipSync, strToU8 } from "fflate";

// Fast path — minimal runtime
const fast = gzipSync(data, { level: 1 });

// Storage path — smallest possible output
const compact = gzipSync(data, { level: 9 });
```

---

## Reference Documentation

- [CRE Documentation](https://docs.chain.link/cre)
- [fflate GitHub](https://github.com/101arrowz/fflate)
- [fflate npm](https://www.npmjs.com/package/fflate)

---

## License

MIT - see the repository's [LICENSE](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE).
