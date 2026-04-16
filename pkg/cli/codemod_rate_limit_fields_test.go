//go:build !integration

package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateLimitFieldsCodemod(t *testing.T) {
	codemod := getRateLimitFieldsCodemod()

	t.Run("renames max and window under rate-limit", func(t *testing.T) {
		content := `---
engine: copilot
rate-limit:
  max: 5
  window: 60
---

# Test Workflow
`
		frontmatter := map[string]any{
			"engine": "copilot",
			"rate-limit": map[string]any{
				"max":    5,
				"window": 60,
			},
		}

		result, applied, err := codemod.Apply(content, frontmatter)
		require.NoError(t, err, "Should not error")
		assert.True(t, applied, "Should have applied the codemod")
		assert.Contains(t, result, "max-runs: 5", "Should rename max to max-runs")
		assert.Contains(t, result, "max-runs-window: 60", "Should rename window to max-runs-window")
		assert.NotContains(t, result, "\n  max: ", "Should not contain old max field")
		assert.NotContains(t, result, "\n  window: ", "Should not contain old window field")
	})

	t.Run("renames only max when max-runs-window already exists", func(t *testing.T) {
		content := `---
engine: copilot
rate-limit:
  max: 5
  max-runs-window: 60
---

# Test Workflow
`
		frontmatter := map[string]any{
			"engine": "copilot",
			"rate-limit": map[string]any{
				"max":             5,
				"max-runs-window": 60,
			},
		}

		result, applied, err := codemod.Apply(content, frontmatter)
		require.NoError(t, err, "Should not error")
		assert.True(t, applied, "Should have applied the codemod")
		assert.Contains(t, result, "max-runs: 5", "Should rename max to max-runs")
		assert.NotContains(t, result, "\n  max: ", "Should not contain old max field")
		assert.Contains(t, result, "max-runs-window: 60", "Should keep max-runs-window unchanged")
	})

	t.Run("does not modify already migrated fields", func(t *testing.T) {
		content := `---
engine: copilot
rate-limit:
  max-runs: 5
  max-runs-window: 60
---

# Test Workflow
`
		frontmatter := map[string]any{
			"engine": "copilot",
			"rate-limit": map[string]any{
				"max-runs":        5,
				"max-runs-window": 60,
			},
		}

		result, applied, err := codemod.Apply(content, frontmatter)
		require.NoError(t, err, "Should not error")
		assert.False(t, applied, "Should not have applied the codemod")
		assert.Equal(t, content, result, "Content should be unchanged")
	})
}
