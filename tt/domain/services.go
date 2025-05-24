package domain

import (
	"fmt"
	"math"
	"time"
)

func FormatDuration(d float64) string {
	h := int(d)
	d -= float64(h)
	m := int(d * 60)

	return fmt.Sprintf("%d:%02d", h, m)
}

// Defines how the tasks time is rounded
func timeRounding(time float64) float64 {
	fact := float64(1) / 60
	return math.Round(time/fact) * fact
}

type Selectable interface {
	GetElement(int) string
	Size() int
}

type DateState interface {
	Get() time.Time
	Set(time.Time)
	IsToday() bool
}

type DateInMemory struct {
	date *time.Time
}

func NewDateInMemory() *DateInMemory {
	return &DateInMemory{}
}

func (d *DateInMemory) Get() time.Time {
	if d.date == nil {
		return time.Now()
	}

	return *d.date
}

func (d *DateInMemory) Set(date time.Time) {
	now := time.Now()
	if date.Year() == now.Year() && date.Month() == now.Month() && date.Day() == now.Day() {
		d.date = nil
		return
	}
	d.date = &date
}

func (d *DateInMemory) IsToday() bool {
	return d.date == nil
}

func TodayIsWeekend() bool {
	return time.Now().Weekday() == time.Saturday || time.Now().Weekday() == time.Sunday
}

type Displayer interface {
	Show(string)
}
