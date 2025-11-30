package main

import (
	"log"

	database "github.com/PIPILaPUPU/finalTest/pkg/db"
	"github.com/PIPILaPUPU/finalTest/pkg/server"
)

func main() {
	if err := database.Init("scheduler.db"); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer database.Close()

	server.StartServer()
}
