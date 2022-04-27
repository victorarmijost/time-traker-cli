package main

import (
	"fmt"
	"math"
	"strconv"
	"time"
	"varmijo/time-tracker/bairestt"
)

func formatWorkHours(shours string) (float32, error) {
	hours, err := strconv.ParseFloat(shours, 32)

	if err != nil {
		return 0, err
	}

	return float32(math.Ceil(hours/0.25) * 0.25), nil
}

func parseHour(shour string) (*time.Time, error) {
	date := time.Now()
	hour, err := time.Parse("15:04", shour)
	if err != nil {
		return nil, err
	}

	newDate := time.Date(date.Year(), date.Month(), date.Day(),
		hour.Hour(), hour.Minute(), 0, 0, date.Location())

	return &newDate, nil
}

func validateDescriptionId(tasks []bairestt.TaskInfo, sid string) (int64, error) {
	id, err := strconv.ParseInt(sid, 10, 64)

	if err != nil {
		return 0, err
	}

	for _, t := range tasks {
		if t.Id == id {
			break
		}
	}

	if id == 0 {
		return 0, fmt.Errorf("not found")
	}

	return id, nil
}

func getTaskDetails(tasks []bairestt.TaskInfo, id int64) (*bairestt.TaskInfo, error) {
	for _, t := range tasks {
		if t.Id == id {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("not found")

}
