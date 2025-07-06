# viberules

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux-lightgrey.svg)]()

[English](README.md) | [한국어](README.ko.md)

> AI assistant rules management tool using symlinks

⚠️ **Note**: Windows is not supported due to symlink limitations. Consider using WSL2 on Windows.

viberules is a CLI tool for managing AI coding assistant rules (Claude Code, Amazon Q Developer, Gemini Code Assist, etc.) using symlinks for real-time synchronization.

## ✨ Features

- 🎯 **Unified Management**: Manage all AI tool rules from a single file in .viberules/ folder
- 🔄 **Real-time Sync**: Changes automatically reflected via symlinks
- 🛠️ **Individual Target Control**: Enable/disable specific AI tools
- 🏠 **Flexible Modes**: Public mode (team sharing) or local mode (personal rules)
- 📝 **Smart .gitignore**: Automatic gitignore management with mode-aware policies
- 🌍 **Cross-platform**: Works on macOS and Linux

## 🚀 Quick Start

### Installation

```bash
# Install from GitHub
go install github.com/sky1core/viberules@latest
```

### Initialize Project

```bash
# Run in your project root
viberules init

# Force reinitialize (preserves existing rules.md)
viberules init --force
```

This creates:
- `.viberules/rules.md` - Single rules file for all AI tools
- Symlinks for each AI tool (CLAUDE.md, GEMINI.md, AGENTS.md, .amazonq/rules/AMAZONQ.md)
- Updated `.gitignore` with mode-aware policies

### Manage Targets

```bash
# List enabled targets
viberules list

# Remove unnecessary targets
viberules remove amazonq

# Add targets back
viberules add amazonq

# Set project mode
viberules mode public   # Share .viberules with team
viberules mode local    # Keep .viberules private
```

## 📋 Supported AI Tools

| AI Tool | Target Name | Output Files |
|---------|-------------|--------------|
| Claude Code | `claude` | `CLAUDE.md` |
| Amazon Q Developer | `amazonq` | `.amazonq/rules/AMAZONQ.md` |
| Gemini Code Assist | `gemini` | `GEMINI.md` |
| Generic AI Tools/Codex | `codex` | `AGENTS.md` |

## 🛠️ Commands

```bash
# Initialize project
viberules init

# Reinitialize existing project (preserves rules.md)
viberules init --force

# List enabled targets
viberules list

# Add/remove targets
viberules add claude
viberules remove amazonq

# Manage project mode
viberules mode          # Show current mode
viberules mode public   # Set to public mode (team sharing)
viberules mode local    # Set to local mode (private)

# Get help
viberules --help
```

## 📁 Project Structure

After running `viberules init`:

```
your-project/
├── .viberules/              # Configuration directory
│   ├── rules.md             # Single rules file for all AI tools
│   └── .config.yaml         # Configuration file (mode & targets, ignored by git)
├── .gitignore               # Updated automatically based on mode
├── CLAUDE.md                # Symlink to .viberules/rules.md
├── GEMINI.md                # Symlink to .viberules/rules.md
├── AGENTS.md                # Symlink to .viberules/rules.md
└── .amazonq/
    └── rules/
        └── AMAZONQ.md       # Symlink to ../../.viberules/rules.md
```

### Mode-based .gitignore Behavior

viberules supports two modes to control how rules are shared:

**Local Mode** (default, for personal rules):
- Entire `.viberules/` directory is ignored by git
- All rules remain private to your local machine
- Output files (CLAUDE.md, etc.) are ignored
- Use this when rules contain personal preferences or sensitive information

**Public Mode** (for team collaboration):
- `.viberules/rules.md` is tracked by git (shared with team)
- `.viberules/.config.yaml` is always ignored (personal config)
- Output files (CLAUDE.md, etc.) are ignored
- Use this when you want to share AI assistant rules with your team

## ⚙️ How It Works

1. **Edit Rules**: Modify `.viberules/rules.md` (single source of truth)
2. **Instant Sync**: Changes automatically appear in all AI tools via symlinks
3. **Mode-aware Git**: Public mode shares rules with team, local mode keeps everything private
4. **Smart Targeting**: Enable only the AI tools you use

## 🔧 Advanced Usage

### Project Modes

**Public Mode** (recommended for team projects):
```bash
viberules mode public
```
- `.viberules/rules.md` is committed and shared with team
- Personal settings (target configuration) remain local

**Local Mode** (for personal projects):
```bash
viberules mode local
```
- Entire `.viberules/` directory is ignored by git
- Rules stay completely private

### Writing Effective Rules

Edit `.viberules/rules.md`:

```markdown
# AI Assistant Rules

## Project Overview
This is a TypeScript React project using Next.js and Tailwind CSS.

## Coding Standards
- Use TypeScript with strict mode
- Follow ESLint configuration
- Write unit tests for all functions
- Use descriptive variable names

## Architecture Guidelines
- Follow clean architecture principles
- Separate business logic from UI components
- Use custom hooks for state management

## API Guidelines
- Use REST APIs with proper HTTP methods
- Implement proper error handling
- Use TypeScript interfaces for API responses
```

### Target Management

```bash
# Start with only Claude
viberules remove amazonq
viberules remove gemini

# Add others later
viberules add gemini
```

## 🧪 Development

### Prerequisites

- Go 1.21 or later
- macOS or Linux (Windows not supported)

### Build

```bash
# Clone and build
git clone https://github.com/sky1core/viberules.git
cd viberules
go build .
```

### Test

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run specific test
go test -v -run TestCompleteViberulesWorkflow .
```

---

<p align="center">
  Created by <a href="https://github.com/sky1core">sky1core</a>
</p>