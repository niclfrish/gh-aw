# Agent Performance — 2026-05-02
Run: §25252418996 | Q:74→74 E:71→71

## Ecosystem Overview (May 2)
- Overall quality: 74/100 (→ stable day 6), effectiveness: 71/100 (→ stable)
- 58 runs analyzed: 57 completed (98.3%), 1 in-progress
- Total cost: $7.89, 23.5M tokens, 419 turns, 28 safe outputs
- Engines: copilot(18), claude(9), codex(4), crush(1), gemini(1), opencode(1), unknown(24 quick skips)

## Top Performers
1. **Test Quality Sentinel** (Q:90 E:92) — 4 runs, 0 errors, ~5-6m copilot
2. **Smoke Copilot** (Q:88 E:87) — infrastructure gatekeeper, clean run 10.9m
3. **Package Specification Enforcer** (Q:85 E:84) — Claude, 0 errors, 8.9m
4. **Daily Go Function Namer** (Q:84 E:83) — Claude, 0 errors, 5m
5. **Draft PR Cleanup** (Q:82 E:80) — 0 errors (minor: 1 missing tool security block)

## Active Failures (May 2)
- **Design Decision Gate** (P1): 50% failure rate — 2/4 runs with errors
- **AI Moderator (Codex)** (P1): Codex run failed (1/2), errors logged
- **Smoke Claude** (P0 ongoing): 1 error, 13.3m duration (timeout-adjacent)
- **Smoke Crush** (P0): 67% firewall block rate (4/6 requests blocked), 1 error
- **Smoke Gemini** (P0): 1 error (API_KEY_INVALID ongoing, #29459)
- **Smoke Codex** (P1): missing `web-fetch` tool + cache_memory miss

## Behavioral Patterns
- **Q drift** (P1): Execution drift 0→72 turns (avg 16) — unstable prompt
- **24 quick-skip PR review runs** (~1-2s): Grumpy/Nitpick/Scout/Security/cloclo all skipping without AI — normal (no relevant PR)
- **GitHub API Consumption Report**: 25.6m — longest non-smoke run (watch for timeout risk)

## 7-day Quality Trend
- Quality:      72→73→74→74→74→74→74 (→ stable)
- Effectiveness: 68→69→70→71→71→71→71 (→ stable)
- Success rate: 93%→94%→95%→57%→73%→85%→98.3%

## New Issues This Run
- No new issues (existing active: #29459 smoke engines, CI build broken from May 2 refactor)

## Recommendations
1. **Design Decision Gate**: Investigate Claude error patterns in 50% failing runs
2. **Q workflow**: Review prompt stability — 0 to 72 turn variance is excessive
3. **GitHub API Consumption Report**: Add timeout guard — 25.6m duration at risk

Last updated: 2026-05-02T13:00Z by agent-performance-manager
