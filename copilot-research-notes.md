# Copilot CLI Research Notes

## Analysis History

### 2026-04-19 (Run 24639070790)
- 196 total MD workflows, 87 using engine: copilot (explicitly)
- **Major improvement**: cache-memory adoption: 55→80 (+45%)
- **Improvement**: custom agent file use: 3→7 (+133%)
- **Improvement**: bare mode: 2→7 (+250%)
- **Persistent gaps**: version pinning (0%), api-target (0%), engine.args/env (0%), mcp-gateway (0%), blocked-domains (0%), observability (0%)
- 5/11 custom agent files still unused: grumpy-reviewer, w3c-specification-writer, create-safe-output-type, custom-engine-implementation, interactive-agent-designer

### 2026-04-17 (Run 24586698669)
- 85 copilot workflows tracked
- Notable: engine.args/env at 0%, mcp-gateway at 0%

### 2026-04-16 (Run 24534029243)
- 192 total, 90 explicit copilot + 26 default = 116 effective
- playwright regression: 20→12 (-40%)
- strict_mode: 111→126 (+13%)

## Persistent Opportunities (Not Addressed in 3+ Runs)

1. **engine.version**: Never used → stability risk
2. **engine.api-target**: Never used → GHEC/GHES teams can't use this
3. **token-weights**: Never used → no custom cost modeling
4. **block-domains**: Never used → missed defense-in-depth
5. **mcp-scripts**: 1 workflow → underutilized dynamic MCP capability
6. **max-continuations**: 2-3 only → Copilot-unique autopilot underused

## Recommendations Tracking

| Recommendation | Status | Date Added |
|---|---|---|
| Use engine.version pinning for reproducibility | ⏳ Pending | 2026-04-16 |
| Expand max-continuations for complex workflows | ⏳ Pending | 2026-04-16 |
| Use bare:true for simple/creative workflows | ⏳ Pending | 2026-04-16 |
| Add network.blocked for defense-in-depth | ⏳ Pending | 2026-04-17 |
| Activate unused agent files | ⏳ Pending | 2026-04-16 |
| Model override for cost optimization | ⏳ Pending | 2026-04-17 |
| Cache-memory adoption growing | ✅ Improving | 2026-04-16 |
| Custom agent adoption | ✅ Improving | 2026-04-17 |
