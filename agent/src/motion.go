package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

func StartMotionAndKeepAlive(m *Motion) {
	log.Printf("Starting motion with streaming on port %d", m.Port)
	c := exec.Command(m.Executable, "-c", m.ConfPath)
	e := c.Start()
	if e != nil {
		Fatal(fmt.Errorf("error starting motion: %s", e))
	}
	log.Printf("Motion started with pid %d", c.Process.Pid)
	s, e := c.Process.Wait() // will block here until the process terminates
	if e == nil {
		log.Printf("Motion terminated with %s", s)
	} else {
		log.Printf("Error polling motion process: %s", e)
	}
	log.Printf("Will restart motion in %s...", m.KeepAliveInterval)
	time.Sleep(m.KeepAliveInterval)
	go StartMotionAndKeepAlive(m)
	return
}
