package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ========== 数据结构定义 ==========

// Direction 优化方向
type Direction struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

// UserFeedback 用户反馈
type UserFeedback struct {
	SelectedDirections []string `json:"selectedDirections"`
	UserInput          string   `json:"userInput"`
}

// IterationData 当前迭代数据 (不含用户反馈)
type IterationData struct {
	IterationID         string      `json:"iterationId"`
	OptimizedPrompt     string      `json:"optimizedPrompt"`
	ReviewReport        string      `json:"reviewReport"`
	EvaluationReport    string      `json:"evaluationReport"`
	Score               int         `json:"score"`
	SuggestedDirections []Direction `json:"suggestedDirections,omitempty"`
}

// HistoryItem 历史迭代记录 (含用户反馈)
type HistoryItem struct {
	IterationID      string       `json:"iterationId"`
	OptimizedPrompt  string       `json:"optimizedPrompt"`
	ReviewReport     string       `json:"reviewReport"`
	EvaluationReport string       `json:"evaluationReport"`
	Score            int          `json:"score"`
	UserFeedback     UserFeedback `json:"userFeedback"`
}

// SessionData v4 session 数据结构 (唯一真相源)
type SessionData struct {
	Version        int           `json:"version"`
	SessionID      string        `json:"sessionId"`
	CreatedAt      string        `json:"createdAt"`
	UpdatedAt      string        `json:"updatedAt"`
	ProjectPath    string        `json:"projectPath"`
	Lang           string        `json:"lang"`
	Mode           string        `json:"mode"`
	OriginalPrompt string        `json:"originalPrompt"`
	Current        IterationData `json:"current"`
	History        []HistoryItem `json:"history"`
	LastAction     string        `json:"lastAction,omitempty"`
	Status         string        `json:"status"`
}

// InputData 旧版输入数据结构 (v1-v3 兼容)
type InputData struct {
	Version        int           `json:"version"`
	OriginalPrompt string        `json:"originalPrompt"`
	Current        IterationData `json:"current"`
	History        []HistoryItem `json:"history"`
}

// Result 输出结果结构
type Result struct {
	Action              string   `json:"action"`
	SelectedDirections  []string `json:"selectedDirections"`
	UserInput           string   `json:"userInput"`
	RollbackToIteration string   `json:"rollbackToIteration,omitempty"`
}

// ========== App 结构 ==========

// QuitFunc 退出函数类型 (用于测试时 mock)
type QuitFunc func(ctx context.Context)

// App 应用主结构
type App struct {
	ctx         context.Context
	inputFile   string
	outputFile  string
	timeout     int
	startTime   time.Time
	sessionData *SessionData
	writeOnce   sync.Once
	resultChan  chan Result
	quitFunc    QuitFunc // 可注入的退出函数，默认为 runtime.Quit
}

// NewApp 创建新应用实例
func NewApp(inputFile, outputFile string, timeout int) (*App, error) {
	// 读取输入文件
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %w", err)
	}

	sessionData, err := loadSessionData(data)
	if err != nil {
		return nil, err
	}

	return &App{
		inputFile:   inputFile,
		outputFile:  outputFile,
		timeout:     timeout,
		startTime:   time.Now(),
		sessionData: sessionData,
		resultChan:  make(chan Result, 1),
		quitFunc:    runtime.Quit, // 默认使用 Wails 的 Quit
	}, nil
}

// loadSessionData 加载 session 数据，支持 v1-v4 格式
func loadSessionData(data []byte) (*SessionData, error) {
	// 先检测版本
	var versionCheck struct {
		Version int `json:"version"`
	}
	if err := json.Unmarshal(data, &versionCheck); err != nil {
		return nil, fmt.Errorf("failed to parse version: %w", err)
	}

	// v4 格式直接解析
	if versionCheck.Version >= 4 {
		var sessionData SessionData
		if err := json.Unmarshal(data, &sessionData); err != nil {
			return nil, fmt.Errorf("failed to parse session JSON: %w", err)
		}
		return &sessionData, nil
	}

	// v1-v3 格式需要迁移
	var inputData InputData
	if err := json.Unmarshal(data, &inputData); err != nil {
		return nil, fmt.Errorf("failed to parse input JSON: %w", err)
	}

	return migrateFromLegacy(&inputData), nil
}

