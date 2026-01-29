---
name: monorepo-patterns
description: Turborepo, Nx, and workspace management patterns for monorepo development.
allowed-tools: Read, Write, Edit, Glob, Grep
---

# Monorepo Patterns

> Principles for managing monorepo projects.
> **Shared code, independent deployments.**

---

## 1. Tool Selection

| Tool | Best For |
|------|----------|
| **Turborepo** | Simple, fast, Vercel projects |
| **Nx** | Enterprise, plugins, advanced caching |
| **pnpm workspaces** | Simple dependency management |
| **Lerna** | Legacy, npm publishing |

### Decision

```
Simple monorepo? → Turborepo
Enterprise/Angular? → Nx
Just workspaces? → pnpm/npm/yarn workspaces
```

---

## 2. Directory Structure

### Standard Layout

```
monorepo/
├── apps/
│   ├── web/
│   ├── mobile/
│   └── api/
├── packages/
│   ├── ui/           # Shared components
│   ├── utils/        # Shared utilities
│   └── config/       # Shared configs
├── turbo.json        # Turborepo config
├── package.json      # Root package
└── pnpm-workspace.yaml
```

---

## 3. Dependency Management

### Workspace Protocol

```json
// packages/web/package.json
{
  "dependencies": {
    "@repo/ui": "workspace:*",
    "@repo/utils": "workspace:*"
  }
}
```

### Shared Dependencies

```json
// Root package.json
{
  "devDependencies": {
    "typescript": "^5.0.0",  // Shared tooling
    "eslint": "^8.0.0"
  }
}
```

---

## 4. Turborepo Configuration

### turbo.json

```json
{
  "$schema": "https://turbo.build/schema.json",
  "globalDependencies": ["**/.env.*local"],
  "pipeline": {
    "build": {
      "dependsOn": ["^build"],
      "outputs": ["dist/**", ".next/**"]
    },
    "test": {
      "dependsOn": ["build"],
      "inputs": ["src/**", "test/**"]
    },
    "lint": {},
    "dev": {
      "cache": false,
      "persistent": true
    }
  }
}
```

### Key Concepts

| Concept | Meaning |
|---------|---------|
| `^build` | Run dependencies' build first |
| `outputs` | What to cache |
| `inputs` | What triggers rebuild |
| `persistent` | Long-running tasks |

---

## 5. Shared Configurations

### TypeScript

```json
// packages/config/tsconfig.base.json
{
  "compilerOptions": {
    "strict": true,
    "esModuleInterop": true,
    "moduleResolution": "bundler"
  }
}

// apps/web/tsconfig.json
{
  "extends": "@repo/config/tsconfig.base.json",
  "include": ["src"]
}
```

### ESLint

```javascript
// packages/config/eslint.config.js
module.exports = {
  extends: ["next", "turbo"],
  rules: {
    // Shared rules
  }
};
```

---

## 6. Internal Packages

### Package Setup

```json
// packages/ui/package.json
{
  "name": "@repo/ui",
  "version": "0.0.0",
  "private": true,
  "main": "./src/index.ts",
  "types": "./src/index.ts",
  "exports": {
    ".": "./src/index.ts",
    "./button": "./src/button.tsx"
  }
}
```

### Transpilation

| Strategy | When |
|----------|------|
| **Just-in-time** | Consumer transpiles (Next.js) |
| **Pre-built** | Build before publish |

---

## 7. CI/CD Optimization

### Affected Commands

```bash
# Only build changed packages
turbo run build --filter=[HEAD^1]

# Build specific package and deps
turbo run build --filter=web...
```

### Remote Caching

```bash
# Enable Vercel Remote Cache
npx turbo login
npx turbo link
```

---

## ✅ Checklist

- [ ] **Workspace protocol used?** (`workspace:*`)
- [ ] **Shared configs extracted?** (TypeScript, ESLint)
- [ ] **Build pipeline correct?** (`^` for dependencies)
- [ ] **Outputs specified?** (for caching)
- [ ] **Remote caching enabled?** (CI optimization)
- [ ] **Clear ownership?** (CODEOWNERS per package)

---

> **Remember:** Monorepos are about sharing code while maintaining independence. Get the boundaries right.
