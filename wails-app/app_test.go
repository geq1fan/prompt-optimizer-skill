package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// atomicBool 兼容旧版 Go 的原子布尔类型
type atomicBool struct {
	value int32
}

func (b *atomicBool) Store(val bool) {
	if val {
		atomic.StoreInt32(&b.value, 1)
	} else {
		atomic.StoreInt32(&b.value, 0)
	}
}

func (b *atomicBool) Load() bool {
	return atomic.LoadInt32(&b.value) != 0
}

// ========== 测试辅助函数 ==========

// createTestApp 创建用于测试的 App 实例 (带 mock quitFunc)
func createTestApp(t *testing.T, quitFunc QuitFunc) (*App, string) {
	t.Helper()
	dir := t.TempDir()
	inputData := createValidInputData()
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}

	// 注入 mock quitFunc
	app.quitFunc = quitFunc
	app.ctx = context.Background()

	return app, outputPath
}

// createTestInputFile 创建测试用的输入 JSON 文件
func createTestInputFile(t *testing.T, dir string, data InputData) string {
	t.Helper()
	inputPath := filepath.Join(dir, "input.json")
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal test input: %v", err)
	}
	if err := os.WriteFile(inputPath, jsonData, 0644); err != nil {
		t.Fatalf("failed to write test input file: %v", err)
	}
	return inputPath
}

// createValidInputData 创建有效的测试输入数据
func createValidInputData() InputData {
	return InputData{
		Version:        1,
		OriginalPrompt: "Test original prompt",
		Current: IterationData{
			IterationID:      "iter-001",
			OptimizedPrompt:  "Test optimized prompt",
			ReviewReport:     "Test review report",
			EvaluationReport: "Test evaluation report",
			Score:            85,
			SuggestedDirections: []Direction{
				{ID: "examples", Label: "添加示例", Description: "补充使用案例"},
				{ID: "constraints", Label: "增强约束", Description: "明确边界条件"},
			},
		},
		History: []HistoryItem{},
	}
}

// ========== NewApp 测试 ==========

func TestNewApp_ValidInput(t *testing.T) {
	dir := t.TempDir()
	inputData := createValidInputData()
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)

	if err != nil {
		t.Fatalf("NewApp failed with valid input: %v", err)
	}
	if app == nil {
		t.Fatal("NewApp returned nil app")
	}
	if app.sessionData == nil {
		t.Fatal("NewApp did not load input data")
	}
	if app.sessionData.Version != 4 {
		t.Errorf("expected version 4 (migrated from v1), got %d", app.sessionData.Version)
	}
	if app.sessionData.OriginalPrompt != "Test original prompt" {
		t.Errorf("unexpected original prompt: %s", app.sessionData.OriginalPrompt)
	}
	if app.timeout != 600 {
		t.Errorf("expected timeout 600, got %d", app.timeout)
	}
}

func TestNewApp_FileNotExists(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "nonexistent.json")
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)

	if err == nil {
		t.Fatal("NewApp should fail for nonexistent file")
	}
	if app != nil {
		t.Fatal("NewApp should return nil app on error")
	}
}

func TestNewApp_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "invalid.json")
	outputPath := filepath.Join(dir, "result.json")

	// 写入无效 JSON
	if err := os.WriteFile(inputPath, []byte("{ invalid json }"), 0644); err != nil {
		t.Fatalf("failed to write invalid json: %v", err)
	}

	app, err := NewApp(inputPath, outputPath, 600)

	if err == nil {
		t.Fatal("NewApp should fail for invalid JSON")
	}
	if app != nil {
		t.Fatal("NewApp should return nil app on error")
	}
}

func TestNewApp_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "empty.json")
	outputPath := filepath.Join(dir, "result.json")

	// 写入空文件
	if err := os.WriteFile(inputPath, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write empty file: %v", err)
	}

	app, err := NewApp(inputPath, outputPath, 600)

	if err == nil {
		t.Fatal("NewApp should fail for empty file")
	}
	if app != nil {
		t.Fatal("NewApp should return nil app on error")
	}
}

