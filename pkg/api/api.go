package api

import (
	"net/http"

	"github.com/PIPILaPUPU/finalTest/pkg/auth"
)

func Init() error {
	if err := auth.Init(); err != nil {
		return err
	}

	http.HandleFunc("/api/signin", signInHandler)

	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", AuthMiddleware(taskHandler))
	http.HandleFunc("/api/tasks", AuthMiddleware(returnTasksHandler))
	http.HandleFunc("/api/task/done", AuthMiddleware(taskDoneHandler))

	return nil
}
