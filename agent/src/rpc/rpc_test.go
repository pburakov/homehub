package rpc

import (
	hh "io.pburakov/homehub/agent/schema"
	"reflect"
	"testing"
)

func TestBuildRequest(t *testing.T) {
	type args struct {
		aID      string
		eIP      string
		pWeb     uint
		pStream  uint
		pSensors uint
	}
	tests := []struct {
		name string
		args args
		want *hh.CheckInRequest
	}{
		{"noname",
			args{
				aID:      "testid",
				eIP:      "testip",
				pWeb:     123,
				pStream:  456,
				pSensors: 789,
			},
			&hh.CheckInRequest{
				AgentId:     "testid",
				Address:     "testip",
				WebPort:     123,
				StreamPort:  456,
				SensorsPort: 789,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildRequest(tt.args.aID, tt.args.eIP, tt.args.pWeb, tt.args.pStream, tt.args.pSensors); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
