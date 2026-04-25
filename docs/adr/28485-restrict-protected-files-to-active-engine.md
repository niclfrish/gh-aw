# ADR-28485: Restrict Protected Agent Config Files Collection to the Active Engine

**Date**: 2026-04-25
**Status**: Draft
**Deciders**: pelikhan

---

## Part 1 — Narrative (Human-Friendly)

### Context

Each compiled agentic workflow contains save/restore steps that snapshot and re-apply agent config directories (e.g., `.claude/`, `.codex/`, `.gemini/`, `.crush/`, `.opencode/`) from the base branch. This mechanism prevents fork PRs from injecting malicious instruction files into the workflow's execution environment. Previously, these steps collected files and directories for **all** registered engines regardless of which engine the workflow actually used, resulting in unnecessary sparse-checkout paths and extraneous environment variables being set in every compiled workflow.

### Decision

We will restrict the collection of protected agent config files and folders to the **active engine** for the workflow being compiled. Two new `EngineRegistry` methods — `GetEngineAgentManifestFolders(engineID)` and `GetEngineAgentManifestFiles(engineID)` — replace the previous calls to `GetAllAgentManifestFolders()` and `GetAllAgentManifestFiles()` at every call site in the compiler. The engine ID is resolved from `WorkflowData.EngineConfig` and passed through the compilation pipeline. The platform-owned `.agents` directory is always included regardless of engine.

### Alternatives Considered

#### Alternative 1: Keep Collecting All Engine Files (Status Quo)

The existing `GetAllAgentManifestFolders()` / `GetAllAgentManifestFiles()` approach collects from every registered engine unconditionally. This is the simplest code path — no engine ID needs to be threaded through the compiler — but it causes every workflow to protect directories that are irrelevant to its engine. For example, a Claude workflow would still check out and restore `.gemini/`, `.codex/`, and `.opencode/`. This wastes I/O during sparse checkout and pollutes the environment with unnecessary `GH_AW_AGENT_FOLDERS`/`GH_AW_AGENT_FILES` values. It was not chosen because reducing scope is both a performance improvement and a security improvement.

#### Alternative 2: Static Engine-to-Files Configuration Map

A static map (e.g., a `map[string][]string` keyed by engine ID) could provide the per-engine file/folder lists without introducing new methods on `EngineRegistry`. This would separate configuration data from the registry object. However, it would duplicate information already encoded in each engine's `AgentFileProvider` implementation and would need to be kept in sync manually. The new instance methods on `EngineRegistry` are preferred because they derive the data from the existing `AgentFileProvider` interface, keeping the source of truth in one place.

### Consequences

#### Positive
- Compiled workflows produce narrower sparse-checkout paths, reducing the number of files fetched from the base branch during activation.
- The protection boundary is correctly scoped: only the configuration directories that could actually influence the active engine are snapshotted and restored, reducing the blast radius of any future misconfiguration.
- The `GH_AW_AGENT_FOLDERS` and `GH_AW_AGENT_FILES` environment variables in generated YAML reflect only the active engine, making the compiled output easier to audit.

#### Negative
- The engine ID must be threaded through three compiler entry points (`generateCheckoutGitHubFolderForActivation`, `addActivationRepositoryAndOutputSteps`, `generateMainJobSteps`), adding boilerplate `engineID := ""` guard clauses that assume an empty string is a safe fallback.
- When `data.EngineConfig` is `nil` or missing, the fallback `engineID = ""` silently produces a minimal result (only `.agents`). This edge case is not validated; a missing engine config will not produce a compiler error but will emit a workflow that protects fewer files than expected.
- All 202 workflow lock files were recompiled, producing a large diff that obscures the actual logic change and complicates code review.

#### Neutral
- The `.agents` platform directory is unconditionally included by `GetEngineAgentManifestFolders` regardless of engine, preserving the prior behavior for that directory.
- Unit tests were added for the two new `EngineRegistry` methods; golden test files were updated to reflect the new engine-scoped output.

---

## Part 2 — Normative Specification (RFC 2119)

> The key words **MUST**, **MUST NOT**, **REQUIRED**, **SHALL**, **SHALL NOT**, **SHOULD**, **SHOULD NOT**, **RECOMMENDED**, **MAY**, and **OPTIONAL** in this section are to be interpreted as described in [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119).

### Engine-Scoped File Collection

1. Implementations **MUST** resolve the active engine ID from `WorkflowData.EngineConfig.ID` before calling any manifest-folder or manifest-file collection method.
2. Implementations **MUST** use `EngineRegistry.GetEngineAgentManifestFolders(engineID)` and `EngineRegistry.GetEngineAgentManifestFiles(engineID)` in the save, restore, and sparse-checkout steps instead of `GetAllAgentManifestFolders()` / `GetAllAgentManifestFiles()`.
3. Implementations **MUST NOT** include manifest folders or files belonging to engines other than the active engine in the `GH_AW_AGENT_FOLDERS` or `GH_AW_AGENT_FILES` environment variables of compiled workflows.
4. The `.agents` platform directory **MUST** always be included in the result of `GetEngineAgentManifestFolders`, regardless of the engine ID passed.

### EngineRegistry API

1. `GetEngineAgentManifestFolders` **MUST** accept an engine ID string and return only the manifest folder prefixes registered for that engine, plus `.agents`.
2. `GetEngineAgentManifestFiles` **MUST** accept an engine ID string and return only the manifest files registered for that engine.
3. Both methods **MUST** return a sorted, deduplicated list.
4. When the engine ID is empty or refers to an unknown engine, `GetEngineAgentManifestFolders` **MUST** return at minimum `[".agents"]`; `GetEngineAgentManifestFiles` **MAY** return `nil` or an empty slice.
5. Implementations **SHOULD NOT** pass an empty engine ID to these methods in production call sites; callers **SHOULD** log or return an error if `EngineConfig` is unexpectedly nil.

### Conformance

An implementation is considered conformant with this ADR if it satisfies all **MUST** and **MUST NOT** requirements above. Failure to meet any **MUST** or **MUST NOT** requirement constitutes non-conformance.

---

*ADR created by [adr-writer agent]. Review and finalize before changing status from Draft to Accepted.*
