[English](README.md) | 中文

# Prompt Optimizer

一个将模糊提示词转化为高效 AI 指令的 Claude Code 技能。

## 为什么使用？

**问题**：模糊的提示词导致模糊的回答。"帮我写一篇博客" 只能得到泛泛的结果。

**解决方案**：Prompt Optimizer 分析你的提示词，识别薄弱环节，将其转化为精确、结构化的指令，从而获得更好的 AI 回答。

### 优化前后对比

| 优化前 | 优化后 |
|--------|-------|
| "帮我写一篇关于AI的博客" | 包含明确目标、目标受众、语气指南和具体输出格式的结构化提示词 |
| "分析这段代码" | 指定分析维度、预期结果格式和可执行建议的详细提示词 |

## 功能特性

- **智能优化**：自动检测提示词复杂度，应用合适的优化策略
- **清晰评估**：获得 0-100 分的评分，以及具体的优缺点反馈
- **迭代改进**：通过针对性指令持续优化提示词
- **即用输出**：优化后的提示词格式化呈现，可直接复制使用

## 安装

### 一键安装

**macOS/Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/geq1fan/prompt-optimizer-skill/main/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/geq1fan/prompt-optimizer-skill/main/install.ps1 | iex
```

### 手动安装

```bash
# 克隆到 Claude Code skills 目录
git clone https://github.com/geq1fan/prompt-optimizer-skill ~/.claude/skills/prompt-optimizer-skill
```

### 更新

```bash
# macOS/Linux
~/.claude/skills/prompt-optimizer-skill/install.sh update

# Windows
& "$env:USERPROFILE\.claude\skills\prompt-optimizer-skill\install.ps1" -Action update
```

## 使用方法

### 优化提示词

```
/optimize-prompt 写一个解析JSON的函数
```

### 迭代改进

```
/optimize-prompt iterate 添加错误处理要求
```

## 示例

**输入:**
```
/optimize-prompt 帮我调试这个React组件
```

**输出:**

### 优化后的提示词
```
分析以下 React 组件，识别并解决问题。

**上下文**: [将提供组件代码]

**分析要求**:
1. 识别语法错误和拼写问题
2. 检查 React 反模式（如：缺少 key、错误的 hook 用法）
3. 评估状态管理和 props 处理
4. 审查性能影响

**预期输出**:
- 带行号的问题列表
- 解释每个问题为何有问题
- 修正后的代码，附带解释变更的注释
- 预防问题的最佳实践建议
```

### 评估结果

**综合评分: 88/100 (良好)**

| 维度 | 分数 | 评价 |
|-----------|-------|------------|
| 清晰度 | 90/100 | 调试目标明确 |
| 完整性 | 85/100 | 覆盖面好，可指定组件类型 |
| 具体性 | 90/100 | 分析步骤具体 |
| 结构性 | 88/100 | 组织良好 |

**想要进一步改进？**
尝试: `/optimize-prompt iterate 指定这是一个带表单验证的组件`

## 应用场景

| 场景 | 如何帮助 |
|----------|--------------|
| **代码审查** | 获得结构化、全面的审查标准 |
| **文档编写** | 明确全面文档的要求 |
| **数据分析** | 具体的方法论和输出格式 |
| **创意写作** | 定义语气、受众和结构 |
| **问题解决** | 带约束条件的分步方法 |

## 工作原理

1. **分析**：检查提示词的清晰度、完整性和结构
2. **策略选择**：根据复杂度选择优化方案
3. **增强**：在保持原意的同时进行针对性改进
4. **评估**：提供可操作的反馈和评分

## 模板

该技能使用精心设计的模板：

| 模板 | 用途 |
|----------|---------|
| `optimize.md` | 自适应策略的主优化模板 |
| `iterate/general.md` | 基于指令的定向改进 |
| `evaluation/user.md` | 全面的质量评估 |

## 贡献

欢迎贡献！你可以：
- 报告问题
- 提出改进建议
- 提交 Pull Request

## 许可证

MIT License

## 致谢

灵感来源于 [linshenkx/prompt-optimizer](https://github.com/linshenkx/prompt-optimizer)。
