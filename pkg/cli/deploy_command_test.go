//go:build !integration

package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDeployCommand_BasicShape(t *testing.T) {
	cmd := NewDeployCommand(func(string) error { return nil })
	require.NotNil(t, cmd)
	assert.Equal(t, "deploy <workflow>...", cmd.Use)
	assert.Equal(t, "deploy", cmd.Name())
}

func TestNewDeployCommand_RequiresWorkflowArg(t *testing.T) {
	cmd := NewDeployCommand(func(string) error { return nil })
	require.NotNil(t, cmd)

	err := cmd.Args(cmd, []string{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing workflow specification")
}

func TestNewDeployCommand_RegistersCoreFlags(t *testing.T) {
	cmd := NewDeployCommand(func(string) error { return nil })
	require.NotNil(t, cmd)

	expectedFlags := []string{
		"repo",
		"name",
		"engine",
		"force",
		"append",
		"no-gitattributes",
		"dir",
		"no-stop-after",
		"stop-after",
		"disable-security-scanner",
	}

	for _, flagName := range expectedFlags {
		t.Run(flagName, func(t *testing.T) {
			flag := cmd.Flags().Lookup(flagName)
			require.NotNil(t, flag, "expected flag %q to be registered", flagName)
		})
	}
}
