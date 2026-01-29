---
description: Production incident response and post-mortem workflow.
---

# /incident - Incident Response

$ARGUMENTS

---

## Purpose

Structured approach to handling production incidents and creating post-mortems.

---

## 🔴 CRITICAL RULES

1. **Triage first** - Understand scope and impact
2. **Communicate** - Keep stakeholders informed
3. **Mitigate before root cause** - Stop the bleeding
4. **Document everything** - For post-mortem

---

## Severity Levels

| Level | Impact | Response Time |
|-------|--------|---------------|
| **SEV1** | Complete outage | Immediate |
| **SEV2** | Major functionality broken | < 15 min |
| **SEV3** | Minor functionality broken | < 1 hour |
| **SEV4** | Cosmetic/Low impact | Next business day |

---

## Task Flow

```
/incident [description]

1. TRIAGE
   └── Assess severity
   └── Identify affected systems
   
2. COMMUNICATE
   └── Alert stakeholders
   └── Set up war room if needed
   
3. INVESTIGATE
   └── Invoke debugger agent
   └── Check logs, metrics, traces
   
4. MITIGATE
   └── Invoke devops-engineer if rollback needed
   └── Apply temporary fix
   
5. RESOLVE
   └── Implement proper fix
   └── Verify resolution
   
6. POST-MORTEM
   └── Document incident
   └── Identify improvements
```

---

## Investigation Checklist

### Quick Checks
- [ ] What changed recently? (deploys, config)
- [ ] What do errors say? (logs, exceptions)
- [ ] What do metrics show? (latency, error rate)
- [ ] Is it widespread or isolated?

### System Checks
- [ ] Application logs
- [ ] Database metrics
- [ ] External dependency status
- [ ] Infrastructure metrics

---

## Output: Post-Mortem Template

```markdown
# Incident Post-Mortem: [Title]

**Date:** [date]
**Duration:** [start] - [end] ([X hours/minutes])
**Severity:** SEV[X]
**Author:** [name]

---

## Summary

[One paragraph summary of what happened]

---

## Timeline

| Time (UTC) | Event |
|------------|-------|
| 14:23 | Monitoring alert fired |
| 14:25 | On-call acknowledged |
| 14:30 | Root cause identified |
| 14:45 | Rollback initiated |
| 14:50 | Service restored |

---

## Impact

- Users affected: [number]
- Revenue impact: [estimate]
- Data loss: [none/details]

---

## Root Cause

[Detailed explanation of why this happened]

---

## Resolution

[What was done to fix it]

---

## Action Items

| Action | Owner | Due Date | Status |
|--------|-------|----------|--------|
| Add monitoring for X | @dev | 2025-02-05 | TODO |
| Improve rollback process | @ops | 2025-02-10 | TODO |

---

## What Went Well

- Quick detection
- Effective communication

## What Went Wrong

- Delayed rollback
- Missing runbook

---

## Lessons Learned

1. [Lesson 1]
2. [Lesson 2]
```

---

## Usage Examples

```
/incident "API returning 500s"
/incident --severity 2 "Payment processing failed"
/incident --post-mortem    # Generate post-mortem template
```
