# Session Analysis Notes

## Recurring Observation: data_quality = "infrastructure-only"

Across 8 consecutive runs (2026-05-06 → 2026-05-16), conversation transcripts have NEVER been available from the `copilot-session-data-fetch` shared module. Today's `20260516-conversation.txt` is a 1-line OAuth error from `gh auth login`, confirming the auth-gated access path is the blocker. This is a known, persistent limitation rather than a regression. The workflow is effectively a CI gate-sweep / orphan-detection monitor, not a behavioral analyzer.

## Persistent Patterns (7-day window)

- Orphan rate: 0% every day for 5+ consecutive days (baseline is 40% per spec)
- Action_required share is the volatile metric: 8% → 16% → 100% (2026-05-13 spike)
- Branch concentration drives queue size: 2-7 unique branches; max gates on one branch has ranged 16-35
- Copilot agent success rate (when transcripts surface): 100% on most days, 57.9% on 2026-05-07

## Run: 2026-05-13

- 100% action_required (highest in window) — unusual spike, was 58% on 2026-05-12
- Workflow concentration: Scout/Q/Agentic Commands fire 12x each = 72% of queue
- 4 active branches: refactor-extract-safe-outputs-config (18), aw-failures-fix-daily-report-generator (14), 2x (9)
- All branches in warning band (1-2h wait), none critical yet

## Run: 2026-05-16

- 92% action_required (46/50) — back near the May-13 pattern after May 14/15 gap
- 6 active branches; top 2 absorb 66% of queue: update-claude-code-and-mcp-gateway (17), refactor-git-patch-header-parsing (16)
- 4 success runs: 2x "Running Copilot cloud agent" (top 2 branches) + 2x "Haiku Printer"
- 0 spec-orphans; 5 chaos/* PRs have no assignee but zero gates so they don't escalate
- Workflow fingerprint: Agentic Commands + Q = 26 firings (52%) — same concentration as May 13 but spread across 2 agent branches instead of 4
- Notable: both top branches got Copilot success runs AND a 5-fire gate burst at the success timestamp — sweep-after-success pattern

## Run: 2026-05-17

- 92% action_required (46/50) — identical share to 2026-05-16 (two-day-in-a-row run-stuck pattern)
- 6 active branches; queue is flatter than yesterday: top branch (scan-repeated-permission-denied-issues) holds only 11 of 50 (22%), no branch dominates
- 4 success runs are all real Copilot agent sessions on 4 distinct PR branches — every "investigate-safe-output*" branch + the agentic-workflow + the comment-addressing run on PR #32759 all completed
- Copilot agent duration spread: 8 / 11 / 15 / 20 min (avg 13.5 min) — wider than recent days but all in the >5-min "high-success" band per [historical_trend_regression strategy](session-analysis-strategies.json)
- Workflow fingerprint: Agentic Commands + Q = 24 firings (48%) — concentration relaxed slightly vs May-13/16
- 0 spec-orphans: only 2 in-progress runs system-wide and both on `main` (this workflow + Failure Investigator). The 3 chaos/* PRs from 05:54 have no assignee but zero in-progress gates, so they don't trip escalation. The `docs/update-dictation-glossary-*` PR is brand-new (06:37) and below the 1h waiting threshold.
- Notable: 3 of 4 success runs were close-paired (within ~6 min) with their `action_required` gate bursts on the same branch — same sweep-after-success pattern as 05-16 but tighter time spacing
- data_quality stays `infrastructure-only` for the 9th consecutive run (2026-05-06 → 2026-05-17). Conversation logs dir empty again.

## Open Action Items

- [ ] Investigate why conversation transcripts have never been delivered to /tmp/gh-aw/session-data/logs/
- [ ] Consider an "approval bottleneck" severity tier — strict orphan filter misses the dominant failure mode (gates stuck despite agent assignment)
- [ ] Once transcripts arrive, retroactively backfill prompt-quality scoring
