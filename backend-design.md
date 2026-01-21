# Prompt Optimizer Skill - 后端实现方案

## Session 管理

### 核心策略

**一次 `/optimize-prompt` 调用 = 一个完整 Session**

| 阶段 | 动作 |
|-----|------|
| `/optimize-prompt` 调用 | 创建 session，生成 session_id |
| 首次 WebView 弹出 | 使用 session 目录存储 input.json |
| submit/rollback | 更新 state.json，再次弹出 WebView |
| cancel/timeout | 关闭 WebView，Session 结束 |

### 目录结构

```
~/.prompt-optimizer/
├── sessions/                           # Session 数据目录
│   └── {session_id}/                   # 每个 session 一个目录
│       ├── meta.json                   # Session 元数据
│       ├── state.json                  # 当前状态 (核心文件)
│       ├── input.json                  # WebView 输入 (临时)
│       └── result.json                 # WebView 输出 (临时)
├── index.json                          # Session 索引
└── config.json                         # 全局配置 (可选)
```

### Session ID 生成

```
session_id = {project_hash}_{terminal_hash}_{timestamp}
```

| 组成部分 | 来源 | 作用 |
|---------|------|------|
| project_hash | pwd 的 SHA256 前 8 位 | 项目隔离 |
| terminal_hash | PPID+TTY 的 SHA256 前 8 位 | 终端隔离 |
| timestamp | Unix 毫秒时间戳 | 唯一性 |

**示例**: `a1b2c3d4_e5f6g7h8_1705632000000`

### state.json (核心状态文件)

```json
{
  "version": 3,
  "sessionId": "a1b2c3d4_e5f6g7h8_1705632000000",
  "originalPrompt": "用户的原始提示词",
  "current": {
    "iterationId": "iter-003",
    "optimizedPrompt": "...",
    "reviewReport": "...",
    "evaluationReport": "...",
    "score": 85,
    "suggestedDirections": [...]
  },
  "history": [...],
  "lastAction": "submit",
  "lastActionAt": "2024-01-19T10:30:00.000Z"
}
```

### meta.json (元数据)

```json
{
  "sessionId": "a1b2c3d4_e5f6g7h8_1705632000000",
  "createdAt": "2024-01-19T10:00:00.000Z",
  "updatedAt": "2024-01-19T10:30:00.000Z",
  "projectPath": "/Users/user/project",
  "status": "active",
  "currentIteration": 3
}
```

### 职责划分

| 组件 | 职责 |
|------|------|
| **SKILL.md (Agent)** | 创建 session，读写 state.json，构造 input.json |
| **WebView** | 展示数据，收集用户输入，返回 result.json (不参与状态管理) |

### 清理策略

| 规则 | 触发条件 | 动作 |
|------|---------|------|
| 过期清理 | 超过 30 天未更新 | 删除 session 目录 |
| 临时文件清理 | input/result.json 超过 1 小时 | 删除 |
| 容量限制 | sessions/ 超过 1GB | 删除最旧的 |

---

## 交互模式

> **核心约束**:
> - 每个终端的每个优化/迭代任务弹出独立的 **Wails webview 弹窗**
> - 一次 CLI 调用 = 一个 webview 实例
> - 交互完成后自动关闭 webview
> - 调用方通过**命令行参数 + 文件**传递数据，**同步阻塞**等待结果

---

## 技术栈与依赖

| 组件 | 技术选型 | 版本要求 | 说明 |
|------|----------|----------|------|
| 桌面框架 | Wails | v2.x | Go + webview |
| 后端语言 | Go | >= 1.21 | Wails 应用开发 |
| 前端 | HTML/CSS/JS | - | 复用 frontend-design.md |

### 分发方式

- **预编译 + 嵌入**: Wails 应用预编译为各平台二进制，嵌入到项目中
- 无需用户安装 Go 或 Wails CLI
- SKILL.md 通过 Shell 命令直接调用二进制

---

## 目录结构

