package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PIPILaPUPU/finalTest/tests"
)

const webDir = "./web"

func main() {
	port := tests.GetPort()

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
