package main

import (
	"amdalert/internal/adl"
	"amdalert/internal/config"
	"amdalert/internal/daemon"
)

func main() {
	adl.Init()
	defer adl.Destroy()

	config.Load()
	daemon.Run()
}
