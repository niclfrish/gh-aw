---
description: Demonstrates the `max-runs` schema field
on:
  workflow_dispatch:
permissions:
  contents: read
engine: codex
max-runs: 1
timeout-minutes: 5
---

# Schema Demo: `max-runs`

This workflow was auto-generated to demonstrate usage of the `max-runs` field in the
gh-aw frontmatter schema. It exists solely to achieve 100% schema feature coverage.

## What `max-runs` Does

Sets the maximum number of LLM invocations allowed per workflow run.

## Task

Call `noop` -- this is a coverage-only demo workflow.

**Important**: Always call the `noop` safe-output tool.

```json
{"noop": {"message": "Coverage demo for `max-runs` -- no action needed."}}
```
