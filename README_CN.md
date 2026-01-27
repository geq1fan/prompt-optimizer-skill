[English](README.md) | 中文

# Prompt Optimizer 🚀

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Claude Code](https://img.shields.io/badge/Built%20for-Claude%20Code-d97757)](https://claude.ai)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Windows%20%7C%20Linux-blue)](https://github.com/geq1fan/prompt-optimizer-skill/releases)

**专业的 Claude Code 技能，利用对抗性评估和项目上下文，将简单的指令转化为生产就绪的提示词。**

![Prompt Optimization Workflow](assets/demo.gif)
*一键将模糊的需求转化为结构化、无漏洞的专业提示词。*

## 功能特性

- **感知项目上下文 (Context-Aware)**：自动分析您的*整个项目结构*（文件、设计文档），生成深度相关的提示词，而非通用模板。
- **对抗性测试 (Adversarial Testing)**：自动模拟 "红队测试"，在运行前发现提示词中的逻辑漏洞和边界情况。
- **量化评估 (Quantitative Evaluation)**：基于清晰度、特异性和鲁棒性提供 0-100 分的详细评分。
- **交互式评审 (WebView)**：原生桌面 UI (Wails)，让您无需离开工作流即可对比、编辑和确认更改。
- **即用输出**：直接返回格式化好的提示词，可直接复制到 Claude 使用。

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

## 工作原理

![Prompt Optimizer Architecture](assets/architecture.png)

1. **分析**：检查提示词的清晰度、完整性和结构
2. **策略选择**：根据复杂度选择优化方案
3. **增强**：在保持原意的同时进行针对性改进
4. **评估**：提供可操作的反馈和评分
5. **交互式确认**：使用 WebView 应用查看和确认优化结果

## 贡献

欢迎贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解如何提交 PR、报告问题和建议改进。

## 许可证

本项目基于 MIT 许可证开源 - 详见 [LICENSE](LICENSE) 文件。

## 致谢

灵感来源于 [linshenkx/prompt-optimizer](https://github.com/linshenkx/prompt-optimizer)。
