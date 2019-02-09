package rpc

import (
	hh "github.com/pburakov/homehub/schema"
	"github.com/pburakov/homehub/util"
)

func BuildRequest(appID string, pWeb uint, pStream uint, pMeta uint) (*hh.CheckInRequest, error) {
	eip, e := util.GetExternalIP()
	if e != nil {
		return nil, e
	}
	return &hh.CheckInRequest{
		HubId:      util.MachineId(appID),
		Address:    eip,
		WebPort:    int32(pWeb),
		StreamPort: int32(pStream),
		MetaPort:   int32(pMeta),
	}, nil
}
