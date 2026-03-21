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
		cmd = "windres assets/res.rc -O coff -o cmd/amdalert/res.syso && " +
			"windres assets/res.rc -O coff -o cmd/amdalert-daemon/res.syso && " +
			"go build -o AMDAlert.exe ./cmd/amdalert && " +
			"go build -ldflags='-H=windowsgui' -o AMDAlertDaemon.exe ./cmd/amdalert-daemon"
	case slices.Contains(os.Args, "-r"):
		cmd = "go run ./cmd/amdalert"
	case slices.Contains(os.Args, "-f"):
		cmd = "go fmt ./..."
	case slices.Contains(os.Args, "-l"):
		cmd = "bunx relion -b assets/res.rc internal/cli/cli.go"
	}

	execCmd := exec.Command("bash", "-c", cmd)
	execCmd.Stdout, execCmd.Stderr, execCmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	execCmd.Env = append(os.Environ(), "GOOS=windows")
	execCmd.Run()
}
