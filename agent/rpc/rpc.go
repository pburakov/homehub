package rpc

import (
	"context"
	"google.golang.org/grpc"
	hh "io.pburakov/homehub/agent/schema"
	"log"
	"time"
)

func SetUpConnection(address string) *grpc.ClientConn {
	log.Printf("Connecting to %s", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	return conn
}

func BuildRequest(aID string, eIP string, pWeb uint, pStream uint, pSensors uint) *hh.CheckInRequest {
	return &hh.CheckInRequest{
		AgentId:     aID,
		Address:     eIP,
		WebPort:     int32(pWeb),
		StreamPort:  int32(pStream),
		SensorsPort: int32(pSensors),
	}
}

func CheckIn(c hh.HomeHubClient, timeout time.Duration, req *hh.CheckInRequest) (*hh.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r, e := c.CheckIn(ctx, req)
	if e != nil {
		return nil, e
	}
	return &r.Result, nil
}
