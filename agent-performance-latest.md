# Agent Performance — 2026-05-10
Run: §25629331070 | Q:74→74 E:71→71

## Ecosystem Overview (May 10)
- Overall quality: 74/100 (→ stable plateau, day 9), effectiveness: 71/100 (→ stable)
- 218 workflows (stable), health: 61/100 (→ stable, day 5)
- Engines: copilot (140), claude (60), codex (12), pi (2), opencode/gemini/crush (1 each)
- PR-review cluster: Scout/Archie/Q/cloclo/Grumpy/Security Review/PR Nitpick/PR Code Quality — ~272 wasted run-attempts/day (0% success)
- **P0 ongoing**: Smoke Gemini 100% failure (35+ days), fetch TypeError
- **Pattern**: Under-creation dominant (8/19 profiled agents = 42%)

## Top Performers (May 10)
1. **Agentic Maintenance** (Q:90 E:92) — Stable top performer ✅
2. **Issue Monster** (Q:85 E:87) — Active and effective ✅
3. **Auto-Close Parent Issues** (Q:82 E:85) — 100% success rate ✅
4. **Bot Detection** (Q:80 E:80) — Stable ✅
5. **PR Triage Agent** (Q:80 E:80) — Stable ✅

## Key Patterns Detected (May 10)
- `under-creation` (8 agents, 42%): PR-review cluster (8 sub-agents), Smoke Gemini, Smoke Pi, Smoke Codex, Resource Summarizer, Doc Build Deploy, Deployment Incident Monitor, Daily Fact
- `inconsistency` (7 agents): PR-review cluster, AI Moderator, Content Moderation, Daily Fact, Dev, Stale PR Cleanup, Weekly Editors Health Check
- `scope-creep` (improving): AI Moderator, Content Moderation (2/3 success, recovering)
- `over-creation`+`repetition`: Plan Command — 5 issues in <60s (#31207-#31211)

## Active Issues (May 10, unchanged from May 9)
- **P0 ongoing**: Smoke Gemini (35+ days) — #30175 fix ineffective
- **P0 ongoing**: Smoke CI CGO/EROFS — #29666
- **P0 ongoing**: APM unpack systemic — #30252
- **P0 ongoing**: Daily Model Inventory Checker — #30043
- **P0 ongoing**: config.models unsupported field — #30307
- **P1 ongoing**: Smoke macOS ARM64 — filed 2026-05-07
- **P1 ongoing**: CI TestStrictModePermissions
- **P1 ongoing**: MCP gateway session timeout — #23153
- **P1 ongoing**: Performance Regression — #30180

## 7-day Quality Trend
- Quality:      74→74→74→74→74→74→74→74→74 (→ stable plateau, day 9)
- Effectiveness: 71→71→71→71→71→71→71→71→71 (→ stable plateau, day 9)

## Actions This Run
- Discussion created: Agent Performance Report — Week of 2026-05-10
- Pattern analysis: pattern-detector classified 19 agents
- PR-review cluster identified as highest-priority (272 wasted runs/day)
- No new P0/P1 issues filed (all active items already tracked)

Last updated: 2026-05-10T13:00Z by agent-performance-manager
