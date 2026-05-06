package workflow

import "github.com/github/gh-aw/pkg/constants"

// GetSecretValidationStep returns the secret validation step for the Gemini engine.
// Returns an empty step if custom command is specified.
func (e *GeminiEngine) GetSecretValidationStep(workflowData *WorkflowData) GitHubActionStep {
	return BuildDefaultSecretValidationStep(
		workflowData,
		[]string{"GEMINI_API_KEY"},
		"Gemini CLI",
		"https://geminicli.com/docs/get-started/authentication/",
	)
}

func (e *GeminiEngine) GetInstallationSteps(workflowData *WorkflowData) []GitHubActionStep {
	geminiLog.Printf("Generating installation steps for Gemini engine: workflow=%s", workflowData.Name)

	// Skip installation if custom command is specified
	if workflowData.EngineConfig != nil && workflowData.EngineConfig.Command != "" {
		geminiLog.Printf("Skipping installation steps: custom command specified (%s)", workflowData.EngineConfig.Command)
		return []GitHubActionStep{}
	}

	npmSteps := BuildStandardNpmEngineInstallSteps(
		"@google/gemini-cli",
		string(constants.DefaultGeminiVersion),
		"Install Gemini CLI",
		"gemini",
		workflowData,
	)
	return BuildNpmEngineInstallStepsWithAWF(npmSteps, workflowData)
}