func TestNewApp_WithHistory(t *testing.T) {
	dir := t.TempDir()
	inputData := createValidInputData()
	inputData.Version = 3
	inputData.History = []HistoryItem{
		{
			IterationID:      "iter-001",
			OptimizedPrompt:  "First version",
			ReviewReport:     "First review",
			EvaluationReport: "First evaluation",
			Score:            65,
			UserFeedback: UserFeedback{
				SelectedDirections: []string{"structure"},
				UserInput:          "希望结构更清晰",
			},
		},
		{
			IterationID:      "iter-002",
			OptimizedPrompt:  "Second version",
			ReviewReport:     "Second review",
			EvaluationReport: "Second evaluation",
			Score:            75,
			UserFeedback: UserFeedback{
				SelectedDirections: []string{"examples", "constraints"},
				UserInput:          "需要更多示例",
			},
		},
	}
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)

	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}
	if len(app.sessionData.History) != 2 {
		t.Errorf("expected 2 history items, got %d", len(app.sessionData.History))
	}
	if app.sessionData.History[0].IterationID != "iter-001" {
		t.Errorf("unexpected first history iteration ID: %s", app.sessionData.History[0].IterationID)
	}
}

// ========== GetInputData 测试 ==========

func TestGetInputData_ReturnsCorrectData(t *testing.T) {
	dir := t.TempDir()
	inputData := createValidInputData()
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}

	result := app.GetInputData()

	if result == nil {
		t.Fatal("GetInputData returned nil")
	}
	if result.Version != 4 {
		t.Errorf("expected version 4 (migrated), got %d", result.Version)
	}
	if result.Current.Score != inputData.Current.Score {
		t.Errorf("expected score %d, got %d", inputData.Current.Score, result.Current.Score)
	}
}

// ========== GetRemainingSeconds 测试 ==========

func TestGetRemainingSeconds_WithinTimeout(t *testing.T) {
	dir := t.TempDir()
	inputData := createValidInputData()
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}

	remaining := app.GetRemainingSeconds()

	// 刚创建的 app，剩余时间应该接近 600 秒
	if remaining < 598 || remaining > 600 {
		t.Errorf("expected remaining ~600, got %d", remaining)
	}
}

func TestGetRemainingSeconds_AfterSomeTime(t *testing.T) {
	dir := t.TempDir()
	inputData := createValidInputData()
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 10)
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}

	// 模拟时间流逝
	app.startTime = time.Now().Add(-5 * time.Second)

	remaining := app.GetRemainingSeconds()

	if remaining < 4 || remaining > 6 {
		t.Errorf("expected remaining ~5, got %d", remaining)
	}
}

func TestGetRemainingSeconds_AfterTimeout(t *testing.T) {
	dir := t.TempDir()
	inputData := createValidInputData()
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 10)
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}

	// 模拟超时后
	app.startTime = time.Now().Add(-15 * time.Second)

	remaining := app.GetRemainingSeconds()

	if remaining != 0 {
		t.Errorf("expected remaining 0 after timeout, got %d", remaining)
	}
}

func TestGetRemainingSeconds_ExactlyAtTimeout(t *testing.T) {
	dir := t.TempDir()
	inputData := createValidInputData()
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 10)
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}

	// 模拟恰好超时
	app.startTime = time.Now().Add(-10 * time.Second)

	remaining := app.GetRemainingSeconds()

	if remaining != 0 {
		t.Errorf("expected remaining 0 at exact timeout, got %d", remaining)
	}
}

// ========== writeResult 测试 ==========

func TestWriteResult_SuccessfulWrite(t *testing.T) {
	dir := t.TempDir()
	inputData := createValidInputData()
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}

	result := Result{
		Action:             "submit",
		SelectedDirections: []string{"examples", "constraints"},
		UserInput:          "Test input",
	}

	err = app.writeResult(result)

	if err != nil {
		t.Fatalf("writeResult failed: %v", err)
	}

	// 验证文件内容
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	var written Result
	if err := json.Unmarshal(data, &written); err != nil {
		t.Fatalf("failed to parse output file: %v", err)
	}

	if written.Action != "submit" {
		t.Errorf("expected action 'submit', got '%s'", written.Action)
	}
	if len(written.SelectedDirections) != 2 {
		t.Errorf("expected 2 directions, got %d", len(written.SelectedDirections))
	}
	if written.UserInput != "Test input" {
		t.Errorf("expected user input 'Test input', got '%s'", written.UserInput)
	}
}

