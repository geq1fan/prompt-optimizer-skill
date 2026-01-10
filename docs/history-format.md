# 历史文件格式

所有优化结果自动保存到工作目录的 `.prompt_history.jsonl` 文件。

## 文件格式

使用 [JSON Lines](https://jsonlines.org/) 格式，每行是一个独立的 JSON 对象：

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-15T10:30:00.000Z",
  "type": "user",
  "mode": "professional",
  "original": "写一篇关于AI的博客",
  "optimized": "请撰写一篇面向技术爱好者的博客文章，主题为人工智能在日常生活中的应用...",
  "evaluation": "综合评分: 85/100。任务表达清晰，信息完整度良好，格式规范。"
}
```

## 字段说明

| 字段 | 类型 | 说明 |
|-----|------|-----|
| `id` | string | UUID v4 格式的唯一标识符 |
| `timestamp` | string | ISO 8601 格式的时间戳 |
| `type` | string | 优化类型：`user`、`system` 或 `iterate` |
| `mode` | string | 优化模式：`basic`、`professional`、`planning`、`general`、`analytical` |
| `original` | string | 原始提示词内容 |
| `optimized` | string | 优化后的提示词内容 |
| `evaluation` | string | 评估摘要 |

## 使用场景

### 自动上下文加载

当用户调用 `/optimize-prompt iterate` 时，Agent 自动执行：

1. 读取 `.prompt_history.jsonl` 文件
2. 获取最后一行（最近一次优化记录）
3. 解析 JSON 并提取 `optimized` 字段
4. 将其作为 `lastOptimizedPrompt` 用于迭代优化

### 手动查看历史

```bash
# 查看所有历史记录
cat .prompt_history.jsonl

# 查看最近一条记录
tail -1 .prompt_history.jsonl

# 格式化显示最近一条
tail -1 .prompt_history.jsonl | jq .

# 统计优化次数
wc -l .prompt_history.jsonl
```

### 版本控制

`.prompt_history.jsonl` 是纯文本文件，可以：

- 提交到 Git 追踪提示词演进
- 使用 `git diff` 对比不同版本的优化结果
- 通过 `git log` 查看提示词优化的时间线

## 注意事项

1. **文件位置**: 始终在当前工作目录创建/读取
2. **追加写入**: 新记录追加到文件末尾，不覆盖历史
3. **编码格式**: UTF-8 编码，支持多语言内容
4. **空文件处理**: 如果文件不存在或为空，iterate 操作会提示用户提供原始提示词
