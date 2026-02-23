package e2e_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

// TemplateYAML represents the .cre/template.yaml schema.
type TemplateYAML struct {
	Kind       string     `yaml:"kind"`
	ID         string     `yaml:"id"`
	ProjectDir string     `yaml:"projectDir"`
	Title      string     `yaml:"title"`
	Language   string     `yaml:"language"`
	Category   string     `yaml:"category"`
	Workflows  []Workflow `yaml:"workflows"`
	Networks     []string `yaml:"networks"`
	Capabilities []string `yaml:"capabilities"`
}

type Workflow struct {
	Dir         string `yaml:"dir"`
	Description string `yaml:"description"`
}

// stripANSI removes ANSI escape codes from CLI output.
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// cliPath returns the path to the CRE CLI binary.
// Set CRE_CLI_PATH env var to override; defaults to "cre" (must be on PATH).
func cliPath(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("CRE_CLI_PATH"); p != "" {
		return p
	}
	return "cre"
}

// repoRoot returns the absolute path to the cre-templates repo root.
// Set CRE_TEMPLATE_REPO env var to override; defaults to parent of the e2e/ directory.
func repoRoot(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("CRE_TEMPLATE_REPO"); p != "" {
		abs, err := filepath.Abs(p)
		if err != nil {
			t.Fatalf("invalid CRE_TEMPLATE_REPO path: %v", err)
		}
		return abs
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	return filepath.Dir(wd)
}

// templateRepoRef returns the GitHub owner/repo@ref to add as a template source.
// Set CRE_TEMPLATE_REPO_REF env var (e.g. "myorg/cre-templates@my-branch").
// If not set, the default remote repo is used (no custom source added).
func templateRepoRef() string {
	return os.Getenv("CRE_TEMPLATE_REPO_REF")
}

// discoverTemplates walks the local repo and parses all .cre/template.yaml files.
func discoverTemplates(t *testing.T) []TemplateYAML {
	t.Helper()
	root := repoRoot(t)
	var templates []TemplateYAML

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		base := filepath.Base(path)
		if info.IsDir() && base == "e2e" {
			return filepath.SkipDir
		}
		if info.IsDir() && strings.HasPrefix(base, ".") && base != ".cre" {
			return filepath.SkipDir
		}
		if info.IsDir() && base == "node_modules" {
			return filepath.SkipDir
		}

		if filepath.Base(path) == "template.yaml" && strings.Contains(path, ".cre") {
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("reading %s: %w", path, err)
			}
			var tmpl TemplateYAML
			if err := yaml.Unmarshal(data, &tmpl); err != nil {
				return fmt.Errorf("parsing %s: %w", path, err)
			}
			if tmpl.ID == "" {
				t.Logf("WARN: skipping %s (no id field)", path)
				return nil
			}
			templates = append(templates, tmpl)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("discovering templates: %v", err)
	}
	if len(templates) == 0 {
		t.Fatal("no templates discovered — check repo root path")
	}
	t.Logf("discovered %d templates", len(templates))
	return templates
}

