---
mcp-servers:
  agentdb:
    command: "npx"
    args: ["agentdb@alpha", "mcp", "start"]
    env:
      AGENTDB_PATH: "/tmp/gh-aw/agentdb/discussions.rvf"
    allowed: ["*"]
---

<!--
## AgentDB MCP Server

Shared MCP configuration for AgentDB vector memory/search.

- Docs: https://github.com/ruvnet/agentdb/blob/main/docs/README-full.md
- Launch command from docs: `npx agentdb mcp start`
- Default store path in this workflow runtime: `/tmp/gh-aw/agentdb/discussions.rvf`

Usage:

```yaml
imports:
  - shared/mcp/agentdb.md
```
-->
