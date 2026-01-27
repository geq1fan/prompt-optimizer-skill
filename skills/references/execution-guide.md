# 执行指南

本文档提供提示词优化的完整执行流程，供 Agent 逐步遵循。

## 目录

- [Phase 1: Session 创建](#phase-1-session-创建)
- [Phase 2: 模式选择](#phase-2-模式选择)
- [Phase 3: 上下文采集](#phase-3-上下文采集)
- [Phase 4: 优化生成](#phase-4-优化生成)
- [Phase 5: 深度评审](#phase-5-深度评审)
- [Phase 6: 综合评估](#phase-6-综合评估)
- [Phase 7: 输出展示](#phase-7-输出展示)
- [Phase 8: WebView 交互循环](#phase-8-webview-交互循环)

---

## Phase 1: Session 创建

### 触发条件

```
用户输入: /optimize-prompt [内容]
```

每次 `/optimize-prompt` 调用都会创建新 Session。

### Session ID 生成

```
session_id = session_{timestamp}

示例: session_1705632000000
```

**生成方式**:
```bash
session_id="session_$(date +%s%3N)"
```

### 目录初始化

Session 目录由 WebView 自动创建，无需手动创建：

```
.claude/prompt-optimizer/sessions/{session_id}/
```

### 清理旧 Session（可选）

在创建新 Session 前，可清理过期数据：

```bash
# 删除 30 天未更新的 Session
find .claude/prompt-optimizer/sessions -type d -mtime +30 -exec rm -rf {} \;
```

### 语言检测

```
lang: "cn" | "en"  # 基于输入是否包含中文字符
```

检测规则：包含中文字符使用 `cn`，否则使用 `en`。

---

## Phase 2: 模式选择

### 判定流程

根据任务内容选择优化模板：

```
1. 分析用户输入内容
   ├─ 涉及步骤、计划、流程、roadmap → planning
   ├─ 涉及代码、分析、学术、技术文档 → professional
   └─ 其他通用场景 → basic
```

### 模式对应模板

| 特征 | 模式 | 模板 |
|-----|------|-----|
| 涉及步骤、计划、流程、roadmap | planning | `user-optimize/planning.md` |
| 涉及代码、分析、学术、技术文档 | professional | `user-optimize/professional.md` |
| 其他通用场景 | basic | `user-optimize/basic.md` |

### 输出变量

```
mode: "basic" | "professional" | "planning"
lang: "cn" | "en"
```

---

## Phase 3: 上下文采集

### 判断条件

任务是否依赖当前项目的具体信息？

### 执行逻辑

```
IF 任务涉及代码/架构/特定文件:
    1. 使用 glob 扫描相关文件结构
    2. 使用 read 读取关键代码/配置
    3. 使用 grep 搜索特定模式
    4. 提取元数据:
       - 技术栈 (语言、框架、版本)
       - 代码规范 (命名、结构)
       - 依赖关系
ELSE:
    跳过此阶段
```

### 目的

确保生成的 Prompt 基于项目事实，而非通用假设。

---

## Phase 4: 优化生成

### 模板路径

```
templates/{lang}/user-optimize/{mode}.md

示例:
- templates/cn/user-optimize/basic.md
- templates/cn/user-optimize/professional.md
- templates/en/user-optimize/planning.md
```

### 占位符替换

**首次优化**:
| 占位符 | 值 |
|-------|---|
| `{{originalPrompt}}` | 用户原始输入 |

**迭代优化**（WebView submit/rollback 触发）:
| 占位符 | 值 |
|-------|---|
| `{{lastOptimizedPrompt}}` | 上次优化结果（从 session.json 读取） |
| `{{lastEvaluationReport}}` | 上次评估报告（从 session.json 读取） |
| `{{iterateInput}}` | 用户改进指令（从 result.json 读取） |

### 执行步骤

1. 读取模板文件
2. 替换占位符
3. 按模板指令生成优化后的提示词
4. 存储结果为 `optimizedPrompt`

---

## Phase 5: 深度评审

### 模板

```
templates/{lang}/evaluation/critical-review.md
```

### 占位符

| 占位符 | 值 |
|-------|---|
| `{{optimizedPrompt}}` | Phase 4 生成的优化结果 |

### 检测维度

1. **歧义表达** - 可能被误解的措辞
2. **边界盲区** - 未覆盖的边界情况
3. **逻辑冲突** - 相互矛盾的要求

### 输出

生成 `reviewReport`，包含发现的问题和风险评估。

---

## Phase 6: 综合评估

### 模板

```
templates/{lang}/evaluation/user.md
```

### 占位符

| 占位符 | 值 |
|-------|---|
| `{{optimizedPrompt}}` | 优化后的提示词 |
| `{{reviewReport}}` | Phase 5 的评审报告 |
| `{{hasOriginalPrompt}}` | 布尔值，用于 Mustache 条件渲染 |

### 评分维度

| 维度 | 权重 | 评估点 |
|-----|------|-------|
| 任务表达 | 30% | 意图和目标是否清晰 |
| 信息完整性 | 30% | 关键信息和约束是否齐全 |
| 格式规范性 | 20% | 结构是否清晰易懂 |
| 改进程度 | 20% | 相比原始的提升程度 |

### 评分等级

| 分数 | 等级 |
|-----|------|
| 90-100 | 优秀 |
| 80-89 | 良好 |
| 70-79 | 中等 |
| 60-69 | 及格 |
| 0-59 | 不及格 |

### 输出

生成 `evaluationReport`，包含分数和详细评价。

---

## Phase 7: 输出展示

### 输出格式

简要提示优化完成，不输出完整内容：

```
优化完成，评分: XX/100 ([等级])
正在打开交互界面...
```

**注意**: 所有详细信息（优化后提示词、评审报告、评估报告）通过 WebView 展示，不在终端输出完整的提示词和报告。

### 创建 session.json

在 Session 目录下创建 `session.json`（v4 格式）：

```json
{
  "version": 4,
  "sessionId": "session_1705632000000",
  "createdAt": "2026-01-22T10:00:00Z",
  "updatedAt": "2026-01-22T10:00:00Z",
  "projectPath": "/path/to/project",
  "lang": "cn",
  "mode": "professional",
  "originalPrompt": "用户原始输入",
  "current": {
    "iterationId": "iter-001",
    "optimizedPrompt": "[optimizedPrompt]",
    "reviewReport": "[reviewReport]",
    "evaluationReport": "[evaluationReport]",
    "score": 85,
    "suggestedDirections": [
      {"id": "xxx", "label": "标签", "description": "描述"}
    ]
  },
  "history": [],
  "lastAction": null,
  "status": "active"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `version` | number | 是 | 协议版本，当前为 4 |
| `sessionId` | string | 是 | 唯一标识符 |
| `createdAt` | string | 是 | ISO 8601 时间戳 |
| `updatedAt` | string | 是 | ISO 8601 时间戳 |
| `projectPath` | string | 是 | 项目根目录路径 |
| `lang` | string | 是 | "cn" 或 "en" |
| `mode` | string | 是 | "basic"/"professional"/"planning" |
| `originalPrompt` | string | 是 | 用户原始输入 |
| `current` | object | 是 | 当前迭代数据 |
| `history` | array | 是 | 历史记录（初始为空） |
| `lastAction` | string | 否 | 上次操作类型 |
| `status` | string | 是 | "active"/"completed"/"cancelled" |

---

## Phase 8: WebView 交互循环 ⚠️ 必须执行

**重要**: 此步骤是必须的用户确认环节，不可跳过。每次优化完成后必须调用 WebView，等待用户做出选择。

### 文件路径

WebView 自动使用项目目录下的 `.claude/prompt-optimizer/sessions/{session_id}/`：

```
SESSION_DIR=.claude/prompt-optimizer/sessions/{session_id}
SESSION_FILE="$SESSION_DIR/session.json"
OUTPUT_FILE="$SESSION_DIR/result.json"
```

### 调用命令

使用 `--session-id` 参数调用 WebView，目录自动创建：

```bash
# Windows
bin/prompt-optimizer-webview.exe --session-id {session_id}

# macOS
bin/prompt-optimizer-webview.app/Contents/MacOS/prompt-optimizer-webview --session-id {session_id}

# Linux
bin/prompt-optimizer-webview --session-id {session_id}
```

可选参数：
- `--timeout 600` - 超时秒数（默认 600）

### 结果处理

读取 `result.json`，根据 action 执行后续操作：

| action | 后续操作 |
|--------|---------|
| `submit` | 1. WebView 已更新 session.json（current 移入 history）<br>2. 从 history 最后一项读取 userFeedback<br>3. 返回 Phase 4 开始新一轮迭代 |
| `rollback` | 1. WebView 已更新 session.json（回滚到指定版本）<br>2. 读取 userFeedback 和 rollbackToIteration<br>3. 返回 Phase 4 开始新一轮迭代 |
| `cancel` | 结束 Session，更新 status 为 "cancelled" |
| `timeout` | 提示用户超时，结束 Session |

### Session 状态更新

WebView 会自动更新 `session.json`：

**submit 时**:
```json
{
  "current": null,
  "history": [
    ...existingHistory,
    {
      ...previousCurrent,
      "userFeedback": {
        "selectedDirections": ["id1", "id2"],
        "userInput": "用户输入"
      }
    }
  ],
  "lastAction": "submit",
  "updatedAt": "2026-01-22T10:05:00Z"
}
```

**rollback 时**:
```json
{
  "current": null,
  "history": [
    ...trimmedHistory,
    {
      ...rollbackTarget,
      "userFeedback": {
        "selectedDirections": ["id1"],
        "userInput": "用户输入"
      }
    }
  ],
  "lastAction": "rollback",
  "updatedAt": "2026-01-22T10:05:00Z"
}
```

### ⚠️ 重要警告

- **每次优化完成后必须调用 WebView** - 不要直接结束流程
- **等待用户在 WebView 中做出选择** - submit/rollback/cancel/timeout
- **仅当 WebView 二进制不存在时才可跳过** - 此时直接输出结果
- **WebView 会自动更新 session.json** - Agent 读取更新后的数据继续迭代

---

## 继续优化的判断逻辑

### 对话中记住 Session

Agent 在对话上下文中记住当前 `session_id`，用于判断后续请求是否为继续优化：

```
current_session_id: "session_1705632000000"
session_dir: "~/.prompt-optimizer/sessions/session_1705632000000/"
```

### 判断流程

```
用户输入
    │
    ▼
┌─────────────────────┐
│ 包含 /optimize-prompt │──Yes──▶ 新建 Session（Phase 1 开始）
└─────────────────────┘
    │ No
    ▼
┌─────────────────────┐
│ 意图是"继续/改进"？  │──No───▶ 非优化任务，正常处理
└─────────────────────┘
    │ Yes
    ▼
┌─────────────────────┐
│ 对话中有活跃 session │──Yes──▶ 恢复该 Session，继续迭代
└─────────────────────┘
    │ No
    ▼
询问用户：新建还是从历史恢复？
```

**继续优化的触发词示例**：
- "再优化一下"
- "继续改进"
- "根据反馈修改"
- "换个方向试试"
- "分数太低，重新优化"

### 恢复 Session 的流程

1. 读取 `session.json`
2. 检查 `status` 是否为 `active`
3. 从 `history` 最后一项获取 `userFeedback`
4. 执行迭代优化流程（Phase 4 开始）

---

## 错误处理

| 情况 | 处理 |
|-----|-----|
| 无有效内容 | 提示用户提供提示词内容 |
| 模板文件不存在 | 回退到英文模板，通知用户 |
| WebView 应用不存在 | 跳过交互确认，直接输出结果 |
| Session 目录创建失败 | 输出错误信息，终止流程 |
| session.json 读写失败 | 输出错误信息，终止流程 |


