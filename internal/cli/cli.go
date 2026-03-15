package cli

import (
	"bufio"
	"fmt"
	"amdalert/internal/adl"
	"amdalert/internal/config"
	"amdalert/internal/daemon"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const version = "0.3.0"

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
			fmt.Print("Enter max temperature (°C): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			val, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid input")
			} else {
				config.Config.MaxTemp = val
				config.Save()
			}
		case "3":
			fmt.Print("Enter max fan-off temperature (°C): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			val, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid input")
			} else {
				config.Config.MaxFanOffTemp = val
				config.Save()
			}
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

func printMenu() {
	temp, fan := adl.ReadGPU()

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
		config.Config.MaxTemp,
		config.Config.MaxFanOffTemp,
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
