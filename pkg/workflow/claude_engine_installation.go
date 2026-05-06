package workflow

import "github.com/github/gh-aw/pkg/constants"

func (e *ClaudeEngine) GetSecretValidationStep(workflowData *WorkflowData) GitHubActionStep {
	return BuildDefaultSecretValidationStep(
		workflowData,
		[]string{"ANTHROPIC_API_KEY"},
		"Claude Code",
		"https://github.github.com/gh-aw/reference/engines/#anthropic-claude-code",
	)
}

func (e *ClaudeEngine) GetInstallationSteps(workflowData *WorkflowData) []GitHubActionStep {
	claudeLog.Printf("Generating installation steps for Claude engine: workflow=%s", workflowData.Name)

	// Skip installation if custom command is specified
	if workflowData.EngineConfig != nil && workflowData.EngineConfig.Command != "" {
		claudeLog.Printf("Skipping installation steps: custom command specified (%s)", workflowData.EngineConfig.Command)
		return []GitHubActionStep{}
	}

	// Use version from engine config if provided, otherwise default to pinned version
	version := string(constants.DefaultClaudeCodeVersion)
	if workflowData.EngineConfig != nil && workflowData.EngineConfig.Version != "" {
		version = workflowData.EngineConfig.Version
	}

	// Claude Code requires post-install scripts (native binaries) so --ignore-scripts must
	// NOT be passed. This is intentionally different from other engine installs.
	npmSteps := GenerateNpmInstallSteps(
		"@anthropic-ai/claude-code",
		version,
		"Install Claude Code CLI",
		"claude",
		true, // Include Node.js setup
		true, // Claude Code requires post-install scripts for native binaries
	)
	return BuildNpmEngineInstallStepsWithAWF(npmSteps, workflowData)
}
