# Shared Alerts — 2026-04-26T12:03Z

## P0 (Critical)
- **aw-gpu-runner-T4: node not found** (NEW issue #aw_gpunode, Apr 26): 3 workflows completely failing — Daily News, Daily Fact About gh-aw, Daily Issues Report Generator. All use `copilot` engine + `aw-gpu-runner-T4`. Node.js missing from runner. Fix: remove `runs-on: aw-gpu-runner-T4` or install node on that runner. Prior issue #27534 was closed "not planned".

## P1 (High)
- **AWF binary CDN 502** (#28529 OPEN Apr 26): Intermittent HTTP 502 from GitHub releases CDN for gh-aw-firewall binary. Affects Daily Go Function Namer, Smoke CI, Design Decision Gate.
- **GitHub Remote MCP Auth Test** (#27965 OPEN): Persistent model not supported (gpt-5.1-codex-mini). Day 5+.
- **Daily Community Attribution model not supported** (#28025/#28235 OPEN): Recurring 400. Duplicate issues daily.
- **Smoke Copilot** (#27028 OPEN) + **Smoke Claude** (#27030 OPEN): Ongoing.
- **Safe outputs "session not found" at 37min** (#27755 OPEN): Long-running workflows at risk.
- **dependabot-go-checker compilation failure**: `vulnerability-alerts: read` not at job level.
- **awf-api-proxy sidecar unhealthy** (#27888 OPEN): Docker compose failures.
- **GitHub App rate limit exhaustion** (#27251 OPEN).
- **CODEX_HOME variable collision** (#27512 OPEN).

## P2 (Watch)
- **THREAT_DETECTION_RESULT parse failure**: Recurring.
- **Safe Outputs SEC-004** (#27235 OPEN).
- **Daily Documentation Updater protected files** (#27801 OPEN).
- **Performance regressions** (#27280/#27279/#27278 OPEN).
- **MCP gateway long-running drops** (#23153 OPEN).

## Trends (Apr 26)
- 203 workflows, 0 missing lock files
- P0 escalation: GPU runner node failure now explicitly tracked
- ~93% scheduled success rate (27/29 today excl. CI)

## Resolved (Recent)
- Stale lock files ✅ RESOLVED Apr 23
- CLI Version Checker Docker failure ✅ RESOLVED Apr 24
