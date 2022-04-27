package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ProjectId    int64   `json:"projectId"`
	FocalPointId int64   `json:"focalPointId"`
	WorkingTime  float32 `json:"workingTime"`
	Email        string  `json:"email"`
	Password     string  `json:"password"`
}

func NewConfig() *Config {
	return &Config{}
}

func (s *Config) IsComplete() bool {
	if s.ProjectId == 0 {
		return false
	}

	if s.FocalPointId == 0 {
		return false
	}

	if s.Email == "" {
		return false
	}

	return true
}

func (s *Config) Load() error {
	configData, err := os.ReadFile("config.json")

	if err != nil {
		return err
	}

	c := Config{}
	err = json.Unmarshal(configData, &c)

	if err != nil {
		return err
	}

	*s = c

	return nil
}

func (s *Config) Save() error {
	configData, err := json.MarshalIndent(s, "", "\t")

	if err != nil {
		return err
	}

	return os.WriteFile("config.json", configData, 0644)
}
