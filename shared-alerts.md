# Shared Alerts — 2026-04-14T12:10Z

## P2 (High)
- **Smoke Claude schedule failure** (ongoing, #25727): Still failing on daily schedule, passes on PR runs. Environment-specific — schedule vs PR context divergence.
- **Smoke Cross-Repo PR Create/Update** (#25221, #25217, Apr 8): STALE 6 days. No fix applied. Needs escalation.
- **Schema Feature Coverage Checker** (#25992, Apr 13): Protected-files config blocks PR creation to .github/workflows/schema-demo-*.md. Fix: add `protected-files: fallback-to-issue` to frontmatter.
- **Documentation Unbloat inconsistent** (ongoing): ~$55/week Claude, 50% success. Cost gate needed.
- **Daily Firewall Logs** (#25456): safe_outputs process failure.
- **Smoke Multi PR** (#25415): persistent failure.
- **GitHub Remote MCP Auth Test**: 100% failure — #24829 closed not_planned. Test still failing.

## P3 (Watch)
- **Smoke Gemini**: 100% failure. Gemini 0.37.2 now available (#26158) — may fix after CLI bump.
- **~16 other daily workflow failures**: From Apr 8-13, mostly engine/copilot crashes.
- **Daily Issues Report recurring failure** (#25265, #25503): Copilot agent crash pattern.

## Copilot Version Status
- v1.0.25 NOW AVAILABLE (new --remote/--no-remote flags; see #26158)
- v1.0.21 ACTIVE (current in production as of Apr 14)
- Claude Code 2.1.105 available (#26158)
- Codex 0.120.0 available (#26158)
- Gemini 0.37.2 available (#26158) — may fix Smoke Gemini

## Recoveries (Apr 11-14)
- ✅ Smoke Copilot: RECOVERED
- ✅ Contribution Check: RECOVERED
- ✅ Agent Persona Explorer: IMPROVED (safe-output instructions strengthened, #26152)
- ✅ 20 PRs merged by Copilot bot (OTel, security, workflow fixes)

## Ecosystem State
- 191 compiled workflows (+4 since Apr 13). Health: 73/100 (↓1 Apr 14 12:10Z)
- Engine split: ~124 copilot, ~41 claude, ~18 codex, ~4 others
- v1.0.21 currently active, v1.0.25 available for upgrade

## Orchestrator Summaries (Apr 14)
- Agent Performance (Apr 14 04:37Z): Q:74↑1 E:66↑1. CLI Version Checker standout (4 upgrades).
- Workflow Health (Apr 14 12:10Z): Score 73/100 (↓1). 23 open failure issues. 191 workflows all compiled.
- Campaign Manager (last known: Mar 16 17:41Z): Status unknown — no recent update

Last updated: 2026-04-14T12:10Z by workflow-health-manager

## Update 2026-04-15T04:37Z (Agent Performance)
- CLI Version Checker: 5 tools upgraded (Claude 2.1.109, Copilot 1.0.27, Codex 0.120.0, Gemini 0.38.0, GitHub MCP v0.33.1) → #26367 open
- GitHub Remote MCP Auth Test: 100% failure continues — consider deprecation (see #24829 closed not_planned)
- Auto-Triage Issues: Newly failed — needs investigation
- Documentation Unbloat: 1/1 success today (improved from 50% historical)
- MCP Rate-Limit: #26239 open (circuit breaker request) — P2 risk

Last updated: 2026-04-15T04:37Z by agent-performance-analyzer

## Update 2026-04-15T12:09Z (Workflow Health Manager)
- **18 stale lock files** (NEW, Apr 15): PR #26372 (context propagation) merged without recompile. Run `make recompile`.
- **Daily Issues Report:** `node: command not found` recurring (#26393, Apr 15) — Node.js missing in Copilot runner
- **Auto-Triage Issues** (NEW failure Apr 15): #26364, 2x fails today — was recovered Apr 13-14
- **Daily Fact / Daily News:** new failures Apr 15 (#26405, #26388) — may be stale lock related
- **Smoke Gemini (#26351):** still failing — Gemini 0.38.0 now available, awaiting CLI upgrade

### Recoveries (Apr 15)
- ✅ Contribution Check: successful report #26376 (4 PRs reviewed)
- ✅ Architecture Diagram: successful #26389

### Version Status (Apr 15)
- Copilot v1.0.27 available (upgrade PR #26367 open)
- Claude 2.1.109, Codex 0.120.0, Gemini 0.38.0, GitHub MCP v0.33.1 available
- v1.0.21 still active in production

### Ecosystem State (Apr 15 12:09Z)
- 191 workflows, 18 stale locks ⚠️. Health: 72/100 (↓1). ~27 open failure issues.
- Dashboard: see safeoutputs #aw_dash15

Last updated: 2026-04-15T12:09Z by workflow-health-manager

## Update 2026-04-16T04:41Z (Agent Performance)
- **CLI Version Checker failed Apr 16** (§24492016365): MCP `github-agentic-workflows` server failed → 40 turns, $1.19 cost, no PR created. Infrastructure issue.
- **GitHub Remote MCP Auth Test**: Deprecation still pending — #24829 closed not_planned, workflow still consuming compute daily.
- **Auto-Triage Issues**: ~30-40% failure rate continues. Intermittent MCP/engine issues.
- **Documentation Unbloat**: Failed again Apr 16 (#26554). Historical 50% success. Cost gate needed.
- **18 stale lock files**: Still unresolved from PR #26372 (Apr 15). No automated recompile trigger on merge.
- **Smoke Codex**: New failure #26523 (Apr 16). 
- **CLI upgrade PR #26367 still open**: Would fix Smoke Gemini if merged.

### Ecosystem: Q:76↑1 E:69↑1 | Health: 72/100

Last updated: 2026-04-16T04:41Z by agent-performance-analyzer
