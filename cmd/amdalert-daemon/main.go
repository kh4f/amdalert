package main

import (
	"amdalert/internal/adl"
	"amdalert/internal/daemon"
	"amdalert/internal/settings"
)

func main() {
	adl.Init()
	defer adl.Destroy()

	settings.Load()
	daemon.Run()
}
