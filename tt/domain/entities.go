package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Hours float64

func NewHours(hours float64) (Hours, error) {
	if hours <= 0 {
		return 0, fmt.Errorf("hours must be greater than 0")
	}

	return Hours(hours), nil
}

func (r Hours) Float() float64 {
	return float64(r)
}

type Record struct {
	id    string
	date  time.Time
	hours Hours
}

func RecreateRecord(id string, date time.Time, hours float64) (*Record, error) {
	hours = timeRounding(hours)

	hoursValue, err := NewHours(hours)
	if err != nil {
		return nil, err
	}

	return &Record{
		id:    id,
		date:  date,
		hours: hoursValue,
	}, nil
}

func NewCloseRecord(date time.Time, hours float64) (*Record, error) {
	return RecreateRecord(uuid.New().String(), date, hours)
}

func (r *Record) ID() string {
	return r.id
}

func (r *Record) Hours() float64 {
	return r.hours.Float()
}

func (r *Record) Date() time.Time {
	return r.date
}

func (r *Record) UpdateHours(hours float64) error {
	hours = timeRounding(hours)

	hoursValue, err := NewHours(hours)
	if err != nil {
		return err
	}

	r.hours = hoursValue

	return nil
}

type OpenRecord struct {
	startDate time.Time
}

func NewOpenRecord(date time.Time) *OpenRecord {
	return &OpenRecord{
		startDate: date,
	}
}

func (r *OpenRecord) CloseRecord(endDate time.Time) (*Record, error) {
	hours := endDate.Sub(r.startDate).Hours()
	if hours <= 0 {
		return nil, fmt.Errorf("record is empty")
	}

	return NewCloseRecord(r.startDate, hours)
}

func (r *OpenRecord) IsEmpty(endDate time.Time) bool {
	return r.Hours() <= 0
}

func (r *OpenRecord) Date() time.Time {
	return r.startDate
}

func (r *OpenRecord) Hours() float64 {
	return timeRounding(time.Since(r.startDate).Hours())
}

type PromptData interface {
	RefreshData()
	Wt() float64
	Tt() float64
	Dt() float64
	IsWorking() bool
	IsToday() bool
	GetDate() time.Time
}
