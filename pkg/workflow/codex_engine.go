package workflow

import (
	"maps"
	"regexp"
	"sort"
	"strings"

	"github.com/github/gh-aw/pkg/constants"
	"github.com/github/gh-aw/pkg/logger"
)

var codexEngineLog = logger.New("workflow:codex_engine")

// Pre-compiled regexes for Codex log parsing (performance optimization)
var (
	codexToolCallOldFormat    = regexp.MustCompile(`\] tool ([^(]+)\(`)
	codexToolCallNewFormat    = regexp.MustCompile(`^tool ([^(]+)\(`)
	codexExecCommandOldFormat = regexp.MustCompile(`\] exec (.+?) in`)
	codexExecCommandNewFormat = regexp.MustCompile(`^exec (.+?) in`)
	codexDurationPattern      = regexp.MustCompile(`in\s+(\d+(?:\.\d+)?)\s*s`)
	codexTokenUsagePattern    = regexp.MustCompile(`(?i)tokens\s+used[:\s]+(\d+)`)
	codexTotalTokensPattern   = regexp.MustCompile(`total_tokens:\s*(\d+)`)
)

// CodexEngine represents the Codex agentic engine
type CodexEngine struct {
	BaseEngine
}

func NewCodexEngine() *CodexEngine {
	return &CodexEngine{
		BaseEngine: BaseEngine{
			id:                       "codex",
			displayName:              "Codex",
			description:              "Uses OpenAI Codex CLI with MCP server support",
			experimental:             false,
			supportsToolsAllowlist:   true,
			supportsMaxTurns:         false, // Codex does not support max-turns feature
			supportsMaxContinuations: false, // Codex does not support --max-autopilot-continues-style continuation mode
			supportsWebSearch:        true,  // Codex has built-in web-search support
			supportsNativeAgentFile:  false, // Codex does not support agent file natively; the compiler prepends the agent file content to prompt.txt
			dedicatedLLMGatewayPort:  constants.CodexLLMGatewayPort,
		},
	}
}

// GetModelEnvVarName returns an empty string because the Codex CLI does not support
// selecting the model via a native environment variable. Model selection for Codex
// is done via the -c model=... configuration override in the shell command.
func (e *CodexEngine) GetModelEnvVarName() string {
	return ""
}

// GetRequiredSecretNames returns the list of secrets required by the Codex engine
// This includes CODEX_API_KEY, OPENAI_API_KEY, and optionally MCP_GATEWAY_API_KEY and mcp-scripts secrets
func (e *CodexEngine) GetRequiredSecretNames(workflowData *WorkflowData) []string {
	return append([]string{"CODEX_API_KEY", "OPENAI_API_KEY"}, collectCommonMCPSecrets(workflowData)...)
}

// GetSecretValidationStep is implemented in codex_engine_installation.go

// GetInstallationSteps is implemented in codex_engine_installation.go

// GetExecutionSteps is implemented in codex_engine_execution.go
// GetDeclaredOutputFiles returns the output files that Codex may produce.
// Use /tmp/gh-aw for Codex runtime logs because ${RUNNER_TEMP}/gh-aw is
// mounted read-only inside the AWF chroot sandbox.
func (e *CodexEngine) GetDeclaredOutputFiles() []string {
	// Return the Codex log directory for artifact collection.
	return []string{
		"/tmp/gh-aw/mcp-config/logs/",
	}
}

// GetAgentManifestFiles returns Codex-specific instruction files that should be
// treated as security-sensitive manifests.  AGENTS.md is the primary OpenAI
// Codex agent-instruction file; modifying it can redirect agent behaviour.
// CLAUDE.md and GEMINI.md are also listed because repositories often use multiple
// engines and Codex runs alongside them.
func (e *CodexEngine) GetAgentManifestFiles() []string {
	return []string{"AGENTS.md", "CLAUDE.md", "GEMINI.md"}
}

// GetAgentManifestPathPrefixes returns Codex-specific config directory prefixes.
// The .codex/ directory can contain agent configuration and task-specific settings.
func (e *CodexEngine) GetAgentManifestPathPrefixes() []string {
	return []string{".codex/"}
}

