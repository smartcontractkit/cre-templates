<div style="text-align:center" align="center">
    <a href="https://chain.link" target="_blank">
        <img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
    </a>

[![License](https://img.shields.io/badge/license-MIT-blue)](https://github.com/smartcontractkit/cre-templates/blob/main/LICENSE)
[![CRE Home](https://img.shields.io/static/v1?label=CRE&message=Home&color=blue)](https://chain.link/chainlink-runtime-environment)
[![CRE Documentation](https://img.shields.io/static/v1?label=CRE&message=Docs&color=blue)](https://docs.chain.link/cre)

</div>

# Key Value Storage using AWS S3 - CRE Building Blocks

This example workflow demonstrates how to use the Chainlink Runtime Environment (CRE) Go SDK to create a simple, stateful counter that uses an AWS S3 object as its persistent storage.

## Description

This workflow runs on a schedule defined by a cron expression. On each run, it:

1. Securely fetches AWS credentials (`AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`) from the CRE Secrets API.
2. Authenticates with AWS S3 by generating a SigV4 signature using the consensus-based `nodeRuntime.Now()` timestamp.
3. Reads the current value from a specified S3 object (e.g., `cre-counter.txt`).
4. If the object doesn't exist or is empty, it initializes the value to "0".
5. Increments the numeric value by 1.
6. Writes the new value back to the S3 object, overwriting the previous content.

This example showcases CRE capabilities, including using the `http` capability, fetching secrets, and performing stateful business logic using a 3rd party service.

## How to Use

### 1. Configure the Workflow

Open the `my-go-workflow/config.json` file and set the S3 object details.

```json
{
  "schedule": "* * */1 * * *",
  "aws_region": "my-aws-region",
  "s3_bucket": "product-release-bucket",
  "s3_key": "cre-counter.txt"
}
```

  * **`schedule`**: The cron expression for how often to run (e.g., "`* * */1 * * *`" for every hour).
  * **`aws_region`**: The AWS region where your S3 bucket is located.
  * **`s3_bucket`**: The name of your S3 bucket.
  * **`s3_key`**: The path to the object (file) that will store the counter.