// migrateFromLegacy 将旧版 InputData 迁移到 SessionData
func migrateFromLegacy(input *InputData) *SessionData {
	now := time.Now().Format(time.RFC3339)
	return &SessionData{
		Version:        4,
		SessionID:      fmt.Sprintf("session_%d", time.Now().UnixMilli()),
		CreatedAt:      now,
		UpdatedAt:      now,
		ProjectPath:    "",
		Lang:           detectLang(input.OriginalPrompt),
		Mode:           "basic",
		OriginalPrompt: input.OriginalPrompt,
		Current:        input.Current,
		History:        input.History,
		Status:         "active",
	}
}

// detectLang 检测语言
func detectLang(text string) string {
	for _, r := range text {
		if r >= '\u4e00' && r <= '\u9fff' {
			return "cn"
		}
	}
	return "en"
}

// startup Wails 启动回调
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.startTime = time.Now()

	// 启动超时监控
	go a.watchTimeout()
}

// beforeClose 窗口关闭前回调
func (a *App) beforeClose(ctx context.Context) bool {
	// 用户点击 X 按钮，写入 cancel 结果
	a.writeResult(Result{
		Action:             "cancel",
		SelectedDirections: []string{},
		UserInput:          "",
	})
	return false // 允许关闭
}

// watchTimeout 超时监控
func (a *App) watchTimeout() {
	timer := time.NewTimer(time.Duration(a.timeout) * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		// 超时，写入 timeout 结果并关闭窗口
		a.writeResult(Result{
			Action:             "timeout",
			SelectedDirections: []string{},
			UserInput:          "",
		})
		if a.quitFunc != nil {
			a.quitFunc(a.ctx)
		}
	case <-a.resultChan:
		// 已经有结果写入，停止监控
		return
	}
}

// writeResult 原子写入结果文件 (只执行一次)
func (a *App) writeResult(result Result) error {
	var writeErr error
	a.writeOnce.Do(func() {
		// 通知超时监控停止
		select {
		case a.resultChan <- result:
		default:
		}

		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			writeErr = fmt.Errorf("failed to marshal result: %w", err)
			return
		}

		// 原子写入：先写临时文件，再 rename
		dir := filepath.Dir(a.outputFile)
		tmpFile, err := os.CreateTemp(dir, "result-*.tmp")
		if err != nil {
			writeErr = fmt.Errorf("failed to create temp file: %w", err)
			return
		}
		tmpPath := tmpFile.Name()

		if _, err := tmpFile.Write(data); err != nil {
			tmpFile.Close()
			os.Remove(tmpPath)
			writeErr = fmt.Errorf("failed to write temp file: %w", err)
			return
		}

		if err := tmpFile.Close(); err != nil {
			os.Remove(tmpPath)
			writeErr = fmt.Errorf("failed to close temp file: %w", err)
			return
		}

		if err := os.Rename(tmpPath, a.outputFile); err != nil {
			os.Remove(tmpPath)
			writeErr = fmt.Errorf("failed to rename temp file: %w", err)
			return
		}

		// 如果是 submit 或 rollback，更新 session.json
		if result.Action == "submit" || result.Action == "rollback" {
			if err := a.updateInputForNextIteration(result); err != nil {
				// 记录错误但不影响主流程
				fmt.Fprintf(os.Stderr, "Warning: failed to update session.json: %v\n", err)
			}
		}
	})
	return writeErr
}