func TestWriteResult_OnlyOnce(t *testing.T) {
	dir := t.TempDir()
	inputData := createValidInputData()
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}

	// 第一次写入
	result1 := Result{
		Action:             "submit",
		SelectedDirections: []string{"examples"},
		UserInput:          "First write",
	}
	err = app.writeResult(result1)
	if err != nil {
		t.Fatalf("first writeResult failed: %v", err)
	}

	// 第二次写入 (应该被忽略)
	result2 := Result{
		Action:             "cancel",
		SelectedDirections: []string{},
		UserInput:          "Second write",
	}
	err = app.writeResult(result2)
	// 第二次调用不应报错，只是被忽略
	if err != nil {
		t.Fatalf("second writeResult should not fail: %v", err)
	}

	// 验证文件内容仍是第一次的
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	var written Result
	if err := json.Unmarshal(data, &written); err != nil {
		t.Fatalf("failed to parse output file: %v", err)
	}

	if written.Action != "submit" {
		t.Errorf("expected action 'submit' (first write), got '%s'", written.Action)
	}
	if written.UserInput != "First write" {
		t.Errorf("expected 'First write', got '%s'", written.UserInput)
	}
}

func TestWriteResult_ConcurrentWrites(t *testing.T) {
	dir := t.TempDir()
	inputData := createValidInputData()
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}

	// 并发写入测试
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			result := Result{
				Action:             "submit",
				SelectedDirections: []string{},
				UserInput:          "Concurrent write",
			}
			app.writeResult(result)
		}(i)
	}
	wg.Wait()

	// 验证文件存在且内容有效
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	var written Result
	if err := json.Unmarshal(data, &written); err != nil {
		t.Fatalf("failed to parse output file: %v", err)
	}

	if written.Action != "submit" {
		t.Errorf("expected action 'submit', got '%s'", written.Action)
	}
}

func TestWriteResult_RollbackAction(t *testing.T) {
	dir := t.TempDir()
	inputData := createValidInputData()
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}

	result := Result{
		Action:              "rollback",
		SelectedDirections:  []string{"structure"},
		UserInput:           "Rollback test",
		RollbackToIteration: "iter-001",
	}

	err = app.writeResult(result)
	if err != nil {
		t.Fatalf("writeResult failed: %v", err)
	}

	// 验证文件内容
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	var written Result
	if err := json.Unmarshal(data, &written); err != nil {
		t.Fatalf("failed to parse output file: %v", err)
	}

	if written.Action != "rollback" {
		t.Errorf("expected action 'rollback', got '%s'", written.Action)
	}
	if written.RollbackToIteration != "iter-001" {
		t.Errorf("expected rollback to 'iter-001', got '%s'", written.RollbackToIteration)
	}
}

// ========== Result 结构测试 ==========

func TestResult_JSONSerialization(t *testing.T) {
	tests := []struct {
		name   string
		result Result
	}{
		{
			name: "submit action",
			result: Result{
				Action:             "submit",
				SelectedDirections: []string{"examples", "constraints"},
				UserInput:          "Test input",
			},
		},
		{
			name: "cancel action",
			result: Result{
				Action:             "cancel",
				SelectedDirections: []string{},
				UserInput:          "",
			},
		},
		{
			name: "timeout action",
			result: Result{
				Action:             "timeout",
				SelectedDirections: []string{},
				UserInput:          "",
			},
		},
		{
			name: "rollback action",
			result: Result{
				Action:              "rollback",
				SelectedDirections:  []string{"structure"},
				UserInput:           "Rollback",
				RollbackToIteration: "iter-002",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 序列化
			data, err := json.Marshal(tt.result)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}

			// 反序列化
			var decoded Result
			if err := json.Unmarshal(data, &decoded); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}

			// 验证
			if decoded.Action != tt.result.Action {
				t.Errorf("action mismatch: expected %s, got %s", tt.result.Action, decoded.Action)
			}
			if len(decoded.SelectedDirections) != len(tt.result.SelectedDirections) {
				t.Errorf("directions count mismatch")
			}
			if decoded.UserInput != tt.result.UserInput {
				t.Errorf("user input mismatch")
			}
			if decoded.RollbackToIteration != tt.result.RollbackToIteration {
				t.Errorf("rollback iteration mismatch")
			}
		})
	}
}

// ========== InputData 结构测试 ==========

