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

type SQLiteRecordRepository struct {
	db    *sqlx.DB
	cache *dbCache
}

func NewSQLiteRecordRepository(db *sqlx.DB) *SQLiteRecordRepository {
	return &SQLiteRecordRepository{
		db:    db,
		cache: newDBCache(),
	}
}

func (r *SQLiteRecordRepository) Save(ctx context.Context, record *domain.Record) error {
	return withResetCache(r.cache, func() error {
		dbRecord := DBRecord{
			Id:     record.ID(),
			Date:   record.Date().Format(time.RFC3339),
			Status: record.Status().String(),
			Hours:  record.Hours(),
		}

		_, err := r.db.NamedExecContext(ctx,
			`INSERT INTO records (id, date, status, hours) VALUES (:id, :date, :status, :hours)
		ON CONFLICT(id) DO UPDATE SET date = excluded.date, status = excluded.status, hours = excluded.hours`,
			dbRecord)

		return err
	})
}

func (r *SQLiteRecordRepository) Delete(ctx context.Context, id string) error {
	return withResetCache(r.cache, func() error {
		_, err := r.db.ExecContext(ctx, `DELETE FROM records WHERE id = ?`, id)
		return err
	})
}

func (r *SQLiteRecordRepository) Get(ctx context.Context, id string) (*domain.Record, error) {
	key := fmt.Sprintf("get:%s", id)

	return withCache(r.cache, key, func() (*domain.Record, error) {
		var dbRecord DBRecord

		err := r.db.GetContext(ctx, &dbRecord, `SELECT id, date, status, hours FROM records WHERE id = ?`, id)
		if err != nil {
			return nil, err
		}

		date, err := time.Parse(time.RFC3339, dbRecord.Date)
		if err != nil {
			return nil, err
		}

		return domain.RecreateRecord(dbRecord.Id, date, dbRecord.Status, dbRecord.Hours)
	})
}

func (r *SQLiteRecordRepository) GetAllByDateStatus(ctx context.Context, date time.Time, status domain.RecordStatus) ([]*domain.Record, error) {
	key := fmt.Sprintf("get-all:%s:%s", date.Format("060102"), status.String())

	return withCache(r.cache, key, func() ([]*domain.Record, error) {
		var dbRecords []*DBRecord

		err := r.db.SelectContext(ctx, &dbRecords, `SELECT id, date, status, hours FROM records WHERE SUBSTR(date,1,10) = SUBSTR(?,1,10) AND status = ?`, date.Format(time.RFC3339), status.String())
		if err != nil {
			return nil, err
		}

		records := make([]*domain.Record, len(dbRecords))
		for i, dbRecord := range dbRecords {
			date, err := time.Parse(time.RFC3339, dbRecord.Date)
			if err != nil {
				return nil, err
			}

			record, err := domain.RecreateRecord(dbRecord.Id, date, dbRecord.Status, dbRecord.Hours)
			if err != nil {
				return nil, err
			}
			records[i] = record
		}

		return records, nil
	})
}

func (r *SQLiteRecordRepository) GetAllByStatus(ctx context.Context, status domain.RecordStatus) ([]*domain.Record, error) {
	key := fmt.Sprintf("get-all:%s", status.String())

	return withCache(r.cache, key, func() ([]*domain.Record, error) {
		var dbRecords []*DBRecord
		err := r.db.SelectContext(ctx, &dbRecords, `SELECT id, date, status, hours FROM records WHERE status = ?`, status.String())
		if err != nil {
			return nil, err
		}

		records := make([]*domain.Record, len(dbRecords))
		for i, dbRecord := range dbRecords {
			date, err := time.Parse(time.RFC3339, dbRecord.Date)
			if err != nil {
				return nil, err
			}

			record, err := domain.RecreateRecord(dbRecord.Id, date, dbRecord.Status, dbRecord.Hours)
			if err != nil {
				return nil, err
			}
			records[i] = record
		}

		return records, nil
	})
}

func (r *SQLiteRecordRepository) GetHours(ctx context.Context) (float64, error) {
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
