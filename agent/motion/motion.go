package motion

import (
	"fmt"
	"io.pburakov/homehub/agent/config"
	"io.pburakov/homehub/agent/util"
	"log"
	"os/exec"
	"time"
)

func StartMotionAndKeepAlive(m *config.Motion) {
	log.Print("Starting motion")
	c := exec.Command("motion", "-c", m.ConfPath)
	e := c.Start()
	if e != nil {
		util.Fatal(fmt.Errorf("error starting motion: %s", e))
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
