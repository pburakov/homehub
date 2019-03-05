package main

import (
	"flag"
	"io.pburakov/homehub/agent/config"
	"io.pburakov/homehub/agent/http"
	"io.pburakov/homehub/agent/motion"
	"io.pburakov/homehub/agent/rpc"
	hh "io.pburakov/homehub/agent/schema"
	"io.pburakov/homehub/agent/util"
	"log"
)

func main() {
	c := config.InitConfig()

	// Get flags from command line
	fRemote := flag.String("r", c.RemoteHubAddress, "Remote hub server address (including port)")
	fWebPort := flag.Uint("pw", c.WebServer.Port, "Local web service port to bind")
	fSensorsPort := flag.Uint("pm", c.Sensors.Port, "Local sensor feed port to bind")
	fMotionPort := flag.Uint("ps", c.Motion.Port, "Local streaming port to bind")
	flag.Parse()

	// Start motion detection and video-streaming process
	go motion.StartMotionAndKeepAlive(&c.Motion)

	// Start serving files from motion-detector output
	go http.ServeFolder(c)

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
		req := rpc.BuildRequest(c.AgentId, eIP, *fWebPort, *fMotionPort, *fSensorsPort)
		res, e := rpc.CheckIn(client, c.ConnectionTimeout, req)
		if e != nil {
			log.Printf("Check-in failed: %s", e)
		} else {
			log.Printf("Check-in ACK: %s", res.String())
		}
	}, c.CheckInInterval)

	util.Wait()
}
