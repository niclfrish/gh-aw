---
title: Agentic Optimization Kit
description: Use the built-in Agentic Optimization Kit to audit token usage, select optimization targets, and generate actionable improvement prompts from a single weekly workflow.
---

> [!WARNING]
> **Experimental:** The Agentic Optimization Kit is still experimental! Things may break, change, or be removed without deprecation at any time.

The Agentic Optimization Kit is a weekly workflow that consolidates token auditing, optimization targeting, and agentic observability into one report. It replaces three previously separate daily/weekly workflows (`copilot-token-audit`, `copilot-token-optimizer`, and `agentic-observability-kit`) with a single run that shares data, avoids duplicate log downloads, and produces a unified discussion with five embedded charts and five copyable prompt artifacts.

The kit runs every Monday. In one run it covers seven days of Copilot workflow logs, 30 days of episode history, and the full repository portfolio. It posts one discussion and opens at most one escalation issue when repeated signals warrant owner action.

## What it analyzes

The kit works in three layers, each feeding the next.

**Baseline layer (Audit):** Reads pre-downloaded 7-day Copilot run logs, groups by workflow, and computes per-workflow totals for tokens, effective tokens, cost, turns, action minutes, errors, and warnings. Flags heavy-hitters using three thresholds: workflows that own >30% of total tokens, workflows averaging >100,000 tokens per run, and workflows with more than one incident per two runs.

**Optimization layer:** Selects one high-ROI target (highest total tokens, excluding recently optimized workflows and self-referential candidates), audits its individual runs across four areas — tool usage, token efficiency, reliability, and prompt efficiency — and produces a ranked list of up to five recommendations with estimated token savings per run. The selection and cooldown log are persisted to repo-memory so the same workflow is not targeted every week.

**Observability layer:** Uses `gh aw logs` episode data to analyze the full 30-day DAG lineage, compute episode risk scores, build per-workflow instability and value-proxy scores, and classify each workflow into a portfolio quadrant (`keep`, `optimize`, `simplify`, or `review`). Applies domain-aware interpretation so that `triage` and `issue_response` workflows are judged differently from `research` or `code_fix` workflows.

## Visual report form

The discussion always includes five charts before any text walls:

### 1. Token Usage by Workflow

Horizontal bar chart of the top 15 workflows by total tokens for the 7-day window. Bars are colored by flag: dominant (>30% share), expensive per run (>100k avg), or noisy (>0.5 incidents per run). This is the fastest way to answer: which workflows dominate the token budget?

### 2. Historical Token Trend

Line chart of daily total tokens and total cost from the 90-day rolling summary stored in repo-memory. Shows week-over-week direction. Simplified or skipped if fewer than 2 data points are available.

### 3. Episode Risk–Cost Frontier

Scatter/bubble chart with episodes from `gh aw logs`. The x-axis is episode cost, the y-axis is an episode risk score that weights escalation eligibility (2.0×) and control-degradation signals (1.2–1.4×) above raw counts (1.0×), and the point size reflects run count. Pareto-frontier outliers are annotated. This is the fastest way to answer: which execution chains are both expensive and risky?

### 4. Workflow Stability Matrix

Heatmap with one row per workflow (sorted by instability score, descending) and six signal columns: risky run rate, poor control rate, resource-heavy rate, latest-success fallback rate, blocked request rate, and MCP failure rate. This is the fastest way to answer: which workflows are chronically unstable, and which are noisy only in one dimension?

### 5. Repository Portfolio Map

Scatter chart with one point per workflow. The x-axis is recent cost (or effective tokens if cost is sparse), the y-axis is a value-proxy score derived from successful usage, stability, repeat use, and absence of overkill signals. Point size reflects run count. Quadrants are labeled `keep`, `optimize`, `simplify`, and `review`. This is the fastest way to answer: which workflows deserve investment and which demand a maintainer decision?

## Actionable prompt artifacts

Each run produces five ready-to-paste prompt blocks included in the discussion. These are designed to be copied directly into an AI agent or coding assistant session:

1. **Optimization prompt** — Targets the highest-ROI workflow with the specific recommendations from the current run.
2. **Stability prompt** — Addresses the most unstable workflow identified in the stability matrix.
3. **Consolidation prompt** — Targets the highest-overlap workflow pair for consolidation review.
4. **Right-sizing prompt** — Targets a workflow flagged as overkill for its task domain.
5. **Escalation prompt** — Addresses any workflows that crossed escalation thresholds this run.

These prompts turn the weekly report into a directly actionable to-do list rather than a passive observation.

## Metric glossary

