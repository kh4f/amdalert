package settings

import (
	"encoding/json"
	"os"
	"time"
)

const filePath = "settings.json"

type Settings struct {
	MaxTemp       int `json:"maxTemp"`
	MaxFanOffTemp int `json:"maxFanOffTemp"`
}

var (
	Current = Settings{
		MaxTemp:       60,
		MaxFanOffTemp: 40,
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
	if err := json.Unmarshal(data, &loaded); err != nil {
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
	data, err := json.MarshalIndent(Current, "", "  ")
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
