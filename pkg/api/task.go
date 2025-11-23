package api

import (
	"log"
	"net/http"

	"github.com/ArtyomGaribyan/Task-Scheduler/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func HandleTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		Error := db.Task{Error: "Error fetching tasks: " + err.Error()}
		log.Println(Error)
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, Error)
		return
	}
	log.Println("Tasks fetched in handler:", tasks)
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

func UpdateTaskHandler(w http.ResponseWriter, task db.Task) {
	if task.ID == "" {
		Error := db.Task{Error: "Validation error in updating task: missing ID"}
		log.Println(Error)
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, Error)
		return
	}

	log.Println("Updating task:", task.ID, "\n", task.Title, task.Date, task.Comment, task.Repeat)

	err := checkTask(task)
	if err != nil {
		Error := db.Task{Error: "Validation error in updating task: wrong data: " + err.Error()}
		log.Println(Error)
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, Error)
		return
	}

	err = db.UpdateTask(task)
	if err != nil {
		Error := db.Task{Error: "Error in updating task: " + err.Error()}
		log.Println(Error)
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, Error)
		return
	}

	log.Println("Task updated successfully:", task.ID)
	writeJson(w, db.Task{})
}

func GetTaskHandler(w http.ResponseWriter, id string) {
	task := db.Task{ID: id}

	if task.ID == "" {
		Error := db.Task{Error: "Error getting task: missing ID"}
		log.Println(Error)
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, Error)
		return
	}

	taskResieved, err := db.GetTask(task.ID)
	if err != nil {
		Error := db.Task{Error: "Error getting task from DB: " + err.Error()}
		log.Println(Error)
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, Error)
		return
	}
	log.Println("Task received:", taskResieved.ID, taskResieved.Title, taskResieved.Date, taskResieved.Comment, taskResieved.Repeat)
	writeJson(w, taskResieved)
}

func DeleteTaskHandler(w http.ResponseWriter, id string) {
	if id == "" {
		Error := db.Task{Error: "Validation error in deleting task: missing ID"}
		log.Println(Error)
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, Error)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		Error := db.Task{Error: "Error deleting task: " + err.Error()}
		log.Println(Error)
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, Error)
		return
	}

	log.Printf("Successfully deleted task: %s\n", id)
	writeJson(w, db.Task{})
}
