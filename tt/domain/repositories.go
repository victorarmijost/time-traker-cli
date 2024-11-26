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
	GetAllByDateStatus(ctx context.Context, date time.Time, status RecordStatus) ([]*Record, error)
	GetAllByStatus(ctx context.Context, status RecordStatus) ([]*Record, error)
}

type StatsRepository interface {
	GetHoursByDateStatus(ctx context.Context, date time.Time, status RecordStatus) (float64, error)
	GetHoursByStatus(ctx context.Context, status RecordStatus) (float64, error)
	GetTrackedHours(ctx context.Context) (float64, error)
	GetDebt(ctx context.Context, workingTime float64) (*Debt, error)
}

type TrackRepository interface {
	Save(ctx context.Context, t *OpenRecord) error
	Get(ctx context.Context) (*OpenRecord, error)
	Delete(ctx context.Context) error
	IsWorking(ctx context.Context) bool
}
