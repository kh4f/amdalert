package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const filePath = "config.yml"

type Config struct {
	AlertTemp       int `yaml:"alertTemp"`
	FanOffAlertTemp int `yaml:"fanOffAlertTemp"`
}

var (
	Current = Config{
		AlertTemp:       60,
		FanOffAlertTemp: 40,
	}
	modTime time.Time
)

func Load() error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return Save()
		}
		return err
	}

	loaded := Current
	if err := yaml.Unmarshal(data, &loaded); err != nil {
		return err
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	Current = loaded
	modTime = info.ModTime()
	return nil
}

func Save() error {
	data, err := yaml.Marshal(Current)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return err
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	modTime = info.ModTime()
	return nil
}

func ReloadIfChanged() (bool, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}

	if info.ModTime().Equal(modTime) {
		return false, nil
	}

	return true, Load()
}
