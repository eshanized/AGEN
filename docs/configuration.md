# Configuration

AGEN is designed to be "zero-config" for most users, deriving its context from your project structure. However, there are ways to customize its behavior.

## Global Configuration

AGEN stores global data in your system's user config directory:

- **Linux**: `~/.config/agen/`
- **macOS**: `~/Library/Application Support/agen/`
- **Windows**: `%APPDATA%\agen\`

### Profiles
Saved profiles are stored in the `profiles/` subdirectory as JSON files. You can manually edit these if needed, though using the `agen profile` command is recommended.

## Project Configuration

Once initialized, AGEN's configuration lives inside your project.

### Antigravity (`.agent/`)
This is the most granular configuration. You can directly edit any file in `.agent/agents/` or `.agent/skills/`.
- **Rules**: `.agent/rules/` contains global rules applied to all agents.

### Cursor (`.cursorrules`)
This is a single generated file. **Warning**: If you edit this file manually, `agen update` might overwrite your changes unless you are careful.
- **Tip**: AGEN adds comments to sections. Try to keep your custom rules outside the managed blocks if possible, or use the `agen update` conflict resolution prompts.

## Environment Variables

| Variable | Description |
|----------|-------------|
| `AGEN_NO_COLOR` | Set to `true` to disable colored output. |
| `AGEN_DEBUG` | Set to `true` to enable verbose debug logging (equivalent to `--verbose`). |

## Custom Templates (Advanced)

You can maintain your own library of templates that override the built-in ones.

1. Create `~/.agen-templates/` directory.
2. Create `agents/` and `skills/` subdirectories.
3. Add your Markdown files there.

AGEN will prioritize these local templates over the embedded ones during `init` and `update`.
