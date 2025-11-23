package db

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const DateLayout = "20060102"

func caseD(now, date time.Time, repeatSplitted []string) (time.Time, error) {
	if len(repeatSplitted) != 2 {
		return time.Time{}, fmt.Errorf("invalid date format")
	}

	days, err := strconv.Atoi(repeatSplitted[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %v", err)
	}

	if days <= 0 || days > 400 {
		return time.Time{}, fmt.Errorf("invalid date format: days value too large")
	}

	log.Println("Calculating next date by days:", days, "from", date, "with now =", now)
	date = date.AddDate(0, 0, days)
	for date.Before(now) || date.Equal(now) {
		date = date.AddDate(0, 0, days)
	}
	return date, nil
}

func caseY(now, date time.Time) time.Time {
	date = date.AddDate(1, 0, 0)
	for date.Before(now) || date.Equal(now) {
		date = date.AddDate(1, 0, 0)
	}
	return date
}

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	var nextDate time.Time
	dateStart, err := time.Parse(DateLayout, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}

	repeatSplitted := strings.Fields(repeat)
	if len(repeatSplitted) == 0 || repeatSplitted[0] == "" {
		return "", fmt.Errorf("empty repeat pattern")
	}
	if len(repeatSplitted) > 2 {
		return "", fmt.Errorf("invalid date format")
	}
	log.Println("Repeat pattern: " + strings.Join(repeatSplitted, "~") + "!")

	switch repeatSplitted[0] {
	case "d":
		nextDate, err = caseD(now, dateStart, repeatSplitted)
		if err != nil {
			return "", err
		}
	case "y":
		nextDate = caseY(now, dateStart)
	default:
		return "", fmt.Errorf("invalid date format")
	}

	return nextDate.Format(DateLayout), nil
}
