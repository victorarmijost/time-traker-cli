package app

import (
	"varmijo/time-tracker/tt/domain"
)

type App struct {
	pomodoro domain.PomodoroState
	date     domain.DateState
	config   domain.ConfigRepository
	records  domain.RecordRepository
	track    domain.TrackRepository
}

func NewApp(config domain.ConfigRepository, records domain.RecordRepository, track domain.TrackRepository) *App {
	return &App{
		config:   config,
		records:  records,
		track:    track,
		date:     domain.NewDateInMemory(),
		pomodoro: domain.NewPomodoroInMemory(),
	}
}
