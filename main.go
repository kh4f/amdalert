package main

import ("bufio"; "fmt"; "time"; "os"; "os/exec"; "strings"; "syscall"; "golang.org/x/sys/windows")

var eventPtr = utf16("Global\\HotAMD")

func main() {
	initialize()
	if len(os.Args) > 1 && os.Args[1] == "--daemon" {
		runDaemon()
	} else {
		runCLI()
	}
	close()
}

func runDaemon() {
	h, e := windows.CreateEvent(nil, 1, 0, eventPtr)
	if e != nil { return }
	defer windows.CloseHandle(h)
	go start()
	windows.WaitForSingleObject(h, windows.INFINITE)
}

func runCLI() {
    reader := bufio.NewReader(os.Stdin)

    for {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()

        fmt.Println("♨️ HotAMD")

		isRunning := isRunning()


		if isRunning {
			getAdapters()
			temp, fan := readGPU(0)
			fmt.Printf("\nTemperature: %d°C\n", temp)
			fmt.Printf("Fan: %d RPM", fan)
		}
		fmt.Print("\nDaemon: ")
        if isRunning { fmt.Println("running") } else { fmt.Println("not running") }

		fmt.Print("\n1) ")
		if isRunning { fmt.Print("Stop ") } else { fmt.Print("Start ") }
		fmt.Println("daemon")
        fmt.Println("2) Exit")
        fmt.Print("\n> ")

        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)

        switch input {
			case "1":
				if isRunning { stopDaemon() } else { startDaemon() }
				time.Sleep(300 * time.Millisecond)
			case "2": return
        }
    }
}

func isRunning() bool {
	handle, err := windows.OpenEvent(windows.SYNCHRONIZE, false, eventPtr)
	if err != nil { return false }
	windows.CloseHandle(handle)
	return true
}

func startDaemon() {
	if isRunning() { return }
	exe, _ := os.Executable()
	cmd := exec.Command(exe, "--daemon")
	cmd.SysProcAttr = &syscall.SysProcAttr{ CreationFlags: 0x00000008 }
	cmd.Start()
}

func stopDaemon() {
	event, err := windows.OpenEvent(windows.EVENT_MODIFY_STATE, false, eventPtr)
	if err != nil { return }
	defer windows.CloseHandle(event)
	windows.SetEvent(event)
	fmt.Println("Stopping daemon...")
}