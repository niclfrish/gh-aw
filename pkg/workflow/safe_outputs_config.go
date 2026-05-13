package workflow

import (
	"math"
	"strings"

	"github.com/github/gh-aw/pkg/logger"
	"github.com/github/gh-aw/pkg/typeutil"
)

var safeOutputsConfigLog = logger.New("workflow:safe_outputs_config")

// ========================================
// Safe Output Configuration Extraction
// ========================================
//
// ## Schema Generation Architecture
//
// MCP tool schemas for Safe Outputs are managed through a hybrid approach:
//
// ### Static Schemas (30+ built-in safe output types)
// Defined in: pkg/workflow/js/safe_outputs_tools.json
// - Embedded at compile time via //go:embed directive in pkg/workflow/js.go
// - Contains complete MCP tool definitions with inputSchema for all built-in types
// - Examples: create_issue, create_pull_request, add_comment, update_project, etc.
// - Accessed via GetSafeOutputsToolsJSON() function
//
// ### Dynamic Schema Generation (custom safe-jobs)
// Implemented in: pkg/workflow/safe_outputs_config_generation.go
// - generateCustomJobToolDefinition() builds MCP tool schemas from SafeJobConfig
// - Converts job input definitions to JSON Schema format
// - Supports type mapping (string, boolean, number, choice/enum)
// - Enforces required fields and additionalProperties: false
// - Custom job tools are merged with static tools at runtime
//
// ### Schema Filtering
// Implemented in: pkg/workflow/safe_outputs_config_generation.go
// - generateFilteredToolsJSON() filters tools based on enabled safe-outputs
// - Only includes tools that are configured in the workflow frontmatter
// - Reduces MCP gateway overhead by exposing only necessary tools
//
// ### Validation
// Implemented in: pkg/workflow/safe_outputs_tools_schema_test.go
// - TestSafeOutputsToolsJSONCompliesWithMCPSchema validates against MCP spec
// - TestEachToolHasRequiredMCPFields checks name, description, inputSchema
// - TestNoTopLevelOneOfAllOfAnyOf prevents unsupported schema constructs
//
// This architecture ensures schema consistency by:
// 1. Using embedded JSON for static schemas (single source of truth)
// 2. Programmatic generation for dynamic schemas (type-safe)
// 3. Automated validation in CI (regression prevention)
//

// safeOutputExtractHandler encapsulates the full parse+assign logic for one safe-output
// handler type. Using a single run function avoids the typed-nil-pointer-in-interface
// pitfall that arises when returning *T through any.
//
// Scalar settings (staged, env, github-token, max-patch-size, max-patch-files,
// runs-on, messages, steps, etc.) are handled inline in extractSafeOutputsConfig
// and do not appear in this table.
type safeOutputExtractHandler struct {
	run func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any)
}

