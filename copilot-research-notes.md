# Copilot CLI Research Notes (Trimmed - last 5 runs)

### 2026-05-18 (Run 26014468484) — This Run
- 230 total MD workflows (+1); 126 Copilot (114 simple + 29 extended block = 55%)
- **engine.agent**: 11 workflows (drop from 25 — volatility in count due to AWF vs engine field ambiguity)
- **max-continuations**: 5 workflows (stable: contribution-check:20, test-quality-sentinel:15, mattpocock-skills-reviewer:10, smoke:2, smoke-otel-backends:1)
- **bare mode**: 11 workflows (stable: ab-testing-advisor, constraint-solving-potd, daily-fact, daily-hippo-learn, daily-news, hippo-embed, outcome-collector, poem-bot, smoke-claude, smoke-copilot, smoke-otel-backends)
- **cache-memory**: 73 workflows (32% — stable)
- **repo-memory**: 23 workflows (10% — moderate usage)
- **sandbox AWF**: 16 workflows (  agent: awf patterns)
- **web-search/fetch**: 21 workflows (up significantly from 2!)
- **model overrides**: 48 total (43 model:small, 3 model:large, rest specific — significant use of convenience aliases)
- **engine.args**: 0 (PERSISTENT GAP, 10th+ consecutive run)
- **engine.env**: 0 (persistent gap)
- **engine.api-target**: 0 (PERSISTENT GAP, 11th consecutive run)
- **engine.harness**: 0 (persistent gap)
- **mcp-servers (frontmatter)**: 2 workflows (minimal use)
- **mcp-scripts (frontmatter)**: 0 workflows (used in prompts but not as frontmatter config)
- **BYOK**: 0 (persistent gap)
- **experiments**: 16 workflows (A/B testing moderately used)
- **Unused agent files**: grumpy-reviewer, interactive-agent-designer, w3c-specification-writer, create-safe-output-type, custom-engine-implementation (5 files still unused)
- **max-runs**: 2 workflows only (daily-safe-output-optimizer:200, linter-miner:1000) — severely underused
- **playwright**: 11 workflows (growing browser automation)
- Discussion created: "Copilot CLI Deep Research - 2026-05-18"

### 2026-05-17 (Run 25981819267)
- 229 total MD workflows; 126 Copilot (99 simple + 27 object form = 55%)
- **engine.agent**: 25 workflows (up from 14 last run — growing strongly again; AWF-only + 10 custom agents)
- **max-continuations**: 5 workflows (down from 6; contribution-check:20, test-quality-sentinel:15, mattpocock-skills-reviewer:10, smoke:2, one other)
- **bare mode**: 11 workflows (up from 10 — slow steady growth)
- **cache-memory**: 73 workflows (down from 94 — notable drop, likely measurement difference)
- **sandbox AWF**: 19 workflows (stable)
- **web-search**: 2 workflows (stable)
- **version pinning**: 0 (was 10 last run — alarming drop, verify next run)
- **model overrides**: 18 (stable, mostly smoke tests and experiments)
- **engine.args**: 0 (persistent gap, 9th+ consecutive run)
- **engine.env**: 0 (persistent gap)
- **engine.api-target**: 0 (persistent gap, 10th consecutive run)
- **engine.harness**: 0 (persistent gap)
- **mcp-scripts (frontmatter)**: 0 (mcpscripts tools used in prompts but not as frontmatter tools)
- **BYOK**: 0 (persistent gap)
- **Unused agent files**: grumpy-reviewer, interactive-agent-designer, w3c-specification-writer, create-safe-output-type, custom-engine-implementation (5 files still unused)
- **max-runs**: 1 workflow (daily-safe-output-optimizer:200) — severely underused
- Discussion created: "Copilot CLI Deep Research - 2026-05-17"

### 2026-05-16 (Run 25953071091)
- 229 total MD workflows; 128 Copilot (99 simple + 29 object form = 56%)
- **engine.agent**: 14 workflows (down from 25 — stabilizing; AWF-only agents dominate)
- **max-continuations**: 6 workflows (up from 4 — contribution-check:20, test-quality-sentinel:15, mattpocock-skills-reviewer:10, smoke-copilot:2, smoke-otel:1, smoke-otel-backends:1)
- **bare mode**: 10 workflows (stable)
- **cache-memory**: 94/229 (up from 92, widespread adoption ~41%)
- **sandbox AWF**: 19 workflows (down from 21)
- **web-search**: 2 workflows (up from 0 — new adoption!)
- **version pinning**: 10 workflows (up from 0 last run — strong rebound)
- **model: small**: 6 copilot workflows; **model: large**: 3
- **engine.args**: 0 (persistent gap, 8th+ consecutive run)
- **engine.env**: 0 (persistent gap)
- **engine.api-target**: 0 (persistent gap, 9th consecutive run)
- **engine.harness**: 0 (persistent gap)
- **BYOK**: 0 (persistent gap)
- **Unused agent files**: grumpy-reviewer, interactive-agent-designer, w3c-specification-writer, create-safe-output-type, custom-engine-implementation (5 files)
- **max-runs**: 1 workflow (daily-safe-output-optimizer:200) — severely underused
- Discussion created: "Copilot CLI Deep Research - 2026-05-16"

### 2026-05-14 (Run 25842508637)
- 225 total MD workflows; 121 Copilot (97 simple + 24 object engine form = 54%)
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

### 2026-05-13 (Run 25779191470)
- 223 total MD workflows; 121 Copilot (95 simple + 26 object engine form = 54%)
- **engine.agent**: 15 workflows (grew from 7 last run!) — accelerating adoption
- **max-continuations**: 4 workflows — stable
- **bare mode**: 1 workflow (smoke-copilot only)
- **mcp-scripts**: 4 workflows
- **sandbox AWF**: 11 workflows
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
