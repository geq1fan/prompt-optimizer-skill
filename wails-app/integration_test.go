package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ========== 集成测试 - 使用预制 testdata ==========

// TestIntegration_LoadBasicInput 测试加载基础输入文件 (v1, 无历史记录)
func TestIntegration_LoadBasicInput(t *testing.T) {
	// 创建临时测试数据，避免依赖可能被修改的文件
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.json")
	outputPath := filepath.Join(dir, "result.json")

	// 创建 v1 格式的测试数据
	testData := `{
		"version": 1,
		"originalPrompt": "帮我写一个登录页面",
		"current": {
			"iterationId": "iter-001",
			"optimizedPrompt": "# Role: 前端开发专家",
			"reviewReport": "评审报告",
			"evaluationReport": "评估报告",
			"score": 75,
			"suggestedDirections": [{"id": "a", "label": "A", "description": "Desc"}]
		},
		"history": []
	}`
	if err := os.WriteFile(inputPath, []byte(testData), 0644); err != nil {
		t.Fatalf("failed to write test data: %v", err)
	}

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load basic input: %v", err)
	}

	// v1 迁移后变为 v4
	if app.sessionData.Version != 4 {
		t.Errorf("expected version 4 (migrated), got %d", app.sessionData.Version)
	}

	// 验证原始 Prompt
	if app.sessionData.OriginalPrompt != "帮我写一个登录页面" {
		t.Errorf("unexpected original prompt: %s", app.sessionData.OriginalPrompt)
	}

	// 验证迁移时生成了 sessionId
	if !strings.HasPrefix(app.sessionData.SessionID, "session_") {
		t.Errorf("expected sessionId to start with 'session_', got %s", app.sessionData.SessionID)
	}

	// 验证语言检测
	if app.sessionData.Lang != "cn" {
		t.Errorf("expected lang 'cn' for Chinese prompt, got %s", app.sessionData.Lang)
	}

	// 验证 current 有数据
	if app.sessionData.Current.Score != 75 {
		t.Errorf("expected current score 75, got %d", app.sessionData.Current.Score)
	}

	// 验证无历史记录
	if len(app.sessionData.History) != 0 {
		t.Errorf("expected 0 history items, got %d", len(app.sessionData.History))
	}
}

// TestIntegration_LoadInputWithHistory 测试加载带历史记录的输入文件 (v3)
func TestIntegration_LoadInputWithHistory(t *testing.T) {
	// 创建临时测试数据，避免依赖可能被修改的文件
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.json")
	outputPath := filepath.Join(dir, "result.json")

	// 创建 v3 格式的测试数据
	testData := `{
		"version": 3,
		"originalPrompt": "写一个 API 接口",
		"current": {
			"iterationId": "iter-003",
			"optimizedPrompt": "# Role: API 架构师",
			"reviewReport": "评审报告",
			"evaluationReport": "评估报告",
			"score": 88,
			"suggestedDirections": [{"id": "pagination", "label": "分页", "description": "添加分页"}]
		},
		"history": [
			{"iterationId": "iter-001", "optimizedPrompt": "v1", "reviewReport": "", "evaluationReport": "", "score": 55, "userFeedback": {"selectedDirections": ["structure"], "userInput": "改进"}},
			{"iterationId": "iter-002", "optimizedPrompt": "v2", "reviewReport": "", "evaluationReport": "", "score": 72, "userFeedback": {"selectedDirections": ["security", "format"], "userInput": "继续"}}
		]
	}`
	if err := os.WriteFile(inputPath, []byte(testData), 0644); err != nil {
		t.Fatalf("failed to write test data: %v", err)
	}

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load input with history: %v", err)
	}

	// v3 迁移后变为 v4
	if app.sessionData.Version != 4 {
		t.Errorf("expected version 4 (migrated), got %d", app.sessionData.Version)
	}

	// 验证历史记录数量 (当前测试数据有2个历史记录)
	if len(app.sessionData.History) != 2 {
		t.Errorf("expected 2 history items, got %d", len(app.sessionData.History))
	}

	// 验证历史记录内容
	if len(app.sessionData.History) > 0 {
		history1 := app.sessionData.History[0]
		if history1.IterationID != "iter-001" {
			t.Errorf("unexpected first history iteration ID: %s", history1.IterationID)
		}
		if history1.Score != 55 {
			t.Errorf("expected first history score 55, got %d", history1.Score)
		}
		if len(history1.UserFeedback.SelectedDirections) != 1 {
			t.Errorf("expected 1 direction in first history feedback")
		}
		if len(history1.UserFeedback.SelectedDirections) > 0 && history1.UserFeedback.SelectedDirections[0] != "structure" {
			t.Errorf("unexpected direction: %s", history1.UserFeedback.SelectedDirections[0])
		}
	}

	// 验证当前版本分数
	if app.sessionData.Current.Score != 88 {
		t.Errorf("expected current score 88, got %d", app.sessionData.Current.Score)
	}
}

