# AGEN - AI Agent Template Manager

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/eshanized/agen)](https://github.com/eshanized/agen/releases)

A cross-platform CLI tool for managing AI agent templates. Supports **12 IDEs and AI coding tools** with unified configuration.

## 🚀 Quick Start

```bash
# Install (macOS/Linux)
brew install eshanized/tap/agen

# Install (Windows)
scoop bucket add eshanized https://github.com/eshanized/scoop-bucket
scoop install agen

# Or download binary from releases
```

```bash
# Initialize in your project
cd your-project
agen init

# List available agents
agen list

# Check project health
agen health
```

## ✨ Features

- **🎯 12 IDE Support**: Works with Cursor, VS Code, JetBrains, Neovim, Emacs, and more
- **📦 30+ Agents**: Backend, frontend, DevOps, security, ML, blockchain specialists
- **🛠 46+ Skills**: API patterns, testing, accessibility, Kubernetes, Rust, and more
- **🔄 16 Workflows**: `/plan`, `/deploy`, `/audit`, `/review`, `/migrate`, and more
- **🔍 Fuzzy Search**: Find agents and skills with `agen search`
- **📊 Health Dashboard**: `agen health` shows project status and recommendations
- **🧪 Playground**: Test agents in temporary projects with `agen playground`
- **📁 Profiles**: Save and reuse agent configurations across projects

## 🔧 Supported IDEs

| Category | IDEs/Tools |
|----------|------------|
| **VS Code Extensions** | Cursor, Windsurf, Cline, Continue |
| **Desktop IDEs** | Antigravity, Zed, JetBrains (IntelliJ, PyCharm), Neovim, Emacs |
| **CLI Tools** | Aider |
| **Cloud/Platform** | Claude Code, GitHub Copilot Workspace |

<details>
<summary><strong>View IDE Configuration Formats</strong></summary>

| IDE | Config Format | Detection |
|-----|---------------|-----------|
| Antigravity | `.agent/` folder | `.agent/` directory |
| Cursor | `.cursorrules` | `.cursorrules` file |
| Windsurf | `.windsurfrules` | `.windsurfrules` file |
| Cline | `.clinerules` | `.clinerules` file |
| Continue | `.continue/` | `.continue/` directory |
| Zed | `.zed/prompts/` | `.zed/` directory |
| JetBrains | `.idea/ai-assistant.xml` | `.idea/` directory |
| Neovim | `.nvim/ai-rules.md` | `.nvim/` or `.nvim.lua` |
| Emacs | `.dir-locals.el` | `.dir-locals.el` file |
| Aider | `.aider.conf.yml` | `.aider.conf.yml` |
| Claude Code | `CLAUDE.md` | `CLAUDE.md` file |
| Copilot Workspace | `.github/copilot-instructions.md` | `.github/` |

</details>

## 📖 Commands

| Command | Description |
|---------|-------------|
| `agen init` | Initialize templates in current project |
| `agen list` | List available agents, skills, workflows |
| `agen status` | Check installation status |
| `agen health` | Show project health dashboard |
| `agen verify` | Run verification scripts |
| `agen search` | Fuzzy search agents/skills |
| `agen update` | Update templates to latest version |
| `agen upgrade` | Update agen binary itself |
| `agen profile` | Manage saved configurations |
| `agen playground` | Create temporary test project |

## 📦 Installation

### Homebrew (macOS/Linux)
```bash
brew install eshanized/tap/agen
```

### Scoop (Windows)
```bash
scoop bucket add eshanized https://github.com/eshanized/scoop-bucket
scoop install agen
```

### Arch Linux (AUR)
```bash
yay -S agen-bin
```

### Debian/Ubuntu
```bash
# Download .deb from releases
sudo dpkg -i agen_*.deb
```

### Binary Download
Download from [GitHub Releases](https://github.com/eshanized/agen/releases).

## 🤖 Available Agents

<details>
<summary><strong>View all 30+ agents</strong></summary>

| Agent | Domain |
|-------|--------|
| `orchestrator` | Multi-agent coordination |
| `backend-specialist` | Server-side development |
| `frontend-developer` | UI/UX implementation |
| `mobile-developer` | React Native, Flutter |
| `devops-engineer` | CI/CD, infrastructure |
| `security-auditor` | Security analysis |
| `debugger` | Root cause analysis |
| `test-engineer` | Testing strategies |
| `ai-ml-engineer` | ML pipelines, LLMs |
| `data-engineer` | ETL, data warehousing |
| `cloud-architect` | AWS/GCP/Azure |
| `blockchain-developer` | Smart contracts |
| `tech-lead` | Architecture decisions |
| `accessibility-specialist` | WCAG compliance |
| *...and more* | |

</details>

## 🏗 Building from Source

```bash
git clone https://github.com/eshanized/agen.git
cd agen
go build -o agen ./cmd/agen
```

## 📚 Documentation

Visit the [full documentation](https://eshanized.github.io/agen/) for detailed guides.

## 📜 License

MIT License - Copyright (c) 2026 [Eshan Roy](mailto:eshanized@proton.me)

## 🤝 Contributing

Contributions welcome! Please read the [contributing guidelines](docs/contributing.md) first.

---

Made with ❤️ by [@eshanized](https://github.com/eshanized)