```
prompt-optimizer-skill/
├── webui/
│   └── bin/                     # 预编译的 Wails 二进制
│       ├── prompt-optimizer-webview.exe   # Windows
│       ├── prompt-optimizer-webview       # Linux
│       └── prompt-optimizer-webview.app/  # macOS (应用包)
│
└── wails-app/                   # Wails 源码 (开发/构建用)
    ├── main.go                  # 入口 + CLI 参数解析
    ├── app.go                   # 绑定函数
    ├── wails.json               # Wails 配置
    ├── go.mod
    ├── go.sum
    └── frontend/                # 前端资源
        ├── index.html
        ├── src/
        │   ├── main.js
        │   └── styles.css
        └── dist/                # 构建产物
```

---

## 数据结构设计

### 输入文件格式 (input.json)

```json
{
  "version": 3,
  "originalPrompt": "用户最初输入的原始 prompt 文本",

  "current": {
    "iterationId": "iter-003",
    "optimizedPrompt": "# Role: 专家\n\n## Profile\n- Author: Expert\n- Version: 2.0\n\n## Background\n你是一个专业的...",
    "reviewReport": "## 深度评审报告\n\n### 优点\n1. 结构清晰\n2. 角色定位明确\n\n### 改进建议\n...",
    "evaluationReport": "## 评估报告\n\n### 评分明细\n| 维度 | 得分 |\n|------|------|\n| 清晰度 | 90 |\n| 完整性 | 80 |",
    "score": 85,
    "suggestedDirections": [
      {
        "id": "examples",
        "label": "添加示例",
        "description": "补充具体使用案例，让模型更好理解期望输出"
      },
      {
        "id": "constraints",
        "label": "增强约束",
        "description": "明确边界条件和限制，减少模型偏离预期"
      },
      {
        "id": "structure",
        "label": "优化结构",
        "description": "改进段落组织，提升可读性"
      }
    ]
  },

  "history": [
    {
      "iterationId": "iter-001",
      "optimizedPrompt": "第一次优化后的 prompt 文本...",
      "reviewReport": "第一次的评审报告...",
      "evaluationReport": "第一次的评估报告...",
      "score": 65,
      "userFeedback": {
        "selectedDirections": ["structure"],
        "userInput": "希望结构更清晰"
      }
    },
    {
      "iterationId": "iter-002",
      "optimizedPrompt": "第二次优化后的 prompt 文本...",
      "reviewReport": "第二次的评审报告...",
      "evaluationReport": "第二次的评估报告...",
      "score": 75,
      "userFeedback": {
        "selectedDirections": ["examples", "constraints"],
        "userInput": "需要更多示例和约束"
      }
    }
  ]
}
```

### 输入字段说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `version` | int | 是 | 当前迭代版本号 (从 1 开始) |
| `originalPrompt` | string | 是 | 用户最初输入的原始 prompt (仅展示，不可回滚) |
| `current` | object | 是 | 当前版本的优化数据 |
| `current.iterationId` | string | 是 | 迭代唯一标识 (如 "iter-001") |
| `current.optimizedPrompt` | string | 是 | 优化后的 prompt (Markdown) |
| `current.reviewReport` | string | 是 | 深度评审报告 (Markdown) |
| `current.evaluationReport` | string | 是 | 评估报告 (Markdown) |
| `current.score` | int | 是 | 评分 (0-100) |
| `current.suggestedDirections` | array | 是 | LLM 建议的优化方向 |
| `history` | array | 否 | 历史迭代记录 (按时间正序，最早的在前) |
| `history[].userFeedback` | object | 是 | 该版本用户提交的反馈 |

---

### 输出文件格式 (result.json)

#### 场景 1: 用户提交 (继续优化当前版本)

```json
{
  "action": "submit",
  "selectedDirections": ["examples", "constraints"],
  "userInput": "请增加 3 个具体的使用示例，并明确输入输出格式的约束"
}
```

#### 场景 2: 用户取消

```json
{
  "action": "cancel",
  "selectedDirections": [],
  "userInput": ""
}
```

#### 场景 3: 超时

```json
{
  "action": "timeout",
  "selectedDirections": [],
  "userInput": ""
}
```

#### 场景 4: 用户回滚到历史版本 (基于该版本继续优化)

```json
{
  "action": "rollback",
  "rollbackToIteration": "iter-001",
  "selectedDirections": ["examples"],
  "userInput": "基于第一版重新优化，这次重点添加示例"
}
```

### 输出字段说明