// TestIntegration_LoadLongHistory 测试加载长历史记录 (v5, 4个历史版本)
func TestIntegration_LoadLongHistory(t *testing.T) {
	inputPath := filepath.Join("testdata", "input_v5_long_history.json")
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load long history input: %v", err)
	}

	// v5 >= 4，直接解析
	if app.sessionData.Version != 5 {
		t.Errorf("expected version 5, got %d", app.sessionData.Version)
	}

	// 验证历史记录数量
	if len(app.sessionData.History) != 4 {
		t.Errorf("expected 4 history items, got %d", len(app.sessionData.History))
	}

	// 验证分数递增趋势
	expectedScores := []int{40, 55, 70, 82}
	for i, h := range app.sessionData.History {
		if h.Score != expectedScores[i] {
			t.Errorf("history[%d] expected score %d, got %d", i, expectedScores[i], h.Score)
		}
	}

	// 验证当前分数最高
	if app.sessionData.Current.Score != 92 {
		t.Errorf("expected current score 92, got %d", app.sessionData.Current.Score)
	}
}

// TestIntegration_LoadEmptyInput 测试加载空内容输入文件
func TestIntegration_LoadEmptyInput(t *testing.T) {
	inputPath := filepath.Join("testdata", "input_empty.json")
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load empty input: %v", err)
	}

	// 验证空值处理
	if app.sessionData.OriginalPrompt != "" {
		t.Errorf("expected empty original prompt")
	}
	if app.sessionData.Current.OptimizedPrompt != "" {
		t.Errorf("expected empty optimized prompt")
	}
	if app.sessionData.Current.Score != 0 {
		t.Errorf("expected score 0, got %d", app.sessionData.Current.Score)
	}
	if len(app.sessionData.Current.SuggestedDirections) != 0 {
		t.Errorf("expected no directions")
	}
}

// TestIntegration_LoadUnicodeInput 测试加载 Unicode 内容
func TestIntegration_LoadUnicodeInput(t *testing.T) {
	inputPath := filepath.Join("testdata", "input_unicode.json")
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load unicode input: %v", err)
	}

	// 验证中文
	if app.sessionData.OriginalPrompt != "中文测试 🎉 日本語 한국어 العربية" {
		t.Errorf("unicode original prompt not preserved: %s", app.sessionData.OriginalPrompt)
	}

	// 验证 Emoji 在优化后的 Prompt 中
	optimized := app.sessionData.Current.OptimizedPrompt
	if len(optimized) == 0 {
		t.Error("optimized prompt is empty")
	}

	// 验证方向标签包含中文
	if len(app.sessionData.Current.SuggestedDirections) > 0 {
		dir := app.sessionData.Current.SuggestedDirections[0]
		if dir.Label == "" {
			t.Error("direction label is empty")
		}
	}
}

// TestIntegration_SubmitAndWriteResult 测试提交并写入结果
func TestIntegration_SubmitAndWriteResult(t *testing.T) {
	// 复制测试数据到临时目录，避免修改原始文件
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.json")
	outputPath := filepath.Join(dir, "result.json")

	// 使用 v1 格式测试数据
	testData := `{
		"version": 1,
		"originalPrompt": "帮我写一个登录页面",
		"current": {
			"iterationId": "iter-001",
			"optimizedPrompt": "# Role: 前端开发专家",
			"reviewReport": "评审报告",
			"evaluationReport": "评估报告",
			"score": 75,
			"suggestedDirections": [
				{"id": "error-handling", "label": "错误处理", "description": "添加错误处理"}
			]
		},
		"history": []
	}`
	if err := os.WriteFile(inputPath, []byte(testData), 0644); err != nil {
		t.Fatalf("failed to write test data: %v", err)
	}

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load input: %v", err)
	}

	// 模拟提交
	result := Result{
		Action:             "submit",
		SelectedDirections: []string{"error-handling", "ui-style"},
		UserInput:          "请添加错误处理和现代化 UI 风格",
	}

	err = app.writeResult(result)
	if err != nil {
		t.Fatalf("failed to write result: %v", err)
	}

	// 读取并验证结果
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read result file: %v", err)
	}

	var written Result
	if err := json.Unmarshal(data, &written); err != nil {
		t.Fatalf("failed to parse result: %v", err)
	}

	if written.Action != "submit" {
		t.Errorf("expected action 'submit', got '%s'", written.Action)
	}
	if len(written.SelectedDirections) != 2 {
		t.Errorf("expected 2 directions, got %d", len(written.SelectedDirections))
	}
	if written.SelectedDirections[0] != "error-handling" {
		t.Errorf("unexpected first direction: %s", written.SelectedDirections[0])
	}
	if written.UserInput != "请添加错误处理和现代化 UI 风格" {
		t.Errorf("unexpected user input: %s", written.UserInput)
	}
}

