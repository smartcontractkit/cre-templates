# Hello Confidential Workflows

Quickstart confidential workflow. Run a handler's callback inside a secure enclave: fetch a secret from the Vault DON, call an API from inside the enclave, execute decision logic over the confidential data such as Vault DON secrets or HTTP response payloads, then cross back to the Workflow DON for consensus and DON capability calls. 

Note that a confidential workflow, despite running inside the enclave, is part of the binary the Workflow DON provides to the enclave. That binary — including the logic to be executed in the enclave — is revealed; what the enclave keeps confidential is the data that logic computes over.

## Available Languages

| Language | Directory |
|----------|-----------|
| TypeScript | [hello-confidential-workflows-ts](./hello-confidential-workflows-ts) |
| Go | [hello-confidential-workflows-go](./hello-confidential-workflows-go) |

Both implementations perform the same four steps and return the same result for the same input.
