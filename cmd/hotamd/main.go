package main

import (
	"os"
	"hotamd/internal/config"
	"hotamd/internal/monitor"
	"hotamd/internal/cli"
	"hotamd/internal/daemon"
)

func main() {
	monitor.InitADL()
	config.LoadConfig()
	if len(os.Args) > 1 && os.Args[1] == "--daemon" {
		daemon.Run()
	} else {
		cli.Run()
	}
	monitor.DestroyADL()
}