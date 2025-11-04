<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE&message=Home&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE&message=Docs&color=blue)](https://docs.chain.link/cre)

</div>

# Multi-Chain Token Manager - CRE Template

A template multi-chain token manager that maximize lending yields via automated cross-chain rebalancing using the **Chainlink Runtime Environment (CRE)** integrated with the **Cross-Chain Interoperability Protocol (CCIP)**.

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

This template provides an end-to-end staring point for writing your own multi-chain token manager that maximizes supply APY by rebalancing tokens cross-chain with the **Chainlink Runtime Environment (CRE)** via the **Cross-Chain Interoperability Protocol (CCIP)**.

The template consists of two components:
- [Contracts](./contracts/README.md) (deployed on multiple chains)
  - [MockPool](./contracts/src/multi-chain-token-manager/MockPool.sol) contract that mimics an AAVE liquidity pool
  - [ProtocolSmartWallet](./contracts/src/multi-chain-token-manager/ProtocolSmartWallet.sol) contract that manages lending positions
- CRE Workflow that monitors supply APY of a target asset on each chain and automatically rebalances lending positions (to maximize supply APY)
  - [Golang](./workflow-go/workflow/workflow.go) workflow targeting the **Golang CRE SDK**
  - [Typescript](./workflow-ts/workflow/main.ts) workflow targeting the **Typescript CRE SDK**

**Key Technologies:**
- **CRE (Chainlink Runtime Environment)** - Orchestrates workflow with DON consensus
- **CCIP (Cross-Chain Interoperability Protocol)** - Secure token bridging with instructions

<img width="1600" height="900" alt="mctm" src="https://github.com/user-attachments/assets/591441d3-a41f-4ee9-b354-f2d23cf1e329" />

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
2. **Contracts are examples** - Write your own audited contracts for your use case
3. **Use your own RPC for stability** - For stable deployment and chainwrite operations it is advised to use your own private RPCs
4. **Secrets hygiene** – Keep real secrets out of version control; use secure secret managers for `.env` values.

---

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.
