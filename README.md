# Prompt Optimizer Skill for Claude Code

基于 [prompt-optimizer](https://github.com/linshenkx/prompt-optimizer) 项目的 Claude Code Skill，提供 AI 提示词优化能力。

## 功能特性

- **用户提示词优化**: 消除模糊表达，补充关键信息，提升表达清晰度
- **系统提示词优化**: 按标准结构重组角色定义、技能和规则
- **迭代优化**: 基于反馈改进现有提示词

## 安装

将此目录复制到 Claude Code 的 skills 目录：

```bash
# macOS/Linux
cp -r prompt-optimizer-skill ~/.claude/skills/

# Windows
xcopy /E prompt-optimizer-skill %USERPROFILE%\.claude\skills\
```

## 使用方法

### 快速开始

```bash
# 优化用户提示词
/optimize-prompt 帮我写一篇关于人工智能的文章

# 优化系统提示词
/optimize-prompt system 你是一个客服助手

# 迭代优化现有提示词
/optimize-prompt iterate
```

### 完整语法

```
/optimize-prompt [type] [mode] [提示词内容]
```

| 参数 | 可选值 | 默认值 | 说明 |
|------|--------|--------|------|
| type | user, system, iterate | user | 优化类型 |
| mode | basic, professional, planning, general, analytical | basic/general | 优化模式 |

### 优化模式详解

**用户提示词模式：**

| 模式 | 说明 | 适用场景 |
|------|------|----------|
| basic | 基础优化 | 日常对话、简单问答 |
| professional | 专业优化 | 技术文档、专业需求 |
| planning | 规划优化 | 复杂任务、项目规划 |

**系统提示词模式：**

| 模式 | 说明 | 适用场景 |
|------|------|----------|
| general | 通用优化 | 常规角色定义 |
| analytical | 分析式优化 | 复杂业务场景 |

## 示例

### 用户提示词优化

```bash
# 基础优化（默认）
/optimize-prompt 帮我分析这段代码

# 专业优化
/optimize-prompt user professional 帮我写一份技术方案，关于微服务架构设计

# 规划优化
/optimize-prompt user planning 开发一个电商网站的后台管理系统
```

### 系统提示词优化

```bash
# 通用优化（默认）
/optimize-prompt system 你是一个编程助手

# 分析式优化
/optimize-prompt system analytical 你是一个代码审查专家，负责review团队的代码
```

### 迭代优化

```bash
/optimize-prompt iterate
# 然后按提示输入：
# 1. 现有提示词
# 2. 优化需求
```

## 自定义模板

模板文件位于 `templates/` 目录，可以直接修改：

```
templates/
├── cn/                     # 中文模板
│   ├── user-optimize/
│   │   ├── basic.md
│   │   ├── professional.md
│   │   └── planning.md
│   ├── system-optimize/
│   │   ├── general.md
│   │   └── analytical.md
│   └── iterate/
│       └── general.md
└── en/                     # 英文模板
    └── ...
```

### 模板变量

| 变量 | 说明 |
|------|------|
| `{{originalPrompt}}` | 原始提示词 |
| `{{lastOptimizedPrompt}}` | 待迭代的提示词 |
| `{{iterateInput}}` | 优化需求描述 |

## 完整模板存档

`archived_templates/` 目录包含更多模板：

- **图像优化**: text2image, image2image
- **上下文优化**: 带上下文的高级优化
- **评估模板**: 提示词质量评估

## 致谢

本项目基于 [linshenkx/prompt-optimizer](https://github.com/linshenkx/prompt-optimizer) 的模板系统。

## 许可证

MIT License
