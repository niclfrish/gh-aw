# Workflow Health — 2026-04-29T12:24Z

Score: 73/100 (↑+16 from 57 Apr 28). 204 workflows. Run: §25108329742

## KEY FINDINGS

### Compilation Status
- 204/204 lock files present ✅
- **0 missing lock files** ✅
- **0 stale lock files** ✅

### Today's Failures (Apr 29)
8 scheduled runs failed out of 30 (73% success rate — improvement from 57% yesterday)

**Category 1: Codex engine crash (P0 ongoing)**
- **Daily Fact About gh-aw** — `codex: command not found` (auto-issue #29088)

**Category 2: CI integration test failures (P1)**
- **CI** — 4 jobs failed: js-integration-live-api (`ERR_API: fetch file audit-workflows.md failed`), Integration Release Availability, DIFC Proxy sh Integration Test, Integration Update - Preserve Local Imports

**Category 3: Safe outputs failures (likely agent crashes)**
- **Daily Rendering Scripts Verifier** — safe_outputs failed, no OTEL (Docker/Playwright env)
- **Daily Go Function Namer** — safe_outputs failed
- **Developer Documentation Consolidator** — safe_outputs failed (Docker env)
- **Daily AstroStyleLite Markdown Spellcheck** — safe_outputs failed, no OTEL
- **Instructions Janitor** — safe_outputs failed, no OTEL
- **Daily AW Cross-Repo Compile Check** — safe_outputs failed

### Improvements vs Yesterday
- THREAT_DETECTION_RESULT parse failure: NOT appearing in today's failures (was systemic yesterday - resolved or intermittent)
- Total failure count: 8 (down from 13)
- Success rate: 73% (up from 57%)

### P0 Issues (Active)
- **Daily Fact About gh-aw codex failure** (#29088 auto-created today): codex binary not found. Daily recurring.

### P1 Issues (Carry from Apr 28)
- **CI integration tests failing** (new today): js-integration-live-api `ERR_API: fetch file audit-workflows.md failed`
- **Node.js 20 deprecation warning** in CI: actions/setup-go using deprecated Node.js 20 (end-of-life Sep 2026)
- **Documentation Unbloat claude auth failure** (#28659 OPEN)
- **GitHub Remote MCP Authentication Test** (#27965 OPEN): day 8+
- **Safe outputs session not found** (#23153 OPEN)
- **awf-api-proxy sidecar unhealthy** (#27888 OPEN)
- **GitHub App rate limit** (#27251 OPEN)
- **CODEX_HOME collision** (#27512 OPEN)

### P2 Issues
- **Safe Outputs SEC-004** (#27235 OPEN)
- **Performance regressions** (#27280/#27279/#27278 OPEN)
- **Daily Documentation Updater protected files** (#27801 OPEN)
- **MCP gateway long-running drops** (#23153 OPEN)

## Issues Created This Run
- None (codex failure tracked in auto-issue #29088; CI issues tracked in existing #28659)

## Issues Updated
- None

## Positive Notes
- 204/204 workflows compiled, all lock files present
- THREAT_DETECTION_RESULT systemic failure from Apr 28 did NOT recur today
- Success rate improved significantly: 57% → 73%