// TestIntegration_RollbackAndWriteResult 测试回滚并写入结果
func TestIntegration_RollbackAndWriteResult(t *testing.T) {
	// 复制测试数据到临时目录，避免修改原始文件
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.json")
	outputPath := filepath.Join(dir, "result.json")

	// 使用 v3 格式测试数据（含历史记录）
	testData := `{
		"version": 3,
		"originalPrompt": "写一个 API 接口",
		"current": {
			"iterationId": "iter-002",
			"optimizedPrompt": "# Role: API 架构师",
			"reviewReport": "评审报告",
			"evaluationReport": "评估报告",
			"score": 72
		},
		"history": [
			{
				"iterationId": "iter-001",
				"optimizedPrompt": "# Role: API 开发者",
				"reviewReport": "初版评审",
				"evaluationReport": "初版评估",
				"score": 55,
				"userFeedback": {
					"selectedDirections": ["structure"],
					"userInput": "需要更详细的结构"
				}
			}
		]
	}`
	if err := os.WriteFile(inputPath, []byte(testData), 0644); err != nil {
		t.Fatalf("failed to write test data: %v", err)
	}

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load input: %v", err)
	}

	// 模拟回滚到第一个版本
	result := Result{
		Action:              "rollback",
		SelectedDirections:  []string{"examples"},
		UserInput:           "基于第一版重新优化",
		RollbackToIteration: "iter-001",
	}

	err = app.writeResult(result)
	if err != nil {
		t.Fatalf("failed to write result: %v", err)
	}

	// 读取并验证结果
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read result file: %v", err)
	}

	var written Result
	if err := json.Unmarshal(data, &written); err != nil {
		t.Fatalf("failed to parse result: %v", err)
	}

	if written.Action != "rollback" {
		t.Errorf("expected action 'rollback', got '%s'", written.Action)
	}
	if written.RollbackToIteration != "iter-001" {
		t.Errorf("expected rollback to 'iter-001', got '%s'", written.RollbackToIteration)
	}
}

// TestIntegration_CancelAndWriteResult 测试取消并写入结果
func TestIntegration_CancelAndWriteResult(t *testing.T) {
	// 复制测试数据到临时目录，避免修改原始文件
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.json")
	outputPath := filepath.Join(dir, "result.json")

	// 使用 v1 格式测试数据
	testData := `{
		"version": 1,
		"originalPrompt": "帮我写一个登录页面",
		"current": {
			"iterationId": "iter-001",
			"optimizedPrompt": "# Role: 前端开发专家",
			"reviewReport": "评审报告",
			"evaluationReport": "评估报告",
			"score": 75
		},
		"history": []
	}`
	if err := os.WriteFile(inputPath, []byte(testData), 0644); err != nil {
		t.Fatalf("failed to write test data: %v", err)
	}

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load input: %v", err)
	}

	// 模拟取消
	result := Result{
		Action:             "cancel",
		SelectedDirections: []string{},
		UserInput:          "",
	}

	err = app.writeResult(result)
	if err != nil {
		t.Fatalf("failed to write result: %v", err)
	}

	// 读取并验证结果
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read result file: %v", err)
	}

	var written Result
	if err := json.Unmarshal(data, &written); err != nil {
		t.Fatalf("failed to parse result: %v", err)
	}

	if written.Action != "cancel" {
		t.Errorf("expected action 'cancel', got '%s'", written.Action)
	}
	if len(written.SelectedDirections) != 0 {
		t.Errorf("expected empty directions for cancel")
	}
}

