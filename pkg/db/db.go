package db

import (
	"database/sql"
	"log"

	"github.com/ArtyomGaribyan/Task-Scheduler/tests"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	log.Println("Opening database file:", tests.DBFile)

	var err error
	DB, err = sql.Open(tests.SQL, tests.DBFile)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
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

func Close() {
	log.Println("Closing database")
	DB.Close()
}
