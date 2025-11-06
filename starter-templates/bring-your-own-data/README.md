<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE&message=Home&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE&message=Docs&color=blue)](https://docs.chain.link/cre)

</div>

# Bring Your Own Data - CRE Templates

A template for bringing your own off-chain data on chain with the **Chainlink Runtime Environment (CRE)**.

---

**⚠️ DISCLAIMER**

This tutorial represents an educational example to use a Chainlink system, product, or service and is provided to demonstrate how to interact with Chainlink's systems, products, and services to integrate them into your own. This template is provided "AS IS" and "AS AVAILABLE" without warranties of any kind, it has not been audited, and it may be missing key checks or error handling to make the usage of the system, product or service more clear. Do not use the code in this example in a production environment without completing your own audits and application of best practices. Neither Chainlink Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for unintended outputs that are generated due to errors in code.

---

## Table of Contents

* [What This Template Does](#what-this-template-does)
* [Getting Started](#getting-started)
* [Security Considerations](#security-considerations)
* [License](#license)

---

## What This Template Does

This template provides an end-to-end staring point for bringing your own off-chain Proof-of-Reserve (PoR) or Net-Asset-Value (NAV) data on-chain with the **Chainlink Runtime Environment (CRE)**.

The template consists of two components:
- [Contracts](./contracts/README.md) (deployed on multiple chains)
  - [DataFeedsCache](https://github.com/smartcontractkit/chainlink-evm/blob/88d90433a15f1c34bb5fabc29be192400fad396c/contracts/src/v0.8/data-feeds/DataFeedsCache.sol) contract that receives data on-chain
  - [DecimalAggregatorProxy](./contracts/src/nav/DecimalAggregatorProxy.sol) proxy contract for PoR data stored in the DataFeedsCache contract
  - [BundleAggregatorProxy](https://github.com/smartcontractkit/chainlink-evm/blob/88d90433a15f1c34bb5fabc29be192400fad396c/contracts/src/v0.8/data-feeds/BundleAggregatorProxy.sol) proxy contract for NAV data stored in the DataFeedsCache contract
- CRE Workflows
  - [Golang](./workflow-go/README.md) workflows targeting the **Golang CRE SDK**
    - [PoR Workflow](./workflow-go/por/workflow.go) Proof-of-Reserve (PoR)
    - [NAV Workflow](./workflow-go/nav/workflow.go) Net-Asset-Value (NAV)
  - [Typescript](./workflow-ts/README.md) workflows targeting the **Typescript CRE SDK**
    - [PoR Workflow](./workflow-ts/por/main.ts) Proof-of-Reserve (PoR)
    - [NAV Workflow](./workflow-ts/nav/main.ts) Net-Asset-Value (NAV)

**Key Technologies:**
- **CRE (Chainlink Runtime Environment)** - Orchestrates workflows with DON consensus

<img width="1600" height="900" alt="por" src="https://github.com/user-attachments/assets/39f811b5-7afe-4944-822f-cb2d97af7156" />

---

## Getting Started

You can build with the **CRE SDK** in either **Golang** or **TypeScript**.

### Golang SDK

Start [here](./workflow-go/README.md)

### TypeScript SDK

Start [here](./workflow-ts/README.md)

---

## Security Considerations

**⚠️ Important Notes:**

1. **This is a demo project** - Not production-ready
2. **DecimalAggregatorProxy contract is demo contract** - Not audited and not production-ready
3. **Use your own RPC for stability** - For stable deployment and chainwrite operations it is advised to use your own private RPCs
4. **Secrets hygiene** – Keep real secrets out of version control; use secure secret managers for `.env` values.

---

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.
