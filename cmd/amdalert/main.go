package main

import (
	"amdalert/internal/adl"
	"amdalert/internal/cli"
	"amdalert/internal/config"
	"amdalert/internal/daemon"
	"os"
)

func main() {
	adl.Init()
	config.Load()
	if len(os.Args) > 1 && os.Args[1] == "--daemon" {
		daemon.Run()
	} else {
		cli.Run()
	}
	adl.Destroy()
}
