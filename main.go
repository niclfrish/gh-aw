package main

import (
	"fmt"
	"os"

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

// newAPIClient creates a new GitHub API client, respecting GHE configuration
func newAPIClient(opts ...api.ClientOption) (*api.RESTClient, error) {
	client, err := api.NewRESTClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}
	return client, nil
}
