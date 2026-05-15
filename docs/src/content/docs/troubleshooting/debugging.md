---
title: Debugging Workflows
description: How to run, debug, and investigate agentic workflow failures using the Copilot CLI, gh aw audit, and log analysis.
sidebar:
  order: 250
---

This guide shows how to debug agentic workflow failures on **github.com** using the Copilot CLI, `gh aw` commands, and manual investigation.

## Debugging with an AI Agent

The fastest path to a fix is to let an agent investigate. Launch `copilot`, run `/agent` and select **agentic-workflows**, then paste the failing run URL:

```text
Debug this workflow run: https://github.com/OWNER/REPO/actions/runs/RUN_ID
```

The agent downloads logs, identifies the root cause (missing tools, permission errors, network blocks), and suggests a fix or opens a PR. Follow up with targeted questions like *"What domains were blocked?"* or *"Why did the MCP server fail?"*.

**Alternatives**: on GitHub.com with [agentic authoring](/gh-aw/guides/agentic-authoring/) configured, use `/agent agentic-workflows debug <run-url>` in Copilot Chat. For any other coding agent, point it at the standalone prompt at `https://raw.githubusercontent.com/github/gh-aw/main/debug.md` along with the run URL.

## Debugging with CLI Commands

### Auditing a Specific Run

`gh aw audit` produces a full breakdown of a run — failure analysis, behavior fingerprint, tool usage, MCP status, firewall analysis, token/cost metrics, and safe-outputs. Pass a run ID or any URL form (run, job, or step), and add `--parse` for shareable markdown:

```bash
gh aw audit 12345678
gh aw audit https://github.com/OWNER/REPO/actions/runs/123/job/456#step:7:1
```

Pass multiple IDs to compare runs and detect regressions; use `gh aw logs --format markdown` for trends across many runs:

```bash
gh aw audit 12345678 12345679 --format markdown
gh aw logs my-workflow --format markdown --count 10
```

See [Audit Commands](/gh-aw/reference/audit/) for complete flag documentation.

### Analyzing Workflow Logs

`gh aw logs` downloads and analyzes logs across runs (tool usage, network patterns, errors, warnings). Common flags: `-c` count, `--start-date`, `--firewall`, `--safe-output`, `--json`.

```bash
gh aw logs my-workflow -c 10 --start-date -1w --firewall
```

Results are cached locally for 10–100× speedup on subsequent runs.

### Checking Workflow Health

`gh aw health` gives a quick overview of workflow status across all workflows in a repository:

```bash
gh aw health
```

### Inspecting MCP Configuration

When MCP servers are suspect: list workflows with MCP servers (`gh aw mcp list`), inspect one (`gh aw mcp inspect my-workflow`), or open the web inspector (`gh aw mcp inspect my-workflow --inspector`).

## Common Errors

### "Authentication failed"

The Copilot token is missing, expired, or lacks required permissions. Verify an active Copilot subscription, ensure the token has **Copilot Requests** permission (fine-grained PATs), and check validity with `gh auth status`. See [Authentication Reference](/gh-aw/reference/auth/).

### "Tool not found" or Missing Tool Calls

A referenced tool isn't configured, or the MCP server failed to connect. Run `gh aw mcp inspect my-workflow` to verify tool configuration and `gh aw audit <run-id>` to see which tools were requested vs. available.

### Network / Firewall Blocks

A line like `DENIED CONNECT registry.npmjs.org:443` means the agent reached a domain not in `network.allowed`. Add the domain, or an ecosystem shorthand (`node`, `python`, …) that bundles common registries:

```aw
network:
  allowed:
    - defaults
    - node
    - registry.npmjs.org
```

See [Network Configuration](/gh-aw/guides/network-configuration/) for common domain configurations.

### Safe-Outputs Not Creating Issues / Comments

The safe-outputs job failed, the agent didn't produce expected output, or permissions are missing. Run `gh aw audit <run-id>` and inspect its safe-outputs section. See [Safe Outputs Reference](/gh-aw/reference/safe-outputs/).

### Compilation Errors

Frontmatter has schema validation errors or unsupported fields. Run `gh aw compile my-workflow --verbose` for details, `gh aw fix --write` to auto-correct, or `gh aw compile --validate` to validate without compiling. See [Error Reference](/gh-aw/troubleshooting/errors/).

## Advanced Debugging

### Enable Debug Logging

Set `DEBUG` to scope verbose logging for any `gh aw` command — use `*` for everything or a comma-separated package list like `workflow:*,cli:*`. Output goes to `stderr` (capture with `2>&1 | tee debug.log`).

```bash
DEBUG=workflow:*,cli:* gh aw compile my-workflow
```

### Enable GitHub Actions Debug Logging

Add a repository secret `ACTIONS_STEP_DEBUG=true` (Settings → Secrets and variables → Actions) and re-run the workflow for verbose step-level logs in the Actions UI.

### Inspecting Firewall Logs

Run artifacts include `sandbox/firewall/logs/access.log`. Each line shows whether a domain was allowed (`TCP_TUNNEL/200 api.github.com:443`) or blocked (`DENIED CONNECT blocked-domain.com:443`). For the same data via CLI, use `gh aw logs my-workflow --firewall` or `gh aw audit <run-id>` (which includes firewall analysis).

### Inspecting Artifacts

Workflow runs produce several artifacts useful for debugging:

| Artifact | Location | Contents |
|----------|----------|----------|
| `prompt.txt` | `/tmp/gh-aw/aw-prompts/` | The full prompt sent to the AI agent |
| `agent_output.json` | `/tmp/gh-aw/safeoutputs/` | Structured safe-output data |
| `agent-stdio.log` | `/tmp/gh-aw/` | Raw agent stdin/stdout log |
| `firewall-logs/` | `/tmp/gh-aw/firewall-logs/` | Network access logs |

Download artifacts from the GitHub Actions run page or via the CLI:

```bash
gh run download <run-id> --repo OWNER/REPO
```

### Recompiling for a Quick Fix

After editing the `.md` file, run `gh aw compile my-workflow`, then commit both the `.md` and the regenerated `.lock.yml`.
