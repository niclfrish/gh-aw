// @ts-check
import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import fs from "node:fs/promises";
import path from "node:path";
import os from "node:os";

describe("run_trace_indexer", () => {
  let originalGlobals;
  let originalEnv;
  let mockCore;
  let mockContext;
  let mockExec;
  let tempOutputDir;

  beforeEach(() => {
    originalEnv = { ...process.env };
    tempOutputDir = path.join(os.tmpdir(), `run-trace-indexer-${Date.now()}-${Math.random().toString(36).slice(2)}`);
    process.env.GH_AW_CMD_PREFIX = "gh aw";
    process.env.GH_AW_TRACE_INDEX_OUTPUT_DIR = tempOutputDir;

    originalGlobals = {
      core: global.core,
      context: global.context,
      exec: global.exec,
    };

    mockCore = {
      info: vi.fn(),
      warning: vi.fn(),
    };
    mockContext = {
      repo: {
        owner: "testowner",
        repo: "testrepo",
      },
    };
    mockExec = {
      getExecOutput: vi.fn(),
    };

    global.core = mockCore;
    global.context = mockContext;
    global.exec = mockExec;
  });

  afterEach(async () => {
    process.env = originalEnv;
    global.core = originalGlobals.core;
    global.context = originalGlobals.context;
    global.exec = originalGlobals.exec;
    await fs.rm(tempOutputDir, { recursive: true, force: true });
    vi.clearAllMocks();
  });

  it("runs 24h and 7d trace indexing and stores markdown sections", async () => {
    mockExec.getExecOutput.mockResolvedValueOnce({ stdout: "## 24h report\nok", stderr: "", exitCode: 0 }).mockResolvedValueOnce({ stdout: "## 7d report\nok", stderr: "", exitCode: 0 });

    const { main } = await import("./run_trace_indexer.cjs");
    await main();

    expect(mockExec.getExecOutput).toHaveBeenCalledTimes(2);
    expect(mockExec.getExecOutput).toHaveBeenNthCalledWith(
      1,
      "gh",
      expect.arrayContaining(["aw", "logs", "--repo", "testowner/testrepo", "--start-date", "-1d", "--count", "1000", "--output", tempOutputDir, "--format", "markdown"]),
      expect.objectContaining({ ignoreReturnCode: true })
    );
    expect(mockExec.getExecOutput).toHaveBeenNthCalledWith(
      2,
      "gh",
      expect.arrayContaining(["aw", "logs", "--repo", "testowner/testrepo", "--start-date", "-1w", "--count", "1000", "--output", tempOutputDir, "--format", "markdown"]),
      expect.objectContaining({ ignoreReturnCode: true })
    );

    const report24h = await fs.readFile(path.join(tempOutputDir, "activity-report", "24h.md"), "utf8");
    const report7d = await fs.readFile(path.join(tempOutputDir, "activity-report", "7d.md"), "utf8");
    expect(report24h).toContain("## 24h report");
    expect(report7d).toContain("## 7d report");
  });

  it("writes fallback report text and fails when trace indexing has range failures", async () => {
    mockExec.getExecOutput.mockResolvedValueOnce({ stdout: "", stderr: "API rate limit exceeded", exitCode: 1 }).mockResolvedValueOnce({ stdout: "## 7d report\nok", stderr: "", exitCode: 0 });

    const { main } = await import("./run_trace_indexer.cjs");
    await expect(main()).rejects.toThrow("Trace indexing completed with one or more range failures");

    const report24h = await fs.readFile(path.join(tempOutputDir, "activity-report", "24h.md"), "utf8");
    expect(report24h).toContain("rate limiting");
  });
});
