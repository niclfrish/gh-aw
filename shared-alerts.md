# Shared Alerts — 2026-05-10T13:00Z

## P0 (Critical)
- **Smoke Gemini** (#29666 OPEN): 100% failure, proxy/API-key blocks. 35+ days. #30175 closed ineffective. Needs fresh investigation or formal deprecation.
- **Smoke CI** (#29666 OPEN): CGO/EROFS persistent, 100% action_required.
- **Daily Model Inventory Checker** (#30043 OPEN): Copilot CLI silent startup crash.
- **APM Unpack Systemic Failure** (#30252 OPEN): apm-default.tar.gz exits code 1.
- **config.models** (#30307 OPEN): unsupported AWF config field, blocks smoke runs.

## P1 (High)
- **Smoke macOS ARM64**: 100% failure since 2026-02-20 (79 days). Issue filed 2026-05-07 ✅
- **CI regression on main**: TestStrictModePermissions failing. Issue filed 2026-05-06.
- **MCP gateway session timeout** (#23153 OPEN): Long-running workflows at risk.
- **Performance Regression in Validation** (#30180): 82.1% slower.
- **Daily Fact About gh-aw**: 3+ consecutive failures — escalate to P1 if continues.

## P2 (Watch)
- **PR-review cluster** (Q, cloclo, Archie, Scout, Grumpy, Security Review, PR Nitpick, PR Code Quality): ~272 wasted run-attempts/day. Trigger gate fix or consolidation needed. HIGHEST WASTE item.
- **Plan Command over-creation**: 5 [plan] issues in one batch (#31207-#31211) — dedup gap.
- **Deployment Incident Monitor**: zombie pattern — 8x skipped today; consider deprecation.
- **Resource Summarizer Agent**: chronic skips, zero outputs.
- **Doc Build - Deploy**: action_required persistent (deployment stalled).
- **Smoke Pi**: noop violation (no safe outputs called); needs compliance fix.
- **Smoke Codex**: missing web-fetch MCP tool.
- **Node.js 20 deprecation** in CI: deadline Sep 16, 2026.

## Resolved (Do Not Re-File)
- #29863 Smoke Copilot regression → RECOVERED ✅
- #30205 Auto-Triage Issues → CLOSED ✅
- #30188 Documentation Unbloat → CLOSED ✅
- #30233 Daily Documentation Healer → CLOSED ✅
- #30069 Step Name Alignment → CLOSED ✅
- #30241 Smoke Claude → CLOSED ✅
- #30244 Smoke Codex → CLOSED ✅
- #30347, #30144 GitHub MCP Structural Analysis → CLOSED ✅
- #30085, #30086, #30087 Safe Outputs Conformance → CLOSED ✅
- #30102 Schema Consistency Checker → CLOSED ✅
- #29109 Dashboard issue (active, updated periodically)

## Trends (May 10)
- 218 workflows (stable), 0 missing lock files
- Quality: 74/100 (→ stable plateau, day 9)
- Health: 61/100 (→ stable, day 5)
- AI Moderator recovering (2/3 success)
- Under-creation dominant pattern (42% of profiled agents)
