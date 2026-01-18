---
name: optimize-prompt
description: 专业用户提示词优化工具，执行「优化 → 对抗测试 → 评估」工作流。
---

# 提示词优化器

## 激活条件

用户只需输入 `/optimize-prompt [内容]`，系统会自动判断是**新任务**还是**迭代指令**。

- `/optimize-prompt [内容]` — 智能识别意图（推荐）
  - 如果内容是"更短一点"、"添加示例"等指令 → 自动进入 **Iterate 模式**
  - 如果内容是"写个周报"、"分析代码"等新话题 → 自动进入 **User 模式** (Basic/Pro/Plan)
- `/optimize-prompt iterate [指令]` — 强制进入迭代模式
- `/optimize-prompt user [模式] [内容]` — 强制进入新任务模式

## 快速参考

| 类型 | 可用模式 | 模板路径 |
|-----|---------|---------|
| user | basic, professional, planning | `templates/{lang}/user-optimize/{mode}.md` |
| iterate | general | `templates/{lang}/iterate/general.md` |
| review | default | `templates/{lang}/evaluation/critical-review.md` |

**语言检测**: 中文输入使用 `cn`，英文输入使用 `en`。

## 执行步骤

### 1. 意图识别与路由 (Dispatcher)

系统首先分析用户的输入内容和当前的对话上下文：

**优先级判定逻辑**：
1. **显式指定**: 若用户输入了 `iterate` 或 `user` 关键字，直接进入对应模式。
2. **冷启动**: 若当前对话历史中**不存在** `lastOptimizedPrompt`，默认为 **User 模式**。
3. **语义分析**: 若存在历史记录，分析输入内容与上文的关系：
   - **Iterate 模式 (迭代)**: 输入是对上文的修改、补充、反馈。
     - *特征词*: "更..."、"有点..."、"添加..."、"删除..."、"修改..."、"不..."、"太..."。
   - **User 模式 (新任务)**: 输入是一个全新的、独立的话题。
     - *特征*: 完整的任务描述，与上文无明显指代关系。

### 2. 模式执行

#### A. User 模式 (新任务)
若判定为新任务，进一步分析内容特征以选择模板：

1. **Planning (规划模式)**: 涉及步骤、计划、流程、roadmap。
2. **Professional (专业模式)**: 涉及代码、分析、学术、专业输出。
3. **Basic (基础模式)**: 其他通用场景。

#### B. Iterate 模式 (迭代)
若判定为迭代指令：

1. **自动检索上下文**: 
   - 获取 `lastOptimizedPrompt` (上次结果) 和 `lastEvaluationReport` (上次报告)。
2. **提取改进指令**: 
   - 将用户输入作为 `iterateInput`。

### 3. 上下文采集 (Context Gathering) - Agent 核心能力

**这是区别于普通优化器的关键步骤。**

在生成 Prompt 之前，Agent **必须**判断：完成该任务是否依赖当前项目的具体信息？

- **如果是**（例如涉及代码、架构、特定文件）：
  - **主动探索**: 立即使用工具 (`ls`, `read`, `glob`, `grep`) 扫描相关文件结构、读取关键代码或配置文件。
  - **信息提取**: 识别项目的技术栈、代码规范、依赖版本等关键元数据。
- **如果否**（通用型任务）：
  - 跳过此步骤。

**目的**: 确保最终生成的 Prompt 是**基于事实**的，而非基于幻觉或通用模板的。

### 4. 优化生成

1. 读取模板: `templates/{lang}/{type}-optimize/{mode}.md`
   - iterate 类型: `templates/{lang}/iterate/general.md`
2. 替换占位符后，按模板指令生成优化后的提示词

### 4. 深度评审 (Critical Review)

1. 读取评审模板: `templates/{lang}/evaluation/critical-review.md`
2. 将优化后的提示词传入，生成**深度评审报告** (Critical Review Report)
3. 重点检测：歧义表达、边界盲区、逻辑冲突

### 5. 综合评估 (Final Evaluation)

1. 读取评估模板: `templates/{lang}/evaluation/user.md`
2. **关键**: 将「深度评审报告」作为 `{{reviewReport}}` 传入评估模板
3. 结合评审发现的逻辑盲区，对优化结果进行客观评分
4. 生成详细评估报告

### 6. 输出

向用户展示：
1. **优化后的提示词** — 使用代码块格式
2. **深度评审摘要** — 简要展示发现的潜在理解偏差
3. **评估报告** — 包含分数、维度评价和改进建议
4. **当前模式** — 在报告中注明使用的优化模式

### 7. 交互式确认 (可选)

当需要用户确认优化结果时，可调用 WebView 桌面应用进行交互：

```bash
# 调用 WebView 应用
webui/bin/prompt-optimizer-webview --input <input.json> --output <result.json> --timeout 600
```

**输入文件格式 (input.json)**:
```json
{
  "version": 3,
  "originalPrompt": "原始提示词",
  "current": {
    "iterationId": "iter-003",
    "optimizedPrompt": "优化后的提示词 (Markdown)",
    "reviewReport": "评审报告 (Markdown)",
    "evaluationReport": "评估报告 (Markdown)",
    "score": 85,
    "suggestedDirections": [
      {"id": "examples", "label": "添加示例", "description": "补充使用案例"}
    ]
  },
  "history": []
}
```

**输出结果格式 (result.json)**:
```json
// 用户确认
{"action": "submit", "selectedDirections": ["examples"], "userInput": "补充说明"}

// 用户取消
{"action": "cancel", "selectedDirections": [], "userInput": ""}

// 超时
{"action": "timeout", "selectedDirections": [], "userInput": ""}

// 回滚到历史版本
{"action": "rollback", "rollbackToIteration": "iter-001", "selectedDirections": [], "userInput": ""}
```

## 错误处理

| 情况 | 处理方式 |
|-----|---------|
| iterate 找不到历史提示词 | 提示用户：无法自动获取上下文，请直接提供提示词内容 |
| 无效的模式 | 列出可用模式，请用户选择 |
| 模板文件不存在 | 回退到英文模板，通知用户 |

## 详细文档

- [完整工作流程](docs/workflow.md)
