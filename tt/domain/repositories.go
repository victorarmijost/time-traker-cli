package domain

import (
	"context"
	"time"
)

type ConfigRepository interface {
	GetLogLevel() string
	GetWorkTime() float64
}

type RecordRepository interface {
	Save(ctx context.Context, r *Record) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*Record, error)
	GetAllByDate(ctx context.Context, date time.Time) ([]*Record, error)
}

type StatsRepository interface {
	GetHoursByDate(ctx context.Context, date time.Time) (float64, error)
	//GetHoursByStatus(ctx context.Context, status RecordStatus) (float64, error)
	GetTrackedHours(ctx context.Context) (float64, error)
	GetDebt(ctx context.Context, workingTime float64) (float64, error)
}

type TrackRepository interface {
	Save(ctx context.Context, t *OpenRecord) error
	Get(ctx context.Context) (*OpenRecord, error)
	Delete(ctx context.Context) error
	IsWorking(ctx context.Context) bool
}
