---
name: otlp-env-vars
description: Configure and troubleshoot OpenTelemetry SDK environment variables safely.
---

# Skill: Configure OpenTelemetry SDK via Environment Variables

Use this skill when configuring, reviewing, or troubleshooting OpenTelemetry SDK environment variables.

Source: OpenTelemetry Environment Variable Specification, OTel 1.56.0: https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/

## Core rules

- Treat an empty environment variable exactly like an unset variable.
- Boolean values are true only when set to `"true"` case-insensitively; all other values are false.
- Invalid numeric values should be warned about and ignored.
- Enum values should be case-insensitive; unknown enum values must warn and be ignored.

## Common variables

### General SDK

- `OTEL_SDK_DISABLED`: disable SDK when `true`
- `OTEL_SERVICE_NAME`: sets `service.name`; overrides `service.name` in `OTEL_RESOURCE_ATTRIBUTES`
- `OTEL_RESOURCE_ATTRIBUTES`: comma-separated resource attributes
- `OTEL_LOG_LEVEL`: SDK internal log level, default `info`
- `OTEL_PROPAGATORS`: default `tracecontext,baggage`
- `OTEL_TRACES_SAMPLER`: default `parentbased_always_on`
- `OTEL_TRACES_SAMPLER_ARG`: argument for selected sampler

### Exporter selection

- `OTEL_TRACES_EXPORTER`: default `otlp`
- `OTEL_METRICS_EXPORTER`: default `otlp`
- `OTEL_LOGS_EXPORTER`: default `otlp`

Valid common exporter values include `otlp`, `console`, `none`; traces also support `zipkin`, and metrics support `prometheus`. `logging` is deprecated.

### Batch processors

Span processor:

- `OTEL_BSP_SCHEDULE_DELAY`: default `5000`
- `OTEL_BSP_EXPORT_TIMEOUT`: default `30000`
- `OTEL_BSP_MAX_QUEUE_SIZE`: default `2048`
- `OTEL_BSP_MAX_EXPORT_BATCH_SIZE`: default `512`

LogRecord processor:

- `OTEL_BLRP_SCHEDULE_DELAY`: default `1000`
- `OTEL_BLRP_EXPORT_TIMEOUT`: default `30000`
- `OTEL_BLRP_MAX_QUEUE_SIZE`: default `2048`
- `OTEL_BLRP_MAX_EXPORT_BATCH_SIZE`: default `512`

### Limits

- `OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT`
- `OTEL_ATTRIBUTE_COUNT_LIMIT`: default `128`
- `OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT`: default `128`
- `OTEL_SPAN_EVENT_COUNT_LIMIT`: default `128`
- `OTEL_SPAN_LINK_COUNT_LIMIT`: default `128`
- `OTEL_LOGRECORD_ATTRIBUTE_COUNT_LIMIT`: default `128`

### Metrics

- `OTEL_METRICS_EXEMPLAR_FILTER`: default `trace_based`
- `OTEL_METRIC_EXPORT_INTERVAL`: default `60000`
- `OTEL_METRIC_EXPORT_TIMEOUT`: default `30000`

### Declarative config

- `OTEL_CONFIG_FILE`: path to SDK config file
- If set, all other SDK environment variables are ignored unless referenced through environment-variable substitution in the config file.

## Troubleshooting checklist

1. Check whether `OTEL_CONFIG_FILE` is set; it overrides flat SDK env vars.
2. Check empty values; they behave as unset.
3. Validate booleans: only `true` enables true behavior.
4. Validate enum spelling: unsupported values are ignored.
5. Check signal-specific exporter settings: traces, metrics, and logs are configured independently.
6. Confirm SDK/language support; implementations may choose whether to support these env vars.
