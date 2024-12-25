package api

import (
	"context"
	"modular_chassis/echo/pkg/services/carpay"
	mediatorAPIs "modular_chassis/mediator/pkg/api"
	"sync"
)

type impls struct{}

var (
	once     sync.Once
	implsIns *impls
)

func getService() *impls {
	once.Do(func() {
		if implsIns == nil {
			implsIns = &impls{}
		}
	})
	return implsIns
}

func init() {
	mediatorAPIs.GetMediatorAPI().RegisterServiceFunc(getService())
}

func (impls) ListVehicles(ctx context.Context, req carpay.ListVehiclesRequest) (carpay.ListVehiclesResponse, error) {
	return carpay.ListVehiclesResponse{Vehicles: []carpay.Vehicle{{"86768"}, {"68778798"}}}, nil
}

func (impls) Test(ctx context.Context, request carpay.TestRequest) (carpay.TestResponse, error) {
	return carpay.TestResponse{Output: request.Input}, nil
}
