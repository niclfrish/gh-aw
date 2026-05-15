# ADR-32239: Embed `SafeOutputTargetConfig` in Safe-Output Configs to Eliminate Target-Field Duplication

**Date**: 2026-05-15
**Status**: Draft
**Deciders**: pelikhan, Copilot

---

## Part 1 — Narrative (Human-Friendly)

### Context

`pkg/workflow` already defines a shared `SafeOutputTargetConfig` struct (`safe_outputs_parser.go:9`) that models the three cross-repository targeting fields used by every safe-output handler: `Target` (`triggering` / `*` / explicit ID), `TargetRepoSlug` (`owner/repo`), and `AllowedRepos` (additional allowed repositories). Despite this, ten safe-output config structs — `AddCommentsConfig`, `CommentMemoryConfig`, `CreateAgentSessionConfig`, `CreateCodeScanningAlertsConfig`, `CreateDiscussionsConfig`, `CreateIssuesConfig`, `CreatePullRequestReviewCommentsConfig`, `CreatePullRequestsConfig`, `PushToPullRequestBranchConfig`, and `UpdateProjectConfig` — each redeclared `TargetRepoSlug` and `AllowedRepos` (and most also redeclared `Target`) inline with identical YAML tags and identical semantics. The pattern of embedding a shared sibling config alongside `BaseSafeOutputConfig` is already established in the package (every config embeds `BaseSafeOutputConfig` with `yaml:",inline"`), so the duplication was not a constraint, just historical drift. The duplication made cross-repo targeting behavior harder to audit and caused any change to target-field semantics to require touching ten structs in lockstep.

### Decision

We will embed `SafeOutputTargetConfig` with the `yaml:",inline"` tag into each of the ten safe-output configs that carry cross-repository targeting semantics, and delete the duplicated inline `Target`, `TargetRepoSlug`, and `AllowedRepos` fields from those structs. The YAML/JSON surface is preserved exactly because the embedded fields keep the same `yaml:"target,omitempty"`, `yaml:"target-repo,omitempty"`, and `yaml:"allowed-repos,omitempty"` tags they had inline. Test fixtures and struct literals that previously initialized the duplicated fields directly are updated to initialize the embedded `SafeOutputTargetConfig` block explicitly. This is a pure refactor: no field is added or removed from any handler's YAML contract, and no compiled-workflow output changes.

### Alternatives Considered

#### Alternative 1: Leave the Duplication In Place

The ten structs could continue redeclaring the three target fields. This was rejected because the duplication actively impeded auditability — any change to target-field semantics (validation, new option, renamed YAML tag) had to be repeated across ten structs and ten test fixtures, and concurrent edits to different handlers regularly produced near-identical diffs that obscured the actual change. The shared struct already existed; not using it left readers wondering whether the inline fields had subtly different semantics from the shared form.

#### Alternative 2: Introduce a New Interface Instead of Embedding

A `TargetedSafeOutput` interface with `GetTarget()`, `GetTargetRepoSlug()`, `GetAllowedRepos()` accessor methods could have abstracted access to the three fields without changing struct layout. This was rejected because the fields are populated directly from YAML decoding and consumed by struct-literal-driven code paths (test fixtures, handler config builders), so an interface would force every call site to switch from direct field access to method calls without removing the underlying duplication. Embedding solves the duplication at the source; an interface would only paper over it.

#### Alternative 3: Migrate to a Map-Based Config Representation

The safe-output configs could have been collapsed into a single generic config with a `map[string]any` payload, eliminating the typed structs entirely. This was rejected as far out of scope: it would discard the type-safe handler signatures throughout the package, require rewriting every consumer, and was never the problem being solved. The actual problem was three fields duplicated across ten structs, which is exactly what Go struct embedding is designed for.

### Consequences

#### Positive
- `TargetRepoSlug` and `AllowedRepos` (and `Target` where applicable) are now defined in exactly one place (`SafeOutputTargetConfig` in `safe_outputs_parser.go`), so a future change to cross-repo targeting semantics touches one struct instead of ten.
- The embedded pattern visibly signals at the struct definition that a handler participates in cross-repo targeting — handlers that do not embed `SafeOutputTargetConfig` are unambiguously target-agnostic.
- Aligns target-field handling with the existing embedded-base-config style (`BaseSafeOutputConfig` is already embedded with `yaml:",inline"` across the same structs), making the safe-output config surface internally consistent.
- The YAML/JSON contract is bit-identical for every affected handler — no migration is required for workflow authors, and no compiled workflow output changes.

