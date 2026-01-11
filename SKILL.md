---
name: optimize-prompt
description: Transform vague prompts into highly effective AI instructions with intelligent optimization and clear evaluation feedback.
context: fork
---

# Prompt Optimizer

A Claude Code skill that transforms your prompts into highly effective AI instructions.

## Activation

Invoke with `/optimize-prompt`:

- `/optimize-prompt <prompt>` — Optimize a prompt
- `/optimize-prompt iterate` — Refine the last optimized prompt

## Workflow

### 1. Parse Input

From the user command, extract:
- `type`: "optimize" (default) or "iterate"
- `content`: The prompt to optimize (or improvement instructions for iterate)

### 2. Optimize

**For new prompts:**
1. Read template: `templates/optimize.md`
2. Replace `{{originalPrompt}}` with user's prompt
3. Generate optimized prompt following template instructions

**For iterations:**
1. Read template: `templates/iterate/general.md`
2. Ask user for:
   - The prompt to improve (if not provided)
   - Improvement instructions
3. Replace placeholders and generate refined prompt

### 3. Evaluate

1. Read evaluation template: `templates/evaluation/user.md`
2. Evaluate the optimized prompt
3. Generate evaluation report with:
   - Overall score (0-100)
   - Dimension scores (Clarity, Completeness, Specificity, Structure)
   - Key strengths and areas for improvement
   - Specific iteration recommendations

### 4. Output

Present to the user:

1. **Optimized Prompt** — In a code block for easy copying
2. **Evaluation Report** — Scores and actionable feedback
3. **Next Steps** — Suggestions for further iteration if needed

## Example Output Format

```
## Optimized Prompt

\`\`\`
[The optimized prompt text]
\`\`\`

## Evaluation Results

### Overall Score: 85/100 (Good)

### Dimension Scores
| Dimension | Score | Assessment |
|-----------|-------|------------|
| Clarity | 90/100 | Intent is crystal clear |
| Completeness | 80/100 | Good context, could add constraints |
| Specificity | 85/100 | Requirements are concrete |
| Structure | 85/100 | Well-organized |

### Key Strengths
- Clear objective statement
- Good use of specific parameters

### Areas for Improvement
- Add output format specification
- Include error handling instructions

### Want to Improve Further?
Try: `/optimize-prompt iterate` with instructions like:
- "Add JSON output format"
- "Make it more concise"
```

## Error Handling

| Situation | Action |
|-----------|--------|
| No prompt provided | Ask user for the prompt |
| Iterate without previous prompt | Ask user to provide the prompt to improve |
| Template not found | Report error and suggest checking installation |

## Templates

| File | Purpose |
|------|---------|
| `templates/optimize.md` | Main optimization template |
| `templates/iterate/general.md` | Iteration refinement template |
| `templates/evaluation/user.md` | Evaluation criteria and format |
