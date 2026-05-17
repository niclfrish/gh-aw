# Agent Performance — 2026-05-17
Run: §25991507629 | Q:74→74 E:71→71 H:67/100

## Ecosystem Overview (May 17)
- Overall quality: 74/100 (plateau day 16 — breaking: PR-review cluster fixed ✅)
- Effectiveness: 71/100 (expect uplift next run)
- 229 workflows (stable), health: 67/100 (↑3 from PR cluster fix)
- Engines: copilot (140+), claude (60+), codex (12), others
- **PR-review cluster #31724 RESOLVED**: ~272 wasted runs/day eliminated ✅
- CGO/CJS: still failing every push to main (P1, #29669)

## Top Performers (May 17)
1. **Agentic Maintenance** (Q:90 E:92) — 100% success ✅
2. **Issue Monster** (Q:85 E:87) — effective, ~6m39s runtime ✅
3. **Auto-Triage Issues** (Q:82 E:85) — 100% success ✅
4. **Bot Detection** (Q:82 E:83) — 100% success, 9s runtime ✅
5. **License Compliance Check** (Q:80 E:82) — ~98% success ✅
6. **PR Sous Chef** (Q:80 E:82) — 100% success ✅
7. **Copilot SWE Agent** (Q:78 E:85) — 56% PR merge rate ✅

## Pattern Classification (May 17)
- RESOLVED: PR-review cluster (was P0, 0% success, 272 wasted runs/day) ✅
- P1 (2): CGO/CJS regression (#29669), Codex OPENAI_API_KEY (#32446), Smoke CI (#32690)
- P2 (3): ET budget exhaustion (#32717), engine-fail-after-completion (#32736), Daily Fact (#31432/#31524)
- NEW: Clustered transient failures May 17 — Sergo (#32755), Step Name (#32754), Linter Miner (#32748), Outcome Collector (#32728)
- OK: Agentic Maintenance, Issue Monster, Auto-Close, Bot Detection, License, PR Triage, Auto-Triage, PR Sous Chef

## Active Issues (May 17)
- **CGO/CJS**: #29669 open, failing every push — P1
- **Smoke CI**: #32690 open — P1
- **Codex OPENAI_API_KEY**: #32446 open — P1
- **Daily Observability Report ET exhaustion**: #32717 — P2
- **Engine-fail-after-completion**: #32736 — P2
- **Daily Fact**: #31432, #31524 open — P2
- **MCP gateway timeout**: #23153 open — P2
- **Performance Regression**: #30180 open — P2
- 21 open [aw] failure issues (↑2 from May 16; PR cluster fix reduces future count)

## 16-day Quality Trend
- Quality:       74 (plateau day 16 — breaking point: cluster fix removes drag)
- Effectiveness: 71 (plateau day 16 — expect 76-78 / 73-75 next run)
- Primary driver resolved: PR-review cluster waste (~272 runs/day at 0%) ✅

## Actions This Run
- Discussion created: Agent Performance Report — Week of 2026-05-17
- No new issues filed (existing issues cover all active items)
- Updated shared memory + shared-alerts

Last updated: 2026-05-17T13:00Z by agent-performance-manager

## Pattern Detector Results (May 17)
- inconsistency: 11 agents
- resource-waste: 4 agents (PR-review cluster RESOLVED, CGO/CJS, Observability Report, Deployment Monitor)
- zombie: 2 (PR-review cluster RESOLVED, Deployment Monitor active)
- over-creation: 2 (same as zombie)
- scope-creep: 1 (Codex Smoke Test — OPENAI_API_KEY sandbox access)
- Healthy (no patterns): 7 agents
