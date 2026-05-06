package workflow

import "github.com/github/gh-aw/pkg/constants"

func (e *CodexEngine) GetSecretValidationStep(workflowData *WorkflowData) GitHubActionStep {
	return BuildDefaultSecretValidationStep(
		workflowData,
		[]string{"CODEX_API_KEY", "OPENAI_API_KEY"},
		"Codex",
		"https://github.github.com/gh-aw/reference/engines/#openai-codex",
	)
}

func (e *CodexEngine) GetInstallationSteps(workflowData *WorkflowData) []GitHubActionStep {
	codexEngineLog.Printf("Generating installation steps for Codex engine: workflow=%s", workflowData.Name)

	// Skip installation if custom command is specified
	if workflowData.EngineConfig != nil && workflowData.EngineConfig.Command != "" {
		codexEngineLog.Printf("Skipping installation steps: custom command specified (%s)", workflowData.EngineConfig.Command)
		return []GitHubActionStep{}
	}

	// Use base installation steps (npm install only; secret validation is in the activation job)
	steps := GetBaseInstallationSteps(EngineInstallConfig{
		Secrets:         []string{"CODEX_API_KEY", "OPENAI_API_KEY"},
		DocsURL:         "https://github.github.com/gh-aw/reference/engines/#openai-codex",
		NpmPackage:      "@openai/codex",
		Version:         string(constants.DefaultCodexVersion),
		Name:            "Codex CLI",
		InstallStepName: "Install Codex CLI",
		CliName:         "codex",
	}, workflowData)

	// Add AWF installation step if firewall is enabled
	if isFirewallEnabled(workflowData) {
		firewallConfig := getFirewallConfig(workflowData)
		agentConfig := getAgentConfig(workflowData)
		var awfVersion string
		if firewallConfig != nil {
			awfVersion = firewallConfig.Version
		}

		// Install AWF binary (or skip if custom command is specified)
		awfInstall := generateAWFInstallationStep(awfVersion, agentConfig)
		if len(awfInstall) > 0 {
			steps = append(steps, awfInstall)
		}
	}

	return steps
}
