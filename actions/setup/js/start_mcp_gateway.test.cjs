import { describe, expect, it } from "vitest";
import { createRequire } from "node:module";

const require = createRequire(import.meta.url);
const { normalizeGatewayConfigForRuntime } = require("./start_mcp_gateway.cjs");

describe("start_mcp_gateway normalizeGatewayConfigForRuntime", () => {
  it("removes gateway.opentelemetry when endpoint is empty", () => {
    const config = {
      gateway: {
        port: 8080,
        domain: "localhost",
        apiKey: "key",
        opentelemetry: { endpoint: "", headers: "Authorization=token" },
      },
    };

    const result = normalizeGatewayConfigForRuntime(config);

    expect(result.removedEmptyOtelEndpoint).toBe(true);
    expect(config.gateway).not.toHaveProperty("opentelemetry");
  });

  it("keeps gateway.opentelemetry when endpoint is non-empty", () => {
    const config = {
      gateway: {
        port: 8080,
        domain: "localhost",
        apiKey: "key",
        opentelemetry: { endpoint: "https://otlp.example.com", headers: "Authorization=token" },
      },
    };

    const result = normalizeGatewayConfigForRuntime(config);

    expect(result.removedEmptyOtelEndpoint).toBe(false);
    expect(config.gateway).toHaveProperty("opentelemetry");
  });
});
