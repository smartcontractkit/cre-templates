#!/bin/bash
# scripts/check-templates.sh
# Validates and compiles changed CRE templates (TypeScript + Go).
#
# Usage (CI):
#   BASE_REF=main ./scripts/check-templates.sh
#
# Usage (local — specific templates):
#   CHANGED_TEMPLATES="building-blocks/indexer-block-trigger/block-trigger-ts" \
#     ./scripts/check-templates.sh --verbose
#
# Usage (local — all templates):
#   ./scripts/check-templates.sh --verbose
#
# Environment variables:
#   BASE_REF            Git base branch to diff against (e.g. "main").
#                       Used to find changed templates in a PR.
#   CHANGED_TEMPLATES   Newline-separated list of template root dirs
#                       relative to the repo root. Overrides BASE_REF detection.
#   VERBOSE             Set to 1 for verbose output (same as --verbose/-v).
#   RESULTS_FILE        Path to write structured JSON results.
#                       Default: /tmp/template-results.json

# --------------------------------------------------------------------------
# Flags

VERBOSE=false
for arg in "$@"; do
  case "$arg" in
    -v|--verbose) VERBOSE=true ;;
  esac
done
[[ "${VERBOSE:-0}" == "1" ]] && VERBOSE=true

# --------------------------------------------------------------------------
# Logging helpers

info()  { echo "$@"; }
vlog()  { $VERBOSE && echo "$@" || true; }

# Capture stdout+stderr of a command into a variable, streaming in verbose mode.
# Usage: run_captured <output_var> <cmd> [args...]
run_captured() {
  local _outvar="$1"; shift
  local _tmpfile; _tmpfile=$(mktemp)
  if $VERBOSE; then
    "$@" 2>&1 | tee "$_tmpfile"; local _rc="${PIPESTATUS[0]}"
  else
    "$@" > "$_tmpfile" 2>&1; local _rc=$?
  fi
  printf -v "$_outvar" '%s' "$(cat "$_tmpfile")"
  rm -f "$_tmpfile"
  return "$_rc"
}

# --------------------------------------------------------------------------
# Setup

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
RESULTS_FILE="${RESULTS_FILE:-/tmp/template-results.json}"
JSONL_FILE=$(mktemp /tmp/cre-check-XXXXXX.jsonl)

# Lockfile tracking — populated by check_ts_workflow before each npm install.
# Backed-up lockfiles are restored; generated ones are deleted on exit.
LOCKFILE_BACKUPS=()    # "backup_path:original_path" pairs
GENERATED_LOCKFILES=() # paths to delete on exit

cleanup() {
  rm -f "$JSONL_FILE" "${_FAIL_OUT:-}" 2>/dev/null || true

  # Restore lockfiles that existed before we ran
  local entry backup original
  for entry in "${LOCKFILE_BACKUPS[@]+"${LOCKFILE_BACKUPS[@]}"}"; do
    backup="${entry%%:*}"
    original="${entry##*:}"
    [[ -f "$backup" ]] && mv -f "$backup" "$original" 2>/dev/null || true
  done

  # Remove lockfiles that we generated fresh
  local f
  for f in "${GENERATED_LOCKFILES[@]+"${GENERATED_LOCKFILES[@]}"}"; do
    rm -f "$f" 2>/dev/null || true
  done
}
trap cleanup EXIT INT TERM

# --------------------------------------------------------------------------
# YAML helpers (require python3 + PyYAML, always available on ubuntu-latest)

# Print the value of a top-level YAML field to stdout.
yaml_field() {
  local yaml_file="$1" field="$2"
  python3 -c "
import yaml, sys
with open(sys.argv[1]) as f:
    d = yaml.safe_load(f) or {}
val = d.get(sys.argv[2])
if val is not None:
    print(str(val).strip())
" "$yaml_file" "$field"
}

# Print workflow dirs (one per line) from .cre/template.yaml.
yaml_workflow_dirs() {
  local yaml_file="$1"
  python3 -c "
import yaml, sys
with open(sys.argv[1]) as f:
    d = yaml.safe_load(f) or {}
for w in d.get('workflows', []):
    print(w['dir'])
" "$yaml_file"
}

