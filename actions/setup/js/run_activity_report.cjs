// @ts-check
/// <reference types="@actions/github-script" />

const fs = require("node:fs/promises");
const path = require("node:path");
const { resolveExecutionOwnerRepo } = require("./repo_helpers.cjs");
const { sanitizeContent } = require("./sanitize_content.cjs");

const ISSUE_TITLE = "[aw] agentic status report";
const HEADING_DEMOTION_LEVELS = 2;
const DEFAULT_REPORT_OUTPUT_DIR = "./.cache/gh-aw/agentic-workflow-logs";
const REPORT_SECTION_DIR = "activity-report";

/** @typedef {{ key: string, heading: string }} ActivityRange */

/** @type {ActivityRange[]} */
const REPORT_RANGES = [
  { key: "24h", heading: "Last 24 hours" },
  { key: "7d", heading: "Last 7 days" },
];

/**
 * Read pre-indexed report markdown from the cache directory.
 *
 * @param {ActivityRange} range
 * @param {string} outputDir
 * @returns {Promise<{ heading: string, body: string }>}
 */
async function readCachedRangeReport(range, outputDir) {
  const rangeReportPath = path.join(outputDir, REPORT_SECTION_DIR, `${range.key}.md`);
  try {
    const markdown = await fs.readFile(rangeReportPath, "utf8");
    return {
      heading: range.heading,
      body: normalizeReportMarkdown(sanitizeContent(markdown.trim())),
    };
  } catch {
    core.warning(`Missing cached report for ${range.heading}: ${rangeReportPath}`);
    return {
      heading: range.heading,
      body: "_No cached trace index is available for this range yet._",
    };
  }
}

/**
 * Normalize report markdown for issue rendering.
 * Demotes headings so top-level report headings start at H3.
 *
 * @param {string} markdown
 * @returns {string}
 */
function normalizeReportMarkdown(markdown) {
  return markdown.replace(/^(#{1,6})\s+/gm, (_, hashes) => {
    const headingLevel = hashes.length;
    const demotedHeadingLevel = Math.min(6, headingLevel + HEADING_DEMOTION_LEVELS);
    return `${"#".repeat(demotedHeadingLevel)} `;
  });
}

/**
 * Generate an agentic workflow activity report issue.
 * @returns {Promise<void>}
 */
async function main() {
  const reportOutputDir = process.env.GH_AW_ACTIVITY_REPORT_OUTPUT_DIR || DEFAULT_REPORT_OUTPUT_DIR;
  const { owner, repo } = resolveExecutionOwnerRepo();
  const repoSlug = `${owner}/${repo}`;

  core.info(`Generating agentic workflow activity report for ${repoSlug} from cached trace index data`);

  const sections = [];
  for (const range of REPORT_RANGES) {
    sections.push(await readCachedRangeReport(range, reportOutputDir));
  }

  const headerLines = ["### Agentic workflow activity report", "", `Repository: \`${repoSlug}\``, `Generated at: ${new Date().toISOString()}`, ""];
  const sectionLines = sections.flatMap(section => ["<details>", `<summary>${section.heading}</summary>`, "", section.body, "", "</details>", ""]);
  const body = [...headerLines, ...sectionLines].join("\n");

  const createdIssue = await github.rest.issues.create({
    owner,
    repo,
    title: ISSUE_TITLE,
    body,
    labels: ["agentic-workflows"],
  });

  core.info(`Created issue #${createdIssue.data.number}: ${createdIssue.data.html_url}`);
}

module.exports = { main, readCachedRangeReport, normalizeReportMarkdown };
