package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend
var assets embed.FS

func main() {
	// CLI 参数解析
	inputFile := flag.String("input", "", "Input JSON file path (required)")
	outputFile := flag.String("output", "", "Output JSON file path (required)")
	timeout := flag.Int("timeout", 600, "Timeout in seconds (default: 600)")
	flag.Parse()

	// 检查是否为 Wails 绑定生成模式 (无参数运行)
	// Wails 构建时会先运行程序来生成绑定，此时不传入参数
	isBindingMode := *inputFile == "" && *outputFile == ""

	var app *App
	var err error

	if isBindingMode {
		// 绑定模式：创建空的 App 实例用于生成绑定
		app = &App{
			timeout:    600,
			resultChan: make(chan Result, 1),
		}
	} else {
		// 正常模式：验证参数并加载数据
		app, err = NewApp(*inputFile, *outputFile, *timeout)
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
