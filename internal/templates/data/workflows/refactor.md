---
description: Systematic refactoring with safety nets and verification.
---

# /refactor - Safe Refactoring

$ARGUMENTS

---

## Purpose

Perform systematic refactoring with proper testing and verification at each step.

---

## 🔴 CRITICAL RULES

1. **Tests first** - Never refactor without test coverage
2. **Small steps** - Each change should be independently valid
3. **Always green** - Tests must pass after every step
4. **No behavior change** - Refactoring != adding features

---

## Task Flow

```
/refactor [target] [goal]

1. ASSESS
   └── Identify what to refactor
   └── Check existing test coverage
   
2. PROTECT
   └── If coverage < 80%, write tests first
   └── Create characterization tests
   
3. PLAN
   └── Break into small, safe steps
   └── Define success criteria
   
4. EXECUTE
   └── Apply transformations one at a time
   └── Run tests after each change
   
5. VERIFY
   └── All tests passing
   └── Behavior unchanged
   └── Code is better
   
6. CLEANUP
   └── Remove dead code
   └── Update documentation
```

---

## Common Refactorings

| Goal | Approach |
|------|----------|
| **Extract function** | Identify reusable logic |
| **Rename** | Improve clarity |
| **Remove duplication** | DRY principle |
| **Simplify conditionals** | Extract, invert, combine |
| **Split large files** | Single responsibility |
| **Add types** | Incremental TypeScript |

---

## Output Format

```markdown
# Refactoring: [Target]

## Goal
[What we're trying to achieve]

## Before
```[language]
// Current code
```

## After
```[language]
// Improved code
```

## Steps Taken
1. [x] Extracted helper function `parseInput`
2. [x] Renamed `data` to `userPreferences`
3. [x] Removed duplicate validation logic

## Tests
- [x] All existing tests pass
- [x] Added 2 new tests for edge cases

## Metrics
| Metric | Before | After |
|--------|--------|-------|
| Lines | 247 | 189 |
| Complexity | 23 | 12 |
| Functions | 3 | 7 |
```

---

## Usage Examples

```
/refactor src/utils/parser.ts "reduce complexity"
/refactor OrderService "extract validation logic"
/refactor --dry-run          # Show plan, don't execute
```

---

## Key Principles

### The "Campsite Rule"
Leave code cleaner than you found it.

### The "Boy Scout Rule"
Make one small improvement every time you touch code.

### Safe Transformations
- Rename (IDE-assisted)
- Extract (method, variable, class)
- Move (to better location)
- Inline (remove unnecessary indirection)

### Dangerous Transformations (Require Extra Care)
- Change signature
- Combine classes
- Split classes
- Change inheritance
