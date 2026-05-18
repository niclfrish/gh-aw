# Shared Alerts — 2026-05-18T05:52Z

## P1 (High)
- **CGO/CJS regression** (#29669): failing every push to main (90+ days, 3/3 runs this week)
- **Agentic Maintenance compile failure** (NEW): compile-workflows down 2 runs — orchestrator impaired
- **Codex OPENAI_API_KEY sandbox exclusion** (#32446): blocking Codex workflows
- **MCP gateway session timeout** (#23153): long-running workflows at risk
- **Performance Regression in Validation** (#30180): 82.1% slower

## P2 (Watch)
- **UK AI Operational Resilience** (NEW): OTLP header masking failing activation (run 26012832575)
- **ET budget exhaustion**: Multiple daily workflows at risk. Audit `max-effective-tokens`.
- **Engine-fail-after-completion** pattern persists
- **Step Name Alignment recurring daily**: #32754 closed May 17, #32955 opened May 18 same day — structural fix needed
- **[aw-compat] Cross-repo warnings** (#32528): P2

## Resolved (Do Not Re-File)
- PR-review cluster #31724: CLOSED ✅ (was ~272 wasted runs/day)
- May 14 mass failure (#32045-#32119): resolved ✅
- AWF Firewall v0.25.47 #32522: CONTAINED ✅
- Sergo #32755: CLOSED ✅
- Step Name #32754: CLOSED ✅ (but recurred as #32955)

## Outlook
- Agentic Maintenance down is new critical blocker — affects all workflow management
- Step Name Alignment recurring daily needs structural investigation
- CGO/CJS need urgent attention — 90+ days unresolved

Last updated: 2026-05-18T05:52Z by workflow-health-manager
