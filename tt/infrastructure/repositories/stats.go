package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"varmijo/time-tracker/tt/domain"

	"github.com/jmoiron/sqlx"
)

type SQLiteStatsRepository struct {
	db    *sqlx.DB
	cache *dbCache
}

func NewSQLiteStatsRepository(db *sqlx.DB) *SQLiteStatsRepository {
	return &SQLiteStatsRepository{
		db:    db,
		cache: newDBCache(),
	}
}

func (r *SQLiteStatsRepository) GetHoursByDate(ctx context.Context, date time.Time) (float64, error) {
	key := fmt.Sprintf("get-hours:%s", date.Format("060102"))

	return withCache(r.cache, key, func() (float64, error) {
		var totalHours *float64

		err := r.db.GetContext(ctx, &totalHours, `SELECT SUM(hours) FROM records WHERE SUBSTR(date,1,10) = SUBSTR(?,1,10)`, date.Format(time.RFC3339))
		if err != nil {
			return 0, err
		}

		if totalHours == nil {
			return 0, nil
		}

		return *totalHours, nil
	})
}

func (r *SQLiteStatsRepository) GetDebt(ctx context.Context, workingTime float64) (float64, error) {
	key := "get-debts"

	dbDebt, err := withCache(r.cache, key, func() (*DBDebt, error) {
		var dbDebt DBDebt
		err := r.db.GetContext(ctx, &dbDebt, `
			SELECT min(date) date,sum(hours) hours FROM records
		`, workingTime)
		if err != nil {
			return nil, err
		}

		return &dbDebt, nil
	})

	if err != nil {
		return 0, err
	}

	startDate, err := time.Parse(time.RFC3339, dbDebt.StartDate)
	if err != nil {
		return 0, err
	}

	expectedHours := getExpectedHours(workingTime, startDate)
	debt := expectedHours - dbDebt.Hours

	trackedHours, err := r.GetTrackedHours(ctx)
	if err != nil {
		return 0, err
	}

	return debt - trackedHours, nil
}

func getExpectedHours(workingTime float64, startDate time.Time) float64 {
	now := time.Now()

	total := 0
	for date := startDate; date.Before(now); date = date.AddDate(0, 0, 1) {
		if domain.IsWeekend(date) {
			continue
		}

		total++
	}

	return float64(total) * workingTime
}

func (r *SQLiteStatsRepository) GetTrackedHours(ctx context.Context) (float64, error) {
	key := "get:hours"
	openRecord, err := withCache(r.cache, key, func() (*domain.OpenRecord, error) {
		var dbOpenRecord DBOpenRecord
		err := r.db.GetContext(ctx, &dbOpenRecord, `SELECT value as date FROM state_variables WHERE key = 'open_record_start_time'`)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		if err != nil {
			return nil, err
		}

		date, err := time.Parse(time.RFC3339, dbOpenRecord.Date)
		if err != nil {
			return nil, err
		}

		return domain.NewOpenRecord(date), nil
	})

	if err != nil {
		return 0, err
	}

	if openRecord == nil {
		return 0, nil
	}

	return openRecord.Hours(), nil
}
