# Shared Alerts — 2026-05-17T13:00Z

## P1 (High)
- **CGO/CJS regression** (#29669): failing every push to main (90+ days)
- **Smoke CI** (#32690): open
- **Codex OPENAI_API_KEY sandbox exclusion** (#32446): blocking Codex workflows
- **MCP gateway session timeout** (#23153): long-running workflows at risk
- **Performance Regression in Validation** (#30180): 82.1% slower

## P2 (Watch)
- **ET budget exhaustion** (#32717): Daily Observability Report hit 80M token limit. 5-10 other daily workflows at risk. Audit `max-effective-tokens` across all.
- **Engine-fail-after-completion** (#32736): Workflow completes but safe-output not sent. New pattern.
- **Deployment Monitor zombie**: 100 runs/day × 8% success = 92 wasted invocations/day. Deprecate.
- **Daily Fact parse failures** (#31432, #31524): still failing post-PR#31411 merge — second root cause needed
- **May 17 clustered transient failures**: Sergo #32755, Step Name #32754, Linter Miner #32748, Outcome Collector #32728 — same 01:00-05:00Z window
- **[aw-compat] Cross-repo warnings** (#32528): P2

## Resolved (Do Not Re-File)
- PR-review cluster #31724: CLOSED ✅ (was ~272 wasted runs/day, primary Q/E plateau driver)
- May 14 mass failure (#32045-#32119): resolved by PR #32070 ✅
- AWF Firewall v0.25.47 #32522: CONTAINED ✅
- APM Unpack #30252: CLOSED ✅
- Smoke CI + Gemini #29666: CLOSED ✅
- Daily Model Inventory Checker #30043: CLOSED ✅

## Outlook
- Quality/Effectiveness plateau (Day 16 at Q:74/E:71) should break next run — cluster fix removes drag
- Expect Q→76-78, E→73-75 on May 18

Last updated: 2026-05-17T13:00Z by agent-performance-manager
