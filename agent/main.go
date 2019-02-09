package main

import (
	"flag"
	"github.com/pburakov/homehub/conf"
	"github.com/pburakov/homehub/rpc"
	hh "github.com/pburakov/homehub/schema"
	"github.com/pburakov/homehub/util"
	"log"
)

func main() {
	config := conf.Init()

	// Get flags from command line
	fRemote := flag.String("r", config.RemoteHubAddress, "Remote hub server address (including port)")
	fWebPort := flag.Uint("pw", config.Ports.Web, "Local web service port to bind")
	fMetaPort := flag.Uint("pm", config.Ports.Meta, "Local meta port to bind")
	// TODO: uncomment when supported
	// fStreamPort := flag.Uint("ps", config.Ports.Stream, "Local streaming port to bind")
	fStreamPort := &config.Ports.Stream
	flag.Parse()

	conn := rpc.SetUpConnection(*fRemote)
	defer conn.Close()

	client := hh.NewHomeHubClient(conn)

	util.Schedule(func() {
		req, _ := rpc.BuildRequest(config.AppId, *fWebPort, *fStreamPort, *fMetaPort)
		res, e := rpc.CheckIn(client, config.ConnectionTimeout, req)
		if e != nil {
			log.Printf("Check-in failed: %s", e)
		} else {
			log.Printf("Check-in ACK: %s", res.String())
		}
	}, config.CheckInInterval)

	util.Wait()
}
