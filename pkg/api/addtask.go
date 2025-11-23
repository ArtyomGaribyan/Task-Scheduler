package api

import (
	"fmt"
	"time"

	"github.com/ArtyomGaribyan/Task-Scheduler/pkg/db"
)

func checkTask (task *db.Task) error {
if task == nil {
		return fmt.Errorf("task is nil")
	}
	if task.Title == "" {
		return fmt.Errorf("title is required")
	}
	now := time.Now()
	today := now.Format(db.DateLayout)
	if task.Date == "" {
		task.Date = today
	} else {
		_, err := time.Parse(db.DateLayout, task.Date)
		if err != nil {
			return fmt.Errorf("invalid date format: %v", err)
		}
	}

	var err error
	if task.Date < today {
		if task.Repeat != "" {
			task.Date, err = db.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
		} else {
			task.Date = today
		}
	}
	return nil
}

func addTaskHandler(task *db.Task) (int64, error) {
	err := checkTask(task)
	if err != nil {
		return 0, err
	}

	id, err := db.AddTask(task)
	if err != nil {
		return 0, err
	}

	return id, nil
}
