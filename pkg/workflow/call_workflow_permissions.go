package workflow

import (
	"fmt"

	"github.com/github/gh-aw/pkg/logger"
)

var callWorkflowPermissionsLog = logger.New("workflow:call_workflow_permissions")

// extractJobPermissionsFromParsedWorkflow extracts and merges all job-level permissions
// from a parsed GitHub Actions workflow map. Returns the union of all jobs' permissions.
func extractJobPermissionsFromParsedWorkflow(workflow map[string]any) *Permissions {
	merged := NewPermissions()

	jobsSection, ok := workflow["jobs"]
	if !ok {
		return merged
	}

	jobsMap, ok := jobsSection.(map[string]any)
	if !ok {
		return merged
	}

	for jobName, jobConfig := range jobsMap {
		jobMap, ok := jobConfig.(map[string]any)
		if !ok {
			continue
		}

		permsValue, hasPerms := jobMap["permissions"]
		if !hasPerms {
			callWorkflowPermissionsLog.Printf("Job '%s' has no permissions block, skipping", jobName)
			continue
		}

		jobPerms := NewPermissionsParserFromValue(permsValue).ToPermissions()
		callWorkflowPermissionsLog.Printf("Merging permissions from job '%s'", jobName)
		merged.Merge(jobPerms)
	}

	return merged
}

// extractCallWorkflowPermissions returns the permission superset required by the worker
// workflow identified by workflowName. It resolves the file in priority order:
// .lock.yml > .yml > .md (same-batch compilation target).
//
// For compiled files (.lock.yml / .yml), permissions are extracted from each job's
// permissions block and unioned together. For .md sources, the frontmatter-level
// permissions field is used as a proxy (the compiler will turn it into per-job
// permissions when the worker is eventually compiled).
//
// Returns nil when no workflow file is found or no permissions are declared.
// The caller should omit the permissions block on the call-* job in that case.
func extractCallWorkflowPermissions(workflowName, markdownPath string) (*Permissions, error) {
	fileResult, err := findWorkflowFile(workflowName, markdownPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find workflow file for '%s': %w", workflowName, err)
	}

	// Priority: .lock.yml > .yml > .md
	if fileResult.lockExists {
		return extractPermissionsFromYAMLFile(fileResult.lockPath)
	}

	if fileResult.ymlExists {
		return extractPermissionsFromYAMLFile(fileResult.ymlPath)
	}

	if fileResult.mdExists {
		return extractPermissionsFromWorkflowSource(fileResult.mdPath)
	}

	// No file found — return nil so the caller omits the permissions block.
	callWorkflowPermissionsLog.Printf("No workflow file found for '%s', skipping permissions", workflowName)
	return nil, nil
}

// extractPermissionsFromYAMLFile reads a .lock.yml or .yml workflow file, parses it,
// and returns the merged permissions from all its jobs.
func extractPermissionsFromYAMLFile(filePath string) (*Permissions, error) {
	workflow, err := loadParsedWorkflow(filePath)
	if err != nil {
		return nil, err
	}

	perms := extractJobPermissionsFromParsedWorkflow(workflow)
	callWorkflowPermissionsLog.Printf("Extracted permissions from YAML file %s", filePath)
	return perms, nil
}

// extractPermissionsFromWorkflowSource reads a workflow source file and extracts permissions.
// For Markdown sources, it uses frontmatter-level permissions as a proxy for the job
// permissions that will be generated when the worker is compiled.
func extractPermissionsFromWorkflowSource(workflowPath string) (*Permissions, error) {
	workflow, err := loadParsedWorkflow(workflowPath)
	if err != nil {
		return nil, err
	}

	permsValue, hasPerms := workflow["permissions"]
	if !hasPerms {
		callWorkflowPermissionsLog.Printf("No permissions in workflow source %s", workflowPath)
		return nil, nil
	}

	perms := NewPermissionsParserFromValue(permsValue).ToPermissions()
	callWorkflowPermissionsLog.Printf("Extracted permissions from workflow source %s", workflowPath)
	return perms, nil
}
