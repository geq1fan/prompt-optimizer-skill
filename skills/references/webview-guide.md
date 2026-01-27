# WebView 交互指南

本文档提供 WebView 桌面应用的完整使用指南，这是提示词优化器的核心交互组件。

## 目录

- [概述](#概述)
- [调用流程](#调用流程)
- [输入格式详解](#输入格式详解)
- [输出格式详解](#输出格式详解)
- [完整示例](#完整示例)
- [错误处理](#错误处理)

---

## 概述

WebView 应用提供图形化界面，让用户可以：
- 查看优化后的提示词
- 阅读评审报告和评估报告
- 选择预设的改进方向
- 输入自定义改进指令
- 查看历史版本并回滚

---

## 调用流程

### 1. 确定二进制路径

```bash
# 相对于 skill 安装目录
SKILL_DIR=$(dirname "$0")  # 或从环境变量获取

# Windows
WEBVIEW_BIN="$SKILL_DIR/bin/prompt-optimizer-webview.exe"

# macOS
WEBVIEW_BIN="$SKILL_DIR/bin/prompt-optimizer-webview.app/Contents/MacOS/prompt-optimizer-webview"

# Linux
WEBVIEW_BIN="$SKILL_DIR/bin/prompt-optimizer-webview"
```

### 2. 准备 Session 目录

Session 目录由 WebView 自动创建，无需手动创建：

```bash
# Session ID 格式简化为时间戳
session_id="session_{timestamp}"

# Session 目录路径（相对于项目根目录）
SESSION_DIR=.claude/prompt-optimizer/sessions/{session_id}
SESSION_FILE="$SESSION_DIR/session.json"
OUTPUT_FILE="$SESSION_DIR/result.json"

# Agent 创建 session.json (见下文格式)
```

### 3. 调用 WebView

使用 `--session-id` 参数调用 WebView，目录和路径自动推断：

```bash
"$WEBVIEW_BIN" --session-id "$session_id"
```

可选参数：
- `--timeout 600` - 超时秒数（默认 600）

### 4. 读取结果

```bash
# 等待进程退出后读取结果
RESULT=$(cat "$OUTPUT_FILE")

# WebView 已自动更新 session.json
# Agent 读取更新后的 session.json 继续迭代
```

---

## 输入格式详解

### 完整结构 (v4)

```json
{
  "version": 4,
  "sessionId": "session_1705632000000",
  "createdAt": "2026-01-22T10:00:00Z",
  "updatedAt": "2026-01-22T10:00:00Z",
  "projectPath": "/path/to/project",
  "lang": "cn",
  "mode": "professional",
  "originalPrompt": "用户的原始提示词",
  "current": {
    "iterationId": "iter-001",
    "optimizedPrompt": "优化后的提示词 (Markdown 格式)",
    "reviewReport": "深度评审报告 (Markdown 格式)",
    "evaluationReport": "综合评估报告 (Markdown 格式)",
    "score": 85,
    "suggestedDirections": [
      {
        "id": "examples",
        "label": "添加示例",
        "description": "补充具体的使用案例"
      },
      {
        "id": "constraints",
        "label": "明确约束",
        "description": "添加更多限制条件"
      }
    ]
  },
  "history": [
    {
      "iterationId": "iter-000",
      "optimizedPrompt": "初版优化结果",
      "reviewReport": "初版评审报告",
      "evaluationReport": "初版评估报告",
      "score": 72,
      "userFeedback": {
        "selectedDirections": ["structure"],
        "userInput": "第一次用户输入"
      }
    }
  ],
  "lastAction": null,
  "status": "active"
}
```

### 字段说明

#### 根字段

| 字段 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| `version` | number | 是 | 协议版本，当前为 4 |
| `sessionId` | string | 是 | 唯一标识符，格式 `session_{timestamp}` |
| `createdAt` | string | 是 | ISO 8601 时间戳 |
| `updatedAt` | string | 是 | ISO 8601 时间戳 |
| `projectPath` | string | 是 | 项目根目录路径 |
| `lang` | string | 是 | "cn" 或 "en" |
| `mode` | string | 是 | "basic"/"professional"/"planning" |
| `originalPrompt` | string | 是 | 用户最初的原始提示词 |
| `current` | object | 是 | 当前迭代的数据 |
| `history` | array | 是 | 历史迭代记录（初始为空） |
| `lastAction` | string | 否 | 上次操作类型 |
| `status` | string | 是 | "active"/"completed"/"cancelled" |

#### current 对象

| 字段 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| `iterationId` | string | 是 | 迭代标识，格式 `iter-XXX` |
| `optimizedPrompt` | string | 是 | 优化后的提示词，支持 Markdown |
| `reviewReport` | string | 是 | 深度评审报告，支持 Markdown |
| `evaluationReport` | string | 是 | 综合评估报告，支持 Markdown |
| `score` | number | 是 | 评分 (0-100) |
| `suggestedDirections` | array | 否 | 建议的改进方向列表 |

#### suggestedDirections 数组项

| 字段 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| `id` | string | 是 | 唯一标识符 |
| `label` | string | 是 | 显示标签 (简短) |
| `description` | string | 是 | 详细描述 |

#### history 数组项

与 current 结构类似，额外包含：

| 字段 | 类型 | 说明 |
|-----|------|------|
| `userFeedback` | object | 用户反馈信息 |
| `userFeedback.selectedDirections` | array | 用户选择的改进方向 ID |
| `userFeedback.userInput` | string | 用户输入的自定义改进指令 |

---

## 输出格式详解

### 用户确认 (submit)

用户点击"确认"或"继续优化"按钮：

```json
{
  "action": "submit",
  "selectedDirections": ["examples", "constraints"],
  "userInput": "请添加 3 个具体的代码示例"
}
```

| 字段 | 类型 | 说明 |
|-----|------|------|
| `action` | string | 固定为 `"submit"` |
| `selectedDirections` | array | 用户选择的改进方向 ID 列表 |
| `userInput` | string | 用户输入的自定义改进指令 |

**WebView 自动更新 session.json**：
1. 将 `current` 移入 `history`（附带 `userFeedback`）
2. 清空 `current`
3. 更新 `lastAction` 为 "submit"
4. 更新 `updatedAt`

**Agent 后续处理**：
1. 从 history 最后一项读取 `userFeedback`
2. 将 `selectedDirections` 转换为改进描述
3. 结合 `userInput` 作为 `iterateInput`
4. 进入下一轮迭代

### 用户取消 (cancel)

用户点击"取消"或关闭窗口：

```json
{
  "action": "cancel",
  "selectedDirections": [],
  "userInput": ""
}
```

**后续处理**：
1. 结束优化流程
2. 更新 `status` 为 "cancelled"

### 超时 (timeout)

用户超过指定时间未操作：

```json
{
  "action": "timeout",
  "selectedDirections": [],
  "userInput": ""
}
```

**后续处理**：
1. 提示用户已超时
2. 询问是否继续或结束

### 回滚 (rollback)

用户选择回滚到历史版本：

```json
{
  "action": "rollback",
  "rollbackToIteration": "iter-001",
  "selectedDirections": ["examples"],
  "userInput": "基于这个版本重新优化"
}
```

| 字段 | 类型 | 说明 |
|-----|------|------|
| `rollbackToIteration` | string | 目标迭代的 ID |
| `selectedDirections` | array | 用户选择的改进方向（可选） |
| `userInput` | string | 用户输入的改进指令（可选） |

**WebView 自动更新 session.json**：
1. 从 history 找到目标迭代
2. 截断 history 到目标迭代（包含）
3. 将目标迭代附带 `userFeedback` 保留在 history 末尾
4. 清空 `current`
5. 更新 `lastAction` 为 "rollback"

**Agent 后续处理**：
1. 从 history 最后一项读取迭代数据和 `userFeedback`
2. 基于该版本继续迭代

---

## 完整示例

### 场景：首次优化

**输入 (session.json)**：

```json
{
  "version": 4,
  "sessionId": "session_1705632000000",
  "createdAt": "2026-01-22T10:00:00Z",
  "updatedAt": "2026-01-22T10:00:00Z",
  "projectPath": "/Users/demo/my-project",
  "lang": "cn",
  "mode": "professional",
  "originalPrompt": "帮我写个登录功能",
  "current": {
    "iterationId": "iter-001",
    "optimizedPrompt": "## 任务\n\n实现用户登录功能\n\n## 要求\n\n1. 支持用户名/密码登录\n2. 实现 JWT token 认证\n3. 添加登录失败次数限制\n\n## 技术栈\n\n- 后端: Node.js + Express\n- 数据库: PostgreSQL\n\n## 输出\n\n- 完整的登录 API 代码\n- 单元测试\n- API 文档",
    "reviewReport": "## 评审摘要\n\n**健壮性**: 需完善\n\n### 发现的问题\n\n1. 未指定密码加密方式\n2. JWT 过期时间未明确\n3. 登录失败次数限制的具体数值未说明",
    "evaluationReport": "## 评估报告\n\n**总分**: 78/100 (中等)\n\n### 维度评分\n\n- 任务表达: 85/100\n- 信息完整性: 70/100\n- 格式规范性: 80/100\n- 改进程度: 75/100",
    "score": 78,
    "suggestedDirections": [
      {
        "id": "security",
        "label": "增强安全性",
        "description": "明确密码加密、JWT 配置等安全细节"
      },
      {
        "id": "error-handling",
        "label": "错误处理",
        "description": "添加各种异常情况的处理要求"
      }
    ]
  },
  "history": [],
  "lastAction": null,
  "status": "active"
}
```

**用户操作**：选择"增强安全性"并输入"使用 bcrypt 加密，JWT 有效期 1 小时"

**输出 (result.json)**：

```json
{
  "action": "submit",
  "selectedDirections": ["security"],
  "userInput": "使用 bcrypt 加密，JWT 有效期 1 小时"
}
```

**更新后的 session.json**：

```json
{
  "version": 4,
  "sessionId": "session_1705632000000",
  "createdAt": "2026-01-22T10:00:00Z",
  "updatedAt": "2026-01-22T10:05:00Z",
  "projectPath": "/Users/demo/my-project",
  "lang": "cn",
  "mode": "professional",
  "originalPrompt": "帮我写个登录功能",
  "current": null,
  "history": [
    {
      "iterationId": "iter-001",
      "optimizedPrompt": "...(同上)...",
      "reviewReport": "...",
      "evaluationReport": "...",
      "score": 78,
      "userFeedback": {
        "selectedDirections": ["security"],
        "userInput": "使用 bcrypt 加密，JWT 有效期 1 小时"
      }
    }
  ],
  "lastAction": "submit",
  "status": "active"
}
```

### 场景：迭代优化后回滚

**输入**：包含多个 history 迭代记录的 session.json

**用户操作**：查看历史版本，决定回滚到 iter-001

**输出 (result.json)**：

```json
{
  "action": "rollback",
  "rollbackToIteration": "iter-001",
  "selectedDirections": [],
  "userInput": "换个方向，专注于性能优化"
}
```

---

## 错误处理

### WebView 应用不存在

```bash
if [ ! -f "$WEBVIEW_BIN" ]; then
    echo "WebView 应用未找到，跳过交互确认"
    # 直接输出结果，不进行交互
fi
```

### 进程异常退出

检查退出码：

```bash
"$WEBVIEW_BIN" --session-id "$session_id"
EXIT_CODE=$?

if [ $EXIT_CODE -ne 0 ]; then
    echo "WebView 异常退出，退出码: $EXIT_CODE"
    # 回退到非交互模式
fi
```

---

## 命令行参数速查

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|-----|------|------|-------|------|
| `--session-id` | string | 是 | - | Session ID，用于推断文件路径 |
| `--timeout` | number | 否 | 600 | 超时秒数 |

**路径推断规则**：
- 输入文件：`.claude/prompt-optimizer/sessions/{session-id}/session.json`
- 输出文件：`.claude/prompt-optimizer/sessions/{session-id}/result.json`
- 目录不存在时自动创建

---

## Session 目录约定

### 目录结构

```
{project-root}/
└── .claude/
    └── prompt-optimizer/
        └── sessions/                       # Session 数据目录
            └── session_1705632000000/      # 每个 session 一个目录
                ├── session.json            # 核心状态文件（唯一真相源）
                └── result.json             # WebView 输出（临时文件）
```

### Session ID 格式

```
session_id = session_{timestamp}

示例: session_1705632000000
```

**生成方式**:
```bash
session_id="session_$(date +%s%3N)"
```

### 文件职责

| 文件 | 创建者 | 读取者 | 写入者 | 生命周期 |
|------|--------|--------|--------|----------|
| `session.json` | Agent | Agent + WebView | Agent + WebView | Session 存活期间 |
| `result.json` | WebView | Agent | WebView | 每次交互后被覆盖 |

### 清理策略

过期 Session 可按以下规则清理：

```bash
# 删除 30 天未更新的 Session
find .claude/prompt-optimizer/sessions -type d -mtime +30 -exec rm -rf {} \;
```

---

## 向后兼容

WebView 应用支持 v1-v3 格式的自动迁移。加载旧版本数据时，会自动填充 v4 新增字段：

| 字段 | 迁移策略 |
|------|---------|
| `sessionId` | 从文件路径推断或生成新 ID |
| `createdAt` | 使用当前时间 |
| `updatedAt` | 使用当前时间 |
| `projectPath` | 使用当前工作目录 |
| `lang` | 根据 originalPrompt 内容检测 |
| `mode` | 默认 "basic" |
| `status` | 默认 "active" |
