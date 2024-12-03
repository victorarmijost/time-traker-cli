package repositories

import (
	"context"
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

func (r *SQLiteRecordRepository) GetHoursByDateStatus(ctx context.Context, date time.Time, status domain.RecordStatus) (float64, error) {
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

func (r *SQLiteRecordRepository) GetHoursByStatus(ctx context.Context, status domain.RecordStatus) (float64, error) {
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

func (r *SQLiteRecordRepository) GetDebts(ctx context.Context) ([]*domain.Debt, error) {
	key := "get-debts"

	return withCache(r.cache, key, func() ([]*domain.Debt, error) {
		var dbDebts []*DBDebt
		err := r.db.SelectContext(ctx, &dbDebts, `
			SELECT date,hours
			FROM(
			SELECT DATE(date,'localtime') as date,8-sum(hours) as hours 
			FROM records r 
			WHERE strftime('%u',DATE(date,'localtime')) BETWEEN '1' and '5'
			group by DATE(date,'localtime')
			)
			WHERE hours > 0
			ORDER BY date DESC;
		`)
		if err != nil {
			return nil, err
		}

		debts := make([]*domain.Debt, len(dbDebts))
		for i, dbDebt := range dbDebts {
			date, err := time.Parse("2006-01-02", dbDebt.Date)
			if err != nil {
				return nil, err
			}

			debt, err := domain.NewDebt(date, dbDebt.Hours)
			if err != nil {
				return nil, err
			}

			debts[i] = debt
		}

		return debts, nil
	})
}
