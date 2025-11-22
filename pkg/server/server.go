package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ArtyomGaribyan/Task-Scheduler/pkg/api"
	"github.com/ArtyomGaribyan/Task-Scheduler/tests"
)

func Run() {
	log.Println("Starting server on port", tests.Port)

	api.Init()

	port := fmt.Sprintf(":%d", tests.Port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
