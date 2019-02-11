package main

import (
	"flag"
	"io.pburakov/homehub/agent/rpc"
	hh "io.pburakov/homehub/agent/schema"
	"io.pburakov/homehub/agent/util"
	"log"
)

func main() {
	config := util.InitConfig()

	// Get flags from command line
	fRemote := flag.String("r", config.RemoteHubAddress, "Remote hub server address (including port)")
	fWebPort := flag.Uint("pw", config.WebServer.Port, "Local web service port to bind")
	fSensorsPort := flag.Uint("pm", config.Sensors.Port, "Local sensor feed port to bind")
	fMotionPort := flag.Uint("ps", config.Motion.Port, "Local streaming port to bind")
	flag.Parse()

	// Prepare motion startup
	util.DumpMotionConf(&config.Motion)

	// Create RPC connection and schedule RPC check-in
	conn := rpc.SetUpConnection(*fRemote)
	defer conn.Close()
	client := hh.NewHomeHubClient(conn)
	util.Schedule(func() {
		eIP, e := util.GetExternalIP()
		if e != nil {
			log.Printf("Error obtaining external address: %s", e)
			return
		}
		req := rpc.BuildRequest(config.AgentId, eIP, *fWebPort, *fMotionPort, *fSensorsPort)
		res, e := rpc.CheckIn(client, config.ConnectionTimeout, req)
		if e != nil {
			log.Printf("Check-in failed: %s", e)
		} else {
			log.Printf("Check-in ACK: %s", res.String())
		}
	}, config.CheckInInterval)

	util.Wait()
}
