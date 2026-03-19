package app

import (
	"amdalert/internal/adl"
	"amdalert/internal/cli"
	"amdalert/internal/config"
	"amdalert/internal/daemon"
	"amdalert/internal/win"
)

type App struct {
	gpu     *adl.Client
	config  *config.Store
	daemon  *daemon.Service
	console *cli.UI
}

func New() *App {
	gpu := adl.NewClient()
	store := config.NewStore("config.json")
	notifier := win.MessageBoxNotifier{}
	service := daemon.NewService(gpu, notifier, store)

	return &App{
		gpu:     gpu,
		config:  store,
		daemon:  service,
		console: cli.NewUI(gpu, store, service),
	}
}

func (a *App) Run(args []string) {
	a.gpu.Init()
	defer a.gpu.Destroy()

	_ = a.config.Load()

	if len(args) > 0 && args[0] == "--daemon" {
		a.daemon.Run()
		return
	}

	a.console.Run()
}
