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
	c := exec.Command("motion", "-c", m.ConfPath)
	e := c.Start()
	if e != nil {
		Fatal(fmt.Errorf("error starting motion: %s", e))
	}
	pid := c.Process.Pid
	log.Printf("Motion started with pid %d", pid)
	ticker := time.NewTicker(m.KeepAliveInterval)
	for range ticker.C {
		_, e := os.FindProcess(pid)
		if e != nil {
			log.Print("Motion appears to be down")
			go StartMotionAndKeepAlive(m)
			return
		}
	}
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
