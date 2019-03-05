package util

import (
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func MustGetMachineId(appID string) string {
	m, e := machineid.ProtectedID(appID)
	if e != nil {
		Fatal(fmt.Errorf("error generating machine id: %s", e))
	}
	return m
}

func GetExternalIP() (string, error) {
	r, e := http.Get("http://ifconfig.me")
	if e != nil {
		return "unknown", e
	}
	defer r.Body.Close()
	ip, e := ioutil.ReadAll(r.Body)
	if e != nil || len(ip) == 0 {
		return "unknown", e
	}
	return string(ip), nil
}

func MustCreateMotionDir() string {
	d, e := os.Getwd()
	if e != nil {
		Fatal(fmt.Errorf("error getting current working directory: %s", e))
	}
	md := d + "/motion"
	e = os.Mkdir(md, 0666)
	if e != nil && !os.IsExist(e) {
		Fatal(fmt.Errorf("unable to create dir for motion output: %s", e))
	}
	return md
}

func StartMotion() {
	c := exec.Command("motion")
	e := c.Start()
	if e != nil {
		Fatal(fmt.Errorf("error starting motion: %s", e))
	}
	log.Printf("Motion started with pid %d", c.Process.Pid)
}
