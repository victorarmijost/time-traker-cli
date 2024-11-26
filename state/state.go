package state

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"varmijo/time-tracker/utils"
)

const (
	PomodoroWorkTime      = 25
	PomodoroBreakTime     = 5
	PomodoroLongBreakTime = 15
)

type TimeRounder func(float64) float64

type State struct {
	Date            *time.Time  `json:"date"`
	CurrentTask     *Task       `json:"currentTask"`
	Pomodoro        *Pomodoro   `json:"pomodoro"`
	TaskTimeRounder TimeRounder `json:"-"`
}

type Pomodoro struct {
	State       string    `json:"state"`
	AccRestTime float64   `json:"accRestTime"`
	Count       float64   `json:"count"`
	StartTime   time.Time `json:"startTime"`
}

type Task struct {
	TaskName  string    `json:"taskName"`
	Comment   string    `json:"comment"`
	StartTime time.Time `json:"startTime"`
}

func NewState(round TimeRounder) *State {
	return &State{
		TaskTimeRounder: round,
	}
}

const statePath = ".tmp"
const stateFile = "state.json"

func (s *State) Load() error {
	stateData, err := os.ReadFile(utils.GeAppPath(statePath) + "/" + stateFile)

	if err != nil {
		return err
	}

	newState := &State{}

	err = json.Unmarshal(stateData, newState)

	if err != nil {
		return err
	}

	rounder := s.TaskTimeRounder

	*s = *newState

	s.TaskTimeRounder = rounder

	return nil
}

func (s *State) Save() error {
	stateData, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = os.WriteFile(utils.GeAppPath(statePath)+"/"+stateFile, stateData, 0644)

	if err != nil {
		return err
	}

	return nil

}

func (s *State) StartRecord(taskName string, comment string, st *time.Time) error {
	if s.IsWorking() {
		return fmt.Errorf("already working on a task")
	}

	s.CurrentTask = &Task{}

	var stime time.Time
	if st != nil {
		stime = *st
	} else {
		stime = time.Now()
	}

	s.CurrentTask.StartTime = stime
	s.CurrentTask.Comment = comment
	s.CurrentTask.TaskName = taskName

	s.StartPomodoro(stime)

	return nil

}

func (s *State) EndRecord(et *time.Time) (float64, error) {
	if !s.IsWorking() {
		return 0, fmt.Errorf("not working")
	}

	duration := s.GetTaskTime(et)

	if duration < 0 {
		return 0, fmt.Errorf("wrong end time")
	}

	s.CurrentTask = nil

	var etime time.Time
	if et != nil {
		etime = *et
	} else {
		etime = time.Now()
	}

	s.EndPomodoro(etime)

	return duration, nil
}

func (s *State) GetTaskTime(et *time.Time) float64 {
	if s.IsWorking() {
		if et == nil {
			return s.TaskTimeRounder(float64(time.Since(s.CurrentTask.StartTime).Hours()))
		}

		return s.TaskTimeRounder(float64(et.Sub(s.CurrentTask.StartTime).Hours()))
	}

	return 0
}

func (s *State) IsWorking() bool {
	return !(s.CurrentTask == nil)
}

func (s *State) GetCurrentTask() (*Task, error) {
	if !s.IsWorking() {
		return nil, fmt.Errorf("not working")
	}

	return s.CurrentTask, nil
}

func (s *State) HasPomodoro() bool {
	return s.Pomodoro != nil
}

func (s *State) GetPomodoroState() string {
	if s.HasPomodoro() {
		return s.Pomodoro.State
	}

	return ""
}

func (s *State) StartPomodoro(stime time.Time) {
	if s.HasPomodoro() {
		if s.Pomodoro.State == "w" {
			return
		}

		if s.GetTimer() > s.Pomodoro.AccRestTime {
			s.Pomodoro.AccRestTime = 0
		} else {
			s.Pomodoro.AccRestTime -= s.GetTimer()
		}

		s.Pomodoro.State = "w"
		s.Pomodoro.StartTime = stime

		return
	}

	s.Pomodoro = &Pomodoro{
		State:       "w",
		AccRestTime: 0,
		StartTime:   stime,
	}
}

func (s *State) GetTimer() float64 {
	if !s.HasPomodoro() {
		return 0
	}

	return time.Since(s.Pomodoro.StartTime).Minutes()
}

func (s *State) EndPomodoro(etime time.Time) {
	if !s.HasPomodoro() {
		return
	}

	if s.Pomodoro.State != "w" {
		return
	}

	wp := s.GetTimer() / PomodoroWorkTime

	s.Pomodoro.Count += wp
	s.Pomodoro.AccRestTime = wp * PomodoroBreakTime

	for s.Pomodoro.Count >= 4 {
		s.Pomodoro.AccRestTime += PomodoroLongBreakTime - PomodoroBreakTime
		s.Pomodoro.Count -= 4
	}

	s.Pomodoro.StartTime = etime
	s.Pomodoro.State = "b"
}

func (s *State) GetStatusProgress() int {
	switch s.Pomodoro.State {
	case "w":
		return int(s.GetTimer() / PomodoroWorkTime * 100)
	case "b":
		t := s.GetTimer()
		if t < s.Pomodoro.AccRestTime {
			return int(t / s.Pomodoro.AccRestTime * 100)
		}
		return 100
	}

	return 0
}

func (s *State) GetBreakTime() float64 {
	if s.HasPomodoro() {
		return s.Pomodoro.AccRestTime
	}

	return 0
}

func (s *State) ClearPomodoro() {
	s.Pomodoro = nil
}
