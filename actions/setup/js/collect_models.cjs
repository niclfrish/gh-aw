// @ts-check
/// <reference types="@actions/github-script" />

const fs = require("fs");
const path = require("path");
const http = require("http");
const https = require("https");

function ensureParent(filePath) {
  fs.mkdirSync(path.dirname(filePath), { recursive: true });
}

function normalizeURL(host, route) {
  const base = host.startsWith("http://") || host.startsWith("https://") ? host : `https://${host}`;
  const url = new URL(base);
  const normalizedRoute = route.startsWith("/") ? route : `/${route}`;
  url.pathname = normalizedRoute;
  url.search = "";
  return url;
}

function requestJSON(url, headers) {
  const client = url.protocol === "http:" ? http : https;
  return new Promise((resolve, reject) => {
    const req = client.request(
      url,
      {
        method: "GET",
        headers,
        timeout: 15000,
      },
      res => {
        let body = "";
        res.on("data", chunk => {
          body += chunk;
        });
        res.on("end", () => {
          resolve({ statusCode: res.statusCode || 0, body });
        });
      }
    );

    req.on("error", reject);
    req.on("timeout", () => req.destroy(new Error("request timeout")));
    req.end();
  });
}

function extractModels(payload) {
  if (Array.isArray(payload?.data)) {
    return payload.data;
  }
  if (Array.isArray(payload?.models)) {
    return payload.models;
  }
  return [];
}

function modelDisplayName(model) {
  if (typeof model === "string") {
    return model;
  }
  if (model && typeof model === "object") {
    return model.id || model.name || model.model || JSON.stringify(model);
  }
  return String(model);
}

function buildHeaders(authType, token) {
  const headers = {
    Accept: "application/json",
    "User-Agent": "gh-aw-models-collector",
  };

  if (!token) {
    return headers;
  }

  if (authType === "x-api-key") {
    headers["x-api-key"] = token;
    headers["anthropic-version"] = "2023-06-01";
    return headers;
  }

  if (authType === "x-goog-api-key") {
    headers["x-goog-api-key"] = token;
    return headers;
  }

  headers.Authorization = `Bearer ${token}`;
  return headers;
}

async function writeSummary(engineId, models, warning) {
  if (warning) {
    core.summary.addDetails(`Available Models (${engineId})`, `\n\n⚠️ ${warning}\n`);
    await core.summary.write();
    return;
  }

  const names = models.map(modelDisplayName);
  const lines = names
    .slice(0, 100)
    .map(name => `- \`${name}\``)
    .join("\n");
  const hidden = names.length > 100 ? `\n\n_...and ${names.length - 100} more._` : "";
  core.summary.addDetails(`Available Models (${engineId}, ${names.length})`, `\n\n${lines}${hidden}\n`);
  await core.summary.write();
}

async function main() {
  const engineId = process.env.GH_AW_ENGINE_ID || "unknown";
  const modelsHost = process.env.GH_AW_MODELS_HOST || "";
  const modelsRoute = process.env.GH_AW_MODELS_ROUTE || "";
  const modelsFile = process.env.GH_AW_MODELS_FILE || "";
  const artifactFile = process.env.GH_AW_MODELS_ARTIFACT_FILE || "";
  const authType = process.env.GH_AW_MODELS_AUTH_TYPE || "bearer";
  const token = process.env.GH_AW_MODELS_TOKEN || "";

  if (!modelsHost || !modelsRoute || !modelsFile) {
    core.info("Skipping model collection: host, route, or output path missing");
    return;
  }

  try {
    const url = normalizeURL(modelsHost, modelsRoute);
    const headers = buildHeaders(authType, token);
    const response = await requestJSON(url, headers);

    if (response.statusCode === 401 || response.statusCode === 403) {
      const warning = `Could not fetch models from ${url.host}${url.pathname}: invalid or unauthorized token (HTTP ${response.statusCode})`;
      core.warning(warning);
      await writeSummary(engineId, [], warning);
      return;
    }

    if (response.statusCode < 200 || response.statusCode >= 300) {
      const warning = `Could not fetch models from ${url.host}${url.pathname}: HTTP ${response.statusCode}`;
      core.warning(warning);
      await writeSummary(engineId, [], warning);
      return;
    }

    let payload = {};
    try {
      payload = JSON.parse(response.body);
    } catch (error) {
      const warning = `Could not parse models response JSON: ${error instanceof Error ? error.message : String(error)}`;
      core.warning(warning);
      await writeSummary(engineId, [], warning);
      return;
    }

    const models = extractModels(payload);
    ensureParent(modelsFile);
    fs.writeFileSync(modelsFile, JSON.stringify(payload, null, 2) + "\n", "utf8");
    core.info(`Saved models response to ${modelsFile}`);

    if (artifactFile) {
      ensureParent(artifactFile);
      fs.copyFileSync(modelsFile, artifactFile);
      core.info(`Copied models response to artifact path ${artifactFile}`);
    }

    await writeSummary(engineId, models, models.length === 0 ? "No models were returned by the provider." : "");
  } catch (error) {
    const warning = `Model collection failed: ${error instanceof Error ? error.message : String(error)}`;
    core.warning(warning);
    await writeSummary(engineId, [], warning);
  }
}

if (typeof module !== "undefined" && module.exports) {
  module.exports = {
    main,
    normalizeURL,
    extractModels,
    modelDisplayName,
    buildHeaders,
  };
}

if (require.main === module) {
  main();
}
