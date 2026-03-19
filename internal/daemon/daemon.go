package daemon

import (
	"amdalert/internal/adl"
	"amdalert/internal/config"
	"amdalert/internal/win"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const defaultEventName = "Global\\AMDAlert"

type Service struct {
	gpu       adl.Reader
	notifier  win.Notifier
	config    *config.Store
	eventName *uint16
}

func NewService(gpu adl.Reader, notifier win.Notifier, store *config.Store) *Service {
	return &Service{
		gpu:       gpu,
		notifier:  notifier,
		config:    store,
		eventName: windows.StringToUTF16Ptr(defaultEventName),
	}
}

func (s *Service) Run() {
	handle, err := windows.CreateEvent(nil, 1, 0, s.eventName)
	if err != nil {
		return
	}
	defer windows.CloseHandle(handle)

	go s.monitorGPU()
	windows.WaitForSingleObject(handle, windows.INFINITE)
}

func (s *Service) IsRunning() bool {
	handle, err := windows.OpenEvent(windows.SYNCHRONIZE, false, s.eventName)
	if err != nil {
		return false
	}
	windows.CloseHandle(handle)
	return true
}

func (s *Service) Start() {
	if s.IsRunning() {
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

func (s *Service) Stop() {
	handle, err := windows.OpenEvent(windows.EVENT_MODIFY_STATE, false, s.eventName)
	if err != nil {
		return
	}
	defer windows.CloseHandle(handle)

	fmt.Println("Stopping daemon...")
	windows.SetEvent(handle)
}

func (s *Service) AddToStartup() {
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

func (s *Service) RemoveFromStartup() {
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

func (s *Service) IsInStartup() bool {
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

func (s *Service) monitorGPU() {
	for {
		_, _ = s.config.ReloadIfChanged()

		settings := s.config.Current()
		temp, rpm := s.gpu.ReadGPU()

		if int(temp) > settings.MaxFanOffTemp && rpm == 0 {
			s.notifier.Alert(fmt.Sprintf("GPU temperature > %d°C and fan is not spinning", settings.MaxFanOffTemp))
		} else if int(temp) > settings.MaxTemp {
			s.notifier.Alert(fmt.Sprintf("GPU temperature > %d°C", settings.MaxTemp))
		}

		time.Sleep(10 * time.Second)
	}
}
