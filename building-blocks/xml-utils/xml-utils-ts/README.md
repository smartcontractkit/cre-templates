# XML Utils - CRE Building Block (TypeScript)

**⚠️ DISCLAIMER**

This tutorial represents an educational example to use a Chainlink system, product, or service and is provided to demonstrate how to interact with Chainlink's systems, products, and services to integrate them into your own. This template is provided "AS IS" and "AS AVAILABLE" without warranties of any kind, it has not been audited, and it may be missing key checks or error handling to make the usage of the system, product or service more clear. Do not use the code in this example in a production environment without completing your own audits and application of best practices. Neither Chainlink Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for unintended outputs that are generated due to errors in code.

---

This building block demonstrates how to parse, validate, and build XML in CRE TypeScript workflows. Since the CRE TypeScript SDK runs on QuickJS (a lightweight JavaScript engine), the standard browser `DOMParser` and Node.js `xml2js` are not available. This template shows how to use **fast-xml-parser** as a drop-in alternative.

## Why fast-xml-parser?

[fast-xml-parser](https://github.com/NaturalIntelligence/fast-xml-parser) is:
- **Pure JavaScript** - No native dependencies, works in any JS environment including QuickJS
- **Bidirectional** - Parse XML to JS objects and build XML from JS objects
- **Attribute-aware** - Exposes XML attributes as typed properties, not just text nodes
- **Namespace-aware** - Can strip namespace prefixes from element names automatically
- **Standards-compliant** - Handles real-world XML including SDMX, RSS, Atom, and custom schemas

## Features Demonstrated

This workflow fetches live ECB EUR FX reference rates from the [ECB Data API](https://data-api.ecb.europa.eu) (SDMX format) and demonstrates:

| Demo | Description |
|------|-------------|
| **Validation** | Validate XML before parsing — get structured error info instead of a thrown exception |
| **Parsing** | Parse XML with attributes, namespace stripping, and automatic type casting |
| **isArray** | Force consistent array shape for elements that may appear once or many times |
| **XMLBuilder** | Reconstruct a clean, flat XML document from a parsed and filtered object |

---

## Node.js / Browser API to fast-xml-parser Mapping

Use this table to migrate from `DOMParser` or `xml2js` to fast-xml-parser:

| DOMParser / xml2js | fast-xml-parser | Code Reference |
|--------------------|-----------------|----------------|
| `try { new DOMParser().parseFromString(xml, ...) }` | `XMLValidator.validate(xml)` | [main.ts:84](./workflow/main.ts#L84) |
| `new DOMParser().parseFromString(xml, 'text/xml')` | `new XMLParser(opts).parse(xml)` | [main.ts:102](./workflow/main.ts#L102) |
| `element.getAttribute('rate')` (always string) | `parseAttributeValue: true` | [main.ts:106](./workflow/main.ts#L106) |
| `element.localName` (strips namespace manually) | `removeNSPrefix: true` | [main.ts:105](./workflow/main.ts#L105) |
| `querySelectorAll('entry')` (always a NodeList) | `isArray: (name) => name === 'entry'` | [main.ts:142](./workflow/main.ts#L142) |
| `xml2js` `explicitArray: true` (default) | `isArray: (name) => ...` | [main.ts:142](./workflow/main.ts#L142) |
| `document.createElement` + `XMLSerializer` | `new XMLBuilder(opts).build(obj)` | [main.ts:190](./workflow/main.ts#L190) |
| `xml2js.Builder.buildObject(obj)` | `new XMLBuilder(opts).build(obj)` | [main.ts:190](./workflow/main.ts#L190) |

---

## Dependencies

```json
{
  "fast-xml-parser": "^4.4.0"
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
   cd building-blocks/xml-utils/xml-utils-ts/workflow
   bun install
   ```

---

## Running the Workflow

### Simulate the workflow

From the project root directory (`xml-utils-ts`):

```bash
cre workflow simulate workflow
```

---

## Example Output

When the workflow runs, it logs the output of each XML operation:

```
========================================
FAST-XML-PARSER DEMO
Alternative to xml2js / DOMParser
Compatible with QuickJS / CRE Workflows
========================================


=== VALIDATION — XMLValidator.validate() ===
Alternative to: try/catch around DOMParser.parseFromString()

Validating ECB XML...
ECB XML valid: true

Validating malformed XML...
Malformed XML — code: InvalidAttr, line: 1
Attributes for 'entry' have open quote.

=== PARSING — attributes + namespace removal + type casting ===
Alternative to: xml2js / DOMParser with manual attribute handling and type casting

Date:       2026-04-17
Series:     7

1 EUR = 1.6438 AUD
1 EUR = 1.6129 CAD
1 EUR = 0.9231 CHF
1 EUR = 8.0483 CNY
1 EUR = 0.87168 GBP
1 EUR = 187.72 JPY
1 EUR = 1.1797 USD

=== isArray — force consistent array shape ===
Alternative to: DOMParser querySelectorAll() / xml2js explicitArray:true — both always return a collection

1 entry,  isArray off -> Array.isArray=false (object — .map() will throw)
1 entry,  isArray on  -> Array.isArray=true (always safe)
2 entries,isArray off -> Array.isArray=true (naturally an array)

=== XML BUILDER — parse -> filter -> rebuild ===
Alternative to: DOM createElement/setAttribute + XMLSerializer, or xml2js.Builder.buildObject()

Compact   (356 chars): <FXRates source="ECB" date="2026-04-17" base="EUR"><Rate currency="AUD" rate="1.6438"></Rate><Rate currency="CAD" rate="...
Formatted (379 chars):
<FXRates source="ECB" date="2026-04-17" base="EUR">
<Rate currency="AUD" rate="1.6438"></Rate>
<Rate currency="CAD" rate="1.6129"></Rate>
<Rate currency="CHF" rate="0.9231"></Rate>
<Rate currency="CNY" rate="8.0483"></Rate>
<Rate currency="GBP" rate="0.87168"></Rate>
<Rate currency="JPY" rate="187.72"></Rate>
<Rate currency="USD" rate="1.1797"></Rate>

========================================
DEMO COMPLETE
========================================
```

---

## Use Cases

### 1. Parsing Financial Data Feeds
Consume XML-native price feeds from central banks, financial data providers, or government agencies:
```typescript
const parser = new XMLParser({
  ignoreAttributes: false,
  attributeNamePrefix: "@_",
  removeNSPrefix: true,
  parseAttributeValue: true,
});
const result = parser.parse(rawXml) as Record<string, unknown>;
```

### 2. Validating External XML Before Use
Check XML validity before parsing to surface structured errors instead of silent failures:
```typescript
const valid = XMLValidator.validate(rawXml);
if (valid !== true) {
  throw new Error(`Invalid XML at line ${valid.err.line}: ${valid.err.msg}`);
}
```

### 3. Transforming and Subsetting XML
Parse a verbose upstream XML response, filter it down, and rebuild a clean output:
```typescript
const parser  = new XMLParser({ ignoreAttributes: false, attributeNamePrefix: "@_" });
const builder = new XMLBuilder({ ignoreAttributes: false, attributeNamePrefix: "@_", format: true });

const parsed   = parser.parse(rawXml);
const filtered = { FXRates: { Rate: parsed.rates.filter(...) } };
const output   = builder.build(filtered);
```

---

## Reference Documentation

- [CRE Documentation](https://docs.chain.link/cre)
- [fast-xml-parser GitHub](https://github.com/NaturalIntelligence/fast-xml-parser)
- [fast-xml-parser docs](https://naturalintelligence.github.io/fast-xml-parser/)
- [ECB Data API](https://data-api.ecb.europa.eu)

---

## License

MIT - see the repository's [LICENSE](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE).
