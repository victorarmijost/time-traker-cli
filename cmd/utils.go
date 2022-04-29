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

func getDateFromText(sdate string) (*time.Time, error) {
	var date time.Time

	switch sdate {
	case "now", "today", "", time.Now().Format("06-01-02"):
		return nil, nil
	case "yesterday":
		date = time.Now().AddDate(0, 0, -1)
	default:
		val, err := strconv.Atoi(sdate)

		if err == nil {
			pdate := time.Now().AddDate(0, 0, val)
			return &pdate, nil
		}

		pdate, err := time.Parse("06-01-02", sdate)

		if err != nil {
			return nil, err
		}

		date = pdate
	}

	return &date, nil
}
