package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/github/gh-aw/pkg/logger"
	"github.com/github/gh-aw/pkg/parser"
)

var workflowLoaderLog = logger.New("workflow:workflow_loader")

// loadParsedWorkflow loads a workflow file into a parsed map.
//
// For .yml/.yaml files, this returns the parsed YAML workflow.
// For .md files, this returns parsed frontmatter (or an empty map when
// frontmatter is absent/invalid).
//
// Callers are responsible for passing trusted, repository-bounded paths
// (for example from findWorkflowFile()).
func loadParsedWorkflow(workflowPath string) (map[string]any, error) {
	ext := strings.ToLower(filepath.Ext(workflowPath))

	switch ext {
	case ".yml", ".yaml":
		return readWorkflowYAML(workflowPath)
	case ".md":
		return loadParsedMarkdownWorkflow(workflowPath)
	default:
		return nil, fmt.Errorf("unsupported workflow file extension for %s: %s", workflowPath, ext)
	}
}

func loadParsedMarkdownWorkflow(mdPath string) (map[string]any, error) {
	absPath, err := filepath.Abs(mdPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve workflow path %s: %w", mdPath, err)
	}

	content, err := os.ReadFile(absPath) // #nosec G304 -- callers must validate repository-bounded paths via findWorkflowFile
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow source %s: %w", mdPath, err)
	}

	result, err := parser.ExtractFrontmatterFromContent(string(content))
	if err != nil {
		workflowLoaderLog.Printf("Failed to parse frontmatter from %s: %v", mdPath, err)
		return make(map[string]any), nil
	}
	if result == nil {
		workflowLoaderLog.Printf("No frontmatter found in %s", mdPath)
		return make(map[string]any), nil
	}

	return result.Frontmatter, nil
}
