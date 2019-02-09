package util

import (
	"github.com/denisbrodbeck/machineid"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Schedule initiates action and then repeat with intervals
func Schedule(action func(), interval time.Duration) {
	action()
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			action()
		}
	}()
}

func MachineId(appID string) string {
	m, e := machineid.ProtectedID(appID)
	if e != nil {
		Fatal(e)
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

func Fatal(e error) {
	println(e)
	os.Exit(1)
}

func Wait() {
	select {}
}