# Validate a YAML file. For kind='template', also checks required fields.
# Prints an error message to stdout and returns 1 on failure.
validate_yaml() {
  local yaml_file="$1" kind="${2:-}"
  python3 -c "
import yaml, sys
kind = sys.argv[2] if len(sys.argv) > 2 else ''
try:
    with open(sys.argv[1]) as f:
        d = yaml.safe_load(f)
    if d is None:
        print('empty YAML file')
        sys.exit(1)
    if kind == 'template':
        required = ['kind', 'id', 'title', 'language', 'category', 'workflows']
        missing = [f for f in required if f not in d]
        if missing:
            print(f'missing required fields: {missing}')
            sys.exit(1)
except yaml.YAMLError as e:
    print(f'YAML parse error: {e}')
    sys.exit(1)
" "$yaml_file" "$kind" 2>&1
}

# --------------------------------------------------------------------------
# Results recording
#
# Each result is one JSON line in JSONL_FILE.
# Failure output is read from _FAIL_OUT (a temp file written by run_captured).

_FAIL_OUT=$(mktemp /tmp/cre-failout-XXXXXX.txt)

# record_pass <name> <language> <sdk_range> <sdk_resolved>
record_pass() {
  _write_result "$1" "$2" "$3" "$4" "pass" "" ""
}

# record_fail <name> <language> <sdk_range> <sdk_resolved> <step> <output_file>
record_fail() {
  _write_result "$1" "$2" "$3" "$4" "fail" "$5" "${6:-}"
}

_write_result() {
  local name="$1" lang="$2" sdk_range="$3" sdk_resolved="$4"
  local status="$5" step="$6" output_file="$7"
  python3 -c "
import json, sys, os
name, lang, sdk_range, sdk_resolved, status, step = sys.argv[1:7]
output_file = sys.argv[7] if len(sys.argv) > 7 else ''
output = None
if output_file and os.path.isfile(output_file):
    with open(output_file) as f:
        text = f.read(3000).strip()
    if text:
        output = text
print(json.dumps({
    'name':         name,
    'language':     lang,
    'sdk_range':    sdk_range,
    'sdk_resolved': sdk_resolved,
    'status':       status,
    'failure_step': step or None,
    'failure_output': output,
}))
" "$name" "$lang" "$sdk_range" "$sdk_resolved" "$status" "$step" "$output_file" >> "$JSONL_FILE"
}

finalize_results() {
  python3 -c "
import json, sys, os
jsonl = sys.argv[1]
out   = sys.argv[2]
templates = []
if os.path.isfile(jsonl):
    with open(jsonl) as f:
        for line in f:
            line = line.strip()
            if line:
                try:
                    templates.append(json.loads(line))
                except json.JSONDecodeError:
                    pass
passed = sum(1 for t in templates if t['status'] == 'pass')
failed = len(templates) - passed
with open(out, 'w') as f:
    json.dump({'templates': templates, 'total': len(templates),
               'passed': passed, 'failed': failed}, f, indent=2)
" "$JSONL_FILE" "$RESULTS_FILE"
}

# --------------------------------------------------------------------------
# Detect changed template roots

detect_template_roots() {
  if [[ -n "${CHANGED_TEMPLATES:-}" ]]; then
    # Use explicitly provided list (newline-separated)
    while IFS= read -r dir; do
      [[ -n "$dir" && -f "$REPO_ROOT/$dir/.cre/template.yaml" ]] && echo "$dir"
    done <<< "$CHANGED_TEMPLATES"
    return
  fi

  if [[ -n "${BASE_REF:-}" ]]; then
    # Derive from git diff against the PR base branch
    local changed_files
    changed_files=$(git -C "$REPO_ROOT" diff --name-only "origin/${BASE_REF}...HEAD" 2>/dev/null || true)
    if [[ -z "$changed_files" ]]; then
      return
    fi

    while IFS= read -r file; do
      [[ -z "$file" ]] && continue
      local check_dir="$REPO_ROOT/$(dirname "$file")"
      while [[ "$check_dir" != "$REPO_ROOT" && "$check_dir" != "/" ]]; do
        if [[ -f "$check_dir/.cre/template.yaml" ]]; then
          echo "${check_dir#"$REPO_ROOT/"}"
          break
        fi
        check_dir="$(dirname "$check_dir")"
      done
    done <<< "$changed_files"
    return
  fi

  # Fallback: all templates (useful for local dev)
  while IFS= read -r yaml; do
    local template_dir
    template_dir="$(dirname "$(dirname "$yaml")")"
    echo "${template_dir#"$REPO_ROOT/"}"
  done < <(find "$REPO_ROOT" -name "template.yaml" -path "*/.cre/template.yaml" \
             -not -path "*/node_modules/*" 2>/dev/null)
}

# --------------------------------------------------------------------------
# TypeScript workflow check