`episode_risk_score`
Composite risk score for one execution episode. Combines risky nodes (1.0×), poor-control nodes (1.2×), MCP failures (1.2×), blocked requests (1.0×), new MCP failure trends (1.4×), blocked request increase trends (1.4×), and escalation eligibility (2.0×). Escalation eligibility receives the highest weight because it already aggregates multiple threshold crossings.

`workflow_instability_score`
Workflow-level instability derived from six signal rates: risky run rate (0.25), poor control rate (0.20), resource-heavy rate (0.20), latest-success fallback rate (0.15), blocked request rate (0.10), and MCP failure rate (0.10). Weights sum to 1.0 after min-max normalization per column.

`workflow_value_proxy`
Repository-local proxy for workflow value. Combines successful recent usage (0.35), inverse instability (0.25), repeat use (0.20), and absence of overkill signals (0.20). Not a business KPI — designed to rank workflows into `keep`, `optimize`, `simplify`, and `review` rather than claim objective value.

`workflow_overlap_score`
Approximate similarity between two workflows. Blends task domain (0.30), schedule family (0.25), behavior cluster (0.20), name similarity (0.15), and assessment similarity (0.10). Values ≥ 0.55 are strong consolidation candidates.

## Shared repo-memory

The kit shares the `memory/token-audit` repo-memory branch with the `copilot-token-audit` workflow, so daily snapshots and the 90-day rolling summary are preserved across both workflows. When the two overlap in any given day, writes are merged non-destructively. The optimization cooldown log is stored separately in `optimization-log.json` on the same branch and keeps the last 30 entries.

## Relationship to source workflows

The kit consolidates three source workflows that previously ran independently:

| Source workflow | Schedule | Contribution |
|---|---|---|
| `copilot-token-audit` | Daily weekdays | Baseline per-workflow aggregates, rolling 90-day history, heavy-hitter flags |
| `copilot-token-optimizer` | Daily weekdays | Target selection, deep 4-area analysis, ranked recommendations, cooldown log |
| `agentic-observability-kit` | Weekly Monday | Episode/DAG lineage, risk/control signals, portfolio map, escalation gating |

The kit runs once per week and covers all three areas in one pass. The source workflows remain active but serve different operational frequencies. Use the kit for the weekly executive review; use the daily source workflows when you need more frequent sampling of individual signals.

## When to use it

This kit is a good fit when:

- A repository has multiple Copilot workflows and maintainers want one weekly report covering token spend, optimization targets, stability patterns, and portfolio health.
- Three separate weekly or daily workflows produce too much noise and overlap.
- The team wants five ready-to-use optimization prompts at the end of each weekly cycle.
- Escalation issues should be consolidated (one issue max) rather than filed per workflow.

This kit is a poor fit when:

- You need daily or hourly sampling — the source workflows run more frequently.
- Your repository has only one or two workflows — the portfolio analysis is not useful at that scale.
- You need exact billing reconciliation — all cost figures are run-level estimates, not invoice-grade data.

## Accuracy and cost caveats

`action_minutes` is an estimate derived from workflow duration rounded to billable minutes. Useful for relative comparison and trend detection, but not equivalent to a GitHub invoice line item.

`estimated_cost` comes from structured log fields emitted by the engine runtime. Sufficient for portfolio analysis and prioritization, but should be treated as a run-level estimate rather than finance-grade accounting.

Effective tokens are a normalization layer, not a billing unit. They account for token class weighting and model differences, making cross-run and cross-model comparisons more useful than raw token totals.

## Relationship to other tools

The kit sits above the lower-level debugging and auditing tools.

- Use [`gh aw logs`](/gh-aw/reference/audit/#gh-aw-logs---format-fmt) to inspect cross-run trends directly.
- Use [`gh aw audit`](/gh-aw/reference/audit/#gh-aw-audit-run-id-or-url) for a detailed single-run report.
- Use [Cost Management](/gh-aw/reference/cost-management/) to understand Actions minutes, inference spend, and optimization levers.
- Use [Agentic Observability Kit](/gh-aw/patterns/agentic-observability-kit/) for a separate, more detailed treatment of the observability and portfolio components.

## Source workflow

The built-in workflow lives at [`/.github/workflows/agentic-optimization-kit.md`](https://github.com/github/gh-aw/blob/main/.github/workflows/agentic-optimization-kit.md).

> [!NOTE]
> The workflow uses `tracker-id: agentic-optimization-kit`, which means the daily-audit-discussion import posts at most one active discussion per tracker ID. Previous discussions from the same tracker are expired after 7 days.
