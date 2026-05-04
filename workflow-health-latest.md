# Workflow Health — 2026-05-04T05:39Z

Score: 65/100 (→ stable from 65). 211 workflows. Run: §25302920193

## KEY FINDINGS

### Compilation Status
- 211/211 lock files present ✅ (+2 new workflows)
- 0 missing lock files ✅

### P0 Issues (Active)
- **Smoke Gemini** (#29852, #29816, #29459 OPEN): 100% failure, proxy blocks traffic (chronic 30+ days)
- **Smoke CI** (#29666 OPEN): 100% action_required, EROFS crash (chronic)
- **Smoke macOS ARM64**: 100% failure since Feb 2026, NO ISSUE FILED — needs one
- **Daily Model Inventory Checker** (#30043 OPEN): New P0 — Copilot CLI silent startup crash

### P1 Issues (Active)
- **Smoke Claude** (#29974 OPEN): 1 failure at 01:49 UTC May 4 (transient wave)
- **Smoke Pi** (#29973 OPEN): 1 failure at 01:49 UTC May 4 (transient wave)
- **Smoke Codex**: 1 failure at 01:49 UTC May 4 (transient wave, no dedicated issue yet)
- **Smoke Copilot ARM64**: 1 failure at 01:49 UTC May 4 (transient wave)
- **Smoke OpenCode**: 1 failure at 01:49 UTC May 4 (transient wave)
- **Metrics Collector** (#30050 OPEN)
- **Auto-Triage Issues** (#30039 OPEN)
- **Dev** (#29969 OPEN)
- **Documentation Unbloat** (#29964 OPEN)
- **Step Name Alignment** (#30069 OPEN)

### P2 Issues
- Node.js 20 deprecation deadline Sep 16, 2026
- MCP gateway session timeout (#23153)
- 6 PR-review agents with approval queue backlog

### Recovery
- Smoke Copilot: SUCCESS at 00:56 UTC ✅ (was failing May 3)

### Actions Taken This Run
- Created health dashboard issue
- Updated shared memory

### Trends
- Score: 65/100 (→ stable)
- 01:49 UTC transient wave: 5 smoke tests failed simultaneously
- Gemini still completely broken (30+ days, P0 unresolved)
- macOS ARM64 chronic failure since Feb 2026 — needs attention
