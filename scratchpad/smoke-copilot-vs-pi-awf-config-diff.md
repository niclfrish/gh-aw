# AWF configuration diff: smoke-copilot vs smoke-pi (`*.lock.yml`)

This note compares:

- `/home/runner/work/gh-aw/gh-aw/.github/workflows/smoke-copilot.lock.yml`
- `/home/runner/work/gh-aw/gh-aw/.github/workflows/smoke-pi.lock.yml`

## Top-level workflow env

No differences were found in top-level `env:`. Both workflows set the same OTEL/GH_AW telemetry variables.

## AWF execution env differences (agent job)

Copilot-specific env present in `smoke-copilot.lock.yml` but not in `smoke-pi.lock.yml`:

- `GH_AW_MCP_CONFIG` (`smoke-copilot.lock.yml:1801`)
- `GITHUB_MCP_SERVER_TOKEN` (`smoke-copilot.lock.yml:1811`)
- `COPILOT_AGENT_RUNNER_TYPE` (`smoke-copilot.lock.yml:1797`)
- `GITHUB_COPILOT_INTEGRATION_ID` (`smoke-copilot.lock.yml:1809`)
- Also present only for copilot execute path:
  - `AWF_REFLECT_ENABLED`
  - `COPILOT_API_KEY`
  - `COPILOT_MODEL`
  - `GITHUB_API_URL`
  - `GITHUB_HEAD_REF`
  - `GITHUB_REF_NAME`
  - `GITHUB_SERVER_URL`
  - `XDG_CONFIG_HOME`

Pi-specific env present in `smoke-pi.lock.yml` but not in `smoke-copilot.lock.yml`:

- `GH_AW_PI_MODEL` (`smoke-pi.lock.yml:848`)
- `PI_CODING_AGENT_DIR` (`smoke-pi.lock.yml:860`)

## AWF command/CLI argument differences

`awf` invocation in copilot includes an extra artifact mount not present in pi:

- Copilot has `--mount "${RUNNER_TEMP}/gh-aw/safeoutputs/upload-artifacts:...:rw"` (`smoke-copilot.lock.yml:1793`)
- Pi does not include this mount (`smoke-pi.lock.yml:843`)

Inner engine command executed via `awf` differs significantly:

- Copilot runs `copilot_harness.cjs` + `/usr/local/bin/copilot` with copilot flags:
  - `--disable-builtin-mcps`
  - `--autopilot`
  - `--max-autopilot-continues 2`
  - `--allow-all-tools`
  - `--allow-all-paths`
  - `--no-custom-instructions`
  - (`smoke-copilot.lock.yml:1794`)
- Pi runs `pi --print --mode json --no-session --model ... --extension ...`:
  - no copilot-equivalent autopilot/allow-all-paths/no-custom-instructions flags
  - (`smoke-pi.lock.yml:844`)

## Guard/integrity policy differences around AWF execution

- Copilot path configures `GH_AW_MIN_INTEGRITY: approved` (`smoke-copilot.lock.yml:619`) and `CLI_PROXY_POLICY` with `"min-integrity":"approved"` (`smoke-copilot.lock.yml:1772`)
- Pi path configures `GH_AW_MIN_INTEGRITY: none` (`smoke-pi.lock.yml:480`) and `CLI_PROXY_POLICY` with `"min-integrity":"none"` (`smoke-pi.lock.yml:826`)

## Conclusion (focused on env/CLI mismatch)

The strongest configuration mismatches tied to AWF behavior are:

1. Pi execute env lacks copilot-path variables used for MCP wiring/session integration (`GH_AW_MCP_CONFIG`, `GITHUB_MCP_SERVER_TOKEN`, integration vars).
2. Pi executes a different engine command (`pi ...`) with a different CLI contract than copilot harness flags.
3. Pi runs under a less strict integrity policy (`none` vs `approved`), which changes guard behavior around tool access/proxying.

These are the primary environment/CLI differences that explain why AWF can be correctly wired for copilot while pi behaves differently.
