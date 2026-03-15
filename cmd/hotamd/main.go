package main

import (
	"bufio"; "fmt"; "time"; "os"; "strconv"; "os/exec"; "strings"; "syscall"
	"golang.org/x/sys/windows"
	"hotamd/internal/config"
	"hotamd/internal/monitor"
	"hotamd/internal/win"
)

var (
	version = "0.3.0"
	eventPtr = win.Utf16("Global\\HotAMD")
)

func main() {
	monitor.InitADL()
	config.LoadConfig()
	if len(os.Args) > 1 && os.Args[1] == "--daemon" { runDaemon() } else { runCLI() }
	monitor.DestroyADL()
}

func runDaemon() {
	handle, err := windows.CreateEvent(nil, 1, 0, eventPtr)
	if err != nil { return }
	defer windows.CloseHandle(handle)
	go monitor.Start()
	windows.WaitForSingleObject(handle, windows.INFINITE)
}

func runCLI() {
    reader := bufio.NewReader(os.Stdin)

    for {
		clearConsole()

        fmt.Println("♨️ HotAMD v" + version)

		temp, fan := monitor.ReadGPU()
		fmt.Printf("\nGPU: %v°C / %v RPM\n", temp, fan)

		fmt.Printf("Max temperature: %v°C\n", config.Config.MaxTemp)
		fmt.Printf("Max fan-off temperature: %v°C\n", config.Config.MaxFanOffTemp)

		fmt.Print("Alerts: ")
		isRunning := isRunning()
        if isRunning { fmt.Println("on") } else { fmt.Println("off") }

		fmt.Print("\n1) ")
		if isRunning { fmt.Print("Disable ") } else { fmt.Print("Enable ") }
		fmt.Println("alerts")

		fmt.Println("2) Set max temp")
		fmt.Println("3) Set max fan-off temp")
		fmt.Println("4) Exit")

        fmt.Print("\n> ")

        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)

        switch input {
			case "1":
				if isRunning { stopDaemon() } else { startDaemon() }
				time.Sleep(300 * time.Millisecond)
			case "2":
				fmt.Print("Enter max temperature (°C): ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				val, err := strconv.Atoi(input)
				if err != nil {
					fmt.Println("Invalid input")
				} else {
					config.Config.MaxTemp = val
					config.SaveConfig()
				}
			case "3":
				fmt.Print("Enter max fan-off temperature (°C): ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				val, err := strconv.Atoi(input)
				if err != nil {
					fmt.Println("Invalid input")
				} else {
					config.Config.MaxFanOffTemp = val
					config.SaveConfig()
				}
			case "4": return
        }
    }
}

func clearConsole() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func isRunning() bool {
	handle, err := windows.OpenEvent(windows.SYNCHRONIZE, false, eventPtr)
	if err != nil { return false }
	windows.CloseHandle(handle)
	return true
}

func startDaemon() {
	if isRunning() { return }
	fmt.Println("Starting daemon...")
	exe, _ := os.Executable()
	cmd := exec.Command(exe, "--daemon")
	cmd.SysProcAttr = &syscall.SysProcAttr{ CreationFlags: 0x00000008 }
	cmd.Start()
}

func stopDaemon() {
	handle, err := windows.OpenEvent(windows.EVENT_MODIFY_STATE, false, eventPtr)
	if err != nil { return }
	fmt.Println("Stopping daemon...")
	defer windows.CloseHandle(handle)
	windows.SetEvent(handle)
}