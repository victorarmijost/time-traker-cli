package repositories

import (
	"context"
	"time"
	"varmijo/time-tracker/tt/domain"

	"github.com/jmoiron/sqlx"
)

type SQLiteRecordRepository struct {
	db *sqlx.DB
}

func NewSQLiteRecordRepository(db *sqlx.DB) *SQLiteRecordRepository {
	return &SQLiteRecordRepository{db: db}
}

func (r *SQLiteRecordRepository) Save(ctx context.Context, record *domain.Record) error {
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
}

func (r *SQLiteRecordRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM records WHERE id = ?`, id)
	return err
}

func (r *SQLiteRecordRepository) Get(ctx context.Context, id string) (*domain.Record, error) {
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
}

func (r *SQLiteRecordRepository) GetAllByDateStatus(ctx context.Context, date time.Time, status domain.RecordStatus) ([]*domain.Record, error) {
	var dbRecords []*DBRecord
	err := r.db.SelectContext(ctx, &dbRecords, `SELECT id, date, status, hours FROM records WHERE DATE(date) = DATE(?) AND status = ?`, date.Format(time.RFC3339), status.String())
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
}

func (r *SQLiteRecordRepository) GetAllByStatus(ctx context.Context, status domain.RecordStatus) ([]*domain.Record, error) {
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
}

func (r *SQLiteRecordRepository) GetHoursByDateStatus(ctx context.Context, date time.Time, status domain.RecordStatus) (float64, error) {
	var totalHours *float64
	err := r.db.GetContext(ctx, &totalHours, `SELECT SUM(hours) FROM records WHERE DATE(date) = DATE(?) AND status = ?`, date.Format(time.RFC3339), status.String())
	if err != nil {
		return 0, err
	}

	if totalHours == nil {
		return 0, nil
	}

	return *totalHours, nil
}

func (r *SQLiteRecordRepository) GetHoursByStatus(ctx context.Context, status domain.RecordStatus) (float64, error) {
	var totalHours *float64
	err := r.db.GetContext(ctx, &totalHours, `SELECT SUM(hours) FROM records WHERE status = ?`, status.String())
	if err != nil {
		return 0, err
	}

	if totalHours == nil {
		return 0, nil
	}

	return *totalHours, nil
}
