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

const version = "3.0.2"

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
			updateTemperature(reader, "Enter alert temperature (°C): ", func(current *config.Config, value int) {
				current.AlertTemp = value
			})
		case "3":
			updateTemperature(reader, "Enter fan-off alert temperature (°C): ", func(current *config.Config, value int) {
				current.FanOffAlertTemp = value
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

func updateTemperature(reader *bufio.Reader, prompt string, apply func(*config.Config, int)) {
	fmt.Print(prompt)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	value, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input")
		return
	}

	current := config.Current
	apply(&current, value)
	config.Current = current
	_ = config.Save()
}

func printMenu() {
	temp, fan := adl.ReadGPU()
	current := config.Current

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
		current.AlertTemp,
		current.FanOffAlertTemp,
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
