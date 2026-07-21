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

import { XMLParser, XMLBuilder, XMLValidator } from "fast-xml-parser";

export type Config = {
  schedule: string;
};

// ============================================================================
// HTTP Fetch — ECB Daily EUR FX Reference Rates (SDMX format)
//
// XML structure (simplified, after removeNSPrefix strips message:/generic:):
//   <GenericData>
//     <DataSet action="Replace" validFromDate="...">
//       <Series>
//         <SeriesKey>
//           <Value id="CURRENCY" value="USD"/>
//           ...
//         </SeriesKey>
//         <Obs>
//           <ObsDimension value="2026-04-17"/>
//           <ObsValue value="1.1797"/>
//         </Obs>
//       </Series>
//       ...
//     </DataSet>
//   </GenericData>
// ============================================================================

const fetchECBRates = (sendRequester: HTTPSendRequester): string => {
  const resp = sendRequester
    .sendRequest({
      url: "https://data-api.ecb.europa.eu/service/data/EXR/D.USD+GBP+JPY+CHF+CNY+AUD+CAD.EUR.SP00.A?format=genericData&lastNObservations=1",
      method: "GET",
    })
    .result();
  return text(resp);
};

type SDMXKeyValue = { "@_id": string; "@_value": string | number };

type SDMXSeries = {
  SeriesKey: { Value: SDMXKeyValue | SDMXKeyValue[] };
  Obs: {
    ObsDimension: { "@_value": string };
    ObsValue: { "@_value": number };
  };
};

type ECBData = {
  Header: { Prepared: string };
  DataSet: { Series: SDMXSeries | SDMXSeries[] };
};

function getDataSet(result: Record<string, unknown>): ECBData {
  const data = result["GenericData"] as ECBData | undefined;
  if (!data) throw new Error(`ECB XML: GenericData not found. Keys: ${Object.keys(result).join(", ")}`);
  return data;
}

// SDMX stores currency as a Value node with id="CURRENCY" rather than as a direct attribute.
function getCurrency(series: SDMXSeries): string {
  const vals = Array.isArray(series.SeriesKey.Value)
    ? series.SeriesKey.Value
    : [series.SeriesKey.Value];
  return String(vals.find((v) => v["@_id"] === "CURRENCY")?.["@_value"] ?? "???");
}

function demoValidation(runtime: Runtime<Config>, rawXml: string): void {
  runtime.log("=== VALIDATION — XMLValidator.validate() ===");
  runtime.log("Alternative to: try/catch around DOMParser.parseFromString()");
  runtime.log("");

  runtime.log("Validating ECB XML...");
  const valid = XMLValidator.validate(rawXml);
  runtime.log(`ECB XML valid: ${valid === true}`);
  runtime.log("");

  runtime.log("Validating malformed XML...");
  const malformed = `<rates><entry currency="USD" rate="1.13</rates>`;
  const err = XMLValidator.validate(malformed);
  if (err !== true) {
    runtime.log(`Malformed XML — code: ${err.err.code}, line: ${err.err.line}`);
    runtime.log(`  ${err.err.msg}`);
  }
}

function demoParsing(runtime: Runtime<Config>, rawXml: string): void {
  runtime.log("=== PARSING — attributes + namespace removal + type casting ===");
  runtime.log("Alternative to: xml2js / DOMParser with manual attribute handling and type casting");
  runtime.log("");

  const parser = new XMLParser({
    ignoreAttributes: false,      // expose XML attributes as properties
    attributeNamePrefix: "@_",    // prefix attrs to distinguish from child elements
    removeNSPrefix: true,         // message:GenericData -> GenericData, generic:Series -> Series
    parseAttributeValue: true,    // ObsValue value="1.1797" -> 1.1797 (number, not string)
  });

  const result  = parser.parse(rawXml) as Record<string, unknown>;
  const ecb     = getDataSet(result);
  const rawSeries = ecb.DataSet.Series;
  const series  = Array.isArray(rawSeries) ? rawSeries : [rawSeries];
  const date    = series[0]?.Obs?.ObsDimension["@_value"] ?? "unknown";

  runtime.log(`Date:       ${date}`);
  runtime.log(`Series:     ${series.length}`);
  runtime.log("");

  for (const s of series) {
    const ccy  = getCurrency(s);
    const rate = s.Obs.ObsValue["@_value"];
    runtime.log(`  1 EUR = ${rate} ${ccy}`);
  }
}