| 字段 | 类型 | 条件 | 说明 |
|------|------|------|------|
| `action` | string | 必填 | `"submit"` \| `"cancel"` \| `"timeout"` \| `"rollback"` |
| `selectedDirections` | array | action=submit/rollback | 用户选择的优化方向 ID 列表 |
| `userInput` | string | action=submit/rollback | 用户输入的补充说明 |
| `rollbackToIteration` | string | action=rollback | 回滚目标的 iterationId |

---

## 历史记录功能

### 功能概述

- **范围**: 当前 prompt 的迭代历史 (非全局历史)
- **存储**: 调用方 (SKILL.md) 管理，每次通过 input.json 传入
- **数量限制**: 不限制

### 用户可执行的操作

| 操作 | 说明 | result.json action |
|------|------|-------------------|
| 查看历史版本 | 浏览之前的迭代版本内容 | - (仅 UI 交互) |
| 版本对比 | 并排对比两个版本的差异 | - (仅 UI 交互) |
| 基于历史版本继续优化 | 选择某个历史版本作为新的起点 | `"rollback"` |

### UI 交互设计

#### 主界面 - 历史记录入口

```
┌─────────────────────────────────────────────────────────────────┐
│                      Prompt Optimizer                    [9:45] │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  [当前版本 v3]  [历史记录 ▼]  [版本对比]                          │
│       ↑              ↑            ↑                             │
│    Tab 切换     下拉选择历史   打开对比视图                       │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                                                           │  │
│  │  (Prompt / 评审报告 / 评估报告 内容区)                      │  │
│  │                                                           │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                 │
│  评分: ████████░░ 85/100                                        │
│                                                                 │
│  原始 Prompt: [查看原始输入]                                     │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### 历史记录下拉面板

```
┌─ 历史记录 ─────────────────────────────────────────────────────┐
│                                                                │
│  ● v3 (当前) - 85分                                            │
│    ├─ 用户反馈: (当前版本，待提交)                              │
│                                                                │
│  ○ v2 - 75分                                                   │
│    ├─ 选择方向: 添加示例, 增强约束                              │
│    └─ 用户输入: "需要更多示例和约束"                            │
│    [查看详情]  [基于此版本继续优化]                             │
│                                                                │
│  ○ v1 - 65分                                                   │
│    ├─ 选择方向: 优化结构                                       │
│    └─ 用户输入: "希望结构更清晰"                                │
│    [查看详情]  [基于此版本继续优化]                             │
│                                                                │
└────────────────────────────────────────────────────────────────┘
```

#### 版本对比视图

```
┌─────────────────────────────────────────────────────────────────┐
│  版本对比                                          [← 返回]     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  对比: [v1 ▼]  ↔  [v3 ▼]                                        │
│                                                                 │
├────────────────────────────┬────────────────────────────────────┤
│         v1 (65分)          │           v3 (85分)                │
├────────────────────────────┼────────────────────────────────────┤
│                            │                                    │
│ # Role: 助手               │ # Role: 专家                       │
│                            │                                    │
│ 你是一个AI助手...          │ ## Profile                         │
│                            │ - Author: Expert                   │
│                            │ - Version: 2.0                     │
│                            │                                    │
│                            │ ## Background                      │
│                            │ 你是一个专业的...                   │
│                            │                                    │
└────────────────────────────┴────────────────────────────────────┘
│                                                                 │
│                        [基于 v1 继续优化]                        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 核心模块职责

### 1. main.go - Wails 入口 + CLI 解析

**职责**: 解析命令行参数，初始化 Wails 应用，配置窗口参数

**CLI 参数**:
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `--input` | string | 是 | - | 输入 JSON 文件路径 |
| `--output` | string | 是 | - | 输出 JSON 文件路径 |
| `--timeout` | int | 否 | 600 | 超时时间 (秒) |

**窗口配置**:
- 窗口大小: 1000 x 700
- 窗口标题: "Prompt Optimizer"
- 窗口居中显示
- `OnBeforeClose`: 用户点击 X 按钮时返回 cancel

---

### 2. app.go - Wails 绑定函数

**职责**: 提供前端可调用的 Go 函数，管理数据读写

