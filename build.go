//go:build ignore

package main

import ("slices"; "os"; "os/exec"; "strings")

func main() {
    var cmd string

    switch {
	case slices.Contains(os.Args, "--build"):
		cmd = "go build -o HotAMD.exe ./cmd/hotamd"
	case slices.Contains(os.Args, "--run"):
		cmd = "go run ./cmd/hotamd"
	case slices.Contains(os.Args, "--release"):
		cmd = "bunx relion -b cmd/hotamd/main.go"
    }

    args := strings.Fields(cmd)
    execCmd := exec.Command(args[0], args[1:]...)
    execCmd.Stdout, execCmd.Stderr, execCmd.Stdin = os.Stdout, os.Stderr, os.Stdin
    execCmd.Run()
}