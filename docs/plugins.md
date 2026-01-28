# Plugin System

AGEN supports plugins to extend its functionality with custom agents, skills, and workflows.

## Overview

Plugins can be installed from:

- **GitHub**: `github.com/user/repo`
- **Local Path**: `/path/to/plugin`
- **URL**: `https://example.com/plugin.zip`

---

## Installing Plugins

### From GitHub

```bash
# Install from GitHub repository
agen plugin install github.com/username/agen-security-pack
```

### From Local Path

```bash
# Install from local directory
agen plugin install /path/to/my-plugin
```

### From URL

```bash
# Install from zip URL
agen plugin install https://example.com/plugins/my-plugin.zip
```

---

## Managing Plugins

### List Installed Plugins

```bash
agen plugin list
```

Output:
```
📦 Installed Plugins

NAME                    VERSION   TYPE      AGENTS  SKILLS
security-pack           1.2.0     bundle    3       5
my-custom-agent         1.0.0     agent     1       0
react-patterns          2.1.0     skill     0       4
```

### Uninstall Plugin

```bash
agen plugin uninstall security-pack
```

### Get Plugin Info

```bash
agen plugin info security-pack
```

---

## Creating Plugins

### Plugin Structure

```
my-plugin/
├── plugin.json          # Plugin manifest (optional but recommended)
├── agents/              # Custom agents
│   └── my-agent.md
├── skills/              # Custom skills
│   └── my-skill/
│       └── SKILL.md
└── workflows/           # Custom workflows
    └── my-workflow.md
```

### Plugin Manifest (plugin.json)

```json
{
  "name": "my-plugin",
  "version": "1.0.0",
  "description": "My custom AGEN plugin",
  "author": "Your Name",
  "type": "bundle",
  "agents": ["my-agent"],
  "skills": ["my-skill"],
  "workflows": ["my-workflow"],
  "metadata": {
    "homepage": "https://github.com/user/my-plugin",
    "license": "MIT"
  }
}
```

### Plugin Types

| Type | Description |
|------|-------------|
| `agent` | Contains only agents |
| `skill` | Contains only skills |
| `workflow` | Contains only workflows |
| `bundle` | Contains multiple types |

---

## Creating a Plugin Project

Use the CLI to scaffold a new plugin:

```bash
# Create an agent plugin
agen plugin create my-agent --type agent

# Create a skill plugin
agen plugin create my-skill --type skill

# Create a bundle plugin
agen plugin create my-bundle --type bundle
```

This creates a starter structure:

```
my-agent/
├── plugin.json
├── agents/
│   └── my-agent.md
└── README.md
```

---

## Plugin Registry

Plugins are stored in your AGEN config directory:

| Platform | Location |
|----------|----------|
| Linux | `~/.config/agen/plugins/` |
| macOS | `~/Library/Application Support/agen/plugins/` |
| Windows | `%APPDATA%\agen\plugins\` |

A registry file (`registry.json`) tracks installed plugins.

---

## Publishing Plugins

### GitHub (Recommended)

1. Create a GitHub repository
2. Add plugin structure and `plugin.json`
3. Create releases with semantic versioning
4. Users install via `agen plugin install github.com/user/repo`

### Distribution

You can also distribute plugins as:

- ZIP archives
- Tarballs
- Direct file copies

---

## Agent Plugin Example

Create `agents/custom-reviewer.md`:

```markdown
---
name: custom-reviewer
description: Code review specialist with custom rules
tools: Read, Grep, Edit
skills: clean-code, code-review-checklist
---

# Custom Code Reviewer

You are a code reviewer following our team's specific standards...

## Review Checklist

- [ ] Variable naming follows team conventions
- [ ] Error handling is comprehensive
- [ ] Tests cover edge cases
```

---

## Skill Plugin Example

Create `skills/docker-expert/SKILL.md`:

```markdown
---
name: docker-expert
description: Advanced Docker and containerization patterns
version: 1.0
---

# Docker Expert Skill

## Container Best Practices

- Multi-stage builds for smaller images
- Non-root user execution
- Proper layer caching
- Health checks

## Common Commands

| Command | Purpose |
|---------|---------|
| `docker build` | Build image |
| `docker compose up` | Start services |
```

---

## Best Practices

1. **Version Your Plugins**: Use semantic versioning
2. **Include README**: Document usage and requirements
3. **Test Locally**: Use `agen plugin install /local/path` first
4. **Keep Dependencies Minimal**: Plugins should be self-contained
5. **Follow Naming Conventions**: Use kebab-case for names
