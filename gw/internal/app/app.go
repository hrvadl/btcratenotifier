package app

import "github.com/hrvadl/gw/internal/cfg"

func New() *App {
	return &App{}
}

type App struct {
	cfg cfg.Config
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	return nil
}
