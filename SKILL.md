---
name: optimize-prompt
description: 专业提示词优化工具，执行「优化 → 评估 → 持久化」工作流。支持用户/系统提示词优化、迭代改进和历史追踪。
---

# 提示词优化器

## 激活条件

当用户调用 `/optimize-prompt` 命令时激活此技能：

- `/optimize-prompt [内容]` — 优化用户提示词（basic 模式）
- `/optimize-prompt user [模式] [内容]` — 指定模式优化用户提示词
- `/optimize-prompt system [模式] [内容]` — 优化系统提示词
- `/optimize-prompt iterate [改进指令]` — 迭代改进上次优化结果

## 快速参考

| 类型 | 可用模式 | 模板路径 |
|-----|---------|---------|
| user | basic, professional, planning | `templates/{lang}/user-optimize/{mode}.md` |
| system | general, analytical | `templates/{lang}/system-optimize/{mode}.md` |
| iterate | general | `templates/{lang}/iterate/general.md` |

**语言检测**: 中文输入使用 `cn`，英文输入使用 `en`。

## 执行步骤

### 1. 解析输入

从用户命令中提取：
- `type`: user（默认）、system 或 iterate
- `mode`: 见快速参考（默认值：user→basic, system→general）
- `content`: 待优化的提示词（iterate 时为改进指令）

### 2. 处理迭代（type=iterate 时）

1. 读取工作目录下 `.prompt_history.jsonl` 的最后一行
2. 提取 `optimized` 字段作为 `lastOptimizedPrompt`
3. 用户 `content` 作为 `iterateInput`
4. 若无 content，询问用户改进指令

### 3. 优化

1. 读取模板: `templates/{lang}/{type}-optimize/{mode}.md`
   - iterate 类型: `templates/{lang}/iterate/general.md`
2. 替换占位符后，按模板指令生成优化后的提示词

### 4. 评估

1. 读取评估模板: `templates/{lang}/evaluation/{type}.md`
2. 按模板对优化结果进行评分（0-100）
3. 生成简要评估摘要

### 5. 持久化

将记录追加到 `.prompt_history.jsonl`。
详细格式见: [docs/history-format.md](docs/history-format.md)

### 6. 输出

向用户展示：
1. **优化后的提示词** — 使用代码块格式
2. **评估报告** — 分数和简要评价
3. **保存确认** — 确认已保存到历史文件

## 错误处理

| 情况 | 处理方式 |
|-----|---------|
| iterate 时历史文件不存在 | 创建空文件，询问用户原始提示词 |
| 无效的模式 | 列出可用模式，请用户选择 |
| 模板文件不存在 | 回退到英文模板，通知用户 |

## 详细文档

- [完整工作流程](docs/workflow.md)
- [历史文件格式](docs/history-format.md)
