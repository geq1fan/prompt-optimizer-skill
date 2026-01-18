package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// ========== 集成测试 - 使用预制 testdata ==========

// TestIntegration_LoadBasicInput 测试加载基础输入文件 (v1, 无历史)
func TestIntegration_LoadBasicInput(t *testing.T) {
	inputPath := filepath.Join("testdata", "input_v1_basic.json")
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load basic input: %v", err)
	}

	// 验证版本
	if app.inputData.Version != 1 {
		t.Errorf("expected version 1, got %d", app.inputData.Version)
	}

	// 验证原始 Prompt
	if app.inputData.OriginalPrompt != "帮我写一个登录页面" {
		t.Errorf("unexpected original prompt: %s", app.inputData.OriginalPrompt)
	}

	// 验证当前迭代
	if app.inputData.Current.IterationID != "iter-001" {
		t.Errorf("unexpected iteration ID: %s", app.inputData.Current.IterationID)
	}
	if app.inputData.Current.Score != 75 {
		t.Errorf("expected score 75, got %d", app.inputData.Current.Score)
	}

	// 验证建议方向
	if len(app.inputData.Current.SuggestedDirections) != 3 {
		t.Errorf("expected 3 directions, got %d", len(app.inputData.Current.SuggestedDirections))
	}

	// 验证无历史记录
	if len(app.inputData.History) != 0 {
		t.Errorf("expected no history, got %d items", len(app.inputData.History))
	}
}

// TestIntegration_LoadInputWithHistory 测试加载带历史记录的输入文件 (v3)
func TestIntegration_LoadInputWithHistory(t *testing.T) {
	inputPath := filepath.Join("testdata", "input_v3_with_history.json")
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load input with history: %v", err)
	}

	// 验证版本
	if app.inputData.Version != 3 {
		t.Errorf("expected version 3, got %d", app.inputData.Version)
	}

	// 验证历史记录数量
	if len(app.inputData.History) != 2 {
		t.Errorf("expected 2 history items, got %d", len(app.inputData.History))
	}

	// 验证历史记录内容
	history1 := app.inputData.History[0]
	if history1.IterationID != "iter-001" {
		t.Errorf("unexpected first history iteration ID: %s", history1.IterationID)
	}
	if history1.Score != 55 {
		t.Errorf("expected first history score 55, got %d", history1.Score)
	}
	if len(history1.UserFeedback.SelectedDirections) != 1 {
		t.Errorf("expected 1 direction in first history feedback")
	}
	if history1.UserFeedback.SelectedDirections[0] != "structure" {
		t.Errorf("unexpected direction: %s", history1.UserFeedback.SelectedDirections[0])
	}

	history2 := app.inputData.History[1]
	if history2.IterationID != "iter-002" {
		t.Errorf("unexpected second history iteration ID: %s", history2.IterationID)
	}
	if history2.Score != 72 {
		t.Errorf("expected second history score 72, got %d", history2.Score)
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

	// 验证版本
	if app.inputData.Version != 5 {
		t.Errorf("expected version 5, got %d", app.inputData.Version)
	}

	// 验证历史记录数量
	if len(app.inputData.History) != 4 {
		t.Errorf("expected 4 history items, got %d", len(app.inputData.History))
	}

	// 验证分数递增趋势
	expectedScores := []int{40, 55, 70, 82}
	for i, h := range app.inputData.History {
		if h.Score != expectedScores[i] {
			t.Errorf("history[%d] expected score %d, got %d", i, expectedScores[i], h.Score)
		}
	}

	// 验证当前分数最高
	if app.inputData.Current.Score != 92 {
		t.Errorf("expected current score 92, got %d", app.inputData.Current.Score)
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
	if app.inputData.OriginalPrompt != "" {
		t.Errorf("expected empty original prompt")
	}
	if app.inputData.Current.OptimizedPrompt != "" {
		t.Errorf("expected empty optimized prompt")
	}
	if app.inputData.Current.Score != 0 {
		t.Errorf("expected score 0, got %d", app.inputData.Current.Score)
	}
	if len(app.inputData.Current.SuggestedDirections) != 0 {
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
	if app.inputData.OriginalPrompt != "中文测试 🎉 日本語 한국어 العربية" {
		t.Errorf("unicode original prompt not preserved: %s", app.inputData.OriginalPrompt)
	}

	// 验证 Emoji 在优化后的 Prompt 中
	optimized := app.inputData.Current.OptimizedPrompt
	if len(optimized) == 0 {
		t.Error("optimized prompt is empty")
	}

	// 验证方向标签包含中文
	if len(app.inputData.Current.SuggestedDirections) > 0 {
		dir := app.inputData.Current.SuggestedDirections[0]
		if dir.Label == "" {
			t.Error("direction label is empty")
		}
	}
}

// TestIntegration_SubmitAndWriteResult 测试提交并写入结果
func TestIntegration_SubmitAndWriteResult(t *testing.T) {
	inputPath := filepath.Join("testdata", "input_v1_basic.json")
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "result.json")

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
	inputPath := filepath.Join("testdata", "input_v3_with_history.json")
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "result.json")

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
	inputPath := filepath.Join("testdata", "input_v1_basic.json")
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "result.json")

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
	inputPath := filepath.Join("testdata", "input_v3_with_history.json")
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load input: %v", err)
	}

	// 调用 GetInputData (前端 API)
	data := app.GetInputData()

	if data == nil {
		t.Fatal("GetInputData returned nil")
	}

	// 验证返回的数据与原始数据一致
	if data.Version != 3 {
		t.Errorf("expected version 3, got %d", data.Version)
	}
	if data.OriginalPrompt != "写一个 API 接口" {
		t.Errorf("unexpected original prompt: %s", data.OriginalPrompt)
	}
	if len(data.History) != 2 {
		t.Errorf("expected 2 history items, got %d", len(data.History))
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
	inputPath := filepath.Join("testdata", "input_v3_with_history.json")
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load input: %v", err)
	}

	// 验证 Markdown 标题被保留
	optimized := app.inputData.Current.OptimizedPrompt
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
	if !containsSubstring(optimized, "```json") {
		t.Error("Markdown code block not preserved")
	}
}

// TestIntegration_DirectionFieldsComplete 测试方向字段完整性
func TestIntegration_DirectionFieldsComplete(t *testing.T) {
	inputPath := filepath.Join("testdata", "input_v1_basic.json")
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to load input: %v", err)
	}

	for i, dir := range app.inputData.Current.SuggestedDirections {
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

	for i, h := range app.inputData.History {
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
