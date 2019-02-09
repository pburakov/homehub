package conf

import (
	"encoding/json"
	"fmt"
	"github.com/pburakov/homehub/util"
	"io/ioutil"
	"os"
	"time"
)

const (
	AppID    = "homehub"
	confFile = "conf/agent.json"
)

type Ports struct {
	Web    uint
	Stream uint
	Meta   uint
}

type Configuration struct {
	AgentId           string
	Ports             Ports
	RemoteHubAddress  string        `json:"remote_hub_address"`
	CheckInInterval   time.Duration `json:"check_in_interval_seconds"`
	ConnectionTimeout time.Duration `json:"connection_timeout_seconds"`
}

func Init() *Configuration {
	f, e := os.Open(confFile)
	if e != nil {
		util.Fatal(fmt.Errorf("error loading configuration from %s: %s", confFile, e))
	}
	defer f.Close()
	b, e := ioutil.ReadAll(f)
	if e != nil {
		util.Fatal(fmt.Errorf("error reading configuration file: %s", e))
	}
	c := new(Configuration)
	if e := json.Unmarshal(b, c); e != nil {
		util.Fatal(fmt.Errorf("invalid configuration file: %s", e))
	}

	// Populate auto-generated fields and convert durations
	c.AgentId = util.MustGetMachineId(AppID)
	c.CheckInInterval = c.CheckInInterval * time.Second
	c.ConnectionTimeout = c.ConnectionTimeout * time.Second

	return c
}
