package workflow

import (
	"encoding/json"

	"github.com/github/gh-aw/pkg/logger"
	"github.com/github/gh-aw/pkg/parser"
)

var workflowGitHubAppLog = logger.New("workflow:workflow_github_app")

// resolveGitHubAppFromImports parses a GitHubAppConfig from a JSON-encoded
// merged import string. Returns nil if the string is empty, unparseable, or incomplete.
func resolveGitHubAppFromImports(mergedJSON string) *GitHubAppConfig {
	if mergedJSON == "" {
		return nil
	}
	var appMap map[string]any
	if err := json.Unmarshal([]byte(mergedJSON), &appMap); err != nil {
		return nil
	}
	app := parseAppConfig(appMap)
	if app.AppID != "" && app.PrivateKey != "" {
		return app
	}
	return nil
}

// extractTopLevelGitHubApp extracts the 'github-app' field from the top-level frontmatter.
// This provides a single GitHub App configuration that serves as a fallback for all nested
// github-app token minting operations (on, safe-outputs, checkout, tools.github, dependencies).
func extractTopLevelGitHubApp(frontmatter map[string]any) *GitHubAppConfig {
	appAny, ok := frontmatter["github-app"]
	if !ok {
		return nil
	}
	appMap, ok := appAny.(map[string]any)
	if !ok {
		return nil
	}
	app := parseAppConfig(appMap)
	if app.AppID == "" || app.PrivateKey == "" {
		return nil
	}
	return app
}

// resolveTopLevelGitHubApp resolves the top-level github-app for token minting fallback.
// Precedence:
//  1. Current workflow's top-level github-app (explicit override wins)
//  2. First top-level github-app found across imported shared workflows
//  3. Nil (no fallback configured)
func resolveTopLevelGitHubApp(frontmatter map[string]any, importsResult *parser.ImportsResult) *GitHubAppConfig {
	if app := extractTopLevelGitHubApp(frontmatter); app != nil {
		return app
	}
	if importsResult != nil {
		if app := resolveGitHubAppFromImports(importsResult.MergedTopLevelGitHubApp); app != nil {
			workflowGitHubAppLog.Print("Using top-level github-app from imported shared workflow")
			return app
		}
	}
	return nil
}

// topLevelFallbackNeeded reports whether the top-level github-app should be applied as a
// fallback for a given section. It returns true when the section has neither an explicit
// github-app nor an explicit github-token already configured.
//
// Rules (consistent across all sections):
//   - If a section-specific github-app is set → keep it, no fallback needed.
//   - If a section-specific github-token is set → keep it, no fallback needed (a token
//     already provides the auth; injecting a github-app would silently change precedence).
//   - Otherwise → apply the top-level fallback.
func topLevelFallbackNeeded(app *GitHubAppConfig, token string) bool {
	return app == nil && token == ""
}

// applyTopLevelGitHubAppFallbacks applies the top-level github-app as a fallback for all
// nested github-app token minting operations when no section-specific github-app is configured.
// Precedence: section-specific github-app > section-specific github-token > top-level github-app.
//
// Every section uses topLevelFallbackNeeded to decide whether the fallback is required,
// ensuring consistent behaviour across all token-minting sites.
func applyTopLevelGitHubAppFallbacks(data *WorkflowData) {
	fallback := data.TopLevelGitHubApp
	if fallback == nil {
		return
	}

	// Fallback for activation (on.github-app / on.github-token)
	if topLevelFallbackNeeded(data.ActivationGitHubApp, data.ActivationGitHubToken) {
		workflowGitHubAppLog.Print("Applying top-level github-app fallback for activation")
		data.ActivationGitHubApp = fallback
	}

	// Fallback for safe-outputs (safe-outputs.github-app / safe-outputs.github-token)
	if data.SafeOutputs != nil && topLevelFallbackNeeded(data.SafeOutputs.GitHubApp, data.SafeOutputs.GitHubToken) {
		workflowGitHubAppLog.Print("Applying top-level github-app fallback for safe-outputs")
		data.SafeOutputs.GitHubApp = fallback
	}

	// Fallback for checkout configs (checkout.github-app / checkout.github-token per entry)
	for _, cfg := range data.CheckoutConfigs {
		if topLevelFallbackNeeded(cfg.GitHubApp, cfg.GitHubToken) {
			workflowGitHubAppLog.Print("Applying top-level github-app fallback for checkout")
			cfg.GitHubApp = fallback
		}
	}

	// Fallback for tools.github (tools.github.github-app / tools.github.github-token).
	// Also skipped when tools.github is explicitly disabled (github: false) — do not re-enable it.
	if data.ParsedTools != nil && data.ParsedTools.GitHub != nil &&
		topLevelFallbackNeeded(data.ParsedTools.GitHub.GitHubApp, data.ParsedTools.GitHub.GitHubToken) &&
		data.Tools["github"] != false {
		workflowGitHubAppLog.Print("Applying top-level github-app fallback for tools.github")
		data.ParsedTools.GitHub.GitHubApp = fallback
		// Also update the raw tools map so applyDefaultTools (called from applyDefaults in
		// processOnSectionAndFilters) does not lose the fallback when it rebuilds ParsedTools
		// from the map.
		appMap := map[string]any{
			"client-id":   fallback.AppID,
			"private-key": fallback.PrivateKey,
		}
		if fallback.Owner != "" {
			appMap["owner"] = fallback.Owner
		}
		if len(fallback.Repositories) > 0 {
			repos := make([]any, len(fallback.Repositories))
			for i, r := range fallback.Repositories {
				repos[i] = r
			}
			appMap["repositories"] = repos
		}
		// Normalize data.Tools["github"] to a map so the github-app survives re-parsing.
		// Configurations like `github: true` are normalized here rather than losing the fallback.
		if github, ok := data.Tools["github"].(map[string]any); ok {
			// Already a map; inject into existing settings.
			github["github-app"] = appMap
		} else {
			// Non-map value (e.g. true) — create a fresh map.
			data.Tools["github"] = map[string]any{"github-app": appMap}
		}
	}
}

