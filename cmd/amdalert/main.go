package main

import (
	"amdalert/internal/app"
	"os"
)

func main() {
	app.New().Run(os.Args[1:])
}
