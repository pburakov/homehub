package rpc

import (
	"google.golang.org/grpc"
	"log"
)

func SetUpConnection(address string) *grpc.ClientConn {
	log.Printf("Connecting to %s", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	return conn
}
