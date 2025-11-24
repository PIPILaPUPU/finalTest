package tests

import (
	"os"
	"strconv"
)

var Port = GetPort()
var DBFile = "../scheduler.db"
var FullNextDate = false
var Search = false
var Token = ``

func GetPort() int {
	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			return port
		}
	}

	return 7540
}
