package app

import (
	"context"
	"fmt"
	"xiaohongshu/app/infra/app_context"
)

// App struct
type App struct {
	appContext *app_context.AppContext
	ctx        context.Context
}

// NewApp creates a new App application struct
func NewApp(appContext *app_context.AppContext) *App {
	return &App{
		appContext: appContext,
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", a.appContext.GetRootPath())
}
