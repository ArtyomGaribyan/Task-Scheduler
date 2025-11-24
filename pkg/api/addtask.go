package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ArtyomGaribyan/Task-Scheduler/pkg/db"
)

func checkTask(task *db.Task) error {
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
			return fmt.Errorf("invalid date format: %w", err)
		}
	}

	var err error
	if task.Date < today {
		if task.Repeat != "" {
			task.Date, err = db.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("next date calculation error: %w", err)
			}
		} else {
			task.Date = today
		}
	}
	return nil
}

func AddTaskHandler(w http.ResponseWriter, task db.Task) {
	log.Println("Adding task:", task)
	err := checkTask(&task)
	if err != nil {
		Error := db.Task{Error: "Validation error for adding task: " + err.Error()}
		log.Println(Error)
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, Error)
		return
	}

	id, err := db.AddTask(task)
	if err != nil {
		Error := db.Task{Error: "Error for adding task: " + err.Error()}
		log.Println(Error)
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, Error)
		return
	}
	task.ID = strconv.Itoa(int(id))
	log.Println("Added task ID:", task.ID)
	writeJson(w, task)
}
