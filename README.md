English | [中文](README_CN.md)

# Prompt Optimizer

A Claude Code skill that transforms your prompts into highly effective AI instructions.

## Why Use This?

**The Problem**: Vague prompts lead to vague responses. "Help me write a blog post" gets generic results.

**The Solution**: Prompt Optimizer analyzes your prompt, identifies weaknesses, and transforms it into a precise, well-structured instruction that gets better AI responses.

### Before & After

| Before | After |
|--------|-------|
| "Help me write a blog post about AI" | A structured prompt with clear objectives, target audience, tone guidelines, and specific output format |
| "Analyze this code" | A detailed prompt specifying analysis dimensions, expected findings format, and actionable recommendations |

## Features

- **Intelligent Optimization**: Automatically detects prompt complexity and applies the appropriate enhancement strategy
- **Clear Evaluation**: Get a score (0-100) with specific feedback on what's strong and what can improve
- **Iterative Refinement**: Keep improving your prompts with targeted instructions
- **Copy-Ready Output**: Optimized prompts are formatted for immediate use

## Installation

### One-Line Install

**macOS/Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/geq1fan/prompt-optimizer-skill/main/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/geq1fan/prompt-optimizer-skill/main/install.ps1 | iex
```

### Manual Install

```bash
# Clone to Claude Code skills directory
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

### Optimize a Prompt

```
/optimize-prompt Write a function to parse JSON
```

### Iterate on Results

```
/optimize-prompt iterate Add error handling requirements
```

## Example

**Input:**
```
/optimize-prompt Help me debug this React component
```

**Output:**

### Optimized Prompt
```
Analyze the following React component to identify and resolve issues.

**Context**: [Component code will be provided]

**Analysis Required**:
1. Identify syntax errors and typos
2. Check for React anti-patterns (e.g., missing keys, incorrect hook usage)
3. Evaluate state management and prop handling
4. Review performance implications

**Expected Output**:
- List of identified issues with line numbers
- Explanation of why each is problematic
- Corrected code with comments explaining changes
- Best practice recommendations for prevention
```

### Evaluation Results

**Overall Score: 88/100 (Good)**

| Dimension | Score | Assessment |
|-----------|-------|------------|
| Clarity | 90/100 | Clear debugging objective |
| Completeness | 85/100 | Good coverage, could specify component type |
| Specificity | 90/100 | Concrete analysis steps |
| Structure | 88/100 | Well-organized sections |

**Want to improve further?**
Try: `/optimize-prompt iterate Specify this is a form component with validation`

## Use Cases

| Scenario | How It Helps |
|----------|--------------|
| **Code Reviews** | Get structured, thorough review criteria |
| **Documentation** | Clear requirements for comprehensive docs |
| **Data Analysis** | Specific methodology and output format |
| **Creative Writing** | Defined tone, audience, and structure |
| **Problem Solving** | Step-by-step approach with constraints |

## How It Works

1. **Analysis**: Examines your prompt for clarity, completeness, and structure
2. **Strategy Selection**: Chooses optimization approach based on complexity
3. **Enhancement**: Applies targeted improvements while preserving intent
4. **Evaluation**: Provides actionable feedback and scores

## Templates

The skill uses carefully crafted templates:

| Template | Purpose |
|----------|---------|
| `optimize.md` | Main optimization with adaptive strategy |
| `iterate/general.md` | Targeted refinement based on instructions |
| `evaluation/user.md` | Comprehensive quality assessment |

## Contributing

Contributions welcome! Feel free to:
- Report issues
- Suggest improvements
- Submit pull requests

## License

MIT License

## Acknowledgments

Inspired by [linshenkx/prompt-optimizer](https://github.com/linshenkx/prompt-optimizer).
