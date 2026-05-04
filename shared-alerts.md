# Shared Alerts — 2026-05-04T05:39Z

## P0 (Critical)
- **Smoke Gemini** (#29852, #29816, #29459 OPEN): 100% failure, proxy architecture blocks all agent traffic. 30+ days unresolved.
- **Smoke CI** (#29666 OPEN): CGO/EROFS persistent, 100% action_required.
- **Smoke macOS ARM64**: 100% failure since 2026-02-20. NO ISSUE FILED — needs one urgently.
- **Daily Model Inventory Checker** (#30043 OPEN): Copilot CLI silent startup crash (new P0).

## P1 (High)
- **01:49 UTC transient wave May 4**: Smoke Claude, Pi, Codex, Copilot ARM64, OpenCode all failed simultaneously — likely shared runner/API issue
- **Metrics Collector failed** (#30050 OPEN)
- **MCP gateway session timeout** (#23153 OPEN): Long-running workflows at risk.

## P2 (Watch)
- **Node.js 20 deprecation** in CI: deadline Sep 16, 2026. Migrate to Node.js 22.
- **Safe Outputs SEC-004** (#27235 OPEN).
- **6 PR-review agents** on same triggers — evaluate redundancy (Scout, Archie, /cloclo, Q, AI Moderator, Content Moderation)

## Resolved (Do Not Re-File)
- #29863 Smoke Copilot regression → RECOVERED ✅ (success May 4 00:56)
- #29088 Codex crash → CLOSED
- #28659 Doc Unbloat claude auth → CLOSED
- #27965 GitHub Remote MCP Auth → CLOSED

## Trends
- 211 workflows (+2), 0 missing lock files
- Health: 65/100 (stable)
- Gemini still completely broken (30+ days, P0 unresolved)
- macOS ARM64 chronic failure since Feb 2026 — needs attention
