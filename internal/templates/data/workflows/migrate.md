---
description: Technology migration workflow for upgrading or replacing frameworks and languages.
---

# /migrate - Technology Migration

$ARGUMENTS

---

## Purpose

Safely migrate from one technology to another with minimal disruption.

---

## 🔴 CRITICAL RULES

1. **Incremental changes** - Never big bang
2. **Tests first** - Ensure coverage before migration
3. **Parallel running** - Old and new coexist
4. **Rollback plan** - Always have an exit

---

## Common Migrations

| From | To | Workflow |
|------|-----|----------|
| JavaScript | TypeScript | Incremental typing |
| REST | GraphQL | Facade pattern |
| Monolith | Microservices | Strangler fig |
| Class components | Hooks | Component by component |
| Webpack | Vite | Build system swap |

---

## Task Flow

```
/migrate [from] to [to]

1. ANALYSIS
   └── Invoke explorer-agent to map codebase
   └── Identify dependencies and impact
   
2. PLANNING
   └── Invoke project-planner
   └── Create migration plan with phases
   
3. PREPARATION
   └── Ensure test coverage
   └── Set up both systems in parallel
   
4. EXECUTION
   └── Migrate incrementally
   └── Verify after each step
   
5. CLEANUP
   └── Remove old code
   └── Update documentation
```

---

## Output Format

### Migration Plan

```markdown
# Migration: [From] → [To]

## Phase 1: Preparation
- [ ] Add TypeScript configuration
- [ ] Ensure test coverage > 80%
- [ ] Create compatibility layer

## Phase 2: Gradual Migration
- [ ] Convert utils/ directory
- [ ] Convert components/ directory
- [ ] Convert pages/ directory

## Phase 3: Cleanup
- [ ] Remove JavaScript files
- [ ] Update build scripts
- [ ] Update documentation

## Rollback Plan
[How to revert at each phase]
```

---

## Usage Examples

```
/migrate js to typescript
/migrate rest to graphql
/migrate cra to vite
/migrate pages to app router
```

---

## Best Practices

### JavaScript → TypeScript

1. Add `tsconfig.json` with `allowJs: true`
2. Rename files one at a time (`.js` → `.ts`)
3. Fix type errors incrementally
4. Enable strict mode last

### REST → GraphQL

1. Create GraphQL schema matching REST
2. Add GraphQL endpoint alongside REST
3. Migrate clients one by one
4. Deprecate REST endpoints

### Webpack → Vite

1. Add Vite config
2. Run both build systems
3. Compare outputs
4. Switch dev to Vite first
5. Switch build after validation
