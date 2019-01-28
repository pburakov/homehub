package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"github.com/joshb/pi-camera-go/server"
	hh "github.com/pburakov/homehub/schema"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	// Default HomeHub "mothership" server (local configuration)
	defaultMothershipServer = "localhost:31321"

	// Default address and port the hub is listening to
	defaultAddress = "0.0.0.0"
	defaultPort    = 31322

	// App identifier used to generate unique hub id
	appID = "homehub"

	checkInInterval   = time.Second * 10 // Defines frequency of check-ins
	connectionTimeout = time.Second * 2  // How long to wait for check-in to succeed¬
)

func main() {
	// Get flags from command line
	mothershipAddress := flag.String("r", defaultMothershipServer, "Remote address (including port) of a mothership server")
	localAddress := flag.String("a", defaultAddress, "Local address to bind to")
	localPort := flag.Uint("p", defaultPort, "Local port to bind to")
	useHTTPS := flag.Bool("https", false, "Use HTTPS")

	externalIP := mustGetExternalIP()

	c := setUpConnection(*mothershipAddress)
	defer c.Close()

	client := hh.NewHomeHubClient(c)

	// Initiate check-in and then repeat with intervals
	checkIn(client, connectionTimeout, externalIP, *localPort)
	ticker := time.NewTicker(checkInInterval)
	go func() {
		for range ticker.C {
			checkIn(client, connectionTimeout, externalIP, *localPort)
		}
	}()

	mustStartVideoServer(fmt.Sprintf("%s:%d", *localAddress, *localPort), *useHTTPS)
}

// mustGetExternalIP resolves external IP using server ifconfig.me
func mustGetExternalIP() string {
	resp, err := http.Get("http://ifconfig.me")
	if err != nil {
		log.Fatalf("Could not determine machine address: %s", err)
	}
	defer resp.Body.Close()
	externalIp, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(externalIp) == 0 {
		log.Fatalf("Could not parse machine address '%s'", externalIp)
	}
	return string(externalIp)
}

func setUpConnection(address string) *grpc.ClientConn {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	return conn
}

func getMachineId(appID string) string {
	m, err := machineid.ProtectedID(appID)
	if err != nil {
		log.Fatalf("Could not determine machine id: %s", err)
	}
	return m
}

func checkIn(c hh.HomeHubClient, t time.Duration, address string, port uint) {
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	r, err := c.CheckIn(ctx, &hh.CheckInRequest{
		HubId:   getMachineId(appID),
		Address: address,
		Port:    int32(port),
	})
	if err == nil {
		log.Printf("Check-in ACK: %s", r.Result)
	} else {
		log.Printf("Could not check-in: %v", err)
	}
}

func mustStartVideoServer(address string, useHTTPS bool) {
	flag.Parse()

	s, err := server.New(useHTTPS)
	if err != nil {
		log.Fatalf("Unable to create server: %s", err)
	}

	if err := s.Start(address); err != nil {
		log.Fatalf("Unable to start server: %s", err)
	}
}
