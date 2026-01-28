# Advanced Usage

This guide covers advanced AGEN features for power users.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `AGEN_NO_COLOR` | Disable colored output | `false` |
| `AGEN_DEBUG` | Enable verbose debug logging | `false` |
| `AGEN_CONFIG_DIR` | Override config directory location | System default |
| `AGEN_CACHE_DIR` | Override cache directory location | System default |

**Example:**
```bash
AGEN_DEBUG=true agen init --verbose
```

---

## Profile Management

Profiles save specific configurations (IDE + Agents + Skills) for reuse.

### Save and Load

```bash
# Save current config as profile
agen profile save my-web-stack

# Apply profile to new project
agen profile load my-web-stack
```

### Export and Share

```bash
# Export profile as JSON
agen profile export my-web-stack > web-stack.json

# Share with team, they import with:
agen profile import web-stack.json
```

See [Profiles](profiles.md) for detailed documentation.

---

## Custom Templates

Override built-in templates or add your own.

### Local Template Directory

Create a local template directory:

```bash
mkdir -p ~/.agen-templates/agents
mkdir -p ~/.agen-templates/skills
```

### Adding Custom Agent

Create `~/.agen-templates/agents/my-agent.md`:

```markdown
---
name: my-agent
description: My custom agent for specific tasks
tools: Read, Grep, Edit, Write
skills: clean-code
---

# My Custom Agent

You are a specialized agent for...
```

### Priority Order

AGEN loads templates in priority order:
1. Local overrides (`~/.agen-templates/`)
2. Installed plugins
3. Embedded templates (built-in)

First match wins.

---

## Conflict Resolution

When running `agen update`, AGEN handles conflicts intelligently.

### Behavior

| Scenario | Action |
|----------|--------|
| Template unchanged locally | Updated automatically |
| Template modified locally | Skipped (warning shown) |
| New template available | Added |
| `--force` flag used | All templates overwritten |

### Checking Conflicts

```bash
# See what would be updated
agen update --dry-run
```

---

## CI/CD Integration

### GitHub Actions

```yaml
name: AGEN Validation

on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install AGEN
        run: |
          curl -L https://github.com/eshanized/agen/releases/latest/download/agen_linux_amd64.tar.gz | tar xz
          sudo mv agen /usr/local/bin/
      
      - name: Validate Configuration
        run: agen team validate --strict
      
      - name: Run Verification
        run: agen verify
```

### GitLab CI

```yaml
agen-validate:
  stage: test
  script:
    - curl -L https://github.com/eshanized/agen/releases/latest/download/agen_linux_amd64.tar.gz | tar xz
    - ./agen team validate
    - ./agen verify
```

### Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

if command -v agen &> /dev/null; then
  if ! agen team validate --quiet 2>/dev/null; then
    echo "❌ Team validation failed"
    exit 1
  fi
fi

exit 0
```

---

## Multi-Project Setup

### Monorepo Configuration

For monorepos with multiple projects:

```
my-monorepo/
├── .agen-team.json       # Shared team config
├── packages/
│   ├── frontend/
│   │   └── .agent/       # Frontend-specific agents
│   ├── backend/
│   │   └── .agent/       # Backend-specific agents
│   └── shared/
│       └── .agent/       # Shared utilities
```

### Per-Project Initialization

```bash
# Initialize each package separately
cd packages/frontend
agen init --agents frontend-specialist,test-engineer

cd ../backend
agen init --agents backend-specialist,security-auditor
```

---

## Custom Verification Scripts

### Adding Project-Specific Scripts

Create scripts in `.agent/scripts/`:

```bash
mkdir -p .agent/scripts
```

Add `my-check.py`:

```python
#!/usr/bin/env python3
"""Custom verification script."""

import sys
import os

def main():
    errors = []
    
    # Your custom checks here
    if not os.path.exists("README.md"):
        errors.append("Missing README.md")
    
    if errors:
        print("❌ Errors found:")
        for e in errors:
            print(f"  • {e}")
        return 1
    
    print("✅ All checks passed")
    return 0

if __name__ == "__main__":
    sys.exit(main())
```

### Running Custom Scripts

```bash
python .agent/scripts/my-check.py
```

---

## Plugin Development

### Creating a Plugin

```bash
agen plugin create my-plugin --type bundle
```

This generates:

```
my-plugin/
├── plugin.json
├── agents/
│   └── my-agent.md
├── skills/
│   └── my-skill/
│       └── SKILL.md
└── README.md
```

### Testing Locally

```bash
# Install from local path
agen plugin install ./my-plugin

# Verify it's installed
agen plugin list

# Test in a project
agen init --agents my-agent
```

See [Plugin System](plugins.md) for detailed documentation.

---

## Remote Template Sources

### Adding Custom Sources

```bash
# Add a custom template repository
agen remote add company https://github.com/company/agen-templates

# Fetch templates
agen remote fetch

# List remotes
agen remote list
```

### Private Repositories

For private repos, configure authentication:

```bash
# Using GitHub token
export GITHUB_TOKEN=ghp_xxxxx
agen remote add private https://github.com/company/private-templates
```

---

## Debug Mode

### Verbose Output

```bash
agen init --verbose
```

### Debug Logging

```bash
AGEN_DEBUG=true agen init 2>&1 | tee debug.log
```

### Checking Configuration

```bash
# View current config
agen config list

# Check profiles
agen profile list

# Check plugins
agen plugin list
```

---

## Performance Tips

1. **Use `--dry-run`**: Preview changes before applying
2. **Save Profiles**: Avoid re-selecting agents each time
3. **Local Plugins**: Keep frequently-used custom templates as plugins
4. **Selective Updates**: Use `agen update` instead of full re-init

---

## Troubleshooting

### Common Issues

| Issue | Solution |
|-------|----------|
| IDE not detected | Use `--ide` flag explicitly |
| Templates not found | Check `agen list` for available templates |
| Update conflicts | Use `--force` or resolve manually |
| Plugin not working | Check `agen plugin info <name>` |

### Reset Configuration

```bash
# Reset global config
agen config reset

# Remove all profiles
rm -rf ~/.config/agen/profiles/

# Remove all plugins
rm -rf ~/.config/agen/plugins/
```
