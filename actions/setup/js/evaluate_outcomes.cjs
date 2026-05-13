// @ts-check

/**
 * evaluate_outcomes.cjs
 *
 * Evaluates safe output outcomes for recent successful workflow runs.
 * Replaces the shell-based evaluation logic in the outcome-collector workflow.
 *
 * Responsibilities:
 * - Load previously evaluated run IDs from cache-memory
 * - Fetch recent successful runs via `gh run list`
 * - Download safe-outputs-items artifacts via `gh run download`
 * - Classify each item (accepted/rejected/pending/noop) using the GitHub API
 * - Extract time-to-resolution, PR quality signals, pending age
 * - Write per-item evaluations to outcome-evaluations.jsonl
 * - Compute and write fleet summary to outcome-summary.json
 * - Update the seen-runs cache
 *
 * Outputs:
 *   /tmp/gh-aw/outcome-evaluations.jsonl  — per-item JSONL
 *   /tmp/gh-aw/outcome-summary.json       — fleet summary
 *   /tmp/gh-aw/outcomes/run-*.json        — per-run data
 *
 * Errors in individual run/item evaluation are non-fatal and logged to stderr.
 */

const fs = require("fs");
const path = require("path");
const { execFileSync } = require("child_process");

// ---------------------------------------------------------------------------
// Paths
// ---------------------------------------------------------------------------
const CACHE_DIR = "/tmp/gh-aw/cache-memory/outcome-collector";
const SEEN_FILE = path.join(CACHE_DIR, "seen-runs.json");
const OUTCOMES_DIR = "/tmp/gh-aw/outcomes";
const EVAL_JSONL = "/tmp/gh-aw/outcome-evaluations.jsonl";
const SUMMARY_PATH = "/tmp/gh-aw/outcome-summary.json";
let ghCommandRunner = args => execFileSync("gh", args, { encoding: "utf8", stdio: ["pipe", "pipe", "pipe"] });

// ---------------------------------------------------------------------------
// Noop types that are tracked but not counted as actionable
// ---------------------------------------------------------------------------
const NOOP_TYPES = new Set(["noop", "missing_tool", "missing_data", "report_incomplete"]);

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/**
 * Run a `gh` CLI command and capture stdout/stderr.
 * @param {string[]} args
 * @returns {{ok: true, stdout: string, stderr: string, status: number} | {ok: false, stdout: string, stderr: string, status: number, error: string}}
 */
function ghResult(args) {
  try {
    return {
      ok: true,
      stdout: ghCommandRunner(args),
      stderr: "",
      status: 0,
    };
  } catch (error) {
    return {
      ok: false,
      stdout: typeof error?.stdout === "string" ? error.stdout : String(error?.stdout || ""),
      stderr: typeof error?.stderr === "string" ? error.stderr : String(error?.stderr || ""),
      status: typeof error?.status === "number" ? error.status : 1,
      error: error?.message || "gh command failed",
    };
  }
}

/**
 * Run a `gh` CLI command, returning stdout as a string.
 * Returns null on failure.
 * @param {string[]} args
 * @returns {string | null}
 */
function gh(args) {
  const result = ghResult(args);
  return result.ok ? result.stdout.trim() : null;
}

/**
 * Run a `gh api` call, returning parsed JSON plus error metadata.
 * @param {string} endpoint
 * @returns {{ok: true, data: any} | {ok: false, status: number, error: string}}
 */
function ghAPI(endpoint) {
  const result = ghResult(["api", endpoint]);
  if (!result.ok) {
    return {
      ok: false,
      status: result.status,
      error: [result.stderr, result.stdout, result.error].filter(Boolean).join(" ").trim() || "gh api failed",
    };
  }
  try {
    return { ok: true, data: JSON.parse(result.stdout) };
  } catch {
    return { ok: false, status: 0, error: `failed to parse JSON for ${endpoint}` };
  }
}

/**
 * Read a JSON file, returning a default value on failure.
 * @param {string} filePath
 * @param {any} fallback
 * @returns {any}
 */
