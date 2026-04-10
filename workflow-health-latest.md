# Workflow Health - 2026-04-10T12:05Z

Score: 75/100 (↑5 from 70 yesterday). 187 workflows. Run: §24241980102

## KEY FINDING: Copilot v1.0.22 Regression (RESOLVED)

Copilot CLI v1.0.22 had MCP blocking + /models silent failure bugs. 
- PR #25577 (05:26 UTC): bumped to v1.0.22 — BROKE workflows
- PR #25623 (11:50 UTC): re-pinned to v1.0.20 (stable)
- Current: v1.0.20 pinned, STABLE

5 workflows failed in regression window (10:35-11:27 UTC) — self-healing

## P1 Issues (RESOLVED)
- Copilot Engine Crash: RESOLVED. v1.0.20 is pinned stable.

## P2 Issues (Active)
- Copilot CLI Upgrade Blocked: v1.0.22 has MCP+models bugs. Need v1.0.23+ fix.
- Design Decision Gate failures (#25548): 2/3 runs failing
- Documentation Unbloat cost: ~$55/week (Claude), no safe outputs
- Contribution Check report_incomplete: Every run

## Open Failure Issues (15)
Key: #25215 (Auto-Triage, 25 comments), #25374 (Smoke Copilot), #25371 (Agent Container Smoke Test)
Most from Apr 8-9 crash window, should auto-close as workflows recover with v1.0.20

## Score Breakdown
- Compilation: 187/187: +35
- v1.0.20 stable: +20
- 5 temporary failures (self-healing): -5
- 15 open failure issues: -8
- Copilot upgrade blocked: -7
- Net: 75/100

## Score Trend
68 → 71 → 73 → 71 → 70 → 75 (Apr 5–10)

## Dashboard Issue
#25470 (updated this run)

Last updated: 2026-04-10T12:05Z
