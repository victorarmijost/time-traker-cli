package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"varmijo/time-tracker/tt/domain"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteTrackRepository struct {
	db *sqlx.DB
}

func NewSQLiteTrackRepository(db *sqlx.DB) *SQLiteTrackRepository {
	return &SQLiteTrackRepository{db: db}
}

func (r *SQLiteTrackRepository) Save(ctx context.Context, openRecord *domain.OpenRecord) error {
	dbOpenRecord := DBOpenRecord{
		Date: openRecord.Date().Format(time.RFC3339),
	}

	_, err := r.db.NamedExecContext(ctx,
		`INSERT INTO state_variables (key, value) VALUES ('open_record_start_time', :date)`,
		dbOpenRecord)

	return err
}

func (r *SQLiteTrackRepository) Get(ctx context.Context) (*domain.OpenRecord, error) {
	var dbOpenRecord DBOpenRecord
	err := r.db.GetContext(ctx, &dbOpenRecord, `SELECT value as date FROM state_variables WHERE key = 'open_record_start_time'`)
	if err != nil {
		return nil, err
	}

	date, err := time.Parse(time.RFC3339, dbOpenRecord.Date)
	if err != nil {
		return nil, err
	}

	return domain.NewOpenRecord(date), nil
}

func (r *SQLiteTrackRepository) Delete(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM state_variables WHERE key = 'open_record_start_time'`)
	return err
}

func (r *SQLiteTrackRepository) IsWorking(ctx context.Context) bool {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM state_variables WHERE key = 'open_record_start_time'`)
	if err != nil {
		return false
	}
	return count > 0
}

func (r *SQLiteTrackRepository) GetHours(ctx context.Context) (float64, error) {
	var dbOpenRecord DBOpenRecord
	err := r.db.GetContext(ctx, &dbOpenRecord, `SELECT value as date FROM state_variables WHERE key = 'open_record_start_time'`)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	date, err := time.Parse(time.RFC3339, dbOpenRecord.Date)
	if err != nil {
		return 0, err
	}

	openRecord := domain.NewOpenRecord(date)

	return openRecord.Hours(), nil
}
