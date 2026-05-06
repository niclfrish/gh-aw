package workflow

import (
	"github.com/github/gh-aw/pkg/constants"
	"github.com/github/gh-aw/pkg/logger"
)

var geminiLog = logger.New("workflow:gemini_engine")

// GeminiEngine represents the Google Gemini CLI agentic engine
type GeminiEngine struct {
	BaseEngine
}

func NewGeminiEngine() *GeminiEngine {
	return &GeminiEngine{
		BaseEngine: BaseEngine{
			id:                       "gemini",
			displayName:              "Google Gemini CLI",
			description:              "Google Gemini CLI with headless mode and LLM gateway support",
			experimental:             false,
			supportsToolsAllowlist:   true,
			supportsMaxTurns:         false,
			supportsMaxContinuations: false, // Gemini CLI does not support --max-autopilot-continues-style continuation mode
			supportsWebSearch:        false,
			supportsNativeAgentFile:  false, // Gemini does not support agent file natively; the compiler prepends the agent file content to prompt.txt
			dedicatedLLMGatewayPort:  constants.GeminiLLMGatewayPort,
		},
	}
}

// GetModelEnvVarName returns the native environment variable name that the Gemini CLI uses
// for model selection. Setting GEMINI_MODEL is equivalent to passing --model to the CLI.
func (e *GeminiEngine) GetModelEnvVarName() string {
	return constants.GeminiCLIModelEnvVar
}

// GetRequiredSecretNames returns the list of secrets required by the Gemini engine
// This includes GEMINI_API_KEY and optionally MCP_GATEWAY_API_KEY, GITHUB_MCP_SERVER_TOKEN,
// HTTP MCP header secrets, and mcp-scripts secrets
func (e *GeminiEngine) GetRequiredSecretNames(workflowData *WorkflowData) []string {
	geminiLog.Print("Collecting required secrets for Gemini engine")
	secrets := []string{"GEMINI_API_KEY"}

	// Add common MCP secrets (MCP_GATEWAY_API_KEY if MCP servers present, mcp-scripts secrets)
	secrets = append(secrets, collectCommonMCPSecrets(workflowData)...)

	// Add GitHub token for GitHub MCP server if present
	if hasGitHubTool(workflowData.ParsedTools) {
		geminiLog.Print("Adding GITHUB_MCP_SERVER_TOKEN secret")
		secrets = append(secrets, "GITHUB_MCP_SERVER_TOKEN")
	}

	// Add HTTP MCP header secret names
	headerSecrets := collectHTTPMCPHeaderSecrets(workflowData.Tools)
	for varName := range headerSecrets {
		secrets = append(secrets, varName)
	}
	if len(headerSecrets) > 0 {
		geminiLog.Printf("Added %d HTTP MCP header secrets", len(headerSecrets))
	}

	return secrets
}

// GetSecretValidationStep is implemented in gemini_engine_installation.go

// GetInstallationSteps is implemented in gemini_engine_installation.go

// GetDeclaredOutputFiles returns the output files that Gemini may produce.
// Gemini CLI writes structured error reports to /tmp/gemini-client-error-*.json
// with a timestamp in the filename (e.g. gemini-client-error-Turn.run-sendMessageStream-2026-02-21T20-45-59-824Z.json).
// These files provide detailed diagnostics when the Gemini API call fails.
// GetPreBundleSteps moves these files into /tmp/gh-aw/ so all artifact paths share a common
// ancestor under /tmp/gh-aw/ and the actions/upload-artifact LCA calculation stays correct.
func (e *GeminiEngine) GetDeclaredOutputFiles() []string {
	return []string{
		"/tmp/gh-aw/gemini-client-error-*.json",
	}
}

// GetAgentManifestFiles returns Gemini-specific instruction files that should be
// treated as security-sensitive manifests.  A fork PR that modifies these files
// can redirect the agent's behaviour or expand which files it treats as instructions.
// GEMINI.md is the primary per-project context file; AGENTS.md is the cross-engine
// convention that Gemini CLI also reads.
func (e *GeminiEngine) GetAgentManifestFiles() []string {
	return []string{"GEMINI.md", "AGENTS.md"}
}

// GetAgentManifestPathPrefixes returns Gemini-specific config directory prefixes.
// The .gemini/ directory contains settings.json and other configuration that could
// expand which files are treated as instructions or alter agent behaviour.
// Protecting this directory prevents fork PRs from injecting malicious configuration.
func (e *GeminiEngine) GetAgentManifestPathPrefixes() []string {
	return []string{".gemini/"}
}

// GetPreBundleSteps returns a step that moves Gemini CLI error reports from /tmp/ into
// /tmp/gh-aw/ before the unified artifact upload. This keeps all artifact paths under
// /tmp/gh-aw/ so that actions/upload-artifact computes the correct least-common-ancestor
// path and downstream jobs find files at the expected locations.
func (e *GeminiEngine) GetPreBundleSteps(workflowData *WorkflowData) []GitHubActionStep {
	return []GitHubActionStep{
		{
			"      - name: Move Gemini error files to artifact directory",
			"        if: always()",
			"        run: mv /tmp/gemini-client-error-*.json /tmp/gh-aw/ 2>/dev/null || true",
		},
	}
}

// GetExecutionSteps is implemented in gemini_engine_execution.go
