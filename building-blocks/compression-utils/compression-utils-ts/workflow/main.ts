import {
  CronCapability,
  consensusIdenticalAggregation,
  handler,
  HTTPClient,
  type HTTPSendRequester,
  Runner,
  text,
  type Runtime,
} from "@chainlink/cre-sdk";

import {
  deflateSync,
  decompressSync,
  gzipSync,
  gunzipSync,
  inflateSync,
  strFromU8,
  strToU8,
  unzlibSync,
  unzipSync,
  zlibSync,
  zipSync,
} from "fflate";

type Config = {
  schedule: string;
};

// ============================================================================
// HTTP Fetch — JSONPlaceholder /comments?_limit=300 (~92KB, 300 records)
// Used as the input payload for all compression demos below.
// ============================================================================

const fetchComments = (sendRequester: HTTPSendRequester): string => {
  const resp = sendRequester
    .sendRequest({
      url: "https://jsonplaceholder.typicode.com/comments?_limit=300",
      method: "GET",
    })
    .result();
  return text(resp);
};

// ============================================================================
// Helpers
// ============================================================================

function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  return `${(bytes / 1024).toFixed(2)} KB`;
}

function reductionPct(original: number, compressed: number): string {
  return `${((1 - compressed / original) * 100).toFixed(1)}% smaller`;
}

function demoGzip(runtime: Runtime<Config>, jsonBytes: Uint8Array): void {
  runtime.log("=== GZIP / GUNZIP ===");
  runtime.log("Alternative to: zlib.gzipSync() / zlib.gunzipSync()");

  const compressed = gzipSync(jsonBytes);
  runtime.log(`Original:   ${formatBytes(jsonBytes.length)}`);
  runtime.log(`Compressed: ${formatBytes(compressed.length)} (${reductionPct(jsonBytes.length, compressed.length)})`);

  const decompressed = gunzipSync(compressed);
  const records = (JSON.parse(strFromU8(decompressed)) as unknown[]).length;
  runtime.log(`Verified:   ${records} records restored via gunzipSync`);
}

function demoDeflate(runtime: Runtime<Config>, jsonBytes: Uint8Array): void {
  runtime.log("=== DEFLATE / INFLATE (raw) ===");
  runtime.log("Alternative to: zlib.deflateRawSync() / zlib.inflateRawSync()");

  const compressed = deflateSync(jsonBytes);
  runtime.log(`Original:   ${formatBytes(jsonBytes.length)}`);
  runtime.log(`Compressed: ${formatBytes(compressed.length)} (${reductionPct(jsonBytes.length, compressed.length)})`);

  const decompressed = inflateSync(compressed);
  const records = (JSON.parse(strFromU8(decompressed)) as unknown[]).length;
  runtime.log(`Verified:   ${records} records restored via inflateSync`);
}

function demoZlib(runtime: Runtime<Config>, jsonBytes: Uint8Array): void {
  runtime.log("=== ZLIB / UNZLIB ===");
  runtime.log("Alternative to: zlib.deflateSync() / zlib.inflateSync()");

  const compressed = zlibSync(jsonBytes);
  runtime.log(`Original:   ${formatBytes(jsonBytes.length)}`);
  runtime.log(`Compressed: ${formatBytes(compressed.length)} (${reductionPct(jsonBytes.length, compressed.length)})`);

  const decompressed = unzlibSync(compressed);
  const records = (JSON.parse(strFromU8(decompressed)) as unknown[]).length;
  runtime.log(`Verified:   ${records} records restored via unzlibSync`);
}

function demoAutoDetect(runtime: Runtime<Config>, jsonBytes: Uint8Array): void {
  runtime.log("=== AUTO-DETECT DECOMPRESSION ===");
  runtime.log("decompressSync() detects gzip / zlib / deflate automatically.");

  const formats: Array<{ label: string; compressed: Uint8Array }> = [
    { label: "gzip",    compressed: gzipSync(jsonBytes) },
    { label: "zlib",    compressed: zlibSync(jsonBytes) },
    { label: "deflate", compressed: deflateSync(jsonBytes) },
  ];

  for (const { label, compressed } of formats) {
    const decompressed = decompressSync(compressed);
    const records = (JSON.parse(strFromU8(decompressed)) as unknown[]).length;
    runtime.log(`From ${label.padEnd(7)}: ${records} records (${formatBytes(compressed.length)})`);
  }
}

