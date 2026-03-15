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
	case slices.Contains(os.Args, "--build"):
		cmd = "cd cmd/hotamd && " +
			"windres res.rc -O coff -o res.syso && " +
			"go build -o HotAMD.exe ./cmd/hotamd"
	case slices.Contains(os.Args, "--run"):
		cmd = "go run ./cmd/hotamd"
	case slices.Contains(os.Args, "--release"):
		cmd = "bunx relion -b internal/cli/cli.go"
	}

	execCmd := exec.Command("cmd", "/C", cmd)
	execCmd.Stdout, execCmd.Stderr, execCmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	execCmd.Run()
}
