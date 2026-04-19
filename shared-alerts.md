# Shared Alerts — 2026-04-19T12:00Z

## P0 (Critical)
- **Codex engine 401 auth** (#27127, OPEN): OPENAI_API_KEY missing/expired. Duplicate Code Detector confirmed same error today (#27177). Affects AI Moderator, Daily Observability Report, Duplicate Code Detector. Needs admin credential rotation.

## P1 (High)
- **Copilot CLI upgrade critical** (#27143 OPEN): v1.0.21 active; v1.0.32 available (+11 versions). Claude Code 2.1.114 also pending. AWF must bump versions.
- **17 stale lock files** (#27140 OPEN): After AWF v0.25.25 bump. Recompile needed for: auto-triage-issues, bot-detection, code-simplifier, copilot-pr-nlp-analysis, daily-repo-chronicle, daily-safe-output-integrator, daily-security-red-team, daily-workflow-updater, gpclean, plan, prompt-clustering-analysis, refiner, scout, smoke-copilot-arm, static-analysis-report, tidy, update-astro
- **Daily Community Attribution** (#27173 TODAY): Copilot engine crash - "Permission denied" executing Python script
- **Daily Issues Report Generator** (#27165 TODAY): Copilot crash - `node: command not found`
- **Artifacts Summary** (#27155 Apr 19): Copilot stuck in Read loop
- **Smoke Claude** (#27030 OPEN): Failing since Apr 14
- **Smoke Copilot** (#27028 OPEN): Ongoing issues group

## P2 (Watch)
- **dev-hawk github-env** (#26933): High severity zizmor finding - GITHUB_ENV write from GITHUB_SERVER_URL
- **PR Triage Agent** (#26778 OPEN): 67% success rate
- **Auto-Triage Issues** (#26364 OPEN): 67% success, intermittent
- **MCP Rate Limit** (#26239 OPEN): Circuit breaker needed

## Ecosystem State
- 196 workflows total (stable)
- Schedule success rate: ~85%
- P0 failures: 1 (Codex 401 auth - ongoing from Apr 18)
- Overall quality trend: ↑ (Q:75 vs 73 Apr 17)

## Orchestrator Summaries
- Workflow Health (Apr 19 12:00Z): Score 75/100. 196 workflows. 17 stale locks. Codex P0. Copilot engine P1 failures.
- Agent Performance (Apr 19 04:41Z): Q:76 E:75. 85% success. Codex P0. Copilot upgrade critical.
- Workflow Health (Apr 17 12:10Z): Score 73/100. 194 workflows. 0 stale lock files.

Last updated: 2026-04-19T12:00Z by workflow-health-manager
