# ADR-31136: Consolidate Micro Codemod Files, Inline Per-Entity Update Helpers, and Generalize sortedMapKeys

**Date**: 2026-05-09
**Status**: Draft
**Deciders**: Unknown (PR author: pelikhan; auto-generated from PR diff)

---

## Part 1 — Narrative (Human-Friendly)

### Context

Over time, `pkg/cli` accumulated seven codemod files (`codemod_grep_tool.go`, `codemod_byok_copilot.go`, `codemod_inline_agents.go`, `codemod_mcp_scripts.go`, `codemod_schema_file.go`, `codemod_bots.go`, `codemod_roles.go`) where each file held one logger variable and one call to `newFieldRemovalCodemod` or `newMoveTopLevelKeyToOnBlockCodemod`. In parallel, `pkg/workflow` hosted three thin per-entity update helpers (`update_issue_helpers.go`, `update_discussion_helpers.go`, `update_pull_request_helpers.go`) of 37–45 lines each, all calling the same `parseUpdateEntityConfigTyped` generic. Two more issues compounded this fragmentation: the dual `compiler_safe_outputs_*` vs. `safe_outputs_*` naming convention was undocumented, and `sortedMapKeys` was hard-coded to `map[string]string` even though the same `MapToSlice + sort.Strings` pattern was repeated for other map value types. This PR addresses all five findings of a semantic function-clustering analysis without changing behavior.

### Decision

We will consolidate over-fragmented files when each constituent unit is purely data (a single struct literal or single delegating call) and shares its core type/infrastructure with siblings. Specifically: (1) collapse the seven micro codemod files into `codemod_field_removals.go`, (2) inline the three per-entity update parsers into the existing `update_entity_helpers.go` (mirroring the `close_entity_helpers.go` precedent), (3) document the `compiler_safe_outputs_*` vs. `safe_outputs_*` split in a file header, (4) generalize `sortedMapKeys` to `map[string]V` via Go generics and remove the now-redundant `sliceutil.MapToSlice + sort.Strings` call sites, and (5) rename single-function `util.go` to `mcp_server_helpers.go` to match its existing test file. The driver is reducing per-reader file-open count and eliminating boilerplate without sacrificing semantic clarity, consistent with prior ADRs on semantic function clustering (ADR-29952, ADR-29336, ADR-27325).

### Alternatives Considered

#### Alternative 1: Keep One File per Codemod / Per Entity (Status Quo)

Retain the existing one-file-per-thing layout under a strict interpretation of "single responsibility per file." Rejected because each split file held only data (a single struct literal) plus boilerplate (package, imports, logger var). The split added grep noise and indirection without any structural benefit — a reader investigating "which fields can be removed" had to open seven files instead of one section. The same logic applied to the three thin entity-update parsers, which all delegate to one shared generic.

#### Alternative 2: Collapse Everything into One Mega-File

