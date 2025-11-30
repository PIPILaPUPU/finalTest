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
	api.Init()

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
