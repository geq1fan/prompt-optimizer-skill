---
name: prompt-optimizer
description: |
  专业用户提示词优化工具，执行「优化 → 对抗测试 → 评估」工作流。
  触发条件：(1) 用户输入 `/optimize-prompt [内容]` 或 `-enhancer` 标记；(2) 用户请求优化、改进、增强提示词。
  每次调用创建新 Session，迭代只能通过 WebView 交互触发。
---

# 提示词优化器

## 命令格式

```
/optimize-prompt [内容]
```

每次调用 `/optimize-prompt` 都会创建一个新的 Session，开始一轮完整的优化流程。

**迭代方式**：通过 WebView 交互（submit/rollback）触发，无需再次输入命令。

## 模板参考

| 类型 | 模式 | 路径 |
|-----|------|-----|
| user | basic / professional / planning | `templates/{lang}/user-optimize/{mode}.md` |
| iterate | general | `templates/{lang}/iterate/general.md` |
| review | critical-review | `templates/{lang}/evaluation/critical-review.md` |
| evaluation | user | `templates/{lang}/evaluation/user.md` |

**语言检测**：中文输入 → `cn/`；其他 → `en/`

## 执行流程

### 1. Session 创建

每次 `/optimize-prompt` 调用都会创建新 Session：
- 生成 `session_id`：`session_{timestamp}`（毫秒时间戳）
- Session 目录：`.claude/prompt-optimizer/sessions/{session_id}/`（由 WebView 自动创建）
- 创建 `session.json`（v4 格式，唯一数据源）

### 2. 模式选择

根据任务内容选择优化模板：

| 特征 | 模式 | 模板 |
|-----|------|-----|
| 涉及步骤、计划、流程、roadmap | planning | `user-optimize/planning.md` |
| 涉及代码、分析、学术、技术文档 | professional | `user-optimize/professional.md` |
| 其他通用场景 | basic | `user-optimize/basic.md` |

### 3. 上下文采集

判断任务是否依赖项目信息：
- **是**：使用 `ls`/`read`/`glob`/`grep` 探索项目结构，提取技术栈和规范
- **否**：跳过

### 4. 优化生成

读取模板并替换占位符，生成优化后的提示词。

### 5. 深度评审

读取 `critical-review.md`，检测歧义、边界盲区、逻辑冲突。

### 6. 综合评估

读取 `evaluation/user.md`，将评审报告作为 `{{reviewReport}}` 传入，生成评分（0-100）。

### 7. 输出展示

简要提示优化完成，不输出完整内容：

```
优化完成，评分: XX/100 ([等级])
正在打开交互界面...
```

**注意**: 所有详细信息通过 WebView 展示，不在终端输出完整的提示词和报告。

### 8. WebView 交互循环 ⚠️ 必须执行

**重要**: 此步骤是必须的用户确认环节，不可跳过。

调用 WebView 应用进行交互确认，进入优化循环：

```
┌─────────────────────────────────────────┐
│  优化循环 (在同一个 Session 内)          │
│                                         │
│  显示当前结果 → 用户操作：               │
│  ├─ submit → 更新 session.json → 迭代优化│
│  ├─ rollback → 恢复历史 → 迭代优化       │
│  └─ cancel/timeout → 结束 Session        │
│                                         │
└─────────────────────────────────────────┘
```

⚠️ **警告**:
- 每次优化完成后必须调用 WebView
- 不要直接结束流程，等待用户在 WebView 中做出选择
- 仅当 WebView 二进制不存在时才可跳过

详见 [WebView 指南](references/webview-guide.md)。

## 错误处理

| 情况 | 处理 |
|-----|-----|
| 无有效内容 | 提示用户提供提示词内容 |
| 模板不存在 | 回退英文模板 |
| WebView 应用不存在 | 直接输出结果，跳过交互 |
| Session 目录创建失败 | 输出错误信息，终止流程 |

## 参考文档

| 文档 | 说明 |
|-----|------|
| [执行指南](references/execution-guide.md) | 完整执行流程，Phase 1-8 详解 |
| [WebView 指南](references/webview-guide.md) | 交互确认核心组件，输入输出格式 |
| [模板规范](references/templates-spec.md) | 占位符、目录结构、评估维度 |
