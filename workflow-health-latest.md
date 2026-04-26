# Workflow Health — 2026-04-26T12:03Z

Score: 73/100 (↑+1 from 72 Apr 24). 203 workflows. Run: §24956132230

## KEY FINDINGS

### Compilation Status
- 203/203 lock files present ✅
- **0 missing lock files** ✅
- **0 stale lock files** ✅

### Today's Failures (Apr 26)
- **Daily Issues Report Generator**: `node: command not found` on aw-gpu-runner-T4 → auto-issue #28568
- **Daily Go Function Namer**: AWF binary CDN 502 (transient) → covered by #28529
- **CI Integration Tests**: failing (non-agentic)

### P0 Issues (NEW)
- **aw-gpu-runner-T4: node not found** (NEW issue #aw_gpunode created): 3 workflows 100% failing — Daily News, Daily Fact About gh-aw, Daily Issues Report Generator. All use copilot engine + aw-gpu-runner-T4. Prior issue #27534 closed "not planned" Apr 21 without fix. All still failing.

### P1 Issues (Active/Ongoing)
- **dependabot-go-checker compilation failure**: `vulnerability-alerts: read` not at job level
- **AWF binary CDN 502** (#28529 OPEN): Intermittent gh-aw-firewall download failures
- **Daily Community Attribution model not supported** (#28025/#28235 OPEN): Recurring 400
- **Daily Fact About gh-aw** — subsumed by P0 aw-gpu-runner-T4 issue
- **Safe outputs "session not found" at 37min** (#27755 OPEN)
- **Design Decision Gate push bundle failure** (#27756 OPEN)
- **Smoke Claude** (#27030 OPEN): Ongoing
- **Smoke Copilot** (#27028 OPEN): Ongoing
- **awf-api-proxy sidecar unhealthy** (#27888 OPEN)
- **GitHub Remote MCP Auth Test REGRESSION** (#27965 OPEN)
- **GitHub App rate limit exhaustion** (#27251 OPEN)
- **CODEX_HOME variable collision** (#27512 OPEN)

### P2 Issues (Watch)
- **THREAT_DETECTION_RESULT parse failure**: Recurring
- **Safe Outputs SEC-004** (#27235 OPEN)
- **Daily Documentation Updater protected files** (#27801 OPEN)
- **Performance regressions** (#27280/#27279/#27278 OPEN)
- **MCP gateway long-running drops** (#23153 OPEN)

## Issues Created This Run
- #aw_gpunode: [P0] aw-gpu-runner-T4: node not found — 3 workflows failing 100%

## Issues Updated
- None (no resolved issues confirmed)
