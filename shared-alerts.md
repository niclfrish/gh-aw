# Shared Alerts — 2026-04-19T04:41Z

## P0 (Critical)
- **Codex engine 401 auth** (#27127, OPEN): OPENAI_API_KEY missing/expired or gpt-5.3-codex model access issue. Affects AI Moderator (#27122), Daily Observability Report (#27115). Full RCA in #27127. Needs admin credential rotation.
- **Copilot CLI 11→32 versions behind** (#27143 OPEN TODAY): v1.0.21 active; v1.0.32 available. Claude Code 2.1.114 also available. Upgrade PR #27143 open.

## P1 (High)
- **Performance regressions**: BenchmarkFindIncludesInContent +51.4% (#26995), BenchmarkValidation +24% (#26993) — both assigned to Copilot; raised Apr 18
- **Smoke Claude** (#27030 OPEN): Failing since Apr 14; no sub-issues with error details
- **Smoke Copilot** (#27028 OPEN): Issue group created; investigate specific failures
- **GitHub Remote MCP Auth Test**: New failure today run #24620886472; related to #26458
- **Workflows out of sync** (#27140 OPEN TODAY): Lock files need recompile post-AWF v0.25.25 bump

## P2 (Watch)
- **PR Triage Agent** (#26778 OPEN): 67% success rate
- **Auto-Triage Issues** (#26364 OPEN): 67% success, intermittent
- **MCP Rate Limit** (#26239 OPEN): Circuit breaker needed
- **Agent Persona Explorer**: Recovered today after AWF v0.25.25. Monitor 3+ more runs.

## Recoveries (Apr 18-19)
- ✅ AWF bumped to v0.25.25 (#27101, closed)
- ✅ Agent Persona Explorer: Successful run Apr 19 (was 100% failure Apr 18)
- ✅ Many transient failure issues auto-closed (20+ in 24h)
- ✅ Lock file drift #27140 detected immediately by Agentic Maintenance

## Engine/Tool Status
- Copilot v1.0.21 active / v1.0.32 available (#27143 open)
- Claude Code 2.1.114 available (#27143 same PR)
- Codex: 401 auth failures (OPENAI_API_KEY) — avoid for new workflows
- AWF: v0.25.25 ✅; MCP Gateway: v0.2.25 ✅

## Ecosystem State
- 196 workflows total (stable)
- Schedule success rate: ~85% (↑2%)
- P0 failures: 1 (Codex 401 auth)
- Overall quality trend: ↑ recovering (Q:76 vs 74 Apr 18)

## Orchestrator Summaries
- Agent Performance (Apr 19 04:41Z): Q:76 E:75. 85% success. Codex P0. Copilot upgrade critical.
- Workflow Health (Apr 17 12:10Z): Score 73/100. 194 workflows. 0 stale lock files.

Last updated: 2026-04-19T04:41Z by agent-performance-manager
