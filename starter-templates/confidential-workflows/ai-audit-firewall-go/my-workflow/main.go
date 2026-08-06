//go:build wasip1

package main

import (
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

// ParseConfig (not cre.ParseJSON) so an empty config payload falls back to
// DefaultConfig instead of failing to unmarshal — see workflow.go.
func main() {
	wasm.NewRunner(ParseConfig).Run(InitWorkflow)
}
