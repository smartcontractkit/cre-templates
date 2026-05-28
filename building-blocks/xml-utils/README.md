<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE&message=Home&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE&message=Docs&color=blue)](https://docs.chain.link/cre)

</div>

# XML Utils - CRE Building Block

This building block demonstrates how to parse, validate, and build XML in CRE workflows using **fast-xml-parser** as an alternative to the browser `DOMParser` and Node.js `xml2js`.

## The Problem

The CRE TypeScript SDK runs on **QuickJS**, a lightweight JavaScript engine that does not support browser or Node.js XML APIs. This means the standard `DOMParser` and `xml2js` modules are **not available** in CRE workflows.

## The Solution

[fast-xml-parser](https://github.com/NaturalIntelligence/fast-xml-parser) is a pure JavaScript XML library that works in QuickJS. It provides:

- **Pure JavaScript** - No native dependencies, works in any JS environment
- **Bidirectional** - Parse XML to JS objects and build XML from JS objects
- **Attribute-aware** - Exposes XML attributes as typed properties
- **Namespace-aware** - Can strip namespace prefixes from element names automatically
- **Standards-compliant** - Handles real-world XML including SDMX, RSS, Atom, and custom schemas

---

## What's Covered

This template demonstrates:

| Operation | API | Use Case |
|-----------|-----|----------|
| **Validation** | `XMLValidator.validate()` | Validate XML before parsing; get structured error info instead of thrown exceptions |
| **Parsing** | `XMLParser.parse()` | Parse XML with attributes, namespace stripping, and automatic type casting |
| **Array Consistency** | `isArray` option | Force single child elements to parse as arrays instead of objects |
| **XML Building** | `XMLBuilder.build()` | Reconstruct or transform XML from a parsed and filtered object |

## Get Started

- **TypeScript**: See the [TypeScript README](./xml-utils-ts/README.md) for detailed setup, API mapping table, example output, and use cases.

## Quick Example

```typescript
// Instead of DOMParser:
// const doc = new DOMParser().parseFromString(xml, 'text/xml');
// const rate = doc.querySelector('rate')?.getAttribute('value');

// Use fast-xml-parser:
import { XMLParser } from "fast-xml-parser";

const parser = new XMLParser({
  ignoreAttributes: false,
  attributeNamePrefix: "@_",
  parseAttributeValue: true,
});
const result = parser.parse(xml);
const rate = result.root.rate["@_value"];
```

## Reference Documentation

- [fast-xml-parser GitHub](https://github.com/NaturalIntelligence/fast-xml-parser)
- [fast-xml-parser Documentation](https://naturalintelligence.github.io/fast-xml-parser/)
- [CRE Documentation](https://docs.chain.link/cre)
