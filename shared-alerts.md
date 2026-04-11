# Shared Alerts — 2026-04-11T04:31Z

## P2 (High)
- **Design Decision Gate broken** (#25548, root cause in #25670): Empty prompt when --print flag used. Architecture decisions blocked. Fix documented, awaiting PR.
- **Documentation Unbloat zero-output** (ongoing): Claude workflow ~$55/week, 0 safe outputs. Agent runs successfully but never calls GitHub output tools. Investigation needed.
- **Smoke Copilot recovering**: v1.0.24 PR #25752 in progress. Was 21/30 failures during v1.0.21/1.0.22 regression. Should auto-resolve post-merge.
- **Smoke Gemini failing**: 100% failure today (1/1). Gemini CLI updated to 0.37.0 in recent CA run.
- **Smoke Cross-Repo Create/Update failing**: 100% failure today. Likely Copilot version-related.

## P3 (Watch)
- **Contribution Check report_incomplete**: Every run. Permission/network issue. #25215-era problem.
- **GitHub Remote MCP Auth Test**: 100% failure today. #24829 closed not_planned but test still failing.
- **Workflow Normalizer deduplication gap**: Created 3 identical "Normalize report formatting" issues (#25412, #25554, #25724) in 24h.

## Recent Fixes
- Copilot v1.0.21 crash RESOLVED by pinning to v1.0.20 (Apr 8-10 saga)
- v1.0.22 regression (05:26-11:50 UTC Apr 10): 5 workflow failures, now self-healing
- v1.0.24 bump in progress (#25751/PR #25752) - will restore Copilot smoke tests

## Ecosystem State
- 187 compiled workflows. Health: ~75/100 (recovering). 20/25 scheduled healthy.
- Engine split: ~124 copilot, ~41 claude, ~18 codex, ~4 others
- v1.0.20 currently pinned as stable Copilot version (v1.0.24 bump in progress)
- Claude/Codex engines: 100% resilience

## Orchestrator Summaries (Apr 11)
- Agent Performance: Q:70↑5, E:60↓6. Recovery trend continuing. See discussion.
- Workflow Health (Apr 10): Score 75/100 ↑5. v1.0.20 stable, failures self-healing.

Last updated: 2026-04-11T04:31Z by agent-performance-analyzer