func TestInputData_JSONSerialization(t *testing.T) {
	original := createValidInputData()
	original.History = []HistoryItem{
		{
			IterationID:      "iter-001",
			OptimizedPrompt:  "History prompt",
			ReviewReport:     "History review",
			EvaluationReport: "History evaluation",
			Score:            70,
			UserFeedback: UserFeedback{
				SelectedDirections: []string{"examples"},
				UserInput:          "Feedback",
			},
		},
	}

	// 序列化
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// 反序列化
	var decoded InputData
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// 验证
	if decoded.Version != original.Version {
		t.Errorf("version mismatch")
	}
	if decoded.OriginalPrompt != original.OriginalPrompt {
		t.Errorf("original prompt mismatch")
	}
	if decoded.Current.IterationID != original.Current.IterationID {
		t.Errorf("current iteration ID mismatch")
	}
	if len(decoded.History) != 1 {
		t.Errorf("history count mismatch")
	}
	if decoded.History[0].UserFeedback.UserInput != "Feedback" {
		t.Errorf("history feedback mismatch")
	}
}

// ========== Direction 测试 ==========

func TestDirection_JSONSerialization(t *testing.T) {
	dir := Direction{
		ID:          "examples",
		Label:       "添加示例",
		Description: "补充具体使用案例",
	}

	data, err := json.Marshal(dir)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded Direction
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.ID != dir.ID {
		t.Errorf("ID mismatch: expected %s, got %s", dir.ID, decoded.ID)
	}
	if decoded.Label != dir.Label {
		t.Errorf("Label mismatch")
	}
	if decoded.Description != dir.Description {
		t.Errorf("Description mismatch")
	}
}

// ========== 边界条件测试 ==========

func TestNewApp_ZeroTimeout(t *testing.T) {
	dir := t.TempDir()
	inputData := createValidInputData()
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 0)

	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}
	if app.timeout != 0 {
		t.Errorf("expected timeout 0, got %d", app.timeout)
	}

	// 零超时时，剩余时间应该为 0
	remaining := app.GetRemainingSeconds()
	if remaining != 0 {
		t.Errorf("expected remaining 0 with zero timeout, got %d", remaining)
	}
}

func TestNewApp_EmptyDirections(t *testing.T) {
	dir := t.TempDir()
	inputData := InputData{
		Version:        1,
		OriginalPrompt: "Test",
		Current: IterationData{
			IterationID:         "iter-001",
			OptimizedPrompt:     "Test",
			Score:               50,
			SuggestedDirections: []Direction{}, // 空数组
		},
		History: nil, // nil history
	}
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)

	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}
	if len(app.sessionData.Current.SuggestedDirections) != 0 {
		t.Errorf("expected empty directions")
	}
	if app.sessionData.History != nil {
		t.Errorf("expected nil history")
	}
}

func TestNewApp_UnicodeContent(t *testing.T) {
	dir := t.TempDir()
	inputData := InputData{
		Version:        1,
		OriginalPrompt: "中文测试 🎉 日本語 한국어",
		Current: IterationData{
			IterationID:      "iter-001",
			OptimizedPrompt:  "优化后的提示词 with émojis 🚀",
			ReviewReport:     "评审报告 contains 特殊字符 <>&\"'",
			EvaluationReport: "评估报告",
			Score:            88,
			SuggestedDirections: []Direction{
				{ID: "中文ID", Label: "中文标签", Description: "中文描述"},
			},
		},
	}
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	app, err := NewApp(inputPath, outputPath, 600)

	if err != nil {
		t.Fatalf("NewApp failed with unicode content: %v", err)
	}
	if app.sessionData.OriginalPrompt != inputData.OriginalPrompt {
		t.Errorf("unicode original prompt not preserved")
	}
	if app.sessionData.Current.SuggestedDirections[0].Label != "中文标签" {
		t.Errorf("unicode direction label not preserved")
	}
}

// ========== Submit/Rollback/Cancel API 测试 ==========