**结构体**:
```go
type App struct {
    ctx        context.Context
    inputData  InputData         // 从输入文件读取的完整数据
    outputPath string            // 输出文件路径
    timeout    int               // 超时时间 (秒)
    startTime  time.Time         // 启动时间
    writeOnce  sync.Once         // 保证 result.json 只写一次
}

// 输入数据结构
type InputData struct {
    Version        int              `json:"version"`
    OriginalPrompt string           `json:"originalPrompt"`
    Current        IterationData    `json:"current"`
    History        []HistoryItem    `json:"history"`
}

type IterationData struct {
    IterationId         string      `json:"iterationId"`
    OptimizedPrompt     string      `json:"optimizedPrompt"`
    ReviewReport        string      `json:"reviewReport"`
    EvaluationReport    string      `json:"evaluationReport"`
    Score               int         `json:"score"`
    SuggestedDirections []Direction `json:"suggestedDirections"`
}

type HistoryItem struct {
    IterationId      string       `json:"iterationId"`
    OptimizedPrompt  string       `json:"optimizedPrompt"`
    ReviewReport     string       `json:"reviewReport"`
    EvaluationReport string       `json:"evaluationReport"`
    Score            int          `json:"score"`
    UserFeedback     UserFeedback `json:"userFeedback"`
}

type UserFeedback struct {
    SelectedDirections []string `json:"selectedDirections"`
    UserInput          string   `json:"userInput"`
}

type Direction struct {
    ID          string `json:"id"`
    Label       string `json:"label"`
    Description string `json:"description"`
}

// 输出数据结构
type Result struct {
    Action              string   `json:"action"`
    SelectedDirections  []string `json:"selectedDirections"`
    UserInput           string   `json:"userInput"`
    RollbackToIteration string   `json:"rollbackToIteration,omitempty"`
}
```

**绑定方法**:
| 方法 | 说明 |
|------|------|
| `GetInputData() InputData` | 前端调用获取完整输入数据 (含历史) |
| `GetRemainingSeconds() int` | 获取剩余超时时间 (秒) |
| `Submit(directions []string, userInput string)` | 用户提交 (继续优化当前版本) |
| `Rollback(iterationId string, directions []string, userInput string)` | 用户回滚到历史版本继续优化 |
| `Cancel()` | 用户点击取消 |

---

## CLI 接口定义

### 调用方式

```bash
# 基本调用
./prompt-optimizer-webview --input /tmp/input.json --output /tmp/result.json

# 指定超时时间
./prompt-optimizer-webview --input /tmp/input.json --output /tmp/result.json --timeout 300
```

### Exit Code

| Exit Code | 含义 |
|-----------|------|
| 0 | 正常退出 (submit/cancel/timeout/rollback) |
| 1 | 输入文件不存在或读取失败 |
| 2 | JSON 解析失败 |
| 3 | 输出文件写入失败 |
| 4 | 其他错误 |

### 错误输出 (stderr)

解析失败或其他错误时，错误信息写入 stderr：
```
Error: failed to parse input file: invalid JSON at line 5
```

---

## 执行时序图

```
┌──────────┐     ┌──────────────────┐     ┌─────────┐
│ SKILL.md │     │  Wails 二进制     │     │  User   │
└────┬─────┘     └────────┬─────────┘     └────┬────┘
     │                    │                    │
     │ 写入 input.json    │                    │
     │ (含 history)       │                    │
     ├──────────┐         │                    │
     │          │         │                    │
     │<─────────┘         │                    │
     │                    │                    │
     │ ./webview --input input.json --output result.json
     ├───────────────────>│                    │
     │                    │                    │
     │                    │ 读取 input.json    │
     │                    ├──────────┐         │
     │                    │          │         │
     │                    │<─────────┘         │
     │                    │                    │
     │                    │ 创建窗口 (显示 loading)
     │                    ├──────────┐         │
     │                    │          │         │
     │                    │<─────────┘         │
     │                    │                    │
     │                    │ 显示 UI + 倒计时   │
     │                    ├───────────────────>│
     │                    │                    │
     │  [进程阻塞等待]     │                    │
     │  ................  │                    │
     │                    │                    │
     │                    │  [用户操作]         │
     │                    │  - 查看历史版本     │
     │                    │  - 版本对比         │
     │                    │  - 选择优化方向     │
     │                    │<───────────────────│
     │                    │                    │
     │                    │  点击确定/回滚      │
     │                    │<───────────────────│
     │                    │                    │
     │                    │ 显示成功动画 (0.5s) │
     │                    ├──────────┐         │
     │                    │          │         │
     │                    │<─────────┘         │
     │                    │                    │
     │                    │ 写入 result.json   │
     │                    ├──────────┐         │
     │                    │          │         │
     │                    │<─────────┘         │
     │                    │                    │
     │                    │ 窗口关闭           │
     │                    ├──────────┐         │
     │                    │          │         │
     │                    │<─────────┘         │
     │                    │                    │
     │  exit code 0       │                    │
     │<───────────────────│                    │
     │                    │                    │
     │ 读取 result.json   │                    │
     ├──────────┐         │                    │
     │          │         │                    │
     │<─────────┘         │                    │
     │                    │                    │
     │ 根据 action 处理:  │                    │
     │ - submit: 基于当前版本继续优化          │
     │ - rollback: 基于历史版本继续优化        │
     │ - cancel/timeout: 结束流程              │
```

