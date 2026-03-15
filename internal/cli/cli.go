package cli

import (
	"bufio"; "fmt"; "os"; "strconv"; "strings"; "time"; "os/exec"
	"hotamd/internal/config"
	"hotamd/internal/adl"
	"hotamd/internal/daemon"
)

const version = "0.3.0"

func Run() {
    reader := bufio.NewReader(os.Stdin)

    for {
		clearConsole()

        fmt.Println("♨️ HotAMD v" + version)

		temp, fan := adl.ReadGPU()
		fmt.Printf("\nGPU: %v°C / %v RPM\n", temp, fan)

		fmt.Printf("Max temperature: %v°C\n", config.Config.MaxTemp)
		fmt.Printf("Max fan-off temperature: %v°C\n", config.Config.MaxFanOffTemp)

		fmt.Print("Alerts: ")
		isRunning := daemon.IsRunning()
        if isRunning { fmt.Println("on") } else { fmt.Println("off") }

		fmt.Print("\n1) ")
		if isRunning { fmt.Print("Disable ") } else { fmt.Print("Enable ") }
		fmt.Println("alerts")

		fmt.Println("2) Set max temp")
		fmt.Println("3) Set max fan-off temp")
		fmt.Println("4) Exit")

        fmt.Print("\n> ")

        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)

        switch input {
		case "1":
			if isRunning { daemon.Stop() } else { daemon.Start() }
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
				config.SaveConfig()
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
				config.SaveConfig()
			}
		case "4": return
        }
    }
}

func clearConsole() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}