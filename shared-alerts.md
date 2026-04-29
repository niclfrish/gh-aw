# Shared Alerts — 2026-04-29T12:24Z

## P0 (Critical)
- **Daily Fact About gh-aw codex crash** (auto-issue #29088 created today): `codex: command not found`. Recurring daily.

## P1 (High)
- **CI integration tests failing** (NEW TODAY): js-integration-live-api `ERR_API: fetch file .github/workflows/audit-workflows.md failed`; 4 jobs failed total.
- **Node.js 20 deprecation** in CI: actions/setup-go on deprecated Node.js 20 (removal Sep 16, 2026).
- **Documentation Unbloat claude auth failure** (#28659 OPEN): Claude OAuth token issue. Recurring.
- **GitHub Remote MCP Authentication Test** (#27965 OPEN): Day 8+ of model-not-supported error.
- **Safe outputs session not found** (#23153 OPEN): Long-running workflows at risk.
- **awf-api-proxy sidecar unhealthy** (#27888 OPEN): Docker compose failures.
- **GitHub App rate limit exhaustion** (#27251 OPEN).
- **CODEX_HOME variable collision** (#27512 OPEN).

## P2 (Watch)
- **Safe Outputs SEC-004** (#27235 OPEN).
- **Daily Documentation Updater protected files** (#27801 OPEN).
- **Performance regressions** (#27280/#27279/#27278 OPEN).
- **MCP gateway long-running drops** (#23153 OPEN).

## Trends (Apr 29)
- 204 workflows, 0 missing lock files
- Scheduled success rate: 73% (22/30) — IMPROVEMENT from 57% yesterday
- THREAT_DETECTION_RESULT systemic failure from Apr 28 did NOT recur today
- 6 safe_outputs job failures (agent crashes without OTEL - pattern to watch)
- Daily Fact codex crash: persistent P0

## Resolved (Recent)
- **THREAT_DETECTION_RESULT parse failure** (#28866): Was systemic Apr 28 (3+ workflows), appears resolved/intermittent Apr 29 — watch.
