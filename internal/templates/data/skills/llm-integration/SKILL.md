---
name: llm-integration
description: LLM API patterns, prompt engineering, RAG pipelines, and AI feature integration principles.
allowed-tools: Read, Write, Edit, Glob, Grep
---

# LLM Integration

> Principles for integrating Large Language Models into applications.
> **Think first, pattern second.**

## 🎯 Selective Reading Rule

**Read ONLY sections relevant to your current task.**

---

## 1. LLM Provider Selection

| Provider | Best For | Considerations |
|----------|----------|----------------|
| **OpenAI (GPT-4)** | General purpose, mature | Cost, rate limits |
| **Anthropic (Claude)** | Complex reasoning, safety | Context length, cost |
| **Google (Gemini)** | Multimodal, Google integration | API stability |
| **Local (Ollama)** | Privacy, no API costs | Hardware requirements |

### Decision Framework

```
Privacy Critical? → Local model (Ollama, llama.cpp)
Complex Reasoning? → Claude/GPT-4
Cost Sensitive? → GPT-3.5/Claude Haiku/Local
Multimodal (images)? → GPT-4V/Gemini/Claude
```

---

## 2. API Best Practices

### Error Handling

| Error | Strategy |
|-------|----------|
| Rate Limit (429) | Exponential backoff with jitter |
| Timeout | Retry with shorter prompt |
| Context Overflow | Truncate or summarize |
| Invalid Response | Parse error handling, retry |

### Request Patterns

```typescript
// Structured output with retry
async function llmCall<T>(prompt: string, schema: z.Schema<T>): Promise<T> {
  for (let attempt = 0; attempt < 3; attempt++) {
    try {
      const response = await client.chat.completions.create({...});
      return schema.parse(JSON.parse(response.content));
    } catch (e) {
      if (attempt === 2) throw e;
      await sleep(exponentialBackoff(attempt));
    }
  }
}
```

---

## 3. Prompt Engineering Principles

### Effective Prompts

| Principle | Implementation |
|-----------|----------------|
| **Be Specific** | Clear, unambiguous instructions |
| **Show Examples** | Few-shot > zero-shot for consistency |
| **Structure Output** | Request JSON/structured format |
| **Set Constraints** | Length, format, tone |
| **Use System Prompts** | Persona and context |

### Prompt Templates

```
System: You are an expert {role} helping with {task}.

Context:
{relevant_context}

Task: {specific_task}

Requirements:
- {requirement_1}
- {requirement_2}

Output Format:
{format_specification}
```

---

## 4. RAG (Retrieval Augmented Generation)

### Architecture

```
User Query
    ↓
Embed Query → Vector Search → Top-K Chunks
    ↓                              ↓
    └──────── Context + Query ─────┘
                    ↓
              LLM Generation
                    ↓
               Response
```

### Best Practices

| Stage | Recommendation |
|-------|----------------|
| **Chunking** | 200-500 tokens, overlap 10-20% |
| **Embeddings** | OpenAI ada-002, Cohere, sentence-transformers |
| **Vector DB** | Pinecone (managed), pgvector (PostgreSQL) |
| **Retrieval** | Start with 3-5 chunks, adjust based on accuracy |
| **Reranking** | Cross-encoder for relevance |

---

## 5. Caching Strategies

### When to Cache

| Cache Level | Use Case |
|-------------|----------|
| **Exact Match** | Same prompt → same response |
| **Semantic** | Similar prompts → reuse response |
| **Embedding** | Expensive to generate, stable |
| **RAG Context** | Document chunks, infrequently updated |

---

## 6. Security Considerations

### Prompt Injection Prevention

| Risk | Mitigation |
|------|------------|
| **Direct Injection** | Input validation, sanitization |
| **Indirect Injection** | Don't include untrusted data in prompts |
| **Output Attacks** | Validate LLM output before use |

### Data Privacy

- Never send PII to external APIs without consent
- Log prompts and responses carefully
- Use local models for sensitive data

---

## 7. Cost Optimization

### Strategies

| Strategy | Savings |
|----------|---------|
| **Prompt Optimization** | Shorter = cheaper |
| **Model Selection** | GPT-3.5 when sufficient |
| **Caching** | Avoid duplicate calls |
| **Batch Processing** | Combine related requests |

---

## ✅ Checklist

Before integrating LLM:

- [ ] **Provider selected?** (based on requirements)
- [ ] **Error handling implemented?** (retries, fallbacks)
- [ ] **Rate limiting handled?** (backoff, queuing)
- [ ] **Output validated?** (schema, parsing)
- [ ] **Caching considered?** (where appropriate)
- [ ] **Costs estimated?** (tokens, calls)
- [ ] **Security reviewed?** (injection, privacy)

---

> **Remember:** LLMs are probabilistic, not deterministic. Always validate output and handle failures gracefully.
