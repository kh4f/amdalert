//go:build ignore

package main

import (
	"os"
	"os/exec"
	"slices"
)

func main() {
	var cmd string

	switch {
	case slices.Contains(os.Args, "-b"):
		cmd = "cd cmd/amdalert && " +
			"windres res.rc -O coff -o res.syso && " +
			"go build -o ../../AMDAlert.exe ."
	case slices.Contains(os.Args, "-r"):
		cmd = "go run ./cmd/amdalert"
	case slices.Contains(os.Args, "-f"):
		cmd = "go fmt ./..."
	case slices.Contains(os.Args, "-l"):
		cmd = "bunx relion -b cmd/amdalert/res.rc internal/cli/cli.go"
	}

	execCmd := exec.Command("bash", "-c", cmd)
	execCmd.Stdout, execCmd.Stderr, execCmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	execCmd.Env = append(os.Environ(), "GOOS=windows")
	execCmd.Run()
}