// extractActivationGitHubToken extracts the 'github-token' field from the 'on:' section of frontmatter.
// This token is used for pre-activation reactions and activation status comments.
func (c *Compiler) extractActivationGitHubToken(frontmatter map[string]any) string {
	if onValue, exists := frontmatter["on"]; exists {
		if onMap, ok := onValue.(map[string]any); ok {
			if tokenValue, hasToken := onMap["github-token"]; hasToken {
				if tokenStr, ok := tokenValue.(string); ok {
					workflowGitHubAppLog.Printf("Extracted activation github-token from on section")
					return tokenStr
				}
			}
		}
	}
	return ""
}

// extractActivationGitHubApp extracts the 'github-app' field from the 'on:' section of frontmatter.
// When configured, a GitHub App installation access token is minted for use in reactions and status comments.
func (c *Compiler) extractActivationGitHubApp(frontmatter map[string]any) *GitHubAppConfig {
	if onValue, exists := frontmatter["on"]; exists {
		if onMap, ok := onValue.(map[string]any); ok {
			if appValue, hasApp := onMap["github-app"]; hasApp {
				if appMap, ok := appValue.(map[string]any); ok {
					workflowGitHubAppLog.Printf("Extracted activation github-app from on section")
					return parseAppConfig(appMap)
				}
			}
		}
	}
	return nil
}

// resolveActivationGitHubToken returns the GitHub token to use for activation operations
// (reactions, status comments, skip-if checks). Precedence:
//  1. Current workflow's on.github-token (explicit override wins)
//  2. First on.github-token found across imported shared workflows
//  3. Empty string (use default GITHUB_TOKEN)
func (c *Compiler) resolveActivationGitHubToken(frontmatter map[string]any, importsResult *parser.ImportsResult) string {
	if token := c.extractActivationGitHubToken(frontmatter); token != "" {
		return token
	}
	if importsResult != nil && importsResult.MergedActivationGitHubToken != "" {
		workflowGitHubAppLog.Print("Using on.github-token from imported shared workflow")
		return importsResult.MergedActivationGitHubToken
	}
	return ""
}

// resolveActivationGitHubApp returns the GitHub App config to use for activation operations
// (reactions, status comments, skip-if checks). Precedence:
//  1. Current workflow's on.github-app (explicit override wins)
//  2. First on.github-app found across imported shared workflows
//  3. Nil (use default GITHUB_TOKEN)
func (c *Compiler) resolveActivationGitHubApp(frontmatter map[string]any, importsResult *parser.ImportsResult) *GitHubAppConfig {
	if app := c.extractActivationGitHubApp(frontmatter); app != nil {
		return app
	}
	if importsResult != nil {
		if app := resolveGitHubAppFromImports(importsResult.MergedActivationGitHubApp); app != nil {
			workflowGitHubAppLog.Print("Using on.github-app from imported shared workflow")
			return app
		}
	}
	return nil
}
