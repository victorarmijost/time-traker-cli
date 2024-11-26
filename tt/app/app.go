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
	stats    domain.StatsRepository
}

func NewApp(config domain.ConfigRepository, records domain.RecordRepository, track domain.TrackRepository, stats domain.StatsRepository) *App {
	return &App{
		config:   config,
		records:  records,
		track:    track,
		stats:    stats,
		date:     domain.NewDateInMemory(),
		pomodoro: domain.NewPomodoroInMemory(),
	}
}
