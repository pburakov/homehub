package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	AppID              = "homehub"
	agentConfFile      = "conf/agent.json"
	motionConfTemplate = "conf/motion.template"
	motionConfFile     = "motion.conf"
)

type WebServer struct {
	Port uint
}

type Motion struct {
	// These values are used in motion.conf template
	Port     uint
	Username string
	Password string
	Dir      string

	// These values are used by runtime
	ConfPath          string
	KeepAliveInterval time.Duration `json:"keepalive_ping_interval_seconds"`
}

type Sensors struct {
	Port uint
}

type Configuration struct {
	AgentId           string
	Motion            Motion
	Sensors           Sensors
	CWD               string
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

	// Prepare directories, populate auto-generated fields and convert durations
	c.Motion.Dir = MustCreateMotionDir()
	c.Motion.ConfPath = MustDumpMotionConf(&c.Motion)
	c.Motion.KeepAliveInterval = c.Motion.KeepAliveInterval * time.Second
	c.AgentId = MustGetMachineId(AppID)
	c.CheckInInterval = c.CheckInInterval * time.Second
	c.ConnectionTimeout = c.ConnectionTimeout * time.Second

	return c
}

func MustCreateMotionDir() string {
	d := MustGetCWD()
	md := d + "/motion"
	e := os.Mkdir(md, 0766)
	if e != nil && !os.IsExist(e) {
		Fatal(fmt.Errorf("unable to create dir for motion output: %s", e))
	}
	return md
}
