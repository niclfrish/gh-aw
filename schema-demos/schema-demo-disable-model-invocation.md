---
description: Demonstrates the `disable-model-invocation` schema field
on:
  workflow_dispatch:
permissions:
  contents: read
engine: codex
disable-model-invocation: false
timeout-minutes: 5
---

# Schema Demo: `disable-model-invocation`

This workflow was auto-generated to demonstrate usage of the `disable-model-invocation` field in the
gh-aw frontmatter schema. It exists solely to achieve 100% schema feature coverage.

## What `disable-model-invocation` Does

Controls whether the custom agent disables additional model invocation.

## Task

Call `noop` -- this is a coverage-only demo workflow.

**Important**: Always call the `noop` safe-output tool.

```json
{"noop": {"message": "Coverage demo for `disable-model-invocation` -- no action needed."}}
```
