//go:build !integration

package workflow

import (
	"strings"
	"testing"

	"github.com/github/gh-aw/pkg/constants"
)

func TestGetVersionForSetup(t *testing.T) {
	tests := []struct {
		name            string
		data            *WorkflowData
		expectedVersion string
	}{
		{
			name:            "nil data returns empty string",
			data:            nil,
			expectedVersion: "",
		},
		{
			name:            "no engine config returns empty string",
			data:            &WorkflowData{},
			expectedVersion: "",
		},
		{
			name: "explicit version in EngineConfig takes priority",
			data: &WorkflowData{
				EngineConfig: &EngineConfig{ID: "copilot", Version: "1.2.3"},
			},
			expectedVersion: "1.2.3",
		},
		{
			name: "copilot engine uses default version",
			data: &WorkflowData{
				EngineConfig: &EngineConfig{ID: "copilot"},
			},
			expectedVersion: string(constants.DefaultCopilotVersion),
		},
		{
			name: "claude engine uses default version",
			data: &WorkflowData{
				EngineConfig: &EngineConfig{ID: "claude"},
			},
			expectedVersion: string(constants.DefaultClaudeCodeVersion),
		},
		{
			name: "codex engine uses default version",
			data: &WorkflowData{
				EngineConfig: &EngineConfig{ID: "codex"},
			},
			expectedVersion: string(constants.DefaultCodexVersion),
		},
		{
			name: "AI field used when EngineConfig.ID is empty",
			data: &WorkflowData{
				AI: "copilot",
			},
			expectedVersion: string(constants.DefaultCopilotVersion),
		},
		{
			name: "EngineConfig.ID takes priority over AI field",
			data: &WorkflowData{
				AI:           "copilot",
				EngineConfig: &EngineConfig{ID: "claude"},
			},
			expectedVersion: string(constants.DefaultClaudeCodeVersion),
		},
		{
			name: "custom engine returns empty string",
			data: &WorkflowData{
				EngineConfig: &EngineConfig{ID: "custom"},
			},
			expectedVersion: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getVersionForSetup(tt.data)
			if result != tt.expectedVersion {
				t.Errorf("getVersionForSetup() = %q, want %q", result, tt.expectedVersion)
			}
		})
	}
}

func TestGenerateSetupStepIncludesVersion(t *testing.T) {
	tests := []struct {
		name          string
		data          *WorkflowData
		expectVersion string
		noVersionLine bool
	}{
		{
			name: "copilot engine injects default version",
			data: &WorkflowData{
				Name:         "my-workflow",
				EngineConfig: &EngineConfig{ID: "copilot"},
			},
			expectVersion: string(constants.DefaultCopilotVersion),
		},
		{
			name: "explicit version is injected",
			data: &WorkflowData{
				Name:         "my-workflow",
				EngineConfig: &EngineConfig{ID: "copilot", Version: "1.2.3"},
			},
			expectVersion: "1.2.3",
		},
		{
			name: "custom engine without version does not inject GH_AW_INFO_VERSION",
			data: &WorkflowData{
				Name:         "my-workflow",
				EngineConfig: &EngineConfig{ID: "custom"},
			},
			noVersionLine: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCompiler()
			lines := c.generateSetupStep(tt.data, "github/gh-aw/actions/setup@abc123", "${{ runner.temp }}/gh-aw", false, "", "")
			combined := strings.Join(lines, "")

			if tt.noVersionLine {
				if strings.Contains(combined, "GH_AW_INFO_VERSION") {
					t.Errorf("expected no GH_AW_INFO_VERSION in setup step, but found it:\n%s", combined)
				}
				return
			}

			expectedLine := `GH_AW_INFO_VERSION: "` + tt.expectVersion + `"`
			if !strings.Contains(combined, expectedLine) {
				t.Errorf("expected setup step to contain %q, got:\n%s", expectedLine, combined)
			}
		})
	}
}

