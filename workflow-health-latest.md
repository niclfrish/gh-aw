# Workflow Health - 2026-04-19T12:00Z

Score: 75/100 (↑2 from 73 Apr 17). 196 workflows. Run: §24628574565

## KEY FINDINGS

### Compilation Status
- 196/196 lock files present ✅ (stable)
- 17 stale lock files ⚠️ (AWF v0.25.24→v0.25.25 bump, tracked in #27140)
  Stale: auto-triage-issues, bot-detection, code-simplifier, copilot-pr-nlp-analysis, daily-repo-chronicle, daily-safe-output-integrator, daily-security-red-team, daily-workflow-updater, gpclean, plan, prompt-clustering-analysis, refiner, scout, smoke-copilot-arm, static-analysis-report, tidy, update-astro

### P0 Persistent Failures
- **Codex 401 Auth** (#27127 OPEN): OPENAI_API_KEY missing/expired
  - Duplicate Code Detector failed today (run #24628420682) - same 401 auth error confirmed
  - Affects: all Codex engine workflows (duplicate-code-detector, ai-moderator, daily-observability-report)

### P1 Issues Today
- **Daily Community Attribution** (#27173 TODAY): Copilot engine crash: "Permission denied" running Python script
- **Daily Issues Report Generator** (#27165 TODAY): Copilot crash - `node: command not found`
- **Artifacts Summary** (#27155 Apr 19 06:27): Copilot stuck in Read loop
- **Copilot CLI upgrade** (#27143 OPEN): v1.0.21 active; v1.0.32 available (+11 versions)
  - Claude Code 2.1.112→2.1.114 also pending

### Schedule Success Rate (Today, Apr 19)
- ~85% success (6 failures / ~41 observed schedule runs)
- Healthy: Issue Monster (10/11), most daily workflows ✅
- New failures: Daily Community Attribution, Daily Issues Report, Duplicate Code Detector, Artifacts Summary

## Open Issues (workflow-related)
- #27127 Codex 401 auth (P0) - CRITICAL
- #27143 CLI version updates: Copilot 1.0.21→1.0.32, Claude Code 2.1.114 (P1) - UPGRADE NEEDED
- #27140 17 stale lock files (P1) - AWF version bump
- #27177 Duplicate Code Detector failed (auto-generated, TODAY)
- #27173 Daily Community Attribution failed (auto-generated, TODAY)
- #27165 Daily Issues Report Generator failed (auto-generated, TODAY)
- #27155 Artifacts Summary failed (auto-generated, Apr 19 early)
- #27030 Smoke Claude (P1) - ongoing
- #27028 Smoke Copilot (P1) - ongoing
- #27128 Failure Investigator group (tracking)
- #26933 Static Analysis Report - dev-hawk github-env High severity

## Actions This Run
- No new P0/P1 issues created (all auto-tracked)
- Dashboard issue created (see GitHub)
- Memory files updated

## Engine/Tool Status
- Copilot v1.0.21 active / v1.0.32 available (#27143 open)
- Claude Code 2.1.114 available (#27143 same)
- Codex: 401 auth failures (OPENAI_API_KEY) - P0 blocked
- AWF: v0.25.25 ✅; MCP Gateway: v0.2.25 ✅
