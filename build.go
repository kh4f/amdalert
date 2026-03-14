//go:build ignore

package main

import ("slices"; "os"; "os/exec")

func main() {
    entry := "./cmd/hotamd"
    var args []string

    switch {
		case slices.Contains(os.Args, "--build"):
			args = []string{"build", "-o", "HotAMD.exe", entry}
		case slices.Contains(os.Args, "--run"):
			args = []string{"run", entry}
    }

    cmd := exec.Command("go", args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    cmd.Run()
}