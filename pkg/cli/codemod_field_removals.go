// This file consolidates small single-call codemod definitions.
//
// Each function below wraps a single call to newFieldRemovalCodemod,
// newMoveTopLevelKeyToOnBlockCodemod, or a plain Codemod literal. They are
// grouped here because each definition is purely data — one struct literal
// per function — and splitting them across individual files adds boilerplate
// without structural benefit.
//
// For codemods with non-trivial logic (custom Apply functions, multiple
// helper calls, or shared state), keep them in dedicated files.

package cli

import (
	"strings"

	"github.com/github/gh-aw/pkg/logger"
)

// ── Field-removal codemods ───────────────────────────────────────────────────

var mcpScriptsModeCodemodLog = logger.New("cli:codemod_mcp_scripts")

// getMCPScriptsModeCodemod creates a codemod for removing the deprecated mcp-scripts.mode field
func getMCPScriptsModeCodemod() Codemod {
	return newFieldRemovalCodemod(fieldRemovalCodemodConfig{
		ID:           "mcp-scripts-mode-removal",
		Name:         "Remove deprecated mcp-scripts.mode field",
		Description:  "Removes the deprecated 'mcp-scripts.mode' field (HTTP is now the only supported mode)",
		IntroducedIn: "0.2.0",
		ParentKey:    "mcp-scripts",
		FieldKey:     "mode",
		LogMsg:       "Applied mcp-scripts.mode removal",
		Log:          mcpScriptsModeCodemodLog,
	})
}

var grepToolCodemodLog = logger.New("cli:codemod_grep_tool")

// getGrepToolRemovalCodemod creates a codemod for removing the deprecated tools.grep field
func getGrepToolRemovalCodemod() Codemod {
	return newFieldRemovalCodemod(fieldRemovalCodemodConfig{
		ID:           "grep-tool-removal",
		Name:         "Remove deprecated tools.grep field",
		Description:  "Removes 'tools.grep' field as grep is now always enabled as part of default bash tools",
		IntroducedIn: "0.7.0",
		ParentKey:    "tools",
		FieldKey:     "grep",
		LogMsg:       "Applied grep tool removal",
		Log:          grepToolCodemodLog,
	})
}

var byokCopilotCodemodLog = logger.New("cli:codemod_byok_copilot")

// getByokCopilotFeatureRemovalCodemod removes deprecated features.byok-copilot.
func getByokCopilotFeatureRemovalCodemod() Codemod {
	return newFieldRemovalCodemod(fieldRemovalCodemodConfig{
		ID:           "features-byok-copilot-removal",
		Name:         "Remove deprecated features.byok-copilot",
		Description:  "Removes deprecated features.byok-copilot. Copilot now uses BYOK behavior by default.",
		IntroducedIn: "1.0.0",
		ParentKey:    "features",
		FieldKey:     "byok-copilot",
		LogMsg:       "Removed deprecated features.byok-copilot",
		Log:          byokCopilotCodemodLog,
	})
}

var inlineAgentsCodemodLog = logger.New("cli:codemod_inline_agents")

// getInlineAgentsFeatureRemovalCodemod removes deprecated features.inline-agents.
func getInlineAgentsFeatureRemovalCodemod() Codemod {
	return newFieldRemovalCodemod(fieldRemovalCodemodConfig{
		ID:           "features-inline-agents-removal",
		Name:         "Remove deprecated features.inline-agents",
		Description:  "Removes deprecated features.inline-agents. Inline sub-agents are now enabled by default.",
		IntroducedIn: "1.0.0",
		ParentKey:    "features",
		FieldKey:     "inline-agents",
		LogMsg:       "Removed deprecated features.inline-agents",
		Log:          inlineAgentsCodemodLog,
	})
}

// ── Move-to-on-block codemods ────────────────────────────────────────────────

var botsCodemodLog = logger.New("cli:codemod_bots")

// getBotsToOnBotsCodemod creates a codemod for moving top-level 'bots' to 'on.bots'
func getBotsToOnBotsCodemod() Codemod {
	return newMoveTopLevelKeyToOnBlockCodemod(moveToOnBlockConfig{
		ID:           "bots-to-on-bots",
		Name:         "Move bots to on.bots",
		Description:  "Moves the top-level 'bots' field to 'on.bots' as per the new frontmatter structure",
		IntroducedIn: "0.10.0",
		FieldKey:     "bots",
		IsInlineSingle: func(v string) bool {
			return strings.HasPrefix(v, "[")
		},
		Log: botsCodemodLog,
	})
}

var rolesCodemodLog = logger.New("cli:codemod_roles")

// getRolesToOnRolesCodemod creates a codemod for moving top-level 'roles' to 'on.roles'
func getRolesToOnRolesCodemod() Codemod {
	return newMoveTopLevelKeyToOnBlockCodemod(moveToOnBlockConfig{
		ID:           "roles-to-on-roles",
		Name:         "Move roles to on.roles",
		Description:  "Moves the top-level 'roles' field to 'on.roles' as per the new frontmatter structure",
		IntroducedIn: "0.10.0",
		FieldKey:     "roles",
		IsInlineSingle: func(v string) bool {
			return v == "all" || strings.HasPrefix(v, "[")
		},
		Log: rolesCodemodLog,
	})
}

// ── Miscellaneous single-definition codemods ─────────────────────────────────

// getDeleteSchemaFileCodemod creates a codemod for deleting deprecated schema files
func getDeleteSchemaFileCodemod() Codemod {
	return Codemod{
		ID:           "delete-schema-file",
		Name:         "Delete deprecated schema file",
		Description:  "Deletes .github/aw/schemas/agentic-workflow.json which is no longer written by init command",
		IntroducedIn: "0.6.0",
		Apply: func(content string, frontmatter map[string]any) (string, bool, error) {
			// This codemod is handled by the fix command itself (see runFixCommand)
			// It doesn't modify workflow files, so we just return content unchanged
			return content, false, nil
		},
	}
}
