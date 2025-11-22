package main

import (
	"log"

	"github.com/ArtyomGaribyan/Task-Scheduler/pkg/db"
	"github.com/ArtyomGaribyan/Task-Scheduler/pkg/server"
)

func main() {
	go func() {
		err := db.InitDB()
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
	}()
	
	server.Run()
}
