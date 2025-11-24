package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PIPILaPUPU/finalTest/tests"
)

const (
	webPath = "./web"
)

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func main() {
	port := tests.GetPort()

	if !dirExists(webPath) {
		log.Fatalf("there is no such dir %s", webPath)
	}

	fs := http.FileServer(http.Dir(webPath))
	http.Handle("/", fs)

	addr := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(addr, fs); err != nil {
		log.Fatal(err)
	}
}
