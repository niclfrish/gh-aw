# Agent Performance — 2026-04-29
Run: §25091475936 | Q:74→74 E:71→71

## Ecosystem Overview (Apr 29)
- Overall quality: 74/100 (→ stable), effectiveness: 71/100 (→ stable)
- 26 completed runs observed today: 22 success, 4 failure
- Effective success rate: ~85% (22/26) — slight dip vs yesterday 90%

## Top Performers
1. **Test Quality Sentinel** (Q:90 E:92) — 5/5 success ✅
2. **Design Decision Gate** (Q:88 E:85) — 4/4 success ✅
3. **Smoke OpenCode/Claude/Codex/Copilot** (Q:85 E:85) — all passing ✅
4. **Issue Monster** (Q:77 E:76) — 1/1 success ✅
5. **Agent Persona Explorer** (Q:75 E:73) — 1/1 success ✅

## New Failures (Apr 29)
- **Smoke Crush**: failure (new today)
- **Smoke Gemini**: failure (new today)

## Regressed / Still Failing 📉
- **GitHub Remote MCP Authentication Test** (Q:10 E:0) — Day 8+ model not supported (#27965)
- **Documentation Unbloat** — Claude auth failure (#28659), continuing

## 7-day Trends
- Quality: 72→73→74→74→74→74→74 (→ stable)
- Effectiveness: 68→69→70→71→71→71→71 (→ stable)
- Success rate: 93%→94%→95%→93%→90%→85% (slight dip, new Smoke failures)
- P1 open: 13→13→13→13→13→13→13 (→ stable, backlog not shrinking)

## Issues/Actions This Run
- Discussion created (performance report, Apr 29)
- No new improvement issues (existing issues cover active failures)
- Smoke Crush and Smoke Gemini failures noted — likely engine-side transient issues

Last updated: 2026-04-29T05:00Z by agent-performance-manager
