---
name: Outcome Collector
description: Periodic evaluation of safe output outcomes to measure workflow value and acceptance rates
on:
  schedule:
    - cron: every 6 hours
  workflow_dispatch:
permissions:
  contents: read
  issues: read
  pull-requests: read
  actions: read
  discussions: read
tracker-id: outcome-collector
engine:
  id: copilot
  model: claude-haiku-4.5
  bare: true
strict: true
timeout-minutes: 20
network:
  allowed:
    - defaults
    - github
tools:
  bash: true
  cache-memory: true
  github:
    mode: gh-proxy
    toolsets: [default]
safe-outputs:
  create-issue:
    title-prefix: "[Outcome Report]"
    labels: [automation, observability, outcomes]
    close-older-issues: true
    group-by-day: true
    expires: 7d
  noop:
  messages:
    footer: "> 📊 *Measured by [{workflow_name}]({run_url})*{effective_tokens_suffix}"
    run-started: "📊 [{workflow_name}]({run_url}) is evaluating safe output outcomes..."
    run-success: "📊 [{workflow_name}]({run_url}) outcome evaluation complete!"
    run-failure: "📊 [{workflow_name}]({run_url}) {status}"
imports:
  - shared/observability-otlp.md
pre-agent-steps:
  - name: Evaluate outcomes for recent runs
    env:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    run: |
      node "${RUNNER_TEMP}/gh-aw/actions/evaluate_outcomes.cjs"
  - name: Export outcome telemetry
    run: |
      if [ -f /tmp/gh-aw/outcome-evaluations.jsonl ] && [ -s /tmp/gh-aw/outcome-evaluations.jsonl ]; then
        node "${RUNNER_TEMP}/gh-aw/actions/emit_outcome_spans.cjs"
      else
        echo "No outcome evaluations to export"
      fi
---

# Outcome Collector

You are the Outcome Collector. Your job is to create a concise report of safe output outcomes.

## Input

The pre-agent step has already evaluated outcomes for recent workflow runs. Results are in:

- `/tmp/gh-aw/outcome-summary.json` — fleet-wide summary
- `/tmp/gh-aw/outcomes/run-*.json` — per-run outcome details

The summary now also includes structured breakdowns:

- `workflows[]` — per-workflow outcome counts and rates
- `types[]` — per safe-output-type outcome counts and rates
- `events[]` — per trigger/event outcome counts and rates
- `median_resolution_sec`, `median_pending_age_sec`, `zero_touch`, `zero_touch_rate`

## Task

1. Read `/tmp/gh-aw/outcome-summary.json`
2. If `total_outcomes` is 0, call `noop` with "No new safe output outcomes to report"
3. Otherwise, create a report issue with the summary

## Report Format

Create an issue with this structure:

```markdown
Safe Output Outcomes — {date}

### Fleet Summary

| Metric | Value |
|--------|-------|
| Runs checked | {runs_checked} |
| Total outcomes | {total_outcomes} |
| Accepted | {accepted} |
| Rejected | {rejected} |
| Ignored | {ignored} |
| Pending | {pending} |
| Lifecycle | {lifecycle} |
| Zero-touch accepted | {zero_touch} |
| **Acceptance rate** | **{acceptance_rate}%** |
| Waste rate | {waste_rate}% |
| Zero-touch rate | {zero_touch_rate}% |
| Median time to outcome | {median_resolution_sec} sec |

### Per-Workflow Breakdown

Use `workflows[]` from the summary. For each workflow with outcomes, show:
- Workflow name
- Outcomes: accepted / rejected / ignored / pending / lifecycle
- Acceptance rate
- Waste rate
- Zero-touch count if non-zero

If there are many workflows, keep only the highest-volume ones visible and put the rest in `<details>`.

### Outcome Type Breakdown

Use `types[]` from the summary to show which safe output types are generating the most accepted, ignored, and rejected outcomes.

### Key Observations

- Which workflows have the highest acceptance rate?
- Which workflows have the highest waste rate?
- Any workflows with all outcomes ignored (noise signal)?
- Which output types have the strongest zero-touch acceptance?
- Is median pending age increasing for any high-volume workflow or type?
```

## Guidelines

- Keep the report factual — numbers only, no speculation
- Do not re-evaluate outcomes — use the pre-computed data
- If no outcomes exist, use `noop`
- Stop immediately after creating the issue
