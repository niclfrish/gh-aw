// @ts-check
import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import fs from "node:fs/promises";
import path from "node:path";
import os from "node:os";

describe("run_activity_report", () => {
  let originalGlobals;
  let originalEnv;
  let mockCore;
  let mockGithub;
  let mockContext;
  let tempOutputDir;

  beforeEach(() => {
    originalEnv = { ...process.env };
    tempOutputDir = path.join(os.tmpdir(), `run-activity-report-${Date.now()}-${Math.random().toString(36).slice(2)}`);
    process.env.GH_AW_ACTIVITY_REPORT_OUTPUT_DIR = tempOutputDir;

    originalGlobals = {
      core: global.core,
      github: global.github,
      context: global.context,
    };

    mockCore = {
      info: vi.fn(),
      warning: vi.fn(),
    };
    mockGithub = {
      rest: {
        issues: {
          create: vi.fn().mockResolvedValue({
            data: { number: 42, html_url: "https://github.com/testowner/testrepo/issues/42" },
          }),
        },
      },
    };
    mockContext = {
      repo: {
        owner: "testowner",
        repo: "testrepo",
      },
    };

    global.core = mockCore;
    global.github = mockGithub;
    global.context = mockContext;
  });

  afterEach(async () => {
    process.env = originalEnv;
    global.core = originalGlobals.core;
    global.github = originalGlobals.github;
    global.context = originalGlobals.context;
    await fs.rm(tempOutputDir, { recursive: true, force: true });
    vi.clearAllMocks();
  });

  it("creates an activity report issue using cached 24h and 7d reports", async () => {
    await fs.mkdir(path.join(tempOutputDir, "activity-report"), { recursive: true });
    await fs.writeFile(path.join(tempOutputDir, "activity-report", "24h.md"), "## 24h report\nok\n", "utf8");
    await fs.writeFile(path.join(tempOutputDir, "activity-report", "7d.md"), "## 7d report\nok\n", "utf8");

    const { main } = await import("./run_activity_report.cjs");
    await main();

    expect(mockGithub.rest.issues.create).toHaveBeenCalledWith(
      expect.objectContaining({
        owner: "testowner",
        repo: "testrepo",
        title: "[aw] agentic status report",
        labels: ["agentic-workflows"],
      })
    );

    const issueBody = mockGithub.rest.issues.create.mock.calls[0][0].body;
    expect(issueBody).toContain("### Agentic workflow activity report");
    expect(issueBody).toContain("<details>");
    expect(issueBody).toContain("<summary>Last 24 hours</summary>");
    expect(issueBody).toContain("<summary>Last 7 days</summary>");
    expect(issueBody).not.toContain("<summary>Last 30 days</summary>");
    expect(issueBody).toContain("#### 24h report");
  });

  it("uses fallback text when cached range reports are missing", async () => {
    const { main } = await import("./run_activity_report.cjs");
    await main();

    const issueBody = mockGithub.rest.issues.create.mock.calls[0][0].body;
    expect(issueBody).toContain("<summary>Last 24 hours</summary>");
    expect(issueBody).toContain("_No cached trace index is available for this range yet._");
    expect(mockCore.warning).toHaveBeenCalled();
  });

  it("demotes report headings by two levels", async () => {
    const { normalizeReportMarkdown } = await import("./run_activity_report.cjs");
    const transformed = normalizeReportMarkdown("# H1\n## H2\n### H3");
    expect(transformed).toContain("### H1");
    expect(transformed).toContain("#### H2");
    expect(transformed).toContain("##### H3");
  });
});