// GetHarnessScriptName returns the filename of the JavaScript harness script that wraps
// Codex CLI execution with retry logic for transient OpenAI API errors.
func (e *CodexEngine) GetHarnessScriptName() string {
	return "codex_harness.cjs"
}

// GetSquidLogsSteps returns the steps for uploading and parsing Squid logs (after secret redaction)
func (e *CodexEngine) GetSquidLogsSteps(workflowData *WorkflowData) []GitHubActionStep {
	return defaultGetSquidLogsSteps(workflowData, codexEngineLog)
}

// computeCodexToolArguments converts neutral tools to Codex-specific tool arguments.
// This ensures that playwright tools get the same allowlist as the copilot agent.
func (e *CodexEngine) computeCodexToolArguments(toolsConfig *ToolsConfig) *ToolsConfig {
	if toolsConfig == nil {
		return &ToolsConfig{
			Custom: make(map[string]MCPServerConfig),
			raw:    make(map[string]any),
		}
	}

	// Create a copy of the tools config
	result := &ToolsConfig{
		GitHub:           toolsConfig.GitHub,
		Bash:             toolsConfig.Bash,
		WebFetch:         toolsConfig.WebFetch,
		WebSearch:        toolsConfig.WebSearch,
		Edit:             toolsConfig.Edit,
		Playwright:       toolsConfig.Playwright,
		AgenticWorkflows: toolsConfig.AgenticWorkflows,
		CacheMemory:      toolsConfig.CacheMemory,
		Timeout:          toolsConfig.Timeout,
		StartupTimeout:   toolsConfig.StartupTimeout,
		Custom:           make(map[string]MCPServerConfig),
		raw:              make(map[string]any),
	}

	// Copy custom tools
	maps.Copy(result.Custom, toolsConfig.Custom)

	// Copy raw map
	maps.Copy(result.raw, toolsConfig.raw)

	// Handle playwright tool by converting it to an MCP tool configuration with copilot agent tools
	if toolsConfig.Playwright != nil {
		// Create an updated Playwright config preserving all fields including Mode
		playwrightConfig := &PlaywrightToolConfig{
			Version: toolsConfig.Playwright.Version,
			Args:    toolsConfig.Playwright.Args,
			Mode:    toolsConfig.Playwright.Mode,
		}

		result.Playwright = playwrightConfig

		// In CLI mode, playwright is not an MCP server — remove from raw map and skip MCP config entry.
		// result.raw is populated by maps.Copy(result.raw, toolsConfig.raw) earlier in this function,
		// so delete is safe regardless of whether the key was originally present.
		if playwrightConfig.IsCLIMode() {
			delete(result.raw, "playwright")
		} else {
			// Also update the Custom map entry for playwright with allowed tools list
			playwrightMCP := map[string]any{
				"allowed": GetPlaywrightTools(),
			}
			if playwrightConfig.Version != "" {
				playwrightMCP["version"] = playwrightConfig.Version
			}
			if len(playwrightConfig.Args) > 0 {
				playwrightMCP["args"] = playwrightConfig.Args
			}

			// Update raw map for backward compatibility
			result.raw["playwright"] = playwrightMCP
		}
	}

	return result
}

// computeCodexToolArgumentsFromMap is a backward-compatible wrapper that accepts
// map[string]any instead of *ToolsConfig.
func (e *CodexEngine) computeCodexToolArgumentsFromMap(tools map[string]any) map[string]any {
	toolsConfig, _ := ParseToolsConfig(tools)
	result := e.computeCodexToolArguments(toolsConfig)
	return result.ToMap()
}

// expandNeutralToolsToCodexTools is a backward-compatible wrapper around
// computeCodexToolArguments.
func (e *CodexEngine) expandNeutralToolsToCodexTools(toolsConfig *ToolsConfig) *ToolsConfig {
	return e.computeCodexToolArguments(toolsConfig)
}

// expandNeutralToolsToCodexToolsFromMap is a backward-compatible wrapper around
// computeCodexToolArgumentsFromMap.
func (e *CodexEngine) expandNeutralToolsToCodexToolsFromMap(tools map[string]any) map[string]any {
	return e.computeCodexToolArgumentsFromMap(tools)
}

