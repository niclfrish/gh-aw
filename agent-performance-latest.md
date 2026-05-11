# Agent Performance — 2026-05-11
Run: §25673457208 | Q:74→74 E:71→71

## Ecosystem Overview (May 11)
- Overall quality: 74/100 (→ stable plateau, day 10), effectiveness: 71/100 (→ stable)
- 218 workflows (stable), health: 62/100 (↑ +1)
- Engines: copilot (140), claude (60), codex (12), pi (2), opencode/gemini/crush (1 each)
- PR-review cluster: Scout/Archie/Q/cloclo/Grumpy/Security Review/PR Nitpick/PR Code Quality — ~272 wasted run-attempts/day (0% success)
- **P0 RESOLVED today**: Smoke CI/Gemini (#29666) ✅, Daily Model Inventory (#30043) ✅, config.models (#30307) ✅
- **P0 remaining**: APM Unpack (#30252 OPEN)
- **Pattern**: Under-creation dominant (12/23 profiled agents = 52%)

## Top Performers (May 11)
1. **Agentic Maintenance** (Q:90 E:92) — Stable top performer ✅
2. **Issue Monster** (Q:85 E:87) — Active and effective ✅
3. **Auto-Close Parent Issues** (Q:82 E:85) — 100% success rate ✅
4. **Bot Detection** (Q:80 E:80) — Stable ✅
5. **PR Triage Agent** (Q:80 E:80) — Stable ✅

## Key Patterns Detected (May 11)
- `under-creation` (12 agents, 52%): PR-review cluster (8), Daily Fact, Deployment Incident Monitor, Resource Summarizer, Doc Build - Deploy
- `inconsistency` (4 agents): AI Moderator, Content Moderation, Resource Summarizer, Copilot cloud agent
- `scope-creep` (2 agents): AI Moderator, Content Moderation (recovering)
- `repetition` (1 agent): Daily Fact About gh-aw (15+ consecutive failures, no circuit-breaker)
- `over-creation` (1 agent): Plan Command (5 issues burst)

## Active Issues (May 11)
- **P0 remaining**: APM Unpack (#30252 OPEN)
- **P1**: Daily Fact parse failures (issue filed today), Smoke macOS ARM64 (filed 2026-05-07), CI TestStrictModePermissions, MCP gateway session timeout (#23153), Performance Regression (#30180)
- **P2**: PR-review cluster (highest waste ~272/day), Deployment Incident Monitor (zombie), PR #31411/#31418 pending merge

## 10-day Quality Trend
- Quality:      74→74→74→74→74→74→74→74→74→74→74 (→ stable plateau, day 10)
- Effectiveness: 71→71→71→71→71→71→71→71→71→71→71 (→ stable plateau, day 10)

## Actions This Run
- Discussion created: Agent Performance Report — Week of 2026-05-11
- Pattern analysis: pattern-detector classified 23 agents
- 3 P0 closures noted: #29666, #30043, #30307 ✅
- No new P0/P1 issues filed (P1 for Daily Fact already filed by Workflow Health)

Last updated: 2026-05-11T13:36Z by agent-performance-manager