function readJSON(filePath, fallback) {
  try {
    return JSON.parse(fs.readFileSync(filePath, "utf8"));
  } catch {
    return fallback;
  }
}

/**
 * Read a JSONL file, returning an array of parsed objects.
 * @param {string} filePath
 * @returns {object[]}
 */
function readJSONL(filePath) {
  try {
    return fs
      .readFileSync(filePath, "utf8")
      .split("\n")
      .filter(l => l.trim())
      .map(l => {
        try {
          return JSON.parse(l);
        } catch {
          return null;
        }
      })
      .filter(Boolean);
  } catch {
    return [];
  }
}

/**
 * Atomically write JSON to a file using a tmp+rename swap.
 * @param {string} filePath
 * @param {any} data
 */
function writeJSONAtomic(filePath, data) {
  const tmp = filePath + ".tmp";
  fs.writeFileSync(tmp, JSON.stringify(data, null, 2) + "\n");
  fs.renameSync(tmp, filePath);
}

/**
 * Parse an ISO-8601 timestamp to epoch seconds. Returns null on failure.
 * @param {string} ts
 * @returns {number | null}
 */
function isoToEpoch(ts) {
  if (!ts) return null;
  const ms = Date.parse(ts);
  return Number.isFinite(ms) ? Math.floor(ms / 1000) : null;
}

/**
 * Compute seconds between two ISO timestamps. Returns null if either is invalid.
 * @param {string} from
 * @param {string} to
 * @returns {number | null}
 */
function secondsBetween(from, to) {
  const a = isoToEpoch(from);
  const b = isoToEpoch(to);
  if (a === null || b === null) return null;
  return b - a;
}

/**
 * @param {string} url
 * @returns {string}
 */
