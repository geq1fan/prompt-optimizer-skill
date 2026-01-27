package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend
var assets embed.FS

func main() {
	// CLI 参数解析
	sessionID := flag.String("session-id", "", "Session ID (required)")
	timeout := flag.Int("timeout", 600, "Timeout in seconds (default: 600)")
	flag.Parse()

	// 检查是否为 Wails 绑定生成模式 (无参数运行)
	// Wails 构建时会先运行程序来生成绑定，此时不传入参数
	isBindingMode := *sessionID == ""

	var app *App
	var err error

	if isBindingMode {
		// 绑定模式：创建空的 App 实例用于生成绑定
		app = &App{
			timeout:    600,
			resultChan: make(chan Result, 1),
		}
	} else {
		// 正常模式：基于 session-id 推断路径
		baseDir := filepath.Join(".", ".claude", "prompt-optimizer", "sessions", *sessionID)

		// 自动创建目录
		if err := os.MkdirAll(baseDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to create session directory: %v\n", err)
			os.Exit(1)
		}

		inputFile := filepath.Join(baseDir, "session.json")
		outputFile := filepath.Join(baseDir, "result.json")

		app, err = NewApp(inputFile, outputFile, *timeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	// 创建 Wails 应用
	err = wails.Run(&options.App{
		Title:       "Prompt Optimizer",
		Width:       1000,
		Height:      700,
		MinWidth:    800,
		MinHeight:   600,
		AlwaysOnTop: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		OnStartup:        app.startup,
		OnBeforeClose:    app.beforeClose,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running app: %v\n", err)
		os.Exit(4)
	}
}
