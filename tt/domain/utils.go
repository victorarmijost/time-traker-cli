package domain

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseHour(shour string) (time.Time, error) {
	return time.Parse("15:04", shour)
}

func SetDate(hour time.Time, date time.Time) time.Time {
	newDate := time.Date(date.Year(), date.Month(), date.Day(),
		hour.Hour(), hour.Minute(), 0, 0, date.Location())

	return newDate
}

func ParseDuration(sd string) (float64, error) {
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

func GetDateFromText(sdate string) (time.Time, error) {
	var date time.Time

	switch sdate {
	case "now", "today", "", time.Now().Format("06-01-02"):
		return time.Now(), nil
	case "yesterday":
		date = time.Now().AddDate(0, 0, -1)
	default:
		val, err := strconv.Atoi(sdate)

		if err == nil {
			pdate := time.Now().AddDate(0, 0, val)
			return pdate, nil
		}

		pdate, err := time.Parse("06-01-02", sdate)

		if err != nil {
			return pdate, err
		}

		date = pdate
	}

	return date, nil
}

func SprintList(list []string) string {
	buf := bytes.NewBufferString("")

	if len(list) == 0 {
		fmt.Fprint(buf, "Nothing to show!")
	}

	for i, l := range list {
		if i == len(list)-1 {
			fmt.Fprintf(buf, "%d. %s", i+1, l)
		} else {
			fmt.Fprintf(buf, "%d. %s\n", i+1, l)
		}
	}

	return buf.String()
}

func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}

	return val
}
