package main

import ("os"; "encoding/json")

type Config struct {
	MaxTemp     int `json:"maxTemp"`
	MaxFanOffTemp int `json:"maxFanOffTemp"`
}

var config = Config{
	MaxTemp:     60,
	MaxFanOffTemp: 40,
}

const configFile = "config.json"

func loadConfig() {
	data, err := os.ReadFile(configFile)
	if err != nil {
		saveConfig()
	} else {
		json.Unmarshal(data, &config)
	}
}

func saveConfig() {
	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(configFile, data, 0644)
}