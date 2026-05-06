package workflow

import (
	"fmt"
	"maps"

	"github.com/github/gh-aw/pkg/constants"
)

// GetExecutionSteps returns the GitHub Actions steps for executing Gemini
func (e *GeminiEngine) GetExecutionSteps(workflowData *WorkflowData, logFile string) []GitHubActionStep {
	geminiLog.Printf("Generating execution steps for Gemini engine: workflow=%s, firewall=%v", workflowData.Name, isFirewallEnabled(workflowData))

	var steps []GitHubActionStep

	// Write .gemini/settings.json with context.includeDirectories and tools.core.
	// This step runs after the MCP gateway setup (which may have written mcpServers config)
	// and merges the context/tools settings into any existing settings.json.
	settingsStep := e.generateGeminiSettingsStep(workflowData)
	steps = append(steps, settingsStep)

	// Build gemini CLI arguments based on configuration
	var geminiArgs []string

	// Model is passed via the native GEMINI_MODEL environment variable only when explicitly
	// configured. When not configured, the Gemini CLI uses its built-in default model.
	// This avoids embedding the value directly in the shell command (which fails template injection
	// validation for GitHub Actions expressions like ${{ inputs.model }}).
	modelConfigured := workflowData.EngineConfig != nil && workflowData.EngineConfig.Model != ""

	// Gemini CLI reads MCP config from .gemini/settings.json (project-level)
	// The conversion script (convert_gateway_config_gemini.sh) writes settings.json
	// during the MCP setup step, so no --mcp-config flag is needed here.

	// Auto-approve all tool executions (equivalent to Codex's --dangerously-bypass-approvals-and-sandbox)
	// Without this, Gemini CLI's default approval mode rejects tool calls with "Tool execution denied by policy"
	geminiArgs = append(geminiArgs, "--yolo")

	// Skip the workspace trust check so --yolo is not overridden to "default" approval mode.
	// Gemini CLI v1.x checks whether the working directory is trusted and overrides --yolo
	// with "default" approval mode (exit code 55) when the folder is untrusted.
	// GEMINI_CLI_TRUST_WORKSPACE=true (also set in the step env) handles the same case via
	// environment variable, but --skip-trust is more reliable when AWF's sandbox does not
	// forward all host environment variables into the container.
	geminiArgs = append(geminiArgs, "--skip-trust")

	// Add streaming JSON output (JSONL format, compatible with the log parser)
	geminiArgs = append(geminiArgs, "--output-format", "stream-json")

	// Note: the --prompt argument is appended raw after shellJoinArgs below because it contains
	// a shell command substitution ("$(cat ...)") that must NOT go through shellEscapeArg —
	// single-quoting it would prevent shell expansion at runtime.

	// Build the command
	commandName := "gemini"
	if workflowData.EngineConfig != nil && workflowData.EngineConfig.Command != "" {
		commandName = workflowData.EngineConfig.Command
	}

	// Append the prompt arg raw (not through shellJoinArgs) to preserve shell expansion
	geminiCommand := fmt.Sprintf(`%s %s --prompt "$(cat /tmp/gh-aw/aw-prompts/prompt.txt)"`, commandName, shellJoinArgs(geminiArgs))

	// Build the full command with AWF wrapping if enabled
	var command string
	firewallEnabled := isFirewallEnabled(workflowData)
	if firewallEnabled {
		// Get allowed domains: prefer the pre-warmed cache on WorkflowData to avoid
		// re-running the expensive map+sort operation.
		var allowedDomains string
		if workflowData.CachedAllowedDomainsComputed {
			allowedDomains = workflowData.CachedAllowedDomainsStr
		} else {
			allowedDomains = GetAllowedDomainsForEngine(constants.GeminiEngine,
				workflowData.NetworkPermissions,
				workflowData.Tools,
				workflowData.Runtimes,
			)
		}
		// Add GHES/custom API target domains to the firewall allow-list when engine.api-target is set
		if workflowData.EngineConfig != nil && workflowData.EngineConfig.APITarget != "" {
			allowedDomains = mergeAPITargetDomains(allowedDomains, workflowData.EngineConfig.APITarget)
		}

		npmPathSetup := GetNpmBinPathSetup()
		geminiCommandWithPath := fmt.Sprintf("%s && %s", npmPathSetup, geminiCommand)
		// Add MCP CLI bin directory to PATH when cli-proxy is enabled
		if mcpCLIPath := GetMCPCLIPathSetup(workflowData); mcpCLIPath != "" {
			geminiCommandWithPath = fmt.Sprintf("%s && %s", mcpCLIPath, geminiCommandWithPath)
		}

		command = BuildAWFCommand(AWFCommandConfig{
			EngineName:     "gemini",
			EngineCommand:  geminiCommandWithPath,
			LogFile:        logFile,
			WorkflowData:   workflowData,
			UsesTTY:        false,
			AllowedDomains: allowedDomains,
			// Create the agent step summary file before AWF starts so it is accessible
			// inside the sandbox. The agent writes its step summary content here, and the
			// file is appended to $GITHUB_STEP_SUMMARY after secret redaction.
			PathSetup: "touch " + AgentStepSummaryPath,
			// Exclude every env var whose step-env value is a secret so the agent
			// cannot read raw token values via bash tools (env / printenv).
			ExcludeEnvVarNames: ComputeAWFExcludeEnvVarNames(workflowData, []string{"GEMINI_API_KEY"}),
		})
	} else {
		command = fmt.Sprintf(`set -o pipefail
touch %s
(umask 177 && touch %s)
%s 2>&1 | tee -a %s`, AgentStepSummaryPath, logFile, geminiCommand, logFile)
	}

	// Build environment variables
	env := map[string]string{
		"GEMINI_API_KEY": "${{ secrets.GEMINI_API_KEY }}",
		"GH_AW_PROMPT":   "/tmp/gh-aw/aw-prompts/prompt.txt",
		// Tag the step as a GitHub AW agentic execution for discoverability by agents
		"GITHUB_AW":        "true",
		"GITHUB_WORKSPACE": "${{ github.workspace }}",
		// Override GITHUB_STEP_SUMMARY with a path that exists inside the sandbox.
		// The runner's original path is unreachable within the AWF isolated filesystem;
		// we create this file before the agent starts and append it to the real
		// $GITHUB_STEP_SUMMARY after secret redaction.
		"GITHUB_STEP_SUMMARY": AgentStepSummaryPath,
		// Enable verbose debug logging from Gemini CLI for better diagnostics.
		// Gemini CLI uses the npm 'debug' package, and 'gemini-cli:*' enables all
		// internal Gemini CLI debug channels (see: https://gemini-cli-docs.pages.dev/cli/configuration).
		// Non-JSON debug lines are gracefully skipped by ParseLogMetrics.
		"DEBUG": "gemini-cli:*",
		// Trust the workspace to prevent Gemini CLI v1.x from overriding --yolo to default
		// approval mode when the workspace is untrusted, which causes exit code 55.
		"GEMINI_CLI_TRUST_WORKSPACE": "true",
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

	// Add MCP config env var if needed (points to .gemini/settings.json for Gemini)
	if HasMCPServers(workflowData) {
		env["GH_AW_MCP_CONFIG"] = "${{ github.workspace }}/.gemini/settings.json"
	}

	// When the firewall (AWF) is enabled with --enable-api-proxy, point Gemini CLI at the
	// LLM gateway sidecar instead of the real googleapis.com endpoint.
	if firewallEnabled {
		env["GEMINI_API_BASE_URL"] = fmt.Sprintf("http://host.docker.internal:%d", constants.GeminiLLMGatewayPort)

		// Set git identity environment variables so the first git commit succeeds inside the
		// container. AWF's --env-all forwards these to the container, ensuring git does not
		// rely on the host-side ~/.gitconfig which is not visible in the sandbox.
		maps.Copy(env, getGitIdentityEnvVars())
	}

	// Add safe outputs env
	applySafeOutputEnvToMap(env, workflowData)

	// Set the model environment variable only when explicitly configured.
	// When model is configured, use the native GEMINI_MODEL env var - the Gemini CLI reads it
	// directly, avoiding the need to embed the value in the shell command (which would fail
	// template injection validation for GitHub Actions expressions like ${{ inputs.model }}).
	// When model is not configured, let the Gemini CLI use its built-in default model.
	if modelConfigured {
		geminiLog.Printf("Setting %s env var for model: %s", constants.GeminiCLIModelEnvVar, workflowData.EngineConfig.Model)
		env[constants.GeminiCLIModelEnvVar] = workflowData.EngineConfig.Model
	}

	// Add custom environment variables from engine config.
	// This allows users to override the default engine token expression (e.g.
	// GEMINI_API_KEY: ${{ secrets.MY_ORG_GEMINI_KEY }}) via engine.env.
	if workflowData.EngineConfig != nil && len(workflowData.EngineConfig.Env) > 0 {
		maps.Copy(env, workflowData.EngineConfig.Env)
	}

	// Add custom environment variables from agent config
	agentConfig := getAgentConfig(workflowData)
	if agentConfig != nil && len(agentConfig.Env) > 0 {
		maps.Copy(env, agentConfig.Env)
		geminiLog.Printf("Added %d custom env vars from agent config", len(agentConfig.Env))
	}

	// Generate the execution step
	stepLines := []string{
		"      - name: Execute Gemini CLI",
		"        id: agentic_execution",
	}

	// Filter environment variables for security
	allowedSecrets := e.GetRequiredSecretNames(workflowData)
	filteredEnv := FilterEnvForSecrets(env, allowedSecrets)

	// Inject GH_TOKEN for CLI proxy (added after filtering since it uses a special
	// fallback expression that is always allowed when cli-proxy is enabled)
	addCliProxyGHTokenToEnv(filteredEnv, workflowData)

	// Format step with command and env
	stepLines = FormatStepWithCommandAndEnv(stepLines, command, filteredEnv)

	steps = append(steps, GitHubActionStep(stepLines))
	return steps
}
