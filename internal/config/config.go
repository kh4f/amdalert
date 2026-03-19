package config

import (
	"encoding/json"
	"os"
	"time"
)

type Settings struct {
	MaxTemp       int `json:"maxTemp"`
	MaxFanOffTemp int `json:"maxFanOffTemp"`
}

func DefaultSettings() Settings {
	return Settings{
		MaxTemp:       60,
		MaxFanOffTemp: 40,
	}
}

type Store struct {
	path    string
	modTime time.Time
	current Settings
}

func NewStore(path string) *Store {
	return &Store{
		path:    path,
		current: DefaultSettings(),
	}
}

func (s *Store) Load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return s.Save(s.current)
		}
		return err
	}

	settings := DefaultSettings()
	if err := json.Unmarshal(data, &settings); err != nil {
		return err
	}

	info, err := os.Stat(s.path)
	if err != nil {
		return err
	}

	s.current = settings
	s.modTime = info.ModTime()
	return nil
}

func (s *Store) Save(settings Settings) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return err
	}

	info, err := os.Stat(s.path)
	if err != nil {
		return err
	}

	s.current = settings
	s.modTime = info.ModTime()
	return nil
}

func (s *Store) Current() Settings {
	return s.current
}

func (s *Store) ReloadIfChanged() (bool, error) {
	info, err := os.Stat(s.path)
	if err != nil {
		return false, err
	}

	if info.ModTime().Equal(s.modTime) {
		return false, nil
	}

	return true, s.Load()
}
