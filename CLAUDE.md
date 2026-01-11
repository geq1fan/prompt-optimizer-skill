# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **Claude Code Skill** that optimizes user prompts. It transforms vague prompts into precise, well-structured AI instructions with evaluation feedback.

**Skill activation**: `/optimize-prompt`

## Architecture

```
prompt-optimizer-skill/
├── SKILL.md                      # Skill definition (workflow, activation, output format)
├── templates/
│   ├── optimize.md               # Main optimization template (auto-detects complexity)
│   ├── iterate/general.md        # Iterative refinement template
│   └── evaluation/user.md        # Evaluation criteria and scoring
├── install.sh                    # macOS/Linux installer
└── install.ps1                   # Windows installer
```

### Template Placeholders

| Placeholder | Used In | Description |
|-------------|---------|-------------|
| `{{originalPrompt}}` | optimize.md | User's input prompt |
| `{{lastOptimizedPrompt}}` | iterate/general.md | Previous optimized prompt |
| `{{iterateInput}}` | iterate/general.md | User's improvement instructions |
| `{{optimizedPrompt}}` | evaluation/user.md | Prompt to evaluate |
| `{{hasOriginalPrompt}}` | evaluation/user.md | Conditional for comparison mode |

### Workflow

1. **Parse**: Extract type (optimize/iterate) and content from user command
2. **Optimize**: Apply template to transform the prompt
3. **Evaluate**: Score on 4 dimensions (Clarity, Completeness, Specificity, Structure)
4. **Output**: Present optimized prompt + evaluation report

## Key Files

- **SKILL.md**: Defines the skill behavior, workflow steps, and output format. This is what Claude Code reads to execute the skill.
- **templates/optimize.md**: Contains three strategies (A: Clarity, B: Precision, C: Structured Planning) - Claude auto-selects based on prompt complexity.

## Installation Scripts

Both scripts support: `install`, `update`, `uninstall`

```bash
# Test install script locally
./install.sh install    # or: ./install.sh update
```

```powershell
# Windows
.\install.ps1 -Action install
```

## Conventions

- All templates output only the result (no meta-commentary or explanations)
- Evaluation scores use 0-100 scale with levels: Excellent (90+), Good (80-89), Adequate (70-79), Needs Work (60-69), Poor (<60)
- Templates are English-only; README has bilingual support (README.md / README_CN.md)
