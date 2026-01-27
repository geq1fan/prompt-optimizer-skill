English | [中文](README_CN.md)

# Prompt Optimizer 🚀

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Claude Code](https://img.shields.io/badge/Built%20for-Claude%20Code-d97757)](https://claude.ai)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Windows%20%7C%20Linux-blue)](https://github.com/geq1fan/prompt-optimizer-skill/releases)

**A professional Claude Code skill that turns simple instructions into production-ready prompts using adversarial evaluation and project context.**

![Prompt Optimization Workflow](assets/demo.gif)
*Watch how a vague request becomes a structured, edge-case-proof prompt in seconds.*

## Features

- **Context-Aware Optimization**: Analyzes your *entire project structure* (files, design docs) to generate relevant prompts, not generic ones.
- **Adversarial Testing**: Automatically simulates "Red Teaming" to find logical loopholes in your prompt before you run it.
- **Quantitative Evaluation**: Detailed scoring system (0-100) based on clarity, specificity, and robustness.
- **Interactive Review (WebView)**: A native desktop UI (Wails) to diff, edit, and confirm changes without leaving your workflow.
- **Ready-to-use output**: Returns a formatted prompt you can copy directly into Claude.

## Installation

### One-line install

**macOS/Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/geq1fan/prompt-optimizer-skill/main/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/geq1fan/prompt-optimizer-skill/main/install.ps1 | iex
```

### Manual install

```bash
# Clone into the Claude Code skills directory
git clone https://github.com/geq1fan/prompt-optimizer-skill ~/.claude/skills/prompt-optimizer-skill
```

### Update

```bash
# macOS/Linux
~/.claude/skills/prompt-optimizer-skill/install.sh update

# Windows
& "$env:USERPROFILE\.claude\skills\prompt-optimizer-skill\install.ps1" -Action update
```

## Usage

### Optimize a prompt

```
/optimize-prompt Write a function to parse JSON
```

### Iterate

```
/optimize-prompt iterate Add error handling requirements
```

## How it works

![Prompt Optimizer Architecture](assets/architecture.png)

1. **Analyze**: Checks prompt clarity, completeness, and structure
2. **Choose a strategy**: Selects an optimization approach based on complexity
3. **Enhance**: Improves the prompt while preserving intent
4. **Evaluate**: Provides actionable feedback and a score
5. **Interactive Review**: Use the WebView app to review and confirm the optimized prompt

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to submit pull requests, report issues, and suggest improvements.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

Inspired by [linshenkx/prompt-optimizer](https://github.com/linshenkx/prompt-optimizer).
