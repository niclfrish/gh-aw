# Workflow Health — 2026-05-17T05:41Z

Score: 67/100 (↑ from 64). 229 workflows. Run: §25982599702

## KEY FINDINGS

### Today's New Issues (May 17)
- **Sergo - Serena Go Expert failed** (#32755): open
- **Step Name Alignment failed** (#32754): open
- **Linter Miner failed** (#32748): open
- **Daily Compiler Quality Check failed** (#32736): engine failure after completing work
- **Outcome Collector failed** (#32728): open
- **Daily Observability Report failed** (#32717): ET budget exhausted (80M limit)

### 🎉 Resolved Since May 16
- **PR-review cluster #31724 CLOSED** ✅ — trigger gates fixed for 8 workflows (Q, Agentic Commands, CGO, CJS, Smoke CI, Doc Build-Deploy, AI Moderator, Content Moderation). Was ~272 wasted runs/day at 0% success.

### Persistent Issues (Unchanged)
- **CGO/CJS regression** (#29669): still open, failing on every push (P1)
- **Smoke CI** (#32690): still open (P1)
- **Codex OPENAI_API_KEY sandbox exclusion** (#32446): P1
- **MCP gateway session timeout** (#23153): P2
- **Performance Regression** (#30180): P2
- **[aw-compat] Cross-repo warnings** (#32528): P2

### Systemic Patterns
- **ET budget exhaustion**: Multiple daily workflows hitting token limits (Daily Observability Report at 80M)
- **Engine failure after completion**: Workflow completes but safe-output not sent (#32736)

### Open [aw] failures
21 open (↑2 from 19 on May 16; PR cluster fix reduces future waste)

### Actions Taken This Run
- Added comment to dashboard issue #29109
- Updated shared memory

### Trends
- Score: 67/100 (↑3 from PR-review cluster fix)
- PR-review cluster: RESOLVED ✅ (~272 wasted runs/day eliminated)
- Top performers stable: Issue Monster, PR Sous Chef, Daily Semgrep, CI
