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
	id     string
	date   time.Time
	status RecordStatus
	hours  Hours
}

func RecreateRecord(id string, date time.Time, status string, hours float64) (*Record, error) {
	statusType, err := NewRecordStatus(status)
	if err != nil {
		return nil, err
	}

	hours = timeRounding(hours)

	hoursValue, err := NewHours(hours)
	if err != nil {
		return nil, err
	}

	return &Record{
		id:     id,
		date:   date,
		status: statusType,
		hours:  hoursValue,
	}, nil
}

func NewCloseRecord(date time.Time, hours float64) (*Record, error) {
	return RecreateRecord(uuid.New().String(), date, StatusPending.String(), hours)
}

func (r *Record) ID() string {
	return r.id
}

func (r *Record) Hours() float64 {
	return r.hours.Float()
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

func (r *Record) Pour(date time.Time) error {
	if r.status != StatusPool {
		return fmt.Errorf("record is not in pool, can't pour")
	}

	r.updateDate(date)
	r.updateStatus(StatusPending)

	return nil
}

func (r *Record) Commit() error {
	if r.status != StatusPending {
		return fmt.Errorf("record is not pending, can't commit")
	}

	r.updateStatus(StatusCommitted)

	return nil
}

func (r *Record) UpdateHours(hours float64) error {
	if r.status != StatusPending {
		return fmt.Errorf("record is not pending, can't update hours")
	}

	hours = timeRounding(hours)

	hoursValue, err := NewHours(hours)
	if err != nil {
		return err
	}

	r.hours = hoursValue

	return nil
}

func (r *Record) updateStatus(status RecordStatus) {
	r.status = status
}

func (r *Record) updateDate(date time.Time) {
	r.date = date
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

type RecordStatus string

const (
	StatusPending   RecordStatus = "pending"   //Records pending to commit
	StatusCommitted RecordStatus = "committed" //Records that are committed
	StatusPool      RecordStatus = "pool"      //Records that are not attached to a date
)

func NewRecordStatus(status string) (RecordStatus, error) {
	switch s := RecordStatus(status); s {
	case StatusCommitted, StatusPending, StatusPool:
		return s, nil
	default:
		return "", fmt.Errorf("invalid status")
	}
}

func (s RecordStatus) String() string {
	return string(s)
}

type PromptData interface {
	RefreshData()
	Wt() float64
	Ct() float64
	Pt() float64
	Tt() float64
	Dt() float64
	IsWorking() bool
	IsToday() bool
	GetDate() time.Time
}

type Debt struct {
	debt        map[time.Time]float64
	adjustments []float64
}

func NewDebt() *Debt {
	return &Debt{
		debt: make(map[time.Time]float64),
	}
}

func (d *Debt) truckDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
}

func (d *Debt) Set(date time.Time, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	d.debt[d.truckDate(date)] = amount

	return nil
}

func (d *Debt) Adjust(amount float64) {
	d.adjustments = append(d.adjustments, amount)
}

func (d *Debt) Length() int {
	return len(d.debt)
}

func (d *Debt) Do(f func(time.Time, float64)) {
	for date, amount := range d.debt {
		f(date, amount)
	}
}

func (d *Debt) Total() float64 {
	var total float64
	for _, amount := range d.debt {
		total += amount
	}

	for _, amount := range d.adjustments {
		if total < 0 {
			break
		}
		total -= amount
	}

	return total
}
