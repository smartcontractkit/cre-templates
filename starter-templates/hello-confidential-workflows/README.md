# Hello Confidential Workflows

Run part of a workflow inside a TEE (Trusted Execution Environment) so that computation over sensitive data — Vault DON secrets such as API keys, the request and response payloads of HTTP calls made from the enclave, and other intermediate data in the computation — is kept confidential from node operators.

Note that a confidential workflow, despite running inside the enclave, is part of the binary the Workflow DON provides to the enclave. That binary — including the logic to be executed in the enclave — is revealed; what the enclave keeps confidential is the data that logic computes over.

## Available Languages

| Language | Directory |
|----------|-----------|
| TypeScript | [hello-confidential-workflows-ts](./hello-confidential-workflows-ts) |
| Go | [hello-confidential-workflows-go](./hello-confidential-workflows-go) |

Both implementations perform the same four steps and return the same result for the same input.
