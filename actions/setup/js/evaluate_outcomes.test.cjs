import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import fs from "fs";
import path from "path";
import { createRequire } from "module";

const req = createRequire(import.meta.url);

const OUTCOMES_DIR = "/tmp/gh-aw/outcomes";
const SUMMARY_PATH = "/tmp/gh-aw/outcome-summary.json";
const EVAL_JSONL = "/tmp/gh-aw/outcome-evaluations.jsonl";
const SEEN_FILE = "/tmp/gh-aw/cache-memory/outcome-collector/seen-runs.json";

describe("evaluate_outcomes.cjs", () => {
  /** @type {any} */
  let moduleUnderTest;

  beforeEach(() => {
    vi.clearAllMocks();
    fs.rmSync("/tmp/gh-aw", { recursive: true, force: true });
    fs.mkdirSync(OUTCOMES_DIR, { recursive: true });
    process.env.GITHUB_REPOSITORY = "github/gh-aw";
    const modulePath = req.resolve("./evaluate_outcomes.cjs");
    delete req.cache[modulePath];
    moduleUnderTest = req("./evaluate_outcomes.cjs");
  });

  afterEach(() => {
    moduleUnderTest.resetGHCommandRunnerForTest();
    delete process.env.GITHUB_REPOSITORY;
    fs.rmSync("/tmp/gh-aw", { recursive: true, force: true });
  });

  it("classifies bot-closed issues as lifecycle outcomes", () => {
    moduleUnderTest.setGHCommandRunnerForTest(args => {
      const [, endpoint] = args;
      if (endpoint === "repos/github/gh-aw/issues/42") {
        return JSON.stringify({
          state: "closed",
          state_reason: "not_planned",
          closed_at: "2026-05-13T11:00:00Z",
        });
      }
      if (endpoint === "repos/github/gh-aw/issues/42/comments") {
        return JSON.stringify([]);
      }
      if (endpoint === "repos/github/gh-aw/issues/42/events") {
        return JSON.stringify([{ event: "closed", actor: { login: "github-actions[bot]" } }]);
      }
      throw new Error(`unexpected gh args: ${args.join(" ")}`);
    });

    const result = moduleUnderTest.evaluateItem(
      {
        type: "create_issue",
        repo: "github/gh-aw",
        number: 42,
        timestamp: "2026-05-13T10:00:00Z",
      },
      "github/gh-aw"
    );

    expect(result.result).toBe("lifecycle");
    expect(result.closed_by).toBe("github-actions[bot]");
    expect(result.closed_by_bot).toBe(true);
    expect(result.resolution_sec).toBe(3600);
  });

  it("writes richer workflow, type, and zero-touch summary data", () => {
    fs.mkdirSync(path.join(OUTCOMES_DIR, "run-101"), { recursive: true });
    fs.writeFileSync(
      path.join(OUTCOMES_DIR, "run-101", "safe-output-items.jsonl"),
      [
        JSON.stringify({
          type: "create_issue",
          repo: "github/gh-aw",
          number: 11,
          timestamp: "2026-05-13T09:00:00Z",
        }),
        JSON.stringify({
          type: "create_pull_request",
          repo: "github/gh-aw",
          number: 22,
          timestamp: "2026-05-13T09:30:00Z",
        }),
      ].join("\n") + "\n"
    );

    moduleUnderTest.setGHCommandRunnerForTest(args => {
      if (args[0] === "run" && args[1] === "list") {
        return JSON.stringify([{ databaseId: 101, workflowName: "triage", event: "issues" }]);
      }
      if (args[0] === "run" && args[1] === "download") {
        return "";
      }

      const endpoint = args[1];
      if (endpoint === "repos/github/gh-aw/issues/11") {
        return JSON.stringify({ state: "open", comments: 0 });
      }
      if (endpoint === "repos/github/gh-aw/issues/11/comments") {
        return JSON.stringify([]);
      }
      if (endpoint === "repos/github/gh-aw/pulls/22") {
        return JSON.stringify({
          state: "closed",
          merged: true,
          merged_at: "2026-05-13T10:30:00Z",
          review_comments: 0,
          changed_files: 2,
          additions: 20,
          deletions: 5,
        });
      }
      if (endpoint === "repos/github/gh-aw/issues/22/comments") {
        return JSON.stringify([]);
      }
      if (endpoint === "repos/github/gh-aw/pulls/22/reviews") {
        return JSON.stringify([]);
      }

      throw new Error(`unexpected gh args: ${args.join(" ")}`);
    });

    moduleUnderTest.main();

    const summary = JSON.parse(fs.readFileSync(SUMMARY_PATH, "utf8"));
    expect(summary.total_outcomes).toBe(2);
    expect(summary.accepted).toBe(1);
    expect(summary.ignored).toBe(1);
    expect(summary.zero_touch).toBe(1);
    expect(summary.zero_touch_rate).toBe(1);
    expect(summary.workflows).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          name: "triage",
          outcomes: 2,
          accepted: 1,
          ignored: 1,
        }),
      ])
    );
    expect(summary.types).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ type: "create_issue", ignored: 1 }),
        expect.objectContaining({ type: "create_pull_request", accepted: 1, zero_touch: 1 }),
      ])
    );
    expect(summary.events).toEqual(expect.arrayContaining([expect.objectContaining({ event: "issues", outcomes: 2 })]));

    const evaluations = fs
      .readFileSync(EVAL_JSONL, "utf8")
      .trim()
      .split("\n")
      .map(line => JSON.parse(line));
    expect(evaluations).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ type: "create_issue", result: "ignored" }),
        expect.objectContaining({ type: "create_pull_request", result: "accepted", zero_touch: true }),
      ])
    );
    expect(fs.existsSync(SEEN_FILE)).toBe(true);
  });
});
