# Shared Alerts — 2026-05-18T14:00Z

## P1 (High)
- **Agentic Maintenance compile failure** (NEW): compile-workflows down 2 runs — orchestrator impaired, all downstream jobs skipped
- **CGO/CJS regression** (#29669): failing every push to main (90+ days, 3/3 runs this week) — CRITICAL: 90 days unresolved
- **Codex OPENAI_API_KEY sandbox exclusion** (#32446): blocking all Codex workflows (12 workflows)
- **MCP gateway session timeout** (#23153): long-running workflows at risk
- **Performance Regression in Validation** (#30180): 82.1% slower

## P2 (Watch)
- **UK AI Operational Resilience** (NEW): OTLP header masking failing activation (run 26012832575)
- **ET budget exhaustion**: Multiple daily workflows at risk. Audit `max-effective-tokens`. #32717
- **Engine-fail-after-completion** pattern persists (#32736) — systemic engine lifecycle bug
- **Step Name Alignment recurring daily**: #32955 opened May 18 (was #32754 closed May 17) — structural fix needed
- **[aw-compat] Cross-repo warnings** (#32528): P2

## Resolved (Do Not Re-File)
- PR-review cluster #31724: CLOSED ✅ (was ~272 wasted runs/day)
- May 14 mass failure (#32045-#32119): resolved ✅
- AWF Firewall v0.25.47 #32522: CONTAINED ✅
- Sergo #32755: CLOSED ✅
- Step Name #32754: CLOSED ✅ (recurred as #32955)

## Outlook
- **Critical blocker**: Agentic Maintenance compile failure — restore immediately to unblock orchestration
- CGO/CJS: 90+ days at P1, must be escalated with dedicated engineering time
- Quality/effectiveness plateau (day 17) at 74/71 — expected breakout to 76-78/73-75 once Agentic Maintenance restored
- Step Name Alignment: structural fix needed (daily noise)

Last updated: 2026-05-18T14:00Z by agent-performance-manager
