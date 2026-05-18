# Workflow Health — 2026-05-18T05:52Z

Score: 63/100 (↓4 from 67). ~229 workflows. Run: §26016026133

## KEY FINDINGS

### New Issues (May 18)
- **Agentic Maintenance compile failure** (NEW P1): `compile-workflows` step failing 2 consecutive runs. All downstream jobs skipped. Issue created this run.
- **UK AI Operational Resilience** (NEW P2): OTLP header masking failing in activation job (run 26012832575)
- **Step Name Alignment** (recurred): #32955 opened — #32754 was closed May 17 but already recurred

### Auto-created issues detected today
- #32955: [aw] Step Name Alignment failed (auto-created 05:22Z)
- #32946: [aw] Weekly Blog Post Writer failed

### 🎉 Resolved Since May 17
- #32755 Sergo ✅ CLOSED
- #32754 Step Name Alignment ✅ CLOSED (but recurred same day)

### Persistent Issues (Unchanged)
- **CGO/CJS regression** (#29669): still open, 3/3 CGO runs failed this week
- **Codex OPENAI_API_KEY sandbox exclusion** (#32446): P1
- **MCP gateway session timeout** (#23153): P2
- **Performance Regression** (#30180): P2

### Systemic Patterns
- **ET budget exhaustion**: Daily Observability Report + other token-heavy workflows at risk
- **Agentic Maintenance is now DOWN**: compile step broken — orchestrator capacity impaired
- **Step Name Alignment recurring**: same day close/reopen pattern — structural fix needed

### Open [aw] failures
~21 open (↑1 from 21 due to Agentic Maintenance P1 created + Step Name #32955)

### Actions Taken This Run
- Created 1 new issue: Agentic Maintenance compile failure (P1)
- Added comment to dashboard issue #29109
- Updated shared memory

### Trends
- Score: 63/100 (↓4 — Agentic Maintenance new drag)
- CGO/CJS: still at 0% success, no fix shipped
- Step Name Alignment: recurring daily — needs root cause fix
