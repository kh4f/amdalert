package daemon

import (
	"fmt"
	"hotamd/internal/adl"
	"hotamd/internal/config"
	"hotamd/internal/win"
	"os"
	"os/exec"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

var eventPtr = win.Utf16("Global\\HotAMD")

func Run() {
	handle, err := windows.CreateEvent(nil, 1, 0, eventPtr)
	if err != nil {
		return
	}
	defer windows.CloseHandle(handle)
	go monitorGPU()
	windows.WaitForSingleObject(handle, windows.INFINITE)
}

func IsRunning() bool {
	handle, err := windows.OpenEvent(windows.SYNCHRONIZE, false, eventPtr)
	if err != nil {
		return false
	}
	windows.CloseHandle(handle)
	return true
}

func Start() {
	if IsRunning() {
		return
	}
	fmt.Println("Starting daemon...")
	exe, _ := os.Executable()
	cmd := exec.Command(exe, "--daemon")
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x00000008}
	cmd.Start()
}

func Stop() {
	handle, err := windows.OpenEvent(windows.EVENT_MODIFY_STATE, false, eventPtr)
	if err != nil {
		return
	}
	defer windows.CloseHandle(handle)
	fmt.Println("Stopping daemon...")
	windows.SetEvent(handle)
}

func monitorGPU() {
	for {
		config.ReloadIfChanged()

		temp, rpm := adl.ReadGPU()
		if int(temp) > config.Config.MaxFanOffTemp && rpm == 0 {
			win.Alert("GPU temperature > " + fmt.Sprintf("%v", config.Config.MaxFanOffTemp) + "°C and fan is not spinning")
		} else if int(temp) > config.Config.MaxTemp {
			win.Alert("GPU temperature > " + fmt.Sprintf("%v", config.Config.MaxTemp) + "°C")
		}

		time.Sleep(10 * time.Second)
	}
}
