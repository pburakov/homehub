package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"
	"time"
)

func StartMotionAndKeepAlive(m *Motion) {
	log.Print("Starting motion")
	c := exec.Command("motion", "-c", m.ConfPath)
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

// MustDumpMotionConf generates motion.conf file and returns its path.
func MustDumpMotionConf(m *Motion) string {
	p := MustGetCWD() + "/" + motionConfFile
	f, e := template.ParseFiles(motionConfTemplate)
	if e != nil {
		Fatal(fmt.Errorf("error reading motion config template: %s", e))
	}
	w, e := os.Create(p)
	if e != nil {
		Fatal(fmt.Errorf("error creating motion config: %s", e))
	}
	e = f.Execute(w, m)
	if e != nil {
		Fatal(fmt.Errorf("error writing motion config: %s", e))
	}
	return p
}
