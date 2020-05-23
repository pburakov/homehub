package main

import (
	"context"
	"google.golang.org/grpc"
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

func BuildRequest(aID string, eIP string, pWeb uint, pStream uint, pSensors uint) *CheckInRequest {
	return &CheckInRequest{
		AgentId:     aID,
		Address:     eIP,
		WebPort:     int32(pWeb),
		StreamPort:  int32(pStream),
		SensorsPort: int32(pSensors),
	}
}

func CheckIn(c HomeHubClient, timeout time.Duration, req *CheckInRequest) (*Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r, e := c.CheckIn(ctx, req)
	if e != nil {
		return nil, e
	}
	return &r.Result, nil
}
