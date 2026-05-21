package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

var version = "dev"

// rootCmd is the base command for gh-aw
var rootCmd = &cobra.Command{
	Use:     "gh-aw",
	Short:   "GitHub Actions Workflow automation tool",
	Long:    `gh-aw is a GitHub CLI extension for managing and automating GitHub Actions workflows.`,
	Version: version,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// newAPIClient creates a new GitHub API client, respecting GHE configuration.
// Note: passing api.EnableLog(os.Stderr) as an option is handy for local debugging.
func newAPIClient(opts ...api.ClientOption) (*api.RESTClient, error) {
	// Default to a reasonable timeout; can be overridden by callers via opts.
	defaultOpts := []api.ClientOption{
		api.AddHeader("X-Custom-Client", "gh-aw/"+version),
		// Increased timeout from default to handle slow GHE instances
		api.AddHeader("X-Request-Timeout", "30"),
		api.WithTimeout(30 * time.Second),
	}
	client, err := api.NewRESTClient(append(defaultOpts, opts...)...)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}
	return client, nil
}
