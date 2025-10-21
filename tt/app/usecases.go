package app

import (
	"context"
	"fmt"
	"time"

	"varmijo/time-tracker/tt/domain"
)

// Add a new record
func (kern *App) AddRecord(ctx context.Context, hours float64) error {
	record, err := domain.NewCloseRecord(kern.date.Get(), hours)
	if err != nil {
		return fmt.Errorf("error creating new record, %w", err)
	}

	err = kern.records.Save(ctx, record)
	if err != nil {
		return fmt.Errorf("new record can't be inserted, %w", err)
	}

	return nil
}

func (kern *App) startRecordWithDate(ctx context.Context, recTime time.Time) error {
	if kern.track.IsWorking(ctx) {
		return fmt.Errorf("record already started")
	}

	record := domain.NewOpenRecord(recTime)

	err := kern.track.Save(ctx, record)
	if err != nil {
		return fmt.Errorf("error saving new record, %w", err)
	}

	return nil
}

// Add a new record
func (kern *App) StartRecord(ctx context.Context) error {
	if !kern.date.IsToday() {
		return fmt.Errorf("wrong date, change back to today")
	}

	recTime := time.Now()

	return kern.startRecordWithDate(ctx, recTime)
}

func (kern *App) StartRecordAt(ctx context.Context, hour time.Time) error {
	recTime := domain.SetDate(hour, kern.date.Get())

	return kern.startRecordWithDate(ctx, recTime)
}

func (kern *App) stopRecordWithDate(ctx context.Context, endTime time.Time) (float64, error) {
	if !kern.track.IsWorking(ctx) {
		return 0, fmt.Errorf("record not started")
	}

	openRecord, err := kern.track.Get(ctx)
	if err != nil {
		return 0, err
	}

	hours := 0.0

	if !openRecord.IsEmpty(endTime) {
		record, err := openRecord.CloseRecord(endTime)
		if err != nil {
			return 0, fmt.Errorf("can't close record, %w", err)
		}

		hours = record.Hours()

		err = kern.records.Save(ctx, record)
		if err != nil {
			return 0, fmt.Errorf("error inserting new record, %w", err)
		}
	}

	err = kern.track.Delete(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting open record, %w", err)
	}

	return hours, nil
}

func (kern *App) StopRecord(ctx context.Context) (float64, error) {
	endTime := time.Now()

	return kern.stopRecordWithDate(ctx, endTime)
}

func (kern *App) StopRecordAt(ctx context.Context, hour time.Time) (float64, error) {
	endTime := domain.SetDate(hour, kern.date.Get())

	return kern.stopRecordWithDate(ctx, endTime)
}

func (kern *App) ChangeDate(ctx context.Context, date time.Time) error {
	kern.date.Set(date)

	return nil
}

func (kern *App) DropRecord(ctx context.Context) (float64, error) {
	if !kern.track.IsWorking(ctx) {
		return 0, fmt.Errorf("record not started")
	}

	openRecord, err := kern.track.Get(ctx)
	if err != nil {
		return 0, err
	}

	record, err := openRecord.CloseRecord(time.Now())
	if err != nil {
		return 0, fmt.Errorf("can't close record, %w", err)
	}

	err = kern.track.Delete(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting open record, %w", err)
	}

	return record.Hours(), nil
}