func TestGenerateSetupStepIncludesParentSpanID(t *testing.T) {
	c := NewCompiler()
	data := &WorkflowData{Name: "my-workflow"}
	parentExpr := "${{ needs.activation.outputs.setup-span-id }}"

	lines := c.generateSetupStep(data, "github/gh-aw/actions/setup@abc123", "${{ runner.temp }}/gh-aw", false, "", parentExpr)
	combined := strings.Join(lines, "")

	if !strings.Contains(combined, "parent-span-id: "+parentExpr) {
		t.Fatalf("expected setup step to include parent-span-id input, got:\n%s", combined)
	}
}

func TestGenerateSetupStepIncludesEngineID(t *testing.T) {
	t.Run("action mode injects engine id", func(t *testing.T) {
		c := NewCompiler()
		data := &WorkflowData{
			Name:         "my-workflow",
			EngineConfig: &EngineConfig{ID: "codex"},
		}
		lines := c.generateSetupStep(data, "github/gh-aw/actions/setup@abc123", "${{ runner.temp }}/gh-aw", false, "", "")
		combined := strings.Join(lines, "")

		if !strings.Contains(combined, `GH_AW_INFO_ENGINE_ID: "codex"`) {
			t.Fatalf("expected setup step to include GH_AW_INFO_ENGINE_ID for action mode, got:\n%s", combined)
		}
	})

	t.Run("script mode injects engine id", func(t *testing.T) {
		c := NewCompiler()
		c.actionMode = ActionModeScript
		data := &WorkflowData{
			Name: "my-workflow",
			AI:   "copilot",
		}
		lines := c.generateSetupStep(data, "./actions/setup", "${{ runner.temp }}/gh-aw", false, "", "")
		combined := strings.Join(lines, "")

		if !strings.Contains(combined, `GH_AW_INFO_ENGINE_ID: "copilot"`) {
			t.Fatalf("expected setup step to include GH_AW_INFO_ENGINE_ID for script mode, got:\n%s", combined)
		}
	})
}

func TestGetEngineIDForSetup(t *testing.T) {
	t.Run("reads engine from raw frontmatter string", func(t *testing.T) {
		data := &WorkflowData{
			RawFrontmatter: map[string]any{"engine": "claude"},
			EngineConfig:   &EngineConfig{ID: "copilot"},
			AI:             "copilot",
		}
		if got := getEngineIDForSetup(data); got != "claude" {
			t.Fatalf("expected raw frontmatter engine to win, got %q", got)
		}
	})

	t.Run("reads engine.id from raw frontmatter object", func(t *testing.T) {
		data := &WorkflowData{
			RawFrontmatter: map[string]any{"engine": map[string]any{"id": "codex"}},
			EngineConfig:   &EngineConfig{ID: "copilot"},
		}
		if got := getEngineIDForSetup(data); got != "codex" {
			t.Fatalf("expected raw frontmatter engine.id, got %q", got)
		}
	})

	t.Run("reads engine.runtime.id from raw frontmatter object", func(t *testing.T) {
		data := &WorkflowData{
			RawFrontmatter: map[string]any{"engine": map[string]any{"runtime": map[string]any{"id": "custom-engine"}}},
			EngineConfig:   &EngineConfig{ID: "copilot"},
		}
		if got := getEngineIDForSetup(data); got != "custom-engine" {
			t.Fatalf("expected raw frontmatter engine.runtime.id, got %q", got)
		}
	})

	t.Run("falls back to merged engine fields when raw frontmatter has no engine", func(t *testing.T) {
		data := &WorkflowData{
			RawFrontmatter: map[string]any{"name": "wf"},
			EngineConfig:   &EngineConfig{ID: "copilot"},
			AI:             "claude",
		}
		if got := getEngineIDForSetup(data); got != "copilot" {
			t.Fatalf("expected EngineConfig fallback, got %q", got)
		}
	})
}