func TestSubmit_WritesResultAndCallsQuit(t *testing.T) {
	var quitCalled atomicBool
	mockQuit := func(ctx context.Context) {
		quitCalled.Store(true)
	}

	app, outputPath := createTestApp(t, mockQuit)

	err := app.Submit([]string{"examples", "constraints"}, "Test submit")
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	// 验证结果文件
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	var result Result
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to parse result: %v", err)
	}

	if result.Action != "submit" {
		t.Errorf("expected action 'submit', got '%s'", result.Action)
	}
	if len(result.SelectedDirections) != 2 {
		t.Errorf("expected 2 directions, got %d", len(result.SelectedDirections))
	}
	if result.UserInput != "Test submit" {
		t.Errorf("expected user input 'Test submit', got '%s'", result.UserInput)
	}

	// 等待 goroutine 调用 quit (0.5s + 余量)
	time.Sleep(600 * time.Millisecond)
	if !quitCalled.Load() {
		t.Error("expected quitFunc to be called")
	}
}

func TestSubmit_NilDirections(t *testing.T) {
	mockQuit := func(ctx context.Context) {}
	app, outputPath := createTestApp(t, mockQuit)

	// 传入 nil directions
	err := app.Submit(nil, "Test with nil")
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	var result Result
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to parse result: %v", err)
	}

	// nil 应该被转为空数组
	if result.SelectedDirections == nil {
		t.Error("expected empty array, got nil")
	}
	if len(result.SelectedDirections) != 0 {
		t.Errorf("expected empty directions, got %d", len(result.SelectedDirections))
	}
}

func TestSubmit_EmptyDirections(t *testing.T) {
	mockQuit := func(ctx context.Context) {}
	app, outputPath := createTestApp(t, mockQuit)

	err := app.Submit([]string{}, "Empty directions")
	if err != nil {
		t.Fatalf("Submit failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	var result Result
	json.Unmarshal(data, &result)

	if len(result.SelectedDirections) != 0 {
		t.Errorf("expected empty directions")
	}
}

func TestRollback_WritesResultWithIterationId(t *testing.T) {
	var quitCalled atomicBool
	mockQuit := func(ctx context.Context) {
		quitCalled.Store(true)
	}

	app, outputPath := createTestApp(t, mockQuit)

	err := app.Rollback("iter-002", []string{"structure"}, "Rollback test")
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	var result Result
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to parse result: %v", err)
	}

	if result.Action != "rollback" {
		t.Errorf("expected action 'rollback', got '%s'", result.Action)
	}
	if result.RollbackToIteration != "iter-002" {
		t.Errorf("expected rollback to 'iter-002', got '%s'", result.RollbackToIteration)
	}
	if len(result.SelectedDirections) != 1 {
		t.Errorf("expected 1 direction, got %d", len(result.SelectedDirections))
	}

	// 等待 quit 调用
	time.Sleep(600 * time.Millisecond)
	if !quitCalled.Load() {
		t.Error("expected quitFunc to be called")
	}
}

func TestRollback_NilDirections(t *testing.T) {
	mockQuit := func(ctx context.Context) {}
	app, outputPath := createTestApp(t, mockQuit)

	err := app.Rollback("iter-001", nil, "Rollback with nil")
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	data, _ := os.ReadFile(outputPath)
	var result Result
	json.Unmarshal(data, &result)

	if result.SelectedDirections == nil {
		t.Error("expected empty array, got nil")
	}
}

func TestCancel_WritesResultAndCallsQuit(t *testing.T) {
	var quitCalled atomicBool
	mockQuit := func(ctx context.Context) {
		quitCalled.Store(true)
	}

	app, outputPath := createTestApp(t, mockQuit)

	err := app.Cancel()
	if err != nil {
		t.Fatalf("Cancel failed: %v", err)
	}

	// Cancel 是同步调用 quit，不需要等待
	if !quitCalled.Load() {
		t.Error("expected quitFunc to be called immediately")
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	var result Result
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to parse result: %v", err)
	}

	if result.Action != "cancel" {
		t.Errorf("expected action 'cancel', got '%s'", result.Action)
	}
	if len(result.SelectedDirections) != 0 {
		t.Errorf("expected empty directions")
	}
	if result.UserInput != "" {
		t.Errorf("expected empty user input")
	}
}

func TestCancel_NilQuitFunc(t *testing.T) {
	// 测试 quitFunc 为 nil 时不 panic
	app, outputPath := createTestApp(t, nil)

	err := app.Cancel()
	if err != nil {
		t.Fatalf("Cancel failed: %v", err)
	}

	// 验证结果仍然写入
	data, _ := os.ReadFile(outputPath)
	var result Result
	json.Unmarshal(data, &result)

	if result.Action != "cancel" {
		t.Errorf("expected action 'cancel', got '%s'", result.Action)
	}
}

// ========== startup 和 beforeClose 测试 ==========

func TestStartup_SetsContextAndStartTime(t *testing.T) {
	app, _ := createTestApp(t, nil)

	ctx := context.Background()
	oldStartTime := app.startTime

	// 等待一小段时间
	time.Sleep(10 * time.Millisecond)

	app.startup(ctx)

	if app.ctx != ctx {
		t.Error("startup should set context")
	}
	if !app.startTime.After(oldStartTime) {
		t.Error("startup should update startTime")
	}
}

func TestBeforeClose_WritesCancel(t *testing.T) {
	app, outputPath := createTestApp(t, nil)
	ctx := context.Background()

	result := app.beforeClose(ctx)

	// 应该返回 false 允许关闭
	if result != false {
		t.Error("beforeClose should return false")
	}

	// 验证写入了 cancel 结果
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	var written Result
	json.Unmarshal(data, &written)

	if written.Action != "cancel" {
		t.Errorf("expected action 'cancel', got '%s'", written.Action)
	}
}

// ========== watchTimeout 测试 ==========

func TestWatchTimeout_TimesOut(t *testing.T) {
	dir := t.TempDir()
	inputData := createValidInputData()
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	// 创建一个 1 秒超时的 app
	app, err := NewApp(inputPath, outputPath, 1)
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}

	var quitCalled atomicBool
	app.quitFunc = func(ctx context.Context) {
		quitCalled.Store(true)
	}
	app.ctx = context.Background()

	// 手动启动 watchTimeout
	go app.watchTimeout()

	// 等待超时 (1s + 余量)
	time.Sleep(1500 * time.Millisecond)

	if !quitCalled.Load() {
		t.Error("expected quitFunc to be called on timeout")
	}

	// 验证写入了 timeout 结果
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	var result Result
	json.Unmarshal(data, &result)

	if result.Action != "timeout" {
		t.Errorf("expected action 'timeout', got '%s'", result.Action)
	}
}

