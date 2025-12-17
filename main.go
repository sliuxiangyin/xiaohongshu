package main

import (
	"context"
	"embed"
	"xiaohongshu/app/binds/app"
	"xiaohongshu/app/binds/xiaohongshu"
	"xiaohongshu/app/infra/app_context"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed all:scripts
var script embed.FS

func main() {
	// Create an instance of the app structure
	appContext := app_context.NewContext()
	appBind := app.NewApp(appContext)
	xiaohongshuBind := xiaohongshu.NewXiaohongshu(appContext, script)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "xiaohongshu",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			appContext.OnStartup(ctx)
			appBind.Startup(ctx)
			xiaohongshuBind.Startup(ctx)
		},
		OnShutdown: func(ctx context.Context) {

		},
		OnDomReady: func(ctx context.Context) {
			go func() {
				appContext.OnDomReady(ctx)
				xiaohongshuBind.OnDomReady(ctx)
				appContext.OnMount(ctx)
			}()
		},
		Bind: []interface{}{
			appBind,
			xiaohongshuBind,
		},

		Debug: options.Debug{
			OpenInspectorOnStartup: true, // 开发时设为true
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