func (e *CodexEngine) getShellEnvironmentPolicyVars(tools map[string]any, mcpTools []string) []string {
	// Collect all environment variables needed by MCP servers
	envVars := make(map[string]bool)

	// Always include core environment variables
	envVars["PATH"] = true
	envVars["HOME"] = true

	// Add CODEX_API_KEY for authentication
	envVars["CODEX_API_KEY"] = true
	envVars["OPENAI_API_KEY"] = true // Fallback for CODEX_API_KEY

	// Check each MCP tool for required environment variables
	for _, toolName := range mcpTools {
		switch toolName {
		case "github":
			// GitHub MCP server needs GITHUB_PERSONAL_ACCESS_TOKEN
			envVars["GITHUB_PERSONAL_ACCESS_TOKEN"] = true
		case "agentic-workflows":
			// Agentic workflows MCP server needs GITHUB_TOKEN
			envVars["GITHUB_TOKEN"] = true
		case "safe-outputs":
			// Safe outputs MCP server needs several environment variables
			envVars["GH_AW_SAFE_OUTPUTS"] = true
			envVars["GH_AW_ASSETS_BRANCH"] = true
			envVars["GH_AW_ASSETS_MAX_SIZE_KB"] = true
			envVars["GH_AW_ASSETS_ALLOWED_EXTS"] = true
			envVars["GITHUB_REPOSITORY"] = true
			envVars["GITHUB_SERVER_URL"] = true
		default:
			// For custom MCP tools, check if they have env configuration
			if toolValue, ok := tools[toolName]; ok {
				if toolConfig, ok := toolValue.(map[string]any); ok {
					// Extract environment variable names from env configuration
					if env, hasEnv := toolConfig["env"].(map[string]any); hasEnv {
						for envKey := range env {
							envVars[envKey] = true
						}
					}
				}
			}
		}
	}

	var sortedEnvVars []string
	for envVar := range envVars {
		sortedEnvVars = append(sortedEnvVars, envVar)
	}
	sort.Strings(sortedEnvVars)

	return sortedEnvVars
}

// renderShellEnvironmentPolicy generates the [shell_environment_policy] section for config.toml
// This controls which environment variables are passed through to MCP servers for security
func (e *CodexEngine) renderShellEnvironmentPolicy(yaml *strings.Builder, tools map[string]any, mcpTools []string) {
	sortedEnvVars := e.getShellEnvironmentPolicyVars(tools, mcpTools)

	// Render [shell_environment_policy] section
	yaml.WriteString("          \n")
	yaml.WriteString("          [shell_environment_policy]\n")
	yaml.WriteString("          inherit = \"core\"\n")
	yaml.WriteString("          include_only = [")
	for i, envVar := range sortedEnvVars {
		if i > 0 {
			yaml.WriteString(", ")
		}
		yaml.WriteString("\"" + envVar + "\"")
	}
	yaml.WriteString("]\n")
}

func (e *CodexEngine) renderShellEnvironmentPolicyToml(yaml *strings.Builder, tools map[string]any, mcpTools []string, indent string) {
	sortedEnvVars := e.getShellEnvironmentPolicyVars(tools, mcpTools)

	yaml.WriteString(indent + "[shell_environment_policy]\n")
	yaml.WriteString(indent + "inherit = \"core\"\n")
	yaml.WriteString(indent + "include_only = [")
	for i, envVar := range sortedEnvVars {
		if i > 0 {
			yaml.WriteString(", ")
		}
		yaml.WriteString("\"" + envVar + "\"")
	}
	yaml.WriteString("]\n")
}

// RenderMCPConfig is implemented in codex_mcp.go

// renderCodexMCPConfig is implemented in codex_mcp.go

// ParseLogMetrics is implemented in codex_logs.go

// parseCodexToolCallsWithSequence is implemented in codex_logs.go

// updateMostRecentToolWithDuration is implemented in codex_logs.go

// extractCodexTokenUsage is implemented in codex_logs.go

// GetLogParserScriptId is implemented in codex_logs.go
