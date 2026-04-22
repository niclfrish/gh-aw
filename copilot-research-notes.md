# Copilot CLI Research Notes

## Analysis History

### 2026-04-21 (Run 24746483988)
- 197 total MD workflows; 87 explicit copilot + 24 default = 111 total Copilot effective
- 46 Claude, 10 Codex workflows
- **Corrections from prev**: prev counted `agent: awf` as custom agent files (inflated to 21); actual custom agent file use = 7 (stable)
- **Stable gaps** (persistent 5–7 days): engine.version (0%), api-target (0%), blocked-domains (0%), mcp-gateway (0%), mcp-scripts (1 workflow), max-continuations (2 workflows)
- 5/11 custom agent files still unused: grumpy-reviewer, w3c-specification-writer, create-safe-output-type, custom-engine-implementation, interactive-agent-designer
- 45 Copilot workflows without any network config (no AWF, no network: block)
- 45 workflows using only `toolsets: [default]` (over-provisioned GitHub access)

### 2026-04-20 (Run 24690376692)
- 197 total MD workflows, 90 using engine: copilot (explicitly)
- **Major improvement**: engine.agent adoption: 7→21 (+200%) - more workflows using custom agent files [NOTE: this was inflated due to awf sandbox agent refs]
- **Improvement**: cache-memory: 80→99 (+24%) - persistent data usage growing
- **Improvement**: strict_mode: 115→131 (+14%) - more workflows using security mode
- **First adoption**: engine.args: 0%→5% and engine.env: 0%→2% (custom CLI args now used)
- **Persistent gaps (5+ days)**: engine.version (0%), api-target (0%), mcp-gateway (0%), blocked-domains (0%)
- 5/11 custom agent files still unused: grumpy-reviewer, w3c-specification-writer, create-safe-output-type, custom-engine-implementation, interactive-agent-designer

### 2026-04-17 (Run 24586698669)
- 85 copilot workflows tracked
- Notable: engine.args/env at 0%, mcp-gateway at 0%

### 2026-04-16 (Run 24534029243)
- 192 total, 90 explicit copilot + 26 default = 116 effective
- playwright regression: 20→12 (-40%)
- strict_mode: 111→126 (+13%)

## Persistent Opportunities (Not Addressed in 5+ Runs)

1. **engine.version**: Never used → stability risk for critical workflows
2. **engine.api-target**: Never used → GHEC/GHES teams can't use this
3. **token-weights**: Never used → no custom cost modeling
4. **blocked-domains**: Never used → missed defense-in-depth layer
5. **mcp-scripts**: 1 workflow (security-review.md) → underutilized dynamic MCP capability
6. **max-continuations**: 2 only → Copilot-unique autopilot for complex tasks underused
7. **5/11 custom agent files unused**: grumpy-reviewer, w3c-specification-writer, create-safe-output-type, custom-engine-implementation, interactive-agent-designer

## Recommendations Tracking

| Recommendation | Status | Date Added |
|---|---|---|
| Use engine.version pinning for reproducibility | ⏳ Pending | 2026-04-16 |
| Expand max-continuations for complex/long workflows | ⏳ Pending | 2026-04-16 |
| Use bare:true for simple/creative/analytical workflows | ⏳ Pending | 2026-04-16 |
| Add network.blocked for defense-in-depth | ⏳ Pending | 2026-04-17 |
| Activate unused agent files | ⏳ Pending | 2026-04-16 |
| Model override for cost optimization | ⏳ Pending | 2026-04-17 |
| Add network config to the 45 unrestricted workflows | ⏳ Pending | 2026-04-21 |
| Tighten toolsets beyond [default] | ⏳ Pending | 2026-04-21 |
| Cache-memory adoption growing | ✅ Improving | 2026-04-16 |
| Custom agent adoption | ✅ Stable at 7 | 2026-04-21 |
| engine.args/env adoption | ✅ Achieved (5%) | 2026-04-20 |

### 2026-04-22 (Run 24802849397)
- 197 total workflows; 87 explicit copilot; 111 total Copilot effective
- **Trending up**: cache-memory (+4%), strict mode (+3%), mcp-cli (+2%), AWF sandbox (+6%), bare mode (+2%)
- **Stable/persistent gaps**: engine.version (0%), api-target (0%), startup-timeout (0%), tool-timeout (0%), network.blocked (0%), max-continuations (1%)
- **Confirmed: 5/11 custom agent files still unused**: grumpy-reviewer, w3c-specification-writer, create-safe-output-type, custom-engine-implementation, interactive-agent-designer
- 34 workflows (39%) have no network config at all
- **New finding**: startup-timeout and tool-timeout have been features for multiple releases with 0% adoption
- **test-quality-sentinel.md**: uses max-continuations: 40 (extremely high) - unique outlier