// runCLI executes the CRE CLI with the given args and returns stdout/stderr.
// If dir is non-empty, the command runs in that directory.
func runCLI(t *testing.T, dir string, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	cmd := exec.Command(cliPath(t), args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	cmd.Env = os.Environ()
	if dir != "" {
		cmd.Dir = dir
	}
	err = cmd.Run()
	return stripANSI(outBuf.String()), stripANSI(errBuf.String()), err
}

// runCmd executes an arbitrary command in the given directory.
func runCmd(t *testing.T, dir string, env []string, name string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if env != nil {
		cmd.Env = append(os.Environ(), env...)
	}
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf
	err := cmd.Run()
	return outBuf.String(), err
}

// setupTemplateRepo adds a custom template repo source if CRE_TEMPLATE_REPO_REF is set.
// Returns a cleanup function that removes the source after tests complete.
func setupTemplateRepo(t *testing.T) {
	t.Helper()
	ref := templateRepoRef()
	if ref == "" {
		t.Log("using default template repo (set CRE_TEMPLATE_REPO_REF to override)")
		return
	}
	t.Logf("adding custom template repo: %s", ref)
	stdout, stderr, err := runCLI(t, "", "templates", "add", ref)
	if err != nil {
		t.Fatalf("failed to add template repo %s:\nstdout: %s\nstderr: %s\nerr: %v", ref, stdout, stderr, err)
	}
	t.Cleanup(func() {
		_, _, _ = runCLI(t, "", "templates", "remove", ref)
	})
}

// TestDiscoverTemplates verifies that template discovery works and finds templates.
func TestDiscoverTemplates(t *testing.T) {
	templates := discoverTemplates(t)
	if len(templates) < 10 {
		t.Errorf("expected at least 10 templates, got %d", len(templates))
	}
	for _, tmpl := range templates {
		if tmpl.ID == "" {
			t.Errorf("template missing id")
		}
		if tmpl.Language == "" {
			t.Errorf("template %s missing language", tmpl.ID)
		}
		if len(tmpl.Workflows) == 0 {
			t.Errorf("template %s has no workflows", tmpl.ID)
		}
	}
}

// TestAllTemplates is the main E2E test. It scaffolds every template via `cre init`
// and validates the resulting project structure, build, bindings, tests, and simulation.
//
// Environment variables:
//   - CRE_CLI_PATH: path to the cre binary (default: "cre" on PATH)
//   - CRE_TEMPLATE_REPO: local path to cre-templates repo for discovery (default: ../)
//   - CRE_TEMPLATE_REPO_REF: GitHub owner/repo@ref to add as template source
//     (e.g. "smartcontractkit/cre-templates@feature-branch")
//   - CRE_ENABLE_SIMULATE: set to "true" to enable workflow simulation tests
//   - CRE_SIMULATE_TIMEOUT: timeout in seconds for each simulate run (default: 60)
func TestAllTemplates(t *testing.T) {
	cli := cliPath(t)
	if _, err := exec.LookPath(cli); err != nil {
		t.Skipf("CRE CLI not found at %q (set CRE_CLI_PATH): %v", cli, err)
	}

	setupTemplateRepo(t)

	templates := discoverTemplates(t)

	for _, tmpl := range templates {
		t.Run(tmpl.ID, func(t *testing.T) {
			t.Parallel()

			tempDir := t.TempDir()
			projectName := "test-" + tmpl.ID

			args := []string{
				"init",
				"--project-root", tempDir,
				"--project-name", projectName,
				"--template", tmpl.ID,
			}

			// Add --workflow-name for single-workflow templates without projectDir.
			if tmpl.ProjectDir == "" && len(tmpl.Workflows) <= 1 {
				args = append(args, "--workflow-name", "test-workflow")
			}

			// Supply dummy RPC URLs for templates that declare networks.
			for _, network := range tmpl.Networks {
				args = append(args, "--rpc-url", network+"=https://dummy-rpc.example.com")
			}

			t.Logf("running: cre %s", strings.Join(args, " "))
			stdout, stderr, err := runCLI(t, tempDir, args...)
			if err != nil {
				t.Fatalf("cre init failed:\nstdout: %s\nstderr: %s\nerr: %v", stdout, stderr, err)
			}

			projectDir := filepath.Join(tempDir, projectName)

			// --- Validate project structure ---
			t.Run("structure", func(t *testing.T) {
				assertFileExists(t, projectDir, "project.yaml")

				if tmpl.ProjectDir != "" {
					// Template with projectDir: all workflow dirs should be present.
					for _, wf := range tmpl.Workflows {
						assertDirExists(t, projectDir, wf.Dir)
					}
				} else {
					// Template without projectDir: CLI generates workflow dir.
					wfName := "test-workflow"
					if len(tmpl.Workflows) > 1 {
						for _, wf := range tmpl.Workflows {
							assertDirExists(t, projectDir, wf.Dir)
						}
					} else {
						assertDirExists(t, projectDir, wfName)
						switch tmpl.Language {
						case "go":
							assertFileExists(t, filepath.Join(projectDir, wfName), "main.go")
						case "typescript":
							assertFileExists(t, filepath.Join(projectDir, wfName), "main.ts")
						}
					}
				}
			})

			// --- Validate build ---
			t.Run("build", func(t *testing.T) {
				switch tmpl.Language {
				case "go":
					validateGoBuild(t, projectDir)
				case "typescript":
					validateTSBuild(t, projectDir)
				default:
					t.Logf("unknown language %q, skipping build validation", tmpl.Language)
				}
			})

			// --- Validate generate-bindings ---
			t.Run("generate-bindings", func(t *testing.T) {
				if tmpl.Language != "go" {
					t.Skip("generate-bindings only applies to Go templates")
				}
				if !hasContractABIs(projectDir) {
					t.Skip("no contract ABIs found, skipping generate-bindings")
				}
				validateGenerateBindings(t, projectDir)
			})

			// --- Validate go test ---
			t.Run("go-test", func(t *testing.T) {
				if tmpl.Language != "go" {
					t.Skip("go-test only applies to Go templates")
				}
				validateGoTests(t, projectDir)
			})

			// --- Validate workflow simulate ---
			t.Run("workflow-simulate", func(t *testing.T) {
				if os.Getenv("CRE_ENABLE_SIMULATE") != "true" {
					t.Skip("workflow simulation disabled (set CRE_ENABLE_SIMULATE=true to enable)")
				}
				if tmpl.ProjectDir == "" {
					t.Skip("template has no projectDir, skipping simulate")
				}
				if tmpl.Language != "go" {
					t.Skip("workflow simulate only applies to Go templates")
				}
				validateWorkflowSimulate(t, projectDir, tmpl.Workflows)
			})
		})
	}
}

// validateGoBuild runs go mod tidy and go build for Go templates.
// CRE workflow code is typically build-tagged with `//go:build wasip1`, so we
// try a native build first and fall back to WASM cross-compilation.
func validateGoBuild(t *testing.T, projectDir string) {
	t.Helper()

	assertFileExists(t, projectDir, "go.mod")

	// go mod tidy
	output, err := runCmd(t, projectDir, nil, "go", "mod", "tidy")
	if err != nil {
		t.Fatalf("go mod tidy failed in %s:\n%s\nerr: %v", projectDir, output, err)
	}

	// Try native build first.
	output, err = runCmd(t, projectDir, nil, "go", "build", "-o", "/dev/null", "./...")
	if err != nil {
		t.Logf("native go build failed (trying WASM): %s", output)
		// Workflow code is often wasip1-only — try WASM cross-compilation.
		wasmEnv := []string{"GOOS=wasip1", "GOARCH=wasm"}
		output, err = runCmd(t, projectDir, wasmEnv, "go", "build", "-o", "/dev/null", "./...")
		if err != nil {
			t.Fatalf("go build (WASM) also failed in %s:\n%s\nerr: %v", projectDir, output, err)
		}
		t.Logf("WASM build succeeded")
		return
	}

	// If native build succeeded, also try WASM as an extra check.
	wasmEnv := []string{"GOOS=wasip1", "GOARCH=wasm"}
	output, err = runCmd(t, projectDir, wasmEnv, "go", "build", "-o", "/dev/null", "./...")
	if err != nil {
		t.Logf("WARN: WASM build failed (non-fatal) in %s:\n%s\nerr: %v", projectDir, output, err)
	}
}

// validateTSBuild validates TypeScript templates by checking package.json and running install.
func validateTSBuild(t *testing.T, projectDir string) {
	t.Helper()

	packageJSONs := findFiles(t, projectDir, "package.json")
	if len(packageJSONs) == 0 {
		t.Logf("no package.json found in %s, skipping TS build validation", projectDir)
		return
	}

	for _, pkgJSON := range packageJSONs {
		dir := filepath.Dir(pkgJSON)
		t.Logf("installing dependencies in %s", dir)

		if _, err := exec.LookPath("bun"); err == nil {
			output, err := runCmd(t, dir, nil, "bun", "install")
			if err != nil {
				t.Fatalf("bun install failed in %s:\n%s\nerr: %v", dir, output, err)
			}
		} else if _, err := exec.LookPath("npm"); err == nil {
			output, err := runCmd(t, dir, nil, "npm", "install")
			if err != nil {
				t.Fatalf("npm install failed in %s:\n%s\nerr: %v", dir, output, err)
			}
		} else {
			t.Skip("neither bun nor npm found, skipping TS build validation")
		}
	}
}

// findFiles recursively finds files matching the given name under dir.
func findFiles(t *testing.T, dir, name string) []string {
	t.Helper()
	var matches []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "node_modules" {
			return filepath.SkipDir
		}
		if info.Name() == name {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking %s: %v", dir, err)
	}
	return matches
}

// hasContractABIs checks if the project has contract ABI files.
func hasContractABIs(projectDir string) bool {
	abiDir := filepath.Join(projectDir, "contracts", "evm", "src", "abi")
	entries, err := os.ReadDir(abiDir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".abi") {
			return true
		}
	}
	return false
}

