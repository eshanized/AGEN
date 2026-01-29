---
description: Generate developer onboarding documentation for a codebase.
---

# /onboard - Developer Onboarding

$ARGUMENTS

---

## Purpose

Generate comprehensive onboarding documentation for new developers joining the project.

---

## Task Flow

```
/onboard

1. DISCOVER
   └── Invoke explorer-agent to map codebase
   └── Identify key directories and patterns
   
2. ANALYZE
   └── Determine tech stack
   └── Find configuration files
   └── Identify dependencies
   
3. GENERATE
   └── Invoke documentation-writer
   └── Create onboarding guide
   
4. VALIDATE
   └── Verify setup instructions work
   └── Test common workflows
```

---

## Output: docs/ONBOARDING.md

```markdown
# Developer Onboarding Guide

Welcome to [Project Name]! This guide will get you up and running.

---

## Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Node.js | 20.x | Runtime |
| pnpm | 8.x | Package manager |
| Docker | Latest | Services |

---

## Quick Start

```bash
# 1. Clone the repository
git clone [repo-url]
cd [project-name]

# 2. Install dependencies
pnpm install

# 3. Set up environment
cp .env.example .env.local

# 4. Start development
pnpm dev
```

---

## Project Structure

```
├── src/
│   ├── app/          # Next.js app router
│   ├── components/   # UI components
│   ├── lib/          # Utilities
│   └── api/          # API routes
├── tests/            # Test files
└── docs/             # Documentation
```

---

## Key Concepts

### [Concept 1]
[Explanation]

### [Concept 2]
[Explanation]

---

## Development Workflow

### Running Tests
```bash
pnpm test
```

### Building
```bash
pnpm build
```

### Creating a PR
1. Create branch: `git checkout -b feature/my-feature`
2. Make changes
3. Run tests: `pnpm test`
4. Push and create PR

---

## Common Tasks

### Adding a New Component
1. Create file in `src/components/`
2. Export from `index.ts`
3. Add tests in `tests/`

### Adding an API Endpoint
1. Create route in `src/app/api/`
2. Add validation
3. Add tests

---

## Troubleshooting

### "Module not found"
```bash
rm -rf node_modules && pnpm install
```

### "Port already in use"
```bash
lsof -i :3000
kill -9 [PID]
```

---

## Resources

- [Architecture Docs](./architecture.md)
- [API Reference](./api.md)
- [Contributing Guide](./CONTRIBUTING.md)

---

## Getting Help

- Slack: #[project-channel]
- Wiki: [wiki-url]
- Tech Lead: @[name]
```

---

## Usage Examples

```
/onboard                    # Full onboarding docs
/onboard --quick            # Quick start only
/onboard --for-contributors # Include contribution guide
```
