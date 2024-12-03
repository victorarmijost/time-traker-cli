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

	kern.pomodoro.Start(recTime)

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

	kern.pomodoro.End(endTime)

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

func (kern *App) CommitAll(ctx context.Context, pamount *float64) error {
	var amount float64
	if pamount == nil {
		amount = kern.config.GetWorkTime()
	} else {
		if *pamount < 1 {
			return fmt.Errorf("amount must be greater than 1")
		}
		amount = *pamount
	}

	records, err := kern.records.GetAllByDateStatus(ctx, kern.date.Get(), domain.StatusPending)
	if err != nil {
		return err
	}

	commitedTime, err := kern.records.GetHoursByDateStatus(ctx, kern.date.Get(), domain.StatusCommited)
	if err != nil {
		return err
	}

	if commitedTime >= amount {
		return fmt.Errorf("amount already commited")
	}

	remTime := amount - commitedTime

	for _, record := range records {
		if remTime < record.Hours() {
			newRecordHours := record.Hours() - remTime

			newRecord, err := domain.NewCloseRecord(record.Date(), newRecordHours)
			if err != nil {
				return err
			}

			err = newRecord.SendToPool()
			if err != nil {
				return err
			}

			err = kern.records.Save(ctx, newRecord)
			if err != nil {
				return err
			}

			if remTime == 0 {
				err = kern.records.Delete(ctx, record.ID())
				if err != nil {
					return err
				}

				//If there is not remaining time the time is not sent to the TT
				continue

			} else {
				record.UpdateHours(remTime)

				err = kern.records.Save(ctx, record)
				if err != nil {
					return err
				}
			}
		}

		//If the record is 0 size, is marked as commited, but is not sent
		if record.Hours() <= 0 {
			continue
		}

		err = record.Commit()
		if err != nil {
			return err
		}

		err = kern.records.Save(ctx, record)
		if err != nil {
			return err
		}

		remTime -= record.Hours()
	}

	kern.pomodoro.Clear()

	return nil
}

func (kern *App) SendToPool(ctx context.Context) error {
	records, err := kern.records.GetAllByDateStatus(ctx, kern.date.Get(), domain.StatusPending)
	if err != nil {
		return err
	}

	for _, record := range records {
		err := record.SendToPool()
		if err != nil {
			return err
		}

		err = kern.records.Save(ctx, record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (kern *App) ListLocal(ctx context.Context) ([]string, error) {
	pendingRecords, err := kern.records.GetAllByDateStatus(ctx, kern.date.Get(), domain.StatusPending)

	if err != nil {
		err = fmt.Errorf("can't get records, %w", err)
		return nil, err
	}

	list := []string{}
	for _, r := range pendingRecords {
		list = append(list, fmt.Sprintf("%.2f", r.Hours()))
	}

	commitedRecords, err := kern.records.GetAllByDateStatus(ctx, kern.date.Get(), domain.StatusCommited)

	if err != nil {
		err = fmt.Errorf("can't get records, %w", err)
		return nil, err
	}

	for _, r := range commitedRecords {
		list = append(list, fmt.Sprintf("[%.2f] ✔️", r.Hours()))
	}

	return list, nil
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

func (kern *App) PourePool(ctx context.Context) error {
	records, err := kern.records.GetAllByStatus(ctx, domain.StatusPool)
	if err != nil {
		return err
	}

	for _, record := range records {
		err := record.Poure(kern.date.Get())
		if err != nil {
			return err
		}

		err = kern.records.Save(ctx, record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (kern *App) GetEditableRecords(ctx context.Context) ([]*domain.Record, error) {
	return kern.records.GetAllByDateStatus(ctx, kern.date.Get(), domain.StatusPending)
}
