package config

import (
	"encoding/json"
	"os"
)

type ConfigT struct {
	MaxTemp       int `json:"maxTemp"`
	MaxFanOffTemp int `json:"maxFanOffTemp"`
}

var Config = ConfigT{
	MaxTemp:       60,
	MaxFanOffTemp: 40,
}

const configFile = "config.json"

var configModTime int64

func Load() {
	data, err := os.ReadFile(configFile)
	if err != nil {
		Save()
		return
	}
	json.Unmarshal(data, &Config)

	info, err := os.Stat(configFile)
	if err != nil {
		return
	}
	configModTime = info.ModTime().Unix()
}

func Save() {
	data, _ := json.MarshalIndent(Config, "", "  ")
	os.WriteFile(configFile, data, 0644)

	info, err := os.Stat(configFile)
	if err == nil {
		configModTime = info.ModTime().Unix()
	}
}

func ReloadIfChanged() {
	info, err := os.Stat(configFile)
	if err != nil {
		return
	}

	if info.ModTime().Unix() != configModTime {
		Load()
	}
}
