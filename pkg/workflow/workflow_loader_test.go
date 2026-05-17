//go:build !integration

package workflow

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadParsedWorkflow_FromYAML(t *testing.T) {
	tmpDir := t.TempDir()
	workflowPath := filepath.Join(tmpDir, "worker.lock.yml")
	content := `name: Worker
on:
  workflow_call: {}
jobs:
  work:
    runs-on: ubuntu-latest
`
	require.NoError(t, os.WriteFile(workflowPath, []byte(content), 0o644), "Should write test workflow file")

	workflow, err := loadParsedWorkflow(workflowPath)
	require.NoError(t, err, "Should load YAML workflow")
	assert.Contains(t, workflow, "jobs", "Should include parsed YAML keys")
}

func TestLoadParsedWorkflow_FromMarkdown(t *testing.T) {
	tmpDir := t.TempDir()
	workflowPath := filepath.Join(tmpDir, "worker.md")
	content := `---
on:
  workflow_call:
    inputs:
      payload:
        type: string
permissions:
  contents: read
---

# Worker
`
	require.NoError(t, os.WriteFile(workflowPath, []byte(content), 0o644), "Should write test markdown workflow")

	workflow, err := loadParsedWorkflow(workflowPath)
	require.NoError(t, err, "Should load markdown workflow frontmatter")
	assert.Contains(t, workflow, "on", "Should include frontmatter keys")
	assert.Contains(t, workflow, "permissions", "Should include permissions from frontmatter")
}

func TestLoadParsedWorkflow_FromMarkdownWithoutFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	workflowPath := filepath.Join(tmpDir, "worker.md")
	require.NoError(t, os.WriteFile(workflowPath, []byte("# No frontmatter"), 0o644), "Should write markdown without frontmatter")

	workflow, err := loadParsedWorkflow(workflowPath)
	require.NoError(t, err, "Should not error when frontmatter is missing")
	assert.Empty(t, workflow, "Should return an empty parsed map")
}

func TestLoadParsedWorkflow_UnsupportedExtension(t *testing.T) {
	tmpDir := t.TempDir()
	workflowPath := filepath.Join(tmpDir, "worker.txt")
	require.NoError(t, os.WriteFile(workflowPath, []byte("text"), 0o644), "Should write file with unsupported extension")

	workflow, err := loadParsedWorkflow(workflowPath)
	require.Error(t, err, "Should fail for unsupported extension")
	assert.Nil(t, workflow, "Should return nil workflow on unsupported extension")
}
