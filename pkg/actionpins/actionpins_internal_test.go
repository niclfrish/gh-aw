//go:build !integration

package actionpins

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildByRepoIndex_GroupsByRepoAndSortsDescending(t *testing.T) {
	pins := []ActionPin{
		{Repo: "actions/checkout", Version: "v4.0.0", SHA: "sha-v4"},
		{Repo: "actions/checkout", Version: "v5.0.0", SHA: "sha-v5"},
		{Repo: "actions/setup-go", Version: "v5.1.0", SHA: "sha-go-v5-1"},
		{Repo: "actions/setup-go", Version: "v5.0.0", SHA: "sha-go-v5-0"},
	}

	byRepo := buildByRepoIndex(pins)

	require.Len(t, byRepo["actions/checkout"], 2, "Expected checkout pins to be grouped")
	assert.Equal(t, "v5.0.0", byRepo["actions/checkout"][0].Version, "Expected checkout pins sorted by newest version first")
	assert.Equal(t, "v4.0.0", byRepo["actions/checkout"][1].Version, "Expected checkout pins sorted by newest version first")

	require.Len(t, byRepo["actions/setup-go"], 2, "Expected setup-go pins to be grouped")
	assert.Equal(t, "v5.1.0", byRepo["actions/setup-go"][0].Version, "Expected setup-go pins sorted by newest version first")
	assert.Equal(t, "v5.0.0", byRepo["actions/setup-go"][1].Version, "Expected setup-go pins sorted by newest version first")
}

func TestCountPinKeyMismatches_ReturnsOnlyVersionMismatches(t *testing.T) {
	entries := map[string]ActionPin{
		"actions/checkout@v5": {Repo: "actions/checkout", Version: "v5", SHA: "sha-1"},
		"actions/setup-go@v5": {Repo: "actions/setup-go", Version: "v4", SHA: "sha-2"},
		"invalid-key":         {Repo: "actions/cache", Version: "v4", SHA: "sha-3"},
	}

	count := countPinKeyMismatches(entries)

	assert.Equal(t, 1, count, "Expected only one key/version mismatch to be counted")
}

func TestInitWarnings_InitializesAndPreservesMap(t *testing.T) {
	t.Run("initializes nil warnings map", func(t *testing.T) {
		ctx := &PinContext{}

		initWarnings(ctx)

		require.NotNil(t, ctx.Warnings, "Expected warnings map to be initialized")
		assert.Empty(t, ctx.Warnings, "Expected initialized warnings map to be empty")
	})

	t.Run("preserves existing warnings map", func(t *testing.T) {
		existing := map[string]bool{"actions/checkout@v5": true}
		ctx := &PinContext{Warnings: existing}

		initWarnings(ctx)

		require.NotNil(t, ctx.Warnings, "Expected warnings map to remain initialized")
		assert.Equal(t, existing, ctx.Warnings, "Expected existing warnings entries to be preserved")
	})
}

func TestGetContainerPin_ReturnsPinnedImage(t *testing.T) {
	pin, ok := GetContainerPin("node:lts-alpine")
	require.True(t, ok, "Expected embedded container pin for node:lts-alpine")
	assert.Equal(t, "node:lts-alpine", pin.Image, "Expected image name to match key")
	assert.NotEmpty(t, pin.Digest, "Expected digest to be populated")
	assert.Contains(t, pin.PinnedImage, "@sha256:", "Expected pinned image to include digest")
}

func TestDispatchResolutionFailure(t *testing.T) {
	t.Run("records failure when callback is present", func(t *testing.T) {
		var recorded ResolutionFailure
		called := false
		ctx := &PinContext{
			RecordResolutionFailure: func(f ResolutionFailure) {
				called = true
				recorded = f
			},
		}

		dispatchResolutionFailure(ctx, "actions/checkout", "v5", ResolutionErrorTypePinNotFound)

		require.True(t, called, "Expected resolution failure callback to be called")
		assert.Equal(t, "actions/checkout", recorded.Repo, "Expected repo to be forwarded to callback")
		assert.Equal(t, "v5", recorded.Ref, "Expected ref to be forwarded to callback")
		assert.Equal(t, ResolutionErrorTypePinNotFound, recorded.ErrorType, "Expected error type to be forwarded to callback")
	})

	t.Run("no-op when context or callback is nil", func(t *testing.T) {
		assert.NotPanics(t, func() {
			dispatchResolutionFailure(nil, "actions/checkout", "v5", ResolutionErrorTypePinNotFound)
		}, "Expected nil context to be safely ignored")

		assert.NotPanics(t, func() {
			dispatchResolutionFailure(&PinContext{}, "actions/checkout", "v5", ResolutionErrorTypePinNotFound)
		}, "Expected missing callback to be safely ignored")
	})
}
