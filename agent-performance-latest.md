# Agent Performance — 2026-04-28
Run: §25034542635 | Q:74→74 E:71→71

## Ecosystem Overview (Apr 28)
- Overall quality: 74/100 (→ stable), effectiveness: 71/100 (→ stable)
- 30 runs observed (past 2 days): 27 success, 3 failure
- Effective success rate: ~90% (27/30) — slight dip vs yesterday 93%

## Top Performers
1. **Test Quality Sentinel** (Q:90 E:92) — 7/7 success ✅
2. **Design Decision Gate** (Q:88 E:85) — 7/7 success ✅
3. **Smoke CI** (Q:85 E:88) — 5/5 success ✅
4. **Issue Monster** (Q:77 E:76) — 2/2 success ✅
5. **Agent Persona Explorer** (Q:75 E:73) — 1/1 success ✅

## New Failures (Apr 28)
- **CLI Version Checker**: docker compose failure (`docker compose up -d --pull never` exit 1) — related to awf-api-proxy sidecar P1 #27888

## Regressed / Still Failing 📉
- **GitHub Remote MCP Authentication Test** (Q:10 E:0) — Day 7+ model not supported (#27965)
- **Documentation Unbloat** — Claude auth failure (#28659), continuing

## 6-day Trends
- Quality: 68→72→73→74→74→74 (→ stable)
- Effectiveness: 62→68→69→70→71→71 (→ leveling off)
- Success rate: 47%→93%→94%→95%→93%→90% (slight dip, docker sidecar impact)
- P1 open: 12→13→13→13→13→13 (→ stable)

## Issues/Actions This Run
- Discussion created (performance report, Apr 28)
- No new improvement issues (existing issues cover active failures)
- CLI Version Checker docker failure → linked to P1 #27888

Last updated: 2026-04-28T05:00Z by agent-performance-manager
