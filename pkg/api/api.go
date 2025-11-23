package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ArtyomGaribyan/Task-Scheduler/pkg/db"
)

func writeJson(w http.ResponseWriter, data any) {
	fmt.Printf("\n")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	var parsedNow time.Time
	var err error
	if now == "" {
		parsedNow = time.Now()
	} else {
		parsedNow, err = time.Parse(db.DateLayout, now)
		if err != nil {
			w.Write([]byte(""))
			return
		}
	}

	nextDate, err := db.NextDate(parsedNow, date, repeat)
	if err != nil {
		w.Write([]byte(""))
		return
	}

	w.Write([]byte(nextDate))
}

func HandleTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	method := r.Method
	log.Println("Method:", method)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		task.ID = r.URL.Query().Get("id")
		task.Title = r.URL.Query().Get("title")
		task.Date = r.URL.Query().Get("date")
		task.Comment = r.URL.Query().Get("comment")
		task.Repeat = r.URL.Query().Get("repeat")
		if task != (db.Task{}) {
			goto If_close
		}

		task.Error = "Invalid request body: " + err.Error()
		log.Println(task.Error)
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, task)
		return
	}
If_close:

	log.Println("Received task:", task)

	switch method {
	case http.MethodGet:
		if task.ID == "" {
			task.Error = "Error getting task: missing ID"
			log.Println(task.Error)
			w.WriteHeader(http.StatusInternalServerError)
			writeJson(w, task)
			return
		}

		taskResieved, err := db.GetTask(task.ID)
		log.Println("Task received:", taskResieved.ID, taskResieved.Title, taskResieved.Date, taskResieved.Comment, taskResieved.Repeat)
		if err != nil {
			task.Error = "Error getting task: " + err.Error()
			log.Println(task.Error)
			w.WriteHeader(http.StatusInternalServerError)
			writeJson(w, task)
			return
		}
		writeJson(w, taskResieved)
		return

	case http.MethodPost:
		log.Println("Adding task:", task)
		id, err := addTaskHandler(&task)
		if err != nil {
			task.Error = "Error adding task: " + err.Error()
			log.Println(task.Error)
			w.WriteHeader(http.StatusInternalServerError)
			writeJson(w, task)
			return
		}
		task.ID = strconv.Itoa(int(id))
		log.Println("Added task ID:", task.ID)
		writeJson(w, task)
		return

	case http.MethodPut:
		if task.ID == "" {
			task.Error = "Error updating task: missing ID"
			log.Println(task.Error)
			w.WriteHeader(http.StatusInternalServerError)
			writeJson(w, task)
			return
		}

		log.Println("Updating task:", task.ID, "\n", task.Title, task.Date, task.Comment, task.Repeat)
		err := UpdateTaskHandler(&task)
		if err != nil {
			task.Error = "Error updating task: " + err.Error()
			log.Println(task.Error)
			w.WriteHeader(http.StatusInternalServerError)
			writeJson(w, task)
			return
		}
		writeJson(w, db.Task{})
	case http.MethodDelete:
		if task.ID == "" {
			task.Error = "Error deleting task: missing ID"
			log.Println(task.Error)
			w.WriteHeader(http.StatusInternalServerError)
			writeJson(w, task)
			return
		}

		err := db.DeleteTask(task.ID)
		if err != nil {
			task.Error = "Error deleting task: " + err.Error()
			log.Println(task.Error)
			w.WriteHeader(http.StatusInternalServerError)
			writeJson(w, task)
			return
		}

		log.Printf("Successfully deleted task: %s\n", task.ID)
		writeJson(w, db.Task{})

	default:
		task.Error = "Method not allowed"
		log.Println(task.Error)
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, task)
		return
	}
}

func HandleTaskDone(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	method := r.Method
	log.Println("Method:", method)
	if method != http.MethodPost {
		log.Println("Method not allowed")
		w.WriteHeader(http.StatusInternalServerError)
		task.Error = "Method not allowed"
		writeJson(w, task)
		return
	}

	task.ID = r.URL.Query().Get("id")
	if task.ID == "" {
		task.Error = "Error marking task as done: missing ID"
		log.Println(task.Error)
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, task)
		return
	}
	log.Println("Marking task as done", task.ID)

	err := TaskDoneHandler(task.ID)
	if err != nil {
		task.Error = "Error marking task as done: " + err.Error()
		log.Println(task.Error)
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, task)
		return
	}

	log.Println("Successfully marked as done: ", task.ID)
	writeJson(w, db.Task{})
}

func Init() {
	http.Handle("/", http.FileServer(http.Dir("web")))
	http.HandleFunc("/api/nextdate", HandleNextDate)
	http.HandleFunc("api/nextdate", HandleNextDate)
	http.HandleFunc("/api/task", HandleTask)
	http.HandleFunc("api/task", HandleTask)
	http.HandleFunc("/api/tasks", HandleTasks)
	http.HandleFunc("api/tasks", HandleTasks)
	http.HandleFunc("/api/task/done", HandleTaskDone)
	http.HandleFunc("api/task/done", HandleTaskDone)
}
