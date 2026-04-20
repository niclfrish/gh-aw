# Workflow Health - 2026-04-20T12:14Z

Score: 73/100 (→ stable from 75 Apr 19). 197 workflows. Run: §24665804498

## KEY FINDINGS

### Compilation Status
- 197/197 lock files present ✅ (stable)
- 0 stale lock files ✅ (resolved from 17 last run, #27140 fixed)
- 1 new workflow added since last run (total: 197 vs 196)

### P0 Persistent Failures
- **Codex 401 Auth** (#27127 OPEN): OPENAI_API_KEY missing/expired — **day 3**
  - Schema Feature Coverage Checker (run #24653363185 → #27286)
  - Duplicate Code Detector (run #24665358726 → #27328)
  - Affects: all Codex engine workflows

### P1 Issues Today
- **node: command not found** (#aw_node404 NEW): Recurring on aw-gpu-runner-T4
  - Daily News (#27295): run #24658469657
  - Daily Issues Report (#27301): run #24662146658
  - Also observed Apr 19 (#27165)
- **MCP Gateway startup failure** (#27317 NEW): Daily Fact About gh-aw
  - Run #24664362189 — "Start MCP Gateway" step failed
  - Uses mempalace MCP CLI server
- **GitHub App rate limit** (#27251 OPEN): Co-scheduled at 23:44 UTC
  - First observed Apr 19; may recur tonight if not fixed

### Resolved
- CLI version updates (#27143) CLOSED Apr 20 ✅
- Stale lock files (#27140) resolved ✅

### Schedule Success Rate (Today, Apr 20)
- ~85% success (5 confirmed failures from auto-generated issues)
- Healthy: most daily/weekly workflows ✅

## Open Issues (workflow-related)
- #27127 Codex 401 auth (P0) - CRITICAL — day 3, needs OPENAI_API_KEY rotation
- #aw_node404 node not found on GPU runner (P1) - NEW TODAY
- #27251 Rate limit exhaustion co-scheduled workflows (P1) - OPEN
- #27317 Daily Fact MCP Gateway failure (P1) - NEW TODAY
- #27235 Safe Outputs SEC-004 conformance (P2) - OPEN
- #27030 Smoke Claude (P1) - ongoing
- #27028 Smoke Copilot (P1) - ongoing
- #23153 MCP gateway drops in long-running jobs (bug) - OPEN

## Engine/Tool Status
- Copilot v1.0.32 ✅ (updated Apr 20)
- Claude Code 2.1.114 ✅ (updated Apr 20)
- Codex: 401 auth failures (OPENAI_API_KEY) - P0 blocked
- AWF: v0.25.25 ✅; MCP Gateway: v0.2.25 ✅

## Actions This Run
- Created #aw_node404 (node not found P1 tracker)
- Created #aw_dashboard (health dashboard)
- Memory files updated

Last updated: 2026-04-20T12:14Z by workflow-health-manager
