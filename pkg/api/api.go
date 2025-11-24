package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ArtyomGaribyan/Task-Scheduler/pkg/db"
)

func writeJson(w http.ResponseWriter, data any) {
	fmt.Printf("\n") // for better readability in logs
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
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
			Error := "Next date parse error: " + err.Error()
			log.Println(Error)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(Error))
			return
		}
	}

	nextDate, err := db.NextDate(parsedNow, date, repeat)
	if err != nil {
		Error := "Next date calculation error: " + err.Error()
		log.Println(Error)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(Error))
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
	if err != nil && err.Error() != "EOF" {
		Error := db.Task{Error: "Invalid request body: " + err.Error()}
		log.Println(Error)
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, Error)
		return
	}
	if task.ID == "" {
		task.ID = r.URL.Query().Get("id")
	}

	log.Println("Received task:", task)

	switch method {
	case http.MethodGet:
		GetTaskHandler(w, task.ID)
	case http.MethodPost:
		AddTaskHandler(w, task)
	case http.MethodPut:
		UpdateTaskHandler(w, task)
	case http.MethodDelete:
		DeleteTaskHandler(w, task.ID)
	default:
		Error := db.Task{Error: "Method not allowed"}
		log.Println(Error)
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJson(w, Error)
	}
}

func HandleTaskDone(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	method := r.Method
	log.Println("Method:", method)
	if method != http.MethodPost {
		Error := db.Task{Error: "Method not allowed"}
		log.Println(Error)
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJson(w, Error)
		return
	}

	task.ID = r.URL.Query().Get("id")
	if task.ID == "" {
		Error := db.Task{Error: "Validation error for marking task as done: missing ID"}
		log.Println(Error)
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, Error)
		return
	}
	log.Println("Marking task as done", task.ID)

	task, err := db.GetTask(task.ID)
	if err != nil {
		Error := db.Task{Error: "Error for marking task as done: " + err.Error()}
		log.Println(Error)
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, Error)
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(task.ID)
		if err != nil {
			if err.Error() == "incorrect id for deleting task" {
				Error := db.Task{Error: "Error for marking task as done: " + err.Error()}
				log.Println(Error)
				w.WriteHeader(http.StatusBadRequest)
				writeJson(w, Error)
				return
			}
			Error := db.Task{Error: "Error for marking task as done: " + err.Error()}
			log.Println(Error)
			w.WriteHeader(http.StatusInternalServerError)
			writeJson(w, Error)
			return
		}
		log.Println("Task", task.ID, "was successfully removed")
		writeJson(w, db.Task{})
		return
	}

	log.Println("Calculating next date for task:", task)
	task.Date, err = db.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		Error := db.Task{Error: "Error for calculating next date: " + err.Error()}
		log.Println(Error)
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, Error)
		return
	}
	log.Println("Next date for task:", task.Date)

	err = db.UpdateDate(&task)
	if err != nil {
		if err.Error() == "incorrect id for updating date" {
			Error := db.Task{Error: "Error for marking task as done: " + err.Error()}
			log.Println(Error)
			w.WriteHeader(http.StatusBadRequest)
			writeJson(w, Error)
			return
		}
		Error := db.Task{Error: "Error for marking task as done: " + err.Error()}
		log.Println(Error)
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, Error)
		return
	}
	
	log.Println("Successfully marked as done: ", task.ID)
	writeJson(w, db.Task{})
}

func Init() {
	http.Handle("/", http.FileServer(http.Dir("web")))
	http.HandleFunc("/api/nextdate", HandleNextDate)
	http.HandleFunc("/api/task", HandleTask)
	http.HandleFunc("/api/tasks", HandleTasks)
	http.HandleFunc("/api/task/done", HandleTaskDone)
}
