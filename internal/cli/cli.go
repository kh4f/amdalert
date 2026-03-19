package cli

import (
	"amdalert/internal/adl"
	"amdalert/internal/config"
	"amdalert/internal/daemon"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const version = "1.0.0"

type UI struct {
	gpu     adl.Reader
	config  *config.Store
	service *daemon.Service
}

func NewUI(gpu adl.Reader, store *config.Store, service *daemon.Service) *UI {
	return &UI{
		gpu:     gpu,
		config:  store,
		service: service,
	}
}

func (ui *UI) Run() {
	reader := bufio.NewReader(os.Stdin)

	for {
		clearConsole()
		ui.printMenu()

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			if ui.service.IsRunning() {
				ui.service.Stop()
			} else {
				ui.service.Start()
			}
			time.Sleep(300 * time.Millisecond)
		case "2":
			ui.updateTemperature(reader, "Enter max temperature (°C): ", func(settings *config.Settings, value int) {
				settings.MaxTemp = value
			})
		case "3":
			ui.updateTemperature(reader, "Enter max fan-off temperature (°C): ", func(settings *config.Settings, value int) {
				settings.MaxFanOffTemp = value
			})
		case "4":
			if ui.service.IsInStartup() {
				ui.service.RemoveFromStartup()
			} else {
				ui.service.AddToStartup()
			}
		case "5":
			return
		}
	}
}

func (ui *UI) updateTemperature(reader *bufio.Reader, prompt string, apply func(*config.Settings, int)) {
	fmt.Print(prompt)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	value, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input")
		return
	}

	settings := ui.config.Current()
	apply(&settings, value)
	_ = ui.config.Save(settings)
}

func (ui *UI) printMenu() {
	temp, fan := ui.gpu.ReadGPU()
	settings := ui.config.Current()

	alertsStatus := "off"
	alertsAction := "Enable alerts"
	if ui.service.IsRunning() {
		alertsStatus = "on"
		alertsAction = "Disable alerts"
	}

	autostartStatus := "off"
	autostartAction := "Add to startup"
	if ui.service.IsInStartup() {
		autostartStatus = "on"
		autostartAction = "Remove from startup"
	}

	fmt.Printf(`AMDAlert v%s

Status
  GPU temp:       %d°C
  Fan speed:      %d RPM
  Alerts:         %s
  Alert temp:     %d°C
  Fan-off alert:  %d°C
  Autostart:      %s

Actions
  1) %s
  2) Set alert temp
  3) Set fan-off alert temp
  4) %s
  5) Exit

> `,
		version,
		temp,
		fan,
		alertsStatus,
		settings.MaxTemp,
		settings.MaxFanOffTemp,
		autostartStatus,
		alertsAction,
		autostartAction,
	)
}

func clearConsole() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
