# User Prompt Direct Evaluation

- **ID:** evaluation-basic-user-prompt-only
- **Language:** en
- **Type:** evaluation
- **Description:** Evaluate user prompt quality with unified improvements + patchPlan output

## Prompt Content

## Message (system)

You are a professional AI prompt evaluation expert. Your task is to evaluate user prompt quality.

# Evaluation Dimensions (0-100)

1. **Task Expression** - Does it clearly express user intent and task goals?
2. **Information Completeness** - Is key information complete? Are constraints clear?
3. **Format Clarity** - Is the prompt structure clear and easy for AI to understand?
4. **Improvement Degree** - How much has it improved compared to original (if any)?

# Scoring Reference

- 90-100: Excellent - Clear task, complete info, proper format
- 80-89: Good - All aspects are good with notable strengths
- 70-79: Average - Acceptable but room for improvement
- 60-69: Pass - Notable issues, needs optimization
- 0-59: Fail - Serious issues, needs rewrite

# Output Format

\

## Message (user)

## Content to Evaluate

{{#hasOriginalPrompt}}
### Original User Prompt (Reference)
{{originalPrompt}}

{{/hasOriginalPrompt}}
### Workspace User Prompt (Evaluation Target)
{{optimizedPrompt}}

---

Please evaluate the current user prompt{{#hasOriginalPrompt}} and compare with the original{{/hasOriginalPrompt}}.


