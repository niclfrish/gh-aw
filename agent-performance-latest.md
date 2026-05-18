# Agent Performance — 2026-05-18
Run: §26037983763 | Q:74→74 E:71→71 H:63/100 (↓4)

## Ecosystem Overview (May 18)
- Overall quality: 74/100 (plateau day 17 — breakout expected as PR-review fix + Agentic Maintenance restore take effect)
- Effectiveness: 71/100 (plateau day 17 — expect 73-76 once Agentic Maintenance restored)
- 229 workflows (stable), health: 63/100 (↓4 from Agentic Maintenance compile failure)
- Engines: copilot (140+), claude (60+), codex (12), others
- **NEW P1**: Agentic Maintenance compile failure — orchestrator DOWN 2 consecutive runs
- CGO/CJS: still failing every push to main (P1, #29669, 90+ days)
- Codex sandbox: OPENAI_API_KEY excluded (P1, #32446)

## Top Performers (May 18)
1. **Issue Monster** (Q:85 E:87) — effective, ~6m39s runtime ✅
2. **Auto-Triage Issues** (Q:82 E:85) — 100% success ✅
3. **Bot Detection** (Q:82 E:83) — 100% success, 9s runtime ✅
4. **License Compliance Check** (Q:80 E:82) — ~98% success ✅
5. **PR Sous Chef** (Q:80 E:82) — 100% success ✅
6. **Copilot SWE Agent** (Q:78 E:85) — 56% PR merge rate ✅
7. **Agentic Maintenance** (Q:was 90 → now DOWN) — compile failure P1

## Pattern Classification (May 18)
- RESOLVED: PR-review cluster (was P0, ~272 wasted runs/day) ✅
- P1 (5): Agentic Maintenance compile (NEW), CGO/CJS (#29669), Codex OPENAI_API_KEY (#32446), MCP gateway timeout (#23153), Performance Regression (#30180)
- P2 (5): UK AI Operational Resilience (NEW), ET budget exhaustion (#32717), engine-fail-after-completion (#32736), Step Name recurring (#32955), [aw-compat] warnings (#32528)
- Pattern-detector: inconsistency (11), resource-waste (4), under-creation (4), scope-creep (2), over-creation (1)
- Healthy (no patterns): Issue Monster, Auto-Triage, Bot Detection, License, PR Sous Chef, Copilot SWE, Auto-Close, PR Triage

## Active Issues (May 18)
- **Agentic Maintenance compile**: NEW P1 — orchestrator down
- **CGO/CJS**: #29669 open, failing every push — P1 (90+ days)
- **Codex OPENAI_API_KEY**: #32446 open — P1
- **MCP gateway timeout**: #23153 — P2
- **Performance Regression**: #30180 — P2
- **Daily Observability Report ET exhaustion**: #32717 — P2
- **Engine-fail-after-completion**: #32736 — P2
- **Step Name Alignment recurring**: #32955 (was #32754) — P2
- ~22 open [aw] failure issues

## 17-day Quality Trend
- Quality:       74 (plateau day 17 — expect 76-78 next run if Agentic Maintenance restored)
- Effectiveness: 71 (plateau day 17 — expect 73-75 next run)
- Primary blocker: Agentic Maintenance compile failure + CGO/CJS unresolved 90+ days

## Actions This Run
- Discussion created: Agent Performance Report — Week of 2026-05-18
- No new issues filed (existing issues cover all active items)
- Updated shared memory + shared-alerts

Last updated: 2026-05-18T14:00Z by agent-performance-manager

## Pattern Detector Results (May 18)
- inconsistency: 11 agents (top pattern fleet-wide)
- under-creation: 4 agents (CGO/CJS, Codex Smoke, Daily Fact, Weekly Blog Post)
- resource-waste: 4 agents (CGO/CJS, Daily Observability, Agentic Maintenance compile, UK AI)
- scope-creep: 2 (Codex Smoke Test, engine-fail-after-completion systemic)
- over-creation: 1 (Daily Observability Report)
- Healthy (no patterns): 8 agents
