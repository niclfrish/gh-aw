//go:build !integration

package cli

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// connectInProcess creates an in-process MCP server/client pair using in-memory transports.
// The returned session must be closed by the caller.
func connectInProcess(t *testing.T, server *mcp.Server) *mcp.ClientSession {
	t.Helper()
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0.0"}, nil)
	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	ctx := context.Background()
	_, err := server.Connect(ctx, serverTransport, nil)
	require.NoError(t, err, "server.Connect should succeed")

	session, err := client.Connect(ctx, clientTransport, nil)
	require.NoError(t, err, "client.Connect should succeed")
	t.Cleanup(func() { session.Close() })
	return session
}

// TestStatusToolStructuredContent verifies that the status tool returns
// non-nil structuredContent with the typed StatusOutput shape.
func TestStatusToolStructuredContent(t *testing.T) {
	// Create a temp dir with a workflow file so GetWorkflowStatuses has something to return.
	tmpDir := t.TempDir()
	workflowsDir := filepath.Join(tmpDir, ".github", "workflows")
	require.NoError(t, os.MkdirAll(workflowsDir, 0755), "create workflows dir")

	workflowContent := "---\non: push\nengine: copilot\n---\n# Test Workflow\n"
	require.NoError(t,
		os.WriteFile(filepath.Join(workflowsDir, "test.md"), []byte(workflowContent), 0644),
		"write workflow file",
	)

	// Change working directory so GetWorkflowStatuses discovers the temp workflow files.
	oldDir, err := os.Getwd()
	require.NoError(t, err, "get working directory")
	require.NoError(t, os.Chdir(tmpDir), "change to temp dir")
	t.Cleanup(func() { os.Chdir(oldDir) })

	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	registerStatusTool(server)

	session := connectInProcess(t, server)

	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "status",
		Arguments: map[string]any{},
	})
	require.NoError(t, err, "CallTool should succeed")

	// Text content must be present and non-empty (backward compatibility).
	require.NotEmpty(t, result.Content, "result should have content")
	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok, "first content should be TextContent")
	assert.NotEmpty(t, textContent.Text, "text content should not be empty")

	// Structured content must be present and parseable as StatusOutput.
	assert.NotNil(t, result.StructuredContent, "StructuredContent should be non-nil")

	// Round-trip via JSON to verify the shape.
	raw, marshalErr := json.Marshal(result.StructuredContent)
	require.NoError(t, marshalErr, "marshal StructuredContent")

	var output StatusOutput
	require.NoError(t, json.Unmarshal(raw, &output), "unmarshal StructuredContent as StatusOutput")
	// The temp dir has one .md file, so Workflows should contain at least one entry.
	assert.NotEmpty(t, output.Workflows, "Workflows should contain at least one entry")
}

// TestCompileToolStructuredContent verifies that the compile tool returns
// non-nil structuredContent with the typed CompileOutput shape when the
// underlying compile command succeeds and produces valid JSON.
func TestCompileToolStructuredContent(t *testing.T) {
	// Mock execCmd that writes valid JSON compile output to stdout.
	jsonOutput := `[{"workflow":"test.md","valid":true,"errors":[],"warnings":[]}]`
	mockExecCmd := func(ctx context.Context, args ...string) *exec.Cmd {
		// Use "echo" to write the JSON output to stdout, simulating a successful compile.
		return exec.CommandContext(ctx, "echo", jsonOutput)
	}

	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	err := registerCompileTool(server, mockExecCmd, "")
	require.NoError(t, err, "registerCompileTool should succeed")

	session := connectInProcess(t, server)

	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "compile",
		Arguments: map[string]any{},
	})
	require.NoError(t, err, "CallTool should succeed")

	// Text content must be present (backward compatibility).
	require.NotEmpty(t, result.Content, "result should have content")

	// Structured content must be present.
	assert.NotNil(t, result.StructuredContent, "StructuredContent should be non-nil")

	// Verify the shape is CompileOutput.
	raw, marshalErr := json.Marshal(result.StructuredContent)
	require.NoError(t, marshalErr, "marshal StructuredContent")

	var output CompileOutput
	require.NoError(t, json.Unmarshal(raw, &output), "unmarshal StructuredContent as CompileOutput")
}

// TestCompileToolStructuredContent_WithResults verifies that when the compile
// command returns a valid JSON array, the structured output contains parsed results.
func TestCompileToolStructuredContent_WithResults(t *testing.T) {
	jsonOutput := `[{"workflow":"a.md","valid":true,"errors":[],"warnings":[]},{"workflow":"b.md","valid":false,"errors":[{"type":"parse_error","message":"syntax error"}],"warnings":[]}]`
	mockExecCmd := func(ctx context.Context, args ...string) *exec.Cmd {
		return exec.CommandContext(ctx, "echo", jsonOutput)
	}

	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	require.NoError(t, registerCompileTool(server, mockExecCmd, ""), "registerCompileTool")

	session := connectInProcess(t, server)

	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "compile",
		Arguments: map[string]any{},
	})
	require.NoError(t, err, "CallTool should succeed")

	assert.NotNil(t, result.StructuredContent, "StructuredContent should be non-nil")

	raw, marshalErr := json.Marshal(result.StructuredContent)
	require.NoError(t, marshalErr, "marshal StructuredContent")

	var output CompileOutput
	require.NoError(t, json.Unmarshal(raw, &output), "unmarshal StructuredContent as CompileOutput")
	assert.Len(t, output.Results, 2, "should have 2 results")
	assert.Equal(t, "a.md", output.Results[0].Workflow, "first result workflow name")
	assert.True(t, output.Results[0].Valid, "first result should be valid")
	assert.Equal(t, "b.md", output.Results[1].Workflow, "second result workflow name")
	assert.False(t, output.Results[1].Valid, "second result should be invalid")

	// Text content must also match (backward compatibility).
	require.NotEmpty(t, result.Content, "text content must be present")
}

// TestStatusOutput_JSONShape verifies that StatusOutput serialises to the expected JSON shape.
func TestStatusOutput_JSONShape(t *testing.T) {
	out := StatusOutput{
		Workflows: []WorkflowStatus{
			{
				Workflow:      "my-workflow.md",
				EngineID:      "copilot",
				Compiled:      "Yes",
				Status:        "active",
				TimeRemaining: "",
			},
		},
	}

	data, err := json.Marshal(out)
	require.NoError(t, err, "marshal StatusOutput")

	var m map[string]any
	require.NoError(t, json.Unmarshal(data, &m), "unmarshal as map")

	workflows, ok := m["workflows"].([]any)
	require.True(t, ok, "workflows field should be a JSON array")
	assert.Len(t, workflows, 1, "should contain one workflow")
}

// TestCompileOutput_JSONShape verifies that CompileOutput serialises to the expected JSON shape.
func TestCompileOutput_JSONShape(t *testing.T) {
	out := CompileOutput{
		Results: []ValidationResult{
			{
				Workflow: "workflow.md",
				Valid:    true,
				Errors:   []CompileValidationError{},
				Warnings: []CompileValidationError{},
			},
		},
	}

	data, err := json.Marshal(out)
	require.NoError(t, err, "marshal CompileOutput")

	var m map[string]any
	require.NoError(t, json.Unmarshal(data, &m), "unmarshal as map")

	results, ok := m["results"].([]any)
	require.True(t, ok, "results field should be a JSON array")
	assert.Len(t, results, 1, "should contain one result")
}