function parseRepoFromURL(url) {
  const match = url.match(/github\.com\/([^/]+\/[^/#?]+)/);
  return match ? match[1] : "";
}

/**
 * @param {string} url
 * @returns {number}
 */
function parseNumberFromURL(url) {
  const match = url.match(/\/(?:issues|pull|discussions)\/(\d+)/);
  return match ? Number(match[1]) : 0;
}

/**
 * @param {string} url
 * @returns {string}
 */
function extractCommentID(url) {
  const hashMatch = url.match(/#issuecomment-(\d+)/);
  if (hashMatch) return hashMatch[1];
  const pathMatch = url.match(/\/issues\/comments\/(\d+)/);
  return pathMatch ? pathMatch[1] : "";
}

/**
 * @param {string} login
 * @returns {boolean}
 */
function isBotUser(login) {
  return /\[bot\]$/.test(login) || login === "github-actions" || login.startsWith("copilot-");
}

/**
 * @param {string} endpoint
 * @returns {Array<Record<string, any>>}
 */
function ghAPIArray(endpoint) {
  const result = ghAPI(endpoint);
  return result.ok && Array.isArray(result.data) ? result.data : [];
}

/**
 * @param {Array<Record<string, any>>} comments
 * @returns {number}
 */
function countHumanComments(comments) {
  let count = 0;
  for (const comment of comments) {
    const user = comment && typeof comment.user === "object" ? comment.user : null;
    const login = user && typeof user.login === "string" ? user.login : "";
    if (login && !isBotUser(login)) {
      count++;
    }
  }
  return count;
}

/**
 * @param {number[]} values
 * @returns {number | null}
 */
function median(values) {
  if (values.length === 0) return null;
  const sorted = [...values].sort((a, b) => a - b);
  const mid = Math.floor(sorted.length / 2);
  if (sorted.length % 2 === 1) return sorted[mid];
  return Math.round((sorted[mid - 1] + sorted[mid]) / 2);
}

/**
 * @param {number} value
 * @returns {number}
 */
function round4(value) {
  return Math.round(value * 10000) / 10000;
}

/**
 * @param {string} name
 * @returns {{name: string, runs: number, outcomes: number, accepted: number, rejected: number, ignored: number, pending: number, lifecycle: number, noop: number, zero_touch: number, acceptance_rate?: number, waste_rate?: number, zero_touch_rate?: number}}
 */
function createWorkflowBucket(name) {
  return {
    name,
    runs: 0,
    outcomes: 0,
    accepted: 0,
    rejected: 0,
    ignored: 0,
    pending: 0,
    lifecycle: 0,
    noop: 0,
    zero_touch: 0,
  };
}

/**
 * @param {string} type
 * @returns {{type: string, total: number, accepted: number, rejected: number, ignored: number, pending: number, lifecycle: number, zero_touch: number, acceptance_rate?: number, waste_rate?: number, zero_touch_rate?: number}}
 */
function createTypeBucket(type) {
  return {
    type,
    total: 0,
    accepted: 0,
    rejected: 0,
    ignored: 0,
    pending: 0,
    lifecycle: 0,
    zero_touch: 0,
  };
}

/**
 * @param {string} event
 * @returns {{event: string, runs: number, outcomes: number, accepted: number, rejected: number, ignored: number, pending: number, lifecycle: number, noop: number, zero_touch: number, acceptance_rate?: number, waste_rate?: number, zero_touch_rate?: number}}
 */
function createEventBucket(event) {
  return {
    event,
    runs: 0,
    outcomes: 0,
    accepted: 0,
    rejected: 0,
    ignored: 0,
    pending: 0,
    lifecycle: 0,
    noop: 0,
    zero_touch: 0,
  };
}

/**
 * @param {{accepted: number, rejected: number, zero_touch: number, acceptance_rate?: number, waste_rate?: number, zero_touch_rate?: number}} bucket
 * @param {number} total
 */
function finalizeRates(bucket, total) {
  const resolved = bucket.accepted + bucket.rejected;
  if (resolved > 0) {
    bucket.acceptance_rate = round4(bucket.accepted / resolved);
  }
  if (total > 0) {
    bucket.waste_rate = round4(bucket.rejected / total);
  }
  if (bucket.accepted > 0) {
    bucket.zero_touch_rate = round4(bucket.zero_touch / bucket.accepted);
  }
}

// ---------------------------------------------------------------------------
// Item evaluation
// ---------------------------------------------------------------------------

/**
 * @typedef {object} EvalResult
 * @property {string} result
 * @property {string} detail
 * @property {number | null} resolution_sec
 * @property {number | null} pending_age_sec
 * @property {number | null} review_comments
 * @property {number | null} human_comments
 * @property {number | null} human_reviews
 * @property {number | null} changed_files
 * @property {number | null} additions
 * @property {number | null} deletions
 * @property {number | null} reactions
 * @property {number | null} replies
 * @property {boolean} zero_touch
 * @property {string} state_reason
 * @property {string} closed_by
 * @property {boolean | null} closed_by_bot
 */

/**
 * Evaluate a single safe-output item against the GitHub API.
 * @param {object} item
 * @param {string} defaultRepo
 * @returns {EvalResult}
 */
function evaluateItem(item, defaultRepo) {
  const url = item.url || "";
  const itemRepo = item.repo || parseRepoFromURL(url) || defaultRepo;
  const itemNumber = typeof item.number === "number" ? item.number : parseNumberFromURL(url);
  const timestamp = item.timestamp || "";
  const type = item.type || "";

  /** @type {EvalResult} */
  const out = {
    result: "pending",
    detail: "",
    resolution_sec: null,
    pending_age_sec: null,
    review_comments: null,
    human_comments: null,
    human_reviews: null,
    changed_files: null,
    additions: null,
    deletions: null,
    reactions: null,
    replies: null,
    zero_touch: false,
    state_reason: "",
    closed_by: "",
    closed_by_bot: null,
  };

  if (!itemRepo) {
    out.detail = "missing repo";
    setPendingAge(out, timestamp);
    return out;
  }

  if (type === "create_issue") {
    const issue = ghAPI(`repos/${itemRepo}/issues/${itemNumber}`);
    if (!issue.ok || !issue.data || !issue.data.state) {
      out.detail = "api error";
      setPendingAge(out, timestamp);
      return out;
    }

    const data = issue.data;
    const comments = ghAPIArray(`repos/${itemRepo}/issues/${itemNumber}/comments`);
    out.human_comments = countHumanComments(comments);
    out.state_reason = typeof data.state_reason === "string" ? data.state_reason : "";

    if (data.state === "closed") {
      if (out.state_reason === "not_planned") {
        const events = ghAPIArray(`repos/${itemRepo}/issues/${itemNumber}/events`);
        for (let i = events.length - 1; i >= 0; i--) {
          const event = events[i];
          if (event && event.event === "closed") {
            const actor = event.actor && typeof event.actor === "object" ? event.actor : null;
            const login = actor && typeof actor.login === "string" ? actor.login : "";
            out.closed_by = login;
            out.closed_by_bot = login ? isBotUser(login) : null;
            break;
          }
        }
        out.result = out.closed_by_bot ? "lifecycle" : "rejected";
        out.detail = out.closed_by_bot ? "closed by bot" : "closed as not planned";
      } else {
        out.result = "accepted";
        out.detail = out.state_reason || "closed";
      }
      if (timestamp && data.closed_at) {
        out.resolution_sec = secondsBetween(timestamp, data.closed_at);
      }
      return out;
    }

    if ((out.human_comments || 0) > 0) {
      out.result = "pending";
      out.detail = `${out.human_comments} human comments`;
    } else {
      out.result = "ignored";
      out.detail = "open, no engagement";
    }
    return out;
  }

  if (type === "create_pull_request" || type === "push_to_pull_request_branch" || type === "mark_pull_request_as_ready_for_review") {
    const pr = ghAPI(`repos/${itemRepo}/pulls/${itemNumber}`);
    if (!pr.ok || !pr.data || !pr.data.state) {
      out.detail = "api error";
      setPendingAge(out, timestamp);
      return out;
    }

    const data = pr.data;
    const comments = ghAPIArray(`repos/${itemRepo}/issues/${itemNumber}/comments`);
    const reviews = ghAPIArray(`repos/${itemRepo}/pulls/${itemNumber}/reviews`);
    out.human_comments = countHumanComments(comments);
    out.human_reviews = reviews.filter(review => {
      const user = review && typeof review.user === "object" ? review.user : null;
      const login = user && typeof user.login === "string" ? user.login : "";
      return login && !isBotUser(login);
    }).length;

    out.review_comments = typeof data.review_comments === "number" ? data.review_comments : null;
    out.changed_files = typeof data.changed_files === "number" ? data.changed_files : null;
    out.additions = typeof data.additions === "number" ? data.additions : null;
    out.deletions = typeof data.deletions === "number" ? data.deletions : null;

    if (data.merged === true) {
      out.result = "accepted";
      out.detail = "merged";
      out.zero_touch = (out.human_comments || 0) === 0 && (out.human_reviews || 0) === 0;
      if (timestamp && data.merged_at) {
        out.resolution_sec = secondsBetween(timestamp, data.merged_at);
      }
    } else if (data.state === "closed") {
      out.result = "rejected";
      out.detail = "closed";
      if (timestamp && data.closed_at) {
        out.resolution_sec = secondsBetween(timestamp, data.closed_at);
      }
    } else if (data.state === "open") {
      out.result = "pending";
      out.detail = "open";
      setPendingAge(out, timestamp);
    } else {
      out.detail = "api error";
      setPendingAge(out, timestamp);
    }
    return out;
  }

  if (type === "add_comment") {
    const commentId = extractCommentID(url);
    if (!commentId || !itemNumber) {
      out.detail = "missing comment reference";
      setPendingAge(out, timestamp);
      return out;
    }
    const comment = ghAPI(`repos/${itemRepo}/issues/comments/${commentId}`);
    if (!comment.ok) {
      out.result = comment.status === 404 ? "rejected" : "pending";
      out.detail = comment.status === 404 ? "deleted" : "api error";
      setPendingAge(out, timestamp);
      return out;
    }
    const data = comment.data || {};
    const reactions = data.reactions && typeof data.reactions === "object" ? data.reactions : null;
    out.reactions = reactions && typeof reactions.total_count === "number" ? reactions.total_count : 0;

    const comments = ghAPIArray(`repos/${itemRepo}/issues/${itemNumber}/comments`);
    const createdAt = typeof data.created_at === "string" ? data.created_at : "";
    let replies = 0;
    for (const reply of comments) {
      const replyCreatedAt = reply && typeof reply.created_at === "string" ? reply.created_at : "";
      const user = reply && typeof reply.user === "object" ? reply.user : null;
      const login = user && typeof user.login === "string" ? user.login : "";
      if (createdAt && replyCreatedAt > createdAt && login && !isBotUser(login)) {
        replies++;
      }
    }
    out.replies = replies;
    out.human_comments = replies;
    if ((out.reactions || 0) > 0 || replies > 0) {
      out.result = "accepted";
      out.detail = `${out.reactions || 0} reactions, ${replies} replies`;
    } else {
      out.result = "ignored";
      out.detail = "no engagement";
    }
    return out;
  }

  if (type === "add_labels") {
    const labels = ghAPIArray(`repos/${itemRepo}/issues/${itemNumber}/labels`);
    if (labels.length === 0) {
      out.result = "rejected";
      out.detail = "all labels removed";
    } else {
      out.result = "pending";
      out.detail = "labels still present";
    }
    return out;
  }

  if (type === "close_issue" || type === "close_pull_request") {
    const itemPath = type === "close_pull_request" ? "pulls" : "issues";
    const objectData = ghAPI(`repos/${itemRepo}/${itemPath}/${itemNumber}`);
    if (!objectData.ok || !objectData.data || !objectData.data.state) {
      out.detail = "api error";
      setPendingAge(out, timestamp);
      return out;
    }
    out.result = objectData.data.state === "closed" ? "accepted" : "rejected";
    out.detail = objectData.data.state === "closed" ? "still closed" : "reopened";
    return out;
  }

  if (type === "assign_milestone") {
    const issue = ghAPI(`repos/${itemRepo}/issues/${itemNumber}`);
    if (!issue.ok || !issue.data) {
      out.detail = "api error";
      setPendingAge(out, timestamp);
      return out;
    }
    out.result = issue.data.milestone ? "accepted" : "rejected";
    out.detail = issue.data.milestone ? "milestone still assigned" : "milestone removed";
    return out;
  }

  if (itemNumber) {
    const issue = ghAPI(`repos/${itemRepo}/issues/${itemNumber}`);
    if (issue.ok && issue.data) {
      out.result = "accepted";
      out.detail = "object still exists";
      return out;
    }
  }

  out.result = url ? "accepted" : "pending";
  out.detail = url ? "object exists" : "no object reference";
  return out;
}

/**
 * Set pending_age_sec on the result if the item has a timestamp.
 * @param {EvalResult} out
 * @param {string} timestamp
 */
function setPendingAge(out, timestamp) {
  if (!timestamp) return;
  const itemEpoch = isoToEpoch(timestamp);
  if (itemEpoch === null) return;
  out.pending_age_sec = Math.floor(Date.now() / 1000) - itemEpoch;
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

function main() {
  const repo = process.env.GITHUB_REPOSITORY || "";
  if (!repo) {
    console.error("GITHUB_REPOSITORY is not set");
    process.exit(1);
  }

  // Ensure directories exist
  fs.mkdirSync(CACHE_DIR, { recursive: true });
  fs.mkdirSync(OUTCOMES_DIR, { recursive: true });

  // Load seen-runs cache
  const seenIds = new Set(readJSON(SEEN_FILE, []));

  // Fetch recent successful runs
  const runsRaw = gh(["run", "list", "--repo", repo, "--limit", "200", "--json", "databaseId,conclusion,workflowName,event", "--jq", '[.[] | select(.conclusion == "success")] | .[0:150]']);

  if (!runsRaw || runsRaw === "[]" || runsRaw === "null") {
    console.log("No recent successful runs found");
    writeJSONAtomic(SUMMARY_PATH, { runs_checked: 0, total_outcomes: 0 });
    process.exit(0);
  }

  /** @type {Array<{databaseId: number, workflowName: string, event: string}>} */
  let runs;
  try {
    runs = JSON.parse(runsRaw);
  } catch {
    console.error("Failed to parse run list");
    writeJSONAtomic(SUMMARY_PATH, { runs_checked: 0, total_outcomes: 0 });
    process.exit(0);
  }

  // Counters
  let checked = 0;
  let accepted = 0;
  let rejected = 0;
  let ignored = 0;
  let pending = 0;
  let lifecycle = 0;
  let total = 0;
  let noop = 0;
  let zeroTouch = 0;
  /** @type {number[]} */
  const resolutionTimes = [];
  /** @type {number[]} */
  const pendingAges = [];
  /** @type {Map<string, ReturnType<typeof createWorkflowBucket>>} */
  const workflowBuckets = new Map();
  /** @type {Map<string, ReturnType<typeof createTypeBucket>>} */
  const typeBuckets = new Map();
  /** @type {Map<string, ReturnType<typeof createEventBucket>>} */
  const eventBuckets = new Map();

  // Clear the evaluations file
  fs.writeFileSync(EVAL_JSONL, "");

  /** @type {number[]} */
  const evaluatedIds = [];

  for (const run of runs) {
    const runId = run.databaseId;
    const workflow = run.workflowName || "";
    const event = run.event || "";

    // Skip previously evaluated
    if (seenIds.has(runId)) continue;

    // Download artifact
    const itemDir = path.join(OUTCOMES_DIR, `run-${runId}`);
    const dlResult = gh(["run", "download", String(runId), "--repo", repo, "--name", "safe-outputs-items", "--dir", itemDir]);
    if (dlResult === null) continue;

    const manifestPath = path.join(itemDir, "safe-output-items.jsonl");
    if (!fs.existsSync(manifestPath)) continue;

    const manifest = readJSONL(manifestPath);
    if (manifest.length === 0) continue;

    // Separate actionable items from noops
    const actionable = manifest.filter(m => m.type && !NOOP_TYPES.has(m.type));
    const noops = manifest.filter(m => m.type && NOOP_TYPES.has(m.type));
    const runNoops = noops.length;
    const runItems = actionable.length;

    if (runItems === 0 && runNoops === 0) continue;

    noop += runNoops;

    console.log(`Run ${runId} (${workflow}): ${runItems} item(s), ${runNoops} noop(s) [trigger: ${event}]`);
    checked++;
    total += runItems;
    const workflowBucket = workflowBuckets.get(workflow) || createWorkflowBucket(workflow);
    workflowBucket.runs++;
    workflowBucket.outcomes += runItems;
    workflowBucket.noop += runNoops;
    workflowBuckets.set(workflow, workflowBucket);
    const eventBucket = eventBuckets.get(event) || createEventBucket(event);
    eventBucket.runs++;
    eventBucket.outcomes += runItems;
    eventBucket.noop += runNoops;
    eventBuckets.set(event, eventBucket);
    const runSummary = {
      workflow,
      run_id: runId,
      items: runItems,
      noops: runNoops,
      event,
      accepted: 0,
      rejected: 0,
      ignored: 0,
      pending: 0,
      lifecycle: 0,
      zero_touch: 0,
      types: {},
    };

    // Write noop entries
    for (const n of noops) {
      fs.appendFileSync(
        EVAL_JSONL,
        JSON.stringify({
          type: n.type,
          url: "",
          repo,
          result: "noop",
          detail: n.type,
          workflow,
          run_id: runId,
          timestamp: "",
          event,
        }) + "\n"
      );
    }

    if (runItems === 0) {
      // Only noops — still mark as evaluated
      writeJSONAtomic(path.join(OUTCOMES_DIR, `run-${runId}.json`), runSummary);
      evaluatedIds.push(runId);
      continue;
    }

    // Evaluate each actionable item
    for (const item of actionable) {
      const itemURL = item.url || "";
      const itemRepo = item.repo || parseRepoFromURL(itemURL) || repo;
      const itemNumber = typeof item.number === "number" ? item.number : parseNumberFromURL(itemURL);
      const evalResult = evaluateItem(item, repo);

      switch (evalResult.result) {
        case "accepted":
          accepted++;
          runSummary.accepted++;
          workflowBucket.accepted++;
          eventBucket.accepted++;
          break;
        case "rejected":
          rejected++;
          runSummary.rejected++;
          workflowBucket.rejected++;
          eventBucket.rejected++;
          break;
        case "ignored":
          runSummary.ignored++;
          ignored++;
          workflowBucket.ignored++;
          eventBucket.ignored++;
          break;
        case "lifecycle":
          lifecycle++;
          runSummary.lifecycle++;
          workflowBucket.lifecycle++;
          eventBucket.lifecycle++;
          break;
        default:
          pending++;
          runSummary.pending++;
          workflowBucket.pending++;
          eventBucket.pending++;
          break;
      }

      if (evalResult.zero_touch) {
        zeroTouch++;
        runSummary.zero_touch++;
        workflowBucket.zero_touch++;
        eventBucket.zero_touch++;
      }
      if (typeof evalResult.resolution_sec === "number" && evalResult.resolution_sec > 0) {
        resolutionTimes.push(evalResult.resolution_sec);
      }
      if (typeof evalResult.pending_age_sec === "number" && evalResult.pending_age_sec > 0) {
        pendingAges.push(evalResult.pending_age_sec);
      }

      const typeBucket = typeBuckets.get(item.type || "unknown") || createTypeBucket(item.type || "unknown");
      typeBucket.total++;
      if (evalResult.result === "accepted") typeBucket.accepted++;
      if (evalResult.result === "rejected") typeBucket.rejected++;
      if (evalResult.result === "ignored") typeBucket.ignored++;
      if (evalResult.result === "pending") typeBucket.pending++;
      if (evalResult.result === "lifecycle") typeBucket.lifecycle++;
      if (evalResult.zero_touch) typeBucket.zero_touch++;
      typeBuckets.set(typeBucket.type, typeBucket);
      runSummary.types[typeBucket.type] = {
        total: (runSummary.types[typeBucket.type]?.total || 0) + 1,
        accepted: (runSummary.types[typeBucket.type]?.accepted || 0) + (evalResult.result === "accepted" ? 1 : 0),
        rejected: (runSummary.types[typeBucket.type]?.rejected || 0) + (evalResult.result === "rejected" ? 1 : 0),
        ignored: (runSummary.types[typeBucket.type]?.ignored || 0) + (evalResult.result === "ignored" ? 1 : 0),
        pending: (runSummary.types[typeBucket.type]?.pending || 0) + (evalResult.result === "pending" ? 1 : 0),
        lifecycle: (runSummary.types[typeBucket.type]?.lifecycle || 0) + (evalResult.result === "lifecycle" ? 1 : 0),
      };

      fs.appendFileSync(
        EVAL_JSONL,
        JSON.stringify({
          type: item.type || "",
          url: itemURL,
          repo: item.repo || itemRepo || repo,
          number: item.number || itemNumber || 0,
          result: evalResult.result,
          detail: evalResult.detail,
          workflow,
          run_id: runId,
          timestamp: item.timestamp || "",
          event,
          resolution_sec: evalResult.resolution_sec,
          pending_age_sec: evalResult.pending_age_sec,
          review_comments: evalResult.review_comments,
          human_comments: evalResult.human_comments,
          human_reviews: evalResult.human_reviews,
          changed_files: evalResult.changed_files,
          additions: evalResult.additions,
          deletions: evalResult.deletions,
          reactions: evalResult.reactions,
          replies: evalResult.replies,
          zero_touch: evalResult.zero_touch,
          state_reason: evalResult.state_reason,
          closed_by: evalResult.closed_by,
          closed_by_bot: evalResult.closed_by_bot,
        }) + "\n"
      );
    }

    // Save per-run data
    finalizeRates(runSummary, runSummary.items);
    writeJSONAtomic(path.join(OUTCOMES_DIR, `run-${runId}.json`), runSummary);

    evaluatedIds.push(runId);
  }

  // Compute fleet summary
  const resolved = accepted + rejected;
  const acceptanceRate = resolved > 0 ? accepted / resolved : 0;
  const wasteRate = total > 0 ? rejected / total : 0;
  const noopRate = total + noop > 0 ? noop / (total + noop) : 0;
  const zeroTouchRate = accepted > 0 ? zeroTouch / accepted : 0;
  const medianResolutionSec = median(resolutionTimes);
  const medianPendingAgeSec = median(pendingAges);

  const workflows = [...workflowBuckets.values()]
    .map(bucket => {
      finalizeRates(bucket, bucket.outcomes);
      return bucket;
    })
    .sort((a, b) => b.outcomes - a.outcomes || a.name.localeCompare(b.name));
  const types = [...typeBuckets.values()]
    .map(bucket => {
      finalizeRates(bucket, bucket.total);
      return bucket;
    })
    .sort((a, b) => b.total - a.total || a.type.localeCompare(b.type));
  const events = [...eventBuckets.values()]
    .map(bucket => {
      finalizeRates(bucket, bucket.outcomes);
      return bucket;
    })
    .sort((a, b) => b.outcomes - a.outcomes || a.event.localeCompare(b.event));

  writeJSONAtomic(SUMMARY_PATH, {
    runs_checked: checked,
    total_outcomes: total,
    accepted,
    rejected,
    ignored,
    pending,
    lifecycle,
    noop,
    zero_touch: zeroTouch,
    acceptance_rate: Math.round(acceptanceRate * 10000) / 10000,
    waste_rate: Math.round(wasteRate * 10000) / 10000,
    noop_rate: Math.round(noopRate * 10000) / 10000,
    zero_touch_rate: Math.round(zeroTouchRate * 10000) / 10000,
    median_resolution_sec: medianResolutionSec,
    median_pending_age_sec: medianPendingAgeSec,
    workflows,
    types,
    events,
    date: new Date().toISOString().slice(0, 10),
  });

  // Update seen-runs cache: merge old + new, keep last 500
  const merged = [...new Set([...seenIds, ...evaluatedIds])].sort((a, b) => a - b).slice(-500);
  writeJSONAtomic(SEEN_FILE, merged);

  console.log(`✓ Checked ${checked} runs, ${total} outcomes`);
  console.log(`  Accepted: ${accepted}, Rejected: ${rejected}, Ignored: ${ignored}, Pending: ${pending}, Lifecycle: ${lifecycle}, Noop: ${noop}`);
  console.log(`  Acceptance rate: ${acceptanceRate.toFixed(4)}`);
  console.log(JSON.stringify(readJSON(SUMMARY_PATH, {}), null, 2));
}

/**
 * @param {(args: string[]) => string} runner
 */
function setGHCommandRunnerForTest(runner) {
  ghCommandRunner = runner;
}

function resetGHCommandRunnerForTest() {
  ghCommandRunner = args => execFileSync("gh", args, { encoding: "utf8", stdio: ["pipe", "pipe", "pipe"] });
}

if (require.main === module) {
  main();
}

module.exports = {
  main,
  evaluateItem,
  readJSONL,
  secondsBetween,
  isoToEpoch,
  parseRepoFromURL,
  parseNumberFromURL,
  extractCommentID,
  isBotUser,
  median,
  setGHCommandRunnerForTest,
  resetGHCommandRunnerForTest,
};
