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

type PomodoroState interface {
	Clear()
	GetState() string
	Has() bool
	GetBreakTime() float64
	GetProgress() int
	Start(time.Time)
	End(time.Time)
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

const (
	PomodoroWorkTime      = 25
	PomodoroBreakTime     = 5
	PomodoroLongBreakTime = 15
)

type PomodoroInMemory struct {
	has         bool
	state       string
	accRestTime float64
	count       float64
	startTime   time.Time
}

func NewPomodoroInMemory() *PomodoroInMemory {
	return &PomodoroInMemory{}
}

// TODO: temporary disable pomodoro
func (s *PomodoroInMemory) Has() bool {
	return false
}

func (s *PomodoroInMemory) GetState() string {
	if s.Has() {
		return s.state
	}

	return ""
}

func (s *PomodoroInMemory) Start(stime time.Time) {
	if s.Has() {
		if s.state == "w" {
			return
		}

		if s.GetTimer() > s.accRestTime {
			s.accRestTime = 0
		} else {
			s.accRestTime -= s.GetTimer()
		}

		s.state = "w"
		s.startTime = stime

		s.has = true

		return
	}

	s.state = "w"
	s.startTime = stime

	s.has = true
}

func (s *PomodoroInMemory) GetTimer() float64 {
	if !s.Has() {
		return 0
	}

	return time.Since(s.startTime).Minutes()
}

func (s *PomodoroInMemory) End(etime time.Time) {
	if !s.Has() {
		return
	}

	if s.state != "w" {
		return
	}

	wp := s.GetTimer() / PomodoroWorkTime

	s.count += wp
	s.accRestTime = wp * PomodoroBreakTime

	for s.count >= 4 {
		s.accRestTime += PomodoroLongBreakTime - PomodoroBreakTime
		s.count -= 4
	}

	s.startTime = etime
	s.state = "b"
}

func (s *PomodoroInMemory) GetProgress() int {
	switch s.state {
	case "w":
		return int(s.GetTimer() / PomodoroWorkTime * 100)
	case "b":
		t := s.GetTimer()
		if t < s.accRestTime {
			return int(t / s.accRestTime * 100)
		}
		return 100
	}

	return 0
}

func (s *PomodoroInMemory) GetBreakTime() float64 {
	if s.Has() {
		return s.accRestTime
	}

	return 0
}

func (s *PomodoroInMemory) Clear() {
	s.has = false
}

func TodayIsWeekend() bool {
	return time.Now().Weekday() == time.Saturday || time.Now().Weekday() == time.Sunday
}

type Displayer interface {
	Show(string)
}
