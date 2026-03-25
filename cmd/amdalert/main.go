package main

import (
	"amdalert/internal/adl"
	"amdalert/internal/cli"
	"amdalert/internal/config"
)

func main() {
	adl.Init()
	defer adl.Destroy()

	config.Load()
	cli.Run()
}
