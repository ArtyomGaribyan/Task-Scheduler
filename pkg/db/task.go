package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ArtyomGaribyan/Task-Scheduler/tests"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Date    string `json:"date,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
	Error   string `json:"error,omitempty"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	db, err := sql.Open(tests.SQL, tests.DBFile)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	res, err := db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func Tasks(limit int) ([]*Task, error) {
	db, err := sql.Open(tests.SQL, tests.DBFile)
	if err != nil {
		return []*Task{}, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", limit)
	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return []*Task{}, err
		}
		tasks = append(tasks, &task)
	}

	if tasks == nil {
		tasks = []*Task{}
	}
	log.Println("tasks:", tasks)

	return tasks, nil
}

func GetTask(id string) (Task, error) {
	db, err := sql.Open(tests.SQL, tests.DBFile)
	if err != nil {
		return Task{}, err
	}
	defer db.Close()

	var task Task
	err = db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func UpdateTask(task *Task) error {
	db, err := sql.Open(tests.SQL, tests.DBFile)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Println("Updating task in DB:", task.ID, task.Title, task.Date, task.Comment, task.Repeat)
	res, err := db.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}

	log.Println("Rows affected:", count)
	return nil
}

func UpdateDate(task *Task) error {
	db, err := sql.Open(tests.SQL, tests.DBFile)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Println("Updating date in DB:", task)
	res, err := db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, task.Date, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating date")
	}

	log.Println("Rows affected:", count)
	return nil
}

func DeleteTask(id string) error {
	db, err := sql.Open(tests.SQL, tests.DBFile)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Println("Deleting task with ID:", id)
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	log.Println("Rows affected:", count)
	return nil
}
