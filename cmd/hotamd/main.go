package main

import (
	"os"
	"hotamd/internal/config"
	"hotamd/internal/adl"
	"hotamd/internal/cli"
	"hotamd/internal/daemon"
)

func main() {
	adl.Init()
	config.LoadConfig()
	if len(os.Args) > 1 && os.Args[1] == "--daemon" {
		daemon.Run()
	} else {
		cli.Run()
	}
	adl.Destroy()
}