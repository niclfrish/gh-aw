#!/bin/bash
set +o histexpand

# validate_gemini_api_key.sh - Validate the Gemini API key with a lightweight API call
#
# This script performs a minimal test call to the Google Generative Language API to detect
# expired, revoked, or invalid GEMINI_API_KEY values early in the workflow — before
# installing the Gemini CLI or running the agent. When the key is invalid the script fails
# immediately with a clear, actionable error message and rotation instructions.
#
# Usage: validate_gemini_api_key.sh
#
# Environment:
#   GEMINI_API_KEY - The Gemini API key to validate (required)
#
# Exit codes:
#   0 - API key is valid, or could not be validated (network unavailable — non-fatal)
#   1 - API key is invalid or expired

set -e

# If the key is not set, let validate_multi_secret.sh handle the empty-key case.
if [ -z "${GEMINI_API_KEY:-}" ]; then
  echo "GEMINI_API_KEY is not set — skipping API key validation" >&2
  exit 0
fi

echo "Validating GEMINI_API_KEY with Generative Language API…" >&2

TEMP_FILE=$(mktemp) || { echo "Error: Failed to create temporary file" >&2; exit 1; }
# Use the x-goog-api-key request header (instead of the ?key= query param) to avoid
# embedding the secret in a URL that may appear in server logs or process listings.
# Use pageSize=1 to minimise response payload — we only need the HTTP status.
HTTP_CODE=$(curl -s -o "${TEMP_FILE}" -w "%{http_code}" \
  "https://generativelanguage.googleapis.com/v1beta/models?pageSize=1" \
  -H "x-goog-api-key: ${GEMINI_API_KEY}" \
  --connect-timeout 10 --max-time 30 2>/dev/null || echo "000")

if [ "${HTTP_CODE}" = "000" ]; then
  # Network error or timeout — do not block the workflow; the agent step will catch this.
  echo "Warning: could not reach generativelanguage.googleapis.com (network timeout or unavailable)." >&2
  echo "Skipping API key pre-flight check; the agent will fail if the key is truly invalid." >&2
  rm -f "${TEMP_FILE}"
  exit 0
fi

if [ "${HTTP_CODE}" = "200" ]; then
  echo "✅ GEMINI_API_KEY is valid" >&2
  rm -f "${TEMP_FILE}"
  exit 0
fi

# The key is invalid — extract error details from the JSON response.
ERROR_CONTENT=$(cat "${TEMP_FILE}" 2>/dev/null || echo "{}")
rm -f "${TEMP_FILE}"

# Use jq for reliable JSON parsing when available; fall back to grep+sed.
if command -v jq >/dev/null 2>&1; then
  ERROR_MESSAGE=$(jq -r '.error.message // empty' <<<"${ERROR_CONTENT}" 2>/dev/null || true)
  ERROR_REASON=$(jq -r '.error.details[0].reason // empty' <<<"${ERROR_CONTENT}" 2>/dev/null || true)
else
  ERROR_MESSAGE=$(echo "${ERROR_CONTENT}" | grep -o '"message":"[^"]*"' | head -1 | sed 's/"message":"//;s/"$//')
  ERROR_REASON=$(echo "${ERROR_CONTENT}" | grep -o '"reason":"[^"]*"' | head -1 | sed 's/"reason":"//;s/"$//')
fi

{
  echo "❌ Error: GEMINI_API_KEY is invalid or expired (HTTP ${HTTP_CODE})"
  echo ""
  if [ -n "${ERROR_MESSAGE}" ]; then
    echo "API error: ${ERROR_MESSAGE}"
  fi
  if [ -n "${ERROR_REASON}" ]; then
    echo "Reason: ${ERROR_REASON}"
  fi
  echo ""
  echo "**How to fix:**"
  echo "1. Generate a new API key at: https://aistudio.google.com/apikey"
  echo "2. Go to your repository **Settings → Secrets and variables → Actions**"
  echo "3. Update the \`GEMINI_API_KEY\` secret with the new key"
  echo ""
  echo "**Note:** API keys can expire or be revoked. Rotate them regularly to avoid workflow failures."
} >> "${GITHUB_STEP_SUMMARY:-/dev/null}"

echo "Error: GEMINI_API_KEY is invalid or expired (HTTP ${HTTP_CODE})" >&2
if [ -n "${ERROR_MESSAGE}" ]; then
  echo "API error: ${ERROR_MESSAGE}" >&2
fi
if [ -n "${ERROR_REASON}" ]; then
  echo "Reason: ${ERROR_REASON}" >&2
fi
echo "" >&2
echo "Generate a new API key at: https://aistudio.google.com/apikey" >&2
echo "Then update the GEMINI_API_KEY repository secret." >&2

exit 1
