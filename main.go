package main

import (
	"bufio"; "fmt"; "time"; "os"; "strconv"; "os/exec"; "strings"; "syscall";
	w "golang.org/x/sys/windows"
)

var eventPtr = utf16("Global\\HotAMD")

func main() {
	initADL()
	loadConfig()
	if len(os.Args) > 1 && os.Args[1] == "--daemon" { runDaemon() } else { runCLI() }
	destroyADL()
}

func runDaemon() {
	handle, err := w.CreateEvent(nil, 1, 0, eventPtr)
	if err != nil { return }
	defer w.CloseHandle(handle)
	go start()
	w.WaitForSingleObject(handle, w.INFINITE)
}

func runCLI() {
    reader := bufio.NewReader(os.Stdin)

    for {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()

        fmt.Println("♨️ HotAMD")

		temp, fan := readGPU()
		fmt.Printf("\nGPU: %v°C / %v RPM\n", temp, fan)

		fmt.Printf("Max temperature: %v°C\n", config.MaxTemp)
		fmt.Printf("Max fan-off temperature: %v°C\n", config.MaxFanOffTemp)

		fmt.Print("Alerts: ")
		isRunning := isRunning()
        if isRunning { fmt.Println("on") } else { fmt.Println("off") }

		fmt.Print("\n1) ")
		if isRunning { fmt.Print("Disable ") } else { fmt.Print("Enable ") }
		fmt.Println("alerts")

		fmt.Println("2) Set max temperature")
		fmt.Println("3) Set max fan-off temperature")
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
					config.MaxTemp = val
					saveConfig()
				}
			case "3":
				fmt.Print("Enter max fan-off temperature (°C): ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				val, err := strconv.Atoi(input)
				if err != nil {
					fmt.Println("Invalid input")
				} else {
					config.MaxFanOffTemp = val
					saveConfig()
				}
			case "4": return
        }
    }
}

func isRunning() bool {
	handle, err := w.OpenEvent(w.SYNCHRONIZE, false, eventPtr)
	if err != nil { return false }
	w.CloseHandle(handle)
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
	handle, err := w.OpenEvent(w.EVENT_MODIFY_STATE, false, eventPtr)
	if err != nil { return }
	fmt.Println("Stopping daemon...")
	defer w.CloseHandle(handle)
	w.SetEvent(handle)
}