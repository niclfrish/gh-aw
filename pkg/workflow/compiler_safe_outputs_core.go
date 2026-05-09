// This file is the entry point for the compiler-side safe-outputs subsystem.
//
// # File-naming convention for the safe-outputs subsystem
//
// There are two parallel file groups in this package that both deal with safe
// outputs.  They are intentionally kept separate:
//
//   - compiler_safe_outputs_*.go — methods with a (*Compiler) receiver.
//     These files orchestrate the workflow compilation pipeline: they call
//     helpers, build YAML job/step/env strings, and write results into the
//     compiled workflow.  Functions here are tightly coupled to Compiler
//     state and the compilation lifecycle.
//
//   - safe_outputs_*.go — standalone (receiver-free) helper functions.
//     These files contain pure or near-pure helpers: config structs,
//     parsers, env-var builders, and validation logic that does not depend
//     on Compiler internals.  They are reusable across compilation and
//     non-compilation contexts (e.g. tests, validation-only paths).
//
// When adding new code:
//   - If the function needs Compiler state or calls other Compiler methods,
//     place it in the appropriate compiler_safe_outputs_*.go file.
//   - If the function is self-contained and testable in isolation, place it
//     in the matching safe_outputs_*.go file (env, config, jobs, steps …).
//
// # Module layout
//
// compiler_safe_outputs_core.go    — shared types (SafeOutputStepConfig)
// compiler_safe_outputs_builder.go — top-level compilation entry points
// compiler_safe_outputs_config.go  — addHandlerManagerConfigEnvVar
// compiler_safe_outputs_env.go     — addAllSafeOutputConfigEnvVars
// compiler_safe_outputs_handlers.go — per-handler compilation helpers
// compiler_safe_outputs_job.go     — buildConsolidatedSafeOutputsJob
// compiler_safe_outputs_steps.go   — buildConsolidatedSafeOutputStep

package workflow

import (
	"github.com/github/gh-aw/pkg/logger"
)

var consolidatedSafeOutputsLog = logger.New("workflow:compiler_safe_outputs_consolidated")

// SafeOutputStepConfig holds configuration for building a single safe output step
// within the consolidated safe-outputs job
type SafeOutputStepConfig struct {
	StepName                   string            // Human-readable step name (e.g., "Create Issue")
	StepID                     string            // Step ID for referencing outputs (e.g., "create_issue")
	Script                     string            // JavaScript script to execute (for inline mode)
	ScriptName                 string            // Name of the script in the registry (for file mode)
	CustomEnvVars              []string          // Environment variables specific to this step
	Condition                  ConditionNode     // Step-level condition (if clause)
	Token                      string            // GitHub token for this step
	UseCopilotRequestsToken    bool              // Whether to use Copilot requests token preference chain
	UseCopilotCodingAgentToken bool              // Whether to use Copilot coding agent token preference chain
	PreSteps                   []string          // Optional steps to run before the script step
	PostSteps                  []string          // Optional steps to run after the script step
	Outputs                    map[string]string // Outputs from this step
	ContinueOnError            bool              // Whether to continue the job even if this step fails (continue-on-error: true)
}

// Note: The implementation functions have been moved to focused module files:
// - buildConsolidatedSafeOutputsJob, buildJobLevelSafeOutputEnvVars, buildDetectionSuccessCondition
//   are in compiler_safe_outputs_job.go
// - buildConsolidatedSafeOutputStep, buildSharedPRCheckoutSteps, buildHandlerManagerStep
//   are in compiler_safe_outputs_steps.go
// - addHandlerManagerConfigEnvVar is in compiler_safe_outputs_config.go
// - addAllSafeOutputConfigEnvVars is in compiler_safe_outputs_env.go
