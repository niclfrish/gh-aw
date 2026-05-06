package workflow

import (
	"fmt"
	"maps"
	"strings"

	"github.com/github/gh-aw/pkg/constants"
)

func (e *CodexEngine) GetExecutionSteps(workflowData *WorkflowData, logFile string) []GitHubActionStep {
	modelConfigured := workflowData.EngineConfig != nil && workflowData.EngineConfig.Model != ""
	firewallEnabled := isFirewallEnabled(workflowData)
	codexEngineLog.Printf("Building Codex execution steps: workflow=%s, modelConfigured=%v, firewall=%v",
		workflowData.Name, modelConfigured, firewallEnabled)

	var steps []GitHubActionStep

	// Codex does not support a native model environment variable, so model selection
	// always uses GH_AW_MODEL_AGENT_CODEX or GH_AW_MODEL_DETECTION_CODEX with shell expansion.
	// This also correctly handles GitHub Actions expressions like ${{ inputs.model }}.
	isDetectionJob := workflowData.SafeOutputs == nil
	var modelEnvVar string
	if isDetectionJob {
		modelEnvVar = constants.EnvVarModelDetectionCodex
	} else {
		modelEnvVar = constants.EnvVarModelAgentCodex
	}
	modelParam := fmt.Sprintf(`${%s:+-c model="$%s" }`, modelEnvVar, modelEnvVar)

	// Build search parameter: disable web search by default, enable only if web-search tool is present.
	// Codex enables web search by default, so we must explicitly set web_search="disabled" to disable it.
	// The --no-search flag does not exist; use the -c web_search="disabled" config option instead.
	// See https://developers.openai.com/codex/cli/features#web-search
	// Leading space is intentional: the format string concatenates this directly after "exec" with no space separator.
	webSearchParam := ` -c web_search="disabled"`
	if workflowData.ParsedTools != nil && workflowData.ParsedTools.WebSearch != nil {
		// Web search is enabled by default in Codex; no extra flag needed.
		webSearchParam = ""
	}

	// Build fetch parameter: disable the native fetch tool by default, enable only if web-fetch tool is present.
	// Codex enables the fetch tool by default, so we must explicitly set fetch="disabled" to disable it.
	// See https://developers.openai.com/api/docs/mcp#fetch-tool
	// Leading space is intentional: the format string concatenates this directly after webSearchParam with no space separator.
	webFetchParam := ` -c fetch="disabled"`
	if workflowData.ParsedTools != nil && workflowData.ParsedTools.WebFetch != nil {
		// Fetch is enabled by default in Codex; no extra flag needed.
		webFetchParam = ""
	}

	// See https://github.com/github/gh-aw/issues/892
	// --dangerously-bypass-approvals-and-sandbox: Skips all confirmation prompts and disables sandboxing
	// This is safe because AWF already provides a container-level sandbox layer
	// --skip-git-repo-check: Allows running in directories without a git repo
	fullAutoParam := " --dangerously-bypass-approvals-and-sandbox --skip-git-repo-check "

	// Build custom args parameter if specified in engineConfig
	var customArgsParam string
	if workflowData.EngineConfig != nil && len(workflowData.EngineConfig.Args) > 0 {
		var customArgsParamSb strings.Builder
		for _, arg := range workflowData.EngineConfig.Args {
			customArgsParamSb.WriteString(arg + " ")
		}
		customArgsParam += customArgsParamSb.String()
	}

	// Build the Codex command
	// Determine which command to use
	var commandName string
	if workflowData.EngineConfig != nil && workflowData.EngineConfig.Command != "" {
		commandName = workflowData.EngineConfig.Command
		codexEngineLog.Printf("Using custom command: %s", commandName)
	} else {
		// Use regular codex command - PATH is inherited via --env-all in AWF mode
		commandName = "codex"
	}

	// Determine harness script to wrap codex execution.
	// The built-in harness provides retry logic for transient OpenAI API errors
	// (rate limits, server errors).  A custom engine.harness overrides the built-in one.
	harnessScriptName := e.GetHarnessScriptName()
	if workflowData.EngineConfig != nil && workflowData.EngineConfig.HarnessScript != "" {
		harnessScriptName = workflowData.EngineConfig.HarnessScript
		codexEngineLog.Printf("Using custom harness script: %s", harnessScriptName)
	}

	// Build the Codex command.
	// The default harness (codex_harness.cjs) wraps execution with retry logic and reads the
	// prompt via --prompt-file.  The else branch is a defensive fallback for the case where
	// harnessScriptName is empty (e.g. a future code path that does not set a harness).
	var codexCommand string
	if harnessScriptName != "" {
		// Harness-wrapped execution: the harness reads --prompt-file and passes its content
		// as the last positional arg.  The harness also provides retry logic.
		execPrefix := fmt.Sprintf(`%s %s/%s %s`, nodeRuntimeResolutionCommand, SetupActionDestinationShell, harnessScriptName, commandName)
		codexCommand = fmt.Sprintf("%s %sexec%s%s%s%s--prompt-file /tmp/gh-aw/aw-prompts/prompt.txt",
			execPrefix, modelParam, webSearchParam, webFetchParam, fullAutoParam, customArgsParam)
	} else {
		// Without harness: use shell expansion for the prompt (no retry logic).
		codexCommand = fmt.Sprintf("%s %sexec%s%s%s%s\"$INSTRUCTION\"",
			commandName, modelParam, webSearchParam, webFetchParam, fullAutoParam, customArgsParam)
	}

	// Build the full command with agent file handling and AWF wrapping if enabled
	var command string
	if firewallEnabled {
		// Build AWF-wrapped command using helper function
		// Get allowed domains: prefer the pre-warmed cache on WorkflowData to avoid
		// re-running the expensive map+sort operation.
		var allowedDomains string
		if workflowData.CachedAllowedDomainsComputed {
			allowedDomains = workflowData.CachedAllowedDomainsStr
		} else {
			allowedDomains = GetAllowedDomainsForEngine(constants.CodexEngine, workflowData.NetworkPermissions, workflowData.Tools, workflowData.Runtimes)
		}
		// Add GHES/custom API target domains to the firewall allow-list when engine.api-target is set
		if workflowData.EngineConfig != nil && workflowData.EngineConfig.APITarget != "" {
			allowedDomains = mergeAPITargetDomains(allowedDomains, workflowData.EngineConfig.APITarget)
		}

		// AWF v0.15.0+ with --env-all handles most PATH setup natively (chroot mode is default):
		// - GOROOT, JAVA_HOME, etc. are handled via AWF_HOST_PATH and entrypoint.sh
		// However, npm-installed CLIs (like codex) need hostedtoolcache bin directories in PATH.
		npmPathSetup := GetNpmBinPathSetup()

		// Build the codex command with PATH setup inside the AWF container.
		// For engines that do not support native agent-file handling (including Codex),
		// the compiler prepends the agent file content to prompt.txt.
		// When using the harness, --prompt-file is passed directly; otherwise the prompt
		// is read via shell variable expansion.
		var codexCommandWithSetup string
		if harnessScriptName != "" {
			// Harness handles prompt reading via --prompt-file; no INSTRUCTION variable needed.
			codexCommandWithSetup = fmt.Sprintf(`%s && %s`, npmPathSetup, codexCommand)
		} else {
			codexCommandWithSetup = fmt.Sprintf(`%s && INSTRUCTION="$(cat /tmp/gh-aw/aw-prompts/prompt.txt)" && %s`, npmPathSetup, codexCommand)
		}
		// Add MCP CLI bin directory to PATH when cli-proxy is enabled
		if mcpCLIPath := GetMCPCLIPathSetup(workflowData); mcpCLIPath != "" {
			codexCommandWithSetup = fmt.Sprintf("%s && %s", mcpCLIPath, codexCommandWithSetup)
		}

		command = BuildAWFCommand(AWFCommandConfig{
			EngineName:     "codex",
			EngineCommand:  codexCommandWithSetup,
			LogFile:        logFile,
			WorkflowData:   workflowData,
			UsesTTY:        false, // Codex is not a TUI, outputs to stdout/stderr
			AllowedDomains: allowedDomains,
			// Create logs directory and agent step summary file before AWF.
			// The agent writes its step summary content to AgentStepSummaryPath, which is
			// appended to $GITHUB_STEP_SUMMARY after secret redaction.
			PathSetup: "mkdir -p \"$CODEX_HOME/logs\" && touch " + AgentStepSummaryPath,
			// Exclude every env var whose step-env value is a secret so the agent
			// cannot read raw token values via bash tools (env / printenv).
			ExcludeEnvVarNames: ComputeAWFExcludeEnvVarNames(workflowData, []string{"CODEX_API_KEY", "OPENAI_API_KEY"}),
		})
	} else {
		// Build the command without AWF wrapping.
		// For engines that do not support native agent-file handling (including Codex),
		// the compiler prepends the agent file content to prompt.txt so no special
		// shell variable juggling is needed here.
		if harnessScriptName != "" {
			// Harness handles prompt reading via --prompt-file; no INSTRUCTION variable needed.
			command = fmt.Sprintf(`set -o pipefail
touch %s
(umask 177 && touch %s)
mkdir -p "$CODEX_HOME/logs"
%s 2>&1 | tee %s`, AgentStepSummaryPath, logFile, codexCommand, logFile)
		} else {
			command = fmt.Sprintf(`set -o pipefail
touch %s
(umask 177 && touch %s)
INSTRUCTION="$(cat "$GH_AW_PROMPT")"
mkdir -p "$CODEX_HOME/logs"
%s 2>&1 | tee %s`, AgentStepSummaryPath, logFile, codexCommand, logFile)
		}
	}

	// Get effective GitHub token based on precedence: custom token > default
	effectiveGitHubToken := getEffectiveGitHubToken("")

	env := map[string]string{
		"CODEX_API_KEY": "${{ secrets.CODEX_API_KEY || secrets.OPENAI_API_KEY }}",
		// Override GITHUB_STEP_SUMMARY with a path that exists inside the sandbox.
		// The runner's original path is unreachable within the AWF isolated filesystem;
		// we create this file before the agent starts and append it to the real
		// $GITHUB_STEP_SUMMARY after secret redaction.
		"GITHUB_STEP_SUMMARY": AgentStepSummaryPath,
		"GH_AW_PROMPT":        "/tmp/gh-aw/aw-prompts/prompt.txt",
		// Tag the step as a GitHub AW agentic execution for discoverability by agents
		"GITHUB_AW":        "true",
		"GH_AW_MCP_CONFIG": "${{ runner.temp }}/gh-aw/mcp-config/config.toml",
		// Keep Codex runtime state in /tmp/gh-aw because ${RUNNER_TEMP}/gh-aw is
		// mounted read-only inside the AWF chroot sandbox.
		"CODEX_HOME":                   "/tmp/gh-aw/mcp-config",
		"RUST_LOG":                     "trace,hyper_util=info,mio=info,reqwest=info,os_info=info,codex_otel=warn,codex_core=debug,ocodex_exec=debug",
		"GH_AW_GITHUB_TOKEN":           effectiveGitHubToken,
		"GITHUB_PERSONAL_ACCESS_TOKEN": effectiveGitHubToken,                                     // Used by GitHub MCP server via env_vars
		"OPENAI_API_KEY":               "${{ secrets.CODEX_API_KEY || secrets.OPENAI_API_KEY }}", // Fallback for CODEX_API_KEY
	}
	// Indicate the phase: "agent" for the main run, "detection" for threat detection
	// Include the compiler version so agents can identify which gh-aw version generated the workflow
	if workflowData.IsDetectionRun {
		env["GH_AW_PHASE"] = "detection"
	} else {
		env["GH_AW_PHASE"] = "agent"
	}
	if IsRelease() {
		env["GH_AW_VERSION"] = GetVersion()
	} else {
		env["GH_AW_VERSION"] = "dev"
	}

	// Add GH_AW_SAFE_OUTPUTS if output is needed
	applySafeOutputEnvToMap(env, workflowData)

	// In sandbox (AWF) mode, set git identity environment variables so the first git commit
	// succeeds inside the container. AWF's --env-all forwards these to the container, ensuring
	// git does not rely on the host-side ~/.gitconfig which is not visible in the sandbox.
	if firewallEnabled {
		maps.Copy(env, getGitIdentityEnvVars())
	}

	// Add GH_AW_STARTUP_TIMEOUT environment variable (in seconds) if startup-timeout is specified
	// Supports both literal integers and GitHub Actions expressions (e.g. "${{ inputs.startup-timeout }}")
	if workflowData.ToolsStartupTimeout != "" {
		env["GH_AW_STARTUP_TIMEOUT"] = workflowData.ToolsStartupTimeout
	}

	// Add GH_AW_TOOL_TIMEOUT environment variable (in seconds) if timeout is specified
	// Supports both literal integers and GitHub Actions expressions (e.g. "${{ inputs.tool-timeout }}")
	if workflowData.ToolsTimeout != "" {
		env["GH_AW_TOOL_TIMEOUT"] = workflowData.ToolsTimeout
	}

	// Set the model environment variable.
	// Codex has no native model env var, so model selection always goes through
	// GH_AW_MODEL_AGENT_CODEX / GH_AW_MODEL_DETECTION_CODEX with shell expansion.
	// When model is configured (static or GitHub Actions expression), set the env var directly.
	// When not configured, use the GitHub variable fallback so users can set a default.
	if modelConfigured {
		codexEngineLog.Printf("Setting %s env var for model: %s", modelEnvVar, workflowData.EngineConfig.Model)
		env[modelEnvVar] = workflowData.EngineConfig.Model
	} else {
		env[modelEnvVar] = fmt.Sprintf("${{ vars.%s || '' }}", modelEnvVar)
	}

	// Add custom environment variables from engine config
	if workflowData.EngineConfig != nil && len(workflowData.EngineConfig.Env) > 0 {
		maps.Copy(env, workflowData.EngineConfig.Env)
	}

	// Add custom environment variables from agent config
	agentConfig := getAgentConfig(workflowData)
	if agentConfig != nil && len(agentConfig.Env) > 0 {
		maps.Copy(env, agentConfig.Env)
		codexEngineLog.Printf("Added %d custom env vars from agent config", len(agentConfig.Env))
	}

	// Add mcp-scripts secrets to env for passthrough to MCP servers
	if IsMCPScriptsEnabled(workflowData.MCPScripts) {
		mcpScriptsSecrets := collectMCPScriptsSecrets(workflowData.MCPScripts)
		for varName, secretExpr := range mcpScriptsSecrets {
			// Only add if not already in env
			if _, exists := env[varName]; !exists {
				env[varName] = secretExpr
			}
		}
	}

	// Generate the step for Codex execution
	stepName := "Execute Codex CLI"
	var stepLines []string

	stepLines = append(stepLines, "      - name: "+stepName)
	stepLines = append(stepLines, "        id: agentic_execution")

	// Filter environment variables to only include allowed secrets
	// This is a security measure to prevent exposing unnecessary secrets to the AWF container
	allowedSecrets := e.GetRequiredSecretNames(workflowData)
	filteredEnv := FilterEnvForSecrets(env, allowedSecrets)

	// Inject GH_TOKEN for CLI proxy (added after filtering since it uses a special
	// fallback expression that is always allowed when cli-proxy is enabled)
	addCliProxyGHTokenToEnv(filteredEnv, workflowData)

	// Format step with command and filtered environment variables using shared helper
	stepLines = FormatStepWithCommandAndEnv(stepLines, command, filteredEnv)

	steps = append(steps, GitHubActionStep(stepLines))

	return steps
}
