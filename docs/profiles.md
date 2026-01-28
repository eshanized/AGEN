# Profile Management

Profiles allow you to save and reuse agent configurations across projects. Think of them as "presets" for your AI development environment.

## Overview

A profile captures:

- Selected IDE
- Installed agents
- Installed skills
- Installed workflows
- Timestamp metadata

---

## Saving Profiles

### Save Current Configuration

```bash
agen profile save frontend-stack
```

This saves the current project's configuration as a named profile.

### What Gets Saved

The profile captures:
- Which IDE format is being used
- List of installed agents
- List of installed skills
- List of installed workflows

---

## Loading Profiles

### Apply Profile to Project

```bash
agen profile load frontend-stack
```

This installs all agents and skills from the saved profile.

### Load with IDE Override

```bash
agen profile load frontend-stack --ide cursor
```

---

## Listing Profiles

### View All Profiles

```bash
agen profile list
```

Output:
```
📋 Saved Profiles

NAME              IDE           AGENTS  SKILLS  CREATED
frontend-stack    antigravity   4       8       2026-01-15
backend-api       cursor        3       5       2026-01-10
fullstack         antigravity   6       12      2026-01-05
security-audit    windsurf      2       4       2026-01-01
```

---

## Deleting Profiles

### Remove a Profile

```bash
agen profile delete old-profile
```

---

## Exporting and Importing

### Export Profile to JSON

```bash
# Export to stdout
agen profile export frontend-stack

# Export to file
agen profile export frontend-stack > frontend-stack.json
```

### Import Profile from JSON

```bash
agen profile import frontend-stack.json
```

---

## Profile File Format

Profiles are stored as JSON files:

```json
{
  "name": "frontend-stack",
  "ide": "antigravity",
  "agents": [
    "frontend-specialist",
    "test-engineer",
    "performance-optimizer",
    "seo-specialist"
  ],
  "skills": [
    "clean-code",
    "nextjs-react-expert",
    "tailwind-patterns",
    "testing-patterns",
    "performance-profiling",
    "seo-fundamentals"
  ],
  "workflows": [
    "orchestrate",
    "deploy"
  ],
  "created_at": "2026-01-15T10:30:00Z",
  "modified_at": "2026-01-20T14:15:00Z"
}
```

---

## Profile Storage Location

Profiles are stored in your AGEN config directory:

| Platform | Location |
|----------|----------|
| Linux | `~/.config/agen/profiles/` |
| macOS | `~/Library/Application Support/agen/profiles/` |
| Windows | `%APPDATA%\agen\profiles\` |

Each profile is saved as `<name>.json`.

---

## Common Profile Examples

### Frontend Development

```bash
agen init --agents frontend-specialist,test-engineer,seo-specialist \
          --skills clean-code,nextjs-react-expert,tailwind-patterns
agen profile save frontend-dev
```

### Backend API

```bash
agen init --agents backend-specialist,security-auditor,database-architect \
          --skills api-patterns,nodejs-best-practices,database-design
agen profile save backend-api
```

### Full Stack

```bash
agen init --agents frontend-specialist,backend-specialist,devops-engineer,test-engineer \
          --skills clean-code,nextjs-react-expert,api-patterns,deployment-procedures
agen profile save fullstack
```

### Security Audit

```bash
agen init --agents security-auditor,penetration-tester \
          --skills vulnerability-scanner,red-team-tactics
agen profile save security
```

### Mobile Development

```bash
agen init --agents mobile-developer,test-engineer,performance-optimizer \
          --skills mobile-design,testing-patterns,performance-profiling
agen profile save mobile
```

---

## Sharing Profiles

### Export for Team

```bash
# Export profile
agen profile export frontend-stack > frontend-stack.json

# Share the file with team members
# They import with:
agen profile import frontend-stack.json
```

### Version Control

You can commit profile JSON files to a repository:

```
my-team-profiles/
├── frontend.json
├── backend.json
├── fullstack.json
└── security.json
```

Team members clone and import as needed.

---

## Best Practices

1. **Descriptive Names**: Use clear, descriptive profile names
2. **Keep Updated**: Recreate profiles when you find better agent combinations
3. **Document Profiles**: Note what each profile is intended for
4. **Share With Team**: Export and share successful configurations
5. **Project-Specific**: Create profiles for different project types
