package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
	"time"
)

const (
	AppID              = "homehub"
	agentConfFile      = "conf/agent.json"
	motionConfFile     = "motion.conf"
	motionConfTemplate = "conf/motion.template"
)

type WebServer struct {
	Port uint
}

type Motion struct {
	Port     uint
	Username string
	Password string
	Dir      string
}

type Sensors struct {
	Port uint
}

type Configuration struct {
	AgentId           string
	Motion            Motion
	Sensors           Sensors
	WebServer         WebServer     `json:"web_server"`
	RemoteHubAddress  string        `json:"remote_hub_address"`
	CheckInInterval   time.Duration `json:"check_in_interval_seconds"`
	ConnectionTimeout time.Duration `json:"connection_timeout_seconds"`
}

// InitConfig initializes program configuration
func InitConfig() *Configuration {
	f, e := os.Open(agentConfFile)
	if e != nil {
		Fatal(fmt.Errorf("error loading configuration from %s: %s", agentConfFile, e))
	}
	defer f.Close()
	b, e := ioutil.ReadAll(f)
	if e != nil {
		Fatal(fmt.Errorf("error reading configuration file: %s", e))
	}
	c := new(Configuration)
	if e := json.Unmarshal(b, c); e != nil {
		Fatal(fmt.Errorf("invalid configuration file: %s", e))
	}

	// Populate auto-generated fields and convert durations
	c.AgentId = MustGetMachineId(AppID)
	c.CheckInInterval = c.CheckInInterval * time.Second
	c.ConnectionTimeout = c.ConnectionTimeout * time.Second
	c.Motion.Dir = MustCreateMotionDir()

	return c
}

func DumpMotionConf(m *Motion) {
	f, e := template.ParseFiles(motionConfTemplate)
	if e != nil {
		Fatal(fmt.Errorf("error reading motion config template: %s", e))
	}
	w, e := os.Create(motionConfFile)
	if e != nil {
		Fatal(fmt.Errorf("error creating motion config: %s", e))
	}
	e = f.Execute(w, m)
	if e != nil {
		Fatal(fmt.Errorf("error writing motion config: %s", e))
	}
}
