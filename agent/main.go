package main

import (
	"context"
	"flag"
	"github.com/denisbrodbeck/machineid"
	hh "github.com/pburakov/homehub/schema"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	// Default "mothership" server address
	defaultHub = "localhost:8000"

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

	c := setUpConnection(*hubAddress)
	defer c.Close()

	client := hh.NewHomeHubClient(c)

	// Initiate check-in and then repeat with intervals
	checkIn(client, connectionTimeout, *localPort)
	ticker := time.NewTicker(checkInInterval)
	go func() {
		for range ticker.C {
			checkIn(client, connectionTimeout, *localPort)
		}
	}()

	wait()
}

// mustGetExternalIP resolves external IP using server ifconfig.me or 'unknown'
func getExternalIP() string {
	resp, err := http.Get("http://ifconfig.me")
	if err != nil {
		log.Printf("Could not determine machine address: %s", err)
		return "unknown"
	}
	defer resp.Body.Close()
	externalIp, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(externalIp) == 0 {
		log.Printf("Could not parse machine address '%s'", externalIp)
		return "unknown"
	}
	return string(externalIp)
}

func setUpConnection(address string) *grpc.ClientConn {
	log.Printf("Connecting to %s", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	return conn
}

func mustGetMachineId(appID string) string {
	m, err := machineid.ProtectedID(appID)
	if err != nil {
		log.Fatalf("Could not determine machine id: %s", err)
	}
	return m
}

func checkIn(c hh.HomeHubClient, t time.Duration, port uint) {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	r, err := c.CheckIn(ctx, &hh.CheckInRequest{
		HubId:   mustGetMachineId(appID),
		Address: getExternalIP(),
		Port:    int32(port),
	})
	if err == nil {
		log.Printf("Check-in ACK: %s", r.Result)
	} else {
		log.Printf("Could not check-in: %v", err)
	}
}

func wait() {
	select {}
}
