package main

import (
	"flag"
	"github.com/pburakov/homehub/rpc"
	hh "github.com/pburakov/homehub/schema"
	"github.com/pburakov/homehub/service"
	"github.com/pburakov/homehub/util"
	"log"
	"time"
)

const (
	// Default "mothership" server address
	defaultHub = "localhost:31321"

	// Default port the agent is listening to
	defaultPort = 31322

	// App identifier used to generate unique hub id
	appID = "homehub"

	checkInInterval   = time.Second * 10 // Defines frequency of check-ins
	connectionTimeout = time.Second * 2  // How long to wait for check-in to succeed¬
)

func main() {
	// Get flags from command line
	hubAddress := flag.String("r", defaultHub, "Remote hub address (including port), mothership server to check in with")
	localPort := flag.Uint("p", defaultPort, "Local port to bind to")
	flag.Parse()

	conn := rpc.SetUpConnection(*hubAddress)
	defer conn.Close()

	client := hh.NewHomeHubClient(conn)

	util.Schedule(func() {
		req, _ := buildRequest(*localPort)
		res, e := service.CheckIn(client, connectionTimeout, req)
		if e != nil {
			log.Printf("Check-in failed: %s", e)
		} else {
			log.Printf("Check-in ACK: %s", res.String())
		}
	}, checkInInterval)

	util.Wait()
}

func buildRequest(port uint) (*hh.CheckInRequest, error) {
	eip, e := util.GetExternalIP()
	if e != nil {
		return nil, e
	}
	return &hh.CheckInRequest{
		HubId:   util.MachineId(appID),
		Address: eip,
		Port:    int32(port),
	}, nil
}
