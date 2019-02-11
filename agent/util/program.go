package util

import (
	"os"
	"time"
)

// Schedule initiates action and then schedules repeats with intervals
func Schedule(action func(), interval time.Duration) {
	action()
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			action()
		}
	}()
}

func Fatal(e error) {
	println(e.Error())
	os.Exit(1)
}

func Wait() {
	select {}
}
