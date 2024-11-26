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

func (r *SQLiteStatsRepository) GetHoursByDateStatus(ctx context.Context, date time.Time, status domain.RecordStatus) (float64, error) {
	key := fmt.Sprintf("get-hours:%s:%s", date.Format("060102"), status.String())

	return withCache(r.cache, key, func() (float64, error) {
		var totalHours *float64

		err := r.db.GetContext(ctx, &totalHours, `SELECT SUM(hours) FROM records WHERE SUBSTR(date,1,10) = SUBSTR(?,1,10) AND status = ?`, date.Format(time.RFC3339), status.String())
		if err != nil {
			return 0, err
		}

		if totalHours == nil {
			return 0, nil
		}

		return *totalHours, nil
	})
}

func (r *SQLiteStatsRepository) GetHoursByStatus(ctx context.Context, status domain.RecordStatus) (float64, error) {
	key := fmt.Sprintf("get-hours:%s", status.String())

	return withCache(r.cache, key, func() (float64, error) {
		var totalHours *float64
		err := r.db.GetContext(ctx, &totalHours, `SELECT SUM(hours) FROM records WHERE status = ?`, status.String())
		if err != nil {
			return 0, err
		}

		if totalHours == nil {
			return 0, nil
		}

		return *totalHours, nil
	})
}

func (r *SQLiteStatsRepository) GetDebt(ctx context.Context, workingTime float64) (*domain.Debt, error) {
	key := "get-debts"

	dbDebts, err := withCache(r.cache, key, func() ([]*DBDebt, error) {
		var dbDebts []*DBDebt
		err := r.db.SelectContext(ctx, &dbDebts, `
			SELECT date,hours
			FROM(
			SELECT DATE(date,'localtime') as date,?-sum(hours) as hours 
			FROM records r 
			WHERE strftime('%u',DATE(date,'localtime')) BETWEEN '1' and '5'
			group by DATE(date,'localtime')
			)
			WHERE hours > 0
			ORDER BY date DESC;
		`, workingTime)
		if err != nil {
			return nil, err
		}

		return dbDebts, nil
	})

	if err != nil {
		return nil, err
	}

	debt := domain.NewDebt()
	for _, dbDebt := range dbDebts {
		date, err := time.Parse("2006-01-02", dbDebt.Date)
		if err != nil {
			return nil, err
		}

		err = debt.Set(date, dbDebt.Hours)
		if err != nil {
			return nil, err
		}
	}

	trackedHours, err := r.GetTrackedHours(ctx)
	if err != nil {
		return nil, err
	}

	debt.Adjust(trackedHours)

	pool, err := r.GetHoursByStatus(ctx, domain.StatusPool)
	if err != nil {
		return nil, err
	}

	debt.Adjust(pool)

	return debt, nil
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
