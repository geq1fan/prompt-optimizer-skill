# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Prompt Optimizer 是一个 Claude Code skill，用于基于**项目上下文**优化用户提示词。核心工作流程：分析 → 策略选择 → 优化 → 对抗评审 → 综合评估。

## Commands

### Wails 应用（Go 后端 + 原生 WebView）

```bash
# 进入 wails-app 目录
cd wails-app

# 运行测试
go test -v ./...

# 运行单个测试
go test -v -run TestName ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...

# 构建应用
wails build
```

### 前端测试

```bash
cd wails-app/frontend

# 安装依赖
npm install

# 自定义端口
npm run dev:port 3000

# 使用自定义测试数据
node dev-server.js --input ..\testdata\input_v5_long_history.json

# 运行测试
npm test

# 运行测试（一次性）
npm run test:run

# 运行覆盖率测试
npm run test:coverage
```

### 安装脚本

```bash
# macOS/Linux
./install.sh

# Windows PowerShell
.\install.ps1

# 更新
./install.sh update
```

## Architecture

### 目录结构

```
prompt-optimizer-skill/
├── skills/               # 发布目录（用户安装后获得的内容）
│   ├── SKILL.md          # Skill 定义文件（入口点）
│   ├── templates/        # 优化模板
│   │   ├── cn/           # 中文模板
│   │   │   ├── user-optimize/   # 新任务模板 (basic/professional/planning)
│   │   │   ├── iterate/         # 迭代模板
│   │   │   └── evaluation/      # 评估模板
│   │   └── en/           # 英文模板（结构相同）
│   └── bin/              # 预编译二进制（构建时填充）
├── wails-app/            # WebView 桌面应用源码
│   ├── app.go            # Go 后端逻辑
│   ├── main.go           # 应用入口
│   └── frontend/         # 前端代码
├── .github/              # CI/CD 配置
├── install.sh            # macOS/Linux 安装脚本
├── install.ps1           # Windows 安装脚本
└── docs/                 # 详细文档
```

### 核心流程

1. **意图识别 (Dispatcher)**：分析用户输入，判断是新任务还是迭代
   - 冷启动（无历史）→ User 模式
   - 有历史 + 修改类指令 → Iterate 模式

2. **模式选择**：
   - User 模式：basic（通用）/ professional（代码/分析）/ planning（规划/流程）
   - Iterate 模式：基于上次结果进行增量优化

3. **上下文采集**：Agent 主动探索项目文件，提取技术栈、代码规范等元数据

4. **优化生成**：读取 `skills/templates/{lang}/{type}/{mode}.md`，替换占位符生成优化后的提示词

5. **深度评审**：使用 `critical-review.md` 模板检测歧义、边界盲区、逻辑冲突

6. **综合评估**：使用 `evaluation/user.md` 模板，结合评审报告生成最终评分（0-100）

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

Wails v2 应用，用于交互式确认优化结果。通过命令行参数传入输入文件，用户操作后输出结果文件。

```bash
prompt-optimizer-webview --input <input.json> --output <result.json> --timeout 600
```

### CI/CD

项目使用 GitHub Actions 进行持续集成和发布。

**CI 工作流** (`.github/workflows/ci.yml`)：
- 触发条件：push 到 main、PR 到 main
- 运行 Go 后端测试、前端测试、三平台构建验证

**发布工作流** (`.github/workflows/release.yml`)：
- 触发条件：推送 `v*` 标签
- 构建三平台版本，打包 `skills/` 目录（含二进制）
- 上传到 GitHub Releases

**发布产物**：
- `prompt-optimizer-skill-linux-amd64.tar.gz`
- `prompt-optimizer-skill-windows-amd64.zip`
- `prompt-optimizer-skill-darwin-universal.tar.gz`

```bash
# 发布新版本
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

**依赖更新** (`.github/dependabot.yml`)：
- 每周一自动检查 Go、npm、GitHub Actions 依赖更新

### 其他

禁止使用 **take_screenshot** tool，使用 **take_snapshot** 进行替代
