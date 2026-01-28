# AGEN - AI Agent Template Manager

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/eshanized/agen)](https://github.com/eshanized/agen/releases)

A cross-platform CLI tool for managing AI agent templates. Supports multiple IDEs: **Antigravity**, **Cursor**, **Windsurf**, and **Zed**.

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

- **🎯 Multi-IDE Support**: Auto-detects and configures for Antigravity, Cursor, Windsurf, Zed
- **📦 Embedded Templates**: Works offline with built-in agent templates
- **🔄 Smart Updates**: Conflict resolution when updating templates you've modified
- **🔍 Fuzzy Search**: Find agents and skills with `agen search`
- **📊 Health Dashboard**: `agen health` shows project status and recommendations
- **🧪 Playground**: Test agents in temporary projects with `agen playground`
- **📁 Profiles**: Save and reuse agent configurations across projects
- **🔐 Verification**: Built-in security, lint, UX, and SEO checks

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

## 🔧 IDE Support

| IDE | Format | Auto-detection |
|-----|--------|----------------|
| **Antigravity/Claude** | `.agent/` folder | ✅ `.agent/` directory |
| **Cursor** | `.cursorrules` file | ✅ `.cursorrules` file |
| **Windsurf** | `.windsurfrules` file | ✅ `.windsurfrules` file |
| **Zed** | `.zed/prompts/` folder | ✅ `.zed/` directory |

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

## 🏗 Building from Source

```bash
git clone https://github.com/eshanized/agen.git
cd agen
go build -o agen ./cmd/agen
```

## 📜 License

MIT License - Copyright (c) 2026 [Eshan Roy](mailto:eshanized@proton.me)

## 🤝 Contributing

Contributions welcome! Please read the contributing guidelines first.

---

Made with ❤️ by [@eshanized](https://github.com/eshanized)
