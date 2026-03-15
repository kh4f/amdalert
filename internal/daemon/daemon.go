package daemon

import (
	"fmt"
	"amdalert/internal/adl"
	"amdalert/internal/config"
	"amdalert/internal/win"
	"os"
	"os/exec"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var eventPtr = win.Utf16("Global\\AMDAlert")

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

func AddToStartup() {
	exePath, err := os.Executable()
	if err != nil {
		return
	}

	key, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE,
	)
	if err != nil {
		return
	}
	defer key.Close()

	key.SetStringValue("AMDAlert", exePath+" --daemon")
}

func RemoveFromStartup() {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE,
	)
	if err != nil {
		return
	}
	defer key.Close()

	key.DeleteValue("AMDAlert")
}

func IsInStartup() bool {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return false
	}
	defer key.Close()

	_, _, err = key.GetStringValue("AMDAlert")
	return err == nil
}
