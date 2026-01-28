# Command Reference

Comprehensive guide to all `agen` CLI commands, flags, and usage examples.

## Global Flags

These flags are available for all commands:
- `-v, --verbose`: Enable verbose output for debugging.
- `--no-color`: Disable colored output (useful for scripts).
- `--version`: Show version information.

---

## `agen init`

Initialize agent templates in the current or specified project directory. This is the main bootstrapping command.

**Usage:**
```bash
agen init [path] [flags]
```

**Flags:**
- `-i, --ide string`: Force specific IDE configuration.
    - Supported values: `antigravity`, `cursor`, `windsurf`, `zed`
- `-a, --agents strings`: Comma-separated list of agents to install (e.g., `frontend,backend`).
- `-s, --skills strings`: Comma-separated list of skills to install (e.g., `docker,git`).
- `-f, --force`: Overwrite existing files without prompting.
- `--dry-run`: Show what would be done without making legitimate changes.
- `--no-wizard`: Skip the interactive TUI wizard even if no flags are provided (defaults to Antigravity if no IDE detected).

**Examples:**
```bash
# Interactive mode (wizard)
agen init

# Quick setup for a specific stack
agen init --ide cursor --agents frontend,backend --skills react,node --force

# Just verify what would happen
agen init --dry-run
```

---

## `agen list`

List all available templates that can be installed.

**Usage:**
```bash
agen list [flags]
```

**Flags:**
- `-a, --agents`: Only list agents.
- `-s, --skills`: Only list skills.
- `-w, --workflows`: Only list workflows.
- `--json`: Output in JSON format (useful for scripting/tools).

---

## `agen health`

Analyze the current project's configuration health. Checks for proper installation, file integrity, and missing recommended agents based on project type.

**Usage:**
```bash
agen health [path]
```

**Checks Performed:**
- **IDE Config**: Validates `.agent/`, `.cursorrules`, etc.
- **Version Status**: Checks if templates are up-to-date with the binary.
- **Local Modifications**: Detects if you've customized templates.
- **Recommendations**: Scans project files (`package.json`, `go.mod`, etc.) to recommend relevant agents (e.g., suggesting `mobile-developer` for React Native projects).

---

## `agen status`

Check the installation status of AGEN in the current directory. Similar to `health` but more focused on simple verification of installed components.

**Usage:**
```bash
agen status
```

---

## `agen search`

Perform a fuzzy search across all agents, skills, and workflows.

**Usage:**
```bash
agen search <query> [flags]
```

**Flags:**
- `-n, --limit int`: Maximum number of results to show (default 10).

**Examples:**
```bash
agen search security  # Finds security-auditor, penetration-tester, etc.
agen search react     # Finds frontend-specialist, react-skills
```

---

## `agen update`

Update installed templates to the latest version available in the `agen` binary. 

**Smart Updates**: AGEN attempts to respect your local changes. If a file has been modified locally, it will skip updating it unless you force it.

**Usage:**
```bash
agen update [flags]
```

**Flags:**
- `-f, --force`: Overwrite local modifications with the latest version.
- `--dry-run`: Show what files would be updated.

---

## `agen upgrade`

Upgrade the `agen` binary itself to the latest version from GitHub Releases.

**Usage:**
```bash
agen upgrade
```

---

## `agen verify`

Run a suite of verification scripts to ensure your project adheres to the standards defined in the templates.

**Usage:**
```bash
agen verify
```

---

## `agen profile`

Manage configuration profiles. Profiles allow you to save a specific combination of agents, skills, and IDE settings to reuse later.

**Subcommands:**

### `save`
Save current project configuration as a profile.
```bash
agen profile save <name>
```

### `load`
Apply a saved profile to the current project.
```bash
agen profile load <name>
```

### `list`
List all saved profiles.
```bash
agen profile list
```

### `delete`
Delete a saved profile.
```bash
agen profile delete <name>
```

### `export`
Export a profile configuration as JSON (stdout).
```bash
agen profile export <name> > my-profile.json
```

### `import`
Import a profile from a JSON file.
```bash
agen profile import <file>
```
