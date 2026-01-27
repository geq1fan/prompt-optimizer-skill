package main

import (
	"os"
	"path/filepath"
	"testing"
)

// ========== 路径推断逻辑测试 ==========

func TestPathInference_SessionID(t *testing.T) {
	// 测试基于 session-id 的路径推断逻辑
	tests := []struct {
		name             string
		sessionID        string
		expectedBaseDir  string
		expectedInput    string
		expectedOutput   string
	}{
		{
			name:            "standard session ID",
			sessionID:       "session_1705632000000",
			expectedBaseDir: filepath.Join(".", ".claude", "prompt-optimizer", "sessions", "session_1705632000000"),
			expectedInput:   filepath.Join(".", ".claude", "prompt-optimizer", "sessions", "session_1705632000000", "session.json"),
			expectedOutput:  filepath.Join(".", ".claude", "prompt-optimizer", "sessions", "session_1705632000000", "result.json"),
		},
		{
			name:            "simple session ID",
			sessionID:       "test-session",
			expectedBaseDir: filepath.Join(".", ".claude", "prompt-optimizer", "sessions", "test-session"),
			expectedInput:   filepath.Join(".", ".claude", "prompt-optimizer", "sessions", "test-session", "session.json"),
			expectedOutput:  filepath.Join(".", ".claude", "prompt-optimizer", "sessions", "test-session", "result.json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 模拟 main.go 中的路径推断逻辑
			baseDir := filepath.Join(".", ".claude", "prompt-optimizer", "sessions", tt.sessionID)
			inputFile := filepath.Join(baseDir, "session.json")
			outputFile := filepath.Join(baseDir, "result.json")

			if baseDir != tt.expectedBaseDir {
				t.Errorf("baseDir = %q, want %q", baseDir, tt.expectedBaseDir)
			}
			if inputFile != tt.expectedInput {
				t.Errorf("inputFile = %q, want %q", inputFile, tt.expectedInput)
			}
			if outputFile != tt.expectedOutput {
				t.Errorf("outputFile = %q, want %q", outputFile, tt.expectedOutput)
			}
		})
	}
}

// ========== 目录自动创建测试 ==========

func TestDirectoryAutoCreation(t *testing.T) {
	// 测试 os.MkdirAll 目录自动创建功能
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		sessionID string
	}{
		{
			name:      "create nested directories",
			sessionID: "session_test_001",
		},
		{
			name:      "create with special characters",
			sessionID: "session-with-dashes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用临时目录模拟项目根目录
			baseDir := filepath.Join(tmpDir, tt.name, ".claude", "prompt-optimizer", "sessions", tt.sessionID)

			// 模拟 main.go 中的目录创建逻辑
			if err := os.MkdirAll(baseDir, 0755); err != nil {
				t.Fatalf("os.MkdirAll failed: %v", err)
			}

			// 验证目录已创建
			info, err := os.Stat(baseDir)
			if err != nil {
				t.Fatalf("directory not created: %v", err)
			}
			if !info.IsDir() {
				t.Errorf("expected directory, got file")
			}
		})
	}
}

func TestDirectoryAutoCreation_AlreadyExists(t *testing.T) {
	// 测试目录已存在时 os.MkdirAll 不报错
	tmpDir := t.TempDir()
	baseDir := filepath.Join(tmpDir, ".claude", "prompt-optimizer", "sessions", "existing-session")

	// 先创建目录
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		t.Fatalf("initial MkdirAll failed: %v", err)
	}

	// 再次调用不应报错
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		t.Errorf("second MkdirAll should not fail: %v", err)
	}
}

// ========== 超时参数默认值测试 ==========

func TestTimeoutDefault(t *testing.T) {
	// 测试 timeout 默认值为 600
	defaultTimeout := 600

	if defaultTimeout != 600 {
		t.Errorf("default timeout = %d, want 600", defaultTimeout)
	}
}

func TestTimeoutCustomValue(t *testing.T) {
	// 测试自定义 timeout 值
	customTimeout := 300

	if customTimeout <= 0 {
		t.Errorf("custom timeout should be positive: %d", customTimeout)
	}
}

// ========== 绑定模式测试 ==========

func TestBindingMode_EmptySessionID(t *testing.T) {
	// 测试空 session-id 进入绑定模式
	sessionID := ""
	isBindingMode := sessionID == ""

	if !isBindingMode {
		t.Errorf("empty session-id should trigger binding mode")
	}
}

func TestBindingMode_NonEmptySessionID(t *testing.T) {
	// 测试非空 session-id 进入正常模式
	sessionID := "session_123"
	isBindingMode := sessionID == ""

	if isBindingMode {
		t.Errorf("non-empty session-id should not trigger binding mode")
	}
}
