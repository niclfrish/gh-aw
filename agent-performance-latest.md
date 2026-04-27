# Agent Performance — 2026-04-27
Run: §24977099842 | Q:74→74 E:70→71

## Ecosystem Overview (Apr 27)
- Overall quality: 74/100 (→ stable), effectiveness: 71/100 (↑+1 from 70)
- 30 runs observed (past 2 days): ~22 success, 2 failure, 6 in-progress/skipped
- Effective success rate: ~92% (22/24 completed) — slight dip vs yesterday 95%

## Top Performers
1. **Test Quality Sentinel** (Q:90 E:92) — stable top performer
2. **Design Decision Gate** (Q:88 E:85) — watch period continues
3. **CLI Version Checker** (Q:82 E:84) — success ✅
4. **Auto-Triage Issues** (Q:78 E:77) — success ✅
5. **Contribution Check** (Q:78 E:77) — success ✅
6. **Issue Monster** (Q:77 E:76) — 2 runs, both success ✅
7. **Agentic Maintenance** (Q:75 E:73) — success ✅
8. **Copilot Maintenance** (Q:75 E:73) — success ✅

## New Failures (Apr 27)
- **Documentation Unbloat** (Q:45 E:40): New failure today — needs investigation

## Regressed / Still Failing 📉
- **GitHub Remote MCP Authentication Test** (Q:10 E:0) — STILL failing (#27965 P1 Day 5+). Model not supported.

## P1 Active (Apr 27)
- **GitHub Remote MCP Auth Test** (#27965): Day 5+ — model not supported (persistent)
- **Smoke Copilot** (#27028), **Smoke Claude** (#27030): Not observed today
- **Daily Community Attribution** (#28025/#28235): Not observed today
- **Documentation Unbloat**: New failure — needs investigation
- **Safe outputs session not found 37min** (#27755)
- **dependabot-go-checker compilation** (#aw_deplck)
- **awf-api-proxy sidecar** (#27888)
- **aw-gpu-runner-T4: node not found** (#aw_gpunode P0)

## 5-day Trends
- Quality: 68→72→73→74→74 (→ leveling off)
- Effectiveness: 62→68→69→70→71 (↑ slow improvement)
- Success rate: 47%→93%→94%→95%→92% (slight dip)
- P1 open: 12→13→13→13→13 (→ stable)

## Issues/Actions This Run
- Discussion created (performance report, Apr 27)
- No new improvement issues (existing issues cover open P1s)
- Documentation Unbloat failure flagged for investigation

Last updated: 2026-04-27T05:00Z by agent-performance-manager