# check_ts_workflow <template_dir_rel> <workflow_subdir>
check_ts_workflow() {
  local template_dir="$1" workflow_subdir="$2"
  local abs_wf="$REPO_ROOT/$template_dir/$workflow_subdir"
  local pkg_json="$abs_wf/package.json"
  local display="$template_dir / $workflow_subdir"

  if [[ ! -f "$pkg_json" ]]; then
    vlog "    ⚠️  No package.json in $workflow_subdir — skipping"
    return 0
  fi

  info "  📦 $display"

  # Resolve the SDK semver range from package.json
  local sdk_range sdk_resolved="unknown"
  sdk_range=$(node -e "
const p = require('$pkg_json');
const deps = Object.assign({}, p.dependencies, p.devDependencies);
process.stdout.write(deps['@chainlink/cre-sdk'] || 'unknown');
" 2>/dev/null || echo "unknown")

  local _out=""
  cd "$abs_wf"

  # 1. Validate workflow.yaml if present
  if [[ -f "workflow.yaml" ]]; then
    local yaml_err
    if ! yaml_err=$(validate_yaml "workflow.yaml"); then
      info "    ❌ workflow.yaml invalid: $yaml_err"
      printf '%s' "$yaml_err" > "$_FAIL_OUT"
      record_fail "$display" "typescript" "$sdk_range" "$sdk_resolved" "yaml-validation" "$_FAIL_OUT"
      return 0
    fi
    vlog "    ✅ workflow.yaml valid"
  fi

  # 2. npm install — track package-lock.json so we can restore it on exit.
  # Use git to distinguish committed lockfiles (restore) from untracked/absent
  # ones (delete). This correctly handles leftovers from previous script runs.
  local _wf_rel="${abs_wf#"$REPO_ROOT/"}"
  if git -C "$REPO_ROOT" ls-files --error-unmatch "$_wf_rel/package-lock.json" &>/dev/null; then
    # File is committed to git — back up its current content to restore later
    cp package-lock.json package-lock.json.__cre_bak
    LOCKFILE_BACKUPS+=("$(pwd)/package-lock.json.__cre_bak:$(pwd)/package-lock.json")
  else
    # File is not tracked — mark for deletion after the check
    GENERATED_LOCKFILES+=("$(pwd)/package-lock.json")
  fi
  vlog "    Installing dependencies (SDK range: $sdk_range)..."
  if ! run_captured _out npm install --no-audit --fund=false; then
    printf '%s' "$_out" > "$_FAIL_OUT"
    info "    ❌ npm install failed"
    record_fail "$display" "typescript" "$sdk_range" "$sdk_resolved" "npm install" "$_FAIL_OUT"
    return 0
  fi

  # 2b. npm install in template contracts/ when present (generated bindings import viem / cre-sdk).
  local contracts_dir="$REPO_ROOT/$template_dir/contracts"
  if [[ -f "$contracts_dir/package.json" ]]; then
    vlog "    Installing contracts dependencies..."
    cd "$contracts_dir"
    local _contracts_rel="${contracts_dir#"$REPO_ROOT/"}"
    if git -C "$REPO_ROOT" ls-files --error-unmatch "$_contracts_rel/package-lock.json" &>/dev/null; then
      cp package-lock.json package-lock.json.__cre_bak
      LOCKFILE_BACKUPS+=("$(pwd)/package-lock.json.__cre_bak:$(pwd)/package-lock.json")
    else
      GENERATED_LOCKFILES+=("$(pwd)/package-lock.json")
    fi
    if ! run_captured _out npm install --no-audit --fund=false; then
      printf '%s' "$_out" > "$_FAIL_OUT"
      info "    ❌ npm install (contracts) failed"
      record_fail "$display" "typescript" "$sdk_range" "$sdk_resolved" "npm install (contracts)" "$_FAIL_OUT"
      cd "$abs_wf"
      return 0
    fi
    cd "$abs_wf"
  fi

  # 3. Capture resolved SDK version
  if [[ -f "node_modules/@chainlink/cre-sdk/package.json" ]]; then
    sdk_resolved=$(node -e \
      "process.stdout.write(require('./node_modules/@chainlink/cre-sdk/package.json').version)" \
      2>/dev/null || echo "unknown")
  fi
  vlog "    SDK: $sdk_range → resolved $sdk_resolved"

  # 4. Typecheck (if script exists)
  local has_typecheck
  has_typecheck=$(node -e \
    "const s=(require('./package.json').scripts||{}); process.stdout.write(s.typecheck?'yes':'no')" \
    2>/dev/null || echo "no")
  if [[ "$has_typecheck" == "yes" ]]; then
    vlog "    Running typecheck..."
    if ! run_captured _out npm run typecheck --silent; then
      printf '%s' "$_out" > "$_FAIL_OUT"
      info "    ❌ typecheck failed (SDK $sdk_resolved)"
      record_fail "$display" "typescript" "$sdk_range" "$sdk_resolved" "typecheck" "$_FAIL_OUT"
      return 0
    fi
    vlog "    ✅ typecheck passed"
  else
    vlog "    ⚠️  No typecheck script — skipping"
  fi

  # 5. cre-compile (if main.ts exists)
  # Prefer `bun x` so the Bun runtime is used explicitly (cre-compile's bin uses a bun shebang;
  # `npx` on some Linux/npm combinations does not invoke it reliably).
  if [[ -f "main.ts" ]]; then
    vlog "    Running cre-compile..."
    local _compile_cmd
    if command -v bun >/dev/null 2>&1; then
      _compile_cmd=(bun x cre-compile main.ts)
    else
      _compile_cmd=(npx --no cre-compile main.ts)
    fi
    if ! run_captured _out "${_compile_cmd[@]}"; then
      printf '%s' "$_out" > "$_FAIL_OUT"
      info "    ❌ cre-compile failed (SDK $sdk_resolved)"
      record_fail "$display" "typescript" "$sdk_range" "$sdk_resolved" "cre-compile" "$_FAIL_OUT"
      # Clean up any partial build artifacts
      rm -f main.js main.wasm
      return 0
    fi
    rm -f main.js main.wasm
    vlog "    ✅ cre-compile passed"
    info "  ✅ $display (SDK $sdk_resolved)"
  else
    vlog "    ⚠️  No main.ts — typecheck only"
    info "  ✅ $display (SDK $sdk_resolved, typecheck only)"
  fi

  record_pass "$display" "typescript" "$sdk_range" "$sdk_resolved"
}

# --------------------------------------------------------------------------
# Go template check

# check_go_template <template_dir_rel> <workflow_dirs_newline_sep>
check_go_template() {
  local template_dir="$1"
  local abs_template="$REPO_ROOT/$template_dir"
  local display="$template_dir"

  info "  🔧 $display"

  if [[ ! -f "$abs_template/go.mod" ]]; then
    vlog "    ⚠️  No go.mod — skipping"
    return 0
  fi

  # Resolve cre-sdk-go version from go.mod (first require line for the core module)
  local sdk_version
  sdk_version=$(grep -m1 'github\.com/smartcontractkit/cre-sdk-go ' "$abs_template/go.mod" \
    | awk '{print $2}')
  sdk_version="${sdk_version:-unknown}"

  local _out=""
  cd "$abs_template"

  # 1. Validate workflow.yaml files (one per workflow dir)
  local wf_dir yaml_err
  while IFS= read -r wf_dir; do
    [[ -z "$wf_dir" ]] && continue
    local wf_yaml="$wf_dir/workflow.yaml"
    if [[ -f "$wf_yaml" ]]; then
      if ! yaml_err=$(validate_yaml "$wf_yaml"); then
        info "    ❌ $wf_yaml invalid: $yaml_err"
        printf '%s' "$yaml_err" > "$_FAIL_OUT"
        record_fail "$display" "go" "$sdk_version" "$sdk_version" "yaml-validation" "$_FAIL_OUT"
        return 0
      fi
      vlog "    ✅ $wf_yaml valid"
    fi
  done <<< "$(yaml_workflow_dirs "$abs_template/.cre/template.yaml")"

  # 2. go mod verify
  vlog "    Running go mod verify..."
  if ! run_captured _out go mod verify; then
    printf '%s' "$_out" > "$_FAIL_OUT"
    info "    ❌ go mod verify failed"
    record_fail "$display" "go" "$sdk_version" "$sdk_version" "go mod verify" "$_FAIL_OUT"
    return 0
  fi
  vlog "    ✅ go mod verify passed"

  # 3. go vet (with wasip1 build tags so constrained files are included)
  vlog "    Running go vet..."
  if ! run_captured _out env GOOS=wasip1 GOARCH=wasm go vet ./...; then
    printf '%s' "$_out" > "$_FAIL_OUT"
    info "    ❌ go vet failed"
    record_fail "$display" "go" "$sdk_version" "$sdk_version" "go vet" "$_FAIL_OUT"
    return 0
  fi
  vlog "    ✅ go vet passed"

  # 4. go build each workflow dir to WASM
  local failed_build=false
  while IFS= read -r wf_dir; do
    [[ -z "$wf_dir" ]] && continue
    if [[ ! -d "$wf_dir" ]]; then
      vlog "    ⚠️  Workflow dir '$wf_dir' not found — skipping"
      continue
    fi

    local wasm_out; wasm_out=$(mktemp /tmp/cre-wasm-XXXXXX.wasm)
    vlog "    Building ./$wf_dir → WASM..."
    if ! run_captured _out env GOOS=wasip1 GOARCH=wasm go build -o "$wasm_out" "./$wf_dir"; then
      rm -f "$wasm_out"
      printf '%s' "$_out" > "$_FAIL_OUT"
      info "    ❌ go build ./$wf_dir failed (SDK $sdk_version)"
      record_fail "$display" "go" "$sdk_version" "$sdk_version" "go build ./$wf_dir" "$_FAIL_OUT"
      failed_build=true
      break
    fi
    rm -f "$wasm_out"
    vlog "    ✅ go build ./$wf_dir passed"
  done <<< "$(yaml_workflow_dirs "$abs_template/.cre/template.yaml")"

  if ! $failed_build; then
    info "  ✅ $display (SDK $sdk_version)"
    record_pass "$display" "go" "$sdk_version" "$sdk_version"
  fi
}

# --------------------------------------------------------------------------
# Main

cd "$REPO_ROOT"

TEMPLATE_ROOTS=()
while IFS= read -r root; do
  [[ -n "$root" ]] && TEMPLATE_ROOTS+=("$root")
done < <(detect_template_roots | sort -u)

if [[ ${#TEMPLATE_ROOTS[@]} -eq 0 ]]; then
  info "No changed templates detected — nothing to check."
  echo '{"templates":[],"total":0,"passed":0,"failed":0}' > "$RESULTS_FILE"
  exit 0
fi

TOTAL=${#TEMPLATE_ROOTS[@]}
info "Found $TOTAL changed template(s) to check."
info ""

IDX=0
for template_dir in "${TEMPLATE_ROOTS[@]}"; do
  IDX=$((IDX + 1))
  abs_template="$REPO_ROOT/$template_dir"
  template_yaml="$abs_template/.cre/template.yaml"

  if [[ ! -f "$template_yaml" ]]; then
    vlog "[$IDX/$TOTAL] Skipping $template_dir (no .cre/template.yaml)"
    continue
  fi

  # Validate .cre/template.yaml before anything else
  local_yaml_err=$(validate_yaml "$template_yaml" "template" || true)
  if [[ -n "$local_yaml_err" ]]; then
    info "[$IDX/$TOTAL] $template_dir"
    info "  ❌ .cre/template.yaml invalid: $local_yaml_err"
    printf '%s' "$local_yaml_err" > "$_FAIL_OUT"
    record_fail "$template_dir" "unknown" "unknown" "unknown" "template-yaml-validation" "$_FAIL_OUT"
    continue
  fi

  language=$(yaml_field "$template_yaml" "language")
  info "[$IDX/$TOTAL] $template_dir ($language)"

  case "$language" in
    typescript)
      # Check each workflow dir listed in .cre/template.yaml
      while IFS= read -r wf_dir; do
        [[ -z "$wf_dir" ]] && continue
        check_ts_workflow "$template_dir" "$wf_dir"
      done <<< "$(yaml_workflow_dirs "$template_yaml")"
      ;;
    go)
      check_go_template "$template_dir"
      ;;
    *)
      info "  ⚠️  Unknown language '$language' — skipping"
      ;;
  esac
  info ""
done

# --------------------------------------------------------------------------
# Summary

finalize_results

PASS_COUNT=0 FAIL_COUNT=0
if [[ -f "$RESULTS_FILE" ]]; then
  PASS_COUNT=$(python3 -c "import json; d=json.load(open('$RESULTS_FILE')); print(d['passed'])" 2>/dev/null || echo 0)
  FAIL_COUNT=$(python3 -c "import json; d=json.load(open('$RESULTS_FILE')); print(d['failed'])" 2>/dev/null || echo 0)
fi

info "========================================================"
info "Results: $PASS_COUNT passed, $FAIL_COUNT failed"
info "========================================================"

if [[ "$FAIL_COUNT" -gt 0 ]]; then
  info ""
  info "Failed templates:"
  python3 -c "
import json, sys
data = json.load(open('$RESULTS_FILE'))
for t in data['templates']:
    if t['status'] != 'pass':
        step = t.get('failure_step') or ''
        print(f\"  ❌ {t['name']} — {step}\")
        out = t.get('failure_output') or ''
        if out:
            for line in out.strip().split('\n')[:20]:
                print(f'     {line}')
"
  exit 1
else
  info ""
  info "All templates passed!"
  exit 0
fi