func TestWatchTimeout_CancelledByResult(t *testing.T) {
	dir := t.TempDir()
	inputData := createValidInputData()
	inputPath := createTestInputFile(t, dir, inputData)
	outputPath := filepath.Join(dir, "result.json")

	// 创建一个 5 秒超时的 app
	app, err := NewApp(inputPath, outputPath, 5)
	if err != nil {
		t.Fatalf("NewApp failed: %v", err)
	}

	var quitCalled atomicBool
	app.quitFunc = func(ctx context.Context) {
		quitCalled.Store(true)
	}
	app.ctx = context.Background()

	// 启动 watchTimeout
	go app.watchTimeout()

	// 立即发送结果，应该取消超时
	time.Sleep(100 * time.Millisecond)
	app.writeResult(Result{
		Action:             "submit",
		SelectedDirections: []string{},
		UserInput:          "Manual submit",
	})

	// 等待一段时间，确认没有因超时调用 quit
	time.Sleep(200 * time.Millisecond)

	// quit 不应该被调用 (因为 watchTimeout 收到 resultChan 后退出)
	if quitCalled.Load() {
		t.Error("quitFunc should not be called when result is written before timeout")
	}
}

// ========== 并发安全测试 ==========

func TestSubmitRollbackCancel_OnlyFirstWrites(t *testing.T) {
	app, outputPath := createTestApp(t, nil)

	var wg sync.WaitGroup

	// 同时调用 Submit, Rollback, Cancel
	wg.Add(3)
	go func() {
		defer wg.Done()
		app.Submit([]string{"a"}, "submit")
	}()
	go func() {
		defer wg.Done()
		app.Rollback("iter-001", []string{"b"}, "rollback")
	}()
	go func() {
		defer wg.Done()
		app.Cancel()
	}()

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	// 只应该有一个结果被写入
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	var result Result
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to parse result: %v", err)
	}

	// 确认 action 是有效的
	validActions := map[string]bool{"submit": true, "rollback": true, "cancel": true}
	if !validActions[result.Action] {
		t.Errorf("unexpected action: %s", result.Action)
	}
}
