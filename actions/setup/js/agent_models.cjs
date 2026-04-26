// @ts-check
/// <reference types="@actions/github-script" />

"use strict";

/**
 * agent_models.cjs
 *
 * Queries the agentic engine's models endpoint before agent execution and stores
 * the results in /tmp/gh-aw/agents.json for inclusion in the agent artifact.
 *
 * The JSON file follows the structure:
 *   { "<engineId>": { "version": "<version>", "models": <models-data> } }
 *
 * Primary API — usable from any Node.js context (driver, github-script, standalone):
 *   queryModels({ endpoint, token, engineId, engineVersion, agentsJsonPath?, stepSummaryPath?, logFn? })
 *
 * github-script convenience wrapper (uses core.* globals):
 *   main()   — reads env vars and delegates to queryModels()
 *
 * Required environment variables when using main():
 *   - GH_AW_MODELS_ROUTE:       Route path for the models endpoint (e.g. "/models")
 *   - GITHUB_COPILOT_BASE_URL:  API base URL; assembled with GH_AW_MODELS_ROUTE to form the full URL
 *   - COPILOT_GITHUB_TOKEN:     Bearer token for authentication
 *   - GH_AW_ENGINE_ID:          Agentic engine identifier (e.g. "copilot")
 *   - GH_AW_ENGINE_VERSION:     Version string of the engine CLI
 */

const fs = require("fs");
const { getErrorMessage } = require("./error_helpers.cjs");

/** Default Copilot API base URL when GITHUB_COPILOT_BASE_URL is not configured. */
const DEFAULT_COPILOT_BASE_URL = "https://api.githubcopilot.com";

/** Path where model data is written so it is bundled in the agent artifact. */
const AGENTS_JSON_PATH = "/tmp/gh-aw/agents.json";

/** Request timeout in milliseconds for the models HTTP call. */
const REQUEST_TIMEOUT_MS = 15_000;

/**
 * Perform an HTTP GET request to the models endpoint and return the parsed JSON body.
 * Uses the Node.js built-in fetch API (available since Node 18, stable in Node 21+).
 *
 * @param {string} endpointUrl - Full URL of the models endpoint
 * @param {string} authToken   - Bearer token for the Authorization header
 * @returns {Promise<unknown>}  Parsed JSON response body
 */
