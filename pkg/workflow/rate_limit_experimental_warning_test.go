//go:build integration

package workflow

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/github/gh-aw/pkg/testutil"
)

// TestRateLimitNoExperimentalWarning tests that the rate-limit feature
// does not emit an experimental warning.
func TestRateLimitNoExperimentalWarning(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expectWarning bool
	}{
		{
			name: "rate-limit enabled does not produce experimental warning",
			content: `---
on: workflow_dispatch
engine: copilot
rate-limit:
  max-runs: 5
  max-runs-window: 60
permissions:
  contents: read
  issues: read
  pull-requests: read
---

# Test Workflow
`,
			expectWarning: false,
		},
		{
			name: "no rate-limit does not produce experimental warning",
			content: `---
on: workflow_dispatch
engine: copilot
permissions:
  contents: read
  issues: read
  pull-requests: read
---

# Test Workflow
`,
			expectWarning: false,
		},
		{
			name: "rate-limit with custom ignored roles does not produce experimental warning",
			content: `---
on: workflow_dispatch
engine: copilot
rate-limit:
  max-runs: 3
  max-runs-window: 30
  ignored-roles:
    - admin
    - maintain
permissions:
  contents: read
  issues: read
  pull-requests: read
---

# Test Workflow
`,
			expectWarning: false,
		},
		{
			name: "rate-limit with events does not produce experimental warning",
			content: `---
on:
  workflow_dispatch:
  issue_comment:
    types: [created]
engine: copilot
rate-limit:
  max-runs: 5
  max-runs-window: 60
  events: [workflow_dispatch, issue_comment]
permissions:
  contents: read
  issues: read
  pull-requests: read
---

# Test Workflow
`,
			expectWarning: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := testutil.TempDir(t, "rate-limit-experimental-warning-test")

			testFile := filepath.Join(tmpDir, "test-workflow.md")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			// Capture stderr to check for warnings
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			compiler := NewCompiler()
			compiler.SetStrictMode(false)
			err := compiler.CompileWorkflow(testFile)

			// Restore stderr
			w.Close()
			os.Stderr = oldStderr
			var buf bytes.Buffer
			io.Copy(&buf, r)
			stderrOutput := buf.String()

			if err != nil {
				t.Errorf("Expected compilation to succeed but it failed: %v", err)
				return
			}

			expectedMessage := "Using experimental feature: rate-limit"

			if tt.expectWarning {
				if !strings.Contains(stderrOutput, expectedMessage) {
					t.Errorf("Expected warning containing '%s', got stderr:\n%s", expectedMessage, stderrOutput)
				}
			} else {
				if strings.Contains(stderrOutput, expectedMessage) {
					t.Errorf("Did not expect warning '%s', but got stderr:\n%s", expectedMessage, stderrOutput)
				}
			}

			// Verify warning count does not include rate-limit warning
			if !tt.expectWarning {
				if compiler.GetWarningCount() > 0 && strings.Contains(stderrOutput, expectedMessage) {
					t.Errorf("Did not expect rate-limit warning count, got stderr:\n%s", stderrOutput)
				}
			}
		})
	}
}
