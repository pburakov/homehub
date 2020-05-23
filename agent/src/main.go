package main

import (
	"flag"
	"log"
)

func main() {
	c := InitConfig()

	// Get flags from command line
	fRemote := flag.String("r", c.RemoteHubAddress, "Remote hub server address (including port)")
	fWebPort := flag.Uint("pw", c.WebServer.Port, "Local web service port to bind")
	fSensorsPort := flag.Uint("pm", c.Sensors.Port, "Local sensor feed port to bind")
	fMotionPort := flag.Uint("ps", c.Motion.Port, "Local streaming port to bind")
	flag.Parse()

	// Start motion detection and video-streaming process
	go StartMotionAndKeepAlive(&c.Motion)

	// Start serving files from motion-detector output
	go ServeFolder(c)

	// Create RPC connection and schedule RPC check-in
	conn := SetUpConnection(*fRemote)
	defer conn.Close()
	client := NewHomeHubClient(conn)
	Schedule(func() {
		eIP, e := GetExternalIP()
		if e != nil {
			log.Printf("Error obtaining external address: %s", e)
			return
		}
		req := BuildRequest(c.AgentId, eIP, *fWebPort, *fMotionPort, *fSensorsPort)
		res, e := CheckIn(client, c.ConnectionTimeout, req)
		if e != nil {
			log.Printf("Check-in failed: %s", e)
		} else {
			log.Printf("Check-in ACK: %s", res.String())
		}
	}, c.CheckInInterval)

	Wait()
}