// TestIntegration_GetInputDataAPI 测试 GetInputData API
func TestIntegration_GetInputDataAPI(t *testing.T) {
	// 创建临时测试数据，避免依赖可能被修改的文件
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.json")
	outputPath := filepath.Join(dir, "result.json")

	// 创建 v3 格式的测试数据
	testData := `{
		"version": 3,
		"originalPrompt": "写一个 API 接口",
		"current": {
			"iterationId": "iter-003",
			"optimizedPrompt": "# Role: API 架构师",
			"reviewReport": "评审报告",
			"evaluationReport": "评估报告",
			"score": 88,
			"suggestedDirections": [{"id": "pagination", "label": "分页", "description": "添加分页"}]
		},
		"history": [
			{"iterationId": "iter-001", "optimizedPrompt": "v1", "score": 55, "userFeedback": {"selectedDirections": ["a"], "userInput": "改进"}},
			{"iterationId": "iter-002", "optimizedPrompt": "v2", "score": 72, "userFeedback": {"selectedDirections": ["b"], "userInput": "继续"}}
		]
	}`
	if err := os.WriteFile(inputPath, []byte(testData), 0644); err != nil {
		t.Fatalf("failed to write test data: %v", err)
	}

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load input: %v", err)
	}

	// 调用 GetInputData (前端 API)
	data := app.GetInputData()

	if data == nil {
		t.Fatal("GetInputData returned nil")
	}

	// 验证返回的数据与迁移后的版本一致
	if data.Version != 4 {
		t.Errorf("expected version 4 (migrated), got %d", data.Version)
	}
	if data.OriginalPrompt != "写一个 API 接口" {
		t.Errorf("unexpected original prompt: %s", data.OriginalPrompt)
	}
	// 当前测试数据有2个历史记录
	if len(data.History) != 2 {
		t.Errorf("expected 2 history items, got %d", len(data.History))
	}
	// 验证当前版本分数
	if data.Current.Score != 88 {
		t.Errorf("expected current score 88, got %d", data.Current.Score)
	}
}

// TestIntegration_GetSessionDataAPI 测试 GetSessionData API (v4)
func TestIntegration_GetSessionDataAPI(t *testing.T) {
	// 创建临时测试数据
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.json")
	outputPath := filepath.Join(dir, "result.json")

	// 创建 v3 格式的测试数据（包含中文，用于测试语言检测）
	testData := `{
		"version": 3,
		"originalPrompt": "写一个 API 接口",
		"current": {
			"iterationId": "iter-001",
			"optimizedPrompt": "优化后",
			"score": 80
		},
		"history": []
	}`
	if err := os.WriteFile(inputPath, []byte(testData), 0644); err != nil {
		t.Fatalf("failed to write test data: %v", err)
	}

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load input: %v", err)
	}

	// 调用 GetSessionData (v4 API)
	data := app.GetSessionData()

	if data == nil {
		t.Fatal("GetSessionData returned nil")
	}

	// 验证 v4 特有字段
	if !strings.HasPrefix(data.SessionID, "session_") {
		t.Errorf("expected sessionId to start with 'session_', got %s", data.SessionID)
	}
	if data.Status != "active" {
		t.Errorf("expected status 'active', got %s", data.Status)
	}
	if data.Lang != "cn" {
		t.Errorf("expected lang 'cn', got %s", data.Lang)
	}
}

// TestIntegration_GetRemainingSecondsAPI 测试 GetRemainingSeconds API
func TestIntegration_GetRemainingSecondsAPI(t *testing.T) {
	inputPath := filepath.Join("testdata", "input_v1_basic.json")
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 300)
	if err != nil {
		t.Fatalf("failed to load input: %v", err)
	}

	// 调用 GetRemainingSeconds (前端 API)
	remaining := app.GetRemainingSeconds()

	// 验证剩余时间在合理范围内
	if remaining < 298 || remaining > 300 {
		t.Errorf("expected remaining ~300, got %d", remaining)
	}
}

