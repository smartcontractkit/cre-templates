<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE\&message=Home\&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE\&message=Docs\&color=blue)](https://docs.chain.link/cre)

</div>

# Key-Value Store (AWS S3) - CRE Building Blocks

A minimal example that, on a cron schedule, reads a value from an **AWS S3 object**, increments it, and writes it back—using **CRE (Chainlink Runtime Environment)** with SigV4-signed HTTP requests.

The workflow:
- Retrieves `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` from CRE Secrets
- Signs the request with AWS **SigV4** (timestamp sourced from the consensus runtime)
- Reads the S3 object (initializes to `0` if missing)
- Aggregates the **current value** across nodes (median), increments once, and writes the **agreed next value** back

---

## Setup & Run

### 1) Provide AWS credentials as secrets

Add the following secrets to your CRE environment:

```

AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY

````

(For local testing, ensure these are available to the simulator via your environment variables.)

### 2) Configure the workflow

Create/update `my-workflow/config.json` with your S3 details:

```json
{
  "schedule": "* * */1 * * *",
  "aws_region": "my-aws-region",
  "s3_bucket": "product-release-bucket",
  "s3_key": "cre-counter.txt"
}
````

* `schedule` uses a **6-field** cron expression with seconds (e.g., `* * */1 * * *` runs every hour).
* `aws_region` is the AWS region of your bucket.
* `s3_bucket` is the bucket name.
* `s3_key` is the object path that stores the counter.

### 3) Run a local simulation

From your project root:

```bash
cre workflow simulate my-workflow
```

You should see output similar to:

```
Workflow compiled
2025-11-03T16:31:45Z [SIMULATION] Simulator Initialized

2025-11-03T16:31:45Z [SIMULATION] Running trigger trigger=cron-trigger@1.0.0
2025-11-03T16:31:45Z [USER LOG] msg="Cron trigger fired. Fetching AWS credentials..."
2025-11-03T16:31:45Z [USER LOG] msg="AWS credentials fetched. Performing consensus read, then write."
2025-11-03T16:31:45Z [USER LOG] msg="Consensus old value computed. Incrementing." old=2 new=3
2025-11-03T16:31:45Z [USER LOG] msg="Workflow finished successfully." old=2 new=3

Workflow Simulation Result:
 {
  "NewValue": "3",
  "OldValue": "2"
}

2025-11-03T16:31:45Z [SIMULATION] Execution finished signal received
2025-11-03T16:31:45Z [SIMULATION] Skipping WorkflowEngineV2
```