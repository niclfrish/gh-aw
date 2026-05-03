# Shared Alerts — 2026-05-03T05:38Z

## P0 (Critical)
- **Smoke Gemini** (#29459, #29852 OPEN): 100% failure, API_KEY_INVALID. Every PR sees red.
- **Smoke Copilot/Claude regression** (#29863, #29864 OPEN): New regression ~03:37 UTC May 3. 2-3/5 failures.
- **Smoke CI** (#29666 OPEN): Crush EROFS, 4/5 action_required.
- **CGO build** (#29669 OPEN): ongoing failures.

## P1 (High)
- **Smoke Codex**: 1/5 failures (minor regression today)
- **MCP gateway session timeout** (#23153 OPEN): Long-running workflows at risk.

## P2 (Watch)
- **Node.js 20 deprecation** in CI: deadline Sep 16, 2026. Migrate to Node.js 22.
- **YAMLGeneration regression** (#29779): 21.7% slower.
- **Safe Outputs SEC-004** (#27235 OPEN).

## Resolved (Do Not Re-File)
- #29088 Codex crash → CLOSED
- #28659 Doc Unbloat claude auth → CLOSED
- #27965 GitHub Remote MCP Auth → CLOSED
- #27888 awf-api-proxy sidecar → CLOSED
- #27251 GitHub App rate limit → CLOSED
- #27512 CODEX_HOME collision → CLOSED

## Trends
- 209 workflows, 0 missing lock files (+2 new)
- Score: 65/100 ↓3
- New smoke regression wave starting 03:37 UTC May 3 (Copilot, Claude, Codex)
- Gemini still completely broken (100% failure, ongoing P0)
