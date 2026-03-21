package cli

import (
	"amdalert/internal/adl"
	"amdalert/internal/daemon"
	"amdalert/internal/settings"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const version = "1.0.0"

func Run() {
	reader := bufio.NewReader(os.Stdin)

	for {
		clearConsole()
		printMenu()

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			if daemon.IsRunning() {
				daemon.Stop()
			} else {
				daemon.Start()
			}
			time.Sleep(300 * time.Millisecond)
		case "2":
			updateTemperature(reader, "Enter max temperature (°C): ", func(current *settings.Settings, value int) {
				current.MaxTemp = value
			})
		case "3":
			updateTemperature(reader, "Enter max fan-off temperature (°C): ", func(current *settings.Settings, value int) {
				current.MaxFanOffTemp = value
			})
		case "4":
			if daemon.IsInStartup() {
				daemon.RemoveFromStartup()
			} else {
				daemon.AddToStartup()
			}
		case "5":
			return
		}
	}
}

func updateTemperature(reader *bufio.Reader, prompt string, apply func(*settings.Settings, int)) {
	fmt.Print(prompt)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	value, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input")
		return
	}

	current := settings.Current
	apply(&current, value)
	settings.Current = current
	_ = settings.Save()
}

func printMenu() {
	temp, fan := adl.ReadGPU()
	current := settings.Current

	alertsStatus := "off"
	alertsAction := "Enable alerts"
	if daemon.IsRunning() {
		alertsStatus = "on"
		alertsAction = "Disable alerts"
	}

	autostartStatus := "off"
	autostartAction := "Add to startup"
	if daemon.IsInStartup() {
		autostartStatus = "on"
		autostartAction = "Remove from startup"
	}

	fmt.Printf(`🚨 AMDAlert v%s

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
		current.MaxTemp,
		current.MaxFanOffTemp,
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
