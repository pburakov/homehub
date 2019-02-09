package service

import (
	"context"
	hh "github.com/pburakov/homehub/schema"
	"time"
)

func CheckIn(c hh.HomeHubClient, timeout time.Duration, req *hh.CheckInRequest) (*hh.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r, e := c.CheckIn(ctx, req)
	if e != nil {
		return nil, e
	}
	return &r.Result, nil
}