// safeOutputExtractHandlers is the registry that drives extractSafeOutputsConfig.
// Each Shape A entry assigns a parsed config when the key is present and non-false.
// Each Shape B entry additionally falls back to a default when the key is absent.
var safeOutputExtractHandlers = []safeOutputExtractHandler{
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseIssuesConfig(m); v != nil {
			cfg.CreateIssues = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseAgentSessionConfig(m); v != nil {
			cfg.CreateAgentSessions = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseUpdateProjectConfig(m); v != nil {
			cfg.UpdateProjects = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseCreateProjectsConfig(m); v != nil {
			cfg.CreateProjects = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseCreateProjectStatusUpdateConfig(m); v != nil {
			cfg.CreateProjectStatusUpdates = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseDiscussionsConfig(m); v != nil {
			cfg.CreateDiscussions = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseCloseDiscussionsConfig(m); v != nil {
			cfg.CloseDiscussions = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseCloseIssuesConfig(m); v != nil {
			cfg.CloseIssues = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseClosePullRequestsConfig(m); v != nil {
			cfg.ClosePullRequests = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseMarkPullRequestAsReadyForReviewConfig(m); v != nil {
			cfg.MarkPullRequestAsReadyForReview = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseCommentsConfig(m); v != nil {
			cfg.AddComments = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parsePullRequestsConfig(m); v != nil {
			cfg.CreatePullRequests = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parsePullRequestReviewCommentsConfig(m); v != nil {
			cfg.CreatePullRequestReviewComments = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseSubmitPullRequestReviewConfig(m); v != nil {
			cfg.SubmitPullRequestReview = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseReplyToPullRequestReviewCommentConfig(m); v != nil {
			cfg.ReplyToPullRequestReviewComment = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseResolvePullRequestReviewThreadConfig(m); v != nil {
			cfg.ResolvePullRequestReviewThread = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseCodeScanningAlertsConfig(m); v != nil {
			cfg.CreateCodeScanningAlerts = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseAutofixCodeScanningAlertConfig(m); v != nil {
			cfg.AutofixCodeScanningAlert = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseAddLabelsConfig(m); v != nil {
			cfg.AddLabels = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseRemoveLabelsConfig(m); v != nil {
			cfg.RemoveLabels = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseAddReviewerConfig(m); v != nil {
			cfg.AddReviewer = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseAssignMilestoneConfig(m); v != nil {
			cfg.AssignMilestone = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseAssignToAgentConfig(m); v != nil {
			cfg.AssignToAgent = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseAssignToUserConfig(m); v != nil {
			cfg.AssignToUser = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseUnassignFromUserConfig(m); v != nil {
			cfg.UnassignFromUser = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseUpdateIssuesConfig(m); v != nil {
			cfg.UpdateIssues = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseUpdateDiscussionsConfig(m); v != nil {
			cfg.UpdateDiscussions = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseUpdatePullRequestsConfig(m); v != nil {
			cfg.UpdatePullRequests = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseMergePullRequestConfig(m); v != nil {
			cfg.MergePullRequest = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parsePushToPullRequestBranchConfig(m); v != nil {
			cfg.PushToPullRequestBranch = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseUploadAssetConfig(m); v != nil {
			cfg.UploadAssets = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseUploadArtifactConfig(m); v != nil {
			cfg.UploadArtifact = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseUpdateReleaseConfig(m); v != nil {
			cfg.UpdateRelease = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseLinkSubIssueConfig(m); v != nil {
			cfg.LinkSubIssue = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseHideCommentConfig(m); v != nil {
			cfg.HideComment = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseSetIssueTypeConfig(m); v != nil {
			cfg.SetIssueType = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseSetIssueFieldConfig(m); v != nil {
			cfg.SetIssueField = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseDispatchWorkflowConfig(m); v != nil {
			cfg.DispatchWorkflow = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseDispatchRepositoryConfig(m); v != nil {
			cfg.DispatchRepository = v
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseCallWorkflowConfig(m); v != nil {
			cfg.CallWorkflow = v
		}
	}},
	// Shape B: "default-on" handlers — enabled automatically when the key is absent.
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseMissingToolConfig(m); v != nil {
			cfg.MissingTool = v
		} else if _, exists := m["missing-tool"]; !exists {
			trueVal := "true"
			cfg.MissingTool = &MissingToolConfig{CreateIssue: &trueVal, TitlePrefix: "", Labels: nil}
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseMissingDataConfig(m); v != nil {
			cfg.MissingData = v
		} else if _, exists := m["missing-data"]; !exists {
			trueVal := "true"
			cfg.MissingData = &MissingDataConfig{CreateIssue: &trueVal, TitlePrefix: "", Labels: nil}
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseNoOpConfig(m); v != nil {
			cfg.NoOp = v
		} else if _, exists := m["noop"]; !exists {
			trueVal := "true"
			cfg.NoOp = &NoOpConfig{}
			cfg.NoOp.Max = defaultIntStr(1)
			cfg.NoOp.ReportAsIssue = &trueVal
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseReportIncompleteConfig(m); v != nil {
			cfg.ReportIncomplete = v
		} else if _, exists := m["report-incomplete"]; !exists {
			trueVal := "true"
			cfg.ReportIncomplete = &ReportIncompleteConfig{CreateIssue: &trueVal, TitlePrefix: "", Labels: nil}
		}
	}},
	{run: func(c *Compiler, cfg *SafeOutputsConfig, m map[string]any) {
		if v := c.parseThreatDetectionConfig(m); v != nil {
			cfg.ThreatDetection = v
		}
	}},
}

// extractSafeOutputsConfig extracts output configuration from frontmatter
func (c *Compiler) extractSafeOutputsConfig(frontmatter map[string]any) *SafeOutputsConfig {
	safeOutputsConfigLog.Print("Extracting safe-outputs configuration from frontmatter")

	var config *SafeOutputsConfig

	if output, exists := frontmatter["safe-outputs"]; exists {
		if outputMap, ok := output.(map[string]any); ok {
			safeOutputsConfigLog.Printf("Processing safe-outputs configuration with %d top-level keys", len(outputMap))
			config = &SafeOutputsConfig{}

			// Apply all handler-shaped entries (Shape A + Shape B) from the registry table.
			for _, h := range safeOutputExtractHandlers {
				h.run(c, config, outputMap)
			}

			// Parse allowed-domains configuration (additional domains, unioned with network.allowed; supports ecosystem identifiers)
			if allowedDomains, exists := outputMap["allowed-domains"]; exists {
				if domainsArray, ok := allowedDomains.([]any); ok {
					var domainStrings []string
					for _, domain := range domainsArray {
						if domainStr, ok := domain.(string); ok {
							domainStrings = append(domainStrings, domainStr)
						}
					}
					config.AllowedDomains = domainStrings
					safeOutputsConfigLog.Printf("Configured allowed-domains with %d domain(s)", len(domainStrings))
				}
			}

			// Parse allowed-github-references configuration
			if allowGitHubRefs, exists := outputMap["allowed-github-references"]; exists {
				if refsArray, ok := allowGitHubRefs.([]any); ok {
					refStrings := []string{} // Initialize as empty slice, not nil
					for _, ref := range refsArray {
						if refStr, ok := ref.(string); ok {
							refStrings = append(refStrings, refStr)
						}
					}
					config.AllowGitHubReferences = refStrings
				}
			}

			// Handle staged flag
			if staged, exists := outputMap["staged"]; exists {
				if stagedBool, ok := staged.(bool); ok {
					config.Staged = stagedBool
				}
			}

			// Handle env configuration
			if env, exists := outputMap["env"]; exists {
				if envMap, ok := env.(map[string]any); ok {
					config.Env = make(map[string]string)
					for key, value := range envMap {
						if valueStr, ok := value.(string); ok {
							config.Env[key] = valueStr
						}
					}
				}
			}

			// Handle github-token configuration
			if githubToken, exists := outputMap["github-token"]; exists {
				if githubTokenStr, ok := githubToken.(string); ok {
					config.GitHubToken = githubTokenStr
				}
			}

			// Handle max-patch-size configuration
			if maxPatchSize, exists := outputMap["max-patch-size"]; exists {
				switch v := maxPatchSize.(type) {
				case int:
					if v >= 1 {
						config.MaximumPatchSize = v
					}
				case int64:
					if v >= 1 {
						config.MaximumPatchSize = int(v)
					}
				case uint64:
					if v >= 1 {
						config.MaximumPatchSize = int(v)
					}
				case float64:
					intVal := int(v)
					// Warn if truncation occurs (value has fractional part)
					if v != float64(intVal) {
						safeOutputsConfigLog.Printf("max-patch-size: float value %.2f truncated to integer %d", v, intVal)
					}
					if intVal >= 1 {
						config.MaximumPatchSize = intVal
					}
				}
			}

			// Set default value if not specified or invalid
			if config.MaximumPatchSize == 0 {
				config.MaximumPatchSize = 1024 // Default to 1MB = 1024 KB
			}

			// Handle max-patch-files configuration (maximum unique files allowed in
			// a create-pull-request patch). Mirrors max-patch-size handling above,
			// with explicit bounds checks before narrowing to int so that very
			// large source values can't overflow/wrap into a negative or wrapped
			// number that would silently fall back to the default.
			if maxPatchFiles, exists := outputMap["max-patch-files"]; exists {
				switch v := maxPatchFiles.(type) {
				case int:
					if v >= 1 {
						config.MaximumPatchFiles = v
					}
				case int64:
					if v >= 1 {
						if v > int64(math.MaxInt) {
							safeOutputsConfigLog.Printf("max-patch-files: int64 value %d exceeds platform int range, clamping to %d", v, math.MaxInt)
							config.MaximumPatchFiles = math.MaxInt
						} else {
							config.MaximumPatchFiles = int(v)
						}
					}
				case uint64:
					if v >= 1 {
						if v > uint64(math.MaxInt) {
							safeOutputsConfigLog.Printf("max-patch-files: uint64 value %d exceeds platform int range, clamping to %d", v, math.MaxInt)
							config.MaximumPatchFiles = math.MaxInt
						} else {
							config.MaximumPatchFiles = int(v)
						}
					}
				case float64:
					// Reject NaN/Inf and clamp out-of-range floats before
					// narrowing — `int(NaN)` and `int(±Inf)` are
					// implementation-defined and can produce surprising
					// values (including 0, which would silently fall back
					// to the default).
					if v != v || v > float64(math.MaxInt) || v < float64(math.MinInt) {
						safeOutputsConfigLog.Printf("max-patch-files: float value %.2f is out of range, ignoring", v)
						break
					}
					intVal := int(v)
					if v != float64(intVal) {
						safeOutputsConfigLog.Printf("max-patch-files: float value %.2f truncated to integer %d", v, intVal)
					}
					if intVal >= 1 {
						config.MaximumPatchFiles = intVal
					}
				}
			}

			// Set default value if not specified or invalid
			if config.MaximumPatchFiles == 0 {
				config.MaximumPatchFiles = 100 // Default to 100 unique files
			}

			// Handle runs-on configuration
			if runsOn, exists := outputMap["runs-on"]; exists {
				if runsOnStr, ok := runsOn.(string); ok {
					config.RunsOn = runsOnStr
				}
			}

			// Handle messages configuration
			if messages, exists := outputMap["messages"]; exists {
				if messagesMap, ok := messages.(map[string]any); ok {
					config.Messages = parseMessagesConfig(messagesMap)
				}
			}

			// Handle activation-comments at safe-outputs top level (templatable boolean)
			if err := preprocessBoolFieldAsString(outputMap, "activation-comments", safeOutputsConfigLog); err != nil {
				safeOutputsConfigLog.Printf("activation-comments: %v", err)
			}
			if activationComments, exists := outputMap["activation-comments"]; exists {
				if activationCommentsStr, ok := activationComments.(string); ok && activationCommentsStr != "" {
					if config.Messages == nil {
						config.Messages = &SafeOutputMessagesConfig{}
					}
					config.Messages.ActivationComments = activationCommentsStr
				}
			}

			// Handle mentions configuration
			if mentions, exists := outputMap["mentions"]; exists {
				config.Mentions = parseMentionsConfig(mentions)
			}

			// Handle global footer flag
			if footer, exists := outputMap["footer"]; exists {
				if footerBool, ok := footer.(bool); ok {
					config.Footer = &footerBool
					safeOutputsConfigLog.Printf("Global footer control: %t", footerBool)
				}
			}

			// Handle group-reports flag
			if groupReports, exists := outputMap["group-reports"]; exists {
				if groupReportsBool, ok := groupReports.(bool); ok {
					config.GroupReports = groupReportsBool
					safeOutputsConfigLog.Printf("Group reports control: %t", groupReportsBool)
				}
			}

			// Handle report-failure-as-issue flag
			if reportFailureAsIssue, exists := outputMap["report-failure-as-issue"]; exists {
				if reportFailureAsIssueBool, ok := reportFailureAsIssue.(bool); ok {
					config.ReportFailureAsIssue = &reportFailureAsIssueBool
					safeOutputsConfigLog.Printf("Report failure as issue: %t", reportFailureAsIssueBool)
				}
			}

			// Handle failure-issue-repo (repository for failure issues, format: "owner/repo")
			if failureIssueRepo, exists := outputMap["failure-issue-repo"]; exists {
				if failureIssueRepoStr, ok := failureIssueRepo.(string); ok && failureIssueRepoStr != "" {
					config.FailureIssueRepo = failureIssueRepoStr
					safeOutputsConfigLog.Printf("Failure issue repo: %s", failureIssueRepoStr)
				}
			}

			// Handle max-bot-mentions (templatable integer)
			if err := preprocessIntFieldAsString(outputMap, "max-bot-mentions", safeOutputsConfigLog); err != nil {
				safeOutputsConfigLog.Printf("max-bot-mentions: %v", err)
			} else if maxBotMentions, exists := outputMap["max-bot-mentions"]; exists {
				if maxBotMentionsStr, ok := maxBotMentions.(string); ok {
					config.MaxBotMentions = &maxBotMentionsStr
				}
			}

			// Handle steps (user-provided steps injected after checkout/setup, before safe-output code)
			if steps, exists := outputMap["steps"]; exists {
				if stepsList, ok := steps.([]any); ok {
					config.Steps = stepsList
					safeOutputsConfigLog.Printf("Configured %d user-provided steps for safe-outputs", len(stepsList))
				}
			}

			// Handle id-token permission override ("write" to force-add, "none" to disable auto-detection)
			if idToken, exists := outputMap["id-token"]; exists {
				if idTokenStr, ok := idToken.(string); ok {
					if idTokenStr == "write" || idTokenStr == "none" {
						config.IDToken = &idTokenStr
						safeOutputsConfigLog.Printf("Configured id-token permission override: %s", idTokenStr)
					} else {
						safeOutputsConfigLog.Printf("Warning: unrecognized safe-outputs id-token value %q (expected \"write\" or \"none\"); ignoring", idTokenStr)
					}
				}
			}

			// Handle concurrency-group configuration
			if concurrencyGroup, exists := outputMap["concurrency-group"]; exists {
				if concurrencyGroupStr, ok := concurrencyGroup.(string); ok && concurrencyGroupStr != "" {
					config.ConcurrencyGroup = concurrencyGroupStr
					safeOutputsConfigLog.Printf("Configured concurrency-group for safe-outputs job: %s", concurrencyGroupStr)
				}
			}

			// Handle needs configuration
			if needsValue, exists := outputMap["needs"]; exists {
				if needsArray, ok := needsValue.([]any); ok {
					for _, need := range needsArray {
						if needStr, ok := need.(string); ok && needStr != "" {
							config.Needs = append(config.Needs, needStr)
						}
					}
					if len(config.Needs) > 0 {
						safeOutputsConfigLog.Printf("Configured %d explicit safe-outputs needs dependency(ies)", len(config.Needs))
					}
				}
			}

			// Handle environment configuration (override for safe-outputs job; falls back to top-level environment)
			config.Environment = c.extractTopLevelYAMLSection(outputMap, "environment")
			if config.Environment != "" {
				safeOutputsConfigLog.Printf("Configured environment override for safe-outputs job: %s", config.Environment)
			}

			// Handle jobs (safe-jobs must be under safe-outputs)
			if jobs, exists := outputMap["jobs"]; exists {
				if jobsMap, ok := jobs.(map[string]any); ok {
					c := &Compiler{} // Create a temporary compiler instance for parsing
					config.Jobs = c.parseSafeJobsConfig(jobsMap)
				}
			}

			// Handle scripts (inline handlers that run in the safe-output handler loop)
			if scripts, exists := outputMap["scripts"]; exists {
				if scriptsMap, ok := scripts.(map[string]any); ok {
					config.Scripts = parseSafeScriptsConfig(scriptsMap)
					safeOutputsConfigLog.Printf("Configured %d custom safe-output script(s)", len(config.Scripts))
				}
			}

			// Handle actions (custom GitHub Actions mounted as safe output tools)
			if actions, exists := outputMap["actions"]; exists {
				if actionsMap, ok := actions.(map[string]any); ok {
					config.Actions = parseActionsConfig(actionsMap)
					safeOutputsConfigLog.Printf("Configured %d custom safe-output action(s)", len(config.Actions))
				}
			}

			// Handle app configuration for GitHub App token minting
			if app, exists := outputMap["github-app"]; exists {
				if appMap, ok := app.(map[string]any); ok {
					config.GitHubApp = parseAppConfig(appMap)
				}
			}
		}
	}

	// Apply default threat detection whenever safe-outputs are configured and threat-detection
	// is not explicitly disabled. Detection is always on unless threat-detection is false.
	if config != nil && config.ThreatDetection == nil {
		if output, exists := frontmatter["safe-outputs"]; exists {
			if outputMap, ok := output.(map[string]any); ok {
				if _, exists := outputMap["threat-detection"]; !exists {
					// Only apply default if threat-detection key doesn't exist
					safeOutputsConfigLog.Print("Applying default threat-detection configuration")
					config.ThreatDetection = &ThreatDetectionConfig{}
				}
			}
		}
	}

	if config != nil {
		safeOutputsConfigLog.Print("Successfully extracted safe-outputs configuration")
	} else {
		safeOutputsConfigLog.Print("No safe-outputs configuration found in frontmatter")
	}

	return config
}

// parseBaseSafeOutputConfig parses common fields (max, github-token, staged) from a config map.
// If defaultMax is provided (> 0), it will be set as the default value for config.Max
// before parsing the max field from configMap. Supports both integer values and GitHub
// Actions expression strings (e.g. "${{ inputs.max }}").
func (c *Compiler) parseBaseSafeOutputConfig(configMap map[string]any, config *BaseSafeOutputConfig, defaultMax int) {
	// Set default max if provided
	if defaultMax > 0 {
		safeOutputsConfigLog.Printf("Setting default max: %d", defaultMax)
		config.Max = defaultIntStr(defaultMax)
	}

	// Parse max (this will override the default if present in configMap)
	if max, exists := configMap["max"]; exists {
		switch v := max.(type) {
		case string:
			// Accept GitHub Actions expression strings
			if strings.HasPrefix(v, "${{") && strings.HasSuffix(v, "}}") {
				safeOutputsConfigLog.Printf("Parsed max as GitHub Actions expression: %s", v)
				config.Max = &v
			}
		default:
			// Convert integer/float64/etc to string via typeutil.ParseIntValue
			if maxInt, ok := typeutil.ParseIntValue(max); ok {
				safeOutputsConfigLog.Printf("Parsed max as integer: %d", maxInt)
				s := defaultIntStr(maxInt)
				config.Max = s
			}
		}
	}

	// Parse github-token
	if githubToken, exists := configMap["github-token"]; exists {
		if githubTokenStr, ok := githubToken.(string); ok {
			safeOutputsConfigLog.Print("Parsed custom github-token from config")
			config.GitHubToken = githubTokenStr
		}
	}

	// Parse staged flag (per-handler staged mode)
	if staged, exists := configMap["staged"]; exists {
		if stagedBool, ok := staged.(bool); ok {
			safeOutputsConfigLog.Printf("Parsed staged flag: %t", stagedBool)
			config.Staged = stagedBool
		}
	}
}
