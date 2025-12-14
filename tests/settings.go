package tests

import (
	"os"
	"strconv"
)

var Port = GetPort()
var DBFile = "../scheduler.db"
var FullNextDate = false
var Search = true
var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwd2RfaGFzaCI6IjU5OTQ0NzFhYmIwMTExMmFmY2MxODE1OWY2Y2M3NGI0ZjUxMWI5OTgwNmRhNTliM2NhZjVhOWMxNzNjYWNmYzUiLCJleHAiOjE3NjU3NTc1MzgsImlhdCI6MTc2NTcyODczOH0.u5kKrEPEZaa2v13aVSGMJRsoXxEb7y1ScrQJHOj1OUQ`

func GetPort() int {
	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			return port
		}
	}

	return 7540
}