---

## Go 代码实现

### main.go

```go
package main

import (
    "context"
    "embed"
    "flag"
    "fmt"
    "os"

    "github.com/wailsapp/wails/v2"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    // 解析命令行参数
    inputPath := flag.String("input", "", "Input JSON file path (required)")
    outputPath := flag.String("output", "", "Output JSON file path (required)")
    timeout := flag.Int("timeout", 600, "Timeout in seconds")
    flag.Parse()

    // 验证必填参数
    if *inputPath == "" || *outputPath == "" {
        fmt.Fprintln(os.Stderr, "Error: --input and --output are required")
        flag.Usage()
        os.Exit(4)
    }

    // 创建 App 实例
    app, err := NewApp(*inputPath, *outputPath, *timeout)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // 启动 Wails 应用
    err = wails.Run(&options.App{
        Title:     "Prompt Optimizer",
        Width:     1000,
        Height:    700,
        MinWidth:  800,
        MinHeight: 600,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        OnStartup: app.startup,
        OnBeforeClose: func(ctx context.Context) bool {
            app.Cancel()
            return false
        },
        Bind: []interface{}{
            app,
        },
    })

    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(4)
    }
}
```

### app.go

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "sync"
    "time"

    "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
    ctx        context.Context
    inputData  InputData
    outputPath string
    timeout    int
    startTime  time.Time
}

// 输入数据结构
type InputData struct {
    Version        int           `json:"version"`
    OriginalPrompt string        `json:"originalPrompt"`
    Current        IterationData `json:"current"`
    History        []HistoryItem `json:"history"`
}

type IterationData struct {
    IterationId         string      `json:"iterationId"`
    OptimizedPrompt     string      `json:"optimizedPrompt"`
    ReviewReport        string      `json:"reviewReport"`
    EvaluationReport    string      `json:"evaluationReport"`
    Score               int         `json:"score"`
    SuggestedDirections []Direction `json:"suggestedDirections"`
}

type HistoryItem struct {
    IterationId      string       `json:"iterationId"`
    OptimizedPrompt  string       `json:"optimizedPrompt"`
    ReviewReport     string       `json:"reviewReport"`
    EvaluationReport string       `json:"evaluationReport"`
    Score            int          `json:"score"`
    UserFeedback     UserFeedback `json:"userFeedback"`
}

type UserFeedback struct {
    SelectedDirections []string `json:"selectedDirections"`
    UserInput          string   `json:"userInput"`
}

type Direction struct {
    ID          string `json:"id"`
    Label       string `json:"label"`
    Description string `json:"description"`
}

// 输出数据结构
type Result struct {
    Action              string   `json:"action"`
    SelectedDirections  []string `json:"selectedDirections"`
    UserInput           string   `json:"userInput"`
    RollbackToIteration string   `json:"rollbackToIteration,omitempty"`
}

func NewApp(inputPath, outputPath string, timeout int) (*App, error) {
    // 读取输入文件
    data, err := os.ReadFile(inputPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read input file: %w", err)
    }

    var inputData InputData
    if err := json.Unmarshal(data, &inputData); err != nil {
        return nil, fmt.Errorf("failed to parse input JSON: %w", err)
    }

    return &App{
        inputData:  inputData,
        outputPath: outputPath,
        timeout:    timeout,
        startTime:  time.Now(),
    }, nil
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx

    // 启动超时监控 goroutine
    go a.watchTimeout()
}

