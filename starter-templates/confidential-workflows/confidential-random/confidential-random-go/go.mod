module confidential-random-go

go 1.25.3

// cre-sdk-go v1.22.0-beta.0 is the first release whose TeeRuntime.Rand() uses
// TEE-local (node-mode) randomness. It is tagged at commit
// c9308ff292c98d4067518ec9387d2ec89b446560; earlier releases returned DON
// randomness from TeeRuntime.Rand().
//
// The capabilities/scheduler/cron submodule has not had a tagged release since
// that commit, so it is pinned by pseudo-version.
require (
	github.com/smartcontractkit/cre-sdk-go v1.22.0-beta.0
	github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron v1.3.1-0.20260806122424-24d22c8dc5bd
)

require (
	github.com/ProjectZKM/Ziren/crates/go-runtime/zkvm_runtime v0.0.0-20251001021608-1fe7b43fc4d6 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/ethereum/go-ethereum v1.17.2 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/holiman/uint256 v1.3.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/smartcontractkit/chainlink-protos/cre/go v0.0.0-20260911202725-254c3ed37098 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	golang.org/x/sys v0.40.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
