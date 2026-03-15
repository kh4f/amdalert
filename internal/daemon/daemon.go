package daemon

import (
	"fmt"; "os"; "os/exec"; "syscall"
	"golang.org/x/sys/windows"
	"hotamd/internal/monitor"
	"hotamd/internal/win"
)

var eventPtr = win.Utf16("Global\\HotAMD")

func Run() {
	handle, err := windows.CreateEvent(nil, 1, 0, eventPtr)
	if err != nil { return }
	defer windows.CloseHandle(handle)
	go monitor.Start()
	windows.WaitForSingleObject(handle, windows.INFINITE)
}

func IsRunning() bool {
	handle, err := windows.OpenEvent(windows.SYNCHRONIZE, false, eventPtr)
	if err != nil { return false }
	windows.CloseHandle(handle)
	return true
}

func Start() {
	if IsRunning() { return }
	fmt.Println("Starting daemon...")
	exe, _ := os.Executable()
	cmd := exec.Command(exe, "--daemon")
	cmd.SysProcAttr = &syscall.SysProcAttr{ CreationFlags: 0x00000008 }
	cmd.Start()
}

func Stop() {
	handle, err := windows.OpenEvent(windows.EVENT_MODIFY_STATE, false, eventPtr)
	if err != nil { return }
	defer windows.CloseHandle(handle)
	fmt.Println("Stopping daemon...")
	windows.SetEvent(handle)
}