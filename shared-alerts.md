# Shared Alerts — 2026-04-10T12:05Z

## P2 (High)
- **Copilot CLI Upgrade Blocked**: v1.0.22 broke MCPs + /models. v1.0.20 re-pinned. Need v1.0.23+ verified before upgrade. Checklist in #25623.
- **Documentation Unbloat cost** (ongoing): Claude workflow ~$55/week, 0 safe outputs.
- **Design Decision Gate failures** (#25548): 2/3 runs failing. Architecture decisions blocked.

## P3 (Watch)
- **Contribution Check report_incomplete**: Every run. Permission/network investigation needed.

## Recent Fixes
- Copilot v1.0.21 crash RESOLVED by pinning to v1.0.20 (Apr 8–10 saga)
- v1.0.22 regression (05:26–11:50 UTC Apr 10): 5 workflow failures, now self-healing
- #25022 AI Moderator missing_data: CLOSED not_planned Apr 9
- #24718 Duplicate Code Detector: CLOSED not_planned Apr 6
- #24829 GitHub Remote MCP Auth: CLOSED not_planned Apr 7

## Active Failure Issues (15)
Key open: #25215, #25374, #25371, #25443, #25276, #25372, #25259, #25478, #25480, #25415, #25395, #25432, #25456, #25469, #25470
Dashboard: #25470 (updated Apr 10)

## Ecosystem State
- 187 compiled workflows. Engine split: ~124 copilot, ~41 claude, ~18 codex, ~4 others
- v1.0.20 pinned as stable Copilot version
- Claude/Codex engines showing 100% resilience
- Expect Copilot workflow recovery over next 24h

Last updated: 2026-04-10T12:05Z by workflow-health-manager
