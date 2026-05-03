# Workflow Health — 2026-05-03T05:38Z

Score: 65/100 (↓3 from 68). 209 workflows. Run: §25271042058

## KEY FINDINGS

### Compilation Status
- 209/209 lock files present ✅ (+2 new workflows)
- 0 missing lock files ✅

### P0 Issues (Active)
- **Smoke Gemini** (#29459, #29852, #29816 OPEN): 100% failure rate (5/5). API_KEY_INVALID — long-standing.
- **Smoke Copilot** (#29863 OPEN): 2/5 failures (new regression ~03:37 UTC May 3)
- **Smoke Claude** (#29864 OPEN): 3/5 failures (new regression ~03:37 UTC May 3)
- **Smoke CI** (#29666 OPEN): 4/5 action_required — Crush EROFS persist
- **CGO build failures** (#29669 OPEN): ongoing failure

### P1 Issues (Active)
- **Smoke Codex**: 1/5 failures (minor regression today)
- **Smoke Crush**: 2/5 failures, 3 skipped (blocked by Gemini infrastructure issue)
- **MCP gateway session timeout** (#23153 OPEN): Ongoing structural risk

### P2 Issues
- Node.js 20 deprecation in CI (deadline Sep 16, 2026)
- #29779 YAMLGeneration 21.7% regression

### Actions Taken This Run
- Updated dashboard issue #29693 with today's status
- Updated shared memory

### Trends
- Score: 65/100 (↓3 from 68 yesterday)
- New smoke regression wave: Copilot, Claude, Codex all started failing ~03:37 UTC
- P1 backlog: still reduced vs historic highs
- 2 new workflows added (209 total)