// TestIntegration_MarkdownContentPreserved 测试 Markdown 内容保持完整
func TestIntegration_MarkdownContentPreserved(t *testing.T) {
	// 使用 v5 测试数据，它有完整的 current 数据
	inputPath := filepath.Join("testdata", "input_v5_long_history.json")
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load input: %v", err)
	}

	// 验证 Markdown 标题被保留
	optimized := app.sessionData.Current.OptimizedPrompt
	if len(optimized) == 0 {
		t.Fatal("optimized prompt is empty")
	}

	// 检查是否包含 Markdown 格式
	if !containsSubstring(optimized, "# Role:") {
		t.Error("Markdown h1 header not preserved")
	}
	if !containsSubstring(optimized, "## Goals") {
		t.Error("Markdown h2 header not preserved")
	}
	// v5 测试数据中包含 markdown 代码块
	if !containsSubstring(optimized, "```markdown") && !containsSubstring(optimized, "```python") {
		t.Error("Markdown code block not preserved")
	}
}

// TestIntegration_DirectionFieldsComplete 测试方向字段完整性
func TestIntegration_DirectionFieldsComplete(t *testing.T) {
	// 使用 v5 测试数据，它有完整的 current 数据（包含 directions）
	inputPath := filepath.Join("testdata", "input_v5_long_history.json")
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load input: %v", err)
	}

	for i, dir := range app.sessionData.Current.SuggestedDirections {
		if dir.ID == "" {
			t.Errorf("direction[%d] ID is empty", i)
		}
		if dir.Label == "" {
			t.Errorf("direction[%d] Label is empty", i)
		}
		if dir.Description == "" {
			t.Errorf("direction[%d] Description is empty", i)
		}
	}
}

// TestIntegration_HistoryFeedbackComplete 测试历史反馈完整性
func TestIntegration_HistoryFeedbackComplete(t *testing.T) {
	inputPath := filepath.Join("testdata", "input_v5_long_history.json")
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load input: %v", err)
	}

	for i, h := range app.sessionData.History {
		if h.IterationID == "" {
			t.Errorf("history[%d] IterationID is empty", i)
		}
		if len(h.UserFeedback.SelectedDirections) == 0 {
			t.Errorf("history[%d] has no selected directions", i)
		}
		if h.UserFeedback.UserInput == "" {
			t.Errorf("history[%d] user input is empty", i)
		}
	}
}

// TestIntegration_LoadSessionV4 测试加载 v4 格式的 session.json
func TestIntegration_LoadSessionV4(t *testing.T) {
	inputPath := filepath.Join("testdata", "session_v4.json")
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load session v4: %v", err)
	}

	// 验证版本为 4（不需要迁移）
	if app.sessionData.Version != 4 {
		t.Errorf("expected version 4, got %d", app.sessionData.Version)
	}

	// 验证 v4 特有字段
	if app.sessionData.SessionID != "session_1705632000000" {
		t.Errorf("unexpected sessionId: %s", app.sessionData.SessionID)
	}
	if app.sessionData.Lang != "cn" {
		t.Errorf("unexpected lang: %s", app.sessionData.Lang)
	}
	if app.sessionData.Mode != "professional" {
		t.Errorf("unexpected mode: %s", app.sessionData.Mode)
	}
	if app.sessionData.Status != "active" {
		t.Errorf("unexpected status: %s", app.sessionData.Status)
	}
	if app.sessionData.LastAction != "submit" {
		t.Errorf("unexpected lastAction: %s", app.sessionData.LastAction)
	}

	// 验证当前迭代
	if app.sessionData.Current.Score != 85 {
		t.Errorf("expected current score 85, got %d", app.sessionData.Current.Score)
	}

	// 验证历史记录
	if len(app.sessionData.History) != 1 {
		t.Errorf("expected 1 history item, got %d", len(app.sessionData.History))
	}
	if app.sessionData.History[0].Score != 65 {
		t.Errorf("expected history score 65, got %d", app.sessionData.History[0].Score)
	}
}

// TestIntegration_DetectLang 测试语言检测功能
func TestIntegration_DetectLang(t *testing.T) {
	tests := []struct {
		text     string
		expected string
	}{
		{"Hello world", "en"},
		{"你好世界", "cn"},
		{"Hello 世界", "cn"},
		{"", "en"},
		{"1234567890", "en"},
		{"日本語テスト", "cn"}, // 日语汉字在 CJK 范围内，返回 cn
	}

	for _, tt := range tests {
		result := detectLang(tt.text)
		if result != tt.expected {
			t.Errorf("detectLang(%q) = %q, want %q", tt.text, result, tt.expected)
		}
	}
}

// ========== 辅助函数 ==========

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstringHelper(s, substr))
}

func containsSubstringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
