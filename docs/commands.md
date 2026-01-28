# Command Reference

Comprehensive guide to all `agen` CLI commands, flags, and usage examples.

## Global Flags

These flags are available for all commands:

| Flag | Description |
|------|-------------|
| `-v, --verbose` | Enable verbose output for debugging |
| `--no-color` | Disable colored output (useful for scripts) |
| `--version` | Show version information |
| `-h, --help` | Show help for any command |

---

## Core Commands

### `agen init`

Initialize agent templates in the current or specified project directory. This is the main bootstrapping command.

**Usage:**
```bash
agen init [path] [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `-i, --ide string` | Force specific IDE (antigravity, cursor, windsurf, zed) |
| `-a, --agents strings` | Comma-separated list of agents to install |
| `-s, --skills strings` | Comma-separated list of skills to install |
| `-f, --force` | Overwrite existing files without prompting |
| `--dry-run` | Show what would be done without making changes |
| `--no-wizard` | Skip the interactive TUI wizard |

**Examples:**
```bash
# Interactive mode (wizard)
agen init

# Quick setup for a specific stack
agen init --ide cursor --agents frontend,backend --skills react,node --force

# Just verify what would happen
agen init --dry-run

# Initialize a specific directory
agen init /path/to/project --ide antigravity
```

---

### `agen list`

List all available templates that can be installed.

**Usage:**
```bash
agen list [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `-a, --agents` | Only list agents |
| `-s, --skills` | Only list skills |
| `-w, --workflows` | Only list workflows |
| `--json` | Output in JSON format |

**Examples:**
```bash
# List everything
agen list

# List only agents
agen list --agents

# Export as JSON for scripting
agen list --json > templates.json
```

---

### `agen status`

Check the installation status of AGEN in the current directory.

**Usage:**
```bash
agen status
```

**Output includes:**
- Detected IDE
- Installed agents
- Installed skills
- Configuration health

---

### `agen health`

Analyze the current project's configuration health with recommendations.

**Usage:**
```bash
agen health [path]
```

**Checks Performed:**
- **IDE Config**: Validates `.agent/`, `.cursorrules`, etc.
- **Version Status**: Checks if templates are up-to-date
- **Local Modifications**: Detects customized templates
- **Recommendations**: Suggests agents based on project type

**Example:**
```bash
agen health
# Shows health score and recommendations
```

---

### `agen search`

Perform a fuzzy search across all agents, skills, and workflows.

**Usage:**
```bash
agen search <query> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `-n, --limit int` | Maximum number of results (default 10) |

**Examples:**
```bash
agen search security    # Finds security-auditor, penetration-tester
agen search react       # Finds frontend-specialist, react-related skills
agen search test -n 20  # Show up to 20 results
```

---

## Update Commands

### `agen update`

Update installed templates to the latest version from the binary.

**Usage:**
```bash
agen update [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `-f, --force` | Overwrite local modifications |
| `--dry-run` | Show what files would be updated |

**Smart Updates:** AGEN respects local changes. Modified files are skipped unless `--force` is used.

---

### `agen upgrade`

Upgrade the `agen` binary itself to the latest version.

**Usage:**
```bash
agen upgrade
```

Downloads and installs the latest release from GitHub.

---

## Profile Commands

### `agen profile`

Manage configuration profiles for reuse across projects.

**Subcommands:**

| Command | Description |
|---------|-------------|
| `save <name>` | Save current config as a profile |
| `load <name>` | Apply a saved profile |
| `list` | List all saved profiles |
| `delete <name>` | Delete a saved profile |
| `export <name>` | Export profile as JSON (stdout) |
| `import <file>` | Import profile from JSON file |

**Examples:**
```bash
# Save current configuration
agen profile save frontend-stack

# Load profile into new project
agen profile load frontend-stack

# Export for sharing
agen profile export frontend-stack > frontend.json

# Import shared profile
agen profile import frontend.json
```

See [Profiles](profiles.md) for detailed documentation.

---

## Team Commands

### `agen team`

Manage team collaboration and shared configurations.

**Subcommands:**

| Command | Description |
|---------|-------------|
| `init <name>` | Initialize team configuration |
| `require <type> <name>` | Add required agent/skill |
| `remove <type> <name>` | Remove requirement |
| `lock <name> <version>` | Lock template version |
| `unlock <name>` | Remove version lock |
| `sync` | Sync project with team config |
| `validate` | Validate against team requirements |
| `config <key> <value>` | Modify team settings |

**Examples:**
```bash
# Initialize team config
agen team init my-team

# Require security-auditor for all team members
agen team require agent security-auditor

# Sync project with team requirements
agen team sync

# Validate project meets team requirements
agen team validate --strict
```

See [Team Collaboration](team.md) for detailed documentation.

---

## Plugin Commands

### `agen plugin`

Manage AGEN plugins.

**Subcommands:**

| Command | Description |
|---------|-------------|
| `install <source>` | Install plugin from GitHub/URL/path |
| `uninstall <name>` | Remove installed plugin |
| `list` | List installed plugins |
| `info <name>` | Show plugin details |
| `create <name>` | Create new plugin project |

**Examples:**
```bash
# Install from GitHub
agen plugin install github.com/user/agen-security-pack

# Install from local path
agen plugin install /path/to/my-plugin

# Create new plugin
agen plugin create my-plugin --type bundle
```

See [Plugin System](plugins.md) for detailed documentation.

---

## AI Commands

### `agen ai`

AI-powered features for intelligent suggestions.

**Subcommands:**

| Command | Description |
|---------|-------------|
| `suggest` | Analyze project and recommend agents |
| `explain <name>` | Get detailed explanation of agent/skill |
| `compose <name>` | Create custom agent from description |

**Examples:**
```bash
# Get agent suggestions for current project
agen ai suggest

# Learn about an agent
agen ai explain frontend-specialist

# Create a custom agent
agen ai compose my-reviewer --description "React code reviewer"
```

See [AI Features](ai-features.md) for detailed documentation.

---

## Verification Command

### `agen verify`

Run verification checks on your project.

**Usage:**
```bash
agen verify [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--verbose` | Show detailed output |
| `--security` | Only run security checks |
| `--lint` | Only run lint checks |

**Example:**
```bash
agen verify
# Runs security, lint, UX, and SEO checks
```

See [Verification](verification.md) for detailed documentation.

---

## Other Commands

### `agen create`

Create custom agents or skills.

**Usage:**
```bash
agen create <type> <name>
```

**Types:** `agent`, `skill`, `workflow`

**Example:**
```bash
agen create agent my-custom-agent
# Creates .agent/agents/my-custom-agent.md
```

---

### `agen remote`

Manage remote template sources.

**Subcommands:**

| Command | Description |
|---------|-------------|
| `add <name> <url>` | Add remote source |
| `remove <name>` | Remove remote source |
| `list` | List configured remotes |
| `fetch` | Fetch templates from remotes |

---

### `agen config`

Manage global AGEN configuration.

**Subcommands:**

| Command | Description |
|---------|-------------|
| `get <key>` | Get config value |
| `set <key> <value>` | Set config value |
| `list` | Show all config |
| `reset` | Reset to defaults |

**Example:**
```bash
agen config set default_ide cursor
agen config set auto_check_updates false
```

---

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | General error |
| `2` | Configuration error |
| `3` | Template not found |
| `4` | IDE not detected |
| `5` | Validation failed |

---

## Shell Completion

Generate shell completion scripts:

```bash
# Bash
agen completion bash > /etc/bash_completion.d/agen

# Zsh
agen completion zsh > "${fpath[1]}/_agen"

# Fish
agen completion fish > ~/.config/fish/completions/agen.fish
```
