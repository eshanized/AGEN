---
name: observability-patterns
description: OpenTelemetry, distributed tracing, metrics, and logging patterns.
allowed-tools: Read, Write, Edit, Glob, Grep
---

# Observability Patterns

> Principles for understanding production systems.
> **Logs tell you what. Metrics tell you how much. Traces tell you why.**

---

## 1. Three Pillars

| Pillar | Purpose | Tools |
|--------|---------|-------|
| **Logs** | What happened | ELK, Loki, CloudWatch |
| **Metrics** | Measurements over time | Prometheus, Datadog |
| **Traces** | Request flow | Jaeger, Tempo, X-Ray |

### When to Use Each

| Question | Use |
|----------|-----|
| "What error occurred?" | Logs |
| "Is latency increasing?" | Metrics |
| "Why was this request slow?" | Traces |

---

## 2. Structured Logging

### Format

```json
{
  "timestamp": "2025-01-29T12:00:00Z",
  "level": "error",
  "message": "Failed to process order",
  "service": "order-service",
  "trace_id": "abc123",
  "span_id": "def456",
  "order_id": "order-789",
  "error": "Database connection timeout"
}
```

### Best Practices

| Do | Don't |
|----|-------|
| JSON format | Plain text in production |
| Include context (IDs) | Log sensitive data |
| Log at boundaries | Log everything |
| Consistent field names | Different names per service |

---

## 3. Metrics

### Metric Types

| Type | Use | Example |
|------|-----|---------|
| **Counter** | Cumulative count | `http_requests_total` |
| **Gauge** | Current value | `memory_usage_bytes` |
| **Histogram** | Distribution | `request_duration_seconds` |
| **Summary** | Percentiles | `request_latency_p99` |

### Golden Signals

| Signal | Measures |
|--------|----------|
| **Latency** | Time to serve request |
| **Traffic** | Demand on system |
| **Errors** | Failed requests |
| **Saturation** | Resource utilization |

### Naming Convention

```
# Prometheus style
http_requests_total{method="GET", status="200"}
http_request_duration_seconds{endpoint="/api/users"}
```

---

## 4. Distributed Tracing

### OpenTelemetry Setup

```typescript
import { trace } from '@opentelemetry/api';

const tracer = trace.getTracer('my-service');

async function handleRequest(req) {
  return tracer.startActiveSpan('handle-request', async (span) => {
    try {
      span.setAttribute('user.id', req.userId);
      
      const result = await processOrder(req.orderId);
      
      span.setStatus({ code: SpanStatusCode.OK });
      return result;
    } catch (error) {
      span.setStatus({ 
        code: SpanStatusCode.ERROR, 
        message: error.message 
      });
      throw error;
    } finally {
      span.end();
    }
  });
}
```

### Span Attributes

| Attribute | Example |
|-----------|---------|
| `http.method` | GET, POST |
| `http.url` | /api/users |
| `http.status_code` | 200, 500 |
| `db.system` | postgresql |
| `db.statement` | SELECT * FROM... |

---

## 5. Alerting

### Alert Design

| Principle | Implementation |
|-----------|----------------|
| **Actionable** | Someone should do something |
| **Relevant** | Affects users or SLOs |
| **Concise** | Clear what's wrong |
| **Routed** | Right team notified |

### SLO-Based Alerts

```yaml
# Alert on error budget burn
groups:
- name: slo
  rules:
  - alert: HighErrorRate
    expr: |
      sum(rate(http_requests_total{status=~"5.."}[5m]))
      / sum(rate(http_requests_total[5m])) > 0.01
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "Error rate above 1% SLO"
```

---

## 6. Context Propagation

### HTTP Headers

| Standard | Header |
|----------|--------|
| **W3C Trace Context** | `traceparent`, `tracestate` |
| **B3** | `X-B3-TraceId`, `X-B3-SpanId` |

### Correlation IDs

```typescript
// Generate at edge
const correlationId = req.headers['x-request-id'] || generateId();

// Pass through all services
res.setHeader('x-request-id', correlationId);

// Include in all logs
logger.info({ correlationId, message: 'Processing request' });
```

---

## ✅ Checklist

- [ ] **Structured logs?** (JSON, consistent fields)
- [ ] **Metrics exported?** (golden signals)
- [ ] **Traces instrumented?** (OpenTelemetry)
- [ ] **Context propagated?** (trace IDs, correlation IDs)
- [ ] **Alerts actionable?** (SLO-based)
- [ ] **Dashboards created?** (key services)
- [ ] **Sensitive data excluded?** (PII, tokens)

---

> **Remember:** Observability is about asking arbitrary questions of your system. Instrument for understanding, not just monitoring.
