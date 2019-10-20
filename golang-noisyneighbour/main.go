package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	var sleepInterval = 1 // default 1 microsecond
	if sleepIntervalEnv := os.Getenv("SLEEP_INTERVAL_MS"); sleepIntervalEnv != "" {
		sleepInterval, _ = strconv.Atoi(sleepIntervalEnv)
	}

	for {
		fmt.Println("log message")
		time.Sleep(time.Microsecond * time.Duration(sleepInterval))
	}
}
