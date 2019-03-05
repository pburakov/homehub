package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"
)

func StartMotion(m *Motion) {
	config := MustDumpMotionConf(m)
	c := exec.Command("motion", "-c", config)
	e := c.Start()
	if e != nil {
		Fatal(fmt.Errorf("error starting motion: %s", e))
	}
	log.Printf("Motion started with pid %d", c.Process.Pid)
	log.Printf("Motion status is %s", c.ProcessState.String())
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
