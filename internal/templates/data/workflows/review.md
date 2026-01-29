---
description: Multi-agent code review for thorough PR analysis.
---

# /review - Code Review

$ARGUMENTS

---

## Purpose

Perform a comprehensive code review using multiple specialized agents.

---

## 🔴 CRITICAL RULES

1. **Be constructive** - Suggest improvements, don't just criticize
2. **Prioritize feedback** - Not all issues are equal
3. **Context matters** - Understand the PR's purpose
4. **Automate what's automatable** - Focus on creative review

---

## Task Flow

```
/review [files/PR]

1. UNDERSTAND
   └── Read PR description / commit messages
   └── Understand the goal
   
2. SECURITY REVIEW
   └── Invoke security-auditor
   └── Check auth, inputs, secrets
   
3. ARCHITECTURE REVIEW
   └── Check patterns, structure
   └── Look for anti-patterns
   
4. DOMAIN REVIEW
   └── Invoke appropriate specialist
   └── Check best practices
   
5. TEST REVIEW
   └── Invoke test-engineer
   └── Check coverage, quality
   
6. SYNTHESIZE
   └── Combine all feedback
   └── Prioritize by importance
```

---

## Review Categories

| Category | Agent | Focus |
|----------|-------|-------|
| **Security** | security-auditor | Vulnerabilities, auth |
| **Performance** | performance-optimizer | Efficiency, complexity |
| **Quality** | test-engineer | Tests, coverage |
| **Architecture** | tech-lead | Patterns, structure |
| **Domain** | backend/frontend/etc. | Best practices |

---

## Output Format

```markdown
# 📝 Code Review: [PR Title / Files]

## Summary
[Brief overview of the changes and review verdict]

---

## 🔴 Must Fix

### 1. [Issue Title]
**File:** `path/to/file.ts:42`
**Category:** Security

[Description of issue]

**Suggestion:**
```diff
- badCode();
+ goodCode();
```

---

## 🟡 Should Consider

### 1. [Suggestion Title]
**File:** `path/to/file.ts:78`
**Category:** Performance

[Description and reasoning]

---

## 🟢 Nitpicks

- Line 23: Consider using `const` instead of `let`
- Line 45: Typo in comment

---

## ✅ What's Good

- Clean separation of concerns
- Good test coverage
- Clear variable naming

---

## Verdict

[ ] ✅ Approve
[ ] 🔄 Request Changes
[ ] ❓ Questions/Discussion Needed
```

---

## Usage Examples

```
/review                     # Review staged changes
/review src/components/     # Review specific directory
/review --security-only     # Security focus only
/review --quick             # Only critical issues
```

---

## Review Checklist

### Always Check
- [ ] No secrets in code
- [ ] Input validation present
- [ ] Error handling exists
- [ ] Tests cover new code
- [ ] No obvious performance issues

### For PRs
- [ ] PR description clear
- [ ] Commits well-organized
- [ ] Breaking changes documented
- [ ] Migration guide if needed
