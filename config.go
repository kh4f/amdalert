package main

import ("os"; "encoding/json")

type Config struct {
	MaxTemp int `json:"maxTemp"`
	MaxFanOffTemp int `json:"maxFanOffTemp"`
}

var config = Config{
	MaxTemp: 60,
	MaxFanOffTemp: 40,
}

const configFile = "config.json"
var configModTime int64

func loadConfig() {
	data, err := os.ReadFile(configFile)
	if err != nil {
		saveConfig()
		return
	}
	json.Unmarshal(data, &config)

	info, err := os.Stat(configFile)
	if err != nil { return }
	configModTime = info.ModTime().Unix()
}

func saveConfig() {
	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(configFile, data, 0644)

	info, err := os.Stat(configFile)
	if err == nil {
		configModTime = info.ModTime().Unix()
	}
}

func reloadConfigIfChanged() {
	info, err := os.Stat(configFile)
	if err != nil { return }

	if info.ModTime().Unix() != configModTime {
		loadConfig()
	}
}