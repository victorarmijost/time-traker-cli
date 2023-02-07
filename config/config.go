package config

import (
	"encoding/json"
	"os"
	"varmijo/time-tracker/utils"
)

type Config struct {
	LogLevel    string  `json:"log_level"`
	WorkingTime float32 `json:"workingTime"`
}

const ConfigFileName = "config.json"

func NewConfig() *Config {
	return &Config{
		LogLevel: "error",
	}
}

func (s *Config) IsComplete() bool {
	if s.WorkingTime == 0 {
		return false
	}

	return true
}

func (s *Config) Load() error {
	configData, err := os.ReadFile(utils.GeAppPath(ConfigFileName))

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

	return os.WriteFile(utils.GeAppPath(ConfigFileName), configData, 0644)
}
