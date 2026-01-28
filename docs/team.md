# Team Collaboration

AGEN supports team collaboration features that help standardize AI agent configurations across team members and projects.

## Overview

Team features allow you to:

- Define required agents and skills for all team members
- Sync configurations across projects
- Enforce consistent agent usage
- Lock specific template versions

---

## Team Configuration

### Initialize Team Config

```bash
agen team init my-team
```

This creates `.agen-team.json` in your project:

```json
{
  "name": "my-team",
  "version": "1.0.0",
  "created_at": "2026-01-29T00:00:00Z",
  "required_agents": [],
  "required_skills": [],
  "version_locks": {},
  "settings": {
    "enforce_agents": false,
    "enforce_skills": false,
    "allow_plugins": true
  }
}
```

---

## Managing Requirements

### Add Required Agents

```bash
# Require specific agents for all team members
agen team require agent security-auditor
agen team require agent test-engineer
```

### Add Required Skills

```bash
# Require specific skills
agen team require skill clean-code
agen team require skill testing-patterns
```

### Remove Requirements

```bash
agen team remove agent security-auditor
agen team remove skill clean-code
```

---

## Version Locking

Lock specific template versions to ensure consistency:

```bash
# Lock a template to a specific version
agen team lock frontend-specialist 1.2.0
agen team lock clean-code 2.0.0

# Unlock
agen team unlock frontend-specialist
```

---

## Team Settings

### Configuration Options

| Setting | Description | Default |
|---------|-------------|---------|
| `enforce_agents` | Fail validation if required agents missing | `false` |
| `enforce_skills` | Fail validation if required skills missing | `false` |
| `allow_plugins` | Allow third-party plugins | `true` |
| `default_ide` | Default IDE for new team members | `""` |
| `template_source` | Custom template repository | `""` |
| `sync_interval` | Auto-sync interval | `""` |

### Modify Settings

Edit `.agen-team.json` directly or use CLI:

```bash
# Enable enforcement
agen team config enforce_agents true
agen team config enforce_skills true
```

---

## Synchronization

### Sync Project with Team Config

```bash
agen team sync
```

This will:
1. Check for missing required agents
2. Check for missing required skills
3. Install any missing components
4. Report what was added/updated

### Sync Results

```
🔄 Syncing with team configuration...

✅ Added:
  • agents/security-auditor.md
  • skills/testing-patterns/

⚠️ Updated:
  • agents/frontend-specialist.md (version locked to 1.2.0)

📊 Summary:
  • 2 agents added
  • 1 agent updated
  • 0 errors
```

---

## Validation

### Validate Project Against Team Config

```bash
agen team validate
```

### Validation Output

```
🔍 Validating against team configuration...

❌ Missing Required Agents:
  • security-auditor

❌ Missing Required Skills:
  • testing-patterns

⚠️ Warnings:
  • frontend-specialist version differs from lock (1.1.0 vs 1.2.0)

Result: INVALID (2 missing, 1 warning)
```

---

## Configuration File Format

### Complete Example

```json
{
  "name": "acme-engineering",
  "version": "1.0.0",
  "created_at": "2026-01-01T00:00:00Z",
  "required_agents": [
    "frontend-specialist",
    "backend-specialist",
    "security-auditor",
    "test-engineer"
  ],
  "required_skills": [
    "clean-code",
    "testing-patterns",
    "api-patterns"
  ],
  "version_locks": {
    "frontend-specialist": "1.2.0",
    "clean-code": "2.0.0"
  },
  "settings": {
    "enforce_agents": true,
    "enforce_skills": true,
    "allow_plugins": true,
    "default_ide": "antigravity",
    "template_source": "",
    "sync_interval": "weekly"
  },
  "members": [
    {
      "name": "John Doe",
      "email": "john@example.com",
      "role": "lead",
      "joined_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

---

## Workflow Integration

### CI/CD Integration

Add team validation to your CI pipeline:

```yaml
# .github/workflows/ci.yml
- name: Validate AGEN Config
  run: agen team validate --strict
```

### Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

if ! agen team validate --quiet; then
  echo "Team configuration validation failed"
  exit 1
fi
```

---

## Best Practices

1. **Start Small**: Begin with essential agents only
2. **Discuss Changes**: Team-wide agent changes affect everyone
3. **Version Control**: Commit `.agen-team.json` to your repository
4. **Regular Syncs**: Run `agen team sync` after pulling changes
5. **Lock Carefully**: Only lock versions when necessary for stability
