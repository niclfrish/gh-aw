---
title: CodeGraph Code Search
description: Configure a semantic code knowledge graph for AI-powered code search, dependency analysis, and architecture exploration in agentic workflows.
sidebar:
  order: 821
---

CodeGraph Code Search provides a semantically searchable knowledge graph over your codebase. It runs [Jakedismo/codegraph-rust](https://github.com/Jakedismo/codegraph-rust) as an MCP server so agents can search, reason about, and navigate code by meaning rather than by text pattern.

Unlike text search or basic semantic embedding tools, CodeGraph builds a **real knowledge graph**: AST nodes and edges enriched with relationships (calls, defines, uses, returns, mutates) that allow agents to traverse dependencies, trace call chains, and understand architecture — not just find matching strings.

The index is built in a dedicated indexing step (with a separate `contents: read` job or pre-agent steps) and shared with the agent via GitHub Actions cache, so the agent can re-use the same index across runs without re-indexing.

:::caution[Experimental]
CodeGraph Code Search is an experimental integration. The `codegraph-rust` project does not yet publish pre-built binaries, so the binary is compiled from source on first use (cached thereafter). The MCP API may change as the project evolves.
:::

## Prerequisites

CodeGraph's internal reasoning agents require an LLM API key. Configure at least one of the following repository secrets:

| Secret | Provider |
|--------|----------|
| `ANTHROPIC_API_KEY` | Claude (recommended) |
| `OPENAI_API_KEY` | GPT models |
| `JINA_API_KEY` | Jina AI embeddings (free tier available) |

You can also use `CODEGRAPH_ANTHROPIC_API_KEY` as a scoped alias for `ANTHROPIC_API_KEY` if you want to isolate the key used by CodeGraph from other workflows.

## Basic Configuration

```aw wrap
---
engine: copilot
permissions:
  contents: read
imports:
  - uses: shared/mcp/codegraph.md
    with:
      index-tier: fast
---
```

## Configuration Options

### `index-tier`

Controls the depth of analysis during indexing. Defaults to `fast`.

| Tier | What it enables | Typical use |
|------|-----------------|-------------|
| `fast` | AST nodes + core edges only | Quick CI runs, low storage |
| `balanced` | LSP symbols + docs/enrichment + module linking | Best agentic accuracy |
| `full` | All analyzers + LSP definitions + dataflow + architecture | Maximum richness |

```aw wrap
---
imports:
  - uses: shared/mcp/codegraph.md
    with:
      index-tier: balanced
---
```

The `balanced` tier requires language server tools to be available on the runner:

| Language | Requires |
|----------|---------|
| Rust | `rust-analyzer` |
| TypeScript/JavaScript | `node` + `typescript-language-server` |
| Python | `node` + `pyright-langserver` |
| Go | `gopls` |

### `cache-key`

A GitHub Actions cache key used to persist the SurrealDB index across workflow runs. When set and a cache hit occurs, the indexing step is skipped (read-only mode), making the agent job significantly faster.

```aw wrap
---
imports:
  - uses: shared/mcp/codegraph.md
    with:
      cache-key: "codegraph-${{ github.repository }}-${{ hashFiles('src/**', '*.rs') }}"
---
```

When the cache key does not match (cache miss), CodeGraph re-indexes and saves the new index under the given key. Use content-based hashing in the key to invalidate the cache when source files change.

**Read-only mode** (restore without re-indexing): set `cache-key` to a fixed key that was previously saved:

```aw wrap
---
imports:
  - uses: shared/mcp/codegraph.md
    with:
      cache-key: "codegraph-my-project-stable"
---
```

If the cache key does not exist at all, indexing runs normally and the result is saved.

## Agentic Tools

CodeGraph exposes four consolidated agentic tools. Each tool runs an internal reasoning agent that plans, searches the graph, and synthesizes an answer — not just a list of files.

| Tool | Focus values | Best for |
|------|-------------|----------|
| `agentic_context` | `search`, `builder`, `question` | Finding code by concept; pre-implementation context |
| `agentic_impact` | `dependencies`, `call_chain` | Impact analysis before refactoring; tracing callers |
| `agentic_architecture` | `structure`, `api_surface` | Architecture overview; public API enumeration |
| `agentic_quality` | `complexity`, `coupling`, `hotspots` | Risk assessment; refactoring targets |

Each tool accepts an optional `focus` parameter. Without it, the tool auto-selects the most appropriate reasoning strategy based on the query.

## Example: Code Search on Pull Request

```aw wrap
---
on:
  pull_request:
engine: copilot
permissions:
  contents: read
  pull-requests: write
imports:
  - uses: shared/mcp/codegraph.md
    with:
      index-tier: fast
      cache-key: "codegraph-${{ github.repository }}-${{ hashFiles('**/*.go', '**/*.ts', '**/*.rs') }}"
safe-outputs:
  add-comment:
    hide-older-comments: true
---

Call `read_initial_instructions` from the codegraph MCP server, then analyze the pull request diff
to answer: which functions and modules are affected by these changes?

Use `agentic_impact` with `focus: "dependencies"` for each changed file to understand
the blast radius, then use `agentic_context` with `focus: "search"` to find any tests
or related patterns that should be updated.

Post a concise comment summarizing the impact.
```

## Example: Daily Architecture Analysis

```aw wrap
---
on:
  schedule: weekly on monday around 09:00
  workflow_dispatch:
engine: copilot
permissions:
  contents: read
  issues: write
imports:
  - uses: shared/mcp/codegraph.md
    with:
      index-tier: balanced
      cache-key: "codegraph-${{ github.repository }}-arch"
safe-outputs:
  create-issue:
    title-prefix: "[arch-analysis] "
    labels: [architecture, automated-analysis]
    max: 1
    close-older-issues: true
---

Call `read_initial_instructions`, then perform a weekly architecture health review:

1. Use `agentic_quality` with `focus: "hotspots"` to identify the top 3 complexity hotspots
2. Use `agentic_quality` with `focus: "coupling"` to find tightly coupled modules
3. Use `agentic_architecture` to describe the overall layer structure

Create a GitHub issue summarizing findings and actionable improvement suggestions.
```

## Example: Separate Indexing Job

For large codebases or workflows where the agent job should not need `contents: read`,
use a custom indexing job to build the index and pass it as an artifact:

```aw wrap
---
on:
  push:
    branches: [main]
  workflow_dispatch:
engine: copilot
permissions:
  contents: read
  issues: write
jobs:
  codegraph-index:
    runs-on: ubuntu-latest
    needs: [activation]
    permissions:
      contents: read
    steps:
      - name: Checkout repository
        uses: actions/checkout@v6.0.2
        with:
          persist-credentials: false

      - name: Restore codegraph binary from cache
        id: binary-cache
        uses: actions/cache/restore@v4
        with:
          path: ~/.cargo/bin/codegraph
          key: codegraph-bin-${{ runner.os }}-${{ runner.arch }}

      - name: Install Rust and build codegraph
        if: steps.binary-cache.outputs.cache-hit != 'true'
        run: |
          curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
          source "$HOME/.cargo/env"
          cargo install \
            --git https://github.com/Jakedismo/codegraph-rust \
            --bin codegraph --all-features --quiet

      - name: Save codegraph binary to cache
        if: steps.binary-cache.outputs.cache-hit != 'true'
        uses: actions/cache/save@v4
        with:
          path: ~/.cargo/bin/codegraph
          key: codegraph-bin-${{ runner.os }}-${{ runner.arch }}

      - name: Install SurrealDB
        run: |
          curl -sSf https://install.surrealdb.com | sh
          echo "$HOME/.surrealdb" >> "$GITHUB_PATH"

      - name: Start SurrealDB and index repository
        run: |
          mkdir -p /tmp/codegraph-data /tmp/gh-aw/codegraph/logs
          surreal start \
            --bind 0.0.0.0:3004 \
            --user root --pass root \
            file:///tmp/codegraph-data/surreal.db \
            > /tmp/gh-aw/codegraph/logs/surrealdb.log 2>&1 &
          for i in $(seq 1 30); do
            surreal is-ready --endpoint http://localhost:3004 2>/dev/null && break
            sleep 2
          done
          curl -fsSL https://raw.githubusercontent.com/Jakedismo/codegraph-rust/main/schema/codegraph.surql \
            | surreal sql \
              --endpoint ws://localhost:3004 \
              --namespace ouroboros --database codegraph \
              --username root --password root
          codegraph index --path "$GITHUB_WORKSPACE" --index-tier fast \
            2>&1 | tee /tmp/gh-aw/codegraph/logs/index.log
        env:
          SURREAL_URL: ws://localhost:3004
          SURREAL_USER: root
          SURREAL_PASS: root
          SURREAL_NAMESPACE: ouroboros
          SURREAL_DATABASE: codegraph
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}

      - name: Upload codegraph index artifact
        uses: actions/upload-artifact@v7.0.1
        with:
          name: codegraph-index
          path: /tmp/codegraph-data
          retention-days: 1

steps:
  - name: Download codegraph index artifact
    uses: actions/download-artifact@v8.0.1
    with:
      name: codegraph-index
      path: /tmp/codegraph-data

  - name: Install SurrealDB and start with indexed data
    run: |
      curl -sSf https://install.surrealdb.com | sh
      echo "$HOME/.surrealdb" >> "$GITHUB_PATH"
      mkdir -p /tmp/gh-aw/codegraph/logs
      surreal start \
        --bind 0.0.0.0:3004 \
        --user root --pass root \
        file:///tmp/codegraph-data/surreal.db \
        > /tmp/gh-aw/codegraph/logs/surrealdb.log 2>&1 &
      for i in $(seq 1 30); do
        surreal is-ready --endpoint http://localhost:3004 2>/dev/null && break
        sleep 2
      done

  - name: Restore codegraph binary from cache
    uses: actions/cache/restore@v4
    with:
      path: ~/.cargo/bin/codegraph
      key: codegraph-bin-${{ runner.os }}-${{ runner.arch }}

mcp-servers:
  codegraph:
    command: "codegraph"
    args: ["start", "stdio"]
    env:
      SURREAL_URL: ws://localhost:3004
      SURREAL_USER: root
      SURREAL_PASS: root
      SURREAL_NAMESPACE: ouroboros
      SURREAL_DATABASE: codegraph
      ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
      CODEGRAPH_ARCH_BOOTSTRAP: "true"
---

Call `read_initial_instructions` first, then perform your code analysis task.
```

In the separate-job pattern the indexing job runs with `contents: read` while the
agent job can operate without repository access, downloading only the pre-built index artifact.

## Telemetry

Use `shared/observability-otlp.md` to record CodeGraph index and query metrics alongside the workflow's distributed trace:

```yaml wrap title=".github/workflows/shared/codegraph-otlp.md"
---
# Shared import: emit codegraph index stats after the agent job.

steps:
  - name: Record codegraph telemetry
    id: codegraph-otlp
    uses: actions/github-script@v8
    with:
      script: |
        const fs   = require('fs');
        const otlp = require('/tmp/gh-aw/actions/otlp.cjs');

        // codegraph writes index stats to /tmp/gh-aw/codegraph/logs/index.log
        let nodesIndexed = 0;
        try {
          const log = fs.readFileSync('/tmp/gh-aw/codegraph/logs/index.log', 'utf8');
          const match = log.match(/nodes[:\s]+(\d+)/i);
          if (match) nodesIndexed = parseInt(match[1], 10);
        } catch { /* index not available */ }

        await otlp.logSpan('codegraph', {
          'codegraph.nodes.indexed': nodesIndexed,
        });
---
```

## Related Documentation

- [Tools](/gh-aw/reference/tools/) - Overview of all available tools and configuration
- [Imports](/gh-aw/reference/imports/) - Importing shared workflow components
- [QMD Documentation Search](/gh-aw/reference/qmd/) - Vector search over documentation files
- [Cache Memory](/gh-aw/reference/cache-memory/) - Persistent memory across workflow runs
- [GitHub Tools](/gh-aw/reference/github-tools/) - GitHub API operations
- [Custom OTLP Attributes](/gh-aw/guides/custom-otlp-attributes/) - Emit telemetry from shared imports
- [Serena](/gh-aw/guides/serena/) - LSP-based semantic code analysis (alternative/complement)
