package api

import (
	"log"
	"net/http"
	"time"

	"github.com/ArtyomGaribyan/Task-Scheduler/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func HandleTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, TasksResp{
			Tasks: []*db.Task{},
		})
		return
	}
	log.Println("Tasks fetched in handler:", tasks)
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

func UpdateTaskHandler(task *db.Task) error {
	err := checkTask(task)
	if err != nil {
		return err
	}

	err = db.UpdateTask(task)
	if err != nil {
		return err
	}

	return nil
}
func TaskDoneHandler(id string) error {
	task, err := db.GetTask(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			return err
		}
		return nil
	}
	
	log.Println("Calculating next date for task:", task)
	task.Date, err = db.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return err
	}
	log.Println("Next date for task:", task.Date)

	err = db.UpdateDate(&task)
	if err != nil {
		return err
	}

	return nil
}
