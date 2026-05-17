//go:build !integration

package workflow_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/github/gh-aw/pkg/stringutil"
	"github.com/github/gh-aw/pkg/testutil"
	"github.com/github/gh-aw/pkg/workflow"
)

func TestCompileWorkflowWithSharedObservabilityOTLPImport_UsesHeadersSecrets(t *testing.T) {
	tempDir := testutil.TempDir(t, "test-*")
	sharedDir := filepath.Join(tempDir, "shared")
	if err := os.MkdirAll(sharedDir, 0o755); err != nil {
		t.Fatalf("Failed to create shared directory: %v", err)
	}

	sharedSourcePath := filepath.Join("..", "..", ".github", "workflows", "shared", "observability-otlp.md")
	sharedContent, err := os.ReadFile(sharedSourcePath)
	if err != nil {
		t.Fatalf("Failed to read shared observability import: %v", err)
	}

	sharedPath := filepath.Join(sharedDir, "observability-otlp.md")
	if err := os.WriteFile(sharedPath, sharedContent, 0o644); err != nil {
		t.Fatalf("Failed to write shared observability import: %v", err)
	}

	workflowPath := filepath.Join(tempDir, "test-workflow.md")
	workflowContent := `---
on: issues
permissions:
  contents: read
  issues: read
  pull-requests: read
engine: copilot
imports:
  - shared/observability-otlp.md
---

# Test Workflow
`
	if err := os.WriteFile(workflowPath, []byte(workflowContent), 0o644); err != nil {
		t.Fatalf("Failed to write test workflow: %v", err)
	}

	compiler := workflow.NewCompiler()
	if err := compiler.CompileWorkflow(workflowPath); err != nil {
		t.Fatalf("CompileWorkflow failed: %v", err)
	}

	lockFilePath := stringutil.MarkdownToLockFile(workflowPath)
	lockFileContent, err := os.ReadFile(lockFilePath)
	if err != nil {
		t.Fatalf("Failed to read lock file: %v", err)
	}

	compiled := string(lockFileContent)
	if !strings.Contains(compiled, "GH_AW_OTEL_SENTRY_HEADERS") {
		t.Error("Expected compiled workflow to reference GH_AW_OTEL_SENTRY_HEADERS")
	}
	if !strings.Contains(compiled, "GH_AW_OTEL_GRAFANA_HEADERS") {
		t.Error("Expected compiled workflow to reference GH_AW_OTEL_GRAFANA_HEADERS")
	}
	if strings.Contains(compiled, "GH_AW_OTEL_SENTRY_AUTHORIZATION") {
		t.Error("Compiled workflow should not reference deprecated GH_AW_OTEL_SENTRY_AUTHORIZATION secret")
	}
	if strings.Contains(compiled, "GH_AW_OTEL_GRAFANA_AUTHORIZATION") {
		t.Error("Compiled workflow should not reference deprecated GH_AW_OTEL_GRAFANA_AUTHORIZATION secret")
	}
}