Inline all codemods (including ones with non-trivial `Apply` functions) into a single file, and similarly collapse all entity helpers regardless of complexity. Rejected because codemods with custom `Apply` logic (e.g., `update_release.go`'s release-update parser, or codemods with multiple helper calls) genuinely earn their own files — they are *behavior*, not *data*. The chosen rule is narrower: consolidate only when every member is "one struct literal" or "one delegating call" and shares core types with siblings.

#### Alternative 3: Introduce a Registry/Table-Driven Approach

Replace each codemod constructor function with a single slice of `Codemod` values built declaratively. Rejected as out of scope for this PR — it would change the public API surface of how codemods are registered (currently each is a distinct exported `getXxxCodemod()` function). The current consolidation preserves all call sites and changes only file layout, deferring any deeper API redesign.

### Consequences

#### Positive
- Net −124 lines with zero behavior change and `go vet` clean.
- One file (`codemod_field_removals.go`) now answers "what field-removal codemods exist?" instead of seven.
- Generic `sortedMapKeys[V any]` removes three repetitions of the `sliceutil.MapToSlice + sort.Strings` two-line pattern and removes three now-unused `sliceutil` imports.
- The `compiler_safe_outputs_*` vs. `safe_outputs_*` split is now self-documenting via a file header in `compiler_safe_outputs_core.go`.
- `mcp_server_helpers.go` filename now matches `mcp_server_helpers_test.go`, removing a stale name mismatch (`util.go` previously held a single `boolPtr`).

#### Negative
- Larger consolidated files (`codemod_field_removals.go` is now 137 lines; `update_entity_helpers.go` grew by ~127 lines). Future contributors adding a non-trivial codemod must judge whether to add it inline or split it back out — the file header in `codemod_field_removals.go` documents this rule but it is a soft guideline, not a compiler-enforced one.
- Git blame granularity is coarsened for the moved code; `git log --follow` still works but per-codemod history requires `--all -- pkg/cli/codemod_*.go` patterns.
- Generalizing `sortedMapKeys` to `map[string]V` slightly widens its surface area; future callers can use it for any value type, which is the intent but increases the reachable signature space.

#### Neutral
- The `close_entity_helpers.go` comment had to be updated in lockstep because it referenced the now-removed `update_*_helpers.go` files as a rationale for splitting. This cross-file invariant must be maintained if either pattern is revisited.
- The boundary "consolidate when each unit is data, split when each unit is behavior" is judgment-based and not mechanically enforceable. Reviewers must apply it case by case.

---

## Part 2 — Normative Specification (RFC 2119)

> The key words **MUST**, **MUST NOT**, **REQUIRED**, **SHALL**, **SHALL NOT**, **SHOULD**, **SHOULD NOT**, **RECOMMENDED**, **MAY**, and **OPTIONAL** in this section are to be interpreted as described in [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119).

### File Consolidation Rules

1. A new codemod whose entire body is a single call to `newFieldRemovalCodemod`, `newMoveTopLevelKeyToOnBlockCodemod`, or a plain `Codemod` struct literal **MUST** be added to `pkg/cli/codemod_field_removals.go` rather than to a new dedicated file.
2. A new codemod with a non-trivial custom `Apply` function, multiple helper calls, or shared state **SHOULD** live in its own dedicated `pkg/cli/codemod_<name>.go` file.
3. New entity-update parsers that delegate to `parseUpdateEntityConfigTyped` with inline field specs and share the `UpdateEntityConfig` infrastructure **SHOULD** be added to `pkg/workflow/update_entity_helpers.go` rather than a new per-entity file.
4. An entity-update parser with a workflow that does not match the field-spec pattern (as is the case for `update_release.go`) **MAY** remain in its own file.

### Naming Convention for Safe-Outputs Subsystem

1. Files named `compiler_safe_outputs_*.go` in `pkg/workflow` **MUST** contain only methods with a `(*Compiler)` receiver that are coupled to the compilation pipeline.
2. Files named `safe_outputs_*.go` in `pkg/workflow` **MUST NOT** contain methods with a `(*Compiler)` receiver and **MUST** hold standalone, receiver-free helpers.
3. The file header in `pkg/workflow/compiler_safe_outputs_core.go` **MUST** be kept in sync with the actual layout if files are added, removed, or renamed within either group.

### Generic Sorted Map Key Helper

1. Code that needs the sorted keys of a `map[string]V` **MUST** call `sortedMapKeys(m)` from `pkg/workflow/map_helpers.go`.
2. Code **MUST NOT** reintroduce the `sliceutil.MapToSlice(m)` followed by `sort.Strings(...)` two-line pattern for `map[string]V` key extraction within `pkg/workflow`.
3. The signature of `sortedMapKeys` **MUST** remain `func sortedMapKeys[V any](m map[string]V) []string` to preserve generic applicability across value types.

### File Naming and Test Co-location

1. A `pkg/cli/<name>.go` source file whose accompanying tests live in `pkg/cli/<name>_test.go` **SHOULD** keep its base name aligned with the test file's base name; renaming one without renaming the other is **NOT RECOMMENDED**.
2. A file named `util.go` containing a single function **SHOULD** be renamed to a name describing its actual subsystem when the function is logically associated with a specific domain (e.g., MCP server helpers).

### Conformance

An implementation is considered conformant with this ADR if it satisfies all **MUST** and **MUST NOT** requirements above. Failure to meet any **MUST** or **MUST NOT** requirement constitutes non-conformance. **SHOULD**-level recommendations may be deviated from when a documented justification is provided in the file header or commit message.

---

*This is a DRAFT ADR generated by the [Design Decision Gate](https://github.com/github/gh-aw/actions/runs/25591127204) workflow. The PR author must review, complete, and finalize this document before the PR can merge.*
