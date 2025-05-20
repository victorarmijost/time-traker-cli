package config

import (
	"context"
	"encoding/json"
	"os"
	"time"
	"varmijo/time-tracker/tt/infrastructure/utils"
)

type config struct {
	LogLevel    string  `json:"logLevel"`
	WorkingTime float64 `json:"workingTime"`
}

const ConfigFileName = "config.json"

func MustNewConfig() *config {
	c := &config{
		LogLevel: "error",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.load(ctx)
	if err != nil {
		panic(err)
	}

	return c
}

func (s *config) load(_ context.Context) error {
	configData, err := os.ReadFile(utils.GeAppPath(ConfigFileName))

	if err != nil {
		return err
	}

	c := config{}
	err = json.Unmarshal(configData, &c)

	if err != nil {
		return err
	}

	*s = c

	return nil
}

func (s *config) GetWorkTime() float64 {
	return s.WorkingTime
}

func (s *config) GetLogLevel() string {
	return s.LogLevel
}
