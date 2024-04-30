package state

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"varmijo/time-tracker/utils"
)

type TimeRounder func(float64) float64

type State struct {
	Date            *time.Time  `json:"date"`
	CurrentTask     *Task       `json:"currentTask"`
	TaskTimeRounder TimeRounder `json:"-"`
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

func (s *State) StartRecord(taskName string, comment string, stime *time.Time) error {
	if s.IsWorking() {
		return fmt.Errorf("already working on a task")
	}

	s.CurrentTask = &Task{}

	if stime != nil {
		s.CurrentTask.StartTime = *stime
	} else {
		s.CurrentTask.StartTime = time.Now()
	}

	s.CurrentTask.Comment = comment
	s.CurrentTask.TaskName = taskName

	return nil

}

func (s *State) EndRecord(et *time.Time) (float64, error) {
	if !s.IsWorking() {
		return 0, fmt.Errorf("not working")
	}

	time := s.GetTaskTime(et)

	if time < 0 {
		return 0, fmt.Errorf("wrong end time")
	}

	s.CurrentTask = nil

	return time, nil
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
