#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")" && pwd)"
export GOCACHE="${GOCACHE:-/tmp/cre-auth-template-go-build-cache}"
export GOTMPDIR="${GOTMPDIR:-/tmp/cre-auth-template-go-tmp}"
export GOMODCACHE="${GOMODCACHE:-/tmp/cre-auth-template-go-mod-cache}"
mkdir -p "$GOCACHE" "$GOTMPDIR" "$GOMODCACHE"

echo "===== GOFMT ====="
if [ -n "$(gofmt -l "$ROOT/gateway")" ]; then
  echo "GOFMT=FAIL"
  gofmt -l "$ROOT/gateway"
  exit 1
fi
echo "GOFMT=PASS"

echo
echo "===== GO TEST ====="
(cd "$ROOT/gateway" && go test -count=1 ./...)

echo
echo "===== GO RACE ====="
(cd "$ROOT/gateway" && go test -count=1 -race ./internal/authorization ./internal/httpapi ./internal/store)

echo
echo "===== GO BUILD ====="
(cd "$ROOT/gateway" && go build -o /tmp/cre-confidential-authorization-gateway ./cmd/gateway)
echo "GO_BUILD=PASS"

if command -v bun >/dev/null 2>&1; then
  echo
  echo "===== WORKFLOW INSTALL ====="
  (cd "$ROOT/external-authorization-gateway-ts" && bun install --frozen-lockfile 2>/dev/null || bun install)
  echo
  echo "===== WORKFLOW TYPECHECK ====="
  (cd "$ROOT/external-authorization-gateway-ts" && bun run typecheck)
  echo
  echo "===== WORKFLOW TEST ====="
  (cd "$ROOT/external-authorization-gateway-ts" && bun test)
else
  echo
  echo "WORKFLOW_CHECK=SKIP (bun not installed)"
fi

echo
echo "LOCAL_TEMPLATE_VERIFICATION=PASS"