async function fetchModels(endpointUrl, authToken) {
  const response = await fetch(endpointUrl, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${authToken}`,
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    signal: AbortSignal.timeout(REQUEST_TIMEOUT_MS),
  });

  if (!response.ok) {
    const body = await response.text();
    throw new Error(`HTTP ${response.status} from ${endpointUrl}: ${body.slice(0, 200)}`);
  }

  return response.json();
}

/**
 * Extract the flat models list from the API response, which may arrive in several shapes:
 *   - { models: [...] }  — most common
 *   - { data: [...] }    — OpenAI-compatible list
 *   - [...]              — bare array
 *
 * @param {unknown} responseBody
 * @returns {unknown[]}
 */
function extractModelsList(responseBody) {
  if (Array.isArray(responseBody)) {
    return responseBody;
  }
  if (responseBody && typeof responseBody === "object") {
    const obj = /** @type {Record<string, unknown>} */ responseBody;
    if (Array.isArray(obj["models"])) {
      return /** @type {unknown[]} */ obj["models"];
    }
    if (Array.isArray(obj["data"])) {
      return /** @type {unknown[]} */ obj["data"];
    }
  }
  return [];
}

/**
 * Build a markdown table from a flat models list for the step summary.
 *
 * @param {unknown[]} models
 * @returns {string} Markdown table or a plain message when the list is empty
 */
function buildModelsMarkdown(models) {
  if (models.length === 0) {
    return "No models returned by the endpoint.";
  }

  const rows = ["| ID | Display name | Vendor |", "| --- | --- | --- |"];
  for (const m of models) {
    if (!m || typeof m !== "object") continue;
    const entry = /** @type {Record<string, unknown>} */ m;
    const id = String(entry["id"] || "");
    const name = String(entry["display_name"] || entry["name"] || "");
    const vendor = String(entry["vendor"] || entry["owned_by"] || "");
    rows.push(`| ${id} | ${name} | ${vendor} |`);
  }
  return rows.join("\n");
}

/**
 * Log individual models using the provided logging function.
 *
 * @param {unknown[]} models
 * @param {string}    engineId
 * @param {(msg: string) => void} logFn
 */
function logModels(models, engineId, logFn) {
  logFn(`[${engineId}] Available models (${models.length}):`);
  for (const m of models) {
    if (!m || typeof m !== "object") continue;
    const entry = /** @type {Record<string, unknown>} */ m;
    const id = String(entry["id"] || "?");
    const name = String(entry["display_name"] || entry["name"] || "");
    logFn(`  - ${id}${name ? ": " + name : ""}`);
  }
}

/**
 * Query the models endpoint, persist results to agents.json, and optionally append
 * a summary section to the step-summary file.  Callable from any Node.js context
 * (driver harness, github-script, standalone) without depending on global `core.*`.
 *
 * @param {{
 *   endpoint: string,
 *   token: string,
 *   engineId: string,
 *   engineVersion: string,
 *   agentsJsonPath?: string,
 *   stepSummaryPath?: string | null,
 *   logFn?: (msg: string) => void,
 * }} options
 * @returns {Promise<void>}
 */
async function queryModels({ endpoint, token, engineId, engineVersion, agentsJsonPath = AGENTS_JSON_PATH, stepSummaryPath = null, logFn = () => {} }) {
  logFn(`querying models from: ${endpoint} (engine=${engineId} version=${engineVersion})`);

  let modelsData;
  try {
    modelsData = await fetchModels(endpoint, token);
  } catch (error) {
    logFn(`warning: failed to query models endpoint: ${getErrorMessage(error)}`);
    return;
  }

  const modelsList = extractModelsList(modelsData);
  logModels(modelsList, engineId, logFn);

  // Write agents.json so the data is bundled in the agent artifact
  const agentsInfo = {
    [engineId]: {
      version: engineVersion,
      models: modelsData,
    },
  };

  try {
    fs.mkdirSync("/tmp/gh-aw", { recursive: true });
    fs.writeFileSync(agentsJsonPath, JSON.stringify(agentsInfo, null, 2) + "\n");
    logFn(`wrote models info to ${agentsJsonPath}`);
  } catch (error) {
    logFn(`warning: failed to write ${agentsJsonPath}: ${getErrorMessage(error)}`);
  }

  // Append a collapsible section to the step summary file
  if (stepSummaryPath) {
    const markdown = buildModelsMarkdown(modelsList);
    const section = `\n<details>\n<summary>Available Models (${engineId} ${engineVersion})</summary>\n\n${markdown}\n</details>\n`;
    try {
      fs.appendFileSync(stepSummaryPath, section);
    } catch (error) {
      logFn(`warning: failed to write models step summary: ${getErrorMessage(error)}`);
    }
  }
}

/**
 * Main entry point — called by the compiler-generated github-script step.
 * Exits cleanly (non-fatal) when required env vars are absent.
 */
async function main() {
  const modelsRoute = process.env.GH_AW_MODELS_ROUTE;
  if (!modelsRoute) {
    core.info("GH_AW_MODELS_ROUTE is not set — skipping models query");
    return;
  }

  const authToken = process.env.COPILOT_GITHUB_TOKEN;
  if (!authToken) {
    core.info("COPILOT_GITHUB_TOKEN is not set — skipping models query");
    return;
  }

  const baseUrl = (process.env.GITHUB_COPILOT_BASE_URL || DEFAULT_COPILOT_BASE_URL).replace(/\/$/, "");
  const endpoint = baseUrl + modelsRoute;

  const engineId = process.env.GH_AW_ENGINE_ID || "copilot";
  const engineVersion = process.env.GH_AW_ENGINE_VERSION || "unknown";

  await queryModels({
    endpoint,
    token: authToken,
    engineId,
    engineVersion,
    stepSummaryPath: process.env.GITHUB_STEP_SUMMARY || null,
    logFn: msg => core.info(msg),
  });
}

module.exports = {
  main,
  queryModels,
  fetchModels,
  extractModelsList,
  buildModelsMarkdown,
  logModels,
  AGENTS_JSON_PATH,
  REQUEST_TIMEOUT_MS,
  DEFAULT_COPILOT_BASE_URL,
};

if (require.main === module) {
  main().catch(err => {
    // eslint-disable-next-line no-console
    console.error(`[agent_models] unexpected error: ${err.message}`);
    process.exit(1);
  });
}
