package workflow

import (
	"github.com/github/gh-aw/pkg/constants"
	"github.com/github/gh-aw/pkg/logger"
)

var claudeLog = logger.New("workflow:claude_engine")

// ClaudeEngine represents the Claude Code agentic engine
type ClaudeEngine struct {
	BaseEngine
}

func NewClaudeEngine() *ClaudeEngine {
	return &ClaudeEngine{
		BaseEngine: BaseEngine{
			id:                       "claude",
			displayName:              "Claude Code",
			description:              "Uses Claude Code with full MCP tool support and allow-listing",
			experimental:             false,
			supportsToolsAllowlist:   true,
			supportsMaxTurns:         true,  // Claude supports max-turns feature
			supportsMaxContinuations: false, // Claude Code does not support --max-autopilot-continues-style continuation
			supportsWebSearch:        true,  // Claude has built-in WebSearch support
			supportsNativeAgentFile:  false, // Claude does not support agent file natively; the compiler prepends the agent file content to prompt.txt
			supportsBareMode:         true,  // Claude CLI supports --bare
			dedicatedLLMGatewayPort:  constants.ClaudeLLMGatewayPort,
		},
	}
}

// GetModelEnvVarName returns the native environment variable name that the Claude Code CLI uses
// for model selection. Setting ANTHROPIC_MODEL is equivalent to passing --model to the CLI.
func (e *ClaudeEngine) GetModelEnvVarName() string {
	return constants.ClaudeCLIModelEnvVar
}

// GetAPMTarget returns "claude" so that apm-action packs Claude-specific primitives.
func (e *ClaudeEngine) GetAPMTarget() string {
	return "claude"
}

// GetRequiredSecretNames returns the list of secrets required by the Claude engine
// This includes ANTHROPIC_API_KEY and optionally MCP_GATEWAY_API_KEY and mcp-scripts secrets
func (e *ClaudeEngine) GetRequiredSecretNames(workflowData *WorkflowData) []string {
	return append([]string{"ANTHROPIC_API_KEY"}, collectCommonMCPSecrets(workflowData)...)
}

// GetSecretValidationStep is implemented in claude_engine_installation.go

// GetInstallationSteps is implemented in claude_engine_installation.go

// GetExecutionSteps is implemented in claude_engine_execution.go
// GetDeclaredOutputFiles returns the output files that Claude may produce
func (e *ClaudeEngine) GetDeclaredOutputFiles() []string {
	return []string{}
}

// GetAgentManifestFiles returns Claude-specific instruction files that should be
// treated as security-sensitive manifests.  Modifying these files can change the
// agent's instructions, guidelines, or permissions on the next run.
// CLAUDE.md is the primary per-project instruction file; AGENTS.md is the
// cross-engine convention that Claude Code also reads.
func (e *ClaudeEngine) GetAgentManifestFiles() []string {
	return []string{"CLAUDE.md", "AGENTS.md"}
}

// GetAgentManifestPathPrefixes returns Claude-specific config directory prefixes.
// The .claude/ directory contains settings, custom commands, hooks, and other
// engine configuration that could affect agent behaviour.
func (e *ClaudeEngine) GetAgentManifestPathPrefixes() []string {
	return []string{".claude/"}
}

// GetLogParserScriptId returns the JavaScript script name for parsing Claude logs
func (e *ClaudeEngine) GetLogParserScriptId() string {
	return "parse_claude_log"
}

// GetHarnessScriptName returns the filename of the JavaScript harness script that wraps
// the Claude Code CLI with retry logic for transient Anthropic API errors (overload, rate limit).
func (e *ClaudeEngine) GetHarnessScriptName() string {
	return "claude_harness.cjs"
}

// GetSquidLogsSteps returns the steps for uploading and parsing Squid logs (after secret redaction)
func (e *ClaudeEngine) GetSquidLogsSteps(workflowData *WorkflowData) []GitHubActionStep {
	return defaultGetSquidLogsSteps(workflowData, claudeLog)
}
