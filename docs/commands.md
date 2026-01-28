# Command Reference

## `agen init`

Initializes AGEN in the current directory.

```bash
agen init [flags]
```

**Flags:**
- `-i, --ide string`: Force specific IDE (antigravity, cursor, windsurf, zed)
- `-a, --agents strings`: Comma-separated list of agents to install
- `-s, --skills strings`: Comma-separated list of skills to install
- `-f, --force`: Overwrite existing files without prompting
- `--dry-run`: Show what would be done without making changes
- `--no-wizard`: Skip interactive wizard

## `agen list`

Lists all available agents, skills, and workflows.

```bash
agen list
```

## `agen health`

Checks the health of the installed configuration.

```bash
agen health
```

**Checks performed:**
- IDE configuration validity
- Missing files
- Broken references

## `agen verify`

Runs verification scripts to ensure standards.

```bash
agen verify
```

## `agen search`

Fuzzy search for templates.

```bash
agen search <query>
```

## `agen update`

Updates installed templates to the latest version.

```bash
agen update
```

## `agen upgrade`

Upgrades the AGEN binary itself.

```bash
agen upgrade
```
