# Agent Performance - 2026-04-11
Run: §24274655501 | Q:70↑5 E:60↓6

Top: CLI Version Checker (Q:90 E:88 ↑2), AI Moderator (Q:88 E:92 - Codex), Smoke Claude (Q:88 E:90), Issue Monster (Q:82 E:85 - 5 runs/100%), Agentic Maintenance (Q:80 E:78 - 6 deep-report issues)

Watch: Documentation Unbloat (0 safe outputs, ~$55/wk), Design Decision Gate (root cause: --print flag empty prompt in #25670, P2), Contribution Check (report_incomplete every run), GitHub Remote MCP Auth Test (100% fail today), Workflow Normalizer (3 duplicate issues in 24h - deduplication gap)

Smoke Tests: Claude/Codex/Multi-PR ✅ | Copilot ⚠️ recovering (21/30 fail, v1.0.24 PR #25752 in progress) | Gemini ❌ | Cross-Repo Create/Update ❌

New: v1.0.24 bump issued today (#25751/PR #25752) - will unblock Copilot smoke tests

Issues this week: 30 total | deps:5 (Dependabot), deep-report:6 (Agentic Maintenance), refactoring:4 (Skill Extractor), workflow-style:3 (Normalizer-DUPE), ca:3 (CLI Version), plan:3 (Plan Command)
PRs merged: ~30 (mostly by Copilot bot) | Notable lgtm: #25721, #25693, #25628, #25618

Engine split: copilot:~124wf, claude:~41wf, codex:~18wf, others:~4wf
187 compiled workflows. Health: 75/100 (recovering). 20/25 scheduled workflows healthy today.

Actions: Weekly discussion created. No new improvement issues (existing tracking: #25548 DDG, alerts in shared-alerts.md).
