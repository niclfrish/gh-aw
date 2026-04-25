# Agent Performance — 2026-04-25
Run: §24922665935 | Q:73↑1 E:69↑1

## Ecosystem Overview (Apr 25)
- Overall quality: 73/100 (↑+1 from 72), effectiveness: 69/100 (↑+1 from 68)
- 40 runs observed (Apr 23–25): 34 success, 1 failure, 1 cancelled, 4 in-progress
- Effective success rate: 94% (34/36 completed) — stable improvement trend day 3

## Top Performers
1. **Test Quality Sentinel** (Q:90 E:92) — 8/8 success, high-frequency sentinel
2. **Design Decision Gate** (Q:88 E:85) — 8/8 success, stable after prior P1 watch
3. **CLI Version Checker** (Q:80 E:82) — full recovery ✅
4. **AI Moderator** (Q:80 E:80) — success
5. **Auto-Triage Issues / Issue Monster / Contribution Check** (Q:78 E:77)

## Regressed / Still Failing 📉
- **GitHub Remote MCP Authentication Test** (Q:10 E:0) — STILL failing (#27965 P1). Persistent model unavailability (gpt-5.1-codex-mini not supported).

## P1 Active (Apr 25)
- **GitHub Remote MCP Auth Test** (#27965): Persistent model not supported
- **Daily Community Attribution** (#28025/#28235): Model unavailable, duplicate issues
- **Smoke Copilot** (#27028), **Smoke Claude** (#27030): Not run today (cancelled/not observed)
- **Daily Fact About gh-aw MCP Gateway** (#28245): `Start MCP Gateway` step failing
- **Safe outputs session not found 37min** (#27755)
- **dependabot-go-checker compilation** (#aw_deplck)
- **awf-api-proxy sidecar** (#27888)
- **Design Decision Gate max_turns=5** (#27470) — watch period, currently passing

## Issues/Actions This Run
- Discussion created (performance report, Apr 25)
- No new improvement issues (existing issues cover open P1s)

## 3-day Trends
- Quality: 68→72→73 (↑ steady improvement)
- Effectiveness: 62→68→69 (↑ steady improvement)  
- Success rate: 47%→93%→94% (↑ major recovery from Apr 23 bulk cancel noise)
- P1 open: 12→13→13 (→ stable)

Last updated: 2026-04-25T04:36Z by agent-performance-manager
