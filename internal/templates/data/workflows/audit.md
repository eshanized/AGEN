---
description: Comprehensive project audit for security, performance, accessibility, and code quality.
---

# /audit - Project Audit

$ARGUMENTS

---

## Purpose

Perform a multi-dimensional audit of the codebase using specialized agents.

---

## 🔴 CRITICAL RULES

1. **Use multiple agents** - This is a multi-agent task
2. **Report all findings** - Don't hide issues
3. **Prioritize by impact** - Critical issues first
4. **Provide remediation** - Not just problems, but solutions

---

## Task Flow

```
/audit [scope]

1. SECURITY AUDIT
   └── Invoke security-auditor agent
   
2. PERFORMANCE AUDIT
   └── Invoke performance-optimizer agent
   
3. ACCESSIBILITY AUDIT (if applicable)
   └── Check WCAG compliance
   
4. CODE QUALITY AUDIT
   └── Invoke test-engineer for coverage
   └── Check for code smells
   
5. SYNTHESIZE
   └── Combine all findings
   └── Prioritize by severity
```

---

## Agent Invocations

### Security

```
CONTEXT:
- Task: Security vulnerability assessment
- Focus: OWASP Top 10, dependency audit, secrets exposure
- Output: List of vulnerabilities with severity
```

### Performance

```
CONTEXT:
- Task: Performance analysis
- Focus: Bundle size, lighthouse scores, slow queries
- Output: Performance metrics and recommendations
```

---

## Output Format

```markdown
# 📊 Audit Report: [Project Name]

**Date:** [date]
**Scope:** [what was audited]

---

## 🔴 Critical Issues (0)

| Issue | Category | Location | Fix |
|-------|----------|----------|-----|

---

## 🟠 High Priority (0)

| Issue | Category | Location | Fix |
|-------|----------|----------|-----|

---

## 🟡 Medium Priority (0)

| Issue | Category | Location | Fix |
|-------|----------|----------|-----|

---

## 🟢 Low Priority (0)

| Issue | Category | Location | Fix |
|-------|----------|----------|-----|

---

## Summary

- Security: X issues (X critical)
- Performance: Score XX/100
- Accessibility: X issues
- Code Quality: XX% coverage

## Next Steps

1. [Action item]
2. [Action item]
```

---

## Usage Examples

```
/audit                    # Full audit
/audit security           # Security only
/audit performance        # Performance only
/audit --quick            # Critical issues only
```
