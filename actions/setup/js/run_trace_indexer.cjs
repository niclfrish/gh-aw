// @ts-check
/// <reference types="@actions/github-script" />

const fs = require("node:fs/promises");
const path = require("node:path");

const { getErrorMessage, isRateLimitError } = require("./error_helpers.cjs");
const { resolveExecutionOwnerRepo } = require("./repo_helpers.cjs");
const { sanitizeContent } = require("./sanitize_content.cjs");

const REPORT_COUNT = 1000;
const DEFAULT_TRACE_OUTPUT_DIR = "./.cache/gh-aw/agentic-workflow-logs";
const REPORT_SECTION_DIR = "activity-report";

/** @typedef {{ key: string, heading: string, startDate: string }} TraceRange */

/** @type {TraceRange[]} */
const TRACE_RANGES = [
  { key: "24h", heading: "Last 24 hours", startDate: "-1d" },
  { key: "7d", heading: "Last 7 days", startDate: "-1w" },
];

/**
 * @param {string} text
 * @returns {boolean}
 */
function hasRateLimitText(text) {
  return /\bapi rate limit\b|\brate limit exceeded\b|\bsecondary rate limit\b|\b429\b/i.test(text);
}

/**
 * @param {string} filePath
 * @param {string} content
 * @returns {Promise<void>}
 */
async function writeReportSection(filePath, content) {
  await fs.mkdir(path.dirname(filePath), { recursive: true });
  await fs.writeFile(filePath, `${content.trim()}\n`, "utf8");
}

/**
 * @param {string} bin
 * @param {string[]} prefixArgs
 * @param {string} repoSlug
 * @param {string} outputDir
 * @param {TraceRange} range
 * @returns {Promise<boolean>}
 */
async function runTraceRange(bin, prefixArgs, repoSlug, outputDir, range) {
  const rangeReportPath = path.join(outputDir, REPORT_SECTION_DIR, `${range.key}.md`);
  const args = [...prefixArgs, "logs", "--repo", repoSlug, "--start-date", range.startDate, "--count", String(REPORT_COUNT), "--output", outputDir, "--format", "markdown"];
  core.info(`Running trace indexer: ${bin} ${args.join(" ")}`);

  try {
    const result = await exec.getExecOutput(bin, args, { ignoreReturnCode: true });
    const stdout = (result.stdout || "").trim();
    const stderr = (result.stderr || "").trim();
    const output = `${stdout}\n${stderr}`.trim();
    const rateLimited = hasRateLimitText(output);

    if (result.exitCode === 0 && stdout) {
      await writeReportSection(rangeReportPath, sanitizeContent(stdout));
      return true;
    }

    if (rateLimited) {
      await writeReportSection(rangeReportPath, "_Could not refresh this range due to GitHub API rate limiting._");
      return false;
    }

    await writeReportSection(rangeReportPath, `_Trace indexing failed (exit code ${result.exitCode})._\n\n\`\`\`\n${sanitizeContent(output || "No command output was captured.")}\n\`\`\``);
    return false;
  } catch (error) {
    const errorMessage = getErrorMessage(error);
    const rateLimited = isRateLimitError(error) || hasRateLimitText(errorMessage);
    if (rateLimited) {
      await writeReportSection(rangeReportPath, "_Could not refresh this range due to GitHub API rate limiting._");
      return false;
    }
    await writeReportSection(rangeReportPath, `_Trace indexing failed: ${sanitizeContent(errorMessage)}_`);
    return false;
  }
}

/**
 * Refresh cached logs and report sections for activity reporting.
 * @returns {Promise<void>}
 */
async function main() {
  const cmdPrefixStr = process.env.GH_AW_CMD_PREFIX || "gh aw";
  const traceOutputDir = process.env.GH_AW_TRACE_INDEX_OUTPUT_DIR || DEFAULT_TRACE_OUTPUT_DIR;
  const [bin, ...prefixArgs] = cmdPrefixStr.split(" ").filter(Boolean);
  const { owner, repo } = resolveExecutionOwnerRepo();
  const repoSlug = `${owner}/${repo}`;

  core.info(`Refreshing agentic workflow logs cache for ${repoSlug}`);

  let allRangesSucceeded = true;
  for (const range of TRACE_RANGES) {
    const ok = await runTraceRange(bin, prefixArgs, repoSlug, traceOutputDir, range);
    if (!ok) {
      allRangesSucceeded = false;
    }
  }

  if (!allRangesSucceeded) {
    throw new Error("Trace indexing completed with one or more range failures");
  }
}

module.exports = { main, hasRateLimitText, runTraceRange };
