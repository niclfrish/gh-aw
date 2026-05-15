package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/github/gh-aw/pkg/constants"
	"github.com/github/gh-aw/pkg/gitutil"
	"github.com/github/gh-aw/pkg/logger"
	"github.com/spf13/cobra"
)

var deployLog = logger.New("cli:deploy_command")

// NewDeployCommand creates the deploy command.
func NewDeployCommand(validateEngine func(string) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy <workflow>...",
		Short: "Deploy agentic workflows to a target repository using a pull request",
		Long: `Deploy one or more workflows to a target repository by combining clone, update, add, compile, and pull request creation.

The command clones the target repository, updates existing workflows from source, adds the specified workflows, recompiles lock files with purge enabled, and opens a pull request.

Examples:
  ` + string(constants.CLIExtensionPrefix) + ` deploy githubnext/agentics/ci-doctor --repo owner/repo
  ` + string(constants.CLIExtensionPrefix) + ` deploy githubnext/agentics/repo-assist githubnext/agentics/ci-doctor --repo owner/repo --force
  ` + string(constants.CLIExtensionPrefix) + ` deploy ./my-workflow.md --repo owner/repo`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("missing workflow specification\n\nUsage:\n  %s <workflow>...\n\nExamples:\n  %[1]s githubnext/agentics/ci-doctor --repo owner/repo\n\nRun '%[1]s --help' for more information", cmd.CommandPath())
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			workflows := args
			targetRepo, _ := cmd.Flags().GetString("repo")
			if strings.TrimSpace(targetRepo) == "" {
				return errors.New("--repo flag is required (target repository in owner/repo format)")
			}

			engineOverride, _ := cmd.Flags().GetString("engine")
			nameFlag, _ := cmd.Flags().GetString("name")
			forceFlag, _ := cmd.Flags().GetBool("force")
			appendText, _ := cmd.Flags().GetString("append")
			verbose, _ := cmd.Flags().GetBool("verbose")
			noGitattributes, _ := cmd.Flags().GetBool("no-gitattributes")
			workflowDir, _ := cmd.Flags().GetString("dir")
			noStopAfter, _ := cmd.Flags().GetBool("no-stop-after")
			stopAfter, _ := cmd.Flags().GetString("stop-after")
			disableSecurityScanner, _ := cmd.Flags().GetBool("disable-security-scanner")

			if nameFlag != "" && len(workflows) > 1 {
				return errors.New("--name flag cannot be used when adding multiple workflows at once")
			}

			if err := validateEngine(engineOverride); err != nil {
				return err
			}

			coolDown, err := parseCoolDownFlag("7d")
			if err != nil {
				return err
			}

			opts := AddOptions{
				Verbose:                verbose,
				EngineOverride:         engineOverride,
				Name:                   nameFlag,
				Force:                  forceFlag,
				AppendText:             appendText,
				NoGitattributes:        noGitattributes,
				WorkflowDir:            workflowDir,
				NoStopAfter:            noStopAfter,
				StopAfter:              stopAfter,
				DisableSecurityScanner: disableSecurityScanner,
			}

			return runDeploy(cmd.Context(), targetRepo, workflows, opts, coolDown)
		},
	}

	addRepoFlag(cmd)
	cmd.Flags().StringP("name", "n", "", "Specify name for the added workflow (without .md extension)")
	addEngineFlag(cmd)
	cmd.Flags().BoolP("force", "f", false, "Overwrite existing workflow files without confirmation")
	cmd.Flags().String("append", "", "Append extra content to the end of agentic workflow on installation")
	cmd.Flags().Bool("no-gitattributes", false, "Skip updating .gitattributes file")
	cmd.Flags().StringP("dir", "d", "", "Workflow directory (default: .github/workflows)")
	cmd.Flags().Bool("no-stop-after", false, "Remove any stop-after field from the workflow")
	cmd.Flags().String("stop-after", "", "Override stop-after value in the workflow (e.g., '+48h', '2025-12-31 23:59:59')")
	cmd.Flags().Bool("disable-security-scanner", false, "Disable security scanning of workflow markdown content")

	RegisterEngineFlagCompletion(cmd)
	RegisterDirFlagCompletion(cmd, "dir")

	return cmd
}

func runDeploy(ctx context.Context, targetRepo string, workflows []string, addOpts AddOptions, coolDown time.Duration) error {
	gitRoot, err := gitutil.FindGitRoot()
	if err != nil {
		return fmt.Errorf("--repo requires running inside a git repository: %w", err)
	}

	updatesDir, err := ensureUpdateTargetRepoGitignore(gitRoot)
	if err != nil {
		return err
	}

	checkoutDir := filepath.Join(updatesDir, sanitizeRepoPath(targetRepo))
	if err := shallowCloneTargetRepo(ctx, targetRepo, checkoutDir); err != nil {
		return err
	}

	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to read current directory: %w", err)
	}
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	if err := os.Chdir(checkoutDir); err != nil {
		return fmt.Errorf("failed to change directory to checkout %s: %w", checkoutDir, err)
	}

	if err := PreflightCheckForCreatePR(addOpts.Verbose); err != nil {
		return err
	}

	updateOpts := UpdateWorkflowsOptions{
		Verbose:                addOpts.Verbose,
		EngineOverride:         addOpts.EngineOverride,
		WorkflowsDir:           addOpts.WorkflowDir,
		NoStopAfter:            addOpts.NoStopAfter,
		StopAfter:              addOpts.StopAfter,
		DisableSecurityScanner: addOpts.DisableSecurityScanner,
		CoolDown:               coolDown,
	}
	if err := RunUpdateWorkflows(ctx, updateOpts); err != nil {
		return fmt.Errorf("failed to update existing workflows: %w", err)
	}

	if _, err := AddWorkflows(ctx, workflows, addOpts); err != nil {
		return fmt.Errorf("failed to add workflows: %w", err)
	}

	compileConfig := CompileConfig{
		Verbose:        addOpts.Verbose,
		EngineOverride: addOpts.EngineOverride,
		WorkflowDir:    addOpts.WorkflowDir,
		Purge:          true,
	}
	if _, err := CompileWorkflows(ctx, compileConfig); err != nil {
		return fmt.Errorf("failed to compile workflows with purge: %w", err)
	}

	workflowLabel := workflows[0]
	if len(workflows) > 1 {
		workflowLabel = fmt.Sprintf("%d workflows", len(workflows))
	}
	commitMessage := "chore: deploy agentic workflows"
	prTitle := "chore: deploy agentic workflows"
	prBody := fmt.Sprintf("Deploy %s to %s.\n\nThis PR was created by `gh aw deploy` after running update, add, and compile --purge in the target repository.", workflowLabel, targetRepo)

	_, err = CreatePRWithChanges("deploy-workflows", commitMessage, prTitle, prBody, addOpts.Verbose)
	if err != nil {
		return fmt.Errorf("failed to create deploy pull request: %w", err)
	}

	deployLog.Printf("Successfully deployed workflows to %s", targetRepo)
	return nil
}
