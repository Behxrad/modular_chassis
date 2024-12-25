package api

import (
	"context"
	"errors"
	"fmt"
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
	fmt.Println(">>> ListVehicles called")
	return carpay.ListVehiclesResponse{Vehicles: make([]carpay.Vehicle, 0)}, errors.New("sample error")
}

func (impls) Test(ctx context.Context, request carpay.TestRequest) (carpay.TestResponse, error) {
	fmt.Println(">>> Test called")
	return carpay.TestResponse{}, nil
}