// updateInputForNextIteration 更新 session.json 为下次迭代做准备
func (a *App) updateInputForNextIteration(result Result) error {
	if a.sessionData == nil || a.inputFile == "" {
		return nil
	}

	// 构造历史记录项
	historyItem := HistoryItem{
		IterationID:      a.sessionData.Current.IterationID,
		OptimizedPrompt:  a.sessionData.Current.OptimizedPrompt,
		ReviewReport:     a.sessionData.Current.ReviewReport,
		EvaluationReport: a.sessionData.Current.EvaluationReport,
		Score:            a.sessionData.Current.Score,
		UserFeedback: UserFeedback{
			SelectedDirections: result.SelectedDirections,
			UserInput:          result.UserInput,
		},
	}

	if result.Action == "rollback" {
		// 回滚：截断 history 到目标版本
		targetIdx := a.findHistoryIndex(result.RollbackToIteration)
		if targetIdx >= 0 && targetIdx < len(a.sessionData.History) {
			a.sessionData.History = a.sessionData.History[:targetIdx+1]
		}
	} else {
		// 提交：将当前版本加入历史
		a.sessionData.History = append(a.sessionData.History, historyItem)
	}

	// 清空 current（等待 Agent 填充新结果）
	a.sessionData.Current = IterationData{}

	// 更新元数据
	a.sessionData.LastAction = result.Action
	a.sessionData.UpdatedAt = time.Now().Format(time.RFC3339)

	// 写入更新后的 session.json
	data, err := json.MarshalIndent(a.sessionData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}

	return os.WriteFile(a.inputFile, data, 0644)
}

// findHistoryIndex 查找历史记录中的迭代 ID 索引
func (a *App) findHistoryIndex(iterationID string) int {
	for i, item := range a.sessionData.History {
		if item.IterationID == iterationID {
			return i
		}
	}
	return -1
}

// ========== 前端绑定方法 ==========

// GetInputData 获取输入数据 (兼容旧前端)
func (a *App) GetInputData() *InputData {
	if a.sessionData == nil {
		return nil
	}
	// 转换为旧格式供前端使用
	return &InputData{
		Version:        a.sessionData.Version,
		OriginalPrompt: a.sessionData.OriginalPrompt,
		Current:        a.sessionData.Current,
		History:        a.sessionData.History,
	}
}

// GetSessionData 获取完整 session 数据 (v4)
func (a *App) GetSessionData() *SessionData {
	return a.sessionData
}

// GetRemainingSeconds 获取剩余超时时间
func (a *App) GetRemainingSeconds() int {
	elapsed := time.Since(a.startTime).Seconds()
	remaining := a.timeout - int(elapsed)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetTimeoutSeconds 获取总超时时间
func (a *App) GetTimeoutSeconds() int {
	return a.timeout
}

// Submit 提交操作
func (a *App) Submit(selectedDirections []string, userInput string) error {
	if selectedDirections == nil {
		selectedDirections = []string{}
	}

	err := a.writeResult(Result{
		Action:             "submit",
		SelectedDirections: selectedDirections,
		UserInput:          userInput,
	})

	if err != nil {
		return err
	}

	// 延迟关闭，让前端显示成功动画 (0.5s)
	go func() {
		time.Sleep(500 * time.Millisecond)
		if a.quitFunc != nil {
			a.quitFunc(a.ctx)
		}
	}()

	return nil
}

// Rollback 回滚操作
func (a *App) Rollback(iterationId string, selectedDirections []string, userInput string) error {
	if selectedDirections == nil {
		selectedDirections = []string{}
	}

	err := a.writeResult(Result{
		Action:              "rollback",
		SelectedDirections:  selectedDirections,
		UserInput:           userInput,
		RollbackToIteration: iterationId,
	})

	if err != nil {
		return err
	}

	// 延迟关闭，让前端显示成功动画 (0.5s)
	go func() {
		time.Sleep(500 * time.Millisecond)
		if a.quitFunc != nil {
			a.quitFunc(a.ctx)
		}
	}()

	return nil
}

// Cancel 取消操作
func (a *App) Cancel() error {
	err := a.writeResult(Result{
		Action:             "cancel",
		SelectedDirections: []string{},
		UserInput:          "",
	})

	if err != nil {
		return err
	}

	if a.quitFunc != nil {
		a.quitFunc(a.ctx)
	}
	return nil
}
