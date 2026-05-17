# Shared Alerts — 2026-05-11T13:36Z

## P0 (Critical)
- **APM Unpack Systemic Failure** (#30252 OPEN): apm-default.tar.gz exits code 1. Last updated 2026-05-05. 3 workflows blocked.

## P0 — RESOLVED (Do Not Re-File)
- Smoke CI + Gemini (#29666): CLOSED ✅ 2026-05-11
- Daily Model Inventory Checker (#30043): CLOSED ✅ 2026-05-11
- config.models (#30307): CLOSED ✅ 2026-05-11

## P1 (High)
- **Daily Fact About gh-aw**: 15+ push-time parse failures. Issue created 2026-05-11. PR #31411 merge should help.
- **Smoke macOS ARM64**: 100% failure since 2026-02-20 (81+ days). Issue filed 2026-05-07 ✅
- **CI regression on main**: TestStrictModePermissions failing. Issue filed 2026-05-06.
- **MCP gateway session timeout** (#23153 OPEN): Long-running workflows at risk.
- **Performance Regression in Validation** (#30180): 82.1% slower.

## P2 (Watch)
- **PR-review cluster** (Q, cloclo, Archie, Scout, Grumpy, Security Review, PR Nitpick, PR Code Quality): ~272 wasted run-attempts/day. HIGHEST WASTE. Trigger gate fix or consolidation needed.
- **on.labels push-time failures**: PR #31411 open fix. Merge unblocks systemic issue.
- **engine.max-runs migration**: PR #31418 open. Watch for compilation regressions after merge.
- **Deployment Incident Monitor**: zombie pattern — 8x skipped per 100 runs; consider deprecation.
- **Resource Summarizer Agent**: chronic skips, zero outputs.
- **Doc Build - Deploy**: action_required persistent (deployment stalled).
- **Node.js 20 deprecation** in CI: deadline Sep 16, 2026.
- **Quality/Effectiveness plateau**: 10 days flat (Q:74, E:71) — structural bottleneck suspected.

## Resolved (Do Not Re-File)
- #29863 Smoke Copilot regression → RECOVERED ✅
- #30205 Auto-Triage Issues → CLOSED ✅
- #30188 Documentation Unbloat → CLOSED ✅

---
## Update — 2026-05-12T05:39Z (Workflow Health Manager)

### RESOLVED
- APM Unpack #30252: CLOSED ✅
- PR #31411 (on.labels fix): MERGED ✅
- PR #31418 (engine.max-runs): MERGED ✅

### NEW/ESCALATED
- **Daily Fact parse failures**: still occurring post-PR#31411 merge; issues #31432 #31524 open.
- **Smoke Gemini fetch failed**: issue #31575 open — 100% failure.
- **Firewall reporting broken** (#31607, #31620): no safe outputs from agent job.
- **Multiple workflow failures** today: Design Decision Gate #31626, Go Logger Enhancement #31628, Step Name Alignment #31636, jsweep #31637.

---
## Update — 2026-05-12T13:19Z (Agent Performance Manager)

### RESOLVED (since May 11)
- APM Unpack #30252: CLOSED ✅
- PR #31411 (on.labels fix): MERGED ✅
- PR #31418 (engine.max-runs): MERGED ✅

### NEW/ESCALATED
- **4 same-day workflow failures** (2026-05-12): Design Decision Gate #31626, Go Logger Enhancement #31628, Step Name Alignment #31636, jsweep #31637 — possible shared root cause (engine availability or PR #31418 side-effect). Needs investigation.
- **Daily Fact still failing** post-PR#31411 merge (#31432, #31524 still open)
- **Quality/Effectiveness plateau**: Day 11 flat at Q:74/E:71 — structural bottleneck suspected (PR-review cluster waste dragging averages)
- **PR-review cluster waste escalated**: ~272 wasted run-attempts/day confirmed — highest waste in ecosystem; trigger gate fix is highest-ROI action

---
## Update — 2026-05-13T05:45Z (Workflow Health Manager)

### NEW (since May 12)
- **CI integration test failure**: Fix failing "Integration: Workflow Misc Part 2" (#31860) — CGO failing 3/4 runs. New regression.
- **Semantic Function Refactoring** (#31827): new agentic failure
- **Daily Security Red Team Agent** (#31817): new failure
- **Scout** (#31811): failed
- **Daily Cache Strategy Analyzer** (#31773): new failure
- **4 new workflows added** (219→223): no compilation issues

### RESOLVED
- Daily Firewall Logs Collector: auto-close ran ✅
- Smoke Copilot dispatch_workflow: auto-close ran ✅

### WATCH
- CI failures trending up (scheduled CI + CGO): possible integration test regression around #31860
- Deep-report triage issue #31729: 18 stale [aw]-failed issues — recommend bulk close

---
## Update — 2026-05-13T13:26Z (Agent Performance Manager)

### NEW (since May 13 morning)
- **5 new workflow failures**: PR Sous Chef #31931, Draft PR Cleanup #31929, Semantic Function Refactoring #31827, Daily Security Red Team #31817, Daily Cache Strategy #31773
- **Daily agents batch failure**: Daily Fact + Security Red Team + Cache Strategy all failing — suspected shared cron/runtime root cause; investigate together
- **CI cluster compound friction**: CGO (22%) + CJS (25%) + Smoke CI (50%) all below threshold — contributing to PR friction
- **Moderation twin failure**: Content Moderation + AI Moderator both at exactly 57% success — shared upstream instability suspected

### WATCH
- Quality/Effectiveness plateau: **Day 12** at Q:74/E:71 — primary driver is PR-review cluster structural waste (~272 runs/day at 0%)
- Open [aw] failed issues increasing: 8+ open
- Deep-report triage #31729 recommends bulk-closing 18 stale [aw]-failed issues

### KEY ACTION
- Highest ROI: Fix PR-review cluster trigger gate (#31724) — would recover ~272 wasted runs/day and likely break quality plateau

---
## Update — 2026-05-14T05:41Z (Workflow Health Manager)

### NEW/ESCALATED
- **CGO/CJS failing on every push**: persistent regression, failing across all recent merges to main. Issue #29669 (cgo-failure) is expired but open. 
- **Safe Output Health Monitor** (#32063): first failure after 9 consecutive successes — watch for recurrence.
- **Daily Grafana OTel** (#32066): isolated failure, likely transient.
- **30 open [aw] failure issues**: elevated noise level (was ~18 last week).

### STABLE/UNCHANGED
- MCP gateway session timeout (#23153): still open
- Performance Regression (#30180): still open
- Daily Fact parse failures (#31432, #31524): still open
- PR-review cluster waste (#31724): still open

---
## Update — 2026-05-14T13:26Z (Agent Performance Manager)

### P0 — NEW (May 14 Mass Failure Event)
- **33 [aw] failure issues created TODAY** (#32045–#32119): highest single-day count. Affected workflows include daily agents, code quality, smoke tests, moderators, and the Failure Investigator itself. Total open: 36.
- **Possible causes**: safe-output bundle infra disruption, engine availability window, or accumulated failures triggering batch issue creation
- **PR #32070** (safe output bundle fix) merged today — monitor for recovery

### STABLE/UNCHANGED
- CGO/CJS push regression (#29669): still failing every push to main
- Q + PR-review cluster (#31724): 0% success, ~272 wasted runs/day
- Quality/Effectiveness plateau: Day 13 at Q:74/E:71
- Daily Fact parse failures (#31432, #31524): still open
- MCP gateway session timeout (#23153): still open
- Performance Regression (#30180): still open

### KEY ACTION
- Human investigation needed: why did 33 workflows fail simultaneously on May 14?
- PR-review cluster fix (#31724) remains highest-ROI action to break quality plateau

---
## Update — 2026-05-15T05:43Z (Workflow Health Manager)

### RECOVERED (since May 14 mass failure)
- AI Moderator, Content Moderation, Agentic Commands: SUCCESS today ✅
- Smoke CI, Doc Build-Deploy, Safe Output Health Monitor: SUCCESS today ✅
- PR #32070 (safe output bundle fix) appears effective

### PERSISTENT
- **CGO/CJS regression** (#29669): still failing every push to main (P1)
- Daily Fact parse failures (#31432, #31524): not tested today (no cron run observed)
- MCP gateway session timeout (#23153): open
- Performance Regression (#30180): open
- PR-review cluster waste (#31724): structural, open

### WATCH
- 229 workflows (+4 new since last count); no compilation issues expected
- action_required rate high (70%): mostly no-trigger, not true failures
- Mass failure event (May 14, 33 issues) largely resolved by PR #32070

---
## Update — 2026-05-15T13:14Z (Agent Performance Manager)

### RECOVERED (since May 14 mass failure)
- AI Moderator, Content Moderation, Agentic Commands, Smoke CI, Doc Build-Deploy, Safe Output Health Monitor: all ✅ SUCCESS on May 15
- PR #32070 (safe output bundle fix) confirmed effective — health score 62→64/100

### PERSISTENT
- **CGO/CJS push regression** (#29669): still failing every push to main (P1)
- **Daily Fact parse failures** (#31432, #31524): still open
- **PR-review cluster** (#31724): 0% success, ~272 wasted runs/day — highest-ROI action
- MCP gateway timeout (#23153), Performance Regression (#30180): open

### ECOSYSTEM STATUS (May 15)
- Quality/Effectiveness plateau: **Day 14** at Q:74/E:71 — primary driver is PR-review cluster
- Workflows: 229 (+4 from 225)
- Open [aw] failure issues: ~30+ (recovering from 36 peak May 14)
- No new P0 events

### KEY ACTION
- Fix PR-review cluster trigger gate (#31724) remains highest-ROI action
- Schedule freeze-and-fix sprint to break 14-day quality plateau

---
## Update — 2026-05-16T05:35Z (Workflow Health Manager)

### NEW (since May 15 evening)
- **AWF Firewall v0.25.47 broken** (#32522): oidc-token-provider-base missing. CONTAINED — PR #32503 closed. Main unaffected. ✅
- **Smoke cluster failure at 01:00Z**: All smokes failed due to AWF firewall issue on closed PR. Recovering (Smoke OTEL ✅ at 05:28Z).
- **Smoke Codex** (#32561) + **Smoke Pi** (#32553): Still open as of 05:35Z.
- **AW Compat**: 18/19 repos green; microsoft/aspire hard failure (#32526).
- **[aw] Dev failed** (#32519): isolated failure, May 16 01:22Z.

### STABLE/PERSISTENT
- CGO/CJS push regression (#29669): still failing (P1)
- Daily Fact parse failures (#31432, #31524): open (P2)
- PR-review cluster waste (#31724): 0% success, ~272 wasted runs/day (P2, highest ROI)
- MCP gateway timeout (#23153): open
- Performance Regression (#30180): open

### SCORE: 64/100 (→ stable)

---
## Update — 2026-05-16T13:00Z (Agent Performance Manager)

### RECOVERED
- May 14 mass failure (33 [aw] issues): largely resolved by PR #32070 ✅
- Open [aw] failures: 19 (↓ from 36 peak)

### CONTAINED
- AWF Firewall v0.25.47 (#32522): isolated to PR branch, not merged to main ✅

### PERSISTENT
- **PR-review cluster** (#31724): ~272 wasted runs/day, 0% success — **Day 15 of quality plateau**
- **CGO/CJS push regression** (#29669): still failing every push to main (P1)
- **Daily Fact parse failures** (#31432, #31524): still open (P2)
- **MCP gateway timeout** (#23153): open
- **Performance Regression** (#30180): open
- **Smoke Codex** (#32561), **Smoke Pi** (#32553): open (new May 16)

### KEY ACTION
- Fix PR-review cluster trigger gate (#31724) remains highest-ROI action to break 15-day quality plateau

---
## Update — 2026-05-17T05:41Z (Workflow Health Manager)

### RESOLVED
- PR-review cluster #31724: CLOSED ✅ (~272 wasted runs/day eliminated)

### NEW/ESCALATED
- **ET budget exhaustion systemic**: Daily Observability Report (#32717) hit 80M ET limit. Multiple token-heavy daily workflows at risk. Recommend auditing `max-effective-tokens` config.
- **Engine-failure-after-completion** (#32736): Daily Compiler Quality Check completes work but engine exits without safe-output. Pattern to watch.
- **Codex OPENAI_API_KEY excluded** (#32446): P1, blocking Codex workflows.

### PERSISTENT
- CGO/CJS #29669 (P1), Smoke CI #32690 (P1), MCP timeout #23153 (P2), Perf #30180 (P2)
