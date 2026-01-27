# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Prompt Optimizer 是一个 Claude Code skill，用于基于**项目上下文**优化用户提示词。核心工作流程：分析 → 策略选择 → 优化 → 对抗评审 → 综合评估。

## Commands

### Go 后端（Wails 应用）

```bash
cd wails-app

# 运行测试
go test -v ./...

# 运行单个测试
go test -v -run TestName ./...

# 生成覆盖率报告
go test -v -race -coverprofile=coverage.out ./...

# 构建应用（需要 Wails CLI）
wails build
```

### 前端

```bash
cd wails-app/frontend

npm install                                       # 安装依赖
npm run dev                                       # 开发服务器
npm run dev:port 3000                             # 自定义端口
node dev-server.js --input ..\testdata\xxx.json  # 使用测试数据
npm test                                          # 运行测试（watch 模式）
npm run test:run                                  # 运行测试（一次性）
npm run test:coverage                             # 覆盖率报告
```

### 发布

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
# GitHub Actions 自动构建三平台版本并创建 Release
```

## Architecture

### 目录结构

```
prompt-optimizer-skill/
├── skills/               # 发布目录（用户安装后获得的内容）
│   ├── SKILL.md          # Skill 定义文件（入口点）
│   ├── templates/        # 优化模板
│   │   ├── cn/           # 中文模板
│   │   └── en/           # 英文模板
│   ├── references/       # 执行指南、模板规范等参考文档
│   └── bin/              # 预编译二进制（CI 构建时填充）
├── wails-app/            # WebView 桌面应用源码
│   ├── main.go           # CLI 参数解析 + Wails 初始化
│   ├── app.go            # Go 绑定函数（前端调用）
│   └── frontend/         # 前端代码（纯 HTML/CSS/JS）
├── .github/workflows/    # CI（测试+构建验证）和 Release（三平台打包）
├── install.sh            # macOS/Linux 安装脚本
└── install.ps1           # Windows 安装脚本
```

### 核心流程

1. **意图识别**：分析输入，判断 User 模式（新任务）还是 Iterate 模式（迭代）
2. **模式选择**：basic（通用）/ professional（代码/分析）/ planning（规划/流程）
3. **上下文采集**：探索项目文件，提取技术栈和代码规范
4. **优化生成**：读取模板 `templates/{lang}/{type}/{mode}.md`，替换占位符
5. **深度评审**：使用 `critical-review.md` 检测歧义、边界盲区、逻辑冲突
6. **综合评估**：使用 `evaluation/user.md` 生成评分（0-100）
7. **WebView 交互**：调用桌面应用，用户确认后可继续迭代或回滚

### 模板占位符

| 占位符 | 说明 |
|-------|------|
| `{{originalPrompt}}` | 用户原始提示词 |
| `{{lastOptimizedPrompt}}` | 上次优化结果（iterate 用） |
| `{{lastEvaluationReport}}` | 上次评估报告（iterate 用） |
| `{{iterateInput}}` | 用户改进指令（iterate 用） |
| `{{optimizedPrompt}}` | 优化后提示词（评估用） |
| `{{reviewReport}}` | 深度评审报告（评估用） |

### 语言检测

- 包含中文字符 → 使用 `cn/` 模板
- 不包含中文字符 → 使用 `en/` 模板
- 模板不存在时回退到英文

### WebView 应用

Wails v2 应用，通过 `--session-id` 参数自动推断文件路径：

```bash
# Windows
bin/prompt-optimizer-webview.exe --session-id <session_id>

# macOS
bin/prompt-optimizer-webview.app/Contents/MacOS/prompt-optimizer-webview --session-id <session_id>

# Linux
bin/prompt-optimizer-webview --session-id <session_id>
```

路径自动推断为 `.claude/prompt-optimizer/sessions/{session_id}/session.json` 和 `result.json`。
目录不存在时自动创建。默认 timeout 为 600 秒，可通过 `--timeout` 覆盖。

### CI/CD

- **CI** (`.github/workflows/ci.yml`)：push/PR 到 main 触发，运行后端测试、前端测试、三平台构建验证
- **Release** (`.github/workflows/release.yml`)：推送 `v*` 标签触发，构建并打包 `skills/` 目录（含二进制）到 GitHub Releases

### 其他

禁止使用 **take_screenshot** tool，使用 **take_snapshot** 进行替代
