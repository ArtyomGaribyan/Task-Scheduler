package db

import (
	"database/sql"
	"log"

	"github.com/ArtyomGaribyan/Task-Scheduler/tests"

	_ "modernc.org/sqlite"
)

func InitDB() error {
	log.Println("Opening database file:", tests.DBFile)
	
	db, err := sql.Open(tests.SQL, tests.DBFile)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	date CHAR(8) NOT NULL DEFAULT '',
		title VARCHAR(50) NOT NULL DEFAULT '',
		comment TEXT NOT NULL DEFAULT '',
		repeat VARCHAR(20) NOT NULL DEFAULT ''
	)`)
	if err != nil {
		return err
	}
	log.Println("Database initialized successfully")

	return nil
}
