package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func parseHour(shour string, pdate *time.Time) (*time.Time, error) {
	date := time.Now()
	if pdate != nil {
		date = *pdate
	}

	hour, err := time.Parse("15:04", shour)
	if err != nil {
		return nil, err
	}

	newDate := time.Date(date.Year(), date.Month(), date.Day(),
		hour.Hour(), hour.Minute(), 0, 0, date.Location())

	return &newDate, nil
}

func parseDuration(sd string) (float64, error) {
	// Split the string into hours and minutes
	parts := strings.Split(sd, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("wrong duration format")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	return float64(hours) + float64(minutes)/60, nil
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
