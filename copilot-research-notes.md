# Copilot CLI Research Notes (Trimmed - last 5 runs)

### 2026-05-14 (Run 25842508637) — This Run
- 225 total MD workflows; 121 Copilot (97 simple + 24 object form = 54%)
- **engine.agent**: 25 workflows (growing: 15→18→25) — strong adoption acceleration
- **max-continuations**: 4 workflows (stable: contribution-check:20, test-quality-sentinel:15, mattpocock-skills-reviewer:10, smoke-copilot:2)
- **bare mode**: 10 workflows (up from 1 last run, big jump)
- **mcp-scripts**: 12 workflows (significant growth: 1→4→12)
- **sandbox**: 21 workflows
- **token-steering**: 32 workflows (growing)
- **cache-memory**: 92/225 (41% — widespread adoption)
- **pre-agent-steps**: 9 workflows (new stat)
- **ab-testing**: 12 workflows (new stat)
- **inline-sub-agents**: 0 (persistent gap)
- **engine.api-target**: 0 (persistent gap, 7th consecutive run)
- **engine.harness**: 0 (persistent gap)
- **BYOK**: 0 (persistent gap)
- **version pinning**: 0 (persistent gap — dropped from 2 last run!)
- **Unused agent files**: grumpy-reviewer, interactive-agent-designer, w3c-specification-writer, create-safe-output-type, custom-engine-implementation (5 files)
- Discussion created: "Copilot CLI Deep Research - 2026-05-14"

### 2026-05-13 (Run 25779191470)
- 223 total MD workflows; 121 Copilot (95 simple + 26 object engine form = 54%)
- **engine.agent**: 15 workflows (grew from 7 last run!) — accelerating adoption
- **max-continuations**: 4 workflows (contribution-check:20, test-quality-sentinel:15, mattpocock-skills-reviewer:10, smoke-copilot:2) — stable
- **bare mode**: 1 workflow (smoke-copilot only)
- **mcp-scripts**: 4 workflows
- **sandbox AWF**: 11 workflows
- **strict mode**: 63/96 copilot (66%)
- **cli-proxy**: 86/96 copilot (89%)
- **repo-memory tool**: 14/96 (15%)
- **engine.api-target**: 0 (6th consecutive run)
- **engine.harness**: 0; **BYOK**: 0; **version pinning**: 0

### 2026-05-12 (Run 25714049123)
- 219 total MD workflows; 96 Copilot (44%)
- **max-continuations**: 4; **model overrides**: 27; **engine.agent**: 7
- **mcp-scripts**: 5; **web-fetch**: 20; **cache-memory**: 10; **sandbox AWF**: 20; **bare**: 9
- **engine.api-target**: 0 (5th consecutive run); **engine.harness**: 0; **BYOK**: 0

### 2026-05-11 (Run 25651194663)
- 218 total MD workflows; ~115 Copilot (95 simple + 20 block with id: copilot)
- **engine.agent**: 18; **model overrides**: 13; **max-continuations**: 2
- **mcp-scripts**: 1; **version pinning**: 2; **bare**: 9; **sandbox AWF**: 19; **cache-memory**: 89

### 2026-05-10 (Run 25620196538)
- 218 total MD workflows; 96 Copilot (44%)
- **max-continuations**: 2; **engine.api-target**: 0; **engine.harness**: 0
- **cache-memory**: 89/218 (massive growth); **sandbox AWF**: 19/218 (+73%)

## Persistent Gaps (7+ consecutive runs with 0 usage)
1. **engine.api-target** — Enterprise API endpoint override — never used
2. **engine.harness** — Custom harness script replacement — never used
3. **BYOK (Bring Your Own Key)** — Custom provider keys — never used
4. **inline-sub-agents** (engine config level) — never used

## Trending Positively
- engine.agent: 7 → 15 → 25 (strong growth)
- bare mode: 1 → 10 (sudden jump)
- mcp-scripts: 1 → 4 → 12 (fast growth)
- cache-memory: stable at ~90 (widespread)
- token-steering: new tracking, 32 workflows already
