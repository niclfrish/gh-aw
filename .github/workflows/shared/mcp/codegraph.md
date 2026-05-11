---
# CodeGraph MCP Server - Semantic Code Knowledge Graph
# Transforms your codebase into a semantically searchable knowledge graph
# that AI agents can reason about using graph traversal + vector embeddings.
#
# Documentation: https://github.com/Jakedismo/codegraph-rust
#
# Prerequisites (secrets):
#   ANTHROPIC_API_KEY (or OPENAI_API_KEY / JINA_API_KEY) for the internal
#   reasoning agents that power CodeGraph's agentic tools.
#
# Usage:
#   imports:
#     - uses: shared/mcp/codegraph.md
#       with:
#         index-tier: fast          # fast | balanced | full (default: fast)
#         cache-key: "codegraph-${{ github.repository }}"   # optional; enables caching
#
# The shared workflow:
#   1. Installs SurrealDB (graph + vector database)
#   2. Builds/restores the codegraph binary (Rust; cached after first build)
#   3. Applies the CodeGraph schema to SurrealDB
#   4. Restores the code index from cache (if cache-key is set and cache exists)
#   5. Indexes the workspace (skipped on cache hit)
#   6. Saves the index to cache (if cache-key is set and cache was missed)
#   7. Starts the codegraph MCP server in stdio mode

import-schema:
  index-tier:
    type: string
    enum: [fast, balanced, full]
    default: fast
    description: >
      Indexing tier controlling speed vs. graph richness.
      fast: AST nodes + core edges only (no LSP or enrichment). Fastest, lowest storage.
      balanced: LSP symbols + docs/enrichment + module linking. Best for agentic workflows.
      full: All analyzers + LSP definitions + dataflow + architecture. Maximum accuracy.
  cache-key:
    type: string
    default: ""
    description: >
      GitHub Actions cache key for the SurrealDB index. When set, the index is cached
      across workflow runs. On a cache hit the indexing step is skipped (read-only mode).
      Use a key that includes the content hash to invalidate on source changes, e.g.:
      "codegraph-${{ github.repository }}-${{ hashFiles('**/*.rs', '**/*.go', '**/*.ts') }}"

steps:
  - name: Restore codegraph binary from cache
    id: binary-cache
    uses: actions/cache/restore@v4
    with:
      path: ~/.cargo/bin/codegraph
      key: codegraph-bin-${{ runner.os }}-${{ runner.arch }}

  - name: Install Rust toolchain
    if: steps.binary-cache.outputs.cache-hit != 'true'
    uses: dtolnay/rust-toolchain@stable

  - name: Build and install codegraph from source
    if: steps.binary-cache.outputs.cache-hit != 'true'
    run: |
      set -euo pipefail
      echo "Building codegraph from source (this takes a few minutes on first run)..."
      cargo install \
        --git https://github.com/Jakedismo/codegraph-rust \
        --bin codegraph \
        --all-features \
        --quiet
      echo "codegraph installed: $(codegraph --version 2>/dev/null || echo 'version unknown')"

  - name: Save codegraph binary to cache
    if: steps.binary-cache.outputs.cache-hit != 'true'
    uses: actions/cache/save@v4
    with:
      path: ~/.cargo/bin/codegraph
      key: codegraph-bin-${{ runner.os }}-${{ runner.arch }}

  - name: Install SurrealDB
    run: |
      set -euo pipefail
      if ! command -v surreal &>/dev/null; then
        curl -sSf https://install.surrealdb.com | sh
        echo "$HOME/.surrealdb" >> "$GITHUB_PATH"
        export PATH="$HOME/.surrealdb:$PATH"
      fi
      echo "SurrealDB version: $(surreal version 2>/dev/null || surreal --version 2>/dev/null)"

  - name: Restore codegraph index from cache
    id: index-cache
    if: "${{ github.aw.import-inputs['cache-key'] != '' }}"
    uses: actions/cache/restore@v4
    with:
      path: /tmp/codegraph-data
      key: "${{ github.aw.import-inputs['cache-key'] }}"
      restore-keys: "codegraph-index-${{ github.repository }}-"

  - name: Start SurrealDB (in-memory or file-backed)
    run: |
      set -euo pipefail
      mkdir -p /tmp/codegraph-data /tmp/gh-aw/codegraph/logs

      SURREAL_BIN="${HOME}/.surrealdb/surreal"
      if ! command -v "$SURREAL_BIN" &>/dev/null; then
        SURREAL_BIN="surreal"
      fi

      # Use file-backed mode when caching is enabled for persistence between steps
      if [ -n "$CODEGRAPH_CACHE_KEY" ]; then
        DB_BACKEND="file:///tmp/codegraph-data/surreal.db"
      else
        DB_BACKEND="memory"
      fi

      "$SURREAL_BIN" start \
        --bind "0.0.0.0:3004" \
        --user root \
        --pass root \
        "$DB_BACKEND" \
        > /tmp/gh-aw/codegraph/logs/surrealdb.log 2>&1 &

      SURREAL_PID=$!
      echo "$SURREAL_PID" > /tmp/gh-aw/codegraph/surrealdb.pid

      # Wait for SurrealDB to become ready
      for i in $(seq 1 30); do
        if "$SURREAL_BIN" is-ready --endpoint "http://localhost:3004" 2>/dev/null; then
          echo "SurrealDB is ready (PID: $SURREAL_PID)"
          break
        fi
        if [ "$i" -eq 30 ]; then
          echo "SurrealDB failed to start. Logs:"
          cat /tmp/gh-aw/codegraph/logs/surrealdb.log
          exit 1
        fi
        echo "Waiting for SurrealDB... ($i/30)"
        sleep 2
      done
    env:
      CODEGRAPH_CACHE_KEY: "${{ github.aw.import-inputs['cache-key'] }}"

  - name: Apply CodeGraph schema
    run: |
      set -euo pipefail
      SURREAL_BIN="${HOME}/.surrealdb/surreal"
      if ! command -v "$SURREAL_BIN" &>/dev/null; then
        SURREAL_BIN="surreal"
      fi

      echo "Applying CodeGraph schema..."
      curl -fsSL \
        https://raw.githubusercontent.com/Jakedismo/codegraph-rust/main/schema/codegraph.surql \
        | "$SURREAL_BIN" sql \
          --endpoint "ws://localhost:3004" \
          --namespace ouroboros \
          --database codegraph \
          --username root \
          --password root
      echo "Schema applied successfully"

  - name: Index repository with CodeGraph
    if: steps.index-cache.outputs.cache-hit != 'true'
    run: |
      set -euo pipefail
      TIER="${CODEGRAPH_INDEX_TIER:-fast}"
      echo "Indexing ${GITHUB_WORKSPACE} with tier: ${TIER}"

      codegraph index \
        --path "${GITHUB_WORKSPACE}" \
        --index-tier "${TIER}" \
        2>&1 | tee /tmp/gh-aw/codegraph/logs/index.log

      echo "Indexing complete. Log tail:"
      tail -5 /tmp/gh-aw/codegraph/logs/index.log
    env:
      SURREAL_URL: "ws://localhost:3004"
      SURREAL_USER: root
      SURREAL_PASS: root
      SURREAL_NAMESPACE: ouroboros
      SURREAL_DATABASE: codegraph
      CODEGRAPH_INDEX_TIER: "${{ github.aw.import-inputs['index-tier'] }}"
      ANTHROPIC_API_KEY: "${{ secrets.CODEGRAPH_ANTHROPIC_API_KEY || secrets.ANTHROPIC_API_KEY }}"
      OPENAI_API_KEY: "${{ secrets.OPENAI_API_KEY }}"
      JINA_API_KEY: "${{ secrets.JINA_API_KEY }}"

  - name: Save codegraph index to cache
    if: |
      steps.index-cache.outputs.cache-hit != 'true' &&
      github.aw.import-inputs['cache-key'] != ''
    uses: actions/cache/save@v4
    with:
      path: /tmp/codegraph-data
      key: "${{ github.aw.import-inputs['cache-key'] }}"

