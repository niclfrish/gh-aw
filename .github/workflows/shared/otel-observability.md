---
env:
  OTEL_BACKEND_URL: ${{ secrets.GH_AW_OTEL_ENDPOINT }}
  OTEL_BACKEND_HEADERS: ${{ secrets.GH_AW_OTEL_HEADERS }}
observability:
  otlp:
    endpoint: ${{ secrets.GH_AW_OTEL_ENDPOINT }}
    headers: ${{ secrets.GH_AW_OTEL_HEADERS }}
mcp-servers:
  otel:
    command: npx
    args: ["@your-org/otel-query-mcp"]
    env:
      OTEL_BACKEND_URL: ${{ env.OTEL_BACKEND_URL }}
      OTEL_BACKEND_HEADERS: ${{ env.OTEL_BACKEND_HEADERS }}
---