#### Negative
- Struct-literal call sites (primarily tests and the env-var/step config builders) must now wrap target-field initialization in a nested `SafeOutputTargetConfig{...}` block instead of setting `TargetRepoSlug:` directly, which adds two lines per fixture and is the bulk of the diff.
- Field access via reflection or by direct name in JSON marshaling tools must traverse the embedded struct; the `yaml:",inline"` tag makes this transparent for YAML, but any external tooling reading the Go struct shape sees a different layout.
- `git blame` on the three target fields now points to the embedding commit rather than the original handler-by-handler authorship, slightly increasing friction for tracing how each handler originally adopted the target fields.

#### Neutral
- No new exported API and no removed exported API: `SafeOutputTargetConfig` was already exported and `Target`/`TargetRepoSlug`/`AllowedRepos` remain accessible on each config via Go's promoted-field rules.
- The change extends the same "consolidate by concern" thread already running through [ADR-26297](26297-split-compiler-safe-outputs-config-by-concern.md) and [ADR-29230](29230-parameterize-safe-output-policy-fields-for-workflow-call.md), applied at the struct-field level rather than the file or workflow-call level.
- The `Target` field promotion means handlers that previously used `cfg.Target` directly continue to compile unchanged; only struct literals are affected by the layout change.

---

## Part 2 — Normative Specification (RFC 2119)

> The key words **MUST**, **MUST NOT**, **REQUIRED**, **SHALL**, **SHALL NOT**, **SHOULD**, **SHOULD NOT**, **RECOMMENDED**, **MAY**, and **OPTIONAL** in this section are to be interpreted as described in [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119).

### Target-Field Declaration in Safe-Output Configs (`pkg/workflow`)

1. Safe-output config structs that support cross-repository targeting **MUST** embed `SafeOutputTargetConfig` with the `yaml:",inline"` struct tag.
2. Safe-output config structs **MUST NOT** declare inline `Target`, `TargetRepoSlug`, or `AllowedRepos` fields when those fields carry cross-repository targeting semantics; the fields **MUST** be provided exclusively through the embedded `SafeOutputTargetConfig`.
3. The YAML tags on the embedded fields **MUST** remain `yaml:"target,omitempty"`, `yaml:"target-repo,omitempty"`, and `yaml:"allowed-repos,omitempty"` so that the on-disk YAML surface is preserved.
4. Safe-output config structs that genuinely do not participate in cross-repository targeting (e.g., handlers with no notion of a target repository) **MUST NOT** embed `SafeOutputTargetConfig`; embedding is reserved for handlers that semantically own the three target fields.

### Single-Source-of-Truth for Target Fields

1. `SafeOutputTargetConfig` in `pkg/workflow/safe_outputs_parser.go` **MUST** be the only declaration site for the `Target`, `TargetRepoSlug`, and `AllowedRepos` fields when those fields carry cross-repository targeting semantics across safe-output handlers.
2. New safe-output configs that need cross-repository targeting **MUST** acquire the three target fields by embedding `SafeOutputTargetConfig` and **MUST NOT** redeclare them inline.

### Struct Literal Initialization

1. Struct literals initializing a safe-output config's target fields **MUST** initialize the embedded `SafeOutputTargetConfig` block explicitly (e.g., `SafeOutputTargetConfig: SafeOutputTargetConfig{TargetRepoSlug: "org/repo"}`).
2. Struct literals **MUST NOT** rely on a top-level `TargetRepoSlug:` or `AllowedRepos:` field on the affected config types, as those fields no longer exist at that level after embedding.

### Behavior Preservation

1. The refactor **MUST NOT** alter the YAML or JSON serialization output for any affected handler.
2. The refactor **MUST NOT** alter the compiled workflow output for any input that was valid before the change.
3. Tests **MUST** be updated to use the embedded-struct initialization form and **MUST NOT** depend on the old top-level field layout.

### Conformance

An implementation is considered conformant with this ADR if it satisfies all **MUST** and **MUST NOT** requirements above. Failure to meet any **MUST** or **MUST NOT** requirement — in particular, redeclaring `Target`/`TargetRepoSlug`/`AllowedRepos` inline on a targeted safe-output config, omitting the `yaml:",inline"` tag on the embedded `SafeOutputTargetConfig`, or changing the YAML tag names of the target fields — constitutes non-conformance.

---

*This is a DRAFT ADR generated by the [Design Decision Gate](https://github.com/github/gh-aw/actions/runs/25895128163) workflow. The PR author must review, complete, and finalize this document before the PR can merge.*