func (a *App) watchTimeout() {
    time.Sleep(time.Duration(a.timeout) * time.Second)
    a.writeResult(Result{
        Action:             "timeout",
        SelectedDirections: []string{},
        UserInput:          "",
    })
    runtime.Quit(a.ctx)
}

// GetInputData 前端调用获取完整输入数据 (含历史)
func (a *App) GetInputData() InputData {
    return a.inputData
}

// GetRemainingSeconds 获取剩余超时时间
func (a *App) GetRemainingSeconds() int {
    elapsed := int(time.Since(a.startTime).Seconds())
    remaining := a.timeout - elapsed
    if remaining < 0 {
        return 0
    }
    return remaining
}

// Submit 用户提交 (继续优化当前版本)
func (a *App) Submit(directions []string, userInput string) {
    a.writeResult(Result{
        Action:             "submit",
        SelectedDirections: directions,
        UserInput:          userInput,
    })
    runtime.Quit(a.ctx)
}

// Rollback 用户回滚到历史版本继续优化
func (a *App) Rollback(iterationId string, directions []string, userInput string) {
    a.writeResult(Result{
        Action:              "rollback",
        SelectedDirections:  directions,
        UserInput:           userInput,
        RollbackToIteration: iterationId,
    })
    runtime.Quit(a.ctx)
}

// Cancel 用户点击取消
func (a *App) Cancel() {
    a.writeResult(Result{
        Action:             "cancel",
        SelectedDirections: []string{},
        UserInput:          "",
    })
    runtime.Quit(a.ctx)
}

func (a *App) writeResult(result Result) error {
    var writeErr error
    a.writeOnce.Do(func() {
        data, err := json.MarshalIndent(result, "", "  ")
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: failed to marshal result: %v\n", err)
            writeErr = err
            return
        }

        // 原子写入：先写临时文件，再 rename，避免半截文件
        tmpPath := a.outputPath + ".tmp"
        if err := os.WriteFile(tmpPath, data, 0644); err != nil {
            fmt.Fprintf(os.Stderr, "Error: failed to write temp file: %v\n", err)
            writeErr = err
            return
        }

        if err := os.Rename(tmpPath, a.outputPath); err != nil {
            fmt.Fprintf(os.Stderr, "Error: failed to rename output file: %v\n", err)
            writeErr = err
            return
        }
    })
    return writeErr
}
```

---

## 前端适配 (Wails 绑定)

### main.js

```javascript
// 使用 Wails 运行时调用 Go 绑定函数

let countdownInterval = null;
let inputData = null;

// 获取初始数据
async function loadData() {
    try {
        inputData = await window.go.main.App.GetInputData();
        hideLoading();
        renderUI(inputData);
        startCountdown();
    } catch (error) {
        showError("加载数据失败: " + error.message);
    }
}

// 渲染主界面
function renderUI(data) {
    // 渲染当前版本
    renderCurrentVersion(data.current);

    // 渲染历史记录
    renderHistory(data.history);

    // 渲染原始 prompt (仅展示)
    renderOriginalPrompt(data.originalPrompt);

    // 渲染优化方向选择
    renderDirections(data.current.suggestedDirections);
}

// 渲染历史记录列表
function renderHistory(history) {
    const historyList = document.getElementById('history-list');
    historyList.innerHTML = '';

    history.forEach((item, index) => {
        const version = index + 1;
        const element = createHistoryItem(item, version);
        historyList.appendChild(element);
    });
}

// 查看历史版本详情
function viewHistoryVersion(iterationId) {
    const item = inputData.history.find(h => h.iterationId === iterationId);
    if (item) {
        showVersionDetail(item);
    }
}

// 版本对比
function compareVersions(leftId, rightId) {
    const left = findVersion(leftId);
    const right = findVersion(rightId);
    showComparisonView(left, right);
}

// 启动倒计时
async function startCountdown() {
    updateCountdown();
    countdownInterval = setInterval(updateCountdown, 1000);
}

async function updateCountdown() {
    const remaining = await window.go.main.App.GetRemainingSeconds();
    const minutes = Math.floor(remaining / 60);
    const seconds = remaining % 60;

    const display = document.getElementById('countdown');
    display.textContent = `${minutes}:${seconds.toString().padStart(2, '0')}`;

    // 小于 60 秒时高亮警告
    if (remaining < 60) {
        display.classList.add('warning');
    }

    if (remaining <= 0) {
        clearInterval(countdownInterval);
    }
}

