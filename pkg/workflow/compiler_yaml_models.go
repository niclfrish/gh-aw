package workflow

import (
	"fmt"
	"strings"

	"github.com/github/gh-aw/pkg/constants"
)

// generateEngineModelsCollectionStep emits a step that fetches the engine's available models
// when the engine exposes a models route.
func (c *Compiler) generateEngineModelsCollectionStep(yaml *strings.Builder, data *WorkflowData, engine CodingAgentEngine) bool {
	routeProvider, ok := engine.(ModelsRouteProvider)
	if !ok {
		return false
	}

	modelsRoute := strings.TrimSpace(routeProvider.GetModelsRoute())
	if modelsRoute == "" {
		return false
	}

	modelsHost := getEngineModelsHost(engine, data)
	if modelsHost == "" {
		return false
	}

	yaml.WriteString("      - name: Collect available models\n")
	yaml.WriteString("        if: always()\n")
	yaml.WriteString("        continue-on-error: true\n")
	fmt.Fprintf(yaml, "        uses: %s\n", getCachedActionPin("actions/github-script", data))
	yaml.WriteString("        env:\n")
	yaml.WriteString("          GH_AW_MODELS_FILE: ${{ runner.temp }}/gh-aw/models.json\n")
	yaml.WriteString("          GH_AW_MODELS_ARTIFACT_FILE: /tmp/gh-aw/models.json\n")
	writeYAMLEnv(yaml, "          ", "GH_AW_ENGINE_ID", engine.GetID())
	writeYAMLEnv(yaml, "          ", "GH_AW_MODELS_HOST", modelsHost)
	writeYAMLEnv(yaml, "          ", "GH_AW_MODELS_ROUTE", modelsRoute)
	writeYAMLEnv(yaml, "          ", "GH_AW_MODELS_AUTH_TYPE", getEngineModelsAuthType(engine))
	modelToken := getEngineModelsTokenExpression(engine, data)
	if modelToken != "" {
		writeYAMLEnv(yaml, "          ", "GH_AW_MODELS_TOKEN", modelToken)
	}
	yaml.WriteString("        with:\n")
	yaml.WriteString("          script: |\n")
	yaml.WriteString("            const { setupGlobals } = require('${{ runner.temp }}/gh-aw/actions/setup_globals.cjs');\n")
	yaml.WriteString("            setupGlobals(core, github, context, exec, io, getOctokit);\n")
	yaml.WriteString("            const { main } = require('${{ runner.temp }}/gh-aw/actions/collect_models.cjs');\n")
	yaml.WriteString("            await main();\n")
	return true
}

func getEngineModelsHost(engine CodingAgentEngine, workflowData *WorkflowData) string {
	switch engine.GetID() {
	case "copilot":
		if host := GetCopilotAPITarget(workflowData); host != "" {
			return host
		}
		return "api.githubcopilot.com"
	case "claude":
		if host := extractAPITargetHost(workflowData, "ANTHROPIC_BASE_URL"); host != "" {
			return host
		}
		return "api.anthropic.com"
	case "codex":
		if host := extractAPITargetHost(workflowData, "OPENAI_BASE_URL"); host != "" {
			return host
		}
		return "api.openai.com"
	case "gemini":
		if host := GetGeminiAPITarget(workflowData, engine.GetID()); host != "" {
			return host
		}
		return ""
	case "crush":
		if host := extractAPITargetHost(workflowData, "OPENAI_BASE_URL"); host != "" {
			return host
		}
		return "api.openai.com"
	default:
		return ""
	}
}

func getEngineModelsTokenExpression(engine CodingAgentEngine, workflowData *WorkflowData) string {
	switch engine.GetID() {
	case "copilot":
		if isFeatureEnabled(constants.CopilotRequestsFeatureFlag, workflowData) {
			return "${{ github.token }}"
		}
		return "${{ secrets.COPILOT_GITHUB_TOKEN }}"
	case "claude":
		return "${{ secrets.ANTHROPIC_API_KEY }}"
	case "codex":
		return "${{ secrets.CODEX_API_KEY || secrets.OPENAI_API_KEY }}"
	case "gemini":
		return "${{ secrets.GEMINI_API_KEY }}"
	case "crush":
		if isFeatureEnabled(constants.CopilotRequestsFeatureFlag, workflowData) {
			return "${{ github.token }}"
		}
		return "${{ secrets.COPILOT_GITHUB_TOKEN }}"
	default:
		return ""
	}
}

func getEngineModelsAuthType(engine CodingAgentEngine) string {
	switch engine.GetID() {
	case "claude":
		return "x-api-key"
	case "gemini":
		return "x-goog-api-key"
	default:
		return "bearer"
	}
}