function demoCompressionLevels(runtime: Runtime<Config>, jsonBytes: Uint8Array): void {
  runtime.log("=== COMPRESSION LEVELS ===");
  runtime.log("Level 0 = store only  |  Level 4 |  Level 9 = max");

  const levels = [
    { level: 0 as const, label: "Level 0 (store)" },
    { level: 4 as const, label: "Level 4" },
    { level: 9 as const, label: "Level 9 (max)" },
  ];

  for (const { level, label } of levels) {
    const compressed = gzipSync(jsonBytes, { level });
    runtime.log(`${label}: ${formatBytes(compressed.length)} (${reductionPct(jsonBytes.length, compressed.length)})`);
  }
}

function demoZip(runtime: Runtime<Config>, rawJson: string): void {
  runtime.log("=== ZIP ARCHIVES ===");
  runtime.log("Multi-file archiving — no Node.js zlib equivalent.");

  const allComments = JSON.parse(rawJson) as unknown[];
  const third = Math.floor(allComments.length / 3);

  const files = {
    "batch-1.json": strToU8(JSON.stringify(allComments.slice(0, third))),
    "batch-2.json": strToU8(JSON.stringify(allComments.slice(third, 2 * third))),
    "batch-3.json": strToU8(JSON.stringify(allComments.slice(2 * third))),
    "manifest.json": strToU8(JSON.stringify({
      total: allComments.length,
      files: ["batch-1.json", "batch-2.json", "batch-3.json"],
    })),
  };

  const zipped = zipSync(files);
  runtime.log(`Archive:   ${formatBytes(zipped.length)} (4 files, ${allComments.length} total records)`);

  const unzipped = unzipSync(zipped);
  const batch1 = JSON.parse(strFromU8(unzipped["batch-1.json"])) as unknown[];
  const batch2 = JSON.parse(strFromU8(unzipped["batch-2.json"])) as unknown[];
  const batch3 = JSON.parse(strFromU8(unzipped["batch-3.json"])) as unknown[];
  runtime.log(`Extracted: ${Object.keys(unzipped).join(", ")}`);
  runtime.log(`Verified:  ${batch1.length + batch2.length + batch3.length} records across 3 batches`);
}

function demoStringUtils(runtime: Runtime<Config>, rawJson: string): void {
  runtime.log("=== STRING UTILITIES ===");
  runtime.log("strToU8 / strFromU8 — alternative to Buffer.from() / TextEncoder / TextDecoder");

  const bytes = strToU8(rawJson);
  runtime.log(`strToU8(rawJson):  ${formatBytes(bytes.length)} Uint8Array`);

  const restored = strFromU8(bytes);
  runtime.log(`strFromU8(bytes):  ${(JSON.parse(restored) as unknown[]).length} records decoded`);

  const hello = strToU8("Hello, CRE Workflow!");
  runtime.log(`strToU8("Hello, CRE Workflow!"): [${hello.join(", ")}]`);
  runtime.log(`strFromU8(...):    "${strFromU8(hello)}"`);
}

// ============================================================================
// Main Workflow Handler
// ============================================================================

const onCronTrigger = (runtime: Runtime<Config>): string => {
  runtime.log("========================================");
  runtime.log("FFLATE COMPRESSION LIBRARIES DEMO");
  runtime.log("Alternative to Node.js zlib module");
  runtime.log("Compatible with QuickJS / CRE Workflows");
  runtime.log("========================================");
  runtime.log("");

  runtime.log("Fetching payload from JSONPlaceholder /comments...");
  const httpClient = new HTTPClient();
  const rawJson = httpClient
    .sendRequest(runtime, fetchComments, consensusIdenticalAggregation<string>())()
    .result();

  const jsonBytes = strToU8(rawJson);
  const recordCount = (JSON.parse(rawJson) as unknown[]).length;
  runtime.log(`Fetched:   ${formatBytes(jsonBytes.length)} — ${recordCount} records`);
  runtime.log("");

  demoGzip(runtime, jsonBytes);
  runtime.log("");

  demoDeflate(runtime, jsonBytes);
  runtime.log("");

  demoZlib(runtime, jsonBytes);
  runtime.log("");

  demoAutoDetect(runtime, jsonBytes);
  runtime.log("");

  demoCompressionLevels(runtime, jsonBytes);
  runtime.log("");

  demoZip(runtime, rawJson);
  runtime.log("");

  demoStringUtils(runtime, rawJson);
  runtime.log("");

  runtime.log("========================================");
  runtime.log("DEMO COMPLETE");
  runtime.log("========================================");

  return "Compression demo completed successfully";
};

export const initWorkflow = (config: Config) => {
  const cron = new CronCapability();

  return [
    handler(
      cron.trigger({ schedule: config.schedule }),
      onCronTrigger
    ),
  ];
};

export async function main() {
  const runner = await Runner.newRunner<Config>();
  await runner.run(initWorkflow);
}
