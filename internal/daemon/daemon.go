package daemon

import (
	"amdalert/internal/adl"
	"amdalert/internal/settings"
	"amdalert/internal/win"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var eventName = windows.StringToUTF16Ptr("Global\\AMDAlert")

func Run() {
	handle, err := windows.CreateEvent(nil, 1, 0, eventName)
	if err != nil {
		return
	}
	defer windows.CloseHandle(handle)

	go monitorGPU()
	windows.WaitForSingleObject(handle, windows.INFINITE)
}

func IsRunning() bool {
	handle, err := windows.OpenEvent(windows.SYNCHRONIZE, false, eventName)
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
	exe, err := os.Executable()
	if err != nil {
		return
	}

	cmd := exec.Command(exe, "--daemon")
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x00000008}
	cmd.Start()
}

func Stop() {
	handle, err := windows.OpenEvent(windows.EVENT_MODIFY_STATE, false, eventName)
	if err != nil {
		return
	}
	defer windows.CloseHandle(handle)

	fmt.Println("Stopping daemon...")
	windows.SetEvent(handle)
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

func monitorGPU() {
	for {
		_, _ = settings.ReloadIfChanged()

		current := settings.Current
		temp, rpm := adl.ReadGPU()

		if int(temp) > current.MaxFanOffTemp && rpm == 0 {
			win.Alert(fmt.Sprintf("GPU temperature > %d°C and fan is not spinning", current.MaxFanOffTemp))
		} else if int(temp) > current.MaxTemp {
			win.Alert(fmt.Sprintf("GPU temperature > %d°C", current.MaxTemp))
		}

		time.Sleep(10 * time.Second)
	}
}