mcp-servers:
  codegraph:
    command: "codegraph"
    args: ["start", "stdio"]
    env:
      SURREAL_URL: "ws://localhost:3004"
      SURREAL_USER: "root"
      SURREAL_PASS: "root"
      SURREAL_NAMESPACE: "ouroboros"
      SURREAL_DATABASE: "codegraph"
      ANTHROPIC_API_KEY: "${{ secrets.CODEGRAPH_ANTHROPIC_API_KEY || secrets.ANTHROPIC_API_KEY }}"
      OPENAI_API_KEY: "${{ secrets.OPENAI_API_KEY }}"
      JINA_API_KEY: "${{ secrets.JINA_API_KEY }}"
      CODEGRAPH_ARCH_BOOTSTRAP: "true"
      CODEGRAPH_INDEX_TIER: "${{ github.aw.import-inputs['index-tier'] }}"
---

## CodeGraph Semantic Code Analysis

The CodeGraph MCP server is connected. The `${{ github.repository }}` codebase has been
indexed using the **${{ github.aw.import-inputs.index-tier }}** tier.

- **Workspace**: `${{ github.workspace }}`
- **Index tier**: `${{ github.aw.import-inputs.index-tier }}`
- **Caching**: ${{ github.aw.import-inputs['cache-key'] != '' && 'enabled' || 'disabled (single-run)' }}
- **Logs**: `/tmp/gh-aw/codegraph/logs/`

### Initialize Session

Load CodeGraph's context and guidance before your first tool call:

```
Tool: read_initial_instructions
```

### Available Agentic Tools

| Tool | Best For |
|------|----------|
| `agentic_context` | Finding code, exploring patterns, building pre-implementation context |
| `agentic_impact` | Dependency analysis, impact mapping before refactoring |
| `agentic_architecture` | Big-picture structure, API surface analysis |
| `agentic_quality` | Complexity hotspots, coupling metrics, refactoring priorities |

### Analysis Strategy

1. **Always call `read_initial_instructions` first** — loads per-tool guidance
2. **Use `agentic_context` with `focus: "search"`** for concept-based code search
3. **Use `agentic_impact` with `focus: "dependencies"`** to trace what a change would break
4. **Use `agentic_context` with `focus: "builder"`** to gather pre-implementation context
5. **Use `agentic_quality`** to identify where complexity accumulates before refactoring
