<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE\&message=Home\&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE\&message=Docs\&color=blue)](https://docs.chain.link/cre)

</div>

# Key-Value Store (AWS S3) - CRE Building Block (TypeScript)

A minimal example that, on a cron schedule, reads a value from an **AWS S3 object**, increments it, and writes it back—using **CRE (Chainlink Runtime Environment)** with SigV4-signed HTTP requests.

---

**⚠️ DISCLAIMER**

This template is an educational example to demonstrate how to interact with Chainlink systems, products, and services. It is provided **"AS IS"** and **"AS AVAILABLE"** without warranties of any kind, has **not** been audited, and may omit checks or error handling for clarity. **Do not use this code in production** without performing your own audits and applying best practices. Neither Chainlink Labs, the Chainlink Foundation, nor Chainlink node operators are responsible for unintended outputs generated due to errors in code.

---

## What This Example Does

The workflow:
- Retrieves `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` from CRE Secrets
- Signs each request with AWS **SigV4** using [`@noble/hashes`](https://github.com/paulmillr/noble-hashes) (pure JS — no AWS SDK needed)
- Reads the S3 object (initializes to `0` if missing)
- Aggregates the **current value** across nodes (median), increments once, and writes the **agreed next value** back

### Why not use the AWS SDK?

CRE TypeScript workflows run on **QuickJS compiled to WASM**, which does not support Node.js APIs. The AWS SDK v3 (`@aws-sdk/client-s3`) depends on Node.js `http`, `crypto`, and `stream` modules. Instead, this example implements SigV4 signing manually using `@noble/hashes`, which is pure JavaScript and fully compatible with the CRE runtime.

---

## Setup & Run

### 1) Install dependencies

```bash
cd my-workflow
bun install
```

### 2) Provide AWS credentials as secrets

Add the following secrets to your CRE environment:

```
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
```

For local testing, export them as environment variables:

```bash
export AWS_ACCESS_KEY_ID_ENV="your-access-key-id"
export AWS_SECRET_ACCESS_KEY_ENV="your-secret-access-key"
```

### 3) Configure the workflow

Update `my-workflow/config/config.staging.json` with your S3 details:

```json
{
  "schedule": "* * */1 * * *",
  "aws_region": "us-east-1",
  "s3_bucket": "my-bucket-name",
  "s3_key": "cre-counter.txt"
}
```

* `schedule` uses a **6-field** cron expression with seconds (e.g., `* * */1 * * *` runs every hour).
* `aws_region` is the AWS region of your bucket.
* `s3_bucket` is the bucket name.
* `s3_key` is the object path that stores the counter.

### 4) Run a local simulation

From the `kv-store-ts` directory:

```bash
cre workflow simulate my-workflow
```

You should see output similar to:

```
Workflow compiled
2025-11-03T16:31:45Z [SIMULATION] Simulator Initialized

2025-11-03T16:31:45Z [SIMULATION] Running trigger trigger=cron-trigger@1.0.0
2025-11-03T16:31:45Z [USER LOG] Cron trigger fired. Fetching AWS credentials...
2025-11-03T16:31:45Z [USER LOG] AWS credentials fetched. Performing consensus read, then write.
2025-11-03T16:31:45Z [USER LOG] Consensus old value computed. Incrementing. old=2 new=3
2025-11-03T16:31:45Z [USER LOG] Workflow finished successfully. old=2 new=3

Workflow Simulation Result:
 {
  "oldValue": 2,
  "newValue": "3"
}

2025-11-03T16:31:45Z [SIMULATION] Execution finished signal received
2025-11-03T16:31:45Z [SIMULATION] Skipping WorkflowEngineV2
```
