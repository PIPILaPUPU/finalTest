package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PIPILaPUPU/finalTest/pkg/api"
	"github.com/PIPILaPUPU/finalTest/tests"
)

const webDir = "./web"

func StartServer() {
	port := tests.GetPort()

	password := os.Getenv("TODO_PASSWORD")
	if password != "" {
		log.Println("Аутентификация включена")
	} else {
		log.Println("Аутентификация отключена (TODO_PASSWORD не установлен)")
	}

	if err := api.Init(); err != nil {
		log.Fatal("Failed to initialize API:", err)
	}

	info, err := os.Stat(webDir)
	if err != nil || !info.IsDir() {
		log.Fatalf("Directory %s not found", webDir)
	}

	fileServer := http.FileServer(http.Dir(webDir))

	http.Handle("/", fileServer)

	addr := fmt.Sprintf(":%d", port)

	log.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