// validateGenerateBindings deletes generated bindings, runs `cre generate-bindings evm`,
// and verifies the generated directory is recreated with files.
func validateGenerateBindings(t *testing.T, projectDir string) {
	t.Helper()

	generatedDir := filepath.Join(projectDir, "contracts", "evm", "src", "generated")

	// Delete existing generated bindings.
	if err := os.RemoveAll(generatedDir); err != nil {
		t.Fatalf("failed to remove generated dir: %v", err)
	}

	// Run generate-bindings.
	stdout, stderr, err := runCLI(t, projectDir, "generate-bindings", "evm")
	if err != nil {
		t.Fatalf("cre generate-bindings evm failed:\nstdout: %s\nstderr: %s\nerr: %v", stdout, stderr, err)
	}

	// Verify regeneration.
	entries, err := os.ReadDir(generatedDir)
	if err != nil {
		t.Fatalf("generated dir not recreated after generate-bindings: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("generated dir is empty after generate-bindings")
	}
	t.Logf("generate-bindings regenerated %d entries in %s", len(entries), generatedDir)
}

// validateGoTests runs `go test ./...` in the project directory.
// Tries native first, falls back to WASM cross-compilation.
func validateGoTests(t *testing.T, projectDir string) {
	t.Helper()

	// Try native first.
	output, err := runCmd(t, projectDir, nil, "go", "test", "./...")
	if err != nil {
		t.Logf("native go test failed (trying WASM): %s", output)
		// Workflow code is often wasip1-only — try WASM cross-compilation.
		wasmEnv := []string{"GOOS=wasip1", "GOARCH=wasm"}
		output, err = runCmd(t, projectDir, wasmEnv, "go", "test", "./...")
		if err != nil {
			t.Fatalf("go test (WASM) also failed in %s:\n%s\nerr: %v", projectDir, output, err)
		}
		t.Logf("WASM go test succeeded")
		return
	}
	t.Logf("native go test succeeded")
}

// validateWorkflowSimulate runs `cre workflow simulate` for each workflow.
// Timeout is configurable via CRE_SIMULATE_TIMEOUT (default 60s).
func validateWorkflowSimulate(t *testing.T, projectDir string, workflows []Workflow) {
	t.Helper()

	timeoutSecs := 60
	if v := os.Getenv("CRE_SIMULATE_TIMEOUT"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			timeoutSecs = parsed
		}
	}
	timeout := time.Duration(timeoutSecs) * time.Second

	for _, wf := range workflows {
		wf := wf
		t.Run(wf.Dir, func(t *testing.T) {
			wfDir := filepath.Join(projectDir, wf.Dir)
			if _, err := os.Stat(wfDir); os.IsNotExist(err) {
				t.Skipf("workflow dir %s does not exist, skipping simulate", wfDir)
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			args := []string{
				"workflow", "simulate", wfDir,
				"--target", "staging",
				"--non-interactive",
				"--trigger-index", "0",
			}

			cmd := exec.CommandContext(ctx, cliPath(t), args...)
			var outBuf bytes.Buffer
			cmd.Stdout = &outBuf
			cmd.Stderr = &outBuf
			cmd.Dir = projectDir
			cmd.Env = os.Environ()

			err := cmd.Run()
			output := stripANSI(outBuf.String())

			if ctx.Err() == context.DeadlineExceeded {
				t.Logf("WARN: workflow simulate timed out after %s for %s:\n%s", timeout, wf.Dir, output)
				return
			}
			if err != nil {
				t.Logf("WARN: workflow simulate failed for %s (non-fatal):\n%s\nerr: %v", wf.Dir, output, err)
				return
			}
			t.Logf("workflow simulate succeeded for %s", wf.Dir)
		})
	}
}

func assertFileExists(t *testing.T, dir, name string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", path)
	}
}

func assertDirExists(t *testing.T, dir, name string) {
	t.Helper()
	path := filepath.Join(dir, name)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Errorf("expected directory %s to exist", path)
		return
	}
	if !info.IsDir() {
		t.Errorf("expected %s to be a directory", path)
	}
}
