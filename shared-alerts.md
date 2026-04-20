# Shared Alerts — 2026-04-20T12:14Z

## P0 (Critical)
- **Codex engine 401 auth** (#27127, OPEN, day 3): OPENAI_API_KEY missing/expired. Affects all Codex workflows: AI Moderator, Daily Observability Report, Duplicate Code Detector, Schema Feature Coverage. Needs admin credential rotation. New failures today: #27328, #27286.

## P1 (High)
- **node: command not found on aw-gpu-runner-T4** (#aw_node404 NEW): Recurring across Daily News and Daily Issues Report 2+ days. Node.js PATH not available in bash execution context on GPU runners. Created Apr 20.
- **GitHub App rate limit exhaustion** (#27251 OPEN): Co-scheduled workflows at 23:44 UTC exhausting installation rate limit. May recur tonight. Stagger schedules recommended.
- **Daily Fact MCP Gateway startup failure** (#27317 NEW): "Start MCP Gateway" step failed Apr 20. Uses mempalace MCP CLI server. May be transient.
- **Smoke Claude** (#27030 OPEN): Failing since Apr 14 — ongoing
- **Smoke Copilot** (#27028 OPEN): Ongoing issues group

## P2 (Watch)
- **Safe Outputs SEC-004** (#27235 OPEN): 4 handler files need sanitization
- **Performance regressions** (#27280/#27279/#27278 OPEN): CompileComplexWorkflow +29%, CompileSimpleWorkflow +39%, Validation +96%
- **dev-hawk github-env** (#26933): High severity zizmor finding
- **PR Triage Agent** (#26778 OPEN): 67% success rate
- **MCP gateway long-running drops** (#23153 OPEN): Session not found after 30-45min

## Resolved (Recent)
- CLI updates (#27143) CLOSED Apr 20 ✅
- Stale lock files (#27140) resolved ✅

## Ecosystem State
- 197 workflows total (+1 from last check)
- 0 stale lock files (↓ from 17)
- Schedule success rate: ~85%
- P0 failures: 1 (Codex 401 auth — day 3)
- Overall quality trend: → stable (Q:73)

## Orchestrator Summaries
- Workflow Health (Apr 20 12:14Z): Score 73/100. 197 workflows. 0 stale locks. Codex P0 (day 3). New P1: node not found on GPU runner. Rate limit exhaustion ongoing.
- Workflow Health (Apr 19 12:00Z): Score 75/100. 196 workflows. 17 stale locks. Codex P0. Copilot engine P1 failures.
- Agent Performance (Apr 20 04:46Z): Q:73 E:70. 18 workflows, 33 runs. Codex P0 ongoing.

Last updated: 2026-04-20T12:14Z by workflow-health-manager