// 提交 - 继续优化当前版本 (带成功动画)
async function handleSubmit() {
    const submitBtn = document.getElementById('submit-btn');
    submitBtn.disabled = true;

    // 显示成功动画
    showSuccessAnimation();

    // 0.5 秒后调用 Submit
    setTimeout(async () => {
        const selectedDirections = getSelectedDirections();
        const userInput = document.getElementById('user-input').value;
        await window.go.main.App.Submit(selectedDirections, userInput);
    }, 500);
}

// 回滚 - 基于历史版本继续优化
async function handleRollback(iterationId) {
    const submitBtn = document.getElementById('rollback-btn');
    submitBtn.disabled = true;

    // 显示成功动画
    showSuccessAnimation();

    // 0.5 秒后调用 Rollback
    setTimeout(async () => {
        const selectedDirections = getSelectedDirections();
        const userInput = document.getElementById('user-input').value;
        await window.go.main.App.Rollback(iterationId, selectedDirections, userInput);
    }, 500);
}

// 取消
async function handleCancel() {
    await window.go.main.App.Cancel();
}

// 页面加载时显示 loading 并获取数据
document.addEventListener('DOMContentLoaded', () => {
    showLoading();
    loadData();
});
```

### 前端状态管理

#### Loading 状态
- 初始 HTML 包含 loading 骨架屏
- 调用 `GetInputData()` 成功后切换到内容视图

#### 超时倒计时
- 从 `GetRemainingSeconds()` 获取剩余时间
- 每秒更新一次显示
- 剩余时间 < 60s 时添加 `.warning` 样式高亮

#### 提交成功动画
- 点击确定/回滚后立即禁用按钮 (防重复提交)
- 显示 checkmark 成功动画
- 0.5s 后调用对应方法关闭窗口

#### 历史记录交互
- 点击历史版本可查看详情
- 支持任意两个版本对比
- 可基于任意历史版本继续优化 (触发 rollback)

---

## 窗口行为

| 场景 | 行为 |
|------|------|
| 用户点击确定 | 禁用按钮 → 成功动画 (0.5s) → `Submit()` → 写入 result.json → 窗口关闭 |
| 用户点击回滚 | 禁用按钮 → 成功动画 (0.5s) → `Rollback()` → 写入 result.json → 窗口关闭 |
| 用户点击取消 | `Cancel()` → 写入 result.json → 窗口关闭 |
| 用户点击 X 按钮 | `OnBeforeClose` → `Cancel()` → 写入 result.json → 窗口关闭 |
| 超时 | `watchTimeout()` goroutine → 写入 result.json → 窗口关闭 |
| 输入文件解析失败 | 写入 stderr → exit code 1 或 2 |

---

## 构建与分发

### 构建命令

```bash
# Windows
wails build -platform windows/amd64

# macOS
wails build -platform darwin/universal

# Linux
wails build -platform linux/amd64
```

### CI/CD 构建流程

```yaml
# .github/workflows/build.yml
name: Build Wails App

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    strategy:
      matrix:
        os: [windows-latest, macos-latest, ubuntu-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Install Wails
        run: go install github.com/wailsapp/wails/v2/cmd/wails@latest
      - name: Build
        run: cd wails-app && wails build
      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: webview-${{ runner.os }}
          path: wails-app/build/bin/
```

---

## 实现优先级

| 优先级 | 模块 | 说明 |
|--------|------|------|
| P0 | `main.go` | CLI 参数解析 + Wails 初始化 |
| P0 | `app.go` | 文件读写 + 绑定函数 (含历史) |
| P0 | `frontend/` | 前端 UI (复用 frontend-design.md) |
| P0 | 数据结构 | input.json / result.json 完整实现 |
| P1 | 历史记录 | 查看历史 + 版本对比 + 回滚 |
| P1 | 超时倒计时 | 前端显示 + Go 端监控 |
| P1 | Loading 状态 | 骨架屏 + 数据加载 |
| P1 | 成功动画 | 提交反馈 + 防重复点击 |
| P1 | `OnBeforeClose` | 用户关闭窗口处理 |
| P2 | 窗口置顶 | 确保新窗口获得焦点 |
| P2 | 跨平台构建 | CI/CD 配置 |
