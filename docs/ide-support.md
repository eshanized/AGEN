# IDE Support

AGEN supports **12 IDEs and AI coding tools** by adapting its internal template format to the specific requirements of each environment.

## Quick Reference

| IDE/Tool | Config File | Auto-detection |
|----------|-------------|----------------|
| [Antigravity](#antigravity) | `.agent/` folder | ✅ `.agent/` directory |
| [Cursor](#cursor) | `.cursorrules` | ✅ `.cursorrules` file |
| [Windsurf](#windsurf) | `.windsurfrules` | ✅ `.windsurfrules` file |
| [Cline](#cline) | `.clinerules` | ✅ `.clinerules` file |
| [Continue](#continue) | `.continue/` | ✅ `.continue/` directory |
| [Zed](#zed) | `.zed/` folder | ✅ `.zed/` directory |
| [JetBrains](#jetbrains) | `.idea/` | ✅ `.idea/` directory |
| [Neovim](#neovim) | `.nvim/` | ✅ `.nvim/` or `.nvim.lua` |
| [Emacs](#emacs) | `.dir-locals.el` | ✅ `.dir-locals.el` file |
| [Aider](#aider) | `.aider.conf.yml` | ✅ `.aider.conf.yml` |
| [Claude Code](#claude-code) | `CLAUDE.md` | ✅ `CLAUDE.md` file |
| [Copilot Workspace](#github-copilot-workspace) | `.github/copilot-instructions.md` | ✅ `.github/` |

---

## VS Code Extensions

### Cursor

Cursor uses a single rules file to define AI behavior.

- **Detection**: Checks for `.cursorrules`
- **Format**: Consolidated `.cursorrules` file with agents, skills, and workflows
- **Install**: `agen init --ide cursor`

### Windsurf

Windsurf follows a similar pattern to Cursor.

- **Detection**: Checks for `.windsurfrules`
- **Format**: Single `.windsurfrules` file containing all context
- **Install**: `agen init --ide windsurf`

### Cline

Cline (formerly Roo-Cline) is a popular VS Code AI coding extension.

- **Detection**: Checks for `.clinerules` or `.cline/`
- **Format**: Single `.clinerules` file with project rules
- **Install**: `agen init --ide cline`
- **Website**: [github.com/cline/cline](https://github.com/cline/cline)

### Continue

Continue is an open-source AI coding assistant for VS Code.

- **Detection**: Checks for `.continue/` or `.continuerules`
- **Format**: `.continue/config.json` + `.continuerules`
- **Install**: `agen init --ide continue`
- **Website**: [continue.dev](https://continue.dev)

---

## Desktop IDEs

### Antigravity

This is the native format for AGEN (Google Gemini/Claude).

- **Detection**: Checks for `.agent/` directory or `GEMINI.md`
- **Structure**:
  ```text
  .agent/
  ├── agents/       # Agent personas (.md)
  ├── skills/       # Reusable skills (folders with SKILL.md)
  ├── workflows/    # Workflows (.md)
  └── rules/        # Global rules
  ```
- **Features**: Full support for all AGEN features

### Zed

Zed uses a folder-based approach for prompts and context.

- **Detection**: Checks for `.zed/` directory
- **Structure**:
  ```text
  .zed/
  ├── settings.json   # Configuration
  └── prompts/        # Context files
  ```
- **Install**: `agen init --ide zed`

### JetBrains

Supports IntelliJ IDEA, PyCharm, WebStorm, GoLand, and other JetBrains IDEs.

- **Detection**: Checks for `.idea/` directory
- **Format**: `.idea/ai-assistant.xml` + `.jbrules.md`
- **Install**: `agen init --ide jetbrains`
- **Plugins**: Works with JetBrains AI Assistant

### Neovim

Supports various Neovim AI plugins.

- **Detection**: Checks for `.nvim/` or `.nvim.lua`
- **Format**: `.nvim/ai-rules.md` + `.nvim.lua`
- **Install**: `agen init --ide neovim`
- **Plugins**: codecompanion.nvim, avante.nvim, ChatGPT.nvim

### Emacs

Supports Emacs AI plugins like gptel and ellama.

- **Detection**: Checks for `.dir-locals.el`
- **Format**: `.emacs-project/ai-context.md` + `.dir-locals.el`
- **Install**: `agen init --ide emacs`
- **Plugins**: gptel, ellama, copilot.el

---

## CLI Tools

### Aider

Aider is a command-line AI pair programming tool.

- **Detection**: Checks for `.aider.conf.yml` or `.aider/`
- **Format**: `.aider.conf.yml` + `.aider-context.md`
- **Install**: `agen init --ide aider`
- **Website**: [aider.chat](https://aider.chat)

---

## Cloud/Platform

### Claude Code

Anthropic's Claude Code uses a similar format to Antigravity.

- **Detection**: Checks for `CLAUDE.md`
- **Format**: Single `CLAUDE.md` file in project root
- **Install**: `agen init --ide claudecode`

### GitHub Copilot Workspace

GitHub Copilot can use custom instructions per repository.

- **Detection**: Checks for `.github/copilot-instructions.md`
- **Format**: `.github/copilot-instructions.md`
- **Install**: `agen init --ide copilotworkspace`
- **Docs**: [GitHub Copilot Custom Instructions](https://docs.github.com/en/copilot/customizing-copilot)

---

## Specifying an IDE

You can explicitly specify which IDE format to use:

```bash
# Auto-detect (default)
agen init

# Explicit IDE selection
agen init --ide cursor
agen init --ide continue
agen init --ide jetbrains
```

## Multiple IDEs

AGEN can generate configurations for multiple IDEs simultaneously:

```bash
agen init --ide cursor --ide continue --ide copilotworkspace
```
