English | [中文](README_CN.md)

# Prompt Optimizer

A Claude Code skill that optimizes prompts based on **Project Context**.

## Features

- **Smart optimization**: Automatically detects prompt complexity and applies the right strategy
- **Clear evaluation**: Get a 0–100 score plus concrete strengths/weaknesses feedback
- **Iterative improvements**: Continuously refine prompts with targeted instructions
- **Interactive confirmation**: Desktop WebView app for reviewing and confirming optimization results
- **Ready-to-use output**: Returns a formatted prompt you can copy directly

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

Contributions are welcome! You can:
- Report issues
- Suggest improvements
- Submit a Pull Request

## License

MIT License

## Acknowledgements

Inspired by [linshenkx/prompt-optimizer](https://github.com/linshenkx/prompt-optimizer).