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
	// mdPath originates from findWorkflowFile(), which validates paths via
	// isPathWithinDir() to prevent directory traversal before returning them.
	content, err := os.ReadFile(mdPath) // #nosec G304 -- path pre-validated by findWorkflowFile() via isPathWithinDir()
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow source %s: %w", mdPath, err)
	}

	result, err := parser.ExtractFrontmatterFromContent(string(content))
	if err != nil || result == nil {
		workflowLoaderLog.Printf("Failed to extract frontmatter from %s: %v", mdPath, err)
		return make(map[string]any), nil
	}

	return result.Frontmatter, nil
}
