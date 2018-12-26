package main

import (
	"context"
	"github.com/denisbrodbeck/machineid"
	"github.com/joshb/pi-camera-go/server/recorder"
	hh "github.com/pburakov/homehub/schema"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	serverAddress     = "localhost:31321" // HomeHub "mothership" server
	hubPort           = 84526             // Port the hub is listening to
	appID             = "homehub"         // App identifier used to generate unique hub id
	checkInInterval   = time.Second * 10  // Defines frequency of check-ins
	connectionTimeout = time.Second * 2   // How long to wait for check-in to succeed¬
)

func main() {
	c := setUpConnection(serverAddress)
	defer c.Close()

	client := hh.NewHomeHubClient(c)

	ticker := time.NewTicker(checkInInterval)

	// Initiate check-in and repeat with intervals
	checkIn(client, connectionTimeout)
	go func() {
		for range ticker.C {
			checkIn(client, connectionTimeout)
		}
	}()

	r := startRecording()
	defer stopRecording(r)

	// Wait forever until interrupted
	wait()

	// TODO: implement graceful stop
}

func getExternalIP() string {
	resp, err := http.Get("http://ifconfig.me")
	if err != nil {
		log.Fatalf("[error] could not determine machine address: %s", err)
	}
	defer resp.Body.Close()
	externalIp, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(externalIp) == 0 {
		log.Fatalf("[error] could not parse machine address '%s'", externalIp)
	}
	return string(externalIp)
}

func setUpConnection(address string) *grpc.ClientConn {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[error] could not connect: %v", err)
	}
	return conn
}

func getMachineId(appID string) string {
	m, err := machineid.ProtectedID(appID)
	if err != nil {
		log.Fatalf("[error] could not determine machine id: %s", err)
	}
	return m
}

func checkIn(c hh.HomeHubClient, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r, err := c.CheckIn(ctx, &hh.CheckInRequest{
		HubId:   getMachineId(appID),
		Address: getExternalIP(),
		Port:    hubPort,
	})
	if err == nil {
		log.Printf("Check-in ACK: %s", r.Result)
	} else {
		log.Printf("[error] could not check-in: %v", err)
	}
}

func wait() {
	select {}
}

func startRecording() recorder.Recorder {
	var r recorder.Recorder
	r, e := recorder.New()
	if e != nil {
		log.Printf("[error] could not initiate recorder: %s", e)
	}
	e = r.Start()
	if e != nil {
		log.Printf("[error] could not start recording: %s", e)
		log.Print("Using mock recorder")
		r = recorder.NewMock()
	}
	log.Print("Started recording")
	return r
}

func stopRecording(r recorder.Recorder) {
	e := r.Stop()
	if e != nil {
		log.Fatalf("[error] could not stop recording: %s", e)
	}
	log.Print("Stopped recording")
}
