package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Record struct {
	id     string
	date   time.Time
	status RecordStatus
	hours  float64
}

func RecreateRecord(id string, date time.Time, status string, hours float64) (*Record, error) {
	statusType, err := NewRecordStatus(status)
	if err != nil {
		return nil, err
	}

	hours = timeRounding(hours)

	return &Record{
		id:     id,
		date:   date,
		status: statusType,
		hours:  hours,
	}, nil
}

func NewCloseRecord(date time.Time, hours float64) (*Record, error) {
	return RecreateRecord(uuid.New().String(), date, StatusPending.String(), hours)
}

func (r *Record) ID() string {
	return r.id
}

func (r *Record) Hours() float64 {
	return r.hours
}

func (r *Record) Status() RecordStatus {
	return r.status
}

func (r *Record) Date() time.Time {
	return r.date
}

func (r *Record) SendToPool() error {
	if r.status != StatusPending {
		return fmt.Errorf("record is not pending, can't send to pool")
	}

	r.updateStatus(StatusPool)

	return nil
}

func (r *Record) Poure(date time.Time) error {
	if r.status != StatusPool {
		return fmt.Errorf("record is not in pool, can't poure")
	}

	r.updateDate(date)
	r.updateStatus(StatusPending)

	return nil
}

func (r *Record) Commit() error {
	if r.status != StatusPending {
		return fmt.Errorf("record is not pending, can't commit")
	}

	r.updateStatus(StatusCommited)

	return nil
}

func (r *Record) UpdateHours(hours float64) error {
	if r.status != StatusPending {
		return fmt.Errorf("record is not pending, can't update hours")
	}

	r.hours = timeRounding(hours)

	return nil
}

func (r *Record) updateStatus(status RecordStatus) {
	r.status = status
}

func (r *Record) updateDate(date time.Time) {
	r.date = date
}

type OpenRecord struct {
	r *Record
}

func RecreateOpenRecord(id string, date time.Time) (*OpenRecord, error) {
	r, err := RecreateRecord(id, date, StatusOpen.String(), 0)
	if err != nil {
		return nil, err
	}

	return &OpenRecord{
		r: r,
	}, nil
}

func NewOpenRecord(date time.Time) (*OpenRecord, error) {
	r, err := RecreateRecord(uuid.New().String(), date, StatusOpen.String(), 0)
	if err != nil {
		return nil, err
	}

	return &OpenRecord{
		r: r,
	}, nil
}

func (r *OpenRecord) CloseRecord(date time.Time) (*Record, error) {
	record := r.r

	if record.Status() != StatusOpen {
		return nil, fmt.Errorf("record is not open")
	}

	record.updateStatus(StatusPending)

	err := record.UpdateHours(date.Sub(record.date).Hours())
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (r *OpenRecord) ID() string {
	return r.r.id
}

func (r *OpenRecord) Date() time.Time {
	return r.r.date
}

func (r *OpenRecord) Hours() float64 {
	return timeRounding(time.Since(r.r.date).Hours())
}

type RecordStatus string

const (
	StatusOpen     RecordStatus = "open"     //Records that are open
	StatusPending  RecordStatus = "pending"  //Records pending to commit
	StatusCommited RecordStatus = "commited" //Records that are commited
	StatusPool     RecordStatus = "pool"     //Records that are not attached to a date
)

func NewRecordStatus(status string) (RecordStatus, error) {
	switch s := RecordStatus(status); s {
	case StatusOpen, StatusCommited, StatusPending, StatusPool:
		return s, nil
	default:
		return "", fmt.Errorf("invalid status")
	}
}

func (s RecordStatus) String() string {
	return string(s)
}

type PromptType uint32

const (
	FULL_UPDATE PromptType = iota
	SOFT_UPDATE
)

type Prompt func(PromptType) string
