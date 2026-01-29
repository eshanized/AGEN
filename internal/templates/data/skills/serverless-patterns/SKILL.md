---
name: serverless-patterns
description: Lambda, Edge Functions, and serverless architecture patterns.
allowed-tools: Read, Write, Edit, Glob, Grep
---

# Serverless Patterns

> Principles for serverless and edge computing.
> **Pay for what you use. Scale to zero. Cold starts matter.**

---

## 1. Platform Selection

| Platform | Best For |
|----------|----------|
| **AWS Lambda** | Full AWS integration |
| **Vercel Functions** | Next.js, simple deploy |
| **Cloudflare Workers** | Edge, low latency |
| **Netlify Functions** | Jamstack |

### Decision

```
Need AWS services? → Lambda
Next.js app? → Vercel Functions
Global low latency? → Cloudflare Workers
Static site + API? → Netlify Functions
```

---

## 2. Cold Start Optimization

### Causes & Fixes

| Cause | Solution |
|-------|----------|
| **Large bundle** | Tree shaking, minimize deps |
| **Slow initialization** | Lazy loading, connection pooling |
| **Wrong runtime** | Use smaller runtimes |

### Techniques

```typescript
// Lazy import
let db: Database;
async function getDb() {
  if (!db) {
    const { createDb } = await import('./db');
    db = await createDb();
  }
  return db;
}

// Connection reuse (module-level)
const pool = createPool(); // Reused across invocations
```

---

## 3. Function Design

### Best Practices

| Practice | Why |
|----------|-----|
| **Single responsibility** | One function, one purpose |
| **Stateless** | No in-memory state between calls |
| **Idempotent** | Safe to retry |
| **Fast response** | Avoid timeouts |

### Handler Pattern

```typescript
// Good: Focused, fast
export async function handler(event) {
  const { id } = event.pathParameters;
  const item = await db.get(id);
  return { statusCode: 200, body: JSON.stringify(item) };
}

// Bad: Too much logic
export async function handler(event) {
  // 500 lines of code...
}
```

---

## 4. Database Connections

### The Problem

```
Cold start → New connection
1000 concurrent → 1000 connections 💀
```

### Solutions

| Solution | Platform |
|----------|----------|
| **Connection pooling** | PlanetScale, Neon, Supabase |
| **HTTP/REST API** | Database as service |
| **RDS Proxy** | AWS Lambda |

```typescript
// Use pooled connections
import { Pool } from '@neondatabase/serverless';

const pool = new Pool({ connectionString: process.env.DATABASE_URL });
```

---

## 5. Edge Functions

### When to Use

| Use Case | Edge? |
|----------|-------|
| **Auth/redirect** | ✅ Yes |
| **Personalization** | ✅ Yes |
| **API calls** | ⚠️ Maybe (if origin is far) |
| **Heavy compute** | ❌ No |
| **Database writes** | ❌ No (use origin) |

### Limitations

- Limited runtime (no Node.js APIs)
- Short timeouts (50ms-30s)
- Limited memory
- No persistent connections

---

## 6. Error Handling

### Retry Strategy

```typescript
// Wrap with retries
async function withRetry<T>(fn: () => Promise<T>, retries = 3): Promise<T> {
  for (let i = 0; i < retries; i++) {
    try {
      return await fn();
    } catch (e) {
      if (i === retries - 1) throw e;
      await sleep(100 * 2 ** i);
    }
  }
}
```

### Dead Letter Queues

```yaml
# AWS SAM
DeadLetterQueue:
  Type: AWS::SQS::Queue

MyFunction:
  DeadLetterConfig:
    TargetArn: !GetAtt DeadLetterQueue.Arn
```

---

## 7. Observability

### Essential Metrics

| Metric | Why |
|--------|-----|
| **Duration** | Performance |
| **Cold starts** | User experience |
| **Errors** | Reliability |
| **Concurrency** | Scaling |

### Structured Logging

```typescript
console.log(JSON.stringify({
  level: 'info',
  message: 'Request processed',
  requestId: context.awsRequestId,
  duration: Date.now() - start
}));
```

---

## ✅ Checklist

- [ ] **Bundle optimized?** (tree shaking, minimal deps)
- [ ] **Cold starts acceptable?** (measured, optimized)
- [ ] **Connections pooled?** (database, HTTP)
- [ ] **Timeouts set?** (function and downstream)
- [ ] **Retries implemented?** (with backoff)
- [ ] **Logging structured?** (JSON, request IDs)
- [ ] **Idempotent?** (safe to retry)

---

> **Remember:** Serverless means you don't manage servers, not that there are no servers. Understand the execution model.
