package main

import (
	"amdalert/internal/adl"
	"amdalert/internal/cli"
	"amdalert/internal/daemon"
	"amdalert/internal/settings"
	"amdalert/internal/win"
	"os"
)

func main() {
	adl.Init()
	defer adl.Destroy()

	settings.Load()

	if len(os.Args) > 1 && os.Args[1] == "--daemon" {
		win.HideConsole()
		daemon.Run()
	} else {
		cli.Run()
	}
}
