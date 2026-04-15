# Workflow Health - 2026-04-15T12:09Z

Score: 72/100 (↓1 from 73 Apr 14). 191 workflows. Run: §24453586573

## KEY FINDINGS

### Stale Lock Files (NEW regression)
- 18/191 stale lock files (md newer than lock) — needs `make recompile`
- Likely caused by PR #26372 merged Apr 15 06:04Z (context propagation changes)
- List: agent-performance-analyzer, artifacts-summary, changeset, code-scanning-fixer,
  copilot-pr-merged-report, daily-fact, daily-performance-summary, daily-security-red-team,
  duplicate-code-detector, github-mcp-structural-analysis, lockfile-stats, notion-issue-summary,
  refactoring-cadence, slide-deck-maintainer, smoke-temporary-id, spec-enforcer,
  step-name-alignment, workflow-health-manager

### New Failures Today (Apr 15)
- Daily Issues Report Generator: `node: command not found` (#26393) - RECURRING
- Auto-Triage Issues: 2/2 failed today (#26364, existing issue)
- Daily Fact About gh-aw: failed, no error msg (#26405)
- Daily News: failed (#26388)
- Smoke Gemini: ongoing (#26351) - Gemini 0.38.0 available may fix after CLI bump

### Ongoing P2 Issues (from previous runs)
- Smoke Claude: fails on SCHEDULE (#25727)
- Smoke Multi PR: persistent (#25415)
- Smoke Cross-Repo PR Create: stale 7 days (#25221)
- Smoke Cross-Repo PR Update: stale 7 days (#25217)
- Daily Firewall Logs: safe_outputs process failure (#25456)
- Schema Feature Coverage Checker: protected-files config blocks PR (#25992)
- ~16 other open workflow failures (Go Logger, Multi-Device Docs, Refactoring Cadence, etc.)
- GitHub Remote MCP Auth Test: #24829 closed not_planned (still failing)

### Healthy Workflows (Today's Successes)
- Contribution Check: successful report #26376 (4 PRs reviewed)
- Architecture Diagram: #26389 (successful incremental update)
- Terminal Stylist: successful schedule run
- No-Op Tracker: 364+ comments active (#25214)
- 40+ other workflows running successfully

### Performance
- CompileComplexWorkflow regression: +18.9% slower (#26378)
- MCP Rate Limit event: #26239 (circuit breaker needed)

## Compilation Status
- 191/191 lock files present ✅
- 18 stale lock files ⚠️ (new this run)

## Engine/Tool Status
- Copilot v1.0.27 available (was v1.0.21 active) → #26367
- Claude Code 2.1.109 available
- Codex 0.120.0 available
- Gemini 0.38.0 available (may fix Smoke Gemini)
- GitHub MCP v0.33.1 available

## Score Breakdown
- Compilation 191/191 ✅: +35
- 18 stale lock files: -3
- Multiple healthy workflows today: +18
- ~27 open failure issues: -14
- Smoke suite partial failures continuing: -4
- Net: ~72/100

## Score Trend
68 → 71 → 73 → 71 → 70 → 75 → 73 → 74 → 74 → 73 → 72
Apr5  Apr6  Apr7  Apr8  Apr9  Apr10 Apr11 Apr12 Apr13 Apr14 Apr15

## Dashboard Issue
Created #aw_dash15 (see safeoutputs)

## Note: GitHub API Read-Only
gh CLI not authenticated; used GitHub MCP for all reads.
