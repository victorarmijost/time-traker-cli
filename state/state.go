package state

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"
	"varmijo/time-tracker/utils"
)

type State struct {
	Date        *time.Time `json:"date"`
	CurrentTask *Task      `json:"currentTask"`
}

type Task struct {
	Id        int64     `json:"id"`
	Comment   string    `json:"comment"`
	StartTime time.Time `json:"startTime"`
}

func NewState() *State {
	return &State{}
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

	*s = *newState

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

func (s *State) StartRecord(id int64, comment string, stime *time.Time) error {
	if s.IsWorking() {
		return fmt.Errorf("already working on a task")
	}

	s.CurrentTask = &Task{}

	if stime != nil {
		s.CurrentTask.StartTime = *stime
	} else {
		s.CurrentTask.StartTime = time.Now()
	}

	s.CurrentTask.Id = id
	s.CurrentTask.Comment = comment

	return nil

}

func (s *State) EndRecord(et *time.Time) (float32, error) {
	if !s.IsWorking() {
		return 0, fmt.Errorf("not working")
	}

	time := s.GetTaskTime(et)
	s.CurrentTask = nil

	return time, nil
}

func (s *State) GetTaskTime(et *time.Time) float32 {
	if s.IsWorking() {
		if et == nil {
			return float32(math.Ceil(time.Since(s.CurrentTask.StartTime).Hours()/0.25) * 0.25)
		}

		return float32(math.Ceil(et.Sub(s.CurrentTask.StartTime).Hours()/0.25) * 0.25)
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