// ============================================================================
// isArray — force consistent array shape
// Without isArray, a single child element parses as an object, not an array.
// This silently breaks any code that calls .map()/.filter() on the result.
// ============================================================================

function demoIsArray(runtime: Runtime<Config>): void {
  runtime.log("=== isArray — force consistent array shape ===");
  runtime.log("Alternative to: DOMParser querySelectorAll() / xml2js explicitArray:true — both always return a collection");
  runtime.log("");

  const single = `<feed><entry id="1"><title>Only Entry</title></entry></feed>`;
  const multi  = `<feed><entry id="1"><title>A</title></entry><entry id="2"><title>B</title></entry></feed>`;

  const opts   = { ignoreAttributes: false, attributeNamePrefix: "@_" };
  const base   = new XMLParser(opts);
  const forced = new XMLParser({ ...opts, isArray: (name: string) => name === "entry" });

  const singleBase   = (base.parse(single) as { feed: { entry: unknown } }).feed.entry;
  const singleForced = (forced.parse(single) as { feed: { entry: unknown[] } }).feed.entry;
  const multiBase    = (base.parse(multi) as { feed: { entry: unknown[] } }).feed.entry;

  runtime.log(`1 entry,  isArray off -> Array.isArray=${Array.isArray(singleBase)} (object — .map() will throw)`);
  runtime.log(`1 entry,  isArray on  -> Array.isArray=${Array.isArray(singleForced)} (always safe)`);
  runtime.log(`2 entries,isArray off -> Array.isArray=${Array.isArray(multiBase)} (naturally an array)`);
}

// ============================================================================
// XML BUILDER — parse -> filter -> rebuild
// Use when a workflow needs to transform or subset XML before passing it on.
// ============================================================================

function demoXMLBuilder(runtime: Runtime<Config>, rawXml: string): void {
  runtime.log("=== XML BUILDER — parse -> filter -> rebuild ===");
  runtime.log("Alternative to: DOM createElement/setAttribute + XMLSerializer, or xml2js.Builder.buildObject()");
  runtime.log("");

  const parser = new XMLParser({
    ignoreAttributes: false,
    attributeNamePrefix: "@_",
    removeNSPrefix: true,
    parseAttributeValue: true,
  });

  const result    = parser.parse(rawXml) as Record<string, unknown>;
  const ecb       = getDataSet(result);
  const rawSeries = ecb.DataSet.Series;
  const series    = Array.isArray(rawSeries) ? rawSeries : [rawSeries];
  const date      = series[0]?.Obs?.ObsDimension["@_value"] ?? "unknown";

  const outputObj = {
    FXRates: {
      "@_source": "ECB",
      "@_date":   date,
      "@_base":   "EUR",
      Rate: series.map((s: SDMXSeries) => ({
        "@_currency": getCurrency(s),
        "@_rate":     s.Obs.ObsValue["@_value"],
      })),
    },
  };

  const builderOpts = { ignoreAttributes: false, attributeNamePrefix: "@_" };

  const compact   = new XMLBuilder(builderOpts).build(outputObj) as string;
  runtime.log(`Compact   (${compact.length} chars): ${compact.slice(0, 120)}...`);

  const formatted = new XMLBuilder({ ...builderOpts, format: true, indentBy: "  " }).build(outputObj) as string;
  runtime.log(`Formatted (${formatted.length} chars):`);
  formatted.split("\n").slice(0, 8).forEach((l: string) => runtime.log(`  ${l}`));
}

// ============================================================================
// Main Workflow Handler
// ============================================================================

export const onCronTrigger = (runtime: Runtime<Config>): string => {
  runtime.log("========================================");
  runtime.log("FAST-XML-PARSER DEMO");
  runtime.log("Alternative to xml2js / DOMParser");
  runtime.log("Compatible with QuickJS / CRE Workflows");
  runtime.log("========================================");
  runtime.log("");

  const httpClient = new HTTPClient();
  const rawXml = httpClient
    .sendRequest(runtime, fetchECBRates, consensusIdenticalAggregation<string>())()
    .result();

  demoValidation(runtime, rawXml);
  runtime.log("");

  demoParsing(runtime, rawXml);
  runtime.log("");

  demoIsArray(runtime);
  runtime.log("");

  demoXMLBuilder(runtime, rawXml);
  runtime.log("");

  runtime.log("========================================");
  runtime.log("DEMO COMPLETE");
  runtime.log("========================================");

  return "XML parsing demo completed successfully";
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
